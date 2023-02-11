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

package tproxy

import (
	"errors"
	"net"

	cache "github.com/patrickmn/go-cache"
)

type Tproxy struct {
	CIDR4      []*net.IPNet
	CIDR6      []*net.IPNet
	GeoIP      []string
	Cache      *cache.Cache
	TCPTimeout int
	UDPTimeout int
}

func NewTproxy(cidr4List, cidr6List string, geoIP []string, tcptimeout, udptimeout int) (*Tproxy, error) {
	return nil, errors.New("Only support Linux")
}

func (p *Tproxy) TouchBrook() {
}
