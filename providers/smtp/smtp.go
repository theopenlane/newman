package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"github.com/theopenlane/newman"
)

const (
	defaultConnectionMethod = "IMPLICIT"
	TLSConnection           = "TLS"
	CRAMMD5Auth             = "CRAM-MD5"
)

// smtpEmailSender is responsible for sending emails using SMTP
type smtpEmailSender struct {
	// The SMTP server host
	host string
	// The SMTP server port
	port int
	// The username for authentication
	user string
	// The password for authentication
	password string
	// The authentication method to use
	authMethod string
	// The connection method to use (by default implicit)
	connectionMethod string
}

// New creates a new instance of smtpEmailSender
func New(host string, port int, user, password string, authMethod string) (*smtpEmailSender, error) {
	return NewWithConnMethod(host, port, user, password, authMethod, defaultConnectionMethod)
}

// NewWithConnMethod creates a new instance of smtpEmailSender with the specified connection method
func NewWithConnMethod(host string, port int, user, password string, authMethod string, connectionMethod string) (*smtpEmailSender, error) {
	return &smtpEmailSender{
		host,
		port,
		user,
		password,
		authMethod,
		connectionMethod,
	}, nil
}

// SendEmail satisfies the EmailSender interface
func (s *smtpEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendEmailWithContext satisfies the EmailSender interface
func (s *smtpEmailSender) SendEmailWithContext(ctx context.Context, message *newman.EmailMessage) error {
	sendMailTo := message.GetTo()
	sendMailTo = append(sendMailTo, message.GetCC()...)
	sendMailTo = append(sendMailTo, message.GetBCC()...)

	msg, err := newman.BuildMimeMessage(message)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	if s.authMethod == CRAMMD5Auth {
		auth = smtp.CRAMMD5Auth(s.user, s.password)
	}

	if s.connectionMethod == TLSConnection {
		return s.secureSend(auth, message.GetFrom(), sendMailTo, msg)
	}

	return s.send(auth, message.GetFrom(), sendMailTo, msg)
}

func (s *smtpEmailSender) send(auth smtp.Auth, from string, to []string, message []byte) error {
	return smtp.SendMail(fmt.Sprintf("%s:%d", s.host, s.port), auth, from, to, message)
}

func (s *smtpEmailSender) secureSend(auth smtp.Auth, from string, to []string, message []byte) error {
	// TODO: this should be updated to not use environment variables
	skipInsecure := false
	if os.Getenv("APP_ENV") == "development" {
		skipInsecure = true
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: skipInsecure, // nolint: gosec
		ServerName:         s.host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return err
	}

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	if _, err = w.Write(message); err != nil {
		return err
	}

	if err = w.Close(); err != nil {
		return err
	}

	return client.Quit()
}
