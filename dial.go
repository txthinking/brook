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
)

var Resolve func(network string, addr string) (net.Addr, error) = func(network string, addr string) (net.Addr, error) {
	if network == "tcp" {
		return net.ResolveTCPAddr("tcp", addr)
	}
	return net.ResolveUDPAddr("udp", addr)
}

var DialTCP func(network string, laddr, raddr string) (net.Conn, error) = func(network string, laddr, raddr string) (net.Conn, error) {
	var la, ra *net.TCPAddr
	if laddr != "" {
		var err error
		la, err = net.ResolveTCPAddr(network, laddr)
		if err != nil {
			return nil, err
		}
	}
	a, err := Resolve(network, raddr)
	if err != nil {
		return nil, err
	}
	ra = a.(*net.TCPAddr)
	return net.DialTCP(network, la, ra)
}

var DialUDP func(network string, laddr, raddr string) (net.Conn, error) = func(network string, laddr, raddr string) (net.Conn, error) {
	var la, ra *net.UDPAddr
	if laddr != "" {
		var err error
		la, err = net.ResolveUDPAddr(network, laddr)
		if err != nil {
			return nil, err
		}
	}
	a, err := Resolve(network, raddr)
	if err != nil {
		return nil, err
	}
	ra = a.(*net.UDPAddr)
	return net.DialUDP(network, la, ra)
}

var ListenUDP func(network string, laddr *net.UDPAddr) (*net.UDPConn, error) = net.ListenUDP
