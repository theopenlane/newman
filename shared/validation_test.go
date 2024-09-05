package shared_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/theopenlane/newman/shared"
)

func TestValidateValidateEmailMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  *shared.EmailMessage
		err  error
	}{
		{
			name: "valid email message",
			msg: &shared.EmailMessage{
				From: "mitb@example.com",
				To:   []string{"funky@funk.com", "waters@kitb.com"},
			},
			err: nil,
		},
		{
			name: "missing from",
			msg: &shared.EmailMessage{
				To: []string{"funky@funk.com"},
			},
			err: &shared.MissingRequiredFieldError{
				RequiredField: "from",
			},
		},
		{
			name: "missing to",
			msg: &shared.EmailMessage{
				From: "mitb@example.com",
			},
			err: &shared.MissingRequiredFieldError{
				RequiredField: "to",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resultErr := shared.ValidateEmailMessage(test.msg)
			assert.Equal(t, test.err, resultErr)
		})
	}
}
func TestValidateEmailAddress(t *testing.T) {
	tests := []struct {
		email    string
		expected string
	}{
		{"newman@usps.com", "newman@usps.com"},
		{"test@domain_name.com", "test@domain_name.com"},
		{"test@domain-name.com", "test@domain-name.com"},
		{"test@subdomain.mitbindustries.com", "test@subdomain.mitbindustries.com"},
		{"test_name@subdomain.mitbindustries.com", "test_name@subdomain.mitbindustries.com"},
		{"test.name@subdomain.mitbindustries.com", "test.name@subdomain.mitbindustries.com"},
		{"test-name@subdomain.mitbindustries.com", "test-name@subdomain.mitbindustries.com"},
		{"  newman@usps.com  ", "newman@usps.com"},
		{"invalid-email", ""},
		{"test@.com", ""},
		{"@mitbindustries.com", ""},
		{"test@com", ""},
		{"test@com.", ""},
		{"test@sub.mitbindustries.com", "test@sub.mitbindustries.com"},
		{"test+alias@mitbindustries.com", "test+alias@mitbindustries.com"},
		{"test.email@mitbindustries.com", "test.email@mitbindustries.com"},
		{"test-email@mitbindustries.com", "test-email@mitbindustries.com"},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			result := shared.ValidateEmailAddress(test.email)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestValidateEmailAddresses(t *testing.T) {
	tests := []struct {
		emails   []string
		expected []string
	}{
		{[]string{"newman@usps.com"}, []string{"newman@usps.com"}},
		{[]string{"newman@usps.com", "invalid-email"}, []string{"newman@usps.com"}},
		{[]string{" newman@usps.com ", "test2@mitbindustries.com"}, []string{"newman@usps.com", "test2@mitbindustries.com"}},
		{[]string{"invalid-email", "@mitbindustries.com"}, []string{}},
		{[]string{"newman@usps.com", "test2@sub.mitbindustries.com"}, []string{"newman@usps.com", "test2@sub.mitbindustries.com"}},
	}

	for _, test := range tests {
		t.Run(strings.Join(test.emails, ","), func(t *testing.T) {
			result := shared.ValidateEmailAddresses(test.emails)
			assert.Equal(t, test.expected, result)
		})
	}
}

func ExampleValidateEmailAddress() {
	email := "newman@usps.com"
	result := shared.ValidateEmailAddress(email)
	fmt.Println(result)
}

func ExampleValidateEmailAddress_not() {
	email := "test@com"
	result := shared.ValidateEmailAddress(email)
	fmt.Println(result)
}

func ExampleValidateEmailAddress_trim() {
	email := "  newman@usps.com  "
	result := shared.ValidateEmailAddress(email)
	fmt.Println(result)
}

func ExampleValidateEmailAddresses() {
	emails := []string{"newman@usps.com", "test@domain_name.com"}
	result := shared.ValidateEmailAddresses(emails)
	fmt.Println(result)
}
func ExampleValidateEmailAddresses_partial() {
	emails := []string{"newman@usps.com", "test@com"}
	result := shared.ValidateEmailAddresses(emails)
	fmt.Println(result)
}
