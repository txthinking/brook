package brook

import (
	"net"
	"net/url"
	"regexp"
)

func GetAddressFromURL(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	if HasPort(u.Host) {
		return u.Host, nil
	}
	return net.JoinHostPort(u.Host, "80"), nil
}

func HasPort(host string) bool {
	r := regexp.MustCompile(`.+:\d+$`)
	return r.MatchString(host)
}
