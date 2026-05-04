package render

import (
	"bytes"
	"fmt"
	"html/template"
	texttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

// Renderer generates themed HTML and plain text email output from an EmailContent value
type Renderer struct {
	// theme provides the HTML and text template strings
	theme *Theme
	// inlineCSS enables premailer CSS inlining when true
	inlineCSS bool
}

// RendererOption is a function that configures a Renderer
type RendererOption func(*Renderer)

// WithTheme sets the theme used for email rendering
func WithTheme(t *Theme) RendererOption {
	return func(r *Renderer) {
		r.theme = t
	}
}

// NewRenderer creates a Renderer with the given options applied over sensible defaults
func NewRenderer(opts ...RendererOption) *Renderer {
	r := &Renderer{
		inlineCSS: true,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// GenerateHTML renders the EmailContent into a complete HTML string using the configured theme.
// The content value is passed directly as the template root context, so template authors
// reference fields via flat dotted paths (e.g. {{ .Body.Title }}, {{ .Config.LogoURL }})
func (r *Renderer) GenerateHTML(content EmailContent) (string, error) {
	if r.theme == nil {
		return "", ErrThemeRequired
	}

	tmpl, err := templateBase().Parse(r.theme.HTML)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateParsingFailed, err)
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, content); err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateExecutionFailed, err)
	}

	result := buf.String()

	if !r.inlineCSS {
		return result, nil
	}

	return InlineCSS(result)
}

// GeneratePlainText renders the EmailContent into a plain text string using the configured theme
func (r *Renderer) GeneratePlainText(content EmailContent) (string, error) {
	if r.theme == nil {
		return "", ErrThemeRequired
	}

	tmpl, err := textTemplateBase().Parse(r.theme.Text)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateParsingFailed, err)
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, content); err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateExecutionFailed, err)
	}

	return HTMLToPlainText(buf.String())
}

// Bold wraps text in a <strong> tag with HTML-escaped content
func Bold(s string) template.HTML {
	return template.HTML("<strong>" + template.HTMLEscapeString(s) + "</strong>")
}

// Link renders an anchor tag with HTML-escaped href and text using the default link color
func Link(href, text string) template.HTML {
	return LinkWithColor(href, text, "rgb(63,118,255)")
}

// LinkWithColor renders an anchor tag with HTML-escaped href, text, and a custom inline color
func LinkWithColor(href, text, color string) template.HTML {
	return template.HTML(`<a href="` + template.HTMLEscapeString(href) +
		`" style="color:` + template.HTMLEscapeString(color) + `;text-decoration-line:none" target="_blank">` +
		template.HTMLEscapeString(text) + `</a>`)
}

// templateBase returns a base template pre-loaded with sprig functions and custom
// helpers for URL, CSS, and raw HTML handling in email templates
func templateBase() *template.Template {
	return template.New("").Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"url": func(s string) template.URL {
			return template.URL(s)
		},
		"css": func(in any) template.CSS {
			s, ok := in.(string)
			if !ok {
				return ""
			}

			return template.CSS(s)
		},
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
		"bold": func(s string) template.HTML { return Bold(s) },
		"link": func(href, text string) template.HTML { return Link(href, text) },
	})
}

// textTemplateBase returns a base text/template pre-loaded with sprig functions
// and plain-text equivalents of the HTML helpers, for use in plain text rendering
func textTemplateBase() *texttemplate.Template {
	return texttemplate.New("").Funcs(sprig.TxtFuncMap()).Funcs(texttemplate.FuncMap{
		"url":  func(s string) string { return s },
		"bold": func(s string) string { return s },
		"link": func(href, text string) string { return text + " (" + href + ")" },
		"css": func(in any) string {
			s, ok := in.(string)
			if !ok {
				return ""
			}

			return s
		},
		"safe": func(s string) string { return s },
	})
}
