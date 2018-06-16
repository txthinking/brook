package sysproxy

import (
	"bytes"
	"errors"
	"os/exec"
	"syscall"

	"github.com/txthinking/ant"
)

// TurnOnSystemProxy used to enable system pac proxy, pac is a URL.
func TurnOnSystemProxy(pac string) error {
	c := exec.Command(`reg`, `add`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `AutoConfigURL`, `/t`, `REG_SZ`, `/d`, pac, `/f`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		u, e := ant.GBK2UTF8(out)
		if e != nil {
			return errors.New("reg add acu g2u:" + e.Error())
		}
		return errors.New("reg add acu:" + string(u) + ":" + err.Error())
	}
	c = exec.Command(`reg`, `add`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `ProxyEnable`, `/t`, `REG_DWORD`, `/d`, `0`, `/f`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		u, e := ant.GBK2UTF8(out)
		if e != nil {
			return errors.New("reg add pe g2u:" + e.Error())
		}
		return errors.New("reg add pe:" + string(u) + ":" + err.Error())
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
	if out, err := c.CombinedOutput(); err != nil {
		u, e := ant.GBK2UTF8(out)
		if e != nil {
			return errors.New("reg query acu g2u:" + e.Error())
		}
		return errors.New("reg query acu:" + string(u) + ":" + err.Error())
	} else {
		u, e := ant.GBK2UTF8(out)
		if e != nil {
			return errors.New("reg query acu g2u:" + e.Error())
		}
		if bytes.Contains(bytes.ToLower(u), []byte("autoconfigurl")) {
			c := exec.Command(`reg`, `delete`, `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, `/v`, `AutoConfigURL`, `/f`)
			c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			if out, err := c.CombinedOutput(); err != nil {
				u, e := ant.GBK2UTF8(out)
				if e != nil {
					return errors.New("reg delete acu g2u:" + e.Error())
				}
				return errors.New("reg delete acu:" + string(u) + ":" + err.Error())
			}
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
