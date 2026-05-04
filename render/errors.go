package render

import "errors"

var (
	// ErrTemplateParsingFailed is returned when a Go template fails to parse
	ErrTemplateParsingFailed = errors.New("template parsing failed")

	// ErrTemplateExecutionFailed is returned when a Go template fails to execute
	ErrTemplateExecutionFailed = errors.New("template execution failed")

	// ErrCSSInliningFailed is returned when the premailer CSS inlining transform fails
	ErrCSSInliningFailed = errors.New("css inlining failed")

	// ErrPlainTextConversionFailed is returned when html-to-plaintext conversion fails
	ErrPlainTextConversionFailed = errors.New("plain text conversion failed")

	// ErrThemeRequired is returned when a Renderer is created or invoked with a nil theme
	ErrThemeRequired = errors.New("theme is required")
)
