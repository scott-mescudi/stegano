package jpeg

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

type rgbChannel struct {
	r, g, b uint32
}

func extractRGBChannelsFromJpeg(img image.Image) []rgbChannel {
	bounds := img.Bounds()
	var pixels []rgbChannel

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, rgbChannel{r: r >> 8, g: g >> 8, b: b >> 8})
		}
	}

	return pixels
}

func SaveImage(embeddedRGBChannels []rgbChannel, filename string, height, width int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := y*width + x
			rgb := embeddedRGBChannels[i]

			img.Set(x, y, color.RGBA{
				R: uint8(rgb.r),
				G: uint8(rgb.g),
				B: uint8(rgb.b),
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
