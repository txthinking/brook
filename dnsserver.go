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
	"time"

	"github.com/miekg/dns"
)

var DNSGate func(addr *net.UDPAddr, m *dns.Msg, l1 *net.UDPConn) (done bool, err error) = func(addr *net.UDPAddr, m *dns.Msg, l1 *net.UDPConn) (done bool, err error) {
	if m.Question[0].Qtype == dns.TypeHTTPS || m.Question[0].Qtype == dns.TypeSVCB {
		m1 := &dns.Msg{}
		m1.SetReply(m)
		m1.Authoritative = true
		m1.Answer = append(m1.Answer, &dns.SOA{
			Hdr:     dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 60},
			Ns:      "txthinking.com.",
			Mbox:    "cloud.txthinking.com.",
			Serial:  uint32((time.Now().Year() * 10000) + (int(time.Now().Month()) * 100) + (time.Now().Day())*100),
			Refresh: 21600,
			Retry:   3600,
			Expire:  259200,
			Minttl:  300,
		})
		m1b, err := m1.PackBuffer(nil)
		if err != nil {
			return false, err
		}
		if _, err := l1.WriteToUDP(m1b, addr); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
