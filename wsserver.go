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
	"encoding/binary"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/txthinking/brook/limits"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/socks5"
	"github.com/urfave/negroni"
	"golang.org/x/crypto/acme/autocert"
)

// WSServer.
type WSServer struct {
	Password      []byte
	Domain        string
	TCPAddr       *net.TCPAddr
	HTTPServer    *http.Server
	HTTPSServer   *http.Server
	TCPDeadline   int
	TCPTimeout    int
	UDPDeadline   int
	ServerAuthman plugin.ServerAuthman
}

// NewWSServer.
func NewWSServer(addr, password, domain string, tcpTimeout, tcpDeadline, udpDeadline int) (*WSServer, error) {
	var taddr *net.TCPAddr
	var err error
	if domain == "" {
		taddr, err = net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			return nil, err
		}
	}
	if err := limits.Raise(); err != nil {
		log.Println("Try to raise system limits, got", err)
	}
	s := &WSServer{
		Password:    []byte(password),
		Domain:      domain,
		TCPAddr:     taddr,
		TCPTimeout:  tcpTimeout,
		TCPDeadline: tcpDeadline,
		UDPDeadline: udpDeadline,
	}
	return s, nil
}

// SetServerAuthman sets authman plugin.
func (s *WSServer) SetServerAuthman(m plugin.ServerAuthman) {
	s.ServerAuthman = m
}

// Run server.
func (s *WSServer) ListenAndServe() error {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		return
	})
	r.Methods("GET").Path("/ws").Handler(s)

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
		_, _ = c.Get("https://" + s.Domain + "/ws")
	}()
	return s.HTTPSServer.ListenAndServeTLS("", "")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024*2 + 16,
	WriteBufferSize: 1024*2 + 16,
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
	if s.TCPDeadline != 0 {
		if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
			log.Println(err)
			return
		}
	}
	b := make([]byte, 12+16+10+2)
	if _, err := io.ReadFull(c, b); err != nil {
		return
	}
	l, err := DecryptLength(s.Password, b)
	if err != nil {
		log.Println(err)
		return
	}
	if l-12-16-10 == 1 {
		if err := s.TCPHandle(c); err != nil {
			log.Println(err)
		}
	}
	if l-12-16-10 == 2 {
		if err := s.UDPHandle(c); err != nil {
			log.Println(err)
		}
	}
}

// TCPHandle handles request.
func (s *WSServer) TCPHandle(c net.Conn) error {
	cn := make([]byte, 12)
	if _, err := io.ReadFull(c, cn); err != nil {
		return err
	}
	ck, err := GetKey(s.Password, cn)
	if err != nil {
		return err
	}
	var b []byte
	b, cn, err = ReadFrom(c, ck, cn, true)
	if err != nil {
		return err
	}
	address := socks5.ToAddress(b[0], b[1:len(b)-2], b[len(b)-2:])
	a := b[0]

	var ai plugin.Internet
	if s.ServerAuthman != nil {
		b, cn, err = ReadFrom(c, ck, cn, false)
		if err != nil {
			return err
		}
		ai, err = s.ServerAuthman.VerifyToken(b, "tcp", a, address, nil)
		if err != nil {
			return err
		}
		defer ai.Close()
	}

	if Debug {
		log.Println("Dial TCP", address)
	}
	tmp, err := Dial.Dial("tcp", address)
	if err != nil {
		return err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if s.TCPTimeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(s.TCPTimeout) * time.Second); err != nil {
			return err
		}
	}
	if s.TCPDeadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
			return err
		}
	}

	go func() {
		k, n, err := PrepareKey(s.Password)
		if err != nil {
			log.Println(err)
			return
		}
		i, err := c.Write(n)
		if err != nil {
			return
		}
		if ai != nil {
			if err := ai.TCPEgress(i); err != nil {
				log.Println(err)
				return
			}
		}
		var b [1024 * 2]byte
		for {
			if s.TCPDeadline != 0 {
				if err := rc.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
					return
				}
			}
			i, err := rc.Read(b[:])
			if err != nil {
				return
			}
			n, i, err = WriteTo(c, b[0:i], k, n, false)
			if err != nil {
				return
			}
			if ai != nil {
				if err := ai.TCPEgress(i); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}()

	for {
		if s.TCPDeadline != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.TCPDeadline) * time.Second)); err != nil {
				return nil
			}
		}
		b, cn, err = ReadFrom(c, ck, cn, false)
		if err != nil {
			return nil
		}
		i, err := rc.Write(b)
		if err != nil {
			return nil
		}
		if ai != nil {
			if err := ai.TCPEgress(i); err != nil {
				return err
			}
		}
	}
	return nil
}

// UDPHandle handles packet.
func (s *WSServer) UDPHandle(c net.Conn) error {
	var rc *net.UDPConn
	var ai plugin.Internet
	for {
		if s.UDPDeadline != 0 {
			if err := c.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
				return nil
			}
		}
		b := make([]byte, 12+16+10+2)
		if _, err := io.ReadFull(c, b); err != nil {
			return nil
		}
		l, err := DecryptLength(s.Password, b)
		if err != nil {
			return err
		}
		b = make([]byte, l)
		if _, err := io.ReadFull(c, b); err != nil {
			return nil
		}
		a, h, p, data, err := Decrypt(s.Password, b)
		if err != nil {
			return err
		}
		if rc != nil {
			l := int(binary.BigEndian.Uint16(data[len(data)-2:]))
			data = data[0 : len(data)-2-l]
			i, err := rc.Write(data)
			if err != nil {
				return nil
			}
			if ai != nil {
				if err := ai.UDPEgress(i); err != nil {
					return err
				}
			}
			continue
		}
		address := socks5.ToAddress(a, h, p)
		if s.ServerAuthman != nil {
			l := int(binary.BigEndian.Uint16(data[len(data)-2:]))
			ai, err = s.ServerAuthman.VerifyToken(data[len(data)-2-l:len(data)-2], "udp", a, address, data[0:len(data)-2-l])
			if err != nil {
				return err
			}
			defer ai.Close()
			data = data[0 : len(data)-2-l]
		}
		if Debug {
			log.Println("Dial UDP", address)
		}
		conn, err := Dial.Dial("udp", address)
		if err != nil {
			return err
		}
		rc = conn.(*net.UDPConn)
		i, err := rc.Write(data)
		if err != nil {
			return nil
		}
		if ai != nil {
			if err := ai.UDPEgress(i); err != nil {
				return err
			}
		}
		go func() {
			defer rc.Close()
			var b [65535]byte
			for {
				if s.UDPDeadline != 0 {
					if err := rc.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
						break
					}
				}
				n, err := rc.Read(b[:])
				if err != nil {
					break
				}
				a, addr, port, err := socks5.ParseAddress(c.RemoteAddr().String()) // fake
				if err != nil {
					log.Println(err)
					break
				}
				d := make([]byte, 0, 7)
				d = append(d, a)
				d = append(d, addr...)
				d = append(d, port...)
				d = append(d, b[0:n]...)
				cd, err := EncryptLength(s.Password, d)
				if err != nil {
					log.Println(err)
					break
				}
				i, err := c.Write(cd)
				if err != nil {
					break
				}
				if ai != nil {
					if err := ai.UDPEgress(i); err != nil {
						log.Println(err)
						break
					}
				}
				cd, err = Encrypt(s.Password, d)
				if err != nil {
					log.Println(err)
					break
				}
				i, err = c.Write(cd)
				if err != nil {
					break
				}
				if ai != nil {
					if err := ai.UDPEgress(i); err != nil {
						log.Println(err)
						break
					}
				}
			}
		}()
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
