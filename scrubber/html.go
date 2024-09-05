package scrubber

import "github.com/microcosm-cc/bluemonday"

// defaultHTMLScrubber provides a basic implementation of the Scrubber interface for HTML content
type defaultHTMLScrubber struct{}

// Scrub scrubs HTML content by removing potentially dangerous tags and attributes
func (s *defaultHTMLScrubber) Scrub(htmlContent string) string {
	return bluemonday.UGCPolicy().Sanitize(htmlContent)
}

// defaultHTMLScrubberInstance is the singleton instance of defaultHTMLScrubber
var defaultHTMLScrubberInstance = &defaultHTMLScrubber{}

// DefaultHTMLScrubber returns an instance of defaultHTMLScrubber
func DefaultHTMLScrubber() Scrubber {
	return defaultHTMLScrubberInstance
}
