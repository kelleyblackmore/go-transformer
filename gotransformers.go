// Package gotransformers provides a Go interface for transformer models
// supporting both local inference (ONNX/GGUF) and cloud inference (Hugging Face API)
package gotransformers

import (
	"context"

	"github.com/kelleyblackmore/go-transformer/pkg/api"
	"github.com/kelleyblackmore/go-transformer/pkg/models"
)

// NewHFModel creates a new Hugging Face model instance
// This uses the Hugging Face Inference API for cloud-based inference
func NewHFModel(modelName string) models.Model {
	return api.NewHFModel(modelName)
}

// NewHFModelWithToken creates a new Hugging Face model with explicit API token
func NewHFModelWithToken(modelName, apiToken string) models.Model {
	return api.NewHFModelWithToken(modelName, apiToken)
}

// QuickClassify is a convenience function for quick text classification
// using a default model without needing to manage model instances
func QuickClassify(ctx context.Context, text string) (*models.ClassificationResult, error) {
	model := NewHFModel("distilbert-base-uncased-finetuned-sst-2-english")
	return model.Classify(ctx, text)
}

// QuickGenerate is a convenience function for quick text generation
// using a default model without needing to manage model instances
func QuickGenerate(ctx context.Context, prompt string) (*models.GenerationResult, error) {
	model := NewHFModel("gpt2")
	return model.Generate(ctx, prompt, &models.GenerationOptions{
		MaxLength:   50,
		Temperature: 0.7,
		DoSample:    true,
	})
}
