package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	BaseURL string
	Model   string
	client  *http.Client
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type CommitMessage struct {
	Title       string
	Description string
}

func NewClient(baseURL, model string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2"
	}

	return &Client{
		BaseURL: baseURL,
		Model:   model,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) HealthCheck(ctx context.Context) error {
	httpReq, err := http.NewRequestWithContext(ctx, "GET", c.BaseURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama at %s: %w", c.BaseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama health check failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) GenerateCommitMessage(ctx context.Context, diff string, tone string) (CommitMessage, error) {
	toneInstruction := getToneInstruction(tone)

	prompt := fmt.Sprintf(`Based on the git diff below, generate a commit message with both a title and description.

Respond with exactly this format:
TITLE: [short commit title]
DESCRIPTION: [detailed description]

Title requirements:
- Present tense (Add, Fix, Update, Remove)
- Under 50 characters
- Conventional commit format
%s

Description requirements:
- 2-3 sentences explaining what was changed and why
- Include technical details about the implementation
- Mention any test cases or validation added

Git diff:
%s`, toneInstruction, diff)

	req := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return CommitMessage{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.BaseURL + "/api/generate"
	logrus.Debug("Making request to: %s\n", url)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return CommitMessage{}, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return CommitMessage{}, fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := json.Marshal(resp.Body)
		return CommitMessage{}, fmt.Errorf("ollama request to %s failed with status: %d, response: %s", url, resp.StatusCode, string(body))
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return CommitMessage{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return parseCommitMessage(result.Response), nil
}

func parseCommitMessage(response string) CommitMessage {
	lines := strings.Split(strings.TrimSpace(response), "\n")
	
	var title, description string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TITLE:") {
			title = strings.TrimSpace(strings.TrimPrefix(line, "TITLE:"))
		} else if strings.HasPrefix(line, "DESCRIPTION:") {
			description = strings.TrimSpace(strings.TrimPrefix(line, "DESCRIPTION:"))
		}
	}
	
	// Fallback if the LLM didn't follow the format
	if title == "" && description == "" {
		// Use the first line as title and rest as description
		if len(lines) > 0 {
			title = lines[0]
		}
		if len(lines) > 1 {
			description = strings.Join(lines[1:], " ")
		}
	}
	
	// If still no description, generate a basic one
	if description == "" {
		description = "Code changes as shown in the git diff."
	}
	
	return CommitMessage{
		Title:       title,
		Description: description,
	}
}

func getToneInstruction(tone string) string {
	switch tone {
	case "fun":
		return "- Use a fun, playful tone with emojis and creative language while keeping it professional"
	case "pirate":
		return "- Write the commit message in pirate speak with nautical terminology (e.g., 'Hoist', 'Plunder', 'Navigate')"
	case "haiku":
		return "- Write the commit message as a single-line haiku with 5-7-5 syllable structure, separating each line with ' / ', capturing the essence of the code change"
	case "serious":
		return "- Use a very serious, formal tone with technical precision and no casual language"
	case "professional":
		fallthrough
	default:
		return "- Use a professional, clear tone"
	}
}
