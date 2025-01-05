package stegano

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

type AudioEmbedHandler struct {}
type AudioExtractHandler struct {}

func NewAudioEmbedder() *AudioEmbedHandler {
	return &AudioEmbedHandler{}
}

func NewAudioExtractHandler() *AudioExtractHandler {
	return &AudioExtractHandler{}
}


// NewEmbedHandlerWithConcurrency initializes an EmbedHandler with a specified concurrency
func NewEmbedHandlerWithConcurrency(concurrency int) (*EmbedHandler, error) {
	if concurrency <= 0 {
		return nil, ErrInvalidGoroutines
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
		return nil, ErrInvalidGoroutines
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
		return nil, ErrInvalidGoroutines
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
		return nil, ErrInvalidGoroutines
	}

	return &SecureExtractHandler{concurrency: concurrency}, nil
}

// NewSecureExtractHandler initializes a SecureExtractHandler with default concurrency
func NewSecureExtractHandler() *SecureExtractHandler {
	return &SecureExtractHandler{concurrency: 1}
}
