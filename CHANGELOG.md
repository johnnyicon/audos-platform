# Changelog — Audos SDK

A dated, versioned record of what actually **changed in this SDK's own guidance** — corrected docs, new
capability findings, new backlog items, new playbooks. **Not** blog hosting, domain moves, or the site's
own build/UI — see `blog/HOSTING-CHANGELOG.md` for that; it isn't Audos-platform knowledge, so it doesn't
belong here even though it's a real, dated change. Distinct from the other three records, on purpose:

- **`ACTIVITY-LOG.md`** — raw, unharvested actions (what we just did), trimmed/archived once harvested.
- **`blog/`** — narrative write-ups, for a human reading the story of how something was found.
- **`CHANGELOG.md`** (this file) — the terse "what's different in the SDK now" summary, one entry per
  harvest pass, written for someone who wants to know what changed without reading the narrative.
- **`blog/HOSTING-CHANGELOG.md`** — the same idea, scoped to the blog's own hosting/build/UI instead of
  the SDK's guidance content.

Newest first.

---

## 2026-08-05 — pgvector re-verified live, plus two findings about the platform's own tooling

Before repeating the 16 July pgvector finding to Audos directly, re-checked it — a three-week-old claim
stated confidently is worse than no claim. **It holds:** `pg_extension` returns exactly one row,
`plpgsql`. `vector` is not installed. `CAPABILITY-MATRIX.md` updated with the re-verification date and
`bugs/0015` annotated.

Also corrected a proportion problem in `bugs/0015`'s original wording: it called this "the single biggest
blocker for any retrieval-grounded generation pipeline," which `BACKLOG.md #13` later disproved by
benchmarking the JSON-array fallback at real embedding size. It's a ceiling on headroom, not a blocker
today. Overstating a real gap is its own kind of inaccuracy.

Getting the answer took four routes and surfaced two new findings, both filed to `BACKLOG.md`:

- **#40** — the Otto reachable over the *external* onboarding API has no `manage_server_functions` tool,
  contradicting `docs/platform/07`, which documents it as the creation path "if you're building tools or
  automating." Another instance of the tool-surface split in `docs/platform/29`. Doc 07 corrected in place.
  New: `blog/bugs/0038`.
- **#41** — the jobs list truncates a job's output **mid-JSON**, with no marker, and deleting a temporary
  hook destroys the `get_hook_logs` fallback. Same shape as `#4`. New: `blog/bugs/0039`.

Narrative: `blog/experiments/0033-pgvector-still-not-enabled-reverified.md` — worth reading for the
meta-point, that answering "is this extension installed" required dispatching a background build job to a
third-party agent runner.

## 2026-07-23 (even later) — capability matrix expanded with 3 new areas, live-verified by asking Otto directly

Asked Otto (via the onboarding API's `/chat`) what Audos capabilities the SDK hadn't captured yet, then
independently verified every claim it made rather than recording the description as fact — new standing
pattern for anything Otto self-reports. Free/no-side-effect tools were run for real; anything that would
spend real ad money or reach a real person (launching a campaign, sending DMs) was deliberately left
unrun; media generation (real cost, no human impact) was run with explicit approval.

Two new capability-matrix categories (**Analytics & Reporting**, **Ads & Marketing**) and a **Media
Generation** section added, each backed by a new architecture doc
(`docs/platform/29-otto-tool-surface-vs-app-callable-hooks.md` — these are Otto-chat-orchestrated tools,
not `platform.*` hooks an app can call) plus three index docs (`docs/platform/30`–`32`) and one focused
reference file per real capability under a new `docs/platform/capabilities/{analytics,ads,media}/`
folder — the same progressive-disclosure shape as `otto-pilot`'s own `references/`.

Two new bugs filed to `BACKLOG.md` (#38, #39): `query_data_source(funnel-events)` returns a degenerate,
all-blank result where `query_events` correctly reports real data for the same workspace/window; and
`list_voiceover_voices` returns `undefined` IDs for ~82 of ~85 listed voices (only 3 hardcoded fallbacks
are actually selectable). Also refined the existing AI-model-lock row: the `gpt-4o-mini` lock applies only
to the in-app AI proxy, not to build agents (Cursor) or Audos Code, which use other models. Video and
voiceover generation were independently verified by downloading the returned files and probing them
(`ffprobe`) rather than trusting the URL — both real. Narrative: `blog/0019-asking-otto-what-else-audos-can-do.md`.

## 2026-07-23 (later still) — new durable doc 28, and the onboarding skill restructured as a portable package: otto-pilot

New doc: `docs/platform/28-otto-onboarding-api-auth-and-chatid.md`, consolidating `bugs/0029` and
`feature-requests/0018` (Bearer-vs-body-auth, `chatId` correlation behavior) into one reference, since
both are about the same onboarding API surface. The onboarding agent skill itself — corrected and
extended from Audos's own original file — moved out of `docs/audos-api/` and into a proper, portable
skill package: `sdk/skills/otto-pilot/` (`SKILL.md`, `original.md` for comparison, and a `references/`
folder the skill's own progressive-disclosure table points to, self-contained rather than linking back
into this repo). Renamed from a working title ("audos-onboarding") once it became clear the skill covers
the whole lifecycle now, not just the first ten minutes.

---

## 2026-07-23 (later) — mined the Throughline repos: 9 bugs, 10 feature requests, 6 experiments, 2 posts

A research pass across all four real Throughline repos (`throughline`, `throughline-daemon`,
`throughline-forge`, `throughline-tala`) — prioritizing `AUDOS.md`, `docs/audos-services.md`, ADRs,
decision journals, and a stash of dated friction reports in `throughline-forge/tmp/` that never made it
into this SDK — surfaced 28 candidate findings never previously captured. Wrote up 9 bugs (`0029`-`0037`),
10 feature requests (`0011`-`0020`), 6 experiments (`0023`-`0028`), and 2 narrative posts (`0017`-`0018`:
the Railway/Neon migration story, the Otto API delegation experiment). One candidate (new-app
registration compile timing) was written up as *inconclusive*, not a confirmed bug, since the original
source explicitly flagged it as unresolved. One (external backend hosting) was dropped entirely — an
unanswered question with no actual finding behind it. Checked two items against existing posts `0001`/
`0002` for overlap/contradiction before writing: no contradiction found, `blog/experiments/0024` is
careful to distinguish its finding (a *different* global, used for a *different* purpose) from `0002`'s
`__spaceContext?.username` claim, which still stands. Fixed a stale label in `blog/INDEX.md` for post
`0016` while in there. `BACKLOG.md` updated with rows `29`-`37`.

---

## 2026-07-22 — Audos claimed three DB bugs fixed; one independently confirmed false

Checked our filed bugs for movement. Audos's own Priority Support panel showed three (`bugs/0019`
orphaned table, `bugs/0020` serial-id/uuid rejection, `bugs/0023` DB credentials one-shot) marked
**Completed**, each claiming "verified and published to production" on 2026-07-17 evening. Given this
platform's own recent history of false "Complete" reports (`bugs/0027`, `0028`), tested live before
believing it. The credentials fix (`bugs/0023`) is **confirmed false** — direct network check shows
`409`/`401` unchanged, no regenerate control added. The other two are mid-verification via a scoped
create-then-drop probe. New blog entry: `blog/0016`.

---

## 2026-07-17 (latest) — new durable doc on the mandatory EmailGate mechanism, blog entry, docs pointer

Promoted the unified-space EmailGate finding (`bugs/0024`, `0027`, `0028`) from scattered bug write-ups
into durable platform knowledge: new `docs/platform/26-unified-space-signed-out-view-is-mandatory.md`
covers the mechanism (every unified-space workspace's signed-out view is unconditionally
`EmailGate.tsx`; custom auth is always a second, downstream gate), the CRM/ad-pixel side effect, and
what both rewrite attempts actually taught us. Added a pointer at the top of
`docs/platform/17-emailgate-otp-configuration.md` (which assumes EmailGate is an optional layer) and a
new lookup row in `skill/SKILL.md`. Wrote `blog/0015-no-live-path-to-our-own-content.md` narrating the
`bugs/0028` verification end to end.

---

## 2026-07-17 (still later) — checked the workaround path too: no live path currently reaches real content

With the invite-link fix confirmed broken (`bugs/0027`), checked whether the "old" email-then-password
workaround still worked as a fallback. It doesn't: a cold visitor now clears the Audos email gate with no
OTP, then lands on a generic, unbranded assistant shell reading an unrelated empty table — never reaching
the site's own password gate or the real Findings feed. Direct DB check confirmed the actual content (6
rows) is intact; this is an access-path problem, not data loss. Net: right now there is no live path by
which an outside visitor reaches field-notes' real content. New bug: `bugs/0028`, `BACKLOG.md #28`. Per
explicit instruction, the full 73-item content sync is **not** being run until this is fixed and
independently reverified.

---

## 2026-07-17 (later) — new bug: the EmailGate fix itself reproduced the self-report/reality gap

While shipping the fix for `bugs/0024` (replacing EmailGate with an invite-link system), the fix job
itself marked "Complete" with a specific, confident outcome — and the live site kept serving the old
code regardless. Confirmed with three independent live checks plus a direct read of the served source,
not taken on the job's word. New bug: `bugs/0027`, `BACKLOG.md #27`.

---

## 2026-07-17 — Architecture audit for a co-founder, Audos Code discovered, field-notes actually populated, EmailGate architecture exposed

**Red-team architecture audit.** A synthesis-only audit (no new tests, reads across the existing corpus)
answering whether Audos is architecturally built for hardcore agentic development. Went through multiple
adversarial review passes before being cleared for delivery to an Audos co-founder — caught and corrected
an overstated "vector search works" claim, hunted down unearned absolute language, and required two live
verification passes (identity defaults, Audos Code) before the staleness flags on several claims could be
dropped. Full report + technical companion: `docs/platform/reports/architecture-audit-2026-07-17.md`.

**Mid-audit discovery: Audos Code.** A second, actively-developed editing surface (`portfolio/code`, Beta
0.3.0, its own July 2026 changelog) exists and wasn't in the corpus at all. Live-tested it directly:
genuine token-by-token streaming output during execution (`experiments/0020`), running on a Claude Code
backend (not the Cursor pipeline documented elsewhere), with a routing-visible model picker. But the
core self-report/reality gap this whole project keeps finding reproduced live in this newer surface too
— a confidently-reported edit that never actually landed. Also used Audos Code's DB-write path to finally
populate field-notes' `content_items` table (`experiments/0021`), which direct-Postgres access (0-for-3,
`bugs/0023`) had been blocking since it was built.

**New finding: EmailGate wraps custom auth in every workspace, undisclosed CRM/pixel side effects.** A
real visitor to field-notes hit Audos's own email prompt instead of our custom password gate. Had Otto
read the served source directly: the platform's `EmailGate.tsx` is unconditionally the signed-out view
of any workspace, and submitting its email step registers a CRM contact and fires Meta/Reddit ad
pixels — a private internal tool was silently feeding visitor data into Audos's own lead-capture and
third-party ad platforms (`bugs/0024`, `experiments/0022`). Fix in progress: replacing `EmailGate.tsx`
with a closed-alpha invite-link auto-sign-in system.

New this pass: 3 bugs (`0024`–`0026`), 3 experiments (`0020`–`0022`), 1 feature request (`0010`).

---

## 2026-07-16 (later) — new finding: Database Access "Generate Credentials" is one-shot with no recovery path

**Source:** `field-notes/ACTIVITY-LOG.md`, field-notes workspace Developer panel. Attempting to use direct
Postgres access (rather than another build job) to populate `content_items` hit a wall: credentials had
been generated once earlier in the session and the connection string was never captured; every later
click of "Generate Credentials" 409s with `"Credentials already exist. Use regenerate to rotate them"`,
but no regenerate/view control exists anywhere in the UI — confirmed via a full interactive-element scan
of the panel, not just visual inspection. Filed as `BACKLOG.md #19` / `blog/bugs/0023`. **Filed with
Audos** the same day, via Otto chat inside the workspace — but only after Otto independently confirmed,
from the platform's own side, that it has no tool to view/rotate these credentials either and that no
"view" path exists for anyone. Priority bug submitted, engineering notified, $27 bounty attached if
confirmed.

## 2026-07-16 (evening) — field-notes v0: shell-escape confirmed YES, independently verified — correcting our own false alarm along the way

**Source:** `field-notes/ACTIVITY-LOG.md` (Job #83055, then a scoped retry Job #83114, both in field-notes'
own workspace, `1d30572d-2ced-4dd1-872f-3e67a74891dd`). This entry replaces the interim "left inconclusive"
version written earlier the same day — the missing piece turned out to be a mistake in our own
verification, not a real gap in the platform.

**Resolved: the shell-escape experiment — YES, confirmed independently.** Briefed upfront in a genuinely
fresh workspace whose `Desktop.tsx` had never been touched by any prior fix, the app was born full-screen
from the first paint — no chat-shell flash, no retrofit needed. This is not the job's word for it: we
checked the workspace's own Preview panel ourselves, Draft vs Live side by side, logged in with the test
credentials by hand, reached the real content feed, and clicked the chat-affordance ourselves to confirm
it opens a small popup rather than tearing down to the old three-pane view. Full detail in
`docs/platform/22-eliminating-the-chat-shell-playbook.md`.

**Retracted the same day: `BACKLOG.md #18` (file writes not durable on crash) was a false alarm.** The
interim entry reported that a build job's file writes vanish if the job errors before some final step.
They don't. We'd checked `GET /api/space/{id}/files` right after a job errored, got an empty file list,
and concluded the work was lost. That endpoint only ever reflects the *published* bundle — it has no
visibility into draft/unpublished work. The app was sitting in Draft the whole time; a follow-up job
confirmed the original job's app file, config, and hooks were already present when it started, rather than
rebuilding from scratch. We'd already filed this with Audos and it had been escalated to a human engineer
— followed up on the same ticket to correct the record before real engineering time got spent on it.
**Genuine, smaller finding that survives:** there's no obvious way to tell "nothing was built" apart from
"something was built but only exists as an unpublished draft" without knowing to check the Preview panel's
Draft toggle specifically — worth Audos documenting, but it isn't data loss.

**Still standing from the interim pass:**
- `BACKLOG.md #17` — a possible raw-SQL restriction bypass in `db.rawQuery` (validator may only check the
  leading keyword, not the full statement). Unconfirmed whether it actually executed a blocked write;
  flagged for Audos to verify directly.
- `BACKLOG.md #15`/`#16` — the orphaned-table bug and the undocumented serial-id behavior. Both real,
  both filed, unrelated to the file-durability retraction.

**Net effect:** `field-notes` now has a real, working v0 in Draft form — password gate, content feed,
verified end-to-end by hand, not yet published. Database and app layers both confirmed. Next step is
populating `content_items` via the planned sync script and publishing.

**Narrative:** worth a follow-up to `blog/0013-building-field-notes-in-the-open.md` — the retry, the
independent verification, and the correction, told the same honest way as the rest of that entry.

---

## 2026-07-16 (later) — Round 3 complete: embedding path found (correcting round 3's own first pass), dimensionality closed, no conventions file exists

**Source:** `doknow-kb/audos/ACTIVITY-LOG.md` (round-3 dispatch, job #82974, split into two sub-jobs by
Otto); `field-notes/ACTIVITY-LOG.md` (TEST G directly informs the field-notes build this round was
dispatched alongside).

Round 3 ran as two sub-jobs. The first (TEST E) is corrected in place below by the second
(TEST F + a second, more thorough embedding pass) — this entry replaces, not appends to, the interim
finding written earlier the same day.

**TEST E, corrected: a real embedding path exists on Audos — the first pass's "confirmed absent" was a
false negative.** The first sub-job swept `platform.integrations.isAvailable()` with guessed provider
names (`openai-embeddings`, `vector-search`, `pinecone`, etc.) — all `false` — and concluded no embedding
capability existed. It never called `platform.integrations.proxy()` itself and never tried the bare
provider name. The second sub-job did: `isAvailable('openai')` → `true`, and
`platform.integrations.proxy('openai', '/v1/embeddings', {...})` **returned a real 1,536-float OpenAI
embedding, no API key supplied by the calling hook** — Audos holds and injects the credential through a
generic authenticated passthrough (`proxy`) to a fixed allowlist (`openai, stripe, twilio, heygen`).
There is still no dedicated `platform.generateEmbedding()` convenience method, but the underlying
capability is real, working, and simply undocumented.

**TEST F: dimensionality re-benchmark, real embedding size (1,536 floats).** Same brute-force cosine
similarity method as round 2, same 50/300/1,000 row counts, now at real size instead of 5-float
placeholders: ~0.003ms/row scan time at 300–1,000 rows, full scan 1–3ms even at 1,000 rows. ~307x more
arithmetic per comparison than the original test barely moved the needle — wall-clock is dominated by the
database fetch (~250–540ms), not the comparison math. **Dimensionality is no longer an open question.**

**TEST G: no AGENTS.md-equivalent exists, and a static conventions file is proven not to reach at least
one generation path.** The workspace file tree has no `AGENTS.md`, `.cursorrules`, `.cursor/`, or
`CLAUDE.md`. A behavioral test — create an `AGENTS.md` with an explicit convention, then generate a hook
via `platform.generateText` with a neutral prompt — showed the convention was **ignored**; a direct probe
asked the same generation endpoint if it could see any conventions file, and it answered "NO WORKSPACE
FILE CONTEXT AVAILABLE." This is proven only for the `generateText` path — the editor-agent/Cursor-lane
flow (the one that does most of DoKnow's actual app building) couldn't be behaviorally tested from within
a single run, since the file has to exist before the job starts. Scoped finding, not a full answer.

**Corrected:**
- `docs/platform/06-capabilities-reference.md` — Vector/Embedding Storage section rewritten: both
  dimensionality and embedding-generation gaps now closed, with the corrected `proxy()` finding front and
  center and the false-negative explicitly flagged.
- `BACKLOG.md #14` — corrected in place from "feature request: build embedding support" to "feature
  request: document the existing `platform.integrations`/`proxy()` path" — a materially smaller, easier
  ask than what was filed a few hours earlier the same day.
- `BACKLOG.md #13` — updated: both benchmark gaps closed, scope narrowed to "native vector column/ANN
  index for headroom," not a current blocker.

**Net effect on the DoKnow build:** the ingestion pipeline's embedding step has a real, working, no-key
path today via `platform.integrations.proxy('openai', '/v1/embeddings', ...)`. Combined with round 2's
scheduler-workaround finding, there is no longer a credible architectural reason to build any part of
ingestion off-platform.

**Narrative:** `blog/0012-what-audos-can-actually-do.md`, full round-3 update section.

---

## 2026-07-16 — Capability test round 2: scheduler resolved, vector-search question reopened (not answered)

**Source:** `doknow-kb/audos/ACTIVITY-LOG.md` (round-2 dispatch and results, job #82919).

**Corrected:**
- `docs/platform/06-capabilities-reference.md` — Vector/Embedding Storage section: pgvector is
  "available but not installed" (not "installed but not enabled" as first reported); `vectorscale`
  0.9.0 (DiskANN) found in the same state, missed in round 1. JSON-array brute-force similarity
  benchmarked at 50/300/1000 rows — sub-millisecond throughout (0.02–0.28ms avg). This is a real,
  valid result for what it measured (row-count scaling; cosine arithmetic costs the same regardless
  of what the numbers mean) — but it used fake 5-float vectors, not real embeddings, so it does **not**
  establish that semantic search is viable here. Two gaps remain open: realistic dimensionality
  (~300x the tested size) was extrapolated, not measured; and no test has ever generated a real
  embedding on-platform. See "Correction, 2026-07-16 (later same day)" below.
- `docs/platform/06-capabilities-reference.md` — Scheduler Integration section: round 1's "scheduled
  hooks never fire" is now "genuinely inconsistent" — a daily schedule fired correctly on 2026-07-16
  (~9s of nextRun) after an hourly schedule failed to fire at all on 2026-07-14. Root cause unknown.
  Client-orchestrated sequential hook calls (5/5 succeeded, zero friction) proven as a reliable
  substitute regardless of scheduler behavior — this finding holds, not affected by the correction below.
- `BACKLOG.md #13` (vector storage) — status: **open, pending round 3** (not downgraded — see correction).
- `BACKLOG.md #12` (scheduler) — reframed from "confirmed broken" to "confirmed inconsistent," with
  the client-chaining workaround noted as removing it as a hard blocker. This finding stands.
- `blog/0012-what-audos-can-actually-do.md` — updated in place with a dated correction section (status
  chip changed `open` → `pass`), rather than left standing as the pre-round-2 conclusion.

**New artifact:** `docs/platform/reports/2026-07-16-vector-search-experiment.html` — full experiment
report (hypothesis, both rounds' tests, results table, revised architecture), written for a non-technical
read-through as well as the raw numbers.

**Net effect on the DoKnow build (as first written, same day):** the ingestion pipeline could stay
entirely on Audos. **This was overstated — see correction below.**

**Narrative:** `blog/0012-what-audos-can-actually-do.md`'s "Update, 2026-07-16" section.

---

**Correction, 2026-07-16 (later same day):** the "can stay entirely on Audos" conclusion above conflated
two different things — that comparing arbitrary numeric arrays is fast (true, proven) and that real
semantic vector search is viable on Audos (not tested at all). The scheduler finding is unaffected and
stands as corrected above; only the vector-search conclusion is walked back. Status changed from
"answered" to **reopened, pending a round-3 test**: (1) whether Audos exposes a native embedding
capability the way it exposes `generateText`, and (2) a re-benchmark at realistic (~1536-float)
dimensionality. All four touched files (`docs/platform/06`, `BACKLOG.md #13`, this entry,
`blog/0012.md`, and the HTML report) corrected in place same day.

---

## 2026-07-15 — Capability test: ingestion pipeline verdict, plus the terminology/escalation-path findings

**Source:** `archive/ACTIVITY-LOG-2026-07-15T1904Z.md` (full raw log for this period).

**Corrected:**
- `docs/platform/06-capabilities-reference.md` — three claims fixed against live verification:
  storage endpoint (`/api/storage/upload` doesn't exist; real endpoints are `/api/upload/file` and
  `/api/upload/image`), scheduler behavior (one-time hook scheduling doesn't exist; recurring schedules
  don't actually fire), and `platform.generateText`'s real signature/default model. Added an explicit
  "unverified until proven otherwise" caveat to the rest of the file, which predates this SDK by ~3
  months and had never been checked against a live workspace before this pass.
- `docs/platform/22-eliminating-the-chat-shell-playbook.md` — corrected an overclaim: the "One Good
  Thing" experiment's "all four issues avoided from the start" conclusion was three-quarters
  inheritance from an already-fixed shared shell, not a real test. Only the JSON-write finding was
  genuinely unconfounded.
- `blog/0010-building-it-right-the-first-time.md` — same correction, applied in place with a visible
  callout rather than a silent rewrite.

**New findings, durable docs added:**
- `docs/platform/06` — new "Vector/Embedding Storage" section: pgvector is installed on the server but
  not enabled, no path exists to enable it. JSON-array + brute-force similarity is the only fallback.
- `docs/platform/18-genesis-space-and-ui-ceiling.md` — new terminology section making workspace-vs-app
  explicit (a recurring source of confusion in conversation, not in the underlying facts).
- `docs/platform/25-escalation-and-support-paths.md` — new doc: `request_human_help` is a paid VA queue,
  not a bug/support channel; Priority Support (browser UI Help panel) is the real escalation path, with
  a confirmed tiered automated-fix-then-engineer flow.

**New backlog items:** `BACKLOG.md #10` (no way to cancel a running `cursor_delegation` job — later
corrected: the specific instance wasn't actually hung, just slow), `#11` (`/api/upload/image` silently
corrupts non-image payloads — a real, dangerous bug), `#12` (scheduled hooks confirmed to never fire),
`#13` (no native vector storage, no way to enable pgvector).

**Net effect on the DoKnow build:** the ingestion pipeline (upload → transcribe → chunk → embed →
retrieve) is not buildable on Audos as originally specced — the plumbing (upload, external fetch, AI
generation) works, but the two pieces the whole design depends on (vector search, background
processing) don't. `docknow-kb/docs/mvp-gap-analysis.md` updated with the concrete implication: either a
small-scale synchronous-only version, or the ingestion/embedding layer built as a separate off-platform
service.

**Narrative:** `blog/0012-what-audos-can-actually-do.md`.
