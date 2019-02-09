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

// RunTproxy used to start a tproxy.
func RunTproxy(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	c, err := NewTproxy(address, server, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return err
	}
	return c.ListenAndServe()
}
