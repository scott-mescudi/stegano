package utils

import (
	"fmt"
	"testing"
)

func TestBytesToBinary(t *testing.T) {
	tests := []struct {
		input    []byte
		expected []int
	}{
		{[]byte{0x01}, []int{0, 0, 0, 0, 0, 0, 0, 1}},
		{[]byte{0xFF}, []int{1, 1, 1, 1, 1, 1, 1, 1}},
		{[]byte{0x00}, []int{0, 0, 0, 0, 0, 0, 0, 0}},
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
			result := FlipLSB(tt.input)
			if result != tt.expected {
				t.Errorf("For input %032b, expected %032b but got %032b", tt.input, tt.expected, result)
			}
		})
	}
}

func TestInt32ToBinary(t *testing.T) {
	tests := []struct {
		input    int32
		expected []int
	}{
		{69, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1}},
		{1, []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{-1, []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
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
