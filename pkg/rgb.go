package pkg

import "fmt"


type RgbChannel struct {
	R, G, B uint32
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

func EmbedIntoRGBchannelsWithDepth(RGBchannels []RgbChannel, data []byte, depth uint8) ([]RgbChannel, error) {
	if depth > 7 {
		return nil, fmt.Errorf("bit depth exeeds 7")
	}

	if (len(data)*8)+32 > len(RGBchannels) * 3 * int(depth) {
		return nil, fmt.Errorf("data is too big")
	}

	
	lenOfDataInBinary := Int32ToBinary(int32(len(data)))
	binaryData := BytesToBinary(data)
	combinedData := append(lenOfDataInBinary, binaryData...)

	z := splitIntoGroupsOfThree(combinedData)
	
	curbit := depth
	index := 0

	for i := 0; i < len(z); i++ {
		if z[i].r != GetBit(RGBchannels[index].R, curbit) {
			RGBchannels[index].R = FlipBit(RGBchannels[index].R, curbit)
		}

		if z[i].g != GetBit(RGBchannels[index].G, curbit) {
			RGBchannels[index].G = FlipBit(RGBchannels[index].G, curbit)
		}

		if z[i].b != GetBit(RGBchannels[index].B, curbit) {
			RGBchannels[index].B = FlipBit(RGBchannels[index].B, curbit)
		}

		if curbit != 0 {
			curbit--
		} else {
			curbit = depth
			index++
		}
	}

	return RGBchannels, nil
}

func ExtractDataFromRGBchannelsWithDepth(RGBchannels []RgbChannel, depth uint8) ([]byte, error) {
	if depth > 7 {
		return nil, fmt.Errorf("bit depth exeeds 7")
	}

	var byteSlice = make([]byte, 0)
	var currentByte uint8 = 0
	bitCount := 0

	for i := 0; i < len(RGBchannels); i++ {
		for bd := depth+1; bd > 0; bd-- {
			r := GetBit(RGBchannels[i].R, bd-1)
			currentByte = (currentByte << 1) | (r & 1)
			bitCount++

			if bitCount == 8 {
				byteSlice = append(byteSlice, currentByte)
				currentByte = 0
				bitCount = 0
			}

			g := GetBit(RGBchannels[i].G, bd-1)
			currentByte = (currentByte << 1) | (g & 1)
			bitCount++

			if bitCount == 8 {
				byteSlice = append(byteSlice, currentByte)
				currentByte = 0
				bitCount = 0
			}

			b := GetBit(RGBchannels[i].B, bd-1)
			currentByte = (currentByte << 1) | (b & 1)
			bitCount++

			if bitCount == 8 {
				byteSlice = append(byteSlice, currentByte)
				currentByte = 0
				bitCount = 0
			}
		}
	}

	if bitCount > 0 {
		currentByte = currentByte << (8 - bitCount)
		byteSlice = append(byteSlice, currentByte)
	}

	return byteSlice, nil
}