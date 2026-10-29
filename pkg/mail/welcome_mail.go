package mail

import (
	"bytes"
	"context"

	"golang-starter-kit/resources/views/emails"
)

type WelcomeMail struct {
	Name string
}

func NewWelcomeMail(name string) *WelcomeMail {
	return &WelcomeMail{Name: name}
}

func (w *WelcomeMail) Envelope() Envelope {
	return Envelope{
		Subject: "Selamat Datang di Go Starter Kit!",
	}
}

func (w *WelcomeMail) Render(ctx context.Context) (string, error) {
	buf := new(bytes.Buffer)
	component := emails.Welcome(w.Name)
	if err := component.Render(ctx, buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
