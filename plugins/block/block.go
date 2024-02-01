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

package block

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook"
)

type Block struct {
	DomainList string
	CIDR4List  string
	CIDR6List  string
	Domain     map[string]byte
	CIDR4      []*net.IPNet
	CIDR6      []*net.IPNet
	GeoIP      []string
	Cache      *cache.Cache
	Lock       *sync.RWMutex
	Duration   int
	Done       chan byte
}

func NewBlock(domainList, cidr4List, cidr6List string, geoIP []string, update int) (*Block, error) {
	var err error
	var ds map[string]byte
	if domainList != "" {
		ds, err = brook.ReadDomainList(domainList)
		if err != nil {
			return nil, err
		}
	}
	var c4 []*net.IPNet
	if cidr4List != "" {
		c4, err = brook.ReadCIDRList(cidr4List)
		if err != nil {
			return nil, err
		}
	}
	var c6 []*net.IPNet
	if cidr6List != "" {
		c6, err = brook.ReadCIDRList(cidr6List)
		if err != nil {
			return nil, err
		}
	}
	b := &Block{
		DomainList: domainList,
		CIDR4List:  cidr4List,
		CIDR6List:  cidr6List,
		Domain:     ds,
		CIDR4:      c4,
		CIDR6:      c6,
		GeoIP:      geoIP,
		Cache:      cache.New(cache.NoExpiration, cache.NoExpiration),
		Duration:   update,
	}
	if update != 0 {
		b.Lock = &sync.RWMutex{}
		b.Done = make(chan byte)
	}
	return b, nil
}

func (bk *Block) Update() {
	ticker := time.NewTicker(time.Duration(bk.Duration) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-bk.Done:
			return
		case <-ticker.C:
			var err error
			var ds map[string]byte
			if bk.DomainList != "" {
				ds, err = brook.ReadDomainList(bk.DomainList)
				if err != nil {
					log.Println("ReadDomainList", bk.DomainList, err)
					break
				}
			}
			var c4 []*net.IPNet
			if bk.CIDR4List != "" {
				c4, err = brook.ReadCIDRList(bk.CIDR4List)
				if err != nil {
					log.Println("ReadCIDRList", bk.CIDR4List, err)
					break
				}
			}
			var c6 []*net.IPNet
			if bk.CIDR6List != "" {
				c6, err = brook.ReadCIDRList(bk.CIDR6List)
				if err != nil {
					log.Println("ReadCIDRList", bk.CIDR6List, err)
					break
				}
			}
			bk.Lock.Lock()
			bk.Domain = ds
			bk.CIDR4 = c4
			bk.CIDR6 = c6
			if bk.Cache != nil {
				bk.Cache.Flush()
			}
			bk.Lock.Unlock()
		}
	}
}

func (bk *Block) Stop() {
	select {
	case <-bk.Done:
	default:
		close(bk.Done)
	}
}

func (bk *Block) TouchBrook() {
	f := brook.Resolve
	brook.Resolve = func(network string, addr string) (net.Addr, error) {
		var ds map[string]byte
		var c4 []*net.IPNet
		var c6 []*net.IPNet
		if bk.Lock != nil {
			bk.Lock.RLock()
		}
		ds = bk.Domain
		c4 = bk.CIDR4
		c6 = bk.CIDR6
		if bk.Lock != nil {
			bk.Lock.RUnlock()
		}
		h, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		var a net.Addr
		ip := net.ParseIP(h)
		if ip == nil {
			if brook.ListHasDomain(ds, strings.ToLower(h), bk.Cache) {
				return nil, errors.New("block " + addr)
			}
			a, err = f(network, addr)
			if err != nil {
				return nil, err
			}
			v, ok := a.(*net.TCPAddr)
			if ok {
				ip = v.IP
			}
			if !ok {
				ip = a.(*net.UDPAddr).IP
			}
		}
		if brook.ListHasIP(c4, c6, ip, bk.Cache, bk.GeoIP) {
			return nil, errors.New("block " + addr)
		}
		if a != nil {
			return a, nil
		}
		return f(network, addr)
	}
}
