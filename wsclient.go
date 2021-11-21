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
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/socks5"
	x1 "github.com/txthinking/x"
)

// WSClient.
type WSClient struct {
	Server        *socks5.Server
	ServerHost    string
	ServerAddress string
	TLSConfig     *tls.Config
	Password      []byte
	TCPTimeout    int
	UDPTimeout    int
	TCPListen     *net.TCPListener
	Path          string
	UDPExchanges  *cache.Cache
	DialTCP       func(network, addr string) (net.Conn, error)
}

// NewWSClient.
func NewWSClient(addr, ip, server, password string, tcpTimeout, udpTimeout int) (*WSClient, error) {
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
	path := u.Path
	if path == "" {
		path = "/ws"
	}
	cs := cache.New(cache.NoExpiration, cache.NoExpiration)
	x := &WSClient{
		ServerHost:   u.Host,
		Server:       s5,
		Password:     []byte(password),
		TCPTimeout:   tcpTimeout,
		UDPTimeout:   udpTimeout,
		Path:         path,
		UDPExchanges: cs,
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

// ListenAndServe will let client start a socks5 proxy.
func (x *WSClient) ListenAndServe() error {
	return x.Server.ListenAndServe(x)
}

func (x *WSClient) DialWebsocket(src string) (net.Conn, error) {
	a := x.ServerHost
	if x.ServerAddress != "" {
		a = x.ServerAddress
	}
	var laddr *net.TCPAddr
	if src != "" {
		var err error
		laddr, err = net.ResolveTCPAddr("tcp", src)
		if err != nil {
			return nil, err
		}
	}
	raddr, err := net.ResolveTCPAddr("tcp", a)
	if err != nil {
		return nil, err
	}
	var c net.Conn
	if x.DialTCP != nil {
		c, err = x.DialTCP("tcp", a)
	}
	if x.DialTCP == nil {
		c, err = Dial.DialTCP("tcp", laddr, raddr)
	}
	if err != nil {
		return nil, err
	}
	if x.TCPTimeout != 0 {
		if err := c.SetDeadline(time.Now().Add(time.Duration(x.TCPTimeout) * time.Second)); err != nil {
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
	p := x1.BP16.Get().([]byte)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		x1.BP16.Put(p)
		c.Close()
		return nil, err
	}
	k := base64.StdEncoding.EncodeToString(p)
	x1.BP16.Put(p)
	b := make([]byte, 0, 300)
	b = append(b, []byte("GET "+x.Path+" HTTP/1.1\r\n")...)
	b = append(b, []byte(fmt.Sprintf("Host: %s\r\n", x.ServerHost))...)
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
	if r.Cmd == socks5.CmdConnect {
		debug("dial tcp", r.Address())
		rc, err := x.DialWebsocket("")
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

// UDPHandle handles udp request.
func (x *WSClient) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	src := addr.String()
	dst := d.Address()
	any, ok := s.UDPExchanges.Get(src + dst)
	if ok {
		ue := any.(*UDPExchange)
		return ue.Any.(func(b []byte) error)(d.Data)
	}
	debug("dial udp", dst)
	var laddr *net.UDPAddr
	any, ok = s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	la := ""
	if laddr != nil {
		la = laddr.String()
	}
	rc, err := x.DialWebsocket(la)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			// we dont choose lock, so ignore this error
			return nil
		}
		return err
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
	if err != nil {
		return err
	}
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
func (x *WSClient) Shutdown() error {
	return x.Server.Shutdown()
}
