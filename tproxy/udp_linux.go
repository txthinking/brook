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

package tproxy

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/txthinking/x"
)

func ListenUDP(network string, laddr *net.UDPAddr) (*net.UDPConn, error) {
	c, err := net.ListenUDP(network, laddr)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	f, err := c.File()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fd := int(f.Fd())
	if err := syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
		return nil, err
	}
	if err = syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_RECVORIGDSTADDR, 1); err != nil {
		return nil, err
	}
	tmp, err := net.FileConn(f)
	if err != nil {
		return nil, err
	}
	return tmp.(*net.UDPConn), nil
}

func ReadFromUDP(conn *net.UDPConn, b []byte) (int, *net.UDPAddr, *net.UDPAddr, error) {
	oob := make([]byte, 1024)
	n, oobn, _, addr, err := conn.ReadMsgUDP(b, oob)
	if err != nil {
		return 0, nil, nil, err
	}

	msgs, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		return 0, nil, nil, err
	}

	var originalDst *net.UDPAddr
	for _, msg := range msgs {
		if msg.Header.Level != syscall.SOL_IP || msg.Header.Type != syscall.IP_RECVORIGDSTADDR {
			continue
		}
		originalDstRaw := &syscall.RawSockaddrInet4{}
		if err := binary.Read(bytes.NewReader(msg.Data), binary.LittleEndian, originalDstRaw); err != nil {
			return 0, nil, nil, err
		}

		switch originalDstRaw.Family {
		case syscall.AF_INET:
			pp := (*syscall.RawSockaddrInet4)(unsafe.Pointer(originalDstRaw))
			p := (*[2]byte)(unsafe.Pointer(&pp.Port))
			originalDst = &net.UDPAddr{
				IP:   net.IPv4(pp.Addr[0], pp.Addr[1], pp.Addr[2], pp.Addr[3]),
				Port: int(p[0])<<8 + int(p[1]),
			}
		case syscall.AF_INET6:
			pp := (*syscall.RawSockaddrInet6)(unsafe.Pointer(originalDstRaw))
			p := (*[2]byte)(unsafe.Pointer(&pp.Port))
			originalDst = &net.UDPAddr{
				IP:   net.IP(pp.Addr[:]),
				Port: int(p[0])<<8 + int(p[1]),
				Zone: strconv.Itoa(int(pp.Scope_id)),
			}
		default:
			return 0, nil, nil, nil
		}
	}
	if originalDst == nil {
		return 0, nil, nil, nil
	}
	return n, addr, originalDst, nil
}

func DialUDP(network string, laddr *net.UDPAddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	remoteSocketAddress, err := udpAddrToSocketAddr(raddr)
	if err != nil {
		return nil, err
	}

	localSocketAddress, err := udpAddrToSocketAddr(laddr)
	if err != nil {
		return nil, err
	}

	fd, err := syscall.Socket(udpAddrFamily(network, laddr, raddr), syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if err := syscall.Bind(fd, localSocketAddress); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if err := syscall.Connect(fd, remoteSocketAddress); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	f := os.NewFile(uintptr(fd), string(x.RandomNumber()))
	defer f.Close()

	c, err := net.FileConn(f)
	if err != nil {
		return nil, err
	}
	return c.(*net.UDPConn), nil
}

func udpAddrToSocketAddr(addr *net.UDPAddr) (syscall.Sockaddr, error) {
	switch {
	case addr.IP.To4() != nil:
		ip := [4]byte{}
		copy(ip[:], addr.IP.To4())

		return &syscall.SockaddrInet4{Addr: ip, Port: addr.Port}, nil

	default:
		ip := [16]byte{}
		copy(ip[:], addr.IP.To16())

		zoneID, err := strconv.ParseUint(addr.Zone, 10, 32)
		if err != nil {
			return nil, err
		}

		return &syscall.SockaddrInet6{Addr: ip, Port: addr.Port, ZoneId: uint32(zoneID)}, nil
	}
}

func udpAddrFamily(net string, laddr, raddr *net.UDPAddr) int {
	switch net[len(net)-1] {
	case '4':
		return syscall.AF_INET
	case '6':
		return syscall.AF_INET6
	}

	if (laddr == nil || laddr.IP.To4() != nil) &&
		(raddr == nil || raddr.IP.To4() != nil) {
		return syscall.AF_INET
	}
	return syscall.AF_INET6
}
