[![Build status](https://badge.buildkite.com/97ed7beda0c4aca086a7b4d439855bef106e4a7bdac5c32dbd.svg)](https://buildkite.com/theopenlane/newman)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=theopenlane_newman&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=theopenlane_newman)
[![Go Report Card](https://goreportcard.com/badge/github.com/theopenlane/newman)](https://goreportcard.com/report/github.com/theopenlane/newman)
[![Go Reference](https://pkg.go.dev/badge/github.com/theopenlane/newman.svg)](https://pkg.go.dev/github.com/theopenlane/newman)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache2.0-brightgreen.svg)](https://opensource.org/licenses/Apache-2.0)

# newman

Newman is the postal worker that lives down the hall, delivering your email like a ring-tailed lemur. The newman project allows you to send emails using different email providers
such as Gmail, SendGrid, etc.

![newman](img/newman.png)

This project is organized into several sub-packages:
  - providers: managing various email providers
  - credentials: managing email credentials
  - shared: utilities and types
  - scrubber: sanitizing email content

## Features

- Send emails using various providers
- Support for attachments and both plain text and HTML content
- Scrubber / sanitization for not getting hex0rz

## Usage

To use the library, you need to import the desired provider package and create an instance of the email sender then call the `sendmail` function

```go
package main

import (
	"context"
	"log"

	"github.com/theopenlane/newman"
	"github.com/theopenlane/newman/providers/resend"
)

func main() {
    sender, err := resend.New("your_resend_api_token")
    if err != nil {
        log.Fatal(err)
    }

    msg := newman.NewEmailMessageWithOptions(
        newman.WithFrom("no-reply@youremailaddress.com"),
        newman.WithTo([]string{"mitb@emailsendingfun.com"}),
        newman.WithSubject("Isn't sending emails with golang fun?"),
        newman.WithHTML("<p>Oh Yes! Mark my words, Seinfeld! Your day of reckoning is coming</p>"),
    )

    if err := sender.SendEmail(msg); err != nil {
        log.Fatal(err)
    }
}
```

This package supports various email providers and can be extended to include more

### Implemented providers

  - Gmail
  - SendGrid
  - Mailgun
  - Postmark
  - SMTP
  - Resend

## Contributing

See the [.github/CONTRIBUTING.md](.github/CONTRIBUTING.md) guide for more information