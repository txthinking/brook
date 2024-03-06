// Copyright (c) 2016-present Cloud <cloud@txthinking.com>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of version 3 of the GNU General Public
// License as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package brook

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	utls "github.com/refraction-networking/utls"

	"github.com/txthinking/x"
	x1 "github.com/txthinking/x"
)

func WebSocketDial(src, dst, addr, host, path string, tc *tls.Config, timeout int, tlsfingerprint utls.ClientHelloID, fragmentMinLength, fragmentMaxLength, fragmentMinDelay, fragmentMaxDelay int64) (net.Conn, error) {
	var c net.Conn
	var err error
	if src == "" || dst == "" {
		c, err = DialTCP("tcp", "", addr)
	}
	if src != "" && dst != "" {
		c, err = NATDial("tcp", src, dst, addr)
	}
	if err != nil {
		return nil, err
	}
	if timeout != 0 {
		if err := c.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			c.Close()
			return nil, err
		}
	}
	if tc != nil {
		if fragmentMinLength != 0 && fragmentMaxLength != 0 && fragmentMinDelay != 0 && fragmentMaxDelay != 0 {
			c = &TLSFragmentConn{
				Conn:      c,
				MinLength: fragmentMinLength,
				MaxLength: fragmentMaxLength,
				MinDelay:  fragmentMinDelay,
				MaxDelay:  fragmentMaxDelay,
			}
		}
		if tlsfingerprint.Client == "" {
			c1 := tls.Client(c, tc)
			if err := c1.Handshake(); err != nil {
				c1.Close()
				return nil, err
			}
			s := host
			h, _, err := net.SplitHostPort(host)
			if err == nil {
				s = h
			}
			if !tc.InsecureSkipVerify {
				if err := c1.VerifyHostname(s); err != nil {
					c1.Close()
					return nil, err
				}
			}
			c = c1
		}
		if tlsfingerprint.Client != "" {
			c1 := utls.UClient(c, &utls.Config{
				ServerName:         tc.ServerName,
				NextProtos:         tc.NextProtos,
				InsecureSkipVerify: tc.InsecureSkipVerify,
				RootCAs:            tc.RootCAs,
			}, tlsfingerprint)
			s := host
			h, _, err := net.SplitHostPort(host)
			if err == nil {
				s = h
			}
			if err := c1.BuildHandshakeState(); err != nil {
				return nil, err
			}
			for _, v := range c1.Extensions {
				if vv, ok := v.(*utls.ALPNExtension); ok {
					if tlsfingerprint.Client == "Chrome" {
						vv.AlpnProtocols = []string{"http/1.1"}
					}
					break
				}
			}
			if err := c1.BuildHandshakeState(); err != nil {
				return nil, err
			}
			if err := c1.Handshake(); err != nil {
				c1.Close()
				return nil, err
			}
			if !tc.InsecureSkipVerify {
				if err := c1.VerifyHostname(s); err != nil {
					c1.Close()
					return nil, err
				}
			}
			c = c1
		}
	}
	p := x1.BP16.Get().([]byte)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		x1.BP16.Put(p)
		c.Close()
		return nil, err
	}
	k := base64.StdEncoding.EncodeToString(p)
	x1.BP16.Put(p)
	b := make([]byte, 0, 300)
	b = append(b, []byte("GET "+path+" HTTP/1.1\r\n")...)
	b = append(b, []byte(fmt.Sprintf("Host: %s\r\n", host))...)
	b = append(b, []byte("Upgrade: websocket\r\n")...)
	b = append(b, []byte("Connection: Upgrade\r\n")...)
	b = append(b, []byte(fmt.Sprintf("Sec-WebSocket-Key: %s\r\n", k))...)
	b = append(b, []byte("Sec-WebSocket-Version: 13\r\n\r\n")...)
	if _, err := c.Write(b); err != nil {
		c.Close()
		return nil, err
	}
	r := bufio.NewReader(c)
	for {
		b, err = r.ReadBytes('\n')
		if err != nil {
			c.Close()
			return nil, err
		}
		b = bytes.TrimSpace(b)
		if len(b) == 0 {
			break
		}
		if bytes.HasPrefix(b, []byte("HTTP/1.1 ")) {
			if !bytes.Contains(b, []byte("101")) {
				c.Close()
				return nil, errors.New(string(b))
			}
		}
		if bytes.HasPrefix(b, []byte("Sec-WebSocket-Accept: ")) {
			h := sha1.New()
			h.Write([]byte(k))
			h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
			ak := base64.StdEncoding.EncodeToString(h.Sum(nil))
			if string(b[len("Sec-WebSocket-Accept: "):]) != ak {
				c.Close()
				return nil, errors.New(string(b))
			}
		}
	}
	return c, nil
}

type TLSFragmentConn struct {
	net.Conn
	MinLength int64
	MaxLength int64
	MinDelay  int64
	MaxDelay  int64
	Buf       []byte
	L         int
	Finished  bool
}

func (c *TLSFragmentConn) Write(b []byte) (int, error) {
	if c.Finished {
		return c.Conn.Write(b)
	}
	b1 := make([]byte, len(c.Buf)+len(b))
	copy(b1, c.Buf)
	copy(b1[len(c.Buf):], b)
	c.Buf = b1
	if len(c.Buf) < 5 {
		return len(b), nil
	}
	if c.L == 0 {
		c.L = int(binary.BigEndian.Uint16(c.Buf[3:5]))
	}
	if len(c.Buf) < 5+c.L {
		return len(b), nil
	}
	i := 0
	for {
		r, err := x.CryptoRandom(c.MinLength, c.MaxLength)
		if err != nil {
			return 0, err
		}
		l := int(r)
		if i+l > 5+c.L {
			l = 5 + c.L - i
		}
		if _, err := c.Conn.Write(c.Buf[i : i+l]); err != nil {
			return 0, err
		}
		i += l
		if i == 5+c.L {
			break
		}
		t, err := x.CryptoRandom(c.MinDelay, c.MaxDelay)
		if err != nil {
			return 0, err
		}
		time.Sleep(time.Duration(t) * time.Microsecond)
	}
	if len(c.Buf) > i {
		if _, err := c.Conn.Write(c.Buf[i:]); err != nil {
			return 0, err
		}
	}
	c.Finished = true
	return len(b), nil
}
