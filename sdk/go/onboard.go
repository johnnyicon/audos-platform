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
// Auth uses Bearer token (format: aud_live_ + 48 hex chars), not x-api-key.
// The base path (/api/agent/onboard) is entirely separate from the workspace
// hooks API used by Client.
//
// Known limitation: Tokens are workspace-scoped and only issued via the
// Start/Verify flow. Workspaces provisioned directly through the Audos
// platform (not through this API) will not have an email→token mapping
// and will return AUTH_TOKEN_INVALID until linked by the Audos team.
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
	BusinessName   string `json:"businessName,omitempty"`
	TargetCustomer string `json:"targetCustomer,omitempty"`
	CallbackURL    string `json:"callbackUrl,omitempty"`
	CreateNew      bool   `json:"createNew,omitempty"`
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

// WithCreateNew forces creation of a new workspace even if one exists for the email.
func WithCreateNew() StartOption {
	return func(r *startRequest) { r.CreateNew = true }
}

// StartResponse covers both outcomes of /start:
//   - New user: SessionToken is set; AuthToken/WorkspaceID/WorkspaceName/URLs are empty.
//   - Returning user: AuthToken + workspace fields are set; SessionToken is empty.
type StartResponse struct {
	// Returning user fields
	AuthToken     string        `json:"auth_token,omitempty"`
	WorkspaceID   string        `json:"workspaceId,omitempty"`
	WorkspaceName string        `json:"workspaceName,omitempty"`
	URLs          WorkspaceURLs `json:"urls,omitempty"`

	// New user field — present only when OTP email was sent
	SessionToken string `json:"sessionToken,omitempty"`
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
	WorkspaceID           string   `json:"workspaceId"`
	Status                string   `json:"status"`
	Progress              int      `json:"progress"`
	LandingPageReady      bool     `json:"landingPageReady"`
	EstimatedTimeRemaining string  `json:"estimatedTimeRemaining,omitempty"`
	CompletedSteps        []string `json:"completedSteps,omitempty"`
}

// ChatResponse is returned by POST /chat/:workspaceId.
type ChatResponse struct {
	WorkspaceID string `json:"workspaceId"`
	ChatID      string `json:"chatId"`
	Response    string `json:"response"` // Otto's reply text
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

// Status returns the current build status for the given workspace.
// Requires an auth token — call SetAuthToken or use NewOnboardClient.
func (c *OnboardClient) Status(workspaceID string) (StatusResponse, error) {
	httpReq, err := c.newRequest(http.MethodGet, "/status/"+workspaceID, nil)
	if err != nil {
		return StatusResponse{}, err
	}

	var out StatusResponse
	return out, c.do(httpReq, &out)
}

// Chat sends a message to Otto (the Audos AI assistant) for the given workspace.
// Requires an auth token — call SetAuthToken or use NewOnboardClient.
func (c *OnboardClient) Chat(workspaceID, message string) (ChatResponse, error) {
	httpReq, err := c.newRequest(http.MethodPost, "/chat/"+workspaceID, map[string]string{
		"message": message,
	})
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
