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
	"testing"
)

func TestTest(t *testing.T) {
	log.Println(Resolve6("http3.ooo"))
	log.Println(Resolve4("http3.ooo"))
	log.Println(Resolve6("2001:4860:4860::8888"))
	log.Println(Resolve6("8.8.8.8"))
	log.Println(Resolve4("2001:4860:4860::8888"))
	log.Println(Resolve4("8.8.8.8"))
	log.Println(Resolve6("a"))
	log.Println(Resolve4("a"))
	log.Println(Resolve6("mtalk.google.com"))
	log.Println(Resolve4("mtalk.google.com"))
}
