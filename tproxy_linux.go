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
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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
	RemoteTCPAddr *net.TCPAddr
	RemoteUDPAddr *net.UDPAddr
	Password      []byte
	TCPListen     *net.TCPListener
	UDPConn       *net.UDPConn
	Cache         *cache.Cache
	TCPDeadline   int
	TCPTimeout    int
	UDPDeadline   int
	RunnerGroup   *runnergroup.RunnerGroup
}

// NewTproxy.
func NewTproxy(addr, remote, password string, tcpTimeout, tcpDeadline, udpDeadline int) (*Tproxy, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	uaddr, err := net.ResolveUDPAddr("udp", addr)
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
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &Tproxy{
		Password:      []byte(password),
		TCPAddr:       taddr,
		UDPAddr:       uaddr,
		RemoteTCPAddr: rtaddr,
		RemoteUDPAddr: ruaddr,
		Cache:         cs,
		TCPTimeout:    tcpTimeout,
		TCPDeadline:   tcpDeadline,
		UDPDeadline:   udpDeadline,
		RunnerGroup:   runnergroup.New(),
	}
	return s, nil
}

func (s *Tproxy) RunAutoScripts() error {
	hc := &http.Client{
		Timeout: 9 * time.Second,
	}
	r, err := hc.Get("https://txthinking.github.io/blackwhite/white_cidr.list")
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body.Close()
	data = bytes.TrimSpace(data)
	data = bytes.Replace(data, []byte{0x20}, []byte{}, -1)
	data = bytes.Replace(data, []byte{0x0d, 0x0a}, []byte{0x0a}, -1)
	cidrl := strings.Split(string(data), "\n")

	c := exec.Command("sh", "-c", "ip rule add fwmark 1 lookup 100")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "ip route add local 0.0.0.0/0 dev lo table 100")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "echo 1 > /proc/sys/net/ipv4/ip_forward")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "echo 1 > /proc/sys/net/ipv6/conf/all/forwarding")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "modprobe xt_socket")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "modprobe xt_TPROXY")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d 0.0.0.0/8 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d 10.0.0.0/8 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d 127.0.0.0/8 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d 169.254.0.0/16 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d 172.16.0.0/12 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d 192.168.0.0/16 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d 224.0.0.0/4 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d 240.0.0.0/4 -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	for _, v := range cidrl {
		c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d "+v+" -j RETURN")
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -d "+s.RemoteTCPAddr.IP.String()+" -j RETURN")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -p tcp -m socket -j MARK --set-mark 1")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -p tcp -j TPROXY --tproxy-mark 0x1/0x1 --on-port "+strconv.Itoa(s.TCPAddr.Port))
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -p udp -m socket -j MARK --set-mark 1")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -A PREROUTING -p udp -j TPROXY --tproxy-mark 0x1/0x1 --on-port "+strconv.Itoa(s.TCPAddr.Port))
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}

func (s *Tproxy) ClearAutoScripts() error {
	c := exec.Command("sh", "-c", "ip rule del fwmark 1 lookup 100")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "ip route del local 0.0.0.0/0 dev lo table 100")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -F")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("sh", "-c", "iptables -t mangle -X")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
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
				if err := c.SetKeepAlivePeriod(time.Duration(s.TCPTimeout) * time.Second); err != nil {
					log.Println(err)
					return
				}
			}
			if s.TCPDeadline != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
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
		b := make([]byte, 65535)
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

// TCPHandle handles request.
func (s *Tproxy) TCPHandle(c *net.TCPConn) error {
	tmp, err := tproxy.DialTCP("tcp", s.RemoteTCPAddr.String())
	if err != nil {
		return err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(s.TCPTimeout) * time.Second); err != nil {
			return err
		}
	}
	if s.TCPDeadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
			return err
		}
	}

	k, n, err := PrepareKey(s.Password)
	if err != nil {
		return err
	}
	if _, err := rc.Write(n); err != nil {
		return err
	}

	a, address, port, err := socks5.ParseAddress(c.LocalAddr().String())
	if err != nil {
		return err
	}
	ra := make([]byte, 0, 7)
	ra = append(ra, a)
	ra = append(ra, address...)
	ra = append(ra, port...)
	n, _, err = WriteTo(rc, ra, k, n, true)
	if err != nil {
		return err
	}

	go func() {
		n := make([]byte, 12)
		if _, err := io.ReadFull(rc, n); err != nil {
			return
		}
		k, err := GetKey(s.Password, n)
		if err != nil {
			log.Println(err)
			return
		}
		var b []byte
		for {
			if s.TCPDeadline != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
					return
				}
			}
			b, n, err = ReadFrom(rc, k, n, false)
			if err != nil {
				return
			}
			if _, err := c.Write(b); err != nil {
				return
			}
		}
	}()

	var b [1024 * 2]byte
	for {
		if s.TCPDeadline != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
				return nil
			}
		}
		i, err := c.Read(b[:])
		if err != nil {
			return nil
		}
		n, _, err = WriteTo(rc, b[0:i], k, n, false)
		if err != nil {
			return nil
		}
	}
	return nil
}

type TproxyUDPExchange struct {
	RemoteConn *net.UDPConn
	LocalConn  *net.UDPConn
}

func (s *Tproxy) UDPHandle(addr, daddr *net.UDPAddr, b []byte) error {
	a, address, port, err := socks5.ParseAddress(daddr.String())
	if err != nil {
		return err
	}
	ra := make([]byte, 0, 7)
	ra = append(ra, a)
	ra = append(ra, address...)
	ra = append(ra, port...)
	b = append(ra, b...)

	send := func(ue *TproxyUDPExchange, data []byte) error {
		cd, err := Encrypt(s.Password, data)
		if err != nil {
			return err
		}
		_, err = ue.RemoteConn.Write(cd)
		if err != nil {
			return err
		}
		return nil
	}

	var ue *TproxyUDPExchange
	iue, ok := s.Cache.Get(addr.String())
	if ok {
		ue = iue.(*TproxyUDPExchange)
		return send(ue, b)
	}

	rc, err := tproxy.DialUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 0,
	}, s.RemoteUDPAddr)
	if err != nil {
		return err
	}
	c, err := tproxy.DialUDP("udp", daddr, addr)
	if err != nil {
		rc.Close()
		return errors.New(fmt.Sprintf("src: %s dst: %s %s", daddr.String(), addr.String(), err.Error()))
	}
	ue = &TproxyUDPExchange{
		RemoteConn: rc,
		LocalConn:  c,
	}
	if err := send(ue, b); err != nil {
		ue.RemoteConn.Close()
		ue.LocalConn.Close()
		return err
	}
	s.Cache.Set(ue.LocalConn.RemoteAddr().String(), ue, cache.DefaultExpiration)
	go func(ue *TproxyUDPExchange) {
		defer func() {
			s.Cache.Delete(ue.LocalConn.RemoteAddr().String())
			ue.RemoteConn.Close()
			ue.LocalConn.Close()
		}()
		var b [65535]byte
		for {
			if s.UDPDeadline != 0 {
				if err := ue.RemoteConn.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
					break
				}
			}
			n, err := ue.RemoteConn.Read(b[:])
			if err != nil {
				break
			}
			_, _, _, data, err := Decrypt(s.Password, b[0:n])
			if err != nil {
				break
			}
			if _, err := ue.LocalConn.Write(data); err != nil {
				break
			}
		}
	}(ue)
	return nil
}
