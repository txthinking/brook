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

package dialwithip

import (
	"errors"
	"net"

	"github.com/txthinking/brook"
	"github.com/txthinking/socks5"
)

type DialWithIP struct {
	IP4 net.IP
	IP6 net.IP
}

func NewDialWithIP(ip4, ip6 string) (*DialWithIP, error) {
	if ip4 != "" && (net.ParseIP(ip4) == nil || net.ParseIP(ip4).To4() == nil) {
		return nil, errors.New("Invalid dial with IP")
	}
	if ip6 != "" && (net.ParseIP(ip6) == nil || net.ParseIP(ip6).To4() != nil) {
		return nil, errors.New("Invalid dial with IP")
	}
	d := &DialWithIP{}
	if ip4 != "" {
		d.IP4 = net.ParseIP(ip4).To4()
	}
	if ip6 != "" {
		d.IP6 = net.ParseIP(ip6).To16()
	}
	return d, nil
}

func (p *DialWithIP) TouchBrook() {
	brook.DialTCP = func(network string, laddr, raddr string) (net.Conn, error) {
		var la, ra *net.TCPAddr
		if laddr != "" {
			var err error
			la, err = net.ResolveTCPAddr(network, laddr)
			if err != nil {
				return nil, err
			}
		}
		a, err := brook.Resolve(network, raddr)
		if err != nil {
			return nil, err
		}
		ra = a.(*net.TCPAddr)
		if la == nil {
			if ra.IP.To4() != nil && p.IP4 != nil {
				la = &net.TCPAddr{IP: p.IP4}
			}
			if ra.IP.To4() == nil && p.IP6 != nil {
				la = &net.TCPAddr{IP: p.IP6}
			}
		}
		return net.DialTCP(network, la, ra)
	}
	brook.DialUDP = func(network string, laddr, raddr string) (net.Conn, error) {
		var la, ra *net.UDPAddr
		if laddr != "" {
			var err error
			la, err = net.ResolveUDPAddr(network, laddr)
			if err != nil {
				return nil, err
			}
		}
		a, err := brook.Resolve(network, raddr)
		if err != nil {
			return nil, err
		}
		ra = a.(*net.UDPAddr)
		if la == nil {
			if ra.IP.To4() != nil && p.IP4 != nil {
				la = &net.UDPAddr{IP: p.IP4}
			}
			if ra.IP.To4() == nil && p.IP6 != nil {
				la = &net.UDPAddr{IP: p.IP6}
			}
		}
		return net.DialUDP(network, la, ra)
	}
	socks5.DialTCP = brook.DialTCP
	socks5.DialUDP = brook.DialUDP
}
