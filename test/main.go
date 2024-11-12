package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

type rgbChannel struct {
	r, g, b uint8
}

func ExtractRGBChannelsFromImage(img image.Image) []rgbChannel {
	var lsbs []rgbChannel
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Scale down the RGBA values from 16-bit to 8-bit
			lsbs = append(lsbs, rgbChannel{
				r: uint8(r >> 8),
				g: uint8(g >> 8),
				b: uint8(b >> 8),
			})
		}
	}

	return lsbs
}

func SaveImage(embeddedRGBChannels []rgbChannel, filename string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := y*width + x
			rgb := embeddedRGBChannels[i]

			img.Set(x, y, color.RGBA{
				R: rgb.r,
				G: rgb.g,
				B: rgb.b,
				A: 255,
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

func main() {
	// Open the JPEG file
	jpegFile, err := os.Open("../testimages/in/in.jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer jpegFile.Close()

	img, err := jpeg.Decode(jpegFile)
	if err != nil {
		log.Fatal(err)
	}

	chans := ExtractRGBChannelsFromImage(img)
	fmt.Printf("%v => %v\n", len(chans), img.Bounds().Dx()*img.Bounds().Dy())

	// Correct the order of width and height here
	err = SaveImage(chans, "sd.png", img.Bounds().Dx(), img.Bounds().Dy())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Image saved as 'sd.png'")
}
