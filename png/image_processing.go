package png

import (
	"image"

	s "github.com/scott-mescudi/stegano/pkg"
)

func ExtractRGBChannelsFromImage(img image.Image) []s.RgbChannel {
	var lsbs []s.RgbChannel
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lsbs = append(lsbs, s.RgbChannel{R: r, G: g, B: b})
		}
	}

	return lsbs
}
