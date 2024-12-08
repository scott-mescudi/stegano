package stegano

type JpegHandler struct{}
type PngHandler struct{}

// NewJpegHandler initializes and returns a new JPEG Handler instance
// for embedding or extracting data.
func NewJpegHandler() *JpegHandler {
	return &JpegHandler{}
}

// NewPngHandler initializes and returns a new PNG Handler instance
// for embedding or extracting data.
func NewPngHandler() *PngHandler {
	return &PngHandler{}
}
