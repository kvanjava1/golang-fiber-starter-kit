package queue

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

type Client struct {
	client *asynq.Client
}

// NewClient menginisialisasi asynq.Client menggunakan parser Redis URI resmi
func NewClient(redisURL string) (*Client, error) {
	redisOpt, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis uri for asynq: %w", err)
	}

	asynqClient := asynq.NewClient(redisOpt)
	return &Client{client: asynqClient}, nil
}

func (c *Client) Enqueue(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	if c.client == nil {
		return nil, fmt.Errorf("queue client is not connected")
	}
	return c.client.EnqueueContext(ctx, task, opts...)
}

func (c *Client) EnqueueSendEmail(ctx context.Context, to, subject, html string, opts ...asynq.Option) error {
	task, err := NewEmailSendTask(ctx, to, subject, html)
	if err != nil {
		return err
	}

	// Default retry: 3 kali dengan max delay
	defaultOpts := []asynq.Option{
		asynq.MaxRetry(3),
		asynq.Queue("default"),
	}
	defaultOpts = append(defaultOpts, opts...)

	_, err = c.Enqueue(ctx, task, defaultOpts...)
	return err
}

func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
