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

// GetResponse returns response when the request does not equal with the song
func (c *ChinaUnicomSales) GetResponse(request []byte) []byte {
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
