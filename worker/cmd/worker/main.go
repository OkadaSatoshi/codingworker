package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/OkadaSatoshi/codingworker/worker/internal/aider"
	"github.com/OkadaSatoshi/codingworker/worker/internal/config"
	"github.com/OkadaSatoshi/codingworker/worker/internal/github"
	"github.com/OkadaSatoshi/codingworker/worker/internal/retry"
	"github.com/OkadaSatoshi/codingworker/worker/internal/sqs"
)

var (
	configPath  = flag.String("config", "configs/config.yaml", "Path to config file")
	testMessage = flag.String("test-message", "", "Path to test message JSON file (for local testing)")
	logLevel    = flag.String("log-level", "info", "Log level (debug, info, warn, error)")
)

func main() {
	flag.Parse()

	// Setup structured logging
	level := parseLogLevel(*logLevel)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize components
	sqsClient := sqs.NewClient(cfg.SQS)
	aiderRunner := aider.NewRunner(cfg.Aider)
	ghClient := github.NewClient(cfg.GitHub)

	// Inject test message if provided
	if *testMessage != "" {
		if err := injectTestMessage(sqsClient, *testMessage); err != nil {
			slog.Error("Failed to inject test message", "error", err)
			os.Exit(1)
		}
	}

	// Create worker
	w := &Worker{
		sqs:    sqsClient,
		aider:  aiderRunner,
		github: ghClient,
		config: cfg,
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		slog.Info("Received signal, shutting down", "signal", sig)
		cancel()
	}()

	// Start worker
	slog.Info("Starting CodingWorker",
		"config", *configPath,
		"mock_mode", cfg.SQS.UseMock,
		"queue_length", sqsClient.QueueLength(),
	)

	if err := w.Run(ctx); err != nil {
		slog.Error("Worker error", "error", err)
		os.Exit(1)
	}

	slog.Info("Worker stopped")
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func injectTestMessage(client *sqs.Client, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read test message file: %w", err)
	}

	var msg sqs.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("failed to parse test message: %w", err)
	}

	return client.InjectTestMessage(&msg)
}

type Worker struct {
	sqs    *sqs.Client
	aider  *aider.Runner
	github *github.Client
	config *config.Config
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := w.processNextMessage(ctx); err != nil {
				slog.Error("Failed to process message", "error", err)
			}
		}
	}
}

func (w *Worker) processNextMessage(ctx context.Context) error {
	// 1. Receive message from SQS
	msg, err := w.sqs.ReceiveMessage(ctx)
	if err != nil {
		return err
	}
	if msg == nil {
		return nil // No message available
	}

	slog.Info("Processing task",
		"issue_number", msg.IssueNumber,
		"repository", msg.Repository,
		"title", msg.Title,
	)

	// Execute with retry policy (uses config max_retries, fixed 10s backoff)
	policy := retry.NewPolicy(w.config.Worker.MaxRetries)
	var prURL string

	result := policy.Do(ctx, func() error {
		var err error
		prURL, err = w.processTask(ctx, msg)
		return err
	})

	if result.LastErr != nil {
		slog.Error("Task failed after retries",
			"issue_number", msg.IssueNumber,
			"attempts", result.Attempts,
			"error", result.LastErr,
		)

		// Post failure comment to Issue
		comment := w.buildFailureComment(result.LastErr, result.Attempts)
		if err := w.github.AddComment(ctx, msg.Repository, msg.IssueNumber, comment); err != nil {
			slog.Error("Failed to post failure comment", "error", err)
		}

		// Delete message from SQS (don't retry indefinitely)
		if err := w.sqs.DeleteMessage(ctx, msg.ReceiptHandle); err != nil {
			slog.Error("Failed to delete message after failure", "error", err)
		}

		return result.LastErr
	}

	slog.Info("PR created", "url", prURL)

	// Delete message from SQS
	if err := w.sqs.DeleteMessage(ctx, msg.ReceiptHandle); err != nil {
		return fmt.Errorf("message deletion failed: %w", err)
	}

	slog.Info("Task completed successfully",
		"issue_number", msg.IssueNumber,
		"pr_url", prURL,
		"attempts", result.Attempts,
	)

	return nil
}

// buildFailureComment creates a comment body for failed tasks
func (w *Worker) buildFailureComment(err error, attempts int) string {
	return fmt.Sprintf(`## ⚠️ CodingWorker: タスク処理に失敗しました

**試行回数**: %d
**エラー内容**:
`+"```"+`
%v
`+"```"+`

---
このコメントは CodingWorker によって自動生成されました。
`, attempts, err)
}

// processTask executes the actual work (clone, aider, push, PR)
func (w *Worker) processTask(ctx context.Context, msg *sqs.Message) (string, error) {
	// 2. Clone repository and create branch
	workDir, err := w.github.CloneAndBranch(ctx, msg.Repository, msg.IssueNumber)
	if err != nil {
		return "", fmt.Errorf("clone failed: %w", err)
	}
	defer os.RemoveAll(workDir)

	// 3. Run Aider to generate code (2-pass: implementation + tests)
	if err := w.aider.RunWithTests(ctx, workDir, msg.Title, msg.Body); err != nil {
		// Timeout errors are transient (can retry with fresh clone)
		if errors.Is(err, context.DeadlineExceeded) {
			return "", &retry.TransientError{Err: fmt.Errorf("aider timed out: %w", err)}
		}
		// Other Aider failures (after internal fix attempts) are permanent
		return "", &retry.PermanentError{Err: fmt.Errorf("aider failed: %w", err)}
	}

	// 4. Push and create PR
	prURL, err := w.github.PushAndCreatePR(ctx, workDir, msg)
	if err != nil {
		return "", fmt.Errorf("pr creation failed: %w", err)
	}

	return prURL, nil
}
