package aider

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
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

// Run executes Aider with the given task
func (r *Runner) Run(ctx context.Context, workDir, title, body string) error {
	// Build the prompt from issue title and body
	prompt := r.buildPrompt(title, body)

	slog.Info("Running Aider",
		"work_dir", workDir,
		"model", r.config.Model,
		"prompt_length", len(prompt),
	)

	// Create context with timeout
	timeout := time.Duration(r.config.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Build Aider command
	args := []string{
		"--model", r.config.Model,
		"--yes",           // Auto-confirm changes
		"--no-auto-lint",  // Skip auto-linting
		"--message", prompt,
	}

	cmd := exec.CommandContext(ctx, r.config.BinPath, args...)
	cmd.Dir = workDir

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Aider execution failed",
			"error", err,
			"output", string(output),
		)
		return fmt.Errorf("aider execution failed: %w", err)
	}

	slog.Info("Aider completed successfully",
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
