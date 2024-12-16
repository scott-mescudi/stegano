package pkg

import (
	"fmt"
	"testing"
)

func TestBytesToBinary(t *testing.T) {
	tests := []struct {
		input    []byte
		expected []uint8
	}{
		{[]byte{0x01}, []uint8{0, 0, 0, 0, 0, 0, 0, 1}},
		{[]byte{0xFF}, []uint8{1, 1, 1, 1, 1, 1, 1, 1}},
		{[]byte{0x00}, []uint8{0, 0, 0, 0, 0, 0, 0, 0}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input: %v", tt.input), func(t *testing.T) {
			result := BytesToBinary(tt.input)
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("For input %v at index %d, expected %d but got %d", tt.input, i, tt.expected[i], v)
				}
			}
		})
	}
}

func TestFlipLSB(t *testing.T) {
	tests := []struct {
		input    uint32
		expected uint32
	}{
		{0b00000000000000000000000000000000, 0b00000000000000000000000000000001}, // Flip the LSB from 0 to 1
		{0b00000000000000000000000000000001, 0b00000000000000000000000000000000}, // Flip the LSB from 1 to 0
		{0b00000000000000000000000000000010, 0b00000000000000000000000000000011}, // Flip the LSB from 0 to 1
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input: %032b", tt.input), func(t *testing.T) {
			result := FlipBit(tt.input, 0)
			if result != tt.expected {
				t.Errorf("For input %032b, expected %032b but got %032b", tt.input, tt.expected, result)
			}
		})
	}
}

func TestInt32ToBinary(t *testing.T) {
	tests := []struct {
		input    int32
		expected []uint8
	}{
		{69, []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1}},
		{1, []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{-1, []uint8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Input: %d", tt.input), func(t *testing.T) {
			result := Int32ToBinary(tt.input)
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("For input %d at index %d, expected %d but got %d", tt.input, i, tt.expected[i], v)
				}
			}
		})
	}
}


func BenchmarkBytesToBinary(b *testing.B) {
	// Prepare a large input slice of bytes
	largeData := make([]byte, 1000000) // 1 million bytes
	for i := 0; i < len(largeData); i++ {
		largeData[i] = byte(i % 256) // Fill with cyclic values 0-255
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BytesToBinary(largeData)
	}
}
// before
// 7772794 ns/op        41704331 B/op         43 allocs/op
// after 
// 4121559 ns/op         8003606 B/op          1 allocs/op

func BenchmarkInt32ToBinary(b *testing.B) {
	// Test with a random 32-bit integer
	num := int32(123456789)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Int32ToBinary(num)
	}
}
