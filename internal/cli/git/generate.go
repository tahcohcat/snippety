package git

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/tahcohcat/snippety/internal/ollama"
)

// Color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
	ColorRed    = "\033[31m"
)

func GenerateCommitMessage(ollamaURL, ollamaModel string, showDiff bool, tone string, interactive bool, autoStage bool) {
	if autoStage {
		logrus.Debug("Staging all changes...")
		if err := stageAllChanges(); err != nil {
			fmt.Printf("%sError staging changes: %v%s\n", ColorRed, err, ColorReset)
			return
		}
	}

	diff, err := getStagedDiff()
	if err != nil {
		fmt.Printf("%sError getting staged diff: %v%s\n", ColorRed, err, ColorReset)
		return
	}

	if strings.TrimSpace(diff) == "" {
		if autoStage {
			fmt.Printf("%sNo changes found to stage and commit.%s\n", ColorYellow, ColorReset)
		} else {
			fmt.Printf("%sNo staged changes found. Please stage your changes with 'git add' first.%s\n", ColorYellow, ColorReset)
		}
		return
	}

	if showDiff {
		fmt.Printf("%s%sGit diff output:%s\n", ColorBold, ColorBlue, ColorReset)
		fmt.Printf("%s================%s\n", ColorBlue, ColorReset)
		fmt.Println(diff)
		fmt.Printf("%s================%s\n", ColorBlue, ColorReset)
		fmt.Println()
	}

	logrus.
		WithField("llm", "ollama").
		WithField("url", ollamaURL).
		WithField("model", ollamaModel).
		Debug("generating commit message")

	client := ollama.NewClient(ollamaURL, ollamaModel)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Check if Ollama is available
	var commitMsg ollama.CommitMessage
	if err := client.HealthCheck(ctx); err != nil {
		fmt.Printf("Ollama health check failed: %v\n", err)
		fmt.Println("Falling back to basic analysis...")
		title := analyzeAndGenerateMessage(diff)
		commitMsg = ollama.CommitMessage{
			Title:       title,
			Description: "Code changes as analyzed from the git diff.",
		}
	} else {
		var err error
		commitMsg, err = client.GenerateCommitMessage(ctx, diff, tone)
		if err != nil {
			fmt.Printf("Error generating commit message with ollama: %v\n", err)
			fmt.Printf("Falling back to basic analysis...")
			title := analyzeAndGenerateMessage(diff)
			commitMsg = ollama.CommitMessage{
				Title:       title,
				Description: "Code changes as analyzed from the git diff.",
			}
		}
	}

	fmt.Printf("%sGenerated commit message:%s\n", ColorBold+ColorBlue, ColorReset)
	fmt.Printf("%sTitle:%s %s%s%s\n", ColorBold+ColorCyan, ColorReset, ColorGreen, commitMsg.Title, ColorReset)
	fmt.Printf("%sDescription:%s %s%s%s\n", ColorBold+ColorCyan, ColorReset, ColorYellow, commitMsg.Description, ColorReset)

	if interactive {
		fmt.Print("\nDo you want to create a commit with this message? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			if err := createCommit(commitMsg.Title, commitMsg.Description); err != nil {
				fmt.Printf("%sError creating commit: %v%s\n", ColorRed, err, ColorReset)
				return
			}
			fmt.Printf("%s✅ Commit created successfully!%s\n", ColorGreen, ColorReset)

			if err := pushCommit(); err != nil {
				fmt.Printf("%sError pushing commit: %v%s\n", ColorRed, err, ColorReset)
				return
			}
			fmt.Printf("%s🚀Commit pushed successfully!%s\n", ColorCyan, ColorReset)
		} else {
			fmt.Println("Commit not created.")
		}
	}
}

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}
	return string(output), nil
}

func stageAllChanges() error {
	cmd := exec.Command("git", "add", "-A")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git add -A failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func createCommit(title, description string) error {
	cmd := exec.Command("git", "commit", "-m", title, "-m", description)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git commit failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func pushCommit() error {
	cmd := exec.Command("git", "push")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git push failed: %w\nOutput: %s", err, string(output))
	}
	return nil
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
