package brook

import (
	"net"

	"github.com/txthinking/socks5"
)

// Socks5Middleman is a middleman who can intercept and handle request
type Socks5Middleman interface {
	// TCPRequestHandle does not need to close conn,
	// if return true that means the request has been handled
	TCPHandle(*socks5.Server, *net.TCPConn, *socks5.Request) (bool, error)

	// UDPPacketHandle handles udp packet.
	// If return true that means the request has been handled.
	UDPHandle(*socks5.Server, *net.UDPAddr, *socks5.Datagram) (bool, error)
}
