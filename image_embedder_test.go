package stegano

import (
	"errors"
	"image"
	"image/color"
	"os"
	"testing"
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

func TestEncode_Valid(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte("some secret data")
	bitDepth := uint8(3)
	outputFilename := "test_output.png"

	// Execute
	handler := &EmbedHandler{3}
	err := handler.Encode(coverImage, data, bitDepth, outputFilename, false)

	// Test if no error occurred and file was created
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check if output file exists
	if _, err := os.Stat(outputFilename); err != nil {
		t.Fatalf("output file does not exist: %v", err)
	}

	// Cleanup
	defer os.Remove(outputFilename)
}

func TestEncode_EmptyData(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte{}
	bitDepth := uint8(3)
	outputFilename := "test_output.png"

	// Execute
	handler := &EmbedHandler{3}
	err := handler.Encode(coverImage, data, bitDepth, outputFilename, false)

	// Test if the correct error is returned
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if !errors.Is(err, ErrInvalidData) {
		t.Fatalf("expected error: %v, got: %v", ErrInvalidData, err)
	}
}

func TestEncode_DataTooLarge(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := make([]byte, 10000) // This is more data than the image can hold
	bitDepth := uint8(3)
	outputFilename := "test_output.png"

	// Execute
	handler := &EmbedHandler{3}
	err := handler.Encode(coverImage, data, bitDepth, outputFilename, false)

	// Test if the correct error is returned
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
	if !errors.Is(err, ErrDataTooLarge) {
		t.Fatalf("expected error: %v, got: %v", ErrDataTooLarge, err)
	}
}

func TestEncode_CompressedData(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte("compressed test data")
	bitDepth := uint8(3)
	outputFilename := "test_compressed_output.png"

	// Execute with compression enabled
	handler := &EmbedHandler{3}
	err := handler.Encode(coverImage, data, bitDepth, outputFilename, true)

	// Test if no error occurred and file was created
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check if output file exists
	if _, err := os.Stat(outputFilename); err != nil {
		t.Fatalf("output file does not exist: %v", err)
	}

	// Cleanup
	defer os.Remove(outputFilename)
}

func TestEncode_InvalidFileCreation(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte("test data")
	bitDepth := uint8(3)
	outputFilename := "/invalid/path/test_output.png" // Invalid path

	// Execute
	handler := &EmbedHandler{3}
	err := handler.Encode(coverImage, data, bitDepth, outputFilename, false)

	// Test if the correct error is returned
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestEncode_SuccessWithSpecificFilename(t *testing.T) {
	// Setup
	coverImage := createTestImage()
	data := []byte("another test data")
	bitDepth := uint8(3)
	outputFilename := "specific_output.png"

	// Execute
	handler := &EmbedHandler{3}
	err := handler.Encode(coverImage, data, bitDepth, outputFilename, false)

	// Test if no error occurred and file was created
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Check if the output file exists
	if _, err := os.Stat(outputFilename); err != nil {
		t.Fatalf("output file does not exist: %v", err)
	}

	// Cleanup
	defer os.Remove(outputFilename)
}

func TestEncode_NullImage(t *testing.T) {
	// Setup
	var coverImage image.Image = nil
	data := []byte("some data")
	bitDepth := uint8(3)
	outputFilename := "test_output.png"

	// Execute
	handler := &EmbedHandler{3}
	err := handler.Encode(coverImage, data, bitDepth, outputFilename, false)

	// Test if the correct error is returned
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}
