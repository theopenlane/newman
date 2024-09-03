package main

import (
	"context"
	"log"

	"github.com/theopenlane/newman/providers/resend"
	"github.com/theopenlane/newman/shared"
)

func main() {
	sender := resend.NewResendEmailSender("")

	attach, err := shared.NewAttachmentFromFile("cattyping.gif")
	if err != nil {
		log.Fatal(err)
	}

	newmsg := shared.NewEmailMessageWithOptions(shared.WithAttachments([]*shared.Attachment{attach}), shared.WithTo([]string{"nerds-unite-again-aaaan2dzqcf4dvqjrsyavrbb2q@theopenlane.slack.com"}), shared.WithFrom("no-reply@mail.theopenlane.io"), shared.WithSubject("gosh matt is the best!"), shared.WithText("self referencing variadic functions are the best"))

	//	msg := newman.NewEmailMessage("no-reply@mail.theopenlane.io", []string{"nerds-unite-again-aaaan2dzqcf4dvqjrsyavrbb2q@theopenlane.slack.com"}, "Isnt sending emails with golang fun?", "Check out Matt's new fancy pants package working like a fucking charm")
	//
	if err := sender.SendEmail(context.TODO(), *newmsg); err != nil {
		log.Fatal(err)
	}
}
