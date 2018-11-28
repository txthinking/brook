package brook

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/txthinking/brook/sysproxy"
	"github.com/txthinking/gotun2socks"
	"github.com/txthinking/gotun2socks/tun"
)

// VPN
type VPN struct {
	Client         *Client
	Tunnel         *Tunnel
	Tun            *gotun2socks.Tun2Socks
	ServerIP       string
	TunGateway     string
	DefaultGateway string
}

// NewVPN
func NewVPN(addr, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int, tunDevice, tunIP, tunGateway, tunMask, defaultGateway string) (*VPN, error) {
	h, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if h != "127.0.0.1" {
		return nil, errors.New("Must listen on 127.0.0.1")
	}
	h, p, err := net.SplitHostPort(server)
	if err != nil {
		return nil, err
	}
	l, err := net.LookupIP(h)
	if err != nil {
		return nil, err
	}
	s := ""
	for _, v := range l {
		if v.To4() == nil {
			continue
		}
		s = v.String()
		break
	}
	if s == "" {
		return nil, errors.New("Can not find server v4 IP")
	}
	server = net.JoinHostPort(s, p)

	c, err := NewClient(addr, "127.0.0.1", server, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return nil, err
	}
	tl, err := NewTunnel("127.0.0.1:53", "8.8.8.8:53", server, password, tcpTimeout, tcpDeadline, udpDeadline)
	if err != nil {
		return nil, err
	}
	f, err := tun.OpenTunDevice(tunDevice, tunIP, tunGateway, tunMask, []string{"8.8.8.8"})
	if err != nil {
		return nil, err
	}
	t := gotun2socks.New(f, addr, []string{"8.8.8.8"}, true, true)
	return &VPN{
		Client:         c,
		Tunnel:         tl,
		Tun:            t,
		ServerIP:       s,
		TunGateway:     tunGateway,
		DefaultGateway: defaultGateway,
	}, nil
}

// ListenAndServe starts to run VPN
func (v *VPN) ListenAndServe() error {
	if err := sysproxy.SetDNSServer("127.0.0.1"); err != nil {
		return err
	}
	if err := v.AddRoutes(); err != nil {
		return err
	}

	errch := make(chan error)
	go func() {
		errch <- v.Client.ListenAndServe()
	}()
	go func() {
		errch <- v.Tunnel.ListenAndServe()
	}()
	go func() {
		v.Tun.Run()
	}()
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
		<-sigs
		errch <- nil
	}()
	fmt.Println("Ctrl-C to quit")

	err := <-errch
	if err := v.Shutdown(); err != nil {
		return err
	}
	return err
}

// Shutdown stops VPN
func (v *VPN) Shutdown() error {
	fmt.Println("Quitting...")
	if err := sysproxy.SetDNSServer("8.8.8.8"); err != nil {
		log.Println(err)
	}
	if err := v.DeleteRoutes(); err != nil {
		log.Println(err)
	}
	if v.Client != nil {
		if err := v.Client.Shutdown(); err != nil {
			log.Println(err)
		}
	}
	if v.Tunnel != nil {
		if err := v.Tunnel.Shutdown(); err != nil {
			log.Println(err)
		}
	}
	if v.Tun != nil {
		// v.Tun.Stop()
	}
	return nil
}
