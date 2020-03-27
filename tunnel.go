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
	"io"
	"log"
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

// Tunnel.
type Tunnel struct {
	TCPAddr       *net.TCPAddr
	UDPAddr       *net.UDPAddr
	ToAddr        string
	RemoteTCPAddr *net.TCPAddr
	RemoteUDPAddr *net.UDPAddr
	Password      []byte
	TCPListen     *net.TCPListener
	UDPConn       *net.UDPConn
	Cache         *cache.Cache
	TCPDeadline   int
	TCPTimeout    int
	UDPDeadline   int
	RunnerGroup   *runnergroup.RunnerGroup
}

// NewTunnel.
func NewTunnel(addr, to, remote, password string, tcpTimeout, tcpDeadline, udpDeadline int) (*Tunnel, error) {
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
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &Tunnel{
		ToAddr:        to,
		Password:      []byte(password),
		TCPAddr:       taddr,
		UDPAddr:       uaddr,
		RemoteTCPAddr: rtaddr,
		RemoteUDPAddr: ruaddr,
		Cache:         cs,
		TCPTimeout:    tcpTimeout,
		TCPDeadline:   tcpDeadline,
		UDPDeadline:   udpDeadline,
		RunnerGroup:   runnergroup.New(),
	}
	return s, nil
}

// Run server.
func (s *Tunnel) ListenAndServe() error {
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
func (s *Tunnel) RunTCPServer() error {
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
func (s *Tunnel) RunUDPServer() error {
	var err error
	s.UDPConn, err = net.ListenUDP("udp", s.UDPAddr)
	if err != nil {
		return err
	}
	defer s.UDPConn.Close()
	for {
		b := make([]byte, 65535)
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
func (s *Tunnel) Shutdown() error {
	return s.RunnerGroup.Done()
}

// TCPHandle handles request.
func (s *Tunnel) TCPHandle(c *net.TCPConn) error {
	tmp, err := Dial.Dial("tcp", s.RemoteTCPAddr.String())
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

	a, address, port, err := socks5.ParseAddress(s.ToAddr)
	if err != nil {
		return err
	}
	ra := make([]byte, 0, 7)
	ra = append(ra, a)
	ra = append(ra, address...)
	ra = append(ra, port...)
	n, _, err = WriteTo(rc, ra, k, n, true)
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
		n, _, err = WriteTo(rc, b[0:i], k, n, false)
		if err != nil {
			return nil
		}
	}
	return nil
}

// UDPHandle handles packet.
func (s *Tunnel) UDPHandle(addr *net.UDPAddr, b []byte) error {
	a, address, port, err := socks5.ParseAddress(s.ToAddr)
	if err != nil {
		return err
	}
	ra := make([]byte, 0, 7)
	ra = append(ra, a)
	ra = append(ra, address...)
	ra = append(ra, port...)
	b = append(ra, b...)

	send := func(ue *socks5.UDPExchange, data []byte) error {
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

	var ue *socks5.UDPExchange
	iue, ok := s.Cache.Get(addr.String())
	if ok {
		ue = iue.(*socks5.UDPExchange)
		return send(ue, b)
	}

	c, err := Dial.Dial("udp", s.RemoteUDPAddr.String())
	if err != nil {
		return err
	}
	rc := c.(*net.UDPConn)
	ue = &socks5.UDPExchange{
		ClientAddr: addr,
		RemoteConn: rc,
	}
	if err := send(ue, b); err != nil {
		ue.RemoteConn.Close()
		return err
	}
	s.Cache.Set(ue.ClientAddr.String(), ue, cache.DefaultExpiration)
	go func(ue *socks5.UDPExchange) {
		defer func() {
			s.Cache.Delete(ue.ClientAddr.String())
			ue.RemoteConn.Close()
		}()
		var b [65535]byte
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
				log.Println(err)
				break
			}
			if _, err := s.UDPConn.WriteToUDP(data, ue.ClientAddr); err != nil {
				break
			}
		}
	}(ue)
	return nil
}
