package brook

import (
	"net"

	"github.com/txthinking/socks5"
)

func ErrorReply(r *socks5.Request, c *net.TCPConn, e error) error {
	var p *socks5.Reply
	if r.Atyp == socks5.ATYPIPv4 || r.Atyp == socks5.ATYPDomain {
		p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, net.IPv4zero, []byte{0x00, 0x00})
	} else {
		p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv6, net.IPv6zero, []byte{0x00, 0x00})
	}
	if err := p.WriteTo(c); err != nil {
		return err
	}
	return e
}
