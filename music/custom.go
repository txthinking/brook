package music

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Custom is music for customize URL
type Custom struct {
	Song []byte
}

// NewCustom returns a Custom by URL
func NewCustom(u *url.URL) (*Custom, error) {
	if u.Scheme == "http" || u.Scheme == "https" {
		c := http.Client{
			Timeout: time.Second * 10,
		}
		res, err := c.Get(u.String())
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return &Custom{
			Song: data,
		}, nil
	}
	if u.Scheme == "file" {
		f, err := os.Open(u.Path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		return &Custom{
			Song: data,
		}, nil
	}
	return nil, errors.New("Unsupport URL")
}

// Length returns length of song
func (c *Custom) Length() int {
	return len(c.Song)
}

// GetSong returns song of music
func (c *Custom) GetSong() []byte {
	return c.Song
}
