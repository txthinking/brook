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
	"io"
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
	"golang.org/x/crypto/hkdf"
)

type PacketServer struct {
	Client   net.Conn
	Password []byte
	RB       []byte
	WB       []byte
	Timeout  int
	src      string
	dstb     []byte
}

func NewPacketServer(password []byte, src string, client net.Conn, timeout int, dstb []byte) (Exchanger, error) {
	s := &PacketServer{Password: password, Client: client, Timeout: timeout, src: src}
	s.RB = x.BP65507.Get().([]byte)
	s.WB = x.BP65507.Get().([]byte)
	s.dstb = dstb
	return ServerGate(s)
}

func (s *PacketServer) Exchange(remote net.Conn) error {
	go func() {
		for {
			if s.Timeout != 0 {
				if err := remote.SetDeadline(time.Now().Add(time.Duration(s.Timeout) * time.Second)); err != nil {
					return
				}
			}
			copy(s.WB[12:12+len(s.dstb)], s.dstb)
			l, err := remote.Read(s.WB[12+len(s.dstb) : 65507-16])
			if err != nil {
				return
			}
			if _, err := io.ReadFull(rand.Reader, s.WB[:12]); err != nil {
				Log(err)
				return
			}
			sk := x.BP32.Get().([]byte)
			if _, err := io.ReadFull(hkdf.New(sha256.New, s.Password, s.WB[:12], ServerHKDFInfo), sk); err != nil {
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
			sa.Seal(s.WB[:12], s.WB[:12], s.WB[12:12+len(s.dstb)+l], nil)
			_, err = s.Client.Write(s.WB[:12+len(s.dstb)+l+16])
			if err != nil {
				return
			}
		}
	}()
	for {
		if s.Timeout != 0 {
			if err := s.Client.SetDeadline(time.Now().Add(time.Duration(s.Timeout) * time.Second)); err != nil {
				return nil
			}
		}
		l, err := s.Client.Read(s.RB)
		if err != nil {
			return nil
		}
		if _, err := remote.Write(s.RB[:l]); err != nil {
			return nil
		}
	}
	return nil
}

func (s *PacketServer) Clean() {
	x.BP65507.Put(s.RB)
	x.BP65507.Put(s.WB)
}

func (s *PacketServer) Network() string {
	return "udp"
}

func (s *PacketServer) Src() string {
	return s.src
}

func (s *PacketServer) Dst() string {
	return socks5.ToAddress(s.dstb[0], s.dstb[1:len(s.dstb)-2], s.dstb[len(s.dstb)-2:])
}
