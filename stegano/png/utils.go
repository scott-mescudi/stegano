package png

import (
	"fmt"
	"strconv"
)

// stringToBinary converts a string to a slice of bits (0s and 1s).
func bytesToBinary(data []byte) []int {
	var bits []int
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1
			bits = append(bits, int(bit))
		}
	}
	return bits
}

func flipLSB(num uint32) uint32 {
	return num ^ 1 // Flip the LSB using XOR
}

func int32ToBinary(num int32) []int {
	var bits []int
	for i := 31; i >= 0; i-- {
		bit := (num >> i) & 1
		bits = append(bits, int(bit))
	}
	return bits
}

func GetlenOfData(test []byte) (int, error) {
	container := ""
	for i := 0; i < 4; i++ {
		b := test[i]
		binary := fmt.Sprintf("%08b", b)
		container += binary
	}

	n, err := strconv.ParseInt(container, 2, 32)
	if err != nil {
		return 0, fmt.Errorf("error parsing binary to int: %e", err)
	}

	return int(n), nil
}
