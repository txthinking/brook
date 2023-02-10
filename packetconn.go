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
	"errors"
	"net"
	"sync"
	"time"

	"github.com/txthinking/socks5"
)

type PacketConnFactory struct {
	Conns map[string]*PacketConn
	Lock  *sync.Mutex
}

func NewPacketConnFactory() *PacketConnFactory {
	return &PacketConnFactory{
		Conns: make(map[string]*PacketConn),
		Lock:  &sync.Mutex{},
	}
}

func (f *PacketConnFactory) Handle(addr *net.UDPAddr, dstb, data []byte, w func([]byte) (int, error), timeout int) (net.Conn, error) {
	dst := socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:])
	f.Lock.Lock()
	c, ok := f.Conns[addr.String()+dst]
	f.Lock.Unlock()
	if ok {
		_ = c.In(data)
		return nil, nil
	}
	f.Lock.Lock()
	c = NewPacketConn(data, w, timeout, func() {
		f.Lock.Lock()
		delete(f.Conns, addr.String()+dst)
		f.Lock.Unlock()
	})
	f.Conns[addr.String()+dst] = c
	f.Lock.Unlock()
	return c, nil
}

type PacketConn struct {
	First   []byte
	InCh    chan []byte
	Done    chan byte
	W       func([]byte) (int, error)
	Clean   func()
	Timeout int
}

func NewPacketConn(fb []byte, w func([]byte) (int, error), timeout int, clean func()) *PacketConn {
	c := &PacketConn{
		InCh:    make(chan []byte),
		Done:    make(chan byte),
		W:       w,
		First:   fb,
		Clean:   clean,
		Timeout: timeout,
	}
	return c
}

func (c *PacketConn) In(b []byte) error {
	select {
	case c.InCh <- b:
		return nil
	case <-c.Done:
		return errors.New("closed")
	}
	return nil
}

func (c *PacketConn) Read(b []byte) (int, error) {
	if c.First != nil {
		i := copy(b, c.First)
		c.First = nil
		return i, nil
	}
	tm := time.NewTimer(time.Duration(c.Timeout) * time.Second)
	defer tm.Stop()
	select {
	case <-tm.C:
		return 0, errors.New("timeout")
	case bb := <-c.InCh:
		i := copy(b, bb)
		return i, nil
	case <-c.Done:
		return 0, errors.New("closed")
	}
	return 0, nil
}

func (c *PacketConn) Write(b []byte) (int, error) {
	select {
	case <-c.Done:
		return 0, errors.New("closed")
	default:
		return c.W(b)
	}
	return 0, nil
}

func (c *PacketConn) Close() error {
	select {
	case <-c.Done:
	default:
		c.Clean()
		close(c.Done)
	}
	return nil
}

func (c *PacketConn) LocalAddr() net.Addr {
	panic("no")
	return nil
}

func (c *PacketConn) RemoteAddr() net.Addr {
	panic("no")
	return nil
}

func (c *PacketConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *PacketConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *PacketConn) SetWriteDeadline(t time.Time) error {
	return nil
}

type ConnFirst struct {
	*net.UDPConn
	First []byte
}

func (c *ConnFirst) Read(b []byte) (int, error) {
	if c.First != nil {
		i := copy(b, c.First)
		c.First = nil
		return i, nil
	}
	return c.UDPConn.Read(b)
}
