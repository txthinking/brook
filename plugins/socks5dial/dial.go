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

package socks5dial

import (
	"net"

	"github.com/txthinking/brook"
	"github.com/txthinking/socks5"
)

type Socks5Dial struct {
	s5c *socks5.Client
}

func NewSocks5Dial(server, username, password string, tcptimeout, udptimeout int) (*Socks5Dial, error) {
	s5c, err := socks5.NewClient(server, username, password, tcptimeout, udptimeout)
	if err != nil {
		return nil, err
	}
	return &Socks5Dial{
		s5c: s5c,
	}, nil
}

func (p *Socks5Dial) TouchBrook() {
	brook.DialTCP = func(network string, laddr, raddr string) (net.Conn, error) {
		var fake net.Addr
		if network == "tcp" {
			fake = &net.TCPAddr{IP: net.IPv4zero}
		}
		if network == "udp" {
			fake = &net.UDPAddr{IP: net.IPv4zero}
		}
		return p.s5c.DialWithLocalAddr("tcp", laddr, raddr, fake)
	}
	brook.DialUDP = func(network string, laddr, raddr string) (net.Conn, error) {
		var fake net.Addr
		if network == "tcp" {
			fake = &net.TCPAddr{IP: net.IPv4zero}
		}
		if network == "udp" {
			fake = &net.UDPAddr{IP: net.IPv4zero}
		}
		return p.s5c.DialWithLocalAddr("udp", laddr, raddr, fake)
	}
}
