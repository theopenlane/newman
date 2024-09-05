package scrubber

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNonScrubber(t *testing.T) {
	scrubber := NonScrubber()

	t.Run("return input as is", func(t *testing.T) {
		input := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`
		expected := input
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""
		expected := ""
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("plain text", func(t *testing.T) {
		input := "plain text"
		expected := "plain text"
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("scrub documentation content", func(t *testing.T) {
		input := "<script>alert('xss')</script>"
		expected := "<script>alert('xss')</script>"
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})
}
