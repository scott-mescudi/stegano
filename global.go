package stegano

import (
	u "github.com/scott-mescudi/stegano/pkg"
	"image"
)

// GetImageCapacity calculates the maximum amount of data (in bytes)
// that can be embedded in the given image, based on the specified bit depth.
// Returns 0 if the bit depth exceeds 7, as higher depths are unsupported.
func GetImageCapacity(coverImage image.Image, bitDepth uint8) int {
	if bitDepth > 7 {
		return 0
	}

	return ((len(u.ExtractRGBChannelsFromImage(coverImage)) * 3) / 8) * (int(bitDepth) + 1)
}
