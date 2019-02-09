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

import "crypto/md5"

// MakeSSKey used to make shadowsocks aes-256-cfb key.
func MakeSSKey(password string) []byte {
	h := md5.New()
	h.Write([]byte(password))
	tmp := h.Sum(nil) // 16 len

	h = md5.New()
	h.Write(tmp)
	h.Write([]byte(password))
	tmp1 := h.Sum(nil) // 16 len

	return append(tmp, tmp1...) // 32 len
}
