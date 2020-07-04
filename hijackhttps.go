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
	"context"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
	"golang.org/x/net/proxy"
)

// HijackHTTPS.
type HijackHTTPS struct {
	Dial         proxy.Dialer
	ListenDNSIP  string
	DefaultDNS   string
	Domains      map[string]byte
	UDPDNSServer *dns.Server
	TCPDNSServer *dns.Server
	HTTPSServer  net.Listener
	TCPDeadline  int
	TCPTimeout   int
	UDPDeadline  int
	RunnerGroup  *runnergroup.RunnerGroup
}

// NewHijackHTTPS.
func NewHijackHTTPS(socks5addr, socks5username, socks5password, listendnsip, defaultDNS, list string, tcpTimeout, tcpDeadline, udpDeadline int) (*HijackHTTPS, error) {
	ds := make(map[string]byte)
	if list != "" {
		ss, err := readList(list)
		if err != nil {
			return nil, err
		}
		for _, v := range ss {
			ds[v] = 0
		}
	}
	h, p, err := net.SplitHostPort(socks5addr)
	if err != nil {
		return nil, err
	}
	l, err := net.LookupIP(h)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, errors.New("Can't find socks5 server IP")
	}
	var auth *proxy.Auth
	if socks5username != "" || socks5password != "" {
		auth = &proxy.Auth{
			User:     socks5username,
			Password: socks5password,
		}
	}
	dial, err := proxy.SOCKS5("tcp", net.JoinHostPort(l[0].String(), p), auth, Dial)
	if err != nil {
		return nil, err
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &HijackHTTPS{
		Dial:        dial,
		ListenDNSIP: listendnsip,
		DefaultDNS:  defaultDNS,
		Domains:     ds,
		TCPTimeout:  tcpTimeout,
		TCPDeadline: tcpDeadline,
		UDPDeadline: udpDeadline,
		RunnerGroup: runnergroup.New(),
	}
	return s, nil
}

// Run server.
func (s *HijackHTTPS) ListenAndServe() error {
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			return s.RunUDPDNSServer()
		},
		Stop: func() error {
			if s.UDPDNSServer != nil {
				ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
				defer cancel()
				return s.UDPDNSServer.ShutdownContext(ctx)
			}
			return nil
		},
	})
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			return s.RunTCPDNSServer()
		},
		Stop: func() error {
			if s.TCPDNSServer != nil {
				return s.TCPDNSServer.Shutdown()
			}
			return nil
		},
	})
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			return s.RunHTTPSServer()
		},
		Stop: func() error {
			if s.HTTPSServer != nil {
				return s.HTTPSServer.Close()
			}
			return nil
		},
	})
	err := s.RunnerGroup.Wait()
	return err
}

func (s *HijackHTTPS) RunUDPDNSServer() error {
	s.UDPDNSServer = &dns.Server{
		Addr:         net.JoinHostPort(s.ListenDNSIP, "53"),
		Net:          "udp",
		ReadTimeout:  time.Duration(s.UDPDeadline) * time.Second,
		WriteTimeout: time.Duration(s.UDPDeadline) * time.Second,
	}
	s.UDPDNSServer.Handler = s.DNSHandle("udp")
	return s.UDPDNSServer.ListenAndServe()
}

func (s *HijackHTTPS) RunTCPDNSServer() error {
	s.TCPDNSServer = &dns.Server{
		Addr:         net.JoinHostPort(s.ListenDNSIP, "53"),
		Net:          "tcp",
		ReadTimeout:  time.Duration(s.TCPTimeout) * time.Second,
		WriteTimeout: time.Duration(s.TCPTimeout) * time.Second,
	}
	s.TCPDNSServer.Handler = s.DNSHandle("tcp")
	return s.TCPDNSServer.ListenAndServe()
}

func (s *HijackHTTPS) DNSHandle(network string) dns.Handler {
	return dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
		has := false
		for _, v := range r.Question {
			debug("dns query", network, v.Qtype, v.Name)
		}
		if len(r.Question) > 0 && r.Question[0].Qtype == dns.TypeA {
			for _, v := range r.Question {
				if len(v.Name) > 0 && s.Has(v.Name[0:len(v.Name)-1]) {
					has = true
					break
				}
			}
		}
		if !has {
			debug("dns hijack", r.Question[0].Name)
			m := &dns.Msg{}
			m.SetReply(r)
			m.Authoritative = true
			m.Answer = append(m.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(s.ListenDNSIP),
			})
			w.WriteMsg(m)
			return
		}
		conn, err := Dial.Dial(network, s.DefaultDNS)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		if network == "tcp" && s.TCPDeadline != 0 {
			if err := conn.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
				log.Println(err)
				return
			}
		}
		if network == "udp" && s.UDPDeadline != 0 {
			if err := conn.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
				log.Println(err)
				return
			}
		}
		co := &dns.Conn{Conn: conn}
		if err := co.WriteMsg(r); err != nil {
			log.Println(err)
			return
		}
		m, err := co.ReadMsg()
		if err != nil {
			log.Println(err)
			return
		}
		w.WriteMsg(m)
	})
}

func (s *HijackHTTPS) RunHTTPSServer() error {
	var err error
	s.HTTPSServer, err = net.Listen("tcp", net.JoinHostPort(s.ListenDNSIP, "443"))
	if err != nil {
		return err
	}
	for {
		c, err := s.HTTPSServer.Accept()
		if err != nil {
			return err
		}
		go func(c net.Conn) {
			defer c.Close()
			if s.TCPTimeout != 0 {
				if err := c.(*net.TCPConn).SetKeepAlivePeriod(time.Duration(s.TCPTimeout) * time.Second); err != nil {
					log.Println(err)
					return
				}
			}
			if s.TCPDeadline != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
					log.Println(err)
					return
				}
			}
			if err := s.HTTPSHandle(c); err != nil {
				log.Println(err)
				return
			}
		}(c)
	}
}

func (s *HijackHTTPS) HTTPSHandle(c net.Conn) error {
	dm := ""
	d := make([]byte, 0)

	b := make([]byte, 9)
	if _, err := io.ReadFull(c, b); err != nil {
		return err
	}
	d = append(d, b...)

	l := int(binary.BigEndian.Uint32(append([]byte{0x00}, b[6:9]...)))
	b = make([]byte, l)
	if _, err := io.ReadFull(c, b); err != nil {
		return err
	}
	d = append(d, b...)

	if len(b) < 1+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	b = b[1:]
	if len(b) < 1+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	b = b[1:]
	if len(b) < 32+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	b = b[32:]
	l = int(binary.BigEndian.Uint16(append([]byte{0x00}, b[0])))
	b = b[1:]
	if len(b) < l+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	b = b[l:]
	if len(b) < 2+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	l = int(binary.BigEndian.Uint16(b[:2]))
	b = b[2:]
	if len(b) < l+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	b = b[l:]
	if len(b) < 1+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	l = int(binary.BigEndian.Uint16(append([]byte{0x00}, b[0])))
	b = b[1:]
	if len(b) < l+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	b = b[l:]
	if len(b) < 2+1 {
		return errors.New("not tls protocol or no SNI extension")
	}
	l = int(binary.BigEndian.Uint16(b[:2]))
	b = b[2:]
	if len(b) != l {
		return errors.New("not tls protocol or no SNI extension")
	}
	for {
		if len(b) < 2+1 {
			return errors.New("not tls protocol or no SNI extension")
		}
		t := b[:2]
		b = b[2:]
		if len(b) < 2+1 {
			return errors.New("not tls protocol or no SNI extension")
		}
		l = int(binary.BigEndian.Uint16(b[:2]))
		b = b[2:]
		if len(b) < l {
			return errors.New("not tls protocol or no SNI extension")
		}
		b1 := b[0:l]
		b = b[l:]
		if bytes.Compare(t, []byte{0x00, 0x00}) == 0 {
			if len(b1) < 2+1 {
				return errors.New("not tls protocol or no SNI extension")
			}
			l1 := int(binary.BigEndian.Uint16(b1[:2]))
			b1 = b1[2:]
			if len(b1) != l1 {
				return errors.New("not tls protocol or no SNI extension")
			}
			if len(b1) == 0 {
				return errors.New("not tls protocol or no SNI extension")
			}
			if b1[0] != 0x00 {
				return errors.New("not tls protocol or no SNI extension")
			}
			b1 = b1[1:]
			if len(b1) < 2+1 {
				return errors.New("not tls protocol or no SNI extension")
			}
			l1 = int(binary.BigEndian.Uint16(b1[:2]))
			b1 = b1[2:]
			if len(b1) != l1 {
				return errors.New("not tls protocol or no SNI extension")
			}
			dm = string(b1)
			break
		}
		if len(b) == 0 {
			return errors.New("not tls protocol or no SNI extension")
		}
	}

	debug("hijack", dm, "assume that socket is TCP and dst port is 443")
	rc, err := s.Dial.Dial("tcp", net.JoinHostPort(dm, "443"))
	if err != nil {
		return errors.New("not tls protocol or no SNI extension")
	}
	if _, err := rc.Write(d); err != nil {
		return errors.New("not tls protocol or no SNI extension")
	}
	go func() {
		var bf [1024 * 2]byte
		for {
			if s.TCPDeadline != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
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
		if s.TCPDeadline != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
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

func (s *HijackHTTPS) Shutdown() error {
	return s.RunnerGroup.Done()
}

func (s *HijackHTTPS) Has(host string) bool {
	ss := strings.Split(host, ".")
	var s1 string
	for i := len(ss) - 1; i >= 0; i-- {
		if s1 == "" {
			s1 = ss[i]
		} else {
			s1 = ss[i] + "." + s1
		}
		if _, ok := s.Domains[s1]; ok {
			return true
		}
	}
	return false
}
