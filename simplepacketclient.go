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
	"encoding/binary"
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
)

type SimplePacketClient struct {
	Server   net.Conn
	Password []byte
	RB       []byte
	WB       []byte
	Timeout  int
	src      string
	dst      string
	dstl     int
}

func NewSimplePacketClient(password []byte, src string, server net.Conn, timeout int, dst []byte) (Exchanger, error) {
	s := &SimplePacketClient{Password: password, Server: server, Timeout: timeout, src: src}
	s.RB = x.BP65507.Get().([]byte)
	s.WB = x.BP65507.Get().([]byte)
	s.dstl = copy(s.WB[32+4:32+4+len(dst)], dst)
	s.dst = socks5.ToAddress(dst[0], dst[1:s.dstl-2], dst[s.dstl-2:])
	return ClientGate(s)
}

func (c *SimplePacketClient) Exchange(local net.Conn) error {
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
			_, h, _, err := socks5.ParseBytesAddress(c.RB[0:])
			if err != nil {
				Log(err)
				return
			}
			_, err = local.Write(c.RB[1+len(h)+2 : i])
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
		l, err := local.Read(c.WB[32+4+c.dstl:])
		if err != nil {
			return nil
		}
		copy(c.WB[0:32], c.Password)
		binary.BigEndian.PutUint32(c.WB[32:32+4], uint32(time.Now().Unix()))
		_, err = c.Server.Write(c.WB[:32+4+c.dstl+l])
		if err != nil {
			return nil
		}
	}
	return nil
}

func (s *SimplePacketClient) Clean() {
	x.BP65507.Put(s.RB)
	x.BP65507.Put(s.WB)
}

func (s *SimplePacketClient) Network() string {
	return "udp"
}

func (s *SimplePacketClient) Src() string {
	return s.src
}

func (s *SimplePacketClient) Dst() string {
	return s.dst
}
