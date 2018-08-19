package server

import (
	"gopkg.in/mailgun/mailgun-go.v1"
	"os"
)

type EmailSender interface {
	Send(address string, title string, body string)
}

type MailgunEmailSender struct {
	EmailSender
	Mail mailgun.Mailgun
}

func newMailGunSender(domain string) MailgunEmailSender {
	apiKey := os.Getenv("MG_KEY")
	apiPubKey := os.Getenv("MG_VALIDATION_KEY")
	return MailgunEmailSender{Mail:newMailGun(domain,apiKey,apiPubKey)}
}

func newMailGun(domain string, apiKey string, pubApiKey string) mailgun.Mailgun {
	f := mailgun.NewMailgun(domain,apiKey,pubApiKey)

	return f;
}

func (emailer MailgunEmailSender) Send(address string, tital string, body string) {


	message := emailer.Mail.NewMessage("No Reply <noreply@emit.sh>", "Soneone sent you a file!", body,address)

	one, two, err := emailer.Mail.Send(message)

	if err != nil || len(one) < 0 || len(two) < 0 {

	}


}