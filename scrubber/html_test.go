package scrubber

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHTMLScrubber(t *testing.T) {
	scrubber := DefaultHTMLScrubber()

	t.Run("remove potential XSS attack", func(t *testing.T) {
		input := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`
		expected := `<div>XSS</div>`
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("on methods not allowed", func(t *testing.T) {
		input := `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`
		expected := `<a href="http://www.google.com" rel="nofollow">Google</a>`
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("<p> can't have href", func(t *testing.T) {
		input := `<p href="http://www.google.com">Google</p>`
		expected := `<p>Google</p>`
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("allow safe HTML tags", func(t *testing.T) {
		input := `<b>Bold</b> <i>Italic</i> <u>Underline</u>`
		expected := `<b>Bold</b> <i>Italic</i> <u>Underline</u>`
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("scrub mixed content", func(t *testing.T) {
		input := `<div>Hello <script>alert("xss")</script> World</div>`
		expected := `<div>Hello  World</div>`
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("scrub documentation content", func(t *testing.T) {
		input := `<script>alert('xss')</script><b>Bold</b>`
		expected := `<b>Bold</b>`
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})
}
