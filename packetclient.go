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
	password []byte
	rb       []byte
	wb       []byte
}

func NewPacketClient(password []byte) *PacketClient {
	c := &PacketClient{password: password}
	c.wb = x.BP65507.Get().([]byte)
	c.rb = x.BP65507.Get().([]byte)
	return c
}

func (c *PacketClient) LocalToServer(dst, d []byte, server net.Conn, timeout int) error {
	if timeout != 0 {
		if err := server.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return err
		}
	}
	if 12+4+len(dst)+len(d)+16 > 65507 {
		return errors.New("udp data too long")
	}
	if _, err := io.ReadFull(rand.Reader, c.wb[:12]); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(c.wb[12:12+4], uint32(time.Now().Unix()))
	copy(c.wb[12+4:12+4+len(dst)], dst)
	copy(c.wb[12+4+len(dst):12+4+len(dst)+len(d)], d)
	ck := x.BP32.Get().([]byte)
	if _, err := io.ReadFull(hkdf.New(sha256.New, c.password, c.wb[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), ck); err != nil {
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
	ca.Seal(c.wb[:12], c.wb[:12], c.wb[12:12+4+len(dst)+len(d)], nil)
	_, err = server.Write(c.wb[:12+4+len(dst)+len(d)+16])
	if err != nil {
		return err
	}
	return nil
}

func (c *PacketClient) ServerToLocal(server net.Conn, timeout int, tolocal func(dst, d []byte) (int, error)) error {
	for {
		if timeout != 0 {
			if err := server.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
				return err
			}
		}
		i, err := server.Read(c.rb)
		if err != nil {
			return err
		}
		if i < 12+16 {
			return errors.New("data too small")
		}
		sk := x.BP32.Get().([]byte)
		if _, err := io.ReadFull(hkdf.New(sha256.New, c.password, c.rb[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), sk); err != nil {
			x.BP32.Put(sk)
			return err
		}
		sb, err := aes.NewCipher(sk)
		if err != nil {
			x.BP32.Put(sk)
			return err
		}
		x.BP32.Put(sk)
		sa, err := cipher.NewGCM(sb)
		if err != nil {
			return err
		}
		if _, err := sa.Open(c.rb[:12], c.rb[:12], c.rb[12:i], nil); err != nil {
			return err
		}
		_, h, _, err := socks5.ParseBytesAddress(c.rb[12:])
		if err != nil {
			return err
		}
		_, err = tolocal(c.rb[12:12+1+len(h)+2], c.rb[12+1+len(h)+2:i-16])
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *PacketClient) Clean() {
	x.BP65507.Put(c.wb)
	x.BP65507.Put(c.rb)
}
