package stegano

import (
	"fmt"
	"image"
	c "github.com/scott-mescudi/stegano/compression"
	s "github.com/scott-mescudi/stegano/png"
)

type PngEmbedder struct {
}

func NewPngEncoder() PngEmbedder {
	return PngEmbedder{}
}

func (m PngEmbedder) GetImageCapacity(coverImage image.Image) int {
	return (len(s.ExtractRGBChannelsFromImage(coverImage)) * 3) / 8
}


/* 
Checks if a coverimage has data in it by extracting the len of data,
if len data is 0 it returns false since image is empty
*/


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