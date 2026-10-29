package queue_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"golang-starter-kit/pkg/logger"
	"golang-starter-kit/pkg/queue"
)

func TestQueue_NewEmailSendTask(t *testing.T) {
	ctx := logger.WithCorrelationID(context.Background(), "test-corr-id-777")
	task, err := queue.NewEmailSendTask(ctx, "target@test.com", "Halo", "<p>Isi</p>")
	require.NoError(t, err)
	assert.Equal(t, queue.TypeEmailSend, task.Type())

	var payload queue.EmailSendPayload
	err = json.Unmarshal(task.Payload(), &payload)
	require.NoError(t, err)

	assert.Equal(t, "target@test.com", payload.To)
	assert.Equal(t, "Halo", payload.Subject)
	assert.Equal(t, "<p>Isi</p>", payload.HTML)
	assert.Equal(t, "test-corr-id-777", payload.CorrelationID)
}
