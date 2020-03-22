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
	"os/signal"
	"strings"
	"syscall"

	"net/http"
	_ "net/http/pprof"

	"github.com/txthinking/brook"
	"github.com/txthinking/brook/sysproxy"
	"github.com/urfave/cli"
)

var debug bool
var debugAddress string

func main() {
	app := cli.NewApp()
	app.Name = "Brook"
	app.Version = "20200201"
	app.Usage = "A Cross-Platform Proxy/VPN Software"
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
			Usage: "Run as server mode",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Server listen address, like: 0.0.0.0:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				s, err := brook.NewServer(c.String("listen"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
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
			Usage: "Run as multiple servers mode",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "listenpassword",
					Aliases: []string{"l"},
					Usage:   "server and password, like '0.0.0.0:1080 password'",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				l := make([]*brook.Server, 0)
				for _, v := range c.StringSlice("listenpassword") {
					ss := strings.Split(v, " ")
					if len(ss) != 2 {
						return errors.New("Invalid listenpassword")
					}
					s, err := brook.NewServer(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
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
			Usage: "Run as client with brook server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address, like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:    "ip",
					Aliases: []string{"i"},
					Usage:   "Client IP address, like: 127.0.0.1",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpSessionTime",
					Value: 60,
					Usage: "udp session time (s), in most cases need this",
				},
				&cli.BoolFlag{
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
				s, err := brook.NewClient(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				if c.Bool("http") {
					return s.ListenAndServeHTTP()
				}
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "wsserver",
			Usage: "Run as websocket server mode",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Server listen address, like: 0.0.0.0:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.StringFlag{
					Name:  "domain",
					Usage: "If domain is specified, the domain must have been resolved to the external IP, listen will be ignored, 80 and 443 will be used",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if (c.String("listen") == "" && c.String("domain") == "") || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wsserver")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewWSServer(c.String("listen"), c.String("password"), c.String("domain"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
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
			Usage: "Run as websocket client with brook wsserver",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address, like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:    "ip",
					Aliases: []string{"i"},
					Usage:   "Client IP address, like: 127.0.0.1",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: ws://1.2.3.4:1080, wss://google.com:443. Do not omit the port under any circumstances",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpSessionTime",
					Value: 60,
					Usage: "udp session time (s), in most cases need this",
				},
				&cli.BoolFlag{
					Name:  "http",
					Usage: "If true, client start a http(s) proxy. default socks5",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("ip") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "wsclient")
					return nil
				}
				if debug {
					enableDebug()
				}
				s, err := brook.NewWSClient(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				if c.Bool("http") {
					return s.ListenAndServeHTTP()
				}
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "tunnel",
			Usage: "Run as tunnel with brook server on client-site",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address, like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:    "to",
					Aliases: []string{"t"},
					Usage:   "Tunnel to where, like: 8.8.8.8:53",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				s, err := brook.NewTunnel(c.String("listen"), c.String("to"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
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
			Usage: "Run DNS server with brook server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address, like: 127.0.0.1:53",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.StringFlag{
					Name:  "defaultDNSServer",
					Usage: "Default DNS server",
					Value: "8.8.8.8:53",
				},
				&cli.StringFlag{
					Name:  "listDNSServer",
					Usage: "DNS server for resolving domain in list",
					Value: "223.5.5.5:53",
				},
				&cli.StringFlag{
					Name:  "list",
					Usage: "https://, http:// or file://",
					Value: "https://blackwhite.txthinking.com/white.list",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpDeadline",
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
				s, err := brook.NewDNS(c.String("listen"), c.String("server"), c.String("password"), c.String("defaultDNSServer"), c.String("listDNSServer"), c.String("list"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
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
			Usage: "Run as tproxy mode on client-site, transparent proxy, only works on Linux",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address, don't contain IP, just like: :1080",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.BoolFlag{
					Name:  "letBrookDoAllForMe",
					Usage: "See more: https://github.com/txthinking/brook/wiki/How-to-run-transparent-proxy-on-Linux%3F",
				},
				&cli.BoolFlag{
					Name:  "cleanBrookDidForMe",
					Usage: "See more: https://github.com/txthinking/brook/wiki/How-to-run-transparent-proxy-on-Linux%3F",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				s, err := brook.NewTproxy(c.String("listen"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
				if err != nil {
					return err
				}
				if c.Bool("cleanBrookDidForMe") {
					if err := s.ClearAutoScripts(); err != nil {
						return err
					}
					return nil
				}
				var dns *brook.DNS
				if c.Bool("letBrookDoAllForMe") {
					if err := s.RunAutoScripts(); err != nil {
						return err
					}
					dns, err = brook.NewDNS(":53", c.String("server"), c.String("password"), "8.8.8.8:53", "223.5.5.5:53", "https://blackwhite.txthinking.com/white.list", c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
					if err != nil {
						return err
					}
					go func() {
						if err := dns.ListenAndServe(); err != nil {
							log.Println(err)
						}
					}()
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					if c.Bool("letBrookDoAllForMe") {
						if err := s.ClearAutoScripts(); err != nil {
							log.Println(err)
						}
					}
					if dns != nil {
						dns.Shutdown()
					}
					s.Shutdown()
				}()
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "vpn",
			Usage: "Run as VPN mode on client-site",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address, must use 127.0.0.1, like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.StringFlag{
					Name:  "dns",
					Value: "8.8.8.8",
					Usage: "DNS Server, like: 8.8.8.8",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpSessionTime",
					Value: 60,
					Usage: "udp session time (s), in most cases need this",
				},
				&cli.StringFlag{
					Name:  "tunDevice",
					Usage: "tun name",
					Value: "tun0",
				},
				&cli.StringFlag{
					Name:  "tunIP",
					Usage: "tun IP",
					Value: "10.9.9.2",
				},
				&cli.StringFlag{
					Name:  "tunGateway",
					Usage: "tun gateway",
					Value: "10.9.9.1",
				},
				&cli.StringFlag{
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
				s, err := brook.NewVPN(c.String("listen"), c.String("server"), c.String("password"), c.String("dns"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"), c.String("tunDevice"), c.String("tunIP"), c.String("tunGateway"), c.String("tunMask"))
				if err != nil {
					return err
				}
				go func() {
					fmt.Println("Ctrl-C to quit")
					log.Println(s.ListenAndServe())
				}()
				sigs := make(chan os.Signal, 1)
				signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
				<-sigs
				fmt.Println("Quitting...")
				return s.Shutdown()
			},
		},
		&cli.Command{
			Name:  "ssserver",
			Usage: "Run as shadowsocks server mode, fixed method is aes-256-cfb",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Server listen address, like: 0.0.0.0:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				s, err := brook.NewSSServer(c.String("listen"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
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
			Name:  "ssservers",
			Usage: "Run as shadowsocks multiple servers mode, fixed method is aes-256-cfb",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "listenpassword",
					Aliases: []string{"l"},
					Usage:   "server address and password, like '0.0.0.0:1080 password'",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				l := make([]*brook.SSServer, 0)
				for _, v := range c.StringSlice("listenpassword") {
					ss := strings.Split(v, " ")
					if len(ss) != 2 {
						return errors.New("Invalid listenpassword")
					}
					s, err := brook.NewSSServer(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
					if err != nil {
						return err
					}
					l = append(l, s)
				}
				for _, v := range l {
					go func(v *brook.SSServer) {
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
			Name:  "ssclient",
			Usage: "Run as shadowsocks client mode, fixed method is aes-256-cfb",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address, like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:    "ip",
					Aliases: []string{"i"},
					Usage:   "Client IP address, like: 127.0.0.1",
				},
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:1080",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpSessionTime",
					Value: 60,
					Usage: "udp session time (s), in most cases need this",
				},
				&cli.BoolFlag{
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
				s, err := brook.NewSSClient(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
				if err != nil {
					return err
				}
				go func() {
					sigs := make(chan os.Signal, 1)
					signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
					<-sigs
					s.Shutdown()
				}()
				if c.Bool("http") {
					return s.ListenAndServeHTTP()
				}
				return s.ListenAndServe()
			},
		},
		&cli.Command{
			Name:  "socks5",
			Usage: "Run as raw socks5 server",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address, like: 127.0.0.1:1080",
				},
				&cli.StringFlag{
					Name:    "ip",
					Aliases: []string{"i"},
					Usage:   "Client IP address, like: 127.0.0.1",
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
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
					Name:  "udpDeadline",
					Value: 60,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				s, err := brook.NewSocks5Server(c.String("listen"), c.String("ip"), c.String("username"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
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
			Name:  "relay",
			Usage: "Run as relay mode",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Relay server address: 0.0.0.0:1080",
				},
				&cli.StringFlag{
					Name:    "remote",
					Aliases: []string{"r"},
					Usage:   "Server address, like: 1.2.3.4:1080",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				s, err := brook.NewRelay(c.String("listen"), c.String("remote"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
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
			Usage: "Run as multiple relays mode",
			Flags: []cli.Flag{
				&cli.StringSliceFlag{
					Name:    "listenremote",
					Aliases: []string{"l"},
					Usage:   "listen address and server address, like '0.0.0.0:1080 1.2.3.4:1080'",
				},
				&cli.IntFlag{
					Name:  "tcpTimeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
					Name:  "tcpDeadline",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				&cli.IntFlag{
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
				l := make([]*brook.Relay, 0)
				for _, v := range c.StringSlice("listenremote") {
					ss := strings.Split(v, " ")
					if len(ss) != 2 {
						return errors.New("Invalid listenremote")
					}
					s, err := brook.NewRelay(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
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
			Name:  "link",
			Usage: "Print brook link",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:1080, ws://1.2.3.4:1080, wss://google.com:443. Do not omit the port under any circumstances",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
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
		&cli.Command{
			Name:  "qr",
			Usage: "Print brook server QR code",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "server",
					Aliases: []string{"s"},
					Usage:   "Server address, like: 1.2.3.4:1080, ws://1.2.3.4:1080, wss://google.com:443. Do not omit the port under any circumstances",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"p"},
					Usage:   "Server password",
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
		&cli.Command{
			Name:  "socks5tohttp",
			Usage: "Convert socks5 to http proxy",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "listen",
					Aliases: []string{"l"},
					Usage:   "Client listen address: like: 127.0.0.1:8080",
				},
				&cli.StringFlag{
					Name:    "socks5",
					Aliases: []string{"s"},
					Usage:   "Socks5 address",
				},
				&cli.StringFlag{
					Name:  "socks5username",
					Usage: "Socks5 username",
				},
				&cli.StringFlag{
					Name:  "socks5password",
					Usage: "Socks5 password",
				},
				&cli.IntFlag{
					Name:  "timeout",
					Value: 60,
					Usage: "connection tcp keepalive timeout (s)",
				},
				&cli.IntFlag{
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
				s, err := brook.NewSocks5ToHTTP(c.String("listen"), c.String("socks5"), c.String("socks5username"), c.String("socks5password"), c.Int("timeout"), c.Int("deadline"))
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
			Name:  "systemproxy",
			Usage: "Set system proxy with pac url, or remove, only works on macOS/Windows",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "url",
					Aliases: []string{"u"},
					Usage:   "Pac address: like: http://127.0.0.1/pac",
				},
				&cli.BoolFlag{
					Name:    "remove",
					Aliases: []string{"r"},
					Usage:   "Remove pac url from system proxy",
				},
			},
			Action: func(c *cli.Context) error {
				if !c.Bool("remove") && c.String("url") == "" {
					cli.ShowCommandHelp(c, "systemproxy")
					return nil
				}
				if c.Bool("remove") {
					if err := sysproxy.TurnOffSystemProxy(); err != nil {
						return err
					}
					return nil
				}
				if err := sysproxy.TurnOnSystemProxy(c.String("url")); err != nil {
					return err
				}
				return nil
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
