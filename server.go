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
	"log"
	"net"
	"time"

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
)

type Server struct {
	Password    []byte
	TCPAddr     *net.TCPAddr
	UDPAddr     *net.UDPAddr
	TCPListen   *net.TCPListener
	UDPConn     *net.UDPConn
	TCPTimeout  int
	UDPTimeout  int
	RunnerGroup *runnergroup.RunnerGroup
}

func NewServer(addr, password string, tcpTimeout, udpTimeout int) (*Server, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &Server{
		Password:    []byte(password),
		TCPAddr:     taddr,
		UDPAddr:     uaddr,
		TCPTimeout:  tcpTimeout,
		UDPTimeout:  udpTimeout,
		RunnerGroup: runnergroup.New(),
	}
	return s, nil
}

func (s *Server) ListenAndServe() error {
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
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
					log.Println(err)
					return
				}
			}
			ss, err := NewStreamServer(s.Password, c.RemoteAddr().String(), c, s.TCPTimeout, s.UDPTimeout)
			if err != nil {
				log.Println(err)
				return
			}
			defer ss.Clean()
			if ss.Network() == "tcp" {
				if err := s.TCPHandle(ss); err != nil {
					log.Println(err)
				}
			}
			if ss.Network() == "udp" {
				if err := s.UDPOverTCPHandle(ss); err != nil {
					log.Println(err)
				}
			}
		}(c)
	}
	return nil
}

func (s *Server) RunUDPServer() error {
	var err error
	s.UDPConn, err = net.ListenUDP("udp", s.UDPAddr)
	if err != nil {
		return err
	}
	defer s.UDPConn.Close()
	f := NewPacketServerConnFactory()
	for {
		b := make([]byte, 65507)
		n, addr, err := s.UDPConn.ReadFromUDP(b)
		if err != nil {
			return err
		}
		conn, dstb, err := f.Handle(addr, b[0:n], s.Password, func(b []byte) (int, error) {
			return s.UDPConn.WriteToUDP(b, addr)
		}, s.UDPTimeout)
		if err != nil {
			log.Println(err)
			continue
		}
		if conn == nil {
			continue
		}
		go func() {
			defer conn.Close()
			ss, err := NewPacketServer(s.Password, addr.String(), conn, s.UDPTimeout, dstb)
			if err != nil {
				log.Println(err)
				return
			}
			defer ss.Clean()
			if err := s.UDPHandle(ss); err != nil {
				log.Println(err)
			}
		}()
	}
	return nil
}

func (s *Server) TCPHandle(ss Exchanger) error {
	rc, err := DialTCP("tcp", "", ss.Dst())
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

func (s *Server) UDPOverTCPHandle(ss Exchanger) error {
	rc, err := NATDial("udp", ss.Src(), ss.Dst(), ss.Dst())
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

func (s *Server) UDPHandle(ss Exchanger) error {
	rc, err := NATDial("udp", ss.Src(), ss.Dst(), ss.Dst())
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

func (s *Server) Shutdown() error {
	return s.RunnerGroup.Done()
}
