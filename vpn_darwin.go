package brook

import (
	"errors"
	"os/exec"

	"github.com/txthinking/brook/sysproxy"
)

// AddRoutes adds routes
func (v *VPN) AddRoutes() error {
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

// DeleteRoutes deletes routes
func (v *VPN) DeleteRoutes() error {
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
