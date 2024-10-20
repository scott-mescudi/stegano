package imageutils

import (
	"image"
	"image/color"
	"image/png"
	"os"
    s "lsb/internal/imageSegoLib"
)


// saveImage saves the modified RGB channels as a PNG image.
func SaveImage(embeddedRGBChannels []s.RgbChannel, filename string, height, width int) error {
    img := image.NewRGBA(image.Rect(0, 0, width, height))

    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            i := y*width + x
            rgb := embeddedRGBChannels[i]

            img.Set(x, y, color.RGBA{
                R: uint8(rgb.R),
                G: uint8(rgb.B),
                B: uint8(rgb.G),
                A: 255, // Fully opaque
            })
        }
    }

    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    return png.Encode(file, img)
}
