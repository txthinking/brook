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
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/tproxy"
	"github.com/txthinking/socks5"
)

// Tproxy.
type Tproxy struct {
	TCPAddr       *net.TCPAddr
	UDPAddr       *net.UDPAddr
	RemoteTCPAddr *net.TCPAddr
	RemoteUDPAddr *net.UDPAddr
	Password      []byte
	TCPListen     *net.TCPListener
	UDPConn       *net.UDPConn
	UDPExchanges  *cache.Cache
	TCPDeadline   int
	TCPTimeout    int
	UDPDeadline   int
}

// NewTproxy.
func NewTproxy(addr, remote, password string, tcpTimeout, tcpDeadline, udpDeadline int) (*Tproxy, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	rtaddr, err := net.ResolveTCPAddr("tcp", remote)
	if err != nil {
		return nil, err
	}
	ruaddr, err := net.ResolveUDPAddr("udp", remote)
	if err != nil {
		return nil, err
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	s := &Tproxy{
		Password:      []byte(password),
		TCPAddr:       taddr,
		UDPAddr:       uaddr,
		RemoteTCPAddr: rtaddr,
		RemoteUDPAddr: ruaddr,
		UDPExchanges:  cs,
		TCPTimeout:    tcpTimeout,
		TCPDeadline:   tcpDeadline,
		UDPDeadline:   udpDeadline,
	}
	return s, nil
}

// Run server.
func (s *Tproxy) ListenAndServe() error {
	errch := make(chan error)
	go func() {
		errch <- s.RunTCPServer()
	}()
	go func() {
		errch <- s.RunUDPServer()
	}()
	return <-errch
}

// RunTCPServer starts tcp server.
func (s *Tproxy) RunTCPServer() error {
	var err error
	s.TCPListen, err = tproxy.ListenTCP("tcp", s.TCPAddr)
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
				if err := c.SetKeepAlivePeriod(time.Duration(s.TCPTimeout) * time.Second); err != nil {
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
			if err := s.TCPHandle(c); err != nil {
				log.Println(err)
			}
		}(c)
	}
	return nil
}

// RunUDPServer starts udp server.
func (s *Tproxy) RunUDPServer() error {
	var err error
	s.UDPConn, err = tproxy.ListenUDP("udp", s.UDPAddr)
	if err != nil {
		return err
	}
	defer s.UDPConn.Close()
	for {
		b := make([]byte, 65536)
		n, saddr, daddr, err := tproxy.ReadFromUDP(s.UDPConn, b)
		if err != nil {
			return err
		}
		if n == 0 {
			continue
		}
		go func(saddr, daddr *net.UDPAddr, b []byte) {
			if err := s.UDPHandle(saddr, daddr, b); err != nil {
				log.Println(err)
				return
			}
		}(saddr, daddr, b[0:n])
	}
	return nil
}

// Shutdown server.
func (s *Tproxy) Shutdown() error {
	var err, err1 error
	if s.TCPListen != nil {
		err = s.TCPListen.Close()
	}
	if s.UDPConn != nil {
		err1 = s.UDPConn.Close()
	}
	if err != nil {
		return err
	}
	return err1
}

// TCPHandle handles request.
func (s *Tproxy) TCPHandle(c *net.TCPConn) error {
	tmp, err := tproxy.DialTCP("tcp", s.RemoteTCPAddr.String())
	if err != nil {
		return err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(s.TCPTimeout) * time.Second); err != nil {
			return err
		}
	}
	if s.TCPDeadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
			return err
		}
	}

	k, n, err := PrepareKey(s.Password)
	if err != nil {
		return err
	}
	if _, err := rc.Write(n); err != nil {
		return err
	}

	a, address, port, err := socks5.ParseAddress(c.LocalAddr().String())
	if err != nil {
		return err
	}
	rawaddr := make([]byte, 0, 7)
	rawaddr = append(rawaddr, a)
	rawaddr = append(rawaddr, address...)
	rawaddr = append(rawaddr, port...)
	n, err = WriteTo(rc, rawaddr, k, n, true)
	if err != nil {
		return err
	}

	go func() {
		n := make([]byte, 12)
		if _, err := io.ReadFull(rc, n); err != nil {
			return
		}
		k, err := GetKey(s.Password, n)
		if err != nil {
			log.Println(err)
			return
		}
		var b []byte
		for {
			if s.TCPDeadline != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
					return
				}
			}
			b, n, err = ReadFrom(rc, k, n, false)
			if err != nil {
				return
			}
			if _, err := c.Write(b); err != nil {
				return
			}
		}
	}()

	var b [1024 * 2]byte
	for {
		if s.TCPDeadline != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
				return nil
			}
		}
		i, err := c.Read(b[:])
		if err != nil {
			return nil
		}
		n, err = WriteTo(rc, b[0:i], k, n, false)
		if err != nil {
			return nil
		}
	}
	return nil
}

type UDPExchange struct {
	RemoteConn *net.UDPConn
	LocalConn  *net.UDPConn
}

func (s *Tproxy) UDPHandle(addr, daddr *net.UDPAddr, b []byte) error {
	a, address, port, err := socks5.ParseAddress(daddr.String())
	if err != nil {
		return err
	}
	rawaddr := make([]byte, 0, 7)
	rawaddr = append(rawaddr, a)
	rawaddr = append(rawaddr, address...)
	rawaddr = append(rawaddr, port...)
	b = append(rawaddr, b...)

	send := func(ue *UDPExchange, data []byte) error {
		cd, err := Encrypt(s.Password, data)
		if err != nil {
			return err
		}
		_, err = ue.RemoteConn.Write(cd)
		if err != nil {
			return err
		}
		return nil
	}

	var ue *UDPExchange
	iue, ok := s.UDPExchanges.Get(addr.String())
	if ok {
		ue = iue.(*UDPExchange)
		return send(ue, b)
	}

	rc, err := tproxy.DialUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}, s.RemoteUDPAddr)
	if err != nil {
		return err
	}
	c, err := tproxy.DialUDP("udp", daddr, addr)
	if err != nil {
		rc.Close()
		return errors.New(fmt.Sprintf("src: %s dst: %s %s", daddr.String(), addr.String(), err.Error()))
	}
	ue = &UDPExchange{
		RemoteConn: rc,
		LocalConn:  c,
	}
	if err := send(ue, b); err != nil {
		ue.RemoteConn.Close()
		ue.LocalConn.Close()
		return err
	}
	s.UDPExchanges.Set(ue.LocalConn.RemoteAddr().String(), ue, cache.DefaultExpiration)
	go func(ue *UDPExchange) {
		defer func() {
			s.UDPExchanges.Delete(ue.LocalConn.RemoteAddr().String())
			ue.RemoteConn.Close()
			ue.LocalConn.Close()
		}()
		var b [65536]byte
		for {
			if s.UDPDeadline != 0 {
				if err := ue.RemoteConn.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
					break
				}
			}
			n, err := ue.RemoteConn.Read(b[:])
			if err != nil {
				break
			}
			_, _, _, data, err := Decrypt(s.Password, b[0:n])
			if err != nil {
				break
			}
			if _, err := ue.LocalConn.Write(data); err != nil {
				break
			}
		}
	}(ue)
	return nil
}
