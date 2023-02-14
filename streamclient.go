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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
	"golang.org/x/crypto/hkdf"
)

type StreamClient struct {
	Server  net.Conn
	cn      []byte
	ca      cipher.AEAD
	sn      []byte
	sa      cipher.AEAD
	RB      []byte
	WB      []byte
	Timeout int
	network string
	src     string
	dst     string
}

func NewStreamClient(network string, password []byte, src string, server net.Conn, timeout int, dst []byte) (Exchanger, error) {
	if timeout != 0 {
		if err := server.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return nil, err
		}
	}
	if len(dst) > 2048-2-16-4-16 {
		return nil, errors.New("dst too long")
	}
	c := &StreamClient{network: network, Server: server, Timeout: timeout, src: src, dst: socks5.ToAddress(dst[0], dst[1:len(dst)-2], dst[len(dst)-2:])}

	c.cn = x.BP12.Get().([]byte)
	if _, err := io.ReadFull(rand.Reader, c.cn); err != nil {
		x.BP12.Put(c.cn)
		return nil, err
	}
	ck := x.BP32.Get().([]byte)
	if _, err := io.ReadFull(hkdf.New(sha256.New, password, c.cn, []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), ck); err != nil {
		x.BP12.Put(c.cn)
		x.BP32.Put(ck)
		return nil, err
	}
	if _, err := c.Server.Write(c.cn); err != nil {
		x.BP12.Put(c.cn)
		x.BP32.Put(ck)
		return nil, err
	}
	cb, err := aes.NewCipher(ck)
	if err != nil {
		x.BP12.Put(c.cn)
		x.BP32.Put(ck)
		return nil, err
	}
	x.BP32.Put(ck)
	c.ca, err = cipher.NewGCM(cb)
	if err != nil {
		x.BP12.Put(c.cn)
		return nil, err
	}

	c.WB = x.BP2048.Get().([]byte)
	i := time.Now().Unix()
	if c.network == "tcp" && i%2 != 0 {
		i += 1
	}
	if c.network == "udp" && i%2 != 1 {
		i += 1
	}
	binary.BigEndian.PutUint32(c.WB[2+16:2+16+4], uint32(i))
	copy(c.WB[2+16+4:2+16+4+len(dst)], dst)
	if err := c.Write(4 + len(dst)); err != nil {
		x.BP12.Put(c.cn)
		x.BP2048.Put(c.WB)
		return nil, err
	}

	c.sn = x.BP12.Get().([]byte)
	if _, err := io.ReadFull(c.Server, c.sn); err != nil {
		x.BP12.Put(c.cn)
		x.BP2048.Put(c.WB)
		x.BP12.Put(c.sn)
		return nil, err
	}
	sk := x.BP32.Get().([]byte)
	if _, err := io.ReadFull(hkdf.New(sha256.New, password, c.sn, []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), sk); err != nil {
		x.BP12.Put(c.cn)
		x.BP2048.Put(c.WB)
		x.BP12.Put(c.sn)
		x.BP32.Put(sk)
		return nil, err
	}
	sb, err := aes.NewCipher(sk)
	if err != nil {
		x.BP12.Put(c.cn)
		x.BP2048.Put(c.WB)
		x.BP12.Put(c.sn)
		x.BP32.Put(sk)
		return nil, err
	}
	x.BP32.Put(sk)
	c.sa, err = cipher.NewGCM(sb)
	if err != nil {
		x.BP12.Put(c.cn)
		x.BP2048.Put(c.WB)
		x.BP12.Put(c.sn)
		return nil, err
	}
	if c.network == "tcp" {
		c.RB = x.BP2048.Get().([]byte)
	}
	if c.network == "udp" {
		x.BP2048.Put(c.WB)
		c.WB = x.BP65507.Get().([]byte)
		c.RB = x.BP65507.Get().([]byte)
	}
	return ClientGate(c)
}

func (c *StreamClient) Exchange(local net.Conn) error {
	go func() {
		for {
			if c.Timeout != 0 {
				if err := c.Server.SetDeadline(time.Now().Add(time.Duration(c.Timeout) * time.Second)); err != nil {
					return
				}
			}
			l, err := c.Read()
			if err != nil {
				return
			}
			if _, err := local.Write(c.RB[2+16 : 2+16+l]); err != nil {
				return
			}
		}
	}()
	for {
		if c.Timeout != 0 {
			if err := local.SetDeadline(time.Now().Add(time.Duration(c.Timeout) * time.Second)); err != nil {
				return nil
			}
		}
		l, err := local.Read(c.WB[2+16 : len(c.WB)-16])
		if err != nil {
			return nil
		}
		if err := c.Write(l); err != nil {
			return nil
		}
	}
	return nil
}

func (c *StreamClient) Write(l int) error {
	binary.BigEndian.PutUint16(c.WB[:2], uint16(l))
	c.ca.Seal(c.WB[:0], c.cn, c.WB[:2], nil)
	NextNonce(c.cn)
	c.ca.Seal(c.WB[:2+16], c.cn, c.WB[2+16:2+16+l], nil)
	if _, err := c.Server.Write(c.WB[:2+16+l+16]); err != nil {
		return err
	}
	NextNonce(c.cn)
	return nil
}

func (c *StreamClient) Read() (int, error) {
	if _, err := io.ReadFull(c.Server, c.RB[:2+16]); err != nil {
		return 0, err
	}
	if _, err := c.sa.Open(c.RB[:0], c.sn, c.RB[:2+16], nil); err != nil {
		return 0, err
	}
	l := int(binary.BigEndian.Uint16(c.RB[:2]))
	if _, err := io.ReadFull(c.Server, c.RB[2+16:2+16+l+16]); err != nil {
		return 0, err
	}
	NextNonce(c.sn)
	if _, err := c.sa.Open(c.RB[:2+16], c.sn, c.RB[2+16:2+16+l+16], nil); err != nil {
		return 0, err
	}
	NextNonce(c.sn)
	return l, nil
}

func (c *StreamClient) Clean() {
	x.BP12.Put(c.cn)
	x.BP12.Put(c.sn)
	if c.network == "tcp" {
		x.BP2048.Put(c.WB)
		x.BP2048.Put(c.RB)
	}
	if c.network == "udp" {
		x.BP65507.Put(c.WB)
		x.BP65507.Put(c.RB)
	}
}

func (s *StreamClient) Network() string {
	return s.network
}

func (s *StreamClient) Src() string {
	return s.src
}

func (s *StreamClient) Dst() string {
	return s.dst
}
