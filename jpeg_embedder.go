package stegano

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	c "github.com/scott-mescudi/stegano/compression"
	u "github.com/scott-mescudi/stegano/pkg"
)

// EncodeAndSave embeds the provided data into the given JPEG image and saves the modified image to a new file.
// The data is embedded using the specified bit depth. If `defaultCompression` is true, the data is compressed before embedding.
// The output filename can be specified. If it is an empty string, the default output filename (`stegano_out.jpg`) will be used.
// Returns an error if the data exceeds the embedding capacity of the image or if the saving process fails.
//
// Parameters:
// - coverImage: The original JPEG image where data will be embedded.
// - data: The data to embed into the image.
// - bitDepth: The number of bits per channel used for embedding (0-7).
// - outputFilename: The name of the file where the modified image will be saved. If empty, `stegano_out.png` is used.
// - defaultCompression: A flag indicating whether the data should be compressed before embedding.
//
// Returns:
// - An error if the embedding or saving process fails, such as if the data is too large, or if the image dimensions are invalid.
func (m *JpegHandler) EncodeAndSave(coverImage image.Image, data []byte, bitDepth uint8, outputFilename string, defaultCompression bool) error {
	if coverImage == nil {
		return fmt.Errorf("coverimage is nil")
	}

	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()

	if height <= 0 || width <= 0 {
		return fmt.Errorf("image size is invalid: height=%d, width=%d", height, width)
	}

	if len(data) == 0 {
		return errors.New("data is empty")
	}

	RGBchannels := u.ExtractRGBChannelsFromImage(coverImage)
	if len(data)*8 > (((len(RGBchannels)) * 3) / 8) * (int(bitDepth) + 1) {
		return fmt.Errorf("data is too large to embed into the image: required space exceeds available capacity")
	}

	var indata []byte = data
	if defaultCompression {
		compressedData, err := c.CompressZSTD(data)
		if err != nil {
			return fmt.Errorf("failed to compress data: %w", err)
		}
		indata = compressedData
	}

	embeddedRGBChannels, err := u.EmbedIntoRGBchannelsWithDepth(RGBchannels, indata, bitDepth)
	if err != nil {
		return fmt.Errorf("failed to embed data into RGB channels: %w", err)
	}

	imgdata, err := u.SaveImage(embeddedRGBChannels, height, width)
	if err != nil {
		return fmt.Errorf("failed to create embedded image: %w", err)
	}

	if outputFilename == "" {
		outputFilename = DefaultpngOutputFileName
	}

	file, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputFilename, err)
	}
	defer file.Close()

	if err := jpeg.Encode(file, imgdata, nil); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

// Decode extracts data embedded in a JPEG image using the specified bit depth.
// If the embedded data was compressed, it will be decompressed when `isDefaultCompressed` is true.
// Returns the extracted data or an error if the extraction or decompression fails.
//
// Parameters:
// - coverImage: The JPEG image containing embedded data to be extracted.
// - bitDepth: The number of bits per channel used during the embedding process (1-7).
// - isDefaultCompressed: A flag indicating whether the embedded data was compressed before embedding.
//
// Returns:
// - The extracted data if successful, or an error if the extraction or decompression process fails.
func (m *JpegHandler) Decode(coverImage image.Image, bitDepth uint8, isDefaultCompressed bool) ([]byte, error) {
	if bitDepth < 1 || bitDepth > 7 {
		return nil, fmt.Errorf("bitDepth is out of range (1-7): %d", bitDepth)
	}

	RGBchannels := u.ExtractRGBChannelsFromImage(coverImage)
	if RGBchannels == nil {
		return nil, errors.New("failed to extract RGB channels from the image")
	}

	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to extract data from RGB channels: %w", err)
	}

	lenData, err := u.GetlenOfData(data)
	if err != nil || lenData == 0 {
		return nil, errors.New("failed to get length of extracted data or data is empty")
	}

	var moddedData = make([]byte, 0, lenData)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, data[i])
	}

	if isDefaultCompressed {
		outdata, err := c.DecompressZSTD(moddedData)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress extracted data: %w", err)
		}
		return outdata, nil
	}

	return moddedData, nil
}
