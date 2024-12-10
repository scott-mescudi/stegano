package stegano

import "fmt"

type EmbedHandler struct {
	concurrency int
}

type ExtractHandler struct {
	concurrency int
}

// NewEmbedHandler initializes an EmbedHandler
func NewEmbedHandler(concurrency int) (*EmbedHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &EmbedHandler{concurrency: concurrency}, nil
}

// NewExtractHandler initializes an ExtractHandler
func NewExtractHandler(concurrency int) (*ExtractHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &ExtractHandler{concurrency: concurrency}, nil
}
