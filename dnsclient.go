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

type DNSClient struct {
	Server string
}

func (c *DNSClient) Exchange(m *dns.Msg) (*dns.Msg, error) {
	conn, err := DialUDP("udp", "", c.Server)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(time.Duration(60) * time.Second)); err != nil {
		return nil, err
	}
	b, err := m.Pack()
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write(b); err != nil {
		return nil, err
	}
	b = make([]byte, 1024)
	i, err := conn.Read(b)
	if err != nil {
		return nil, err
	}
	r := &dns.Msg{}
	if err := r.Unpack(b[:i]); err != nil {
		return nil, err
	}
	return r, nil
}

// if no AAAA, return nil
func (c *DNSClient) AAAA(domain string) (net.IP, error) {
	m := &dns.Msg{}
	m.SetQuestion(domain+".", dns.TypeAAAA)
	m, err := c.Exchange(m)
	if err != nil {
		return nil, err
	}
	for _, v := range m.Answer {
		if t, ok := v.(*dns.AAAA); ok {
			return t.AAAA, nil
		}
	}
	return nil, nil
}

// if no A, return nil
func (c *DNSClient) A(domain string) (net.IP, error) {
	m := &dns.Msg{}
	m.SetQuestion(domain+".", dns.TypeA)
	m, err := c.Exchange(m)
	if err != nil {
		return nil, err
	}
	for _, v := range m.Answer {
		if t, ok := v.(*dns.A); ok {
			return t.A, nil
		}
	}
	return nil, nil
}
