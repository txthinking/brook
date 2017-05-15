package brook

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

// SSClient is the client of shadowsocks protocol
type SSClient struct {
	Address         string
	Password        string
	Server          string
	Timeout         int
	Deadline        int
	Dial            Dialer
	Listen          net.Listener
	HTTPMiddleman   HTTPMiddleman
	Socks5Middleman Socks5Middleman
}

// NewSSClient returns a new SSClient
func NewSSClient(address, server, password string, timeout, deadline int, dial Dialer) *SSClient {
	if dial == nil {
		dial = &DefaultDial{}
	}
	c := &SSClient{
		Address:  address,
		Server:   server,
		Password: password,
		Timeout:  timeout,
		Deadline: deadline,
		Dial:     dial,
	}
	return c
}

// ListenAndServe will let client start to listen and serve
func (c *SSClient) ListenAndServe(sm Socks5Middleman) error {
	ta, err := net.ResolveTCPAddr("tcp", c.Address)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", ta)
	if err != nil {
		return err
	}
	defer l.Close()
	c.Listen = l
	c.Socks5Middleman = sm

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		go func(conn *net.TCPConn) {
			if err := c.handle(conn); err != nil {
				log.Println(err)
			}
		}(conn)
	}
}

// ListenAndServeHTTP will let client start a http(s) proxy to listen and serve.
// For just a http proxy server, so httpmiddleman can be nil
func (c *SSClient) ListenAndServeHTTP(h HTTPMiddleman) error {
	ta, err := net.ResolveTCPAddr("tcp", c.Address)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", ta)
	if err != nil {
		return err
	}
	defer l.Close()
	c.Listen = l
	c.HTTPMiddleman = h

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		go func(conn *net.TCPConn) {
			if err := c.handleHTTP(conn); err != nil {
				log.Println(err)
			}
		}(conn)
	}
}

// Shutdown used to stop the client
func (c *SSClient) Shutdown() error {
	if c.Listen == nil {
		return nil
	}
	return c.Listen.Close()
}

func (c *SSClient) handle(conn *net.TCPConn) error {
	defer conn.Close()
	if c.Timeout != 0 {
		if err := conn.SetKeepAlivePeriod(time.Duration(c.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if c.Deadline != 0 {
		if err := conn.SetDeadline(time.Now().Add(time.Duration(c.Deadline) * time.Second)); err != nil {
			return err
		}
	}

	s := socks5.NewClassicServer(conn)
	if err := s.Negotiate(); err != nil {
		return err
	}
	request, err := s.GetRequest()
	if err != nil {
		return err
	}

	if c.Socks5Middleman != nil {
		if handled, err := c.Socks5Middleman.HandleSocks5Proxy(request, conn); err != nil || handled {
			return err
		}
	}

	rawaddr := make([]byte, 0)
	rawaddr = append(rawaddr, request.Atyp)
	rawaddr = append(rawaddr, request.DstAddr...)
	rawaddr = append(rawaddr, request.DstPort...)

	rc, err := c.Dial.Dial("tcp", c.Server)
	if err != nil {
		var p *socks5.Reply
		if request.Atyp == socks5.ATYPIPv4 || request.Atyp == socks5.ATYPDomain {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(conn); err != nil {
			return err
		}
		return err
	}
	defer rc.Close()
	if c.Timeout != 0 {
		if rtc, ok := rc.(*net.TCPConn); ok {
			if err := rtc.SetKeepAlivePeriod(time.Duration(c.Timeout) * time.Second); err != nil {
				return err
			}
		}
	}
	if c.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(c.Deadline) * time.Second)); err != nil {
			return err
		}
	}
	crc, err := c.wrapCipherConn(rc)
	if err != nil {
		var p *socks5.Reply
		if request.Atyp == socks5.ATYPIPv4 || request.Atyp == socks5.ATYPDomain {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(conn); err != nil {
			return err
		}
		return err
	}
	if _, err := crc.Write(rawaddr); err != nil {
		var p *socks5.Reply
		if request.Atyp == socks5.ATYPIPv4 || request.Atyp == socks5.ATYPDomain {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(conn); err != nil {
			return err
		}
		return err
	}
	a, addr, port, err := socks5.ParseAddress(rc.LocalAddr().String())
	if err != nil {
		var p *socks5.Reply
		if request.Atyp == socks5.ATYPIPv4 || request.Atyp == socks5.ATYPDomain {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(conn); err != nil {
			return err
		}
		return err
	}
	if err := socks5.NewReply(socks5.RepSuccess, a, addr, port).WriteTo(conn); err != nil {
		return err
	}
	go func() {
		// The first packet from server of shadowsocks contain the iv too
		iv := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(crc, iv); err != nil {
			log.Println(err)
			return
		}
		_, _ = io.Copy(conn, crc)
	}()
	_, _ = io.Copy(crc, conn)
	return nil
}

func (c *SSClient) handleHTTP(conn *net.TCPConn) error {
	defer conn.Close()
	if c.Timeout != 0 {
		if err := conn.SetKeepAlivePeriod(time.Duration(c.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if c.Deadline != 0 {
		if err := conn.SetDeadline(time.Now().Add(time.Duration(c.Deadline) * time.Second)); err != nil {
			return err
		}
	}

	b := make([]byte, 0, 1024)
	for {
		var b1 [1024]byte
		n, err := conn.Read(b1[:])
		if err != nil {
			return err
		}
		b = append(b, b1[:n]...)
		if bytes.Contains(b, []byte{0x0d, 0x0a, 0x0d, 0x0a}) {
			break
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
	} else {
		var err error
		addr, err = GetAddressFromURL(aoru)
		if err != nil {
			return err
		}
	}

	if c.HTTPMiddleman != nil {
		if handled, err := c.HTTPMiddleman.HandleHTTPProxy(method, aoru, b, conn); err != nil || handled {
			return err
		}
	}

	a, h, p, err := socks5.ParseAddress(addr)
	if err != nil {
		return err
	}
	rawaddr := make([]byte, 0)
	rawaddr = append(rawaddr, a)
	rawaddr = append(rawaddr, h...)
	rawaddr = append(rawaddr, p...)

	rc, err := c.Dial.Dial("tcp", c.Server)
	if err != nil {
		return err
	}
	defer rc.Close()
	if c.Timeout != 0 {
		if rtc, ok := rc.(*net.TCPConn); ok {
			if err := rtc.SetKeepAlivePeriod(time.Duration(c.Timeout) * time.Second); err != nil {
				return err
			}
		}
	}
	if c.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(c.Deadline) * time.Second)); err != nil {
			return err
		}
	}
	crc, err := c.wrapCipherConn(rc)
	if err != nil {
		return err
	}
	if _, err := crc.Write(rawaddr); err != nil {
		return err
	}
	if method == "CONNECT" {
		_, err := conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			return err
		}
	}
	if method != "CONNECT" {
		if _, err := crc.Write(b); err != nil {
			return err
		}
	}

	go func() {
		// The first packet from server of shadowsocks contain the iv too
		iv := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(crc, iv); err != nil {
			log.Println(err)
			return
		}
		_, _ = io.Copy(conn, crc)
	}()
	_, _ = io.Copy(crc, conn)
	return nil
}

func (c *SSClient) wrapCipherConn(conn net.Conn) (*CipherConn, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	if _, err := conn.Write(iv); err != nil {
		return nil, err
	}
	return NewCipherConn(conn, MakeSSKey(c.Password), iv)
}
