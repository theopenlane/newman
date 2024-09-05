# Overview

The scrubber package defines an interface for content sanitization (e.g. ensuring the content provided is user generated content only) and creates sane defaults for sanitizing plain text and HTML content. The goal is to allow customizable content sanitization for newman's delivery.

# Usage

To use the scrubber package, you can either use the provided default implementations or create your own custom implementations of the scrubber interface.

Example:

```go
	import (
		"html"
		"strings"

		"github.com/theopenlane/newman/scrubber"
		"github.com/theopenlane/newman"
	)

	func main() {
		email := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Subject", "<p>HTML content</p>")

		customTextScrubber := scrubber.ScrubberFunc(func(content string) string {
			//Implement your custom scrubber logic
			return strings.ToLower(strings.TrimSpace(content))
		})

		customHtmlScrubber := scrubber.ScrubberFunc(func(content string) string {
			//Implement your custom scrubber logic
			return html.EscapeString(content)
		})

		email.SetCustomTextScrubber(customTextScrubber)
		email.SetCustomHTMLScrubber(customHtmlScrubber)
	}
```