package music

import "strings"

// ChinaUnicomIRead is the xxx music
type ChinaUnicomIRead struct {
	Song []byte
}

// NewChinaUnicomIRead returns a new ChinaUnicomIRead
func NewChinaUnicomIRead() *ChinaUnicomIRead {
	ss := make([]string, 0)
	ss = append(ss, "POST http://iread.wo.com.cn/ HTTP/1.1")
	ss = append(ss, "Host: iread.wo.com.cn")
	ss = append(ss, "X-Online-Host: iread.wo.com.cn")
	ss = append(ss, "Connection: keep-alive")
	ss = append(ss, "User-Agent: iread")
	ss = append(ss, "Content-Type: application/octet-stream")
	ss = append(ss, "Referer: http://iread.wo.com.cn/")
	ss = append(ss, "Accept-Encoding: gzip, deflate, br")
	ss = append(ss, "Accept-Language: zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.4")
	s := strings.Join(ss, "\r\n")
	s += "\r\n"
	return &ChinaUnicomIRead{
		Song: []byte(s),
	}
}

// Length returns length of song
func (c *ChinaUnicomIRead) Length() int {
	return len(c.Song)
}

// GetSong returns song of music
func (c *ChinaUnicomIRead) GetSong() []byte {
	return c.Song
}
