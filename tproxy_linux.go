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
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/brook/tproxy"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
)

// Tproxy.
type Tproxy struct {
	TCPAddr       *net.TCPAddr
	UDPAddr       *net.UDPAddr
	ServerTCPAddr *net.TCPAddr
	ServerUDPAddr *net.UDPAddr
	Password      []byte
	TCPListen     *net.TCPListener
	UDPConn       *net.UDPConn
	UDPExchanges  *cache.Cache
	UDPSrc        *cache.Cache
	TCPTimeout    int
	UDPTimeout    int
	RunnerGroup   *runnergroup.RunnerGroup
	EnableIPv6    bool
	Cidr4         []*net.IPNet
	Cidr6         []*net.IPNet
	BypassCache   *cache.Cache
}

// NewTproxy.
func NewTproxy(addr, server, password string, enableIPv6 bool, cidr4url, cidr6url string, tcpTimeout, udpTimeout int) (*Tproxy, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	rtaddr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		return nil, err
	}
	ruaddr, err := net.ResolveUDPAddr("udp", server)
	if err != nil {
		return nil, err
	}
	c4 := make([]*net.IPNet, 0)
	if cidr4url != "" {
		l, err := ReadList(cidr4url)
		if err != nil {
			return nil, err
		}
		for _, v := range l {
			_, in, err := net.ParseCIDR(v)
			if err != nil {
				continue
			}
			c4 = append(c4, in)
		}
	}
	c6 := make([]*net.IPNet, 0)
	if cidr6url != "" {
		l, err := ReadList(cidr6url)
		if err != nil {
			return nil, err
		}
		for _, v := range l {
			_, in, err := net.ParseCIDR(v)
			if err != nil {
				continue
			}
			c6 = append(c6, in)
		}
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs2 := cache.New(cache.NoExpiration, cache.NoExpiration)
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &Tproxy{
		Password:      []byte(password),
		TCPAddr:       taddr,
		UDPAddr:       uaddr,
		ServerTCPAddr: rtaddr,
		ServerUDPAddr: ruaddr,
		UDPExchanges:  cs,
		UDPSrc:        cs2,
		TCPTimeout:    tcpTimeout,
		UDPTimeout:    udpTimeout,
		RunnerGroup:   runnergroup.New(),
		EnableIPv6:    enableIPv6,
		Cidr4:         c4,
		Cidr6:         c6,
		BypassCache:   cache.New(cache.NoExpiration, cache.NoExpiration),
	}
	return s, nil
}

func (s *Tproxy) RunAutoScripts() error {
	c := exec.Command("/bin/sh", "-c", "echo 1 > /proc/sys/net/ipv4/ip_forward")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}

	c = exec.Command("/bin/sh", "-c", "ip route add local 0.0.0.0/0 dev lo table 100")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "ip rule add fwmark 1 lookup 100")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}

	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -d 0.0.0.0/8 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -d 10.0.0.0/8 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -d 127.0.0.0/8 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -d 169.254.0.0/16 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -d 172.16.0.0/12 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -d 192.168.0.0/16 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -d 224.0.0.0/4 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -d 240.0.0.0/4 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}

	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -p tcp -m socket -j MARK --set-mark 1")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port "+strconv.Itoa(s.TCPAddr.Port))
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port "+strconv.Itoa(s.UDPAddr.Port))
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}

	if s.EnableIPv6 {
		c = exec.Command("/bin/sh", "-c", "echo 1 > /proc/sys/net/ipv6/conf/all/forwarding")
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}

		c = exec.Command("/bin/sh", "-c", "ip -6 route add local ::/0 dev lo table 106")
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
		c = exec.Command("/bin/sh", "-c", "ip -6 rule add fwmark 1 table 106")
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}

		c = exec.Command("/bin/sh", "-c", "ip address | grep -w inet6 | awk '{print $2}'")
		out, err := c.CombinedOutput()
		if err != nil {
			return errors.New(string(out) + err.Error())
		}
		l := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, v := range l {
			c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -A PREROUTING -d "+v+" -j RETURN")
			if out, err := c.CombinedOutput(); err != nil {
				return errors.New(string(out) + err.Error())
			}
		}

		c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -A PREROUTING -p tcp -m socket -j MARK --set-mark 1")
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
		c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port "+strconv.Itoa(s.TCPAddr.Port))
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
		c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1")
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
		c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port "+strconv.Itoa(s.UDPAddr.Port))
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
	}
	return nil
}

func (s *Tproxy) ClearAutoScripts() error {
	c := exec.Command("/bin/sh", "-c", "iptables -t mangle -F")
	c.Run()
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -X")
	c.Run()
	c = exec.Command("/bin/sh", "-c", "ip rule del fwmark 1 lookup 100")
	c.Run()
	c = exec.Command("/bin/sh", "-c", "ip route del local 0.0.0.0/0 dev lo table 100")
	c.Run()
	c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -F")
	c.Run()
	c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -X")
	c.Run()
	c = exec.Command("/bin/sh", "-c", "ip -6 rule del fwmark 1 table 106")
	c.Run()
	c = exec.Command("/bin/sh", "-c", "ip -6 route del local ::/0 dev lo table 106")
	c.Run()
	return nil
}

// Run server.
func (s *Tproxy) ListenAndServe() error {
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
func (s *Tproxy) RunTCPServer() error {
	var err error
	s.TCPListen, err = tproxy.ListenTCP("tcp", s.TCPAddr)
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
func (s *Tproxy) RunUDPServer() error {
	var err error
	s.UDPConn, err = tproxy.ListenUDP("udp", s.UDPAddr)
	if err != nil {
		return err
	}
	defer s.UDPConn.Close()
	for {
		b := make([]byte, 65507)
		n, saddr, daddr, err := tproxy.ReadFromUDP(s.UDPConn, b)
		if err != nil {
			return err
		}
		if n == 0 {
			continue
		}
		go func(saddr, daddr *net.UDPAddr, b []byte) {
			if err := s.UDPHandle(saddr, daddr, b); err != nil {
				log.Println(err)
				return
			}
		}(saddr, daddr, b[0:n])
	}
	return nil
}

// Shutdown server.
func (s *Tproxy) Shutdown() error {
	return s.RunnerGroup.Done()
}

func (s *Tproxy) HasIP(i net.IP) bool {
	if i == nil {
		return false
	}
	any, ok := s.BypassCache.Get(i.String())
	if ok {
		return any.(bool)
	}
	if i.To4() != nil {
		ii := i.To4()
		for _, v := range s.Cidr4 {
			if v.Contains(ii) {
				s.BypassCache.Set(i.String(), true, -1)
				return true
			}
		}
	}
	if i.To4() == nil {
		ii := i.To16()
		for _, v := range s.Cidr6 {
			if v.Contains(ii) {
				s.BypassCache.Set(i.String(), true, -1)
				return true
			}
		}
	}
	s.BypassCache.Set(i.String(), false, -1)
	return false
}

// TCPHandle handles request.
func (s *Tproxy) TCPHandle(c *net.TCPConn) error {
	if s.HasIP(c.LocalAddr().(*net.TCPAddr).IP) {
		rc, err := Dial.Dial("tcp", c.LocalAddr().String())
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

	rc, err := tproxy.DialTCP("tcp", s.ServerTCPAddr.String())
	if err != nil {
		return err
	}
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}

	a, h, p, err := socks5.ParseAddress(c.LocalAddr().String())
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

func (s *Tproxy) UDPHandle(addr, daddr *net.UDPAddr, b []byte) error {
	src := addr.String()
	dst := daddr.String()
	var laddr *net.UDPAddr
	any, ok := s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	if laddr == nil {
		if addr.IP.To4() != nil {
			laddr = &net.UDPAddr{
				IP:   net.IPv4zero,
				Port: 0,
				Zone: addr.Zone,
			}
		}
		if addr.IP.To4() == nil {
			laddr = &net.UDPAddr{
				IP:   net.IPv6zero,
				Port: 0,
				Zone: addr.Zone,
			}
		}
	}

	if s.HasIP(daddr.IP) {
		any, ok := s.UDPExchanges.Get(src + dst)
		if ok {
			ue := any.(*UDPExchange)
			if _, err := ue.Conn.Write(b); err != nil {
				return err
			}
			return nil
		}
		rc, err := Dial.DialUDP("udp", laddr, daddr)
		if err != nil {
			if strings.Contains(err.Error(), "address already in use") {
				// we dont choose lock, so ignore this error
				return nil
			}
			return err
		}
		defer rc.Close()
		if laddr.Port == 0 {
			s.UDPSrc.Set(src+dst, rc.LocalAddr().(*net.UDPAddr), -1)
		}
		c, err := tproxy.DialUDP("udp", daddr, addr)
		if err != nil {
			return errors.New(fmt.Sprintf("src: %s dst: %s %s", daddr.String(), addr.String(), err.Error()))
		}
		defer c.Close()
		if _, err := rc.Write(b); err != nil {
			return err
		}
		ue := &UDPExchange{
			Conn: rc,
		}
		s.UDPExchanges.Set(src+dst, ue, -1)
		go func(rc *net.UDPConn) {
			var b [65507]byte
			for {
				if s.UDPTimeout != 0 {
					if err := rc.SetDeadline(time.Now().Add(time.Duration(s.UDPTimeout) * time.Second)); err != nil {
						log.Println(err)
						return
					}
				}
				n, err := rc.Read(b[:])
				if err != nil {
					return
				}
				if _, err := c.Write(b[:n]); err != nil {
					return
				}
			}
		}(rc)
		return nil
	}

	any, ok = s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		return ue.Any.(*PacketClient).LocalToServer(ue.Dst, b, ue.Conn, s.UDPTimeout)
	}
	rc, err := tproxy.DialUDP("udp", laddr, s.ServerUDPAddr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			// we dont choose lock, so ignore this error
			return nil
		}
		return err
	}
	defer rc.Close()
	if laddr.Port == 0 {
		s.UDPSrc.Set(src+dst, rc.LocalAddr().(*net.UDPAddr), -1)
	}
	c, err := tproxy.DialUDP("udp", daddr, addr)
	if err != nil {
		return errors.New(fmt.Sprintf("src: %s dst: %s %s", daddr.String(), addr.String(), err.Error()))
	}
	defer c.Close()

	a, h, p, err := socks5.ParseAddress(dst)
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
		return c.Write(d)
	})
	if err != nil {
		return nil
	}
	return nil
}
