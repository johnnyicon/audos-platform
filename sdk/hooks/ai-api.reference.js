// ─────────────────────────────────────────────────────────────────────────────
// Reference `ai-api` server function (hook) — corrected version.
//
// STATUS: UNVERIFIED ON DEPLOY. Grounded in the documented server-function runtime
// (`platform.integrations.proxy`, see docs/platform/reports/architecture-audit-2026-07-17-
// technical-companion.md I-1), and in the live audit of the existing hook (2026-08-10).
// It has NOT yet been deployed and run — deploying a server function needs
// `manage_server_functions`, i.e. a Cursor delegation job with an Otto token. Treat this
// as a ready-to-deploy template, not a proven artifact, until it's run in a workspace.
//
// WHAT IT FIXES vs the current hook (Throughline `workspace-351699/ai-api`):
//   1. GPT-5 + o-series: sends `max_completion_tokens` for models that require it, instead of
//      always sending `max_tokens` (bugs/0040). Unblocks gpt-5*, o1, o3, o3-mini, o4-mini.
//   2. Errors: returns the real upstream reason, not a generic "AI generation failed".
//   3. Vision: passes OpenAI-style content arrays through unchanged (already works upstream).
//
// WHAT IT CANNOT DO: reach non-OpenAI providers with Audos's credentials. The runtime proxy
// allowlist is openai/stripe/twilio/heygen — NO Anthropic, NO Google. So Claude/Fable/Gemini
// are not reachable this way. (Bring-your-own-key via platform.externalFetch is a separate,
// untested path and is NOT "leveraging Audos's providers".)
//
// RUNTIME GOTCHA baked in: platform.integrations.proxy's `body` MUST be a JSON string. Passing
// a plain object silently returns the platform's HTML index with status 200 (I-1). We stringify.
// ─────────────────────────────────────────────────────────────────────────────

const OPENAI_CHAT = '/v1/chat/completions';

// GPT-5 family and o-series reasoning models require max_completion_tokens.
function needsMaxCompletionTokens(model) {
  return /^(gpt-5|o[1-9])/.test(model || '');
}

module.exports = async function (req, res) {
  try {
    const body = req.body || {};
    const action = body.action;
    const model = body.model || 'gpt-4o-mini';

    // Build the OpenAI request from either action.
    let messages;
    if (action === 'chat') {
      if (!Array.isArray(body.messages) || body.messages.length === 0) {
        return res.status(400).json({ error: 'Missing required field: messages (array of {role, content})' });
      }
      messages = body.messages;
    } else if (action === 'generate') {
      if (typeof body.prompt !== 'string') {
        return res.status(400).json({ error: 'Missing required field: prompt' });
      }
      messages = [{ role: 'user', content: body.prompt }];
    } else {
      return res.status(400).json({ error: 'Invalid action. Use "generate" for single-turn or "chat" for multi-turn.' });
    }

    if (body.systemPrompt) {
      messages = [{ role: 'system', content: body.systemPrompt }, ...messages];
    }

    const payload = { model, messages };
    const limit = body.maxTokens != null ? body.maxTokens : 1200;
    // THE FIX: pick the right token parameter per model family.
    if (needsMaxCompletionTokens(model)) {
      payload.max_completion_tokens = limit;
    } else {
      payload.max_tokens = limit;
    }
    // Older models accept temperature; reasoning models reject non-default values, so only
    // forward it when the caller set it and the model isn't a reasoning model.
    if (body.temperature != null && !needsMaxCompletionTokens(model)) {
      payload.temperature = body.temperature;
    }

    const upstream = await platform.integrations.proxy('openai', OPENAI_CHAT, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload), // MUST be a string — see header note.
    });

    // proxy() returns the raw upstream response; parse defensively.
    let data;
    try {
      data = typeof upstream === 'string' ? JSON.parse(upstream) : upstream;
    } catch {
      return res.status(502).json({ error: 'Upstream returned a non-JSON response', raw: String(upstream).slice(0, 300) });
    }

    if (data && data.error) {
      // Surface the REAL reason (unlike the current hook, which flattens it).
      return res.status(data.error.code === 'model_not_found' ? 404 : 400)
        .json({ error: data.error.message || 'OpenAI error', code: data.error.code, upstream: data.error });
    }

    const choice = data && data.choices && data.choices[0];
    const text = choice && choice.message ? (choice.message.content || '') : '';
    return res.status(200).json({ success: true, text, model: data.model || model, usage: data.usage || null });
  } catch (err) {
    return res.status(500).json({ error: 'Hook error: ' + (err && err.message ? err.message : String(err)) });
  }
};
