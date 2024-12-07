package jpeg

import (
	s "github.com/scott-mescudi/stegano/utils"
	"image"
)

type RgbChannel struct {
	R, G, B uint32
}

func ExtractRGBChannelsFromJpeg(img image.Image) []RgbChannel {
	bounds := img.Bounds()
	var pixels []RgbChannel

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			pixels = append(pixels, RgbChannel{R: r >> 8, G: g >> 8, B: b >> 8})
		}
	}

	return pixels
}

type bin struct {
	r, g, b uint8
}


func splitIntoGroupsOfThree(nums []int) []bin {
	var result []bin
	for i := 0; i < len(nums); i += 3 {
		var b bin
		if i < len(nums) {
			b.r = uint8(nums[i])
		}
		if i+1 < len(nums) {
			b.g = uint8(nums[i+1])
		}
		if i+2 < len(nums) {
			b.b = uint8(nums[i+2])
		}
		result = append(result, b)
	}
	return result
}

func EmbedIntoRGBchannels(RGBchannels []RgbChannel, data []byte) []RgbChannel {
	lenOfDataInBinary := s.Int32ToBinary(int32(len(data)))
	binaryData := s.BytesToBinary(data)
	combinedData := append(lenOfDataInBinary, binaryData...)

	z := splitIntoGroupsOfThree(combinedData)

	for i := 0; i < len(z); i++ {
		if z[i].r != s.GetBit(RGBchannels[i].R, 0) {
			RGBchannels[i].R = s.FlipBit(RGBchannels[i].R, 0)
		}
		if z[i].g != s.GetBit(RGBchannels[i].G, 0) {
			RGBchannels[i].G = s.FlipBit(RGBchannels[i].G, 0)
		}
		if z[i].b != s.GetBit(RGBchannels[i].B, 0) {
			RGBchannels[i].B = s.FlipBit(RGBchannels[i].B, 0)
		}
	}

	return RGBchannels
}

func ExtractDataFromRGBchannels(RGBchannels []RgbChannel) []byte {
	var byteSlice = make([]byte, 0)
	var currentByte uint8 = 0
	bitCount := 0

	for i := 0; i < len(RGBchannels); i++ {
		r := s.GetBit(RGBchannels[i].R, 0)
		currentByte = (currentByte << 1) | (r & 1)
		bitCount++

		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}

		g := s.GetBit(RGBchannels[i].G, 0)
		currentByte = (currentByte << 1) | (g & 1)
		bitCount++

		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}

		b := s.GetBit(RGBchannels[i].B, 0)
		currentByte = (currentByte << 1) | (b & 1)
		bitCount++

		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}
	}

	if bitCount > 0 {
		currentByte = currentByte << (8 - bitCount)
		byteSlice = append(byteSlice, currentByte)
	}

	return byteSlice
}
