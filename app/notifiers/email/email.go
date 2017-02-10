package email

import (
	"crypto/tls"
	"log"

	"github.com/allen13/golerta/app/models"
	"gopkg.in/gomail.v2"
)

type Email struct {
	Addresses     []string `toml:"email"`
	EnabledField  bool     `toml:"enabled"`
	SmtpServer    string   `toml:"smtp_server"`
	SmtpPort      int      `toml:"smtp_port"`
	SmtpUser      string   `toml:"smtp_user"`
	SmtpPassword  string   `toml:"smtp_password"`
	SkipSslVerify bool     `toml:"skip_ssl_verify"`
	From          string   `toml:"from"`
	GolertaUrl    string   `toml:"golerta_url"`
}

func (em *Email) Init() error {
	return nil
}

func (em *Email) Enabled() bool {
	return em.EnabledField
}

func (em *Email) CreateEmailEvent(eventType string, alert models.Alert) error {
	d := gomail.NewDialer(em.SmtpServer, em.SmtpPort, em.SmtpUser, em.SmtpPassword)
	if em.SkipSslVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: em.SkipSslVerify}
	}
	s, err := d.Dial()
	if err != nil {
		return (err)
	}

	m := gomail.NewMessage()

	for _, mail := range em.Addresses {
		m.SetHeader("From", em.From)
		m.SetHeader("To", mail)
		m.SetHeader("Subject", eventType+" "+alert.Severity+" "+alert.Resource+" "+alert.Environment)
		m.SetBody("text/plain", "Alert URL:\n"+em.GolertaUrl+alert.Id+
			"\n\nAlert Status:\n"+alert.Status+
			"\n\nAlert Comment:\n"+alert.History[0].Text+
			"\n\nAlert Info:\n"+alert.String())

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send mail to %q: %v", mail, err)
		}
		m.Reset()
	}

	return (err)
}

func (em *Email) Trigger(alert models.Alert) error {
	return em.CreateEmailEvent("trigger", alert)
}

func (em *Email) Acknowledge(alert models.Alert) error {
	return em.CreateEmailEvent("acknowledge", alert)
}

func (em *Email) Resolve(alert models.Alert) error {
	return em.CreateEmailEvent("resolve", alert)
}
