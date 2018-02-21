package brook

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/txthinking/ant"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/socks5"
)

// SSClient
type SSClient struct {
	Server          *socks5.Server
	RemoteAddr      string
	Password        []byte
	TCPTimeout      int
	TCPDeadline     int // Not refreshed
	UDPDeadline     int
	TCPListen       *net.TCPListener
	Socks5Middleman plugin.Socks5Middleman
	HTTPMiddleman   plugin.HTTPMiddleman
}

// NewSSClient returns a new SSClient
func NewSSClient(addr, ip, server, password string, tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime int) (*SSClient, error) {
	s5, err := socks5.NewClassicServer(addr, ip, "", "", tcpTimeout, tcpDeadline, udpDeadline, udpSessionTime)
	if err != nil {
		return nil, err
	}
	x := &SSClient{
		RemoteAddr:  server,
		Server:      s5,
		Password:    MakeSSKey(password),
		TCPTimeout:  tcpTimeout,
		TCPDeadline: tcpDeadline,
		UDPDeadline: udpDeadline,
	}
	return x, nil
}

// SetSocks5Middleman sets socks5middleman plugin
func (x *SSClient) SetSocks5Middleman(m plugin.Socks5Middleman) {
	x.Socks5Middleman = m
}

// SetHTTPMiddleman sets httpmiddleman plugin
func (x *SSClient) SetHTTPMiddleman(m plugin.HTTPMiddleman) {
	x.HTTPMiddleman = m
}

// ListenAndServe will let client start a socks5 proxy
// sm can be nil
func (x *SSClient) ListenAndServe() error {
	return x.Server.Run(x)
}

// TCPHandle handles tcp request
func (x *SSClient) TCPHandle(s *socks5.Server, c *net.TCPConn, r *socks5.Request) error {
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

	if r.Cmd == socks5.CmdConnect {
		tmp, err := Dial.Dial("tcp", x.RemoteAddr)
		if err != nil {
			return ErrorReply(r, c, err)
		}
		rc := tmp.(*net.TCPConn)
		defer rc.Close()
		if x.TCPTimeout != 0 {
			if err := rc.SetKeepAlivePeriod(time.Duration(x.TCPTimeout) * time.Second); err != nil {
				return ErrorReply(r, c, err)
			}
		}
		if x.TCPDeadline != 0 {
			if err := rc.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
				return ErrorReply(r, c, err)
			}
		}
		crc, err := x.WrapCipherConn(rc)
		if err != nil {
			return ErrorReply(r, c, err)
		}
		rawaddr := make([]byte, 0, 7)
		rawaddr = append(rawaddr, r.Atyp)
		rawaddr = append(rawaddr, r.DstAddr...)
		rawaddr = append(rawaddr, r.DstPort...)
		if _, err := crc.Write(rawaddr); err != nil {
			return ErrorReply(r, c, err)
		}
		a, address, port, err := socks5.ParseAddress(rc.LocalAddr().String())
		if err != nil {
			return ErrorReply(r, c, err)
		}

		rp := socks5.NewReply(socks5.RepSuccess, a, address, port)
		if err := rp.WriteTo(c); err != nil {
			return err
		}

		go func() {
			iv := make([]byte, aes.BlockSize)
			if _, err := io.ReadFull(crc, iv); err != nil {
				log.Println(err)
				return
			}
			_, _ = io.Copy(c, crc)
		}()
		_, _ = io.Copy(crc, c)
		return nil
	}
	if r.Cmd == socks5.CmdUDP {
		caddr, err := r.UDP(c, x.Server.ServerAddr)
		if err != nil {
			return err
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
	return socks5.ErrUnsupportCmd
}

// UDPHandle handles udp request
func (x *SSClient) UDPHandle(s *socks5.Server, addr *net.UDPAddr, d *socks5.Datagram) error {
	if x.Socks5Middleman != nil {
		if done, err := x.Socks5Middleman.UDPHandle(s, addr, d); err != nil || done {
			return err
		}
	}

	send := func(ue *socks5.UDPExchange, data []byte) error {
		cd, err := x.Encrypt(data)
		if err != nil {
			return err
		}
		_, err = ue.RemoteConn.Write(cd)
		if err != nil {
			return err
		}
		return nil
	}

	var ue *socks5.UDPExchange
	iue, ok := s.UDPExchanges.Get(addr.String())
	if ok {
		ue = iue.(*socks5.UDPExchange)
		return send(ue, d.Bytes()[3:])
	}

	c, err := Dial.Dial("udp", x.RemoteAddr)
	if err != nil {
		return err
	}
	rc := c.(*net.UDPConn)
	ue = &socks5.UDPExchange{
		ClientAddr: addr,
		RemoteConn: rc,
	}
	if err := send(ue, d.Bytes()[3:]); err != nil {
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
			_, _, _, data, err := x.Decrypt(b[0:n])
			if err != nil {
				log.Println(err)
				break
			}
			a, addr, port, err := socks5.ParseAddress(ue.ClientAddr.String())
			if err != nil {
				log.Println(err)
				break
			}
			d1 := socks5.NewDatagram(a, addr, port, data)
			if _, err := s.UDPConn.WriteToUDP(d1.Bytes(), ue.ClientAddr); err != nil {
				log.Println(err)
				break
			}
		}
	}(ue)
	return nil
}

// ListenAndServeHTTP will let client start a http proxy
// m can be nil
func (x *SSClient) ListenAndServeHTTP() error {
	var err error
	x.TCPListen, err = net.ListenTCP("tcp", x.Server.TCPAddr)
	if err != nil {
		return nil
	}
	for {
		c, err := x.TCPListen.AcceptTCP()
		if err != nil {
			return err
		}
		go func(c *net.TCPConn) {
			defer c.Close()
			if x.TCPTimeout != 0 {
				if err := c.SetKeepAlivePeriod(time.Duration(x.TCPTimeout) * time.Second); err != nil {
					log.Println(err)
					return
				}
			}
			if x.TCPDeadline != 0 {
				if err := c.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
					log.Println(err)
					return
				}
			}
			if err := x.HTTPHandle(c); err != nil {
				log.Println(err)
				return
			}
		}(c)
	}
}

// HTTPHandle handle http request
func (x *SSClient) HTTPHandle(c *net.TCPConn) error {
	b := make([]byte, 0, 1024)
	for {
		var b1 [1024]byte
		n, err := c.Read(b1[:])
		if err != nil {
			return err
		}
		b = append(b, b1[:n]...)
		if bytes.Contains(b, []byte{0x0d, 0x0a, 0x0d, 0x0a}) {
			break
		}
		if len(b) >= 2083+18 {
			return errors.New("HTTP header too long")
		}
	}
	bb := bytes.SplitN(b, []byte(" "), 3)
	if len(bb) != 3 {
		return errors.New("Invalid Request")
	}
	method, aoru := string(bb[0]), string(bb[1])
	var addr string
	if method == "CONNECT" {
		addr = aoru
	}
	if method != "CONNECT" {
		var err error
		addr, err = ant.GetAddressFromURL(aoru)
		if err != nil {
			return err
		}
	}

	if x.HTTPMiddleman != nil {
		if done, err := x.HTTPMiddleman.Handle(method, addr, b, c); err != nil || done {
			return err
		}
	}

	a, h, p, err := socks5.ParseAddress(addr)
	if err != nil {
		return err
	}
	tmp, err := Dial.Dial("tcp", x.RemoteAddr)
	if err != nil {
		return err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if x.TCPTimeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(x.TCPTimeout) * time.Second); err != nil {
			return err
		}
	}
	if x.TCPDeadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(x.TCPDeadline) * time.Second)); err != nil {
			return err
		}
	}
	crc, err := x.WrapCipherConn(rc)
	if err != nil {
		return err
	}

	rawaddr := make([]byte, 0)
	rawaddr = append(rawaddr, a)
	rawaddr = append(rawaddr, h...)
	rawaddr = append(rawaddr, p...)
	if _, err := crc.Write(rawaddr); err != nil {
		return err
	}
	if method == "CONNECT" {
		_, err := c.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			return err
		}
	}
	if method != "CONNECT" {
		if _, err := crc.Write(b); err != nil {
			return err
		}
	}

	go func() {
		iv := make([]byte, aes.BlockSize)
		if _, err := io.ReadFull(crc, iv); err != nil {
			log.Println(err)
			return
		}
		_, _ = io.Copy(c, crc)
	}()
	_, _ = io.Copy(crc, c)
	return nil
}

// WrapChiperConn make a chiper conn
func (x *SSClient) WrapCipherConn(conn *net.TCPConn) (*CipherConn, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	if _, err := conn.Write(iv); err != nil {
		return nil, err
	}
	return NewCipherConn(conn, x.Password, iv)
}

// Encrypt data
func (x *SSClient) Encrypt(rawdata []byte) ([]byte, error) {
	return ant.AESCFBEncrypt(rawdata, x.Password)
}

// Decrypt data
func (x *SSClient) Decrypt(cd []byte) (a byte, addr, port, data []byte, err error) {
	var bb []byte
	bb, err = ant.AESCFBDecrypt(cd, x.Password)
	if err != nil {
		return
	}
	err = errors.New("Data length error")
	n := len(bb)
	minl := 1
	if n < minl {
		return
	}
	if bb[0] == socks5.ATYPIPv4 {
		minl += 4
		if n < minl {
			return
		}
		addr = bb[minl-4 : minl]
	} else if bb[0] == socks5.ATYPIPv6 {
		minl += 16
		if n < minl {
			return
		}
		addr = bb[minl-16 : minl]
	} else if bb[0] == socks5.ATYPDomain {
		minl += 1
		if n < minl {
			return
		}
		l := bb[1]
		if l == 0 {
			return
		}
		minl += int(l)
		if n < minl {
			return
		}
		addr = bb[minl-int(l) : minl]
		addr = append([]byte{l}, addr...)
	} else {
		return
	}
	minl += 2
	if n <= minl {
		return
	}
	a = bb[0]
	port = bb[minl-2 : minl]
	data = bb[minl:]
	err = nil
	return
}

// Shutdown used to stop the client
func (x *SSClient) Shutdown() error {
	return x.Server.Stop()
}
