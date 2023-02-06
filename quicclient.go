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
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	crypto1 "github.com/txthinking/crypto"
	"github.com/txthinking/socks5"
)

type QUICClient struct {
	Server         *socks5.Server
	ServerHost     string
	ServerAddress  string
	TLSConfig      *tls.Config
	Password       []byte
	TCPTimeout     int
	UDPTimeout     int
	TCPListen      *net.TCPListener
	UDPExchanges   *cache.Cache
	WithoutBrook   bool
	PasswordSha256 []byte
}

func NewQUICClient(addr, ip, server, password string, tcpTimeout, udpTimeout int) (*QUICClient, error) {
	s5, err := socks5.NewClassicServer(addr, ip, "", "", tcpTimeout, udpTimeout)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
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
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	b, err := crypto1.SHA256Bytes([]byte(password))
	if err != nil {
		return nil, err
	}
	x := &QUICClient{
		ServerHost:     u.Host,
		Server:         s5,
		Password:       []byte(password),
		TCPTimeout:     tcpTimeout,
		UDPTimeout:     udpTimeout,
		UDPExchanges:   cs,
		PasswordSha256: b,
	}
	h, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return nil, err
	}
	x.TLSConfig = &tls.Config{ServerName: h, NextProtos: []string{"h3"}}
	return x, nil
}

func (x *QUICClient) ListenAndServe() error {
	return x.Server.ListenAndServe(x)
}

func (x *QUICClient) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if r.Cmd == socks5.CmdConnect {
		if Debug {
			log.Println("TCP", r.Address())
		}
		var err error
		var raddr *net.UDPAddr
		if x.ServerAddress == "" {
			raddr, err = net.ResolveUDPAddr("udp", x.ServerHost)
			if err != nil {
				return ErrorReply(r, c, err)
			}
		}
		if x.ServerAddress != "" {
			raddr, err = net.ResolveUDPAddr("udp", x.ServerAddress)
			if err != nil {
				return ErrorReply(r, c, err)
			}
		}
		rc, err := QUICDialTCP(raddr, x.ServerHost, x.TLSConfig, x.TCPTimeout)
		if err != nil {
			return ErrorReply(r, c, err)
		}
		defer rc.Close()
		if x.TCPTimeout != 0 {
			if err := rc.SetDeadline(time.Now().Add(time.Duration(x.TCPTimeout) * time.Second)); err != nil {
				return ErrorReply(r, c, err)
			}
		}
		dst := make([]byte, 0, 1+len(r.DstAddr)+2)
		dst = append(dst, r.Atyp)
		dst = append(dst, r.DstAddr...)
		dst = append(dst, r.DstPort...)
		var sc Exchanger
		if !x.WithoutBrook {
			sc, err = NewStreamClient("tcp", x.Password, dst, rc, x.TCPTimeout)
		}
		if x.WithoutBrook {
			sc, err = NewSimpleStreamClient("tcp", x.PasswordSha256, dst, rc, x.TCPTimeout)
		}
		if err != nil {
			return ErrorReply(r, c, err)
		}
		defer sc.Clean()
		a, address, port, err := socks5.ParseAddress(rc.LocalAddr().String())
		if err != nil {
			return ErrorReply(r, c, err)
		}
		rp := socks5.NewReply(socks5.RepSuccess, a, address, port)
		if _, err := rp.WriteTo(c); err != nil {
			return err
		}
		if err := sc.Exchange(c); err != nil {
			return nil
		}
		return nil
	}
	if r.Cmd == socks5.CmdUDP {
		_, err := r.UDP(c, x.Server.ServerAddr)
		if err != nil {
			return err
		}
		return nil
	}
	return socks5.ErrUnsupportCmd
}

func (x *QUICClient) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	if 12+4+1+len(d.DstAddr)+2+len(d.Data)+16 > 1197 {
		return errors.New("quic max datagram size is 1197")
	}
	src := addr.String()
	dst := d.Address()
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		return ue.Any.(PacketClientT).LocalToServer(ue.Dst, d.Data, ue.Conn, x.UDPTimeout)
	}
	if Debug {
		log.Println("UDP", dst)
	}
	var laddr, raddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	var err error
	if x.ServerAddress == "" {
		raddr, err = net.ResolveUDPAddr("udp", x.ServerHost)
		if err != nil {
			return err
		}
	}
	if x.ServerAddress != "" {
		raddr, err = net.ResolveUDPAddr("udp", x.ServerAddress)
		if err != nil {
			return err
		}
	}
	rc, err := QUICDialUDP(laddr, raddr, x.ServerHost, x.TLSConfig, x.UDPTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "address already in use") {
			return err
		}
		rc, err = QUICDialUDP(laddr, raddr, x.ServerHost, x.TLSConfig, x.UDPTimeout)
		if err != nil {
			return err
		}
		laddr = nil
	}
	defer rc.Close()
	if laddr == nil {
		s.UDPSrc.Set(src+dst, rc.LocalAddr().(*net.UDPAddr), -1)
	}
	dstb := make([]byte, 0, 1+len(d.DstAddr)+2)
	dstb = append(dstb, d.Atyp)
	dstb = append(dstb, d.DstAddr...)
	dstb = append(dstb, d.DstPort...)
	var pc PacketClientT
	if !x.WithoutBrook {
		pc = NewPacketClient(x.Password)
	}
	if x.WithoutBrook {
		pc = NewSimplePacketClient(x.PasswordSha256)
	}
	defer pc.Clean()
	if err := pc.LocalToServer(dstb, d.Data, rc, x.UDPTimeout); err != nil {
		return err
	}
	ue := &UDPExchange{
		Conn: rc,
		Any:  pc,
		Dst:  dstb,
	}
	s.UDPExchanges.Set(src+dst, ue, -1)
	defer s.UDPExchanges.Delete(src + dst)
	err = pc.RunServerToLocal(rc, s.UDPTimeout, func(dst, b []byte) (int, error) {
		d.Data = b
		return s.UDPConn.WriteToUDP(d.Bytes(), addr)
	})
	if err != nil {
		return nil
	}
	return nil
}

func (x *QUICClient) Shutdown() error {
	return x.Server.Shutdown()
}
