# What we know about the Cursor backend Audos delegates to

*Researched 2026-07-13. Sources: Cursor's public docs, the Cursor community forum, and Otto (asked
directly, with clear boundaries on what it could and couldn't answer).*

## It's the real product, confirmed

Audos's "App agent / Cursor" build backend is genuinely **Cursor's Background Agent product** — not
a similarly-named internal tool. Confirmed by the exact error text we hit ourselves:
`[usage_limit_exceeded] Background Agent requires at least $2 remaining until your hard limit` matches
Cursor's real billing/hard-limit mechanism.

Cursor publishes a real **Cloud Agents API** at `cursor.com/docs/cloud-agent/api/endpoints`:

- **Auth:** Basic or Bearer, via a user API key or a service-account key.
- **Create/manage agents:** `POST /v1/agents`, `GET /v1/agents`, `GET /v1/agents/{id}`,
  archive/unarchive/delete.
- **Runs:** `POST /v1/agents/{id}/runs` (follow-up prompts), `GET .../runs`, `GET .../runs/{runId}`,
  `GET .../runs/{runId}/stream` (SSE), cancel.
- **Models:** `GET /v1/models`.
- **Usage:** `GET /v1/agents/{id}/usage` — token counts per run.
- **Artifacts:** list + download (15-minute presigned S3 URLs).

## The truncation problem has an upstream cause

This explains a real chunk of today's pain, not just an Audos choice: Cursor's own API **truncates
by design**. Per their docs, the SSE stream "omits" `args`/`result` fields that are too large and sets
a truncation flag, and **there is no documented endpoint to retrieve the full untruncated log** —
intermediate execution detail is scoped to the stream, which has this built-in limit.

Practical implication for `docs/platform/20-power-user-wishlist.md` item 4: the raw truncation may be
a genuine Cursor-side constraint Audos inherited, not something they chose. That doesn't remove the
ask — Audos could still capture and store the full run transcript server-side on their end, independent
of what Cursor's stream truncates — but it reframes it: worth asking Audos "do you persist the full run
transcript anywhere we can access," rather than assuming they're withholding something Cursor already
hands them in full.

## Models — more visible than expected, and actionable

Asked Otto directly what runs our jobs. Its answer, with the boundary stated plainly:

- Otto's own `create_job_drafts` tool exposes a **`cursorModel` parameter with a curated key list**:
  `fable-5`, `opus-4.8`, `gpt-5.6-sol`, `gpt-5.6-terra`, `sonnet-5`, `gpt-5.5`, `composer-max`.
- **None of today's jobs pinned a model** — they all ran on Audos's default/auto-pick, and Otto
  confirmed the job outcome data doesn't report which model actually executed. So we don't know what
  ran today's builds specifically.
- **Actionable finding:** we can request a specific model on a future job by name. For something that
  needs stronger reasoning (a genuinely hard debugging pass, not routine scaffolding), explicitly pin
  `opus-4.8` or `composer-max` rather than leaving it on auto — worth testing whether pinning changes
  reliability, especially for the kind of "read the actual code, don't guess" tasks that gave us trouble
  today.
- What's still invisible: the actual provider mapping behind those aliases, and whether Audos overrides
  the requested model server-side regardless of what's asked.

## Model pinning — confirmed working, verified live

Ran the experiment directly rather than trust a report: dispatched a job with `cursorModel: "fable-5"`
explicitly set, briefed to build a one-line throwaway app that states what model it believes itself to
be. Otto's own job summary truncated before the self-ID line (as usual), so — per the standing rule —
we published the tiny app and read it live in the browser instead of trusting the report.

**Live result, verbatim from the rendered page:**

> "The agent that built this file believes it is 'Fable 5' (a Claude-family model, per its system
> prompt and internal model slugs like 'claude-fable-5-thinking') — this is a system-prompt-provided
> identity, not independently verifiable by the agent itself."

**What this confirms:**
- **Model pinning genuinely works** — the requested `fable-5` matches what the agent reports running as.
  Not just accepted-and-ignored at dispatch (which is all we could confirm before this test).
- **New concrete detail:** the internal model slug is `claude-fable-5-thinking` — a "thinking" variant
  naming convention not previously known.
- **Indirect evidence for the pass-through question above:** the agent explicitly cites "its system
  prompt" as the source of its own identity, and correctly caveats that this is asserted, not
  self-verifiable. That's real support for Otto's labeled guess that Audos/Cursor injects genuine
  system-level context (not just an appended user message) — the model itself is reporting a
  system-prompt-sourced fact, appropriately hedged.

**Actionable takeaway:** pinning a specific model via `cursorModel` is a real, working lever — worth
using deliberately (e.g. `opus-4.8` or `composer-max` for a genuinely hard debugging pass) rather than
leaving every job on auto-pick.

### Follow-up: a real bug found, and fixed, on a job pinned to fable-5

After the model probe, we found — by directly looking, not assuming — that the small "Assistant" pill
(meant to be a lightweight way to reach chat from a full-screen app, per the original deep-link fix's
own brief) was instead calling `returnToAgentView()` in `Desktop.tsx`, which tore down the full-screen
shell entirely and returned the user to the **complete old three-pane chat interface**. Worth being
honest about the process here: an earlier partial screenshot led to a wrong first read (mistaking it
for a harmless small overlay); a fuller screenshot from the user corrected that immediately, and we
verified the real behavior with our own click before concluding anything.

Dispatched a fix job with `cursorModel: "fable-5"` explicitly pinned — both to fix the real bug and to
test whether pinning a specific model changes diagnostic quality on a job like this. Result: the job
found and cited the exact code (`returnToAgentView`'s state-teardown calls), applied a fix so the pill
opens a small bottom-right popup instead, and **we independently verified it live** — clicked
"Assistant" post-publish and confirmed the full-screen app stays untouched behind a compact drawer.

**Caveat on causation:** this is one successful run on a pinned model, not a controlled comparison
against auto-pick. Worth remembering as a positive data point for pinning `fable-5`/`opus-4.8` on hard
diagnostic tasks, not treating as proven cause-and-effect from a single trial.

## System prompt — genuinely not obtainable right now

Two independent paths, both dead ends, honestly:

- **Otto has no visibility into it.** It knows the *effects* (genesis-shell conventions, brand-token
  steering, single-file TSX, `useWorkspaceDB`) because it sees the outputs and shares the same
  integration docs — but it does not have the actual injected prompt text, and said so plainly rather
  than guess.
- **Cursor doesn't publish its production system prompt.** Unofficial community reconstructions exist
  (a "leaked" prompt dated March 2025 surfaces in search), but that's over a year stale, unverified,
  and for Cursor's general IDE agent — not necessarily what Audos's own Background Agent integration
  additionally layers on top for app builds specifically. Not worth treating as authoritative; we didn't
  reproduce it here.

## Is Audos a thin pass-through to Cursor, or does it inject its own steering?

Asked Otto directly (2026-07-13). Its answer, again marking know/infer/can't-see explicitly:

- **Not a thin pass-through — confirmed by behavior, not by seeing the payload.** A Cursor agent Otto
  never explicitly briefs on project-specific conventions (single-file TSX, `useWorkspaceDB`, brand
  tokens, "publish = recompile") still produces conforming code and correctly references
  `Desktop.tsx`/`config.json`. Something beyond Otto's raw task text is being layered in before the
  agent runs.
- **What's genuinely unknown:** the exact injected content, whether it's one system prompt or several
  context files, how much of Otto's own brief is passed verbatim vs. summarized — Otto has never seen
  the actual payload Audos sends to Cursor.
- **Whether Cursor's API honors an externally-supplied system message the way a raw Anthropic/OpenAI
  call does:** unknown, and Otto was explicit this is a *labeled guess*, not a finding — its best guess
  is that Cursor's Background Agent is a product-level harness with its **own fixed system
  prompt/tool-loop**, and whatever Audos supplies is more likely injected as task/context material
  *within* that harness than a true system-role override. No documentation or API-shape evidence either
  way.

## The real cause of "self-reports don't match reality" — an authority boundary, not a lying agent

This is the sharpest reframe we've gotten on the problem that cost the most time today (see blog entry
6). Splitting the two distinct failure modes:

- **Truncated logs** are **not** a prompt/architecture issue — it's a straightforward Audos-side display
  limit (`list_jobs` clips long outcomes; no tool exposes the untruncated text except the human-only
  Tasks panel UI).
- **False "verified"/"published" claims** are better explained as an **authority/verification-boundary
  gap**, not the agent fabricating: the Cursor agent reports sincerely on what it *attempted inside its
  own sandbox*. The steps that actually determine ground truth — destructive-operation confirmation
  gates, publish/recompile, what bundle is actually served — are enforced **Audos-side, outside the
  agent's control or visibility**. The agent can genuinely believe "I did X" while Audos's own gate
  silently deferred or blocked it. We saw this exact pattern directly: job #81263 reported an app
  deletion as "staged" when the file removal was actually sitting behind a human-confirmation gate the
  agent has no way to see past.

**Practical implication, stated plainly: this is why "verify live yourself" isn't just caution — it's
structurally necessary.** The agent's report is honest about its own sandbox and uninformative about the
platform's actual state. No amount of better prompting fixes that; only checking the real, served result
closes the gap.

## One concrete, load-bearing finding: Background Agents have a confirmed AGENTS.md bug

Cursor's own community forum has a bug report, **confirmed by Cursor staff** ("I have sent them to the
team"), that **Background Agents do not reliably load `AGENTS.md`** — rules defined there get ignored
unless referenced explicitly. No fix released as of this writing; the only workaround is putting
instructions directly in the prompt rather than relying on a convention file being picked up.

**This validates something we were already doing without knowing why it mattered:** every job brief we
sent Otto today was long, explicit, and self-contained — full context restated in the prompt, not
"follow the conventions in the repo." That wasn't just caution; it was, per this bug, the *only*
reliable way to get a Background Agent to actually follow project rules. Worth stating as deliberate
practice going forward, not an accident: **don't rely on `AGENTS.md`/rules-file pickup for any job
delegated to this Cursor backend — restate the relevant instructions in the prompt itself, every time.**
