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
	"fmt"
)

func EchoClient(server string, times int) error {
	c, err := DialTCP("tcp", "", server)
	if err != nil {
		return err
	}
	defer c.Close()
	var b [1024 * 2]byte
	for i := 0; i < times; i++ {
		if _, err := c.Write([]byte(c.LocalAddr().String())); err != nil {
			return err
		}
		i, err := c.Read(b[:])
		if err != nil {
			return err
		}
		if c.LocalAddr().String() == string(b[:i]) {
			fmt.Printf("TCP: src:%s -> dst:%s\n", c.LocalAddr().String(), c.RemoteAddr().String())
			fmt.Printf("TCP: dst:%s <- src:%s\n", c.LocalAddr().String(), c.RemoteAddr().String())
		}
		if c.LocalAddr().String() != string(b[:i]) {
			fmt.Printf("TCP: src:%s -> dst:proxy -> src:proxy -> dst:%s\n", c.LocalAddr().String(), c.RemoteAddr().String())
			fmt.Printf("TCP: dst:%s <- src:proxy <- dst:%s <- src:%s\n", c.LocalAddr().String(), string(b[:i]), c.RemoteAddr().String())
		}
	}

	c1, err := DialUDP("udp", "", server)
	if err != nil {
		return err
	}
	defer c1.Close()
	for i := 0; i < times; i++ {
		if _, err := c1.Write([]byte(c1.LocalAddr().String())); err != nil {
			return err
		}
		i, err := c1.Read(b[:])
		if err != nil {
			return err
		}
		if c1.LocalAddr().String() == string(b[:i]) {
			fmt.Printf("UDP: src:%s -> dst:%s\n", c1.LocalAddr().String(), c1.RemoteAddr().String())
			fmt.Printf("UDP: dst:%s <- src:%s\n", c1.LocalAddr().String(), c1.RemoteAddr().String())
		}
		if c1.LocalAddr().String() != string(b[:i]) {
			fmt.Printf("UDP: src:%s -> dst:proxy -> src:proxy -> dst:%s\n", c1.LocalAddr().String(), c1.RemoteAddr().String())
			fmt.Printf("UDP: dst:%s <- src:proxy <- dst:%s <- src:%s\n", c1.LocalAddr().String(), string(b[:i]), c1.RemoteAddr().String())
		}
	}
	return nil
}
