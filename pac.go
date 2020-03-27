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
	"github.com/txthinking/x"
)

type PAC struct {
	Addr string
	File string
	// global/white/black
	Proxy      string
	Mode       string
	DomainURL  string
	CidrURL    string
	DomainData []byte
	CidrData   []byte
	HTTPServer *http.Server
}

func NewPAC(addr, file, proxy, mode, domainURL, cidrURL string) *PAC {
	return &PAC{
		Addr:      addr,
		File:      file,
		Proxy:     proxy,
		Mode:      mode,
		DomainURL: domainURL,
		CidrURL:   cidrURL,
	}
}

func (p *PAC) MakeBody() (io.Reader, error) {
	var err error
	if p.DomainURL != "" {
		p.DomainData, err = readURL(p.DomainURL)
		if err != nil {
			return nil, err
		}
	}
	if p.CidrURL != "" {
		p.CidrData, err = readURL(p.CidrURL)
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
	cs := make([]map[string]int64, 0)
	if p.CidrData != nil {
		b := bytes.TrimSpace(p.CidrData)
		b = bytes.Replace(b, []byte{0x20}, []byte{}, -1)
		b = bytes.Replace(b, []byte{0x0d, 0x0a}, []byte{0x0a}, -1)
		ss := strings.Split(string(b), "\n")
		for _, s := range ss {
			c, err := x.CIDR(s)
			if err != nil {
				continue
			}
			first, err := x.IP2Decimal(c.First)
			if err != nil {
				continue
			}
			last, err := x.IP2Decimal(c.Last)
			if err != nil {
				continue
			}
			m := make(map[string]int64)
			m["first"] = first
			m["last"] = last
			cs = append(cs, m)
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
		"mode":    p.Mode,
		"domains": ds,
		"cidrs":   cs,
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

func (p *PAC) ListenAndServe() error {
	r, err := p.MakeBody()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
		w.Write(b)
	})
	p.HTTPServer = &http.Server{
		Addr:           p.Addr,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
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

var mode = "{{.mode}}";

var domains = {
	{{range .domains}}
	"{{.}}": 1,
	{{end}}
};

var cidrs = [
    {{range .cidrs}}
    [{{.first}},{{.last}}],
	{{end}}
];

function ip2decimal(ip) {
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
		if(mode == "global"){
			return proxy;
		}
        var d = ip2decimal(host);
        var l = cidrs.length;
        var min = 0;
        var max = l;
        for(;;){
            if (min+1 > max) {
                break;
            }
            var mid = Math.floor(min+(max-min)/2);
            if(d >= cidrs[mid][0] && d <= cidrs[mid][1]){
				if(mode == "white"){
					return "DIRECT";
				}
				if(mode == "black"){
					return proxy;
				}
            }else if(d < cidrs[mid][0]){
                max = mid;
            }else{
                min = mid+1;
            }
        }
		if(mode == "white"){
			return proxy;
		}
		if(mode == "black"){
			return "DIRECT";
		}
    }

    if (isPlainHostName(host)){
        return "DIRECT";
    }

	if(mode == "global"){
		return proxy;
	}
    var a = host.split(".");
    for(var i=a.length-1; i>=0; i--){
        if (domains.hasOwnProperty(a.slice(i).join("."))){
			if(mode == "white"){
				return "DIRECT";
			}
			if(mode == "black"){
				return proxy;
			}
        }
    }
	if(mode == "white"){
		return proxy;
	}
	if(mode == "black"){
		return "DIRECT";
	}
}
`
