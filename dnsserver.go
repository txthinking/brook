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
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	cache "github.com/patrickmn/go-cache"
	"github.com/phuslu/iploc"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
)

type DNSServer struct {
	TCPAddr           *net.TCPAddr
	UDPAddr           *net.UDPAddr
	DNS               string
	DisableIPv4Domain map[string]byte
	DisableIPv6Domain map[string]byte
	BlockDomain       map[string]byte
	BlockGeo          []string
	DisableIPv4Cache  *cache.Cache
	DisableIPv6Cache  *cache.Cache
	BlockCache        *cache.Cache
	TCPListen         *net.TCPListener
	UDPConn           *net.UDPConn
	UDPExchanges      *cache.Cache
	UDPSrc            *cache.Cache
	TCPTimeout        int
	UDPTimeout        int
	RunnerGroup       *runnergroup.RunnerGroup
}

// NewDNS.
func NewDNSServer(addr, dns, disableIPv4DomainList, disableIPv6DomainList, blockDomainList string, blockGeo []string, tcpTimeout, udpTimeout int) (*DNSServer, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	var d4, d6, d map[string]byte
	if disableIPv4DomainList != "" {
		d4, err = ReadDomainList(disableIPv4DomainList)
		if err != nil {
			return nil, err
		}
	}
	if disableIPv6DomainList != "" {
		d6, err = ReadDomainList(disableIPv6DomainList)
		if err != nil {
			return nil, err
		}
	}
	if blockDomainList != "" {
		d, err = ReadDomainList(blockDomainList)
		if err != nil {
			return nil, err
		}
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &DNSServer{
		TCPAddr:           taddr,
		UDPAddr:           uaddr,
		DNS:               dns,
		DisableIPv4Domain: d4,
		DisableIPv6Domain: d6,
		BlockDomain:       d,
		BlockGeo:          blockGeo,
		DisableIPv4Cache:  cache.New(cache.NoExpiration, cache.NoExpiration),
		DisableIPv6Cache:  cache.New(cache.NoExpiration, cache.NoExpiration),
		BlockCache:        cache.New(cache.NoExpiration, cache.NoExpiration),
		UDPExchanges:      cache.New(cache.NoExpiration, cache.NoExpiration),
		UDPSrc:            cache.New(cache.NoExpiration, cache.NoExpiration),
		TCPTimeout:        tcpTimeout,
		UDPTimeout:        udpTimeout,
		RunnerGroup:       runnergroup.New(),
	}
	return s, nil
}

// Run server.
func (s *DNSServer) ListenAndServe() error {
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			return s.RunTCPServer()
		},
		Stop: func() error {
			if s.TCPListen != nil {
				return s.TCPListen.Close()
			}
			return nil
		},
	})
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			return s.RunUDPServer()
		},
		Stop: func() error {
			if s.UDPConn != nil {
				return s.UDPConn.Close()
			}
			return nil
		},
	})
	return s.RunnerGroup.Wait()
}

// RunTCPServer starts tcp server.
func (s *DNSServer) RunTCPServer() error {
	var err error
	s.TCPListen, err = net.ListenTCP("tcp", s.TCPAddr)
	if err != nil {
		return err
	}
	defer s.TCPListen.Close()
	for {
		c, err := s.TCPListen.AcceptTCP()
		if err != nil {
			return err
		}
		go func(c *net.TCPConn) {
			defer c.Close()
			if s.TCPTimeout != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
					log.Println(err)
					return
				}
			}
			if err := s.TCPHandle(c); err != nil {
				log.Println(err)
			}
		}(c)
	}
	return nil
}

// RunUDPServer starts udp server.
func (s *DNSServer) RunUDPServer() error {
	var err error
	s.UDPConn, err = net.ListenUDP("udp", s.UDPAddr)
	if err != nil {
		return err
	}
	defer s.UDPConn.Close()
	for {
		b := make([]byte, 65507)
		n, addr, err := s.UDPConn.ReadFromUDP(b)
		if err != nil {
			return err
		}
		go func(addr *net.UDPAddr, b []byte) {
			if err := s.UDPHandle(addr, b); err != nil {
				log.Println(err)
				return
			}
		}(addr, b[0:n])
	}
	return nil
}

// Shutdown server.
func (s *DNSServer) Shutdown() error {
	return s.RunnerGroup.Done()
}

// TCPHandle handles request.
func (s *DNSServer) TCPHandle(c *net.TCPConn) error {
	co := &dns.Conn{Conn: c}
	m, err := co.ReadMsg()
	if err != nil {
		return err
	}
	if len(m.Question) == 0 {
		return errors.New("no question")
	}
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
			return err
		}
		lb := make([]byte, 2)
		binary.BigEndian.PutUint16(lb, uint16(len(m1b)))
		m1b = append(lb, m1b...)
		if _, err := c.Write(m1b); err != nil {
			return err
		}
		return nil
	}
	if len(m.Question[0].Name) > 0 && ListHasDomain(s.BlockDomain, m.Question[0].Name[0:len(m.Question[0].Name)-1], s.BlockCache) {
		if Debug {
			log.Println("DNS", "BLOCK", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String())
		}
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
		lb := make([]byte, 2)
		binary.BigEndian.PutUint16(lb, uint16(len(m1b)))
		m1b = append(lb, m1b...)
		if _, err := c.Write(m1b); err != nil {
			return err
		}
		return nil
	}
	if m.Question[0].Qtype == dns.TypeA && len(m.Question[0].Name) > 0 && ListHasDomain(s.DisableIPv4Domain, m.Question[0].Name[0:len(m.Question[0].Name)-1], s.DisableIPv4Cache) {
		if Debug {
			log.Println("DNS", "BLOCK", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String())
		}
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
		lb := make([]byte, 2)
		binary.BigEndian.PutUint16(lb, uint16(len(m1b)))
		m1b = append(lb, m1b...)
		if _, err := c.Write(m1b); err != nil {
			return err
		}
		return nil
	}
	if m.Question[0].Qtype == dns.TypeAAAA && len(m.Question[0].Name) > 0 && ListHasDomain(s.DisableIPv6Domain, m.Question[0].Name[0:len(m.Question[0].Name)-1], s.DisableIPv6Cache) {
		if Debug {
			log.Println("DNS", "BLOCK", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String())
		}
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
		lb := make([]byte, 2)
		binary.BigEndian.PutUint16(lb, uint16(len(m1b)))
		m1b = append(lb, m1b...)
		if _, err := c.Write(m1b); err != nil {
			return err
		}
		return nil
	}

	mb, err := m.Pack()
	if err != nil {
		return err
	}
	lb := make([]byte, 2)
	binary.BigEndian.PutUint16(lb, uint16(len(mb)))
	mb = append(lb, mb...)

	if Debug {
		log.Println("DNS", "QUERY", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String())
	}
	rc, err := Dial.Dial("tcp", s.DNS)
	if err != nil {
		return err
	}
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	if _, err := rc.Write(mb); err != nil {
		return err
	}
	go func() {
		var bf [1024 * 2]byte
		for {
			if s.TCPTimeout != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
					return
				}
			}
			i, err := rc.Read(bf[:])
			if err != nil {
				return
			}
			if _, err := c.Write(bf[0:i]); err != nil {
				return
			}
		}
	}()
	var bf [1024 * 2]byte
	for {
		if s.TCPTimeout != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
				return nil
			}
		}
		i, err := c.Read(bf[:])
		if err != nil {
			return nil
		}
		if _, err := rc.Write(bf[0:i]); err != nil {
			return nil
		}
	}
	return nil
}

// UDPHandle handles packet.
func (s *DNSServer) UDPHandle(addr *net.UDPAddr, b []byte) error {
	m := &dns.Msg{}
	if err := m.Unpack(b); err != nil {
		return err
	}
	if len(m.Question) == 0 {
		return errors.New("no question")
	}
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
			return err
		}
		if _, err := s.UDPConn.WriteToUDP(m1b, addr); err != nil {
			return err
		}
		return nil
	}
	if len(m.Question[0].Name) > 0 && ListHasDomain(s.BlockDomain, m.Question[0].Name[0:len(m.Question[0].Name)-1], s.BlockCache) {
		if Debug {
			log.Println("DNS", "BLOCK", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String())
		}
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
		if _, err := s.UDPConn.WriteToUDP(m1b, addr); err != nil {
			return err
		}
		return nil
	}
	if m.Question[0].Qtype == dns.TypeA && len(m.Question[0].Name) > 0 && ListHasDomain(s.DisableIPv4Domain, m.Question[0].Name[0:len(m.Question[0].Name)-1], s.DisableIPv4Cache) {
		if Debug {
			log.Println("DNS", "BLOCK", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String())
		}
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
		if _, err := s.UDPConn.WriteToUDP(m1b, addr); err != nil {
			return err
		}
		return nil
	}
	if m.Question[0].Qtype == dns.TypeAAAA && len(m.Question[0].Name) > 0 && ListHasDomain(s.DisableIPv6Domain, m.Question[0].Name[0:len(m.Question[0].Name)-1], s.DisableIPv6Cache) {
		if Debug {
			log.Println("DNS", "BLOCK", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String())
		}
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
		if _, err := s.UDPConn.WriteToUDP(m1b, addr); err != nil {
			return err
		}
		return nil
	}

	src := addr.String()
	dst := s.DNS
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		if _, err := ue.Conn.Write(b); err != nil {
			return err
		}
		return nil
	}
	if Debug {
		log.Println("DNS", "QUERY", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String())
	}
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	raddr, err := net.ResolveUDPAddr("udp", s.DNS)
	if err != nil {
		return err
	}
	rc, err := Dial.DialUDP("udp", laddr, raddr)
	if err != nil {
		if !strings.Contains(err.Error(), "address already in use") {
			return err
		}
		rc, err = Dial.DialUDP("udp", nil, raddr)
		if err != nil {
			return err
		}
		laddr = nil
	}
	defer rc.Close()
	if s.UDPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.UDPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	if laddr == nil {
		s.UDPSrc.Set(src+dst, rc.LocalAddr().(*net.UDPAddr), -1)
	}
	if _, err := rc.Write(b); err != nil {
		return err
	}
	ue := &UDPExchange{
		Conn: rc,
	}
	s.UDPExchanges.Set(src+dst, ue, -1)
	defer s.UDPExchanges.Delete(src + dst)
	var b0 [65507]byte
	for {
		if s.UDPTimeout != 0 {
			if err := rc.SetDeadline(time.Now().Add(time.Duration(s.UDPTimeout) * time.Second)); err != nil {
				break
			}
		}
		n, err := rc.Read(b0[:])
		if err != nil {
			break
		}
		got := false
		if len(s.BlockGeo) > 0 {
			m := &dns.Msg{}
			if err := m.Unpack(b0[0:n]); err != nil {
				return err
			}
			if len(m.Question) > 0 && (m.Question[0].Qtype == dns.TypeAAAA || m.Question[0].Qtype == dns.TypeA) {
			GEO:
				for _, v := range m.Answer {
					if t, ok := v.(*dns.A); ok {
						for _, vv := range s.BlockGeo {
							if vv == string(iploc.Country(t.A.To4())) {
								if Debug {
									log.Println("DNS", "BLOCK", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String(), t.A.To4())
								}
								got = true
								break GEO
							}
						}
					}
					if t, ok := v.(*dns.AAAA); ok {
						for _, vv := range s.BlockGeo {
							if vv == string(iploc.Country(t.AAAA.To16())) {
								if Debug {
									log.Println("DNS", "BLOCK", m.Question[0].Name, dns.Type(m.Question[0].Qtype).String(), t.AAAA.To16())
								}
								got = true
								break GEO
							}
						}
					}
				}
			}
		}
		if got {
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
			if _, err := s.UDPConn.WriteToUDP(m1b, addr); err != nil {
				return err
			}
			continue
		}
		if _, err := s.UDPConn.WriteToUDP(b0[0:n], addr); err != nil {
			break
		}
	}
	return nil
}
