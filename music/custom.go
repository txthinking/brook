package music

import (
	"crypto/tls"
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
			Transport: &http.Transport{ // https://github.com/golang/go/issues/14514, it may not be resolved completely
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
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

// GetResponse returns response when the request does not equal with the song
func (c *Custom) GetResponse(request []byte) []byte {
	ss := ""
	ss += "HTTP/1.1 404 Not Found\r\n"
	ss += "Cache-Control: no-cache\r\n"
	ss += "Connection: close\r\n"
	ss += "Content-Type: text/html\r\n"
	ss += "Content-Length: 48\r\n"
	ss += "\r\n"
	ss += "<html><body><h1>404 Not Found</h1></body></html>"
	return []byte(ss)
}
