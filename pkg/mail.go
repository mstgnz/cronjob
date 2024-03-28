package pkg

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

// https://github.com/mstgnz/gomail
type Mail struct {
	From        string
	Name        string
	Host        string
	Port        string
	User        string
	Pass        string
	Subject     string
	Content     string
	To          []string
	Cc          []string
	Bcc         []string
	Attachments map[string][]byte
}

// SetFrom sets the sender's email address
func (m *Mail) SetFrom(from string) *Mail {
	m.From = from
	return m
}

// SetName sets the sender's name
func (m *Mail) SetName(name string) *Mail {
	m.Name = name
	return m
}

// SetHost sets the SMTP server host
func (m *Mail) SetHost(host string) *Mail {
	m.Host = host
	return m
}

// SetPort sets the SMTP server port
func (m *Mail) SetPort(port string) *Mail {
	m.Port = port
	return m
}

// SetUser sets the SMTP server username
func (m *Mail) SetUser(user string) *Mail {
	m.User = user
	return m
}

// SetPass sets the SMTP server password
func (m *Mail) SetPass(pass string) *Mail {
	m.Pass = pass
	return m
}

// SetSubject sets the email subject
func (m *Mail) SetSubject(subject string) *Mail {
	m.Subject = subject
	return m
}

// SetContent sets the email content
func (m *Mail) SetContent(content string) *Mail {
	m.Content = content
	return m
}

// SetTo sets the email recipients
func (m *Mail) SetTo(to ...string) *Mail {
	m.To = to
	return m
}

// SetCc sets the email CC recipients
func (m *Mail) SetCc(cc ...string) *Mail {
	m.Cc = cc
	return m
}

// SetBcc sets the email BCC recipients
func (m *Mail) SetBcc(bcc ...string) *Mail {
	m.Bcc = bcc
	return m
}

// SetAttachment sets the email attachments
func (m *Mail) SetAttachment(attachments map[string][]byte) *Mail {
	m.Attachments = attachments
	return m
}

// SendText sends the email with plain text content
func (m *Mail) SendText() error {
	m.Subject = "text/plain: " + m.Subject
	return m.send()
}

// SendHTML sends the email with HTML content
func (m *Mail) SendHTML() error {
	m.Subject = "text/html: " + m.Subject
	return m.send()
}

// Send sends the email
func (m *Mail) send() error {
	if !m.validate() {
		return errors.New("missing parameter")
	}
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)

	// Create content
	var message strings.Builder
	message.WriteString(fmt.Sprintf("From: %s <%s>\n", m.Name, m.From))
	message.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ", ")))
	message.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.Cc, ", ")))
	message.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.Bcc, ", ")))
	message.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	message.WriteString("MIME-Version: 1.0\n")
	message.WriteString("Content-Type: multipart/mixed; boundary=BOUNDARY\n\n")

	// Add email content
	message.WriteString(fmt.Sprintf("--BOUNDARY\nContent-Type: text/plain\n\n%s\n\n", m.Content))

	// Add attachments
	for filename, data := range m.Attachments {
		message.WriteString(fmt.Sprintf("--BOUNDARY\nContent-Disposition: attachment; filename=\"%s\"\n", filename))
		message.WriteString("Content-Type: application/octet-stream\n\n")
		message.Write(data)
		message.WriteString("\n\n")
	}
	message.WriteString("--BOUNDARY--")

	// TLS configuration for connecting to SMTP server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.Host,
	}

	// Connection timeout setting
	dialer := &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// Connecting to the SMTP server
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, m.Host)
	if err != nil {
		return err
	}
	defer client.Quit()

	// Authentication information
	auth := smtp.PlainAuth("", m.User, m.Pass, m.Host)

	if err := client.Auth(auth); err != nil {
		return err
	}

	// Email sending process
	if err := client.Mail(m.From); err != nil {
		return err
	}

	allRecipients := append(append(m.To, m.Cc...), m.Bcc...)
	for _, recipient := range allRecipients {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	// Start writing email content
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	// Write email header and body
	_, err = w.Write([]byte(message.String()))
	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) validate() bool {
	if m.From == "" || m.Name == "" || m.Host == "" || m.Port == "" || m.User == "" || m.Pass == "" || m.Subject == "" || m.Content == "" || len(m.To) == 0 {
		return false
	}
	for _, email := range m.To {
		if !m.isEmailValid(email) {
			log.Printf("This email %s is not correct.\n", email)
			return false
		}
	}
	for _, email := range m.Cc {
		if !m.isEmailValid(email) {
			log.Printf("This email %s is not correct.\n", email)
			return false
		}
	}
	for _, email := range m.Bcc {
		if !m.isEmailValid(email) {
			log.Printf("This email %s is not correct.\n", email)
			return false
		}
	}
	return true
}

func (m *Mail) isEmailValid(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(regex).MatchString(email)
}
