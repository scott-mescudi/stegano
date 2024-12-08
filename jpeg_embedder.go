package stegano

import (
	"fmt"
	c "github.com/scott-mescudi/stegano/compression"
	s "github.com/scott-mescudi/stegano/jpeg"
	u "github.com/scott-mescudi/stegano/pkg"
	"image"
)

type JpegEmbedder struct {
}

// NewJpegEncoder initializes and returns a new JPEG encoder instance
// for embedding or extracting data.
func NewJpegEncoder() JpegEmbedder {
	return JpegEmbedder{}
}

// GetImageCapacity calculates the maximum amount of data (in bytes)
// that can be embedded into the given JPEG image.
func (m JpegEmbedder) GetImageCapacity(coverImage image.Image, bitDepth uint8) int {
	if bitDepth > 7 {
		return 0
	}

	return ((len(s.ExtractRGBChannelsFromJpeg(coverImage)) * 3) / 8) * (int(bitDepth)+1)
}

// EmbedDataIntoRgbChannels embeds the provided data into the RGB channels
// of the given JPEG image. Compression can be applied if `defaultCompression` is true.
func (m JpegEmbedder) EmbedDataIntoRgbChannels(coverImage image.Image, data []byte, bitDepth uint8, defaultCompression bool) ([]u.RgbChannel, error) {
	RGBchannels := s.ExtractRGBChannelsFromJpeg(coverImage)
	if len(data)*8 > len(RGBchannels)*3 {
		return nil, fmt.Errorf("error: Data too large to embed into the image")
	}

	var indata []byte = data
	if defaultCompression {
		compressedData, err := c.CompressZSTD(data)
		if err != nil {
			return nil, err
		}
		indata = compressedData
	}

	embeddedRGBChannels, err := u.EmbedIntoRGBchannelsWithDepth(RGBchannels, indata, bitDepth)
	if err != nil {
		return nil, err
	}

	return embeddedRGBChannels, nil
}

// ExtractDataFromRgbChannels retrieves embedded data from the RGB channels
// of a JPEG image. Decompression is applied if `isDefaultCompressed` is true.
func (m JpegEmbedder) ExtractDataFromRgbChannels(RGBchannels []u.RgbChannel, bitDepth uint8, isDefaultCompressed bool) ([]byte, error) {
	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil,  err
	}

	lenData, err := u.GetlenOfData(data)
	if err != nil || lenData == 0 {
		return nil, err
	}

	var moddedData = make([]byte, 0)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, data[i])
	}

	if isDefaultCompressed {
		outdata, err := c.DecompressZSTD(moddedData)
		if err != nil {
			return nil, err
		}

		return outdata, nil
	}

	return moddedData, nil
}

// EncodeJPEGImage embeds data into a JPEG image and saves it as a new file.
// Compression can be applied if `defaultCompression` is true.
func (m JpegEmbedder) EncodeJPEGImage(coverImage image.Image, data []byte, bitDepth uint8, outputFilename string, defaultCompression bool) error {
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()

	RGBchannels := s.ExtractRGBChannelsFromJpeg(coverImage)
	if len(data)*8 > len(RGBchannels)*3 {
		return fmt.Errorf("error: Data too large to embed into the image")
	}

	var indata []byte = data
	if defaultCompression {
		compressedData, err := c.CompressZSTD(data)
		if err != nil {
			return err
		}
		indata = compressedData
	}

	embeddedRGBChannels, err := u.EmbedIntoRGBchannelsWithDepth(RGBchannels, indata, bitDepth)
	if err != nil {
		return err
	}

	err = s.SaveImage(embeddedRGBChannels, outputFilename, height, width)
	if err != nil {
		return err
	}

	return nil
}

// DecodeJPEGImage extracts embedded data from a JPEG image.
// Decompression is applied if `isDefaultCompressed` is true.
func (m JpegEmbedder) DecodeJPEGImage(coverImage image.Image, bitDepth uint8, isDefaultCompressed bool) ([]byte, error) {
	RGBchannels := s.ExtractRGBChannelsFromJpeg(coverImage)
	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil, err
	}

	lenData, err := u.GetlenOfData(data)
	if err != nil || lenData == 0 {
		return nil, err
	}

	// Extract the actual embedded data
	var moddedData = make([]byte, 0)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, data[i])
	}

	if isDefaultCompressed {
		outdata, err := c.DecompressZSTD(moddedData)
		if err != nil {
			return nil, err
		}

		return outdata, nil
	}

	return moddedData, nil
}
