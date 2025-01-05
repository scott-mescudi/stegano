package stegano

import (
	c "github.com/scott-mescudi/stegano/compression"
	u "github.com/scott-mescudi/stegano/pkg"
)

// EmbedDataIntoWAVWithDepth embeds compressed data into a WAV file with a specified bit depth.
func (s *AudioEmbedHandler) EmbedDataIntoWAVWithDepth(audioFilename, outputFilename string, data []byte, bitDepth uint8) ( error) {
	if bitDepth >= 8 {
		return ErrDepthOutOfRange
	}

	decoder := GetAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return err
	}

	newdata, err := c.CompressZSTD(data)
	if err != nil {
		return err
	}

	buffer = u.EmbedDataWithDepthAudio(buffer, newdata, bitDepth)

	err = WriteAudioFile(outputFilename, decoder, buffer)
	if err != nil {
		return err
	}

	return nil
}

// ExtractDataFromWAVWithDepth extracts compressed data from a WAV file with a specified bit depth.
func (s *AudioExtractHandler) ExtractDataFromWAVWithDepth(audioFilename string, bitDepth uint8) ([]byte, error) {
	if bitDepth >= 8 {
		return nil, ErrDepthOutOfRange
	} 	
	
	decoder := GetAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	data := u.ExtractDataWithDepthAudio(buffer, bitDepth)
	newdata, err := c.DecompressZSTD(data)
	if err != nil {
		return nil, err
	}

	return newdata, nil
}

// EmbedDataIntoWAVAtDepth embeds compressed data into a WAV file at a specified bit depth.
func (s *AudioEmbedHandler) EmbedDataIntoWAVAtDepth(audioFilename, outputFilename string, data []byte, bitDepth uint8) error {
	if bitDepth >= 8 {
		return ErrDepthOutOfRange
	}

	decoder := GetAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return err
	}

	newdata, err := c.CompressZSTD(data)
	if err != nil {
		return err
	}

	buffer = u.EmbedDataAtDepthAudio(buffer, newdata, bitDepth)

	err = WriteAudioFile(outputFilename, decoder, buffer)
	if err != nil {
		return err
	}

	return nil
}

// ExtractDataFromWAVAtDepth extracts compressed data from a WAV file at a specified bit depth.
func (s *AudioExtractHandler) ExtractDataFromWAVAtDepth(audioFilename string, bitDepth uint8) ([]byte, error) {
	if bitDepth >= 8 {
		return nil, ErrDepthOutOfRange
	}

	decoder := GetAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	data := u.ExtractDataAtDepthAudio(buffer, bitDepth)
	newdata, err := c.DecompressZSTD(data)
	if err != nil {
		return nil, err
	}

	return newdata, nil
}
