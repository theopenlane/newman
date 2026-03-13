package compose

import "errors"

var (
	// ErrInvalidURL is returned when a base URL cannot be parsed for token injection
	ErrInvalidURL = errors.New("invalid url")
	// ErrStructMarshal is returned when a struct cannot be marshaled to JSON for template data conversion
	ErrStructMarshal = errors.New("failed to marshal struct to template data")
)
