package png

import (
	"fmt"

	s "github.com/scott-mescudi/stegano/utils"
)


func EmbedIntoRGBchannelsWithDepth(RGBchannels []RgbChannel, data []byte, depth uint8) ([]RgbChannel, error) {
	if depth > 7 {
		return nil, fmt.Errorf("bit depth exeeds 7")
	}

	if (len(data)*8)+32 > len(RGBchannels) * 3 * int(depth) {
		return nil, fmt.Errorf("data is too big")
	}

	
	lenOfDataInBinary := s.Int32ToBinary(int32(len(data)))
	binaryData := s.BytesToBinary(data)
	combinedData := append(lenOfDataInBinary, binaryData...)

	z := splitIntoGroupsOfThree(combinedData)
	
	curbit := depth
	index := 0

	for i := 0; i < len(z); i++ {
		if z[i].r != s.GetBit(RGBchannels[index].R, curbit) {
			RGBchannels[index].R = s.FlipBit(RGBchannels[index].R, curbit)
		}

		if z[i].g != s.GetBit(RGBchannels[index].G, curbit) {
			RGBchannels[index].G = s.FlipBit(RGBchannels[index].G, curbit)
		}

		if z[i].b != s.GetBit(RGBchannels[index].B, curbit) {
			RGBchannels[index].B = s.FlipBit(RGBchannels[index].B, curbit)
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

func ExtractDataFromRGBchannelsWithDepth(RGBchannels []RgbChannel, depth uint8) []byte {
	var byteSlice = make([]byte, 0)
	var currentByte uint8 = 0
	bitCount := 0

	for i := 0; i < len(RGBchannels); i++ {
		for bd := depth+1; bd > 0; bd-- {
			r := s.GetBit(RGBchannels[i].R, bd-1)
			currentByte = (currentByte << 1) | (r & 1)
			bitCount++

			if bitCount == 8 {
				byteSlice = append(byteSlice, currentByte)
				currentByte = 0
				bitCount = 0
			}

			g := s.GetBit(RGBchannels[i].G, bd-1)
			currentByte = (currentByte << 1) | (g & 1)
			bitCount++

			if bitCount == 8 {
				byteSlice = append(byteSlice, currentByte)
				currentByte = 0
				bitCount = 0
			}

			b := s.GetBit(RGBchannels[i].B, bd-1)
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

	return byteSlice
}