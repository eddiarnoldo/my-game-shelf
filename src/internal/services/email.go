package services

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

type EmailService interface {
	SendInvite(toEmail, inviteCode string) error
}

type SMTPEmailService struct {
	host       string
	port       string
	user       string
	password   string
	from       string
	baseURL    string
	implicitTLS bool // true for port 465, false for 587/STARTTLS
}

// NewMailpitEmailService creates an email service pointing at the local Mailpit
// dev server. No credentials needed — emails are captured at http://localhost:8025.
func NewMailpitEmailService(baseURL string) *SMTPEmailService {
	return &SMTPEmailService{
		host:    "localhost",
		port:    "1025",
		from:    "noreply@mygameshelf.local",
		baseURL: baseURL,
	}
}

// NewGmailEmailService creates an email service using Gmail SMTP on port 465 (implicit TLS).
// appPassword must be a Gmail App Password (not your account password).
// Generate one at: myaccount.google.com → Security → App Passwords
func NewGmailEmailService(from, appPassword, baseURL string) *SMTPEmailService {
	return &SMTPEmailService{
		host:        "smtp.gmail.com",
		port:        "465",
		user:        from,
		password:    appPassword,
		from:        from,
		baseURL:     baseURL,
		implicitTLS: true,
	}
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

	if s.implicitTLS {
		return s.sendWithImplicitTLS(addr, msg, toEmail)
	}

	// STARTTLS (port 587) or plaintext (Mailpit)
	var auth smtp.Auth
	if s.user != "" {
		auth = smtp.PlainAuth("", s.user, s.password, s.host)
	}
	return smtp.SendMail(addr, auth, s.from, []string{toEmail}, []byte(msg))
}

// sendWithImplicitTLS dials port 465 with TLS from the start (no STARTTLS negotiation).
func (s *SMTPEmailService) sendWithImplicitTLS(addr, msg, toEmail string) error {
	tlsConfig := &tls.Config{ServerName: s.host}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}

	host, _, _ := net.SplitHostPort(addr)
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", s.user, s.password, s.host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}
	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("smtp rcpt: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := fmt.Fprint(w, msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	return w.Close()
}
