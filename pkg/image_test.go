package pkg

import (
	"image"
	"image/color"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveImage(t *testing.T) {
	t.Run("EmptyRGBChannels", func(t *testing.T) {
		embeddedRGBChannels := []RgbChannel{}
		height, width := 10, 10

		_, err := SaveImage(embeddedRGBChannels, height, width)
		if err == nil {
			t.Errorf("expected error for empty RGB channels, got nil")
		}
	})

	t.Run("InvalidDimensions", func(t *testing.T) {
		embeddedRGBChannels := []RgbChannel{{R: 255, G: 0, B: 0}}
		height, width := 0, 0

		_, err := SaveImage(embeddedRGBChannels, height, width)
		if err == nil {
			t.Errorf("expected error for invalid dimensions, got nil")
		}
	})

	t.Run("ValidImage", func(t *testing.T) {
		height, width := 2, 2
		embeddedRGBChannels := []RgbChannel{
			{R: 255, G: 0, B: 0},   // Red
			{R: 0, G: 255, B: 0},   // Green
			{R: 0, G: 0, B: 255},   // Blue
			{R: 255, G: 255, B: 0}, // Yellow
		}

		img, err := SaveImage(embeddedRGBChannels, height, width)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
			t.Errorf("expected image dimensions %dx%d, got %dx%d",
				width, height, img.Bounds().Dx(), img.Bounds().Dy())
		}

		// Check pixel values
		expectedColors := []color.RGBA{
			{R: 255, G: 0, B: 0, A: 255},
			{R: 0, G: 255, B: 0, A: 255},
			{R: 0, G: 0, B: 255, A: 255},
			{R: 255, G: 255, B: 0, A: 255},
		}

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				expected := expectedColors[y*width+x]
				actual := img.At(x, y).(color.RGBA)

				if expected != actual {
					t.Errorf("at (%d, %d): expected %v, got %v", x, y, expected, actual)
				}
			}
		}
	})
}

// Helper function to create a 1x1 image with a specific color.
func createTestImage(c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, c)
	return img
}

func TestExtractRGBChannelsFromImageWithConCurrency_SinglePixel(t *testing.T) {
	// Test with a single pixel image with red color.
	img := createTestImage(color.RGBA{R: 255, G: 0, B: 0, A: 255})
	expected := []RgbChannel{{R: 255, G: 0, B: 0}}

	// Call the function.
	result := ExtractRGBChannelsFromImageWithConCurrency(img, runtime.NumCPU())

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromImageWithConCurrency_BlackImage(t *testing.T) {
	// Test with a single pixel black image.
	img := createTestImage(color.RGBA{R: 0, G: 0, B: 0, A: 255})
	expected := []RgbChannel{{R: 0, G: 0, B: 0}}

	// Call the function.
	result := ExtractRGBChannelsFromImageWithConCurrency(img, runtime.NumCPU())

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromImageWithConCurrency_AllSameColor(t *testing.T) {
	// Test with a small 2x2 image where all pixels are green.
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}

	expected := []RgbChannel{
		{R: 0, G: 255, B: 0},
		{R: 0, G: 255, B: 0},
		{R: 0, G: 255, B: 0},
		{R: 0, G: 255, B: 0},
	}

	// Call the function.
	result := ExtractRGBChannelsFromImageWithConCurrency(img, runtime.NumCPU())

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromImageWithConCurrency_VaryingColors(t *testing.T) {
	// Test with a 2x2 image where each pixel is a different color.
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255})   // Red
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 255})   // Green
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 255})   // Blue
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 0, A: 255}) // Yellow

	expected := []RgbChannel{
		{R: 255, G: 0, B: 0},   // Red
		{R: 0, G: 255, B: 0},   // Green
		{R: 0, G: 0, B: 255},   // Blue
		{R: 255, G: 255, B: 0}, // Yellow
	}

	// Call the function.
	result := ExtractRGBChannelsFromImageWithConCurrency(img, runtime.NumCPU())

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromImageWithConCurrency_BlankImage(t *testing.T) {
	// Test with a blank image (transparent).
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			img.Set(x, y, color.RGBA{R: 0, G: 0, B: 0, A: 0}) // Transparent
		}
	}

	expected := []RgbChannel{
		{R: 0, G: 0, B: 0},
		{R: 0, G: 0, B: 0},
		{R: 0, G: 0, B: 0},
		{R: 0, G: 0, B: 0},
	}

	// Call the function.
	result := ExtractRGBChannelsFromImageWithConCurrency(img, runtime.NumCPU())

	// Assert that the result matches the expected output.
	assert.Equal(t, expected, result)
}

func TestExtractRGBChannelsFromImageWithConCurrency_EmptyImage(t *testing.T) {
	// Test with an empty image (0x0).
	img := image.NewRGBA(image.Rect(0, 0, 0, 0))

	// Call the function.
	result := ExtractRGBChannelsFromImageWithConCurrency(img, runtime.NumCPU())

	// Assert that the result is an empty slice.
	assert.Empty(t, result)
}



// Helper function to generate a test image with random colors
func generateTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: uint8((x + y) % 256),
				A: 255,
			})
		}
	}
	return img
}

func BenchmarkExtractRGBChannelsFromImageWithConCurrency(b *testing.B) {
	// Generate a sample image for testing
	width, height := 10000, 10000 // Modify as needed
	img := generateTestImage(width, height)

	// Reset the timer to exclude setup time from the benchmark
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = ExtractRGBChannelsFromImageWithConCurrency(img, runtime.NumCPU())
	}
}