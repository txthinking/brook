package brook

import (
	"crypto/aes"
	"crypto/rand"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/ant"
)

// S5Client is the client of socks5 encrypt protocol
type S5Client struct {
	Address  string
	Password string
	Server   string
	Timeout  int
	Deadline int
	Dial     Dialer
	Listen   net.Listener
}

// NewS5Client returns a new S5Client
func NewS5Client(address, server, password string, timeout, deadline int, dial Dialer) *S5Client {
	if dial == nil {
		dial = &DefaultDial{}
	}
	c := &S5Client{
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
func (c *S5Client) ListenAndServe() error {
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

// Shutdown used to stop the client
func (c *S5Client) Shutdown() error {
	if c.Listen == nil {
		return nil
	}
	return c.Listen.Close()
}

func (c *S5Client) handle(conn *net.TCPConn) error {
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
	go func() {
		_, _ = io.Copy(conn, crc)
	}()
	_, _ = io.Copy(crc, conn)
	return nil
}

func (c *S5Client) wrapCipherConn(conn net.Conn) (*CipherConn, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	if _, err := conn.Write(iv); err != nil {
		return nil, err
	}
	return NewCipherConn(conn, []byte(ant.MD5(c.Password)), iv)
}
