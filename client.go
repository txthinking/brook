package brook

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/ant"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/socks5"
)

// Client
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
	TokenGetter     plugin.TokenGetter
}

// NewClient returns a new Client
func NewClient(addr, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) (*Client, error) {
	s5, err := socks5.NewClassicServer(addr, ip, "", "", tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return nil, err
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

// SetToken sets token plugin
func (x *Client) SetTokenGetter(token plugin.TokenGetter) {
	x.TokenGetter = token
}

// SetSocks5Middleman sets socks5middleman plugin
func (x *Client) SetSocks5Middleman(m plugin.Socks5Middleman) {
	x.Socks5Middleman = m
}

// SetHTTPMiddleman sets httpmiddleman plugin
func (x *Client) SetHTTPMiddleman(m plugin.HTTPMiddleman) {
	x.HTTPMiddleman = m
}

// ListenAndServe will let client start a socks5 proxy
// sm can be nil
func (x *Client) ListenAndServe() error {
	return x.Server.Run(x)
}

// TCPHandle handles tcp request
func (x *Client) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if x.Socks5Middleman != nil {
		done, err := x.Socks5Middleman.TCPHandle(s, c, r)
		if err != nil {
			if done {
				return err
			}
			return ErrorReply(r, c, err)
		}
		if done {
			return nil
		}
	}

	if r.Cmd == socks5.CmdConnect {
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

		rawaddr := make([]byte, 0, 7)
		rawaddr = append(rawaddr, r.Atyp)
		rawaddr = append(rawaddr, r.DstAddr...)
		rawaddr = append(rawaddr, r.DstPort...)
		if x.TokenGetter != nil {
			t, err := x.TokenGetter.Get()
			if err != nil {
				return ErrorReply(r, c, err)
			}
			if len(t) == 0 {
				return ErrorReply(r, c, errors.New("Miss Token"))
			}
			bb := make([]byte, 2)
			binary.BigEndian.PutUint16(bb, uint16(len(t)))
			t = append(bb, t...)
			rawaddr = append(t, rawaddr...)
		}
		n, err = WriteTo(rc, rawaddr, k, n, true)
		if err != nil {
			return ErrorReply(r, c, err)
		}

		a, address, port, err := socks5.ParseAddress(rc.LocalAddr().String())
		if err != nil {
			return ErrorReply(r, c, err)
		}
		rp := socks5.NewReply(socks5.RepSuccess, a, address, port)
		if err := rp.WriteTo(c); err != nil {
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
			n, err = WriteTo(rc, b[0:i], k, n, false)
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

// UDPHandle handles udp request
func (x *Client) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	if x.Socks5Middleman != nil {
		if done, err := x.Socks5Middleman.UDPHandle(s, addr, d); err != nil || done {
			return err
		}
	}

	send := func(ue *socks5.UDPExchange, data []byte) error {
		if x.TokenGetter != nil {
			t, err := x.TokenGetter.Get()
			if err != nil {
				return err
			}
			if len(t) == 0 {
				return errors.New("Miss Token")
			}
			bb := make([]byte, 2)
			binary.BigEndian.PutUint16(bb, uint16(len(t)))
			t = append(bb, t...)
			data = append(t, data...)
		}
		cd, err := Encrypt(x.Password, data)
		if err != nil {
			return err
		}
		_, err = ue.RemoteConn.Write(cd)
		if err != nil {
			return err
		}
		return nil
	}

	var ue *socks5.UDPExchange
	iue, ok := s.UDPExchanges.Get(addr.String())
	if ok {
		ue = iue.(*socks5.UDPExchange)
		return send(ue, d.Bytes()[3:])
	}

	c, err := Dial.Dial("udp", x.RemoteAddr)
	if err != nil {
		return err
	}
	rc := c.(*net.UDPConn)
	ue = &socks5.UDPExchange{
		ClientAddr: addr,
		RemoteConn: rc,
	}
	if err := send(ue, d.Bytes()[3:]); err != nil {
		return err
	}
	s.UDPExchanges.Set(ue.ClientAddr.String(), ue, cache.DefaultExpiration)
	go func(ue *socks5.UDPExchange) {
		defer func() {
			v, ok := s.TCPUDPAssociate.Get(ue.ClientAddr.String())
			if ok {
				ch := v.(chan byte)
				ch <- '0'
			}
			s.UDPExchanges.Delete(ue.ClientAddr.String())
			ue.RemoteConn.Close()
		}()
		var b [65536]byte
		for {
			if s.UDPDeadline != 0 {
				if err := ue.RemoteConn.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
					break
				}
			}
			n, err := ue.RemoteConn.Read(b[:])
			if err != nil {
				break
			}
			_, _, _, data, err := Decrypt(x.Password, b[0:n], nil)
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

// ListenAndServeHTTP will let client start a http proxy
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

// HTTPHandle handle http request
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
		addr, err = ant.GetAddressFromURL(aoru)
		if err != nil {
			return err
		}
	}

	if x.HTTPMiddleman != nil {
		if done, err := x.HTTPMiddleman.Handle(method, addr, b, c); err != nil || done {
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
	rawaddr := make([]byte, 0, 7)
	rawaddr = append(rawaddr, a)
	rawaddr = append(rawaddr, h...)
	rawaddr = append(rawaddr, p...)
	if x.TokenGetter != nil {
		t, err := x.TokenGetter.Get()
		if err != nil {
			return err
		}
		if len(t) == 0 {
			return errors.New("Miss Token")
		}
		bb := make([]byte, 2)
		binary.BigEndian.PutUint16(bb, uint16(len(t)))
		t = append(bb, t...)
		rawaddr = append(t, rawaddr...)
	}
	n, err = WriteTo(rc, rawaddr, k, n, true)
	if err != nil {
		return err
	}

	if method == "CONNECT" {
		_, err := c.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			return err
		}
	}
	if method != "CONNECT" {
		n, err = WriteTo(rc, b, k, n, false)
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
		n, err = WriteTo(rc, bf[0:i], k, n, false)
		if err != nil {
			return nil
		}
	}
	return nil
}

// Shutdown used to stop the client
func (x *Client) Shutdown() error {
	return x.Server.Stop()
}
