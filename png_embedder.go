package stegano

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"

	c "github.com/scott-mescudi/stegano/compression"
	u "github.com/scott-mescudi/stegano/pkg"
	s "github.com/scott-mescudi/stegano/png"
)

// EncodeAndSave embeds the provided data into the given image and saves the modified image to a new file.
// The data is embedded using the specified bit depth. If `defaultCompression` is true, the data is compressed before embedding.
// Returns an error if the data exceeds the embedding capacity of the image or if the saving process fails.

// Parameters:
// - coverImage: The original image where data will be embedded.
// - data: The data to embed into the image.
// - bitDepth: The number of bits per channel used for embedding (0-7).
// - outputFilename: The name of the file where the modified image will be saved.
// - defaultCompression: A flag indicating whether the data should be compressed before embedding.
func (m *PngHandler) EncodeAndSave(coverImage image.Image, data []byte, bitDepth uint8, outputFilename string, defaultCompression bool) error {
	// Validate coverImage dimensions
	if coverImage == nil {
		return errors.New("coverImage is nil")
	}
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()
	if height <= 0 || width <= 0 {
		return fmt.Errorf("image size is invalid: height=%d, width=%d", height, width)
	}

	// Validate bit depth
	if bitDepth < 1 || bitDepth > 7 {
		return fmt.Errorf("bitDepth is out of range (1-7): %d", bitDepth)
	}

	// Validate data
	if len(data) == 0 {
		return errors.New("data is empty")
	}

	// Extract RGB channels
	RGBchannels := s.ExtractRGBChannelsFromImage(coverImage)
	if RGBchannels == nil {
		return errors.New("failed to extract RGB channels from the image")
	}

	// Check if data fits into the image
	maxCapacity := (((len(RGBchannels)) * 3) / 8) * (int(bitDepth) + 1)
	if len(data)*8 > maxCapacity {
		return fmt.Errorf("data is too large to embed into the image: maxCapacity=%d bytes, dataSize=%d bytes", maxCapacity, len(data))
	}

	// Compress data if required
	var indata []byte = data
	if defaultCompression {
		compressedData, err := c.CompressZSTD(data)
		if err != nil {
			return fmt.Errorf("failed to compress data: %w", err)
		}
		indata = compressedData
	}

	// Embed data
	embeddedRGBChannels, err := u.EmbedIntoRGBchannelsWithDepth(RGBchannels, indata, bitDepth)
	if err != nil {
		return fmt.Errorf("failed to embed data into RGB channels: %w", err)
	}

	// Generate image from embedded RGB channels
	imgdata, err := u.SaveImage(embeddedRGBChannels, height, width)
	if err != nil {
		return fmt.Errorf("failed to create embedded image: %w", err)
	}

	// Use default filename if none provided
	if outputFilename == "" {
		outputFilename = DefaultpngOutputFileName
	}

	// Create file for output
	file, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputFilename, err)
	}
	defer file.Close()

	// Encode the PNG
	if err := png.Encode(file, imgdata); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// Decode extracts data embedded in an image using the specified bit depth.
// If the embedded data was compressed, it will be decompressed when `isDefaultCompressed` is true.
// Returns the extracted data or an error if the extraction or decompression fails.
//
// Parameters:
// - coverImage: The image containing embedded data to be extracted.
// - bitDepth: The bit depth used during the embedding process.
// - isDefaultCompressed: A flag indicating whether the embedded data was compressed.
func (m *PngHandler) Decode(coverImage image.Image, bitDepth uint8, isDefaultCompressed bool) ([]byte, error) {
	// Validate coverImage dimensions
	if coverImage == nil {
		return nil, errors.New("coverImage is nil")
	}
	if bitDepth < 1 || bitDepth > 7 {
		return nil, fmt.Errorf("bitDepth is out of range (1-7): %d", bitDepth)
	}

	// Extract RGB channels
	RGBchannels := s.ExtractRGBChannelsFromImage(coverImage)
	if RGBchannels == nil {
		return nil, errors.New("failed to extract RGB channels from the image")
	}

	// Extract data
	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from RGB channels: %w", err)
	}

	// Validate extracted data length
	lenData, err := u.GetlenOfData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to get length of extracted data: %w", err)
	}
	if lenData == 0 {
		return nil, errors.New("extracted data length is zero")
	}

	// Retrieve the actual embedded data
	var moddedData = make([]byte, 0, lenData)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, data[i])
	}

	// Decompress data if required
	if isDefaultCompressed {
		outdata, err := c.DecompressZSTD(moddedData)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress extracted data: %w", err)
		}
		return outdata, nil
	}

	return moddedData, nil
}
