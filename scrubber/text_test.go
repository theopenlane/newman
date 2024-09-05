package scrubber

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTextScrubber(t *testing.T) {
	scrubber := DefaultTextScrubber()

	t.Run("escape special characters", func(t *testing.T) {
		input := `Hello <world> & "everyone"`
		expected := `Hello &lt;world&gt; &amp; &#34;everyone&#34;`
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""
		expected := ""
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("string with spaces", func(t *testing.T) {
		input := "   "
		expected := ""
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("string with no special characters", func(t *testing.T) {
		input := "plain text"
		expected := "plain text"
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})

	t.Run("scrub documentation content", func(t *testing.T) {
		input := " <script>alert('xss')</script> "
		expected := `&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;`
		result := scrubber.Scrub(input)
		assert.Equal(t, expected, result)
	})
}
