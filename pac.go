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
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/js"
)

type PAC struct {
	Addr       string
	File       string
	Proxy      string
	DomainURL  string
	DomainData []byte
	HTTPServer *http.Server
	Body       []byte
}

func NewPAC(addr, file, proxy, domainURL string) *PAC {
	p := &PAC{
		Addr:      addr,
		File:      file,
		Proxy:     proxy,
		DomainURL: domainURL,
	}
	mux := http.NewServeMux()
	mux.Handle("/", p)
	p.HTTPServer = &http.Server{
		Addr:           p.Addr,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}
	return p
}

func (p *PAC) MakeBody() (io.Reader, error) {
	var err error
	l := make([]string, 0)
	if p.DomainURL != "" {
		l, err = ReadList(p.DomainURL)
		if err != nil {
			return nil, err
		}
	}
	t := template.New("pac")
	t, err = t.Parse(tpl)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	if err := t.Execute(b, map[string]interface{}{
		"proxy":   p.Proxy,
		"domains": l,
	}); err != nil {
		return nil, err
	}
	b1 := &bytes.Buffer{}
	m := minify.New()
	m.AddFunc("application/javascript", js.Minify)
	if err := m.Minify("application/javascript", b1, b); err != nil {
		return nil, err
	}
	return b1, nil
}

func (p *PAC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
	w.Write(p.Body)
}

func (p *PAC) ListenAndServe() error {
	r, err := p.MakeBody()
	if err != nil {
		return err
	}
	p.Body, err = ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return p.HTTPServer.ListenAndServe()
}

func (p *PAC) Shutdown() error {
	return p.HTTPServer.Shutdown(context.Background())
}

func (p *PAC) WriteToFile() error {
	r, err := p.MakeBody()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(p.File, b, 0644); err != nil {
		return err
	}
	return nil
}

func (p *PAC) WriteToStdout() error {
	r, err := p.MakeBody()
	if err != nil {
		return err
	}
	if _, err := io.Copy(os.Stdout, r); err != nil {
		return err
	}
	return nil
}

var tpl = `
var proxy="{{.proxy}}";

var domains = {
	{{range .domains}}
	"{{.}}": 1,
	{{end}}
};

function ip4todecimal(ip) {
    var d = ip.split('.');
    return ((((((+d[0])*256)+(+d[1]))*256)+(+d[2]))*256)+(+d[3]);
}

function FindProxyForURL(url, host){
    if(/\d+\.\d+\.\d+\.\d+/.test(host)){
        if (isInNet(dnsResolve(host), "10.0.0.0", "255.0.0.0") ||
                isInNet(dnsResolve(host), "172.16.0.0",  "255.240.0.0") ||
                isInNet(dnsResolve(host), "192.168.0.0", "255.255.0.0") ||
                isInNet(dnsResolve(host), "127.0.0.0", "255.255.255.0")){
            return "DIRECT";
        }
		return "DIRECT";
    }
    if (isPlainHostName(host)){
        return "DIRECT";
    }

    var a = host.split(".");
    for(var i=a.length-1; i>=0; i--){
        if (domains.hasOwnProperty(a.slice(i).join("."))){
			return "DIRECT";
        }
    }
	return proxy;
}
`
