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
	app.Version = "20170826"
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
					cli.ShowCommandHelp(c, "server")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunBKServer(c.String("listen"), c.String("password"), c.Int("timeout"), c.Int("deadline"), c.String("music"))
			},
		},
		cli.Command{
			Name:  "servers",
			Usage: "Run as multiple servers mode",
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
					cli.ShowCommandHelp(c, "servers")
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
			Name:  "client",
			Usage: "Run as client mode",
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
					cli.ShowCommandHelp(c, "client")
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
			Usage: "Run as shadowsocks multiple servers mode, fixed method is aes-256-cfb",
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
			Usage: "Run as shadowsocks client mode, fixed method is aes-256-cfb",
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
		cli.Command{
			Name:  "socks5",
			Usage: "Run as raw socks5 server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen, l",
					Usage: "Client listen address: like: 127.0.0.1:1080",
				},
				cli.StringFlag{
					Name:  "username, u",
					Usage: "User name, optional",
				},
				cli.StringFlag{
					Name:  "password, p",
					Usage: "Password, optional",
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
				if c.String("listen") == "" {
					cli.ShowCommandHelp(c, "socks5")
					return nil
				}
				if debug {
					enableDebug()
				}
				return brook.RunSocks5Server(c.String("listen"), c.String("username"), c.String("password"), c.Int("timeout"), c.Int("deadline"))
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
	}
	if len(os.Args) > 1 {
		os.Args[1] = strings.Replace(os.Args[1], "bk", "", -1)
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
