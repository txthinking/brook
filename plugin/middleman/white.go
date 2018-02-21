package middleman

import (
	"io"
	"net"
	"time"

	"github.com/txthinking/ant"
	"github.com/txthinking/pac/blackwhite"
	"github.com/txthinking/socks5"
)

var Dial ant.Dialer = ant.DefaultDial

// WhiteSocks5 is a Socks5Middleman who only handles domain in white list
type WhiteSocks5 struct {
	Handle *socks5.DefaultHandle
}

// NewWhiteSocks5 returns a WhiteSocks5 which can be used to handle http proxy request
func NewWhiteSocks5() *WhiteSocks5 {
	return &WhiteSocks5{
		Handle: &socks5.DefaultHandle{},
	}
}

// TCPHandle handles tcp request
func (w *WhiteSocks5) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) (bool, error) {
	if r.Cmd == socks5.CmdConnect {
		h, _, err := net.SplitHostPort(r.Address())
		if err != nil {
			return false, err
		}
		if !blackwhite.IsWhite(h) {
			return false, nil
		}
		if err := w.Handle.TCPHandle(s, c, r); err != nil {
			return true, err
		}
		return true, nil
	}
	if r.Cmd == socks5.CmdUDP {
		return false, nil
	}
	return false, socks5.ErrUnsupportCmd
}

// UDPHandle handles udp packet
func (w *WhiteSocks5) UDPHandle(s *socks5.Server, ca *net.UDPAddr, d *socks5.Datagram) (bool, error) {
	h, _, err := net.SplitHostPort(d.Address())
	if err != nil {
		return false, err
	}
	if !blackwhite.IsWhite(h) {
		return false, nil
	}
	if err := w.Handle.UDPHandle(s, ca, d); err != nil {
		return true, err
	}
	return true, nil
}

// WhiteHTTP is a HTTPMiddleman who only handles domain in white list
type WhiteHTTP struct {
	Timeout  int
	Deadline int
}

// NewWhiteHTTP returns a WhiteHTTP which can used to handle http proxy request
func NewWhiteHTTP(timeout, deadline int) *WhiteHTTP {
	return &WhiteHTTP{
		Timeout:  timeout,
		Deadline: deadline,
	}
}

// Handle handles http proxy request, if the domain is in the white list
func (w *WhiteHTTP) Handle(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error) {
	h, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false, err
	}
	if !blackwhite.IsWhite(h) {
		return false, nil
	}

	tmp, err := Dial.Dial("tcp", addr)
	if err != nil {
		return true, err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if w.Timeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(w.Timeout) * time.Second); err != nil {
			return true, err
		}
	}
	if w.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(w.Deadline) * time.Second)); err != nil {
			return true, err
		}
	}
	if method == "CONNECT" {
		_, err := conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			return true, err
		}
	}
	if method != "CONNECT" {
		if _, err := rc.Write(request); err != nil {
			return true, err
		}
	}
	go func() {
		_, _ = io.Copy(rc, conn)
	}()
	_, _ = io.Copy(conn, rc)
	return true, nil
}
