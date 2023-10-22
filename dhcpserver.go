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
	"encoding/json"
	"errors"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/krolaw/dhcp4"
)

type DHCPServer struct {
	Listen   net.PacketConn
	ServerIP net.IP
	Start    net.IP
	Count    int
	Leases   map[int]string
	Options  dhcp4.Options
	Cache    string
}

func NewDHCPServer(iface, serverip, start, mask string, count int, gateway string, dnsserver []string, cache string) (*DHCPServer, error) {
	if net.ParseIP(serverip).To4() == nil || net.ParseIP(start).To4() == nil || net.ParseIP(mask).To4() == nil || net.ParseIP(gateway).To4() == nil {
		return nil, errors.New("Invalid v4 IP")
	}
	dnsserverips := make([]net.IP, 0)
	for _, v := range dnsserver {
		if net.ParseIP(v).To4() == nil {
			return nil, errors.New("Invalid v4 IP")
		}
		dnsserverips = append(dnsserverips, net.ParseIP(v).To4())
	}
	if cache == "" {
		s, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cache = filepath.Join(s, ".brook.dhcpserver")
	}
	b, err := os.ReadFile(cache)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	m := make(map[int]string)
	if err == nil {
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}
	}
	l, err := DHCPListen(iface)
	if err != nil {
		return nil, err
	}
	return &DHCPServer{
		Listen:   l,
		ServerIP: net.ParseIP(serverip).To4(),
		Start:    net.ParseIP(start).To4(),
		Count:    count,
		Leases:   m,
		Options: dhcp4.Options{
			dhcp4.OptionSubnetMask:       []byte(net.ParseIP(mask).To4()),
			dhcp4.OptionRouter:           []byte(net.ParseIP(gateway).To4()),
			dhcp4.OptionDomainNameServer: dhcp4.JoinIPs(dnsserverips),
		},
		Cache: cache,
	}, nil
}

func (h *DHCPServer) ListenAndServe() error {
	return dhcp4.Serve(h.Listen, h)
}

func (h *DHCPServer) Shutdown() error {
	b, err := json.Marshal(h.Leases)
	if err != nil {
		Log(err)
	}
	if err == nil {
		if err := os.WriteFile(h.Cache, b, 0644); err != nil {
			Log(err)
		}
	}
	return h.Listen.Close()
}

var DHCPServerGate func(inmt string, in dhcp4.Packet, outmt string, ip net.IP, err error) = func(inmt string, in dhcp4.Packet, outmt string, ip net.IP, err error) {
}

func (h *DHCPServer) ServeDHCP(p dhcp4.Packet, msgType dhcp4.MessageType, options dhcp4.Options) (d dhcp4.Packet) {
	switch msgType {
	case dhcp4.Discover:
		for i, s := range h.Leases {
			if s == p.CHAddr().String() {
				DHCPServerGate(msgType.String(), p, dhcp4.Offer.String(), dhcp4.IPAdd(h.Start, i), nil)
				return dhcp4.ReplyPacket(p, dhcp4.Offer, h.ServerIP, dhcp4.IPAdd(h.Start, i), 7*24*time.Hour, h.Options.SelectOrderOrAll(options[dhcp4.OptionParameterRequestList]))
			}
		}
		for i := 0; i < h.Count; i++ {
			_, ok := h.Leases[i]
			if !ok {
				DHCPServerGate(msgType.String(), p, dhcp4.Offer.String(), dhcp4.IPAdd(h.Start, i), nil)
				return dhcp4.ReplyPacket(p, dhcp4.Offer, h.ServerIP, dhcp4.IPAdd(h.Start, i), 7*24*time.Hour, h.Options.SelectOrderOrAll(options[dhcp4.OptionParameterRequestList]))
			}
		}
		Log(errors.New("DHCP server is full"))
		DHCPServerGate(msgType.String(), p, "", nil, errors.New("DHCP server is full"))
		return nil
	case dhcp4.Request:
		if server, ok := options[dhcp4.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.ServerIP) {
			return nil
		}
		reqIP := net.IP(options[dhcp4.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(p.CIAddr())
		}
		if len(reqIP) == 4 && !reqIP.Equal(net.IPv4zero) {
			i := dhcp4.IPRange(h.Start, reqIP) - 1
			if i >= 0 && i < h.Count {
				s, ok := h.Leases[i]
				if !ok || s == p.CHAddr().String() {
					h.Leases[i] = p.CHAddr().String()
					DHCPServerGate(msgType.String(), p, dhcp4.ACK.String(), reqIP, nil)
					return dhcp4.ReplyPacket(p, dhcp4.ACK, h.ServerIP, reqIP, 7*24*time.Hour, h.Options.SelectOrderOrAll(options[dhcp4.OptionParameterRequestList]))
				}
			}
		}
		DHCPServerGate(msgType.String(), p, dhcp4.NAK.String(), reqIP, nil)
		return dhcp4.ReplyPacket(p, dhcp4.NAK, h.ServerIP, nil, 0, nil)
	case dhcp4.Release, dhcp4.Decline:
		for i, s := range h.Leases {
			if s == p.CHAddr().String() {
				delete(h.Leases, i)
			}
		}
		DHCPServerGate(msgType.String(), p, "", nil, nil)
		return nil
	}
	DHCPServerGate(msgType.String(), p, "", nil, nil)
	return nil
}
