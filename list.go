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
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
)

func BlockAddress(address string, ds map[string]byte, c4, c6 []*net.IPNet, c *cache.Cache) bool {
	if ds == nil && c4 == nil && c6 == nil {
		return false
	}
	h, _, err := net.SplitHostPort(address)
	if err != nil {
		log.Println(err)
		return false
	}
	i := net.ParseIP(h)
	if i == nil {
		return ListHasDomain(ds, h, c)
	}
	return ListHasIP(c4, c6, i, c)
}

func ListHasDomain(ds map[string]byte, domain string, c *cache.Cache) bool {
	if ds == nil {
		return false
	}
	if c != nil {
		any, ok := c.Get(domain)
		if ok {
			return any.(bool)
		}
	}
	ss := strings.Split(domain, ".")
	var s1 string
	for i := len(ss) - 1; i >= 0; i-- {
		if s1 == "" {
			s1 = ss[i]
		} else {
			s1 = ss[i] + "." + s1
		}
		if _, ok := ds[s1]; ok {
			if c != nil {
				c.Set(domain, true, 24*time.Hour)
			}
			return true
		}
	}
	if c != nil {
		c.Set(domain, false, 24*time.Hour)
	}
	return false
}

func ListHasIP(c4, c6 []*net.IPNet, i net.IP, c *cache.Cache) bool {
	if c4 == nil && c6 == nil {
		return false
	}
	if c != nil {
		any, ok := c.Get(i.String())
		if ok {
			return any.(bool)
		}
	}
	if i.To4() != nil {
		if c4 == nil {
			return false
		}
		ii := i.To4()
		for _, v := range c4 {
			if v.Contains(ii) {
				if c != nil {
					c.Set(i.String(), true, 24*time.Hour)
				}
				return true
			}
		}
	}
	if i.To4() == nil {
		if c6 == nil {
			return false
		}
		ii := i.To16()
		for _, v := range c6 {
			if v.Contains(ii) {
				if c != nil {
					c.Set(i.String(), true, 24*time.Hour)
				}
				return true
			}
		}
	}
	if c != nil {
		c.Set(i.String(), false, 24*time.Hour)
	}
	return false
}

func ReadDomainList(url string) (map[string]byte, error) {
	ds := make(map[string]byte)
	ss, err := ReadList(url)
	if err != nil {
		return nil, err
	}
	for _, v := range ss {
		ds[v] = 0
	}
	return ds, nil
}

func ReadCIDRList(url string) ([]*net.IPNet, error) {
	c := make([]*net.IPNet, 0)
	l, err := ReadList(url)
	if err != nil {
		return nil, err
	}
	for _, v := range l {
		_, in, err := net.ParseCIDR(v)
		if err != nil {
			log.Println("net.ParseCIDR", v, err)
			continue
		}
		c = append(c, in)
	}
	return c, nil
}

func ReadList(url string) ([]string, error) {
	var data []byte
	var err error
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		c := &http.Client{
			Timeout: 9 * time.Second,
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					r := &net.Resolver{
						Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
							c, err := net.Dial(network, "8.8.8.8:53")
							if err != nil {
								c, err = net.Dial(network, "[2001:4860:4860::8888]:53")
							}
							return c, err
						},
					}
					h, p, err := net.SplitHostPort(addr)
					if err != nil {
						return nil, err
					}
					l, err := r.LookupIP(ctx, "ip4", h)
					if err == nil && len(l) > 0 {
						c, err := net.Dial(network, net.JoinHostPort(l[0].String(), p))
						if err == nil {
							return c, nil
						}
					}
					l, err = r.LookupIP(ctx, "ip6", h)
					if err == nil && len(l) > 0 {
						c, err := net.Dial(network, net.JoinHostPort(l[0].String(), p))
						if err == nil {
							return c, nil
						}
					}
					return nil, errors.New("Can not fetch " + addr + ", maybe dns or network error")
				},
			},
		}
		r, err := c.Get(url)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		data, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		data, err = ioutil.ReadFile(url)
		if err != nil {
			return nil, err
		}
	}
	data = bytes.TrimSpace(data)
	data = bytes.Replace(data, []byte{0x20}, []byte{}, -1)
	data = bytes.Replace(data, []byte{0x0d, 0x0a}, []byte{0x0a}, -1)
	ss := strings.Split(string(data), "\n")
	return ss, nil
}
