package sysproxy

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
)

func getNetworkServices() ([]string, error) {
	c := exec.Command("networksetup", "-listallnetworkservices")
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, errors.New("ns lans:" + string(out) + ":" + err.Error())
	}
	nss := make([]string, 0)
	reg := regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`)
	for _, v := range bytes.Split(bytes.TrimSpace(out), []byte("\n")) {
		// An asterisk (*) denotes that a network service is disabled.
		if bytes.Contains(v, []byte("*")) {
			continue
		}
		ns := string(bytes.TrimSpace(v))
		c := exec.Command("networksetup", "-getinfo", ns)
		out, err := c.CombinedOutput()
		if err != nil {
			return nil, errors.New("ns gi:" + string(out) + ":" + err.Error())
		}
		if !reg.MatchString(string(out)) {
			continue
		}
		nss = append(nss, ns)
	}
	if len(nss) == 0 {
		return nil, errors.New("no available network service")
	}
	return nss, nil
}

// TurnOnSystemProxy used to enable system pac proxy, pac is a URL.
func TurnOnSystemProxy(pac string) error {
	nss, err := getNetworkServices()
	if err != nil {
		return err
	}
	for _, v := range nss {
		c := exec.Command("networksetup", "-setautoproxyurl", v, pac)
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New("ns apu:" + string(out) + ":" + err.Error())
		}
		c = exec.Command("networksetup", "-setautoproxystate", v, "on")
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New("ns aps:" + string(out) + ":" + err.Error())
		}
	}
	return nil
}

// TurnOffSystemProxy used to disable system pac proxy
func TurnOffSystemProxy() error {
	nss, err := getNetworkServices()
	if err != nil {
		return err
	}
	for _, v := range nss {
		c := exec.Command("networksetup", "-setautoproxystate", v, "off")
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New("ns aps:" + string(out) + ":" + err.Error())
		}
	}
	return nil
}
