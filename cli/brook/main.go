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

package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/txthinking/brook"
	"github.com/urfave/cli/v2"
)

var debug bool
var debugAddress string

func main() {
	app := cli.NewApp()
	app.Name = "Brook"
	app.Version = "20210601"
	app.Usage = "A cross-platform strong encryption and not detectable proxy"
	app.Authors = []*cli.Author{
		{
			Name:  "Cloud",
			Email: "cloud@txthinking.com",
		},
	}
	app.Copyright = "https://github.com/txthinking/brook"
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "Enable debug",
			Destination: &debug,
		},
		&cli.StringFlag{
			Name:        "listen",
			Aliases:     []string{"l"},
			Usage:       "Listen address for debug",
			Value:       ":6060",
			Destination: &debugAddress,
		},
	}
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "server",
			Usage: "Run as brook server, both TCP and UDP",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, like: ':9999'",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "server")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewServer(c.String("listen"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "servers",
			Usage: "Run as multiple brook servers",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "listenpassword",
					Aliases: []string{"l"},
					Usage:   "Listen address and password, like '0.0.0.0:9999 password'",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.StringSlice("listenpassword")) == 0 {
					cli.ShowCommandHelp(c, "servers")
					return nil
				}
				if debug {
					enableDebug()
				}
				l := make([]*brook.Server, 0)
				for _, v := range c.StringSlice("listenpassword") {
					ss := strings.Split(v, " ")
					if len(ss) != 2 {
						return errors.New("invalid listenpassword")
					}
					s, err := brook.NewServer(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					l = append(l, s)
				}
				for _, v := range l {
					go func(v *brook.Server) {
						log.Println(v.ListenAndServe())
					}(v)
				}
				sigs := make(chan os.Signal, 1)
				signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
				<-sigs
				for _, v := range l {
					v.Shutdown()
				}
				return nil
			},
		},
		&cli.Command{
			Name:  "client",
			Usage: "Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst], [works with $ brook server]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Brook server address, like: 1.2.3.4:9999",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Brook server password",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "http",
					Value: "127.0.0.1:8010",
					Usage: "where to listen for HTTP connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "client")
					return nil
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" {
					return errors.New("socks5 server requires a clear IP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewClient(":"+p, h, c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				http, err := brook.NewSocks5ToHTTP(c.String("http"), c.String("socks5"), "", "", c.Int("tcpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
					http.Shutdown()
				}()
				go func() {
					if err := http.ListenAndServe(); err != nil {
						log.Println(err)
					}
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "map",
			Usage: "Run as mapping, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook <-> to address], works with $ brook server/wsserver/wssserver",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "password",
				},
				&cli.StringFlag{
					Name:    "from",
					Aliases: []string{"f"},
					Usage:   "Listen address, like: 127.0.0.1:83",
				},
				&cli.StringFlag{
					Name:    "to",
					Aliases: []string{"t"},
					Usage:   "Map to where, like: 8.8.8.8:53",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("to") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "map")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewMap(c.String("from"), c.String("to"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "dns",
			Usage: "Run as DNS server, both TCP and UDP, [src <-> $ brook dns <-> $ brook <-> dns server] or [src <-> $ brook dns <-> dns server for bypass], works with $ brook server/wsserver/wssserver",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "password",
				},
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, like: 127.0.0.1:53",
				},
				&cli.StringFlag{
					Name:  "dns",
					Usage: "DNS server for resolving domains NOT in list",
					Value: "8.8.8.8:53",
				},
				&cli.StringFlag{
					Name:  "dnsForBypass",
					Usage: "DNS server for resolving domains in bypass list",
					Value: "223.5.5.5:53",
				},
				&cli.StringFlag{
					Name:  "bypassDomainList",
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/chinadomain.txt",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "dns")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewDNS(c.String("listen"), c.String("server"), c.String("password"), c.String("dns"), c.String("dnsForBypass"), c.String("bypassDomainList"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "tproxy",
			Usage: "Run as transparent proxy, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook <-> dst], works with $ brook server/wsserver/wssserver",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, DO NOT contain IP, just like: ':1080'",
					Value:   ":1080",
				},
				&cli.StringFlag{
					Name:  "dnsListen",
					Usage: "Start a smart DNS server, like: ':53'",
				},
				&cli.StringFlag{
					Name:  "dnsForDefault",
					Usage: "DNS server for resolving domains NOT in list",
					Value: "8.8.8.8:53",
				},
				&cli.StringFlag{
					Name:  "dnsForBypass",
					Usage: "DNS server for resolving domains in bypass list",
					Value: "223.5.5.5:53",
				},
				&cli.StringFlag{
					Name:  "bypassDomainList",
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/chinadomain.txt",
				},
				&cli.StringFlag{
					Name:  "bypassCIDR4List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/chinacidr4.txt",
				},
				&cli.StringFlag{
					Name:  "bypassCIDR6List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/chinacidr6.txt",
				},
				&cli.BoolFlag{
					Name:  "enableIPv6",
					Usage: "Your local and server must support IPv6",
				},
				&cli.BoolFlag{
					Name:  "doNotRunScripts",
					Usage: "This will not change iptables and others",
				},
				&cli.StringFlag{
					Name:  "webListen",
					Usage: "Ignore all other parameters, run web UI, like: ':9999'",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("webListen") != "" {
					if err := os.WriteFile("/root/.brook.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0666); err != nil {
						return err
					}
					web, err := fs.Sub(static, "static")
					if err != nil {
						return err
					}
					var cmd *exec.Cmd
					m := http.NewServeMux()
					m.Handle("/", http.FileServer(http.FS(web)))
					m.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
						s, err := os.Executable()
						if err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						os.WriteFile("/root/.brook.args", []byte(r.FormValue("args")), 0666)
						cmd = exec.Command("/bin/sh", "-c", s+" tproxy "+r.FormValue("args"))
						go func() {
							out, _ := cmd.CombinedOutput()
							os.WriteFile("/root/.brook.tproxy.err", out, 0666)
							cmd = nil
						}()
						w.Write([]byte("Diconnect"))
					})
					m.HandleFunc("/disconnect", func(w http.ResponseWriter, r *http.Request) {
						if cmd == nil {
							w.Write([]byte("Connect"))
							return
						}
						if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						w.Write([]byte("Connect"))
					})
					m.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
						if cmd == nil {
							w.Write([]byte("Connect"))
							return
						}
						w.Write([]byte("Diconnect"))
					})
					s := &http.Server{
						Addr:    c.String("webListen"),
						Handler: m,
					}
					go func() {
						sigs := make(chan os.Signal, 1)
						signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
						<-sigs
						if cmd != nil {
							cmd.Process.Signal(syscall.SIGTERM)
						}
						s.Shutdown(context.Background())
					}()
					return s.ListenAndServe()
				}
				if c.String("listen") == "" || c.String("s") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "tproxy")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewTproxy(c.String("listen"), c.String("server"), c.String("password"), c.Bool("enableIPv6"), c.String("bypassCIDR4List"), c.String("bypassCIDR6List"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				var s1 *brook.DNS
				if c.String("dnsListen") != "" {
					s1, err = brook.NewDNS(c.String("dnsListen"), c.String("server"), c.String("password"), c.String("dnsForDefault"), c.String("dnsForBypass"), c.String("bypassDomainList"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					go func() {
						if err := s1.ListenAndServe(); err != nil {
							log.Println(err)
						}
					}()
				}
				if !c.Bool("doNotRunScripts") {
					s.ClearAutoScripts()
					if err := s.RunAutoScripts(); err != nil {
						return err
					}
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					if !c.Bool("doNotRunScripts") {
						if err := s.ClearAutoScripts(); err != nil {
							log.Println(err)
						}
					}
					s.Shutdown()
					if s1 != nil {
						s1.Shutdown()
					}
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "wsserver",
			Usage: "Run as brook wsserver, both TCP and UDP, it will start a standard http server and websocket server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, like: ':80'",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.StringFlag{
					Name:  "domain",
					Usage: "deprecated, please use $ brook wssserver",
				},
				&cli.StringFlag{
					Name:  "path",
					Usage: "URL path",
					Value: "/ws",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("domain") != "" {
					log.Println("--domain deprecated, please use $ brook wssserver")
				}
				if (c.String("listen") == "" && c.String("domain") == "") || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wsserver")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewWSServer(c.String("listen"), c.String("password"), c.String("domain"), c.String("path"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "wssserver",
			Usage: "Run as brook wssserver, both TCP and UDP, it will start a standard https server and websocket server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "domain",
					Usage: "The domain must have been resolved to the external IP, 80 and 443 ports will be used, TLS certificate will be automatically issued",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.StringFlag{
					Name:  "path",
					Usage: "URL path",
					Value: "/ws",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("domain") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wssserver")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewWSServer("", c.String("password"), c.String("domain"), c.String("path"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "wsclient",
			Usage: "Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst], [works with $ brook wsserver]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "wsserver",
					Aliases: []string{"s"},
					Usage:   "Brook wsserver address, like: ws://1.2.3.4:80, if no path then /ws will be used. Do not omit the port under any circumstances",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Brook wsserver password",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "http",
					Value: "127.0.0.1:8010",
					Usage: "where to listen for HTTP connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" || c.String("wsserver") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wsclient")
					return nil
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" {
					return errors.New("socks5 server requires a clear IP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewWSClient(":"+p, h, c.String("wsserver"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				http, err := brook.NewSocks5ToHTTP(c.String("http"), c.String("socks5"), "", "", c.Int("tcpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
					http.Shutdown()
				}()
				go func() {
					if err := http.ListenAndServe(); err != nil {
						log.Println(err)
					}
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "wssclient",
			Usage: "Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst], [works with $ brook wssserver]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "wssserver",
					Aliases: []string{"s"},
					Usage:   "Brook wssserver address, like: wss://google.com:443, if no path then /ws will be used. Do not omit the port under any circumstances",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Brook wssserver password",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "http",
					Value: "127.0.0.1:8010",
					Usage: "where to listen for HTTP connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" || c.String("wssserver") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wssclient")
					return nil
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" {
					return errors.New("socks5 server requires a clear IP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewWSClient(":"+p, h, c.String("wssserver"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				http, err := brook.NewSocks5ToHTTP(c.String("http"), c.String("socks5"), "", "", c.Int("tcpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
					http.Shutdown()
				}()
				go func() {
					if err := http.ListenAndServe(); err != nil {
						log.Println(err)
					}
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "link",
			Usage: "Print brook link",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Support $ brook server, $ brook wsserver, $ brook wssserver and socks5 server, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://google.com:443/ws, socks5://1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Usage:   "Username, when server is socks5 server",
				},
			},
			Action: func(c *cli.Context) error {
				s := "server"
				if strings.HasPrefix(c.String("s"), "ws://") {
					s = "wsserver"
				}
				if strings.HasPrefix(c.String("s"), "wss://") {
					s = "wssserver"
				}
				if strings.HasPrefix(c.String("s"), "socks5://") {
					s = "socks5"
				}
				fmt.Println(brook.Link(s, c.String("s"), c.String("username"), c.String("password")))
				return nil
			},
		},
		&cli.Command{
			Name:  "qr",
			Usage: "Print brook server QR code",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Support $ brook server,  $ brook wsserver, $ brook wssserver and socks5 server, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://google.com:443, socks5://1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Usage:   "Username, when server is socks5 server",
				},
			},
			Action: func(c *cli.Context) error {
				s := "server"
				if strings.HasPrefix(c.String("s"), "ws://") {
					s = "wsserver"
				}
				if strings.HasPrefix(c.String("s"), "wss://") {
					s = "wssserver"
				}
				if strings.HasPrefix(c.String("s"), "socks5://") {
					s = "socks5"
				}
				brook.QR(s, c.String("s"), c.String("username"), c.String("password"))
				return nil
			},
		},
		&cli.Command{
			Name:  "connect",
			Usage: "Connect via standard sharing link (brook server & brook wsserver & brook wssserver)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "link",
					Aliases: []string{"l"},
					Usage:   "specify the sharing link",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "http",
					Value: "127.0.0.1:8010",
					Usage: "where to listen for HTTP connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("link") == "" {
					cli.ShowCommandHelp(c, "connect")
					return nil
				}
				if debug {
					enableDebug()
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" {
					return errors.New("socks5 server requires a clear IP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				var kind, server, password string
				if !strings.Contains(c.String("link"), "?") {
					kind, server, _, password, err = brook.ParseLinkOld(c.String("link"))
					if err != nil {
						return err
					}
				}
				if strings.Contains(c.String("link"), "?") {
					kind, server, _, password, err = brook.ParseLink(c.String("link"))
					if err != nil {
						return err
					}
				}
				if kind == "socks5" {
					return errors.New("connect doesn't support socks5 link, you may want $ brook socks5tohttp")
				}
				if kind == "server" {
					s, err := brook.NewClient(":"+p, h, server, password, c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					http, err := brook.NewSocks5ToHTTP(c.String("http"), c.String("socks5"), "", "", c.Int("tcpTimeout"))
					if err != nil {
						return err
					}
					go func() {
						sigs := make(chan os.Signal, 1)
						signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
						<-sigs
						s.Shutdown()
						http.Shutdown()
					}()
					go func() {
						if err := http.ListenAndServe(); err != nil {
							log.Println(err)
						}
					}()
					return s.ListenAndServe()
				}
				s, err := brook.NewWSClient(":"+p, h, server, password, c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				http, err := brook.NewSocks5ToHTTP(c.String("http"), c.String("socks5"), "", "", c.Int("tcpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
					http.Shutdown()
				}()
				go func() {
					if err := http.ListenAndServe(); err != nil {
						log.Println(err)
					}
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "relay",
			Usage: "Run as standalone relay, both TCP and UDP, this means access [listen address] is equal to access [to address], [src <-> listen address <-> to address]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "from",
					Aliases: []string{"f"},
					Usage:   "Listen address: like ':9999'",
				},
				&cli.StringFlag{
					Name:    "to",
					Aliases: []string{"t"},
					Usage:   "Address which relay to, like: 1.2.3.4:9999",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("from") == "" || c.String("to") == "" {
					cli.ShowCommandHelp(c, "relay")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewRelay(c.String("from"), c.String("to"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "relays",
			Usage: "Run as multiple standalone relays",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:  "fromto",
					Usage: "Listen address and relay to address, like '0.0.0.0:9999 1.2.3.4:9999'",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.StringSlice("fromto")) == 0 {
					cli.ShowCommandHelp(c, "relays")
					return nil
				}
				if debug {
					enableDebug()
				}
				l := make([]*brook.Relay, 0)
				for _, v := range c.StringSlice("fromto") {
					ss := strings.Split(v, " ")
					if len(ss) != 2 {
						return errors.New("invalid fromto")
					}
					s, err := brook.NewRelay(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					l = append(l, s)
				}
				for _, v := range l {
					go func(v *brook.Relay) {
						log.Println(v.ListenAndServe())
					}(v)
				}
				sigs := make(chan os.Signal, 1)
				signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
				<-sigs
				for _, v := range l {
					v.Shutdown()
				}
				return nil
			},
		},
		&cli.Command{
			Name:  "socks5",
			Usage: "Run as standalone standard socks5 server, both TCP and UDP",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "socks5",
					Usage: "Socks5 server which will be created, like: 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:  "username",
					Usage: "User name, optional",
				},
				&cli.StringFlag{
					Name:  "password",
					Usage: "Password, optional",
				},
				&cli.StringFlag{
					Name:  "bindip",
					Usage: "Default bind to all, such as :PORT",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" {
					cli.ShowCommandHelp(c, "socks5")
					return nil
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" {
					return errors.New("socks5 server requires a clear IP, only port is not enough. You may use public IP or lan IP or other, we can not decide for you")
				}
				if debug {
					enableDebug()
				}
				a := ":" + p
				if c.String("bindip") != "" {
					a = net.JoinHostPort(c.String("bindip"), p)
				}
				s, err := brook.NewSocks5Server(a, h, c.String("username"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "socks5tohttp",
			Usage: "Convert socks5 to http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "socks5",
					Aliases: []string{"s"},
					Usage:   "Socks5 server address, like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:  "socks5username",
					Usage: "Socks5 username, optional",
				},
				&cli.StringFlag{
					Name:  "socks5password",
					Usage: "Socks5 password, optional",
				},
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "HTTP proxy which will be create: like: 127.0.0.1:8010",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection tcp timeout (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("socks5") == "" {
					cli.ShowCommandHelp(c, "socks5tohttp")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewSocks5ToHTTP(c.String("listen"), c.String("socks5"), c.String("socks5username"), c.String("socks5password"), c.Int("tcpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "hijackhttps",
			Usage: "Hijack domains and assume is TCP/TLS/443. Requesting these domains from anywhere in the system will be hijacked . [src <-> $ brook hijackhttps <-> socks5 server] or [src <-> direct]",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "socks5",
					Aliases: []string{"s"},
					Usage:   "Socks5 server address, like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:  "socks5username",
					Usage: "Socks5 username, optional",
				},
				&cli.StringFlag{
					Name:  "socks5password",
					Usage: "Socks5 password, optional",
				},
				&cli.StringFlag{
					Name:  "listenIP",
					Usage: "127.0.0.1 or ::1, will create a DNS server(udp 53 port) with it, and listen tcp 443 port on it",
					Value: "127.0.0.1",
				},
				&cli.StringFlag{
					Name:  "dnsForBypass",
					Usage: "DNS server for resolving domains in bypass list",
					Value: "223.5.5.5:53",
				},
				&cli.StringFlag{
					Name:  "bypassDomainList",
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/chinadomain.txt",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" {
					cli.ShowCommandHelp(c, "hijackhttps")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewHijackHTTPS(c.String("socks5"), c.String("socks5username"), c.String("socks5password"), c.String("listenIP"), c.String("dnsForBypass"), c.String("bypassDomainList"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "pac",
			Usage: "Run as PAC server or save PAC to file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, like: 127.0.0.1:1980",
				},
				&cli.StringFlag{
					Name:    "proxy",
					Aliases: []string{"p"},
					Value:   "SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT",
					Usage:   "Proxy, like: 'SOCKS5 127.0.0.1:1080; SOCKS 127.0.0.1:1080; DIRECT'",
				},
				&cli.StringFlag{
					Name:    "bypassDomainList",
					Aliases: []string{"b"},
					Usage:   "domain list url, http(s):// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/chinadomain.txt",
				},
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "Save PAC to file, this will ignore listen address",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" && c.String("file") == "" {
					cli.ShowCommandHelp(c, "pac")
					return nil
				}
				p := brook.NewPAC(c.String("listen"), c.String("file"), c.String("proxy"), c.String("bypassDomainList"))
				if c.String("file") != "" {
					return p.WriteToFile()
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					p.Shutdown()
				}()
				return p.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "howto",
			Usage: "Print some useful tutorial resources",
			Action: func(c *cli.Context) error {
				fmt.Println("")
				fmt.Println("Brook Github:", "https://github.com/txthinking/brook")
				fmt.Println("Brook Docs:", "https://txthinking.github.io/brook")
				fmt.Println("Brook Community:", "https://github.com/txthinking/brook/discussions")
				fmt.Println("")
				fmt.Println("Blog:", "https://talks.txthinking.com")
				fmt.Println("Youtube:", "https://www.youtube.com/channel/UC5j8-I5Y4lWo4KTa4_0Kx5A")
				fmt.Println("")
				return nil
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}

func enableDebug() {
	go func() {
		log.Println(http.ListenAndServe(debugAddress, nil))
	}()
	brook.EnableDebug()
}
