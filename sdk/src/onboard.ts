/**
 * Audos Onboarding / Otto Chat API client.
 *
 * Uses Bearer token auth and the /api/agent/onboard base path — entirely
 * separate from the workspace hooks API in createClient().
 *
 * Known limitation: Tokens are workspace-scoped and only issued via the
 * Start/Verify flow. Workspaces provisioned directly through the Audos
 * platform (not through this API) will not have an email→token mapping
 * and will return AUTH_TOKEN_INVALID until linked by the Audos team.
 */

const ONBOARD_BASE_URL = 'https://audos.com/api/agent/onboard';

// ─── Types ────────────────────────────────────────────────────────────────────

export interface WorkspaceURLs {
  landingPage: string;
  workspace: string;
  dashboard?: string;
  appSpace?: string;
}

/** Covers both outcomes of /start:
 *  - New user: sessionToken is set; authToken/workspaceId/urls are absent.
 *  - Returning user: authToken + workspace fields are set; sessionToken is absent.
 */
export interface StartResponse {
  // Returning user fields
  auth_token?: string;
  workspaceId?: string;
  workspaceName?: string;
  urls?: WorkspaceURLs;
  // New user field — present only when OTP email was sent
  sessionToken?: string;
}

export interface VerifyResponse {
  authToken: string;
  workspaceId: string;
  workspaceName?: string;
  urls: WorkspaceURLs;
  buildInfo?: unknown;
}

export interface StatusResponse {
  workspaceId: string;
  status: string;
  progress: number;
  landingPageReady: boolean;
  estimatedTimeRemaining?: string;
  completedSteps?: string[];
}

export interface ChatResponse {
  workspaceId: string;
  chatId: string;
  /** Otto's reply text */
  response: string;
}

export interface StartOptions {
  businessName?: string;
  targetCustomer?: string;
  callbackUrl?: string;
  createNew?: boolean;
}

// ─── Factory ──────────────────────────────────────────────────────────────────

/**
 * Creates an Audos onboarding client.
 *
 * Pass an authToken when you already have one from a previous Start/Verify
 * flow. Omit it (or pass undefined) to begin the OTP flow unauthenticated —
 * in that case call `setAuthToken(token)` after a successful `verify()`.
 *
 * @example — new user
 * ```ts
 * const client = createOnboardClient();
 * const { sessionToken } = await client.start('user@example.com', 'An app that does X');
 * const { authToken, workspaceId } = await client.verify(sessionToken!, '1234');
 * client.setAuthToken(authToken);
 * const status = await client.status(workspaceId);
 * ```
 *
 * @example — returning user
 * ```ts
 * const client = createOnboardClient(storedToken);
 * const reply = await client.chat(workspaceId, 'Update my landing page headline');
 * ```
 */
export function createOnboardClient(authToken?: string) {
  let _token = authToken;

  function setAuthToken(token: string): void {
    _token = token;
  }

  function buildHeaders(withBody: boolean): Record<string, string> {
    const headers: Record<string, string> = {};
    if (withBody) headers['Content-Type'] = 'application/json';
    if (_token) headers['Authorization'] = `Bearer ${_token}`;
    return headers;
  }

  async function request<T>(
    method: string,
    path: string,
    body?: unknown,
  ): Promise<T> {
    const response = await fetch(`${ONBOARD_BASE_URL}${path}`, {
      method,
      headers: buildHeaders(body !== undefined),
      body: body !== undefined ? JSON.stringify(body) : undefined,
    });

    const text = await response.text();
    if (!response.ok) {
      throw new Error(`API error ${response.status}: ${text}`);
    }

    // Some endpoints (e.g. rebuild) may return an empty body on success.
    if (!text) return undefined as T;
    return JSON.parse(text) as T;
  }

  /**
   * Begins the onboarding flow.
   *
   * - New user → returns `sessionToken`; OTP is emailed.
   * - Returning user → returns `auth_token` + workspace fields; no OTP needed.
   */
  async function start(
    email: string,
    businessIdea: string,
    options: StartOptions = {},
  ): Promise<StartResponse> {
    return request<StartResponse>('POST', '/start', {
      email,
      businessIdea,
      ...options,
    });
  }

  /**
   * Confirms the OTP code and returns the auth token + workspace details.
   * After success, call `setAuthToken(result.authToken)` to authenticate
   * subsequent calls on this same client instance.
   */
  async function verify(
    sessionToken: string,
    otpCode: string,
  ): Promise<VerifyResponse> {
    return request<VerifyResponse>('POST', '/verify', { sessionToken, otpCode });
  }

  /**
   * Returns the current build status for the given workspace.
   * Requires an auth token.
   */
  async function status(workspaceId: string): Promise<StatusResponse> {
    return request<StatusResponse>('GET', `/status/${workspaceId}`);
  }

  /**
   * Sends a message to Otto (Audos AI assistant) for the given workspace.
   * Requires an auth token.
   */
  async function chat(workspaceId: string, message: string): Promise<ChatResponse> {
    return request<ChatResponse>('POST', `/chat/${workspaceId}`, { message });
  }

  /**
   * Triggers a workspace rebuild. Use when a previous build failed.
   * Requires an auth token.
   */
  async function rebuild(workspaceId: string): Promise<void> {
    await request<void>('POST', `/rebuild/${workspaceId}`);
  }

  return { setAuthToken, start, verify, status, chat, rebuild };
}

export type OnboardClient = ReturnType<typeof createOnboardClient>;
