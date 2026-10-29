package drivers

import (
	"context"
	"fmt"

	gomail "github.com/wneessen/go-mail"

	"golang-starter-kit/pkg/mail"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type SMTPDriver struct {
	cfg SMTPConfig
}

func NewSMTPDriver(cfg SMTPConfig) *SMTPDriver {
	return &SMTPDriver{cfg: cfg}
}

func (d *SMTPDriver) Name() string {
	return "smtp"
}

func (d *SMTPDriver) Send(ctx context.Context, msg *mail.Message) error {
	m := gomail.NewMsg()
	if err := m.From(msg.Envelope.From); err != nil {
		return fmt.Errorf("failed to set from: %w", err)
	}

	for _, to := range msg.Envelope.To {
		if err := m.AddTo(to); err != nil {
			return fmt.Errorf("failed to add recipient %s: %w", to, err)
		}
	}

	m.Subject(msg.Envelope.Subject)
	m.SetBodyString(gomail.TypeTextHTML, msg.HTML)

	opts := []gomail.Option{
		gomail.WithPort(d.cfg.Port),
	}

	if d.cfg.Username != "" && d.cfg.Password != "" {
		opts = append(opts,
			gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
			gomail.WithUsername(d.cfg.Username),
			gomail.WithPassword(d.cfg.Password),
		)
	}

	client, err := gomail.NewClient(d.cfg.Host, opts...)
	if err != nil {
		return fmt.Errorf("failed to create smtp client: %w", err)
	}

	return client.DialAndSendWithContext(ctx, m)
}
