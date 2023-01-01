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
	"fmt"
	"log"
	"net"
	"time"

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
)

type EchoServer struct {
	TCPAddr     *net.TCPAddr
	UDPAddr     *net.UDPAddr
	TCPListen   *net.TCPListener
	UDPConn     *net.UDPConn
	RunnerGroup *runnergroup.RunnerGroup
}

func NewEchoServer(addr string) (*EchoServer, error) {
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
	s := &EchoServer{
		TCPAddr:     taddr,
		UDPAddr:     uaddr,
		RunnerGroup: runnergroup.New(),
	}
	return s, nil
}

// Run server.
func (s *EchoServer) ListenAndServe() error {
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
func (s *EchoServer) RunTCPServer() error {
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
			if err := s.TCPHandle(c); err != nil {
				log.Println(err)
			}
		}(c)
	}
	return nil
}

// RunUDPServer starts udp server.
func (s *EchoServer) RunUDPServer() error {
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
func (s *EchoServer) Shutdown() error {
	return s.RunnerGroup.Done()
}

// TCPHandle handles request.
func (s *EchoServer) TCPHandle(c *net.TCPConn) error {
	var b [1024 * 2]byte
	for {
		if err := c.SetDeadline(time.Now().Add(60 * time.Second)); err != nil {
			return err
		}
		i, err := c.Read(b[:])
		if err != nil {
			return nil
		}
		if _, err := c.Write([]byte(c.RemoteAddr().String())); err != nil {
			return err
		}
		if c.RemoteAddr().String() == string(b[:i]) {
			fmt.Printf("TCP: dst:%s <- src:%s\n", c.LocalAddr().String(), c.RemoteAddr().String())
			fmt.Printf("TCP: src:%s -> dst:%s\n", c.LocalAddr().String(), c.RemoteAddr().String())
		}
		if c.RemoteAddr().String() != string(b[:i]) {
			fmt.Printf("TCP: dst:%s <- src:%s <- dst:proxy <- src:%s\n", c.LocalAddr().String(), c.RemoteAddr().String(), string(b[:i]))
			fmt.Printf("TCP: src:%s -> dst:%s -> src:proxy -> dst:%s\n", c.LocalAddr().String(), c.RemoteAddr().String(), string(b[:i]))
		}
	}
	return nil
}

// UDPHandle handles packet.
func (s *EchoServer) UDPHandle(addr *net.UDPAddr, b []byte) error {
	if _, err := s.UDPConn.WriteToUDP([]byte(addr.String()), addr); err != nil {
		return err
	}
	if addr.String() == string(b) {
		fmt.Printf("UDP: dst:%s <- src:%s\n", s.UDPConn.LocalAddr().String(), addr.String())
		fmt.Printf("UDP: src:%s -> dst:%s\n", s.UDPConn.LocalAddr().String(), addr.String())
	}
	if addr.String() != string(b) {
		fmt.Printf("UDP: dst:%s <- src:%s <- dst:proxy <- src:%s\n", s.UDPConn.LocalAddr().String(), addr.String(), string(b))
		fmt.Printf("UDP: src:%s -> dst:%s -> src:proxy -> dst:%s\n", s.UDPConn.LocalAddr().String(), addr.String(), string(b))
	}
	return nil
}
