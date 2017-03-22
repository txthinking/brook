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
	systray.Run(run)
}

func run() {
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
	systray.AddMenuItem("---------", "").Disable()
	mStart := systray.AddMenuItem("Start", "")
	mStop := systray.AddMenuItem("Stop", "")
	mSetting := systray.AddMenuItem("Setting", "")
	systray.AddMenuItem("---------", "").Disable()
	mGithub := systray.AddMenuItem("Github", "")
	mEmail := systray.AddMenuItem("Author: cloud@txthinking.com", "")
	systray.AddMenuItem("Version: 20170322", "")
	systray.AddMenuItem("---------", "").Disable()
	mQuit := systray.AddMenuItem("Quit", "")

	showNotice("Status", "stoped")
	mStop.Disable()
	var bk *brook.BKClient
	var ss *brook.SSClient
	var s5 *brook.S5Client
	var quitTimes int
	quit := make(chan struct{})

	start := func() {
		defer func() {
			mStart.Enable()
			mStop.Disable()
		}()
		mStart.Disable()
		mStop.Enable()
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
			bk, err = brook.NewBKClient(st.Local, st.Server, st.Password, st.Timeout, st.Deadline, st.Music, nil)
			if err != nil {
				showNotice("Error", err.Error())
				return
			}
			if err := bk.ListenAndServe(); err != nil {
				showNotice("Status", "stoped")
				return
			}
		}
		if st.Type == "ss" {
			ss = brook.NewSSClient(st.Local, st.Server, st.Password, st.Timeout, st.Deadline, nil)
			if err := ss.ListenAndServe(); err != nil {
				showNotice("Status", "stoped")
				return
			}
		}
		if st.Type == "s5" {
			s5 = brook.NewS5Client(st.Local, st.Server, st.Password, st.Timeout, st.Deadline, nil)
			if err := ss.ListenAndServe(); err != nil {
				showNotice("Status", "stoped")
				return
			}
		}
		showNotice("Error", "error type")
	}

	stop := func() error {
		if err := sysproxy.TurnOffSystemProxy(); err != nil {
			return err
		}
		if bk != nil {
			bk.Shutdown()
			bk = nil
		}
		if ss != nil {
			ss.Shutdown()
			ss = nil
		}
		if s5 != nil {
			s5.Shutdown()
			s5 = nil
		}
		return nil
	}

	go func() {
		if err := RunHTTPServer("127.0.0.1:1980"); err != nil {
			showNotice("Error", err.Error())
		}
	}()
	go func() {
		for {
			select {
			case <-mStart.ClickedCh:
				go start()
			case <-mStop.ClickedCh:
				if err := stop(); err != nil {
					showNotice("Error", err.Error())
				}
			case <-mSetting.ClickedCh:
				if err := open.Run("http://127.0.0.1:1980"); err != nil {
					showNotice("Error", err.Error())
				}
			case <-mGithub.ClickedCh:
				if err := open.Run("https://github.com/txthinking/brook"); err != nil {
					showNotice("Error", err.Error())
				}
			case <-mEmail.ClickedCh:
				if err := open.Run("mailto:cloud@txthinking.com"); err != nil {
					showNotice("Error", err.Error())
				}
			case <-mQuit.ClickedCh:
				quitTimes++
				err := stop()
				if err == nil {
					quit <- struct{}{}
				}
				if err != nil {
					showNotice("Error", err.Error())
					if quitTimes > 1 {
						quit <- struct{}{}
					}
				}
			}
		}
	}()
	<-quit
	showNotice("Status", "Quiting")
	systray.Quit()
	os.Exit(0)
}
