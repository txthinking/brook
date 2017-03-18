package main

import (
	"os"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"github.com/txthinking/brook"
	"github.com/txthinking/brook/sysproxy"
)

const pac = "https://pac.nixisall.com/white/SOCKS5%20127.0.0.1:1080;%20SOCKS%20127.0.0.1:1080"

func main() {
	systray.Run(func() {
		systray.SetIcon(iconData)
		notice := systray.AddMenuItem("", "")
		showNotice := func(tp, msg string) {
			if tp == "Error" {
				systray.SetTitle("Error")
				systray.SetTooltip(msg)
				notice.SetTitle("[Brook] " + msg)
				return
			}
			systray.SetTitle("")
			if msg != "running" {
				systray.SetTitle(msg)
			}
			systray.SetTooltip(msg)
			notice.SetTitle("[Brook] " + msg)
		}
		showNotice("status", "stoped")
		mStart := systray.AddMenuItem("Start", "")
		mSetting := systray.AddMenuItem("Setting", "")
		mGithub := systray.AddMenuItem("Github", "")
		mEmail := systray.AddMenuItem("Author: cloud@txthinking.com", "")
		mEmail := systray.AddMenuItem("Version: 20170316", "")
		mQuit := systray.AddMenuItem("Quit", "")

		start := func() {
			defer mStart.Enable()
			mStart.Disable()
			showNotice("Status", "running")
			st, err := RetrieveSetting()
			if err != nil {
				showNotice("Error", err.Error())
				return
			}
			if err := sysproxy.TurnOnSystemProxy(pac); err != nil {
				showNotice("Error", err.Error())
				return
			}
			if st.Type == "bk" {
				if err := brook.RunBKClient(st.Local, st.Server, st.Password, st.Timeout, st.Deadline, st.Music); err != nil {
					showNotice("Error", err.Error())
					return
				}
			} else if st.Type == "ss" {
				if err := brook.RunSSClient(st.Local, st.Server, st.Password, st.Timeout, st.Deadline); err != nil {
					showNotice("Error", err.Error())
					return
				}
			} else if st.Type == "s5" {
				if err := brook.RunS5Client(st.Local, st.Server, st.Password, st.Timeout, st.Deadline); err != nil {
					showNotice("Error", err.Error())
					return
				}
			} else {
				showNotice("Error", "error type")
			}
		}
		go func() {
			if err := RunHTTPServer("127.0.0.1:1980"); err != nil {
				showNotice("Error", err.Error())
			}
		}()
		go func() {
			for {
				select {
				case <-mSetting.ClickedCh:
					if err := open.Run("http://127.0.0.1:1980"); err != nil {
						showNotice("Error", err.Error())
					}
				case <-mGithub.ClickedCh:
					if err := open.Run("https://github.com/txthinking/brook"); err != nil {
						showNotice("Error", err.Error())
					}
				case <-mStart.ClickedCh:
					go start()
				case <-mEmail.ClickedCh:
					if err := open.Run("mailto:cloud@txthinking.com"); err != nil {
						showNotice("Error", err.Error())
					}
				}
			}
		}()
		<-mQuit.ClickedCh
		showNotice("Status", "Quiting")
		if err := sysproxy.TurnOffSystemProxy(); err != nil {
			showNotice("Error", err.Error())
			return
		}
		systray.Quit()
		os.Exit(0)
	})
}
