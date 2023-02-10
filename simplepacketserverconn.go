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
	"net"
	"sync"
	"time"

	"github.com/txthinking/socks5"
)

type SimplePacketServerConnFactory struct {
	Conns map[string]*PacketConn
	Lock  *sync.Mutex
}

func NewSimplePacketServerConnFactory() *SimplePacketServerConnFactory {
	return &SimplePacketServerConnFactory{
		Conns: make(map[string]*PacketConn),
		Lock:  &sync.Mutex{},
	}
}

func (f *SimplePacketServerConnFactory) Handle(addr *net.UDPAddr, b, p []byte, w func([]byte) (int, error), timeout int) (net.Conn, []byte, error) {
	if len(b) < 32+4 {
		return nil, nil, errors.New("data too small")
	}
	if bytes.Compare(p, b[:32]) != 0 {
		return nil, nil, errors.New("Password is wrong")
	}
	i := int64(binary.BigEndian.Uint32(b[32 : 32+4]))
	if time.Now().Unix()-i > 60 {
		return nil, nil, errors.New("Expired request")
	}
	a, h, p, err := socks5.ParseBytesAddress(b[32+4:])
	if err != nil {
		return nil, nil, err
	}
	dst := socks5.ToAddress(a, h, p)
	f.Lock.Lock()
	c, ok := f.Conns[addr.String()+dst]
	f.Lock.Unlock()
	if ok {
		_ = c.In(b[32+4+1+len(h)+2:])
		return nil, nil, nil
	}
	f.Lock.Lock()
	c = NewPacketConn(b[32+4+1+len(h)+2:], w, timeout, func() {
		f.Lock.Lock()
		delete(f.Conns, addr.String()+dst)
		f.Lock.Unlock()
	})
	f.Conns[addr.String()+dst] = c
	f.Lock.Unlock()
	return c, b[32+4 : 32+4+1+len(h)+2], nil
}
