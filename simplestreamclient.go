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
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
)

type SimpleStreamClient struct {
	Server  net.Conn
	Timeout int
	RB      []byte
	WB      []byte
	network string
	src     string
	dst     string
}

func NewSimpleStreamClient(network string, password []byte, src string, server net.Conn, timeout int, dst []byte) (Exchanger, error) {
	c := &SimpleStreamClient{network: network, Server: server, Timeout: timeout, src: src, dst: socks5.ToAddress(dst[0], dst[1:len(dst)-2], dst[len(dst)-2:])}
	if len(dst) > 2048-32-2-4 {
		return nil, errors.New("dst too long")
	}
	b := x.BP2048.Get().([]byte)
	binary.BigEndian.PutUint16(b[32:32+2], uint16(4+len(dst)))
	i := time.Now().Unix()
	if c.network == "tcp" && i%2 != 0 {
		i += 1
	}
	if c.network == "udp" && i%2 != 1 {
		i += 1
	}
	binary.BigEndian.PutUint32(b[32+2:32+2+4], uint32(i))
	copy(b[:32], password)
	copy(b[32+2+4:], dst)
	if _, err := server.Write(b[:32+2+4+len(dst)]); err != nil {
		x.BP2048.Put(b)
		return nil, err
	}
	if c.network == "tcp" {
		c.RB = b
		c.WB = x.BP2048.Get().([]byte)
	}
	if c.network == "udp" {
		x.BP2048.Put(b)
		c.RB = x.BP65507.Get().([]byte)
		c.WB = x.BP65507.Get().([]byte)
	}
	return ClientGate(c)
}

func (c *SimpleStreamClient) Exchange(local net.Conn) error {
	go func() {
		if c.Timeout == 0 && c.network == "tcp" {
			io.Copy(local, c.Server)
			return
		}
		for {
			if c.Timeout != 0 {
				if err := c.Server.SetDeadline(time.Now().Add(time.Duration(c.Timeout) * time.Second)); err != nil {
					return
				}
			}
			if c.network == "tcp" {
				l, err := c.Server.Read(c.RB)
				if err != nil {
					return
				}
				if _, err := local.Write(c.RB[:l]); err != nil {
					return
				}
			}
			if c.network == "udp" {
				if _, err := io.ReadFull(c.Server, c.RB[:2]); err != nil {
					return
				}
				l := int(binary.BigEndian.Uint16(c.RB[:2]))
				if l > 65507-2 {
					Log(&Error{"from": c.src, "dst": c.dst, "error": "read from server but packet too long"})
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
	if c.Timeout == 0 && c.network == "tcp" {
		io.Copy(c.Server, local)
		return nil
	}
	for {
		if c.Timeout != 0 {
			if err := local.SetDeadline(time.Now().Add(time.Duration(c.Timeout) * time.Second)); err != nil {
				return nil
			}
		}
		if c.network == "tcp" {
			l, err := local.Read(c.WB)
			if err != nil {
				return nil
			}
			if _, err := c.Server.Write(c.WB[:l]); err != nil {
				return nil
			}
		}
		if c.network == "udp" {
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
	if s.network == "tcp" {
		x.BP2048.Put(s.WB)
		x.BP2048.Put(s.RB)
	}
	if s.network == "udp" {
		x.BP65507.Put(s.WB)
		x.BP65507.Put(s.RB)
	}
}

func (s *SimpleStreamClient) Network() string {
	return s.network
}

func (s *SimpleStreamClient) Src() string {
	return s.src
}

func (s *SimpleStreamClient) Dst() string {
	return s.dst
}
