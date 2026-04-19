package services

import (
	"fmt"
	"net/smtp"
)

type EmailService interface {
	SendInvite(toEmail, inviteCode string) error
}

type SMTPEmailService struct {
	host     string
	port     string
	user     string
	password string
	from     string
	baseURL  string
}

func NewSMTPEmailService(host, port, user, password, from, baseURL string) *SMTPEmailService {
	return &SMTPEmailService{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		from:     from,
		baseURL:  baseURL,
	}
}

func (s *SMTPEmailService) SendInvite(toEmail, inviteCode string) error {
	inviteLink := fmt.Sprintf("%s/join?code=%s", s.baseURL, inviteCode)

	subject := "You've been invited to My Game Shelf"
	body := fmt.Sprintf(
		"You have been invited to join My Game Shelf.\n\nClick the link below to set up your account:\n\n%s\n\nThis invite expires in 7 days.\n",
		inviteLink,
	)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		s.from, toEmail, subject, body,
	)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	// Skip auth when no credentials are set (e.g. local Mailpit dev server)
	var auth smtp.Auth
	if s.user != "" {
		auth = smtp.PlainAuth("", s.user, s.password, s.host)
	}

	return smtp.SendMail(addr, auth, s.from, []string{toEmail}, []byte(msg))
}
