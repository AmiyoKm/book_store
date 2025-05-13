package mailer

import (
	"bytes"
	"errors"
	"text/template"

	gomail "gopkg.in/gomail.v2"
)

type goMailClient struct {
	fromEmail string
	apiKey    string
}

func NewGoMailClient(apiKey, fromEmail string) (goMailClient, error) {
	if apiKey == "" {
		return goMailClient{}, errors.New("mail api key required")
	}

	return goMailClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
	}, nil
}

func (m goMailClient) Send(templateFile, username, email string, data any, isSandBox bool) (int, error) {
 	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
 	if err != nil {
 		return -1, err
 	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("smtp.gmail.com", 587, m.fromEmail, m.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return -1, err
	}
	return 200, nil
}
