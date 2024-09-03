package gmail

import "errors"

var (
	// ErrFailedToSendEmail is returned when an email fails to send
	ErrFailedToSendEmail = errors.New("failed to send email")
	// ErrNoUsersMessagesService is returned when no UsersMessagesService is initiated
	ErrNoUsersMessagesService = errors.New("no UsersMessagesService initiated")
	// ErrUnableToBuildMIMEMessage is returned when a MIME message fails to be built
	ErrUnableToBuildMIMEMessage = errors.New("unable to build MIME message")
	// ErrUnableToStartGmailService is returned when a Gmail service fails to start
	ErrUnableToStartGmailService = errors.New("unable to start Gmail service")
	// ErrUnableToParseServiceAccount is returned when a service account fails to be parsed
	ErrUnableToParseServiceAccount = errors.New("unable to parse service account")
	// ErrUnableToParseJWTCredentials is returned when JWT credentials fail to be parsed
	ErrUnableToParseJWTCredentials = errors.New("unable to parse JWT credentials")
	// ErrInvalidToken is returned when a token is invalid
	ErrInvalidToken = errors.New("invalid token")
	// ErrMockServiceError is returned when a mock service encounters an error
	ErrMockServiceError = errors.New("mock service error")
)
