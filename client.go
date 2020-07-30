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
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/socks5"
)

// Client.
type Client struct {
	Server          *socks5.Server
	RemoteAddr      string
	Password        []byte
	TCPTimeout      int
	TCPDeadline     int
	UDPDeadline     int
	TCPListen       *net.TCPListener
	Socks5Middleman plugin.Socks5Middleman
	HTTPMiddleman   plugin.HTTPMiddleman
	ClientAuthman   plugin.ClientAuthman
}

// NewClient returns a new Client.
func NewClient(addr, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) (*Client, error) {
	s5, err := socks5.NewClassicServer(addr, ip, "", "", tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return nil, err
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	x := &Client{
		RemoteAddr:  server,
		Server:      s5,
		Password:    []byte(password),
		TCPTimeout:  tcpTimeout,
		TCPDeadline: tcpDeadline,
		UDPDeadline: udpDeadline,
	}
	return x, nil
}

// SetSocks5Middleman sets socks5middleman plugin.
func (x *Client) SetSocks5Middleman(m plugin.Socks5Middleman) {
	x.Socks5Middleman = m
}

// SetHTTPMiddleman sets httpmiddleman plugin.
func (x *Client) SetHTTPMiddleman(m plugin.HTTPMiddleman) {
	x.HTTPMiddleman = m
}

// SetClientAuthman sets authman plugin.
func (x *Client) SetClientAuthman(m plugin.ClientAuthman) {
	x.ClientAuthman = m
}

// ListenAndServe will let client start a socks5 proxy.
func (x *Client) ListenAndServe() error {
	return x.Server.ListenAndServe(x)
}

// TCPHandle handles tcp request.
func (x *Client) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if x.Socks5Middleman != nil {
		done, err := x.Socks5Middleman.TCPHandle(s, c, r)
		if done {
			return err
		}
		if err != nil {
			return ErrorReply(r, c, err)
		}
	}

	if r.Cmd == socks5.CmdConnect {
		debug("dial tcp", r.Address())
		tmp, err := Dial.Dial("tcp", x.RemoteAddr)
		if err != nil {
			return ErrorReply(r, c, err)
		}
		rc := tmp.(*net.TCPConn)
		defer rc.Close()
		if x.TCPTimeout != 0 {
			if err := rc.SetKeepAlivePeriod(time.Duration(x.TCPTimeout) * time.Second); err != nil {
				return ErrorReply(r, c, err)
			}
		}
		if x.TCPDeadline != 0 {
			if err := rc.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
				return ErrorReply(r, c, err)
			}
		}

		k, n, err := PrepareKey(x.Password)
		if err != nil {
			return ErrorReply(r, c, err)
		}
		if _, err := rc.Write(n); err != nil {
			return ErrorReply(r, c, err)
		}

		ra := make([]byte, 0, 7)
		ra = append(ra, r.Atyp)
		ra = append(ra, r.DstAddr...)
		ra = append(ra, r.DstPort...)
		n, _, err = WriteTo(rc, ra, k, n, true)
		if err != nil {
			return ErrorReply(r, c, err)
		}

		if x.ClientAuthman != nil {
			b, err := x.ClientAuthman.GetToken()
			if err != nil {
				return ErrorReply(r, c, err)
			}
			n, _, err = WriteTo(rc, b, k, n, false)
			if err != nil {
				return ErrorReply(r, c, err)
			}
		}

		a, address, port, err := socks5.ParseAddress(rc.LocalAddr().String())
		if err != nil {
			return ErrorReply(r, c, err)
		}
		rp := socks5.NewReply(socks5.RepSuccess, a, address, port)
		if _, err := rp.WriteTo(c); err != nil {
			return err
		}

		go func() {
			n := make([]byte, 12)
			if _, err := io.ReadFull(rc, n); err != nil {
				return
			}
			k, err := GetKey(x.Password, n)
			if err != nil {
				log.Println(err)
				return
			}
			var b []byte
			for {
				if x.TCPDeadline != 0 {
					if err := rc.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
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
			if x.TCPDeadline != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
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
	if r.Cmd == socks5.CmdUDP {
		caddr, err := r.UDP(c, x.Server.ServerAddr)
		if err != nil {
			return err
		}
		// TODO
		_ = caddr
		return nil
	}
	return socks5.ErrUnsupportCmd
}

// UDPHandle handles udp request.
func (x *Client) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	if x.Socks5Middleman != nil {
		done, err := x.Socks5Middleman.UDPHandle(s, addr, d)
		if done {
			return err
		}
		if err != nil {
			return err
		}
	}

	src := addr.String()
	// any, ok := s.AssociatedUDP.Get(src)
	// if !ok {
	// return fmt.Errorf("This udp address %s is not associated with tcp", src)
	// }
	// ch := any.(chan byte)
	send := func(ue *socks5.UDPExchange, data []byte) error {
		select {
		// case <-ch:
		// 	return fmt.Errorf("This udp address %s is not associated with tcp", src)
		default:
			if x.ClientAuthman != nil {
				b, err := x.ClientAuthman.GetToken()
				if err != nil {
					return err
				}
				data = append(data, b...)
				bb := make([]byte, 2)
				binary.BigEndian.PutUint16(bb, uint16(len(b)))
				data = append(data, bb...)
			}
			cd, err := Encrypt(x.Password, data)
			if err != nil {
				return err
			}
			_, err = ue.RemoteConn.Write(cd)
			if err != nil {
				return err
			}
		}
		return nil
	}

	dst := d.Address()
	var ue *socks5.UDPExchange
	iue, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue = iue.(*socks5.UDPExchange)
		return send(ue, d.Bytes()[3:])
	}

	debug("dial udp", dst)
	var laddr *net.UDPAddr
	any, ok := s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	raddr, err := net.ResolveUDPAddr("udp", x.RemoteAddr)
	if err != nil {
		return err
	}
	rc, err := Dial.DialUDP("udp", laddr, raddr)
	if err != nil {
		return err
	}
	if laddr == nil {
		s.UDPSrc.Set(src+dst, rc.LocalAddr().(*net.UDPAddr), -1)
	}
	ue = &socks5.UDPExchange{
		ClientAddr: addr,
		RemoteConn: rc,
	}
	if err := send(ue, d.Bytes()[3:]); err != nil {
		ue.RemoteConn.Close()
		return err
	}
	s.UDPExchanges.Set(src+dst, ue, -1)
	go func(ue *socks5.UDPExchange, dst string) {
		defer func() {
			s.UDPExchanges.Delete(ue.ClientAddr.String() + dst)
			ue.RemoteConn.Close()
		}()
		var b [65535]byte
		for {
			select {
			// case <-ch:
			// 	if Debug {
			// 		log.Printf("The tcp that udp address %s associated closed\n", ue.ClientAddr.String())
			// 	}
			// 	return
			default:
				if s.UDPDeadline != 0 {
					if err := ue.RemoteConn.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
						return
					}
				}
				n, err := ue.RemoteConn.Read(b[:])
				if err != nil {
					return
				}
				_, _, _, data, err := Decrypt(x.Password, b[0:n])
				if err != nil {
					log.Println(err)
					return
				}
				a, addr, port, err := socks5.ParseAddress(ue.ClientAddr.String())
				if err != nil {
					log.Println(err)
					return
				}
				d1 := socks5.NewDatagram(a, addr, port, data)
				if _, err := s.UDPConn.WriteToUDP(d1.Bytes(), ue.ClientAddr); err != nil {
					return
				}
			}
		}
	}(ue, dst)
	return nil
}

// ListenAndServeHTTP will let client start a http proxy.
func (x *Client) ListenAndServeHTTP() error {
	var err error
	x.TCPListen, err = net.ListenTCP("tcp", x.Server.TCPAddr)
	if err != nil {
		return nil
	}
	for {
		c, err := x.TCPListen.AcceptTCP()
		if err != nil {
			return err
		}
		go func(c *net.TCPConn) {
			defer c.Close()
			if x.TCPTimeout != 0 {
				if err := c.SetKeepAlivePeriod(time.Duration(x.TCPTimeout) * time.Second); err != nil {
					log.Println(err)
					return
				}
			}
			if x.TCPDeadline != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
					log.Println(err)
					return
				}
			}
			if err := x.HTTPHandle(c); err != nil {
				log.Println(err)
				return
			}
		}(c)
	}
}

// HTTPHandle handles http request.
func (x *Client) HTTPHandle(c *net.TCPConn) error {
	b := make([]byte, 0, 1024)
	for {
		var b1 [1024]byte
		n, err := c.Read(b1[:])
		if err != nil {
			return err
		}
		b = append(b, b1[:n]...)
		if bytes.Contains(b, []byte{0x0d, 0x0a, 0x0d, 0x0a}) {
			break
		}
		if len(b) >= 2083+18 {
			return errors.New("HTTP header too long")
		}
	}
	bb := bytes.SplitN(b, []byte(" "), 3)
	if len(bb) != 3 {
		return errors.New("Invalid Request")
	}
	method, aoru := string(bb[0]), string(bb[1])
	var addr string
	if method == "CONNECT" {
		addr = aoru
	}
	if method != "CONNECT" {
		var err error
		addr, err = GetAddressFromURL(aoru)
		if err != nil {
			return err
		}
	}
	debug("dial http", addr)

	if x.HTTPMiddleman != nil {
		done, err := x.HTTPMiddleman.Handle(method, addr, b, c)
		if done {
			return err
		}
		if err != nil {
			return err
		}
	}

	tmp, err := Dial.Dial("tcp", x.RemoteAddr)
	if err != nil {
		return err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if x.TCPTimeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(x.TCPTimeout) * time.Second); err != nil {
			return err
		}
	}
	if x.TCPDeadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
			return err
		}
	}

	k, n, err := PrepareKey(x.Password)
	if err != nil {
		return err
	}
	if _, err := rc.Write(n); err != nil {
		return err
	}

	a, h, p, err := socks5.ParseAddress(addr)
	if err != nil {
		return err
	}
	ra := make([]byte, 0, 7)
	ra = append(ra, a)
	ra = append(ra, h...)
	ra = append(ra, p...)
	n, _, err = WriteTo(rc, ra, k, n, true)
	if err != nil {
		return err
	}

	if x.ClientAuthman != nil {
		b, err := x.ClientAuthman.GetToken()
		if err != nil {
			return err
		}
		n, _, err = WriteTo(rc, b, k, n, false)
		if err != nil {
			return err
		}
	}

	if method == "CONNECT" {
		_, err := c.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			return err
		}
	}
	if method != "CONNECT" {
		n, _, err = WriteTo(rc, b, k, n, false)
		if err != nil {
			return err
		}
	}

	go func() {
		n := make([]byte, 12)
		if _, err := io.ReadFull(rc, n); err != nil {
			return
		}
		k, err := GetKey(x.Password, n)
		if err != nil {
			log.Println(err)
			return
		}
		var b []byte
		for {
			if x.TCPDeadline != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
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

	var bf [1024 * 2]byte
	for {
		if x.TCPDeadline != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
				return nil
			}
		}
		i, err := c.Read(bf[:])
		if err != nil {
			return nil
		}
		n, _, err = WriteTo(rc, bf[0:i], k, n, false)
		if err != nil {
			return nil
		}
	}
	return nil
}

// Shutdown used to stop the client.
func (x *Client) Shutdown() error {
	var e error
	if x.TCPListen != nil {
		if err := x.TCPListen.Close(); err != nil {
			e = err
		}
	}
	if err := x.Server.Shutdown(); err != nil {
		e = err
	}
	return e
}
