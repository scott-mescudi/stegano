package stegano

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// GetImageCapacity calculates the maximum amount of data (in bytes)
// that can be embedded in the given image, based on the specified bit depth.
// Returns 0 if the bit depth exceeds 7, as higher depths are unsupported.
func GetImageCapacity(coverImage image.Image, bitDepth uint8) int {
    if bitDepth > 7 {
        return 0
    }
 
    return ((coverImage.Bounds().Max.X * coverImage.Bounds().Max.Y * 3) / 8) * (int(bitDepth) + 1)
}

// takes in path to imaeg and returns image.image
func Decodeimage(path string) (image.Image, error) {
	ext := filepath.Ext(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	
	switch ext {
	case ".jpg":
		return jpeg.Decode(file)
	case ".jpeg":
		return jpeg.Decode(file)
	case "png":
		return png.Decode(file)
	}

	return nil, fmt.Errorf("invalid image format")
}