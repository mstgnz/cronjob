package config

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

type Mail struct {
	from        string
	name        string
	host        string
	port        string
	user        string
	pass        string
	subject     string
	content     string
	to          []string
	cc          []string
	bcc         []string
	attachments map[string][]byte
}

// SetFrom sets the sender's email address
func (m *Mail) SetFrom(from string) *Mail {
	m.from = from
	return m
}

// SetName sets the sender's name
func (m *Mail) SetName(name string) *Mail {
	m.name = name
	return m
}

// SetHost sets the SMTP server host
func (m *Mail) SetHost(host string) *Mail {
	m.host = host
	return m
}

// SetPort sets the SMTP server port
func (m *Mail) SetPort(port string) *Mail {
	m.port = port
	return m
}

// SetUser sets the SMTP server username
func (m *Mail) SetUser(user string) *Mail {
	m.user = user
	return m
}

// SetPass sets the SMTP server password
func (m *Mail) SetPass(pass string) *Mail {
	m.pass = pass
	return m
}

// SetSubject sets the email subject
func (m *Mail) SetSubject(subject string) *Mail {
	m.subject = subject
	return m
}

// SetContent sets the email content
func (m *Mail) SetContent(content string) *Mail {
	m.content = content
	return m
}

// SetTo sets the email recipients
func (m *Mail) SetTo(to ...string) *Mail {
	m.to = to
	return m
}

// SetCc sets the email CC recipients
func (m *Mail) SetCc(cc ...string) *Mail {
	m.cc = cc
	return m
}

// SetBcc sets the email BCC recipients
func (m *Mail) SetBcc(bcc ...string) *Mail {
	m.bcc = bcc
	return m
}

// SetAttachment sets the email attachments
func (m *Mail) SetAttachment(attachments map[string][]byte) *Mail {
	m.attachments = attachments
	return m
}

// SendText sends the email with plain text content
func (m *Mail) SendText() error {
	m.subject = "text/plain: " + m.subject
	return m.send()
}

// SendHTML sends the email with HTML content
func (m *Mail) SendHTML() error {
	m.subject = "text/html: " + m.subject
	return m.send()
}

// Send sends the email
func (m *Mail) send() error {
	if !m.validate() {
		return errors.New("missing parameter")
	}
	addr := fmt.Sprintf("%s:%s", m.host, m.port)

	// Create content
	var message strings.Builder
	message.WriteString(fmt.Sprintf("From: %s <%s>\n", m.name, m.from))
	message.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.to, ", ")))
	message.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.cc, ", ")))
	message.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.bcc, ", ")))
	message.WriteString(fmt.Sprintf("Subject: %s\n", m.subject))
	message.WriteString("MIME-Version: 1.0\n")
	message.WriteString("Content-Type: multipart/mixed; boundary=BOUNDARY\n\n")

	// Add email content
	message.WriteString(fmt.Sprintf("--BOUNDARY\nContent-Type: text/plain\n\n%s\n\n", m.content))

	// Add attachments
	for filename, data := range m.attachments {
		message.WriteString(fmt.Sprintf("--BOUNDARY\nContent-Disposition: attachment; filename=\"%s\"\n", filename))
		message.WriteString("Content-Type: application/octet-stream\n\n")
		message.Write(data)
		message.WriteString("\n\n")
	}
	message.WriteString("--BOUNDARY--")

	// TLS configuration for connecting to SMTP server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.host,
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

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return err
	}
	defer client.Quit()

	// Authentication information
	auth := smtp.PlainAuth("", m.user, m.pass, m.host)

	if err := client.Auth(auth); err != nil {
		return err
	}

	// Email sending process
	if err := client.Mail(m.from); err != nil {
		return err
	}

	allRecipients := append(append(m.to, m.cc...), m.bcc...)
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
	if m.from == "" || m.name == "" || m.host == "" || m.port == "" || m.user == "" || m.pass == "" || m.subject == "" || m.content == "" || len(m.to) == 0 {
		return false
	}
	for _, email := range m.to {
		if !m.isEmailValid(email) {
			log.Printf("This email %s is not correct.\n", email)
			return false
		}
	}
	for _, email := range m.cc {
		if !m.isEmailValid(email) {
			log.Printf("This email %s is not correct.\n", email)
			return false
		}
	}
	for _, email := range m.bcc {
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
