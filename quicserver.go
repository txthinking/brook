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
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/quic-go/quic-go"
	"github.com/txthinking/brook/limits"
	crypto1 "github.com/txthinking/crypto"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
	"golang.org/x/crypto/acme/autocert"
)

type QUICServer struct {
	Password       []byte
	Domain         string
	Addr           string
	TCPTimeout     int
	UDPTimeout     int
	UDPSrc         *cache.Cache
	UDPExchanges   *cache.Cache
	BlockDomain    map[string]byte
	BlockCIDR4     []*net.IPNet
	BlockCIDR6     []*net.IPNet
	BlockGeoIP     []string
	BlockCache     *cache.Cache
	BlockLock      *sync.RWMutex
	Done           chan byte
	Cert           []byte
	CertKey        []byte
	Dial           func(network, laddr, raddr string) (net.Conn, error)
	RunnerGroup    *runnergroup.RunnerGroup
	WithoutBrook   bool
	PasswordSha256 []byte
}

func NewQUICServer(addr, password, domain string, tcpTimeout, udpTimeout int, blockDomainList, blockCIDR4List, blockCIDR6List string, updateListInterval int64, blockGeoIP []string) (*QUICServer, error) {
	var err error
	var ds map[string]byte
	if blockDomainList != "" {
		ds, err = ReadDomainList(blockDomainList)
		if err != nil {
			return nil, err
		}
	}
	var c4 []*net.IPNet
	if blockCIDR4List != "" {
		c4, err = ReadCIDRList(blockCIDR4List)
		if err != nil {
			return nil, err
		}
	}
	var c6 []*net.IPNet
	if blockCIDR6List != "" {
		c6, err = ReadCIDRList(blockCIDR6List)
		if err != nil {
			return nil, err
		}
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs2 := cache.New(cache.NoExpiration, cache.NoExpiration)
	cs3 := cache.New(cache.NoExpiration, cache.NoExpiration)
	var lock *sync.RWMutex
	if updateListInterval != 0 {
		lock = &sync.RWMutex{}
	}
	done := make(chan byte)
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	if runtime.GOOS == "linux" {
		c := exec.Command("sysctl", "-w", "net.core.rmem_max=2500000")
		b, err := c.CombinedOutput()
		if err != nil {
			log.Println("Try to raise UDP Receive Buffer Size, got", err, string(b))
		}
	}
	if runtime.GOOS == "darwin" {
		c := exec.Command("sysctl", "-w", "kern.ipc.maxsockbuf=3014656")
		b, err := c.CombinedOutput()
		if err != nil {
			log.Println("Try to raise UDP Receive Buffer Size, got", err, string(b))
		}
	}
	b, err := crypto1.SHA256Bytes([]byte(password))
	if err != nil {
		return nil, err
	}
	s := &QUICServer{
		Password:       []byte(password),
		Domain:         domain,
		Addr:           addr,
		TCPTimeout:     tcpTimeout,
		UDPTimeout:     udpTimeout,
		UDPSrc:         cs2,
		UDPExchanges:   cs,
		BlockDomain:    ds,
		BlockCIDR4:     c4,
		BlockCIDR6:     c6,
		BlockGeoIP:     blockGeoIP,
		BlockCache:     cs3,
		BlockLock:      lock,
		Done:           done,
		RunnerGroup:    runnergroup.New(),
		PasswordSha256: b,
	}
	if updateListInterval != 0 {
		go func() {
			ticker := time.NewTicker(time.Duration(updateListInterval) * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					var ds map[string]byte
					if blockDomainList != "" {
						ds, err = ReadDomainList(blockDomainList)
						if err != nil {
							log.Println("ReadDomainList", blockDomainList, err)
							break
						}
					}
					var c4 []*net.IPNet
					if blockCIDR4List != "" {
						c4, err = ReadCIDRList(blockCIDR4List)
						if err != nil {
							log.Println("ReadCIDRList", blockCIDR4List, err)
							break
						}
					}
					var c6 []*net.IPNet
					if blockCIDR6List != "" {
						c6, err = ReadCIDRList(blockCIDR6List)
						if err != nil {
							log.Println("ReadCIDRList", blockCIDR6List, err)
							break
						}
					}
					lock.Lock()
					s.BlockDomain = ds
					s.BlockCIDR4 = c4
					s.BlockCIDR6 = c6
					if cs3 != nil {
						cs3.Flush()
					}
					lock.Unlock()
				}
			}
		}()
	}
	return s, nil
}

func (s *QUICServer) ListenAndServe() error {
	var t *tls.Config
	if s.Cert == nil || s.CertKey == nil {
		m := autocert.Manager{
			Cache:      autocert.DirCache(".letsencrypt"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(s.Domain),
			Email:      "cloud@txthinking.com",
		}
		server := &http.Server{Addr: ":80", Handler: m.HTTPHandler(nil)}
		s.RunnerGroup.Add(&runnergroup.Runner{
			Start: func() error {
				return server.ListenAndServe()
			},
			Stop: func() error {
				return server.Shutdown(context.Background())
			},
		})
		t = &tls.Config{GetCertificate: m.GetCertificate, ServerName: s.Domain, NextProtos: []string{"h3"}}
	}
	if s.Cert != nil && s.CertKey != nil {
		ct, err := tls.X509KeyPair(s.Cert, s.CertKey)
		if err != nil {
			return err
		}
		t = &tls.Config{Certificates: []tls.Certificate{ct}, ServerName: s.Domain, NextProtos: []string{"h3"}}
	}
	l, err := quic.ListenAddr(s.Addr, t, &quic.Config{MaxIdleTimeout: time.Duration(s.UDPTimeout) * time.Second, EnableDatagrams: true})
	if err != nil {
		return err
	}
	s.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			for {
				c, err := l.Accept(context.Background())
				if err != nil {
					return err
				}
				go func(c quic.Connection) {
					defer c.CloseWithError(0, "defer")
					for {
						st, err := c.AcceptStream(context.Background())
						if err != nil {
							return
						}
						go func(c net.Conn) {
							defer c.Close()
							if s.TCPTimeout != 0 {
								if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
									log.Println(err)
									return
								}
							}
							var ss Exchanger
							var dst []byte
							if !s.WithoutBrook {
								ss, dst, err = NewStreamServer(s.Password, c, s.TCPTimeout)
							}
							if s.WithoutBrook {
								ss, dst, err = NewSimpleStreamServer(s.PasswordSha256, c, s.TCPTimeout)
							}
							if err != nil {
								log.Println(err)
								return
							}
							defer ss.Clean()
							if err := s.TCPHandle(ss, dst); err != nil {
								log.Println(err)
							}
						}(&QUICConn{
							Conn:   c,
							Stream: st,
							LAddr: &net.TCPAddr{
								IP:   c.LocalAddr().(*net.UDPAddr).IP,
								Port: c.LocalAddr().(*net.UDPAddr).Port,
								Zone: c.LocalAddr().(*net.UDPAddr).Zone,
							},
							RAddr: &net.TCPAddr{
								IP:   c.RemoteAddr().(*net.UDPAddr).IP,
								Port: c.RemoteAddr().(*net.UDPAddr).Port,
								Zone: c.RemoteAddr().(*net.UDPAddr).Zone,
							},
						})
					}
				}(c)
				if c.ConnectionState().SupportsDatagrams {
					go func(c quic.Connection) {
						defer c.CloseWithError(0, "defer")
						for {
							b, err := c.ReceiveMessage()
							if err != nil {
								return
							}
							go func(addr net.Addr, b []byte, c quic.Connection) {
								if err := s.UDPHandle(addr, b, c); err != nil {
									log.Println(err)
									return
								}
							}(c.RemoteAddr(), b, c)
						}
					}(c)
				}
			}
			return nil
		},
		Stop: func() error {
			return l.Close()
		},
	})
	go func() {
		time.Sleep(1 * time.Second)
		_, _ = quic.DialAddr(net.JoinHostPort(s.Domain, s.Addr[1:]), &tls.Config{NextProtos: []string{"h3"}}, nil)
	}()
	return s.RunnerGroup.Wait()
}

func (s *QUICServer) TCPHandle(ss Exchanger, dst []byte) error {
	address := socks5.ToAddress(dst[0], dst[1:len(dst)-2], dst[len(dst)-2:])
	if Debug {
		log.Println("TCP", address)
	}
	var ds map[string]byte
	var c4 []*net.IPNet
	var c6 []*net.IPNet
	if s.BlockLock != nil {
		s.BlockLock.RLock()
	}
	ds = s.BlockDomain
	c4 = s.BlockCIDR4
	c6 = s.BlockCIDR6
	if s.BlockLock != nil {
		s.BlockLock.RUnlock()
	}
	if BlockAddress(address, ds, c4, c6, s.BlockCache, s.BlockGeoIP) {
		return errors.New("block " + address)
	}
	var rc net.Conn
	var err error
	if s.Dial == nil {
		rc, err = Dial.Dial("tcp", address)
	}
	if s.Dial != nil {
		rc, err = s.Dial("tcp", "", address)
	}
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

func (s *QUICServer) UDPHandle(addr net.Addr, b []byte, qc quic.Connection) error {
	src := addr.String()
	var dstb, d []byte
	var w WriterFunc
	var err error
	if !s.WithoutBrook {
		dstb, d, w, err = PacketClientToRemote(s.Password, b)
	}
	if s.WithoutBrook {
		dstb, d, w, err = SimplePacketClientToRemote(s.PasswordSha256, b)
	}
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
	if Debug {
		log.Println("UDP", dst)
	}
	var ds map[string]byte
	var c4 []*net.IPNet
	var c6 []*net.IPNet
	if s.BlockLock != nil {
		s.BlockLock.RLock()
	}
	ds = s.BlockDomain
	c4 = s.BlockCIDR4
	c6 = s.BlockCIDR6
	if s.BlockLock != nil {
		s.BlockLock.RUnlock()
	}
	if BlockAddress(dst, ds, c4, c6, s.BlockCache, s.BlockGeoIP) {
		return errors.New("block " + dst)
	}
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	raddr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return err
	}
	var rc net.Conn
	if s.Dial == nil {
		rc, err = Dial.DialUDP("udp", laddr, raddr)
		if err != nil {
			if !strings.Contains(err.Error(), "address already in use") {
				return err
			}
			rc, err = Dial.DialUDP("udp", nil, raddr)
			laddr = nil
		}
	}
	if s.Dial != nil {
		la := ""
		if laddr != nil {
			la = laddr.String()
		}
		rc, err = s.Dial("udp", la, dst)
		if err != nil {
			if !strings.Contains(err.Error(), "address already in use") {
				return err
			}
			rc, err = s.Dial("udp", "", dst)
			laddr = nil
		}
	}
	if err != nil {
		return err
	}
	defer rc.Close()
	if s.UDPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.UDPTimeout) * time.Second)); err != nil {
			return err
		}
	}
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
	var ps PacketServerT
	if !s.WithoutBrook {
		ps = NewPacketServer(s.Password)
	}
	if s.WithoutBrook {
		ps = NewSimplePacketServer(s.PasswordSha256)
	}
	defer ps.Clean()
	err = ps.RemoteToClient(rc, s.UDPTimeout, dstb, w(func(b []byte) (int, error) {
		if len(b) > 1197 {
			err := errors.New("quic max datagram size is 1197")
			log.Println(err)
			return 0, err
		}
		if err := qc.SendMessage(b); err != nil {
			return 0, err
		}
		return len(b), nil
	}))
	if err != nil {
		return err
	}
	return nil
}

func (s *QUICServer) Shutdown() error {
	close(s.Done)
	return s.RunnerGroup.Done()
}
