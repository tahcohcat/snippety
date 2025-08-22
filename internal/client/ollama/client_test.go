package ollama

import (
	"testing"
)

func TestParseCommitMessage(t *testing.T) {
	tests := []struct {
		name                string
		response            string
		expectedTitle       string
		expectedDescription string
	}{
		{
			name: "Properly formatted response",
			response: `TITLE: Add user authentication system
DESCRIPTION: Implemented JWT-based authentication middleware for API routes. Added validation for bearer tokens and user session management. Includes comprehensive unit tests for edge cases.`,
			expectedTitle:       "Add user authentication system",
			expectedDescription: "Implemented JWT-based authentication middleware for API routes. Added validation for bearer tokens and user session management. Includes comprehensive unit tests for edge cases.",
		},
		{
			name: "Response with extra whitespace",
			response: `  TITLE:   Fix database connection pooling  
  DESCRIPTION:   Resolved connection leak issues by properly closing database connections. Updated connection pool configuration for better performance.   `,
			expectedTitle:       "Fix database connection pooling",
			expectedDescription: "Resolved connection leak issues by properly closing database connections. Updated connection pool configuration for better performance.",
		},
		{
			name: "Multiline description",
			response: `TITLE: Update API documentation
DESCRIPTION: Updated OpenAPI specifications for all endpoints.
Added examples for request/response formats.
Fixed validation schemas for user registration.`,
			expectedTitle:       "Update API documentation",
			expectedDescription: "Updated OpenAPI specifications for all endpoints.",
		},
		{
			name: "Response without proper format (fallback to first line)",
			response: `Add new logging functionality
This commit adds structured logging with different levels
and proper error handling throughout the application.`,
			expectedTitle:       "Add new logging functionality",
			expectedDescription: "This commit adds structured logging with different levels and proper error handling throughout the application.",
		},
		{
			name:                "Single line response (fallback)",
			response:            "Fix critical security vulnerability",
			expectedTitle:       "Fix critical security vulnerability",
			expectedDescription: "Code changes as shown in the git diff.",
		},
		{
			name:                "Empty response",
			response:            "",
			expectedTitle:       "",
			expectedDescription: "Code changes as shown in the git diff.",
		},
		{
			name:                "Only title provided",
			response:            "TITLE: Refactor user service layer",
			expectedTitle:       "Refactor user service layer",
			expectedDescription: "Code changes as shown in the git diff.",
		},
		{
			name:                "Only description provided",
			response:            "DESCRIPTION: Updated all dependencies to latest versions and fixed security vulnerabilities.",
			expectedTitle:       "",
			expectedDescription: "Updated all dependencies to latest versions and fixed security vulnerabilities.",
		},
		{
			name: "Mixed case labels (fallback to first line)",
			response: `Title: Add caching layer
Description: Implemented Redis-based caching for frequently accessed data.`,
			expectedTitle:       "Title: Add caching layer",
			expectedDescription: "Description: Implemented Redis-based caching for frequently accessed data.",
		},
		{
			name: "Response with additional text",
			response: `Based on the git diff, here's the commit message:

TITLE: Update configuration management
DESCRIPTION: Replaced hardcoded configuration with environment-based config system. Added validation for required environment variables and default fallbacks.

This should work well for your project.`,
			expectedTitle:       "Update configuration management",
			expectedDescription: "Replaced hardcoded configuration with environment-based config system. Added validation for required environment variables and default fallbacks.",
		},
		{
			name: "Response with multiple TITLE/DESCRIPTION (takes last occurrence)",
			response: `TITLE: First title
DESCRIPTION: First description
TITLE: Second title
DESCRIPTION: Second description`,
			expectedTitle:       "Second title",
			expectedDescription: "Second description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCommitMessage(tt.response)

			if result.Title != tt.expectedTitle {
				t.Errorf("parseCommitMessage().Title = %q, want %q", result.Title, tt.expectedTitle)
			}

			if result.Description != tt.expectedDescription {
				t.Errorf("parseCommitMessage().Description = %q, want %q", result.Description, tt.expectedDescription)
			}
		})
	}
}

func TestGetToneInstruction(t *testing.T) {
	tests := []struct {
		name     string
		tone     string
		expected string
	}{
		{
			name:     "Professional tone (default)",
			tone:     "professional",
			expected: "- Use a professional, clear tone",
		},
		{
			name:     "Fun tone",
			tone:     "fun",
			expected: "- Use a fun, playful tone with emojis and creative language while keeping it professional",
		},
		{
			name:     "Pirate tone",
			tone:     "pirate",
			expected: "- Write the commit message in pirate speak with nautical terminology (e.g., 'Hoist', 'Plunder', 'Navigate')",
		},
		{
			name:     "Haiku tone",
			tone:     "haiku",
			expected: "- Write the commit message as a single-line haiku with 5-7-5 syllable structure, separating each line with ' / ', capturing the essence of the code change",
		},
		{
			name:     "Serious tone",
			tone:     "serious",
			expected: "- Use a very serious, formal tone with technical precision and no casual language",
		},
		{
			name:     "Unknown tone (treated as custom)",
			tone:     "unknown",
			expected: "- Use a unknown tone for the commit message",
		},
		{
			name:     "Empty tone (fallback to custom)",
			tone:     "",
			expected: "- Use a  tone for the commit message",
		},
		{
			name:     "Custom tone - casual",
			tone:     "casual",
			expected: "- Use a casual tone for the commit message",
		},
		{
			name:     "Custom tone - dramatic",
			tone:     "dramatic",
			expected: "- Use a dramatic tone for the commit message",
		},
		{
			name:     "Custom tone - technical",
			tone:     "technical",
			expected: "- Use a technical tone for the commit message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getToneInstruction(tt.tone)
			if result != tt.expected {
				t.Errorf("getToneInstruction(%q) = %q, want %q", tt.tone, result, tt.expected)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name          string
		baseURL       string
		model         string
		expectedURL   string
		expectedModel string
	}{
		{
			name:          "Default values",
			baseURL:       "",
			model:         "",
			expectedURL:   "http://localhost:11434",
			expectedModel: "llama3.2",
		},
		{
			name:          "Custom URL and model",
			baseURL:       "http://custom-server:8080",
			model:         "custom-model",
			expectedURL:   "http://custom-server:8080",
			expectedModel: "custom-model",
		},
		{
			name:          "Custom URL only",
			baseURL:       "http://remote-ollama:9999",
			model:         "",
			expectedURL:   "http://remote-ollama:9999",
			expectedModel: "llama3.2",
		},
		{
			name:          "Custom model only",
			baseURL:       "",
			model:         "codellama",
			expectedURL:   "http://localhost:11434",
			expectedModel: "codellama",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.baseURL, tt.model)

			if client.BaseURL != tt.expectedURL {
				t.Errorf("NewClient().BaseURL = %q, want %q", client.BaseURL, tt.expectedURL)
			}

			if client.Model != tt.expectedModel {
				t.Errorf("NewClient().Model = %q, want %q", client.Model, tt.expectedModel)
			}

			if client.client == nil {
				t.Error("NewClient().client should not be nil")
			}
		})
	}
}
