package compose

// Recipient holds per-recipient data injected into all email template contexts
// JSON tags match the template variable keys (e.g. {{ .Recipient.Email }})
type Recipient struct {
	// Email is the recipient's email address
	Email string `json:"Email" jsonschema:"required,description=Recipient email address"`
	// FirstName is the recipient's first name
	FirstName string `json:"FirstName" jsonschema:"description=Recipient first name"`
	// LastName is the recipient's last name
	LastName string `json:"LastName" jsonschema:"description=Recipient last name"`
}
