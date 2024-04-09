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

package logger

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/krolaw/dhcp4"
	"github.com/miekg/dns"
	"github.com/txthinking/brook"
)

type Logger struct {
	F    *os.File
	File string
	Lock *sync.Mutex
	Tags map[string]string
}

func NewLogger(tags map[string]string, file string) (*Logger, error) {
	if file == "console" {
		return &Logger{Tags: tags}, nil
	}
	f, err := os.OpenFile(file, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{F: f, File: file, Tags: tags, Lock: &sync.Mutex{}}, nil
}

func (p *Logger) Close() error {
	if p.Lock == nil {
		return nil
	}
	p.Lock.Lock()
	defer p.Lock.Unlock()
	return p.F.Close()
}

func (p *Logger) Reset() error {
	if p.Lock == nil {
		return nil
	}
	p.Lock.Lock()
	defer p.Lock.Unlock()
	err := p.F.Close()
	if err != nil {
		return err
	}
	p.F, err = os.OpenFile(p.File, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (p *Logger) TouchBrook() {
	go p.WatchReset()
	brook.Log = func(err error) {
		if _, ok := err.(brook.Error); !ok {
			err = brook.Error{"error": err.Error()}
		}
		err.(brook.Error)["time"] = time.Now().Format(time.RFC3339)
		for k, v := range p.Tags {
			err.(brook.Error)[k] = v
		}
		if p.Lock == nil {
			fmt.Println(err)
			return
		}
		p.Lock.Lock()
		_, err = p.F.Write([]byte(err.Error() + "\n"))
		p.Lock.Unlock()
		if err != nil {
			log.Println(err)
		}
	}
	f := brook.ServerGate
	brook.ServerGate = func(ex brook.Exchanger) (brook.Exchanger, error) {
		brook.Log(brook.Error{"network": ex.Network(), "from": ex.Src(), "dst": strings.ToLower(ex.Dst())})
		return f(ex)
	}
	f1 := brook.ClientGate
	brook.ClientGate = func(ex brook.Exchanger) (brook.Exchanger, error) {
		brook.Log(brook.Error{"network": ex.Network(), "from": ex.Src(), "dst": strings.ToLower(ex.Dst())})
		return f1(ex)
	}
	f2 := brook.DNSGate
	brook.DNSGate = func(addr *net.UDPAddr, m *dns.Msg, l1 *net.UDPConn) (bool, error) {
		brook.Log(brook.Error{"from": addr.String(), "dns": dns.Type(m.Question[0].Qtype).String(), "domain": strings.ToLower(m.Question[0].Name[0 : len(m.Question[0].Name)-1])})
		return f2(addr, m, l1)
	}
	f4 := brook.DOHGate
	brook.DOHGate = func(m *dns.Msg, w http.ResponseWriter, r *http.Request) (done bool, err error) {
		s := r.RemoteAddr
		if r.Header.Get("X-Forwarded-For") != "" {
			s = r.Header.Get("X-Forwarded-For")
		}
		brook.Log(brook.Error{"from": s, "dns": dns.Type(m.Question[0].Qtype).String(), "domain": strings.ToLower(m.Question[0].Name[0 : len(m.Question[0].Name)-1])})
		return f4(m, w, r)
	}
	f3 := brook.DHCPServerGate
	brook.DHCPServerGate = func(inmt string, in dhcp4.Packet, outmt string, ip net.IP, err error) {
		e := brook.Error{"in": inmt, "client": in.CHAddr().String(), "out": outmt}
		if ip != nil {
			e["ip"] = ip.String()
		}
		if err != nil {
			e["error"] = err.Error()
		}
		brook.Log(e)
		f3(inmt, in, outmt, ip, err)
	}
}
