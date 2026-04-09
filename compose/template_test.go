package compose

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testConfig() Config {
	return Config{
		CompanyName:  "Acme",
		FromEmail:    "no-reply@acme.io",
		SupportEmail: "support@acme.io",
		Year:         2025,
		URLS: URLConfig{
			Root:             "https://acme.io",
			Product:          "https://app.acme.io",
			Docs:             "https://docs.acme.io",
			Verify:           "https://app.acme.io/verify",
			Invite:           "https://app.acme.io/invite",
			PasswordReset:    "https://app.acme.io/password-reset",
			VerifySubscriber: "https://app.acme.io/subscriber-verify",
			VerifyBilling:    "https://app.acme.io/verify-billing",
			Questionnaire:    "https://app.acme.io/questionnaire",
		},
	}
}

func TestBuildTemplateURLs_AllKeys(t *testing.T) {
	cfg := testConfig()
	urls, err := BuildTemplateURLs(cfg, nil)
	require.NoError(t, err)

	assert.Equal(t, cfg.URLS.Root, urls["Root"])
	assert.Equal(t, cfg.URLS.Product, urls["Product"])
	assert.Equal(t, cfg.URLS.Docs, urls["Docs"])
	assert.Equal(t, cfg.URLS.Verify, urls["Verify"])
	assert.Equal(t, cfg.URLS.Invite, urls["Invite"])
	assert.Equal(t, cfg.URLS.PasswordReset, urls["PasswordReset"])
	assert.Equal(t, cfg.URLS.VerifySubscriber, urls["VerifySubscriber"])
	assert.Equal(t, cfg.URLS.VerifyBilling, urls["VerifyBilling"])
	assert.Equal(t, cfg.URLS.Questionnaire, urls["Questionnaire"])
}

func TestBuildTemplateURLs_OverridesApplied(t *testing.T) {
	cfg := testConfig()
	tokenURL := "https://app.acme.io/verify?token=abc123"
	urls, err := BuildTemplateURLs(cfg, map[string]any{"Verify": tokenURL})
	require.NoError(t, err)

	assert.Equal(t, tokenURL, urls["Verify"])
	assert.Equal(t, cfg.URLS.Root, urls["Root"])
}

func TestBuildTemplateData_BaseFields(t *testing.T) {
	cfg := testConfig()
	r := Recipient{Email: "ada@example.com", FirstName: "Ada", LastName: "Lovelace"}

	data, err := BuildTemplateData(cfg, r, nil)
	require.NoError(t, err)

	assert.Equal(t, "Acme", data["CompanyName"])
	assert.Equal(t, "no-reply@acme.io", data["FromEmail"])
	assert.Equal(t, 2025, data["Year"])

	recipient, ok := data["Recipient"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "ada@example.com", recipient["Email"])
	assert.Equal(t, "Ada", recipient["FirstName"])
	assert.Equal(t, "Lovelace", recipient["LastName"])
}

func TestBuildTemplateData_OverridesApplied(t *testing.T) {
	cfg := testConfig()
	r := Recipient{Email: "ada@example.com"}

	data, err := BuildTemplateData(cfg, r, map[string]any{
		"OrganizationName": "Openlane",
		"InviterName":      "Charles",
	})
	require.NoError(t, err)

	assert.Equal(t, "Openlane", data["OrganizationName"])
	assert.Equal(t, "Charles", data["InviterName"])
	assert.Equal(t, "Acme", data["CompanyName"])
}

func TestResolveYear_UsesConfigured(t *testing.T) {
	assert.Equal(t, 2020, resolveYear(2020))
}

func TestResolveYear_DefaultsToCurrentYear(t *testing.T) {
	year := resolveYear(0)
	assert.Greater(t, year, 2020)
}

func TestAddTokenToURL_AppendsToken(t *testing.T) {
	result, err := AddTokenToURL("https://app.acme.io/verify", "abc123")
	require.NoError(t, err)
	assert.True(t, strings.Contains(result, "token=abc123"))
}

func TestAddTokenToURL_InvalidURL(t *testing.T) {
	_, err := AddTokenToURL("://bad-url", "token")
	assert.ErrorIs(t, err, ErrInvalidURL)
}
