package brook

import (
	"io"
	"log"
	"net"
	"time"
)

// Relay is a relay for two network
type Relay struct {
	Address  string
	Server   string
	Timeout  int
	Deadline int
	Listen   net.Listener
}

// NewRelay returns a new Relay
func NewRelay(address, server string, timeout, deadline int) *Relay {
	r := &Relay{
		Address:  address,
		Server:   server,
		Timeout:  timeout,
		Deadline: deadline,
	}
	return r
}

// ListenAndServe will let relay start to listen and serve
func (r *Relay) ListenAndServe() error {
	ta, err := net.ResolveTCPAddr("tcp", r.Address)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", ta)
	if err != nil {
		return err
	}
	r.Listen = l

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		go func(conn *net.TCPConn) {
			if err := r.handle(conn); err != nil {
				log.Println(err)
			}
		}(conn)
	}
}

// Shutdown used to stop the relay
func (r *Relay) Shutdown() error {
	if r.Listen == nil {
		return nil
	}
	return r.Listen.Close()
}

func (r *Relay) handle(conn *net.TCPConn) error {
	defer conn.Close()
	if r.Timeout != 0 {
		if err := conn.SetKeepAlivePeriod(time.Duration(r.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if r.Deadline != 0 {
		if err := conn.SetDeadline(time.Now().Add(time.Duration(r.Deadline) * time.Second)); err != nil {
			return err
		}
	}
	ta, err := net.ResolveTCPAddr("tcp", r.Server)
	if err != nil {
		return err
	}
	rc, err := net.DialTCP("tcp", nil, ta)
	if err != nil {
		return err
	}
	defer rc.Close()
	if r.Timeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(r.Timeout) * time.Second); err != nil {
			return err
		}
	}
	if r.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(r.Deadline) * time.Second)); err != nil {
			return err
		}
	}
	go func() {
		_, _ = io.Copy(conn, rc)
	}()
	_, _ = io.Copy(rc, conn)
	return nil
}
