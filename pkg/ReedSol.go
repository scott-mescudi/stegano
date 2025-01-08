package pkg

import (
	"fmt"
	"github.com/klauspost/reedsolomon"
)

func intToArr(num int) ([]byte) {
	return []byte{
		byte(num >> 24),
		byte(num >> 16),
		byte(num >> 8),
		byte(num),
	}
}

func arrToInt(data []byte) (int, error) {
	if len(data) != 4 {
		return 0, fmt.Errorf("expected 4 bytes, got %d", len(data))
	}

	return int(data[0])<<24 | int(data[1])<<16 | int(data[2])<<8 | int(data[3]), nil
}

func matrixToSlice(matrix [][]byte) ([]byte, error) {
	if len(matrix) == 0 {
		return nil, fmt.Errorf("matrix cannot be empty")
	}

	rows := len(matrix)
	cols := len(matrix[0])
	for _, row := range matrix {
		if len(row) != cols {
			return nil, fmt.Errorf("non-uniform row lengths in matrix")
		}
	}

	totalSize := rows * cols
	data := make([]byte, 8+totalSize)

	// Add dimensions
	copy(data[0:4], intToArr(cols))
	copy(data[4:8], intToArr(rows))

	// Copy matrix data
	idx := 8
	for _, row := range matrix {
		copy(data[idx:idx+cols], row)
		idx += cols
	}

	return data, nil
}

func sliceToMatrix(fullData []byte) ([][]byte, error) {
	if len(fullData) < 8 {
		return nil, fmt.Errorf("data too small to extract dimensions")
	}

	col, err := arrToInt(fullData[0:4])
	if err != nil {
		return nil, fmt.Errorf("invalid column data: %w", err)
	}

	row, err := arrToInt(fullData[4:8])
	if err != nil {
		return nil, fmt.Errorf("invalid row data: %w", err)
	}

	expectedSize := 8 + row*col
	if len(fullData) < expectedSize {
		return nil, fmt.Errorf("data size mismatch: expected at least %d bytes, got %d", expectedSize, len(fullData))
	}

	matrix := make([][]byte, row)
	offset := 8
	for i := 0; i < row; i++ {
		matrix[i] = fullData[offset : offset+col]
		offset += col
	}

	return matrix, nil
}

// RsEncode encodes data into shards with parity and transfroms reed sol matrix to slice for embbeding
func RsEncode(data []byte, parity int) ([]byte, error) {
	if parity <= 0 {
		return nil, fmt.Errorf("parity must be greater than zero")
	}

	enc, err := reedsolomon.New(1, parity)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoder: %w", err)
	}

	shards := make([][]byte, 1+parity)
	for i := 0; i < 1+parity; i++ {
		if i == 0 {
			shards[i] = data // Assign data to the first shard
		} else {
			shards[i] = make([]byte, len(data))
		}
	}

	if err := enc.Encode(shards); err != nil {
		return nil, fmt.Errorf("failed to encode shards: %w", err)
	}

	packed, err := matrixToSlice(shards)
	if err != nil {
		return nil, fmt.Errorf("failed to convert matrix to slice: %w", err)
	}

	return packed, nil
}

// RsDecode decodes and reconstructs missing shards from packed matrix and returns a slice first elm in matrix
func RsDecode(packedShards []byte, dataShards, parityShards int) ([]byte, error) {
	enc, err := reedsolomon.New(dataShards, parityShards)
	if err != nil {
		return nil, fmt.Errorf("error creating decoder: %w", err)
	}

	shards, err := sliceToMatrix(packedShards)
	if err != nil {
		return nil, fmt.Errorf("failed to convert slice to matrix: %w", err)
	}

	if err := enc.Reconstruct(shards); err != nil {
		return nil, fmt.Errorf("error reconstructing shards: %w", err)
	}

	return shards[0], nil
}
