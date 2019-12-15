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
	"log"
	"net"

	"github.com/txthinking/brook/sysproxy"
	"github.com/txthinking/gotun2socks"
	"github.com/txthinking/gotun2socks/tun"
)

// VPN.
type VPN struct {
	Client             *Client
	Tunnel             *Tunnel
	Tun                *gotun2socks.Tun2Socks
	ServerIP           string
	TunGateway         string
	OriginalDNSServers []string
}

// NewVPN.
func NewVPN(addr, server, password, dns string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int, tunDevice, tunIP, tunGateway, tunMask string) (*VPN, error) {
	ds, err := sysproxy.GetDNSServers()
	if err != nil {
		return nil, err
	}

	h, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if h != "127.0.0.1" {
		return nil, errors.New("Must listen on 127.0.0.1")
	}
	h, p, err := net.SplitHostPort(server)
	if err != nil {
		return nil, err
	}
	l, err := net.LookupIP(h)
	if err != nil {
		return nil, err
	}
	if len(l) == 0 {
		return nil, errors.New("Can not find server IP")
	}
	s := l[0].String()
	server = net.JoinHostPort(s, p)

	c, err := NewClient(addr, "127.0.0.1", server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return nil, err
	}
	dnsserver := net.JoinHostPort(dns, "53")
	tl, err := NewTunnel("127.0.0.1:53", dnsserver, server, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return nil, err
	}
	f, err := tun.OpenTunDevice(tunDevice, tunIP, tunGateway, tunMask, []string{dns})
	if err != nil {
		return nil, err
	}
	t := gotun2socks.New(f, addr, []string{dns}, false, true)
	return &VPN{
		Client:             c,
		Tunnel:             tl,
		Tun:                t,
		ServerIP:           s,
		TunGateway:         tunGateway,
		OriginalDNSServers: ds,
	}, nil
}

// ListenAndServe starts to run VPN.
func (v *VPN) ListenAndServe() error {
	if err := sysproxy.SetDNSServers([]string{"127.0.0.1"}); err != nil {
		return err
	}
	if err := v.AddRoutes(); err != nil {
		return err
	}

	errch := make(chan error)
	go func() {
		errch <- v.Client.ListenAndServe()
	}()
	go func() {
		errch <- v.Tunnel.ListenAndServe()
	}()
	go func() {
		v.Tun.Run()
	}()
	return <-errch
}

// Shutdown stops VPN.
func (v *VPN) Shutdown() error {
	if err := sysproxy.SetDNSServers(v.OriginalDNSServers); err != nil {
		log.Println(err)
	}
	if err := v.DeleteRoutes(); err != nil {
		log.Println(err)
	}
	if v.Client != nil {
		if err := v.Client.Shutdown(); err != nil {
			log.Println(err)
		}
	}
	if v.Tunnel != nil {
		if err := v.Tunnel.Shutdown(); err != nil {
			log.Println(err)
		}
	}
	if v.Tun != nil {
		// v.Tun.Stop()
	}
	return nil
}
