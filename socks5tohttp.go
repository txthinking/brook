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
	"errors"
	"log"
	"net"
	"time"

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/x"
	"golang.org/x/net/proxy"
)

type Socks5ToHTTP struct {
	Addr           *net.TCPAddr
	Socks5Address  string
	Socks5Username string
	Socks5Password string
	Dial           proxy.Dialer
	Timeout        int
	Deadline       int
	Listen         *net.TCPListener
	HTTPMiddleman  plugin.HTTPMiddleman
}

func NewSocks5ToHTTP(addr, socks5addr, socks5username, socks5password string, timeout, deadline int) (*Socks5ToHTTP, error) {
	var auth *proxy.Auth
	if socks5username != "" || socks5password != "" {
		auth = &proxy.Auth{
			User:     socks5username,
			Password: socks5password,
		}
	}
	dial, err := proxy.SOCKS5("tcp", socks5addr, auth, &net.Dialer{
		Timeout:   time.Duration(deadline) * time.Second,
		KeepAlive: time.Duration(timeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}
	ta, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	return &Socks5ToHTTP{
		Addr:           ta,
		Socks5Address:  socks5addr,
		Socks5Username: socks5username,
		Socks5Password: socks5password,
		Dial:           dial,
		Timeout:        timeout,
		Deadline:       deadline,
	}, nil
}

// SetHTTPMiddleman sets httpmiddleman plugin.
func (s *Socks5ToHTTP) SetHTTPMiddleman(m plugin.HTTPMiddleman) {
	s.HTTPMiddleman = m
}

func (s *Socks5ToHTTP) ListenAndServe() error {
	l, err := net.ListenTCP("tcp", s.Addr)
	if err != nil {
		return err
	}
	defer l.Close()
	s.Listen = l
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
		addr, err = x.GetAddressFromURL(address)
		if err != nil {
			return err
		}
	}

	if s.HTTPMiddleman != nil {
		done, err := s.HTTPMiddleman.Handle(method, addr, b, c)
		if done {
			return err
		}
		if err != nil {
			return err
		}
	}

	if Debug {
		log.Println("Dial TCP", addr)
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
		var bf [1024 * 2]byte
		for {
			if s.Deadline != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.Deadline) * time.Second)); err != nil {
					return
				}
			}
			i, err := rc.Read(bf[:])
			if err != nil {
				return
			}
			if _, err := c.Write(bf[0:i]); err != nil {
				return
			}
		}
	}()
	var bf [1024 * 2]byte
	for {
		if s.Deadline != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.Deadline) * time.Second)); err != nil {
				return nil
			}
		}
		i, err := c.Read(bf[:])
		if err != nil {
			return nil
		}
		if _, err := rc.Write(bf[0:i]); err != nil {
			return nil
		}
	}
	return nil
}

func (s *Socks5ToHTTP) Shutdown() error {
	if s.Listen == nil {
		return nil
	}
	return s.Listen.Close()
}
