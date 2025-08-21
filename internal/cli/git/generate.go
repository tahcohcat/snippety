package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/tahcohcat/snippety/internal/ollama"
)

func GenerateCommitMessage(ollamaURL, ollamaModel string, showDiff bool, tone string) {
	diff, err := getStagedDiff()
	if err != nil {
		fmt.Printf("Error getting staged diff: %v\n", err)
		return
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("No staged changes found. Please stage your changes with 'git add' first.")
		return
	}

	if showDiff {
		fmt.Println("Git diff output:")
		fmt.Println("================")
		fmt.Println(diff)
		fmt.Println("================")
		fmt.Println()
	}

	logrus.
		WithField("llm", "ollama").
		WithField("url", ollamaURL).
		WithField("model", ollamaModel).
		Info("generating commit message")

	client := ollama.NewClient(ollamaURL, ollamaModel)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Check if Ollama is available
	if err := client.HealthCheck(ctx); err != nil {
		fmt.Printf("Ollama health check failed: %v\n", err)
		fmt.Println("Falling back to basic analysis...")
		commitMessage := analyzeAndGenerateMessage(diff)
		fmt.Println("Generated commit message:")
		fmt.Println(strings.TrimSpace(commitMessage))
		return
	}

	commitMessage, err := client.GenerateCommitMessage(ctx, diff, tone)
	if err != nil {
		fmt.Printf("Error generating commit message with ollama: %v\n", err)
		fmt.Printf("Falling back to basic analysis...")
		commitMessage = analyzeAndGenerateMessage(diff)
	}

	fmt.Println("Generated commit message:")
	fmt.Println(strings.TrimSpace(commitMessage))
}

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--summary")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}
	return string(output), nil
}

func analyzeAndGenerateMessage(diff string) string {
	lines := strings.Split(diff, "\n")

	var addedFiles []string
	var modifiedFiles []string
	var deletedFiles []string
	var addedLines, deletedLines int

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				fileName := strings.TrimPrefix(parts[3], "b/")
				modifiedFiles = append(modifiedFiles, fileName)
			}
		} else if strings.HasPrefix(line, "new file mode") {
			if len(modifiedFiles) > 0 {
				addedFiles = append(addedFiles, modifiedFiles[len(modifiedFiles)-1])
				modifiedFiles = modifiedFiles[:len(modifiedFiles)-1]
			}
		} else if strings.HasPrefix(line, "deleted file mode") {
			if len(modifiedFiles) > 0 {
				deletedFiles = append(deletedFiles, modifiedFiles[len(modifiedFiles)-1])
				modifiedFiles = modifiedFiles[:len(modifiedFiles)-1]
			}
		} else if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletedLines++
		}
	}

	if len(addedFiles) > 0 {
		if len(addedFiles) == 1 {
			return fmt.Sprintf("Add %s", addedFiles[0])
		}
		return fmt.Sprintf("Add %d new files", len(addedFiles))
	}

	if len(deletedFiles) > 0 {
		if len(deletedFiles) == 1 {
			return fmt.Sprintf("Remove %s", deletedFiles[0])
		}
		return fmt.Sprintf("Remove %d files", len(deletedFiles))
	}

	if len(modifiedFiles) > 0 {
		if len(modifiedFiles) == 1 {
			if addedLines > deletedLines*2 {
				return fmt.Sprintf("Enhance %s", modifiedFiles[0])
			} else if deletedLines > addedLines*2 {
				return fmt.Sprintf("Refactor %s", modifiedFiles[0])
			}
			return fmt.Sprintf("Update %s", modifiedFiles[0])
		}
		return fmt.Sprintf("Update %d files", len(modifiedFiles))
	}

	return "Update project files"
}
