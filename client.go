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

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/socks5"
)

// Client.
type Client struct {
	Server        *socks5.Server
	ServerAddress string
	Password      []byte
	TCPTimeout    int
	UDPTimeout    int
	UDPOverTCP    bool
}

// NewClient returns a new Client.
func NewClient(addr, ip, server, password string, tcpTimeout, udpTimeout int) (*Client, error) {
	s5, err := socks5.NewClassicServer(addr, ip, "", "", tcpTimeout, udpTimeout)
	if err != nil {
		return nil, err
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	x := &Client{
		ServerAddress: server,
		Server:        s5,
		Password:      []byte(password),
		TCPTimeout:    tcpTimeout,
		UDPTimeout:    udpTimeout,
	}
	return x, nil
}

// ListenAndServe will let client start a socks5 proxy.
func (x *Client) ListenAndServe() error {
	return x.Server.ListenAndServe(x)
}

// TCPHandle handles tcp request.
func (x *Client) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if r.Cmd == socks5.CmdConnect {
		if Debug {
			log.Println("TCP", r.Address())
		}
		var rc net.Conn
		var err error
		if Dial1 != nil {
			rc, err = Dial1.Dial("tcp", x.ServerAddress)
		}
		if Dial1 == nil {
			rc, err = Dial.Dial("tcp", x.ServerAddress)
		}
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
		sc, err := NewStreamClient("tcp", x.Password, dst, rc, x.TCPTimeout)
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

type UDPExchange struct {
	Conn net.Conn
	Any  interface{}
	Dst  []byte
}

// UDPHandle handles udp request.
func (x *Client) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	if x.UDPOverTCP {
		return x.UDPOverTCPHandle(s, addr, d)
	}
	src := addr.String()
	dst := d.Address()
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		return ue.Any.(*PacketClient).LocalToServer(ue.Dst, d.Data, ue.Conn, x.UDPTimeout)
	}
	if Debug {
		log.Println("UDP", dst)
	}
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	raddr, err := net.ResolveUDPAddr("udp", x.ServerAddress)
	if err != nil {
		return err
	}
	var rc net.Conn
	if Dial1 != nil {
		rc, err = Dial1.DialUDP("udp", laddr, raddr)
	}
	if Dial1 == nil {
		rc, err = Dial.DialUDP("udp", laddr, raddr)
	}
	if err != nil {
		if !strings.Contains(err.Error(), "address already in use") {
			return err
		}
		if Dial1 != nil {
			rc, err = Dial1.DialUDP("udp", nil, raddr)
		}
		if Dial1 == nil {
			rc, err = Dial.DialUDP("udp", nil, raddr)
		}
		if err != nil {
			return err
		}
		laddr = nil
	}
	defer rc.Close()
	if x.UDPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(x.UDPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	if laddr == nil {
		s.UDPSrc.Set(src+dst, rc.LocalAddr().(*net.UDPAddr), -1)
	}
	dstb := make([]byte, 0, 1+len(d.DstAddr)+2)
	dstb = append(dstb, d.Atyp)
	dstb = append(dstb, d.DstAddr...)
	dstb = append(dstb, d.DstPort...)
	pc := NewPacketClient(x.Password)
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

func (x *Client) UDPOverTCPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	src := addr.String()
	dst := d.Address()
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		return ue.Any.(func(b []byte) error)(d.Data)
	}
	if Debug {
		log.Println("UDP", dst)
	}
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	var laddr1 *net.TCPAddr
	if laddr != nil {
		laddr1 = &net.TCPAddr{
			IP:   laddr.IP,
			Port: laddr.Port,
			Zone: laddr.Zone,
		}
	}
	raddr1, err := net.ResolveTCPAddr("tcp", x.ServerAddress)
	if err != nil {
		return err
	}
	var rc net.Conn
	if Dial1 != nil {
		rc, err = Dial1.DialTCP("tcp", laddr1, raddr1)
	}
	if Dial1 == nil {
		rc, err = Dial.DialTCP("tcp", laddr1, raddr1)
	}
	if err != nil {
		if !strings.Contains(err.Error(), "address already in use") {
			return err
		}
		if Dial1 != nil {
			rc, err = Dial1.DialTCP("tcp", nil, raddr1)
		}
		if Dial1 == nil {
			rc, err = Dial.DialTCP("tcp", nil, raddr1)
		}
		if err != nil {
			return err
		}
		laddr = nil
	}
	defer rc.Close()
	if laddr == nil {
		laddr = &net.UDPAddr{
			IP:   rc.LocalAddr().(*net.TCPAddr).IP,
			Port: rc.LocalAddr().(*net.TCPAddr).Port,
			Zone: rc.LocalAddr().(*net.TCPAddr).Zone,
		}
		s.UDPSrc.Set(src+dst, laddr, -1)
	}
	if x.UDPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(x.UDPTimeout) * time.Second)); err != nil {
			return err
		}
	}

	dstb := make([]byte, 0, 1+len(d.DstAddr)+2)
	dstb = append(dstb, d.Atyp)
	dstb = append(dstb, d.DstAddr...)
	dstb = append(dstb, d.DstPort...)
	sc, err := NewStreamClient("udp", x.Password, dstb, rc, x.UDPTimeout)
	defer sc.Clean()
	ps, pi := NewPacketStream(func(b []byte) (int, error) {
		d.Data = b
		return s.UDPConn.WriteToUDP(d.Bytes(), addr)
	})
	defer ps.Close()
	ue := &UDPExchange{
		Any: pi,
	}
	s.UDPExchanges.Set(src+dst, ue, -1)
	defer s.UDPExchanges.Delete(src + dst)
	go pi(d.Data)
	if err := sc.Exchange(ps); err != nil {
		return nil
	}
	return nil
}

// Shutdown used to stop the client.
func (x *Client) Shutdown() error {
	return x.Server.Shutdown()
}
