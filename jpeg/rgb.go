package jpeg

import (
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


func getLSB(value uint32) uint8 {
	return uint8(value & 1)
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
	lenOfDataInBinary := int32ToBinary(int32(len(data)))
	binaryData := bytesToBinary(data)
	combinedData := append(lenOfDataInBinary, binaryData...)

	z := splitIntoGroupsOfThree(combinedData)

	for i := 0; i < len(z); i++ {
		if z[i].r != getLSB(RGBchannels[i].R) {
			RGBchannels[i].R = flipLSB(RGBchannels[i].R)
		}
		if z[i].g != getLSB(RGBchannels[i].G) {
			RGBchannels[i].G = flipLSB(RGBchannels[i].G)
		}
		if z[i].b != getLSB(RGBchannels[i].B) {
			RGBchannels[i].B = flipLSB(RGBchannels[i].B)
		}
	}

	return RGBchannels
}

func ExtractDataFromRGBchannels(RGBchannels []RgbChannel) []byte {
	var byteSlice = make([]byte, 0)
	var currentByte uint8 = 0
	bitCount := 0

	for i := 0; i < len(RGBchannels); i++ {
		r := getLSB(RGBchannels[i].R)
		currentByte = (currentByte << 1) | (r & 1)
		bitCount++

		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}

		g := getLSB(RGBchannels[i].G)
		currentByte = (currentByte << 1) | (g & 1)
		bitCount++

		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}

		b := getLSB(RGBchannels[i].B)
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