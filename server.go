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
	"io"
	"log"
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/socks5"
)

// Server.
type Server struct {
	Password      []byte
	TCPAddr       *net.TCPAddr
	UDPAddr       *net.UDPAddr
	TCPListen     *net.TCPListener
	UDPConn       *net.UDPConn
	Cache         *cache.Cache
	TCPDeadline   int
	TCPTimeout    int
	UDPDeadline   int
	ServerAuthman plugin.ServerAuthman
}

// NewServer.
func NewServer(addr, password string, tcpTimeout, tcpDeadline, udpDeadline int) (*Server, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &Server{
		Password:    []byte(password),
		TCPAddr:     taddr,
		UDPAddr:     uaddr,
		Cache:       cs,
		TCPTimeout:  tcpTimeout,
		TCPDeadline: tcpDeadline,
		UDPDeadline: udpDeadline,
	}
	return s, nil
}

// SetServerAuthman sets authman plugin.
func (s *Server) SetServerAuthman(m plugin.ServerAuthman) {
	s.ServerAuthman = m
}

// Run server.
func (s *Server) ListenAndServe() error {
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
func (s *Server) RunTCPServer() error {
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
func (s *Server) RunUDPServer() error {
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

// TCPHandle handles request.
func (s *Server) TCPHandle(c *net.TCPConn) error {
	cn := make([]byte, 12)
	if _, err := io.ReadFull(c, cn); err != nil {
		return err
	}
	ck, err := GetKey(s.Password, cn)
	if err != nil {
		return err
	}
	var b []byte
	b, cn, err = ReadFrom(c, ck, cn, true)
	if err != nil {
		return err
	}
	address := socks5.ToAddress(b[0], b[1:len(b)-2], b[len(b)-2:])
	a := b[0]

	var ai plugin.Internet
	if s.ServerAuthman != nil {
		b, cn, err = ReadFrom(c, ck, cn, false)
		if err != nil {
			return err
		}
		ai, err = s.ServerAuthman.VerifyToken(b, "tcp", a, address, nil)
		if err != nil {
			return err
		}
		defer ai.Close()
	}

	if Debug {
		log.Println("Dial TCP", address)
	}
	tmp, err := Dial.Dial("tcp", address)
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

	go func() {
		k, n, err := PrepareKey(s.Password)
		if err != nil {
			log.Println(err)
			return
		}
		i, err := c.Write(n)
		if err != nil {
			return
		}
		if ai != nil {
			if err := ai.TCPEgress(i); err != nil {
				log.Println(err)
				return
			}
		}
		var b [1024 * 2]byte
		for {
			if s.TCPDeadline != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
					return
				}
			}
			i, err := rc.Read(b[:])
			if err != nil {
				return
			}
			n, i, err = WriteTo(c, b[0:i], k, n, false)
			if err != nil {
				return
			}
			if ai != nil {
				if err := ai.TCPEgress(i); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}()

	for {
		if s.TCPDeadline != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
				return nil
			}
		}
		b, cn, err = ReadFrom(c, ck, cn, false)
		if err != nil {
			return nil
		}
		i, err := rc.Write(b)
		if err != nil {
			return nil
		}
		if ai != nil {
			if err := ai.TCPEgress(i); err != nil {
				return err
			}
		}
	}
	return nil
}

type ServerUDPExchange struct {
	ClientAddr *net.UDPAddr
	RemoteConn *net.UDPConn
	Internet   plugin.Internet
}

// UDPHandle handles packet.
func (s *Server) UDPHandle(addr *net.UDPAddr, b []byte) error {
	a, h, p, data, err := Decrypt(s.Password, b)
	if err != nil {
		return err
	}
	send := func(ue *ServerUDPExchange, data []byte) error {
		if s.ServerAuthman != nil {
			l := int(binary.BigEndian.Uint16(data[len(data)-2:]))
			data = data[0 : len(data)-l-2]
		}
		i, err := ue.RemoteConn.Write(data)
		if err != nil {
			return err
		}
		if ue.Internet != nil {
			if err := ue.Internet.UDPEgress(i); err != nil {
				return err
			}
		}
		return nil
	}

	var ue *ServerUDPExchange
	iue, ok := s.Cache.Get(addr.String())
	if ok {
		ue = iue.(*ServerUDPExchange)
		return send(ue, data)
	}

	address := socks5.ToAddress(a, h, p)
	var ai plugin.Internet
	if s.ServerAuthman != nil {
		l := int(binary.BigEndian.Uint16(data[len(data)-2:]))
		ai, err = s.ServerAuthman.VerifyToken(data[len(data)-l-2:len(data)-2], "udp", a, address, data[0:len(data)-l-2])
		if err != nil {
			return err
		}
	}
	if Debug {
		log.Println("Dial UDP", address)
	}
	c, err := Dial.Dial("udp", address)
	if err != nil {
		return err
	}
	rc := c.(*net.UDPConn)
	ue = &ServerUDPExchange{
		ClientAddr: addr,
		RemoteConn: rc,
		Internet:   ai,
	}
	if err := send(ue, data); err != nil {
		ue.RemoteConn.Close()
		ue.Internet.Close()
		return err
	}
	s.Cache.Set(ue.ClientAddr.String(), ue, cache.DefaultExpiration)
	go func(ue *ServerUDPExchange) {
		defer func() {
			s.Cache.Delete(ue.ClientAddr.String())
			ue.RemoteConn.Close()
			ue.Internet.Close()
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
			a, addr, port, err := socks5.ParseAddress(ue.ClientAddr.String())
			if err != nil {
				log.Println(err)
				break
			}
			d := make([]byte, 0, 7)
			d = append(d, a)
			d = append(d, addr...)
			d = append(d, port...)
			d = append(d, b[0:n]...)
			cd, err := Encrypt(s.Password, d)
			if err != nil {
				log.Println(err)
				break
			}
			i, err := s.UDPConn.WriteToUDP(cd, ue.ClientAddr)
			if err != nil {
				break
			}
			if ue.Internet != nil {
				if err := ue.Internet.UDPEgress(i); err != nil {
					log.Println(err)
					break
				}
			}
		}
	}(ue)
	return nil
}

// Shutdown server.
func (s *Server) Shutdown() error {
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
