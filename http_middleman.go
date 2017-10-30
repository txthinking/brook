package brook

import "net"

// HTTPMiddleman is a middleman who can intercept and handle request
type HTTPMiddleman interface {
	// Addr is the absoluteURI, RFC 2396.
	// Request is the http header, don't guarantee it is complete, but must contains the host line
	// Has not written anything to conn.
	// Handle does not need to close conn.
	// Ff handled is true or err is not nil that means the request has been handled
	Handle(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error)
}
