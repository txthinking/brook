package brook

import "github.com/txthinking/socks5"

// EnableDebug will print more logs
func EnableDebug() {
	socks5.Debug = true
}
