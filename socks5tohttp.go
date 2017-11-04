package brook

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/ant"
	"golang.org/x/net/proxy"
)

type Socks5ToHTTP struct {
	Addr          *net.TCPAddr
	Socks5Address string
	Dial          proxy.Dialer
	Timeout       int
	Deadline      int // Not refreshed
	Listen        *net.TCPListener
	Middleman     HTTPMiddleman
}

func NewSocks5ToHTTP(addr, socks5addr string, timeout, deadline int) (*Socks5ToHTTP, error) {
	dial, err := proxy.SOCKS5("tcp", socks5addr, nil, &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	ta, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Socks5ToHTTP{
		Addr:          ta,
		Socks5Address: socks5addr,
		Dial:          dial,
		Timeout:       timeout,
		Deadline:      deadline,
	}, nil
}

func (s *Socks5ToHTTP) ListenAndServe(h HTTPMiddleman) error {
	l, err := net.ListenTCP("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer l.Close()
	s.Listen = l
	s.Middleman = h
	for {
		c, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		go func(c *net.TCPConn) {
			defer c.Close()
			if s.Timeout != 0 {
				if err := c.SetKeepAlivePeriod(time.Duration(s.Timeout) * time.Second); err != nil {
					log.Println(err)
					return
				}
			}
			if s.Deadline != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(s.Deadline) * time.Second)); err != nil {
					log.Println(err)
					return
				}
			}
			if err := s.Handle(c); err != nil {
				log.Println(err)
				return
			}
		}(c)
	}
}

func (s *Socks5ToHTTP) Handle(c *net.TCPConn) error {
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
	method, address := string(bb[0]), string(bb[1])
	var addr string
	if method == "CONNECT" {
		addr = address
	}
	if method != "CONNECT" {
		var err error
		addr, err = ant.GetAddressFromURL(address)
		if err != nil {
			return err
		}
	}

	if s.Middleman != nil {
		if done, err := s.Middleman.Handle(method, addr, b, c); err != nil || done {
			return err
		}
	}

	tmp, err := s.Dial.Dial("tcp", addr)
	if err != nil {
		return err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if s.Timeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(s.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if s.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.Deadline) * time.Second)); err != nil {
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
		if _, err := rc.Write(b); err != nil {
			return err
		}
	}
	go func() {
		_, _ = io.Copy(rc, c)
	}()
	_, _ = io.Copy(c, rc)
	return nil
}

func (s *Socks5ToHTTP) Shutdown() error {
	if s.Listen == nil {
		return nil
	}
	return s.Listen.Close()
}
