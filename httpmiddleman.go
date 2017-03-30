package brook

import (
	"io"
	"net"
	"time"

	"github.com/txthinking/pac/blackwhite"
)

// HTTPMiddleman is a middleman who can intercept and handle request
type HTTPMiddleman interface {
	// HandleHTTPProxy does not need to close conn,
	// if handled is true or err is not nil that means the request has been handled
	HandleHTTPProxy(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error)
}

// WhiteHTTPMiddleman is a HTTPMiddleman who only handle domain in white list
type WhiteHTTPMiddleman struct {
	Timeout  int
	Deadline int
}

// NewWhiteHTTPMiddleman returns a WhiteHTTPMiddleman which can used to handle http proxy request
func NewWhiteHTTPMiddleman(timeout, deadline int) *WhiteHTTPMiddleman {
	blackwhite.InitWhiteList()
	return &WhiteHTTPMiddleman{
		Timeout:  timeout,
		Deadline: deadline,
	}
}

// HandleHTTPProxy handle http proxy request, if the domain is in the white list
func (w *WhiteHTTPMiddleman) HandleHTTPProxy(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false, err
	}
	if !blackwhite.IsWhite(host) {
		return false, nil
	}

	rc, err := net.Dial("tcp", addr)
	if err != nil {
		return true, err
	}
	defer rc.Close()
	if w.Timeout != 0 {
		if rtc, ok := rc.(*net.TCPConn); ok {
			if err := rtc.SetKeepAlivePeriod(time.Duration(w.Timeout) * time.Second); err != nil {
				return true, err
			}
		}
	}
	if w.Deadline != 0 {
		if err := rc.SetDeadline(time.Now().Add(time.Duration(w.Deadline) * time.Second)); err != nil {
			return true, err
		}
	}
	if method == "CONNECT" {
		_, err := conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			return true, err
		}
	}
	if method != "CONNECT" {
		if _, err := rc.Write(request); err != nil {
			return true, err
		}
	}
	go func() {
		_, _ = io.Copy(rc, conn)
	}()
	_, _ = io.Copy(conn, rc)
	return true, nil
}
