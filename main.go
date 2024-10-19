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
func stringToBinary(s string) []int {
    var bits []int

    for _, char := range s {
        // Convert character to binary and append to bits slice
        for i := 7; i >= 0; i-- {
            // Use bitwise AND to extract each bit and append it
            bit := (char >> i) & 1
            bits = append(bits, int(bit))
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

func main() {
	img, err := decodeImage("input.png")
	if err!= nil {
        fmt.Println("Error decoding image:", err)
        return
    }

	// Extract LSB image
	RGBchannels := extractRGBChannels(img)

	z := splitIntoGroupsOfThree(stringToBinary("hi"))
	for _,v := range z{
		fmt.Println(v)
	}
	
	fmt.Println()

	for i := 0; i < len(z); i++ {
		fmt.Printf("{%v %v %v}\n", getLSB(RGBchannels[i].r), getLSB(RGBchannels[i].g), getLSB(RGBchannels[i].b))
	}

	

	for i := 0; i < len(z); i++ {
		if z[i].r != getLSB(RGBchannels[i].r){
			RGBchannels[i].r = flipLSB(RGBchannels[i].r);
		};

		if z[i].g != getLSB(RGBchannels[i].g){
            RGBchannels[i].g = flipLSB(RGBchannels[i].g);
        };
		
        if z[i].b != getLSB(RGBchannels[i].b){
            RGBchannels[i].b = flipLSB(RGBchannels[i].b);
        };

	}

	fmt.Println()

	for i := 0; i < len(z); i++ {
		fmt.Printf("{%v %v %v}\n", getLSB(RGBchannels[i].r), getLSB(RGBchannels[i].g), getLSB(RGBchannels[i].b))
	}
}




// iterate over each rgb channel and get the lsb then check if the bit == the data bit, if it does i leave it but if it dont then flip it?