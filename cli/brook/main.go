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
	"crypto/x509"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"net/http"
	"net/url"

	"github.com/txthinking/brook"
	"github.com/txthinking/brook/plugins/block"
	"github.com/txthinking/brook/plugins/pprof"
	"github.com/txthinking/brook/plugins/socks5dial"
	"github.com/txthinking/brook/plugins/thedns"
	"github.com/txthinking/brook/plugins/tproxy"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
	"github.com/urfave/cli/v2"
)

func main() {
	g := runnergroup.New()
	app := cli.NewApp()
	app.Name = "Brook"
	app.Version = "20230218"
	app.Usage = "A cross-platform network tool designed for developers"
	app.Authors = []*cli.Author{
		{
			Name:  "Cloud",
			Email: "cloud@txthinking.com",
		},
	}
	app.Copyright = "https://github.com/txthinking/brook"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "pprof",
			Usage: "go http pprof listen addr, such as :6060",
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.String("pprof") != "" {
			p, err := pprof.NewPprof(c.String("pprof"))
			if err != nil {
				return err
			}
			g.Add(&runnergroup.Runner{
				Start: func() error {
					return p.ListenAndServe()
				},
				Stop: func() error {
					return p.Shutdown()
				},
			})
		}
		return nil
	}
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "server",
			Usage: "Run as brook server, both TCP and UDP",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
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
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt",
				},
				&cli.StringSliceFlag{
					Name:  "blockGeoIP",
					Usage: "Block IP by Geo country code, such as US",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "Update list interval, second. default 0, only read one time on start",
				},
				&cli.StringFlag{
					Name:  "toSocks5",
					Usage: "Forward to socks5 server, requires your socks5 supports standard socks5 TCP and UDP, such as 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:  "toSocks5Username",
					Usage: "Forward to socks5 server, username",
				},
				&cli.StringFlag{
					Name:  "toSocks5Password",
					Usage: "Forward to socks5 server, password",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("blockCIDR4List") != "" && !strings.HasPrefix(c.String("blockCIDR4List"), "http://") && !strings.HasPrefix(c.String("blockCIDR4List"), "https://") && !filepath.IsAbs(c.String("blockCIDR4List")) {
					return errors.New("--blockCIDR4List must be with absolute path")
				}
				if c.String("blockCIDR6List") != "" && !strings.HasPrefix(c.String("blockCIDR6List"), "http://") && !strings.HasPrefix(c.String("blockCIDR6List"), "https://") && !filepath.IsAbs(c.String("blockCIDR6List")) {
					return errors.New("--blockCIDR6List must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.String("blockCIDR4List") != "" || c.String("blockCIDR6List") != "" || len(c.StringSlice("blockGeoIP")) != 0 {
					p, err := block.NewBlock(c.String("blockDomainList"), c.String("blockCIDR4List"), c.String("blockCIDR6List"), c.StringSlice("blockGeoIP"), c.Int("updateListInterval"))
					if err != nil {
						return err
					}
					p.TouchBrook()
					if c.Int("updateListInterval") != 0 {
						g.Add(&runnergroup.Runner{
							Start: func() error {
								p.Update()
								return nil
							},
							Stop: func() error {
								p.Stop()
								return nil
							},
						})
					}
				}
				s, err := brook.NewServer(c.String("listen"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				if c.String("toSocks5") != "" {
					p, err := socks5dial.NewSocks5Dial(c.String("toSocks5"), c.String("toSocks5Username"), c.String("toSocks5Password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				return nil
			},
		},
		&cli.Command{
			Name:  "client",
			Usage: "Run as brook client, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
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
				&cli.BoolFlag{
					Name:  "udpovertcp",
					Usage: "UDP over TCP",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "Where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "http",
					Usage: "Where to listen for HTTP proxy connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5",
					Usage: "Dial with your socks5 proxy, such as 127.0.0.1:1081",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Username",
					Usage: "Optional",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Password",
					Usage: "Optional",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" || c.String("server") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" && c.String("socks5ServerIP") == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				if c.String("dialWithSocks5") != "" {
					p, err := socks5dial.NewSocks5Dial(c.String("dialWithSocks5"), c.String("dialWithSocks5Username"), c.String("dialWithSocks5Password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
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
				s.UDPOverTCP = c.Bool("udpovertcp")
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
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
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
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "The data will not be encrypted with brook protocol",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt",
				},
				&cli.StringSliceFlag{
					Name:  "blockGeoIP",
					Usage: "Block IP by Geo country code, such as US",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "Update list interval, second. default 0, only read one time on start",
				},
				&cli.StringFlag{
					Name:  "toSocks5",
					Usage: "Forward to socks5 server, requires your socks5 supports standard socks5 TCP and UDP, such as 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:  "toSocks5Username",
					Usage: "Forward to socks5 server, username",
				},
				&cli.StringFlag{
					Name:  "toSocks5Password",
					Usage: "Forward to socks5 server, password",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("blockCIDR4List") != "" && !strings.HasPrefix(c.String("blockCIDR4List"), "http://") && !strings.HasPrefix(c.String("blockCIDR4List"), "https://") && !filepath.IsAbs(c.String("blockCIDR4List")) {
					return errors.New("--blockCIDR4List must be with absolute path")
				}
				if c.String("blockCIDR6List") != "" && !strings.HasPrefix(c.String("blockCIDR6List"), "http://") && !strings.HasPrefix(c.String("blockCIDR6List"), "https://") && !filepath.IsAbs(c.String("blockCIDR6List")) {
					return errors.New("--blockCIDR6List must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.String("blockCIDR4List") != "" || c.String("blockCIDR6List") != "" || len(c.StringSlice("blockGeoIP")) != 0 {
					p, err := block.NewBlock(c.String("blockDomainList"), c.String("blockCIDR4List"), c.String("blockCIDR6List"), c.StringSlice("blockGeoIP"), c.Int("updateListInterval"))
					if err != nil {
						return err
					}
					p.TouchBrook()
					if c.Int("updateListInterval") != 0 {
						g.Add(&runnergroup.Runner{
							Start: func() error {
								p.Update()
								return nil
							},
							Stop: func() error {
								p.Stop()
								return nil
							},
						})
					}
				}
				if c.String("toSocks5") != "" {
					p, err := socks5dial.NewSocks5Dial(c.String("toSocks5"), c.String("toSocks5Username"), c.String("toSocks5Password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				s, err := brook.NewWSServer(c.String("listen"), c.String("password"), "", c.String("path"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.Bool("withoutBrookProtocol"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "wsclient",
			Usage: "Run as brook wsclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
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
					Name:  "address",
					Usage: "Specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "The data will not be encrypted with brook protocol",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "Where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "http",
					Usage: "Where to listen for HTTP proxy connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5",
					Usage: "Dial with your socks5 proxy, such as 127.0.0.1:1081",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Username",
					Usage: "Optional",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Password",
					Usage: "Optional",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" || c.String("wsserver") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" && c.String("socks5ServerIP") == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				if c.String("dialWithSocks5") != "" {
					p, err := socks5dial.NewSocks5Dial(c.String("dialWithSocks5"), c.String("dialWithSocks5Username"), c.String("dialWithSocks5Password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				var ip string
				if h != "" {
					ip = h
				}
				if c.String("socks5ServerIP") != "" {
					ip = c.String("socks5ServerIP")
				}
				s, err := brook.NewWSClient(c.String("socks5"), ip, c.String("wsserver"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.Bool("withoutBrookProtocol"))
				if err != nil {
					return err
				}
				if c.String("address") != "" {
					s.ServerAddress = c.String("address")
				}
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
				return nil
			},
		},
		&cli.Command{
			Name:  "wssserver",
			Usage: "Run as brook wssserver, both TCP and UDP, it will start a standard https server and websocket server",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "domainaddress",
					Usage: "Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.StringFlag{
					Name:  "cert",
					Usage: "The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically",
				},
				&cli.StringFlag{
					Name:  "certkey",
					Usage: "The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically",
				},
				&cli.StringFlag{
					Name:  "path",
					Usage: "URL path",
					Value: "/ws",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "The data will not be encrypted with brook protocol",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt",
				},
				&cli.StringSliceFlag{
					Name:  "blockGeoIP",
					Usage: "Block IP by Geo country code, such as US",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "Update list interval, second. default 0, only read one time on start",
				},
				&cli.StringFlag{
					Name:  "toSocks5",
					Usage: "Forward to socks5 server, requires your socks5 supports standard socks5 TCP and UDP, such as 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:  "toSocks5Username",
					Usage: "Forward to socks5 server, username",
				},
				&cli.StringFlag{
					Name:  "toSocks5Password",
					Usage: "Forward to socks5 server, password",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("domainaddress") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("blockCIDR4List") != "" && !strings.HasPrefix(c.String("blockCIDR4List"), "http://") && !strings.HasPrefix(c.String("blockCIDR4List"), "https://") && !filepath.IsAbs(c.String("blockCIDR4List")) {
					return errors.New("--blockCIDR4List must be with absolute path")
				}
				if c.String("blockCIDR6List") != "" && !strings.HasPrefix(c.String("blockCIDR6List"), "http://") && !strings.HasPrefix(c.String("blockCIDR6List"), "https://") && !filepath.IsAbs(c.String("blockCIDR6List")) {
					return errors.New("--blockCIDR6List must be with absolute path")
				}
				if c.String("cert") != "" && !filepath.IsAbs(c.String("cert")) {
					return errors.New("--cert must be with absolute path")
				}
				if c.String("certkey") != "" && !filepath.IsAbs(c.String("certkey")) {
					return errors.New("--certkey must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.String("blockCIDR4List") != "" || c.String("blockCIDR6List") != "" || len(c.StringSlice("blockGeoIP")) != 0 {
					p, err := block.NewBlock(c.String("blockDomainList"), c.String("blockCIDR4List"), c.String("blockCIDR6List"), c.StringSlice("blockGeoIP"), c.Int("updateListInterval"))
					if err != nil {
						return err
					}
					p.TouchBrook()
					if c.Int("updateListInterval") != 0 {
						g.Add(&runnergroup.Runner{
							Start: func() error {
								p.Update()
								return nil
							},
							Stop: func() error {
								p.Stop()
								return nil
							},
						})
					}
				}
				if c.String("toSocks5") != "" {
					p, err := socks5dial.NewSocks5Dial(c.String("toSocks5"), c.String("toSocks5Username"), c.String("toSocks5Password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				h, p, err := net.SplitHostPort(c.String("domainaddress"))
				if err != nil {
					return err
				}
				s, err := brook.NewWSServer(":"+p, c.String("password"), h, c.String("path"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.Bool("withoutBrookProtocol"))
				if err != nil {
					return err
				}
				if c.String("cert") != "" {
					b, err := ioutil.ReadFile(c.String("cert"))
					if err != nil {
						return err
					}
					s.Cert = b
				}
				if c.String("certkey") != "" {
					b, err := ioutil.ReadFile(c.String("certkey"))
					if err != nil {
						return err
					}
					s.CertKey = b
				}
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				return nil
			},
		},
		&cli.Command{
			Name:  "wssclient",
			Usage: "Run as brook wssclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
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
					Name:  "address",
					Usage: "Specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "Client do not verify the server's certificate chain and host name",
				},
				&cli.StringFlag{
					Name:  "ca",
					Usage: "When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "The data will not be encrypted with brook protocol",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "Where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "http",
					Usage: "Where to listen for HTTP proxy connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5",
					Usage: "Dial with your socks5 proxy, such as 127.0.0.1:1081",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Username",
					Usage: "Optional",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Password",
					Usage: "Optional",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" || c.String("wssserver") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if h == "" && c.String("socks5ServerIP") == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				if c.String("dialWithSocks5") != "" {
					p, err := socks5dial.NewSocks5Dial(c.String("dialWithSocks5"), c.String("dialWithSocks5Username"), c.String("dialWithSocks5Password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				var ip string
				if h != "" {
					ip = h
				}
				if c.String("socks5ServerIP") != "" {
					ip = c.String("socks5ServerIP")
				}
				s, err := brook.NewWSClient(c.String("socks5"), ip, c.String("wssserver"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.Bool("withoutBrookProtocol"))
				if err != nil {
					return err
				}
				if c.String("address") != "" {
					s.ServerAddress = c.String("address")
				}
				if c.Bool("insecure") {
					s.TLSConfig.InsecureSkipVerify = true
				}
				if c.String("ca") != "" {
					b, err := ioutil.ReadFile(c.String("ca"))
					if err != nil {
						return err
					}
					roots := x509.NewCertPool()
					ok := roots.AppendCertsFromPEM(b)
					if !ok {
						return errors.New("failed to parse root certificate")
					}
					s.TLSConfig.RootCAs = roots
				}
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
				return nil
			},
		},
		&cli.Command{
			Name:  "quicserver",
			Usage: "Run as brook quicserver, both TCP and UDP",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "domainaddress",
					Usage: "Such as: domain.com:443. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.StringFlag{
					Name:  "cert",
					Usage: "The cert file absolute path for the domain, such as /path/to/cert.pem. If cert or certkey is empty, a certificate will be issued automatically",
				},
				&cli.StringFlag{
					Name:  "certkey",
					Usage: "The cert key file absolute path for the domain, such as /path/to/certkey.pem. If cert or certkey is empty, a certificate will be issued automatically",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "The data will not be encrypted with brook protocol",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt",
				},
				&cli.StringSliceFlag{
					Name:  "blockGeoIP",
					Usage: "Block IP by Geo country code, such as US",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "Update list interval, second. default 0, only read one time on start",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("domainaddress") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("blockCIDR4List") != "" && !strings.HasPrefix(c.String("blockCIDR4List"), "http://") && !strings.HasPrefix(c.String("blockCIDR4List"), "https://") && !filepath.IsAbs(c.String("blockCIDR4List")) {
					return errors.New("--blockCIDR4List must be with absolute path")
				}
				if c.String("blockCIDR6List") != "" && !strings.HasPrefix(c.String("blockCIDR6List"), "http://") && !strings.HasPrefix(c.String("blockCIDR6List"), "https://") && !filepath.IsAbs(c.String("blockCIDR6List")) {
					return errors.New("--blockCIDR6List must be with absolute path")
				}
				if c.String("cert") != "" && !filepath.IsAbs(c.String("cert")) {
					return errors.New("--cert must be with absolute path")
				}
				if c.String("certkey") != "" && !filepath.IsAbs(c.String("certkey")) {
					return errors.New("--certkey must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.String("blockCIDR4List") != "" || c.String("blockCIDR6List") != "" || len(c.StringSlice("blockGeoIP")) != 0 {
					p, err := block.NewBlock(c.String("blockDomainList"), c.String("blockCIDR4List"), c.String("blockCIDR6List"), c.StringSlice("blockGeoIP"), c.Int("updateListInterval"))
					if err != nil {
						return err
					}
					p.TouchBrook()
					if c.Int("updateListInterval") != 0 {
						g.Add(&runnergroup.Runner{
							Start: func() error {
								p.Update()
								return nil
							},
							Stop: func() error {
								p.Stop()
								return nil
							},
						})
					}
				}
				h, p, err := net.SplitHostPort(c.String("domainaddress"))
				if err != nil {
					return err
				}
				s, err := brook.NewQUICServer(":"+p, c.String("password"), h, c.Int("tcpTimeout"), c.Int("udpTimeout"), c.Bool("withoutBrookProtocol"))
				if err != nil {
					return err
				}
				if c.String("cert") != "" {
					b, err := ioutil.ReadFile(c.String("cert"))
					if err != nil {
						return err
					}
					s.Cert = b
				}
				if c.String("certkey") != "" {
					b, err := ioutil.ReadFile(c.String("certkey"))
					if err != nil {
						return err
					}
					s.CertKey = b
				}
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				return nil
			},
		},
		&cli.Command{
			Name:  "quicclient",
			Usage: "Run as brook quicclient, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook quicclient <-> $ brook quicserver <-> dst]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "quicserver",
					Aliases: []string{"s"},
					Usage:   "Brook quicserver address, like: quic://google.com:443. Do not omit the port under any circumstances",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Brook quicserver password",
				},
				&cli.StringFlag{
					Name:  "address",
					Usage: "Specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.StringFlag{
					Name:  "ca",
					Usage: "When server is brook wssserver, specify ca instead of insecure, such as /path/to/ca.pem",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "The data will not be encrypted with brook protocol",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "Where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "http",
					Usage: "Where to listen for HTTP proxy connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" || c.String("quicserver") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
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
				s, err := brook.NewQUICClient(c.String("socks5"), ip, c.String("quicserver"), c.String("password"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.Bool("withoutBrookProtocol"))
				if err != nil {
					return err
				}
				if c.String("address") != "" {
					s.ServerAddress = c.String("address")
				}
				if c.String("ca") != "" {
					b, err := ioutil.ReadFile(c.String("ca"))
					if err != nil {
						return err
					}
					roots := x509.NewCertPool()
					ok := roots.AppendCertsFromPEM(b)
					if !ok {
						return errors.New("failed to parse root certificate")
					}
					s.TLSConfig.RootCAs = roots
				}
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
				return nil
			},
		},
		&cli.Command{
			Name:  "relayoverbrook",
			Usage: "Run as relay over brook, both TCP and UDP, this means access [from address] is equal to [to address], [src <-> from address <-> $ brook server/wsserver/wssserver/quicserver <-> to address]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "from",
					Aliases: []string{"f", "l"},
					Usage:   "Listen address: like ':9999'",
				},
				&cli.StringFlag{
					Name:    "to",
					Aliases: []string{"t"},
					Usage:   "Address which relay to, like: 1.2.3.4:9999",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws, quic://domain.com:443",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.BoolFlag{
					Name:  "udpovertcp",
					Usage: "When server is brook server, UDP over TCP",
				},
				&cli.StringFlag{
					Name:  "address",
					Usage: "When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol",
				},
				&cli.StringFlag{
					Name:  "ca",
					Usage: "When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("from") == "" || c.String("to") == "" || c.String("server") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				kind := "server"
				if strings.HasPrefix(c.String("server"), "ws://") {
					kind = "wsserver"
				}
				if strings.HasPrefix(c.String("server"), "wss://") {
					kind = "wssserver"
				}
				if strings.HasPrefix(c.String("server"), "quic://") {
					kind = "quicserver"
				}
				v := url.Values{}
				if c.Bool("udpovertcp") {
					v.Set("udpovertcp", "true")
				}
				if c.String("address") != "" {
					v.Set("address", c.String("address"))
				}
				if c.Bool("insecure") {
					v.Set("insecure", "true")
				}
				if c.Bool("withoutBrookProtocol") {
					v.Set("withoutBrookProtocol", "true")
				}
				if c.String("ca") != "" {
					b, err := ioutil.ReadFile(c.String("ca"))
					if err != nil {
						return err
					}
					v.Set("ca", string(b))
				}
				v.Set("password", c.String("password"))
				s, err := brook.NewRelayOverBrook(c.String("from"), brook.Link(kind, c.String("server"), v), c.String("to"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "dnsserveroverbrook",
			Usage: "Run as dns server over brook, both TCP and UDP, [src <-> $ brook dnserversoverbrook <-> $ brook server/wsserver/wssserver/quicserver <-> dns] or [src <-> $ brook dnsserveroverbrook <-> dnsForBypass]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
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
					Usage: "One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.BoolFlag{
					Name:  "disableA",
					Usage: "Disable A query",
				},
				&cli.BoolFlag{
					Name:  "disableAAAA",
					Usage: "Disable AAAA query",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.BoolFlag{
					Name:  "udpovertcp",
					Usage: "When server is brook server, UDP over TCP",
				},
				&cli.StringFlag{
					Name:  "address",
					Usage: "When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name",
				},
				&cli.StringFlag{
					Name:  "ca",
					Usage: "When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("bypassDomainList") != "" && !strings.HasPrefix(c.String("bypassDomainList"), "http://") && !strings.HasPrefix(c.String("bypassDomainList"), "https://") && !filepath.IsAbs(c.String("bypassDomainList")) {
					return errors.New("--bypassDomainList must be with absolute path")
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.String("bypassDomainList") != "" || c.Bool("disableA") || c.Bool("disableAAAA") {
					p, err := thedns.NewTheDNS(c.String("blockDomainList"), c.String("bypassDomainList"), c.String("dnsForBypass"), c.Bool("disableA"), c.Bool("disableAAAA"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				kind := "server"
				if strings.HasPrefix(c.String("server"), "ws://") {
					kind = "wsserver"
				}
				if strings.HasPrefix(c.String("server"), "wss://") {
					kind = "wssserver"
				}
				if strings.HasPrefix(c.String("server"), "quic://") {
					kind = "quicserver"
				}
				v := url.Values{}
				if c.Bool("udpovertcp") {
					v.Set("udpovertcp", "true")
				}
				if c.String("address") != "" {
					v.Set("address", c.String("address"))
				}
				if c.Bool("insecure") {
					v.Set("insecure", "true")
				}
				if c.Bool("withoutBrookProtocol") {
					v.Set("withoutBrookProtocol", "true")
				}
				if c.String("ca") != "" {
					b, err := ioutil.ReadFile(c.String("ca"))
					if err != nil {
						return err
					}
					v.Set("ca", string(b))
				}
				v.Set("password", c.String("password"))
				s, err := brook.NewRelayOverBrook(c.String("listen"), brook.Link(kind, c.String("server"), v), c.String("dns"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				s.IsDNS = true
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				return nil
			},
		},
		&cli.Command{
			Name:  "tproxy",
			Usage: "Run as transparent proxy, both TCP and UDP, only works on Linux, [src <-> $ brook tproxy <-> $ brook server/wsserver/wssserver/quicserver <-> dst]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, DO NOT contain IP, just like: ':1080'. No need to operate iptables by default!",
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
					Name:  "blockDomainList",
					Usage: "One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.StringFlag{
					Name:  "bypassDomainList",
					Usage: "One domain per line, Suffix match mode. https://, http:// or local file absolute path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.BoolFlag{
					Name:  "disableA",
					Usage: "Disable A query",
				},
				&cli.BoolFlag{
					Name:  "disableAAAA",
					Usage: "Disable AAAA query",
				},
				&cli.StringFlag{
					Name:  "bypassCIDR4List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr4.txt",
				},
				&cli.StringFlag{
					Name:  "bypassCIDR6List",
					Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://txthinking.github.io/bypass/example_cidr6.txt",
				},
				&cli.StringSliceFlag{
					Name:  "bypassGeoIP",
					Usage: "Bypass IP by Geo country code, such as US",
				},
				&cli.BoolFlag{
					Name:  "enableIPv6",
					Usage: "deprecated",
				},
				&cli.BoolFlag{
					Name:  "doNotRunScripts",
					Usage: "This will not change iptables and others if you want to do by yourself",
				},
				&cli.StringFlag{
					Name:  "webListen",
					Usage: "Ignore all other parameters, run web UI, like: ':9999'",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "brook server or brook wsserver or brook wssserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.BoolFlag{
					Name:  "udpovertcp",
					Usage: "When server is brook server, UDP over TCP",
				},
				&cli.StringFlag{
					Name:  "address",
					Usage: "When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol",
				},
				&cli.StringFlag{
					Name:  "ca",
					Usage: "When server is brook wssserver or brook quicserver, specify ca instead of insecure, such as /path/to/ca.pem",
				},
				&cli.StringFlag{
					Name:  "link",
					Usage: "brook link. This will ignore server, password, udpovertcp, address, insecure, withoutBrookProtocol, ca",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5",
					Usage: "Dial with your socks5 proxy, such as 127.0.0.1:1081",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Username",
					Usage: "Optional",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Password",
					Usage: "Optional",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("webListen") != "" {
					web, err := fs.Sub(static, "static")
					if err != nil {
						return err
					}
					var cmd *exec.Cmd
					lock := &sync.Mutex{}
					m := http.NewServeMux()
					m.Handle("/", http.FileServer(http.FS(web)))
					m.HandleFunc("/hasp", func(w http.ResponseWriter, r *http.Request) {
						lock.Lock()
						defer lock.Unlock()
						_, err := os.Stat("/root/.brook.web.password")
						if os.IsNotExist(err) {
							w.Write([]byte("no"))
							return
						}
						w.Write([]byte("yes"))
					})
					m.HandleFunc("/setp", func(w http.ResponseWriter, r *http.Request) {
						lock.Lock()
						defer lock.Unlock()
						_, err := os.Stat("/root/.brook.web.password")
						if !os.IsNotExist(err) {
							http.Error(w, "file exsits", 500)
							return
						}
						err = ioutil.WriteFile("/root/.brook.web.password", []byte(r.FormValue("p")), 0600)
						if err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						w.WriteHeader(200)
					})
					m.HandleFunc("/authp", func(w http.ResponseWriter, r *http.Request) {
						lock.Lock()
						defer lock.Unlock()
						b, err := ioutil.ReadFile("/root/.brook.web.password")
						if err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						if string(b) != r.FormValue("p") {
							http.Error(w, "web ui password wrong", 500)
							return
						}
						w.WriteHeader(200)
					})
					m.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
						b, err := ioutil.ReadFile("/root/.brook.web.password")
						if err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						if string(b) != r.FormValue("p") {
							http.Error(w, "web ui password wrong", 500)
							return
						}
						s, err := os.Executable()
						if err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						lock.Lock()
						defer lock.Unlock()
						cmd = exec.Command("/bin/sh", "-c", s+" tproxy "+r.FormValue("args"))
						done := make(chan byte)
						defer close(done)
						errch := make(chan error)
						go func() {
							out, _ := cmd.CombinedOutput()
							select {
							case <-done:
								log.Println(string(out))
							default:
								select {
								case <-done:
									log.Println(string(out))
								case errch <- errors.New(string(out)):
								}
							}
							lock.Lock()
							cmd = nil
							lock.Unlock()
						}()
						select {
						case err := <-errch:
							http.Error(w, err.Error(), 500)
							return
						case <-time.After(6 * time.Second):
							w.Write([]byte("connected"))
						}
					})
					m.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
						b, err := ioutil.ReadFile("/root/.brook.web.password")
						if err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						if string(b) != r.FormValue("p") {
							http.Error(w, "web ui password wrong", 500)
							return
						}
						lock.Lock()
						defer lock.Unlock()
						if cmd == nil {
							w.Write([]byte("disconnected"))
							return
						}
						if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						w.Write([]byte("disconnected"))
					})
					m.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
						b, err := ioutil.ReadFile("/root/.brook.web.password")
						if err != nil {
							http.Error(w, err.Error(), 500)
							return
						}
						if string(b) != r.FormValue("p") {
							http.Error(w, "web ui password wrong", 500)
							return
						}
						lock.Lock()
						defer lock.Unlock()
						if cmd == nil {
							w.Write([]byte("disconnected"))
							return
						}
						w.Write([]byte("connected"))
					})
					s := &http.Server{
						Addr:    c.String("webListen"),
						Handler: m,
					}
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
					return nil
				}
				if c.String("listen") == "" || (c.String("link") == "" && (c.String("server") == "" || c.String("password") == "")) {
					_ = cli.ShowSubcommandHelp(c)
					return errors.New("")
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("bypassDomainList") != "" && !strings.HasPrefix(c.String("bypassDomainList"), "http://") && !strings.HasPrefix(c.String("bypassDomainList"), "https://") && !filepath.IsAbs(c.String("bypassDomainList")) {
					return errors.New("--bypassDomainList must be with absolute path")
				}
				if c.String("bypassCIDR4List") != "" && !strings.HasPrefix(c.String("bypassCIDR4List"), "http://") && !strings.HasPrefix(c.String("bypassCIDR4List"), "https://") && !filepath.IsAbs(c.String("bypassCIDR4List")) {
					return errors.New("--bypassCIDR4List must be with absolute path")
				}
				if c.String("bypassCIDR6List") != "" && !strings.HasPrefix(c.String("bypassCIDR6List"), "http://") && !strings.HasPrefix(c.String("bypassCIDR6List"), "https://") && !filepath.IsAbs(c.String("bypassCIDR6List")) {
					return errors.New("--bypassCIDR6List must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.String("bypassDomainList") != "" || c.Bool("disableA") || c.Bool("disableAAAA") {
					p, err := thedns.NewTheDNS(c.String("blockDomainList"), c.String("bypassDomainList"), c.String("dnsForBypass"), c.Bool("disableA"), c.Bool("disableAAAA"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				if c.String("bypassCIDR4List") != "" || c.String("bypassCIDR6List") != "" || len(c.StringSlice("bypassGeoIP")) != 0 {
					p, err := tproxy.NewTproxy(c.String("bypassCIDR4List"), c.String("bypassCIDR6List"), c.StringSlice("bypassGeoIP"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				if c.String("dialWithSocks5") != "" {
					p, err := socks5dial.NewSocks5Dial(c.String("dialWithSocks5"), c.String("dialWithSocks5Username"), c.String("dialWithSocks5Password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				kind := "server"
				if strings.HasPrefix(c.String("server"), "ws://") {
					kind = "wsserver"
				}
				if strings.HasPrefix(c.String("server"), "wss://") {
					kind = "wssserver"
				}
				if strings.HasPrefix(c.String("server"), "quic://") {
					kind = "quicserver"
				}
				v := url.Values{}
				if c.Bool("udpovertcp") {
					v.Set("udpovertcp", "true")
				}
				if c.String("address") != "" {
					v.Set("address", c.String("address"))
				}
				if c.Bool("insecure") {
					v.Set("insecure", "true")
				}
				if c.Bool("withoutBrookProtocol") {
					v.Set("withoutBrookProtocol", "true")
				}
				if c.String("ca") != "" {
					b, err := ioutil.ReadFile(c.String("ca"))
					if err != nil {
						return err
					}
					v.Set("ca", string(b))
				}
				v.Set("password", c.String("password"))
				link := brook.Link(kind, c.String("server"), v)
				if c.String("link") != "" {
					link = c.String("link")
				}
				s, err := brook.NewTproxy(c.String("listen"), link, c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
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
					s, err := brook.NewRelayOverBrook(c.String("dnsListen"), link, c.String("dnsForDefault"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					s.IsDNS = true
					g.Add(&runnergroup.Runner{
						Start: func() error {
							return s.ListenAndServe()
						},
						Stop: func() error {
							return s.Shutdown()
						},
					})
				}
				return nil
			},
		},
		&cli.Command{
			Name:  "link",
			Usage: "Generate brook link",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Support brook server, brook wsserver, brook wssserver, socks5 server, brook quicserver. Like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://google.com:443/ws, socks5://1.2.3.4:1080, quic://google.com:443",
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
				&cli.BoolFlag{
					Name:  "udpovertcp",
					Usage: "When server is brook server, UDP over TCP",
				},
				&cli.StringFlag{
					Name:  "address",
					Usage: "When server is brook wsserver or brook wssserver or brook quicserver, specify address instead of resolving addresses from host, such as 1.2.3.4:443",
				},
				&cli.BoolFlag{
					Name:  "insecure",
					Usage: "When server is brook wssserver or brook quicserver, client do not verify the server's certificate chain and host name",
				},
				&cli.BoolFlag{
					Name:  "withoutBrookProtocol",
					Usage: "When server is brook wsserver or brook wssserver or brook quicserver, the data will not be encrypted with brook protocol",
				},
				&cli.StringFlag{
					Name:  "ca",
					Usage: "When server is brook wssserver or brook quicserver, specify ca for untrusted cert, such as /path/to/ca.pem",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				s := "server"
				if strings.HasPrefix(c.String("server"), "ws://") {
					s = "wsserver"
				}
				if strings.HasPrefix(c.String("server"), "wss://") {
					s = "wssserver"
				}
				if strings.HasPrefix(c.String("server"), "socks5://") {
					s = "socks5"
				}
				if strings.HasPrefix(c.String("server"), "quic://") {
					s = "quicserver"
				}
				v := url.Values{}
				v.Set("password", c.String("password"))
				if c.String("username") != "" {
					v.Set("username", c.String("username"))
				}
				if c.String("name") != "" {
					v.Set("name", c.String("name"))
				}
				if c.Bool("udpovertcp") {
					v.Set("udpovertcp", "true")
				}
				if c.String("address") != "" {
					v.Set("address", c.String("address"))
				}
				if c.Bool("insecure") {
					v.Set("insecure", "true")
				}
				if c.Bool("withoutBrookProtocol") {
					v.Set("withoutBrookProtocol", "true")
				}
				if c.String("ca") != "" {
					b, err := ioutil.ReadFile(c.String("ca"))
					if err != nil {
						return err
					}
					v.Set("ca", string(b))
				}
				fmt.Println(brook.Link(s, c.String("server"), v))
				return nil
			},
		},
		&cli.Command{
			Name:  "connect",
			Usage: "Run as client and connect to brook link, both TCP and UDP, to start a socks5 proxy, [src <-> socks5 <-> $ brook connect <-> $ brook server/wsserver/wssserver/quicserver <-> dst]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "link",
					Aliases: []string{"l"},
					Usage:   "brook link, you can get it via $ brook link",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:1080",
					Usage: "Where to listen for SOCKS5 connections",
				},
				&cli.StringFlag{
					Name:  "socks5ServerIP",
					Usage: "Only if your socks5 server IP is different from listen IP",
				},
				&cli.StringFlag{
					Name:  "http",
					Usage: "Where to listen for HTTP proxy connections",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5",
					Usage: "If you already have a socks5, such as 127.0.0.1:1081, and want [src <-> listen socks5 <-> $ brook connect <-> dialWithSocks5 <-> $ brook server/wsserver/wssserver <-> dst]",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Username",
					Usage: "Optional",
				},
				&cli.StringFlag{
					Name:  "dialWithSocks5Password",
					Usage: "Optional",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("link") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				h, p, err := net.SplitHostPort(c.String("socks5"))
				if err != nil {
					return err
				}
				if c.String("socks5ServerIP") != "" {
					h = c.String("socks5ServerIP")
				}
				if h == "" {
					return errors.New("socks5 server requires a clear IP for UDP, only port is not enough. You may use loopback IP or lan IP or other, we can not decide for you")
				}
				kind, _, _, err := brook.ParseLink(c.String("link"))
				if err != nil {
					return err
				}
				if kind == "socks5" {
					return errors.New("Looks like you want create socks5 from a socks5, you may want $ brook socks5tohttp?")
				}
				if c.String("dialWithSocks5") != "" {
					p, err := socks5dial.NewSocks5Dial(c.String("dialWithSocks5"), c.String("dialWithSocks5Username"), c.String("dialWithSocks5Password"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				s, err := brook.NewBrookLink(c.String("link"))
				if err != nil {
					return err
				}
				err = s.PrepareSocks5Server(c.String("socks5"), h, c.Int("tcpTimeout"), c.Int("udpTimeout"))
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
				if c.String("http") != "" {
					h, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(h, p), "", "", c.Int("tcpTimeout"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "relay",
			Usage: "Run as standalone relay, both TCP and UDP, this means access [from address] is equal to access [to address], [src <-> from address <-> to address]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "from",
					Aliases: []string{"f", "l"},
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
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("from") == "" || c.String("to") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				s, err := brook.NewRelay(c.String("from"), c.String("to"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "dnsserver",
			Usage: "Run as standalone dns server, both TCP and UDP",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, like: 127.0.0.1:53",
				},
				&cli.StringFlag{
					Name:  "dns",
					Usage: "DNS server which forward to",
					Value: "8.8.8.8:53",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.BoolFlag{
					Name:  "disableA",
					Usage: "Disable A query",
				},
				&cli.BoolFlag{
					Name:  "disableAAAA",
					Usage: "Disable AAAA query",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.Bool("disableA") || c.Bool("disableAAAA") {
					p, err := thedns.NewTheDNS(c.String("blockDomainList"), "", "", c.Bool("disableA"), c.Bool("disableAAAA"))
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				s, err := brook.NewRelay(c.String("listen"), c.String("dns"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				s.IsDNS = true
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				return nil
			},
		},
		&cli.Command{
			Name:  "socks5",
			Usage: "Run as standalone standard socks5 server, both TCP and UDP",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
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
				&cli.BoolFlag{
					Name:  "limitUDP",
					Usage: "The server MAY use this information to limit access to the UDP association. This usually causes connection failures in a NAT environment, where most clients are.",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 0,
					Usage: "Connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpTimeout",
					Value: 60,
					Usage: "Connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" {
					return cli.ShowSubcommandHelp(c)
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
				s.Server.LimitUDP = c.Bool("limitUDP")
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				return nil
			},
		},
		&cli.Command{
			Name:  "socks5tohttp",
			Usage: "Convert socks5 to http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
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
					Usage: "Connection tcp timeout (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("socks5") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				s, err := brook.NewSocks5ToHTTP(c.String("listen"), c.String("socks5"), c.String("socks5username"), c.String("socks5password"), c.Int("tcpTimeout"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "pac",
			Usage: "Run as PAC server or save PAC to file",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
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
					Usage:   "One domain per line, suffix match mode. http(s):// or local absolute file path. Like: https://txthinking.github.io/bypass/example_domain.txt",
				},
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "Save PAC to file, this will ignore listen address",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" && c.String("file") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("bypassDomainList") != "" && !strings.HasPrefix(c.String("bypassDomainList"), "http://") && !strings.HasPrefix(c.String("bypassDomainList"), "https://") && !filepath.IsAbs(c.String("bypassDomainList")) {
					return errors.New("--bypassDomainList must be with absolute path")
				}
				s := brook.NewPAC(c.String("listen"), c.String("file"), c.String("proxy"), c.String("bypassDomainList"))
				if c.String("file") != "" {
					return s.WriteToFile()
				}
				g.Add(&runnergroup.Runner{
					Start: func() error {
						return s.ListenAndServe()
					},
					Stop: func() error {
						return s.Shutdown()
					},
				})
				return nil
			},
		},
		&cli.Command{
			Name:  "testsocks5",
			Usage: "Test UDP and TCP of socks5 server",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "socks5",
					Aliases: []string{"s"},
					Usage:   "Like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Usage:   "Socks5 username",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Socks5 password",
				},
				&cli.StringFlag{
					Name:  "dns",
					Value: "8.8.8.8:53",
					Usage: "DNS server for connecting",
				},
				&cli.StringFlag{
					Name:  "domain",
					Value: "http3.ooo",
					Usage: "Domain for query",
				},
				&cli.StringFlag{
					Name:  "a",
					Value: "137.184.237.95",
					Usage: "The A record of domain",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("socks5") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				socks5.Debug = true
				return brook.Socks5Test(c.String("socks5"), c.String("username"), c.String("password"), c.String("domain"), c.String("a"), c.String("dns"))
			},
		},
		&cli.Command{
			Name:  "testbrook",
			Usage: "Test UDP and TCP of brook server/wsserver/wssserver",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "link",
					Aliases: []string{"l"},
					Usage:   "brook link. Get it via $ brook link",
				},
				&cli.StringFlag{
					Name:  "socks5",
					Value: "127.0.0.1:11080",
					Usage: "Temporarily listening socks5",
				},
				&cli.StringFlag{
					Name:  "dns",
					Value: "8.8.8.8:53",
					Usage: "DNS server for connecting",
				},
				&cli.StringFlag{
					Name:  "domain",
					Value: "http3.ooo",
					Usage: "Domain for query",
				},
				&cli.StringFlag{
					Name:  "a",
					Value: "137.184.237.95",
					Usage: "The A record of domain",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("link") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				socks5.Debug = true
				fmt.Println("Run brook connect to listen", c.String("socks5"))
				var cmd *exec.Cmd
				var err error
				go func() {
					var s string
					s, err = os.Executable()
					if err != nil {
						return
					}
					cmd = exec.Command(s, "connect", "--link", c.String("link"), "--socks5", c.String("socks5"))
					b, _ := cmd.CombinedOutput()
					err = errors.New(string(b))
				}()
				time.Sleep(3 * time.Second)
				if err != nil {
					return err
				}
				err1 := brook.Socks5Test(c.String("socks5"), "", "", c.String("domain"), c.String("a"), c.String("dns"))
				_ = cmd.Process.Signal(syscall.SIGTERM)
				return err1
			},
		},
		&cli.Command{
			Name:  "echoserver",
			Usage: "Echo server, echo UDP and TCP address of routes",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Listen address, like: ':7777'",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				s, err := brook.NewEchoServer(c.String("listen"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "echoclient",
			Usage: "Connect to echoserver, echo UDP and TCP address of routes",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Echo server address, such as 1.2.3.4:7777",
				},
				&cli.IntFlag{
					Name:  "times",
					Value: 1,
					Usage: "Times of interactions",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				return brook.EchoClient(c.String("server"), c.Int("times"))
			},
		},
		&cli.Command{
			Name:  "completion",
			Usage: "Generate shell completions",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "Write to file",
					Value:   "brook_autocomplete",
				},
			},
			Action: func(c *cli.Context) error {
				l := c.App.VisibleCommands()
				if strings.Contains(os.Getenv("SHELL"), "zsh") {
					s := `
_cli_zsh_autocomplete() {
  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi
}
compdef _cli_zsh_autocomplete brook
`
					for _, v := range l {
						s += "compdef _cli_zsh_autocomplete brook " + v.Name + "\n"
					}
					if err := os.WriteFile(c.String("file"), []byte(s), 0644); err != nil {
						return err
					}
					fmt.Println("Generated")
					fmt.Println("\t" + c.String("file"))
					fmt.Println("To enable auto-completion for the current shell session")
					fmt.Println("\t$ source " + c.String("file"))
					fmt.Println("To enable persistent auto-completion, add this line to your .zsh")
					fmt.Println("\t$ source /path/to/" + c.String("file"))
					return nil
				}
				s := `#! /bin/bash
_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}
complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete brook
`
				for _, v := range l {
					s += "complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete brook " + v.Name + "\n"
				}
				if err := os.WriteFile(c.String("file"), []byte(s), 0644); err != nil {
					return err
				}
				fmt.Println("Generated:")
				fmt.Println("\t" + c.String("file"))
				fmt.Println("To enable auto-completion for the current shell session")
				fmt.Println("\t$ source " + c.String("file"))
				if runtime.GOOS == "darwin" {
					fmt.Println("To enable persistent auto-completion, add this line to your .bash_profile")
				}
				if runtime.GOOS != "darwin" {
					fmt.Println("To enable persistent auto-completion, add this line to your .bashrc")
				}
				fmt.Println("\t$ source /path/to/" + c.String("file"))
				return nil
			},
		},
		&cli.Command{
			Name:  "mdpage",
			Usage: "Generate markdown page",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "Write to file, default print to stdout",
				},
			},
			Action: func(c *cli.Context) error {
				s, err := c.App.ToMarkdown()
				if err != nil {
					return err
				}
				if c.String("file") != "" {
					if err := os.WriteFile(c.String("file"), []byte(s), 0644); err != nil {
						return err
					}
					return nil
				}
				fmt.Println(s)
				return nil
			},
		},
		&cli.Command{
			Name:  "manpage",
			Usage: "Generate man.1 page",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "Write to file, default print to stdout. You should put to /path/to/man/man1/brook.1 on linux or /usr/local/share/man/man1/brook.1 on macos",
				},
			},
			Action: func(c *cli.Context) error {
				s, err := c.App.ToManWithSection(1)
				if err != nil {
					return err
				}
				if c.String("file") != "" {
					if err := os.WriteFile(c.String("file"), []byte(s), 0644); err != nil {
						return err
					}
					return nil
				}
				fmt.Println(s)
				return nil
			},
		},
	}
	if os.Getenv("SOCKS5_DEBUG") != "" {
		socks5.Debug = true
	}
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
		return
	}
	if len(g.Runners) == 0 {
		return
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		g.Done()
	}()
	log.Println(g.Wait())
}
