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
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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
}

// NewDNS.
func NewDNS(addr, server, password, dnsServer, dnsServerForBypass, bypassList string, tcpTimeout, udpTimeout int) (*DNS, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	rtaddr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		return nil, err
	}
	ruaddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return nil, err
	}
	ds := make(map[string]byte)
	ss := make([]string, 0)
	if bypassList != "" {
		ss, err = readList(bypassList)
		if err != nil {
			return nil, err
		}
	}
	for _, v := range ss {
		ds[v] = 0
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs2 := cache.New(cache.NoExpiration, cache.NoExpiration)
	s := &DNS{
		TCPAddr:            taddr,
		UDPAddr:            uaddr,
		ServerTCPAddr:      rtaddr,
		ServerUDPAddr:      ruaddr,
		UDPExchanges:       cs,
		Password:           []byte(password),
		BypassDomains:      ds,
		DNSServer:          dnsServer,
		DNSServerForBypass: dnsServerForBypass,
		TCPTimeout:         tcpTimeout,
		UDPTimeout:         udpTimeout,
		RunnerGroup:        runnergroup.New(),
		UDPSrc:             cs2,
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
		debug("dns query", "udp", v.Qtype, v.Name)
		if len(v.Name) > 0 && s.Has(v.Name[0:len(v.Name)-1]) {
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
		debug("in bypass list", "tcp", m.Question[0].Name)
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
	rc, err := Dial.Dial("tcp", s.ServerTCPAddr.String())
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
	sc, err := NewStreamClient("tcp", s.Password, dst, rc, s.TCPTimeout)
	if err != nil {
		return err
	}
	defer sc.Clean()
	i := copy(sc.WB[2+16:], mb)
	if err := sc.WriteL(i); err != nil {
		return err
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
		debug("dns query", "udp", v.Qtype, v.Name)
		if len(v.Name) > 0 && s.Has(v.Name[0:len(v.Name)-1]) {
			has = true
			break
		}
	}
	if has {
		debug("in bypass list", "udp", m.Question[0].Name)
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
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		return ue.Any.(*PacketClient).LocalToServer(ue.Dst, b, ue.Conn, s.UDPTimeout)
	}
	debug("dial udp", dst)
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	rc, err := Dial.DialUDP("udp", laddr, s.ServerUDPAddr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			// we dont choose lock, so ignore this error
			return nil
		}
		return err
	}
	defer rc.Close()
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

func (s *DNS) Has(host string) bool {
	ss := strings.Split(host, ".")
	var s1 string
	for i := len(ss) - 1; i >= 0; i-- {
		if s1 == "" {
			s1 = ss[i]
		} else {
			s1 = ss[i] + "." + s1
		}
		if _, ok := s.BypassDomains[s1]; ok {
			return true
		}
	}
	return false
}

func readList(url string) ([]string, error) {
	var data []byte
	var err error
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		c := &http.Client{
			Timeout: 9 * time.Second,
		}
		r, err := c.Get(url)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		data, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		data, err = ioutil.ReadFile(url)
		if err != nil {
			return nil, err
		}
	}
	data = bytes.TrimSpace(data)
	data = bytes.Replace(data, []byte{0x20}, []byte{}, -1)
	data = bytes.Replace(data, []byte{0x0d, 0x0a}, []byte{0x0a}, -1)
	ss := strings.Split(string(data), "\n")
	return ss, nil
}
