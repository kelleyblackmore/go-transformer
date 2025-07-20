package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/kris/gotransformers"
	"github.com/kris/gotransformers/pkg/models"
)

var (
	modelName   string
	apiToken    string
	outputJSON  bool
	maxLength   int
	temperature float64
	timeout     time.Duration
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gotransformers",
		Short: "A Go CLI for transformer models",
		Long:  `A command-line interface for running transformer models via Hugging Face API or local inference.`,
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&modelName, "model", "", "Model name to use")
	rootCmd.PersistentFlags().StringVar(&apiToken, "token", "", "Hugging Face API token")
	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "Output results in JSON format")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")

	// Add subcommands
	rootCmd.AddCommand(classifyCmd())
	rootCmd.AddCommand(generateCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func classifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "classify [text]",
		Short: "Classify text using a transformer model",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			text := args[0]
			
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			var result *models.ClassificationResult
			var err error

			if modelName != "" {
				var model models.Model
				if apiToken != "" {
					model = gotransformers.NewHFModelWithToken(modelName, apiToken)
				} else {
					model = gotransformers.NewHFModel(modelName)
				}
				result, err = model.Classify(ctx, text)
			} else {
				result, err = gotransformers.QuickClassify(ctx, text)
			}

			if err != nil {
				return fmt.Errorf("classification failed: %w", err)
			}

			if outputJSON {
				output, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(output))
			} else {
				fmt.Printf("Label: %s\n", result.Label)
				fmt.Printf("Score: %.4f\n", result.Score)
			}

			return nil
		},
	}

	return cmd
}

func generateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [prompt]",
		Short: "Generate text using a transformer model",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			prompt := args[0]
			
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			options := &models.GenerationOptions{
				MaxLength:   maxLength,
				Temperature: temperature,
				DoSample:    temperature > 0,
			}

			var result *models.GenerationResult
			var err error

			if modelName != "" {
				var model models.Model
				if apiToken != "" {
					model = gotransformers.NewHFModelWithToken(modelName, apiToken)
				} else {
					model = gotransformers.NewHFModel(modelName)
				}
				result, err = model.Generate(ctx, prompt, options)
			} else {
				result, err = gotransformers.QuickGenerate(ctx, prompt)
			}

			if err != nil {
				return fmt.Errorf("generation failed: %w", err)
			}

			if outputJSON {
				output, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(output))
			} else {
				fmt.Println(result.GeneratedText)
			}

			return nil
		},
	}

	cmd.Flags().IntVar(&maxLength, "max-length", 50, "Maximum length of generated text")
	cmd.Flags().Float64Var(&temperature, "temperature", 0.7, "Temperature for text generation")

	return cmd
}