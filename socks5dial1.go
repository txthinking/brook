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

	"github.com/txthinking/socks5"
)

type Socks5Dial1 struct {
	Client *socks5.Client
}

func NewSocks5Dial1(server, username, password string, tcpTimeout, udpTimeout int) (*Socks5Dial1, error) {
	s5c, err := socks5.NewClient(server, username, password, tcpTimeout, udpTimeout)
	if err != nil {
		return nil, err
	}
	return &Socks5Dial1{
		Client: s5c,
	}, nil
}

func (d *Socks5Dial1) Dial(network, addr string) (net.Conn, error) {
	return d.Client.Dial(network, addr)
}

func (d *Socks5Dial1) DialTCP(network string, laddr, raddr *net.TCPAddr) (net.Conn, error) {
	src := ""
	if laddr != nil {
		src = laddr.String()
	}
	return d.Client.DialWithLocalAddr(network, src, raddr.String(), nil)
}

func (d *Socks5Dial1) DialUDP(network string, laddr, raddr *net.UDPAddr) (net.Conn, error) {
	src := ""
	if laddr != nil {
		src = laddr.String()
	}
	return d.Client.DialWithLocalAddr(network, src, raddr.String(), nil)
}
