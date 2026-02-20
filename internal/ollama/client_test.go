package ollama

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestServer creates a local HTTP test server that mimics the Ollama API.
func newTestServer(t *testing.T, response string, statusCode int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if response != "" {
			resp := GenerateResponse{Response: response, Done: true}
			_ = json.NewEncoder(w).Encode(resp)
		} else {
			_, _ = w.Write([]byte(`{"error":"model not found"}`))
		}
	}))
}

// TestGenerate_HappyPath verifies the client correctly receives an AI response.
func TestGenerate_HappyPath(t *testing.T) {
	srv := newTestServer(t, "This command lists files.", http.StatusOK)
	defer srv.Close()

	client := &Client{
		BaseURL: srv.URL,
		Model:   "testmodel",
		HTTP:    &http.Client{},
	}

	result, err := client.Generate("Explain ls -la")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "This command lists files." {
		t.Errorf("unexpected response: %q", result)
	}
}

// TestGenerate_ModelNotFound verifies that a 404 returns a helpful error message.
func TestGenerate_ModelNotFound(t *testing.T) {
	srv := newTestServer(t, "", http.StatusNotFound)
	defer srv.Close()

	client := &Client{
		BaseURL: srv.URL,
		Model:   "missing-model",
		HTTP:    &http.Client{},
	}

	_, err := client.Generate("Explain ls")
	if err == nil {
		t.Fatal("expected an error for 404, got nil")
	}
	// Verify the error message includes the model name and pull instruction
	errMsg := err.Error()
	if !contains(errMsg, "missing-model") {
		t.Errorf("error should mention the model name, got: %s", errMsg)
	}
	if !contains(errMsg, "ollama pull") {
		t.Errorf("error should mention 'ollama pull', got: %s", errMsg)
	}
}

// TestNewClient_ModelPriority verifies that the modelOverride takes precedence.
func TestNewClient_ModelPriority(t *testing.T) {
	t.Setenv("SSAGE_MODEL", "env-model")

	// When override is provided, it should win
	c := NewClient("override-model")
	if c.Model != "override-model" {
		t.Errorf("expected 'override-model', got %q", c.Model)
	}

	// When override is empty, env var should win
	c2 := NewClient("")
	if c2.Model != "env-model" {
		t.Errorf("expected 'env-model', got %q", c2.Model)
	}
}

// TestNewClient_DefaultModel verifies that llama3 is used when no override/env is set.
func TestNewClient_DefaultModel(t *testing.T) {
	t.Setenv("SSAGE_MODEL", "")
	c := NewClient("")
	if c.Model != DefaultModel {
		t.Errorf("expected default model %q, got %q", DefaultModel, c.Model)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsRune(s, substr))
}

func containsRune(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
