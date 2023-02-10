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
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"

	"github.com/txthinking/socks5"
	"github.com/txthinking/x"
	"golang.org/x/crypto/hkdf"
)

type StreamServer struct {
	Client  net.Conn
	cn      []byte
	ca      cipher.AEAD
	sn      []byte
	sa      cipher.AEAD
	RB      []byte
	WB      []byte
	Timeout int
	network string
	src     string
	dst     string
}

func NewStreamServer(password []byte, src string, client net.Conn, timeout, udptimeout int) (Exchanger, error) {
	if timeout != 0 {
		if err := client.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second)); err != nil {
			return nil, err
		}
	}
	s := &StreamServer{Client: client, Timeout: timeout, src: src}
	s.cn = x.BP12.Get().([]byte)
	if _, err := io.ReadFull(s.Client, s.cn); err != nil {
		x.BP12.Put(s.cn)
		return nil, err
	}
	ck := x.BP32.Get().([]byte)
	if _, err := io.ReadFull(hkdf.New(sha256.New, password, s.cn, []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), ck); err != nil {
		x.BP12.Put(s.cn)
		x.BP32.Put(ck)
		return nil, err
	}
	cb, err := aes.NewCipher(ck)
	if err != nil {
		x.BP12.Put(s.cn)
		x.BP32.Put(ck)
		return nil, err
	}
	x.BP32.Put(ck)
	s.ca, err = cipher.NewGCM(cb)
	if err != nil {
		x.BP12.Put(s.cn)
		return nil, err
	}

	s.RB = x.BP2048.Get().([]byte)
	l, err := s.Read()
	if err != nil {
		x.BP12.Put(s.cn)
		x.BP2048.Put(s.RB)
		return nil, err
	}
	i := int64(binary.BigEndian.Uint32(s.RB[2+16 : 2+16+4]))
	if time.Now().Unix()-i > 60 {
		x.BP12.Put(s.cn)
		x.BP2048.Put(s.RB)
		WaitReadErr(s.Client)
		return nil, errors.New("Expired request")
	}
	if i%2 == 0 {
		s.network = "tcp"
	}
	if i%2 == 1 {
		s.network = "udp"
		s.Timeout = udptimeout
	}

	s.sn = x.BP12.Get().([]byte)
	if _, err := io.ReadFull(rand.Reader, s.sn); err != nil {
		x.BP12.Put(s.cn)
		x.BP2048.Put(s.RB)
		x.BP12.Put(s.sn)
		return nil, err
	}
	sk := x.BP32.Get().([]byte)
	if _, err := io.ReadFull(hkdf.New(sha256.New, password, s.sn, []byte{0x62, 0x72, 0x6f, 0x6f, 0x6b}), sk); err != nil {
		x.BP12.Put(s.cn)
		x.BP2048.Put(s.RB)
		x.BP12.Put(s.sn)
		x.BP32.Put(sk)
		return nil, err
	}
	if _, err := s.Client.Write(s.sn); err != nil {
		x.BP12.Put(s.cn)
		x.BP2048.Put(s.RB)
		x.BP12.Put(s.sn)
		x.BP32.Put(sk)
		return nil, err
	}
	sb, err := aes.NewCipher(sk)
	if err != nil {
		x.BP12.Put(s.cn)
		x.BP2048.Put(s.RB)
		x.BP12.Put(s.sn)
		x.BP32.Put(sk)
		return nil, err
	}
	x.BP32.Put(sk)
	s.sa, err = cipher.NewGCM(sb)
	if err != nil {
		x.BP12.Put(s.cn)
		x.BP2048.Put(s.RB)
		x.BP12.Put(s.sn)
		return nil, err
	}

	if s.network == "tcp" {
		s.WB = x.BP2048.Get().([]byte)
	}
	if s.network == "udp" {
		RB := x.BP65507.Get().([]byte)
		copy(RB[2+16+4:2+16+l], s.RB[2+16+4:2+16+l])
		x.BP2048.Put(s.RB)
		s.RB = RB
		s.WB = x.BP65507.Get().([]byte)
	}
	s.dst = socks5.ToAddress(s.RB[2+16+4], s.RB[2+16+4+1:2+16+l-2], s.RB[2+16+l-2:])
	return ServerGate(s)
}

func (s *StreamServer) Exchange(remote net.Conn) error {
	go func() {
		for {
			if s.Timeout != 0 {
				if err := remote.SetDeadline(time.Now().Add(time.Duration(s.Timeout) * time.Second)); err != nil {
					return
				}
			}
			l, err := remote.Read(s.WB[2+16 : len(s.WB)-16])
			if err != nil {
				return
			}
			if err := s.Write(l); err != nil {
				return
			}
		}
	}()
	for {
		if s.Timeout != 0 {
			if err := s.Client.SetDeadline(time.Now().Add(time.Duration(s.Timeout) * time.Second)); err != nil {
				return nil
			}
		}
		l, err := s.Read()
		if err != nil {
			return nil
		}
		if _, err := remote.Write(s.RB[2+16 : 2+16+l]); err != nil {
			return nil
		}
	}
	return nil
}

func (s *StreamServer) Write(l int) error {
	binary.BigEndian.PutUint16(s.WB[:2], uint16(l))
	s.sa.Seal(s.WB[:0], s.sn, s.WB[:2], nil)
	NextNonce(s.sn)
	s.sa.Seal(s.WB[:2+16], s.sn, s.WB[2+16:2+16+l], nil)
	if _, err := s.Client.Write(s.WB[:2+16+l+16]); err != nil {
		return err
	}
	NextNonce(s.sn)
	return nil
}

func (s *StreamServer) Read() (int, error) {
	if _, err := io.ReadFull(s.Client, s.RB[:2+16]); err != nil {
		return 0, err
	}
	if _, err := s.ca.Open(s.RB[:0], s.cn, s.RB[:2+16], nil); err != nil {
		WaitReadErr(s.Client)
		return 0, err
	}
	l := int(binary.BigEndian.Uint16(s.RB[:2]))
	if _, err := io.ReadFull(s.Client, s.RB[2+16:2+16+l+16]); err != nil {
		return 0, err
	}
	NextNonce(s.cn)
	if _, err := s.ca.Open(s.RB[:2+16], s.cn, s.RB[2+16:2+16+l+16], nil); err != nil {
		return 0, err
	}
	NextNonce(s.cn)
	return l, nil
}

func (s *StreamServer) Clean() {
	x.BP12.Put(s.cn)
	x.BP12.Put(s.sn)
	if s.network == "tcp" {
		x.BP2048.Put(s.WB)
		x.BP2048.Put(s.RB)
	}
	if s.network == "udp" {
		x.BP65507.Put(s.WB)
		x.BP65507.Put(s.RB)
	}
}

func (s *StreamServer) Network() string {
	return s.network
}

func (s *StreamServer) Src() string {
	return s.src
}

func (s *StreamServer) Dst() string {
	return s.dst
}
