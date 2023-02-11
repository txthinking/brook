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
	"net"

	"github.com/miekg/dns"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

type RelayOverBrook struct {
	From        string
	Link        string
	dstb        []byte
	TCPTimeout  int
	UDPTimeout  int
	blk         *BrookLink
	pcf         *PacketConnFactory
	RunnerGroup *runnergroup.RunnerGroup
	IsDNS       bool
}

func NewRelayOverBrook(from, link, to string, tcpTimeout, udpTimeout int) (*RelayOverBrook, error) {
	a, h, p, err := socks5.ParseAddress(to)
	if err != nil {
		return nil, err
	}
	blk, err := NewBrookLink(link)
	if err != nil {
		return nil, err
	}
	s := &RelayOverBrook{
		From:        from,
		Link:        link,
		dstb:        append(append([]byte{a}, h...), p...),
		TCPTimeout:  tcpTimeout,
		UDPTimeout:  udpTimeout,
		blk:         blk,
		pcf:         NewPacketConnFactory(),
		RunnerGroup: runnergroup.New(),
	}
	return s, nil
}

func (s *RelayOverBrook) ListenAndServe() error {
	addr, err := net.ResolveTCPAddr("tcp", s.From)
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
						Log(&Error{"from": c.RemoteAddr().String(), "error": err.Error()})
					}
				}(c)
			}
			return nil
		},
		Stop: func() error {
			return l.Close()
		},
	})
	addr1, err := net.ResolveUDPAddr("udp", s.From)
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
						Log(&Error{"from": addr.String(), "error": err.Error()})
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

func (s *RelayOverBrook) TCPHandle(c *net.TCPConn) error {
	sc, err := s.blk.CreateExchanger("tcp", c.RemoteAddr().String(), s.dstb, s.TCPTimeout, s.UDPTimeout)
	if err != nil {
		return err
	}
	defer sc.Clean()
	if v, ok := sc.(*StreamClient); ok {
		defer v.Server.Close()
	}
	if v, ok := sc.(*SimpleStreamClient); ok {
		defer v.Server.Close()
	}
	if err := sc.Exchange(c); err != nil {
		return nil
	}
	return nil
}

func (s *RelayOverBrook) UDPHandle(addr *net.UDPAddr, b []byte, l1 *net.UDPConn) error {
	if s.IsDNS {
		m := &dns.Msg{}
		if err := m.Unpack(b); err != nil {
			return err
		}
		if len(m.Question) == 0 {
			return errors.New("no question")
		}
		done, err := DNSGate(addr, m, l1)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
	conn, err := s.pcf.Handle(addr, s.dstb, b, func(b []byte) (int, error) {
		return l1.WriteToUDP(b, addr)
	}, s.UDPTimeout)
	if err != nil {
		return err
	}
	if conn == nil {
		return nil
	}
	defer conn.Close()
	sc, err := s.blk.CreateExchanger("udp", addr.String(), s.dstb, s.TCPTimeout, s.UDPTimeout)
	if err != nil {
		return err
	}
	defer sc.Clean()
	if v, ok := sc.(*PacketClient); ok {
		defer v.Server.Close()
	}
	if v, ok := sc.(*StreamClient); ok {
		defer v.Server.Close()
	}
	if v, ok := sc.(*SimplePacketClient); ok {
		defer v.Server.Close()
	}
	if v, ok := sc.(*SimpleStreamClient); ok {
		defer v.Server.Close()
	}
	if err := sc.Exchange(conn); err != nil {
		return nil
	}
	return nil
}

func (s *RelayOverBrook) Shutdown() error {
	return s.RunnerGroup.Done()
}
