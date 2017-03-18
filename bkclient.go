package brook

import (
	"crypto/aes"
	"crypto/rand"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/ant"
	"github.com/txthinking/brook/music"
	"github.com/txthinking/socks5"
)

// BKClient is the client of brook protocol
type BKClient struct {
	Address  string
	Password string
	Server   string
	Timeout  int
	Deadline int
	Music    music.Music
	Dial     Dialer
	Listen   net.Listener
}

// NewBKClient returns a new BKClient, dial can be set to nil
func NewBKClient(address, server, password string, timeout, deadline int, m string, dial Dialer) (*BKClient, error) {
	msc, err := music.GetMusic(m)
	if err != nil {
		return nil, err
	}
	if dial == nil {
		dial = &DefaultDial{}
	}
	c := &BKClient{
		Address:  address,
		Server:   server,
		Password: password,
		Timeout:  timeout,
		Deadline: deadline,
		Music:    msc,
		Dial:     dial,
	}
	return c, nil
}

// ListenAndServe will let client start to listen and serve
func (c *BKClient) ListenAndServe() error {
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
func (c *BKClient) Shutdown() error {
	if c.Listen == nil {
		return nil
	}
	return c.Listen.Close()
}

func (c *BKClient) handle(conn *net.TCPConn) error {
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
	rawaddr := make([]byte, 0)
	rawaddr = append(rawaddr, request.Atyp)
	rawaddr = append(rawaddr, request.DstAddr...)
	rawaddr = append(rawaddr, request.DstPort...)

	rc, err := c.Dial.Dial("tcp", c.Server)
	if err != nil {
		if err := socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x09, 0x0a}).WriteTo(conn); err != nil {
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
	if c.Music != nil {
		if _, err := rc.Write(c.Music.GetSong()); err != nil {
			return err
		}
	}
	crc, err := c.wrapCipherConn(rc)
	if err != nil {
		if err := socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x09, 0x0a}).WriteTo(conn); err != nil {
			return err
		}
		return err
	}
	if _, err := crc.Write(rawaddr); err != nil {
		if err := socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x09, 0x0a}).WriteTo(conn); err != nil {
			return err
		}
		return err
	}
	a, addr, port := socks5.ParseAddress(crc.LocalAddr())
	if err := socks5.NewReply(socks5.RepSuccess, a, addr, port).WriteTo(conn); err != nil {
		return err
	}
	go func() {
		_, _ = io.Copy(conn, crc)
	}()
	_, _ = io.Copy(crc, conn)
	return nil
}

func (c *BKClient) wrapCipherConn(conn net.Conn) (*CipherConn, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	if _, err := conn.Write(iv); err != nil {
		return nil, err
	}
	return NewCipherConn(conn, []byte(ant.MD5(c.Password)), iv)
}
