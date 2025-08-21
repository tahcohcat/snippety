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

func (c *Client) GenerateCommitMessage(ctx context.Context, diff string, tone string) (string, error) {
	toneInstruction := getToneInstruction(tone)

	prompt := fmt.Sprintf(`Based on the git diff below, respond with ONLY the commit message. No explanations, no additional text, no formatting - just the commit message itself.

Requirements:
- Present tense (Add, Fix, Update, Remove)
- Under 50 characters
- Conventional commit format
%s

Git diff:
%s

Commit message:`, toneInstruction, diff)

	req := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := c.BaseURL + "/api/generate"
	logrus.Debug("Making request to: %s\n", url)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to make request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := json.Marshal(resp.Body)
		return "", fmt.Errorf("ollama request to %s failed with status: %d, response: %s", url, resp.StatusCode, string(body))
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return strings.TrimSpace(result.Response), nil
}

func getToneInstruction(tone string) string {
	switch tone {
	case "fun":
		return "- Use a fun, playful tone with emojis and creative language while keeping it professional"
	case "pirate":
		return "- Write the commit message in pirate speak with nautical terminology (e.g., 'Hoist', 'Plunder', 'Navigate')"
	case "serious":
		return "- Use a very serious, formal tone with technical precision and no casual language"
	case "professional":
		fallthrough
	default:
		return "- Use a professional, clear tone"
	}
}
