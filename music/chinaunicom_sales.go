package music

import "strings"

// ChinaUnicomSales is the xxx music
type ChinaUnicomSales struct {
	Song []byte
}

// NewChinaUnicomSales returns a new ChinaUnicomSales
func NewChinaUnicomSales() *ChinaUnicomSales {
	ss := make([]string, 0)
	ss = append(ss, "GET http://sales.wostore.cn:8081/activity/vpn_tpl/in/index.html?tplId=21 HTTP/1.1")
	ss = append(ss, "Host: sales.wostore.cn:8081")
	ss = append(ss, "Proxy-Connection: keep-alive")
	s := strings.Join(ss, "\r\n")
	return &ChinaUnicomSales{
		Song: []byte(s),
	}
}

// Length returns length of song
func (c *ChinaUnicomSales) Length() int {
	return len(c.Song)
}

// GetSong returns song of music
func (c *ChinaUnicomSales) GetSong() []byte {
	return c.Song
}
