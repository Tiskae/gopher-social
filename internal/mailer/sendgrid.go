package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendgrid(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		fromEmail,
		apiKey,
		client,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, isSandbox bool, data any) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// Template parsing logic
	tmpl, err := template.ParseFS(FS, "/templates"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	// retry mechanism
	for i := range MaxRetries {
		response, err := m.client.Send(message)

		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, MaxRetries)
			log.Printf("Error: %v", err.Error())

			// exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))

			continue
		}

		log.Printf("Email sent successfully with status code %v", response.StatusCode)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", MaxRetries)
}
