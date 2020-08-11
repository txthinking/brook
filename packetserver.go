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

type PacketServer struct {
	password []byte
	wb       []byte
}

func NewPacketServer(password []byte) *PacketServer {
	s := &PacketServer{password: password}
	s.wb = x.BP65507.Get().([]byte)
	return s
}

func (s *PacketServer) RemoteToClient(remote net.Conn, timeout int, dst []byte, toclient func(b []byte) (int, error)) error {
	for {
		if timeout != 0 {
			if err := remote.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
				return err
			}
		}
		l, err := remote.Read(s.wb[12+len(dst) : 65507-16])
		if _, err := io.ReadFull(rand.Reader, s.wb[:12]); err != nil {
			return err
		}
		sk := x.BP32.Get().([]byte)
		if _, err := io.ReadFull(hkdf.New(sha256.New, s.password, s.wb[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), sk); err != nil {
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
		copy(s.wb[12:12+len(dst)], dst)
		sa.Seal(s.wb[:12], s.wb[:12], s.wb[12:12+len(dst)+l], nil)
		_, err = toclient(s.wb[:12+len(dst)+l+16])
		if err != nil {
			return err
		}
	}
	return nil
}

func PacketClientToRemote(p, b []byte) ([]byte, []byte, error) {
	if len(b) < 12+4+16 {
		return nil, nil, errors.New("data too small")
	}
	ck := x.BP32.Get().([]byte)
	if _, err := io.ReadFull(hkdf.New(sha256.New, p, b[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), ck); err != nil {
		x.BP32.Put(ck)
		return nil, nil, err
	}

	cb, err := aes.NewCipher(ck)
	if err != nil {
		x.BP32.Put(ck)
		return nil, nil, err
	}
	x.BP32.Put(ck)
	ca, err := cipher.NewGCM(cb)
	if err != nil {
		return nil, nil, err
	}
	if _, err := ca.Open(b[:12], b[:12], b[12:], nil); err != nil {
		return nil, nil, err
	}
	i := int64(binary.BigEndian.Uint32(b[12 : 12+4]))
	if time.Now().Unix()-i > 60 {
		return nil, nil, errors.New("Expired request")
	}
	_, h, _, err := socks5.ParseBytesAddress(b[12+4:])
	if err != nil {
		return nil, nil, err
	}
	return b[12+4 : 12+4+1+len(h)+2], b[12+4+1+len(h)+2 : len(b)-16], nil
}

func (s *PacketServer) Clean() {
	x.BP65507.Put(s.wb)
}
