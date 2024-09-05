package shared_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/theopenlane/newman/shared"
)

func TestStrPtr(t *testing.T) {
	str := "String to test for pointer"
	ptrStr := shared.StrPtr(str)
	assert.Equal(t, ptrStr, &str)
	assert.EqualValues(t, ptrStr, &str)
}

func TestGetMimeType(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"document.pdf", "application/pdf"},
		{"image.png", "image/png"},
		{"archive.zip", "application/zip"},
		{"unknownfile.unknown", ""},
		{"text.txt", "text/plain; charset=utf-8"},
		{"no_extension", ""},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			result := shared.GetMimeType(test.filename)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetMimeTypeEdgeCases(t *testing.T) {
	t.Run("unknown extension", func(t *testing.T) {
		filename := "file.unknownext"
		expected := ""
		result := shared.GetMimeType(filename)
		assert.Equal(t, expected, result)
	})

	t.Run("empty filename", func(t *testing.T) {
		filename := ""
		expected := ""
		result := shared.GetMimeType(filename)
		assert.Equal(t, expected, result)
	})
}

func TestIsHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"<html><body>Hello</body></html>", true},
		{"Just a plain text", false},
		{"<div>HTML content</div>", true},
		{"Plain text with <html> tag", true},
	}

	for _, test := range tests {
		result := shared.IsHTML(test.input)
		assert.Equal(t, result, test.expected)
	}
}

func ExampleStrPtr() {
	name := "Jerry Seinfeld"
	namePtr := shared.StrPtr(name)
	fmt.Println(*namePtr)

	// Output: Jerry Seinfeld
}

func ExampleGetMimeType() {
	filename := "document.pdf"
	mimeType := shared.GetMimeType(filename)
	fmt.Println(mimeType)

	// Output: application/pdf
}

func ExampleIsHTML_true() {
	html := "<html><body>HTML body</body></html>"
	result := shared.IsHTML(html)
	fmt.Println(result)

	// Output: true
}

func ExampleIsHTML_false() {
	plainText := "Just a plain text"
	result := shared.IsHTML(plainText)
	fmt.Println(result)

	// Output: false
}

func ExampleIsHTML_partiallyContainsHTML() {
	combined := "Plain text with <html> tag"
	result := shared.IsHTML(combined)
	fmt.Println(result)

	// Output: true
}
