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
	app.Version = "20170814"
	app.Usage = "A Cross-Platform Proxy Software"
	app.Author = "Cloud"
	app.Email = "cloud@txthinking.com"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug, more logs",
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
			Name:  "bkserver",
			Usage: "Run as brook protocol server mode",
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
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.StringFlag{
					Name:  "music, m",
					Usage: "The music you love to listen: https://github.com/txthinking/brook/wiki/Music-List",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "bkserver")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunBKServer(c.String("listen"), c.String("password"), c.Int("timeout"), c.Int("deadline"), c.String("music"))
			},
		},
		cli.Command{
			Name:  "bkservers",
			Usage: "Run as brook protocol multiple servers mode",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.StringSliceFlag{
					Name:  "listenpasswordmusic, l",
					Usage: "server address and password, like '0.0.0.0:1080 password' or '0.0.0.0:1080 password music'",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.StringSlice("listenpasswordmusic")) == 0 {
					cli.ShowCommandHelp(c, "bkservers")
					return nil
				}
				if debug {
					enableDebug()
				}
				errch := make(chan error)
				go func() {
					for _, v := range c.StringSlice("listenpasswordmusic") {
						ss := strings.Split(v, " ")
						if len(ss) < 2 || len(ss) > 3 {
							errch <- errors.New("Invalid listenpasswordmusic")
							return
						}
						var music string
						if len(ss) == 3 {
							music = ss[2]
						}
						go func() {
							errch <- brook.RunBKServer(ss[0], ss[1], c.Int("timeout"), c.Int("deadline"), music)
						}()
					}
				}()
				return <-errch
			},
		},
		cli.Command{
			Name:  "bkclient",
			Usage: "Run as brook protocol client mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address: like: 127.0.0.1:1080",
				},
				cli.BoolFlag{
					Name:  "http",
					Usage: "If true, client start a http(s) proxy. default socks5",
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
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.StringFlag{
					Name:  "music, m",
					Usage: "The music you love to listen: https://github.com/txthinking/brook/wiki/Music-List or a url(http://a.com/some https://a.com/some file:///path/to/file)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "bkclient")
					return nil
				}
				if debug {
					enableDebug()
				}
				if c.Bool("http") {
					return brook.RunBKHTTPClient(c.String("listen"), c.String("server"), c.String("password"), c.Int("timeout"), c.Int("deadline"), c.String("music"))
				}
				return brook.RunBKClient(c.String("listen"), c.String("server"), c.String("password"), c.Int("timeout"), c.Int("deadline"), c.String("music"))
			},
		},
		cli.Command{
			Name:  "s5server",
			Usage: "Run as socks5 encrypt protocol server mode",
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
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "s5server")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunS5Server(c.String("listen"), c.String("password"), c.Int("timeout"), c.Int("deadline"))
			},
		},
		cli.Command{
			Name:  "s5servers",
			Usage: "Run as socks5 encrypt protocol multiple servers mode",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.StringSliceFlag{
					Name:  "listenpassword, l",
					Usage: "server address and password, like '0.0.0.0:1080 password'",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.StringSlice("listenpassword")) == 0 {
					cli.ShowCommandHelp(c, "s5servers")
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
							errch <- brook.RunS5Server(ss[0], ss[1], c.Int("timeout"), c.Int("deadline"))
						}()
					}
				}()
				return <-errch
			},
		},
		cli.Command{
			Name:  "s5client",
			Usage: "Run as socks5 encrypt protocol client mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address: like: 127.0.0.1:1080",
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
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "s5client")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunS5Client(c.String("listen"), c.String("server"), c.String("password"), c.Int("timeout"), c.Int("deadline"))
			},
		},
		cli.Command{
			Name:  "ssserver",
			Usage: "Run as shadowsocks protocol server mode, fixed method is aes-256-cfb",
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
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
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
				return brook.RunSSServer(c.String("listen"), c.String("password"), c.Int("timeout"), c.Int("deadline"))
			},
		},
		cli.Command{
			Name:  "ssservers",
			Usage: "Run as shadowsocks protocol multiple servers mode, fixed method is aes-256-cfb",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
				cli.StringSliceFlag{
					Name:  "listenpassword, l",
					Usage: "server address and password, like '0.0.0.0:1080 password'",
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
							errch <- brook.RunSSServer(ss[0], ss[1], c.Int("timeout"), c.Int("deadline"))
						}()
					}
				}()
				return <-errch
			},
		},
		cli.Command{
			Name:  "ssclient",
			Usage: "Run as shadowsocks protocol client mode, fixed method is aes-256-cfb",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address: like: 127.0.0.1:1080",
				},
				cli.BoolFlag{
					Name:  "http",
					Usage: "If true, client start a http(s) proxy. default socks5",
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
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "ssclient")
					return nil
				}
				if debug {
					enableDebug()
				}
				if c.Bool("http") {
					return brook.RunSSHTTPClient(c.String("listen"), c.String("server"), c.String("password"), c.Int("timeout"), c.Int("deadline"))
				}
				return brook.RunSSClient(c.String("listen"), c.String("server"), c.String("password"), c.Int("timeout"), c.Int("deadline"))
			},
		},
		cli.Command{
			Name:  "relay",
			Usage: "Run as relay mode",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Relay server address: :1080",
				},
				cli.StringFlag{
					Name:  "server, s",
					Usage: "Server address, like: 1.2.3.4:1080",
				},
				cli.IntFlag{
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("listen") == "" || c.String("server") == "" {
					cli.ShowCommandHelp(c, "relay")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunRelay(c.String("listen"), c.String("server"), c.Int("timeout"), c.Int("deadline"))
			},
		},
		cli.Command{
			Name:  "relays",
			Usage: "Run as multiple relays mode",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "listenserver, l",
					Usage: "listen address and server address, like '0.0.0.0:1080 1.2.3.4:1080'",
				},
				cli.IntFlag{
					Name:  "timeout, t",
					Value: 0,
					Usage: "connection tcp keepalive timeout (s)",
				},
				cli.IntFlag{
					Name:  "deadline, d",
					Value: 0,
					Usage: "connection deadline time (s)",
				},
			},
			Action: func(c *cli.Context) error {
				if len(c.StringSlice("listenserver")) == 0 {
					cli.ShowCommandHelp(c, "relays")
					return nil
				}
				if debug {
					enableDebug()
				}
				errch := make(chan error)
				go func() {
					for _, v := range c.StringSlice("listenserver") {
						ss := strings.Split(v, " ")
						if len(ss) != 2 {
							errch <- errors.New("Invalid listenserver")
							return
						}
						go func() {
							errch <- brook.RunRelay(ss[0], ss[1], c.Int("timeout"), c.Int("deadline"))
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
				cli.StringFlag{
					Name:  "music, m",
					Usage: "Server music",
				},
			},
			Action: func(c *cli.Context) error {
				if c.String("server") == "" || c.String("password") == "" {
					cli.ShowCommandHelp(c, "qr")
					return nil
				}
				brook.QR(c.String("server"), c.String("password"), c.String("music"))
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
