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

type PacketClient struct {
	Server   net.Conn
	Password []byte
	RB       []byte
	WB       []byte
	Timeout  int
	src      string
	dstb     []byte
}

func NewPacketClient(password []byte, src string, server net.Conn, timeout int, dstb []byte) (Exchanger, error) {
	s := &PacketClient{Password: password, Server: server, Timeout: timeout, src: src}
	s.RB = x.BP65507.Get().([]byte)
	s.WB = x.BP65507.Get().([]byte)
	s.dstb = dstb
	return ClientGate(s)
}

func (c *PacketClient) Exchange(local net.Conn) error {
	go func() {
		for {
			if c.Timeout != 0 {
				if err := c.Server.SetDeadline(time.Now().Add(time.Duration(c.Timeout) * time.Second)); err != nil {
					return
				}
			}
			i, err := c.Server.Read(c.RB)
			if err != nil {
				return
			}
			if i < 12+16 {
				Log(errors.New("data too small"))
				return
			}
			sk := x.BP32.Get().([]byte)
			if _, err := io.ReadFull(hkdf.New(sha256.New, c.Password, c.RB[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), sk); err != nil {
				x.BP32.Put(sk)
				Log(err)
				return
			}
			sb, err := aes.NewCipher(sk)
			if err != nil {
				x.BP32.Put(sk)
				Log(err)
				return
			}
			x.BP32.Put(sk)
			sa, err := cipher.NewGCM(sb)
			if err != nil {
				Log(err)
				return
			}
			if _, err := sa.Open(c.RB[:12], c.RB[:12], c.RB[12:i], nil); err != nil {
				Log(err)
				return
			}
			_, h, _, err := socks5.ParseBytesAddress(c.RB[12:])
			if err != nil {
				Log(err)
				return
			}
			_, err = local.Write(c.RB[12+1+len(h)+2 : i-16])
			if err != nil {
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
		copy(c.WB[12+4:12+4+len(c.dstb)], c.dstb)
		l, err := local.Read(c.WB[12+4+len(c.dstb) : 65507-16])
		if err != nil {
			return nil
		}
		if _, err := io.ReadFull(rand.Reader, c.WB[:12]); err != nil {
			return err
		}
		binary.BigEndian.PutUint32(c.WB[12:12+4], uint32(time.Now().Unix()))
		ck := x.BP32.Get().([]byte)
		if _, err := io.ReadFull(hkdf.New(sha256.New, c.Password, c.WB[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), ck); err != nil {
			x.BP32.Put(ck)
			return err
		}
		cb, err := aes.NewCipher(ck)
		if err != nil {
			x.BP32.Put(ck)
			return err
		}
		x.BP32.Put(ck)
		ca, err := cipher.NewGCM(cb)
		if err != nil {
			return err
		}
		ca.Seal(c.WB[:12], c.WB[:12], c.WB[12:12+4+len(c.dstb)+l], nil)
		_, err = c.Server.Write(c.WB[:12+4+len(c.dstb)+l+16])
		if err != nil {
			return nil
		}
	}
	return nil
}

func (s *PacketClient) Clean() {
	x.BP65507.Put(s.RB)
	x.BP65507.Put(s.WB)
}

func (s *PacketClient) Network() string {
	return "udp"
}

func (s *PacketClient) Src() string {
	return s.src
}

func (s *PacketClient) Dst() string {
	return socks5.ToAddress(s.dstb[0], s.dstb[1:len(s.dstb)-2], s.dstb[len(s.dstb)-2:])
}
