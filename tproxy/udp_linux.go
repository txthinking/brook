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
	"errors"
	"net"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/txthinking/x"
	"golang.org/x/sys/unix"
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
	if laddr.IP.To4() != nil {
		if err := syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
			return nil, err
		}
		if err = syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_RECVORIGDSTADDR, 1); err != nil {
			return nil, err
		}
	} else {
		if err := syscall.SetsockoptInt(fd, syscall.SOL_IPV6, unix.IPV6_TRANSPARENT, 1); err != nil {
			return nil, err
		}
		if err = syscall.SetsockoptInt(fd, syscall.SOL_IPV6, unix.IPV6_ORIGDSTADDR, 1); err != nil {
			return nil, err
		}
	}
	tmp, err := net.FileConn(f)
	if err != nil {
		return nil, err
	}
	return tmp.(*net.UDPConn), nil
}

func ReadFromUDP(conn *net.UDPConn, oob, b []byte) (int, *net.UDPAddr, *net.UDPAddr, error) {
	n, oobn, _, addr, err := conn.ReadMsgUDP(b, oob)
	if err != nil {
		return 0, nil, nil, err
	}
	msgs, err := syscall.ParseSocketControlMessage(oob[:oobn])
	if err != nil {
		return 0, nil, nil, err
	}
	for _, msg := range msgs {
		if msg.Header.Level == syscall.SOL_IP && msg.Header.Type == syscall.IP_RECVORIGDSTADDR {
			pp := &unix.RawSockaddrInet4{}
			ai := *(*uint16)(unsafe.Pointer(&msg.Data[unsafe.Offsetof(pp.Family)]))
			if ai != unix.AF_INET {
				return 0, nil, nil, errors.New("!AF_INET")
			}
			if err = binary.Read(bytes.NewReader(msg.Data), binary.BigEndian, pp); err != nil {
				return 0, nil, nil, err
			}
			dst := &net.UDPAddr{
				IP:   net.IPv4(pp.Addr[0], pp.Addr[1], pp.Addr[2], pp.Addr[3]),
				Port: int(pp.Port),
			}
			return n, addr, dst, nil
		}
		if msg.Header.Level == unix.SOL_IPV6 && msg.Header.Type == unix.IPV6_ORIGDSTADDR {
			pp := &unix.RawSockaddrInet6{}
			ai := *(*uint16)(unsafe.Pointer(&msg.Data[unsafe.Offsetof(pp.Family)]))
			if ai != unix.AF_INET6 {
				return 0, nil, nil, errors.New("!AF_INET6")
			}
			scopeId := *(*uint32)(unsafe.Pointer(&msg.Data[unsafe.Offsetof(pp.Scope_id)]))
			if err = binary.Read(bytes.NewReader(msg.Data), binary.BigEndian, pp); err != nil {
				return 0, nil, nil, err
			}
			dst := &net.UDPAddr{
				IP:   net.IP(pp.Addr[:]),
				Port: int(pp.Port),
				Zone: strconv.Itoa(int(scopeId)),
			}
			return n, addr, dst, nil
		}
	}
	return 0, nil, nil, nil
}

func DialUDP(network string, laddr *net.UDPAddr, raddr *net.UDPAddr) (*net.UDPConn, error) {
	var laddrs, raddrs syscall.Sockaddr
	var ai int
	if laddr.IP.To4() != nil {
		lip := [4]byte{}
		copy(lip[:], laddr.IP.To4())
		laddrs = &syscall.SockaddrInet4{Addr: lip, Port: laddr.Port}
		rip := [4]byte{}
		copy(rip[:], raddr.IP.To4())
		raddrs = &syscall.SockaddrInet4{Addr: rip, Port: raddr.Port}
		ai = syscall.AF_INET
	}
	if laddr.IP.To4() == nil {
		lip := [16]byte{}
		copy(lip[:], laddr.IP.To16())
		zoneID, err := strconv.ParseUint(laddr.Zone, 10, 32)
		if err != nil {
			i, err := net.InterfaceByName(laddr.Zone)
			if err != nil {
				return nil, err
			}
			zoneID = uint64(i.Index)
		}
		laddrs = &syscall.SockaddrInet6{Addr: lip, Port: laddr.Port, ZoneId: uint32(zoneID)}
		rip := [16]byte{}
		copy(rip[:], raddr.IP.To16())
		zoneID, err = strconv.ParseUint(raddr.Zone, 10, 32)
		if err != nil {
			i, err := net.InterfaceByName(raddr.Zone)
			if err != nil {
				return nil, err
			}
			zoneID = uint64(i.Index)
		}
		raddrs = &syscall.SockaddrInet6{Addr: rip, Port: raddr.Port, ZoneId: uint32(zoneID)}
		ai = syscall.AF_INET6
	}
	fd, err := syscall.Socket(ai, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		syscall.Close(fd)
		return nil, err
	}
	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, unix.SO_REUSEPORT, 1); err != nil {
		syscall.Close(fd)
		return nil, err
	}
	if ai == syscall.AF_INET {
		if err := syscall.SetsockoptInt(fd, syscall.SOL_IP, syscall.IP_TRANSPARENT, 1); err != nil {
			syscall.Close(fd)
			return nil, err
		}
	}
	if ai == syscall.AF_INET6 {
		if err := syscall.SetsockoptInt(fd, syscall.SOL_IPV6, unix.IPV6_TRANSPARENT, 1); err != nil {
			syscall.Close(fd)
			return nil, err
		}
	}
	if err := syscall.Bind(fd, laddrs); err != nil {
		syscall.Close(fd)
		return nil, err
	}
	if err := syscall.Connect(fd, raddrs); err != nil {
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
