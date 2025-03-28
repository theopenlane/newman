package scrubber

// Scrubber defines a method for sanitizing content
type Scrubber interface {
	// Scrub scrubs the provided content
	Scrub(content string) string
}

// ScrubberFunc is an adapter that allows the use of functions as Scrubbers
type ScrubberFunc func(message string) string //nolint:revive

// Scrub calls the function f with the given message
func (f ScrubberFunc) Scrub(message string) string {
	return f(message)
}
