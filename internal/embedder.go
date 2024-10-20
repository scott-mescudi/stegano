package png

import (
	"fmt"
	"image"
	s "lsb/internal/png"
)

type PngEmbedder struct {
}

func NewPngEncoder() PngEmbedder {
	return PngEmbedder{}
}

func (m PngEmbedder) EncodePngImage(coverImage image.Image, data []byte, outputFilename string) error {
	height := coverImage.Bounds().Dy()
	width := coverImage.Bounds().Dx()

	RGBchannels := s.ExtractRGBChannelsFromImage(coverImage)
	if len(data)*8 > len(RGBchannels)*3 {
		return fmt.Errorf("error: Data too large to embed into the image")
	}

	embeddedRGBChannels := s.EmbedIntoRGBchannels(RGBchannels, data)
	err := s.SaveImage(embeddedRGBChannels, outputFilename, height, width)
	if err != nil {
		return err
	}

	return nil
}

func (m PngEmbedder) DecodePngImage(coverImage image.Image) ([]byte, error) {
	RGBchannels := s.ExtractRGBChannelsFromImage(coverImage)
	data := s.ExtractDataFromRGBchannels(RGBchannels)

	lenData, err := s.GetlenOfData(data)
	if err != nil {
		return nil, err
	}

	// Extract the actual embedded data
	var moddedData = make([]byte, 0)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, data[i])
	}

	return moddedData, nil
}
