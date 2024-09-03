package credentials

import (
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// ParseCredentials parses the OAuth2 credentials JSON byte slice and returns an *oauth2.Config
func ParseCredentials(credentialsJSON []byte) (*oauth2.Config, error) {
	config, err := google.ConfigFromJSON(credentialsJSON, gmail.GmailSendScope)

	if err != nil {
		return nil, ErrFailedToLoadCredentials
	}

	return config, nil
}

// ParseToken parses the OAuth2 token JSON byte slice and returns an *oauth2.Token
func ParseToken(tokenJSON []byte) (*oauth2.Token, error) {
	token := &oauth2.Token{}

	if err := json.Unmarshal(tokenJSON, token); err != nil {
		return nil, ErrFailedToParseToken
	}

	return token, nil
}
