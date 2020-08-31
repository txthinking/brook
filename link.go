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
	"os"

	"github.com/mdp/qrterminal"
	"github.com/txthinking/encrypt"
)

// Link
func Link(server, password, username string) string {
	s := server
	if username == "" && password != "" {
		s += " " + password
	}
	if username != "" && password != "" {
		s += " " + username + " " + password
	}
	s = "brook://" + encrypt.URIEscape(s)
	return s
}

// QR generate and print QR code.
func QR(server, password, username string) {
	s := server
	if username == "" && password != "" {
		s += " " + password
	}
	if username != "" && password != "" {
		s += " " + username + " " + password
	}
	s = "brook://" + encrypt.URIEscape(s)
	qrterminal.GenerateHalfBlock(s, qrterminal.L, os.Stdout)
}
