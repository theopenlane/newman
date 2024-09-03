package scrubber

import (
	"html"
	"strings"
)

// defaultTextScrubber provides a basic implementation of the Scrubber interface for plain text content
type defaultTextScrubber struct{}

// Scrub scrubs plain text content by escaping special characters and trimming whitespace
func (s *defaultTextScrubber) Scrub(text string) string {
	return html.EscapeString(strings.TrimSpace(text))
}

var defaultTextScrubberInstance = &defaultTextScrubber{}

// DefaultTextScrubber returns an instance of defaultTextScrubber
func DefaultTextScrubber() Scrubber {
	return defaultTextScrubberInstance
}
