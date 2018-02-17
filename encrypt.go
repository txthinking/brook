package brook

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/txthinking/ant"
	"github.com/txthinking/brook/plugin"
	"github.com/txthinking/socks5"
)

// IncrementNonce loves your compute to use Litter Endian
func IncrementNonce(n []byte) []byte {
	i := int(binary.LittleEndian.Uint16(n))
	i += 1
	n = make([]byte, 12)
	binary.LittleEndian.PutUint16(n, uint16(i))
	return n
}

// ReadFrom
func ReadFrom(c *net.TCPConn, k, n []byte, hasTime bool) ([]byte, []byte, error) {
	b := make([]byte, 18)
	if _, err := io.ReadFull(c, b); err != nil {
		return nil, nil, err
	}
	n = IncrementNonce(n)
	d, err := ant.AESGCMDecrypt(b, k, n)
	if err != nil {
		return nil, nil, err
	}

	l := int(binary.BigEndian.Uint16(d))
	b = make([]byte, l)
	if _, err := io.ReadFull(c, b); err != nil {
		return nil, nil, err
	}
	n = IncrementNonce(n)
	d, err = ant.AESGCMDecrypt(b, k, n)
	if err != nil {
		return nil, nil, err
	}

	if hasTime {
		i, err := strconv.Atoi(string(d[0:10]))
		if err != nil {
			return nil, nil, err
		}
		if time.Now().Unix()-int64(i) > 90 {
			time.Sleep(time.Duration(ant.Random(1, 60*10)) * time.Second)
			return nil, nil, errors.New("Expired request")
		}
		d = d[10:]
	}
	return d, n, nil
}

// WriteTo
func WriteTo(c *net.TCPConn, d, k, n []byte, needTime bool) ([]byte, error) {
	if needTime {
		d = append(bytes.NewBufferString(strconv.Itoa(int(time.Now().Unix()))).Bytes(), d...)
	}

	i := len(d) + 16
	bb := make([]byte, 2)
	binary.BigEndian.PutUint16(bb, uint16(i))
	n = IncrementNonce(n)
	b, err := ant.AESGCMEncrypt(bb, k, n)
	if err != nil {
		return nil, err
	}
	if _, err := c.Write(b); err != nil {
		return nil, err
	}

	n = IncrementNonce(n)
	b, err = ant.AESGCMEncrypt(d, k, n)
	if err != nil {
		return nil, err
	}
	if _, err := c.Write(b); err != nil {
		return nil, err
	}
	return n, nil
}

// PrepareKey
func PrepareKey(p []byte) ([]byte, []byte, error) {
	return ant.HkdfSha256RandomSalt(p, []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}, 12)
}

// GetKey
func GetKey(p, n []byte) ([]byte, error) {
	return ant.HkdfSha256WithSalt(p, n, []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b})
}

// Encrypt data
func Encrypt(p, b []byte) ([]byte, error) {
	b = append(bytes.NewBufferString(strconv.Itoa(int(time.Now().Unix()))).Bytes(), b...)
	k, n, err := PrepareKey(p)
	if err != nil {
		return nil, err
	}
	b, err = ant.AESGCMEncrypt(b, k, n)
	if err != nil {
		return nil, err
	}
	b = append(n, b...)
	return b, nil
}

// Decrypt data
func Decrypt(p, b []byte, token plugin.TokenChecker) (a byte, addr, port, data []byte, err error) {
	err = errors.New("Data length error")
	if len(b) <= 12+16 {
		return
	}
	k, err := GetKey(p, b[0:12])
	bb, err := ant.AESGCMDecrypt(b[12:], k, b[0:12])
	if err != nil {
		return
	}
	i, err := strconv.Atoi(string(bb[0:10]))
	if err != nil {
		return
	}
	if time.Now().Unix()-int64(i) > 90 {
		time.Sleep(time.Duration(ant.Random(1, 60*10)) * time.Second)
		err = errors.New("Expired request")
		return
	}
	bb = bb[10:]
	if token != nil {
		l := int(binary.BigEndian.Uint16(bb[0:2]))
		t := bb[2 : l+2]
		if err = token.Check(t); err != nil {
			return
		}
		bb = b[l+2:]
	}
	a = bb[0]
	if a == socks5.ATYPIPv4 {
		addr = bb[1:5]
		port = bb[5:7]
		data = bb[7:]
	} else if a == socks5.ATYPIPv6 {
		addr = bb[1:17]
		port = bb[17:19]
		data = bb[19:]
	} else if a == socks5.ATYPDomain {
		l := int(bb[1])
		addr = bb[1 : 1+l]
		port = bb[1+l : 1+l+2]
		data = bb[1+l+2:]
	} else {
		return
	}
	err = nil
	return
}
