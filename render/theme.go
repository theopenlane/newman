package render

// Theme holds the HTML and plain text template strings for an email theme.
// Theme authors reference EmailContent fields directly in the template via
// dotted paths (e.g. {{ .Body.Title }}, {{ .Config.LogoURL }}); any section-
// level Style override precedence is expressed inline in the template
type Theme struct {
	// Name is the human-readable identifier for the theme
	Name string
	// HTML is the raw Go html/template string for the theme
	HTML string
	// Text is the raw Go text/template string for the theme
	Text string
}
