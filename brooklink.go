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
	"net"
	"net/url"

	"github.com/txthinking/brook/limits"
	crypto1 "github.com/txthinking/crypto"
	"github.com/txthinking/socks5"
)

type BrookLink struct {
	Kind     string
	Address  string
	Host     string
	Path     string
	Password []byte
	V        url.Values
	Tc       *tls.Config

	S5         *socks5.Server
	Pcf        *PacketConnFactory
	TCPTimeout int
	UDPTimeout int
}

func NewBrookLink(link string) (*BrookLink, error) {
	var address, host, path string
	kind, server, v, err := ParseLink(link)
	if err != nil {
		return nil, err
	}
	p := []byte(v.Get("password"))
	if kind == "server" {
		address = server
	}
	var tc *tls.Config
	if kind == "socks5" || kind == "wsserver" || kind == "wssserver" || kind == "quicserver" {
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
		if kind == "wsserver" || kind == "wssserver" || kind == "quicserver" {
			if v.Get("address") != "" {
				address = v.Get("address")
			}
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
		if kind == "wsserver" || kind == "wssserver" || kind == "quicserver" {
			if v.Get("withoutBrookProtocol") == "true" {
				p, err = crypto1.SHA256Bytes([]byte(v.Get("password")))
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return &BrookLink{
		Kind:     kind,
		Address:  address,
		Host:     host,
		Path:     path,
		Password: p,
		V:        v,
		Tc:       tc,
	}, nil
}

// The caller is responsible for closing the connection in Exchanger.
func (blk *BrookLink) CreateExchanger(network, src string, dstb []byte, tcptimeout, udptimeout int) (Exchanger, error) {
	if blk.Kind == "server" {
		if network == "tcp" {
			rc, err := DialTCP("tcp", "", blk.Address)
			if err != nil {
				return nil, err
			}
			sc, err := NewStreamClient("tcp", blk.Password, src, rc, tcptimeout, dstb)
			if err != nil {
				rc.Close()
				return nil, err
			}
			return sc, nil
		}
		if blk.V.Get("udpovertcp") == "true" {
			rc, err := NATDial("tcp", src, socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), blk.Address)
			if err != nil {
				return nil, err
			}
			sc, err := NewStreamClient("udp", blk.Password, src, rc, udptimeout, dstb)
			if err != nil {
				rc.Close()
				return nil, err
			}
			return sc, nil
		}
		rc, err := NATDial("udp", src, socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), blk.Address)
		if err != nil {
			return nil, err
		}
		sc, err := NewPacketClient(blk.Password, src, rc, udptimeout, dstb)
		if err != nil {
			rc.Close()
			return nil, err
		}
		return sc, nil
	}
	if blk.Kind == "wsserver" || blk.Kind == "wssserver" {
		if network == "tcp" {
			rc, err := WebSocketDial("", "", blk.Address, blk.Host, blk.Path, blk.Tc, tcptimeout)
			if err != nil {
				return nil, err
			}
			var sc Exchanger
			if blk.V.Get("withoutBrookProtocol") != "true" {
				sc, err = NewStreamClient("tcp", blk.Password, src, rc, tcptimeout, dstb)
			}
			if blk.V.Get("withoutBrookProtocol") == "true" {
				sc, err = NewSimpleStreamClient("tcp", blk.Password, src, rc, tcptimeout, dstb)
			}
			if err != nil {
				rc.Close()
				return nil, err
			}
			return sc, nil
		}
		rc, err := WebSocketDial(src, socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), blk.Address, blk.Host, blk.Path, blk.Tc, tcptimeout)
		if err != nil {
			return nil, err
		}
		var sc Exchanger
		if blk.V.Get("withoutBrookProtocol") != "true" {
			sc, err = NewStreamClient("udp", blk.Password, src, rc, udptimeout, dstb)
		}
		if blk.V.Get("withoutBrookProtocol") == "true" {
			sc, err = NewSimpleStreamClient("udp", blk.Password, src, rc, udptimeout, dstb)
		}
		if err != nil {
			rc.Close()
			return nil, err
		}
		return sc, nil
	}
	if blk.Kind == "quicserver" {
		if network == "tcp" {
			rc, err := QUICDialTCP("", "", blk.Address, blk.Host, blk.Tc, tcptimeout)
			if err != nil {
				return nil, err
			}
			var sc Exchanger
			if blk.V.Get("withoutBrookProtocol") != "true" {
				sc, err = NewStreamClient("tcp", blk.Password, src, rc, tcptimeout, dstb)
			}
			if blk.V.Get("withoutBrookProtocol") == "true" {
				sc, err = NewSimpleStreamClient("tcp", blk.Password, src, rc, tcptimeout, dstb)
			}
			if err != nil {
				rc.Close()
				return nil, err
			}
			return sc, nil
		}
		rc, err := QUICDialUDP(src, socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:]), blk.Address, blk.Host, blk.Tc, udptimeout)
		if err != nil {
			return nil, err
		}
		var sc Exchanger
		if blk.V.Get("withoutBrookProtocol") != "true" {
			sc, err = NewPacketClient(blk.Password, src, rc, udptimeout, dstb)
		}
		if blk.V.Get("withoutBrookProtocol") == "true" {
			sc, err = NewSimplePacketClient(blk.Password, src, rc, udptimeout, dstb)
		}
		if err != nil {
			rc.Close()
			return nil, err
		}
		return sc, nil
	}
	return nil, errors.New("cannot create exchanger from " + blk.Kind)
}

func (x *BrookLink) PrepareSocks5Server(addr, ip string, tcptimeout, udptimeout int) error {
	var err error
	x.S5, err = socks5.NewClassicServer(addr, ip, "", "", tcptimeout, udptimeout)
	if err != nil {
		return err
	}
	if err := limits.Raise(); err != nil {
		Log(&Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	x.Pcf = NewPacketConnFactory()
	x.TCPTimeout = tcptimeout
	x.UDPTimeout = udptimeout
	return nil
}

func (x *BrookLink) ListenAndServe() error {
	return x.S5.ListenAndServe(x)
}

func (x *BrookLink) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if r.Cmd == socks5.CmdConnect {
		dstb := append(append([]byte{r.Atyp}, r.DstAddr...), r.DstPort...)
		sc, err := x.CreateExchanger("tcp", c.RemoteAddr().String(), dstb, x.TCPTimeout, x.UDPTimeout)
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
		_, err := r.UDP(c, x.S5.ServerAddr)
		if err != nil {
			return err
		}
		return nil
	}
	return socks5.ErrUnsupportCmd
}

func (x *BrookLink) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	dstb := append(append([]byte{d.Atyp}, d.DstAddr...), d.DstPort...)
	conn, err := x.Pcf.Handle(addr, dstb, d.Data, func(b []byte) (int, error) {
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
	sc, err := x.CreateExchanger("udp", addr.String(), dstb, x.TCPTimeout, x.UDPTimeout)
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
	return x.S5.Shutdown()
}
