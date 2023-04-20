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

package dialwithnic

import (
	"errors"
	"net"

	"github.com/txthinking/brook"
	"github.com/txthinking/socks5"
)

type DialWithNIC struct {
	NIC string
}

func NewDialWithNIC(nic string) *DialWithNIC {
	d := &DialWithNIC{NIC: nic}
	return d
}

func (p *DialWithNIC) IP(v46 string) (net.IP, error) {
	ief, err := net.InterfaceByName(p.NIC)
	if err != nil {
		return nil, err
	}
	addrs, err := ief.Addrs()
	if err != nil {
		return nil, err
	}
	if v46 == "6" {
		for _, v := range addrs {
			if v.(*net.IPNet).IP.IsGlobalUnicast() && v.(*net.IPNet).IP.To4() == nil {
				return v.(*net.IPNet).IP, nil
			}
		}
		return nil, errors.New("no ipv6 from nic")
	}
	if v46 == "4" {
		for _, v := range addrs {
			if v.(*net.IPNet).IP.IsGlobalUnicast() && !v.(*net.IPNet).IP.IsPrivate() && v.(*net.IPNet).IP.To4() != nil {
				return v.(*net.IPNet).IP, nil
			}
		}
		for _, v := range addrs {
			if v.(*net.IPNet).IP.IsGlobalUnicast() && v.(*net.IPNet).IP.IsPrivate() && v.(*net.IPNet).IP.To4() != nil {
				return v.(*net.IPNet).IP, nil
			}
		}
		return nil, errors.New("no ipv4 from nic")
	}
	return nil, errors.New("black hole")
}

func (p *DialWithNIC) TouchBrook() {
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
			v46 := "6"
			if ra.IP.To4() != nil {
				v46 = "4"
			}
			ip, err := p.IP(v46)
			if err != nil {
				return nil, err
			}
			la = &net.TCPAddr{IP: ip}
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
			v46 := "6"
			if ra.IP.To4() != nil {
				v46 = "4"
			}
			ip, err := p.IP(v46)
			if err != nil {
				return nil, err
			}
			la = &net.UDPAddr{IP: ip}
		}
		return net.DialUDP(network, la, ra)
	}
	socks5.DialTCP = brook.DialTCP
	socks5.DialUDP = brook.DialUDP
}
