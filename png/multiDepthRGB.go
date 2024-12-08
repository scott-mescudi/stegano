package png

import (
	"fmt"

	s "github.com/scott-mescudi/stegano/utils"
)

// how the fuck am i going to do this?????
// maybe try double loop since data is already spilt in to 3 just iterate ofer groups and embed them based on the depths

// so with depth of 2(3 without 0 index) it would iterate over groups three times
// first iteration would embed at bit index 2 
// second iteration would embed on bitindex 1
// third iteratioon would embed on index 0

// but how will i do this?
// kms?

//same process for ewxtarction

// bro  there  is no fucking way to test tjhos
// aaaaaaaaaaahhdLSHREDDLJVFHAELHJFBVGLJHEFBVLJHWBHEFBVHGJEDFWKBGHJVBGJHDASWFVJGUHCDWUGVYFEVCITYGHbbbb

// bok this might work
// havnt tested it yet but if it does then omfg

// teste stuff remoe in prod
var test = []RgbChannel{{255, 255, 255},{255, 255, 255},{255, 255, 255},{255, 255, 255},{255, 255, 255},{255, 255, 255},{255, 255, 255},{255, 255, 255},{255, 255, 255},{255, 255, 255},{255, 255, 255},}

func printit(ch  []RgbChannel) {
	for _, v := range ch {
		fmt.Printf("{%08b, %08b, %08b\n}", v.R, v.G, v.B)
	}
}

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
