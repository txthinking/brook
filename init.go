package brook

import (
	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
)

var Dial x.Dialer = x.DefaultDial

// EnableDebug
func EnableDebug() {
	socks5.Debug = true
}
