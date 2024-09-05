package scrubber

// nonScrubber is an implementation of the Scrubber interface that performs no sanitization
type nonScrubber struct{}

// Scrub returns the input text without any modifications
func (s *nonScrubber) Scrub(text string) string {
	return text
}

var nonScrubberInstance = &nonScrubber{}

// NonScrubber returns an instance of nonScrubber
func NonScrubber() Scrubber {
	return nonScrubberInstance
}
