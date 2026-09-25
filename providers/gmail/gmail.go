package gmail

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/credentials"
)

const (
	defaultUser = "me"
	// insufficientPermissionsReason is the Google API error reason for a token missing a required scope
	insufficientPermissionsReason = "insufficientPermissions"
)

// gmailEmailSender wraps the Gmail UsersService
type gmailEmailSender struct {
	users *gmail.UsersService
	user  string
}

// SendEmail satisfies the EmailSender interface
func (s *gmailEmailSender) SendEmail(message *newman.EmailMessage) error {
	return s.SendEmailWithContext(context.Background(), message)
}

// SendBatchEmail satisfies the EmailSender interface
func (s *gmailEmailSender) SendBatchEmail(_ []*newman.EmailMessage) error {
	return newman.ErrBatchNotImplemented
}

// SendBatchEmailWithContext satisfies the EmailSender interface
func (s *gmailEmailSender) SendBatchEmailWithContext(_ context.Context, _ []*newman.EmailMessage) error {
	return newman.ErrBatchNotImplemented
}

// SendEmailWithContext satisfies the EmailSender interface
func (s *gmailEmailSender) SendEmailWithContext(_ context.Context, message *newman.EmailMessage) error {
	mimeMessage, err := newman.BuildMimeMessage(message)
	if err != nil {
		return ErrUnableToBuildMIMEMessage
	}

	mimeMessage = addBCCRecipients(mimeMessage, message.GetBCC())

	gMessage := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(mimeMessage),
	}

	if _, err = s.send(gMessage); err != nil {
		return ErrFailedToSendEmail
	}

	return nil
}

// send a Gmail message
func (s *gmailEmailSender) send(message *gmail.Message) (*gmail.Message, error) {
	if s.users == nil {
		return nil, ErrNoUsersService
	}

	return s.users.Messages.Send(s.userID(), message).Do()
}

// userID returns the configured user or the default authenticated user
func (s *gmailEmailSender) userID() string {
	if s.user == "" {
		return defaultUser
	}

	return s.user
}

// Verify satisfies the EmailSender interface by fetching the profile of the configured user
func (s *gmailEmailSender) Verify(ctx context.Context) error {
	if s.users == nil {
		return ErrNoUsersService
	}

	_, err := s.users.GetProfile(s.userID()).Context(ctx).Do()

	switch {
	case err == nil, isInsufficientScope(err):
		return nil
	default:
		return fmt.Errorf("%w: %w", ErrVerifyFailed, err)
	}
}

// isInsufficientScope reports whether err is a 403 caused by the token lacking the read scope
func isInsufficientScope(err error) bool {
	var apiErr *googleapi.Error
	if !errors.As(err, &apiErr) || apiErr.Code != http.StatusForbidden {
		return false
	}

	return slices.ContainsFunc(apiErr.Errors, func(item googleapi.ErrorItem) bool {
		return item.Reason == insufficientPermissionsReason
	})
}

// addBCCRecipients adds BCC recipients to the message
func addBCCRecipients(message []byte, bccs []string) []byte {
	if len(bccs) == 0 {
		return message
	}

	var msg bytes.Buffer

	fmt.Fprintf(&msg, "Bcc: %s\r\n", strings.Join(bccs, ","))

	message = append(msg.Bytes(), message...)

	return message
}

// TokenManager defines an interface for obtaining OAuth2 tokens
type TokenManager interface {
	GetToken() ([]byte, error)
}

// NewWithOauth2 initializes a new gmailEmailSenderOauth2 instance using OAuth2 credentials
func NewWithOauth2(ctx context.Context, configJSON []byte, tokenManager TokenManager, user string) (newman.EmailSender, error) {
	config, err := credentials.ParseCredentials(configJSON)
	if err != nil {
		return nil, err
	}

	tokBytes, err := tokenManager.GetToken()
	if err != nil {
		return nil, err
	}

	tok, err := credentials.ParseToken(tokBytes)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, tok)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailEmailSender{users: srv.Users, user: user}, nil
}

// NewWithServiceAccount initializes a new gmailEmailSenderServiceAccount instance using service account JSON credentials
func NewWithServiceAccount(ctx context.Context, jsonCredentials []byte, user string) (newman.EmailSender, error) {
	params := google.CredentialsParams{
		Scopes:  []string{gmail.GmailSendScope},
		Subject: user,
	}

	creds, err := google.CredentialsFromJSONWithTypeAndParams(ctx, jsonCredentials, google.ServiceAccount, params)
	if err != nil {
		return nil, ErrUnableToParseServiceAccount
	}

	srv, err := gmail.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailEmailSender{users: srv.Users, user: user}, nil
}

// NewWithAPIKey initializes a new gmailEmailSenderAPIKey instance using an API key
func NewWithAPIKey(ctx context.Context, apiKey, user string) (newman.EmailSender, error) {
	srv, err := gmail.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailEmailSender{users: srv.Users, user: user}, nil
}

// NewWithJWTConfig initializes a new gmailEmailSenderJWT instance using JWT configuration
func NewWithJWTConfig(ctx context.Context, configJSON []byte, user string) (newman.EmailSender, error) {
	config, err := google.JWTConfigFromJSON(configJSON)
	if err != nil {
		return nil, ErrUnableToParseJWTCredentials
	}

	client := config.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailEmailSender{users: srv.Users, user: user}, nil
}

// NewWithJWTAccess initializes a new gmailEmailSenderJWTAccess instance using a JWT access token
func NewWithJWTAccess(ctx context.Context, jsonCredentials []byte, user string) (newman.EmailSender, error) {
	tokenSource, err := google.JWTAccessTokenSourceFromJSON(jsonCredentials, gmail.GmailSendScope)
	if err != nil {
		return nil, ErrUnableToParseJWTCredentials
	}

	srv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailEmailSender{users: srv.Users, user: user}, nil
}
