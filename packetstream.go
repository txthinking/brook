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
	"time"
)

type PacketStream struct {
	in   chan []byte
	done chan byte
	w    func(b []byte) (int, error)
}

func NewPacketStream(w func(d []byte) (int, error)) (*PacketStream, func(b []byte) error) {
	c := &PacketStream{
		in:   make(chan []byte),
		done: make(chan byte),
		w:    w,
	}
	f := func(b []byte) error {
		select {
		case c.in <- b:
			return nil
		case <-c.done:
			return errors.New("closed")
		}
		return nil
	}
	return c, f
}

func (c *PacketStream) Read(b []byte) (int, error) {
	select {
	case bb := <-c.in:
		i := copy(b, bb)
		return i, nil
	case <-c.done:
		return 0, errors.New("closed")
	}
	return 0, nil
}

func (c *PacketStream) Write(b []byte) (int, error) {
	if _, err := c.w(b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *PacketStream) Close() error {
	close(c.done)
	return nil
}

func (c *PacketStream) LocalAddr() net.Addr {
	panic("no addr")
}

func (c *PacketStream) RemoteAddr() net.Addr {
	panic("no addr")
}

func (c *PacketStream) SetDeadline(t time.Time) error {
	return nil
}

func (c *PacketStream) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *PacketStream) SetWriteDeadline(t time.Time) error {
	return nil
}
