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
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/socks5"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
)

// WSServer.
type WSServer struct {
	Password    []byte
	Domain      string
	TCPAddr     *net.TCPAddr
	HTTPServer  *http.Server
	HTTPSServer *http.Server
	TCPTimeout  int
	UDPTimeout  int
	Path        string
	UDPSrc      *cache.Cache
}

// NewWSServer.
func NewWSServer(addr, password, domain, path string, tcpTimeout, udpTimeout int) (*WSServer, error) {
	var taddr *net.TCPAddr
	var err error
	if domain == "" {
		taddr, err = net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return nil, err
		}
	}
	cs2 := cache.New(cache.NoExpiration, cache.NoExpiration)
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &WSServer{
		Password:   []byte(password),
		Domain:     domain,
		TCPAddr:    taddr,
		TCPTimeout: tcpTimeout,
		UDPTimeout: udpTimeout,
		Path:       path,
		UDPSrc:     cs2,
	}
	return s, nil
}

// Run server.
func (s *WSServer) ListenAndServe() error {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		return
	})
	r.Methods("GET").Path(s.Path).Handler(s)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	if Debug {
		n.Use(negroni.NewLogger())
	}
	n.UseFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		w.Header().Set("Server", "nginx")
		next(w, r)
	})
	n.UseHandler(r)

	if s.Domain == "" {
		s.HTTPServer = &http.Server{
			Addr:           s.TCPAddr.String(),
			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    120 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        n,
		}
		return s.HTTPServer.ListenAndServe()
	}
	m := autocert.Manager{
		Cache:      autocert.DirCache(".letsencrypt"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(s.Domain),
		Email:      "cloud@txthinking.com",
	}
	go http.ListenAndServe(":80", m.HTTPHandler(nil))
	s.HTTPSServer = &http.Server{
		Addr:         ":443",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      n,
		TLSConfig:    &tls.Config{GetCertificate: m.GetCertificate},
	}
	go func() {
		time.Sleep(1 * time.Second)
		c := &http.Client{
			Timeout: 10 * time.Second,
		}
		_, _ = c.Get("https://" + s.Domain + s.Path)
	}()
	return s.HTTPSServer.ListenAndServeTLS("", "")
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
	if s.TCPTimeout != 0 {
		if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
			log.Println(err)
			return
		}
	}
	ss, dst, err := NewStreamServer(s.Password, c, s.TCPTimeout)
	if err != nil {
		log.Println(err)
		return
	}
	defer ss.Clean()
	if ss.Network == "tcp" {
		if err := s.TCPHandle(ss, dst); err != nil {
			log.Println(err)
		}
	}
	if ss.Network == "udp" {
		ss.Timeout = s.UDPTimeout
		if err := s.UDPHandle(ss, c.RemoteAddr().String(), dst); err != nil {
			log.Println(err)
		}
	}
}

// TCPHandle handles request.
func (s *WSServer) TCPHandle(ss *StreamServer, dst []byte) error {
	address := socks5.ToAddress(dst[0], dst[1:len(dst)-2], dst[len(dst)-2:])
	debug("dial tcp", address)
	rc, err := Dial.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPTimeout) * time.Second)); err != nil {
			return err
		}
	}
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

// UDPHandle handles packet.
func (s *WSServer) UDPHandle(ss *StreamServer, src string, dstb []byte) error {
	dst := socks5.ToAddress(dstb[0], dstb[1:len(dstb)-2], dstb[len(dstb)-2:])
	debug("dial udp", dst)
	var laddr *net.UDPAddr
	any, ok := s.UDPSrc.Get(src + dst)
	if ok {
		laddr = any.(*net.UDPAddr)
	}
	raddr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return err
	}
	rc, err := Dial.DialUDP("udp", laddr, raddr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			// we dont choose lock, so ignore this error
			return nil
		}
		return err
	}
	defer rc.Close()
	if laddr == nil {
		s.UDPSrc.Set(src+dst, rc.LocalAddr().(*net.UDPAddr), -1)
	}
	if err := ss.Exchange(rc); err != nil {
		return nil
	}
	return nil
}

// Shutdown server.
func (s *WSServer) Shutdown() error {
	if s.Domain == "" {
		return s.HTTPServer.Shutdown(context.Background())
	}
	return s.HTTPSServer.Shutdown(context.Background())
}
