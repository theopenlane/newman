package shared

import (
	"regexp"
	"strings"
)

// regex for validating email addresses
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9._\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmailMessage(msg *EmailMessage) error {
	from := ValidateEmailAddress(msg.From)
	if from == "" {
		return newMissingRequiredFieldError("from")
	}

	to := ValidateEmailAddresses(msg.To)
	if len(to) == 0 {
		return newMissingRequiredFieldError("to")
	}

	return nil
}

// ValidateEmailAddress trims the email and checks if it is a valid email address
func ValidateEmailAddress(email string) string {
	trimmed := strings.TrimSpace(email)
	if !emailRegex.MatchString(trimmed) {
		return ""
	}

	return trimmed
}

// ValidateEmailAddresses trims and validates each email in the slice
func ValidateEmailAddresses(emails []string) []string {
	validEmails := []string{}

	for _, email := range emails {
		if validEmail := ValidateEmailAddress(email); validEmail != "" {
			validEmails = append(validEmails, validEmail)
		}
	}

	return validEmails
}
