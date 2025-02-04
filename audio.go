package stegano

import (
	"fmt"

	c "github.com/scott-mescudi/stegano/compression"
	u "github.com/scott-mescudi/stegano/pkg"
)

// EmbedDataIntoWAVWithDepth embeds compressed data into a WAV file with a specified bit depth.
func (s *AudioEmbedHandler) EmbedIntoWAVWithDepth(audioFilename, outputFilename string, data []byte, bitDepth uint8) error {
	if bitDepth >= 8 {
		return ErrDepthOutOfRange
	}

	nd, err := c.CompressZSTD(data)
	if err != nil {
		return err
	}

	decoder := LoadAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return err
	}

	buffer, err = u.EmbedDataWithDepthAudio(buffer, nd, bitDepth)
	if err != nil {
		return err
	}

	err = SaveAudioToFile(outputFilename, decoder, buffer)
	if err != nil {
		return err
	}

	return nil
}

// ExtractDataFromWAVWithDepth extracts compressed data from a WAV file with a specified bit depth.
func (s *AudioExtractHandler) ExtractFromWAVWithDepth(audioFilename string, bitDepth uint8) ([]byte, error) {
	if bitDepth >= 8 {
		return nil, ErrDepthOutOfRange
	}

	decoder := LoadAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	data := u.ExtractDataWithDepthAudio(buffer, bitDepth)

	lenData, err := u.GetlenOfData(data)
	if err != nil {
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

	nd, err := c.DecompressZSTD(moddedData)
	if err != nil {
		return nil, err
	}

	return nd, nil
}

// EmbedDataIntoWAVAtDepth embeds compressed data into a WAV file at a specified bit depth.
func (s *AudioEmbedHandler) EmbedIntoWAVAtDepth(audioFilename, outputFilename string, data []byte, bitDepth uint8) error {
	if bitDepth >= 8 {
		return ErrDepthOutOfRange
	}

	decoder := LoadAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return err
	}

	buffer, err = u.EmbedDataAtDepthAudio(buffer, data, bitDepth)
	if err != nil {
		return ErrInvalidData
	}

	err = SaveAudioToFile(outputFilename, decoder, buffer)
	if err != nil {
		return err
	}

	return nil
}

// ExtractDataFromWAVAtDepth extracts compressed data from a WAV file at a specified bit depth.
func (s *AudioExtractHandler) ExtractFromWAVAtDepth(audioFilename string, bitDepth uint8) ([]byte, error) {
	if bitDepth >= 8 {
		return nil, ErrDepthOutOfRange
	}

	decoder := LoadAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	data := u.ExtractDataAtDepthAudio(buffer, bitDepth)

	return data, nil
}
