# Providers

The providers package allows you to send emails using different email providers such as Resend, gmail, whatever. It abstracts the provider-specific
details and provides a simple interface for sending emails.

# Usage

To use the package, you need to create an instance of the email sender and then call the sendemail func, ex:

```go

package main

import (
    "github.com/theopenlane/newman"
    "github.com/theopenlane/newman/providers/resend"
)

func main() {
    client := resend.New(token, opts...)

    err := client.SendEmail(newman.NewEmailMessage([]string{"info@theopenlane.io"},"Hey openlane please have my money","We got hotdogs for sale"))
    if err != nil {
        log.Fatal(err)
    }
}

```