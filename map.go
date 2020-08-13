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
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

// Map.
type Map struct {
	TCPAddr       *net.TCPAddr
	UDPAddr       *net.UDPAddr
	RemoteAddress string
	ServerTCPAddr *net.TCPAddr
	ServerUDPAddr *net.UDPAddr
	Password      []byte
	TCPListen     *net.TCPListener
	UDPConn       *net.UDPConn
	UDPExchanges  *cache.Cache
	TCPTimeout    int
	UDPTimeout    int
	RunnerGroup   *runnergroup.RunnerGroup
	UDPSrc        *cache.Cache
}

// NewMap.
func NewMap(from, to, remote, password string, tcpTimeout, udpTimeout int) (*Map, error) {
	taddr, err := net.ResolveTCPAddr("tcp", from)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", from)
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
	cs2 := cache.New(cache.NoExpiration, cache.NoExpiration)
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &Map{
		RemoteAddress: to,
		Password:      []byte(password),
		TCPAddr:       taddr,
		UDPAddr:       uaddr,
		ServerTCPAddr: rtaddr,
		ServerUDPAddr: ruaddr,
		UDPExchanges:  cs,
		TCPTimeout:    tcpTimeout,
		UDPTimeout:    udpTimeout,
		RunnerGroup:   runnergroup.New(),
		UDPSrc:        cs2,
	}
	return s, nil
}

// Run server.
func (s *Map) ListenAndServe() error {
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
func (s *Map) RunTCPServer() error {
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
func (s *Map) RunUDPServer() error {
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
func (s *Map) Shutdown() error {
	return s.RunnerGroup.Done()
}

// TCPHandle handles request.
func (s *Map) TCPHandle(c *net.TCPConn) error {
	rc, err := Dial.Dial("tcp", s.ServerTCPAddr.String())
	if err != nil {
		return err
	}
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	a, h, p, err := socks5.ParseAddress(s.RemoteAddress)
	if err != nil {
		return err
	}
	dst := make([]byte, 0, 1+len(h)+2)
	dst = append(dst, a)
	dst = append(dst, h...)
	dst = append(dst, p...)
	sc, err := NewStreamClient("tcp", s.Password, dst, rc, s.TCPTimeout)
	if err != nil {
		return err
	}
	defer sc.Clean()
	if err := sc.Exchange(c); err != nil {
		return nil
	}
	return nil
}

// UDPHandle handles packet.
func (s *Map) UDPHandle(addr *net.UDPAddr, b []byte) error {
	src := addr.String()
	dst := s.RemoteAddress
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		return ue.Any.(*PacketClient).LocalToServer(ue.Dst, b, ue.Conn, s.UDPTimeout)
	}
	debug("dial udp", dst)
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	rc, err := Dial.DialUDP("udp", laddr, s.ServerUDPAddr)
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
	a, h, p, err := socks5.ParseAddress(s.RemoteAddress)
	if err != nil {
		return err
	}
	dstb := make([]byte, 0, 1+len(h)+2)
	dstb = append(dstb, a)
	dstb = append(dstb, h...)
	dstb = append(dstb, p...)
	pc := NewPacketClient(s.Password)
	defer pc.Clean()
	if err := pc.LocalToServer(dstb, b, rc, s.UDPTimeout); err != nil {
		return err
	}
	ue := &UDPExchange{
		Conn: rc,
		Any:  pc,
		Dst:  dstb,
	}
	s.UDPExchanges.Set(src+dst, ue, -1)
	defer s.UDPExchanges.Delete(src + dst)
	err = pc.RunServerToLocal(rc, s.UDPTimeout, func(dst, d []byte) (int, error) {
		return s.UDPConn.WriteToUDP(d, addr)
	})
	if err != nil {
		return nil
	}
	return nil
}
