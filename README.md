# 🧠 go-transformer

A Transformers-compatible library in Go for running transformer models with a focus on inference. Provides an ergonomic, idiomatic Go API with support for both cloud (Hugging Face API) and local inference (ONNX/GGUF).

## 🏗️ Features

- **Text Classification** - Sentiment analysis, topic classification, etc.
- **Text Generation** - GPT-style text completion and generation
- **Multiple Backends** - Hugging Face API (✅), ONNX Runtime (🚧), GGUF/llama.cpp (🚧)
- **CLI Tool** - Command-line interface for quick testing
- **Batch Processing** - Process multiple inputs efficiently
- **Context Support** - Full Go context support for timeouts and cancellation

## 🚀 Installation

```bash
go get github.com/kelleyblackmore/go-transformer
```

## 📖 Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/kelleyblackmore/go-transformer"
)

func main() {
    ctx := context.Background()
    
    // Text Classification
    model := gotransformers.NewHFModel("distilbert-base-uncased-finetuned-sst-2-english")
    result, err := model.Classify(ctx, "Go is amazing for backend services!")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Label: %s, Score: %.4f\n", result.Label, result.Score)
    
    // Text Generation
    genModel := gotransformers.NewHFModel("gpt2")
    genResult, err := genModel.Generate(ctx, "The future of AI is", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Generated: %s\n", genResult.GeneratedText)
}
```

### Convenience Functions

```go
// Quick classification with default model
result, err := gotransformers.QuickClassify(ctx, "This movie is great!")

// Quick generation with default model
genResult, err := gotransformers.QuickGenerate(ctx, "Once upon a time")
```

## 🔧 Configuration

### Environment Variables

Set your Hugging Face API token (optional for public models):

```bash
export HUGGINGFACE_API_TOKEN="your_token_here"
# or
export HF_TOKEN="your_token_here"
```

### Using API Token in Code

```go
model := gotransformers.NewHFModelWithToken("model-name", "your-token")
```

## 🖥️ CLI Usage

Build the CLI tool:

```bash
go build -o gotransformers ./cmd
```

### Text Classification

```bash
# Basic classification
./gotransformers classify "This is awesome!"

# With custom model
./gotransformers --model distilbert-base-uncased classify "Hello world"

# JSON output
./gotransformers --json classify "Great product!"
```

### Text Generation

```bash
# Basic generation
./gotransformers generate "The future of programming"

# With parameters
./gotransformers generate "Once upon a time" --max-length 100 --temperature 0.8

# With custom model
./gotransformers --model gpt2-medium generate "In a galaxy far, far away"
```

## 📚 API Reference

### Models Interface

```go
type Model interface {
    Classify(ctx context.Context, text string) (*ClassificationResult, error)
    Generate(ctx context.Context, prompt string, options *GenerationOptions) (*GenerationResult, error)
    GetModelInfo() *ModelInfo
}
```

### Generation Options

```go
type GenerationOptions struct {
    MaxLength    int     `json:"max_length,omitempty"`
    Temperature  float64 `json:"temperature,omitempty"`
    TopP         float64 `json:"top_p,omitempty"`
    TopK         int     `json:"top_k,omitempty"`
    DoSample     bool    `json:"do_sample,omitempty"`
    NumReturn    int     `json:"num_return_sequences,omitempty"`
    Stream       bool    `json:"stream,omitempty"`
}
```

### Results

```go
type ClassificationResult struct {
    Label string  `json:"label"`
    Score float64 `json:"score"`
}

type GenerationResult struct {
    GeneratedText string  `json:"generated_text"`
    Score         float64 `json:"score,omitempty"`
}
```

## 🧪 Testing

Run the test suite:

```bash
go test ./...
```

Run with coverage:

```bash
go test -cover ./...
```

## 🗂️ Project Structure

```
gotransformers/
├── cmd/                  # CLI tool
├── pkg/
│   ├── tokenizers/       # BPE/WordPiece tokenizer implementations
│   ├── models/           # Model interfaces and types
│   ├── inference/        # Inference logic
│   ├── api/              # Hugging Face API support
│   └── utils/            # Downloaders, config parsers, etc.
├── examples/             # Usage examples
├── go.mod
└── README.md
```

## 🛣️ Roadmap

### ✅ Phase 1: MVP with Hugging Face API
- [x] Text classification support
- [x] Text generation support
- [x] CLI tool
- [x] Unit tests
- [x] Documentation

### 🚧 Phase 2: ONNX Runtime Support
- [ ] ONNX model loading
- [ ] Local tokenization (WordPiece, BPE)
- [ ] CPU inference
- [ ] GPU inference (optional)

### 🚧 Phase 3: GGUF / llama.cpp Support
- [ ] GGUF model loading
- [ ] Streaming responses
- [ ] Quantized model support

### 🚧 Phase 4: Advanced Features
- [ ] Model auto-download from HF Hub
- [ ] More model architectures (RoBERTa, DistilBERT)
- [ ] Batch processing optimization
- [ ] WebAssembly builds

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Hugging Face](https://huggingface.co/) for the amazing transformer models and API
- [ONNX Runtime](https://onnxruntime.ai/) for efficient local inference
- [llama.cpp](https://github.com/ggerganov/llama.cpp) for GGUF support

## 💡 Examples

Check out the `examples/` directory for more comprehensive usage examples:

- [Basic Usage](examples/text_gen.go) - Text classification and generation
- [Advanced Usage](examples/advanced.go) - Custom parameters and error handling
- [Batch Processing](examples/batch.go) - Processing multiple inputs

## 📞 Support

- 📧 Email: support@gotransformers.dev
- 💬 Discord: [Join our community](https://discord.gg/gotransformers)
- 🐛 Issues: [GitHub Issues](https://github.com/kelleyblackmore/go-transformer/issues)