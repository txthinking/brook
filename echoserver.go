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
	"net"
	"time"

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
)

type EchoServer struct {
	Addr        string
	RunnerGroup *runnergroup.RunnerGroup
}

func NewEchoServer(addr string) (*EchoServer, error) {
	if err := limits.Raise(); err != nil {
		Log(Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	s := &EchoServer{
		Addr:        addr,
		RunnerGroup: runnergroup.New(),
	}
	return s, nil
}

func (s *EchoServer) ListenAndServe() error {
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
					if err := s.TCPHandle(c); err != nil {
						Log(Error{"from": c.RemoteAddr().String(), "error": err.Error()})
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
			for {
				b := make([]byte, 65507)
				n, addr, err := l1.ReadFromUDP(b)
				if err != nil {
					return err
				}
				go func(addr *net.UDPAddr, b []byte) {
					if err := s.UDPHandle(addr, b, l1); err != nil {
						Log(Error{"from": addr.String(), "error": err.Error()})
						return
					}
				}(addr, b[0:n])
			}
			return nil
		},
		Stop: func() error {
			return l1.Close()
		},
	})
	return s.RunnerGroup.Wait()
}

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

func (s *EchoServer) UDPHandle(addr *net.UDPAddr, b []byte, l1 *net.UDPConn) error {
	if _, err := l1.WriteToUDP([]byte(addr.String()), addr); err != nil {
		return err
	}
	if addr.String() == string(b) {
		fmt.Printf("UDP: dst:%s <- src:%s\n", l1.LocalAddr().String(), addr.String())
		fmt.Printf("UDP: src:%s -> dst:%s\n", l1.LocalAddr().String(), addr.String())
	}
	if addr.String() != string(b) {
		fmt.Printf("UDP: dst:%s <- src:%s <- dst:proxy <- src:%s\n", l1.LocalAddr().String(), addr.String(), string(b))
		fmt.Printf("UDP: src:%s -> dst:%s -> src:proxy -> dst:%s\n", l1.LocalAddr().String(), addr.String(), string(b))
	}
	return nil
}

func (s *EchoServer) Shutdown() error {
	return s.RunnerGroup.Done()
}
