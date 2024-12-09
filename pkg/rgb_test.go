package pkg

import (
	"strings"
	"testing"
)

func TestEmbedAtDepth(t *testing.T) {
	rgbdata := []RgbChannel{
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
	}

	data := []byte("hey")

	for i := 0; i < 7; i++ {
		ec, err := EmbedAtDepth(rgbdata, data, uint8(i))
		if err != nil {
			t.Fatal(err)
		}

		if ec[0].R == 255 {
			t.Fatalf("Data didnt change at index %v", i)
		}
	}

	ec, err := EmbedAtDepth(rgbdata, data, uint8(10))
	if err == nil || ec != nil {
		t.Fatalf("Didnt retrun error when passing in out of bounds bitDepth")
	}
}

func TestExtractAtDepth(t *testing.T) {
	rgbdata := []RgbChannel{
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
		{R: 255, G: 255, B: 255},
	}

	data := []byte("hey")
	var depth uint8 = 7

	EmbedAtDepth(rgbdata, data, depth)
	edata, err := ExtractAtDepth(rgbdata, depth)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(edata), string(data)) {
		t.Fatal("ExtractedData is incomplete")
	}
}
