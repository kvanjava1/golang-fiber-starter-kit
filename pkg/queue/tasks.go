package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"

	"golang-starter-kit/pkg/logger"
)

const (
	TypeEmailSend = "email:send"
)

type EmailSendPayload struct {
	To            string `json:"to"`
	Subject       string `json:"subject"`
	HTML          string `json:"html"`
	CorrelationID string `json:"correlationId,omitempty"`
}

// NewEmailSendTask membuat asynq.Task untuk pengiriman email dengan membawa correlation ID
func NewEmailSendTask(ctx context.Context, to, subject, html string) (*asynq.Task, error) {
	traceID := logger.CorrelationIDFromContext(ctx)

	payload, err := json.Marshal(EmailSendPayload{
		To:            to,
		Subject:       subject,
		HTML:          html,
		CorrelationID: traceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email task payload: %w", err)
	}

	return asynq.NewTask(TypeEmailSend, payload), nil
}
