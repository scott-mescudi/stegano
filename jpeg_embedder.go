package stegano

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"

	c "github.com/scott-mescudi/stegano/compression"
	s "github.com/scott-mescudi/stegano/jpeg"
	u "github.com/scott-mescudi/stegano/pkg"
)

// EncodeAndSave embeds data into a JPEG image and saves it as a new file.
// Compression can be applied if `defaultCompression` is true.
func (m *JpegHandler) EncodeAndSave(coverImage image.Image, data []byte, bitDepth uint8, outputFilename string, defaultCompression bool) error {
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()

	RGBchannels := s.ExtractRGBChannelsFromJpeg(coverImage)
	if len(data)*8 > (((len(RGBchannels))*3)/8)*(int(bitDepth)+1) {
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

	imgdata, err := u.SaveImage(embeddedRGBChannels, height, width)
	if err != nil {
		return err
	}

	file, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	jpeg.Encode(file, imgdata, nil)

	return nil
}

// Decode extracts embedded data from a JPEG image.
// Decompression is applied if `isDefaultCompressed` is true.
func (m *JpegHandler) Decode(coverImage image.Image, bitDepth uint8, isDefaultCompressed bool) ([]byte, error) {
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
