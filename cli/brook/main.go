package main

import (
	"errors"
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
	app.Version = "20180707"
	app.Usage = "A Cross-Platform Proxy Software"
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
			Usage: "Run as tunnel mode",
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
			Usage: "Run as tproxy mode, transparent proxy, only works on Linux",
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
			Name:  "streamserver",
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
					cli.ShowCommandHelp(c, "streamserver")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunStreamServer(c.String("listen"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
			},
		},
		cli.Command{
			Name:  "streamservers",
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
					cli.ShowCommandHelp(c, "streamservers")
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
							errch <- brook.RunStreamServer(ss[0], ss[1], c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"))
						}()
					}
				}()
				return <-errch
			},
		},
		cli.Command{
			Name:  "streamclient",
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
					cli.ShowCommandHelp(c, "streamclient")
					return nil
				}
				if debug {
					enableDebug()
				}
				if c.Bool("http") {
					return brook.RunStreamClientAsHTTP(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
				}
				return brook.RunStreamClient(c.String("listen"), c.String("ip"), c.String("server"), c.String("password"), c.Int("tcpTimeout"), c.Int("tcpDeadline"), c.Int("udpDeadline"), c.Int("udpSessionTime"))
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
				cli.BoolFlag{
					Name:  "stream",
					Usage: "Brook stream or not",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "qr")
					return nil
				}
				brook.QR(c.Bool("stream"), c.String("server"), c.String("password"))
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
