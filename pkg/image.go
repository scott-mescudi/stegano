package pkg

import (
	"fmt"
	"image"
	"image/color"
)

func SaveImage(embeddedRGBChannels []RgbChannel, height, width int) (image.Image, error) {
	if len(embeddedRGBChannels) <= 0 {
		return nil, fmt.Errorf("rgbchannels are empty")
	}

	if height <= 0 || width <= 0 {
		return nil, fmt.Errorf("Inavalid image dimensions")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := y*width + x
			rgb := embeddedRGBChannels[i]

			img.Set(x, y, color.RGBA{
				R: uint8(rgb.R),
				G: uint8(rgb.G),
				B: uint8(rgb.B),
				A: 255,
			})
		}
	}

	return img, nil
}

func ExtractRGBChannelsFromImage(img image.Image) []RgbChannel {
	bounds := img.Bounds()
	var pixels []RgbChannel

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r > 255 || g > 255 || b > 255 {
				pixels = append(pixels, RgbChannel{R: r >> 8, G: g >> 8, B: b >> 8})
			}else {
				pixels = append(pixels, RgbChannel{R: r, G: g, B: b})
			}
			
		}
	}

	return pixels
}