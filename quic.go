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
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

func QUICDialUDP(src, dst, addr, host string, tc *tls.Config, idleTime int) (net.Conn, error) {
	var rc *net.UDPConn
	var err error
	if src == "" || dst == "" {
		rc, err = ListenUDP("udp", nil)
	}
	if src != "" && dst != "" {
		rc, err = NATListenUDP("udp", src, dst)
	}
	if err != nil {
		return nil, err
	}
	raddr, err := Resolve("udp", addr)
	if err != nil {
		rc.Close()
		return nil, err
	}
	rc1, err := quic.Dial(rc, raddr.(*net.UDPAddr), host, tc, &quic.Config{MaxIdleTimeout: time.Duration(idleTime) * time.Second, EnableDatagrams: true})
	if err != nil {
		rc.Close()
		return nil, err
	}
	return &QUICConn{
		UDPConn: rc,
		Conn:    rc1,
		LAddr:   rc1.LocalAddr(),
		RAddr:   rc1.RemoteAddr(),
	}, nil
}

func QUICDialTCP(src, dst, addr, host string, tc *tls.Config, idleTime int) (net.Conn, error) {
	var rc *net.UDPConn
	var err error
	if src == "" || dst == "" {
		rc, err = ListenUDP("udp", nil)
	}
	if src != "" && dst != "" {
		rc, err = NATListenUDP("udp", src, dst)
	}
	if err != nil {
		return nil, err
	}
	raddr, err := Resolve("udp", addr)
	if err != nil {
		rc.Close()
		return nil, err
	}
	rc1, err := quic.Dial(rc, raddr.(*net.UDPAddr), host, tc, &quic.Config{MaxIdleTimeout: time.Duration(idleTime) * time.Second})
	if err != nil {
		rc.Close()
		return nil, err
	}
	s, err := rc1.OpenStreamSync(context.Background())
	if err != nil {
		rc1.CloseWithError(0, err.Error())
		rc.Close()
		return nil, err
	}
	return &QUICConn{
		UDPConn: rc,
		Conn:    rc1,
		Stream:  s,
		LAddr: &net.TCPAddr{
			IP:   rc1.LocalAddr().(*net.UDPAddr).IP,
			Port: rc1.LocalAddr().(*net.UDPAddr).Port,
			Zone: rc1.LocalAddr().(*net.UDPAddr).Zone,
		},
		RAddr: &net.TCPAddr{
			IP:   rc1.RemoteAddr().(*net.UDPAddr).IP,
			Port: rc1.RemoteAddr().(*net.UDPAddr).Port,
			Zone: rc1.RemoteAddr().(*net.UDPAddr).Zone,
		},
	}, nil
}

type QUICConn struct {
	UDPConn *net.UDPConn
	Conn    quic.Connection
	Stream  quic.Stream
	LAddr   net.Addr
	RAddr   net.Addr
}

func (c *QUICConn) Read(b []byte) (int, error) {
	if c.Stream != nil {
		return c.Stream.Read(b)
	}
	b1, err := c.Conn.ReceiveMessage()
	if err != nil {
		return 0, err
	}
	i := copy(b, b1)
	return i, nil
}

func (c *QUICConn) Write(b []byte) (int, error) {
	if c.Stream != nil {
		return c.Stream.Write(b)
	}
	if err := c.Conn.SendMessage(b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *QUICConn) Close() error {
	if c.Stream != nil {
		c.Stream.CancelRead(0)
		c.Stream.Close()
	}
	if c.Conn != nil {
		c.Conn.CloseWithError(0, "close")
	}
	if c.UDPConn != nil {
		c.UDPConn.Close()
	}
	return nil
}

func (c *QUICConn) LocalAddr() net.Addr {
	return c.LAddr
}

func (c *QUICConn) RemoteAddr() net.Addr {
	return c.RAddr
}

func (c *QUICConn) SetDeadline(t time.Time) error {
	if c.Stream != nil {
		return c.Stream.SetDeadline(t)
	}
	// prefer quic.Config.MaxIdleTimeout
	return nil
}

func (c *QUICConn) SetReadDeadline(t time.Time) error {
	if c.Stream != nil {
		return c.Stream.SetReadDeadline(t)
	}
	return nil
}

func (c *QUICConn) SetWriteDeadline(t time.Time) error {
	if c.Stream != nil {
		return c.Stream.SetWriteDeadline(t)
	}
	return nil
}
