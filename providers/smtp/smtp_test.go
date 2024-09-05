package smtp

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/newman"
)

// TestEmailSenderImplementation checks if smtpEmailSender implements the EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*smtpEmailSender)(nil)
}

// newMockSMTPServer creates a mock SMTP server for testing purposes
func newMockSMTPServer(t *testing.T, handler func(conn net.Conn)) *mockSMTPServer {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoErrorf(t, err, "failed to start mock SMTP server")

	server := &mockSMTPServer{
		listener: listener,
		addr:     listener.Addr().String(),
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}

			go handler(conn)
		}
	}()

	return server
}

// mockSMTPServer represents a mock SMTP server
type mockSMTPServer struct {
	listener net.Listener
	addr     string
}

func (s *mockSMTPServer) Close() {
	s.listener.Close()
}

const expectedPassword = "We gotta find that rickshaw" // nolint: gosec

// smtpHandler is the handler for plain SMTP connections
func smtpHandler(conn net.Conn) {
	fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		cmd := strings.TrimSpace(string(buf[:n]))

		switch {
		case strings.HasPrefix(cmd, "EHLO"):
			fmt.Fprintln(conn, "250-Hello")
			fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")

		case strings.HasPrefix(cmd, "AUTH PLAIN"):
			creds := decodeConnectionCommand("AUTH PLAIN", cmd)
			if len(creds) == 2 && creds[1] == expectedPassword {
				fmt.Fprintln(conn, "235 Authentication successful")
			} else {
				fmt.Fprintln(conn, "535 Authentication failed")
			}

		case strings.HasPrefix(cmd, "AUTH CRAM-MD5"):
			fmt.Fprintln(conn, "334 "+base64.StdEncoding.EncodeToString([]byte("challenge")))

			_, err := conn.Read(buf)
			if err != nil {
				break
			}

			fmt.Fprintln(conn, "235 Authentication successful")

		case strings.HasPrefix(cmd, "MAIL FROM"):
			fmt.Fprintln(conn, "250 OK")

		case strings.HasPrefix(cmd, "RCPT TO"):
			fmt.Fprintln(conn, "250 OK")

		case strings.HasPrefix(cmd, "DATA"):
			fmt.Fprintln(conn, "354 Start mail input; end with <CRLF>.<CRLF>")

			_, err := conn.Read(buf)
			if err != nil {
				break
			}

			fmt.Fprintln(conn, "250 OK")

		case strings.HasPrefix(cmd, "QUIT"):
			fmt.Fprintln(conn, "221 Bye")
			conn.Close()

			return
		default:
			fmt.Fprintln(conn, "500 Unrecognized command")
		}
	}
}

// tlsHandler is the handler for TLS connections
func tlsHandler(handler func(conn net.Conn)) func(conn net.Conn) {
	return func(conn net.Conn) {
		cert := generateKeys()

		tlsConn := tls.Server(conn, &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true, // nolint: gosec
		})

		err := tlsConn.Handshake()
		if err != nil {
			tlsConn.Close()
			return
		}

		handler(tlsConn)
	}
}

func decodeConnectionCommand(cmd, message string) []string {
	cleanedCmd := strings.TrimSpace(cmd)
	stringB64 := strings.TrimPrefix(message, cleanedCmd+" ")

	decoded, err := base64.StdEncoding.DecodeString(stringB64)
	if err != nil {
		return nil
	}

	trimmedDecoded := strings.Trim(string(decoded), "\x00")

	return strings.Split(trimmedDecoded, "\x00")
}
func TestNewSMTPEmailSender(t *testing.T) {
	emailSender, err := New("smtp.example.com", 587, "user", "We gotta find that rickshaw", "PLAIN")
	assert.NoError(t, err)
	assert.NotNil(t, emailSender)
}

func TestSendEmailPlainAuth(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := New(host, portInt, "user", "We gotta find that rickshaw", "PLAIN")
	assert.NoError(t, err)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailCramMD5Auth(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := New(host, portInt, "user", "We gotta find that rickshaw", "CRAM-MD5")
	assert.NoError(t, err)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailError(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := New(host, portInt, "user", "wrongWe gotta find that rickshaw", "PLAIN")
	assert.NoError(t, err)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

// //////
func TestSendEmailInvalidServer(t *testing.T) {
	emailSender, err := New("invalid.server.com", 587, "user", "We gotta find that rickshaw", "PLAIN")
	assert.NoError(t, err)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailMissingSettings(t *testing.T) {
	emailSender, err := New("", 0, "", "", "PLAIN")
	assert.NoError(t, err)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendEmailImplicitTLS(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 465 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "IMPLICIT")
	assert.NoError(t, err)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailExplicitTLS(t *testing.T) {
	server := newMockSMTPServer(t, tlsHandler(smtpHandler))
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)

	portInt := 587 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err = emailSender.SendEmail(message)
	assert.NoError(t, err)
}

func TestSendEmailExplicitErrorTLS(t *testing.T) {
	server := newMockSMTPServer(t, smtpHandler)
	defer server.Close()

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 587

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Test Email", "The air is so dewy sweet you dont even have to lick the stamps")

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendTLSEmailConnectionError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				conn.Close()
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{"nonexistent@example.com"},
		"The air is so dewy sweet you dont even have to lick the stamps",
		"The air is so dewy sweet you dont even have to lick the stamps",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
}

func TestSendTLSEmailEHLOError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")

				buf := make([]byte, 1024)

				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}

					cmd := strings.TrimSpace(string(buf[:n]))

					switch {
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()

						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))

	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{"nonexistent@example.com"},
		"Test Email",
		"The air is so dewy sweet you dont even have to lick the stamps",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "500 Unrecognized command", err.Error())
}

func TestSendTLSEmailAUTHError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")

				buf := make([]byte, 1024)

				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}

					cmd := strings.TrimSpace(string(buf[:n]))

					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "403 Forbidden")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()

						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))

	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{"nonexistent@example.com"},
		"The air is so dewy sweet you dont even have to lick the stamps",
		"The air is so dewy sweet you dont even have to lick the stamps",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "403 Forbidden", err.Error())
}

func TestSendTLSEmailMailError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")

				buf := make([]byte, 1024)

				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}

					cmd := strings.TrimSpace(string(buf[:n]))

					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "235 Authentication successful")
					case strings.HasPrefix(cmd, "MAIL FROM"):
						fmt.Fprintln(conn, "550 MAIL ERROR")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()

						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))
	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{"nonexistent@example.com"},
		"The air is so dewy sweet you dont even have to lick the stamps",
		"The air is so dewy sweet you dont even have to lick the stamps",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "550 MAIL ERROR", err.Error())
}

func TestSendTLSEmailRcptError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")

				buf := make([]byte, 1024)

				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}

					cmd := strings.TrimSpace(string(buf[:n]))

					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "235 Authentication successful")
					case strings.HasPrefix(cmd, "MAIL FROM"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "RCPT TO"):
						fmt.Fprintln(conn, "530 No such user")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()

						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))

	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{"nonexistent@example.com"},
		"The air is so dewy sweet you dont even have to lick the stamps",
		"The air is so dewy sweet you dont even have to lick the stamps",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "530 No such user", err.Error())
}

func TestSendTLSEmailDataError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")

				buf := make([]byte, 1024)

				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}

					cmd := strings.TrimSpace(string(buf[:n]))

					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "235 Authentication successful")
					case strings.HasPrefix(cmd, "MAIL FROM"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "RCPT TO"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "DATA"):
						fmt.Fprintln(conn, "554 Service unavailable")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()

						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))

	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{"nonexistent@example.com"},
		"The air is so dewy sweet you dont even have to lick the stamps",
		"The air is so dewy sweet you dont even have to lick the stamps",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "554 Service unavailable", err.Error())
}

func TestSendTLSEmailDataWriteError(t *testing.T) {
	server := newMockSMTPServer(t,
		tlsHandler(
			func(conn net.Conn) {
				fmt.Fprintln(conn, "220 Welcome to the Mock SMTP Server")

				buf := make([]byte, 1024)

				for {
					n, err := conn.Read(buf)
					if err != nil {
						break
					}

					cmd := strings.TrimSpace(string(buf[:n]))

					switch {
					case strings.HasPrefix(cmd, "EHLO"):
						fmt.Fprintln(conn, "250-Hello")
						fmt.Fprintln(conn, "250 AUTH CRAM-MD5 PLAIN LOGIN")
					case strings.HasPrefix(cmd, "AUTH PLAIN"):
						fmt.Fprintln(conn, "235 Authentication successful")
					case strings.HasPrefix(cmd, "MAIL FROM"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "RCPT TO"):
						fmt.Fprintln(conn, "250 OK")
					case strings.HasPrefix(cmd, "DATA"):
						fmt.Fprintln(conn, "354 Start mail input; end with <CRLF>.<CRLF>")

						_, err := conn.Read(buf)
						if err != nil {
							break
						}

						fmt.Fprintln(conn, "552 Message size exceeds fixed limit")
					case strings.HasPrefix(cmd, "QUIT"):
						fmt.Fprintln(conn, "221 Bye")
						conn.Close()

						return
					default:
						fmt.Fprintln(conn, "500 Unrecognized command")
					}
				}
			}))

	defer server.Close()

	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	host, port, _ := net.SplitHostPort(server.addr)
	portInt := 25 // nolint: mnd

	_, err := fmt.Sscanf(port, "%d", &portInt)
	if err != nil {
		t.Errorf("failed to parse port: %v", err)
	}

	emailSender, err := NewWithConnMethod(host, portInt, "user", "We gotta find that rickshaw", "PLAIN", "TLS")
	assert.NoError(t, err)

	message := newman.NewEmailMessage(
		"newman@usps.com",
		[]string{"nonexistent@example.com"},
		"The air is so dewy sweet you dont even have to lick the stamps",
		"The air is so dewy sweet you dont even have to lick the stamps",
	)

	err = emailSender.SendEmail(message)
	assert.Error(t, err)
	assert.Equal(t, "552 Message size exceeds fixed limit", err.Error())
}

func generateKeys() tls.Certificate {
	cert, err := tls.X509KeyPair([]byte(ecdsaCertPEM), []byte(ecdsaKeyPEM))
	if err != nil {
		log.Panicf("error X509KeyPair %v", err)
	}

	return cert
}

func testingKey(s string) string { return strings.ReplaceAll(s, "TESTING KEY", "PRIVATE KEY") }

// keys from https://go.dev/src/crypto/tls/tls_test.go

var ecdsaCertPEM = `-----BEGIN CERTIFICATE-----
MIIB/jCCAWICCQDscdUxw16XFDAJBgcqhkjOPQQBMEUxCzAJBgNVBAYTAkFVMRMw
EQYDVQQIEwpTb21lLVN0YXRlMSEwHwYDVQQKExhJbnRlcm5ldCBXaWRnaXRzIFB0
eSBMdGQwHhcNMTIxMTE0MTI0MDQ4WhcNMTUxMTE0MTI0MDQ4WjBFMQswCQYDVQQG
EwJBVTETMBEGA1UECBMKU29tZS1TdGF0ZTEhMB8GA1UEChMYSW50ZXJuZXQgV2lk
Z2l0cyBQdHkgTHRkMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQBY9+my9OoeSUR
lDQdV/x8LsOuLilthhiS1Tz4aGDHIPwC1mlvnf7fg5lecYpMCrLLhauAc1UJXcgl
01xoLuzgtAEAgv2P/jgytzRSpUYvgLBt1UA0leLYBy6mQQbrNEuqT3INapKIcUv8
XxYP0xMEUksLPq6Ca+CRSqTtrd/23uTnapkwCQYHKoZIzj0EAQOBigAwgYYCQXJo
A7Sl2nLVf+4Iu/tAX/IF4MavARKC4PPHK3zfuGfPR3oCCcsAoz3kAzOeijvd0iXb
H5jBImIxPL4WxQNiBTexAkF8D1EtpYuWdlVQ80/h/f4pBcGiXPqX5h2PQSQY7hP1
+jwM1FGS4fREIOvlBYr/SzzQRtwrvrzGYxDEDbsC0ZGRnA==
-----END CERTIFICATE-----
`

var ecdsaKeyPEM = testingKey(`-----BEGIN EC PARAMETERS-----
BgUrgQQAIw==
-----END EC PARAMETERS-----
-----BEGIN EC TESTING KEY-----
MIHcAgEBBEIBrsoKp0oqcv6/JovJJDoDVSGWdirrkgCWxrprGlzB9o0X8fV675X0
NwuBenXFfeZvVcwluO7/Q9wkYoPd/t3jGImgBwYFK4EEACOhgYkDgYYABAFj36bL
06h5JRGUNB1X/Hwuw64uKW2GGJLVPPhoYMcg/ALWaW+d/t+DmV5xikwKssuFq4Bz
VQldyCXTXGgu7OC0AQCC/Y/+ODK3NFKlRi+AsG3VQDSV4tgHLqZBBus0S6pPcg1q
kohxS/xfFg/TEwRSSws+roJr4JFKpO2t3/be5OdqmQ==
-----END EC TESTING KEY-----
`)
