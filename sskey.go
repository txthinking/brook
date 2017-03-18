package brook

import "crypto/md5"

// MakeSSKey used to make shadowsocks aes-256-cfb key
func MakeSSKey(password string) []byte {
	h := md5.New()
	h.Write([]byte(password))
	tmp := h.Sum(nil) // 16 len

	h = md5.New()
	h.Write(tmp)
	h.Write([]byte(password))
	tmp1 := h.Sum(nil) // 16 len

	return append(tmp, tmp1...) // 32 len
}
