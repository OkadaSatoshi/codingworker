package aider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"time"

	"github.com/OkadaSatoshi/codingworker/worker/internal/config"
)

// Runner executes Aider commands
type Runner struct {
	config config.AiderConfig
}

// NewRunner creates a new Aider runner
func NewRunner(cfg config.AiderConfig) *Runner {
	return &Runner{
		config: cfg,
	}
}

// Run executes Aider with the given task, with model fallback on timeout
func (r *Runner) Run(ctx context.Context, workDir, title, body string) error {
	var lastErr error

	for i, model := range r.config.Models {
		err := r.runWithModel(ctx, workDir, title, body, model)
		if err == nil {
			return nil // 成功
		}
		lastErr = err

		// タイムアウトの場合のみフォールバック
		if !errors.Is(err, context.DeadlineExceeded) {
			// タイムアウト以外のエラー → フォールバックせず終了
			return err
		}

		if i < len(r.config.Models)-1 {
			slog.Warn("Model timed out, trying next",
				"failed_model", model.Name,
				"next_model", r.config.Models[i+1].Name,
			)
		}
	}
	return fmt.Errorf("all models timed out: %w", lastErr)
}

// runWithModel executes Aider with a specific model
func (r *Runner) runWithModel(ctx context.Context, workDir, title, body string, model config.ModelConfig) error {
	prompt := r.buildPrompt(title, body)

	slog.Info("Running Aider",
		"work_dir", workDir,
		"model", model.Name,
		"timeout_seconds", model.Timeout,
		"prompt_length", len(prompt),
	)

	// Create context with model-specific timeout
	timeout := time.Duration(model.Timeout) * time.Second
	modelCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Build Aider command
	args := []string{
		"--model", model.Name,
		"--yes",           // Auto-confirm changes
		"--no-auto-lint",  // Skip auto-linting
		"--map-tokens", strconv.Itoa(r.config.MapTokens),
		"--message", prompt,
	}

	cmd := exec.CommandContext(modelCtx, r.config.BinPath, args...)
	cmd.Dir = workDir

	// Capture output
	output, err := cmd.CombinedOutput()

	// Check for timeout
	if modelCtx.Err() == context.DeadlineExceeded {
		slog.Warn("Aider execution timed out",
			"model", model.Name,
			"timeout_seconds", model.Timeout,
		)
		return context.DeadlineExceeded
	}

	if err != nil {
		slog.Error("Aider execution failed",
			"model", model.Name,
			"error", err,
			"output", string(output),
		)
		return fmt.Errorf("aider execution failed: %w", err)
	}

	slog.Info("Aider completed successfully",
		"model", model.Name,
		"output_length", len(output),
	)

	return nil
}

// buildPrompt creates a prompt from issue title and body
func (r *Runner) buildPrompt(title, body string) string {
	if body == "" {
		return title
	}
	return fmt.Sprintf("%s\n\n%s", title, body)
}

// CheckInstallation verifies Aider is installed and working
func (r *Runner) CheckInstallation(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, r.config.BinPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("aider not found or not working: %w", err)
	}
	slog.Info("Aider installation verified", "version", string(output))
	return nil
}
