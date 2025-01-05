package pkg

import (
	"bytes"
	"fmt"
	"github.com/go-audio/audio"
	"testing"
)

func TestEmbedDataWithDepthAudio(t *testing.T) {
	tests := []struct {
		data     []byte
		bitDepth uint8
	}{
		{
			data:     []byte("Hello, world!"),
			bitDepth: 1,
		},
		{
			data:     []byte("Another test data"),
			bitDepth: 2,
		},
		{
			data:     []byte("Small"),
			bitDepth: 3,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("bitDepth=%d", tt.bitDepth), func(t *testing.T) {
			// Create a dummy audio buffer with arbitrary values
			buffer := &audio.IntBuffer{
				Data: make([]int, len(tt.data)*8000),
				Format: &audio.Format{
					SampleRate:  44100,
					NumChannels: 1, // Adjusted to use NumChannels instead of SampleWidth
				},
			}

			// Embed data in audio buffer with bit depth
			embededBuff, _ := EmbedDataWithDepthAudio(buffer, tt.data, tt.bitDepth)

			// Extract data from the audio buffer to verify embedding worked
			extractedData := ExtractDataWithDepthAudio(embededBuff, tt.bitDepth)

			// Compare extracted data with original data
			if !bytes.Contains(extractedData, tt.data) {
				t.Errorf("Extracted data does not match original data. Expected: %s, got: %s", tt.data, extractedData)
			}
		})
	}
}

func TestExtractDataWithDepthAudio(t *testing.T) {
	// Create a dummy audio buffer with embedded data
	data := []byte("Embed this data with bit depth")
	bitDepth := uint8(2)
	buffer := &audio.IntBuffer{
		Data: make([]int, len(data)*8),
		Format: &audio.Format{
			SampleRate:  44100,
			NumChannels: 1, // Adjusted to use NumChannels instead of SampleWidth
		},
	}

	// Embed data in audio buffer with bit depth
	embededBuffer, _ := EmbedDataWithDepthAudio(buffer, data, bitDepth)

	// Extract data from the audio buffer with specific bit depth
	extractedData := ExtractDataWithDepthAudio(embededBuffer, bitDepth)

	// Compare extracted data with original data
	if !bytes.Contains(extractedData, data) {
		t.Errorf("Extracted data does not match original data. Expected: %s, got: %s", data, extractedData)
	}
}

func TestEmptyDataEmbedding(t *testing.T) {
	// Test with empty data input
	data := []byte("")
	depth := uint8(2)
	buffer := &audio.IntBuffer{
		Data: make([]int, 100),
		Format: &audio.Format{
			SampleRate:  44100,
			NumChannels: 1, // Adjusted to use NumChannels instead of SampleWidth
		},
	}

	// Embed empty data at depth
	_, err := EmbedDataAtDepthAudio(buffer, data, depth)
	if err == nil {
		t.Error("Embedded empty data instead of returning error")
	}
}
