package png

import (
	"fmt"

	s "github.com/scott-mescudi/stegano/utils"
)



func EmbedIntoRGBchannelsWithDepth(RGBchannels []RgbChannel, data []byte, depth uint8) ([]RgbChannel, error) {
	if depth > 7 {
		return nil, fmt.Errorf("bit depth exeeds 7")
	}

	lenOfDataInBinary := s.Int32ToBinary(int32(len(data)))
	binaryData := s.BytesToBinary(data)
	combinedData := append(lenOfDataInBinary, binaryData...)

	z := splitIntoGroupsOfThree(combinedData)

	for i := 0; i < len(z); i+=int(depth)+1 {
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

	return RGBchannels, nil
}

func ExtractDataFromRGBchannelsWithDepth(RGBchannels []RgbChannel, depth uint8) []byte {
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
