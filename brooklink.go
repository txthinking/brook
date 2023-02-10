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
	"crypto/x509"
	"errors"
	"log"
	"net"
	"net/url"

	"github.com/txthinking/brook/limits"
	crypto1 "github.com/txthinking/crypto"
	"github.com/txthinking/socks5"
)

type BrookLink struct {
	kind     string
	address  string
	host     string
	path     string
	password []byte
	v        url.Values
	tc       *tls.Config

	s5  *socks5.Server
	pcf *PacketConnFactory
	tt  int
	ut  int
}

func NewBrookLink(link string) (*BrookLink, error) {
	var address, host, path string
	kind, server, _, password, v, err := ParseLinkExtra(link)
	if err != nil {
		return nil, err
	}
	p := []byte(password)
	if kind == "server" {
		address = server
	}
	var tc *tls.Config
	if kind == "wsserver" || kind == "wssserver" || kind == "quicserver" {
		u, err := url.Parse(server)
		if err != nil {
			return nil, err
		}
		host = u.Host
		path = u.Path
		if path == "" {
			path = "/ws"
		}
		address = host
		if v.Get("address") != "" {
			address = v.Get("address")
		}
		if kind == "wssserver" || kind == "quicserver" {
			h, _, err := net.SplitHostPort(u.Host)
			if err != nil {
				return nil, err
			}
			tc = &tls.Config{ServerName: h}
			if v.Get("insecure") == "true" {
				tc.InsecureSkipVerify = true
			}
			if v.Get("ca") != "" {
				roots := x509.NewCertPool()
				ok := roots.AppendCertsFromPEM([]byte(v.Get("ca")))
				if !ok {
					return nil, errors.New("failed to parse root certificate")
				}
				tc.RootCAs = roots
			}
			if kind == "quicserver" {
				tc.NextProtos = []string{"h3"}
			}
		}
		if v.Get("withoutBrookProtocol") == "true" {
			p, err = crypto1.SHA256Bytes([]byte(password))
			if err != nil {
				return nil, err
			}
		}
	}
	return &BrookLink{
		kind:     kind,
		address:  address,
		host:     host,
		path:     path,
		password: p,
		v:        v,
		tc:       tc,
	}, nil
}

// The caller is responsible for closing the connection in Exchanger.
func (blk *BrookLink) CreateExchanger(network, src string, dstb []byte, tcptimeout, udptimeout int) (Exchanger, error) {
	if blk.kind == "server" {
		if network == "tcp" {
			rc, err := DialTCP("tcp", "", blk.address)
			if err != nil {
				return nil, err
			}
			sc, err := NewStreamClient("tcp", blk.password, src, rc, tcptimeout, dstb)
			if err != nil {
				rc.Close()
				return nil, err
			}
			return sc, nil
		}
		if blk.v.Get("udpovertcp") == "true" {
			rc, err := NATDial("tcp", src, socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), blk.address)
			if err != nil {
				return nil, err
			}
			sc, err := NewStreamClient("udp", blk.password, src, rc, udptimeout, dstb)
			if err != nil {
				rc.Close()
				return nil, err
			}
			return sc, nil
		}
		rc, err := NATDial("udp", src, socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), blk.address)
		if err != nil {
			return nil, err
		}
		sc, err := NewPacketClient(blk.password, src, rc, udptimeout, dstb)
		if err != nil {
			rc.Close()
			return nil, err
		}
		return sc, nil
	}
	if blk.kind == "wsserver" || blk.kind == "wssserver" {
		if network == "tcp" {
			rc, err := WebSocketDial("", "", blk.address, blk.host, blk.path, blk.tc, tcptimeout)
			if err != nil {
				return nil, err
			}
			var sc Exchanger
			if blk.v.Get("withoutBrookProtocol") != "true" {
				sc, err = NewStreamClient("tcp", blk.password, src, rc, tcptimeout, dstb)
			}
			if blk.v.Get("withoutBrookProtocol") == "true" {
				sc, err = NewSimpleStreamClient("tcp", blk.password, src, rc, tcptimeout, dstb)
			}
			if err != nil {
				rc.Close()
				return nil, err
			}
			return sc, nil
		}
		rc, err := WebSocketDial(src, socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), blk.address, blk.host, blk.path, blk.tc, tcptimeout)
		if err != nil {
			return nil, err
		}
		var sc Exchanger
		if blk.v.Get("withoutBrookProtocol") != "true" {
			sc, err = NewStreamClient("udp", blk.password, src, rc, udptimeout, dstb)
		}
		if blk.v.Get("withoutBrookProtocol") == "true" {
			sc, err = NewSimpleStreamClient("udp", blk.password, src, rc, udptimeout, dstb)
		}
		if err != nil {
			rc.Close()
			return nil, err
		}
		return sc, nil
	}
	if blk.kind == "quicserver" {
		if network == "tcp" {
			rc, err := QUICDialTCP("", "", blk.address, blk.host, blk.tc, tcptimeout)
			if err != nil {
				return nil, err
			}
			var sc Exchanger
			if blk.v.Get("withoutBrookProtocol") != "true" {
				sc, err = NewStreamClient("tcp", blk.password, src, rc, tcptimeout, dstb)
			}
			if blk.v.Get("withoutBrookProtocol") == "true" {
				sc, err = NewSimpleStreamClient("tcp", blk.password, src, rc, tcptimeout, dstb)
			}
			if err != nil {
				rc.Close()
				return nil, err
			}
			return sc, nil
		}
		rc, err := QUICDialUDP(src, socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), blk.address, blk.host, blk.tc, udptimeout)
		if err != nil {
			return nil, err
		}
		var sc Exchanger
		if blk.v.Get("withoutBrookProtocol") != "true" {
			sc, err = NewPacketClient(blk.password, src, rc, udptimeout, dstb)
		}
		if blk.v.Get("withoutBrookProtocol") == "true" {
			sc, err = NewSimplePacketClient(blk.password, src, rc, udptimeout, dstb)
		}
		if err != nil {
			rc.Close()
			return nil, err
		}
		return sc, nil
	}
	return nil, errors.New("cannot create exchanger from " + blk.kind)
}

func (x *BrookLink) PrepareSocks5Server(addr, ip string, tcptimeout, udptimeout int) error {
	var err error
	x.s5, err = socks5.NewClassicServer(addr, ip, "", "", tcptimeout, udptimeout)
	if err != nil {
		return err
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	x.pcf = NewPacketConnFactory()
	x.tt = tcptimeout
	x.ut = udptimeout
	return nil
}

func (x *BrookLink) ListenAndServe() error {
	return x.s5.ListenAndServe(x)
}

func (x *BrookLink) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if r.Cmd == socks5.CmdConnect {
		dstb := append(append([]byte{r.Atyp}, r.DstAddr...), r.DstPort...)
		sc, err := x.CreateExchanger("tcp", c.RemoteAddr().String(), dstb, x.tt, x.ut)
		if err != nil {
			return ErrorReply(r, c, err)
		}
		defer sc.Clean()
		laddr := ""
		if v, ok := sc.(*StreamClient); ok {
			defer v.Server.Close()
			laddr = v.Server.LocalAddr().String()
		}
		if v, ok := sc.(*SimpleStreamClient); ok {
			defer v.Server.Close()
			laddr = v.Server.LocalAddr().String()
		}
		a, address, port, err := socks5.ParseAddress(laddr)
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
		_, err := r.UDP(c, x.s5.ServerAddr)
		if err != nil {
			return err
		}
		return nil
	}
	return socks5.ErrUnsupportCmd
}

func (x *BrookLink) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	dstb := append(append([]byte{d.Atyp}, d.DstAddr...), d.DstPort...)
	conn, err := x.pcf.Handle(addr, dstb, d.Data, func(b []byte) (int, error) {
		d.Data = b
		return s.UDPConn.WriteToUDP(d.Bytes(), addr)
	}, x.ut)
	if err != nil {
		return err
	}
	if conn == nil {
		return nil
	}
	defer conn.Close()
	sc, err := x.CreateExchanger("udp", addr.String(), dstb, x.tt, x.ut)
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
	if err := sc.Exchange(conn); err != nil {
		return nil
	}
	return nil
}

func (x *BrookLink) Shutdown() error {
	return x.s5.Shutdown()
}
