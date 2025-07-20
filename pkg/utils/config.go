package utils

import (
	"os"
	"path/filepath"
	"time"
)

// Config represents the library configuration
type Config struct {
	HuggingFaceToken string        `json:"huggingface_token,omitempty"`
	DefaultTimeout   time.Duration `json:"default_timeout,omitempty"`
	CacheDir         string        `json:"cache_dir,omitempty"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	cacheDir := filepath.Join(homeDir, ".cache", "gotransformers")
	
	return &Config{
		HuggingFaceToken: getHuggingFaceToken(),
		DefaultTimeout:   30 * time.Second,
		CacheDir:         cacheDir,
	}
}

// getHuggingFaceToken attempts to retrieve the HF token from environment variables
func getHuggingFaceToken() string {
	// Try multiple environment variable names
	for _, envVar := range []string{"HUGGINGFACE_API_TOKEN", "HF_TOKEN", "HUGGINGFACE_TOKEN"} {
		if token := os.Getenv(envVar); token != "" {
			return token
		}
	}
	return ""
}

// ModelConfig represents configuration for a specific model
type ModelConfig struct {
	Name         string            `json:"name"`
	Provider     string            `json:"provider"`     // "huggingface", "onnx", "gguf"
	Path         string            `json:"path,omitempty"` // Local path for ONNX/GGUF models
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	TokenizerPath string           `json:"tokenizer_path,omitempty"`
}

// PopularModels contains configurations for commonly used models
var PopularModels = map[string]*ModelConfig{
	"sentiment": {
		Name:     "distilbert-base-uncased-finetuned-sst-2-english",
		Provider: "huggingface",
	},
	"gpt2": {
		Name:     "gpt2",
		Provider: "huggingface",
	},
	"gpt2-medium": {
		Name:     "gpt2-medium",
		Provider: "huggingface",
	},
	"bert-base": {
		Name:     "bert-base-uncased",
		Provider: "huggingface",
	},
	"distilbert": {
		Name:     "distilbert-base-uncased",
		Provider: "huggingface",
	},
}