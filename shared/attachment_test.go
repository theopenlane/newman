package shared

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAttachmentGetters(t *testing.T) {
	t.Run("GetFilename", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt"}
		assert.Equal(t, "test.txt", attachment.GetFilename())
		assert.Equal(t, "nil_attachment", (*Attachment)(nil).GetFilename())
	})

	t.Run("GetBase64Content", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt", Content: []byte("hello")}
		expected := []byte(base64.StdEncoding.EncodeToString([]byte("hello")))
		assert.Equal(t, expected, attachment.GetBase64Content())
		assert.Equal(t, []byte{}, (*Attachment)(nil).GetBase64Content())
	})

	t.Run("GetRawContent", func(t *testing.T) {
		attachment := Attachment{Filename: "test.txt", Content: []byte("hello")}
		expected := []byte("hello")
		assert.Equal(t, expected, attachment.GetRawContent())
		assert.Equal(t, []byte{}, (*Attachment)(nil).GetRawContent())
	})
}

func TestNewAttachment(t *testing.T) {
	filename := "testfile.txt"
	content := []byte("This is a test file content.")
	attachment := NewAttachment(filename, content)

	assert.Equal(t, attachment.Filename, filename)
	assert.Equal(t, string(attachment.Content), string(attachment.Content))
}

func TestAttachmentEdgeCases(t *testing.T) {
	t.Run("GetBase64Content with nil content", func(t *testing.T) {
		attachment := &Attachment{}
		assert.Equal(t, []byte{}, attachment.GetBase64Content())
	})

	t.Run("GetBase64StringContent with nil content", func(t *testing.T) {
		attachment := &Attachment{}
		assert.Equal(t, "", attachment.GetBase64StringContent())
	})

	t.Run("SetContent with nil content", func(t *testing.T) {
		attachment := &Attachment{}
		attachment.SetContent(nil)
		assert.Nil(t, attachment.Content)
	})

	t.Run("SetFilename with empty string", func(t *testing.T) {
		attachment := &Attachment{}
		attachment.SetFilename("")
		assert.Equal(t, "", attachment.Filename)
	})
}

func TestGetBase64Content(t *testing.T) {
	tests := []struct {
		attachment *Attachment
		expected   []byte
	}{
		{&Attachment{Filename: "test.txt", Content: []byte("hello")}, []byte(base64.StdEncoding.EncodeToString([]byte("hello")))},
		{&Attachment{Filename: "test.txt", Content: []byte("")}, []byte{}},
		{&Attachment{Filename: "empty.txt", Content: nil}, []byte{}},
	}

	for _, test := range tests {
		t.Run(test.attachment.GetFilename(), func(t *testing.T) {
			result := test.attachment.GetBase64Content()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetBase64StringContent(t *testing.T) {
	tests := []struct {
		attachment *Attachment
		expected   string
	}{
		{&Attachment{Filename: "test.txt", Content: []byte("hello")}, base64.StdEncoding.EncodeToString([]byte("hello"))},
		{&Attachment{Filename: "test.txt", Content: []byte("")}, ""},
		{nil, ""},
	}

	for _, test := range tests {
		t.Run(test.attachment.GetFilename(), func(t *testing.T) {
			result := test.attachment.GetBase64StringContent()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestGetRawContent(t *testing.T) {
	tests := []struct {
		attachment *Attachment
		expected   []byte
	}{
		{&Attachment{Filename: "test.txt", Content: []byte("hello")}, []byte("hello")},
		{&Attachment{Filename: "test.txt", Content: []byte("")}, []byte{}},
		{&Attachment{Filename: "empty.txt", Content: nil}, []byte{}},
	}

	for _, test := range tests {
		t.Run(test.attachment.GetFilename(), func(t *testing.T) {
			result := test.attachment.GetRawContent()
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestExtractFilename(t *testing.T) {
	t.Run("extract filename from valid path", func(t *testing.T) {
		filePath := "/path/to/file/document.pdf"
		expected := "document.pdf"
		result := extractFilename(filePath)
		assert.Equal(t, expected, result)
	})

	t.Run("extract filename from empty path", func(t *testing.T) {
		filePath := ""
		expected := ""
		result := extractFilename(filePath)
		assert.Equal(t, expected, result)
	})

	t.Run("extract filename from path with trailing slash", func(t *testing.T) {
		filePath := "/path/to/directory/"
		expected := ""
		result := extractFilename(filePath)
		assert.Equal(t, expected, result)
	})
}

func TestNewAttachmentFromFile(t *testing.T) {
	const testDataPath = "../testdata"

	t.Run("files in test data", func(t *testing.T) {
		testFiles := []struct {
			filePath        string
			expectedName    string
			expectedContent string
		}{
			{
				filepath.Join(testDataPath, "testfile.txt"),
				"testfile.txt",
				`openlane
https://theopenlane.io/

we do the cat typing

openlane, your trusted partner for cat typing`,
			}, {
				filepath.Join(testDataPath, "testfile.md"),
				"testfile.md",
				`# openlane
**[theopenlane.io](https://theopenlane.io/)**

## we do the cat typing

openlane, your trusted partner for cat typing`,
			},
		}

		for _, testFile := range testFiles {
			attachment, err := NewAttachmentFromFile(testFile.filePath)
			if err != nil {
				t.Fatalf("NewAttachmentFromFile() error = %v, want nil", err)
			}

			assert.Equal(t, attachment.Filename, testFile.expectedName)

			content, err := os.ReadFile(testFile.filePath)
			if err != nil {
				t.Fatalf("Failed to read test file: %v", err)
			}
			assert.Equal(t, string(attachment.Content), string(content))
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		filePath := "nonexistentfile.txt"
		attachment, err := NewAttachmentFromFile(filePath)
		assert.NotNil(t, err)
		assert.Nil(t, attachment)
	})
}

func TestSetFilename(t *testing.T) {
	attachment := &Attachment{}
	t.Run("SetFilename", func(t *testing.T) {
		expected := "newfile.txt"
		attachment.SetFilename(expected)
		assert.Equal(t, expected, attachment.Filename)
	})
}

func TestSetContent(t *testing.T) {
	attachment := &Attachment{}
	t.Run("SetContent", func(t *testing.T) {
		expected := []byte("new content")
		attachment.SetContent(expected)
		assert.Equal(t, expected, attachment.Content)
	})
}

func TestSanitizeFilename(t *testing.T) {
	attachment := &Attachment{}
	t.Run("sanitize Filename with HTML", func(t *testing.T) {
		fileName := "<div>Test</div>"
		expected := "&lt;div&gt;Test&lt;/div&gt;"
		attachment.SetFilename(fileName)
		assert.Equal(t, expected, attachment.GetFilename())
	})

	t.Run("sanitize Filename with spaces", func(t *testing.T) {
		fileName := "  Test  "
		expected := "Test"
		attachment.SetFilename(fileName)
		assert.Equal(t, expected, attachment.GetFilename())
	})
}
