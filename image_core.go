package stegano

import (
	"fmt"
	c "github.com/scott-mescudi/stegano/compression"
	u "github.com/scott-mescudi/stegano/pkg"
	"image"
)

// GetImageCapacity calculates the maximum amount of data (in bytes)
// that can be embedded in the given image, based on the specified bit depth.
// Returns 0 if the bit depth exceeds 7, as higher depths are unsupported.
func (m *ImageHandler) GetImageCapacity(coverImage image.Image, bitDepth uint8) int {
	if bitDepth > 7 {
		return 0
	}

	return ((len(u.ExtractRGBChannelsFromImage(coverImage)) * 3) / 8) * (int(bitDepth) + 1)
}

// EmbedDataIntoImage embeds the given data into the RGB channels of the specified image.
// Supports optional compression via the `defaultCompression` flag. Returns the modified
// image or an error if the data exceeds the embedding capacity of the image.
func (m *ImageHandler) EmbedDataIntoImage(coverImage image.Image, data []byte, bitDepth uint8, defaultCompression bool) (image.Image, error) {
	if coverImage == nil {
		return nil, fmt.Errorf("coverimage is nil")
	}

	RGBchannels := u.ExtractRGBChannelsFromImage(coverImage)
	if len(data)*8 > (((len(RGBchannels))*3)/8)*(int(bitDepth)+1) {
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

	return u.SaveImage(embeddedRGBChannels, coverImage.Bounds().Dy(), coverImage.Bounds().Dx())
}

// ExtractDataFromImage retrieves data embedded in the RGB channels of the specified image.
// Decompresses the data if `isDefaultCompressed` is true. Returns the extracted data
// or an error if the process fails.
func (m *ImageHandler) ExtractDataFromImage(coverImage image.Image, bitDepth uint8, isDefaultCompressed bool) ([]byte, error) {
	if coverImage == nil {
		return nil, fmt.Errorf("coverimage is nil")
	}

	RGBchannels := u.ExtractRGBChannelsFromImage(coverImage)
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

	if isDefaultCompressed {
		outdata, err := c.DecompressZSTD(moddedData)
		if err != nil {
			return nil, err
		}

		return outdata, nil
	}

	return moddedData, nil
}

// EmbedAtDepth embeds the provided data into a specific bit depth of the RGB channels of the image.
// Unlike other embedding methods, this modifies a single bit per channel at the specified depth.
func (m *ImageHandler) EmbedAtDepth(coverimage image.Image, data []byte, depth uint8) (image.Image, error) {
	if coverimage == nil {
		return nil, fmt.Errorf("coverimage is nil")
	}

	channels := u.ExtractRGBChannelsFromImage(coverimage)
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
func (m *ImageHandler) ExtractAtDepth(coverimage image.Image, depth uint8) ([]byte, error) {
	if coverimage == nil {
		return nil, fmt.Errorf("coverimage is nil")
	}

	channels := u.ExtractRGBChannelsFromImage(coverimage)
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
	fmt.Println(lenData)

	var moddedData = make([]byte, 0)
	for i := 4; i < lenData+4; i++ {
		moddedData = append(moddedData, emdata[i])
	}

	return moddedData, nil
}
