package brook

import (
	"errors"
	"net"

	"github.com/miekg/dns"
)

func Resolve6(host string) (string, error) {
	if net.ParseIP(host).To4() != nil {
		return "", errors.New("This is ipv4")
	}
	if net.ParseIP(host).To16() != nil {
		return host, nil
	}
	m := &dns.Msg{}
	m.SetQuestion(host+".", dns.TypeAAAA)
	r, err := dns.Exchange(m, "[2001:4860:4860::8888]:53")
	if err != nil {
		return "", err
	}
	for _, v := range r.Answer {
		if t, ok := v.(*dns.AAAA); ok {
			return t.AAAA.String(), nil
		}
	}
	return "", errors.New("Can not find IP")
}

func Resolve4(host string) (string, error) {
	if net.ParseIP(host).To4() != nil {
		return host, nil
	}
	if net.ParseIP(host).To16() != nil {
		return "", errors.New("This is ipv6")
	}
	m := &dns.Msg{}
	m.SetQuestion(host+".", dns.TypeA)
	r, err := dns.Exchange(m, "8.8.8.8:53")
	if err != nil {
		return "", err
	}
	for _, v := range r.Answer {
		if t, ok := v.(*dns.A); ok {
			return t.A.String(), nil
		}
	}
	return "", errors.New("Can not find IP")
}
