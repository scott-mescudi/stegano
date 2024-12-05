package jpeg

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"strings"
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

func TestSplitIntoGroupsOfThree(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		output []bin
	}{
		{
			name:   "empty input",
			input:  []int{},
			output: []bin{},
		},
		{
			name:   "less than 3 elements",
			input:  []int{1},
			output: []bin{{r: 1, g: 0, b: 0}},
		},
		{
			name:   "exactly 3 elements",
			input:  []int{1, 2, 3},
			output: []bin{{r: 1, g: 2, b: 3}},
		},
		{
			name:   "more than 3 elements",
			input:  []int{1, 2, 3, 4, 5, 6},
			output: []bin{{r: 1, g: 2, b: 3}, {r: 4, g: 5, b: 6}},
		},
		{
			name:   "non-multiple of 3 elements",
			input:  []int{1, 2, 3, 4, 5},
			output: []bin{{r: 1, g: 2, b: 3}, {r: 4, g: 5, b: 0}},
		},
		{
			name:   "just under a multiple of 3",
			input:  []int{1, 2},
			output: []bin{{r: 1, g: 2, b: 0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitIntoGroupsOfThree(tt.input)
			for i := range result {
				if result[i] != tt.output[i] {
					t.Errorf("Test %s failed, expected %v, got %v", tt.name, tt.output[i], result[i])
				}
			}
		})
	}
}

func TestEmbedIntoRGBchannels(t *testing.T) {
	// Create a new blank image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, img, nil); err != nil {
		t.Fatalf("failed to encode image to jpeg: %v", err)
	}

	testimage, err := jpeg.Decode(&buffer)
	if err != nil {
		t.Fatalf("failed to decode jpeg: %v", err)
	}

	testimg := ExtractRGBChannelsFromJpeg(testimage)

	chans := EmbedIntoRGBchannels(testimg, []byte("hello world"))
	dt := ExtractDataFromRGBchannels(chans)

	if !strings.Contains(string(dt), "hello world") {
		t.Fatalf("Failed to get data from image")
	}
}

func TestExtractDataFromRGBchannels(t *testing.T) {
	// Create a new blank image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	var buffer bytes.Buffer
	if err := jpeg.Encode(&buffer, img, nil); err != nil {
		t.Fatalf("failed to encode image to jpeg: %v", err)
	}

	testimage, err := jpeg.Decode(&buffer)
	if err != nil {
		t.Fatalf("failed to decode jpeg: %v", err)
	}

	testimg := ExtractRGBChannelsFromJpeg(testimage)

	chans := EmbedIntoRGBchannels(testimg, []byte("epaovrhbvohebr"))
	dt := ExtractDataFromRGBchannels(chans)

	if !strings.Contains(string(dt), "epaovrhbvohebr") {
		t.Fatalf("Failed to get data from image")
	}
}
