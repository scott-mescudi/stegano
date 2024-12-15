package pkg

import (
	"fmt"
	"image"
	"image/color"
	"sync"
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

type split struct {
	start, end int
}

func splitTask(n int, imgy int) []split {
	if n <= 0 || imgy <= 0 {
		return nil
	}

	sizes := make([]split, n)
	pp := imgy / n
	remainder := imgy % n

	start := 0
	for i := 0; i < n; i++ {
		extra := 0
		if i < remainder {
			extra = 1
		}
		sizes[i] = split{start: start, end: start + pp + extra}
		start = sizes[i].end
	}

	return sizes
}

// use this wth numGoroutines := runtime.NumCPU()
func ExtractRGBChannelsFromImageWithConCurrency(img image.Image, numGoroutines int) []RgbChannel {
	bounds := img.Bounds()
	totalPixels := bounds.Dx() * bounds.Dy()
	pixels := make([]RgbChannel, totalPixels)

	splits := splitTask(numGoroutines, bounds.Max.Y)

	var wg sync.WaitGroup
	for _, s := range splits {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			idx := start * bounds.Dx()
			for y := start; y < end; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, _ := img.At(x, y).RGBA()
					pixels[idx] = RgbChannel{R: r >> 8, G: g >> 8, B: b >> 8}
					idx++
				}
			}
		}(s.start, s.end)
	}

	wg.Wait()
	return pixels
}
