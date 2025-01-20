package stegano

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	u "github.com/scott-mescudi/stegano/pkg"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
)

// Wrapper around rsEncode private function
func RsEncode(data []byte, parity int) ([]byte, error) {
	return u.RsEncode(data, parity)
}

// Wrapper around rsDecode private function
func RsDecode(packedDataShards []byte, parity int) ([]byte, error) {
	return u.RsDecode(packedDataShards, 1, parity)
}
// GetImageCapacity calculates the maximum amount of data (in bytes)
// that can be embedded in the given image, based on the specified bit depth.
// Returns 0 if the bit depth exceeds 7, as higher depths are unsupported.
func GetImageCapacity(coverImage image.Image, bitDepth uint8) int {
	if bitDepth > 7 {
		return 0
	}

	return ((coverImage.Bounds().Max.X * coverImage.Bounds().Max.Y * 3) / 8) * (int(bitDepth) + 1)
}

// takes in path to image and returns image.Image
func Decodeimage(path string) (image.Image, error) {
	ext := filepath.Ext(path)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	switch ext {
	case ".jpg":
		return jpeg.Decode(file)
	case ".jpeg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	}

	return nil, fmt.Errorf("invalid image format")
}

// SaveImage saves the provided image to the specified output file.
//
// Parameters:
//
//	outputfile: The path to the output PNG file. Must not be empty and must have a .png extension.
//	embeddedImage: The image to save. Must not be nil.
//
// Returns:
//
//	An error if the input is invalid or if an issue occurs during the file creation or encoding process.
func SaveImage(outputfile string, embeddedImage image.Image) error {
	if outputfile == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	if filepath.Ext(outputfile) != ".png" {
		return fmt.Errorf("output file must have a .png extension, got '%s'", filepath.Ext(outputfile))
	}

	if embeddedImage == nil {
		return fmt.Errorf("embeddedImage parameter cannot be nil")
	}

	ff, err := os.Create(outputfile)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %v", outputfile, err)
	}
	defer ff.Close()

	encoder := png.Encoder{
		CompressionLevel: png.NoCompression,
	}

	if err := encoder.Encode(ff, embeddedImage); err != nil {
		return fmt.Errorf("failed to encode image to file '%s': %v", outputfile, err)
	}

	return nil
}

// EncryptData encrypts the given data using the provided password.
// It returns the encrypted ciphertext or an error if the encryption fails.
//
// Parameters:
// - data ([]byte): The plaintext data to be encrypted.
// - password (string): The password to be used for encryption.
//
// Returns:
// - ciphertext ([]byte): The encrypted data.
// - err (error): An error if the encryption fails.
func EncryptData(data []byte, password string) (ciphertext []byte, err error) {
	if password == "" {
		return nil, fmt.Errorf("invalid password")
	}

	if len(data) <= 0 {
		return nil, fmt.Errorf("data is empty")
	}

	return u.Encrypt(password, data)
}

// DecryptData decrypts the given ciphertext using the provided password.
// It returns the decrypted plaintext or an error if the decryption fails.
//
// Parameters:
// - ciphertext ([]byte): The encrypted data to be decrypted.
// - password (string): The password to be used for decryption.
//
// Returns:
// - plaintext ([]byte): The decrypted data.
// - err (error): An error if the decryption fails.
func DecryptData(ciphertext []byte, password string) (plaintext []byte, err error) {
	if password == "" {
		return nil, fmt.Errorf("invalid password")
	}

	if len(ciphertext) <= 0 {
		return nil, fmt.Errorf("ciphertext is empty")
	}

	return u.Decrypt(password, ciphertext)
}

// GetAudioData opens the WAV file and returns a decoder
func LoadAudioData(file string) *wav.Decoder {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	decoder := wav.NewDecoder(f)

	// Decode the WAV file header and check if it's valid
	if !decoder.IsValidFile() {
		fmt.Println("Invalid WAV file")
		return nil
	}

	return decoder
}

// WriteAudioFile writes the decoded and modified data to a new WAV file
func SaveAudioToFile(fileName string, decoder *wav.Decoder, buffer *audio.IntBuffer) error {
	// Create a new file for writing the modified WAV data
	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error creating output file: %e\n", err)
	}
	defer outFile.Close()

	// Create a new encoder for the output file
	encoder := wav.NewEncoder(outFile, int(decoder.SampleRate), int(decoder.BitDepth), int(decoder.NumChans), 1)

	// Write the modified buffer to the new file
	if err := encoder.Write(buffer); err != nil {
		return fmt.Errorf("Error encoding WAV file: %e\n", err)
	}

	// Close the encoder to flush the output
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("Error closing encoder: %e\n", err)
	}

	return nil
}
