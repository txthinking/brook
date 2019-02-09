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

package sysproxy

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// GetNetworkInterface returns default interface dev name.
func GetNetworkInterface() (string, error) {
	return "", nil
}

// GetDefaultGateway returns default gateway.
func GetDefaultGateway() (string, error) {
	c := exec.Command("sh", "-c", "ip route | grep default | awk '{print $3}'")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	return strings.TrimSpace(string(out)), nil
}

// GetDNSServers used to get DNS servers.
func GetDNSServers() ([]string, error) {
	c := exec.Command("sh", "-c", "cat /etc/resolv.conf | grep -Pv \"^#\" | grep nameserver | awk '{print $2}'")
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + err.Error())
	}
	if strings.TrimSpace(string(out)) == "" {
		return []string{}, nil
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

// SetDNSServers used to set system DNS servers.
func SetDNSServers(servers []string) error {
	f, err := os.OpenFile("/etc/resolv.conf", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, v := range servers {
		if _, err := f.WriteString("nameserver " + v + "\n"); err != nil {
			return err
		}
	}
	return nil
}

func TurnOnSystemProxy(pac string) error {
	return nil
}

func TurnOffSystemProxy() error {
	return nil
}
