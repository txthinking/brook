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
	"net"
	"net/http"
	"strings"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/phuslu/iploc"
)

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

func ListHasIP(c4, c6 []*net.IPNet, i net.IP, c *cache.Cache, geo []string) bool {
	if c4 == nil && c6 == nil && len(geo) == 0 {
		return false
	}
	if c != nil {
		any, ok := c.Get(i.String())
		if ok {
			return any.(bool)
		}
	}
	if len(geo) != 0 {
		b := iploc.Country(i)
		if b != nil {
			bs := string(b)
			for _, v := range geo {
				if v == bs {
					if c != nil {
						c.Set(i.String(), true, 24*time.Hour)
					}
					return true
				}
			}
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
			Log(&Error{"when": "net.ParseCIDR", "cidr": v, "error": err.Error()})
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
					h, p, err := net.SplitHostPort(addr)
					if err != nil {
						return nil, err
					}
					s, err := Resolve6(h)
					if err == nil {
						c, err := net.Dial(network, net.JoinHostPort(s, p))
						if err == nil {
							return c, nil
						}
					}
					s, err = Resolve4(h)
					if err == nil {
						c, err := net.Dial(network, net.JoinHostPort(s, p))
						if err == nil {
							return c, nil
						}
					}
					return nil, errors.New("Can not fetch " + addr)
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

func ReadData(url string) ([]byte, error) {
	var data []byte
	var err error
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		c := &http.Client{
			Timeout: 9 * time.Second,
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					h, p, err := net.SplitHostPort(addr)
					if err != nil {
						return nil, err
					}
					s, err := Resolve6(h)
					if err == nil {
						c, err := net.Dial(network, net.JoinHostPort(s, p))
						if err == nil {
							return c, nil
						}
					}
					s, err = Resolve4(h)
					if err == nil {
						c, err := net.Dial(network, net.JoinHostPort(s, p))
						if err == nil {
							return c, nil
						}
					}
					return nil, errors.New("Can not fetch " + addr)
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
	return data, nil
}
