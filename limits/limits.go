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

// +build !windows

package limits

import (
	"runtime"
	"syscall"
)

func Raise() error {
	var l syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &l); err != nil {
		return err
	}
	if runtime.GOOS == "darwin" && l.Cur < 10240 {
		l.Cur = 10240
	}
	if runtime.GOOS != "darwin" && l.Cur < 60000 {
		if l.Max < 60000 {
			l.Max = 60000 // with CAP_SYS_RESOURCE capability
		}
		l.Cur = l.Max
	}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &l); err != nil {
		return err
	}
	return nil
}
