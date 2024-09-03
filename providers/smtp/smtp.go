package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"github.com/theopenlane/newman"
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
	// The connection method to use (by default implicit
	connectionMethod string
}

// NewSMTPEmailSender creates a new instance of smtpEmailSender
func NewSMTPEmailSender(host string, port int, user, password string, authMethod string) (*smtpEmailSender, error) {
	return NewSMTPEmailSenderWithConnMethod(host, port, user, password, authMethod, "IMPLICIT")
}

// NewSMTPEmailSenderWithConnMethod creates a new instance of smtpEmailSender
func NewSMTPEmailSenderWithConnMethod(host string, port int, user, password string, authMethod string, connectionMethod string) (*smtpEmailSender, error) {
	return &smtpEmailSender{
		host,
		port,
		user,
		password,
		authMethod,
		connectionMethod,
	}, nil
}

// SendEmail sends an email using the specified SMTP settings and authentication method
func (s *smtpEmailSender) SendEmail(message *newman.EmailMessage) error {
	sendMailTo := message.GetTo()
	sendMailTo = append(sendMailTo, message.GetCC()...)
	sendMailTo = append(sendMailTo, message.GetBCC()...)

	msg, err := newman.BuildMimeMessage(message)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	if s.authMethod == "CRAM-MD5" {
		auth = smtp.CRAMMD5Auth(s.user, s.password)
	}

	if s.connectionMethod == "TLS" {
		skipInsecure := false
		if os.Getenv("APP_ENV") == "development" {
			skipInsecure = true
		}

		tlsconfig := &tls.Config{
			InsecureSkipVerify: skipInsecure, // nolint: gosec
			ServerName:         s.host,
		}

		conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), tlsconfig)
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

		if err = client.Mail(message.GetFrom()); err != nil {
			return err
		}

		for _, addr := range sendMailTo {
			if err = client.Rcpt(addr); err != nil {
				return err
			}
		}

		w, err := client.Data()
		if err != nil {
			return err
		}

		_, err = w.Write(msg)
		if err != nil {
			return err
		}

		err = w.Close()
		if err != nil {
			return err
		}

		return client.Quit()
	} else {
		err = smtp.SendMail(fmt.Sprintf("%s:%d", s.host, s.port), auth, message.GetFrom(), sendMailTo, msg)
	}

	return err
}
