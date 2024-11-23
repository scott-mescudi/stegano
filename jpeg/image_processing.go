package jpeg

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func SaveImage(embeddedRGBChannels []RgbChannel, filename string, height, width int) error {
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
