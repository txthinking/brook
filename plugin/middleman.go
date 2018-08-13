package plugin

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

// HTTPMiddleman is a middleman who can intercept and handle request
type HTTPMiddleman interface {
	// Addr is the absoluteURI, RFC 2396.
	// Request is the http header, don't guarantee it is complete, but contains the host line
	// Has not written anything to conn.
	// Handle does not need to close conn.
	// If handled is true or err is not nil that means the request has been handled
	Handle(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error)
}
