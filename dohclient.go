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
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/miekg/dns"
)

type DOHClient struct {
	Client *http.Client
	DOH    string
}

func NewDOHClient(doh string) (*DOHClient, error) {
	u, err := url.Parse(doh)
	if err != nil {
		return nil, err
	}
	a := u.Query().Get("address")
	if a == "" {
		return nil, errors.New("no address")
	}
	q := u.Query()
	q.Del("address")
	u.RawQuery = q.Encode()
	c := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return DialTCP("tcp", "", a)
			},
		},
	}
	return &DOHClient{
		Client: c,
		DOH:    u.String(),
	}, nil
}

func (c *DOHClient) Exchange(m *dns.Msg) (*dns.Msg, error) {
	b, err := m.Pack()
	if err != nil {
		return nil, err
	}
	hr, err := http.NewRequest("POST", c.DOH, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	hr.Header.Set("Accept", "application/dns-message")
	hr.Header.Set("Content-Type", "application/dns-message")
	res, err := c.Client.Do(hr)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	r := &dns.Msg{}
	if err := r.Unpack(b); err != nil {
		return nil, err
	}
	return r, nil
}

// if no AAAA, return nil
func (c *DOHClient) AAAA(domain string) (net.IP, error) {
	m := &dns.Msg{}
	m.SetQuestion(domain+".", dns.TypeAAAA)
	m, err := c.Exchange(m)
	if err != nil {
		return nil, err
	}
	for _, v := range m.Answer {
		if t, ok := v.(*dns.AAAA); ok {
			return t.AAAA, nil
		}
	}
	return nil, nil
}

// if no A, return nil
func (c *DOHClient) A(domain string) (net.IP, error) {
	m := &dns.Msg{}
	m.SetQuestion(domain+".", dns.TypeA)
	m, err := c.Exchange(m)
	if err != nil {
		return nil, err
	}
	for _, v := range m.Answer {
		if t, ok := v.(*dns.A); ok {
			return t.A, nil
		}
	}
	return nil, nil
}
