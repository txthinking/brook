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
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/txthinking/brook/limits"
	crypto1 "github.com/txthinking/crypto"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
	"golang.org/x/crypto/acme/autocert"
)

type QUICServer struct {
	Password             []byte
	Domain               string
	Addr                 string
	TCPTimeout           int
	UDPTimeout           int
	Cert                 []byte
	CertKey              []byte
	RunnerGroup          *runnergroup.RunnerGroup
	WithoutBrook         bool
	UDPServerConnFactory UDPServerConnFactory
}

func NewQUICServer(addr, password, domain string, tcpTimeout, udpTimeout int, withoutbrook bool) (*QUICServer, error) {
	if err := limits.Raise(); err != nil {
		Log(&Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	if runtime.GOOS == "linux" {
		c := exec.Command("sysctl", "-w", "net.core.rmem_max=2500000")
		b, err := c.CombinedOutput()
		if err != nil {
			Log(&Error{"when": "try to raise UDP Receive Buffer Size", "warning": string(b)})
		}
	}
	if runtime.GOOS == "darwin" {
		c := exec.Command("sysctl", "-w", "kern.ipc.maxsockbuf=3014656")
		b, err := c.CombinedOutput()
		if err != nil {
			Log(&Error{"when": "try to raise UDP Receive Buffer Size", "warning": string(b)})
		}
	}
	var p []byte
	var f UDPServerConnFactory
	if !withoutbrook {
		p = []byte(password)
		f = NewPacketServerConnFactory()
	}
	if withoutbrook {
		var err error
		p, err = crypto1.SHA256Bytes([]byte(password))
		if err != nil {
			return nil, err
		}
		f = NewSimplePacketServerConnFactory()
	}
	s := &QUICServer{
		Password:             p,
		Domain:               domain,
		Addr:                 addr,
		TCPTimeout:           tcpTimeout,
		UDPTimeout:           udpTimeout,
		UDPServerConnFactory: f,
		RunnerGroup:          runnergroup.New(),
		WithoutBrook:         withoutbrook,
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
							var ss Exchanger
							if !s.WithoutBrook {
								ss, err = NewStreamServer(s.Password, c.RemoteAddr().String(), c, s.TCPTimeout, s.UDPTimeout)
							}
							if s.WithoutBrook {
								ss, err = NewSimpleStreamServer(s.Password, c.RemoteAddr().String(), c, s.TCPTimeout, s.UDPTimeout)
							}
							if err != nil {
								Log(&Error{"from": c.RemoteAddr().String(), "error": err.Error()})
								return
							}
							defer ss.Clean()
							if err := s.TCPHandle(ss); err != nil {
								Log(&Error{"from": ss.Src(), "dst": ss.Dst(), "error": err.Error()})
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
							conn, dstb, err := s.UDPServerConnFactory.Handle(c.RemoteAddr().(*net.UDPAddr), b, s.Password, func(b []byte) (int, error) {
								if len(b) > 1197 {
									err := errors.New("when write to client, quic max datagram size is 1197")
									Log(&Error{"from": c.RemoteAddr().String(), "error": err.Error()})
									return 0, err
								}
								if err := c.SendMessage(b); err != nil {
									return 0, err
								}
								return len(b), nil
							}, s.UDPTimeout)
							if err != nil {
								Log(&Error{"from": c.RemoteAddr().String(), "error": err.Error()})
								continue
							}
							if conn == nil {
								continue
							}
							go func() {
								defer conn.Close()
								var ss Exchanger
								if !s.WithoutBrook {
									ss, err = NewPacketServer(s.Password, c.RemoteAddr().String(), conn, s.UDPTimeout, dstb)
								}
								if s.WithoutBrook {
									ss, err = NewSimplePacketServer(s.Password, c.RemoteAddr().String(), conn, s.UDPTimeout, dstb)
								}
								if err != nil {
									Log(&Error{"from": c.RemoteAddr().String(), "dst": socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), "error": err.Error()})
									return
								}
								defer ss.Clean()
								if err := s.UDPHandle(ss); err != nil {
									Log(&Error{"from": c.RemoteAddr().String(), "dst": socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), "error": err.Error()})
								}
							}()
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

func (s *QUICServer) TCPHandle(ss Exchanger) error {
	rc, err := DialTCP("tcp", "", ss.Dst())
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

func (s *QUICServer) UDPHandle(ss Exchanger) error {
	rc, err := NATDial("udp", ss.Src(), ss.Dst(), ss.Dst())
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

func (s *QUICServer) Shutdown() error {
	return s.RunnerGroup.Done()
}
