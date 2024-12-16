package pkg

import (
	"fmt"
)

// stringToBinary converts a string to a slice of bits (0s and 1s).
func BytesToBinary(data []byte) []uint8 {
	var bits = make([]uint8, len(data)*8)
	idx := 0
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1
			bits[idx] = uint8(bit)
			idx++
		}
	}
	return bits[:idx]
}

func Int32ToBinary(num int32) []uint8 {
	var bits []uint8

	for i := 31; i >= 0; i-- {
		bit := (num >> i) & 1
		bits = append(bits, uint8(bit))
	}
	return bits
}

func GetlenOfData(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, fmt.Errorf("insufficient data: expected at least 4 bytes")
	}

	n := int(data[0])<<24 | int(data[1])<<16 | int(data[2])<<8 | int(data[3])
	return n, nil
}

func FlipBit(num uint32, position uint8) uint32 {
	return num ^ (1 << position)
}

func GetBit(value uint32, position uint8) uint8 {
	return uint8((value >> position) & 1)
}
