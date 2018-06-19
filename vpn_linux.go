package brook

import (
	"errors"
	"os/exec"

	"github.com/txthinking/brook/sysproxy"
)

// AddRoutes adds routes
func (v *VPN) AddRoutes() error {
	c := exec.Command("ip", "route", "add", "0.0.0.0/1", "via", v.TunGateway)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("ip", "route", "add", "128.0.0.0/1", "via", v.TunGateway)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	gw, err := sysproxy.GetDefaultGateway()
	if err != nil {
		return err
	}
	c = exec.Command("ip", "route", "add", v.ServerIP, "via", gw)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}

// DeleteRoutes deletes routes
func (v *VPN) DeleteRoutes() error {
	c := exec.Command("ip", "route", "del", "0.0.0.0/1", "via", v.TunGateway)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	c = exec.Command("ip", "route", "del", "128.0.0.0/1", "via", v.TunGateway)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	gw, err := sysproxy.GetDefaultGateway()
	if err != nil {
		return err
	}
	c = exec.Command("ip", "route", "del", v.ServerIP, "via", gw)
	if out, err := c.CombinedOutput(); err != nil {
		return errors.New(string(out) + err.Error())
	}
	return nil
}
