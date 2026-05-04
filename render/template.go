package render

import (
	"bytes"
	"fmt"
	texttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/vanng822/go-premailer/premailer"
)

// ExecuteTextTemplate parses and executes a text/template string with sprig
// template functions against the provided data map
func ExecuteTextTemplate(name, tmplStr string, data map[string]any) (string, error) {
	if tmplStr == "" {
		return "", nil
	}

	tmpl, err := texttemplate.New(name).Funcs(texttemplate.FuncMap(sprig.TxtFuncMap())).Parse(tmplStr)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateParsingFailed, err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("%w: %w", ErrTemplateExecutionFailed, err)
	}

	return buf.String(), nil
}

// InlineCSS transforms CSS style blocks into inline style attributes for
// email client compatibility using premailer
func InlineCSS(html string) (string, error) {
	if html == "" {
		return "", nil
	}

	prem, err := premailer.NewPremailerFromString(html, premailer.NewOptions())
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCSSInliningFailed, err)
	}

	result, err := prem.Transform()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrCSSInliningFailed, err)
	}

	return result, nil
}

// HTMLToPlainText converts HTML content to a plain text representation with
// pretty-printed tables using go-premailer's built-in text conversion
func HTMLToPlainText(html string) (string, error) {
	if html == "" {
		return "", nil
	}

	prem, err := premailer.NewPremailerFromString(html, premailer.NewOptions())
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrPlainTextConversionFailed, err)
	}

	plain, err := prem.TransformText()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrPlainTextConversionFailed, err)
	}

	return plain, nil
}
