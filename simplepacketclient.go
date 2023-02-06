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
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
)

type SimplePacketClient struct {
	Server        net.Conn
	Password      []byte
	RB            []byte
	WB            []byte
	Dst           []byte
	Timeout       int
	RemoteAddress net.Addr
}

func NewSimplePacketClient(password []byte) *SimplePacketClient {
	c := &SimplePacketClient{Password: password}
	c.WB = x.BP65507.Get().([]byte)
	c.RB = x.BP65507.Get().([]byte)
	return c
}

func (c *SimplePacketClient) LocalToServer(dst, d []byte, server net.Conn, timeout int) error {
	if 32+4+len(dst)+len(d) > 65507 {
		return errors.New("packet too big")
	}
	if timeout != 0 {
		if err := server.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return err
		}
	}
	copy(c.WB[0:32], c.Password)
	binary.BigEndian.PutUint32(c.WB[32:32+4], uint32(time.Now().Unix()))
	copy(c.WB[32+4:32+4+len(dst)], dst)
	copy(c.WB[32+4+len(dst):32+4+len(dst)+len(d)], d)
	_, err := server.Write(c.WB[:32+4+len(dst)+len(d)])
	if err != nil {
		return err
	}
	return nil
}

func (c *SimplePacketClient) RunServerToLocal(server net.Conn, timeout int, tolocal func(dst, d []byte) (int, error)) error {
	for {
		if err := c.ServerToLocal(server, timeout, tolocal); err != nil {
			return err
		}
	}
	return nil
}

func (c *SimplePacketClient) ServerToLocal(server net.Conn, timeout int, tolocal func(dst, d []byte) (int, error)) error {
	if timeout != 0 {
		if err := server.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return err
		}
	}
	i, err := server.Read(c.RB)
	if err != nil {
		return err
	}
	_, h, _, err := socks5.ParseBytesAddress(c.RB[0:])
	if err != nil {
		return err
	}
	_, err = tolocal(c.RB[:1+len(h)+2], c.RB[1+len(h)+2:i])
	if err != nil {
		return err
	}
	return nil
}

func (c *SimplePacketClient) Clean() {
	x.BP65507.Put(c.WB)
	x.BP65507.Put(c.RB)
}
