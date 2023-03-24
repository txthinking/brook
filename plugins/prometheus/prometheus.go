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

package prometheus

import (
	"context"
	"net"
	"net/http"

	pm "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/txthinking/brook"
)

type Prometheus struct {
	Server *http.Server
	Tags   map[string]string
}

func NewPrometheus(addr, path string, tags map[string]string) *Prometheus {
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.Handler())
	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return &Prometheus{Server: s, Tags: tags}
}

func (p *Prometheus) ListenAndServe() error {
	return p.Server.ListenAndServe()
}

func (p *Prometheus) Shutdown() error {
	return p.Server.Shutdown(context.Background())
}

func (p *Prometheus) TouchBrook() {
	tags := []string{}
	for k, _ := range p.Tags {
		tags = append(tags, k)
	}

	dstc := pm.NewCounterVec(
		pm.CounterOpts{
			Name: "dst_counter",
			Help: "Number of dst in total",
		},
		append([]string{"network", "from", "dst"}, tags...),
	)
	pm.MustRegister(dstc)
	f := brook.ServerGate
	brook.ServerGate = func(ex brook.Exchanger) (brook.Exchanger, error) {
		from := ex.Src()
		h, _, err := net.SplitHostPort(from)
		if err == nil {
			from = h
		}
		lb := pm.Labels{"network": ex.Network(), "from": from, "dst": ex.Dst()}
		for k, v := range p.Tags {
			lb[k] = v
		}
		dstc.With(lb).Inc()
		return f(ex)
	}
	f1 := brook.ClientGate
	brook.ClientGate = func(ex brook.Exchanger) (brook.Exchanger, error) {
		from := ex.Src()
		h, _, err := net.SplitHostPort(from)
		if err == nil {
			from = h
		}
		lb := pm.Labels{"network": ex.Network(), "from": from, "dst": ex.Dst()}
		for k, v := range p.Tags {
			lb[k] = v
		}
		dstc.With(lb).Inc()
		return f1(ex)
	}
}
