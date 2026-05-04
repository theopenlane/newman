package mock

import (
	"context"
	"fmt"
	"hash/fnv"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/shared"
)

const (
	readWriteMode = 0755
)

// EmailSender is a mock email sender that captures sent messages for test assertion.
// It satisfies newman.EmailSender and additionally exposes Reset and Messages
// for inspecting what was sent
type EmailSender struct {
	logger   *slog.Logger
	mu       sync.Mutex
	messages []*newman.EmailMessage
	storage  string
}

// New creates a mock email sender. If storage is non-empty, sent emails are
// also written to disk as MIME files for manual inspection
func New(storage string) (*EmailSender, error) {
	return &EmailSender{
		logger:  slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		storage: storage,
	}, nil
}

// Reset clears captured messages. Tests that send emails should call Reset
// as part of cleanup so other tests start with a clean slate
func (s *EmailSender) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messages = nil
}

// Messages returns a snapshot of all captured messages
func (s *EmailSender) Messages() []*newman.EmailMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]*newman.EmailMessage, len(s.messages))
	copy(out, s.messages)

	return out
}

// SendEmail sends an email with the given message
func (s *EmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendBatchEmail sends a batch of emails with the given messages
func (s *EmailSender) SendBatchEmail(messages []*newman.EmailMessage) error {
	return s.SendBatchEmailWithContext(context.Background(), messages)
}

// SendBatchEmailWithContext validates and captures each message in the batch
func (s *EmailSender) SendBatchEmailWithContext(ctx context.Context, messages []*newman.EmailMessage) error {
	for _, message := range messages {
		if err := s.SendEmailWithContext(ctx, message); err != nil {
			return err
		}
	}

	return nil
}

// SendEmailWithContext validates and captures the message for later assertion
func (s *EmailSender) SendEmailWithContext(_ context.Context, message *newman.EmailMessage) error {
	if err := shared.ValidateEmailMessage(message); err != nil {
		return err
	}

	s.logger.Info("Sending test email",
		"to", strings.Join(message.To, ","),
		"subject", message.Subject,
		"message", message.Text,
		"html", message.HTML,
	)

	s.mu.Lock()
	s.messages = append(s.messages, message)
	s.mu.Unlock()

	if s.storage != "" {
		if err := s.saveEmailToFile(message); err != nil {
			return err
		}
	}

	return nil
}

// saveEmailToFile for manual inspection
func (s *EmailSender) saveEmailToFile(message *newman.EmailMessage) error {
	// we have already validated the message contains at least one recipient
	firstTo := message.GetTo()[0]
	dir := filepath.Join(s.storage, firstTo)

	if err := os.MkdirAll(dir, readWriteMode); err != nil {
		return err
	}

	mimeMsg, err := shared.BuildMimeMessage(message)
	if err != nil {
		return err
	}

	path, err := generateUniqueFilename(dir, mimeMsg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, mimeMsg, readWriteMode)
}

func generateUniqueFilename(dir string, message []byte) (string, error) {
	// Generate unique filename to avoid overwriting
	ts := time.Now().Format(time.RFC3339)
	h := fnv.New32()

	_, err := h.Write(message)
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, fmt.Sprintf("%s-%d.mim", ts, h.Sum32())), nil
}
