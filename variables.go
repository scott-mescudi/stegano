package stegano

import "errors"

var (
	DefaultOutputFile string = "stegano_out.png"
)

const (
	LSB         uint8 = 0
	MaxBitDepth uint8 = 7
)

// Errors for image_embedder.go and image_core.go
var (
	ErrDepthOutOfRange      = errors.New("bitDepth is out of range (0-7)")
	ErrFailedToExtractRGB   = errors.New("failed to extract RGB channels from the image")
	ErrFailedToExtractData  = errors.New("failed to extract data from RGB channels")
	ErrInvalidDataLength    = errors.New("extracted data length is zero")
	ErrInvalidCoverImage    = errors.New("coverImage is nil or has invalid dimensions")
	ErrInvalidData          = errors.New("data is empty or invalid")
	ErrDataTooLarge         = errors.New("data exceeds the embedding capacity of the image")
	ErrFailedToCompressData = errors.New("failed to compress data")
	ErrFailedToDecryptData  = errors.New("failed to decrypt data")
	ErrFailedToSaveImage    = errors.New("failed to save image")
)

// Errors for methods.go
var (
	ErrInvalidGoroutines = errors.New("invalid number of goroutines")
)
