package integrations

import (
	"context"

	"github.com/mrz1836/postmark"
)

type Email struct {
	To      string
	From    string
	Subject string

	TextBody string
	// HTMLBody string
}

type MailerService interface {
	SendMail(from, to, subject, body string) error
}

type mailerService struct {
	client *postmark.Client
}

func NewMailerClient(key string) MailerService {
	postmark := postmark.NewClient(key, "")
	return &mailerService{client: postmark}
}

func (m *mailerService) SendMail(from, to, subject, body string) error {
	email := postmark.Email{
		From:    from,
		To:      to,
		Subject: subject,
		// HTMLBody:   "",
		TextBody: body,
		// Tag:        "pw-reset",
		TrackOpens: false,
	}
	_, err := m.client.SendEmail(context.Background(), email)
	if err != nil {
		return err
	}
	return nil
}
