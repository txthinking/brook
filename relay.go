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
	"time"

	"github.com/miekg/dns"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

type Relay struct {
	From        string
	To          string
	Dstb        []byte
	TCPTimeout  int
	UDPTimeout  int
	Pcf         *PacketConnFactory
	RunnerGroup *runnergroup.RunnerGroup
	IsDNS       bool
}

func NewRelay(from, to string, tcpTimeout, udpTimeout int) (*Relay, error) {
	if err := limits.Raise(); err != nil {
		Log(Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	a, h, p, err := socks5.ParseAddress(to)
	if err != nil {
		return nil, err
	}
	s := &Relay{
		From:        from,
		To:          to,
		Dstb:        append(append([]byte{a}, h...), p...),
		TCPTimeout:  tcpTimeout,
		UDPTimeout:  udpTimeout,
		Pcf:         NewPacketConnFactory(),
		RunnerGroup: runnergroup.New(),
	}
	return s, nil
}

func (s *Relay) ListenAndServe() error {
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

func (s *Relay) TCPHandle(c *net.TCPConn) error {
	rc, err := DialTCP("tcp", "", s.To)
	if err != nil {
		return err
	}
	defer rc.Close()
	go func() {
		var bf [1024 * 2]byte
		for {
			if s.TCPTimeout != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
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
		if s.TCPTimeout != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
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

func (s *Relay) UDPHandle(addr *net.UDPAddr, b []byte, l1 *net.UDPConn) error {
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
	c, err := s.Pcf.Handle(addr, s.Dstb, b, func(b []byte) (int, error) {
		return l1.WriteToUDP(b, addr)
	}, s.UDPTimeout)
	if err != nil {
		return err
	}
	if c == nil {
		return nil
	}
	defer c.Close()
	rc, err := NATDial("udp", addr.String(), s.To, s.To)
	if err != nil {
		return err
	}
	defer rc.Close()
	go func() {
		var bf [65507]byte
		for {
			if s.UDPTimeout != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.UDPTimeout) * time.Second)); err != nil {
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
	var bf [65507]byte
	for {
		if s.UDPTimeout != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.UDPTimeout) * time.Second)); err != nil {
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

func (s *Relay) Shutdown() error {
	return s.RunnerGroup.Done()
}
