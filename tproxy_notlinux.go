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

// +build !linux

package brook

import (
	"errors"
	"net"

	cache "github.com/patrickmn/go-cache"
)

// Tproxy.
type Tproxy struct {
	TCPAddr       *net.TCPAddr
	UDPAddr       *net.UDPAddr
	RemoteTCPAddr *net.TCPAddr
	RemoteUDPAddr *net.UDPAddr
	Password      []byte
	TCPListen     *net.TCPListener
	UDPConn       *net.UDPConn
	Cache         *cache.Cache
	TCPDeadline   int
	TCPTimeout    int
	UDPDeadline   int
}

// NewTproxy.
func NewTproxy(addr, remote, password string, tcpTimeout, tcpDeadline, udpDeadline int) (*Tproxy, error) {
	return nil, errors.New("Only support Linux")
}

func (s *Tproxy) RunAutoScripts() error {
	return nil
}

func (s *Tproxy) ClearAutoScripts() error {
	return nil
}

// Run server.
func (s *Tproxy) ListenAndServe() error {
	return nil
}

// Shutdown server.
func (s *Tproxy) Shutdown() error {
	return nil
}
