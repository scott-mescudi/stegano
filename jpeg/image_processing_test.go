package jpeg

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert" 
	s "github.com/scott-mescudi/stegano/pkg"
)

// Helper function to create a 1x1 image with a specific color.
func createTestImage(c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, c)
	return img
}

func TestExtractRGBChannelsFromJpeg_SinglePixel(t *testing.T) {
	// Test with a single pixel image with red color.
	img := createTestImage(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	expected := []s.RgbChannel{{R: 255, G: 0, B: 0}}

	// Call the function.
	result := ExtractRGBChannelsFromJpeg(img)

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromJpeg_BlackImage(t *testing.T) {
	// Test with a single pixel black image.
	img := createTestImage(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	expected := []s.RgbChannel{{R: 0, G: 0, B: 0}}

	// Call the function.
	result := ExtractRGBChannelsFromJpeg(img)

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromJpeg_AllSameColor(t *testing.T) {
	// Test with a small 2x2 image where all pixels are green.
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}

	expected := []s.RgbChannel{
		{R: 0, G: 255, B: 0},
		{R: 0, G: 255, B: 0},
		{R: 0, G: 255, B: 0},
		{R: 0, G: 255, B: 0},
	}

	// Call the function.
	result := ExtractRGBChannelsFromJpeg(img)

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromJpeg_VaryingColors(t *testing.T) {
	// Test with a 2x2 image where each pixel is a different color.
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255}) // Red
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255}) // Green
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255}) // Blue
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 0, A: 255}) // Yellow

	expected := []s.RgbChannel{
		{R: 255, G: 0, B: 0},  // Red
		{R: 0, G: 255, B: 0},  // Green
		{R: 0, G: 0, B: 255},  // Blue
		{R: 255, G: 255, B: 0}, // Yellow
	}

	// Call the function.
	result := ExtractRGBChannelsFromJpeg(img)

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromJpeg_BlankImage(t *testing.T) {
	// Test with a blank image (transparent).
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0}) // Transparent
		}
	}

	expected := []s.RgbChannel{
		{R: 0, G: 0, B: 0},
		{R: 0, G: 0, B: 0},
		{R: 0, G: 0, B: 0},
		{R: 0, G: 0, B: 0},
	}

	// Call the function.
	result := ExtractRGBChannelsFromJpeg(img)

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromJpeg_EmptyImage(t *testing.T) {
	// Test with an empty image (0x0).
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	// Call the function.
	result := ExtractRGBChannelsFromJpeg(img)

	// Assert that the result is an empty slice.
	assert.Empty(t, result)
}
