package render

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// EmailContent is the root object rendered against a theme template. Templates
// reference its fields via dotted paths — {{ .Body.Title }}, {{ .Config.CompanyName }}
// style overrides are attached to the section that actually renders them (e.g. Body.Style, Body.Intros.Style)
type EmailContent struct {
	// Request is the operation input value, typed at the dispatch site
	// Templates reach its fields via reflection (e.g. .Request.FirstName, .Request.CampaignID)
	Request any
	// Config is the per-send config supplied by the dispatcher
	Config any
	// Body is the per-send structured body content
	Body ContentBody
}

// ContentBody is the per-send structured body content; templates reference only the slots they render
type ContentBody struct {
	// Preheader is hidden text shown as the email preview in email clients
	Preheader string
	// Header is the top header row rendered by modern themes
	Header HeaderBlock
	// Icon is an image displayed between the logo and content area
	Icon *ContentIcon
	// Name is the recipient name used in the greeting line
	Name string
	// Greeting is the greeting text preceding the recipient name
	Greeting string
	// Title overrides the greeting and name when set, used as the full heading line
	Title string
	// Intros is the introductory block rendered before the main content
	Intros IntrosBlock
	// ContentBlocks are standalone raw HTML blocks rendered after intros and before structured content
	ContentBlocks []template.HTML
	// Callout is a visually distinct block for bulleted or numbered lists with an optional heading
	Callout *Callout
	// Dictionary is a block of key-value pairs rendered as a definition list
	Dictionary Dictionary
	// Tables contains data tables rendered in the body
	Tables []DataTable
	// Actions contains call-to-action buttons and invite codes
	Actions []Action
	// Outros is the closing block rendered after the main content
	Outros OutrosBlock
	// Signature is the sign-off text (e.g. "Best regards")
	Signature string
	// SignatureName is the name appended under the signature line
	SignatureName string
	// FreeMarkdown is a markdown block that replaces all structured body content when set
	FreeMarkdown MarkdownContent
	// BodyWidth overrides the email body width (e.g. "600px")
	BodyWidth string
	// AdditionalCSS is raw CSS appended to the theme styles
	AdditionalCSS string
	// Style is per-Body style overrides
	Style Style
}

// HeaderBlock is the top header row rendered by modern themes; carries the
// compact brand logo and per-section style
type HeaderBlock struct {
	// Logo overrides Config.LogoURL in the compact header slot; nil falls back to Config.LogoURL
	Logo *ContentIcon
	// Style is per-HeaderBlock style overrides
	Style Style
}

// IntrosBlock is the introductory paragraph block. Templates pick one of
// Paragraphs, Markdown, or Unsafe depending on the intended escape path
type IntrosBlock struct {
	// Paragraphs is the list of plain-text paragraphs rendered in order
	Paragraphs []string
	// Markdown is a markdown intro block convertible to HTML via goldmark
	Markdown MarkdownContent
	// Unsafe contains raw HTML intro blocks trusted by the template author
	Unsafe []template.HTML
	// Style is per-Intros style overrides
	Style Style
}

// OutrosBlock is the closing paragraph block. Templates pick one of
// Paragraphs, Markdown, or Unsafe depending on the intended escape path
type OutrosBlock struct {
	// Paragraphs is the list of plain-text paragraphs rendered in order
	Paragraphs []string
	// Markdown is a markdown outro block convertible to HTML via goldmark
	Markdown MarkdownContent
	// Unsafe contains raw HTML outro blocks trusted by the template author
	Unsafe []template.HTML
	// Style is per-Outros style overrides
	Style Style
}

// Action is a call-to-action element in the body — a CTA button or an invite code
type Action struct {
	// Instructions is the descriptive text displayed above the action button or invite code
	Instructions string
	// Button holds the CTA button details
	Button Button
	// InviteCode is an invite or authentication code displayed instead of a button
	InviteCode string
	// Style is per-Action style overrides
	Style Style
}

// Button is a clickable CTA button rendered in an Action
type Button struct {
	// Text is the button label
	Text string
	// Link is the URL the button navigates to
	Link string
	// Color is the hex background color; convenience override of Style.ButtonColor
	Color string
	// TextColor is the hex label color; convenience override of Style.ButtonTextColor
	TextColor string
	// Style is per-Button style overrides
	Style Style
}

// Dictionary is a block of key-value cells rendered as a definition list
type Dictionary struct {
	// Cells is the list of key-value entries
	Cells []Cell
	// Style is per-Dictionary style overrides
	Style Style
}

// Callout is a visually distinct block for bulleted or numbered lists with an optional heading
type Callout struct {
	// Title is the heading displayed above the list items
	Title string
	// Items is the list of entries rendered without escaping; use render.Link / render.Bold for inline formatting
	Items []template.HTML
	// Ordered renders a numbered list when true; unordered (bulleted) when false
	Ordered bool
	// Style is per-Callout style overrides
	Style Style
}

// DataTable is a data table rendered in the body
type DataTable struct {
	// Title is the heading displayed above the table
	Title string
	// TitleUnsafe is a raw HTML title that overrides Title when non-empty
	TitleUnsafe template.HTML
	// Data is a 2D array of cells representing the table rows and columns
	Data [][]Cell
	// Columns holds column width and alignment configuration
	Columns TableColumns
	// Class is the CSS class applied to the table wrapper element
	Class string
	// Footer is the text displayed below the table
	Footer string
	// FooterUnsafe is a raw HTML footer that overrides Footer when non-empty
	FooterUnsafe template.HTML
	// Style is per-DataTable style overrides
	Style Style
}

// TableColumns holds column-level display configuration for a DataTable
type TableColumns struct {
	// CustomWidth maps column keys to width values (e.g. "Price": "15%")
	CustomWidth map[string]string
	// CustomAlignment maps column keys to alignment values (e.g. "Price": "right")
	CustomAlignment map[string]string
}

// Cell is a single key-value entry used in dictionaries and table cells
type Cell struct {
	// Key is the column or property name
	Key string
	// Value is the text value, HTML-escaped during rendering
	Value string
	// UnsafeValue is a raw HTML value used when Value is empty
	UnsafeValue template.HTML
}

// ContentIcon is an icon image displayed between the logo and content area
type ContentIcon struct {
	// Src is the URL of the icon image
	Src string
	// Alt is the alt text for the icon image
	Alt string
}

// Style holds color and font overrides used throughout the email. Every section
// embeds a Style so template authors can override look per section; zero-value
// fields fall through to whatever fallback the theme template elects
type Style struct {
	// PrimaryColor is the headline/emphasis color
	PrimaryColor string
	// SecondaryColor is the secondary accent color
	SecondaryColor string
	// BackgroundColor is the outer page background color
	BackgroundColor string
	// TextColor is the body text color
	TextColor string
	// ButtonColor is the call-to-action button background color
	ButtonColor string
	// ButtonTextColor is the call-to-action button text color
	ButtonTextColor string
	// LinkColor is the anchor link color
	LinkColor string
	// FontFamily is the CSS font-family stack
	FontFamily string
}

// MarkdownContent is markdown source text convertible to HTML via goldmark with
// GFM and table extensions. The ToHTML method is safe to call on the zero value
type MarkdownContent template.HTML

// ToHTML converts the markdown content to HTML using goldmark with GFM and table extensions
func (m MarkdownContent) ToHTML() template.HTML {
	if m == "" {
		return ""
	}

	md := goldmark.New(goldmark.WithExtensions(
		extension.NewTable(),
		extension.GFM,
	))

	var buf bytes.Buffer

	if err := md.Convert([]byte(m), &buf); err != nil {
		return template.HTML(m) //nolint:gosec // fallback: caller-supplied markdown returned as-is
	}

	return template.HTML(buf.String()) //nolint:gosec // goldmark output is trusted
}
