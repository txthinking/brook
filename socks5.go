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

package brook

import (
	"log"

	"github.com/txthinking/brook/limits"
	"github.com/txthinking/socks5"
)

// Socks5Server is raw socks5 server.
type Socks5Server struct {
	Server *socks5.Server
}

// NewSocks5Server returns a new Socks5Server.
func NewSocks5Server(addr, ip, userName, password string, tcpTimeout, udpTimeout int) (*Socks5Server, error) {
	s5, err := socks5.NewClassicServer(addr, ip, userName, password, tcpTimeout, udpTimeout)
	if err != nil {
		return nil, err
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	x := &Socks5Server{
		Server: s5,
	}
	return x, nil
}

// ListenAndServe will let client start to listen and serve.
func (x *Socks5Server) ListenAndServe() error {
	return x.Server.ListenAndServe(nil)
}

// Shutdown used to stop the client.
func (x *Socks5Server) Shutdown() error {
	return x.Server.Shutdown()
}
