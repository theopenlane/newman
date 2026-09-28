package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/theopenlane/newman"
)

const (
	defaultConnectionMethod = "IMPLICIT"
	TLSConnection           = "TLS"
	CRAMMD5Auth             = "CRAM-MD5"
	tcpNetwork              = "tcp"
	startTLSExtension       = "STARTTLS"
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
	// tlsConfig allows custom TLS configuration for testing
	tlsConfig *tls.Config
}

// New creates a new instance of smtpEmailSender
func New(host string, port int, user, password string, authMethod string) (newman.EmailSender, error) {
	return NewWithConnMethod(host, port, user, password, authMethod, defaultConnectionMethod)
}

// NewWithConnMethod creates a new instance of smtpEmailSender with the specified connection method
func NewWithConnMethod(host string, port int, user, password string, authMethod string, connectionMethod string) (newman.EmailSender, error) {
	return &smtpEmailSender{
		host:             host,
		port:             port,
		user:             user,
		password:         password,
		authMethod:       authMethod,
		connectionMethod: connectionMethod,
		tlsConfig:        nil,
	}, nil
}

// SendEmail satisfies the EmailSender interface
func (s *smtpEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendBatchEmail satisfies the EmailSender interface
func (s *smtpEmailSender) SendBatchEmail(_ []*newman.EmailMessage) error {
	return newman.ErrBatchNotImplemented
}

// SendBatchEmailWithContext satisfies the EmailSender interface
func (s *smtpEmailSender) SendBatchEmailWithContext(_ context.Context, _ []*newman.EmailMessage) error {
	return newman.ErrBatchNotImplemented
}

// SendEmailWithContext satisfies the EmailSender interface
func (s *smtpEmailSender) SendEmailWithContext(ctx context.Context, message *newman.EmailMessage) error {
	if ctx == nil {
		ctx = context.Background()
	}

	sendMailTo := message.GetTo()
	sendMailTo = append(sendMailTo, message.GetCC()...)
	sendMailTo = append(sendMailTo, message.GetBCC()...)

	msg, err := newman.BuildMimeMessage(message)
	if err != nil {
		return err
	}

	auth := s.auth()

	if s.connectionMethod == TLSConnection {
		return s.secureSend(ctx, auth, message.GetFrom(), sendMailTo, msg)
	}

	return s.send(auth, message.GetFrom(), sendMailTo, msg)
}

// Verify satisfies the EmailSender interface by connecting and authenticating without sending mail
func (s *smtpEmailSender) Verify(ctx context.Context) error {
	var (
		conn net.Conn
		err  error
	)

	switch s.connectionMethod {
	case TLSConnection:
		dialer := tls.Dialer{Config: s.clientTLSConfig()}
		conn, err = dialer.DialContext(ctx, tcpNetwork, s.address())
	default:
		var dialer net.Dialer
		conn, err = dialer.DialContext(ctx, tcpNetwork, s.address())
	}

	if err != nil {
		return fmt.Errorf("%w: %w", ErrVerifyFailed, err)
	}

	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		conn.Close()

		return verifyError(ctx, err)
	}

	defer client.Close()

	if s.connectionMethod != TLSConnection {
		if ok, _ := client.Extension(startTLSExtension); ok {
			if err = client.StartTLS(s.clientTLSConfig()); err != nil {
				return verifyError(ctx, err)
			}
		}
	}

	if err = client.Auth(s.auth()); err != nil {
		return verifyError(ctx, err)
	}

	if err = client.Quit(); err != nil {
		return verifyError(ctx, err)
	}

	return nil
}

// verifyError wraps err with ErrVerifyFailed, adding the context error when ctx is done
func verifyError(ctx context.Context, err error) error {
	if ctxErr := ctx.Err(); ctxErr != nil {
		return fmt.Errorf("%w: %w: %w", ErrVerifyFailed, ctxErr, err)
	}

	return fmt.Errorf("%w: %w", ErrVerifyFailed, err)
}

// auth returns the smtp.Auth for the configured auth method
func (s *smtpEmailSender) auth() smtp.Auth {
	switch s.authMethod {
	case CRAMMD5Auth:
		return smtp.CRAMMD5Auth(s.user, s.password)
	default:
		return smtp.PlainAuth("", s.user, s.password, s.host)
	}
}

// clientTLSConfig returns the configured TLS settings or a default config for the host
func (s *smtpEmailSender) clientTLSConfig() *tls.Config {
	if s.tlsConfig != nil {
		return s.tlsConfig
	}

	return &tls.Config{
		ServerName: s.host,
		MinVersion: tls.VersionTLS12,
	}
}

// address returns the host:port of the SMTP server
func (s *smtpEmailSender) address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

func (s *smtpEmailSender) send(auth smtp.Auth, from string, to []string, message []byte) error {
	return smtp.SendMail(s.address(), auth, from, to, message)
}

func (s *smtpEmailSender) secureSend(ctx context.Context, auth smtp.Auth, from string, to []string, message []byte) error {
	dialer := tls.Dialer{Config: s.clientTLSConfig()}

	conn, err := dialer.DialContext(ctx, tcpNetwork, s.address())
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
