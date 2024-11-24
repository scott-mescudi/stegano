package stegano

import (
	"fmt"
	"image"
	c "github.com/scott-mescudi/stegano/compression"
	s "github.com/scott-mescudi/stegano/jpeg"
	
)

type JpegEmbedder struct {
}

func NewJpegEncoder() JpegEmbedder {
	return JpegEmbedder{}
}

func (m JpegEmbedder) GetImageCapacity(coverImage image.Image) int {
	return (len(s.ExtractRGBChannelsFromJpeg(coverImage)) * 3) / 8
}

func (m JpegEmbedder) EncodeJPEGImage(coverImage image.Image, data []byte, outputFilename string, defaultCompression bool) error {
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()

	RGBchannels := s.ExtractRGBChannelsFromJpeg(coverImage)
	if len(data)*8 > len(RGBchannels)*3 {
		return fmt.Errorf("error: Data too large to embed into the image")
	}

	var indata []byte = data
	if defaultCompression{
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

func (m JpegEmbedder) DecodeJPEGImage(coverImage image.Image, isDefaultCompressed bool) ([]byte, error) {
	RGBchannels := s.ExtractRGBChannelsFromJpeg(coverImage)
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