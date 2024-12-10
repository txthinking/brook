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
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"net/url"

	"github.com/miekg/dns"
	"github.com/phuslu/iploc"
	"github.com/txthinking/brook"
	"github.com/txthinking/brook/plugins/block"
	"github.com/txthinking/brook/plugins/dialwithdns"
	"github.com/txthinking/brook/plugins/dialwithip"
	"github.com/txthinking/brook/plugins/dialwithnic"
	"github.com/txthinking/brook/plugins/logger"
	"github.com/txthinking/brook/plugins/pprof"
	"github.com/txthinking/brook/plugins/prometheus"
	"github.com/txthinking/brook/plugins/socks5dial"
	"github.com/txthinking/brook/plugins/thedns"
	"github.com/txthinking/runnergroup"
	"github.com/txthinking/socks5"
	"github.com/urfave/cli/v2"
)

func main() {
	g := runnergroup.New()
	df := func() {}
	app := cli.NewApp()
	app.Name = "Brook"
	app.Version = "20250202"
	app.Usage = "A cross-platform programmable network tool"
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
		&cli.StringFlag{
			Name:  "log",
			Usage: "Enable log. A valid value is file path or 'console'. Send SIGUSR1 to me to reset the log file on unix system. If you want to debug SOCKS5 lib, set env SOCKS5_DEBUG=true",
		},
		&cli.StringSliceFlag{
			Name:  "tag",
			Usage: "Tag can be used to the process, will be append into log or serverLog, such as: 'key1:value1'. All tags will also be appended as query parameters one by one to the userAPI",
		},
		&cli.StringFlag{
			Name:  "dialWithDNS",
			Usage: "When a domain name needs to be resolved, use the specified DNS. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required. Note that for client-side commands, this does not affect the client passing the domain address to the server",
		},
		&cli.StringFlag{
			Name:  "dialWithDNSPrefer",
			Usage: "This is used with the dialWithDNS parameter. Prefer A record or AAAA record. Value is A or AAAA",
		},
		&cli.StringFlag{
			Name:  "dialWithIP4",
			Usage: "When the current machine establishes a network connection to the outside IPv4, both TCP and UDP, it is used to specify the IPv4 used",
		},
		&cli.StringFlag{
			Name:  "dialWithIP6",
			Usage: "When the current machine establishes a network connection to the outside IPv6, both TCP and UDP, it is used to specify the IPv6 used",
		},
		&cli.StringFlag{
			Name:  "dialWithNIC",
			Usage: "When the current machine establishes a network connection to the outside, both TCP and UDP, it is used to specify the NIC used",
		},
		&cli.StringFlag{
			Name:  "dialWithSocks5",
			Usage: "When the current machine establishes a network connection to the outside, both TCP and UDP, with your socks5 proxy, such as 127.0.0.1:1081",
		},
		&cli.StringFlag{
			Name:  "dialWithSocks5Username",
			Usage: "If there is",
		},
		&cli.StringFlag{
			Name:  "dialWithSocks5Password",
			Usage: "If there is",
		},
		&cli.IntFlag{
			Name:  "dialWithSocks5TCPTimeout",
			Value: 0,
			Usage: "time (s)",
		},
		&cli.IntFlag{
			Name:  "dialWithSocks5UDPTimeout",
			Value: 60,
			Usage: "time (s)",
		},
		&cli.StringFlag{
			Name:  "prometheus",
			Usage: "prometheus http listen addr, such as :7070. If it is transmitted on the public network, it is recommended to use it with nico",
		},
		&cli.StringFlag{
			Name:  "prometheusPath",
			Usage: "prometheus http path, such as /xxx. If it is transmitted on the public network, a hard-to-guess value is recommended",
		},
		&cli.StringFlag{
			Name:  "clientHKDFInfo",
			Usage: "client HKDF info, most time you don't need to change this, if changed, all and each brook links in client side must be same, I mean each",
			Value: "brook",
		},
		&cli.StringFlag{
			Name:  "serverHKDFInfo",
			Usage: "server HKDF info, most time you don't need to change this, if changed, all and each brook links in client side must be same, I mean each",
			Value: "brook",
		},
		&cli.StringFlag{
			Name:  "blockDomainList",
			Usage: "One domain per line, suffix match mode. https://, http:// or local file absolute path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt. Works with server/wsserver/wssserver/quicserver",
		},
		&cli.StringFlag{
			Name:  "blockCIDR4List",
			Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_cidr4.txt. Works with server/wsserver/wssserver/quicserver",
		},
		&cli.StringFlag{
			Name:  "blockCIDR6List",
			Usage: "One CIDR per line, https://, http:// or local file absolute path, like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_cidr6.txt. Works with server/wsserver/wssserver/quicserver",
		},
		&cli.StringSliceFlag{
			Name:  "blockGeoIP",
			Usage: "Block IP by Geo country code, such as US. Works with server/wsserver/wssserver/quicserver",
		},
		&cli.Int64Flag{
			Name:  "blockListUpdateInterval",
			Usage: "Update list --blockDomainList,--blockCIDR4List,--blockCIDR6List interval, second. default 0, only read one time on start",
		},
		&cli.StringFlag{
			Name:  "pid",
			Usage: "A file path used to store pid. Send SIGUSR1 to me to reset the --serverLog file on unix system",
		},
	}
	app.Before = func(c *cli.Context) error {
		brook.ClientHKDFInfo = []byte(c.String("clientHKDFInfo"))
		brook.ServerHKDFInfo = []byte(c.String("serverHKDFInfo"))
		if c.String("pid") != "" {
			if !filepath.IsAbs(c.String("pid")) {
				return errors.New("--pid must be with absolute path")
			}
			if err := os.WriteFile(c.String("pid"), []byte(strconv.Itoa(os.Getpid())), 0744); err != nil {
				return err
			}
		}
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
		if c.String("log") != "" {
			if c.String("log") != "console" && !filepath.IsAbs(c.String("log")) {
				return errors.New("--log must be with absolute path")
			}
			var m map[string]string
			if len(c.StringSlice("tag")) > 0 {
				m = make(map[string]string)
				for _, v := range c.StringSlice("tag") {
					l := strings.Split(v, ":")
					if len(l) != 2 {
						return errors.New("Invalid tag " + v)
					}
					m[l[0]] = l[1]
				}
			}
			p, err := logger.NewLogger(m, c.String("log"))
			if err != nil {
				return err
			}
			p.TouchBrook()
			f := df
			df = func() {
				p.Close()
				f()
			}
		}
		if c.String("dialWithDNS") != "" {
			p, err := dialwithdns.NewDialWithDNS(c.String("dialWithDNS"), c.String("dialWithDNSPrefer"))
			if err != nil {
				return err
			}
			p.TouchBrook()
		}
		if c.String("dialWithIP4") != "" || c.String("dialWithIP6") != "" {
			p, err := dialwithip.NewDialWithIP(c.String("dialWithIP4"), c.String("dialWithIP6"))
			if err != nil {
				return err
			}
			p.TouchBrook()
		}
		if c.String("dialWithNIC") != "" {
			p := dialwithnic.NewDialWithNIC(c.String("dialWithNIC"))
			p.TouchBrook()
		}
		if c.String("dialWithSocks5") != "" {
			p, err := socks5dial.NewSocks5Dial(c.String("dialWithSocks5"), c.String("dialWithSocks5Username"), c.String("dialWithSocks5Password"), c.Int("dialWithSocks5TCPTimeout"), c.Int("dialWithSocks5UDPTimeout"))
			if err != nil {
				return err
			}
			p.TouchBrook()
		}
		if c.String("prometheus") != "" {
			if c.String("prometheusPath") == "" {
				return errors.New("You forgot the --prometheusPath")
			}
			var m map[string]string
			if len(c.StringSlice("tag")) > 0 {
				m = make(map[string]string)
				for _, v := range c.StringSlice("tag") {
					l := strings.Split(v, ":")
					if len(l) != 2 {
						return errors.New("Invalid tag " + v)
					}
					m[l[0]] = l[1]
				}
			}
			p := prometheus.NewPrometheus(c.String("prometheus"), c.String("prometheusPath"), m)
			p.TouchBrook()
			g.Add(&runnergroup.Runner{
				Start: func() error {
					return p.ListenAndServe()
				},
				Stop: func() error {
					return p.Shutdown()
				},
			})
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
			p, err := block.NewBlock(c.String("blockDomainList"), c.String("blockCIDR4List"), c.String("blockCIDR6List"), c.StringSlice("blockGeoIP"), c.Int("blockListUpdateInterval"))
			if err != nil {
				return err
			}
			p.TouchBrook()
			if c.Int("blockListUpdateInterval") != 0 {
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
		return nil
	}
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:  "server",
			Usage: "Start a brook server that supports tcp and udp",
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
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringSliceFlag{
					Name:  "blockGeoIP",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook server --listen :9999 --password hello")
					return nil
				}
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
			Usage: "Start a brook client that supports tcp and udp. It can open a socks5 proxy, [src <-> socks5 <-> $ brook client <-> $ brook server <-> dst]",
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
				&cli.StringFlag{
					Name:  "link",
					Usage: "brook link, you can get it via $ brook link. The wssserver and password parameters will be ignored",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook client --server 1.2.3.4:9999 --password hello --socks5 127.0.0.1:1080")
					return nil
				}
				if c.String("server") == "" && c.String("link") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				var link = ""
				if c.String("server") != "" {
					v := url.Values{}
					v.Set("password", c.String("password"))
					link = brook.Link("server", c.String("server"), v)
				}
				if c.String("link") != "" {
					link = c.String("link")
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
				kind, _, v, err := brook.ParseLink(link)
				if err != nil {
					return err
				}
				if s := v.Get("clientHKDFInfo"); s != "" {
					brook.ClientHKDFInfo = []byte(s)
				}
				if s := v.Get("serverHKDFInfo"); s != "" {
					brook.ServerHKDFInfo = []byte(s)
				}
				if kind == "socks5" {
					return errors.New("Looks like you want create socks5 from a socks5, you may want $ brook socks5tohttp?")
				}
				s, err := brook.NewBrookLink(link)
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
					s, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(h, p), "", "", c.Int("tcpTimeout"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "wsserver",
			Usage: "Start a brook wsserver that supports tcp and udp. It opens a standard http server and a websocket server",
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
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringSliceFlag{
					Name:  "blockGeoIP",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.BoolFlag{
					Name:  "xForwardedFor",
					Usage: "Replace the from field in --log, note that this may be forged",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook wsserver --listen :9999 --password hello")
					return nil
				}
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
				s, err := brook.NewWSServer(c.String("listen"), c.String("password"), "", c.String("path"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.Bool("withoutBrookProtocol"))
				if err != nil {
					return err
				}
				s.XForwardedFor = c.Bool("xForwardedFor")
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
			Usage: "Start a brook wsclient that supports tcp and udp. It can open a socks5 proxy, [src <-> socks5 <-> $ brook wsclient <-> $ brook wsserver <-> dst]",
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
					Name:  "link",
					Usage: "brook link, you can get it via $ brook link. The wssserver and password parameters will be ignored",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook wsclient --wsserver ws://1.2.3.4:9999 --password hello --socks5 127.0.0.1:1080")
					return nil
				}
				if c.String("wsserver") == "" && c.String("link") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				var link = ""
				if c.String("wsserver") != "" {
					v := url.Values{}
					v.Set("password", c.String("password"))
					link = brook.Link("wsserver", c.String("wsserver"), v)
				}
				if c.String("link") != "" {
					link = c.String("link")
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
				kind, _, v, err := brook.ParseLink(link)
				if err != nil {
					return err
				}
				if s := v.Get("clientHKDFInfo"); s != "" {
					brook.ClientHKDFInfo = []byte(s)
				}
				if s := v.Get("serverHKDFInfo"); s != "" {
					brook.ServerHKDFInfo = []byte(s)
				}
				if kind == "socks5" {
					return errors.New("Looks like you want create socks5 from a socks5, you may want $ brook socks5tohttp?")
				}
				s, err := brook.NewBrookLink(link)
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
					s, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(h, p), "", "", c.Int("tcpTimeout"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "wssserver",
			Usage: "Start a brook wssserver that supports tcp and udp. It opens a standard https server and a websocket server",
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
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringSliceFlag{
					Name:  "blockGeoIP",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook wssserver --domainaddress domain.com:9999 --password hello")
					return nil
				}
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
				s, err := brook.NewWSServer(":"+p, c.String("password"), h, c.String("path"), c.Int("tcpTimeout"), c.Int("udpTimeout"), c.Bool("withoutBrookProtocol"))
				if err != nil {
					return err
				}
				if c.String("cert") != "" {
					b, err := os.ReadFile(c.String("cert"))
					if err != nil {
						return err
					}
					s.Cert = b
				}
				if c.String("certkey") != "" {
					b, err := os.ReadFile(c.String("certkey"))
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
			Usage: "Start a brook wssclient that supports tcp and udp. It can open a socks5 proxy, [src <-> socks5 <-> $ brook wssclient <-> $ brook wssserver <-> dst]",
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
					Name:  "link",
					Usage: "brook link, you can get it via $ brook link. The wssserver and password parameters will be ignored",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook wssclient --wssserver wss://domain.com:9999 --password hello --socks5 127.0.0.1:1080")
					return nil
				}
				if c.String("wssserver") == "" && c.String("link") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				var link = ""
				if c.String("wssserver") != "" {
					v := url.Values{}
					v.Set("password", c.String("password"))
					link = brook.Link("wssserver", c.String("wssserver"), v)
				}
				if c.String("link") != "" {
					link = c.String("link")
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
				kind, _, v, err := brook.ParseLink(link)
				if err != nil {
					return err
				}
				if s := v.Get("clientHKDFInfo"); s != "" {
					brook.ClientHKDFInfo = []byte(s)
				}
				if s := v.Get("serverHKDFInfo"); s != "" {
					brook.ServerHKDFInfo = []byte(s)
				}
				if kind == "socks5" {
					return errors.New("Looks like you want create socks5 from a socks5, you may want $ brook socks5tohttp?")
				}
				s, err := brook.NewBrookLink(link)
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
					s, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(h, p), "", "", c.Int("tcpTimeout"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "quicserver",
			Usage: "Start a brook quicserver that supports tcp and udp.",
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
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringFlag{
					Name:  "blockCIDR4List",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringFlag{
					Name:  "blockCIDR6List",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.StringSliceFlag{
					Name:  "blockGeoIP",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.Int64Flag{
					Name:  "updateListInterval",
					Usage: "This option will be removed in a future version, please use the global option instead",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook quicserver --domainaddress domain.com:9999 --password hello")
					return nil
				}
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
					b, err := os.ReadFile(c.String("cert"))
					if err != nil {
						return err
					}
					s.Cert = b
				}
				if c.String("certkey") != "" {
					b, err := os.ReadFile(c.String("certkey"))
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
			Usage: "Start a brook quicclient that supports tcp and udp. It can open a socks5 proxy, [src <-> socks5 <-> $ brook quicclient <-> $ brook quicserver <-> dst]. (The global-dial-parameter is ignored)",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "link",
					Usage: "brook link, you can get it via $ brook link. The wssserver and password parameters will be ignored",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook quicclient --link 'brook://...' --socks5 127.0.0.1:1080")
					return nil
				}
				if c.String("link") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				link := c.String("link")
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
				kind, _, v, err := brook.ParseLink(link)
				if err != nil {
					return err
				}
				if s := v.Get("clientHKDFInfo"); s != "" {
					brook.ClientHKDFInfo = []byte(s)
				}
				if s := v.Get("serverHKDFInfo"); s != "" {
					brook.ServerHKDFInfo = []byte(s)
				}
				if kind == "socks5" {
					return errors.New("Looks like you want create socks5 from a socks5, you may want $ brook socks5tohttp?")
				}
				s, err := brook.NewBrookLink(link)
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
					s, err := brook.NewSocks5ToHTTP(c.String("http"), net.JoinHostPort(h, p), "", "", c.Int("tcpTimeout"))
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
				return nil
			},
		},
		&cli.Command{
			Name:  "relayoverbrook",
			Usage: "Relay network traffic over brook, which supports TCP and UDP. Accessing [from address] is equal to accessing [to address], [src <-> from address <-> $ brook server/wsserver/wssserver/quicserver <-> to address]",
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
					Usage:   "brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain:443/ws, quic://domain.com:443",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.StringFlag{
					Name:  "link",
					Usage: "brook link, you can get it via $ brook link. The server and password parameters will be ignored",
				},
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook relayoverbrook --server 1.2.3.4:9999 --password hello --from :5353 --to 8.8.8.8:53")
					return nil
				}
				if c.String("from") == "" || c.String("to") == "" || (c.String("server") == "" && c.String("link") == "") {
					return cli.ShowSubcommandHelp(c)
				}
				var link = ""
				if c.String("server") != "" {
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
					if kind == "quicserver" {
						return errors.New("It is recommended to use brook link and specify --udpoverstream")
					}
					v := url.Values{}
					v.Set("password", c.String("password"))
					link = brook.Link(kind, c.String("server"), v)
				}
				if c.String("link") != "" {
					link = c.String("link")
				}
				kind, _, v, err := brook.ParseLink(link)
				if err != nil {
					return err
				}
				if s := v.Get("clientHKDFInfo"); s != "" {
					brook.ClientHKDFInfo = []byte(s)
				}
				if s := v.Get("serverHKDFInfo"); s != "" {
					brook.ServerHKDFInfo = []byte(s)
				}
				if kind == "socks5" {
					return errors.New("Looks like you want create socks5 from a socks5, you may want $ brook socks5tohttp?")
				}
				s, err := brook.NewRelayOverBrook(c.String("from"), link, c.String("to"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
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
			Usage: "Run a dns server over brook, which supports TCP and UDP, [src <-> $ brook dnserversoverbrook <-> $ brook server/wsserver/wssserver/quicserver <-> dns] or [src <-> $ brook dnsserveroverbrook <-> dnsForBypass]",
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
					Usage:   "brook server or brook wsserver or brook wssserver or brook quicserver, like: 1.2.3.4:9999, ws://1.2.3.4:9999, wss://domain.com:443/ws, quic://domain.com:443",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Password",
				},
				&cli.StringFlag{
					Name:  "link",
					Usage: "brook link, you can get it via $ brook link. The server and password parameters will be ignored",
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
					Usage: "DNS server for resolving domains in bypass list. Such as 223.5.5.5:53 or https://dns.alidns.com/dns-query?address=223.5.5.5:443, the address is required",
					Value: "223.5.5.5:53",
				},
				&cli.StringFlag{
					Name:  "bypassDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook dnsserveroverbrook --server 1.2.3.4:9999 --password hello --listen 127.0.0.1:53")
					return nil
				}
				if c.String("listen") == "" || (c.String("server") == "" && c.String("link") == "") {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("bypassDomainList") != "" && !strings.HasPrefix(c.String("bypassDomainList"), "http://") && !strings.HasPrefix(c.String("bypassDomainList"), "https://") && !filepath.IsAbs(c.String("bypassDomainList")) {
					return errors.New("--bypassDomainList must be with absolute path")
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.String("bypassDomainList") != "" || c.Bool("disableA") || c.Bool("disableAAAA") {
					p, err := thedns.NewTheDNS(c.String("blockDomainList"), c.String("bypassDomainList"), c.String("dnsForBypass"), c.Bool("disableA"), c.Bool("disableAAAA"), "")
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				var link = ""
				if c.String("server") != "" {
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
					if kind == "quicserver" {
						return errors.New("It is recommended to use brook link and specify --udpoverstream")
					}
					v := url.Values{}
					v.Set("password", c.String("password"))
					link = brook.Link(kind, c.String("server"), v)
				}
				if c.String("link") != "" {
					link = c.String("link")
				}
				kind, _, v, err := brook.ParseLink(link)
				if err != nil {
					return err
				}
				if s := v.Get("clientHKDFInfo"); s != "" {
					brook.ClientHKDFInfo = []byte(s)
				}
				if s := v.Get("serverHKDFInfo"); s != "" {
					brook.ServerHKDFInfo = []byte(s)
				}
				if kind == "socks5" {
					return errors.New("Looks like you want create socks5 from a socks5, you may want $ brook socks5tohttp?")
				}
				s, err := brook.NewRelayOverBrook(c.String("listen"), link, c.String("dns"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
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
			Name:  "link",
			Usage: "Generate a brook link",
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
				&cli.BoolFlag{
					Name:  "udpoverstream",
					Usage: "When server is brook quicserver, UDP over Stream. Under normal circumstances, you need this parameter because the max datagram size for QUIC is very small. Note: only brook CLI and tun2brook suppport for now",
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
				&cli.StringFlag{
					Name:  "tlsfingerprint",
					Usage: "When server is brook wssserver, select tls fingerprint, value can be: chrome",
				},
				&cli.StringFlag{
					Name:  "clientHKDFInfo",
					Usage: "client HKDF info, most time you don't need to change this, read brook protocol if you don't know what this is",
				},
				&cli.StringFlag{
					Name:  "serverHKDFInfo",
					Usage: "server HKDF info, most time you don't need to change this, read brook protocol if you don't know what this is",
				},
				&cli.StringFlag{
					Name:  "fragment",
					Usage: "When server is brook wssserver, split the ClientHello into multiple fragments and then send them one by one with delays (millisecond). The format is min_length:max_length:min_delay:max_delay, cannot be zero, such as 50:100:10:50",
				},
				&cli.StringFlag{
					Name:  "token",
					Usage: "A token represents a user's identity. A string encoded in hexadecimal. Server needs to have --userAPI enabled. Note that: Only supported by the brook GUI(except for OpenWrt) and tun2brook",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook link --server 1.2.3.4:9999 --password hello --name 'my server'")
					return nil
				}
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
				if c.Bool("udpoverstream") {
					v.Set("udpoverstream", "true")
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
				if c.String("tlsfingerprint") != "" {
					v.Set("tlsfingerprint", c.String("tlsfingerprint"))
				}
				if c.String("ca") != "" {
					b, err := os.ReadFile(c.String("ca"))
					if err != nil {
						return err
					}
					v.Set("ca", string(b))
				}
				if c.String("clientHKDFInfo") != "" {
					v.Set("clientHKDFInfo", c.String("clientHKDFInfo"))
				}
				if c.String("serverHKDFInfo") != "" {
					v.Set("serverHKDFInfo", c.String("serverHKDFInfo"))
				}
				if c.String("fragment") != "" {
					v.Set("fragment", c.String("fragment"))
				}
				if c.String("token") != "" {
					v.Set("token", c.String("token"))
				}
				fmt.Println(brook.Link(s, c.String("server"), v))
				return nil
			},
		},
		&cli.Command{
			Name:  "connect",
			Usage: "Run a client and connect with a brook link, which supports TCP and UDP. It can start a socks5 proxy, [src <-> socks5 <-> $ brook connect <-> $ brook server/wsserver/wssserver/quicserver <-> dst]",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook connect --link 'brook://...' --socks5 127.0.0.1:1080")
					return nil
				}
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
				kind, _, v, err := brook.ParseLink(c.String("link"))
				if err != nil {
					return err
				}
				if s := v.Get("clientHKDFInfo"); s != "" {
					brook.ClientHKDFInfo = []byte(s)
				}
				if s := v.Get("serverHKDFInfo"); s != "" {
					brook.ServerHKDFInfo = []byte(s)
				}
				if kind == "socks5" {
					return errors.New("Looks like you want create socks5 from a socks5, you may want $ brook socks5tohttp?")
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
			Usage: "Run a standalone relay, which supports TCP and UDP. Accessing [from address] is equal to accessing [to address], [src <-> from address <-> to address]",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook relay --from :9999 --to 1.2.3.4:9999")
					return nil
				}
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
			Usage: "Run a standalone dns server",
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
					Usage: "DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required",
					Value: "8.8.8.8:53",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook dnsserver --listen 127.0.0.1:53")
					return nil
				}
				if c.String("listen") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				var doh string
				if strings.HasPrefix(c.String("dns"), "https://") {
					doh = c.String("dns")
				}
				if c.String("blockDomainList") != "" || c.Bool("disableA") || c.Bool("disableAAAA") || doh != "" {
					p, err := thedns.NewTheDNS(c.String("blockDomainList"), "", "", c.Bool("disableA"), c.Bool("disableAAAA"), doh)
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
			Name:  "dnsclient",
			Usage: "Send a dns query",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "dns",
					Aliases: []string{"s"},
					Usage:   "DNS server, such as 8.8.8.8:53",
					Value:   "8.8.8.8:53",
				},
				&cli.StringFlag{
					Name:    "domain",
					Aliases: []string{"d"},
					Usage:   "Domain",
				},
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Usage:   "Type, such as A",
					Value:   "A",
				},
				&cli.BoolFlag{
					Name:  "short",
					Usage: "Short for A/AAAA",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook dnsclient --domain google.com")
					return nil
				}
				if c.String("domain") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				t, ok := dns.StringToType[c.String("type")]
				if !ok {
					return errors.New("invalid type")
				}
				dc := &brook.DNSClient{Server: c.String("dns")}
				m := &dns.Msg{}
				m.SetQuestion(strings.TrimRight(c.String("domain"), ".")+".", t)
				m, err := dc.Exchange(m)
				if err != nil {
					return err
				}
				if c.Bool("short") && (c.String("type") == "A" || c.String("type") == "AAAA") {
					for _, v := range m.Answer {
						if t, ok := v.(*dns.A); ok {
							fmt.Println(t.A)
							return nil
						}
						if t, ok := v.(*dns.AAAA); ok {
							fmt.Println(t.AAAA)
							return nil
						}
					}
					return nil
				}
				fmt.Println(m)
				return nil
			},
		},
		&cli.Command{
			Name:  "dohserver",
			Usage: "Run a standalone doh server",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "listen",
					Usage: "listen address, if you want to create a http server behind nico",
				},
				&cli.StringFlag{
					Name:  "domainaddress",
					Usage: "Such as: domain.com:443, if you want to create a https server. If you choose to automatically issue certificates, the domain must have been resolved to the server IP and 80 port also will be used",
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
					Value: "/dns-query",
				},
				&cli.StringFlag{
					Name:  "dns",
					Usage: "DNS server which forward to. Such as 8.8.8.8:53 or https://dns.google/dns-query?address=8.8.8.8%3A443, the address is required",
					Value: "8.8.8.8:53",
				},
				&cli.StringFlag{
					Name:  "blockDomainList",
					Usage: "One domain per line, suffix match mode. https://, http:// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook dohserver --domainaddress domain.com:443")
					return nil
				}
				if c.String("listen") == "" && c.String("domainaddress") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("cert") != "" && !filepath.IsAbs(c.String("cert")) {
					return errors.New("--cert must be with absolute path")
				}
				if c.String("certkey") != "" && !filepath.IsAbs(c.String("certkey")) {
					return errors.New("--certkey must be with absolute path")
				}
				if c.String("blockDomainList") != "" && !strings.HasPrefix(c.String("blockDomainList"), "http://") && !strings.HasPrefix(c.String("blockDomainList"), "https://") && !filepath.IsAbs(c.String("blockDomainList")) {
					return errors.New("--blockDomainList must be with absolute path")
				}
				if c.String("blockDomainList") != "" || c.Bool("disableA") || c.Bool("disableAAAA") {
					p, err := thedns.NewTheDNS(c.String("blockDomainList"), "", "", c.Bool("disableA"), c.Bool("disableAAAA"), "")
					if err != nil {
						return err
					}
					p.TouchBrook()
				}
				addr := ""
				domain := ""
				if c.String("domainaddress") != "" {
					h, p, err := net.SplitHostPort(c.String("domainaddress"))
					if err != nil {
						return err
					}
					domain = h
					addr = ":" + p
				}
				if c.String("listen") != "" {
					addr = c.String("listen")
				}
				s, err := brook.NewDOHServer(addr, domain, c.String("path"), c.String("dns"), c.Int("tcpTimeout"), c.Int("udpTimeout"))
				if err != nil {
					return err
				}
				if c.String("cert") != "" {
					b, err := os.ReadFile(c.String("cert"))
					if err != nil {
						return err
					}
					s.Cert = b
				}
				if c.String("certkey") != "" {
					b, err := os.ReadFile(c.String("certkey"))
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
			Name:  "dohclient",
			Usage: "Send a dns query",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "doh",
					Aliases: []string{"s"},
					Usage:   "DOH server, the address is required",
					Value:   "https://dns.quad9.net/dns-query?address=9.9.9.9%3A443",
				},
				&cli.StringFlag{
					Name:    "domain",
					Aliases: []string{"d"},
					Usage:   "Domain",
				},
				&cli.StringFlag{
					Name:    "type",
					Aliases: []string{"t"},
					Usage:   "Type, such as A",
					Value:   "A",
				},
				&cli.BoolFlag{
					Name:  "short",
					Usage: "Short for A/AAAA",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook dohclient --domain google.com")
					return nil
				}
				if c.String("domain") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				t, ok := dns.StringToType[c.String("type")]
				if !ok {
					return errors.New("invalid type")
				}
				dc, err := brook.NewDOHClient(c.String("doh"))
				if err != nil {
					return err
				}
				m := &dns.Msg{}
				m.SetQuestion(strings.TrimRight(c.String("domain"), ".")+".", t)
				m, err = dc.Exchange(m)
				if err != nil {
					return err
				}
				if c.Bool("short") && (c.String("type") == "A" || c.String("type") == "AAAA") {
					for _, v := range m.Answer {
						if t, ok := v.(*dns.A); ok {
							fmt.Println(t.A)
							return nil
						}
						if t, ok := v.(*dns.AAAA); ok {
							fmt.Println(t.AAAA)
							return nil
						}
					}
					return nil
				}
				fmt.Println(m)
				return nil
			},
		},
		&cli.Command{
			Name:  "dhcpserver",
			Usage: "Run a standalone dhcp server. Other running dhcp servers need to be stopped.",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "interface",
					Usage: "Select interface on multi interface device. Linux only",
				},
				&cli.StringFlag{
					Name:  "serverip",
					Usage: "DHCP server IP, the IP of the this machine, you shoud set a static IP to this machine before doing this, such as: 192.168.1.10",
				},
				&cli.StringFlag{
					Name:  "start",
					Usage: "Start IP which you want to assign to clients, such as: 192.168.1.100",
				},
				&cli.StringFlag{
					Name:  "netmask",
					Usage: "Subnet netmask which you want to assign to clients",
					Value: "255.255.255.0",
				},
				&cli.IntFlag{
					Name:  "count",
					Usage: "IP range from the start, which you want to assign to clients",
					Value: 100,
				},
				&cli.StringFlag{
					Name:  "gateway",
					Usage: "The router gateway which you want to assign to clients, such as: 192.168.1.1",
				},
				&cli.StringSliceFlag{
					Name:  "dnsserver",
					Usage: "The dns server which you want to assign to clients, such as: 192.168.1.1 or 8.8.8.8",
				},
				&cli.StringFlag{
					Name:  "cache",
					Usage: "Cache file, local absolute file path, default is $HOME/.brook.dhcpserver",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook dhcpserver --serverip 192.168.1.10 --start 192.168.1.100 --gateway 192.168.1.1 --dnsserver 192.168.1.1")
					return nil
				}
				if c.String("serverip") == "" || c.String("start") == "" || c.String("gateway") == "" || len(c.StringSlice("dnsserver")) == 0 {
					return cli.ShowSubcommandHelp(c)
				}
				if c.String("cache") != "" && !filepath.IsAbs(c.String("cache")) {
					return errors.New("--cache must be with absolute path")
				}
				s, err := brook.NewDHCPServer(c.String("interface"), c.String("serverip"), c.String("start"), c.String("netmask"), c.Int("count"), c.String("gateway"), c.StringSlice("dnsserver"), c.String("cache"))
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
			Name:  "socks5",
			Usage: "Run a standalone standard socks5 server, which supports TCP and UDP",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook socks5 --listen :9999 --socks5ServerIP 1.2.3.4")
					return nil
				}
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
			Usage: "Convert a socks5 proxy to a http proxy, [src <-> listen address(http proxy) <-> socks5 address <-> dst]",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook socks5tohttp --listen 127.0.0.1:8010 --socks5 127.0.0.1:1080")
					return nil
				}
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
			Usage: "Run a PAC server or save PAC to a file",
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
					Usage:   "One domain per line, suffix match mode. http(s):// or local absolute file path. Like: https://raw.githubusercontent.com/txthinking/brook/master/programmable/list/example_domain.txt",
				},
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "Save PAC to file, this will ignore listen address",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook pac --listen 127.0.0.1:1980")
					return nil
				}
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
			Usage: "Test a socks5 server to see if it works properly",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook testsocks5 --socks5 127.0.0.1:1080")
					return nil
				}
				if c.String("socks5") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				socks5.Debug = true
				return brook.Socks5Test(c.String("socks5"), c.String("username"), c.String("password"), c.String("domain"), c.String("a"), c.String("dns"))
			},
		},
		&cli.Command{
			Name:  "testbrook",
			Usage: "Test UDP and TCP of a brook server/wsserver/wssserver/quicserver connection. (The global-dial-parameter is ignored)",
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook testbrook --link 'brook://...'")
					return nil
				}
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
					// TODO append global dial options
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook echoserver --listen :9999")
					return nil
				}
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook echoclient --server 1.2.3.4:9999")
					return nil
				}
				if c.String("server") == "" {
					return cli.ShowSubcommandHelp(c)
				}
				return brook.EchoClient(c.String("server"), c.Int("times"))
			},
		},
		&cli.Command{
			Name:  "ipcountry",
			Usage: "Get country of IP",
			BashComplete: func(c *cli.Context) {
				l := c.Command.VisibleFlags()
				for _, v := range l {
					fmt.Println("--" + v.Names()[0])
				}
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "ip",
					Usage: "1.1.1.1",
				},
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook ipcountry --ip 1.2.3.4")
					return nil
				}
				ip := net.ParseIP(c.String("ip"))
				if ip == nil {
					return errors.New(c.String("ip") + " is not IP")
				}
				s := iploc.Country(ip)
				if s == "" {
					return errors.New(c.String("ip") + " unknown")
				}
				fmt.Println(s)
				return nil
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook completion --file brook_autocomplete")
					return nil
				}
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook mdpage")
					return nil
				}
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
				&cli.BoolFlag{
					Name:  "example",
					Usage: "Show a minimal example of usage",
				},
			},
			Action: func(c *cli.Context) error {
				if c.Bool("example") {
					fmt.Println("brook manpage")
					return nil
				}
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
	l := os.Args
	if len(os.Args) == 2 {
		if _, err := os.Stat(os.Args[1]); err == nil {
			if !filepath.IsAbs(os.Args[1]) {
				log.Println("It looks like you passed config file? It must be absolute path")
				os.Exit(1)
				return
			}
			l, err = brook.CAC(os.Args[1])
			if err != nil {
				log.Println(err)
				os.Exit(1)
				return
			}
			bin, err := os.Executable()
			if err != nil {
				log.Println(err)
				os.Exit(1)
				return
			}
			l = append([]string{bin}, l...)
		}
	}
	if err := app.Run(l); err != nil {
		log.Println(err)
		df()
		os.Exit(1)
		return
	}
	if len(g.Runners) == 0 {
		df()
		return
	}
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		g.Done()
	}()
	if err := g.Wait(); err != nil {
		log.Println(err)
		df()
		os.Exit(1)
		return
	}
	df()
}
