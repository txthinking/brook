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

package tproxy

import (
	"net"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook"
)

type Tproxy struct {
	CIDR4      []*net.IPNet
	CIDR6      []*net.IPNet
	GeoIP      []string
	Cache      *cache.Cache
	TCPTimeout int
	UDPTimeout int
}

func NewTproxy(cidr4List, cidr6List string, geoIP []string, tcptimeout, udptimeout int) (*Tproxy, error) {
	var err error
	var c4 []*net.IPNet
	if cidr4List != "" {
		c4, err = brook.ReadCIDRList(cidr4List)
		if err != nil {
			return nil, err
		}
	}
	var c6 []*net.IPNet
	if cidr6List != "" {
		c6, err = brook.ReadCIDRList(cidr6List)
		if err != nil {
			return nil, err
		}
	}
	b := &Tproxy{
		CIDR4:      c4,
		CIDR6:      c6,
		GeoIP:      geoIP,
		Cache:      cache.New(cache.NoExpiration, cache.NoExpiration),
		TCPTimeout: tcptimeout,
		UDPTimeout: udptimeout,
	}
	return b, nil
}

func (p *Tproxy) TouchBrook() {
	f := brook.TproxyGate
	brook.TproxyGate = func(conn net.Conn) (net.Conn, error) {
		var ip net.IP
		network := "tcp"
		timeout := p.TCPTimeout
		size := 1024 * 2
		a, ok := conn.LocalAddr().(*net.TCPAddr)
		if ok {
			ip = a.IP
		}
		if ip == nil {
			ip = conn.LocalAddr().(*net.UDPAddr).IP
			network = "udp"
			timeout = p.UDPTimeout
			size = 65507
		}
		if brook.ListHasIP(p.CIDR4, p.CIDR6, ip, p.Cache, p.GeoIP) {
			var rc net.Conn
			var err error
			if network == "tcp" {
				rc, err = brook.DialTCP(network, "", conn.LocalAddr().String())
			}
			if network == "udp" {
				rc, err = brook.NATDial(network, conn.RemoteAddr().String(), conn.LocalAddr().String(), conn.LocalAddr().String())
			}
			if err != nil {
				return nil, err
			}
			brook.Conn2Conn(conn, rc, size, timeout)
			return nil, nil
		}
		return f(conn)
	}
}
