package shared

import (
	"mime"
	"path/filepath"
	"regexp"
	"strings"
)

// IsHTML checks if a string contains HTML tags
func IsHTML(str string) bool {
	htmlRegex := regexp.MustCompile(`(?i)<\/?[a-z][\s\S]*>`)
	return htmlRegex.MatchString(str)
}

// GetMimeType returns the MIME type based on the file extension
func GetMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	return mime.TypeByExtension(ext)
}

// StrPtr takes a string value and returns a pointer to that string
func StrPtr(str string) *string {
	return &str
}
