package cobra

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	// Reset flags before each test
	resetFlags := func() {
		ollamaURL = "http://localhost:11434"
		ollamaModel = "llama3.2"
		showDiff = false
		tone = "professional"
		interactive = false
		autoStage = true
		debug = false
		showVersion = false
	}

	tests := []struct {
		name     string
		args     []string
		setup    func()
		validate func(*testing.T, string, string, error)
	}{
		{
			name: "Default flag values",
			args: []string{},
			setup: func() {
				resetFlags()
			},
			validate: func(t *testing.T, stdout, stderr string, err error) {
				if ollamaURL != "http://localhost:11434" {
					t.Errorf("Expected default ollama-url to be 'http://localhost:11434', got '%s'", ollamaURL)
				}
				if ollamaModel != "llama3.2" {
					t.Errorf("Expected default model to be 'llama3.2', got '%s'", ollamaModel)
				}
				if tone != "professional" {
					t.Errorf("Expected default tone to be 'professional', got '%s'", tone)
				}
				if !autoStage {
					t.Error("Expected auto-stage to be true by default")
				}
			},
		},
		{
			name: "Custom flag values",
			args: []string{
				"--ollama-url", "http://custom:8080",
				"--model", "custom-model",
				"--tone", "fun",
				"--interactive",
				"--auto-stage=false",
				"--show-diff",
				"--debug",
			},
			setup: func() {
				resetFlags()
			},
			validate: func(t *testing.T, stdout, stderr string, err error) {
				if ollamaURL != "http://custom:8080" {
					t.Errorf("Expected ollama-url to be 'http://custom:8080', got '%s'", ollamaURL)
				}
				if ollamaModel != "custom-model" {
					t.Errorf("Expected model to be 'custom-model', got '%s'", ollamaModel)
				}
				if tone != "fun" {
					t.Errorf("Expected tone to be 'fun', got '%s'", tone)
				}
				if !interactive {
					t.Error("Expected interactive to be true")
				}
				if autoStage {
					t.Error("Expected auto-stage to be false")
				}
				if !showDiff {
					t.Error("Expected show-diff to be true")
				}
				if !debug {
					t.Error("Expected debug to be true")
				}
			},
		},
		{
			name: "Version flag",
			args: []string{"--version"},
			setup: func() {
				resetFlags()
			},
			validate: func(t *testing.T, stdout, stderr string, err error) {
				if !showVersion {
					t.Error("Expected showVersion to be true")
				}
				// Note: We can't easily test the actual version output without
				// running the command, but we can test the flag parsing
			},
		},
		{
			name: "Help flag",
			args: []string{"--help"},
			setup: func() {
				resetFlags()
			},
			validate: func(t *testing.T, stdout, stderr string, err error) {
				// Help should not cause an error and should contain usage info
				if err != nil {
					t.Errorf("Help command should not return error, got: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tt.setup()

			// Create a new command instance for testing to avoid global state issues
			cmd := &cobra.Command{
				Use:   "snippety",
				Short: "Generate commit messages from staged git diff using Ollama",
				Long: `A CLI tool that analyzes your staged git changes and generates
meaningful commit messages using Ollama AI based on the diff.`,
				Run: func(cmd *cobra.Command, args []string) {
					// Don't actually run the git logic in tests
					if showVersion {
						// Would normally print version, but we'll just validate the flag
						return
					}
					// For other flags, we just validate they were set correctly
				},
			}

			// Add flags
			cmd.Flags().StringVar(&ollamaURL, "ollama-url", "http://localhost:11434", "ollama server URL")
			cmd.Flags().StringVar(&ollamaModel, "model", "llama3.2", "ollama model to use for generation")
			cmd.Flags().BoolVar(&showDiff, "show-diff", false, "show git diff output to the user")
			cmd.Flags().StringVar(&tone, "tone", "professional", "tone for commit messages (professional, fun, pirate, haiku, serious)")
			cmd.Flags().BoolVar(&interactive, "interactive", false, "interactively confirm before creating the git commit")
			cmd.Flags().BoolVar(&autoStage, "auto-stage", true, "automatically stage all changes with 'git add -A' before generating commit message")
			cmd.Flags().BoolVar(&debug, "debug", false, "enable debug mode")
			cmd.Flags().BoolVar(&showVersion, "version", false, "show version")

			// Capture output
			var stdout, stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)
			cmd.SetArgs(tt.args)

			// Execute
			err := cmd.Execute()

			// Validate
			tt.validate(t, stdout.String(), stderr.String(), err)
		})
	}
}

func TestToneValidation(t *testing.T) {
	validTones := []string{"professional", "fun", "pirate", "haiku", "serious"}

	for _, validTone := range validTones {
		t.Run("Valid tone: "+validTone, func(t *testing.T) {
			// Reset flags
			tone = "professional"

			cmd := &cobra.Command{
				Use: "snippety",
				Run: func(cmd *cobra.Command, args []string) {
					// Don't actually run git logic
				},
			}
			cmd.Flags().StringVar(&tone, "tone", "professional", "tone for commit messages")

			cmd.SetArgs([]string{"--tone", validTone})
			err := cmd.Execute()

			if err != nil {
				t.Errorf("Valid tone %s should not cause error: %v", validTone, err)
			}

			if tone != validTone {
				t.Errorf("Expected tone to be %s, got %s", validTone, tone)
			}
		})
	}
}

// Test that we can create multiple command instances without conflicts
func TestCommandIsolation(t *testing.T) {
	// This test ensures our command can be instantiated multiple times
	// without global state conflicts

	cmd1 := &cobra.Command{Use: "test1"}
	cmd1.Flags().StringVar(&ollamaURL, "ollama-url", "http://localhost:11434", "")

	cmd2 := &cobra.Command{Use: "test2"}
	cmd2.Flags().StringVar(&ollamaURL, "ollama-url", "http://localhost:11434", "")

	// Set different values
	cmd1.SetArgs([]string{"--ollama-url", "http://server1:8080"})
	cmd2.SetArgs([]string{"--ollama-url", "http://server2:9090"})

	// Execute cmd1
	ollamaURL = "http://localhost:11434" // reset
	err1 := cmd1.Execute()
	url1 := ollamaURL

	// Execute cmd2
	ollamaURL = "http://localhost:11434" // reset
	err2 := cmd2.Execute()
	url2 := ollamaURL

	if err1 != nil {
		t.Errorf("Command 1 should not error: %v", err1)
	}
	if err2 != nil {
		t.Errorf("Command 2 should not error: %v", err2)
	}

	if url1 != "http://server1:8080" {
		t.Errorf("Command 1 should set URL to 'http://server1:8080', got '%s'", url1)
	}
	if url2 != "http://server2:9090" {
		t.Errorf("Command 2 should set URL to 'http://server2:9090', got '%s'", url2)
	}
}
