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
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	x1 "github.com/txthinking/x"
)

func WebSocketDial(src, dst, addr, host, path string, tc *tls.Config, timeout int) (net.Conn, error) {
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
		c1 := tls.Client(c, tc)
		if !tc.InsecureSkipVerify {
			if err := c1.Handshake(); err != nil {
				c1.Close()
				return nil, err
			}
			s := host
			h, _, err := net.SplitHostPort(host)
			if err == nil {
				s = h
			}
			if err := c1.VerifyHostname(s); err != nil {
				c1.Close()
				return nil, err
			}
		}
		c = c1
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
