package scrubber_test

import (
	"fmt"
	"html"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/theopenlane/newman/scrubber"
)

func trimSpaces(message string) string {
	return strings.TrimSpace(message)
}

func toUpper(message string) string {
	return strings.ToUpper(message)
}

func toLower(message string) string {
	return strings.ToLower(message)
}

func escapeString(message string) string {
	return html.EscapeString(message)
}

func TestScrubberFunc_Scrub(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		scrubFunc scrubber.ScrubberFunc
	}{
		{
			name:      "TrimSpaces",
			input:     "  some text  ",
			want:      "some text",
			scrubFunc: scrubber.ScrubberFunc(trimSpaces),
		},
		{
			name:      "ToUpper",
			input:     "some text",
			want:      "SOME TEXT",
			scrubFunc: scrubber.ScrubberFunc(toUpper),
		},
		{
			name:      "ToLower",
			input:     "SOME TEXT",
			want:      "some text",
			scrubFunc: scrubber.ScrubberFunc(toLower),
		},
		{
			name:  "EmptyString",
			input: "",
			want:  "",
			scrubFunc: scrubber.ScrubberFunc(func(message string) string {
				return message
			}),
		},
		{
			name:      "WhitespaceOnly",
			input:     "     ",
			want:      "",
			scrubFunc: scrubber.ScrubberFunc(trimSpaces),
		},
		{
			name:  "NoSanitization",
			input: "no change",
			want:  "no change",
			scrubFunc: scrubber.ScrubberFunc(func(message string) string {
				return message
			}),
		},
		{
			name:      "SpecialCharacters",
			input:     "<script>alert('xss')</script>",
			want:      "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
			scrubFunc: scrubber.ScrubberFunc(escapeString),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.scrubFunc.Scrub(tt.input))
		})
	}
}

func ExampleScrubberFunc() {
	scrubFunc := scrubber.ScrubberFunc(func(message string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(message)), " ", "_")
	})
	scrubdMessage := scrubFunc.Scrub("  some text  ")
	fmt.Println(scrubdMessage)
}
