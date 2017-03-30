package brook

import "net"

// Dialer is a common interface for dialing
type Dialer interface {
	Dial(network, addr string) (net.Conn, error)
}

// DefaultDial is the default dialer which dial with tcp network
type DefaultDial struct {
}

// Dial a remote address
func (d *DefaultDial) Dial(network, addr string) (net.Conn, error) {
	ta, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	return net.DialTCP(network, nil, ta)
}
