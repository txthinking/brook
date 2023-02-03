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

package tproxy

import (
	"net"
	"syscall"

	"golang.org/x/sys/unix"
)

func ListenTCP(network string, laddr *net.TCPAddr) (*net.TCPListener, error) {
	l, err := net.ListenTCP(network, laddr)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	f, err := l.File()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fd := int(f.Fd())
	if laddr.IP.To4() != nil {
		if err := syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
			return nil, err
		}
	} else {
		if err := syscall.SetsockoptInt(fd, syscall.SOL_IPV6, unix.IPV6_TRANSPARENT, 1); err != nil {
			return nil, err
		}
	}
	tmp, err := net.FileListener(f)
	if err != nil {
		return nil, err
	}
	return tmp.(*net.TCPListener), nil
}
