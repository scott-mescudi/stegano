package stegano

import (
	"fmt"

	u "github.com/scott-mescudi/stegano/pkg"
	c "github.com/scott-mescudi/stegano/compression"
)

// EmbedDataIntoWAVWithDepth embeds compressed data into a WAV file with a specified bit depth.
func (s *AudioEmbedHandler) EmbedDataIntoWAVWithDepth(audioFilename, outputFilename string, data []byte, bitDepth uint8) ( error) {
	if bitDepth >= 8 {
		return ErrDepthOutOfRange
	}
	
	nd, err := c.CompressZSTD(data)
	if err != nil {
		return err
	}


	decoder := GetAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return err
	}

	buffer = u.EmbedDataWithDepthAudio(buffer, nd, bitDepth)

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
func (s *AudioEmbedHandler) EmbedDataIntoWAVAtDepth(audioFilename, outputFilename string, data []byte, bitDepth uint8) error {
	if bitDepth >= 8 {
		return ErrDepthOutOfRange
	}

	decoder := GetAudioData(audioFilename)
	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return err
	}


	buffer = u.EmbedDataAtDepthAudio(buffer, data, bitDepth)

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

	return data, nil
}
