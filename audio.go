package stegano

import (
	c "github.com/scott-mescudi/stegano/compression"
	u "github.com/scott-mescudi/stegano/pkg"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// EmbedDataIntoWAVWithDepth embeds compressed data into a WAV file with a specified bit depth.
func (s *AudioEmbedHandler) EmbedDataIntoWAVWithDepth(decoder *wav.Decoder, data []byte, bitDepth uint8) (*audio.IntBuffer, error) {
	if bitDepth >= 8 {
		return nil, ErrDepthOutOfRange
	} 

	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	newdata, err := c.CompressZSTD(data)
	if err != nil {
		return nil, err
	}

	buffer = u.EmbedDataWithDepthAudio(buffer, newdata, bitDepth)

	return buffer, nil
}

// ExtractDataFromWAVWithDepth extracts compressed data from a WAV file with a specified bit depth.
func (s *AudioExtractHandler) ExtractDataFromWAVWithDepth(decoder *wav.Decoder, bitDepth uint8) ([]byte, error) {
	if bitDepth >= 8 {
		return nil, ErrDepthOutOfRange
	} 	
	
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
func (s *AudioEmbedHandler) EmbedDataIntoWAVAtDepth(decoder *wav.Decoder, data []byte, bitDepth uint8) (*audio.IntBuffer, error) {
	if bitDepth >= 8 {
		return nil, ErrDepthOutOfRange
	} 

	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, err
	}

	newdata, err := c.CompressZSTD(data)
	if err != nil {
		return nil, err
	}

	buffer = u.EmbedDataAtDepthAudio(buffer, newdata, bitDepth)

	return buffer, nil
}

// ExtractDataFromWAVAtDepth extracts compressed data from a WAV file at a specified bit depth.
func (s *AudioExtractHandler) ExtractDataFromWAVAtDepth(decoder *wav.Decoder, bitDepth uint8) ([]byte, error) {
	if bitDepth >= 8 {
		return nil, ErrDepthOutOfRange
	} 	
	
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