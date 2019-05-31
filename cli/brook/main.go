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
	"os"
	"strings"

	"net/http"
	_ "net/http/pprof"

	"github.com/txthinking/brook"
	"github.com/urfave/cli"
)

var debug bool
var debugAddress string

func main() {
	app := cli.NewApp()
	app.Name = "Brook"
	app.Version = "20190601"
	app.Usage = "A Cross-Platform Proxy/VPN Software"
	app.Author = "Cloud"
	app.Email = "cloud@txthinking.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug",
			Destination: &debug,
		},
		cli.StringFlag{
			Name:        "listen, l",
			Usage:       "Listen address for debug",
			Value:       ":6060",
			Destination: &debugAddress,
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name:  "server",
			Usage: "Run as server mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Server listen address, like: 0.0.0.0:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
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
				return brook.RunServer(c.String("listen"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
			},
		},
		cli.Command{
			Name:  "servers",
			Usage: "Run as multiple servers mode",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "listenpassword, l",
					Usage: "server and password, like '0.0.0.0:1080 password'",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
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
				errch := make(chan error)
				go func() {
					for _, v := range c.StringSlice("listenpassword") {
						ss := strings.Split(v, " ")
						if len(ss) != 2 {
							errch <- errors.New("Invalid listenpassword")
							return
						}
						go func() {
							errch <- brook.RunServer(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
						}()
					}
				}()
				return <-errch
			},
		},
		cli.Command{
			Name:  "client",
			Usage: "Run as client mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address, like: 127.0.0.1:1080",
				},
				cli.StringFlag{
					Name:  "ip, i",
					Usage: "Client IP address, like: 127.0.0.1",
				},
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpSessionTime",
					Value: 60,
					Usage: "udp session time (s), in most cases need this",
				},
				cli.BoolFlag{
					Name:  "http",
					Usage: "If true, client start a http(s) proxy. default socks5",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("ip") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "client")
					return nil
				}
				if debug {
					enableDebug()
				}
				if c.Bool("http") {
					return brook.RunClientAsHTTP(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
				}
				return brook.RunClient(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
			},
		},
		cli.Command{
			Name:  "tunnel",
			Usage: "Run as tunnel mode on client-site",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address, like: 127.0.0.1:1080",
				},
				cli.StringFlag{
					Name:  "to, t",
					Usage: "Tunnel to where, like: 8.8.8.8:53",
				},
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("to") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "tunnel")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunTunnel(c.String("listen"), c.String("to"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
			},
		},
		cli.Command{
			Name:  "tproxy",
			Usage: "Run as tproxy mode on client-site, transparent proxy, only works on Linux",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address, like: 127.0.0.1:1080",
				},
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "tproxy")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunTproxy(c.String("listen"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
			},
		},
		cli.Command{
			Name:  "vpn",
			Usage: "Run as VPN mode on client-site",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address, must use 127.0.0.1, like: 127.0.0.1:1080",
				},
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
				cli.StringFlag{
					Name:  "dns",
					Value: "8.8.8.8",
					Usage: "DNS Server, like: 8.8.8.8",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpSessionTime",
					Value: 60,
					Usage: "udp session time (s), in most cases need this",
				},
				cli.StringFlag{
					Name:  "tunDevice",
					Usage: "tun name",
					Value: "tun0",
				},
				cli.StringFlag{
					Name:  "tunIP",
					Usage: "tun IP",
					Value: "10.9.9.2",
				},
				cli.StringFlag{
					Name:  "tunGateway",
					Usage: "tun gateway",
					Value: "10.9.9.1",
				},
				cli.StringFlag{
					Name:  "tunMask",
					Usage: "tun mask",
					Value: "255.255.255.0",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "vpn")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunVPN(c.String("listen"), c.String("server"), c.String("password"), c.String("dns"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"), c.String("tunDevice"), c.String("tunIP"), c.String("tunGateway"), c.String("tunMask"))
			},
		},
		cli.Command{
			Name:  "ssserver",
			Usage: "Run as shadowsocks server mode, fixed method is aes-256-cfb",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Server listen address, like: 0.0.0.0:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "ssserver")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunSSServer(c.String("listen"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
			},
		},
		cli.Command{
			Name:  "ssservers",
			Usage: "Run as shadowsocks multiple servers mode, fixed method is aes-256-cfb",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "listenpassword, l",
					Usage: "server address and password, like '0.0.0.0:1080 password'",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.StringSlice("listenpassword")) == 0 {
					cli.ShowCommandHelp(c, "ssservers")
					return nil
				}
				if debug {
					enableDebug()
				}
				errch := make(chan error)
				go func() {
					for _, v := range c.StringSlice("listenpassword") {
						ss := strings.Split(v, " ")
						if len(ss) != 2 {
							errch <- errors.New("Invalid listenpassword")
							return
						}
						go func() {
							errch <- brook.RunSSServer(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
						}()
					}
				}()
				return <-errch
			},
		},
		cli.Command{
			Name:  "ssclient",
			Usage: "Run as shadowsocks client mode, fixed method is aes-256-cfb",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address, like: 127.0.0.1:1080",
				},
				cli.StringFlag{
					Name:  "ip, i",
					Usage: "Client IP address, like: 127.0.0.1",
				},
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpSessionTime",
					Value: 60,
					Usage: "udp session time (s), in most cases need this",
				},
				cli.BoolFlag{
					Name:  "http",
					Usage: "If true, client start a http(s) proxy. default socks5",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("ip") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "ssclient")
					return nil
				}
				if debug {
					enableDebug()
				}
				if c.Bool("http") {
					return brook.RunSSClientAsHTTP(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
				}
				return brook.RunSSClient(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
			},
		},
		cli.Command{
			Name:  "socks5",
			Usage: "Run as raw socks5 server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address, like: 127.0.0.1:1080",
				},
				cli.StringFlag{
					Name:  "ip, i",
					Usage: "Client IP address, like: 127.0.0.1",
				},
				cli.StringFlag{
					Name:  "username",
					Usage: "User name, optional",
				},
				cli.StringFlag{
					Name:  "password",
					Usage: "Password, optional",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpSessionTime",
					Value: 60,
					Usage: "udp session time (s), in most cases need this",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("ip") == "" {
					cli.ShowCommandHelp(c, "socks5")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunSocks5Server(c.String("listen"), c.String("ip"), c.String("username"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
			},
		},
		cli.Command{
			Name:  "relay",
			Usage: "Run as relay mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Relay server address: 0.0.0.0:1080",
				},
				cli.StringFlag{
					Name:  "remote, r",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("remote") == "" {
					cli.ShowCommandHelp(c, "relay")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunRelay(c.String("listen"), c.String("remote"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
			},
		},
		cli.Command{
			Name:  "relays",
			Usage: "Run as multiple relays mode",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "listenremote, l",
					Usage: "listen address and server address, like '0.0.0.0:1080 1.2.3.4:1080'",
				},
				cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.StringSlice("listenremote")) == 0 {
					cli.ShowCommandHelp(c, "relays")
					return nil
				}
				if debug {
					enableDebug()
				}
				errch := make(chan error)
				go func() {
					for _, v := range c.StringSlice("listenremote") {
						ss := strings.Split(v, " ")
						if len(ss) != 2 {
							errch <- errors.New("Invalid listenremote")
							return
						}
						go func() {
							errch <- brook.RunRelay(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
						}()
					}
				}()
				return <-errch
			},
		},
		cli.Command{
			Name:  "link",
			Usage: "Print brook link",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "link")
					return nil
				}
				fmt.Println(brook.Link(c.String("server"), c.String("password")))
				return nil
			},
		},
		cli.Command{
			Name:  "qr",
			Usage: "Print brook server QR code",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Server password",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "qr")
					return nil
				}
				brook.QR(c.String("server"), c.String("password"))
				return nil
			},
		},
		cli.Command{
			Name:  "socks5tohttp",
			Usage: "Convert socks5 to http proxy",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address: like: 127.0.0.1:8080",
				},
				cli.StringFlag{
					Name:  "socks5, s",
					Usage: "Socks5 address",
				},
				cli.IntFlag{
					Name:  "timeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline",
					Value: 0,
					Usage: "connection deadline time (s)",
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
				return brook.RunSocks5ToHTTP(c.String("listen"), c.String("socks5"), c.Int("timeout"), c.Int("deadline"))
			},
		},
		cli.Command{
			Name:  "systemproxy",
			Usage: "Set system proxy with pac url, or remove, only works on MacOS/Windows",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "url, u",
					Usage: "Pac address: like: http://127.0.0.1/pac",
				},
				cli.BoolFlag{
					Name:  "remove, r",
					Usage: "Remove pac url from system proxy",
				},
			},
			Action: func(c *cli.Context) error {
				if !c.Bool("remove") && c.String("url") == "" {
					cli.ShowCommandHelp(c, "systemproxy")
					return nil
				}
				return brook.RunSystemProxy(c.Bool("remove"), c.String("url"))
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func enableDebug() {
	go func() {
		log.Println(http.ListenAndServe(debugAddress, nil))
	}()
	brook.EnableDebug()
}
