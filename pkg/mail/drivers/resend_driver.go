package drivers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang-starter-kit/pkg/mail"
)

type ResendDriver struct {
	apiKey     string
	httpClient *http.Client
}

func NewResendDriver(apiKey string) *ResendDriver {
	return &ResendDriver{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (d *ResendDriver) Name() string {
	return "resend"
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

func (d *ResendDriver) Send(ctx context.Context, msg *mail.Message) error {
	if d.apiKey == "" {
		return fmt.Errorf("resend api key is missing")
	}

	payload := resendRequest{
		From:    msg.Envelope.From,
		To:      msg.Envelope.To,
		Subject: msg.Envelope.Subject,
		HTML:    msg.HTML,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal resend request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create resend http request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
