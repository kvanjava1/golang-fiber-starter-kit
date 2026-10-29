package mail

import "context"

type Envelope struct {
	Subject string
	From    string
	To      []string
}

type Message struct {
	Envelope Envelope
	HTML     string
}

type Mailable interface {
	Envelope() Envelope
	Render(ctx context.Context) (string, error)
}
