package brook

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	if p.DomainURL != "" {
		p.DomainData, err = readURL(p.DomainURL)
		if err != nil {
			return nil, err
		}
	}
	ds := make([]string, 0)
	if p.DomainData != nil {
		b := bytes.TrimSpace(p.DomainData)
		b = bytes.Replace(b, []byte{0x20}, []byte{}, -1)
		b = bytes.Replace(b, []byte{0x0d, 0x0a}, []byte{0x0a}, -1)
		ds = strings.Split(string(b), "\n")
	}

	t := template.New("pac")
	t, err = t.Parse(tpl)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	if err := t.Execute(b, map[string]interface{}{
		"proxy":   p.Proxy,
		"domains": ds,
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

func readURL(url string) ([]byte, error) {
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		c := &http.Client{
			Timeout: 9 * time.Second,
		}
		r, err := c.Get(url)
		if err != nil {
			return nil, err
		}
		defer r.Body.Close()
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
	data, err := ioutil.ReadFile(url)
	if err != nil {
		return nil, err
	}
	return data, nil
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
