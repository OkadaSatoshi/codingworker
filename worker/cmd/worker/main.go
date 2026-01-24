package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/OkadaSatoshi/codingworker/worker/internal/aider"
	"github.com/OkadaSatoshi/codingworker/worker/internal/config"
	"github.com/OkadaSatoshi/codingworker/worker/internal/github"
	"github.com/OkadaSatoshi/codingworker/worker/internal/sqs"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	// Initialize components
	sqsClient := sqs.NewClient(cfg.SQS)
	aiderRunner := aider.NewRunner(cfg.Aider)
	ghClient := github.NewClient(cfg.GitHub)

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
	slog.Info("Starting CodingWorker", "queue", cfg.SQS.QueueURL)
	if err := w.Run(ctx); err != nil {
		slog.Error("Worker error", "error", err)
		os.Exit(1)
	}

	slog.Info("Worker stopped")
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
	)

	// 2. Clone repository and create branch
	workDir, err := w.github.CloneAndBranch(ctx, msg.Repository, msg.IssueNumber)
	if err != nil {
		return err
	}
	defer os.RemoveAll(workDir)

	// 3. Run Aider to generate code
	if err := w.aider.Run(ctx, workDir, msg.Title, msg.Body); err != nil {
		return err
	}

	// 4. Push and create PR
	prURL, err := w.github.PushAndCreatePR(ctx, workDir, msg)
	if err != nil {
		return err
	}

	slog.Info("PR created", "url", prURL)

	// 5. Delete message from SQS
	if err := w.sqs.DeleteMessage(ctx, msg.ReceiptHandle); err != nil {
		return err
	}

	return nil
}
