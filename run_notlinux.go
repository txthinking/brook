// +build !linux

package brook

import "errors"

func RunTproxy(address, server, password string, tcpTimeout, tcpDeadline, udpDeadline int) error {
	return errors.New("Only works on Linux")
}
