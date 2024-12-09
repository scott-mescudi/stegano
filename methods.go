package stegano

type ImageHandler struct{}

// NewImageHandler initializes and returns a new PNG Handler instance
// for embedding or extracting data.
func NewImageHandler() *ImageHandler {
	return &ImageHandler{}
}
