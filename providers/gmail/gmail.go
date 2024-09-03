package gmail

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/credentials"
)

// gmailMessageSenderWrapper wraps the Gmail UsersMessagesService
type gmailMessageSenderWrapper struct {
	messageSender *gmail.UsersMessagesService
	user          string
}

// send sends a Gmail message
func (s *gmailMessageSenderWrapper) send(message *gmail.Message) (*gmail.Message, error) {
	if s.messageSender == nil {
		return nil, ErrNoUsersMessagesService
	}

	user := s.user
	if user == "" {
		user = "me"
	}

	return s.messageSender.Send(user, message).Do()
}

// SendEmail sends an email using the Gmail API
func (s *gmailMessageSenderWrapper) SendEmail(message *newman.EmailMessage) error {
	mimeMessage, err := newman.BuildMimeMessage(message)
	if err != nil {
		return ErrUnableToBuildMIMEMessage
	}

	bccs := message.GetBCC()
	if len(bccs) > 0 {
		var msg bytes.Buffer

		msg.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(bccs, ",")))

		mimeMessage = append(msg.Bytes(), mimeMessage...)
	}

	gMessage := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(mimeMessage),
	}

	_, err = s.send(gMessage)
	if err != nil {
		return ErrFailedToSendEmail
	}

	return nil
}

// GmailTokenManager defines an interface for obtaining OAuth2 tokens
type GmailTokenManager interface {
	GetToken() ([]byte, error)
}

// NewGmailEmailSenderOauth2 initializes a new gmailEmailSenderOauth2 instance using OAuth2 credentials
func NewGmailEmailSenderOauth2(ctx context.Context, configJSON []byte, tokenManager GmailTokenManager, user string) (*gmailMessageSenderWrapper, error) {
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

	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}

// NewGmailEmailSenderServiceAccount initializes a new gmailEmailSenderServiceAccount instance using service account JSON credentials
func NewGmailEmailSenderServiceAccount(ctx context.Context, jsonCredentials []byte, user string) (*gmailMessageSenderWrapper, error) {
	params := google.CredentialsParams{
		Scopes:  []string{gmail.GmailSendScope},
		Subject: user,
	}

	creds, err := google.CredentialsFromJSONWithParams(ctx, jsonCredentials, params)
	if err != nil {
		return nil, ErrUnableToParseServiceAccount
	}

	srv, err := gmail.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}

// NewGmailEmailSenderAPIKey initializes a new gmailEmailSenderAPIKey instance using an API key
func NewGmailEmailSenderAPIKey(ctx context.Context, apiKey, user string) (*gmailMessageSenderWrapper, error) {
	srv, err := gmail.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}

// NewGmailEmailSenderJWT initializes a new gmailEmailSenderJWT instance using JWT configuration
func NewGmailEmailSenderJWT(ctx context.Context, configJSON []byte, user string) (*gmailMessageSenderWrapper, error) {
	config, err := google.JWTConfigFromJSON(configJSON)
	if err != nil {
		return nil, ErrUnableToParseJWTCredentials
	}

	client := config.Client(ctx)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}

// NewGmailEmailSenderJWTAccess initializes a new gmailEmailSenderJWTAccess instance using a JWT access token
func NewGmailEmailSenderJWTAccess(ctx context.Context, jsonCredentials []byte, user string) (*gmailMessageSenderWrapper, error) {
	tokenSource, err := google.JWTAccessTokenSourceFromJSON(jsonCredentials, gmail.GmailSendScope)
	if err != nil {
		return nil, ErrUnableToParseJWTCredentials
	}

	srv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, ErrUnableToStartGmailService
	}

	return &gmailMessageSenderWrapper{messageSender: srv.Users.Messages, user: user}, nil
}
