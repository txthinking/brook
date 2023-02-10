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
	"strings"
	"sync"
)

type nattable struct {
	table map[string]string
	lock  *sync.Mutex
}

var _nattable = &nattable{
	table: map[string]string{},
	lock:  &sync.Mutex{},
}

func (n *nattable) set(src, dst, addr string) {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.table[src+dst] = addr
}

func (n *nattable) get(src, dst string) string {
	n.lock.Lock()
	defer n.lock.Unlock()
	s, _ := n.table[src+dst]
	return s
}

var NATDial func(network string, src, dst, addr string) (net.Conn, error) = func(network string, src, dst, addr string) (net.Conn, error) {
	s := _nattable.get(src, dst)
	var c net.Conn
	var err error
	if network == "tcp" {
		c, err = DialTCP(network, s, addr)
	}
	if network == "udp" {
		c, err = DialUDP(network, s, addr)
	}
	if err != nil {
		if !strings.Contains(err.Error(), "address already in use") {
			return nil, err
		}
		if network == "tcp" {
			c, err = DialTCP(network, "", addr)
		}
		if network == "udp" {
			c, err = DialUDP(network, "", addr)
		}
		s = ""
	}
	if err != nil {
		return nil, err
	}
	if s == "" {
		_nattable.set(src, dst, c.LocalAddr().String())
	}
	return c, nil
}

var NATListenUDP func(network string, src, dst string) (*net.UDPConn, error) = func(network string, src, dst string) (*net.UDPConn, error) {
	var laddr *net.UDPAddr
	s := _nattable.get(src, dst)
	if s != "" {
		a, err := Resolve("udp", s)
		if err != nil {
			return nil, err
		}
		laddr = a.(*net.UDPAddr)
	}
	c, err := ListenUDP("udp", laddr)
	if err != nil {
		if !strings.Contains(err.Error(), "address already in use") {
			return nil, err
		}
		c, err = ListenUDP("udp", nil)
		s = ""
	}
	if err != nil {
		return nil, err
	}
	if s == "" {
		_nattable.set(src, dst, c.LocalAddr().String())
	}
	return c, nil
}
