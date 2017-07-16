package brook

import (
	"io"
	"net"

	"github.com/txthinking/pac/blackwhite"
	"github.com/txthinking/socks5"
)

// WhiteSocks5Middleman is a Socks5Middleman who only handle domain in white list
type WhiteSocks5Middleman struct {
	Timeout  int
	Deadline int
	Dial     Dialer
}

// NewWhiteSocks5Middleman returns a WhiteSocks5Middleman which can used to handle http proxy request
func NewWhiteSocks5Middleman(timeout, deadline int, dial Dialer) *WhiteSocks5Middleman {
	blackwhite.InitWhiteList()
	if dial == nil {
		dial = &DefaultDial{}
	}
	return &WhiteSocks5Middleman{
		Timeout:  timeout,
		Deadline: deadline,
		Dial:     dial,
	}
}

// HandleSocks5Proxy handle http proxy request, if the domain is in the white list
func (w *WhiteSocks5Middleman) HandleSocks5Proxy(request *socks5.Request, c *net.TCPConn) (handled bool, err error) {
	addr := request.Address()
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false, err
	}
	if !blackwhite.IsWhite(host) {
		return false, nil
	}

	rc, err := request.ConnectWithDial(c, w.Dial)
	if err != nil {
		return false, err
	}
	defer rc.Close()
	go func() {
		_, _ = io.Copy(rc, c)
	}()
	_, _ = io.Copy(c, rc)
	return true, nil
}
