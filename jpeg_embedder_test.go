package stegano

import (
	"image"
	"image/color"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock function to create a test image
func createTestImage() image.Image {
	width, height := 100, 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with some color
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	return img
}

func TestEncodeAndSave_Valid(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte("some secret data")
	bitDepth := uint8(3)
	outputFilename := "test_output.jpg"

	// Execute
	handler := &JpegHandler{}
	err := handler.EncodeAndSave(coverImage, data, bitDepth, outputFilename, false)

	// Test if no error occurred and file was created
	require.NoError(t, err)

	// Check if output file exists
	_, err = os.Stat(outputFilename)
	require.NoError(t, err)

	// Cleanup
	defer os.Remove(outputFilename)
}

func TestEncodeAndSave_EmptyData(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte{}
	bitDepth := uint8(3)
	outputFilename := "test_output.jpg"

	// Execute
	handler := &JpegHandler{}
	err := handler.EncodeAndSave(coverImage, data, bitDepth, outputFilename, false)

	// Test if the correct error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data is empty")
}

func TestEncodeAndSave_DataTooLarge(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := make([]byte, 10000) // This is more data than the image can hold
	bitDepth := uint8(3)
	outputFilename := "test_output.jpg"

	// Execute
	handler := &JpegHandler{}
	err := handler.EncodeAndSave(coverImage, data, bitDepth, outputFilename, false)

	// Test if the correct error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data is too large to embed")
}

func TestEncodeAndSave_CompressedData(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte("compressed test data")
	bitDepth := uint8(3)
	outputFilename := "test_compressed_output.jpg"

	// Execute with compression enabled
	handler := &JpegHandler{}
	err := handler.EncodeAndSave(coverImage, data, bitDepth, outputFilename, true)

	// Test if no error occurred and file was created
	require.NoError(t, err)

	// Check if output file exists
	_, err = os.Stat(outputFilename)
	require.NoError(t, err)

	// Cleanup
	defer os.Remove(outputFilename)
}

func TestEncodeAndSave_InvalidFileCreation(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte("test data")
	bitDepth := uint8(3)
	outputFilename := "/invalid/path/test_output.jpg" // Invalid path

	// Execute
	handler := &JpegHandler{}
	err := handler.EncodeAndSave(coverImage, data, bitDepth, outputFilename, false)

	// Test if the correct error is returned
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create output file")
}

func TestEncodeAndSave_SuccessWithSpecificFilename(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte("another test data")
	bitDepth := uint8(3)
	outputFilename := "specific_output.jpg"

	// Execute
	handler := &JpegHandler{}
	err := handler.EncodeAndSave(coverImage, data, bitDepth, outputFilename, false)

	// Test if no error occurred and file was created
	require.NoError(t, err)

	// Check if the output file exists
	_, err = os.Stat(outputFilename)
	require.NoError(t, err)

	// Cleanup
	defer os.Remove(outputFilename)
}

func TestEncodeAndSave_NullImage(t *testing.T) {
	// Setup
	var coverImage image.Image = nil
	data := []byte("some data")
	bitDepth := uint8(3)
	outputFilename := "test_output.jpg"

	// Execute
	handler := &JpegHandler{}
	err := handler.EncodeAndSave(coverImage, data, bitDepth, outputFilename, false)

	// Test if the correct error is returned
	assert.Error(t, err)
}
