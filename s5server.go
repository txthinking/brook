package brook

import (
	"crypto/aes"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/ant"
	"github.com/txthinking/socks5"
)

// S5Server is the server of socks5 encrypt protocol
type S5Server struct {
	Address  string
	Password string
	Timeout  int
	Deadline int
	Listen   net.Listener
}

// NewS5Server returns a new S5Server
func NewS5Server(address, password string, timeout, deadline int) *S5Server {
	s := &S5Server{
		Address:  address,
		Password: password,
		Timeout:  timeout,
		Deadline: deadline,
	}
	return s
}

// ListenAndServe will let server start to listen and serve
func (s *S5Server) ListenAndServe() error {
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
func (s *S5Server) Shutdown() error {
	if s.Listen == nil {
		return nil
	}
	return s.Listen.Close()
}

func (s *S5Server) handle(c *net.TCPConn) error {
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
	s5s := socks5.NewClassicServer(cc)
	if err := s5s.Negotiate(); err != nil {
		return err
	}
	r, err := s5s.GetRequest()
	if err != nil {
		return err
	}
	if r.Cmd == socks5.CmdConnect { // this 'if' is not necessary when SupportedCommands only contain conect method
		rc, err := r.Connect(cc)
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
			_, _ = io.Copy(cc, rc)
		}()
		_, _ = io.Copy(rc, cc)
	}
	return nil
}

func (s *S5Server) wrapCipherConn(conn net.Conn) (*CipherConn, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(conn, iv); err != nil {
		return nil, err
	}
	return NewCipherConn(conn, []byte(ant.MD5(s.Password)), iv)
}
