package sqs

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/OkadaSatoshi/codingworker/worker/internal/config"
)

// Message represents a task message from SQS
type Message struct {
	IssueNumber   int      `json:"issue_number"`
	Repository    string   `json:"repository"`
	Title         string   `json:"title"`
	Body          string   `json:"body"`
	Labels        []string `json:"labels"`
	CreatedAt     string   `json:"created_at"`
	ReceiptHandle string   `json:"-"`
}

// Client handles SQS operations
type Client struct {
	config  config.SQSConfig
	useMock bool
}

// NewClient creates a new SQS client
func NewClient(cfg config.SQSConfig) *Client {
	return &Client{
		config:  cfg,
		useMock: cfg.UseMock,
	}
}

// ReceiveMessage receives a message from SQS (or mock)
func (c *Client) ReceiveMessage(ctx context.Context) (*Message, error) {
	if c.useMock {
		return c.receiveMockMessage(ctx)
	}
	return c.receiveFromAWS(ctx)
}

// DeleteMessage deletes a message from SQS
func (c *Client) DeleteMessage(ctx context.Context, receiptHandle string) error {
	if c.useMock {
		slog.Info("Mock: Deleting message", "receipt_handle", receiptHandle)
		return nil
	}
	return c.deleteFromAWS(ctx, receiptHandle)
}

// Mock implementation for development
func (c *Client) receiveMockMessage(ctx context.Context) (*Message, error) {
	// Simulate long polling delay
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(time.Duration(c.config.WaitTimeSeconds) * time.Second):
		return nil, nil // No message
	}
}

// AWS implementation (placeholder)
func (c *Client) receiveFromAWS(ctx context.Context) (*Message, error) {
	// TODO: Implement real AWS SQS receive
	// This will be implemented when AWS infrastructure is ready
	slog.Warn("AWS SQS not implemented, falling back to mock")
	return c.receiveMockMessage(ctx)
}

func (c *Client) deleteFromAWS(ctx context.Context, receiptHandle string) error {
	// TODO: Implement real AWS SQS delete
	slog.Warn("AWS SQS delete not implemented")
	return nil
}

// SendTestMessage sends a test message (for development/testing)
func (c *Client) SendTestMessage(ctx context.Context, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	slog.Info("Test message created", "message", string(data))
	return nil
}
