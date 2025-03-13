package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Model defines the interface for LLM interaction
type Model interface {
	// Generate produces a response based on the given prompt
	Generate(prompt string) (string, error)
}

// Trainable is an optional interface for models that support training
type Trainable interface {
	// Train trains the model using the provided dataset
	Train(datasetPath, outputPath string) error
}

// OllamaRequest represents the request structure for Ollama API
type OllamaRequest struct {
	Model     string  `json:"model"`
	Prompt    string  `json:"prompt"`
	Stream    bool    `json:"stream,omitempty"`
	MaxTokens int     `json:"max_tokens,omitempty"`
	Temp      float64 `json:"temperature,omitempty"`
}

// OllamaResponse represents the response structure from Ollama API
type OllamaResponse struct {
	Model              string  `json:"model"`
	Response           string  `json:"response"`
	Done               bool    `json:"done"`
	TotalDuration      int64   `json:"total_duration,omitempty"`
	LoadDuration       int64   `json:"load_duration,omitempty"`
	PromptEvalDuration int64   `json:"prompt_eval_duration,omitempty"`
	EvalDuration       int64   `json:"eval_duration,omitempty"`
	EvalCount          int     `json:"eval_count,omitempty"`
	Error              *string `json:"error,omitempty"`
}

// BaseModel provides common functionality for all model implementations
type BaseModel struct {
	endpoint    string
	modelName   string
	maxTokens   int
	temperature float64
}

// Generate implements the Model interface for the BaseModel
func (b *BaseModel) Generate(prompt string) (string, error) {
	// Create request
	reqBody := OllamaRequest{
		Model:     b.modelName,
		Prompt:    prompt,
		Stream:    false,
		MaxTokens: b.maxTokens,
		Temp:      b.temperature,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send request to Ollama
	url := fmt.Sprintf("%s/generate", b.endpoint)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(respBytes, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if ollamaResp.Error != nil {
		return "", fmt.Errorf("Ollama error: %s", *ollamaResp.Error)
	}

	return ollamaResp.Response, nil
}
