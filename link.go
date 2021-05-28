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
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/mdp/qrterminal"
	"github.com/txthinking/encrypt"
)

// kind: brookserver/brookwsserver/brookwssserver/socks5server
func Link(kind, s, username, password string) string {
	v := url.Values{}
	v.Set(kind, s)
	v.Set("username", username)
	v.Set("password", password)
	s = fmt.Sprintf("brook://%s?%s", kind, v.Encode())
	return s
}

func QR(kind, s, username, password string) {
	qrterminal.GenerateHalfBlock(Link(kind, s, username, password), qrterminal.L, os.Stdout)
}

// kind: server/wsserver/wssserver/socks5
func ParseLink(link string) (kind, s, username, password string, err error) {
	var u *url.URL
	u, err = url.Parse(link)
	if err != nil {
		return
	}
	kind = u.Host
	s = u.Query().Get(kind)
	username = u.Query().Get("username")
	password = u.Query().Get("password")
	return
}

// kind: server/wsserver/wssserver/socks5
func ParseLinkOld(link string) (kind, server, username, password string, err error) {
	if !strings.HasPrefix(link, "brook://") {
		err = errors.New("Invalid brook link")
		return
	}
	s := link[8:]
	s, err = encrypt.URIUnescape(s)
	if err != nil {
		return
	}
	l := strings.Split(s, " ")
	if len(l) == 1 {
		kind = "socks5"
		server = l[0]
		return
	}
	if len(l) == 3 {
		kind = "socks5"
		server = l[0]
		username = l[1]
		password = l[2]
		return
	}
	if len(l) == 2 {
		kind = "server"
		if strings.HasPrefix(l[0], "ws://") {
			kind = "wsserver"
		}
		if strings.HasPrefix(l[0], "wss://") {
			kind = "wssserver"
		}
		server = l[0]
		password = l[1]
		return
	}
	err = errors.New("Invalid brook link")
	return
}
