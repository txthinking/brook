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
	"net"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/txthinking/brook/limits"
	crypto1 "github.com/txthinking/crypto"
	"github.com/txthinking/socks5"
)

type QUICClient struct {
	Server            *socks5.Server
	ServerHost        string
	ServerAddress     string
	TLSConfig         *tls.Config
	Password          []byte
	TCPTimeout        int
	UDPTimeout        int
	WithoutBrook      bool
	PacketConnFactory *PacketConnFactory
}

func NewQUICClient(addr, ip, server, password string, tcpTimeout, udpTimeout int, withoutbrook bool) (*QUICClient, error) {
	s5, err := socks5.NewClassicServer(addr, ip, "", "", tcpTimeout, udpTimeout)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
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
	p := []byte(password)
	if withoutbrook {
		p, err = crypto1.SHA256Bytes([]byte(password))
		if err != nil {
			return nil, err
		}
	}
	x := &QUICClient{
		ServerHost:        u.Host,
		Server:            s5,
		Password:          p,
		TCPTimeout:        tcpTimeout,
		UDPTimeout:        udpTimeout,
		WithoutBrook:      withoutbrook,
		PacketConnFactory: NewPacketConnFactory(),
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
		sa := x.ServerAddress
		if sa == "" {
			sa = x.ServerHost
		}
		rc, err := QUICDialTCP("", "", sa, x.ServerHost, x.TLSConfig, x.TCPTimeout)
		if err != nil {
			return ErrorReply(r, c, err)
		}
		defer rc.Close()
		dst := make([]byte, 0, 1+len(r.DstAddr)+2)
		dst = append(dst, r.Atyp)
		dst = append(dst, r.DstAddr...)
		dst = append(dst, r.DstPort...)
		var sc Exchanger
		if !x.WithoutBrook {
			sc, err = NewStreamClient("tcp", x.Password, c.RemoteAddr().String(), rc, x.TCPTimeout, dst)
		}
		if x.WithoutBrook {
			sc, err = NewSimpleStreamClient("tcp", x.Password, c.RemoteAddr().String(), rc, x.TCPTimeout, dst)
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
	dstb := append(append([]byte{d.Atyp}, d.DstAddr...), d.DstPort...)
	conn, err := x.PacketConnFactory.Handle(addr, dstb, d.Data, func(b []byte) (int, error) {
		d.Data = b
		return s.UDPConn.WriteToUDP(d.Bytes(), addr)
	}, x.UDPTimeout)
	if err != nil {
		return err
	}
	if conn == nil {
		return nil
	}
	defer conn.Close()
	sa := x.ServerAddress
	if sa == "" {
		sa = x.ServerHost
	}
	rc, err := QUICDialUDP(addr.String(), d.Address(), sa, x.ServerHost, x.TLSConfig, x.UDPTimeout)
	if err != nil {
		return err
	}
	defer rc.Close()
	var sc Exchanger
	if !x.WithoutBrook {
		sc, err = NewPacketClient(x.Password, addr.String(), rc, x.UDPTimeout, dstb)
	}
	if x.WithoutBrook {
		sc, err = NewSimplePacketClient(x.Password, addr.String(), rc, x.UDPTimeout, dstb)
	}
	if err != nil {
		return err
	}
	defer sc.Clean()
	if err := sc.Exchange(conn); err != nil {
		return nil
	}
	return nil
}

func (x *QUICClient) Shutdown() error {
	return x.Server.Shutdown()
}
