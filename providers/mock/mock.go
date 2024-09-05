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

type mockEmailSender struct {
	logger   *slog.Logger
	emailLog *EmailLog
	storage  string
}

// EmailLog combines email data with a mutex to ensure that tests can safely access the mock concurrently for that email blastertons
type EmailLog struct {
	sync.Mutex
	Data [][]byte
}

// Emails contains all emails sent by the mock which tests can use to verify which emails were sent
var Emails EmailLog

func New(storage string) (newman.EmailSender, error) {
	return &mockEmailSender{
		logger:   slog.New(slog.NewJSONHandler(os.Stdout, nil)),
		emailLog: new(EmailLog),
		storage:  storage,
	}, nil
}

// ResetEmailMock tests that send emails should call ResetEmailMock as part of their cleanup to ensure that other tests can depend on the state of the mock
func (s *mockEmailSender) ResetEmailMock() {
	Emails.Lock()
	defer Emails.Unlock()
	Emails.Data = nil
}

func (s *mockEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

func (s *mockEmailSender) SendEmailWithContext(ctx context.Context, message *newman.EmailMessage) error {
	if err := shared.ValidateEmailMessage(message); err != nil {
		return err
	}

	s.logger.Info("Sending test email",
		"to", strings.Join(message.To, ","),
		"subject", message.Subject,
		"message", message.Text,
		"html", message.HTML,
	)

	s.emailLog.Lock()
	s.emailLog.Data = append(s.emailLog.Data, []byte(message.Text))
	s.emailLog.Unlock()

	if s.storage != "" {
		if err := s.saveEmailToFile(message); err != nil {
			return err
		}
	}

	return nil
}

// saveEmailToFile for manual inspection
func (s *mockEmailSender) saveEmailToFile(message *newman.EmailMessage) error {
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

	path := generateUniqueFilename(dir, mimeMsg)

	return os.WriteFile(path, mimeMsg, readWriteMode)
}

func generateUniqueFilename(dir string, message []byte) string {
	// Generate unique filename to avoid overwriting
	ts := time.Now().Format(time.RFC3339)
	h := fnv.New32()
	h.Write(message)

	return filepath.Join(dir, fmt.Sprintf("%s-%d.mim", ts, h.Sum32()))
}
