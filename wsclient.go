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
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/socks5"
)

// WSClient.
type WSClient struct {
	Server          *socks5.Server
	RemoteAddr      string
	RemoteAddress   string
	TLSConfig       *tls.Config
	Password        []byte
	TCPTimeout      int
	TCPDeadline     int
	UDPDeadline     int
	TCPListen       *net.TCPListener
	Socks5Middleman plugin.Socks5Middleman
	HTTPMiddleman   plugin.HTTPMiddleman
	Cache           *cache.Cache
	ClientAuthman   plugin.ClientAuthman
	Path            string
}

// NewWSClient.
func NewWSClient(addr, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) (*WSClient, error) {
	s5, err := socks5.NewClassicServer(addr, ip, "", "", tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	path := u.Path
	if path == "" {
		path = "/ws"
	}
	x := &WSClient{
		RemoteAddr:  u.Host,
		Server:      s5,
		Password:    []byte(password),
		TCPTimeout:  tcpTimeout,
		TCPDeadline: tcpDeadline,
		UDPDeadline: udpDeadline,
		Cache:       cs,
		Path:        path,
	}
	if u.Scheme == "wss" {
		h, _, err := net.SplitHostPort(u.Host)
		if err != nil {
			return nil, err
		}
		x.TLSConfig = &tls.Config{ServerName: h}
	}
	return x, nil
}

// SetSocks5Middleman sets socks5middleman plugin.
func (x *WSClient) SetSocks5Middleman(m plugin.Socks5Middleman) {
	x.Socks5Middleman = m
}

// SetHTTPMiddleman sets httpmiddleman plugin.
func (x *WSClient) SetHTTPMiddleman(m plugin.HTTPMiddleman) {
	x.HTTPMiddleman = m
}

// SetClientAuthman sets authman plugin.
func (x *WSClient) SetClientAuthman(m plugin.ClientAuthman) {
	x.ClientAuthman = m
}

// ListenAndServe will let client start a socks5 proxy.
func (x *WSClient) ListenAndServe() error {
	return x.Server.ListenAndServe(x)
}

func (x *WSClient) DialWebsocket() (net.Conn, error) {
	a := x.RemoteAddr
	if x.RemoteAddress != "" {
		a = x.RemoteAddress
	}
	c, err := Dial.Dial("tcp", a)
	if err != nil {
		return nil, err
	}
	if x.TCPTimeout != 0 {
		if err := c.(*net.TCPConn).SetKeepAlivePeriod(time.Duration(x.TCPTimeout) * time.Second); err != nil {
			c.Close()
			return nil, err
		}
	}
	if x.TCPDeadline != 0 {
		if err := c.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
			c.Close()
			return nil, err
		}
	}
	if x.TLSConfig != nil {
		tc := tls.Client(c, x.TLSConfig)
		if err := tc.Handshake(); err != nil {
			tc.Close()
			return nil, err
		}
		if err := tc.VerifyHostname(x.TLSConfig.ServerName); err != nil {
			tc.Close()
			return nil, err
		}
		c = tc
	}
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		c.Close()
		return nil, err
	}
	k := base64.StdEncoding.EncodeToString(p)
	b := make([]byte, 0)
	b = append(b, []byte("GET "+x.Path+" HTTP/1.1\r\n")...)
	b = append(b, []byte(fmt.Sprintf("Host: %s\r\n", x.RemoteAddr))...)
	b = append(b, []byte("Upgrade: websocket\r\n")...)
	b = append(b, []byte("Connection: Upgrade\r\n")...)
	b = append(b, []byte(fmt.Sprintf("Sec-WebSocket-Key: %s\r\n", k))...)
	b = append(b, []byte("Sec-WebSocket-Version: 13\r\n\r\n")...)
	if _, err := c.Write(b); err != nil {
		c.Close()
		return nil, err
	}
	r := bufio.NewReader(c)
	for {
		b, err = r.ReadBytes('\n')
		if err != nil {
			c.Close()
			return nil, err
		}
		b = bytes.TrimSpace(b)
		if len(b) == 0 {
			break
		}
		if bytes.HasPrefix(b, []byte("HTTP/1.1 ")) {
			if !bytes.Contains(b, []byte("101")) {
				c.Close()
				return nil, errors.New(string(b))
			}
		}
		if bytes.HasPrefix(b, []byte("Sec-WebSocket-Accept: ")) {
			h := sha1.New()
			h.Write([]byte(k))
			h.Write([]byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
			ak := base64.StdEncoding.EncodeToString(h.Sum(nil))
			if string(b[len("Sec-WebSocket-Accept: "):]) != ak {
				c.Close()
				return nil, errors.New(string(b))
			}
		}
	}
	return c, nil
}

// TCPHandle handles tcp request.
func (x *WSClient) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
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
		rc, err := x.DialWebsocket()
		if err != nil {
			return ErrorReply(r, c, err)
		}
		defer func() {
			rc.Close()
		}()
		// deer as horse
		cd, err := EncryptLength(x.Password, []byte{0x00})
		if err != nil {
			return ErrorReply(r, c, err)
		}
		if _, err := rc.Write(cd); err != nil {
			return ErrorReply(r, c, err)
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
		_, p, err := net.SplitHostPort(caddr.String())
		if err != nil {
			return err
		}
		if p == "0" {
			time.Sleep(time.Duration(x.Server.UDPSessionTime) * time.Second)
			return nil
		}
		ch := make(chan byte)
		x.Server.TCPUDPAssociate.Set(caddr.String(), ch, cache.DefaultExpiration)
		<-ch
		return nil
	}
	return socks5.ErrUnsupportCmd
}

type WSClientUDPExchange struct {
	ClientAddr *net.UDPAddr
	RemoteConn net.Conn
}

// UDPHandle handles udp request.
func (x *WSClient) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	if x.Socks5Middleman != nil {
		done, err := x.Socks5Middleman.UDPHandle(s, addr, d)
		if done {
			return err
		}
		if err != nil {
			v, ok := s.TCPUDPAssociate.Get(addr.String())
			if ok {
				ch := v.(chan byte)
				ch <- 0x00
				s.TCPUDPAssociate.Delete(addr.String())
			}
			return err
		}
	}

	send := func(ue *WSClientUDPExchange, data []byte) error {
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
		cd, err := EncryptLength(x.Password, data)
		if err != nil {
			return err
		}
		if _, err := ue.RemoteConn.Write(cd); err != nil {
			return err
		}
		cd, err = Encrypt(x.Password, data)
		if err != nil {
			return err
		}
		if _, err := ue.RemoteConn.Write(cd); err != nil {
			return err
		}
		return nil
	}

	var ue *WSClientUDPExchange
	iue, ok := x.Cache.Get(addr.String())
	if ok {
		ue = iue.(*WSClientUDPExchange)
		return send(ue, d.Bytes()[3:])
	}

	rc, err := x.DialWebsocket()
	if err != nil {
		v, ok := s.TCPUDPAssociate.Get(addr.String())
		if ok {
			ch := v.(chan byte)
			ch <- 0x00
			s.TCPUDPAssociate.Delete(addr.String())
		}
		return err
	}
	if x.UDPDeadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(x.UDPDeadline) * time.Second)); err != nil {
			v, ok := s.TCPUDPAssociate.Get(addr.String())
			if ok {
				ch := v.(chan byte)
				ch <- 0x00
				s.TCPUDPAssociate.Delete(addr.String())
			}
			rc.Close()
			return err
		}
	}
	// deer as horse
	cd, err := EncryptLength(x.Password, []byte{0x00, 0x00})
	if err != nil {
		v, ok := s.TCPUDPAssociate.Get(addr.String())
		if ok {
			ch := v.(chan byte)
			ch <- 0x00
			s.TCPUDPAssociate.Delete(addr.String())
		}
		rc.Close()
		return err
	}
	if _, err := rc.Write(cd); err != nil {
		v, ok := s.TCPUDPAssociate.Get(addr.String())
		if ok {
			ch := v.(chan byte)
			ch <- 0x00
			s.TCPUDPAssociate.Delete(addr.String())
		}
		rc.Close()
		return err
	}

	ue = &WSClientUDPExchange{
		ClientAddr: addr,
		RemoteConn: rc,
	}
	if err := send(ue, d.Bytes()[3:]); err != nil {
		v, ok := s.TCPUDPAssociate.Get(ue.ClientAddr.String())
		if ok {
			ch := v.(chan byte)
			ch <- 0x00
			s.TCPUDPAssociate.Delete(ue.ClientAddr.String())
		}
		ue.RemoteConn.Close()
		return err
	}
	x.Cache.Set(ue.ClientAddr.String(), ue, cache.DefaultExpiration)
	go func(ue *WSClientUDPExchange) {
		defer func() {
			v, ok := s.TCPUDPAssociate.Get(ue.ClientAddr.String())
			if ok {
				ch := v.(chan byte)
				ch <- 0x00
				s.TCPUDPAssociate.Delete(ue.ClientAddr.String())
			}
			x.Cache.Delete(ue.ClientAddr.String())
			ue.RemoteConn.Close()
		}()
		for {
			if s.UDPDeadline != 0 {
				if err := ue.RemoteConn.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
					break
				}
			}
			b := make([]byte, 12+16+10+2)
			if _, err := io.ReadFull(ue.RemoteConn, b); err != nil {
				break
			}
			l, err := DecryptLength(x.Password, b)
			if err != nil {
				log.Println(err)
				break
			}
			b = make([]byte, l)
			if _, err := io.ReadFull(ue.RemoteConn, b); err != nil {
				break
			}
			_, _, _, data, err := Decrypt(x.Password, b)
			if err != nil {
				log.Println(err)
				break
			}
			a, addr, port, err := socks5.ParseAddress(ue.ClientAddr.String())
			if err != nil {
				log.Println(err)
				break
			}
			d1 := socks5.NewDatagram(a, addr, port, data)
			if _, err := s.UDPConn.WriteToUDP(d1.Bytes(), ue.ClientAddr); err != nil {
				break
			}
		}
	}(ue)
	return nil
}

// ListenAndServeHTTP will let client start a http proxy.
func (x *WSClient) ListenAndServeHTTP() error {
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
func (x *WSClient) HTTPHandle(c *net.TCPConn) error {
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

	if x.HTTPMiddleman != nil {
		done, err := x.HTTPMiddleman.Handle(method, addr, b, c)
		if done {
			return err
		}
		if err != nil {
			return err
		}
	}

	rc, err := x.DialWebsocket()
	if err != nil {
		return err
	}
	defer func() {
		rc.Close()
	}()
	// deer as horse
	cd, err := EncryptLength(x.Password, []byte{0x00})
	if err != nil {
		return err
	}
	if _, err := rc.Write(cd); err != nil {
		return err
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
func (x *WSClient) Shutdown() error {
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
