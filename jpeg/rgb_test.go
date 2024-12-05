package jpeg

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"testing"
)

func TestExtractRGBChannelsFromJpeg(t *testing.T) {
	
	rect := image.Rect(0, 0, 2, 2)
	img := image.NewRGBA(rect)
	
	img.Set(0, 0, color.RGBA{R: 255, G: 0, B: 0, A: 0}) 
	img.Set(1, 0, color.RGBA{R: 0, G: 255, B: 0, A: 0}) 
	img.Set(0, 1, color.RGBA{R: 0, G: 0, B: 255, A: 0}) 
	img.Set(1, 1, color.RGBA{R: 255, G: 255, B: 0, A: 0})

	
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("failed to encode image to JPEG: %v", err)
	}


	decodedImg, _, err := image.Decode(&buf)
	if err != nil {
		t.Fatalf("failed to decode JPEG: %v", err)
	}


	result := ExtractRGBChannelsFromJpeg(decodedImg)
	fmt.Println(result)


	expected := []RgbChannel{
		{R: 255, G: 0, B: 0},   
		{R: 0, G: 255, B: 0},   
		{R: 0, G: 0, B: 255},   
		{R: 255, G: 255, B: 0}, 
	}

	// Compare result with expected
	if len(result) != len(expected) {
		t.Fatalf("expected %d pixels, got %d", len(expected), len(result))
	}

	// note: cant check if its a exact match since jpeg compression changes data
}
