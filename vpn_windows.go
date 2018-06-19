package brook

import (
	"errors"
	"os/exec"
	"syscall"

	"github.com/txthinking/brook/sysproxy"
)

// AddRoutes adds routes
func (v *VPN) AddRoutes() error {
	c := exec.Command("route", "add", "0.0.0.0", "mask", "128.0.0.0", v.TunGateway)
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

// DeleteRoutes deletes routes
func (v *VPN) DeleteRoutes() error {
	c := exec.Command("route", "delete", "0.0.0.0", "mask", "128.0.0.0", v.TunGateway)
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
