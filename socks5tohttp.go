package brook

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

type Socks5ToHTTP struct {
	Address       string
	Socks5Address string
	Dial          proxy.Dialer
	Timeout       int
	Deadline      int
	Listen        net.Listener
	HTTPMiddleman HTTPMiddleman
}

func NewSocks5ToHTTP(addr, socks5addr string, timeout, deadline int) (*Socks5ToHTTP, error) {
	dial, err := proxy.SOCKS5("tcp", socks5addr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	return &Socks5ToHTTP{
		Address:       addr,
		Socks5Address: socks5addr,
		Dial:          dial,
		Timeout:       timeout,
		Deadline:      deadline,
	}, nil
}

func (s *Socks5ToHTTP) ListenAndServe(h HTTPMiddleman) error {
	ta, err := net.ResolveTCPAddr("tcp", s.Address)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", ta)
	if err != nil {
		return err
	}
	defer l.Close()
	s.Listen = l
	s.HTTPMiddleman = h
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		go func(conn *net.TCPConn) {
			if err := s.handle(conn); err != nil {
				log.Println(err)
			}
		}(conn)
	}
}

func (s *Socks5ToHTTP) Shutdown() error {
	if s.Listen == nil {
		return nil
	}
	return s.Listen.Close()
}

func (s *Socks5ToHTTP) handle(conn *net.TCPConn) error {
	defer conn.Close()
	if s.Timeout != 0 {
		if err := conn.SetKeepAlivePeriod(time.Duration(s.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if s.Deadline != 0 {
		if err := conn.SetDeadline(time.Now().Add(time.Duration(s.Deadline) * time.Second)); err != nil {
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
	method, address := string(bb[0]), string(bb[1])
	var addr string
	if method == "CONNECT" {
		addr = address
	} else {
		var err error
		addr, err = GetAddressFromURL(address)
		if err != nil {
			return err
		}
	}

	if s.HTTPMiddleman != nil {
		if handled, err := s.HTTPMiddleman.HandleHTTPProxy(method, addr, b, conn); err != nil || handled {
			return err
		}
	}

	rc, err := s.Dial.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer rc.Close()
	if s.Timeout != 0 {
		if rtc, ok := rc.(*net.TCPConn); ok {
			if err := rtc.SetKeepAlivePeriod(time.Duration(s.Timeout) * time.Second); err != nil {
				return err
			}
		}
	}
	if s.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.Deadline) * time.Second)); err != nil {
			return err
		}
	}
	if method == "CONNECT" {
		_, err := conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			return err
		}
	}
	if method != "CONNECT" {
		if _, err := rc.Write(b); err != nil {
			return err
		}
	}
	go func() {
		_, _ = io.Copy(rc, conn)
	}()
	_, _ = io.Copy(conn, rc)
	return nil
}
