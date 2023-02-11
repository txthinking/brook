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

package thedns

import (
	"net"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/txthinking/brook"
)

type TheDNS struct {
	BlockDomain  map[string]byte
	BypassDomain map[string]byte
	BypassDNS    string
	DisableA     bool
	DisableAAAA  bool
	Cache        *cache.Cache
}

func NewTheDNS(blockDomainList, bypassDomainList, bypassDNS string, disableA, disableAAAA bool) (*TheDNS, error) {
	var err error
	var ds map[string]byte
	if blockDomainList != "" {
		ds, err = brook.ReadDomainList(blockDomainList)
		if err != nil {
			return nil, err
		}
	}
	var ds1 map[string]byte
	if bypassDomainList != "" {
		ds1, err = brook.ReadDomainList(bypassDomainList)
		if err != nil {
			return nil, err
		}
	}
	b := &TheDNS{
		BlockDomain:  ds,
		BypassDomain: ds1,
		BypassDNS:    bypassDNS,
		DisableA:     disableA,
		DisableAAAA:  disableAAAA,
		Cache:        cache.New(cache.NoExpiration, cache.NoExpiration),
	}
	return b, nil
}

func soa(addr *net.UDPAddr, m *dns.Msg, l1 *net.UDPConn) error {
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
		return err
	}
	if _, err := l1.WriteToUDP(m1b, addr); err != nil {
		return err
	}
	return nil
}

func (p *TheDNS) TouchBrook() {
	f := brook.DNSGate
	brook.DNSGate = func(addr *net.UDPAddr, m *dns.Msg, l1 *net.UDPConn) (bool, error) {
		done, err := f(addr, m, l1)
		if err != nil || done {
			return done, err
		}
		if m.Question[0].Qtype == dns.TypeA && p.DisableA {
			err := soa(addr, m, l1)
			return err == nil, err
		}
		if m.Question[0].Qtype == dns.TypeAAAA && p.DisableAAAA {
			err := soa(addr, m, l1)
			return err == nil, err
		}
		if brook.ListHasDomain(p.BlockDomain, m.Question[0].Name[0:len(m.Question[0].Name)-1], p.Cache) {
			err := soa(addr, m, l1)
			return err == nil, err
		}
		if brook.ListHasDomain(p.BypassDomain, m.Question[0].Name[0:len(m.Question[0].Name)-1], p.Cache) {
			m1, err := dns.Exchange(m, p.BypassDNS)
			if err != nil {
				return false, err
			}
			m1b, err := m1.PackBuffer(nil)
			if err != nil {
				return false, err
			}
			if _, err := l1.WriteToUDP(m1b, addr); err != nil {
				return false, err
			}
			return true, nil
		}
		return f(addr, m, l1)
	}
}
