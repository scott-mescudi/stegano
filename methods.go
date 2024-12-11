package stegano

import "fmt"

type embedHandler struct {
	concurrency int
}

type extractHandler struct {
	concurrency int
}

// NewembedHandler initializes an embedHandler
func NewEmbedHandlerWithConcurrency(concurrency int) (*embedHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &embedHandler{concurrency: concurrency}, nil
}

// NewembedHandler initializes an embedHandler
func NewEmbedHandler() (*embedHandler) {
	return &embedHandler{concurrency: 1}
}


// NewextractHandler initializes an extractHandler
func NewExtractHandlerWithConcurrency(concurrency int) (*extractHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &extractHandler{concurrency: concurrency}, nil
}

// NewextractHandler initializes an extractHandler
func NewExtractHandler() (*extractHandler) {
	return &extractHandler{concurrency: 1}
}
