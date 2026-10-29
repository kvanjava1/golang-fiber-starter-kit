package mail

import (
	"context"
	"fmt"

	"golang-starter-kit/pkg/queue"
)

type MailClient struct {
	driver    Driver
	queue     *queue.Client
	fromName  string
	fromEmail string
}

func NewClient(driver Driver, q *queue.Client, fromName, fromEmail string) *MailClient {
	return &MailClient{
		driver:    driver,
		queue:     q,
		fromName:  fromName,
		fromEmail: fromEmail,
	}
}

func (c *MailClient) Driver() Driver {
	return c.driver
}

func (c *MailClient) To(recipient string) *PendingMail {
	return &PendingMail{
		client:    c,
		recipient: recipient,
	}
}

type PendingMail struct {
	client    *MailClient
	recipient string
}

// Send mengirim email secara sinkron via driver yang aktif
func (p *PendingMail) Send(ctx context.Context, m Mailable) error {
	html, err := m.Render(ctx)
	if err != nil {
		return fmt.Errorf("failed to render mailable html: %w", err)
	}

	env := m.Envelope()
	if env.From == "" {
		if p.client.fromName != "" {
			env.From = fmt.Sprintf("%s <%s>", p.client.fromName, p.client.fromEmail)
		} else {
			env.From = p.client.fromEmail
		}
	}
	env.To = []string{p.recipient}

	return p.client.driver.Send(ctx, &Message{
		Envelope: env,
		HTML:     html,
	})
}

// Queue merender HTML terlebih dahulu lalu meng-enqueue task ke Asynq Redis
func (p *PendingMail) Queue(ctx context.Context, m Mailable) error {
	if p.client.queue == nil {
		return fmt.Errorf("queue client is not configured, cannot enqueue email")
	}

	html, err := m.Render(ctx)
	if err != nil {
		return fmt.Errorf("failed to render mailable html: %w", err)
	}

	env := m.Envelope()
	return p.client.queue.EnqueueSendEmail(ctx, p.recipient, env.Subject, html)
}
