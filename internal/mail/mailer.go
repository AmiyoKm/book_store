package mailer

import "embed"

const (
	FromName               = "BookBand"
	maxRetries             = 3
	UserWelcomeTemplate    = "user_invitation.tmpl"
	PasswordChangeTemplate = "password_change.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
	SendPasswordRequestMail(templateFile, username, email string, data any, isSandbox bool) (int, error)
}
