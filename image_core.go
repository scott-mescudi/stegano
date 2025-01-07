package stegano

import (
	"fmt"
	"image"

	u "github.com/scott-mescudi/stegano/pkg"
)

// EmbedDataIntoImage embeds the given data into the RGB channels of the specified image.
func (m *EmbedHandler) EmbedDataIntoImage(coverImage image.Image, data []byte, bitDepth uint8) (image.Image, error) {
	if coverImage == nil {
		return nil, ErrInvalidCoverImage
	}

	if m.concurrency <= 0 {
		m.concurrency = 1
	}

	RGBchannels := u.ExtractRGBChannelsFromImageWithConCurrency(coverImage, m.concurrency)
	if (len(data)*8)+32 > len(RGBchannels)*3*(int(bitDepth)+1) {
		return nil, ErrDataTooLarge
	}

	

	embeddedRGBChannels, err := u.EmbedIntoRGBchannelsWithDepth(RGBchannels, data, bitDepth)
	if err != nil {
		return nil, err
	}

	return u.SaveImage(embeddedRGBChannels, coverImage.Bounds().Dy(), coverImage.Bounds().Dx())
}

// ExtractDataFromImage retrieves data embedded in the RGB channels of the specified image.
func (m *ExtractHandler) ExtractDataFromImage(coverImage image.Image, bitDepth uint8) ([]byte, error) {
	if coverImage == nil {
		return nil, ErrInvalidCoverImage
	}

	if m.concurrency <= 0 {
		m.concurrency = 1
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

	var moddedData = make([]byte, 0, lenData)
	defer func() {
		if r := recover(); r != nil {
			moddedData = nil
			err = fmt.Errorf("fatal error: %v", r)
		}
	}()

	for i := 4; i < lenData+4; i++ {
		if i >= len(data) {
			return nil, fmt.Errorf("index out of range while accessing data: %d", i)
		}
		moddedData = append(moddedData, data[i])
	}

	return moddedData, nil
}

// EmbedAtDepth embeds the provided data into a specific bit depth of the RGB channels of the image.
// Unlike other embedding methods, this modifies a single bit per channel at the specified depth.
func (m *EmbedHandler) EmbedAtDepth(coverimage image.Image, data []byte, depth uint8) (image.Image, error) {
	if coverimage == nil {
		return nil, ErrInvalidCoverImage
	}

	if m.concurrency <= 0 {
		m.concurrency = 1
	}

	channels := u.ExtractRGBChannelsFromImageWithConCurrency(coverimage, m.concurrency)
	if channels == nil {
		return nil, ErrFailedToExtractRGB
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
		return nil, ErrInvalidCoverImage
	}

	if m.concurrency <= 0 {
		m.concurrency = 1
	}

	channels := u.ExtractRGBChannelsFromImageWithConCurrency(coverimage, m.concurrency)
	if channels == nil {
		return nil, ErrFailedToExtractRGB
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
