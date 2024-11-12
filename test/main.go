package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
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


func flipLSB(num uint32) uint32 {
	return num ^ 1 
}	

func main() {
	inputPath := "../testimages/in/in.jpeg"
	outputPath := "output.png"


	file, err := os.Open(inputPath)
	if err != nil {
		return
	}
	defer file.Close()
	img, _ := jpeg.Decode(file)
	pixels := extractRGBChannelsFromJpeg(img)


	fmt.Printf("%08b", pixels[1000])

	pixels[1000].r = flipLSB(pixels[1000].r)
	pixels[1000].g = flipLSB(pixels[1000].g)
	pixels[1000].b = flipLSB(pixels[1000].b)

	fmt.Printf("%08b", pixels[1000])

	err = SaveImage(pixels,"test.png", img.Bounds().Dy(), img.Bounds().Dx())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("PNG image created successfully at", outputPath)
}
