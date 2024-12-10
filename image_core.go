package stegano

import (
	"fmt"
	"image"

	u "github.com/scott-mescudi/stegano/pkg"
)

// EmbedDataIntoImage embeds the given data into the RGB channels of the specified image.
// Supports optional compression via the `defaultCompression` flag. Returns the modified
// image or an error if the data exceeds the embedding capacity of the image.
func (m *EmbedHandler) EmbedDataIntoImage(coverImage image.Image, data []byte, bitDepth uint8) (image.Image, error) {
	if coverImage == nil {
		return nil, fmt.Errorf("coverimage is nil")
	}

	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if len(data)*8 > (((len(RGBchannels))*3)/8)*(int(bitDepth)+1) {
		return nil, fmt.Errorf("error: Data too large to embed into the image")
	}

	var indata []byte = data

	embeddedRGBChannels, err := u.EmbedIntoRGBchannelsWithDepth(RGBchannels, indata, bitDepth)
	if err != nil {
		return nil, err
	}

	return u.SaveImage(embeddedRGBChannels, coverImage.Bounds().Dy(), coverImage.Bounds().Dx())
}

// ExtractDataFromImage retrieves data embedded in the RGB channels of the specified image.
// Decompresses the data if `isDefaultCompressed` is true. Returns the extracted data
// or an error if the process fails.
func (m *ExtractHandler) ExtractDataFromImage(coverImage image.Image, bitDepth uint8) ([]byte, error) {
	if coverImage == nil {
		return nil, fmt.Errorf("coverimage is nil")
	}

	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	data, err := u.ExtractDataFromRGBchannelsWithDepth(RGBchannels, bitDepth)
	if err != nil {
		return nil, err
	}

	lenData, err := u.GetlenOfData(data)
	if err != nil || lenData == 0 {
		return nil, err
	}

	var moddedData = make([]byte, 0)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, data[i])
	}

	return moddedData, nil
}

// EmbedAtDepth embeds the provided data into a specific bit depth of the RGB channels of the image.
// Unlike other embedding methods, this modifies a single bit per channel at the specified depth.
func (m *EmbedHandler) EmbedAtDepth(coverimage image.Image, data []byte, depth uint8) (image.Image, error) {
	if coverimage == nil {
		return nil, fmt.Errorf("coverimage is nil")
	}

	channels := u.ExtractRGBChannelsFromImageWithConCurrency(coverimage, m.concurrency)
	if channels == nil {
		return nil, fmt.Errorf("Failed to extract channels from image")
	}

	ec, err := u.EmbedAtDepth(channels, data, depth)
	if err != nil {
		return nil, err
	}

	return u.SaveImage(ec, coverimage.Bounds().Dy(), coverimage.Bounds().Dx())
}

// ExtractAtDepth extracts data embedded at a specific bit depth from the RGB channels of an image.
// Only retrieves data from the specified bit depth. Returns the extracted data or an error if the process fails.
func (m *ExtractHandler) ExtractAtDepth(coverimage image.Image, depth uint8) ([]byte, error) {
	if coverimage == nil {
		return nil, fmt.Errorf("coverimage is nil")
	}

	channels := u.ExtractRGBChannelsFromImageWithConCurrency(coverimage, m.concurrency)
	if channels == nil {
		return nil, fmt.Errorf("Failed to extract channels from image")
	}

	emdata, err := u.ExtractAtDepth(channels, depth)
	if err != nil {
		return nil, err
	}

	lenData, err := u.GetlenOfData(emdata)
	if err != nil || lenData == 0 {
		return nil, err
	}

	var moddedData = make([]byte, 0)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, emdata[i])
	}

	return moddedData, nil
}
