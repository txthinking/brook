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
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/txthinking/brook/limits"
	crypto1 "github.com/txthinking/crypto"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
)

type WSServer struct {
	Password     []byte
	Domain       string
	Addr         string
	HTTPServer   *http.Server
	TCPTimeout   int
	UDPTimeout   int
	Path         string
	Cert         []byte
	CertKey      []byte
	WithoutBrook bool
}

func NewWSServer(addr, password, domain, path string, tcpTimeout, udpTimeout int, withoutbrook bool) (*WSServer, error) {
	if err := limits.Raise(); err != nil {
		Log(&Error{"when": "try to raise system limits", "warning": err.Error()})
	}
	p := []byte(password)
	if withoutbrook {
		var err error
		p, err = crypto1.SHA256Bytes([]byte(password))
		if err != nil {
			return nil, err
		}
	}
	s := &WSServer{
		Password:     p,
		Addr:         addr,
		Domain:       domain,
		Path:         path,
		TCPTimeout:   tcpTimeout,
		UDPTimeout:   udpTimeout,
		WithoutBrook: withoutbrook,
	}
	return s, nil
}

func (s *WSServer) ListenAndServe() error {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		return
	})
	r.Methods("GET").Path(s.Path).Handler(s)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Server", "nginx")
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
			Log(http.ListenAndServe(":80", m.HTTPHandler(nil)))
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  65507,
	WriteBufferSize: 65507,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := conn.UnderlyingConn()
	defer c.Close()
	var ss Exchanger
	if !s.WithoutBrook {
		ss, err = NewStreamServer(s.Password, c.RemoteAddr().String(), c, s.TCPTimeout, s.UDPTimeout)
	}
	if s.WithoutBrook {
		ss, err = NewSimpleStreamServer(s.Password, c.RemoteAddr().String(), c, s.TCPTimeout, s.UDPTimeout)
	}
	if err != nil {
		Log(&Error{"from": c.RemoteAddr().String(), "error": err.Error()})
		return
	}
	defer ss.Clean()
	if ss.Network() == "tcp" {
		if err := s.TCPHandle(ss); err != nil {
			Log(&Error{"from": c.RemoteAddr().String(), "dst": ss.Dst(), "error": err.Error()})
		}
	}
	if ss.Network() == "udp" {
		if err := s.UDPHandle(ss); err != nil {
			Log(&Error{"from": c.RemoteAddr().String(), "dst": ss.Dst(), "error": err.Error()})
		}
	}
}

func (s *WSServer) TCPHandle(ss Exchanger) error {
	rc, err := DialTCP("tcp", "", ss.Dst())
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

func (s *WSServer) UDPHandle(ss Exchanger) error {
	rc, err := NATDial("udp", ss.Src(), ss.Dst(), ss.Dst())
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

func (s *WSServer) Shutdown() error {
	return s.HTTPServer.Shutdown(context.Background())
}
