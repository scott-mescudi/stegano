package jpeg

import (
        "image"
        s "github.com/scott-mescudi/stegano/pkg"
        
)

func ExtractRGBChannelsFromJpeg(img image.Image) []s.RgbChannel {
	bounds := img.Bounds()
	var pixels []s.RgbChannel

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, s.RgbChannel{R: r >> 8, G: g >> 8, B: b >> 8})
		}
	}

	return pixels
}
