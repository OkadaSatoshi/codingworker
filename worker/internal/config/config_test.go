package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create temporary config file
	content := `
sqs:
  queue_url: "https://sqs.ap-northeast-1.amazonaws.com/123456789/test-queue"
  region: "ap-northeast-1"
  use_mock: true
  wait_time_seconds: 10
  visibility_timeout: 300
aider:
  model: "ollama_chat/qwen2.5-coder:7b"
  bin_path: "/usr/local/bin/aider"
  timeout_seconds: 1800
github:
  token: "test-token"
  clone_base_dir: "/tmp/workdir"
worker:
  max_retries: 5
  worker_id: "test-worker"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify SQS config
	if cfg.SQS.QueueURL != "https://sqs.ap-northeast-1.amazonaws.com/123456789/test-queue" {
		t.Errorf("unexpected queue_url: %s", cfg.SQS.QueueURL)
	}
	if cfg.SQS.Region != "ap-northeast-1" {
		t.Errorf("unexpected region: %s", cfg.SQS.Region)
	}
	if !cfg.SQS.UseMock {
		t.Error("expected use_mock to be true")
	}
	if cfg.SQS.WaitTimeSeconds != 10 {
		t.Errorf("expected wait_time_seconds 10, got %d", cfg.SQS.WaitTimeSeconds)
	}
	if cfg.SQS.VisibilityTimeout != 300 {
		t.Errorf("expected visibility_timeout 300, got %d", cfg.SQS.VisibilityTimeout)
	}

	// Verify Aider config
	if cfg.Aider.Model != "ollama_chat/qwen2.5-coder:7b" {
		t.Errorf("unexpected model: %s", cfg.Aider.Model)
	}
	if cfg.Aider.BinPath != "/usr/local/bin/aider" {
		t.Errorf("unexpected bin_path: %s", cfg.Aider.BinPath)
	}
	if cfg.Aider.Timeout != 1800 {
		t.Errorf("expected timeout 1800, got %d", cfg.Aider.Timeout)
	}

	// Verify GitHub config
	if cfg.GitHub.Token != "test-token" {
		t.Errorf("unexpected token: %s", cfg.GitHub.Token)
	}
	if cfg.GitHub.CloneBaseDir != "/tmp/workdir" {
		t.Errorf("unexpected clone_base_dir: %s", cfg.GitHub.CloneBaseDir)
	}

	// Verify Worker config
	if cfg.Worker.MaxRetries != 5 {
		t.Errorf("expected max_retries 5, got %d", cfg.Worker.MaxRetries)
	}
	if cfg.Worker.WorkerID != "test-worker" {
		t.Errorf("unexpected worker_id: %s", cfg.Worker.WorkerID)
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Minimal config - should use defaults
	content := `
sqs:
  queue_url: "test"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check defaults
	if cfg.SQS.WaitTimeSeconds != 20 {
		t.Errorf("expected default wait_time_seconds 20, got %d", cfg.SQS.WaitTimeSeconds)
	}
	if cfg.SQS.VisibilityTimeout != 3600 {
		t.Errorf("expected default visibility_timeout 3600, got %d", cfg.SQS.VisibilityTimeout)
	}
	if cfg.Aider.Model != "ollama_chat/qwen2.5-coder:1.5b" {
		t.Errorf("expected default model, got %s", cfg.Aider.Model)
	}
	if cfg.Aider.BinPath != "aider" {
		t.Errorf("expected default bin_path 'aider', got %s", cfg.Aider.BinPath)
	}
	if cfg.Aider.Timeout != 3600 {
		t.Errorf("expected default timeout 3600, got %d", cfg.Aider.Timeout)
	}
	if cfg.Worker.MaxRetries != 3 {
		t.Errorf("expected default max_retries 3, got %d", cfg.Worker.MaxRetries)
	}
}

func TestLoad_EnvExpansion(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_GITHUB_TOKEN", "secret-token-from-env")
	defer os.Unsetenv("TEST_GITHUB_TOKEN")

	content := `
github:
  token: "$TEST_GITHUB_TOKEN"
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.GitHub.Token != "secret-token-from-env" {
		t.Errorf("expected token from env, got %s", cfg.GitHub.Token)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	content := `
invalid: yaml: content
  - broken
`
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}
