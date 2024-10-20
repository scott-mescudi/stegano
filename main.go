package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"

	"image/color"
)

type rgbChannel struct {
	r, g, b uint32
}

type bin struct {
	r, g, b uint8
}


func getLSB(value uint32) uint8 {
	return uint8(value & 1)
}

func extractRGBChannelsFromImage(img image.Image) []rgbChannel {
	var lsbs []rgbChannel
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y


	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			lsbs = append(lsbs, rgbChannel{r, g, b})
		}
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

func int32ToBinary(num int32) []int {
    var bits []int

    // Iterate over each bit position from 31 to 0
    for i := 31; i >= 0; i-- {
        bit := (num >> i) & 1 // Extract each bit from the integer
        bits = append(bits, int(bit)) // Append the bit to the slice
    }

    return bits
}

func GetlenOfData(test []byte) (int,error) {
    container := ""
    for i := 0; i < 4; i++ {
        b := test[i]
        binary := fmt.Sprintf("%08b", b)
        container += binary
    }

    n, err :=  strconv.ParseInt(container, 2, 32)
    if err!= nil {
        return 0, fmt.Errorf("Error parsing binary to int: %e", err)
    }

    return int(n), nil
}

func embedIntoRGBchannels(RGBchannels []rgbChannel, data []byte) []rgbChannel {
	lenOfDataInBinary := int32ToBinary(int32(len(data)))
	binaryData := bytesToBinary(data)

	combinedData := append(lenOfDataInBinary, binaryData...)

	z := splitIntoGroupsOfThree(combinedData)
	
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


func extractDataFromRGBchannels(RGBchannels []rgbChannel) []byte {
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

		// Extract LSB from the blue channel
		b := getLSB(RGBchannels[i].b)
		currentByte = (currentByte << 1) | (b & 1)
		bitCount++

		if bitCount == 8 {
			byteSlice = append(byteSlice, currentByte)
			currentByte = 0
			bitCount = 0
		}
	}

	// If there are remaining bits after processing, append the last byte (with zero padding)
	if bitCount > 0 {
		currentByte = currentByte << (8 - bitCount)  // Left-align the remaining bits
		byteSlice = append(byteSlice, currentByte)
	}

	return byteSlice
}

func saveImage(embeddedRGBChannels []rgbChannel, filename string, height, width int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Loop through the rgbSlice and set the pixels
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Get the index in the slice (assuming rgbSlice has enough elements)
			i := y*width + x
			rgb := embeddedRGBChannels[i]

			// Set the pixel at (x, y)
			img.Set(x, y, color.RGBA{
				R: uint8(rgb.r),
				G: uint8(rgb.g),
				B: uint8(rgb.b),
				A: 255, // Fully opaque
			})
		}
	}

	// Save the image to a file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	png.Encode(file, img)

	return nil
}

func main() {
	imagev, err := decodeImage("input.png")
	if err!= nil {
        fmt.Println("Error decoding image:", err)
        return
    }

	height := imagev.Bounds().Dy()
	width := imagev.Bounds().Dx()

	RGBchannels := extractRGBChannelsFromImage(imagev)

	str := "sigma sg=igma"

	embeddedRGBChannels := embedIntoRGBchannels(RGBchannels, []byte(str))
	data := extractDataFromRGBchannels(embeddedRGBChannels)
	lenData, err := GetlenOfData(data)
	if err!= nil {
        fmt.Println("Error getting length of data:", err)
        return
    }

	fmt.Printf("Length of data: %d\n", lenData)

	var moddedData = make([]byte, 0)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, data[i])
	}

	fmt.Println(string(moddedData))

	// Create a new RGBA image
	saveImage(embeddedRGBChannels, "sky.png", height, width)
	
}

//gonna strore datliek this:
// len of bytes to read -- huffman_encoded_data(data)

