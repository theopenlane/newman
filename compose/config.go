package compose

// URLKey identifies a URL field under the .URLS template object
// Values match the JSON field names on URLConfig and the template variable keys (e.g. {{ .URLS.Verify }})
type URLKey string

const (
	// URLKeyRoot identifies the root application URL
	URLKeyRoot URLKey = "Root"
	// URLKeyProduct identifies the product home URL
	URLKeyProduct URLKey = "Product"
	// URLKeyDocs identifies the documentation URL
	URLKeyDocs URLKey = "Docs"
	// URLKeyVerify identifies the email verification URL
	URLKeyVerify URLKey = "Verify"
	// URLKeyInvite identifies the organization invite URL
	URLKeyInvite URLKey = "Invite"
	// URLKeyPasswordReset identifies the password reset URL
	URLKeyPasswordReset URLKey = "PasswordReset"
	// URLKeyVerifySubscriber identifies the subscriber verification URL
	URLKeyVerifySubscriber URLKey = "VerifySubscriber"
	// URLKeyVerifyBilling identifies the billing verification URL
	URLKeyVerifyBilling URLKey = "VerifyBilling"
	// URLKeyQuestionnaire identifies the questionnaire URL
	URLKeyQuestionnaire URLKey = "Questionnaire"
)

// URLConfig holds the URL values injected into all email template contexts
// JSON tags use PascalCase to match the template variable keys (e.g. {{ .URLS.Root }})
// koanf tags use lowercase to match config file conventions
type URLConfig struct {
	// Root is the root application URL
	Root string `json:"Root" koanf:"root" jsonschema:"required,description=Root application URL"`
	// Product is the product home URL
	Product string `json:"Product" koanf:"product" jsonschema:"required,description=Product home URL"`
	// Docs is the documentation URL
	Docs string `json:"Docs" koanf:"docs" jsonschema:"required,description=Documentation URL"`
	// Verify is the email verification URL
	Verify string `json:"Verify" koanf:"verify" jsonschema:"description=Email verification URL"`
	// Invite is the organization invite URL
	Invite string `json:"Invite" koanf:"invite" jsonschema:"description=Organization invite URL"`
	// PasswordReset is the password reset URL
	PasswordReset string `json:"PasswordReset" koanf:"reset" jsonschema:"description=Password reset URL"`
	// VerifySubscriber is the subscriber verification URL
	VerifySubscriber string `json:"VerifySubscriber" koanf:"verifysubscriber" jsonschema:"description=Subscriber verification URL"`
	// VerifyBilling is the billing verification URL
	VerifyBilling string `json:"VerifyBilling" koanf:"verifybilling" jsonschema:"description=Billing verification URL"`
	// Questionnaire is the questionnaire URL
	Questionnaire string `json:"Questionnaire" koanf:"questionnaire" jsonschema:"description=Questionnaire URL"`
}

// ToMap returns the URL values keyed by URLKey constants
// This is the canonical mapping between URLKey values and URLConfig fields
func (u URLConfig) ToMap() map[URLKey]string {
	return map[URLKey]string{
		URLKeyRoot:             u.Root,
		URLKeyProduct:          u.Product,
		URLKeyDocs:             u.Docs,
		URLKeyVerify:           u.Verify,
		URLKeyInvite:           u.Invite,
		URLKeyPasswordReset:    u.PasswordReset,
		URLKeyVerifySubscriber: u.VerifySubscriber,
		URLKeyVerifyBilling:    u.VerifyBilling,
		URLKeyQuestionnaire:    u.Questionnaire,
	}
}

// Config holds the sender and branding configuration injected into all email template contexts
// JSON tags use PascalCase to match the template variable keys (e.g. {{ .CompanyName }})
// Koanf tags use lowercase to match config file conventions
type Config struct {
	// CompanyName is the display name of the sending company
	CompanyName string `json:"CompanyName" koanf:"companyname" jsonschema:"required,description=Company display name"`
	// CompanyAddress is the mailing address of the company, used in footer copy
	CompanyAddress string `json:"CompanyAddress" koanf:"companyaddress" jsonschema:"description=Company mailing address"`
	// Corporation is the legal corporation name, used in copyright notices
	Corporation string `json:"Corporation" koanf:"corporation" jsonschema:"description=Legal corporation name"`
	// Year is the current year for copyright notices
	Year int `json:"Year" koanf:"year" jsonschema:"required,description=Current year for copyright notices"`
	// FromEmail is the sender email address
	FromEmail string `json:"FromEmail" koanf:"fromemail" jsonschema:"required,description=Sender email address"`
	// SupportEmail is the support contact email address
	SupportEmail string `json:"SupportEmail" koanf:"supportemail" jsonschema:"description=Support contact email address"`
	// LogoURL is the company logo image URL
	LogoURL string `json:"LogoURL" koanf:"logourl" jsonschema:"description=Company logo image URL"`
	// URLS holds the URL map for template links
	URLS URLConfig `json:"URLS" koanf:"urls" jsonschema:"required,description=URL map for template links"`
}
