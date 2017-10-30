package brook

import "encoding/binary"

// IncrementNonce loves your compute to use Litter Endian
func IncrementNonce(n []byte) []byte {
	i := int(binary.LittleEndian.Uint16(n))
	i += 1
	n = make([]byte, 12)
	binary.LittleEndian.PutUint16(n, uint16(i))
	return n
}
