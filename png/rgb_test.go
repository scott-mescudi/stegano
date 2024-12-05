package png

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
)

func TestExtractRGBChannelsFromImage(t *testing.T) {
	// Define image dimensions
	width, height := 2, 2

	// Create a new blank image
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Set pixel colors
	img.Set(0, 0, color.RGBA{0, 255, 0, 255})     // Top-left: Green
	img.Set(1, 1, color.RGBA{255, 0, 0, 255})     // Bottom-right: Red
	img.Set(1, 0, color.RGBA{255, 255, 255, 255}) // Top-right: White
	img.Set(0, 1, color.RGBA{255, 255, 255, 255}) // Bottom-left: White

	var buffer bytes.Buffer
	if err := png.Encode(&buffer, img); err != nil {
		t.Fatalf("failed to encode image to PNG: %v", err)
	}

	testimage, err := png.Decode(&buffer)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	//  function returns 16 bit colors
	expected := []RgbChannel{
		{R: 0, G: 65535, B: 0},         // Green
		{R: 65535, G: 65535, B: 65535}, // White
		{R: 65535, G: 65535, B: 65535}, // White
		{R: 65535, G: 0, B: 0},         // Red
	}

	result := ExtractRGBChannelsFromImage(testimage)
	fmt.Println(result)

	if len(result) != len(expected) {
		t.Fatalf("expected %d pixels, got %d", len(expected), len(result))
	}

	for i, pixel := range result {
		if pixel != expected[i] {
			t.Errorf("pixel %d: expected %+v, got %+v", i, expected[i], pixel)
		}
	}
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
	if err := png.Encode(&buffer, img); err != nil {
		t.Fatalf("failed to encode image to PNG: %v", err)
	}

	testimage, err := png.Decode(&buffer)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	testimg := ExtractRGBChannelsFromImage(testimage)

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
	if err := png.Encode(&buffer, img); err != nil {
		t.Fatalf("failed to encode image to PNG: %v", err)
	}

	testimage, err := png.Decode(&buffer)
	if err != nil {
		t.Fatalf("failed to decode PNG: %v", err)
	}

	testimg := ExtractRGBChannelsFromImage(testimage)

	chans := EmbedIntoRGBchannels(testimg, []byte("epaovrhbvohebr"))
	dt := ExtractDataFromRGBchannels(chans)

	if !strings.Contains(string(dt), "epaovrhbvohebr") {
		t.Fatalf("Failed to get data from image")
	}
}
