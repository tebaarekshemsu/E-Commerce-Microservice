package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"

	"notification-service/internal/config"
	"notification-service/internal/models"
)

type EmailService struct {
	config    config.SMTPConfig
	templates map[string]*template.Template
}

func NewEmailService(cfg config.SMTPConfig) *EmailService {
	svc := &EmailService{
		config:    cfg,
		templates: make(map[string]*template.Template),
	}
	svc.loadTemplates()
	return svc
}

func (s *EmailService) loadTemplates() {
	// Load email templates
	templateFiles := []string{
		"order_confirmation",
		"order_shipped",
		"order_delivered",
		"welcome",
		"password_reset",
		"low_stock_alert",
	}

	for _, name := range templateFiles {
		tmpl, err := template.ParseFiles(
			fmt.Sprintf("templates/email/%s.html", name),
		)
		if err != nil {
			continue
		}
		s.templates[name] = tmpl
	}
}

func (s *EmailService) Send(notification *models.Notification) error {
	// Build email content
	var body bytes.Buffer

	if tmpl, ok := s.templates[notification.TemplateID]; ok {
		if err := tmpl.Execute(&body, notification.Metadata); err != nil {
			return fmt.Errorf("template execution failed: %w", err)
		}
	} else {
		body.WriteString(notification.Content)
	}

	// Build email headers
	headers := make(map[string]string)
	headers["From"] = s.config.From
	headers["To"] = notification.Recipient
	headers["Subject"] = notification.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.Write(body.Bytes())

	// Setup TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.config.Host,
	}

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("SMTP client creation failed: %w", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err := client.Mail(s.config.From); err != nil {
		return fmt.Errorf("SMTP mail from failed: %w", err)
	}

	if err := client.Rcpt(notification.Recipient); err != nil {
		return fmt.Errorf("SMTP rcpt to failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP data failed: %w", err)
	}

	if _, err := w.Write(msg.Bytes()); err != nil {
		return fmt.Errorf("SMTP write failed: %w", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("SMTP close failed: %w", err)
	}

	return client.Quit()
}
