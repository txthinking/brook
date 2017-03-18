package music

import "net/url"

// Music is a supplement to the brook protocol,
// you can create your own music
type Music interface {
	Length() int
	GetSong() []byte
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
