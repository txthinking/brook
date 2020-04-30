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

// +build !linux,!darwin,!windows

package brook

import (
	"errors"

	"github.com/txthinking/gotun2socks"
	"github.com/txthinking/runnergroup"
)

// Tun.
type Tun struct {
	Client             *Client
	Tunnel             *Tunnel
	Tun                *gotun2socks.Tun2Socks
	ServerIP           string
	TunGateway         string
	OriginalDNSServers []string
	RunnerGroup        *runnergroup.RunnerGroup
}

// NewTun.
func NewTun(addr, server, password, dns string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int, tunDevice, tunIP, tunGateway, tunMask string) (*Tun, error) {
	return nil, errors.New("Unsupport your OS yet")
}

// ListenAndServe starts to run Tun.
func (v *Tun) ListenAndServe() error {
	return nil
}

// Shutdown stops Tun.
func (v *Tun) Shutdown() error {
	return nil
}
