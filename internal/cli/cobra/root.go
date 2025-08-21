package cobra

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/tahcohcat/snippety/internal/cli/git"
)

var (
	ollamaURL   string
	ollamaModel string
	showDiff    bool
	tone        string
	interactive bool
	autoStage   bool
	debug       bool
)

var rootCmd = &cobra.Command{
	Use:   "snippety",
	Short: "Generate commit messages from staged git diff using Ollama",
	Long: `A CLI tool that analyzes your staged git changes and generates
meaningful commit messages using Ollama AI based on the diff.`,
	Run: func(cmd *cobra.Command, args []string) {

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
			logrus.Debug("debug mode enabled")
		}

		git.GenerateCommitMessage(ollamaURL, ollamaModel, showDiff, tone, interactive, autoStage)
	},
}

func init() {
	rootCmd.Flags().StringVar(&ollamaURL, "ollama-url", "http://localhost:11434", "ollama server URL")
	rootCmd.Flags().StringVar(&ollamaModel, "model", "llama3.2", "ollama model to use for generation")
	rootCmd.Flags().BoolVar(&showDiff, "show-diff", false, "show git diff output to the user")
	rootCmd.Flags().StringVar(&tone, "tone", "professional", "tone for commit messages (professional, fun, pirate, haiku, serious)")
	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "interactively confirm before creating the git commit")
	rootCmd.Flags().BoolVar(&autoStage, "auto-stage", true, "automatically stage all changes with 'git add -A' before generating commit message")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "enable debug mode")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
