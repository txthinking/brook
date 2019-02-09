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
	"os/exec"
	"strings"
)

// GetNetworkInterface returns default interface name, not dev name.
func GetNetworkInterface() (string, error) {
	c := exec.Command("sh", "-c", "networksetup -listnetworkserviceorder | grep $(route -n get default | grep interface | awk '{print $2}') | awk 'BEGIN {FS=\",\"}; {print $1}' | awk 'BEGIN {FS=\": \"}; {print $2}'")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	return strings.TrimSpace(string(out)), nil
}

// GetDefaultGateway returns default gateway.
func GetDefaultGateway() (string, error) {
	c := exec.Command("sh", "-c", "route -n get default | grep gateway | awk '{print $2}'")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	return strings.TrimSpace(string(out)), nil
}

// GetDNSServers used to get DNS servers.
func GetDNSServers() ([]string, error) {
	s, err := GetNetworkInterface()
	if err != nil {
		return nil, err
	}
	c := exec.Command("networksetup", "-getdnsservers", s)
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + err.Error())
	}
	if strings.Contains(string(out), "aren't") {
		return []string{}, nil
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

// SetDNSServers used to set DNS servers.
func SetDNSServers(servers []string) error {
	s, err := GetNetworkInterface()
	if err != nil {
		return err
	}
	if len(servers) != 0 {
		servers = append([]string{"-setdnsservers", s}, servers...)
	} else {
		servers = []string{"-setdnsservers", s, "empty"}
	}
	c := exec.Command("networksetup", servers...)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}

// TurnOnSystemProxy used to enable system pac proxy, pac is a URL.
func TurnOnSystemProxy(pac string) error {
	s, err := GetNetworkInterface()
	if err != nil {
		return err
	}
	c := exec.Command("networksetup", "-setautoproxyurl", s, pac)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("networksetup", "-setautoproxystate", s, "on")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}

// TurnOffSystemProxy used to disable system pac proxy.
func TurnOffSystemProxy() error {
	s, err := GetNetworkInterface()
	if err != nil {
		return err
	}
	c := exec.Command("networksetup", "-setautoproxystate", s, "off")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}
