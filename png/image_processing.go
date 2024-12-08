package png

import (
	"image"
	"image/color"
	"image/png"
	"os"

	s "github.com/scott-mescudi/stegano/pkg"
)

func SaveImage(embeddedRGBChannels []s.RgbChannel, filename string, height, width int) error {
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

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	return err

}

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