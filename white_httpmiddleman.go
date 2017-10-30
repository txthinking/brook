package brook

import (
	"io"
	"net"
	"time"

	"github.com/txthinking/pac/blackwhite"
)

// WhiteHTTPMiddleman is a HTTPMiddleman who only handle domain in white list
type WhiteHTTPMiddleman struct {
	Timeout  int
	Deadline int
}

// NewWhiteHTTPMiddleman returns a WhiteHTTPMiddleman which can used to handle http proxy request
func NewWhiteHTTPMiddleman(timeout, deadline int) *WhiteHTTPMiddleman {
	return &WhiteHTTPMiddleman{
		Timeout:  timeout,
		Deadline: deadline,
	}
}

// Handle handle http proxy request, if the domain is in the white list
func (w *WhiteHTTPMiddleman) Handle(method, addr string, request []byte, conn *net.TCPConn) (handled bool, err error) {
	h, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false, err
	}
	if !blackwhite.IsWhite(h) {
		return false, nil
	}

	tmp, err := Dial.Dial("tcp", addr)
	if err != nil {
		return true, err
	}
	rc := tmp.(*net.TCPConn)
	defer rc.Close()
	if w.Timeout != 0 {
		if err := rc.SetKeepAlivePeriod(time.Duration(w.Timeout) * time.Second); err != nil {
			return true, err
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
