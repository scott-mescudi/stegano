package stegano

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	u "github.com/scott-mescudi/stegano/pkg"
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
	case ".png":
		return png.Decode(file)
	}

	return nil, fmt.Errorf("invalid image format")
}

// SaveImage saves the provided image to the specified output file.
//
// Parameters:
//
//	outputfile: The path to the output PNG file. Must not be empty and must have a .png extension.
//	embeddedImage: The image to save. Must not be nil.
//
// Returns:
//
//	An error if the input is invalid or if an issue occurs during the file creation or encoding process.
func SaveImage(outputfile string, embeddedImage image.Image) error {
	if outputfile == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	if filepath.Ext(outputfile) != ".png" {
		return fmt.Errorf("output file must have a .png extension, got '%s'", filepath.Ext(outputfile))
	}

	if embeddedImage == nil {
		return fmt.Errorf("embeddedImage parameter cannot be nil")
	}

	ff, err := os.Create(outputfile)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %v", outputfile, err)
	}
	defer ff.Close()

	encoder := png.Encoder{
		CompressionLevel: png.NoCompression,
	}

	if err := encoder.Encode(ff, embeddedImage); err != nil {
		return fmt.Errorf("failed to encode image to file '%s': %v", outputfile, err)
	}

	return nil
}

func EncryptData(data []byte, password string) (ciphertext []byte, err error) {
	if password == "" {
		return nil, fmt.Errorf("invalid password")
	}

	if len(data) <= 0 {
		return nil, fmt.Errorf("data is empty")
	}

	return u.Encrypt(password, data)
}

func DecryptData(ciphertext []byte, password string) (plaintext []byte, err error) {
	if password == "" {
		return nil, fmt.Errorf("invalid password")
	}

	if len(ciphertext) <= 0 {
		return nil, fmt.Errorf("ciphertext is empty")
	}

	return u.Decrypt(password, ciphertext)
}
