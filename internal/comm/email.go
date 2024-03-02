package comm

import (
	"bytes"
	"text/template"
)

const signInOtpSubject = "One time password to verify your email."
const emailUpdateOtpSubject = "One time password to update your email."

const (
	signInOtpTemplate      EmailTemplateEnum = "signInOtpTemplate"
	emailUpdateOtpTemplate EmailTemplateEnum = "emailUpdateOtpTemplate"
)

var templates = map[EmailTemplateEnum]string{
	signInOtpTemplate:      "internal/templates/sign_in_otp_template.html",
	emailUpdateOtpTemplate: "internal/templates/email_update_otp_template.html",
}

type EmailTemplateEnum string

type Emailer interface {
	SendSignInOTP(email, otp string) error
	SendEmailUpdateOTP(email, otp string) error
}

type EmailClient interface {
	Send(recipients []string, subject string, body string) error
}

type emailer struct {
	client    EmailClient
	templates map[EmailTemplateEnum]*template.Template
}

type SendOTP struct {
	OTP string
}

func NewEmailer(client EmailClient) (Emailer, error) {
	emailer := &emailer{
		client:    client,
		templates: map[EmailTemplateEnum]*template.Template{},
	}
	for k, v := range templates {
		tmpl, err := template.ParseFiles(v)
		if err != nil {
			return nil, err
		}
		emailer.templates[k] = tmpl
	}

	return emailer, nil
}

func (e *emailer) SendSignInOTP(email, otp string) error {
	html := bytes.Buffer{}
	err := e.templates[signInOtpTemplate].Execute(&html, &SendOTP{OTP: otp})
	if err != nil {
		return nil
	}
	err = e.client.Send([]string{email}, signInOtpSubject, html.String())
	if err != nil {
		return err
	}

	return nil
}

func (e *emailer) SendEmailUpdateOTP(email, otp string) error {
	html := bytes.Buffer{}
	err := e.templates[emailUpdateOtpTemplate].Execute(&html, &SendOTP{OTP: otp})
	if err != nil {
		return nil
	}
	err = e.client.Send([]string{email}, emailUpdateOtpSubject, html.String())
	if err != nil {
		return err
	}

	return nil
}
