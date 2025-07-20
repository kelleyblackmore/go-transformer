package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/kris/gotransformers/pkg/models"
	"github.com/tidwall/gjson"
)

const (
	HuggingFaceAPIBase = "https://api-inference.huggingface.co"
	DefaultTimeout     = 30 * time.Second
)

// HFModel represents a Hugging Face model accessed via API
type HFModel struct {
	ModelName string
	APIToken  string
	Client    *http.Client
	BaseURL   string
}

// NewHFModel creates a new Hugging Face model instance
func NewHFModel(modelName string) *HFModel {
	apiToken := os.Getenv("HUGGINGFACE_API_TOKEN")
	if apiToken == "" {
		apiToken = os.Getenv("HF_TOKEN")
	}

	return &HFModel{
		ModelName: modelName,
		APIToken:  apiToken,
		Client: &http.Client{
			Timeout: DefaultTimeout,
		},
		BaseURL: HuggingFaceAPIBase,
	}
}

// NewHFModelWithToken creates a new Hugging Face model instance with explicit token
func NewHFModelWithToken(modelName, apiToken string) *HFModel {
	return &HFModel{
		ModelName: modelName,
		APIToken:  apiToken,
		Client: &http.Client{
			Timeout: DefaultTimeout,
		},
		BaseURL: HuggingFaceAPIBase,
	}
}

// GetModelInfo returns information about the model
func (hf *HFModel) GetModelInfo() *models.ModelInfo {
	return &models.ModelInfo{
		Name:     hf.ModelName,
		Task:     models.TaskTextClassification, // Default, can be detected
		Provider: "huggingface",
	}
}

// Classify performs text classification using Hugging Face API
func (hf *HFModel) Classify(ctx context.Context, text string) (*models.ClassificationResult, error) {
	payload := map[string]interface{}{
		"inputs": text,
	}

	response, err := hf.makeRequest(ctx, "POST", fmt.Sprintf("/models/%s", hf.ModelName), payload)
	if err != nil {
		return nil, fmt.Errorf("classification request failed: %w", err)
	}

	// Parse the response - HF returns an array of classification results
	if !gjson.Valid(response) {
		return nil, fmt.Errorf("invalid JSON response: %s", response)
	}

	result := gjson.Get(response, "0")
	if !result.Exists() {
		return nil, fmt.Errorf("no classification results in response")
	}

	label := gjson.Get(result.Raw, "label").String()
	score := gjson.Get(result.Raw, "score").Float()

	return &models.ClassificationResult{
		Label: label,
		Score: score,
	}, nil
}

// Generate performs text generation using Hugging Face API
func (hf *HFModel) Generate(ctx context.Context, prompt string, options *models.GenerationOptions) (*models.GenerationResult, error) {
	payload := map[string]interface{}{
		"inputs": prompt,
	}

	// Add generation parameters if provided
	if options != nil {
		parameters := make(map[string]interface{})
		
		if options.MaxLength > 0 {
			parameters["max_length"] = options.MaxLength
		}
		if options.Temperature > 0 {
			parameters["temperature"] = options.Temperature
		}
		if options.TopP > 0 {
			parameters["top_p"] = options.TopP
		}
		if options.TopK > 0 {
			parameters["top_k"] = options.TopK
		}
		if options.DoSample {
			parameters["do_sample"] = options.DoSample
		}
		if options.NumReturn > 0 {
			parameters["num_return_sequences"] = options.NumReturn
		}

		if len(parameters) > 0 {
			payload["parameters"] = parameters
		}
	}

	response, err := hf.makeRequest(ctx, "POST", fmt.Sprintf("/models/%s", hf.ModelName), payload)
	if err != nil {
		return nil, fmt.Errorf("generation request failed: %w", err)
	}

	// Parse the response
	if !gjson.Valid(response) {
		return nil, fmt.Errorf("invalid JSON response: %s", response)
	}

	result := gjson.Get(response, "0")
	if !result.Exists() {
		return nil, fmt.Errorf("no generation results in response")
	}

	generatedText := gjson.Get(result.Raw, "generated_text").String()

	return &models.GenerationResult{
		GeneratedText: generatedText,
	}, nil
}

// makeRequest makes an HTTP request to the Hugging Face API
func (hf *HFModel) makeRequest(ctx context.Context, method, endpoint string, payload interface{}) (string, error) {
	var body io.Reader
	
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	url := hf.BaseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if hf.APIToken != "" {
		req.Header.Set("Authorization", "Bearer "+hf.APIToken)
	}

	resp, err := hf.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return string(responseBody), nil
}