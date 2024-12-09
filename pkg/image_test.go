package pkg

import (
	"image/color"
	"testing"
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
