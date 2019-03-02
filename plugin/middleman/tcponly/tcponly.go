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

package tcponly

import (
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
)

var Dial x.Dialer = x.DefaultDial

// TCPOnly is a middleman.
type TCPOnly struct {
	Timeout      int
	Deadline     int
	Socks5Handle socks5.Handler
}

// NewTCPOnly returns a TCPOnly.
func NewTCPOnly(timeout, deadline int) *TCPOnly {
	return &TCPOnly{
		Timeout:      timeout,
		Deadline:     deadline,
		Socks5Handle: &socks5.DefaultHandle{},
	}
}

// TCPHandle handles tcp request.
func (t *TCPOnly) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) (bool, error) {
	if r.Cmd == socks5.CmdUDP {
		if err := t.Socks5Handle.TCPHandle(s, c, r); err != nil {
			return true, err
		}
		return true, nil
	}
	return false, nil
}

// UDPHandle handles udp packet.
func (t *TCPOnly) UDPHandle(s *socks5.Server, ca *net.UDPAddr, d *socks5.Datagram) (bool, error) {
	if err := t.Socks5Handle.UDPHandle(s, ca, d); err != nil {
		return true, err
	}
	return true, nil
}

// Handle handles http proxy request, if the domain is in the white list.
func (t *TCPOnly) Handle(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error) {
	tmp, err := Dial.Dial("tcp", addr)
	if err != nil {
		return true, err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if t.Timeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(t.Timeout) * time.Second); err != nil {
			return true, err
		}
	}
	if t.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(t.Deadline) * time.Second)); err != nil {
			return true, err
		}
	}
	if method == "CONNECT" {
		_, err := conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			return true, err
		}
	}
	if method != "CONNECT" {
		if _, err := rc.Write(request); err != nil {
			return true, err
		}
	}
	go func() {
		var bf [1024 * 2]byte
		for {
			if t.Deadline != 0 {
				if err := conn.SetDeadline(time.Now().Add(time.Duration(t.Deadline) * time.Second)); err != nil {
					return
				}
			}
			i, err := conn.Read(bf[:])
			if err != nil {
				return
			}
			if _, err := rc.Write(bf[0:i]); err != nil {
				return
			}
		}
	}()
	var bf [1024 * 2]byte
	for {
		if t.Deadline != 0 {
			if err := rc.SetDeadline(time.Now().Add(time.Duration(t.Deadline) * time.Second)); err != nil {
				return true, nil
			}
		}
		i, err := rc.Read(bf[:])
		if err != nil {
			return true, nil
		}
		if _, err := conn.Write(bf[0:i]); err != nil {
			return true, nil
		}
	}
	return true, nil
}
