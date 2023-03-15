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
	"net"

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/socks5"
)

type Client struct {
	Server            *socks5.Server
	ServerAddress     string
	Password          []byte
	TCPTimeout        int
	UDPTimeout        int
	UDPOverTCP        bool
	PacketConnFactory *PacketConnFactory
}

func NewClient(addr, ip, server, password string, tcpTimeout, udpTimeout int) (*Client, error) {
	if err := limits.Raise(); err != nil {
		Log(Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	s5, err := socks5.NewClassicServer(addr, ip, "", "", tcpTimeout, udpTimeout)
	if err != nil {
		return nil, err
	}
	x := &Client{
		ServerAddress:     server,
		Server:            s5,
		Password:          []byte(password),
		TCPTimeout:        tcpTimeout,
		UDPTimeout:        udpTimeout,
		PacketConnFactory: NewPacketConnFactory(),
	}
	return x, nil
}

func (x *Client) ListenAndServe() error {
	return x.Server.ListenAndServe(x)
}

func (x *Client) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if r.Cmd == socks5.CmdConnect {
		rc, err := DialTCP("tcp", "", x.ServerAddress)
		if err != nil {
			return ErrorReply(r, c, err)
		}
		defer rc.Close()
		dst := make([]byte, 0, 1+len(r.DstAddr)+2)
		dst = append(dst, r.Atyp)
		dst = append(dst, r.DstAddr...)
		dst = append(dst, r.DstPort...)
		sc, err := NewStreamClient("tcp", x.Password, c.RemoteAddr().String(), rc, x.TCPTimeout, dst)
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

func (x *Client) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
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
	if x.UDPOverTCP {
		rc, err := NATDial("tcp", addr.String(), d.Address(), x.ServerAddress)
		if err != nil {
			return err
		}
		defer rc.Close()
		sc, err := NewStreamClient("udp", x.Password, addr.String(), rc, x.UDPTimeout, dstb)
		if err != nil {
			return err
		}
		defer sc.Clean()
		if err := sc.Exchange(conn); err != nil {
			return nil
		}
		return nil
	}
	rc, err := NATDial("udp", addr.String(), d.Address(), x.ServerAddress)
	if err != nil {
		return err
	}
	defer rc.Close()
	sc, err := NewPacketClient(x.Password, addr.String(), rc, x.UDPTimeout, dstb)
	if err != nil {
		return err
	}
	defer sc.Clean()
	if err := sc.Exchange(conn); err != nil {
		return nil
	}
	return nil
}

func (x *Client) Shutdown() error {
	return x.Server.Shutdown()
}
