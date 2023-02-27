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

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

type Server struct {
	Addr        string
	Password    []byte
	TCPTimeout  int
	UDPTimeout  int
	RunnerGroup *runnergroup.RunnerGroup
}

func NewServer(addr, password string, tcpTimeout, udpTimeout int) (*Server, error) {
	if err := limits.Raise(); err != nil {
		Log(Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	s := &Server{
		Password:    []byte(password),
		Addr:        addr,
		TCPTimeout:  tcpTimeout,
		UDPTimeout:  udpTimeout,
		RunnerGroup: runnergroup.New(),
	}
	return s, nil
}

func (s *Server) ListenAndServe() error {
	addr, err := net.ResolveTCPAddr("tcp", s.Addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			for {
				c, err := l.AcceptTCP()
				if err != nil {
					return err
				}
				go func(c *net.TCPConn) {
					defer c.Close()
					ss, err := NewStreamServer(s.Password, c.RemoteAddr().String(), c, s.TCPTimeout, s.UDPTimeout)
					if err != nil {
						Log(Error{"from": c.RemoteAddr().String(), "error": err.Error()})
						return
					}
					defer ss.Clean()
					if ss.Network() == "tcp" {
						if err := s.TCPHandle(ss); err != nil {
							Log(Error{"from": c.RemoteAddr().String(), "dst": ss.Dst(), "error": err.Error()})
						}
					}
					if ss.Network() == "udp" {
						if err := s.UDPOverTCPHandle(ss); err != nil {
							Log(Error{"from": c.RemoteAddr().String(), "dst": ss.Dst(), "error": err.Error()})
						}
					}
				}(c)
			}
			return nil
		},
		Stop: func() error {
			return l.Close()
		},
	})
	addr1, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
		l.Close()
		return err
	}
	l1, err := net.ListenUDP("udp", addr1)
	if err != nil {
		l.Close()
		return err
	}
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			f := NewPacketServerConnFactory()
			for {
				b := make([]byte, 65507)
				n, addr, err := l1.ReadFromUDP(b)
				if err != nil {
					return err
				}
				conn, dstb, err := f.Handle(addr, b[0:n], s.Password, func(b []byte) (int, error) {
					return l1.WriteToUDP(b, addr)
				}, s.UDPTimeout)
				if err != nil {
					Log(Error{"from": addr.String(), "error": err.Error()})
					continue
				}
				if conn == nil {
					continue
				}
				go func() {
					defer conn.Close()
					ss, err := NewPacketServer(s.Password, addr.String(), conn, s.UDPTimeout, dstb)
					if err != nil {
						Log(Error{"from": addr.String(), "dst": socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), "error": err.Error()})
						return
					}
					defer ss.Clean()
					if err := s.UDPHandle(ss); err != nil {
						Log(Error{"from": addr.String(), "dst": socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), "error": err.Error()})
					}
				}()
			}
			return nil
		},
		Stop: func() error {
			return l1.Close()
		},
	})
	return s.RunnerGroup.Wait()
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
