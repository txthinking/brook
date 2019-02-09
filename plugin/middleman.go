// Copyright (c) 2016-present Cloud <cloud@txthinking.com>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of version 3 of the GNU General Public
// License as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package plugin

import (
	"net"

	"github.com/txthinking/socks5"
)

// Socks5Middleman is a middleman who can intercept and handle request.
type Socks5Middleman interface {
	// TCPHandle does not need to close conn.
	// If return true or error that means the request has been handled.
	TCPHandle(*socks5.Server, *net.TCPConn, *socks5.Request) (bool, error)

	// UDPHandle handles udp packet.
	// If return true or error that means the request has been handled.
	UDPHandle(*socks5.Server, *net.UDPAddr, *socks5.Datagram) (bool, error)
}

// HTTPMiddleman is a middleman who can intercept and handle request
type HTTPMiddleman interface {
	// Addr is the absoluteURI, RFC 2396.
	// Request is the http header, don't guarantee it is complete, but contains the host line.
	// Has not written anything to conn.
	// Handle does not need to close conn.
	// If return true or error that means the request has been handled.
	Handle(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error)
}
