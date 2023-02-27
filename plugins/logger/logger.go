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
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/txthinking/brook"
)

type Logger struct {
	F    *os.File
	Lock *sync.Mutex
	Tags map[string]string
}

func NewLogger(tags map[string]string, file string) (*Logger, error) {
	if file == "console" {
		return &Logger{Tags: tags}, nil
	}
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &Logger{F: f, Tags: tags, Lock: &sync.Mutex{}}, nil
}

func (p *Logger) Close() error {
	if p.F == nil {
		return nil
	}
	p.Lock.Lock()
	defer p.Lock.Unlock()
	return p.F.Close()
}

func (p *Logger) TouchBrook() {
	brook.Log = func(err error) {
		if _, ok := err.(brook.Error); !ok {
			err = brook.Error{"error": err.Error()}
		}
		err.(brook.Error)["time"] = strconv.FormatInt(time.Now().Unix(), 10)
		for k, v := range p.Tags {
			err.(brook.Error)[k] = v
		}
		if p.F == nil {
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
		brook.Log(brook.Error{"network": ex.Network(), "from": ex.Src(), "dst": ex.Dst()})
		return f(ex)
	}
	f1 := brook.ClientGate
	brook.ClientGate = func(ex brook.Exchanger) (brook.Exchanger, error) {
		brook.Log(brook.Error{"network": ex.Network(), "from": ex.Src(), "dst": ex.Dst()})
		return f1(ex)
	}
}
