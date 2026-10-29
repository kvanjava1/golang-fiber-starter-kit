package mail_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/pkg/mail"
	"golang-starter-kit/pkg/mail/drivers"
)

func TestMailable_WelcomeMailRender(t *testing.T) {
	welcome := mail.NewWelcomeMail("Budi")

	env := welcome.Envelope()
	assert.Equal(t, "Selamat Datang di Go Starter Kit!", env.Subject)

	html, err := welcome.Render(context.Background())
	require.NoError(t, err)
	assert.Contains(t, html, "Selamat Datang, Budi!")
	assert.Contains(t, html, "<!doctype html>")
}

func TestMailClient_SendSynchronous(t *testing.T) {
	logDriver := drivers.NewLogDriver()
	client := mail.NewClient(logDriver, nil, "Go Kit", "no-reply@starter.kit")

	welcome := mail.NewWelcomeMail("Siti")
	err := client.To("siti@example.com").Send(context.Background(), welcome)
	require.NoError(t, err)

	require.NotNil(t, logDriver.LastDispatched)
	assert.Equal(t, []string{"siti@example.com"}, logDriver.LastDispatched.Envelope.To)
	assert.Equal(t, "Go Kit <no-reply@starter.kit>", logDriver.LastDispatched.Envelope.From)
	assert.Contains(t, logDriver.LastDispatched.HTML, "Selamat Datang, Siti!")
}
