// Package stt runs whisper.cpp's `whisper-server` as a local sidecar and
// transcribes audio by posting it over HTTP. Keeping the model resident in the
// sidecar means each utterance skips the ~150 MB model load — the same pattern
// we'll use for the Ollama LLM in Phase 2.
package stt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"voice2prompt/internal/resource"
)

// DefaultServerBin is where scripts/setup-whisper.sh places the built binary.
const DefaultServerBin = "third_party/whisper.cpp/build/bin/whisper-server"

// Server is a running whisper-server process plus an HTTP client to talk to it.
type Server struct {
	cmd     *exec.Cmd
	baseURL string
	client  *http.Client
}

// Config controls how the sidecar is launched.
type Config struct {
	BinPath   string // path to whisper-server; falls back to $PROMPTVOICE_WHISPER_SERVER or DefaultServerBin
	ModelPath string // path to a ggml model
	Host      string // default 127.0.0.1
	Port      int    // default 8642
	Language  string // default "en"
}

func (c *Config) withDefaults() {
	if c.BinPath == "" {
		if env := os.Getenv("PROMPTVOICE_WHISPER_SERVER"); env != "" {
			c.BinPath = env
		} else {
			c.BinPath = DefaultServerBin
		}
	}
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	// Port 0 means "pick a free ephemeral port" — avoids clashing with a stale
	// server that a previous crash may have left running.
	if c.Language == "" {
		c.Language = "en"
	}
}

// freePort asks the OS for an unused TCP port on the loopback interface.
func freePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// Start launches whisper-server and blocks until it is ready to serve (or times out).
func Start(cfg Config) (*Server, error) {
	cfg.withDefaults()

	// Resolve paths against the repo (dev) or the .app bundle Resources (packaged).
	if r, ok := resource.Find(cfg.BinPath); ok {
		cfg.BinPath = r
	}
	if r, ok := resource.Find(cfg.ModelPath); ok {
		cfg.ModelPath = r
	}

	if _, err := os.Stat(cfg.BinPath); err != nil {
		return nil, fmt.Errorf("whisper-server binary not found at %q — run scripts/setup-whisper.sh: %w", cfg.BinPath, err)
	}
	if _, err := os.Stat(cfg.ModelPath); err != nil {
		return nil, fmt.Errorf("model not found at %q: %w", cfg.ModelPath, err)
	}

	if cfg.Port == 0 {
		p, err := freePort()
		if err != nil {
			return nil, fmt.Errorf("could not find a free port: %w", err)
		}
		cfg.Port = p
	}

	cmd := exec.Command(cfg.BinPath,
		"-m", cfg.ModelPath,
		"--host", cfg.Host,
		"--port", fmt.Sprintf("%d", cfg.Port),
		"-l", cfg.Language,
	)
	// Surface server logs on our stderr so failures are visible.
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start whisper-server: %w", err)
	}

	s := &Server{
		cmd:     cmd,
		baseURL: fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
		client:  &http.Client{Timeout: 60 * time.Second},
	}

	if err := s.waitReady(30 * time.Second); err != nil {
		_ = s.Close()
		return nil, err
	}
	return s, nil
}

// waitReady polls the server root until it accepts connections.
func (s *Server) waitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/", nil)
		resp, err := s.client.Do(req)
		cancel()
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}
	return fmt.Errorf("whisper-server did not become ready within %s", timeout)
}

// Result is one transcription plus the round-trip latency (HTTP + inference),
// which is what the app actually experiences.
type Result struct {
	Text    string
	Latency time.Duration
}

// Transcribe sends WAV-encoded audio to the sidecar and returns the transcript.
func (s *Server) Transcribe(wav []byte) (Result, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	part, err := w.CreateFormFile("file", "audio.wav")
	if err != nil {
		return Result{}, err
	}
	if _, err := part.Write(wav); err != nil {
		return Result{}, err
	}
	// Ask for JSON and deterministic decoding.
	_ = w.WriteField("response_format", "json")
	_ = w.WriteField("temperature", "0.0")
	w.Close()

	req, err := http.NewRequest(http.MethodPost, s.baseURL+"/inference", &body)
	if err != nil {
		return Result{}, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	start := time.Now()
	resp, err := s.client.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("inference request failed: %w", err)
	}
	defer resp.Body.Close()
	latency := time.Since(start)

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("whisper-server returned %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	// The server returns {"text": "..."} in JSON mode; fall back to raw text.
	var parsed struct {
		Text string `json:"text"`
	}
	text := ""
	if err := json.Unmarshal(raw, &parsed); err == nil && parsed.Text != "" {
		text = parsed.Text
	} else {
		text = string(raw)
	}

	// whisper separates segments with newlines but each segment already carries its
	// own leading space, so we strip the newlines (not replace with a space — that
	// would break words split across a segment boundary), then normalize spacing.
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\r", "")
	text = stripNonSpeech(text)
	text = strings.Join(strings.Fields(text), " ")

	return Result{Text: text, Latency: latency}, nil
}

// nonSpeech matches whisper's bracketed/parenthesized non-speech annotations
// (e.g. "[BLANK_AUDIO]", "[silence]", "(music)", "[ Inaudible ]") which must never
// be typed into the user's document.
var nonSpeech = regexp.MustCompile(
	`(?i)[\[\(]\s*(blank[_ ]?audio|silence|music|noise|inaudible|pause|applause|laughter|sound|no speech|beep|click)[^\]\)]*[\]\)]`)

// stripNonSpeech removes non-speech markers and musical-note runs. If the whole
// utterance was such a marker, this yields "" — which the engine treats as no speech.
func stripNonSpeech(s string) string {
	s = nonSpeech.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "♪", "")
	return strings.TrimSpace(s)
}

// Close stops the sidecar process.
func (s *Server) Close() error {
	if s.cmd == nil || s.cmd.Process == nil {
		return nil
	}
	_ = s.cmd.Process.Kill()
	_, _ = s.cmd.Process.Wait()
	return nil
}
