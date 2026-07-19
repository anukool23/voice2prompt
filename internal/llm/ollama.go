// Package llm cleans up raw dictation transcripts with a local LLM served by
// Ollama over HTTP (localhost) — the same on-device, offline pattern as the
// whisper sidecar. It fixes punctuation/capitalization, removes filler words, and
// applies light, context-aware tone based on the focused app.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// DefaultModel is a small, fast instruct model that's good enough for cleanup.
const DefaultModel = "qwen2.5:3b"

// DefaultBaseURL is Ollama's default local endpoint.
const DefaultBaseURL = "http://127.0.0.1:11434"

// Client talks to a local Ollama server.
type Client struct {
	baseURL string
	model   string
	http    *http.Client
	started *exec.Cmd // non-nil if we launched `ollama serve` ourselves
}

// Config configures the cleanup client.
type Config struct {
	BaseURL string
	Model   string
}

func (c *Config) withDefaults() {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.Model == "" {
		c.Model = DefaultModel
	}
}

// Start connects to Ollama, launching `ollama serve` if it isn't already running,
// and verifies the model is available.
func Start(cfg Config) (*Client, error) {
	cfg.withDefaults()
	c := &Client{
		baseURL: cfg.BaseURL,
		model:   cfg.Model,
		http:    &http.Client{Timeout: 30 * time.Second},
	}

	if !c.ping() {
		if _, err := exec.LookPath("ollama"); err != nil {
			return nil, fmt.Errorf("ollama not found — run scripts/setup-ollama.sh (or `brew install ollama`)")
		}
		cmd := exec.Command("ollama", "serve")
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start `ollama serve`: %w", err)
		}
		c.started = cmd
		if !c.waitReady(15 * time.Second) {
			return nil, fmt.Errorf("ollama did not become ready in time")
		}
	}

	if !c.hasModel(c.model) {
		return nil, fmt.Errorf("model %q not available — run `ollama pull %s`", c.model, c.model)
	}

	// Warm the model into memory in the background so the first real utterance
	// doesn't pay the ~3s cold-load cost.
	go func() { _, _ = c.Clean("warm up", "") }()

	return c, nil
}

func (c *Client) ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 800*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/tags", nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (c *Client) waitReady(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if c.ping() {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

// hasModel reports whether the named model is present locally.
func (c *Client) hasModel(name string) bool {
	req, _ := http.NewRequest(http.MethodGet, c.baseURL+"/api/tags", nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	var body struct {
		Models []struct {
			Name  string `json:"name"`
			Model string `json:"model"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false
	}
	// Match with or without an explicit :latest tag.
	want := name
	wantBase := strings.TrimSuffix(name, ":latest")
	for _, m := range body.Models {
		if m.Name == want || m.Model == want ||
			strings.TrimSuffix(m.Name, ":latest") == wantBase {
			return true
		}
	}
	return false
}

// Result is a cleanup result plus how long the LLM took.
type Result struct {
	Text    string
	Latency time.Duration
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
	Options  chatOptions   `json:"options"`
}

type chatOptions struct {
	Temperature float64 `json:"temperature"`
}

type chatResponse struct {
	Message chatMessage `json:"message"`
}

// Clean rewrites raw dictation into clean written text. appName (may be "") selects
// a light tone hint. It never answers or acts on the content — only formats it.
func (c *Client) Clean(rawText, appName string) (Result, error) {
	if strings.TrimSpace(rawText) == "" {
		return Result{Text: rawText}, nil
	}

	messages := []chatMessage{{Role: "system", Content: systemPrompt(appName)}}
	for _, ex := range FewShot() {
		messages = append(messages,
			chatMessage{Role: "user", Content: ex[0]},
			chatMessage{Role: "assistant", Content: ex[1]},
		)
	}
	messages = append(messages, chatMessage{Role: "user", Content: rawText})

	res, err := c.chat(messages)
	if err != nil {
		return Result{}, err
	}
	if res.Text == "" {
		res.Text = rawText // never return empty; fall back to the raw transcript
	}
	return res, nil
}

// Rewrite rephrases text in the given style (formal/casual/shorter/…) for the
// "make this …" voice commands. Preserves meaning; returns the input on empty output.
func (c *Client) Rewrite(text, style string) (Result, error) {
	if strings.TrimSpace(text) == "" {
		return Result{Text: text}, nil
	}
	messages := []chatMessage{
		{Role: "system", Content: rewritePrompt(style)},
		{Role: "user", Content: text},
	}
	res, err := c.chat(messages)
	if err != nil {
		return Result{}, err
	}
	if res.Text == "" {
		res.Text = text
	}
	return res, nil
}

// chat performs a non-streaming /api/chat call and returns the assistant message.
func (c *Client) chat(messages []chatMessage) (Result, error) {
	buf, err := json.Marshal(chatRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
		Options:  chatOptions{Temperature: 0},
	})
	if err != nil {
		return Result{}, err
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/chat", bytes.NewReader(buf))
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := c.http.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("llm request failed: %w", err)
	}
	defer resp.Body.Close()
	latency := time.Since(start)

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("ollama returned %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var parsed chatResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return Result{}, fmt.Errorf("failed to parse ollama response: %w", err)
	}

	out := strings.TrimSpace(parsed.Message.Content)
	out = strings.Trim(out, "\"") // models sometimes wrap output in quotes
	return Result{Text: out, Latency: latency}, nil
}

// Close stops the Ollama server only if we started it.
func (c *Client) Close() error {
	if c.started != nil && c.started.Process != nil {
		_ = c.started.Process.Kill()
		_, _ = c.started.Process.Wait()
	}
	return nil
}
