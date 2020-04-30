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
	"os/exec"
	"syscall"

	"github.com/txthinking/brook/sysproxy"
)

// AddRoutes adds routes.
func (v *Tun) AddRoutes() error {
	c := exec.Command("chcp", "65001")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("route", "add", "0.0.0.0", "mask", "128.0.0.0", v.TunGateway)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("route", "add", "128.0.0.0", "mask", "128.0.0.0", v.TunGateway)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	gw, err := sysproxy.GetDefaultGateway()
	if err != nil {
		return err
	}
	c = exec.Command("route", "add", v.ServerIP, "mask", "255.255.255.255", gw, "metric", "1")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}

// DeleteRoutes deletes routes.
func (v *Tun) DeleteRoutes() error {
	c := exec.Command("chcp", "65001")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("route", "delete", "0.0.0.0", "mask", "128.0.0.0", v.TunGateway)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("route", "delete", "128.0.0.0", "mask", "128.0.0.0", v.TunGateway)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	gw, err := sysproxy.GetDefaultGateway()
	if err != nil {
		return err
	}
	c = exec.Command("route", "delete", v.ServerIP, "mask", "255.255.255.255", gw, "metric", "1")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}
