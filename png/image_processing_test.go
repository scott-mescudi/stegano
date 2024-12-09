package png

import (
	"image"
	"image/color"
	"testing"

	s "github.com/scott-mescudi/stegano/pkg"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a 1x1 image with a specific color.
func createTestImage(c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, c)
	return img
}

func TestExtractRGBChannelsFromImage_BlackImage(t *testing.T) {
	// Test with a single pixel black image.
	img := createTestImage(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	expected := []s.RgbChannel{{R: 0, G: 0, B: 0}}

	// Call the function.
	result := ExtractRGBChannelsFromImage(img)

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromImage_BlankImage(t *testing.T) {
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
	result := ExtractRGBChannelsFromImage(img)

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromImage_EmptyImage(t *testing.T) {
	// Test with an empty image (0x0).
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	// Call the function.
	result := ExtractRGBChannelsFromImage(img)

	// Assert that the result is an empty slice.
	assert.Empty(t, result)
}
