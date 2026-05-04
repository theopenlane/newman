package scrubber

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPolicyScrubber_EmailDefaults(t *testing.T) {
	s := NewPolicyScrubber(WithEmailDefaults())

	t.Run("remove potential XSS attack", func(t *testing.T) {
		input := `<div><a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a></div>`
		expected := `<div>XSS</div>`
		result := s.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("on methods not allowed", func(t *testing.T) {
		input := `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`
		expected := `<a href="http://www.google.com" rel="nofollow noopener" target="_blank">Google</a>`
		result := s.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("strip relative URLs", func(t *testing.T) {
		input := `<a href="page.html">Relative</a>`
		result := s.Scrub(input)
		assert.NotContains(t, result, "page.html")
		assert.Contains(t, result, "Relative")
	})

	t.Run("fully qualified links get target blank", func(t *testing.T) {
		input := `<a href="https://example.com">Example</a>`
		result := s.Scrub(input)
		assert.Contains(t, result, `target="_blank"`)
		assert.Contains(t, result, "nofollow")
		assert.Contains(t, result, "noopener")
	})

	t.Run("p cannot have href", func(t *testing.T) {
		input := `<p href="http://www.google.com">Google</p>`
		expected := `<p>Google</p>`
		result := s.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("allow safe HTML tags", func(t *testing.T) {
		input := `<b>Bold</b> <i>Italic</i> <u>Underline</u>`
		expected := `<b>Bold</b> <i>Italic</i> <u>Underline</u>`
		result := s.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("scrub mixed content", func(t *testing.T) {
		input := `<div>Hello <script>alert("xss")</script> World</div>`
		expected := `<div>Hello  World</div>`
		result := s.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("preserve email table layout", func(t *testing.T) {
		input := `<table width="600"><tr><td style="padding:16px">content</td></tr></table>`
		result := s.Scrub(input)
		assert.Contains(t, result, "<table")
		assert.Contains(t, result, "content")
	})

	t.Run("preserve style attributes", func(t *testing.T) {
		input := `<p style="color:#333;font-size:14px">styled text</p>`
		result := s.Scrub(input)
		assert.Contains(t, result, "style=")
		assert.Contains(t, result, "styled text")
	})
}

func TestNewPolicyScrubber_IndividualOptions(t *testing.T) {
	t.Run("styling only", func(t *testing.T) {
		s := NewPolicyScrubber(WithStyling())
		input := `<p style="color:red">text</p>`
		result := s.Scrub(input)
		assert.Contains(t, result, "style=")
	})

	t.Run("tables only", func(t *testing.T) {
		s := NewPolicyScrubber(WithTables())
		input := `<table><tr><td>cell</td></tr></table>`
		result := s.Scrub(input)
		assert.Contains(t, result, "<table>")
	})

	t.Run("images only", func(t *testing.T) {
		s := NewPolicyScrubber(WithImages())
		input := `<img src="https://example.com/logo.png" alt="logo">`
		result := s.Scrub(input)
		assert.Contains(t, result, "<img")
	})

	t.Run("no relative URLs", func(t *testing.T) {
		s := NewPolicyScrubber(WithNoRelativeURLs())
		input := `<a href="../secret.html">Relative</a> <a href="https://example.com">Absolute</a>`
		result := s.Scrub(input)
		assert.NotContains(t, result, "secret.html")
		assert.Contains(t, result, "https://example.com")
	})

	t.Run("target blank on links", func(t *testing.T) {
		s := NewPolicyScrubber(WithTargetBlankOnLinks())
		input := `<a href="https://example.com">Example</a>`
		result := s.Scrub(input)
		assert.Contains(t, result, `target="_blank"`)
	})
}
