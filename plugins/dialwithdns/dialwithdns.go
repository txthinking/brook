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

package dialwithdns

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/txthinking/brook"
	"github.com/txthinking/socks5"
)

type DialWithDNS struct {
	DNSClient *brook.DNSClient
	DOHClient *brook.DOHClient
	Prefer    string
}

func NewDialWithDNS(dns, prefer string) (*DialWithDNS, error) {
	if prefer != "A" && prefer != "AAAA" {
		return nil, errors.New("Invalid prefer")
	}
	if !strings.HasPrefix(dns, "https://") {
		return &DialWithDNS{DNSClient: &brook.DNSClient{Server: dns}, Prefer: prefer}, nil
	}
	dc, err := brook.NewDOHClient(dns)
	if err != nil {
		return nil, err
	}
	return &DialWithDNS{DOHClient: dc, Prefer: prefer}, nil
}

func (p *DialWithDNS) IP(domain string) (net.IP, error) {
	if p.Prefer == "A" {
		if p.DNSClient != nil {
			ip, err := p.DNSClient.A(domain)
			if err != nil {
				return nil, err
			}
			if ip != nil {
				return ip, nil
			}
			ip, err = p.DNSClient.AAAA(domain)
			if err != nil {
				return nil, err
			}
			if ip != nil {
				return ip, nil
			}
			return nil, errors.New("Can not resolve " + domain)
		}
		if p.DOHClient != nil {
			ip, err := p.DOHClient.A(domain)
			if err != nil {
				return nil, err
			}
			if ip != nil {
				return ip, nil
			}
			ip, err = p.DOHClient.AAAA(domain)
			if err != nil {
				return nil, err
			}
			if ip != nil {
				return ip, nil
			}
			return nil, errors.New("Can not resolve " + domain)
		}
	}
	if p.DNSClient != nil {
		ip, err := p.DNSClient.AAAA(domain)
		if err != nil {
			return nil, err
		}
		if ip != nil {
			return ip, nil
		}
		ip, err = p.DNSClient.A(domain)
		if err != nil {
			return nil, err
		}
		if ip != nil {
			return ip, nil
		}
		return nil, errors.New("Can not resolve " + domain)
	}
	if p.DOHClient != nil {
		ip, err := p.DOHClient.AAAA(domain)
		if err != nil {
			return nil, err
		}
		if ip != nil {
			return ip, nil
		}
		ip, err = p.DOHClient.A(domain)
		if err != nil {
			return nil, err
		}
		if ip != nil {
			return ip, nil
		}
		return nil, errors.New("Can not resolve " + domain)
	}
	return nil, errors.New("Can not resolve " + domain)
}

func (p *DialWithDNS) TouchBrook() {
	brook.Resolve = func(network string, addr string) (net.Addr, error) {
		h, p1, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		port, err := strconv.Atoi(p1)
		if err != nil {
			return nil, err
		}
		ip := net.ParseIP(h)
		if ip == nil {
			ip, err = p.IP(h)
			if err != nil {
				return nil, err
			}
		}
		if network == "tcp" {
			return &net.TCPAddr{IP: ip, Port: port}, nil
		}
		return &net.UDPAddr{IP: ip, Port: port}, nil
	}
	socks5.Resolve = brook.Resolve
}
