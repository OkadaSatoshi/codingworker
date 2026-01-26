package aider

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/OkadaSatoshi/codingworker/worker/internal/config"
)

const (
	// maxFixAttempts is the maximum number of times to ask Aider to fix errors
	maxFixAttempts = 3
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

// RunWithTests executes Aider in 2 passes: implementation + test creation
// Each pass includes retry-with-fix logic for build/lint/test failures
func (r *Runner) RunWithTests(ctx context.Context, workDir, title, body string) error {
	// Pass 1: Implementation with build verification
	slog.Info("Pass 1: Running implementation")
	if err := r.runAndVerifyBuild(ctx, workDir, title, body); err != nil {
		return fmt.Errorf("pass 1 (implementation) failed: %w", err)
	}

	// Pass 2: Test creation with full verification
	slog.Info("Pass 2: Running test creation")
	testPrompt := fmt.Sprintf("Add unit tests for the changes made for: %s", title)
	if err := r.runAndVerifyAll(ctx, workDir, testPrompt); err != nil {
		return fmt.Errorf("pass 2 (test creation) failed: %w", err)
	}

	return nil
}

// runAndVerifyBuild runs Aider and verifies build, retrying with fix prompts on failure
func (r *Runner) runAndVerifyBuild(ctx context.Context, workDir, title, body string) error {
	// Initial run
	if err := r.Run(ctx, workDir, title, body); err != nil {
		return err
	}

	// Verify build with retry-fix loop
	for attempt := 1; attempt <= maxFixAttempts; attempt++ {
		buildErr := r.verifyBuildWithOutput(ctx, workDir)
		if buildErr == nil {
			return nil // Success
		}

		if attempt == maxFixAttempts {
			return fmt.Errorf("build failed after %d fix attempts: %w", maxFixAttempts, buildErr)
		}

		slog.Warn("Build failed, asking Aider to fix",
			"attempt", attempt,
			"max_attempts", maxFixAttempts,
		)

		// Ask Aider to fix the build error
		fixPrompt := fmt.Sprintf("Fix the following build error:\n\n%s", buildErr.Error())
		if err := r.Run(ctx, workDir, fixPrompt, ""); err != nil {
			return fmt.Errorf("aider fix attempt failed: %w", err)
		}
	}

	return nil
}

// runAndVerifyAll runs Aider and verifies build+lint+test, retrying with fix prompts on failure
func (r *Runner) runAndVerifyAll(ctx context.Context, workDir, prompt string) error {
	// Initial run
	if err := r.Run(ctx, workDir, prompt, ""); err != nil {
		return err
	}

	// Verify with retry-fix loop
	for attempt := 1; attempt <= maxFixAttempts; attempt++ {
		// Check build
		if buildErr := r.verifyBuildWithOutput(ctx, workDir); buildErr != nil {
			if attempt == maxFixAttempts {
				return fmt.Errorf("build failed after %d fix attempts: %w", maxFixAttempts, buildErr)
			}
			slog.Warn("Build failed, asking Aider to fix", "attempt", attempt)
			fixPrompt := fmt.Sprintf("Fix the following build error:\n\n%s", buildErr.Error())
			if err := r.Run(ctx, workDir, fixPrompt, ""); err != nil {
				return fmt.Errorf("aider fix attempt failed: %w", err)
			}
			continue
		}

		// Check lint
		if lintErr := r.verifyLintWithOutput(ctx, workDir); lintErr != nil {
			if attempt == maxFixAttempts {
				return fmt.Errorf("lint failed after %d fix attempts: %w", maxFixAttempts, lintErr)
			}
			slog.Warn("Lint failed, asking Aider to fix", "attempt", attempt)
			fixPrompt := fmt.Sprintf("Fix the following lint error:\n\n%s", lintErr.Error())
			if err := r.Run(ctx, workDir, fixPrompt, ""); err != nil {
				return fmt.Errorf("aider fix attempt failed: %w", err)
			}
			continue
		}

		// Check tests
		if testErr := r.verifyTestsWithOutput(ctx, workDir); testErr != nil {
			if attempt == maxFixAttempts {
				return fmt.Errorf("tests failed after %d fix attempts: %w", maxFixAttempts, testErr)
			}
			slog.Warn("Tests failed, asking Aider to fix", "attempt", attempt)
			fixPrompt := fmt.Sprintf("Fix the following test failure:\n\n%s", testErr.Error())
			if err := r.Run(ctx, workDir, fixPrompt, ""); err != nil {
				return fmt.Errorf("aider fix attempt failed: %w", err)
			}
			continue
		}

		// All checks passed
		slog.Info("All verifications passed", "attempts", attempt)
		return nil
	}

	return nil
}

// verifyBuildWithOutput runs go build and returns error with output for fix prompts
func (r *Runner) verifyBuildWithOutput(ctx context.Context, workDir string) error {
	cmd := exec.CommandContext(ctx, "go", "build", "./...")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Build failed", "output", string(output))
		// Include output in error for fix prompts
		return fmt.Errorf("go build failed:\n%s", strings.TrimSpace(string(output)))
	}
	slog.Info("Build successful")
	return nil
}

// verifyLintWithOutput runs go fmt and go vet, returns error with output for fix prompts
func (r *Runner) verifyLintWithOutput(ctx context.Context, workDir string) error {
	// go fmt (auto-fixes, so we just run it)
	fmtCmd := exec.CommandContext(ctx, "go", "fmt", "./...")
	fmtCmd.Dir = workDir
	if output, err := fmtCmd.CombinedOutput(); err != nil {
		slog.Error("go fmt failed", "output", string(output))
		return fmt.Errorf("go fmt failed:\n%s", strings.TrimSpace(string(output)))
	}

	// go vet
	vetCmd := exec.CommandContext(ctx, "go", "vet", "./...")
	vetCmd.Dir = workDir
	if output, err := vetCmd.CombinedOutput(); err != nil {
		slog.Error("go vet failed", "output", string(output))
		return fmt.Errorf("go vet failed:\n%s", strings.TrimSpace(string(output)))
	}

	slog.Info("Lint passed")
	return nil
}

// verifyTestsWithOutput runs go test and returns error with output for fix prompts
func (r *Runner) verifyTestsWithOutput(ctx context.Context, workDir string) error {
	cmd := exec.CommandContext(ctx, "go", "test", "./...")
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error("Tests failed", "output", string(output))
		// Include output in error for fix prompts
		return fmt.Errorf("go test failed:\n%s", strings.TrimSpace(string(output)))
	}
	slog.Info("Tests passed", "output", strings.TrimSpace(string(output)))
	return nil
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
