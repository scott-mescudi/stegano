package pkg

import (
	"reflect"
	"strings"
	"testing"
)

func TestEmbedAtDepth(t *testing.T) {
	rgbdata := []RgbChannel{
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
	}

	data := []byte("hey")

	for i := 0; i < 7; i++ {
		data, err := EmbedAtDepth(rgbdata, data, uint8(i))
		if err != nil {
			t.Fatal(err)
		}


		if data[0].R == 255 {
			t.Fatalf("Data didnt change at index %v", i)
		}
	}

	ec, err := EmbedAtDepth(rgbdata, data, uint8(10))
	if err == nil || ec != nil {
		t.Fatalf("Didnt retrun error when passing in out of bounds bitDepth")
	}
}

func TestExtractAtDepth(t *testing.T) {
	rgbdata := []RgbChannel{
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
	}

	data := []byte("hey")
	var depth uint8 = 7

	EmbedAtDepth(rgbdata, data, depth)
	edata, err := ExtractAtDepth(rgbdata, depth)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(edata), string(data)) {
		t.Fatal("ExtractedData is incomplete")
	}
}

func TestSplitIntoGroupsOfThree(t *testing.T) {
	tests := []struct {
		name     string
		input    []uint8
		expected []bin
	}{
		{
			name:  "Exact multiples of 3",
			input: []uint8{1, 2, 3, 4, 5, 6},
			expected: []bin{
				{r: 1, g: 2, b: 3},
				{r: 4, g: 5, b: 6},
			},
		},
		{
			name:  "Non-multiples of 3 (1 extra)",
			input: []uint8{1, 2, 3, 4},
			expected: []bin{
				{r: 1, g: 2, b: 3},
				{r: 4}, // Remaining element
			},
		},
		{
			name:  "Non-multiples of 3 (2 extra)",
			input: []uint8{1, 2, 3, 4, 5},
			expected: []bin{
				{r: 1, g: 2, b: 3},
				{r: 4, g: 5}, // Remaining elements
			},
		},
		{
			name:  "Single element",
			input: []uint8{1},
			expected: []bin{
				{r: 1},
			},
		},
		{
			name:  "Two elements",
			input: []uint8{1, 2},
			expected: []bin{
				{r: 1, g: 2},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := splitIntoGroupsOfThree(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("For input %v, expected %v but got %v", test.input, test.expected, result)
			}
		})
	}
}

func BenchmarkSplitIntoGroupsOfThree(b *testing.B) {
	// Generate a large input slice for benchmarking
	largeInput := make([]uint8, 1000000) // 1 million elements
	for i := 0; i < len(largeInput); i++ {
		largeInput[i] = uint8(i % 256) // Fill with cyclic values from 0 to 255
	}

	b.ResetTimer() // Reset the timer to exclude setup time
	for i := 0; i < b.N; i++ {
		_ = splitIntoGroupsOfThree(largeInput)
	}
}

// before
// 789349 ns/op         5712932 B/op         31 allocs/op
// after 
// 322274 ns/op         1007623 B/op          1 allocs/op
