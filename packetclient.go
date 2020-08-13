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
	Server        net.Conn
	Password      []byte
	RB            []byte
	WB            []byte
	Dst           []byte
	Timeout       int
	RemoteAddress net.Addr
}

func NewPacketClient(password []byte) *PacketClient {
	c := &PacketClient{Password: password}
	c.WB = x.BP65507.Get().([]byte)
	c.RB = x.BP65507.Get().([]byte)
	return c
}

var ClientPacket func([]byte, []byte) ([]byte, []byte, error) = func(dst, d []byte) ([]byte, []byte, error) {
	if 12+4+len(dst)+len(d)+16 > 65507 {
		return nil, nil, errors.New("packet too big")
	}
	return dst, d, nil
}

func (c *PacketClient) LocalToServer(dst, d []byte, server net.Conn, timeout int) error {
	dst, d, err := ClientPacket(dst, d)
	if err != nil {
		return err
	}
	if timeout != 0 {
		if err := server.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return err
		}
	}
	if _, err := io.ReadFull(rand.Reader, c.WB[:12]); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(c.WB[12:12+4], uint32(time.Now().Unix()))
	copy(c.WB[12+4:12+4+len(dst)], dst)
	copy(c.WB[12+4+len(dst):12+4+len(dst)+len(d)], d)
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
	ca.Seal(c.WB[:12], c.WB[:12], c.WB[12:12+4+len(dst)+len(d)], nil)
	_, err = server.Write(c.WB[:12+4+len(dst)+len(d)+16])
	if err != nil {
		return err
	}
	return nil
}

func (c *PacketClient) RunServerToLocal(server net.Conn, timeout int, tolocal func(dst, d []byte) (int, error)) error {
	for {
		if err := c.ServerToLocal(server, timeout, tolocal); err != nil {
			return err
		}
	}
	return nil
}

func (c *PacketClient) ServerToLocal(server net.Conn, timeout int, tolocal func(dst, d []byte) (int, error)) error {
	if timeout != 0 {
		if err := server.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return err
		}
	}
	i, err := server.Read(c.RB)
	if err != nil {
		return err
	}
	if i < 12+16 {
		return errors.New("data too small")
	}
	sk := x.BP32.Get().([]byte)
	if _, err := io.ReadFull(hkdf.New(sha256.New, c.Password, c.RB[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), sk); err != nil {
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
	if _, err := sa.Open(c.RB[:12], c.RB[:12], c.RB[12:i], nil); err != nil {
		return err
	}
	_, h, _, err := socks5.ParseBytesAddress(c.RB[12:])
	if err != nil {
		return err
	}
	_, err = tolocal(c.RB[12:12+1+len(h)+2], c.RB[12+1+len(h)+2:i-16])
	if err != nil {
		return err
	}
	return nil
}

func (c *PacketClient) Clean() {
	x.BP65507.Put(c.WB)
	x.BP65507.Put(c.RB)
}
