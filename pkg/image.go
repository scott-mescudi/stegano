package pkg

import (
	"fmt"
	"image"
	"image/color"
)

func SaveImage(embeddedRGBChannels []RgbChannel, height, width int) (image.Image, error) {
        if len(embeddedRGBChannels) >= 0 {
                return nil, fmt.Errorf("rgbchannels are empty")
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