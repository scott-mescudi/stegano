package pkg

import (
	"reflect"
	"testing"
)

func TestIntToArr(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected []byte
	}{
		{
			name:     "Normal number",
			input:    2345,
			expected: []byte{0, 0, 9, 41},
		},
		{
			name:     "Zero",
			input:    0,
			expected: []byte{0, 0, 0, 0},
		},
		{
			name:     "Maximum 4-byte int",
			input:    2147483647,
			expected: []byte{127, 255, 255, 255},
		},
		{
			name:     "Small positive number",
			input:    42,
			expected: []byte{0, 0, 0, 42},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := intToArr(test.input)
			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("For input %v, expected %v but got %v", test.input, test.expected, res)
			}
		})
	}
}

func TestArrToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected int
	}{
		{
			name:     "Normal number",
			input:    []byte{0, 0, 9, 41},
			expected: 2345,
		},
		{
			name:     "Zero",
			input:    []byte{0, 0, 0, 0},
			expected: 0,
		},
		{
			name:     "Maximum 4-byte int",
			input:    []byte{127, 255, 255, 255},
			expected: 2147483647,
		},
		{
			name:     "Small positive number",
			input:    []byte{0, 0, 0, 42},
			expected: 42,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := arrToInt(test.input)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("For input %v, expected %v but got %v", test.input, test.expected, res)
			}
		})
	}

}

func TestMatrixToSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]byte
		expected []byte
	}{
		{
			name:     "Normal",
			input:    [][]byte{{0, 0, 9, 41}, {0, 0, 9, 41}, {0, 0, 9, 41}, {0, 0, 9, 41}, {0, 0, 9, 41}},
			expected: []byte{0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 9, 41, 0, 0, 9, 41, 0, 0, 9, 41, 0, 0, 9, 41, 0, 0, 9, 41},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := matrixToSlice(test.input)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("For input %v, expected %v but got %v", test.input, test.expected, res)
			}
		})
	}

}

func TestSliceToMatrix(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected [][]byte
	}{
		{
			name:     "Normal",
			input:    []byte{0, 0, 0, 4, 0, 0, 0, 5, 0, 0, 9, 41, 0, 0, 9, 41, 0, 0, 9, 41, 0, 0, 9, 41, 0, 0, 9, 41},
			expected: [][]byte{{0, 0, 9, 41}, {0, 0, 9, 41}, {0, 0, 9, 41}, {0, 0, 9, 41}, {0, 0, 9, 41}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := sliceToMatrix(test.input)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("For input %v, expected %v but got %v", test.input, test.expected, res)
			}
		})
	}

}

func TestRsEncode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Normal",
			input:    []byte("hello world"),
			expected: []byte{0, 0, 0, 11, 0, 0, 0, 5, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := RsEncode(test.input, 4)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("For input %v, expected %v but got %v", test.input, test.expected, res)
			}
		})
	}
}

func TestRsDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "Normal",
			input:    []byte{0, 0, 0, 11, 0, 0, 0, 5, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100},
			expected: []byte("hello world"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := RsDecode(test.input, 1, 4)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(res, test.expected) {
				t.Errorf("For input %v, expected %v but got %v", test.input, test.expected, res)
			}
		})
	}
}
