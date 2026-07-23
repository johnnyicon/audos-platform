package audos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const onboardBaseURL = "https://audos.com/api/agent/onboard"

// OnboardClient calls the Audos onboarding/Otto chat API.
//
// The base path (/api/agent/onboard) is entirely separate from the workspace
// hooks API used by Client.
//
// Production note, 2026-07-08: existing-workspace access currently works
// through body-auth endpoints:
//   - POST /start with workspaceId + createNew=false
//   - POST /status with { authToken }
//   - POST /chat with { authToken, message, chatId? }
//
// Bearer-header routes are retained as explicit *ViaBearer methods for
// compatibility checks, but they may return 401 until Audos publishes the
// bearer-header fix.
type OnboardClient struct {
	authToken string
	baseURL   string
	http      *http.Client
}

// NewOnboardClient creates a client pre-loaded with an auth token — use this
// when you already have a token from a previous Start/Verify flow.
func NewOnboardClient(authToken string) *OnboardClient {
	return &OnboardClient{
		authToken: authToken,
		baseURL:   onboardBaseURL,
		http:      &http.Client{},
	}
}

// NewOnboardClientUnauthenticated creates a client with no token — use this
// to begin the Start/Verify OTP flow before a token exists. Call
// SetAuthToken after a successful Verify to reuse the same client.
func NewOnboardClientUnauthenticated() *OnboardClient {
	return &OnboardClient{
		baseURL: onboardBaseURL,
		http:    &http.Client{},
	}
}

// SetAuthToken stores the auth token returned by Verify so the same client
// instance can proceed to Status, Chat, and Rebuild without re-initialising.
func (c *OnboardClient) SetAuthToken(token string) {
	c.authToken = token
}

// ─── Types ────────────────────────────────────────────────────────────────────

// WorkspaceURLs holds the set of URLs returned for a provisioned workspace.
type WorkspaceURLs struct {
	LandingPage string `json:"landingPage"`
	Workspace   string `json:"workspace"`
	Dashboard   string `json:"dashboard,omitempty"`
	AppSpace    string `json:"appSpace,omitempty"`
}

// StartOption allows optional fields to be passed to Start.
type StartOption func(*startRequest)

type startRequest struct {
	Email          string `json:"email"`
	BusinessIdea   string `json:"businessIdea"`
	WorkspaceID    string `json:"workspaceId,omitempty"`
	BusinessName   string `json:"businessName,omitempty"`
	TargetCustomer string `json:"targetCustomer,omitempty"`
	CallbackURL    string `json:"callbackUrl,omitempty"`
	CreateNew      *bool  `json:"createNew,omitempty"`
}

// WithBusinessName sets the optional business name on a Start call.
func WithBusinessName(name string) StartOption {
	return func(r *startRequest) { r.BusinessName = name }
}

// WithTargetCustomer sets the optional target customer on a Start call.
func WithTargetCustomer(customer string) StartOption {
	return func(r *startRequest) { r.TargetCustomer = customer }
}

// WithCallbackURL sets the optional callback URL on a Start call.
func WithCallbackURL(url string) StartOption {
	return func(r *startRequest) { r.CallbackURL = url }
}

// WithWorkspaceID targets an existing workspace and defaults createNew to false.
func WithWorkspaceID(workspaceID string) StartOption {
	return func(r *startRequest) {
		r.WorkspaceID = workspaceID
		createNew := false
		r.CreateNew = &createNew
	}
}

// WithCreateNew controls whether /start should create a new workspace.
// Calling WithCreateNew() preserves the legacy behavior of forcing creation.
// Calling WithCreateNew(false) is useful when targeting an existing workspace.
func WithCreateNew(value ...bool) StartOption {
	return func(r *startRequest) {
		createNew := true
		if len(value) > 0 {
			createNew = value[0]
		}
		r.CreateNew = &createNew
	}
}

// StartResponse covers both outcomes of /start:
//   - New user: SessionToken is set; AuthToken/WorkspaceID/WorkspaceName/URLs are empty.
//   - Returning user: AuthToken + workspace fields are set; SessionToken is empty.
type StartResponse struct {
	// Returning user fields
	AuthTokenCamel    string            `json:"authToken,omitempty"`
	AuthToken         string            `json:"auth_token,omitempty"`
	WorkspaceID       string            `json:"workspaceId,omitempty"`
	WorkspaceName     string            `json:"workspaceName,omitempty"`
	URLs              WorkspaceURLs     `json:"urls,omitempty"`
	ExistingWorkspace ExistingWorkspace `json:"existingWorkspace,omitempty"`

	// New user field — present only when OTP email was sent
	SessionToken string `json:"sessionToken,omitempty"`
}

type ExistingWorkspace struct {
	AuthTokenCamel string        `json:"authToken,omitempty"`
	AuthToken      string        `json:"auth_token,omitempty"`
	WorkspaceID    string        `json:"workspaceId,omitempty"`
	WorkspaceName  string        `json:"workspaceName,omitempty"`
	URLs           WorkspaceURLs `json:"urls,omitempty"`
}

// Token returns the auth token from either the current or historical /start
// response shape. Do not log or persist this value outside a secret store.
func (r StartResponse) Token() string {
	if r.AuthTokenCamel != "" {
		return r.AuthTokenCamel
	}
	if r.AuthToken != "" {
		return r.AuthToken
	}
	if r.ExistingWorkspace.AuthTokenCamel != "" {
		return r.ExistingWorkspace.AuthTokenCamel
	}
	return r.ExistingWorkspace.AuthToken
}

// VerifyResponse is returned by /verify after OTP confirmation.
type VerifyResponse struct {
	AuthToken     string        `json:"authToken"`
	WorkspaceID   string        `json:"workspaceId"`
	WorkspaceName string        `json:"workspaceName,omitempty"`
	URLs          WorkspaceURLs `json:"urls"`
	BuildInfo     any           `json:"buildInfo,omitempty"`
}

// StatusResponse is returned by GET /status/:workspaceId.
type StatusResponse struct {
	WorkspaceID            string   `json:"workspaceId"`
	Status                 string   `json:"status"`
	Progress               int      `json:"progress"`
	LandingPageReady       bool     `json:"landingPageReady"`
	EstimatedTimeRemaining string   `json:"estimatedTimeRemaining,omitempty"`
	CompletedSteps         []string `json:"completedSteps,omitempty"`
}

// ChatResponse is returned by POST /chat/:workspaceId.
type ChatResponse struct {
	WorkspaceID string `json:"workspaceId"`
	// ChatID is the server-side Otto conversation identifier.
	//
	// Observed 2026-07-08: repeated body-auth /chat calls for the same
	// workspace token returned the same chatId and preserved short-term
	// conversation context. Passing the returned chatId back also returned the
	// same chatId. Creating/selecting arbitrary chat IDs is not confirmed, so
	// callers should reuse only chat IDs returned by Audos.
	ChatID   string `json:"chatId"`
	Response string `json:"response"` // Otto's reply text
}

// ChatOption allows optional chat fields to be passed to Chat and
// ChatViaBearer.
type ChatOption func(*chatRequest)

type chatRequest struct {
	AuthToken string `json:"authToken,omitempty"`
	Message   string `json:"message"`
	ChatID    string `json:"chatId,omitempty"`
}

// WithChatID continues a chat using an ID returned by a previous ChatResponse.
// Do not invent IDs, and do not rely on an empty string to create a new topic
// thread until Audos documents that behavior.
func WithChatID(chatID string) ChatOption {
	return func(r *chatRequest) { r.ChatID = chatID }
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

func (c *OnboardClient) newRequest(method, path string, body any) (*http.Request, error) {
	var reqBody io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal: %w", err)
		}
		reqBody = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}
	return req, nil
}

func (c *OnboardClient) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}

	if out != nil {
		if err := json.Unmarshal(data, out); err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}
	}
	return nil
}

// ─── API methods ──────────────────────────────────────────────────────────────

// Start begins the onboarding flow. For returning users it returns an AuthToken
// directly (no OTP needed). For new users it triggers an OTP email and returns
// a SessionToken to pass to Verify.
func (c *OnboardClient) Start(email, businessIdea string, opts ...StartOption) (StartResponse, error) {
	req := &startRequest{Email: email, BusinessIdea: businessIdea}
	for _, o := range opts {
		o(req)
	}

	httpReq, err := c.newRequest(http.MethodPost, "/start", req)
	if err != nil {
		return StartResponse{}, err
	}

	var out StartResponse
	return out, c.do(httpReq, &out)
}

// Verify confirms the OTP code and returns the auth token + workspace details.
// After a successful Verify, call SetAuthToken(resp.AuthToken) to use this
// same client for subsequent authenticated calls.
func (c *OnboardClient) Verify(sessionToken, otpCode string) (VerifyResponse, error) {
	httpReq, err := c.newRequest(http.MethodPost, "/verify", map[string]string{
		"sessionToken": sessionToken,
		"otpCode":      otpCode,
	})
	if err != nil {
		return VerifyResponse{}, err
	}

	var out VerifyResponse
	return out, c.do(httpReq, &out)
}

// Status returns the current build status for the current auth token.
// workspaceID is kept for source compatibility with the earlier bearer route.
// The working production route authenticates from the body token.
func (c *OnboardClient) Status(workspaceID string) (StatusResponse, error) {
	if c.authToken == "" {
		return StatusResponse{}, fmt.Errorf("auth token required: call Start with WithWorkspaceID or Verify, then SetAuthToken")
	}
	httpReq, err := c.newRequest(http.MethodPost, "/status", map[string]string{
		"authToken": c.authToken,
	})
	if err != nil {
		return StatusResponse{}, err
	}

	var out StatusResponse
	return out, c.do(httpReq, &out)
}

// Chat sends a message to Otto (the Audos AI assistant) for the given workspace.
// workspaceID is kept for source compatibility with the earlier bearer route.
// The working production route authenticates from the body token.
func (c *OnboardClient) Chat(workspaceID, message string, opts ...ChatOption) (ChatResponse, error) {
	if c.authToken == "" {
		return ChatResponse{}, fmt.Errorf("auth token required: call Start with WithWorkspaceID or Verify, then SetAuthToken")
	}
	body := &chatRequest{AuthToken: c.authToken, Message: message}
	for _, o := range opts {
		o(body)
	}
	httpReq, err := c.newRequest(http.MethodPost, "/chat", body)
	if err != nil {
		return ChatResponse{}, err
	}

	var out ChatResponse
	if err := c.do(httpReq, &out); err != nil {
		return ChatResponse{}, err
	}
	if out.WorkspaceID != "" && out.WorkspaceID != workspaceID {
		return ChatResponse{}, fmt.Errorf("response workspace mismatch: expected %s, got %s", workspaceID, out.WorkspaceID)
	}
	return out, nil
}

// StatusViaBearer calls the documented bearer-header status route. Kept for
// Audos QA and backwards compatibility; prefer Status until bearer auth is
// confirmed in production.
func (c *OnboardClient) StatusViaBearer(workspaceID string) (StatusResponse, error) {
	httpReq, err := c.newRequest(http.MethodGet, "/status/"+workspaceID, nil)
	if err != nil {
		return StatusResponse{}, err
	}

	var out StatusResponse
	return out, c.do(httpReq, &out)
}

// ChatViaBearer calls the documented bearer-header chat route. Kept for Audos
// QA and backwards compatibility; prefer Chat until bearer auth is confirmed in
// production.
func (c *OnboardClient) ChatViaBearer(workspaceID, message string, opts ...ChatOption) (ChatResponse, error) {
	body := &chatRequest{Message: message}
	for _, o := range opts {
		o(body)
	}
	httpReq, err := c.newRequest(http.MethodPost, "/chat/"+workspaceID, body)
	if err != nil {
		return ChatResponse{}, err
	}

	var out ChatResponse
	return out, c.do(httpReq, &out)
}

// Rebuild triggers a workspace rebuild. Use when a previous build failed.
// Requires an auth token — call SetAuthToken or use NewOnboardClient.
func (c *OnboardClient) Rebuild(workspaceID string) error {
	httpReq, err := c.newRequest(http.MethodPost, "/rebuild/"+workspaceID, nil)
	if err != nil {
		return err
	}
	return c.do(httpReq, nil)
}
