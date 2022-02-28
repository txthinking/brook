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
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

// DNS.
type DNS struct {
	TCPAddr            *net.TCPAddr
	UDPAddr            *net.UDPAddr
	ServerTCPAddr      *net.TCPAddr
	ServerUDPAddr      *net.UDPAddr
	Password           []byte
	BypassDomains      map[string]byte
	DNSServer          string
	DNSServerForBypass string
	TCPListen          *net.TCPListener
	UDPConn            *net.UDPConn
	UDPExchanges       *cache.Cache
	TCPTimeout         int
	UDPTimeout         int
	RunnerGroup        *runnergroup.RunnerGroup
	UDPSrc             *cache.Cache
	WSClient           *WSClient
	BlockDomain        map[string]byte
	BypassCache        *cache.Cache
	BlockCache         *cache.Cache
}

// NewDNS.
func NewDNS(addr, server, password, dnsServer, dnsServerForBypass, bypassList string, tcpTimeout, udpTimeout int, blockDomainList string) (*DNS, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	var rtaddr *net.TCPAddr
	var ruaddr *net.UDPAddr
	if !strings.HasPrefix(server, "ws://") && !strings.HasPrefix(server, "wss://") {
		rtaddr, err = net.ResolveTCPAddr("tcp", server)
		if err != nil {
			return nil, err
		}
		ruaddr, err = net.ResolveUDPAddr("udp", server)
		if err != nil {
			return nil, err
		}
	}
	var wsc *WSClient
	if strings.HasPrefix(server, "ws://") || strings.HasPrefix(server, "wss://") {
		wsc, err = NewWSClient(":1080", "127.0.0.1", server, password, tcpTimeout, udpTimeout)
		if err != nil {
			return nil, err
		}
	}
	var ds, ds1 map[string]byte
	if bypassList != "" {
		ds, err = ReadDomainList(bypassList)
		if err != nil {
			return nil, err
		}
	}
	if blockDomainList != "" {
		ds1, err = ReadDomainList(blockDomainList)
		if err != nil {
			return nil, err
		}
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs2 := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs3 := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs4 := cache.New(cache.NoExpiration, cache.NoExpiration)
	s := &DNS{
		TCPAddr:            taddr,
		UDPAddr:            uaddr,
		ServerTCPAddr:      rtaddr,
		ServerUDPAddr:      ruaddr,
		UDPExchanges:       cs,
		Password:           []byte(password),
		BypassDomains:      ds,
		BlockDomain:        ds1,
		DNSServer:          dnsServer,
		DNSServerForBypass: dnsServerForBypass,
		TCPTimeout:         tcpTimeout,
		UDPTimeout:         udpTimeout,
		RunnerGroup:        runnergroup.New(),
		UDPSrc:             cs2,
		WSClient:           wsc,
		BypassCache:        cs3,
		BlockCache:         cs4,
	}
	return s, nil
}

// Run server.
func (s *DNS) ListenAndServe() error {
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
func (s *DNS) RunTCPServer() error {
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
func (s *DNS) RunUDPServer() error {
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
func (s *DNS) Shutdown() error {
	return s.RunnerGroup.Done()
}

// TCPHandle handles request.
func (s *DNS) TCPHandle(c *net.TCPConn) error {
	co := &dns.Conn{Conn: c}
	m, err := co.ReadMsg()
	if err != nil {
		return err
	}
	has := false
	for _, v := range m.Question {
		if Debug {
			log.Println("dns query", "tcp", v.Qtype, v.Name)
		}
		if (v.Qtype == dns.TypeAAAA || v.Qtype == dns.TypeA || v.Qtype == dns.TypeHTTPS) && len(v.Name) > 0 && ListHasDomain(s.BlockDomain, v.Name[0:len(v.Name)-1], s.BlockCache) {
			if v.Qtype == dns.TypeA || v.Qtype == dns.TypeHTTPS {
				m1 := &dns.Msg{}
				m1.SetReply(m)
				m1.Authoritative = true
				m1.Answer = append(m1.Answer, &dns.A{
					Hdr: dns.RR_Header{Name: v.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A:   net.IPv4zero,
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
			if v.Qtype == dns.TypeAAAA {
				m1 := &dns.Msg{}
				m1.SetReply(m)
				m1.Authoritative = true
				m1.Answer = append(m1.Answer, &dns.AAAA{
					Hdr:  dns.RR_Header{Name: v.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
					AAAA: net.IPv6zero,
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
		}
		if len(v.Name) > 0 && ListHasDomain(s.BypassDomains, v.Name[0:len(v.Name)-1], s.BypassCache) {
			has = true
			break
		}
	}
	mb, err := m.Pack()
	if err != nil {
		return err
	}
	lb := make([]byte, 2)
	binary.BigEndian.PutUint16(lb, uint16(len(mb)))
	mb = append(lb, mb...)
	if has {
		if Debug {
			log.Println("in bypass list", "tcp", m.Question[0].Name)
		}
		rc, err := Dial.Dial("tcp", s.DNSServerForBypass)
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
	var rc net.Conn
	if s.WSClient == nil {
		rc, err = Dial.Dial("tcp", s.ServerTCPAddr.String())
	}
	if s.WSClient != nil {
		rc, err = s.WSClient.DialWebsocket("")
	}
	if err != nil {
		return err
	}
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	a, h, p, err := socks5.ParseAddress(s.DNSServer)
	if err != nil {
		return err
	}
	dst := make([]byte, 0, 1+len(h)+2+len(mb))
	dst = append(dst, a)
	dst = append(dst, h...)
	dst = append(dst, p...)
	var sc Exchanger
	if s.WSClient == nil || !s.WSClient.WithoutBrook {
		sc, err = NewStreamClient("tcp", s.Password, dst, rc, s.TCPTimeout)
	}
	if s.WSClient != nil && s.WSClient.WithoutBrook {
		sc, err = NewSimpleStreamClient("tcp", s.WSClient.PasswordSha256, dst, rc, s.TCPTimeout)
	}
	if err != nil {
		return err
	}
	defer sc.Clean()
	if v, ok := sc.(*StreamClient); ok {
		i := copy(v.WB[2+16:], mb)
		if err := v.WriteL(i); err != nil {
			return err
		}
	}
	if _, ok := sc.(*SimpleStreamClient); ok {
		if _, err := rc.Write(mb); err != nil {
			return err
		}
	}
	if err := sc.Exchange(c); err != nil {
		return nil
	}
	return nil
}

// UDPHandle handles packet.
func (s *DNS) UDPHandle(addr *net.UDPAddr, b []byte) error {
	m := &dns.Msg{}
	if err := m.Unpack(b); err != nil {
		return err
	}
	has := false
	for _, v := range m.Question {
		if Debug {
			log.Println("dns query", "udp", v.Qtype, v.Name)
		}
		if (v.Qtype == dns.TypeAAAA || v.Qtype == dns.TypeA || v.Qtype == dns.TypeHTTPS) && len(v.Name) > 0 && ListHasDomain(s.BlockDomain, v.Name[0:len(v.Name)-1], s.BlockCache) {
			if v.Qtype == dns.TypeA || v.Qtype == dns.TypeHTTPS {
				m1 := &dns.Msg{}
				m1.SetReply(m)
				m1.Authoritative = true
				m1.Answer = append(m1.Answer, &dns.A{
					Hdr: dns.RR_Header{Name: v.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
					A:   net.IPv4zero,
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
			if v.Qtype == dns.TypeAAAA {
				m1 := &dns.Msg{}
				m1.SetReply(m)
				m1.Authoritative = true
				m1.Answer = append(m1.Answer, &dns.AAAA{
					Hdr:  dns.RR_Header{Name: v.Name, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: 60},
					AAAA: net.IPv6zero,
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
		}
		if len(v.Name) > 0 && ListHasDomain(s.BypassDomains, v.Name[0:len(v.Name)-1], s.BypassCache) {
			has = true
			break
		}
	}
	if has {
		if Debug {
			log.Println("in bypass list", "udp", m.Question[0].Name)
		}
		conn, err := Dial.Dial("udp", s.DNSServerForBypass)
		if err != nil {
			return err
		}
		defer conn.Close()
		if s.UDPTimeout != 0 {
			if err := conn.SetDeadline(time.Now().Add(time.Duration(s.UDPTimeout) * time.Second)); err != nil {
				return err
			}
		}
		co := &dns.Conn{Conn: conn}
		if err := co.WriteMsg(m); err != nil {
			return err
		}
		m1, err := co.ReadMsg()
		if err != nil {
			return err
		}
		m1b, err := m1.Pack()
		if err != nil {
			return err
		}
		if _, err := s.UDPConn.WriteToUDP(m1b, addr); err != nil {
			return err
		}
		return nil
	}

	src := addr.String()
	dst := s.DNSServer
	if s.WSClient == nil {
		any, ok := s.UDPExchanges.Get(src + dst)
		if ok {
			ue := any.(*UDPExchange)
			return ue.Any.(*PacketClient).LocalToServer(ue.Dst, b, ue.Conn, s.UDPTimeout)
		}
		if Debug {
			log.Println("dial udp", dst)
		}
		var laddr *net.UDPAddr
		any, ok = s.UDPSrc.Get(src + dst)
		if ok {
			laddr = any.(*net.UDPAddr)
		}
		rc, err := Dial.DialUDP("udp", laddr, s.ServerUDPAddr)
		if err != nil {
			return err
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
		a, h, p, err := socks5.ParseAddress(s.DNSServer)
		if err != nil {
			return err
		}
		dstb := make([]byte, 0, 1+len(h)+2)
		dstb = append(dstb, a)
		dstb = append(dstb, h...)
		dstb = append(dstb, p...)
		pc := NewPacketClient(s.Password)
		defer pc.Clean()
		if err := pc.LocalToServer(dstb, b, rc, s.UDPTimeout); err != nil {
			return err
		}
		ue := &UDPExchange{
			Conn: rc,
			Any:  pc,
			Dst:  dstb,
		}
		s.UDPExchanges.Set(src+dst, ue, -1)
		defer s.UDPExchanges.Delete(src + dst)
		err = pc.RunServerToLocal(rc, s.UDPTimeout, func(dst, d []byte) (int, error) {
			return s.UDPConn.WriteToUDP(d, addr)
		})
		if err != nil {
			return nil
		}
		return nil
	}
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		return ue.Any.(func(b []byte) error)(b)
	}
	if Debug {
		log.Println("dial udp", dst)
	}
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	la := ""
	if laddr != nil {
		la = laddr.String()
	}
	rc, err := s.WSClient.DialWebsocket(la)
	if err != nil {
		return err
	}
	defer rc.Close()
	if laddr == nil {
		laddr = &net.UDPAddr{
			IP:   rc.LocalAddr().(*net.TCPAddr).IP,
			Port: rc.LocalAddr().(*net.TCPAddr).Port,
			Zone: rc.LocalAddr().(*net.TCPAddr).Zone,
		}
		s.UDPSrc.Set(src+dst, laddr, -1)
	}
	if s.UDPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.UDPTimeout) * time.Second)); err != nil {
			return err
		}
	}

	a, h, p, err := socks5.ParseAddress(s.DNSServer)
	if err != nil {
		return err
	}
	dstb := make([]byte, 0, 1+len(h)+2)
	dstb = append(dstb, a)
	dstb = append(dstb, h...)
	dstb = append(dstb, p...)
	var sc Exchanger
	if !s.WSClient.WithoutBrook {
		sc, err = NewStreamClient("udp", s.Password, dstb, rc, s.UDPTimeout)
	}
	if s.WSClient.WithoutBrook {
		sc, err = NewSimpleStreamClient("udp", s.WSClient.PasswordSha256, dstb, rc, s.UDPTimeout)
	}
	if err != nil {
		return err
	}
	defer sc.Clean()
	ps, pi := NewPacketStream(func(b []byte) (int, error) {
		return s.UDPConn.WriteToUDP(b, addr)
	})
	defer ps.Close()
	ue := &UDPExchange{
		Any: pi,
	}
	s.UDPExchanges.Set(src+dst, ue, -1)
	defer s.UDPExchanges.Delete(src + dst)
	go pi(b)
	if err := sc.Exchange(ps); err != nil {
		return nil
	}
	return nil
}
