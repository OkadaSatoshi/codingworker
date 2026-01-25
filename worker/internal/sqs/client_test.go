package sqs

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/OkadaSatoshi/codingworker/worker/internal/config"
)

func TestNewClient(t *testing.T) {
	cfg := config.SQSConfig{
		UseMock:         true,
		WaitTimeSeconds: 5,
	}

	client := NewClient(cfg)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}
	if !client.useMock {
		t.Error("expected useMock to be true")
	}
	if client.mockQueue == nil {
		t.Error("expected mockQueue to be initialized")
	}
}

func TestInjectTestMessage(t *testing.T) {
	cfg := config.SQSConfig{UseMock: true}
	client := NewClient(cfg)

	msg := &Message{
		IssueNumber: 42,
		Repository:  "test/repo",
		Title:       "Test task",
		Body:        "Test body",
	}

	err := client.InjectTestMessage(msg)
	if err != nil {
		t.Fatalf("InjectTestMessage failed: %v", err)
	}

	if client.QueueLength() != 1 {
		t.Errorf("expected queue length 1, got %d", client.QueueLength())
	}

	// Check that receipt handle was generated
	if msg.ReceiptHandle == "" {
		t.Error("expected ReceiptHandle to be generated")
	}

	// Check that created_at was set
	if msg.CreatedAt == "" {
		t.Error("expected CreatedAt to be set")
	}
}

func TestInjectTestMessage_NotMockMode(t *testing.T) {
	cfg := config.SQSConfig{UseMock: false}
	client := NewClient(cfg)

	msg := &Message{
		IssueNumber: 1,
		Repository:  "test/repo",
		Title:       "Test",
	}

	err := client.InjectTestMessage(msg)
	if err == nil {
		t.Error("expected error when not in mock mode")
	}
}

func TestInjectTestMessageFromJSON(t *testing.T) {
	cfg := config.SQSConfig{UseMock: true}
	client := NewClient(cfg)

	jsonStr := `{
		"issue_number": 123,
		"repository": "owner/repo",
		"title": "JSON test",
		"body": "Body from JSON"
	}`

	err := client.InjectTestMessageFromJSON(jsonStr)
	if err != nil {
		t.Fatalf("InjectTestMessageFromJSON failed: %v", err)
	}

	if client.QueueLength() != 1 {
		t.Errorf("expected queue length 1, got %d", client.QueueLength())
	}
}

func TestInjectTestMessageFromJSON_InvalidJSON(t *testing.T) {
	cfg := config.SQSConfig{UseMock: true}
	client := NewClient(cfg)

	err := client.InjectTestMessageFromJSON("invalid json")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestReceiveMessage_Mock(t *testing.T) {
	cfg := config.SQSConfig{
		UseMock:         true,
		WaitTimeSeconds: 1,
	}
	client := NewClient(cfg)

	// Inject a message first
	msg := &Message{
		IssueNumber: 99,
		Repository:  "test/repo",
		Title:       "Receive test",
	}
	if err := client.InjectTestMessage(msg); err != nil {
		t.Fatalf("failed to inject: %v", err)
	}

	// Receive it
	ctx := context.Background()
	received, err := client.ReceiveMessage(ctx)
	if err != nil {
		t.Fatalf("ReceiveMessage failed: %v", err)
	}

	if received == nil {
		t.Fatal("expected to receive a message")
	}
	if received.IssueNumber != 99 {
		t.Errorf("expected issue 99, got %d", received.IssueNumber)
	}
	if received.Repository != "test/repo" {
		t.Errorf("expected 'test/repo', got %s", received.Repository)
	}
}

func TestReceiveMessage_NoMessage(t *testing.T) {
	cfg := config.SQSConfig{
		UseMock:         true,
		WaitTimeSeconds: 1, // Short timeout for test
	}
	client := NewClient(cfg)

	ctx := context.Background()
	received, err := client.ReceiveMessage(ctx)
	if err != nil {
		t.Fatalf("ReceiveMessage failed: %v", err)
	}

	if received != nil {
		t.Error("expected nil when no message available")
	}
}

func TestReceiveMessage_ContextCancelled(t *testing.T) {
	cfg := config.SQSConfig{
		UseMock:         true,
		WaitTimeSeconds: 10, // Long timeout
	}
	client := NewClient(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	received, err := client.ReceiveMessage(ctx)
	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if received != nil {
		t.Error("expected nil message on cancellation")
	}
}

func TestDeleteMessage_Mock(t *testing.T) {
	cfg := config.SQSConfig{UseMock: true}
	client := NewClient(cfg)

	err := client.DeleteMessage(context.Background(), "mock-receipt-123")
	if err != nil {
		t.Errorf("DeleteMessage failed: %v", err)
	}
}

func TestQueueLength(t *testing.T) {
	cfg := config.SQSConfig{UseMock: true}
	client := NewClient(cfg)

	if client.QueueLength() != 0 {
		t.Errorf("expected 0, got %d", client.QueueLength())
	}

	// Add messages
	for i := 1; i <= 3; i++ {
		msg := &Message{IssueNumber: i, Repository: "test/repo", Title: "Test"}
		client.InjectTestMessage(msg)
	}

	if client.QueueLength() != 3 {
		t.Errorf("expected 3, got %d", client.QueueLength())
	}
}

func TestCreateTestMessage(t *testing.T) {
	msg := CreateTestMessage("owner/repo", 42, "Task title", "Task body")

	if msg.IssueNumber != 42 {
		t.Errorf("expected issue 42, got %d", msg.IssueNumber)
	}
	if msg.Repository != "owner/repo" {
		t.Errorf("expected 'owner/repo', got %s", msg.Repository)
	}
	if msg.Title != "Task title" {
		t.Errorf("expected 'Task title', got %s", msg.Title)
	}
	if msg.Body != "Task body" {
		t.Errorf("expected 'Task body', got %s", msg.Body)
	}
	if len(msg.Labels) != 1 || msg.Labels[0] != LabelTrigger {
		t.Errorf("expected [%s], got %v", LabelTrigger, msg.Labels)
	}
	if msg.CreatedAt == "" {
		t.Error("expected CreatedAt to be set")
	}
}

func TestMessage_JSONMarshal(t *testing.T) {
	msg := &Message{
		IssueNumber:   1,
		Repository:    "test/repo",
		Title:         "Test",
		Body:          "Body",
		Labels:        []string{"ai-task"},
		CreatedAt:     "2024-01-01T00:00:00Z",
		ReceiptHandle: "should-not-appear",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// ReceiptHandle should be omitted (json:"-")
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if _, exists := parsed["ReceiptHandle"]; exists {
		t.Error("ReceiptHandle should not be in JSON output")
	}
}
