package brook

import (
	"net"

	"github.com/txthinking/socks5"
)

// Socks5Middleman is a middleman who can intercept and handle request
type Socks5Middleman interface {
	// HandleSocks5Proxy does not need to close conn,
	// if handled is true or err is not nil that means the request has been handled
	HandleSocks5Proxy(request *socks5.Request, conn *net.TCPConn) (handled bool, err error)
}
