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

%s

Respond with exactly this format:
TITLE: [short commit title]
DESCRIPTION: [detailed description]

Title requirements:
- Present tense (Add, Fix, Update, Remove)
- Under 50 characters
- Conventional commit format

Description requirements:
- 2-3 sentences explaining what was changed and why
- Include technical details about the implementation
- Mention any test cases or validation added
- No prefix needed just the description itself

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
	logrus.
		WithField("request.prompt", req.Prompt).
		Debugf("Making request to:%s", url)

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
	// Debug: log the raw response to understand the format
	logrus.Debugf("Raw LLM response: %q", response)

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
			title = strings.TrimSpace(lines[0])
		}
		if len(lines) > 1 {
			description = strings.TrimSpace(strings.Join(lines[1:], " "))
		}
	}

	// If still no description, generate a basic one
	if description == "" {
		description = "Code changes as shown in the git diff."
	}

	// Debug: log parsed components
	logrus.Debugf("Parsed title: %q", title)
	logrus.Debugf("Parsed description: %q", description)

	return CommitMessage{
		Title:       title,
		Description: description,
	}
}

func getToneInstruction(tone string) string {
	switch tone {
	case "fun":
		return "TONE INSTRUCTION: Write BOTH the title and description using a fun, playful tone with emojis and creative language while keeping it professional."
	case "pirate":
		return "TONE INSTRUCTION: Write BOTH the title and description in pirate speak with nautical terminology (e.g., 'Hoist', 'Plunder', 'Navigate', 'Arrr', 'matey')."
	case "haiku":
		return "TONE INSTRUCTION: Write the TITLE as a single-line haiku with 5-7-5 syllable structure, separating each line with ' / '. Write the description in a poetic, zen-like tone."
	case "serious":
		return "TONE INSTRUCTION: Write BOTH the title and description using a very serious, formal tone with technical precision and no casual language."
	case "professional":
		return "TONE INSTRUCTION: Write BOTH the title and description using a professional, clear tone."
	default:
		// Custom tone provided by user
		return fmt.Sprintf(`TONE INSTRUCTION: Write BOTH the title and description using a %s tone. 

Examples of how to apply this tone:
- If the tone is "like a joke" or "funny": Use humor, puns, wordplay, or amusing language while keeping it understandable
- If the tone is "dramatic": Use intense, theatrical language with strong emotions and vivid descriptions  
- If the tone is "casual": Use relaxed, informal language like you're talking to a friend
- If the tone is "poetic": Use metaphors, rhythm, and beautiful imagery
- If the tone is "sarcastic": Use irony and subtle mockery while still being informative
- If the tone is a specific style (e.g., "like Shakespeare"): Mimic the vocabulary, sentence structure, and mannerisms of that style

Be creative and fully commit to this %s tone in BOTH the title and description. Don't just mention the tone - actually write in that style.`, tone, tone)
	}
}
