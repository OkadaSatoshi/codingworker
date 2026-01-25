package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
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

// Label constants
const (
	LabelTrigger = "ai-task"        // Triggers worker processing
	LabelFailed  = "ai-task-failed" // Added on failure
	LabelDone    = "ai-task-done"   // Added on success (ai-task removed)
)

// Client handles SQS operations
type Client struct {
	config    config.SQSConfig
	useMock   bool
	mockQueue chan *Message
	mu        sync.Mutex
}

// NewClient creates a new SQS client
func NewClient(cfg config.SQSConfig) *Client {
	return &Client{
		config:    cfg,
		useMock:   cfg.UseMock,
		mockQueue: make(chan *Message, 100), // Buffer for test messages
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
		slog.Info("Mock: Message deleted", "receipt_handle", receiptHandle)
		return nil
	}
	return c.deleteFromAWS(ctx, receiptHandle)
}

// InjectTestMessage adds a test message to the mock queue
func (c *Client) InjectTestMessage(msg *Message) error {
	if !c.useMock {
		return fmt.Errorf("cannot inject test message when not in mock mode")
	}

	// Generate a mock receipt handle
	if msg.ReceiptHandle == "" {
		msg.ReceiptHandle = fmt.Sprintf("mock-receipt-%d-%d", msg.IssueNumber, time.Now().UnixNano())
	}

	// Set created_at if not set
	if msg.CreatedAt == "" {
		msg.CreatedAt = time.Now().Format(time.RFC3339)
	}

	select {
	case c.mockQueue <- msg:
		slog.Info("Test message injected",
			"issue_number", msg.IssueNumber,
			"repository", msg.Repository,
			"title", msg.Title,
		)
		return nil
	default:
		return fmt.Errorf("mock queue is full")
	}
}

// InjectTestMessageFromJSON adds a test message from JSON string
func (c *Client) InjectTestMessageFromJSON(jsonStr string) error {
	var msg Message
	if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return c.InjectTestMessage(&msg)
}

// QueueLength returns the number of messages in the mock queue
func (c *Client) QueueLength() int {
	return len(c.mockQueue)
}

// Mock implementation for development
func (c *Client) receiveMockMessage(ctx context.Context) (*Message, error) {
	// First check if there are any test messages in the queue
	select {
	case msg := <-c.mockQueue:
		slog.Info("Mock: Received message from test queue",
			"issue_number", msg.IssueNumber,
			"repository", msg.Repository,
		)
		return msg, nil
	default:
		// No message in queue, do long polling simulation
	}

	// Simulate long polling delay
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-c.mockQueue:
		slog.Info("Mock: Received message from test queue",
			"issue_number", msg.IssueNumber,
			"repository", msg.Repository,
		)
		return msg, nil
	case <-time.After(time.Duration(c.config.WaitTimeSeconds) * time.Second):
		return nil, nil // No message after timeout
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

// CreateTestMessage is a helper to create a test message
func CreateTestMessage(repo string, issueNumber int, title, body string) *Message {
	return &Message{
		IssueNumber: issueNumber,
		Repository:  repo,
		Title:       title,
		Body:        body,
		Labels:      []string{LabelTrigger},
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
}
