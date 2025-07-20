package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kelleyblackmore/go-transformer/pkg/models"
)

func TestNewHFModel(t *testing.T) {
	model := NewHFModel("test-model")
	
	if model.ModelName != "test-model" {
		t.Errorf("Expected ModelName to be 'test-model', got %s", model.ModelName)
	}
	
	if model.BaseURL != HuggingFaceAPIBase {
		t.Errorf("Expected BaseURL to be %s, got %s", HuggingFaceAPIBase, model.BaseURL)
	}
	
	if model.Client.Timeout != DefaultTimeout {
		t.Errorf("Expected timeout to be %v, got %v", DefaultTimeout, model.Client.Timeout)
	}
}

func TestNewHFModelWithToken(t *testing.T) {
	token := "test-token"
	model := NewHFModelWithToken("test-model", token)
	
	if model.APIToken != token {
		t.Errorf("Expected APIToken to be %s, got %s", token, model.APIToken)
	}
}

func TestHFModel_GetModelInfo(t *testing.T) {
	model := NewHFModel("test-model")
	info := model.GetModelInfo()
	
	if info.Name != "test-model" {
		t.Errorf("Expected Name to be 'test-model', got %s", info.Name)
	}
	
	if info.Provider != "huggingface" {
		t.Errorf("Expected Provider to be 'huggingface', got %s", info.Provider)
	}
	
	if info.Task != models.TaskTextClassification {
		t.Errorf("Expected Task to be %s, got %s", models.TaskTextClassification, info.Task)
	}
}

func TestHFModel_Classify(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type to be application/json")
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"label": "POSITIVE", "score": 0.9998}]`))
	}))
	defer server.Close()
	
	model := NewHFModel("test-model")
	model.BaseURL = server.URL
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	result, err := model.Classify(ctx, "This is a test")
	if err != nil {
		t.Fatalf("Classification failed: %v", err)
	}
	
	if result.Label != "POSITIVE" {
		t.Errorf("Expected label to be 'POSITIVE', got %s", result.Label)
	}
	
	if result.Score != 0.9998 {
		t.Errorf("Expected score to be 0.9998, got %f", result.Score)
	}
}

func TestHFModel_Generate(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"generated_text": "Hello world! This is generated text."}]`))
	}))
	defer server.Close()
	
	model := NewHFModel("test-model")
	model.BaseURL = server.URL
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	options := &models.GenerationOptions{
		MaxLength:   50,
		Temperature: 0.7,
		DoSample:    true,
	}
	
	result, err := model.Generate(ctx, "Hello", options)
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}
	
	expected := "Hello world! This is generated text."
	if result.GeneratedText != expected {
		t.Errorf("Expected generated text to be %s, got %s", expected, result.GeneratedText)
	}
}

func TestHFModel_APIError(t *testing.T) {
	// Mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Bad request"}`))
	}))
	defer server.Close()
	
	model := NewHFModel("test-model")
	model.BaseURL = server.URL
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := model.Classify(ctx, "This is a test")
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}
	
	if err.Error() == "" {
		t.Error("Expected a non-empty error message")
	}
}