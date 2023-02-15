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
	"fmt"
	"net/url"
)

func Link(kind, server string, v url.Values) string {
	v.Set(kind, server)
	return fmt.Sprintf("brook://%s?%s", kind, v.Encode())
}

func ParseLink(link string) (kind, server string, v url.Values, err error) {
	var u *url.URL
	u, err = url.Parse(link)
	if err != nil {
		return
	}
	kind = u.Host
	server = u.Query().Get(kind)
	v = u.Query()
	return
}
