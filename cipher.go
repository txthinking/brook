package brook

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"net"
	"time"

	"github.com/txthinking/ant"
)

// CipherConn is the encrypted connection
type CipherConn struct {
	c  net.Conn
	sr cipher.StreamReader
	sw cipher.StreamWriter
}

// NewCipherConn returns a new CipherConn, iv length must be equal aes.BlockSize
func NewCipherConn(c net.Conn, key []byte, iv []byte) (*CipherConn, error) {
	if len(iv) != aes.BlockSize {
		return nil, errors.New("Invalid IV length")
	}
	block, err := aes.NewCipher(ant.AESMake256Key(key))
	if err != nil {
		return nil, err
	}
	return &CipherConn{
		c: c,
		sr: cipher.StreamReader{
			S: cipher.NewCFBDecrypter(block, iv),
			R: c,
		},
		sw: cipher.StreamWriter{
			S: cipher.NewCFBEncrypter(block, iv),
			W: c,
		},
	}, nil
}

// Read is just like net.Conn interface
func (c *CipherConn) Read(b []byte) (n int, err error) {
	return c.sr.Read(b)
}

// Write is just like net.Conn interface
func (c *CipherConn) Write(b []byte) (n int, err error) {
	return c.sw.Write(b)
}

// Close is just like net.Conn interface
func (c *CipherConn) Close() error {
	return c.c.Close()
}

// LocalAddr is just like net.Conn interface
func (c *CipherConn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}

// RemoteAddr is just like net.Conn interface
func (c *CipherConn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}

// SetDeadline is just like net.Conn interface
func (c *CipherConn) SetDeadline(t time.Time) error {
	return c.c.SetDeadline(t)
}

// SetReadDeadline is just like net.Conn interface
func (c *CipherConn) SetReadDeadline(t time.Time) error {
	return c.c.SetReadDeadline(t)
}

// SetWriteDeadline is just like net.Conn interface
func (c *CipherConn) SetWriteDeadline(t time.Time) error {
	return c.c.SetWriteDeadline(t)
}
