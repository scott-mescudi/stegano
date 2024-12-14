package stegano

import "fmt"

type embedHandler struct {
	concurrency int
}

type extractHandler struct {
	concurrency int
}

type secureEmbedHandler struct {
	concurrency int
}

type secureExtractHandler struct {
	concurrency int
}

// NewEmbedHandlerWithConcurrency initializes an embedHandler with a specified concurrency
func NewEmbedHandlerWithConcurrency(concurrency int) (*embedHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &embedHandler{concurrency: concurrency}, nil
}

// NewEmbedHandler initializes an embedHandler with default concurrency
func NewEmbedHandler() (*embedHandler) {
	return &embedHandler{concurrency: 1}
}

// NewExtractHandlerWithConcurrency initializes an extractHandler with a specified concurrency
func NewExtractHandlerWithConcurrency(concurrency int) (*extractHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &extractHandler{concurrency: concurrency}, nil
}

// NewExtractHandler initializes an extractHandler with default concurrency
func NewExtractHandler() (*extractHandler) {
	return &extractHandler{concurrency: 1}
}

// NewSecureEmbedHandlerWithConcurrency initializes a secureEmbedHandler with a specified concurrency
func NewSecureEmbedHandlerWithConcurrency(concurrency int) (*secureEmbedHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &secureEmbedHandler{concurrency: concurrency}, nil
}

// NewSecureEmbedHandler initializes a secureEmbedHandler with default concurrency
func NewSecureEmbedHandler() (*secureEmbedHandler) {
	return &secureEmbedHandler{concurrency: 1}
}

// NewSecureExtractHandlerWithConcurrency initializes a secureExtractHandler with a specified concurrency
func NewSecureExtractHandlerWithConcurrency(concurrency int) (*secureExtractHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &secureExtractHandler{concurrency: concurrency}, nil
}

// NewSecureExtractHandler initializes a secureExtractHandler with default concurrency
func NewSecureExtractHandler() (*secureExtractHandler) {
	return &secureExtractHandler{concurrency: 1}
}
