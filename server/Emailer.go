package server

import (
	"gopkg.in/mailgun/mailgun-go.v1"
)

type EmailSender interface {
	Send(address string, title string, body string)
}

type MailgunEmailSender struct {
	EmailSender
	Mail mailgun.Mailgun
}

func newMailGunSender(domain string, apiKey string, pubApiKey string) MailgunEmailSender {
	return MailgunEmailSender{Mail:newMailGun(domain,apiKey,pubApiKey)}
}

	func newMailGun(domain string, apiKey string, pubApiKey string) mailgun.Mailgun {
	f := mailgun.NewMailgun(domain,apiKey,pubApiKey)

	return f;
}

func (emailer MailgunEmailSender) Send(address string, tital string, body string) {
	
}