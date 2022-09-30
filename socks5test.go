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
	"errors"
	"fmt"
	"log"

	"github.com/miekg/dns"
	"github.com/txthinking/socks5"
)

func Socks5Test(s, u, p, domain, a, ds string) error {
	s5c, err := socks5.NewClient(s, u, p, 0, 60)
	if err != nil {
		return err
	}

	fmt.Println("Testing UDP: query " + domain + " A on " + ds)
	uc, err := s5c.Dial("udp", ds)
	if err != nil {
		return err
	}
	defer uc.Close()
	m := &dns.Msg{}
	m.RecursionDesired = true
	m.SetQuestion(domain+".", dns.TypeA)
	b, err := m.Pack()
	if err != nil {
		return err
	}
	if _, err := uc.Write(b); err != nil {
		return err
	}
	log.Printf("Sent Datagram. %#v\n", b)
	b = make([]byte, 512)
	i, err := uc.Read(b)
	if err != nil {
		return err
	}
	if err := m.Unpack(b[:i]); err != nil {
		return err
	}
	if len(m.Answer) == 0 {
		return errors.New("no answer")
	}
	v, ok := m.Answer[0].(*dns.A)
	if !ok {
		return errors.New("invalid answer")
	}
	if v.A.String() != a {
		fmt.Println("Expect", a, "but got", v.A.String())
	}
	if v.A.String() == a {
		fmt.Println("UDP: OK")
	}

	fmt.Println("Testing TCP: query " + domain + " A on " + ds)
	c := &dns.Client{Net: "tcp"}
	tc, err := s5c.Dial("tcp", ds)
	if err != nil {
		return err
	}
	defer tc.Close()
	m = &dns.Msg{}
	m.RecursionDesired = true
	m.SetQuestion(domain+".", dns.TypeA)
	m, _, err = c.ExchangeWithConn(m, &dns.Conn{Conn: tc})
	if err != nil {
		return err
	}
	if len(m.Answer) == 0 {
		return errors.New("no answer")
	}
	v, ok = m.Answer[0].(*dns.A)
	if !ok {
		return errors.New("invalid answer")
	}
	if v.A.String() != a {
		fmt.Println("Expect", a, "but got", v.A.String())
	}
	if v.A.String() == a {
		fmt.Println("TCP: OK")
	}
	return nil
}
