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

	"github.com/txthinking/brook/sysproxy"
)

// AddRoutes adds routes.
func (v *Tun) AddRoutes() error {
	c := exec.Command("route", "add", "-net", "0.0.0.0", v.TunGateway, "-netmask", "128.0.0.0")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("route", "add", "-net", "128.0.0.0", v.TunGateway, "-netmask", "128.0.0.0")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	gw, err := sysproxy.GetDefaultGateway()
	if err != nil {
		return err
	}
	c = exec.Command("route", "add", "-host", v.ServerIP, gw, "-netmask", "255.255.255.255")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}

// DeleteRoutes deletes routes.
func (v *Tun) DeleteRoutes() error {
	c := exec.Command("route", "delete", "-net", "0.0.0.0", v.TunGateway, "-netmask", "128.0.0.0")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("route", "delete", "-net", "128.0.0.0", v.TunGateway, "-netmask", "128.0.0.0")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	gw, err := sysproxy.GetDefaultGateway()
	if err != nil {
		return err
	}
	c = exec.Command("route", "delete", "-host", v.ServerIP, gw, "-netmask", "255.255.255.255")
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}
