package sysproxy

import (
	"errors"
	"os"
	"os/exec"
	"regexp"
)

// GetNetworkInterfaces returns interface list
func GetNetworkInterfaces() ([]string, error) {
	return []string{}, nil
}

func TurnOnSystemProxy(pac string) error {
	return nil
}

func TurnOffSystemProxy() error {
	return nil
}

// SetDNSServer used to set system DNS server
func SetDNSServer(server string) error {
	f, err := os.OpenFile("/etc/resolv.conf", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("nameserver " + server)
	if err != nil {
		return err
	}
	return nil
}

// GetDefaultGateway returns default gateway
func GetDefaultGateway() (string, error) {
	c := exec.Command("ip", "route")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	r, err := regexp.Compile(`default.*?(\d+.\d+\.\d+\.\d+)`)
	if err != nil {
		return "", err
	}
	ss := r.FindStringSubmatch(string(out))
	if len(ss) == 0 {
		return "", errors.New("Can not find default gateway")
	}
	return ss[1], nil
}
