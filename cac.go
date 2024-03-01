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
	"strings"
)

// TODO: make me more compatible with the shell environment.
// https://github.com/TxThinkingInc/CAC
func CAC(file string) ([]string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	l0 := []string{}
	l1 := strings.Split(string(b), "\n")
	for _, v := range l1 {
		v = strings.TrimSpace(v)
		if strings.HasSuffix(v, "\\") {
			v = v[0 : len(v)-1]
		}
		if strings.HasPrefix(v, "#") || v == "" {
			continue
		}
		l0 = append(l0, v)
	}
	if len(l0) == 0 {
		return l0, nil
	}
	s := strings.Join(l0, " ")
	l0 = []string{}
	l1 = strings.Fields(s)
	for _, v := range l1 {
		if len(v) > 1 && strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			l0 = append(l0, v[1:len(v)-1])
			continue
		}
		if len(v) > 1 && strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"") {
			l0 = append(l0, v[1:len(v)-1])
			continue
		}
		l0 = append(l0, v)
	}
	return l0, nil
}
