package models

import "context"

// Task represents the type of NLP task
type Task string

const (
	TaskTextGeneration    Task = "text-generation"
	TaskTextClassification Task = "text-classification"
	TaskTokenClassification Task = "token-classification"
	TaskQuestionAnswering Task = "question-answering"
	TaskFillMask         Task = "fill-mask"
	TaskSummarization    Task = "summarization"
	TaskTranslation      Task = "translation"
)

// Model represents a transformer model interface
type Model interface {
	// Classify performs text classification
	Classify(ctx context.Context, text string) (*ClassificationResult, error)
	
	// Generate performs text generation
	Generate(ctx context.Context, prompt string, options *GenerationOptions) (*GenerationResult, error)
	
	// GetModelInfo returns information about the model
	GetModelInfo() *ModelInfo
}

// ClassificationResult represents the result of text classification
type ClassificationResult struct {
	Label string  `json:"label"`
	Score float64 `json:"score"`
}

// GenerationOptions configures text generation parameters
type GenerationOptions struct {
	MaxLength    int     `json:"max_length,omitempty"`
	Temperature  float64 `json:"temperature,omitempty"`
	TopP         float64 `json:"top_p,omitempty"`
	TopK         int     `json:"top_k,omitempty"`
	DoSample     bool    `json:"do_sample,omitempty"`
	NumReturn    int     `json:"num_return_sequences,omitempty"`
	Stream       bool    `json:"stream,omitempty"`
}

// GenerationResult represents the result of text generation
type GenerationResult struct {
	GeneratedText string  `json:"generated_text"`
	Score         float64 `json:"score,omitempty"`
}

// ModelInfo contains metadata about a model
type ModelInfo struct {
	Name     string `json:"name"`
	Task     Task   `json:"task"`
	Provider string `json:"provider"` // "huggingface", "onnx", "gguf"
}

// BatchResult represents results for batch processing
type BatchResult[T any] struct {
	Results []T     `json:"results"`
	Errors  []error `json:"errors,omitempty"`
}