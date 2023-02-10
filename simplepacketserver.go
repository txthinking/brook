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
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
)

type SimplePacketServer struct {
	Client   net.Conn
	Password []byte
	RB       []byte
	WB       []byte
	Timeout  int
	src      string
	dst      string
	dstl     int
}

func NewSimplePacketServer(password []byte, src string, client net.Conn, timeout int, dst []byte) (Exchanger, error) {
	s := &SimplePacketServer{Password: password, Client: client, Timeout: timeout, src: src}
	s.RB = x.BP65507.Get().([]byte)
	s.WB = x.BP65507.Get().([]byte)
	s.dstl = copy(s.WB[:len(dst)], dst)
	s.dst = socks5.ToAddress(dst[0], dst[1:s.dstl-2], dst[s.dstl-2:])
	return ServerGate(s)
}

func (s *SimplePacketServer) Exchange(remote net.Conn) error {
	go func() {
		for {
			if s.Timeout != 0 {
				if err := remote.SetDeadline(time.Now().Add(time.Duration(s.Timeout) * time.Second)); err != nil {
					return
				}
			}
			l, err := remote.Read(s.WB[s.dstl:])
			if err != nil {
				return
			}
			_, err = s.Client.Write(s.WB[:s.dstl+l])
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

func (s *SimplePacketServer) Clean() {
	x.BP65507.Put(s.RB)
	x.BP65507.Put(s.WB)
}

func (s *SimplePacketServer) Network() string {
	return "udp"
}

func (s *SimplePacketServer) Src() string {
	return s.src
}

func (s *SimplePacketServer) Dst() string {
	return s.dst
}
