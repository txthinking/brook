package music

import "strings"

// ChinaMobileSDC is the xxx music
type ChinaMobileSDC struct {
	Song []byte
}

// NewChinaMobileSDC returns a new ChinaMobileSDC
func NewChinaMobileSDC() *ChinaMobileSDC {
	ss := make([]string, 0)
	ss = append(ss, "POST http://sdc.10086.cn/ HTTP/1.1")
	ss = append(ss, "Host: sdc.10086.cn")
	ss = append(ss, "X-Online-Host: sdc.10086.cn")
	s := strings.Join(ss, "\r\n")
	return &ChinaMobileSDC{
		Song: []byte(s),
	}
}

// Length returns length of song
func (c *ChinaMobileSDC) Length() int {
	return len(c.Song)
}

// GetSong returns song of music
func (c *ChinaMobileSDC) GetSong() []byte {
	return c.Song
}
