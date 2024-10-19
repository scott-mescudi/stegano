package main

import (
	"fmt"
	"image"
	"image/png"
	"image/jpeg"
	"os"
	"path/filepath"
)

type rgbChannel struct {
	r, g, b uint32
}

type lsb struct {
	r, g, b uint8
}

type bin struct {
	r, g, b uint8
}



func getLSB(value uint32) uint8 {
	return uint8(value & 1)
}

func extractRGBChannels(img image.Image) []rgbChannel {
	var lsbs []rgbChannel
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y


	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get pixel at (x, y)
			r, g, b, _ := img.At(x, y).RGBA()
			lsbs = append(lsbs, rgbChannel{r, g, b})
		}
	}

	return lsbs
}


func getLsbFromChannels(channels []rgbChannel) []lsb {
	var lsbs []lsb

    for _, channel := range channels {
        lsb := lsb{
            r: getLSB(channel.r),
            g: getLSB(channel.g),
            b: getLSB(channel.b),
        }
        lsbs = append(lsbs, lsb)
    }

    return lsbs
}

// stringToBinary converts a string to a slice of bits (0s and 1s).
func bytesToBinary(data []byte) []int {
	var bits []int

	// Iterate over each byte in the byte slice
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bit := (b >> i) & 1 // Extract each bit from the byte
			bits = append(bits, int(bit)) // Append the bit to the slice
		}
	}

	return bits
}

func flipLSB(num uint32) uint32 {
    return num ^ 1 // Flip the LSB using XOR
}


func splitIntoGroupsOfThree(nums []int) []bin {
	var result []bin
	// Iterate through the slice in steps of 3
	for i := 0; i < len(nums); i += 3 {
		// Create a new bin for each group of three
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

func decodeImage(filename string) (image.Image, error){
	ext := filepath.Ext(filename)

    file, err := os.Open(filename)
    if err!= nil {
        fmt.Println("Error opening image:", err)
        return nil, err
    }
    defer file.Close()

	switch ext {
		case ".png":
            return png.Decode(file)
		case ".jpeg":
			return jpeg.Decode(file)
		case ".jpg":
			return jpeg.Decode(file)
        default:
        return nil, fmt.Errorf("unsupported file format: %s", ext)
	}
}
var zx int
func embed(RGBchannels []rgbChannel, data []byte) []rgbChannel {
	z := splitIntoGroupsOfThree(bytesToBinary(data))
	zx = len(z)
	for i := 0; i < len(z); i++ {
		if z[i].r != getLSB(RGBchannels[i].r){
			RGBchannels[i].r = flipLSB(RGBchannels[i].r)
		}

		if z[i].g != getLSB(RGBchannels[i].g){
            RGBchannels[i].g = flipLSB(RGBchannels[i].g)
        }
		
        if z[i].b != getLSB(RGBchannels[i].b){
            RGBchannels[i].b = flipLSB(RGBchannels[i].b)
        }
	}

	return RGBchannels
}

//this doesnt work but has the right idea
func extract(RGBchannels []rgbChannel) []byte {
	var byteSlice []byte
	var currentByte uint8 = 0
	bitCount := 0

	for i := 0; i < len(RGBchannels); i++ {
		// Extract LSB from the red channel
		r := getLSB(RGBchannels[i].r)
		currentByte = (currentByte << 1) | (r & 1)
		bitCount++

		// Extract LSB from the green channel
		g := getLSB(RGBchannels[i].g)
		currentByte = (currentByte << 1) | (g & 1)
		bitCount++

		// Extract LSB from the blue channel
		b := getLSB(RGBchannels[i].b)
		currentByte = (currentByte << 1) | (b & 1)
		bitCount++

		// If we've collected 8 bits, append the current byte
		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}
	}

	// // If there are remaining bits after processing, append the last byte (with zero padding)
	// if bitCount > 0 {
	// 	currentByte = currentByte << (8 - bitCount)  // Left-align the remaining bits
	// 	byteSlice = append(byteSlice, currentByte)
	// }

	return byteSlice
}

func main() {
	img, err := decodeImage("input.png")
	if err!= nil {
        fmt.Println("Error decoding image:", err)
        return
    }

	// Extract LSB image
	RGBchannels := extractRGBChannels(img)
		for i := range 8{
		fmt.Printf("%b\n",RGBchannels[i])
	}
	fmt.Println()

	embeddedRGBChannels := embed(RGBchannels, []byte("LSB"))


	for i := range zx{
		fmt.Printf("%b\n",embeddedRGBChannels[i])
	}




}

// iterate over each rgb channel and get the lsb then check if the bit == the data bit, if it does i leave it but if it dont then flip it?