package ollama

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	DefaultBaseURL = "http://localhost:11434"
	DefaultModel   = "llama3" // Functionality to switch model can be added later
)

type Client struct {
	BaseURL string
	Model   string
	HTTP    *http.Client
}

func NewClient() *Client {
	model := os.Getenv("SSAGE_MODEL")
	if model == "" {
		model = DefaultModel
	}
	return &Client{
		BaseURL: DefaultBaseURL,
		Model:   model,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Model    string `json:"model"`
	Created  string `json:"created_at"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (c *Client) Generate(prompt string) (string, error) {
	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.HTTP.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusNotFound {
			return "", fmt.Errorf("model '%s' not found. Please run 'ollama pull %s' to download it", c.Model, c.Model)
		}
		return "", fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return genResp.Response, nil
}
