# Overview

The common package includes essential utilities such as email validation, sanitization, MIME type determination, and structures for handling email
messages and attachments

# Components

  - `EmailMessage`: constructing and manipulating email messages
  - `Attachment`: managing email attachments, including file handling and base64 encoding
  - `Validation`: validating email addresses and slices of email addresses
  - `Sanitization`: sanitizing input to prevent injection attacks

# Usage

```go

package main

import (
    "github.com/theopenlane/newman/shared"
)

func main() {
    email := shared.NewEmailMessage("newman@usps.com", []string{"jerry@seinfeld.com"}, "Jumbaliyaaaaa", "Crease, crumple, cram. You'll do fine")
    fmt.Println(email.GetSubject())
}

```