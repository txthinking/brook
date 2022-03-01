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
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/x"
)

type SimpleStreamClient struct {
	Server  net.Conn
	Timeout int
	Network string
	RB      []byte
	WB      []byte
}

func NewSimpleStreamClient(network string, password, dst []byte, server net.Conn, timeout int) (*SimpleStreamClient, error) {
	c := &SimpleStreamClient{Network: network, Server: server, Timeout: timeout}
	if len(dst) > 2048-32-2-4 {
		return nil, errors.New("dst too long")
	}
	b := x.BP2048.Get().([]byte)
	binary.BigEndian.PutUint16(b[32:32+2], uint16(4+len(dst)))
	i := time.Now().Unix()
	if c.Network == "tcp" && i%2 != 0 {
		i += 1
	}
	if c.Network == "udp" && i%2 != 1 {
		i += 1
	}
	binary.BigEndian.PutUint32(b[32+2:32+2+4], uint32(i))
	copy(b[:32], password)
	copy(b[32+2+4:], dst)
	if _, err := server.Write(b[:32+2+4+len(dst)]); err != nil {
		x.BP2048.Put(b)
		return nil, err
	}
	if c.Network == "tcp" {
		c.RB = b
		c.WB = x.BP2048.Get().([]byte)
	}
	if c.Network == "udp" {
		x.BP2048.Put(b)
		c.RB = x.BP65507.Get().([]byte)
		c.WB = x.BP65507.Get().([]byte)
	}
	return c, nil
}

func (c *SimpleStreamClient) Exchange(local net.Conn) error {
	go func() {
		if c.Timeout == 0 && c.Network == "tcp" {
			io.Copy(local, c.Server)
			return
		}
		for {
			if c.Timeout != 0 {
				if err := c.Server.SetDeadline(time.Now().Add(time.Duration(c.Timeout) * time.Second)); err != nil {
					return
				}
			}
			if c.Network == "tcp" {
				l, err := c.Server.Read(c.RB)
				if err != nil {
					return
				}
				if _, err := local.Write(c.RB[:l]); err != nil {
					return
				}
			}
			if c.Network == "udp" {
				if _, err := io.ReadFull(c.Server, c.RB[:2]); err != nil {
					return
				}
				l := int(binary.BigEndian.Uint16(c.RB[:2]))
				if l > 65507-2 {
					log.Println(errors.New("packet too long"))
					return
				}
				if _, err := io.ReadFull(c.Server, c.RB[2:2+l]); err != nil {
					return
				}
				if _, err := local.Write(c.RB[2 : 2+l]); err != nil {
					return
				}
			}
		}
	}()
	if c.Timeout == 0 && c.Network == "tcp" {
		io.Copy(c.Server, local)
		return nil
	}
	for {
		if c.Timeout != 0 {
			if err := local.SetDeadline(time.Now().Add(time.Duration(c.Timeout) * time.Second)); err != nil {
				return nil
			}
		}
		if c.Network == "tcp" {
			l, err := local.Read(c.WB)
			if err != nil {
				return nil
			}
			if _, err := c.Server.Write(c.WB[:l]); err != nil {
				return nil
			}
		}
		if c.Network == "udp" {
			l, err := local.Read(c.WB[2:])
			if err != nil {
				return nil
			}
			binary.BigEndian.PutUint16(c.WB[:2], uint16(l))
			if _, err := c.Server.Write(c.WB[:2+l]); err != nil {
				return nil
			}
		}
	}
	return nil
}

func (s *SimpleStreamClient) Clean() {
	if s.Network == "tcp" {
		x.BP2048.Put(s.WB)
		x.BP2048.Put(s.RB)
	}
	if s.Network == "udp" {
		x.BP65507.Put(s.WB)
		x.BP65507.Put(s.RB)
	}
}
func (s *SimpleStreamClient) NetworkName() string {
	return s.Network
}
func (s *SimpleStreamClient) SetTimeout(i int) {
	s.Timeout = i
}
