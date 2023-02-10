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
	"log"
	"net"
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
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
	Update     int
	Done       chan byte
}

func NewBlock(domainList, cidr4List, cidr6List string, update int64, geoIP []string) (*Block, error) {
	var err error
	var ds map[string]byte
	if blockDomainList != "" {
		ds, err = ReadDomainList(blockDomainList)
		if err != nil {
			return nil, err
		}
	}
	var c4 []*net.IPNet
	if blockCIDR4List != "" {
		c4, err = ReadCIDRList(blockCIDR4List)
		if err != nil {
			return nil, err
		}
	}
	var c6 []*net.IPNet
	if blockCIDR6List != "" {
		c6, err = ReadCIDRList(blockCIDR6List)
		if err != nil {
			return nil, err
		}
	}
	return &Block{
		DomainList: domainList,
		CIDR4List:  cidr4List,
		CIDR6List:  cidr6List,
		Domain:     ds,
		CIDR4:      c4,
		CIDR6:      c6,
		GeoIP:      geoIP,
		Cache:      cache.New(cache.NoExpiration, cache.NoExpiration),
		Lock:       &sync.Mutex{},
		Update:     update,
		Done:       make(chan byte),
	}, nil
}

func (block *Block) Update() {
	ticker := time.NewTicker(time.Duration(block.Update) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-block.Done:
			return
		case <-ticker.C:
			var ds map[string]byte
			if block.DomainList != "" {
				ds, err = ReadDomainList(block.DomainList)
				if err != nil {
					log.Println("ReadDomainList", block.DomainList, err)
					break
				}
			}
			var c4 []*net.IPNet
			if block.CIDR4List != "" {
				c4, err = ReadCIDRList(block.CIDR4List)
				if err != nil {
					log.Println("ReadCIDRList", block.CIDR4List, err)
					break
				}
			}
			var c6 []*net.IPNet
			if block.CIDR6List != "" {
				c6, err = ReadCIDRList(block.CIDR6List)
				if err != nil {
					log.Println("ReadCIDRList", block.CIDR6List, err)
					break
				}
			}
			block.Lock.Lock()
			s.Domain = ds
			s.CIDR4 = c4
			s.CIDR6 = c6
			if s.Cache != nil {
				s.Cache.Flush()
			}
			block.Lock.Unlock()
		}
	}
}

func (block *Block) Stop() {
	select {
	case <-block.Done:
	default:
		close(block.Done)
	}
}

func (block *Block) Stop() {
}

func BlockAddress(address string, ds map[string]byte, c4, c6 []*net.IPNet, c *cache.Cache, geo []string) bool {
	if ds == nil && c4 == nil && c6 == nil && len(geo) == 0 {
		return false
	}
	h, _, err := net.SplitHostPort(address)
	if err != nil {
		return true
	}
	i := net.ParseIP(h)
	if i == nil {
		return ListHasDomain(ds, h, c)
	}
	return ListHasIP(c4, c6, i, c, geo)
}
