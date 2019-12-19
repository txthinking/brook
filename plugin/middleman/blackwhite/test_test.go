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

package blackwhite

import (
	"log"
	"testing"
)

func TestTest(t *testing.T) {
	bw, err := NewBlackWhite("white", "https://blackwhite.txthinking.com/white.list", "https://blackwhite.txthinking.com/white_cidr.list", "8.8.8.8:53", "119.29.29.29:53", 0, 6)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(bw.Has("www.baidu.com"))
	log.Println(bw.Has("39.156.69.79"))
	log.Println(bw.Has("220.181.38.148"))
	log.Println(bw.Has("119.29.29.29"))
	log.Println(bw.Has("61.135.186.217"))
}
