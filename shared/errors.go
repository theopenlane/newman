package shared

import "fmt"

// MissingRequiredFieldError is returned when a required field was not provided in a request
type MissingRequiredFieldError struct {
	// RequiredField that is missing
	RequiredField string `json:"required_field"`
}

// Error returns the InvalidEmailConfigError in string format
func (e *MissingRequiredFieldError) Error() string {
	return fmt.Sprintf("%s is required", e.RequiredField)
}

// NewMissingRequiredFieldError returns an error for a missing required field
func newMissingRequiredFieldError(field string) *MissingRequiredFieldError {
	return &MissingRequiredFieldError{
		RequiredField: field,
	}
}
