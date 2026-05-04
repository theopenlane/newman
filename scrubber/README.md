# Overview

The scrubber package defines an interface for content sanitization and provides composable policy construction via functional options. The default scrubber is a no-op; callers opt in to sanitization explicitly, typically at email provider initialization time.

# Usage

```go
import "github.com/theopenlane/newman/scrubber"

// No-op (default) — passes content through unmodified
s := scrubber.DefaultHTMLScrubber()

// Email-safe policy — preserves styles, tables, images, and email layout
// while stripping scripts and event handlers
s = scrubber.NewPolicyScrubber(scrubber.WithEmailDefaults())

// Compose individual options for a custom policy
s = scrubber.NewPolicyScrubber(
    scrubber.WithStyling(),
    scrubber.WithTables(),
    scrubber.WithURLSchemes("http", "https", "mailto"),
)

// Inline function adapter
s = scrubber.ScrubberFunc(func(content string) string {
    return strings.TrimSpace(content)
})
```

## Provider-level scrubbing

Scrubbing is configured at the provider level rather than per-message. Pass a scrubber when constructing the email provider:

```go
import (
    "github.com/theopenlane/newman/providers/resend"
    "github.com/theopenlane/newman/scrubber"
)

sender, err := resend.New(apiKey,
    resend.WithHTMLScrubber(scrubber.NewPolicyScrubber(scrubber.WithEmailDefaults())),
)
```
