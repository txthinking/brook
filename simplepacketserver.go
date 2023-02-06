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
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
)

type SimplePacketServer struct {
	password []byte
	wb       []byte
}

func NewSimplePacketServer(password []byte) *SimplePacketServer {
	s := &SimplePacketServer{password: password}
	s.wb = x.BP65507.Get().([]byte)
	return s
}

func (s *SimplePacketServer) RemoteToClient(remote net.Conn, timeout int, dst []byte, toclient io.Writer) error {
	for {
		if timeout != 0 {
			if err := remote.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
				return err
			}
		}
		copy(s.wb[:len(dst)], dst)
		l, err := remote.Read(s.wb[len(dst):])
		if err != nil {
			return nil
		}
		_, err = toclient.Write(s.wb[:len(dst)+l])
		if err != nil {
			return err
		}
	}
	return nil
}

func SimplePacketClientToRemote(p, b []byte) ([]byte, []byte, WriterFunc, error) {
	if len(b) < 32+4 {
		return nil, nil, nil, errors.New("data too small")
	}
	if bytes.Compare(p, b[:32]) != 0 {
		return nil, nil, nil, errors.New("Password is wrong")
	}
	i := int64(binary.BigEndian.Uint32(b[32 : 32+4]))
	if time.Now().Unix()-i > 60 {
		return nil, nil, nil, errors.New("Expired request")
	}
	_, h, _, err := socks5.ParseBytesAddress(b[32+4:])
	if err != nil {
		return nil, nil, nil, err
	}
	return ServerPacket(b[32+4:32+4+1+len(h)+2], b[32+4+1+len(h)+2:])
}

func (s *SimplePacketServer) Clean() {
	x.BP65507.Put(s.wb)
}
