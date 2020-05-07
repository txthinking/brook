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

// +build linux darwin windows

package brook

import (
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
	"github.com/eycorsican/go-tun2socks/tun"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/brook/sysproxy"
	"github.com/txthinking/runnergroup"
)

// Tun.
type Tun struct {
	Client             *Client
	Tunnel             *Tunnel
	ListenIP           string
	ServerIP           string
	TunGateway         string
	OriginalDNSServers []string
	RunnerGroup        *runnergroup.RunnerGroup
	LwipWriter         io.Writer
	Fd                 io.ReadWriteCloser
	LetBrookDoAllForMe bool
}

// NewTun.
func NewTun(addr, server, password, dns string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int, tunDevice, tunIP, tunGateway, tunMask string) (*Tun, error) {
	h, p, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if h == "" {
		return nil, errors.New("Listen address must contain IP")
	}
	listenIP := h
	listenPort, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		return nil, err
	}
	h1, p1, err := net.SplitHostPort(server)
	if err != nil {
		return nil, err
	}
	l, err := net.LookupIP(h1)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, errors.New("Can not find server IP")
	}
	serverIP := l[0].String()
	server = net.JoinHostPort(serverIP, p1)

	c, err := NewClient(addr, listenIP, server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return nil, err
	}
	d, err := NewTunnel(net.JoinHostPort(listenIP, "53"), net.JoinHostPort(dns, "53"), server, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return nil, err
	}

	fd, err := tun.OpenTunDevice(tunDevice, tunIP, tunGateway, tunMask, []string{dns}, false)
	if err != nil {
		return nil, err
	}
	lw := core.NewLWIPStack().(io.Writer)
	core.RegisterTCPConnHandler(socks.NewTCPHandler(listenIP, uint16(listenPort)))
	core.RegisterUDPConnHandler(socks.NewUDPHandler(listenIP, uint16(listenPort), time.Duration(udpSessionTime)*time.Second))
	core.RegisterOutputFn(func(data []byte) (int, error) {
		return fd.Write(data)
	})

	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	return &Tun{
		Client:      c,
		Tunnel:      d,
		ListenIP:    listenIP,
		ServerIP:    serverIP,
		TunGateway:  tunGateway,
		RunnerGroup: runnergroup.New(),
		LwipWriter:  lw,
		Fd:          fd,
	}, nil
}

// ListenAndServe starts to run Tun.
func (v *Tun) ListenAndServe() error {
	if v.LetBrookDoAllForMe {
		ds, err := sysproxy.GetDNSServers()
		if err != nil {
			return err
		}
		v.OriginalDNSServers = ds

		if err := sysproxy.SetDNSServers([]string{v.ListenIP}); err != nil {
			return err
		}
		if err := v.AddRoutes(); err != nil {
			return err
		}
	}

	v.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			return v.Client.ListenAndServe()
		},
		Stop: func() error {
			return v.Client.Shutdown()
		},
	})
	v.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			return v.Tunnel.ListenAndServe()
		},
		Stop: func() error {
			return v.Tunnel.Shutdown()
		},
	})
	v.RunnerGroup.Add(&runnergroup.Runner{
		Start: func() error {
			_, err := io.CopyBuffer(v.LwipWriter, v.Fd, make([]byte, 1500))
			return err
		},
		Stop: func() error {
			return v.Fd.Close()
		},
	})
	return v.RunnerGroup.Wait()
}

// Shutdown stops Tun.
func (v *Tun) Shutdown() error {
	if v.LetBrookDoAllForMe {
		if err := sysproxy.SetDNSServers(v.OriginalDNSServers); err != nil {
			log.Println(err)
		}
		if err := v.DeleteRoutes(); err != nil {
			log.Println(err)
		}
	}
	return v.RunnerGroup.Done()
}
