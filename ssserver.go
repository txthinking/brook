package brook

import (
	"crypto/aes"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

// SSServer is the server of shadowsocks protocol
type SSServer struct {
	Address  string
	Password string
	Timeout  int
	Deadline int
	Listen   net.Listener
}

// NewSSServer returns a new SSServer
func NewSSServer(address, password string, timeout, deadline int) *SSServer {
	s := &SSServer{
		Address:  address,
		Password: password,
		Timeout:  timeout,
		Deadline: deadline,
	}
	return s
}

// ListenAndServe will let server start to listen and serve
func (s *SSServer) ListenAndServe() error {
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

	for {
		c, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		go func(c *net.TCPConn) {
			if err := s.handle(c); err != nil {
				log.Println(err)
			}
		}(c)
	}
}

// Shutdown used to stop the server
func (s *SSServer) Shutdown() error {
	if s.Listen == nil {
		return nil
	}
	return s.Listen.Close()
}

func (s *SSServer) handle(c *net.TCPConn) error {
	defer c.Close()
	if s.Timeout != 0 {
		if err := c.SetKeepAlivePeriod(time.Duration(s.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if s.Deadline != 0 {
		if err := c.SetDeadline(time.Now().Add(time.Duration(s.Deadline) * time.Second)); err != nil {
			return err
		}
	}
	cc, err := s.wrapCipherConn(c)
	if err != nil {
		return err
	}
	// ss https://shadowsocks.org/en/spec/protocol.html Address Type | Destination Address | Destination Port defined same as SOCKS5
	bb := make([]byte, 1)
	if _, err := io.ReadFull(cc, bb); err != nil {
		return err
	}
	var addr []byte
	if bb[0] == socks5.ATYPIPv4 {
		addr = make([]byte, 4)
		if _, err := io.ReadFull(cc, addr); err != nil {
			return err
		}
	} else if bb[0] == socks5.ATYPIPv6 {
		addr = make([]byte, 16)
		if _, err := io.ReadFull(cc, addr); err != nil {
			return err
		}
	} else if bb[0] == socks5.ATYPDomain {
		dal := make([]byte, 1)
		if _, err := io.ReadFull(cc, dal); err != nil {
			return err
		}
		if dal[0] == 0 {
			return err
		}
		addr = make([]byte, int(dal[0]))
		if _, err := io.ReadFull(cc, addr); err != nil {
			return err
		}
		addr = append(dal, addr...)
	} else {
		return errors.New("Unknown address type")
	}
	port := make([]byte, 2)
	if _, err := io.ReadFull(cc, port); err != nil {
		return err
	}
	address := socks5.ToAddress(bb[0], addr, port)

	ta, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	rc, err := net.DialTCP("tcp", nil, ta)
	if err != nil {
		return err
	}
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
	go func() {
		// The first packet from server of shadowsocks contain the iv too
		iv := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			log.Println(err)
			return
		}
		if _, err := cc.Write(iv); err != nil {
			log.Println(err)
			return
		}
		_, _ = io.Copy(cc, rc)
	}()
	_, _ = io.Copy(rc, cc)
	return nil
}

func (s *SSServer) wrapCipherConn(conn net.Conn) (*CipherConn, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(conn, iv); err != nil {
		return nil, err
	}
	return NewCipherConn(conn, MakeSSKey(s.Password), iv)
}
