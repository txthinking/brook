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

package pprof

import (
	"context"
	"net/http"
	_ "net/http/pprof"
)

type Pprof struct {
	s *http.Server
}

func NewPprof(addr string) (*Pprof, error) {
	s := &http.Server{
		Addr: addr,
	}
	return &Pprof{
		s: s,
	}, nil
}

func (p *Pprof) ListenAndServe() error {
	return p.s.ListenAndServe()
}

func (p *Pprof) Shutdown() error {
	return p.s.Shutdown(context.Background())
}
