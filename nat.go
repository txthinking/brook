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

type NATTable struct {
	Table map[string]string
	Lock  *sync.Mutex
}

var NAT = &NATTable{
	Table: map[string]string{},
	Lock:  &sync.Mutex{},
}

func (n *NATTable) Set(src, dst, addr string) {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	n.Table[src+dst] = addr
}

func (n *NATTable) Get(src, dst string) string {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	s, _ := n.Table[src+dst]
	return s
}

func (n *NATTable) Reset() {
	n.Lock.Lock()
	defer n.Lock.Unlock()
	n.Table = map[string]string{}
}

var NATDial func(network string, src, dst, addr string) (net.Conn, error) = func(network string, src, dst, addr string) (net.Conn, error) {
	s := NAT.Get(src, dst)
	var c net.Conn
	var err error
	if network == "tcp" {
		c, err = DialTCP(network, s, addr)
	}
	if network == "udp" {
		c, err = DialUDP(network, s, addr)
	}
	if err != nil {
		if !strings.Contains(err.Error(), "address already in use") && !strings.Contains(err.Error(), "assign requested address") {
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
		NAT.Set(src, dst, c.LocalAddr().String())
	}
	return c, nil
}

var NATListenUDP func(network string, src, dst string) (*net.UDPConn, error) = func(network string, src, dst string) (*net.UDPConn, error) {
	var laddr *net.UDPAddr
	s := NAT.Get(src, dst)
	if s != "" {
		var err error
		laddr, err = net.ResolveUDPAddr("udp", s)
		if err != nil {
			return nil, err
		}
	}
	c, err := ListenUDP("udp", laddr)
	if err != nil {
		if !strings.Contains(err.Error(), "address already in use") && !strings.Contains(err.Error(), "assign requested address") {
			return nil, err
		}
		c, err = ListenUDP("udp", nil)
		s = ""
	}
	if err != nil {
		return nil, err
	}
	if s == "" {
		NAT.Set(src, dst, c.LocalAddr().String())
	}
	return c, nil
}
