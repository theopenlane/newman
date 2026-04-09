package scrubber

import "html"

// defaultTextScrubber provides a basic implementation of the Scrubber interface for plain text content
type defaultTextScrubber struct{}

// Scrub escapes HTML special characters in plain text content to prevent injection when rendered in an HTML context
func (s *defaultTextScrubber) Scrub(text string) string {
	return html.EscapeString(text)
}

var defaultTextScrubberInstance = &defaultTextScrubber{}

// DefaultTextScrubber returns an instance of defaultTextScrubber
func DefaultTextScrubber() Scrubber {
	return defaultTextScrubberInstance
}
