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
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/socks5"
)

type Socks5Server struct {
	Server *socks5.Server
}

func NewSocks5Server(addr, ip, userName, password string, tcpTimeout, udpTimeout int) (*Socks5Server, error) {
	if err := limits.Raise(); err != nil {
		Log(Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	s5, err := socks5.NewClassicServer(addr, ip, userName, password, tcpTimeout, udpTimeout)
	if err != nil {
		return nil, err
	}
	x := &Socks5Server{
		Server: s5,
	}
	return x, nil
}

func (x *Socks5Server) ListenAndServe() error {
	return x.Server.ListenAndServe(nil)
}

func (x *Socks5Server) Shutdown() error {
	return x.Server.Shutdown()
}
