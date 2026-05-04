package scrubber

import "github.com/microcosm-cc/bluemonday"

// Option configures a bluemonday policy for use with NewPolicyScrubber
type Option func(*bluemonday.Policy)

// NewPolicyScrubber constructs a Scrubber backed by a bluemonday UGCPolicy configured via the provided options
func NewPolicyScrubber(opts ...Option) Scrubber {
	p := bluemonday.UGCPolicy()

	for _, opt := range opts {
		opt(p)
	}

	return ScrubberFunc(p.Sanitize)
}

// WithStyling allows class and inline style attributes on all elements
func WithStyling() Option {
	return func(p *bluemonday.Policy) {
		p.AllowStyling()
		p.AllowAttrs("style").Globally()
	}
}

// WithTables allows table-related elements and attributes
func WithTables() Option {
	return func(p *bluemonday.Policy) {
		p.AllowTables()
	}
}

// WithImages allows img elements with src, alt, width, height attributes
func WithImages() Option {
	return func(p *bluemonday.Policy) {
		p.AllowImages()
	}
}

// WithDocumentStructure allows html, head, body, meta, style, link, and center elements
// along with their common attributes
func WithDocumentStructure() Option {
	return func(p *bluemonday.Policy) {
		p.AllowElements("html", "head", "body", "meta", "style", "link", "center")
		p.AllowAttrs("lang", "dir", "xmlns").OnElements("html")
		p.AllowAttrs("content", "http-equiv", "name", "charset").OnElements("meta")
		p.AllowAttrs("type", "media").OnElements("style")
		p.AllowAttrs("href", "rel", "type", "media").OnElements("link")
	}
}

// WithEmailLayout allows table-layout attributes used by email templates and premailer output
func WithEmailLayout() Option {
	return func(p *bluemonday.Policy) {
		p.AllowAttrs("width", "height", "align", "valign", "bgcolor", "border",
			"cellpadding", "cellspacing", "colspan", "rowspan").Globally()
	}
}

// WithAccessibility allows ARIA and i18n attributes globally
func WithAccessibility() Option {
	return func(p *bluemonday.Policy) {
		p.AllowAttrs("role", "aria-label", "aria-hidden", "dir", "lang").Globally()
	}
}

// WithURLSchemes sets the allowed URL schemes for href and src attributes
func WithURLSchemes(schemes ...string) Option {
	return func(p *bluemonday.Policy) {
		p.AllowURLSchemes(schemes...)
	}
}

// WithNoRelativeURLs disallows relative URLs so that all href and src values
// must be fully qualified. Relative URLs cannot resolve inside email clients
// and may be exploited to link to unexpected destinations
func WithNoRelativeURLs() Option {
	return func(p *bluemonday.Policy) {
		p.AllowRelativeURLs(false)
	}
}

// WithTargetBlankOnLinks adds target="_blank" to fully qualified links so they
// open in a new window or tab when rendered in web-based email clients
func WithTargetBlankOnLinks() Option {
	return func(p *bluemonday.Policy) {
		p.AddTargetBlankToFullyQualifiedLinks(true)
	}
}

// WithEmailDefaults applies all options appropriate for sanitizing rendered email HTML:
// styling, tables, images, document structure, email layout, accessibility,
// common email URL schemes (http, https, mailto, cid, tel), and link security
// settings recommended by bluemonday for email contexts
func WithEmailDefaults() Option {
	return func(p *bluemonday.Policy) {
		opts := []Option{
			WithStyling(),
			WithTables(),
			WithImages(),
			WithDocumentStructure(),
			WithEmailLayout(),
			WithAccessibility(),
			WithURLSchemes("http", "https", "mailto", "cid", "tel"),
			WithNoRelativeURLs(),
			WithTargetBlankOnLinks(),
		}

		for _, opt := range opts {
			opt(p)
		}
	}
}
