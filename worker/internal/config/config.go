package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SQS    SQSConfig    `yaml:"sqs"`
	Aider  AiderConfig  `yaml:"aider"`
	GitHub GitHubConfig `yaml:"github"`
	Worker WorkerConfig `yaml:"worker"`
}

type SQSConfig struct {
	QueueURL          string `yaml:"queue_url"`
	Region            string `yaml:"region"`
	WaitTimeSeconds   int    `yaml:"wait_time_seconds"`
	VisibilityTimeout int    `yaml:"visibility_timeout"`
	UseMock           bool   `yaml:"use_mock"`
}

type AiderConfig struct {
	Model   string `yaml:"model"`
	BinPath string `yaml:"bin_path"`
	Timeout int    `yaml:"timeout_seconds"`
}

type GitHubConfig struct {
	Token        string `yaml:"token"`
	CloneBaseDir string `yaml:"clone_base_dir"`
}

type WorkerConfig struct {
	MaxRetries int `yaml:"max_retries"`
	WorkerID   string `yaml:"worker_id"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Expand environment variables
	expanded := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, err
	}

	// Set defaults
	if cfg.SQS.WaitTimeSeconds == 0 {
		cfg.SQS.WaitTimeSeconds = 20
	}
	if cfg.SQS.VisibilityTimeout == 0 {
		cfg.SQS.VisibilityTimeout = 3600
	}
	if cfg.Aider.Model == "" {
		cfg.Aider.Model = "ollama_chat/qwen2.5-coder:7b"
	}
	if cfg.Aider.BinPath == "" {
		cfg.Aider.BinPath = "aider"
	}
	if cfg.Aider.Timeout == 0 {
		cfg.Aider.Timeout = 3600
	}
	if cfg.Worker.MaxRetries == 0 {
		cfg.Worker.MaxRetries = 3
	}

	return &cfg, nil
}
