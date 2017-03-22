package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"runtime"
)

// Setting is the config of your server/client
type Setting struct {
	Type     string `json:"Type"`
	Server   string `json:"Server"`
	Password string `json:"Password"`
	Timeout  int    `json:"Timeout"`
	Deadline int    `json:"Deadline"`
	Local    string `json:"Local"`
	Music    string `json:"Music"`
}

// PostSetting used to add/update setting
func PostSetting(w http.ResponseWriter, r *http.Request) {
	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	st := &Setting{}
	if err = json.Unmarshal(d, st); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err = StoreSetting(st); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(200)
}

// GetSetting returns the setting
func GetSetting(w http.ResponseWriter, r *http.Request) {
	st, err := RetrieveSetting()
	if err != nil {
		http.Error(w, err.Error(), 422)
		return
	}
	d, err := json.Marshal(st)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(d)
}

// StoreSetting save your setting
func StoreSetting(setting *Setting) error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	var fp string
	if runtime.GOOS == "windows" {
		fp = u.HomeDir + "\\.brook.setting"
	} else {
		fp = u.HomeDir + "/.brook.setting"
	}
	f, err := os.OpenFile(fp, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	d, err := json.Marshal(setting)
	if err != nil {
		return err
	}
	if _, err = f.Write(d); err != nil {
		return err
	}
	return nil
}

// RetrieveSetting returns your settting
func RetrieveSetting() (*Setting, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	var fp string
	if runtime.GOOS == "windows" {
		fp = u.HomeDir + "\\.brook.setting"
	} else {
		fp = u.HomeDir + "/.brook.setting"
	}
	f, err := os.OpenFile(fp, os.O_RDONLY, 0600)
	if os.IsNotExist(err) {
		return nil, errors.New("No setting")
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	setting := &Setting{}
	if err = json.Unmarshal(d, setting); err != nil {
		return nil, err
	}
	return setting, nil
}
