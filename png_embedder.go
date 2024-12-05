package stegano

import (
	"fmt"
	"image"
	c "github.com/scott-mescudi/stegano/compression"
	s "github.com/scott-mescudi/stegano/png"
)

type PngEmbedder struct {
}

// NewPngEncoder initializes and returns a new PNG encoder instance 
// for embedding or extracting data.
func NewPngEncoder() PngEmbedder {
	return PngEmbedder{}
}

// GetImageCapacity calculates the maximum amount of data (in bytes) 
// that can be embedded into the given image.
func (m PngEmbedder) GetImageCapacity(coverImage image.Image) int {
	return (len(s.ExtractRGBChannelsFromImage(coverImage)) * 3) / 8
}

// EmbedDataIntoRgbChannels embeds the provided data into the RGB channels 
// of the given image. Compression can be applied if `defaultCompression` is true.
func (m PngEmbedder) EmbedDataIntoRgbChannels(coverImage image.Image, data []byte, defaultCompression bool) ([]s.RgbChannel, error) {
	RGBchannels := s.ExtractRGBChannelsFromImage(coverImage)
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

	embeddedRGBChannels := s.EmbedIntoRGBchannels(RGBchannels, indata)

	
	return embeddedRGBChannels, nil
}

// ExtractDataFromRgbChannels retrieves embedded data from the RGB channels 
// of an image. Decompression is applied if `isDefaultCompressed` is true.
func (m PngEmbedder) ExtractDataFromRgbChannels(RGBchannels []s.RgbChannel, isDefaultCompressed bool) ([]byte, error) {
	data := s.ExtractDataFromRGBchannels(RGBchannels)

	lenData, err := s.GetlenOfData(data)
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

// HasData checks whether the given image contains any embedded data.
func (m PngEmbedder) HasData(coverImage image.Image) bool {
	var lsbs []s.RgbChannel
	bounds := coverImage.Bounds()

	for x := bounds.Min.X; x < 11; x++ {
		r, g, b, _ := coverImage.At(x, 0).RGBA()
		lsbs = append(lsbs, s.RgbChannel{R:r, G:g, B:b})
	}

	data := s.ExtractDataFromRGBchannels(lsbs)

	lenData, _ := s.GetlenOfData(data)
	fmt.Println(lenData)
	if lenData == 0{
		return false
	}

	return true
}

// EncodePngImage embeds data into an image and saves it as a new file. 
// Compression can be applied if `defaultCompression` is true.
func (m PngEmbedder) EncodePngImage(coverImage image.Image, data []byte, outputFilename string, defaultCompression bool) error {
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()

	RGBchannels := s.ExtractRGBChannelsFromImage(coverImage)
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

	embeddedRGBChannels := s.EmbedIntoRGBchannels(RGBchannels, indata)
	
	err := s.SaveImage(embeddedRGBChannels, outputFilename, height, width)
	if err != nil {
		return err
	}
	
	return nil
}

// DecodePngImage extracts embedded data from an image. 
// Decompression is applied if `isDefaultCompressed` is true.
func (m PngEmbedder) DecodePngImage(coverImage image.Image, isDefaultCompressed bool) ([]byte, error) {
	RGBchannels := s.ExtractRGBChannelsFromImage(coverImage)
	data := s.ExtractDataFromRGBchannels(RGBchannels)

	lenData, err := s.GetlenOfData(data)
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