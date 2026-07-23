/**
 * Audos Onboarding / Otto Chat API client.
 *
 * Uses the /api/agent/onboard base path — entirely separate from the
 * workspace hooks API in createClient().
 *
 * Production note, 2026-07-08: existing-workspace access currently works
 * through body-auth endpoints:
 *   - POST /start with workspaceId + createNew=false
 *   - POST /status with { authToken }
 *   - POST /chat with { authToken, message, chatId? }
 *
 * Bearer-header routes are retained as explicit *ViaBearer methods for
 * compatibility checks, but they may return 401 until Audos publishes the
 * bearer-header fix.
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
  authToken?: string;
  workspaceId?: string;
  workspaceName?: string;
  urls?: WorkspaceURLs;
  existingWorkspace?: ExistingWorkspace;
  // New user field — present only when OTP email was sent
  sessionToken?: string;
}

export interface ExistingWorkspace {
  authToken?: string;
  auth_token?: string;
  workspaceId?: string;
  workspaceName?: string;
  urls?: WorkspaceURLs;
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
  /**
   * Server-side Otto conversation identifier.
   *
   * Observed 2026-07-08: repeated body-auth /chat calls for the same
   * workspace token returned the same chatId and preserved short-term
   * conversation context. Passing the returned chatId back also returned the
   * same chatId. Creating/selecting arbitrary chat IDs is not confirmed, so
   * callers should reuse only chat IDs returned by Audos.
   */
  chatId: string;
  /** Otto's reply text */
  response: string;
}

export interface ChatOptions {
  /**
   * Optional continuation hint from a previous ChatResponse.
   *
   * Use only values returned by Audos. Do not invent UUIDs, and do not rely on
   * an empty string to create a new topic thread until Audos documents that
   * behavior.
   */
  chatId?: string;
}

export interface StartOptions {
  workspaceId?: string;
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
 * const client = createOnboardClient();
 * const start = await client.start('user@example.com', 'Existing workspace work', {
 *   workspaceId,
 *   createNew: false,
 * });
 * client.setAuthToken(client.getAuthTokenFromStart(start)!);
 * const reply = await client.chat(workspaceId, 'Update my landing page headline');
 * ```
 */
export function createOnboardClient(authToken?: string) {
  let _token = authToken;

  function setAuthToken(token: string): void {
    _token = token;
  }

  function getAuthTokenFromStart(result: StartResponse): string | undefined {
    return (
      result.authToken ||
      result.auth_token ||
      result.existingWorkspace?.authToken ||
      result.existingWorkspace?.auth_token
    );
  }

  function requireAuthToken(): string {
    if (!_token) {
      throw new Error(
        '[audos-sdk] auth token required. Call start(..., { workspaceId, createNew: false }) or verify(), then setAuthToken().',
      );
    }
    return _token;
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
   * Returns the current build status for the current auth token.
   *
   * `workspaceId` is kept for source compatibility with the earlier bearer
   * route. The working production route authenticates from the body token.
   */
  async function status(_workspaceId?: string): Promise<StatusResponse> {
    return request<StatusResponse>('POST', '/status', {
      authToken: requireAuthToken(),
    });
  }

  /**
   * Sends a message to Otto (Audos AI assistant) for the given workspace.
   *
   * `workspaceId` is kept for source compatibility with the earlier bearer
   * route. The working production route authenticates from the body token.
   */
  async function chat(
    workspaceId: string,
    message: string,
    options: ChatOptions = {},
  ): Promise<ChatResponse> {
    const body: Record<string, string> = {
      authToken: requireAuthToken(),
      message,
    };
    if (options.chatId) body.chatId = options.chatId;
    const response = await request<ChatResponse>('POST', '/chat', {
      ...body,
    });
    if (response.workspaceId && response.workspaceId !== workspaceId) {
      throw new Error(
        `[audos-sdk] response workspace mismatch: expected ${workspaceId}, got ${response.workspaceId}`,
      );
    }
    return response;
  }

  /**
   * Calls the documented bearer-header status route. Kept for Audos QA and
   * backwards compatibility; prefer `status()` until bearer auth is confirmed.
   */
  async function statusViaBearer(workspaceId: string): Promise<StatusResponse> {
    return request<StatusResponse>('GET', `/status/${workspaceId}`);
  }

  /**
   * Calls the documented bearer-header chat route. Kept for Audos QA and
   * backwards compatibility; prefer `chat()` until bearer auth is confirmed.
   */
  async function chatViaBearer(
    workspaceId: string,
    message: string,
    options: ChatOptions = {},
  ): Promise<ChatResponse> {
    const body: Record<string, string> = { message };
    if (options.chatId) body.chatId = options.chatId;
    return request<ChatResponse>('POST', `/chat/${workspaceId}`, body);
  }

  /**
   * Triggers a workspace rebuild. Use when a previous build failed.
   * Requires an auth token.
   */
  async function rebuild(workspaceId: string): Promise<void> {
    await request<void>('POST', `/rebuild/${workspaceId}`);
  }

  return {
    setAuthToken,
    getAuthTokenFromStart,
    start,
    verify,
    status,
    chat,
    statusViaBearer,
    chatViaBearer,
    rebuild,
  };
}

export type OnboardClient = ReturnType<typeof createOnboardClient>;
