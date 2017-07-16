package brook

import "net"

// HTTPMiddleman is a middleman who can intercept and handle request
type HTTPMiddleman interface {
	// HandleHTTPProxy does not need to close conn,
	// Addr is the absoluteURI, RFC 2396
	// if handled is true or err is not nil that means the request has been handled
	HandleHTTPProxy(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error)
}
