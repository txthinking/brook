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
	"net/url"

	"github.com/txthinking/brook"
	"github.com/txthinking/runnergroup"
	"github.com/urfave/cli/v2"
)

var debug bool
var debugAddress string

func main() {
	app := cli.NewApp()
	app.Name = "Brook"
	app.Version = "20220217.dev"
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
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/sample_block.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/demo_block_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/demo_block_cidr6.txt",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "Update list interval, second. default 0, only read one time on start",
				},
			},
			Action: func(c *cli.Context) error {
				if debug {
					enableDebug()
				}
				if c.String("listen") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "server")
					return nil
				}
				s, err := brook.NewServer(c.String("listen"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.String("blockDomainList"), c.String("blockCIDR4List"), c.String("blockCIDR6List"), c.Int64("updateListInterval"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
			},
		},
		&cli.Command{
			Name:  "client",
			Usage: "Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst]",
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
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "http",
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
				if debug {
					enableDebug()
				}
				if c.String("socks5") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "client")
					return nil
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" && c.String("socks5ServerIP") == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				var ip string
				if h != "" {
					ip = h
				}
				if c.String("socks5ServerIP") != "" {
					ip = c.String("socks5ServerIP")
				}
				s, err := brook.NewClient(c.String("socks5"), ip, c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				if c.String("http") != "" {
					h, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(ip, p), "", "", c.Int("tcpTimeout"))
					if err != nil {
						return err
					}
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return h.ListenAndServe()
						},
						Stop: func() error {
							return h.Shutdown()
						},
					})
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
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
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/sample_block.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/demo_block_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/demo_block_cidr6.txt",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "Update list interval, second. default 0, only read one time on start",
				},
			},
			Action: func(c *cli.Context) error {
				if debug {
					enableDebug()
				}
				if c.String("listen") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wsserver")
					return nil
				}
				s, err := brook.NewWSServer(c.String("listen"), c.String("password"), "", c.String("path"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.String("blockDomainList"), c.String("blockCIDR4List"), c.String("blockCIDR6List"), c.Int64("updateListInterval"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
			},
		},
		&cli.Command{
			Name:  "wsclient",
			Usage: "Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst]",
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
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "http",
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
				&cli.StringFlag{
					Name:  "address",
					Usage: "Specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
			},
			Action: func(c *cli.Context) error {
				if debug {
					enableDebug()
				}
				if c.String("socks5") == "" || c.String("wsserver") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wsclient")
					return nil
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" && c.String("socks5ServerIP") == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				var ip string
				if h != "" {
					ip = h
				}
				if c.String("socks5ServerIP") != "" {
					ip = c.String("socks5ServerIP")
				}
				s, err := brook.NewWSClient(c.String("socks5"), ip, c.String("wsserver"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				if c.String("address") != "" {
					s.ServerAddress = c.String("address")
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				if c.String("http") != "" {
					h, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(ip, p), "", "", c.Int("tcpTimeout"))
					if err != nil {
						return err
					}
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return h.ListenAndServe()
						},
						Stop: func() error {
							return h.Shutdown()
						},
					})
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
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
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/sample_block.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/demo_block_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/demo_block_cidr6.txt",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "Update list interval, second. default 0, only read one time on start",
				},
			},
			Action: func(c *cli.Context) error {
				if debug {
					enableDebug()
				}
				if c.String("domain") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wssserver")
					return nil
				}
				s, err := brook.NewWSServer("", c.String("password"), c.String("domain"), c.String("path"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.String("blockDomainList"), c.String("blockCIDR4List"), c.String("blockCIDR6List"), c.Int64("updateListInterval"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
			},
		},
		&cli.Command{
			Name:  "wssclient",
			Usage: "Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst]",
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
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "http",
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
				&cli.StringFlag{
					Name:  "address",
					Usage: "Specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "Client do not verify the server's certificate chain and host name",
				},
			},
			Action: func(c *cli.Context) error {
				if debug {
					enableDebug()
				}
				if c.String("socks5") == "" || c.String("wssserver") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wssclient")
					return nil
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" && c.String("socks5ServerIP") == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				var ip string
				if h != "" {
					ip = h
				}
				if c.String("socks5ServerIP") != "" {
					ip = c.String("socks5ServerIP")
				}
				s, err := brook.NewWSClient(c.String("socks5"), ip, c.String("wssserver"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				if c.String("address") != "" {
					s.ServerAddress = c.String("address")
				}
				if c.Bool("insecure") {
					s.TLSConfig.InsecureSkipVerify = true
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				if c.String("http") != "" {
					h, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(ip, p), "", "", c.Int("tcpTimeout"))
					if err != nil {
						return err
					}
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return h.ListenAndServe()
						},
						Stop: func() error {
							return h.Shutdown()
						},
					})
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
			},
		},
		&cli.Command{
			Name:  "relayoverbrook",
			Usage: "Run as relay over brook, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook server/wsserver/wssserver <-> to address]",
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
				&cli.StringFlag{
					Name:  "address",
					Usage: "When server is brook wsserver or brook wssserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "When server is brook wssserver, client do not verify the server's certificate chain and host name",
				},
			},
			Action: func(c *cli.Context) error {
				if debug {
					enableDebug()
				}
				if c.String("from") == "" || c.String("to") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "relayoverbrook")
					return nil
				}
				s, err := brook.NewMap(c.String("from"), c.String("to"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				if (strings.HasPrefix(c.String("server"), "ws://") || strings.HasPrefix(c.String("server"), "wss://")) && c.String("address") != "" {
					s.WSClient.ServerAddress = c.String("address")
				}
				if strings.HasPrefix(c.String("server"), "wss://") && c.Bool("insecure") {
					s.WSClient.TLSConfig.InsecureSkipVerify = true
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
			},
		},
		&cli.Command{
			Name:  "dns",
			Usage: "Run as dns server over brook, both TCP and UDP, [src <-> $ brook dns <-> $ brook server/wsserver/wssserver <-> dns] or [src <-> $ brook dns <-> dnsForBypass]",
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
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/china_domain.txt",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/sample_block.txt",
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
				&cli.StringFlag{
					Name:  "address",
					Usage: "When server is brook wsserver or brook wssserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "When server is brook wssserver, client do not verify the server's certificate chain and host name",
				},
			},
			Action: func(c *cli.Context) error {
				if debug {
					enableDebug()
				}
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "dns")
					return nil
				}
				s, err := brook.NewDNS(c.String("listen"), c.String("server"), c.String("password"), c.String("dns"), c.String("dnsForBypass"), c.String("bypassDomainList"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.String("blockDomainList"))
				if err != nil {
					return err
				}
				if (strings.HasPrefix(c.String("server"), "ws://") || strings.HasPrefix(c.String("server"), "wss://")) && c.String("address") != "" {
					s.WSClient.ServerAddress = c.String("address")
				}
				if strings.HasPrefix(c.String("server"), "wss://") && c.Bool("insecure") {
					s.WSClient.TLSConfig.InsecureSkipVerify = true
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
			},
		},
		&cli.Command{
			Name:  "tproxy",
			Usage: "Run as transparent proxy, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook server/wsserver/wssserver <-> dst]",
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
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/china_domain.txt",
				},
				&cli.StringFlag{
					Name:  "bypassCIDR4List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/china_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "bypassCIDR6List",
					Usage: "https://, http:// or local file path, like: https://txthinking.github.io/bypass/china_cidr6.txt",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/sample_block.txt",
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
				&cli.StringFlag{
					Name:  "address",
					Usage: "When server is brook wsserver or brook wssserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "When server is brook wssserver, client do not verify the server's certificate chain and host name",
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
						w.Write([]byte("Disconnect"))
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
						w.Write([]byte("Disconnect"))
					})
					s := &http.Server{
						Addr:    c.String("webListen"),
						Handler: m,
					}
					g := runnergroup.New()
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return s.ListenAndServe()
						},
						Stop: func() error {
							if cmd != nil {
								cmd.Process.Signal(syscall.SIGTERM)
							}
							return s.Shutdown(context.Background())
						},
					})
					go func() {
						sigs := make(chan os.Signal, 1)
						signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
						<-sigs
						g.Done()
					}()
					return g.Wait()
				}
				if debug {
					enableDebug()
				}
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "tproxy")
					return nil
				}
				s, err := brook.NewTproxy(c.String("listen"), c.String("server"), c.String("password"), c.Bool("enableIPv6"), c.String("bypassCIDR4List"), c.String("bypassCIDR6List"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.String("address"), c.Bool("insecure"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						if !c.Bool("doNotRunScripts") {
							s.ClearAutoScripts()
							if err := s.RunAutoScripts(); err != nil {
								return err
							}
						}
						return s.ListenAndServe()
					},
					Stop: func() error {
						if !c.Bool("doNotRunScripts") {
							if err := s.ClearAutoScripts(); err != nil {
								log.Println(err)
							}
						}
						return s.Shutdown()
					},
				})
				if c.String("dnsListen") != "" {
					s1, err := brook.NewDNS(c.String("dnsListen"), c.String("server"), c.String("password"), c.String("dnsForDefault"), c.String("dnsForBypass"), c.String("bypassDomainList"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.String("blockDomainList"))
					if err != nil {
						return err
					}
					if (strings.HasPrefix(c.String("server"), "ws://") || strings.HasPrefix(c.String("server"), "wss://")) && c.String("address") != "" {
						s1.WSClient.ServerAddress = c.String("address")
					}
					if strings.HasPrefix(c.String("server"), "wss://") && c.Bool("insecure") {
						s1.WSClient.TLSConfig.InsecureSkipVerify = true
					}
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return s1.ListenAndServe()
						},
						Stop: func() error {
							return s1.Shutdown()
						},
					})
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
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
				&cli.StringFlag{
					Name:  "name",
					Usage: "Give this server a name",
				},
				&cli.StringFlag{
					Name:  "address",
					Usage: "When server is brook wsserver or brook wssserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "When server is brook wssserver, client do not verify the server's certificate chain and host name",
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
				v := url.Values{}
				v.Set("name", c.String("name"))
				v.Set("address", c.String("address"))
				yn := ""
				if c.Bool("insecure") {
					yn = "true"
				}
				v.Set("insecure", yn)
				fmt.Println(brook.LinkExtra(s, c.String("s"), c.String("username"), c.String("password"), v))
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
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
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
				if debug {
					enableDebug()
				}
				if c.String("link") == "" {
					cli.ShowCommandHelp(c, "connect")
					return nil
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" && c.String("socks5ServerIP") == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				var ip string
				if h != "" {
					ip = h
				}
				if c.String("socks5ServerIP") != "" {
					ip = c.String("socks5ServerIP")
				}
				kind, server, _, password, v, err := brook.ParseLinkExtra(c.String("link"))
				if err != nil {
					return err
				}
				if kind == "socks5" {
					return errors.New("connect doesn't support socks5 link, you may want $ brook socks5tohttp")
				}
				g := runnergroup.New()
				if kind == "server" {
					s, err := brook.NewClient(c.String("socks5"), ip, server, password, c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return s.ListenAndServe()
						},
						Stop: func() error {
							return s.Shutdown()
						},
					})
				}
				if kind == "wsserver" || kind == "wssserver" {
					s, err := brook.NewWSClient(c.String("socks5"), ip, server, password, c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					if v.Get("address") != "" {
						s.ServerAddress = v.Get("address")
					}
					if kind == "wssserver" && v.Get("insecure") == "true" {
						s.TLSConfig.InsecureSkipVerify = true
					}
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return s.ListenAndServe()
						},
						Stop: func() error {
							return s.Shutdown()
						},
					})
				}
				if c.String("http") != "" {
					h, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(ip, p), "", "", c.Int("tcpTimeout"))
					if err != nil {
						return err
					}
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return h.ListenAndServe()
						},
						Stop: func() error {
							return h.Shutdown()
						},
					})
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
			},
		},
		&cli.Command{
			Name:  "relay",
			Usage: "Run as standalone relay, both TCP and UDP, this means access [from address] is equal to access [to address], [src <-> from address <-> to address]",
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
				if debug {
					enableDebug()
				}
				if c.String("from") == "" || c.String("to") == "" {
					cli.ShowCommandHelp(c, "relay")
					return nil
				}
				s, err := brook.NewRelay(c.String("from"), c.String("to"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
			},
		},
		&cli.Command{
			Name:  "socks5",
			Usage: "Run as standalone standard socks5 server, both TCP and UDP",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Socks5 server listen address, like: :1080 or 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "username",
					Usage: "User name, optional",
				},
				&cli.StringFlag{
					Name:  "password",
					Usage: "Password, optional",
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
				if debug {
					enableDebug()
				}
				if c.String("listen") == "" {
					cli.ShowCommandHelp(c, "listen")
					return nil
				}
				h, _, err := net.SplitHostPort(c.String("listen"))
				if err != nil {
					return err
				}
				if h == "" && c.String("socks5ServerIP") == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use public IP or lan IP or other, we can not decide for you")
				}
				var ip string
				if h != "" {
					ip = h
				}
				if c.String("socks5ServerIP") != "" {
					ip = c.String("socks5ServerIP")
				}
				s, err := brook.NewSocks5Server(c.String("listen"), ip, c.String("username"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
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
				if debug {
					enableDebug()
				}
				if c.String("listen") == "" || c.String("socks5") == "" {
					cli.ShowCommandHelp(c, "socks5tohttp")
					return nil
				}
				s, err := brook.NewSocks5ToHTTP(c.String("listen"), c.String("socks5"), c.String("socks5username"), c.String("socks5password"), c.Int("tcpTimeout"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
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
					Usage: "https://, http:// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/china_domain.txt",
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
				if debug {
					enableDebug()
				}
				if c.String("socks5") == "" {
					cli.ShowCommandHelp(c, "hijackhttps")
					return nil
				}
				s, err := brook.NewHijackHTTPS(c.String("socks5"), c.String("socks5username"), c.String("socks5password"), c.String("listenIP"), c.String("dnsForBypass"), c.String("bypassDomainList"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
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
					Usage:   "domain list url, http(s):// or local file path. Suffix match mode. Like: https://txthinking.github.io/bypass/china_domain.txt",
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
				s := brook.NewPAC(c.String("listen"), c.String("file"), c.String("proxy"), c.String("bypassDomainList"))
				if c.String("file") != "" {
					return s.WriteToFile()
				}
				g := runnergroup.New()
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					g.Done()
				}()
				return g.Wait()
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
				fmt.Println("$ brook servers has been removed, you may like joker:")
				fmt.Println("$ joker brook server ...")
				fmt.Println("$ joker brook server ...")
				return nil
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
				fmt.Println("$ brook relays has been removed, you may like joker:")
				fmt.Println("$ joker brook relay ...")
				fmt.Println("$ joker brook relay ...")
				return nil
			},
		},
		&cli.Command{
			Name:  "map",
			Usage: "Run as mapping, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> brook <-> to address]",
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
				fmt.Println("$ brook map was renamed to $ brook relayoverbrook")
				return nil
			},
		},
		&cli.Command{
			Name:  "howto",
			Usage: "Print some useful tutorial resources",
			Action: func(c *cli.Context) error {
				fmt.Println("")
				fmt.Println("Github:", "https://github.com/txthinking/brook")
				fmt.Println("Docs:", "https://txthinking.github.io/brook")
				fmt.Println("")
				fmt.Println("News:", "https://t.me/txthinking_news")
				fmt.Println("Chat:", "https://t.me/brookgroup")
				fmt.Println("Blog:", "https://talks.txthinking.com")
				fmt.Println("Youtube:", "https://www.youtube.com/txthinking")
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
