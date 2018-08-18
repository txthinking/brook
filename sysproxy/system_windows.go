package sysproxy

import (
	"bytes"
	"errors"
	"os/exec"
	"regexp"
	"syscall"
)

// TurnOnSystemProxy used to enable system pac proxy, pac is a URL.
func TurnOnSystemProxy(pac string) error {
	c := exec.Command(`reg`, `add`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `AutoConfigURL`, `/t`, `REG_SZ`, `/d`, pac, `/f`)
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

// TurnOffSystemProxy used to disable system pac proxy
func TurnOffSystemProxy() error {
	c := exec.Command(`reg`, `query`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return errors.New("Can not query proxy settings: " + err.Error())
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

// GetNetworkInterfaces returns interface list
func GetNetworkInterfaces() ([]string, error) {
	c := exec.Command("netsh", "interface", "ipv4", "show", "address")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out) + err.Error())
	}
	r, err := regexp.Compile(`"(.+)".*\n.*Yes`)
	if err != nil {
		return nil, err
	}
	ss := r.FindAllStringSubmatch(string(out), -1)
	is := make([]string, 0)
	for _, v := range ss {
		is = append(is, v[1])
	}
	return is, nil
}

// SetDNSServer used to set system DNS server
func SetDNSServer(server string) error {
	nis, err := GetNetworkInterfaces()
	if err != nil {
		return err
	}
	for _, v := range nis {
		c := exec.Command("netsh", "interface", "ip", "set", "dnsservers", v, "static", server, "primary")
		c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		if out, err := c.CombinedOutput(); err != nil {
			return errors.New(string(out) + err.Error())
		}
	}
	return nil
}

// GetDefaultGateway returns default gateway
func GetDefaultGateway() (string, error) {
	c := exec.Command("netsh", "interface", "ipv4", "show", "address")
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := c.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out) + err.Error())
	}
	r, err := regexp.Compile(`Default Gateway.*?(\d+.\d+\.\d+\.\d+)`)
	if err != nil {
		return "", err
	}
	ss := r.FindStringSubmatch(string(out))
	if len(ss) == 0 {
		return "", errors.New("Can not find default gateway")
	}
	return ss[1], nil
}
