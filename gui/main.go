package main

import (
	"os"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"github.com/txthinking/brook"
	"github.com/txthinking/brook/sysproxy"
)

const pac = "http://local.txthinking.com:1980/pac"

func main() {
	systray.Run(run, func() {})
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
	mGithub := systray.AddMenuItem("Upgrade", "")
	mEmail := systray.AddMenuItem("Contact: cloud@txthinking.com", "")
	systray.AddMenuItem("Version: 20170909", "")
	systray.AddMenuItem("---------", "").Disable()
	mQuit := systray.AddMenuItem("Quit", "")

	showNotice("Status", "stopped")
	mStop.Disable()
	var bk *brook.Client
	var ss *brook.SSClient
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
			bk, err = brook.NewClient(st.Address, st.Server, st.Password, st.TCPTimeout, st.TCPDeadline, st.UDPDeadline, st.UDPSessionTime)
			if err != nil {
				showNotice("Error", err.Error())
				return
			}
			if err := bk.ListenAndServe(nil); err != nil {
				showNotice("Status", "stopped")
				return
			}
		}
		if st.Type == "ss" {
			ss, err = brook.NewSSClient(st.Address, st.Server, st.Password, st.TCPTimeout, st.TCPDeadline, st.UDPDeadline, st.UDPSessionTime)
			if err := ss.ListenAndServe(nil); err != nil {
				showNotice("Status", "stopped")
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
		return nil
	}

	go func() {
		if err := RunHTTPServer("local.txthinking.com:1980"); err != nil {
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
				if err := open.Run("http://local.txthinking.com:1980"); err != nil {
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
