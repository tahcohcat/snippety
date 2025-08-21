package cobra

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/tahcohcat/snippety/internal/cli/git"
)

var (
	ollamaURL   string
	ollamaModel string
)

var rootCmd = &cobra.Command{
	Use:   "snippety",
	Short: "Generate commit messages from staged git diff using Ollama",
	Long: `A CLI tool that analyzes your staged git changes and generates
meaningful commit messages using Ollama AI based on the diff.`,
	Run: func(cmd *cobra.Command, args []string) {
		git.GenerateCommitMessage(ollamaURL, ollamaModel)
	},
}

func init() {
	rootCmd.Flags().StringVar(&ollamaURL, "ollama-url", "http://localhost:11434", "Ollama server URL")
	rootCmd.Flags().StringVar(&ollamaModel, "model", "llama3.2", "Ollama model to use for generation")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
