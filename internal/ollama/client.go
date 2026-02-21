package ollama

import (
	"bufio"
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
	DefaultModel   = "llama3"
)

type Client struct {
	BaseURL string
	Model   string
	HTTP    *http.Client
}

// NewClient creates a new Ollama client. If modelOverride is non-empty it takes
// priority over the SSAGE_MODEL env var and the built-in default.
func NewClient(modelOverride string) *Client {
	model := modelOverride
	if model == "" {
		model = os.Getenv("SSAGE_MODEL")
	}
	if model == "" {
		model = DefaultModel
	}
	return &Client{
		BaseURL: DefaultBaseURL,
		Model:   model,
		HTTP: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for streaming
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

// Generate sends a prompt and waits for the full response (non-streaming).
// Kept for use in tests and stats.
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

	if err := checkStatus(resp, c.Model); err != nil {
		return "", err
	}

	var genResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return genResp.Response, nil
}

// GenerateStream sends a prompt and calls onChunk for every token received,
// allowing the caller to print text as it arrives. It returns the full
// accumulated response string so callers can use it (e.g. for clipboard copy).
func (c *Client) GenerateStream(prompt string, onChunk func(token string)) (string, error) {
	reqBody := GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Stream: true,
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

	if err := checkStatus(resp, c.Model); err != nil {
		return "", err
	}

	var full string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var chunk GenerateResponse
		if err := json.Unmarshal(line, &chunk); err != nil {
			continue // skip malformed lines
		}
		if chunk.Response != "" {
			onChunk(chunk.Response)
			full += chunk.Response
		}
		if chunk.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		return full, fmt.Errorf("error reading stream: %w", err)
	}

	return full, nil
}

// checkStatus returns a descriptive error for non-200 HTTP responses.
func checkStatus(resp *http.Response, model string) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("model '%s' not found. Please run 'ollama pull %s' to download it", model, model)
	}
	return fmt.Errorf("ollama API returned status %d: %s", resp.StatusCode, string(body))
}
