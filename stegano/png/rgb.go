package png

import "image"

type rgbChannel struct {
	r, g, b uint32
}

type bin struct {
	r, g, b uint8
}

func getLSB(value uint32) uint8 {
	return uint8(value & 1)
}

func ExtractRGBChannelsFromImage(img image.Image) []rgbChannel {
	var lsbs []rgbChannel
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lsbs = append(lsbs, rgbChannel{r, g, b})
		}
	}

	return lsbs
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

func EmbedIntoRGBchannels(RGBchannels []rgbChannel, data []byte) []rgbChannel {
	lenOfDataInBinary := int32ToBinary(int32(len(data)))
	binaryData := bytesToBinary(data)
	combinedData := append(lenOfDataInBinary, binaryData...)

	z := splitIntoGroupsOfThree(combinedData)

	for i := 0; i < len(z); i++ {
		if z[i].r != getLSB(RGBchannels[i].r) {
			RGBchannels[i].r = flipLSB(RGBchannels[i].r)
		}
		if z[i].g != getLSB(RGBchannels[i].g) {
			RGBchannels[i].g = flipLSB(RGBchannels[i].g)
		}
		if z[i].b != getLSB(RGBchannels[i].b) {
			RGBchannels[i].b = flipLSB(RGBchannels[i].b)
		}
	}

	return RGBchannels
}

func ExtractDataFromRGBchannels(RGBchannels []rgbChannel) []byte {
	var byteSlice = make([]byte, 0)
	var currentByte uint8 = 0
	bitCount := 0

	for i := 0; i < len(RGBchannels); i++ {
		r := getLSB(RGBchannels[i].r)
		currentByte = (currentByte << 1) | (r & 1)
		bitCount++

		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}

		g := getLSB(RGBchannels[i].g)
		currentByte = (currentByte << 1) | (g & 1)
		bitCount++

		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}

		b := getLSB(RGBchannels[i].b)
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
