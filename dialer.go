package brook

import "net"

// Dialer is a common interface for dialing
type Dialer interface {
	Dial(network, server string) (net.Conn, error)
}

// DefaultDial is the default dialer which dial with tcp network
type DefaultDial struct {
}

// Dial a remote server
func (d *DefaultDial) Dial(network, server string) (net.Conn, error) {
	ta, err := net.ResolveTCPAddr(network, server)
	if err != nil {
		return nil, err
	}
	return net.DialTCP(network, nil, ta)
}
