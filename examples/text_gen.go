package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kelleyblackmore/go-transformer"
	"github.com/kelleyblackmore/go-transformer/pkg/models"
)

func main() {
	ctx := context.Background()

	// Example 1: Text Classification (from your sample usage)
	fmt.Println("=== Text Classification Example ===")
	model := gotransformers.NewHFModel("distilbert-base-uncased-finetuned-sst-2-english")
	result, err := model.Classify(ctx, "Go is amazing for backend services!")
	if err != nil {
		log.Printf("Classification error: %v", err)
	} else {
		fmt.Printf("Label: %s, Score: %.4f\n", result.Label, result.Score)
	}

	// Example 2: Text Generation
	fmt.Println("\n=== Text Generation Example ===")
	genModel := gotransformers.NewHFModel("gpt2")
	genResult, err := genModel.Generate(ctx, "The future of AI is", &models.GenerationOptions{
		MaxLength:   100,
		Temperature: 0.7,
		DoSample:    true,
	})
	if err != nil {
		log.Printf("Generation error: %v", err)
	} else {
		fmt.Printf("Generated text: %s\n", genResult.GeneratedText)
	}

	// Example 3: Quick Classification (convenience function)
	fmt.Println("\n=== Quick Classification Example ===")
	quickResult, err := gotransformers.QuickClassify(ctx, "This movie is terrible!")
	if err != nil {
		log.Printf("Quick classification error: %v", err)
	} else {
		fmt.Printf("Quick Label: %s, Score: %.4f\n", quickResult.Label, quickResult.Score)
	}

	// Example 4: Quick Generation (convenience function)
	fmt.Println("\n=== Quick Generation Example ===")
	quickGenResult, err := gotransformers.QuickGenerate(ctx, "Once upon a time")
	if err != nil {
		log.Printf("Quick generation error: %v", err)
	} else {
		fmt.Printf("Quick Generated: %s\n", quickGenResult.GeneratedText)
	}

	// Example 5: Model Information
	fmt.Println("\n=== Model Information ===")
	info := model.GetModelInfo()
	fmt.Printf("Model: %s, Task: %s, Provider: %s\n", info.Name, info.Task, info.Provider)
}
