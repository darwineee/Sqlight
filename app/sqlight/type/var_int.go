package _type

import "fmt"

type VarInt uint64

func ReadVarInt(data []byte) (VarInt, int, error) {
	var value uint64
	var bytesRead int

	for _, b := range data {
		value = (value << 7) | uint64(b&0x7f) // 0x7f = 0111 1111
		bytesRead++
		if b < 0x80 { // 0x80 = 1000 0000
			return VarInt(value), bytesRead, nil
		}
		if bytesRead >= 9 {
			return 0, 0, fmt.Errorf("varint too large")
		}
	}
	return 0, 0, fmt.Errorf("incomplete varint")
}
