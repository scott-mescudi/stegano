package stegano

import "fmt"

type EmbedHandler struct {
	concurrency int
}

type ExtractHandler struct {
	concurrency int
}

type SecureEmbedHandler struct {
	concurrency int
}

type SecureExtractHandler struct {
	concurrency int
}

// NewEmbedHandlerWithConcurrency initializes an EmbedHandler with a specified concurrency
func NewEmbedHandlerWithConcurrency(concurrency int) (*EmbedHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &EmbedHandler{concurrency: concurrency}, nil
}

// NewEmbedHandler initializes an EmbedHandler with default concurrency
func NewEmbedHandler() *EmbedHandler {
	return &EmbedHandler{concurrency: 1}
}

// NewExtractHandlerWithConcurrency initializes an ExtractHandler with a specified concurrency
func NewExtractHandlerWithConcurrency(concurrency int) (*ExtractHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &ExtractHandler{concurrency: concurrency}, nil
}

// NewExtractHandler initializes an ExtractHandler with default concurrency
func NewExtractHandler() *ExtractHandler {
	return &ExtractHandler{concurrency: 1}
}

// NewSecureEmbedHandlerWithConcurrency initializes a SecureEmbedHandler with a specified concurrency
func NewSecureEmbedHandlerWithConcurrency(concurrency int) (*SecureEmbedHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &SecureEmbedHandler{concurrency: concurrency}, nil
}

// NewSecureEmbedHandler initializes a SecureEmbedHandler with default concurrency
func NewSecureEmbedHandler() *SecureEmbedHandler {
	return &SecureEmbedHandler{concurrency: 1}
}

// NewSecureExtractHandlerWithConcurrency initializes a SecureExtractHandler with a specified concurrency
func NewSecureExtractHandlerWithConcurrency(concurrency int) (*SecureExtractHandler, error) {
	if concurrency <= 0 {
		return nil, fmt.Errorf("invalid number of goroutines")
	}

	return &SecureExtractHandler{concurrency: concurrency}, nil
}

// NewSecureExtractHandler initializes a SecureExtractHandler with default concurrency
func NewSecureExtractHandler() *SecureExtractHandler {
	return &SecureExtractHandler{concurrency: 1}
}
