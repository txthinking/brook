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
	"bytes"
	"errors"
	"os/exec"
	"strings"
	"syscall"
)

// GetNetworkInterface returns interface name.
func GetNetworkInterface() (string, error) {
	c := exec.Command("chcp", "65001")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	c = exec.Command("netsh", "interface", "ip", "show", "address")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	l := strings.Split(string(out), "\n")
	for i, v := range l {
		if !strings.Contains(v, "Default Gateway") {
			continue
		}
		l1 := strings.Split(v, "Gateway:")
		if len(l1) != 2 {
			continue
		}
		for j := i - 1; j >= 0; j-- {
			if !strings.Contains(l[j], "Configuration for interface") {
				continue
			}
			l1 := strings.Split(l[j], "\"")
			if len(l1) != 3 {
				return "", errors.New("Can't find network interface")
			}
			return strings.TrimSpace(l1[1]), nil
		}
		return "", errors.New("Can't find network interface")
	}
	return "", errors.New("Can't find network interface")
}

// GetDefaultGateway returns default gateway.
func GetDefaultGateway() (string, error) {
	c := exec.Command("chcp", "65001")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	c = exec.Command("netsh", "interface", "ip", "show", "address")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	l := strings.Split(string(out), "\n")
	for _, v := range l {
		if !strings.Contains(v, "Default Gateway") {
			continue
		}
		l1 := strings.Split(v, "Gateway:")
		if len(l1) != 2 {
			continue
		}
		if strings.TrimSpace(l1[1]) == "" {
			continue
		}
		return strings.TrimSpace(l1[1]), nil
	}
	return "", errors.New("Can't find default gateway")
}

// GetDNSServers used to get DNS servers.
func GetDNSServers() ([]string, error) {
	s, err := GetNetworkInterface()
	if err != nil {
		return nil, err
	}
	c := exec.Command("netsh", "interface", "ip", "show", "dnsservers", s)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + err.Error())
	}
	if !strings.Contains(string(out), "Statically Configured DNS Servers:") {
		return []string{}, nil
	}
	l := strings.Split(string(out), "Statically Configured DNS Servers:")
	if len(l) != 2 {
		return nil, errors.New("Can't find dns servers")
	}
	l = strings.Split(l[1], "Register with which suffix:")
	if len(l) != 2 {
		return nil, errors.New("Can't find dns servers")
	}
	l = strings.Split(strings.TrimSpace(l[0]), "\n")
	l = []string{}
	for _, v := range l {
		l = append(l, strings.TrimSpace(v))
	}
	return l, nil
}

// SetDNSServers used to set system DNS server.
func SetDNSServers(servers []string) error {
	s, err := GetNetworkInterface()
	if err != nil {
		return err
	}
	if len(servers) == 0 {
		c := exec.Command("netsh", "interface", "ip", "set", "dnsservers", s, "static", "none")
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
		return nil
	}
	for i, v := range servers {
		if i == 0 {
			c := exec.Command("netsh", "interface", "ip", "set", "dnsservers", s, "static", v, "primary")
			c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			if out, err := c.CombinedOutput(); err != nil {
				return errors.New(string(out) + err.Error())
			}
			continue
		}
		c := exec.Command("netsh", "interface", "ip", "add", "dnsservers", s, v)
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
	}
	return nil
}

// TurnOnSystemProxy used to enable system pac proxy, pac is a URL.
func TurnOnSystemProxy(pac string) error {
	c := exec.Command("chcp", "65001")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command(`reg`, `add`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `AutoConfigURL`, `/t`, `REG_SZ`, `/d`, pac, `/f`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command(`reg`, `add`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `ProxyEnable`, `/t`, `REG_DWORD`, `/d`, `0`, `/f`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	if err := reloadWinProxy(); err != nil {
		return err
	}
	return nil
}

// TurnOffSystemProxy used to disable system pac proxy.
func TurnOffSystemProxy() error {
	c := exec.Command("chcp", "65001")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command(`reg`, `query`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return errors.New(string(out) + err.Error())
	}
	if bytes.Contains(bytes.ToLower(out), []byte("autoconfigurl")) {
		c := exec.Command(`reg`, `delete`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `AutoConfigURL`, `/f`)
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
	}
	if err := reloadWinProxy(); err != nil {
		return err
	}
	return nil
}

func reloadWinProxy() error {
	h, err := syscall.LoadLibrary("wininet.dll")
	if err != nil {
		return err
	}
	f, err := syscall.GetProcAddress(h, "InternetSetOptionW")
	if err != nil {
		return err
	}
	ret, _, errno := syscall.Syscall6(uintptr(f), 4, 0, 39, 0, 0, 0, 0)
	if ret != 1 {
		return errors.New(errno.Error())
	}
	ret, _, errno = syscall.Syscall6(uintptr(f), 4, 0, 37, 0, 0, 0, 0)
	if ret != 1 {
		return errors.New(errno.Error())
	}
	return nil
}
