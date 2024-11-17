package stegano

import (
	"fmt"
	"image"
	c "github.com/scott-mescudi/stegano/compression"
	x "github.com/scott-mescudi/stegano/jpeg"
	
)

type JpegEmbedder struct {
}

func NewJpegEncoder() JpegEmbedder {
	return JpegEmbedder{}
}

func (m JpegEmbedder) GetImageCapacity(coverImage image.Image) int {
	return (len(x.ExtractRGBChannelsFromJpeg(coverImage)) * 3) / 8
}

func (m JpegEmbedder) EncodeJPEGImage(coverImage image.Image, data []byte, outputFilename string) error {
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()

	RGBchannels := x.ExtractRGBChannelsFromJpeg(coverImage)
	if len(data)*8 > len(RGBchannels)*3 {
		return fmt.Errorf("error: Data too large to embed into the image")
	}

	compressedData, err := c.CompressZSTD(data)
	if err != nil {
		return err
	}
	
	embeddedRGBChannels := x.EmbedIntoRGBchannels(RGBchannels, compressedData)
	
	err = x.SaveImage(embeddedRGBChannels, outputFilename, height, width)
	if err != nil {
		return err
	}
	
	return nil
}

func (m JpegEmbedder) DecodeJPEGImage(coverImage image.Image) ([]byte, error) {
	RGBchannels := x.ExtractRGBChannelsFromJpeg(coverImage)
	data := x.ExtractDataFromRGBchannels(RGBchannels)

	lenData, err := x.GetlenOfData(data)
	if err != nil {
		return nil, err
	}

	// Extract the actual embedded data
	var moddedData = make([]byte, 0)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, data[i])
	}

	datas, err := c.DecompressZSTD(moddedData)
	if err != nil {
		return nil, err
	}

	return datas, nil
}