package credentials

import (
	"errors"
)

var (
	// ErrFailedToLoadCredentials is returned when credentials fail to load
	ErrFailedToLoadCredentials = errors.New("failed to load credentials")
	// ErrFailedToParseToken is returned when a token fails to parse
	ErrFailedToParseToken = errors.New("failed to parse token")
)
