package gotransformers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kris/gotransformers/pkg/api"
)

func TestNewHFModel(t *testing.T) {
	model := NewHFModel("test-model")
	if model == nil {
		t.Fatal("Expected model to be created, got nil")
	}

	info := model.GetModelInfo()
	if info.Name != "test-model" {
		t.Errorf("Expected model name to be 'test-model', got %s", info.Name)
	}
	
	if info.Provider != "huggingface" {
		t.Errorf("Expected provider to be 'huggingface', got %s", info.Provider)
	}
}

func TestNewHFModelWithToken(t *testing.T) {
	token := "test-token-12345"
	model := NewHFModelWithToken("test-model", token)
	
	// Type assertion to access the underlying HFModel
	hfModel, ok := model.(*api.HFModel)
	if !ok {
		t.Fatal("Expected *api.HFModel type")
	}
	
	if hfModel.APIToken != token {
		t.Errorf("Expected API token to be %s, got %s", token, hfModel.APIToken)
	}
}

func TestQuickClassify_MockServer(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"label": "POSITIVE", "score": 0.9995}]`))
	}))
	defer server.Close()

	// Create a model and override its base URL
	model := NewHFModel("distilbert-base-uncased-finetuned-sst-2-english")
	hfModel := model.(*api.HFModel)
	hfModel.BaseURL = server.URL

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := model.Classify(ctx, "This is a great test!")
	if err != nil {
		t.Fatalf("QuickClassify failed: %v", err)
	}

	if result.Label != "POSITIVE" {
		t.Errorf("Expected label 'POSITIVE', got %s", result.Label)
	}

	if result.Score != 0.9995 {
		t.Errorf("Expected score 0.9995, got %f", result.Score)
	}
}

func TestQuickGenerate_MockServer(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"generated_text": "Hello world! This is a test of text generation."}]`))
	}))
	defer server.Close()

	// Create a model and override its base URL
	model := NewHFModel("gpt2")
	hfModel := model.(*api.HFModel)
	hfModel.BaseURL = server.URL

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := model.Generate(ctx, "Hello", nil)
	if err != nil {
		t.Fatalf("QuickGenerate failed: %v", err)
	}

	expected := "Hello world! This is a test of text generation."
	if result.GeneratedText != expected {
		t.Errorf("Expected generated text %s, got %s", expected, result.GeneratedText)
	}
}