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
	"log"
	"net"
	"os/exec"
	"strings"
	"unsafe"

	"github.com/txthinking/brook/tproxy"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
	"golang.org/x/sys/unix"
)

type Tproxy struct {
	Addr        string
	blk         *BrookLink
	TCPTimeout  int
	UDPTimeout  int
	RunnerGroup *runnergroup.RunnerGroup
}

var TproxyGate func(conn net.Conn) (done bool, err error) = func(conn net.Conn) (done bool, err error) {
	return false, nil
}

func NewTproxy(addr, link string, tcpTimeout, udpTimeout int) (*Tproxy, error) {
	r, err := NewBrookLink(link)
	if err != nil {
		return nil, err
	}
	// convert to ip first, because gateway system may no dns server
	h, p, err := net.SplitHostPort(r.address)
	if err != nil {
		return nil, err
	}
	ip, _ := Resolve6(h)
	if ip == "" {
		ip, _ = Resolve4(h)
	}
	if ip == "" {
		return nil, errors.New("Can not find server IP")
	}
	r.address = net.JoinHostPort(ip, p)

	t := &Tproxy{
		Addr:        addr,
		TCPTimeout:  tcpTimeout,
		UDPTimeout:  udpTimeout,
		RunnerGroup: runnergroup.New(),
		blk:         r,
	}
	return t, nil
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
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-ip 127.0.0.1 --on-port "+s.Addr[1:])
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "iptables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-ip 127.0.0.1 --on-port "+s.Addr[1:])
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}

	// 6
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
	c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-ip ::1 --on-port "+s.Addr[1:])
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("/bin/sh", "-c", "ip6tables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-ip ::1 --on-port "+s.Addr[1:])
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
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

func (s *Tproxy) ListenAndServe() error {
	addr, err := Resolve("tcp", "127.0.0.1"+s.Addr)
	if err != nil {
		return err
	}
	l, err := tproxy.ListenTCP("tcp", addr.(*net.TCPAddr))
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
						log.Println(err)
					}
				}(c)
			}
			return nil
		},
		Stop: func() error {
			return l.Close()
		},
	})
	addr, err = Resolve("tcp", "[::1]"+s.Addr)
	if err != nil {
		l.Close()
		return err
	}
	l1, err := tproxy.ListenTCP("tcp", addr.(*net.TCPAddr))
	if err != nil {
		l.Close()
		return err
	}
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			for {
				c, err := l1.AcceptTCP()
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
		},
		Stop: func() error {
			return l1.Close()
		},
	})
	addr, err = Resolve("udp", "127.0.0.1"+s.Addr)
	if err != nil {
		l.Close()
		l1.Close()
		return err
	}
	l2, err := tproxy.ListenUDP("udp", addr.(*net.UDPAddr))
	if err != nil {
		l.Close()
		l1.Close()
		return err
	}
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			var hdr unix.Cmsghdr
			var addr unix.RawSockaddrInet6
			for {
				oob := make([]byte, int(unsafe.Sizeof(hdr)+unsafe.Sizeof(addr)))
				b := make([]byte, 65507)
				n, src, dst, err := tproxy.ReadFromUDP(l2, oob, b)
				if err != nil {
					return err
				}
				if n == 0 {
					continue
				}
				c, err := tproxy.DialUDP("udp", dst, src)
				if err != nil {
					log.Println(err)
					continue
				}
				go func(c *net.UDPConn, b []byte) {
					defer c.Close()
					if err := s.UDPHandle(c, b); err != nil {
						log.Println(err)
						return
					}
				}(c, b[0:n])
			}
			return nil
		},
		Stop: func() error {
			return l2.Close()
		},
	})
	addr, err = Resolve("udp", "[::1]"+s.Addr)
	if err != nil {
		l.Close()
		l1.Close()
		l2.Close()
		return err
	}
	l3, err := tproxy.ListenUDP("udp", addr.(*net.UDPAddr))
	if err != nil {
		l.Close()
		l1.Close()
		l2.Close()
		return err
	}
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			var hdr unix.Cmsghdr
			var addr unix.RawSockaddrInet6
			for {
				oob := make([]byte, int(unsafe.Sizeof(hdr)+unsafe.Sizeof(addr)))
				b := make([]byte, 65507)
				n, dst, src, err := tproxy.ReadFromUDP(l3, oob, b)
				if err != nil {
					return err
				}
				if n == 0 {
					continue
				}
				c, err := tproxy.DialUDP("udp", dst, src)
				if err != nil {
					log.Println(err)
					continue
				}
				go func(c *net.UDPConn, b []byte) {
					defer c.Close()
					if err := s.UDPHandle(c, b); err != nil {
						log.Println(err)
						return
					}
				}(c, b[0:n])
			}
			return nil
		},
		Stop: func() error {
			return l3.Close()
		},
	})
	return s.RunnerGroup.Wait()
}

func (s *Tproxy) TCPHandle(c *net.TCPConn) error {
	done, err := TproxyGate(c)
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	a, h, p, err := socks5.ParseAddress(c.LocalAddr().String())
	if err != nil {
		return err
	}
	dstb := append(append([]byte{a}, h...), p...)
	sc, err := s.blk.CreateExchanger("tcp", c.RemoteAddr().String(), dstb, s.TCPTimeout, s.UDPTimeout)
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

func (s *Tproxy) UDPHandle(c *net.UDPConn, b []byte) error {
	c1 := &ConnFirst{UDPConn: c, First: b}
	done, err := TproxyGate(c1)
	if err != nil {
		return err
	}
	if done {
		return nil
	}
	a, h, p, err := socks5.ParseAddress(c.LocalAddr().String())
	if err != nil {
		return err
	}
	dstb := append(append([]byte{a}, h...), p...)
	sc, err := s.blk.CreateExchanger("udp", c.RemoteAddr().String(), dstb, s.TCPTimeout, s.UDPTimeout)
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
	if err := sc.Exchange(c1); err != nil {
		return nil
	}
	return nil
}

func (s *Tproxy) Shutdown() error {
	return s.RunnerGroup.Done()
}
