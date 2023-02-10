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
	"log"
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
	dst      string
	dstl     int
}

func NewPacketServer(password []byte, src string, client net.Conn, timeout int, dst []byte) (Exchanger, error) {
	s := &PacketServer{Password: password, Client: client, Timeout: timeout, src: src}
	s.RB = x.BP65507.Get().([]byte)
	s.WB = x.BP65507.Get().([]byte)
	s.dstl = copy(s.WB[12:12+len(dst)], dst)
	s.dst = socks5.ToAddress(dst[0], dst[1:s.dstl-2], dst[s.dstl-2:])
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
			l, err := remote.Read(s.WB[12+s.dstl : 65507-16])
			if err != nil {
				return
			}
			if _, err := io.ReadFull(rand.Reader, s.WB[:12]); err != nil {
				log.Println(err)
				return
			}
			sk := x.BP32.Get().([]byte)
			if _, err := io.ReadFull(hkdf.New(sha256.New, s.Password, s.WB[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), sk); err != nil {
				x.BP32.Put(sk)
				log.Println(err)
				return
			}
			sb, err := aes.NewCipher(sk)
			if err != nil {
				x.BP32.Put(sk)
				log.Println(err)
				return
			}
			x.BP32.Put(sk)
			sa, err := cipher.NewGCM(sb)
			if err != nil {
				log.Println(err)
				return
			}
			sa.Seal(s.WB[:12], s.WB[:12], s.WB[12:12+s.dstl+l], nil)
			_, err = s.Client.Write(s.WB[:12+s.dstl+l+16])
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
	return s.dst
}
