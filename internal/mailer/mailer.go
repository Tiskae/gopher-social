// Package mailer for sending mails
package mailer

import "embed"

const (
	FromName            = "GopherSocial"
	MaxRetries          = 3
	UserWelcomeTemplate = "/user_invitation.tmpl"
)

//go:embed templates
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, isSandbox bool, data any) (int, error)
}
