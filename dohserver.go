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
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
	"github.com/txthinking/brook/limits"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
)

type DOHServer struct {
	Addr       string
	Domain     string
	Path       string
	DNSClient  *DNSClient
	DOHClient  *DOHClient
	HTTPServer *http.Server
	Cert       []byte
	CertKey    []byte
}

func NewDOHServer(addr, domain, path, to string, tcpTimeout, udpTimeout int) (*DOHServer, error) {
	if err := limits.Raise(); err != nil {
		Log(Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	s := &DOHServer{
		Addr:   addr,
		Domain: domain,
		Path:   path,
	}
	if !strings.HasPrefix(to, "https://") {
		s.DNSClient = &DNSClient{Server: to}
	}
	if strings.HasPrefix(to, "https://") {
		c, err := NewDOHClient(to)
		if err != nil {
			return nil, err
		}
		s.DOHClient = c
	}
	return s, nil
}

func (s *DOHServer) ListenAndServe() error {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		return
	})
	r.Methods("POST").Path(s.Path).Handler(s)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Server", "brook")
		next(w, r)
	})
	n.UseHandler(r)

	if s.Domain == "" {
		s.HTTPServer = &http.Server{
			Addr:           s.Addr,
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        n,
		}
		return s.HTTPServer.ListenAndServe()
	}
	var t *tls.Config
	if s.Cert == nil || s.CertKey == nil {
		m := autocert.Manager{
			Cache:      autocert.DirCache(".letsencrypt"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(s.Domain),
			Email:      "cloud@txthinking.com",
		}
		go func() {
			err := http.ListenAndServe(":80", m.HTTPHandler(nil))
			if err != nil {
				Log(err)
			}
		}()
		t = &tls.Config{GetCertificate: m.GetCertificate}
	}
	if s.Cert != nil && s.CertKey != nil {
		ct, err := tls.X509KeyPair(s.Cert, s.CertKey)
		if err != nil {
			return err
		}
		t = &tls.Config{Certificates: []tls.Certificate{ct}, ServerName: s.Domain}
	}
	s.HTTPServer = &http.Server{
		Addr:         s.Addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      n,
		TLSConfig:    t,
	}
	if s.Cert == nil || s.CertKey == nil {
		go func() {
			time.Sleep(1 * time.Second)
			c := &http.Client{
				Timeout: 10 * time.Second,
			}
			_, _ = c.Get("https://" + s.Domain + s.Addr)
		}()
	}
	return s.HTTPServer.ListenAndServeTLS("", "")
}

var DOHGate func(m *dns.Msg, w http.ResponseWriter, r *http.Request) (done bool, err error) = func(m *dns.Msg, w http.ResponseWriter, r *http.Request) (done bool, err error) {
	if m.Question[0].Qtype == dns.TypeHTTPS || m.Question[0].Qtype == dns.TypeSVCB {
		m1 := &dns.Msg{}
		m1.SetReply(m)
		m1.Authoritative = true
		m1.Answer = append(m1.Answer, &dns.SOA{
			Hdr:     dns.RR_Header{Name: m.Question[0].Name, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: 60},
			Ns:      "txthinking.com.",
			Mbox:    "cloud.txthinking.com.",
			Serial:  uint32((time.Now().Year() * 10000) + (int(time.Now().Month()) * 100) + (time.Now().Day())*100),
			Refresh: 21600,
			Retry:   3600,
			Expire:  259200,
			Minttl:  300,
		})
		m1b, err := m1.PackBuffer(nil)
		if err != nil {
			return false, err
		}
		w.Header().Set("Content-Type", "application/dns-message")
		w.Write(m1b)
		return true, nil
	}
	return false, nil
}

func (s *DOHServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	m := &dns.Msg{}
	if err := m.Unpack(b); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	done, err := DOHGate(m, w, r)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if done {
		return
	}
	m1 := &dns.Msg{}
	if s.DNSClient != nil {
		m1, err = s.DNSClient.Exchange(m)
	}
	if s.DOHClient != nil {
		m1, err = s.DOHClient.Exchange(m)
	}
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	b, err = m1.Pack()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "application/dns-message")
	w.Write(b)
}

func (s *DOHServer) Shutdown() error {
	return s.HTTPServer.Shutdown(context.Background())
}
