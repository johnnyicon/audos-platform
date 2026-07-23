---
date: 2026-04-16
area: ai-api
status: fixed
filed: no
label: The AI hook's Generate() call type has a hard ~1000-token response cap that Chat() on the same underlying model doesn't
---

Audos's AI hook exposes (at least) two call shapes — `Generate()` and `Chat()` — routing under the hood
to the same underlying OpenAI Chat Completions proxy. They are not equivalent: `Generate()` has a hard
~1000-token response ceiling; `Chat()`, against the identical model, does not. Anything needing a
longer response (a multi-paragraph summary, an extracted document) silently truncates through
`Generate()` with no error, only a cut-off result.

Confirmed directly in code: a PDF-text-extraction path was deliberately built against `Chat()` instead
of `Generate()`, specifically to route around the cap — the code comment states plainly "Chat is
uncapped so the 1000-token ceiling that afflicts Generate() does not apply." Also confirmed via the same
code path: the underlying proxy only accepts PDF as a MIME type on the `file` content block — a DOCX
sent the same way is rejected outright, since it's OpenAI's own Chat Completions API enforcing that
restriction, not an Audos-specific one.

Status here is "fixed" in the sense that the team's own workaround (call `Chat()`, not `Generate()`, for
anything past ~1000 tokens) is confirmed working and now standard practice — not that Audos changed the
underlying cap.

Source: `throughline-daemon/internal/ai/pdf.go`, `throughline-forge/docs/handoff-2026-04-16-session-close.md`.
