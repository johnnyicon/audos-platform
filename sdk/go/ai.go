package audos

// Standalone AI client — leverage a workspace's ai-api hook (OpenAI text) from any Go app.
// Verified live 2026-08-10.
//
// Separate from Client in client.go on purpose: Client carries the db-api APIKey, whereas the
// AI hook needs no key. A Go app that only wants AI configures NewAI with just a workspace id.
//
// IMPORTANT — ai-api is a *workspace-built server function*, not a platform feature. Verified
// 2026-08-10: it exists in the Throughline workspace but returns "Hook 'ai-api' not found" in
// other workspaces (including freshly-onboarded ones). So creating a workspace is necessary but
// NOT sufficient — the workspace must actually have an ai-api hook. This client turns that 404
// into an explicit, actionable error.
//
// Provider/model reality (verified 2026-08-10): OpenAI only. Working: gpt-4o, gpt-4o-mini
// (default), gpt-4.1. NOT available: gpt-5 / o-series (reasoning), Claude, Gemini. Reasoning/
// effort/temperature/tools/response_format params are silently ignored. Video/voice are a
// different surface (Otto tools), not this hook.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ─── Config & types ───────────────────────────────────────────────────────────

type AIConfig struct {
	// WorkspaceID: "351699" or "workspace-351699"; the prefix is stripped. Required.
	WorkspaceID string
	BaseURL     string // default "https://audos.com"
	HookName    string // default "ai-api"
	DefaultModel string // default "gpt-4o-mini"
	// HookSecret is sent as x-hook-secret. NOTE: the reference ai-api hook does not enforce it
	// (verified 2026-08-10). Pass-through only, for a hook you build that checks it. Not auth.
	HookSecret string
	HTTPClient *http.Client // optional; defaults to http.DefaultClient
}

// ChatMessage.Content is either a plain string, or []ContentPart for vision. Image understanding
// is verified working on gpt-4o/gpt-4.1 through this hook (2026-08-10): pass an image_url part with
// a public URL or a data:image/...;base64,... URI. Use TextMessage / ImageMessage to build them.
type ChatMessage struct {
	Role    string `json:"role"` // "user" | "assistant"
	Content any    `json:"content"`
}

// ContentPart is one element of a multimodal message: {"type":"text",...} or {"type":"image_url",...}.
type ContentPart struct {
	Type     string        `json:"type"`
	Text     string        `json:"text,omitempty"`
	ImageURL *ContentImage `json:"image_url,omitempty"`
}

type ContentImage struct {
	URL string `json:"url"`
}

// TextMessage builds a plain user text message.
func TextMessage(text string) ChatMessage {
	return ChatMessage{Role: "user", Content: text}
}

// ImageMessage builds a user message pairing a prompt with one or more images (URLs or data URIs).
func ImageMessage(text string, imageURLs ...string) ChatMessage {
	parts := []ContentPart{{Type: "text", Text: text}}
	for _, u := range imageURLs {
		parts = append(parts, ContentPart{Type: "image_url", ImageURL: &ContentImage{URL: u}})
	}
	return ChatMessage{Role: "user", Content: parts}
}

type ChatOptions struct {
	Model        string
	SystemPrompt string
	MaxTokens    int // honored by Chat (unlike Generate); defaults to 1200
	Temperature  *float64
}

type AIUsage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}

type AIResult struct {
	Text  string
	Model string
	Usage *AIUsage
}

// AIError carries the HTTP status and the server's message for any non-success outcome.
type AIError struct {
	Msg         string
	Status      int
	ServerError string
}

func (e *AIError) Error() string { return e.Msg }

// AudosAIModels records the model landscape behind the reference ai-api hook, from an explicit
// 43-name live audit (2026-08-10). See the TS SDK's AUDOS_AI_MODELS for the full explanation.
var AudosAIModels = struct {
	Working                 []string
	Vision                  []string
	ReachableButHookBlocked []string
	Unavailable             []string
}{
	Working:                 []string{"gpt-4o", "gpt-4o-mini", "gpt-4o-2024-08-06", "gpt-4.1", "gpt-4.1-mini", "gpt-4.1-nano", "gpt-4-turbo", "gpt-4", "gpt-3.5-turbo"},
	Vision:                  []string{"gpt-4o", "gpt-4o-mini", "gpt-4.1", "gpt-4.1-mini", "gpt-4.1-nano"},
	ReachableButHookBlocked: []string{"gpt-5", "gpt-5-mini", "gpt-5-nano", "gpt-5.5", "gpt-5.6", "o1", "o3", "o3-mini", "o4-mini"},
	Unavailable:             []string{"claude-*", "fable-5", "gemini-*", "o1-mini", "gpt-5.6-sol"},
}

// ─── Client ─────────────────────────────────────────────────────────────────────

type AIClient struct {
	cfg  AIConfig
	http *http.Client
	url  string
}

func NewAI(cfg AIConfig) *AIClient {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://audos.com"
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")
	if cfg.HookName == "" {
		cfg.HookName = "ai-api"
	}
	if cfg.DefaultModel == "" {
		cfg.DefaultModel = "gpt-4o-mini"
	}
	ws := strings.TrimPrefix(cfg.WorkspaceID, "workspace-")
	hc := cfg.HTTPClient
	if hc == nil {
		hc = http.DefaultClient
	}
	return &AIClient{
		cfg:  cfg,
		http: hc,
		url:  fmt.Sprintf("%s/api/hooks/execute/workspace-%s/%s", cfg.BaseURL, ws, cfg.HookName),
	}
}

// Endpoint returns the resolved URL, for logging/debugging.
func (a *AIClient) Endpoint() string { return a.url }

// rawUsage tolerates both shapes: generate returns camelCase, chat passes OpenAI's snake_case.
type rawUsage struct {
	PromptCamel     *int `json:"promptTokens"`
	CompletionCamel *int `json:"completionTokens"`
	TotalCamel      *int `json:"totalTokens"`
	PromptSnake     *int `json:"prompt_tokens"`
	CompletionSnake *int `json:"completion_tokens"`
	TotalSnake      *int `json:"total_tokens"`
}

// ReasoningMinMaxTokens is the floor used for reasoning models (GPT-5 / o-series). They spend part of
// the budget on *internal* reasoning before emitting visible text, and that counts against
// max_completion_tokens. Measured 2026-08-14: gpt-5 on a moderately complex prompt burned all 2000
// tokens reasoning and returned empty text; at 4000 it used ~1150 reasoning plus real output.
const ReasoningMinMaxTokens = 4000

const defaultMaxTokens = 1200

// NeedsMaxCompletionTokens reports whether a model requires max_completion_tokens (GPT-5 + o-series).
func NeedsMaxCompletionTokens(model string) bool {
	return strings.HasPrefix(model, "gpt-5") ||
		(len(model) > 1 && model[0] == 'o' && model[1] >= '1' && model[1] <= '9')
}

type hookResponse struct {
	Text         string          `json:"text"`
	Model        string          `json:"model"`
	Usage        json.RawMessage `json:"usage"`
	Error        string          `json:"error"`
	FinishReason string          `json:"finishReason"`
	Meta         struct {
		Logs []string `json:"logs"`
	} `json:"_meta"`
}

// logDetail pulls the real error out of _meta.logs — the hook hides it there behind the
// generic top-level "AI generation failed" message.
func (r *hookResponse) logDetail() string {
	for _, l := range r.Meta.Logs {
		low := strings.ToLower(l)
		if strings.Contains(low, "error") || strings.Contains(low, "does not exist") || strings.Contains(low, "not found") {
			return strings.TrimPrefix(l, "[ERROR] ")
		}
	}
	return ""
}

func pick(a, b *int) int {
	if a != nil {
		return *a
	}
	if b != nil {
		return *b
	}
	return 0
}

func normalizeUsage(raw json.RawMessage) *AIUsage {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var u rawUsage
	if err := json.Unmarshal(raw, &u); err != nil {
		return nil
	}
	if u.PromptCamel == nil && u.PromptSnake == nil &&
		u.CompletionCamel == nil && u.CompletionSnake == nil &&
		u.TotalCamel == nil && u.TotalSnake == nil {
		return nil
	}
	prompt := pick(u.PromptCamel, u.PromptSnake)
	completion := pick(u.CompletionCamel, u.CompletionSnake)
	total := pick(u.TotalCamel, u.TotalSnake)
	if total == 0 {
		total = prompt + completion
	}
	return &AIUsage{PromptTokens: prompt, CompletionTokens: completion, TotalTokens: total}
}

func (a *AIClient) post(body any) (*hookResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, a.url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if a.cfg.HookSecret != "" {
		req.Header.Set("x-hook-secret", a.cfg.HookSecret)
	}

	resp, err := a.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var parsed hookResponse
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &parsed) // non-JSON bodies handled by status check below
	}

	if resp.StatusCode >= 400 {
		// Reasoning models can spend the entire budget on internal reasoning and return no visible
		// text. Give that its own actionable message instead of a bare "no text".
		if parsed.FinishReason == "length" && parsed.Text == "" {
			return nil, &AIError{
				Msg: fmt.Sprintf("model produced no visible text: the whole token budget went to "+
					"internal reasoning (finish_reason: length). Raise MaxTokens — reasoning models "+
					"need at least ~%d, and more for complex prompts", ReasoningMinMaxTokens),
				Status:      resp.StatusCode,
				ServerError: parsed.Error,
			}
		}
		// The hook flattens its top-level error to a generic string; the real upstream reason
		// (e.g. OpenAI's "model … does not exist or you do not have access") is in _meta.logs.
		serverErr := parsed.logDetail()
		notFoundSignal := parsed.Error
		if serverErr == "" {
			serverErr = parsed.Error
		}
		if serverErr == "" {
			serverErr = string(raw)
		}
		if resp.StatusCode == http.StatusNotFound && strings.Contains(strings.ToLower(notFoundSignal), "not found") {
			ws := strings.TrimPrefix(a.cfg.WorkspaceID, "workspace-")
			return nil, &AIError{
				Msg: fmt.Sprintf("hook %q not found in workspace-%s: the ai-api hook is a "+
					"workspace-built server function, not a platform feature — this workspace "+
					"doesn't have one. Point at a workspace that does, or build the hook there.",
					a.cfg.HookName, ws),
				Status:      resp.StatusCode,
				ServerError: serverErr,
			}
		}
		return nil, &AIError{
			Msg:         fmt.Sprintf("audos AI call failed (HTTP %d): %s", resp.StatusCode, serverErr),
			Status:      resp.StatusCode,
			ServerError: serverErr,
		}
	}

	if parsed.Error != "" {
		detail := parsed.logDetail()
		if detail == "" {
			detail = parsed.Error
		}
		return nil, &AIError{Msg: "audos AI error: " + detail, Status: resp.StatusCode, ServerError: detail}
	}
	return &parsed, nil
}

func (a *AIClient) model(m string) string {
	if m != "" {
		return m
	}
	return a.cfg.DefaultModel
}

// Chat is the reliable workhorse. It honors MaxTokens; use it for anything that may run long.
// messages must be non-empty and end with a user turn.
func (a *AIClient) Chat(messages []ChatMessage, opts *ChatOptions) (*AIResult, error) {
	if len(messages) == 0 || messages[len(messages)-1].Role != "user" {
		return nil, &AIError{Msg: "Chat: messages must be non-empty and end with a user turn", Status: 0}
	}
	if opts == nil {
		opts = &ChatOptions{}
	}
	model := a.model(opts.Model)
	maxTokens := opts.MaxTokens
	if maxTokens == 0 {
		// Reasoning models need headroom for internal reasoning tokens; the normal default
		// silently returns empty text on anything non-trivial.
		if NeedsMaxCompletionTokens(model) {
			maxTokens = ReasoningMinMaxTokens
		} else {
			maxTokens = defaultMaxTokens
		}
	}
	body := map[string]any{
		"action":    "chat",
		"messages":  messages,
		"model":     model,
		"maxTokens": maxTokens,
	}
	if opts.SystemPrompt != "" {
		body["systemPrompt"] = opts.SystemPrompt
	}
	if opts.Temperature != nil {
		body["temperature"] = *opts.Temperature
	}
	d, err := a.post(body)
	if err != nil {
		return nil, err
	}
	return &AIResult{Text: d.Text, Model: d.Model, Usage: normalizeUsage(d.Usage)}, nil
}

// Complete is the method to reach for by default: a single prompt wrapped as a chat turn, so it
// is NOT subject to the generate action's 1000-token cap.
func (a *AIClient) Complete(prompt string, opts *ChatOptions) (*AIResult, error) {
	return a.Chat([]ChatMessage{{Role: "user", Content: prompt}}, opts)
}

// Generate uses the single-turn generate action. Kept for parity, but two hard limits verified
// on the reference hook: (1) output capped at 1000 completion tokens regardless of request;
// (2) on model failure it returns HTTP 200 with empty text and no error — which this converts
// into an AIError. Prefer Complete/Chat.
func (a *AIClient) Generate(prompt, systemPrompt, model string) (*AIResult, error) {
	body := map[string]any{"action": "generate", "prompt": prompt, "model": a.model(model)}
	if systemPrompt != "" {
		body["systemPrompt"] = systemPrompt
	}
	d, err := a.post(body)
	if err != nil {
		return nil, err
	}
	if d.Text == "" {
		return nil, &AIError{
			Msg: fmt.Sprintf("Generate returned empty text for model %q: the generate action "+
				"reports success even when the model call fails — treat as a failed generation "+
				"(often an unavailable model). Use Complete/Chat with gpt-4o, gpt-4o-mini, or gpt-4.1.",
				a.model(model)),
			Status: 200,
		}
	}
	return &AIResult{Text: d.Text, Model: d.Model, Usage: normalizeUsage(d.Usage)}, nil
}
