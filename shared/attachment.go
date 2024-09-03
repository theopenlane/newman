package shared

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"

	"github.com/theopenlane/newman/scrubber"
)

// Attachment represents an email attachment with its filename and content
type Attachment struct {
	// filename is the name of the attachment file
	Filename string
	// content is the binary content of the attachment
	Content []byte
	// contentType is the MIME type of the file
	ContentType string
	// Filepath is the path to the attachment file
	FilePath string
}

// NewAttachment creates a new Attachment instance with the specified filename and content
func NewAttachment(filename string, content []byte) *Attachment {
	return &Attachment{
		Filename: filename,
		Content:  content,
	}
}

// NewAttachmentFromFile creates a new Attachment instance from the specified file path
func NewAttachmentFromFile(filePath string) (*Attachment, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	filename := extractFilename(filePath)

	return NewAttachment(
		filename,
		content,
	), nil
}

// extractFilename extracts the filename from the file path
func extractFilename(filePath string) string {
	parts := strings.Split(filePath, "/")
	return parts[len(parts)-1]
}

// SetFilename sets the filename of the attachment
func (a *Attachment) SetFilename(filename string) {
	a.Filename = filename
}

// GetFilename safely returns the filename of the attachment
func (a *Attachment) GetFilename() string {
	if a == nil {
		return "nil_attachment"
	}

	return scrubber.DefaultTextScrubber().Scrub(a.Filename)
}

// GetBase64StringContent returns the content of the attachment as a base64-encoded string
func (a *Attachment) GetBase64StringContent() string {
	if a == nil {
		return ""
	}

	return string(a.GetBase64Content())
}

// SetContent sets the content of the attachment
func (a *Attachment) SetContent(content []byte) {
	a.Content = content
}

// GetBase64Content returns the content of the attachment as a base64-encoded byte slice
func (a *Attachment) GetBase64Content() []byte {
	if a == nil || len(a.Content) == 0 {
		return []byte{}
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(a.Content)))

	base64.StdEncoding.Encode(buf, a.Content)

	return buf
}

// GetRawContent returns the content of the attachment as its raw byte slice
func (a *Attachment) GetRawContent() []byte {
	if a == nil || len(a.Content) == 0 {
		return []byte{}
	}

	return a.Content
}

// jsonAttachment represents the JSON structure for an email attachment
type jsonAttachment struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

// MarshalJSON custom marshaler for Attachment
func (a Attachment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonAttachment{
		Filename: a.Filename,
		Content:  base64.StdEncoding.EncodeToString(a.Content),
	})
}

// UnmarshalJSON custom unmarshaler for Attachment
func (a *Attachment) UnmarshalJSON(data []byte) error {
	aux := &jsonAttachment{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	a.Filename = aux.Filename

	content, err := base64.StdEncoding.DecodeString(aux.Content)
	if err != nil {
		return err
	}

	a.Content = content

	return nil
}
