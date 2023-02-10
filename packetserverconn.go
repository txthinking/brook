// Copyright (c) 2016-present Cloud <cloud@txthinking.com>
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of version 3 of the GNU General Public
// License as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package brook

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
	"golang.org/x/crypto/hkdf"
)

type PacketServerConnFactory struct {
	Conns map[string]*PacketConn
	Lock  *sync.Mutex
}

func NewPacketServerConnFactory() *PacketServerConnFactory {
	return &PacketServerConnFactory{
		Conns: make(map[string]*PacketConn),
		Lock:  &sync.Mutex{},
	}
}

func (f *PacketServerConnFactory) Handle(addr *net.UDPAddr, b, p []byte, w func([]byte) (int, error), timeout int) (net.Conn, []byte, error) {
	if len(b) < 12+4+16 {
		return nil, nil, errors.New("data too small")
	}
	ck := x.BP32.Get().([]byte)
	if _, err := io.ReadFull(hkdf.New(sha256.New, p, b[:12], []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), ck); err != nil {
		x.BP32.Put(ck)
		return nil, nil, err
	}

	cb, err := aes.NewCipher(ck)
	if err != nil {
		x.BP32.Put(ck)
		return nil, nil, err
	}
	x.BP32.Put(ck)
	ca, err := cipher.NewGCM(cb)
	if err != nil {
		return nil, nil, err
	}
	if _, err := ca.Open(b[:12], b[:12], b[12:], nil); err != nil {
		return nil, nil, err
	}
	i := int64(binary.BigEndian.Uint32(b[12 : 12+4]))
	if time.Now().Unix()-i > 60 {
		return nil, nil, errors.New("Expired request")
	}
	a, h, p, err := socks5.ParseBytesAddress(b[12+4:])
	if err != nil {
		return nil, nil, err
	}
	dst := socks5.ToAddress(a, h, p)
	f.Lock.Lock()
	c, ok := f.Conns[addr.String()+dst]
	f.Lock.Unlock()
	if ok {
		_ = c.In(b[12+4+1+len(h)+2 : len(b)-16])
		return nil, nil, nil
	}
	f.Lock.Lock()
	c = NewPacketConn(b[12+4+1+len(h)+2:len(b)-16], w, timeout, func() {
		f.Lock.Lock()
		delete(f.Conns, addr.String()+dst)
		f.Lock.Unlock()
	})
	f.Conns[addr.String()+dst] = c
	f.Lock.Unlock()
	return c, b[12+4 : 12+4+1+len(h)+2], nil
}
