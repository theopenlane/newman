package shared_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/theopenlane/newman/shared"
)

func TestValidateEmail(t *testing.T) {
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
			result := shared.ValidateEmail(test.email)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestValidateEmailSlice(t *testing.T) {
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
			result := shared.ValidateEmailSlice(test.emails)
			assert.Equal(t, test.expected, result)
		})
	}
}

func ExampleValidateEmail() {
	email := "newman@usps.com"
	result := shared.ValidateEmail(email)
	fmt.Println(result)
}

func ExampleValidateEmail_not() {
	email := "test@com"
	result := shared.ValidateEmail(email)
	fmt.Println(result)
}

func ExampleValidateEmail_trim() {
	email := "  newman@usps.com  "
	result := shared.ValidateEmail(email)
	fmt.Println(result)
}

func ExampleValidateEmailSlice() {
	emails := []string{"newman@usps.com", "test@domain_name.com"}
	result := shared.ValidateEmailSlice(emails)
	fmt.Println(result)
}
func ExampleValidateEmailSlice_partial() {
	emails := []string{"newman@usps.com", "test@com"}
	result := shared.ValidateEmailSlice(emails)
	fmt.Println(result)
}
