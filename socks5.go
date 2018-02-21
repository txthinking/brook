package brook

import (
	"errors"
	"io"
	"log"
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/socks5"
)

// Socks5Server is the client of raw socks5 protocol
type Socks5Server struct {
	Server          *socks5.Server
	Socks5Middleman plugin.Socks5Middleman
	TCPTimeout      int
	TCPDeadline     int // not refreshed
	UDPDeadline     int
	UDPSessionTime  int
	ForwardAddress  string
	ForwardUserName string
	ForwardPassword string
	Cache           *cache.Cache
}

// NewSocks5Server returns a new Socks5Server
func NewSocks5Server(addr, ip, userName, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) (*Socks5Server, error) {
	s5, err := socks5.NewClassicServer(addr, ip, userName, password, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return nil, err
	}
	cs := cache.New(60*time.Minute, 10*time.Minute)
	x := &Socks5Server{
		Server:         s5,
		TCPTimeout:     tcpTimeout,
		TCPDeadline:    tcpDeadline,
		UDPDeadline:    udpDeadline,
		UDPSessionTime: udpSessionTime,
		Cache:          cs,
	}
	return x, nil
}

// SetSocks5Middleman sets socks5middleman plugin
func (x *Socks5Server) SetSocks5Middleman(m plugin.Socks5Middleman) {
	x.Socks5Middleman = m
}

// ListenAndServe will let client start to listen and serve, sm can be nil
func (x *Socks5Server) ListenAndServe() error {
	return x.Server.Run(nil)
}

// ListenAndForward will let client start a proxy to listen and forward to another socks5,
// sm can be nil
func (x *Socks5Server) ListenAndForward(addr, username, password string) error {
	x.ForwardAddress = addr
	x.ForwardUserName = username
	x.ForwardPassword = password
	return x.Server.Run(x)
}

// TCPHandle handles tcp request
func (x *Socks5Server) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
	if x.Socks5Middleman != nil {
		done, err := x.Socks5Middleman.TCPHandle(s, c, r)
		if err != nil {
			if done {
				return err
			}
			return ErrorReply(r, c, err)
		}
		if done {
			return nil
		}
	}

	client, err := socks5.NewClient(x.ForwardAddress, x.ForwardUserName, x.ForwardPassword, x.TCPTimeout, x.TCPDeadline, x.UDPDeadline)
	if err != nil {
		return err
	}
	if err := client.Negotiate(); err != nil {
		return ErrorReply(r, c, err)
	}
	defer client.TCPConn.Close()

	if r.Cmd == socks5.CmdUDP {
		// TODO If client's udp address is not 0, then prepare a local udp address
		if r.Atyp == socks5.ATYPIPv4 || r.Atyp == socks5.ATYPDomain {
			r.DstAddr = net.IPv4zero
		} else {
			r.DstAddr = net.IPv6zero
		}
		r.DstPort = []byte{0x00, 0x00}
	}
	rp, err := client.Request(r)
	if err != nil {
		return ErrorReply(r, c, err)
	}

	// reply ok and choose address according to cmd or something wrong
	if r.Cmd == socks5.CmdConnect {
		a, address, port, err := socks5.ParseAddress(client.TCPConn.LocalAddr().String())
		if err != nil {
			return ErrorReply(r, c, err)
		}
		rp.Atyp = a
		rp.BndAddr = address
		rp.BndPort = port
		if err := rp.WriteTo(c); err != nil {
			return err
		}
		go func() {
			_, _ = io.Copy(c, client.TCPConn)
		}()
		_, _ = io.Copy(client.TCPConn, c)
		return nil
	}
	if r.Cmd == socks5.CmdUDP {
		caddr, err := r.UDP(c, x.Server.ServerAddr)
		if err != nil {
			return err
		}
		// Because we always send zero:0 to remote
		_, ok := x.Cache.Get("RUA")
		if !ok {
			x.Cache.Set("RUA", rp.Address(), cache.DefaultExpiration)
		}
		_, p, err := net.SplitHostPort(caddr.String())
		if err != nil {
			return err
		}
		if p == "0" {
			time.Sleep(time.Duration(x.Server.UDPSessionTime) * time.Second)
			return nil
		}
		ch := make(chan byte)
		x.Server.TCPUDPAssociate.Set(caddr.String(), ch, cache.DefaultExpiration)
		<-ch
		return nil
	}
	return ErrorReply(r, c, socks5.ErrUnsupportCmd)
}

// UDPHandle handles udp request
func (x *Socks5Server) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	if x.Socks5Middleman != nil {
		if done, err := x.Socks5Middleman.UDPHandle(s, addr, d); err != nil || done {
			return err
		}
	}

	send := func(ue *socks5.UDPExchange, data []byte) error {
		_, err := ue.RemoteConn.Write(data)
		if err != nil {
			return err
		}
		return nil
	}

	var ue *socks5.UDPExchange
	iue, ok := s.UDPExchanges.Get(addr.String())
	if ok {
		ue = iue.(*socks5.UDPExchange)
		return send(ue, d.Bytes())
	}

	raddr, ok := x.Cache.Get("RUA")
	if !ok {
		return errors.New("Can not find remote udp address.")
	}
	tmp, err := Dial.Dial("udp", raddr.(string))
	if err != nil {
		return err
	}
	rc := tmp.(*net.UDPConn)
	ue = &socks5.UDPExchange{
		ClientAddr: addr,
		RemoteConn: rc,
	}
	if err := send(ue, d.Bytes()); err != nil {
		return err
	}
	s.UDPExchanges.Set(ue.ClientAddr.String(), ue, cache.DefaultExpiration)
	go func(ue *socks5.UDPExchange) {
		defer func() {
			v, ok := s.TCPUDPAssociate.Get(ue.ClientAddr.String())
			if ok {
				ch := v.(chan byte)
				ch <- '0'
			}
			s.UDPExchanges.Delete(ue.ClientAddr.String())
			ue.RemoteConn.Close()
		}()
		var b [65536]byte
		for {
			if s.UDPDeadline != 0 {
				if err := ue.RemoteConn.SetDeadline(time.Now().Add(time.Duration(s.UDPDeadline) * time.Second)); err != nil {
					log.Println(err)
					break
				}
			}
			n, err := ue.RemoteConn.Read(b[:])
			if err != nil {
				log.Println(err)
				break
			}
			d1, err := socks5.NewDatagramFromBytes(b[0:n])
			if err != nil {
				log.Println(err)
				break
			}
			a, addr, port, err := socks5.ParseAddress(ue.ClientAddr.String())
			if err != nil {
				log.Println(err)
				break
			}
			d1 = socks5.NewDatagram(a, addr, port, d1.Data)
			if _, err := s.UDPConn.WriteToUDP(d1.Bytes(), ue.ClientAddr); err != nil {
				log.Println(err)
				break
			}
		}
	}(ue)
	return nil
}

// Shutdown used to stop the client
func (x *Socks5Server) Shutdown() error {
	return x.Server.Stop()
}
