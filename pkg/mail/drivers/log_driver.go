package drivers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"golang-starter-kit/pkg/mail"
)

type LogDriver struct {
	LastDispatched *mail.Message
}

func NewLogDriver() *LogDriver {
	return &LogDriver{}
}

func (d *LogDriver) Name() string {
	return "log"
}

func (d *LogDriver) Send(ctx context.Context, msg *mail.Message) error {
	d.LastDispatched = msg

	toStr := strings.Join(msg.Envelope.To, ", ")
	banner := "==================== [EMAIL DISPATCHED (LOG DRIVER)] ===================="

	fmt.Printf("\n%s\nFrom:    %s\nTo:      %s\nSubject: %s\nContent:\n%s\n%s\n\n",
		banner, msg.Envelope.From, toStr, msg.Envelope.Subject, msg.HTML, banner)

	slog.InfoContext(ctx, "Email dispatched via log driver",
		"to", toStr,
		"from", msg.Envelope.From,
		"subject", msg.Envelope.Subject,
	)
	return nil
}
