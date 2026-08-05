# `get_ad_campaigns` — list campaigns

**Status: ✅ verified live.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`). Free, read-only.

Returned "No ad campaigns found" for the DoKnow workspace — correct, since this workspace has never
launched a campaign. Worth confirming explicitly: a clean, correctly-typed empty result, not a silent
error or a malformed response — an API that errors instead of returning an empty set is a worse failure
mode, and this one didn't.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0031-ads-and-marketing-live-verification.md`.
