package sysproxy

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

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
	if err := restartWindowsIE(); err != nil {
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
	if err := restartWindowsIE(); err != nil {
		return err
	}
	return nil
}

func restartWindowsIE() error {
	c := exec.Command(`tasklist`, `/fo`, `list`, `/fi`, `imagename eq iexplore.exe`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		u, e := ant.GBK2UTF8(out)
		if e != nil {
			return errors.New("tasklist ie g2u:" + e.Error())
		}
		return errors.New("tasklist ie:" + string(u) + ":" + err.Error())
	} else {
		u, e := ant.GBK2UTF8(out)
		if e != nil {
			return errors.New("tasklist ie g2u:" + e.Error())
		}
		if bytes.Contains(bytes.ToLower(u), []byte("iexplore.exe")) {
			c := exec.Command(`taskkill`, `/f`, `/t`, `/im`, `iexplore.exe`)
			c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			if out, err := c.CombinedOutput(); err != nil {
				u, e := ant.GBK2UTF8(out)
				if e != nil {
					return errors.New("taskkill ie before start g2u:" + e.Error())
				}
				return errors.New("taskkill ie before start:" + string(u) + ":" + err.Error())
			}
		}
	}
	c = exec.Command(fmt.Sprintf(`%s\Internet Explorer\iexplore.exe`, os.Getenv("PROGRAMFILES")), `-nohome`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := c.Start(); err != nil {
		return errors.New("ie start:" + err.Error())
	}
	time.Sleep(2 * time.Second)
	c = exec.Command(`taskkill`, `/f`, `/t`, `/im`, `iexplore.exe`)
	c.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if out, err := c.CombinedOutput(); err != nil {
		u, e := ant.GBK2UTF8(out)
		if e != nil {
			return errors.New("taskkill ie after start g2u:" + e.Error())
		}
		return errors.New("taskkill ie after start:" + string(u) + ":" + err.Error())
	}
	return nil
}
