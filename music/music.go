package music

import "net/url"

// Music is a supplement to the brook protocol,
// you can create your own music
type Music interface {
	// The length of the song
	Length() int
	// The data send first
	GetSong() []byte
	// Length of request is equal with song, but it may be any data and may be not complete data as some protocol, like a part of HTTP request
	GetResponse(request []byte) []byte
}

// GetMusic return a music by the given key/url, or nil if not found
func GetMusic(k string) (Music, error) {
	switch k {
	case "chinamobile_sdc":
		return NewChinaMobileSDC(), nil
	case "chinaunicom_iread":
		return NewChinaUnicomIRead(), nil
	case "chinaunicom_sales":
		return NewChinaUnicomSales(), nil
	}

	if k == "" {
		return nil, nil
	}

	u, err := url.Parse(k)
	if err != nil {
		return nil, err
	}
	return NewCustom(u)
}
