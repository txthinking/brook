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
	"net"

	"github.com/txthinking/socks5"
)

func ErrorReply(r *socks5.Request, c *net.TCPConn, e error) error {
	var p *socks5.Reply
	if r.Atyp == socks5.ATYPIPv4 || r.Atyp == socks5.ATYPDomain {
		p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, net.IPv4zero, []byte{0x00, 0x00})
	} else {
		p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv6, net.IPv6zero, []byte{0x00, 0x00})
	}
	if err := p.WriteTo(c); err != nil {
		return err
	}
	return e
}
