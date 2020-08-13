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
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

// Server.
type Server struct {
	Password     []byte
	TCPAddr      *net.TCPAddr
	UDPAddr      *net.UDPAddr
	TCPListen    *net.TCPListener
	UDPConn      *net.UDPConn
	UDPExchanges *cache.Cache
	TCPTimeout   int
	UDPTimeout   int
	RunnerGroup  *runnergroup.RunnerGroup
	UDPSrc       *cache.Cache
}

// NewServer.
func NewServer(addr, password string, tcpTimeout, udpTimeout int) (*Server, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs2 := cache.New(cache.NoExpiration, cache.NoExpiration)
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &Server{
		Password:     []byte(password),
		TCPAddr:      taddr,
		UDPAddr:      uaddr,
		UDPExchanges: cs,
		TCPTimeout:   tcpTimeout,
		UDPTimeout:   udpTimeout,
		RunnerGroup:  runnergroup.New(),
		UDPSrc:       cs2,
	}
	return s, nil
}

// Run server.
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
func (s *Server) RunUDPServer() error {
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

// TCPHandle handles request.
func (s *Server) TCPHandle(c *net.TCPConn) error {
	ss, dst, err := NewStreamServer(s.Password, c, s.TCPTimeout)
	if err != nil {
		return err
	}
	defer ss.Clean()
	address := socks5.ToAddress(dst[0], dst[1:len(dst)-2], dst[len(dst)-2:])
	debug("dial tcp", address)
	rc, err := Dial.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

// UDPHandle handles packet.
func (s *Server) UDPHandle(addr *net.UDPAddr, b []byte) error {
	src := addr.String()
	dstb, d, w, err := PacketClientToRemote(s.Password, b)
	if err != nil {
		return err
	}
	dst := socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:])
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		if _, err := ue.Any.(io.Writer).Write(d); err != nil {
			return err
		}
		return nil
	}

	debug("dial udp", dst)
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	raddr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return err
	}
	rc, err := Dial.DialUDP("udp", laddr, raddr)
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
	wer := w(rc.Write)
	if _, err := wer.Write(d); err != nil {
		return err
	}
	ue := &UDPExchange{
		Any: wer,
	}
	s.UDPExchanges.Set(src+dst, ue, -1)
	defer s.UDPExchanges.Delete(src + dst)
	ps := NewPacketServer(s.Password)
	defer ps.Clean()
	err = ps.RemoteToClient(rc, s.UDPTimeout, dstb, w(func(b []byte) (int, error) {
		return s.UDPConn.WriteToUDP(b, addr)
	}))
	if err != nil {
		return err
	}
	return nil
}

// Shutdown server.
func (s *Server) Shutdown() error {
	return s.RunnerGroup.Done()
}
