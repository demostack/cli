package provider

import "errors"

var (
	// ErrFileNotFound .
	ErrFileNotFound = errors.New("file is not found")

	// ErrFileUnmarshalError .
	ErrFileUnmarshalError = errors.New("file cannot be unmarshalled")
)
