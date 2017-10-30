package brook

import (
	"net"

	"github.com/txthinking/pac/blackwhite"
	"github.com/txthinking/socks5"
)

// WhiteSocks5Middleman is a Socks5Middleman who only handle domain in white list
type WhiteSocks5Middleman struct {
	Handle *socks5.DefaultHandle
}

// NewWhiteSocks5Middleman returns a WhiteSocks5Middleman which can used to handle http proxy request
func NewWhiteSocks5Middleman() *WhiteSocks5Middleman {
	return &WhiteSocks5Middleman{
		Handle: &socks5.DefaultHandle{},
	}
}

// TCPHandle handles tcp request
func (w *WhiteSocks5Middleman) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) (bool, error) {
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
func (w *WhiteSocks5Middleman) UDPHandle(s *socks5.Server, ca *net.UDPAddr, d *socks5.Datagram) (bool, error) {
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
