/**
 * Audos AI client — standalone, browser-safe access to a workspace's `ai-api` hook.
 *
 * Separate from createClient() in index.ts on purpose. createClient is server-only
 * (it carries the db-api x-api-key and throws in a browser). The AI hook is different:
 * it needs no key and its CORS is open (verified 2026-08-10 — the endpoint reflects the
 * request Origin and answers the preflight), so it is safe to call directly from a
 * browser app. This module therefore has no window guard.
 *
 * ── What this talks to ──────────────────────────────────────────────────────────
 * `POST {baseUrl}/api/hooks/execute/workspace-{workspaceId}/{hookName}`
 *
 * IMPORTANT — `ai-api` is a *workspace-built server function*, not a platform feature.
 * Verified 2026-08-10: it exists in the Throughline workspace but returns
 * `Hook 'ai-api' not found` in other workspaces. So this client only works against a
 * workspace that actually has such a hook. If you point it at a fresh workspace, expect
 * the 404 — which this client turns into an explicit, actionable error.
 *
 * ── Provider / model reality (verified 2026-08-10) ──────────────────────────────
 * The hook proxies **OpenAI only**. Working: gpt-4o, gpt-4o-mini (default), gpt-4.1.
 * NOT available: gpt-5 / o-series (reasoning), Claude, Gemini. Reasoning/effort params
 * (reasoningEffort, effort, temperature, tools, response_format) are silently ignored.
 * Video/voice live on a different surface (Otto tools), not this hook.
 */

// ─── Types ──────────────────────────────────────────────────────────────────────

export interface AudosAIConfig {
  /** Numeric workspace id, e.g. "351699" (with or without the "workspace-" prefix). */
  workspaceId: string;
  /** Defaults to "https://audos.com". */
  baseUrl?: string;
  /** The hook name. Defaults to "ai-api". */
  hookName?: string;
  /** Default model for calls that don't pass one. Defaults to "gpt-4o-mini". */
  defaultModel?: string;
  /**
   * Optional value sent as the `x-hook-secret` header. NOTE: the reference `ai-api`
   * hook does NOT enforce this (verified 2026-08-10 — a bogus value still completed).
   * It is pass-through only, useful solely if *your own* hook implements the check.
   * Do not treat its presence as authentication.
   */
  hookSecret?: string;
}

/**
 * The model landscape behind the reference `ai-api` hook, from an explicit 43-name live audit
 * (2026-08-10). Three tiers:
 *
 * - `working`: usable today through the reference hook (OpenAI GPT-4/3.5 family; vision on the 4o/4.1
 *   families).
 * - `reachableButHookBlocked`: authorized on the proxy but unreachable through the *reference* hook,
 *   which hardcodes `max_tokens` while these models require `max_completion_tokens`. A hook that sends
 *   the right parameter should unlock them (unverified end-to-end — needs that hook).
 * - `unavailable`: genuine upstream `model_not_found`. Not OpenAI models (Claude/Fable/Gemini) or not
 *   provisioned for this account. No hook change reaches these through an OpenAI proxy.
 */
export const AUDOS_AI_MODELS = {
  working: [
    'gpt-4o', 'gpt-4o-mini', 'gpt-4o-2024-08-06',
    'gpt-4.1', 'gpt-4.1-mini', 'gpt-4.1-nano',
    'gpt-4-turbo', 'gpt-4', 'gpt-3.5-turbo',
  ],
  vision: ['gpt-4o', 'gpt-4o-mini', 'gpt-4.1', 'gpt-4.1-mini', 'gpt-4.1-nano'],
  reachableButHookBlocked: ['gpt-5', 'gpt-5-mini', 'gpt-5-nano', 'gpt-5.5', 'gpt-5.6', 'o1', 'o3', 'o3-mini', 'o4-mini'],
  unavailable: ['claude-*', 'fable-5', 'gemini-*', 'o1-mini', 'gpt-5.6-sol'],
} as const;

/** Models that require `max_completion_tokens` instead of `max_tokens` (GPT-5 + o-series). */
export function needsMaxCompletionTokens(model: string): boolean {
  return /^(gpt-5|o[1-9])/.test(model);
}

/**
 * Reasoning models spend part of the token budget on *internal* reasoning before emitting any visible
 * text, and that spend counts against max_completion_tokens. Measured 2026-08-14: gpt-5 on a moderately
 * complex prompt burned all 2000 tokens on reasoning and returned empty text; at 4000 it used ~1150
 * reasoning + real output. So the floor here is deliberately generous — too low silently yields nothing.
 */
export const REASONING_MIN_MAX_TOKENS = 4000;
const DEFAULT_MAX_TOKENS = 1200;

/** A content part for multimodal messages. Text, or an image by URL or data: URI. */
export type ContentPart =
  | { type: 'text'; text: string }
  | { type: 'image_url'; image_url: { url: string } };

export interface ChatMessage {
  role: 'user' | 'assistant';
  /**
   * A plain string, or an array of parts for vision. Image understanding is verified working
   * on `gpt-4o`/`gpt-4.1` through this hook (2026-08-10) — pass an `image_url` part with a
   * public URL or a `data:image/...;base64,...` URI.
   */
  content: string | ContentPart[];
}

/** Convenience: build a user message that pairs a prompt with one or more images. */
export function imageMessage(text: string, ...imageUrls: string[]): ChatMessage {
  return {
    role: 'user',
    content: [
      { type: 'text', text },
      ...imageUrls.map((url) => ({ type: 'image_url' as const, image_url: { url } })),
    ],
  };
}

export interface AIUsage {
  promptTokens: number;
  completionTokens: number;
  totalTokens: number;
}

export interface AIResult {
  text: string;
  model: string;
  usage: AIUsage | null;
}

export interface ChatOptions {
  model?: string;
  systemPrompt?: string;
  /** Honored by chat (unlike generate). Defaults to 1200. */
  maxTokens?: number;
  temperature?: number;
}

/** Thrown for any non-success outcome, carrying the HTTP status and server message. */
export class AudosAIError extends Error {
  status: number;
  serverError?: string;
  constructor(message: string, status: number, serverError?: string) {
    super(message);
    this.name = 'AudosAIError';
    this.status = status;
    this.serverError = serverError;
  }
}

// ─── Implementation ───────────────────────────────────────────────────────────────

/** OpenAI returns either camelCase (generate path) or snake_case (chat path). Normalize. */
function normalizeUsage(u: unknown): AIUsage | null {
  if (!u || typeof u !== 'object') return null;
  const o = u as Record<string, number>;
  const prompt = o.promptTokens ?? o.prompt_tokens;
  const completion = o.completionTokens ?? o.completion_tokens;
  const total = o.totalTokens ?? o.total_tokens;
  if (prompt == null && completion == null && total == null) return null;
  return {
    promptTokens: prompt ?? 0,
    completionTokens: completion ?? 0,
    totalTokens: total ?? (prompt ?? 0) + (completion ?? 0),
  };
}

/** Pull the real error out of `_meta.logs` (the hook hides it there behind a generic message). */
function extractLogDetail(parsed: any): string {
  const logs = parsed?._meta?.logs;
  if (!Array.isArray(logs)) return '';
  const line = logs.find((l: unknown) => typeof l === 'string' && /error|does not exist|not found/i.test(l));
  return typeof line === 'string' ? line.replace(/^\[ERROR]\s*/, '') : '';
}

export function createAudosAI(config: AudosAIConfig) {
  const baseUrl = (config.baseUrl ?? 'https://audos.com').replace(/\/$/, '');
  const hookName = config.hookName ?? 'ai-api';
  const defaultModel = config.defaultModel ?? 'gpt-4o-mini';
  const wsId = String(config.workspaceId).replace(/^workspace-/, '');
  const url = `${baseUrl}/api/hooks/execute/workspace-${wsId}/${hookName}`;

  async function post(body: unknown): Promise<any> {
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (config.hookSecret) headers['x-hook-secret'] = config.hookSecret;

    const res = await fetch(url, { method: 'POST', headers, body: JSON.stringify(body) });

    // Read the body once, as text, then parse — a failing response still carries a
    // useful { error } payload we must not discard.
    const raw = await res.text();
    let parsed: any = null;
    if (raw) {
      try { parsed = JSON.parse(raw); } catch { /* non-JSON body handled below */ }
    }

    // The hook flattens its top-level `error` to a generic "AI generation failed", but the
    // real upstream reason (e.g. OpenAI's "model … does not exist or you do not have access")
    // survives in `_meta.logs`. Dig it out so callers get something actionable.
    const detail = extractLogDetail(parsed) || parsed?.error || raw || '';

    if (!res.ok) {
      // Reasoning models can spend the entire budget on internal reasoning and return no visible
      // text. Give that its own actionable message instead of a bare "no text".
      if (parsed?.finishReason === 'length' && !parsed?.text) {
        const used = parsed?.usage?.completion_tokens ?? parsed?.usage?.completionTokens;
        throw new AudosAIError(
          `Model produced no visible text: the whole token budget${used ? ` (${used})` : ''} went to ` +
          `internal reasoning (finish_reason: length). Raise maxTokens — reasoning models need ` +
          `at least ~${REASONING_MIN_MAX_TOKENS}, and more for complex prompts.`,
          res.status, String(detail),
        );
      }
      // The signature failure for this client: the workspace has no such hook.
      if (res.status === 404 && /not found/i.test(String(parsed?.error ?? raw ?? ''))) {
        throw new AudosAIError(
          `Hook '${hookName}' not found in workspace-${wsId}. The ai-api hook is a ` +
          `workspace-built server function, not a platform feature — this workspace ` +
          `doesn't have one. Either point at a workspace that does, or build the hook there.`,
          res.status, String(detail),
        );
      }
      throw new AudosAIError(
        `Audos AI call failed (HTTP ${res.status}): ${detail || res.statusText}`,
        res.status, String(detail),
      );
    }

    if (parsed?.error) {
      throw new AudosAIError(`Audos AI error: ${detail}`, res.status, String(detail));
    }
    return parsed;
  }

  return {
    /**
     * Multi-turn chat — the reliable workhorse. Honors maxTokens; use this for anything
     * that may run long. `messages` must end with a user turn.
     */
    async chat(messages: ChatMessage[], opts: ChatOptions = {}): Promise<AIResult> {
      if (!messages.length || messages[messages.length - 1].role !== 'user') {
        throw new AudosAIError('chat(): messages must be non-empty and end with a user turn', 0);
      }
      const model = opts.model ?? defaultModel;
      const reasoning = needsMaxCompletionTokens(model);
      // Reasoning models need headroom for internal reasoning tokens; a 1200 default silently
      // returns empty text on anything non-trivial.
      const maxTokens = opts.maxTokens ?? (reasoning ? REASONING_MIN_MAX_TOKENS : DEFAULT_MAX_TOKENS);
      const d = await post({
        action: 'chat',
        messages,
        model,
        systemPrompt: opts.systemPrompt,
        maxTokens,
        temperature: opts.temperature,
      });
      return { text: d.text ?? '', model: d.model ?? '', usage: normalizeUsage(d.usage) };
    },

    /**
     * One-shot completion — the method to reach for by default. Wraps a single prompt as
     * a chat turn, so it is NOT subject to the generate action's 1000-token cap.
     */
    async complete(prompt: string, opts: ChatOptions = {}): Promise<AIResult> {
      return this.chat([{ role: 'user', content: prompt }], opts);
    },

    /**
     * Single-turn `generate` action. Kept for parity, but note two hard limits verified
     * on the reference hook: (1) output is capped at 1000 completion tokens regardless of
     * maxTokens; (2) on model failure it returns HTTP 200 with empty text and no error —
     * which this method converts into a thrown AudosAIError. Prefer complete()/chat().
     */
    async generate(prompt: string, opts: { model?: string; systemPrompt?: string } = {}): Promise<AIResult> {
      const d = await post({
        action: 'generate',
        prompt,
        model: opts.model ?? defaultModel,
        systemPrompt: opts.systemPrompt,
      });
      const text = d.text ?? '';
      if (!text) {
        throw new AudosAIError(
          `generate() returned empty text for model '${opts.model ?? defaultModel}'. The ` +
          `generate action reports success even when the model call fails — treat this as ` +
          `a failed generation (often an unavailable model). Use chat()/complete() and a ` +
          `working model (gpt-4o, gpt-4o-mini, gpt-4.1).`,
          200,
        );
      }
      return { text, model: d.model ?? '', usage: normalizeUsage(d.usage) };
    },

    /** The resolved endpoint, for logging/debugging. */
    endpoint: url,
  };
}
