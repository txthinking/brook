package brook

import (
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/txthinking/socks5"
)

// Socks5Server is the client of raw socks5 protocol
type Socks5Server struct {
	Address         string
	UserName        string
	Password        string
	Timeout         int
	Deadline        int
	Dial            Dialer
	Listen          net.Listener
	Socks5Middleman Socks5Middleman
}

// NewSocks5Server returns a new Socks5Server, dial can be set to nil
func NewSocks5Server(address, userName, password string, timeout, deadline int, dial Dialer) (*Socks5Server, error) {
	if dial == nil {
		dial = &DefaultDial{}
	}
	c := &Socks5Server{
		Address:  address,
		UserName: userName,
		Password: password,
		Timeout:  timeout,
		Deadline: deadline,
		Dial:     dial,
	}
	return c, nil
}

// ListenAndServe will let client start to listen and serve, sm can be nil
func (c *Socks5Server) ListenAndServe(sm Socks5Middleman) error {
	ta, err := net.ResolveTCPAddr("tcp", c.Address)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", ta)
	if err != nil {
		return err
	}
	defer l.Close()
	c.Listen = l
	c.Socks5Middleman = sm
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		go func(conn *net.TCPConn) {
			if err := c.handle(conn); err != nil {
				log.Println(err)
			}
		}(conn)
	}
}

// ListenAndForward will let client start to listen and forward to another socks5, sm can be nil
func (c *Socks5Server) ListenAndForward(sm Socks5Middleman, addr string) error {
	ta, err := net.ResolveTCPAddr("tcp", c.Address)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", ta)
	if err != nil {
		return err
	}
	defer l.Close()
	c.Listen = l
	c.Socks5Middleman = sm
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		go func(conn *net.TCPConn, addr string) {
			if err := c.handleForward(conn, addr); err != nil {
				log.Println(err)
			}
		}(conn, addr)
	}
}

// Shutdown used to stop the client
func (c *Socks5Server) Shutdown() error {
	if c.Listen == nil {
		return nil
	}
	return c.Listen.Close()
}

func (c *Socks5Server) handle(conn *net.TCPConn) error {
	defer conn.Close()
	if c.Timeout != 0 {
		if err := conn.SetKeepAlivePeriod(time.Duration(c.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if c.Deadline != 0 {
		if err := conn.SetDeadline(time.Now().Add(time.Duration(c.Deadline) * time.Second)); err != nil {
			return err
		}
	}

	mt := socks5.MethodNone
	if c.UserName != "" && c.Password != "" {
		mt = socks5.MethodUsernamePassword
	}
	s := &socks5.Server{
		C: conn,
		SelectMethod: func(methods []byte) (method byte, got bool) {
			for _, m := range methods {
				if m == mt {
					method = mt
					got = true
					return
				}
			}
			return
		},
		SupportedCommands: []byte{socks5.CmdConnect},
	}
	if mt == socks5.MethodUsernamePassword {
		s.CheckUserPass = func(user, pass []byte) bool {
			if string(user) == c.UserName && string(pass) == c.Password {
				return true
			}
			return false
		}
	}

	if err := s.Negotiate(); err != nil {
		return err
	}
	request, err := s.GetRequest()
	if err != nil {
		return err
	}

	if c.Socks5Middleman != nil {
		if handled, err := c.Socks5Middleman.HandleSocks5Proxy(request, conn); err != nil || handled {
			return err
		}
	}
	rc, err := request.Connect(conn)
	if err != nil {
		return err
	}
	defer rc.Close()
	go func() {
		_, _ = io.Copy(conn, rc)
	}()
	_, _ = io.Copy(rc, conn)
	return nil
}

func (c *Socks5Server) handleForward(conn *net.TCPConn, addr string) error {
	defer conn.Close()
	if c.Timeout != 0 {
		if err := conn.SetKeepAlivePeriod(time.Duration(c.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if c.Deadline != 0 {
		if err := conn.SetDeadline(time.Now().Add(time.Duration(c.Deadline) * time.Second)); err != nil {
			return err
		}
	}

	s := socks5.NewClassicServer(conn)
	if err := s.Negotiate(); err != nil {
		return err
	}
	request, err := s.GetRequest()
	if err != nil {
		return err
	}

	if c.Socks5Middleman != nil {
		if handled, err := c.Socks5Middleman.HandleSocks5Proxy(request, conn); err != nil || handled {
			return err
		}
	}

	rc, err := c.Dial.Dial("tcp", addr)
	if err != nil {
		var p *socks5.Reply
		if request.Atyp == socks5.ATYPIPv4 || request.Atyp == socks5.ATYPDomain {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
		} else {
			p = socks5.NewReply(socks5.RepConnectionRefused, socks5.ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
		}
		if err := p.WriteTo(conn); err != nil {
			return err
		}
		return err
	}
	defer rc.Close()
	if c.Timeout != 0 {
		if rtc, ok := rc.(*net.TCPConn); ok {
			if err := rtc.SetKeepAlivePeriod(time.Duration(c.Timeout) * time.Second); err != nil {
				return err
			}
		}
	}
	if c.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(c.Deadline) * time.Second)); err != nil {
			return err
		}
	}

	mt := socks5.MethodNone
	if c.UserName != "" && c.Password != "" {
		mt = socks5.MethodUsernamePassword
	}
	rq := socks5.NewNegotiationRequest([]byte{mt})
	if err := rq.WriteTo(rc); err != nil {
		return err
	}
	rp, err := socks5.NewNegotiationReplyFrom(rc)
	if err != nil {
		return err
	}
	if rp.Method != mt {
		return errors.New("Unsupport method")
	}
	if mt == socks5.MethodUsernamePassword {
		urq := socks5.NewUserPassNegotiationRequest([]byte(c.UserName), []byte(c.Password))
		if err := urq.WriteTo(rc); err != nil {
			return err
		}
		urp, err := socks5.NewUserPassNegotiationReplyFrom(rc)
		if err != nil {
			return err
		}
		if urp.Status != socks5.UserPassStatusSuccess {
			return socks5.ErrUserPassAuth
		}
	}
	if err := request.WriteTo(rc); err != nil {
		return err
	}
	reply, err := socks5.NewReplyFrom(rc)
	if err != nil {
		return err
	}
	if reply.Rep != socks5.RepSuccess {
		return errors.New("Host unreachable")
	}
	a, address, port, err := socks5.ParseAddress(rc.LocalAddr().String())
	if err != nil {
		return err
	}
	reply = socks5.NewReply(socks5.RepSuccess, a, address, port)
	if err := reply.WriteTo(conn); err != nil {
		return err
	}

	go func() {
		_, _ = io.Copy(conn, rc)
	}()
	_, _ = io.Copy(rc, conn)
	return nil
}
