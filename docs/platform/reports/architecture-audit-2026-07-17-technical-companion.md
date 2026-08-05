# Audos Platform — Technical Companion to the Architecture Audit

*Engineering annex to `architecture-audit-2026-07-17.md`, dated 2026-07-17. Where the audit answers the
strategic question ("is Audos architected for agentic app development") and ranks root causes by leverage,
this document is the implementation-level reference: exact failure modes, reproductions, root-cause
hypotheses, blast radius, current workaround, and a concrete fix spec for each — subsystem by subsystem.
Written for an engineer (Audos's or a technical evaluator's) who needs to act on the findings, not just
read the verdict. Every entry maps to a cited finding in `audos-platform`.*

---

## 0. How to read this

- **Severity** is operational impact on a technical team building real software, not user-facing polish:
  - **S1 — blocks or silently corrupts.** No safe workaround, or a workaround that itself fails silently.
  - **S2 — structural friction.** Workaround exists but is costly, manual, or fragile; paid on every cycle.
  - **S3 — papercut / documentation.** Real, but cheap to route around once known.
- **Workaround** states what actually unblocks the issue today, and whether it's reliable.
- **Confidence** flags anything not fully confirmed — carried verbatim from the source so nothing here
  reads as more certain than the evidence supports.
- Section 1 is the scannable register. Sections 2.x are the deep-dives. Section 3 is the prioritized
  engineering backlog. Section 4 is the verification methodology (so findings are reproducible). Section 5
  is a correction log, including a precision fix made to the audit itself.

---

## 1. Issue register

| ID | Subsystem | Issue | Sev | Workaround | Filed | Source |
|----|-----------|-------|-----|-----------|-------|--------|
| J-1 | Job dispatch | Job self-reports ("verified") don't match live state | S1 | Independently re-verify every job live | no | `bugs/0007` |
| J-2 | Job dispatch | Completion reports truncate at the critical line | S2 | Read full report from Tasks-board UI, not API | no | `bugs/0008` |
| J-3 | Job dispatch | No progress/cancel; hang indistinguishable from slow | S2 | Dispatch-then-poll SOP; wait it out | yes(#10) | `bugs/0018` |
| J-4 | Job dispatch | Same-app-scope jobs serialize silently, no queue signal | S2 | Never run two jobs on one app scope | no | `bugs/0018`,`0005` |
| J-5 | Job dispatch | `usage_limit_exceeded` on Audos's shared Cursor account | S2 | Retry later; owner can't clear it | no | `bugs/0006` |
| P-1 | State/publish | `GET /files` reflects **published** bundle only; draft invisible | S1 | Check Preview panel Draft/Live toggle by hand | no | `bugs/0022` |
| P-2 | State/publish | "Published: yes" means 3 different things | S2 | Cache-bust + eyeball live | no | `bugs/0005` |
| P-3 | State/publish | Stale served bundle after publish | S2 | `?_cb=<ts>` cache-bust | no | `bugs/0005` |
| P-4 | State/publish | GitHub push→live is 15–30 min; webhooks silently dropped | S2 | Manual "Sync from GitHub"; poll file API | no | `docs/16` |
| D-1 | Database | Failed `CREATE TABLE` leaves orphaned, undroppable table | S1 | Recreate under a new name | yes | `bugs/0019` |
| D-2 | Database | `id` forced to `serial`; `uuid` PK rejected; 2 columns injected | S2 | Design FKs as `integer`; expect extra cols | yes | `bugs/0020` |
| D-3 | Database | No vector type / no pgvector / no `CREATE EXTENSION` | S1 | Brute-force cosine in JS (small scale only) | no | `bugs/0015` |
| D-4 | Database | `json`/`jsonb` writes fail unless `JSON.stringify`'d (×3) | S2 | Brief `JSON.stringify` on every write | n/a | `docs/19` |
| D-5 | Database | Direct-Postgres credentials one-shot, no view/rotate path | S1 | None through UI; blocked | yes(#19) | `bugs/0023` |
| D-6 | Database | Possible `rawQuery` leading-keyword-only validation | S2? | — | yes(#17) | `bugs/0021` |
| R-1 | Runtime | Dual-React: user packages omit React pin → error #31, blank | S1 | Inspect importmap; pin `?deps=react@18.3.1` | no | `docs/16` |
| R-2 | Runtime | Compiles every `.tsx` in app dir regardless of imports | S2 | Keep no dead/broken files in app dir | no | `docs/16` |
| R-3 | Runtime | Sandbox gaps: no `URLSearchParams`/`Buffer`/`setTimeout`/`fs` | S2 | Hand-rolled helpers per hook | n/a | `docs/07`,`08` |
| R-4 | Runtime | Hash routing only; no real browser-URL ownership | S3 | Query-param router workaround | no | `docs/16` |
| I-1 | Integrations | `proxy()` body must be JSON **string**; object → silent 200+HTML | S1 | Always `JSON.stringify` the body | n/a | `experiments/0012` |
| I-2 | Integrations | Undocumented surface; `isAvailable()` gives false negatives | S2 | Try bare provider name + `proxy()` directly | n/a | `experiments/0011`→`0012` |
| I-3 | Integrations | `/api/upload/image` silently corrupts non-image payloads | S1 | Use `/api/upload/file` (multipart) for non-images | no | `bugs/0017` |
| I-4 | Integrations | Can't read platform code (`search_platform_code` empty) | S3 | Empirical probe apps | no | `feature-requests/0009` |
| S-1 | Scheduler | Recurring schedules fire inconsistently; no diagnostics | S2 | Client-orchestrated sequential hook calls | no | `bugs/0016`,`exp/0009` |
| S-2 | Scheduler | No one-time hook scheduling at all | S3 | Email-scheduling endpoint only | no | `bugs/0016` |
| U-1 | UI shell | Chat-shell default; escape needs a 4-step checklist | S2 | The `docs/22` playbook (inherits once fixed) | some | `bugs/0001–0004` |
| U-2 | UI shell | `Desktop.tsx` is workspace-shared & mutable (blast radius) | S2 | Treat shell edits as workspace-global | no | `experiments/0005` |
| SEC-1 | Identity | EmailGate open by default; OTP configurable only via API | S1 | Enable OTP via browser-console `fetch` | no | `docs/17` |
| SEC-2 | Identity | Sessions don't appear to expire; no sign-out primitive | S2 | Implement TTL + sign-out over localStorage | no | `docs/17`* |

*Filed = raised with Audos (Priority Support / BACKLOG item). "n/a" = a usage discipline, not something to
file. Severity is this annex's assessment, not the source's.*

---

## 2. Subsystem deep-dives

### 2.1 Job dispatch & observability — the black-box boundary (audit Gap 1)

**Architecture.** Every structural change (table, hook, app edit) is a natural-language instruction to
Otto, which stages a draft Task and dispatches it to a **Cursor Background Agent**. The dispatcher returns
a task ID; results arrive only on completion — this path has no streaming channel. *(The newer Audos Code
surface does stream a live step-log; see J-1b.)*

**Backend migration — read this before treating any Cursor-specific finding as current.** The path above
is **Cursor-backed**. The newer Audos Code surface runs on a **different backend — Anthropic's Claude Code**
(model visualizer routes through `claude-opus-4-8`; picker lists Opus 4.8 / Fable 5 labeled "Claude Code /
Anthropic"; Fable 5 shows a transparent "0% of daily budget used" quota, not the opaque shared-account
failure of `bugs/0006`). Audos appears **mid-migration** between the two. Consequences for this register:
J-2's Cursor-stream-truncation root cause, J-5's shared-Cursor-account `usage_limit_exceeded`, and the
`cursorModel`-pinning method (`experiments/0004`) are **legacy-path-specific** and may not apply to where
Audos is heading. They're flagged inline.

**J-1 — self-reports vs. reality (S1).** Jobs reported `complete` / `published` / verbatim `"verified —
three lesson rows inserted with valid JSON arrays"`, and the live app reproduced the *identical* failure
(`Insert failed: invalid input syntax for type json`) on immediate re-test — twice, worded the same
(`bugs/0007`, `docs/platform/19`). *Root-cause hypothesis (Otto's own):* an **authority-boundary gap** —
the Cursor agent reports truthfully on its sandbox, but confirmation gates and publish/recompile are
enforced Audos-side, outside the agent's visibility, so sandbox-success is reported as ground-truth.
*Fix spec:* surface the platform-side gate/publish result back into the job's completion context, and/or
run a mandatory post-build smoke check (load route → assert no thrown error → assert target DOM node)
before a job may emit "verified" (`feature-requests/0001`).

**J-1b — same failure mode, newer surface, different backend. Verified live 2026-07-17 (~25 min).** Audos
Code (Beta 0.3.0, `audos.com/portfolio/code`) is a genuine improvement on *observability*: a Settings
toggle ("Assistant output — show token-by-token output") is **on by default**, and during a run the left
panel streams a typed step-log (`Command run`, `Tool call`, `File change`, a live duration counter, a
specific completion line). It also runs on a **different backend**: the model visualizer routes through
`claude-opus-4-8` and the picker lists Opus 4.8 / Fable 5 labeled "Claude Code / Anthropic" — Audos is
migrating off Cursor onto Anthropic's Claude Code (see the §2.1 architecture note). **But the
completion-*trust* problem is unchanged, and that's the point.** Reproduced in under two minutes: asked to
change a headline's `.` to `!`, it streamed real steps, ran 29s, and reported `"Done. The headline now reads
'Stop saving. Start knowing!' … Nothing else was touched."` with an "Up to date" status; the rendered
preview still showed `.`, confirmed by zoom and after a manual refresh. Same defect as J-1 — on the *new*
backend, in the *more observable* UI. The self-report/reality mismatch is not a Cursor artifact and not an
observability artifact; it survives both. *Implication:* the J-1 fix spec (a trustworthy verification gate)
is not obviated by the live log — the log shows *activity*, not *correctness*. Notably, Audos Code's in-app
browser-test runner — exactly that gate — is currently shown *paused* ("disabled … while we make its
results reliable enough to trust"), so this is a fix Audos already appears to be part-way toward.
*(Maturity: Audos Code failed to load **four times across ~40 minutes over two sessions on 2026-07-17** —
"could not be loaded, try again" or a stuck traffic spinner — a reproducible reliability signal, not a
one-off.)*

**J-2 — truncation (S2, legacy-path-specific).** Completion summaries cut off mid-sentence ≥8 times, at the
root-cause/diff/publish line (`bugs/0008`). *Root cause:* upstream — the Cursor Cloud Agents API truncates
its own stream with no documented full-log endpoint. **This attribution is specific to the Cursor backend;
Audos Code streams a live step-log on the Claude Code backend (J-1b), so this particular failure may not
carry forward.** *Note:* still an Audos architectural choice on the legacy path — it could persist the full
transcript server-side. *Workaround:* the untruncated report is readable from the Tasks-board UI
(`get_page_text`), just not from the dispatch/chat API.

**J-3 — no progress/cancel (S2).** `list_jobs` shows status but its outcome/error field stays empty until
`Complete`/`Failed`; `get_hook_logs` covers only server-function hooks; `get_audos_code_status` only Audos
Code threads; the harness's `TaskStop`/`TaskList` cannot see jobs-board tasks (`bugs/0018`). A job ran
45m01s doing 67+ tool actions and was briefly misread as hung — the point being a real hang is
indistinguishable from slow-but-working, with no abort either way. *Fix spec:* live log stream + a cancel
endpoint on the dispatch API.

**J-4 / J-5.** Two same-scope jobs serialize (the second sits `Queued`, silently, no queue-position
signal — `bugs/0018`); dispatch can instant-fail `usage_limit_exceeded` against **Audos's own** shared
Cursor account, unclearable by the workspace owner and with no ETA (`bugs/0006`). *(J-5 is
legacy-path-specific: Audos Code's Claude Code backend showed a transparent per-model "0% of daily budget
used" quota instead — this failure mode looks already-addressed there.)*

---

### 2.2 State, publish & bundle — no authoritative "what exists now" (audit Gap 2)

**P-1 — the published-only file API (S1).** `GET /api/space/{id}/files` and `GET
/api/space/{id}/file/config.json` reflect **only the published bundle**. For an app built but sitting in
Draft, both return empty (`{"files": []}` / `{"content": ""}`) — *identical to "nothing was built."* This
directly caused a false data-loss bug report, filed with Audos and escalated to an engineer before being
retracted the same day; ground truth was only visible via the Preview panel's **Draft vs. Live** toggle
(`bugs/0022`, `field-notes/ACTIVITY-LOG.md`). *Fix spec:* a state API that distinguishes draft from
published (and ideally returns both), so "built-but-unpublished" and "not built" are not the same
response. *(Partial credit, newer surface, live 2026-07-17: Audos Code exposes an explicit
`audos.com/site/{id}?env=live` URL that disambiguates draft vs. live in the UI. The API-level opacity of
`GET /files` is presumably unchanged for programmatic integrators — who are exactly who this fix is for.)*

**P-2 / P-3 — "published" is ambiguous (S2).** `list_apps` can show an app registered + `published: yes`
while the served bundle is stale (app absent from dock, route 404) — observed when a concurrent same-scope
job held the publish (`bugs/0005`, `docs/platform/19`). A manual `?_cb=<epoch-ms>` cache-bust forces the
fresh bundle; the platform itself sometimes emits `?_cb=…&cdn=fallback` under undocumented conditions
(`feature-requests/0004`). *Fix spec:* a publish-status endpoint returning the served bundle's
hash/timestamp (`feature-requests/0003`); auto-cache-bust on every publish.

**P-4 — the GitHub deploy pipeline (S2).** For the Dev-Mode path: `push → webhook → sync → ESBuild → CDN`,
"push to live 15–30 min," and "occasionally a push webhook is silently dropped by Audos — no error, no
indication," remediated only by a manual "Sync from GitHub" click (`docs/platform/16`). No progress
events. The prescribed verification is to poll the raw file API and diff against local.

---

### 2.3 Database & schema — non-transactional, undocumented, indirectly accessed (audit Gap 3)

**D-1 — orphaned tables, no rollback (S1).** `workspace_db_create_table` with an unsatisfiable FK fails
(`Failed to create table: foreign key constraint "fk_app_site_sessions_access_id" cannot be implemented`)
**without rolling back** — leaving a physical `app_site_sessions` relation the catalog doesn't track.
Every retry then fails `relation "app_site_sessions" already exists`; the orphan is invisible to
`db_list_tables`/`db_describe_table`, and no `drop_table` tool exists. Only escape: a new name
(`site_sessions_v2`). The original name stays blocked until Audos clears it server-side (`bugs/0019`,
`BACKLOG.md #15`). *Fix
spec:* atomic DDL (roll back partial `CREATE TABLE` on constraint failure) or, minimum, a `drop_table`
cleanup tool.

**D-2 — silent schema rewriting (S2).** `id` is always a platform-generated `serial` integer; requesting
`id: uuid` returns `column "id" specified more than once`. Every table silently gains `session_id` (text)
and `updated_at` (timestamp default `NOW()`). Column types are enum-restricted to
`text/integer/bigint/decimal/boolean/timestamp/date/json/uuid` — none documented with this caveat
(`bugs/0020`). *Consequence:* any FK must be typed `integer` to reference a serial PK — a uuid FK is what
triggered D-1. *Fix spec:* honor requested PK types (or document the constraint prominently); disclose
injected columns.

**D-3 — no vector search (S1).** *Corrected in the audit; stated precisely here.* There is **no native
vector column type and no pgvector.** `workspace_db_create_table` rejects `vector(N)`
(`invalid_enum_value`); a `::vector` cast fails `type "vector" does not exist`. `pg_available_extensions`
lists `vector` 0.8.1 and `vectorscale` 0.9.0 at `installed_version: null` — **available in the catalog but
not installed** — and `CREATE EXTENSION` can't be run because `db.rawQuery` is `SELECT`/`WITH`/`EXPLAIN`-
only. The *only* mechanism is a **brute-force cosine loop in hook JavaScript** over embeddings stored as
`json` arrays. The benchmark (1,536-float, 50/300/1000 rows, ~1–3ms full scan, dominated by DB fetch not
math — `experiments/0013`) establishes that the **arithmetic is cheap at small row counts** — nothing
more. It does **not** establish retrieval quality, and `experiments/0008` explicitly walked back that
overreach; `bugs/0015` calls the JSON fallback "not a real retrieval index." **Net: viable at DoKnow's
scale (tens–hundreds of chunks/course); not a vector search capability, and not for anything larger.**
*Fix spec:* install/enable pgvector (+ optionally vectorscale/DiskANN) workspace-side, or expose a managed
vector column type.

**D-4 — json serialization footgun (S2).** Every *first* write to a `json`/`jsonb` column in a new code
path failed `Insert failed: invalid input syntax for type json` because a raw JS object/array was sent
instead of a `JSON.stringify`'d string — three times, two unrelated apps (`docs/platform/19`). *Fix spec:*
accept objects and serialize server-side, or return a typed error naming the fix.

**D-5 — one-shot credentials (S1).** The Developer panel's "Generate Credentials" (`POST
/api/workspaces/{id}/db-credentials`) works once; thereafter it returns `409 Conflict: "Credentials
already exist. Use regenerate to rotate them."` while **no regenerate/rotate/view control exists anywhere
in the UI** (full interactive-element scan confirmed one button). `GET` on the same endpoint returns
`401`, so it is not a read-back path. Otto has no tool to view/rotate it either (separate short-lived
token), and confirms the connection string is shown once and never stored retrievably (`bugs/0023`,
`experiments/0019`). **This credential flow has now failed three distinct ways, zero successes, across
three workspaces:** field-notes (2026-07-16) → the 409 loop above; DoKnow's own workspace (2026-07-17,
tried this cycle) → a **silent `404` with no error toast at all** (`POST
/api/workspaces/workspace-156396/db-credentials`), worse than field-notes' visible failure; Throughline
(April, `docs/platform/13`) → "Workspace not found." *Fix spec:* wire the frontend to the regenerate
capability the backend's own 409 references, and fix whatever makes the endpoint 404/409 rather than issue
a usable credential. **Among the highest capability-per-effort fixes we found** — it would restore the one
direct, transactional SQL path that eases D-1/D-2/D-3/D-6 at once, *assuming direct SQL is as full-powered
as the April Throughline doc reports (not re-verified this cycle — 0-for-3 on even reaching it — see §5).*
The 0-for-3 raises priority: "it works once you're in" cannot currently be asserted at all. **Scope,
though:** this blocks the *direct psql* path only — **database writes work** via Audos Code's chat-driven
`db.insert` layer on the Claude Code backend (verified live 2026-07-17: field-notes' `content_items`
populated 0 → 6 rows, confirmed via Preview feed + the live site's Field Scout assistant, no publish step —
see the new "content population" note in §4/§5). So D-5 restores a standard interface and a convenience,
not the only way to get data in.

**D-5b — content population via Audos Code works (positive, S-none).** The credential bug (D-5) does not
gate DB writes in general. Audos Code's `db.insert`/MCP tool layer (Claude Code backend) inserted 6 real
corpus rows into `content_items` from empty, verified two independent ways (Preview feed shows all six with
correct type/status chips; the live site's Field Scout answered "latest bugs" with the exact rows,
immediately, no publish step). Two footnotes: a long free-text paste didn't reach the agent (it flagged the
gap rather than fabricating — good honesty discipline; compact single-line JSON went through), and only 6
of ~64 items were populated in this proof-of-path batch (scales via JSON, but a full sync is worth a hook
that takes a JSON array in one call rather than a conversational round-trip). *Separate S3 quirk:* the
"Findings" quick-link chip inside a Field Scout reply opens a blank side panel (reproduced twice); the feed
renders fine through the normal Preview panel.

**D-6 — possible rawQuery bypass (Confidence: UNCONFIRMED).** Evidence that `db.rawQuery`'s
`SELECT`/`WITH`/`EXPLAIN` guard may validate only the *leading* keyword, so `SELECT 1; DROP …` could chain
a write. ~20 min of probing (transaction-scoped read-only wrappers, `set_config('transaction_read_only',
…)`, a third `options` arg) **never confirmed an actual blocked write executed** — the job abandoned it and
worked around D-1 differently. Flagged for Audos to verify against their own validator; **not** asserted as
an exploit (`bugs/0021`, `BACKLOG.md #17`). *Fix spec:* parse the full statement (reject multi-statement
strings), regardless of whether this specific attempt succeeded.

---

### 2.4 Build pipeline & runtime — a fixed, non-standard target (audit Gaps 6, 9)

**Model.** One file per app (`apps/<Name>/App.tsx`), server-side **ESBuild**, dependencies via
**importmap/CDN (esm.sh)** — no `npm install`, no Vite, no bundler config. **React 18.3.1** shared host
instance. `cdnDependencies` resolves ESM libs with `deps`/`external` pinned to the host React.

**R-1 — dual-React (S1).** The importmap adds `?deps=react@18.3.1` only to pre-configured packages;
user-added packages resolve *their own* React, silently loading two copies → **React error #31**, blank
page. "The importmap looked correct at a glance. The only way to catch it was to inspect the actual network
resources" (`docs/platform/16`, `inbox/INCIDENT-REACT-ERROR-31.md`). *Fix spec:* pin the host React on all
importmap entries, or fail the build on a detected duplicate.

**R-2 — compile-everything (S2).** ESBuild compiles *every* `.tsx` in the app dir regardless of the import
graph, so a stray dead file with a bad import crashes the live app (`docs/platform/16`).

**R-3 — sandbox gaps (S2).** Hooks are sandboxed JS (**not Node**, no `require`), ~30s timeout, ~64KB code
size, ~1MB response; missing `URLSearchParams`, `Buffer`, `process.env`, `setTimeout`/`setInterval`, `fs`,
`crypto`, and a real `Headers` (`response.headers.get` is not a function). Every hook re-implements a
`buildQuery` helper; workspace IDs are hardcoded because env vars don't exist (`docs/platform/07`, `08`).

**R-4 — routing (S3).** `Desktop.tsx` exact-matches the URL hash; apps get in-memory hash routing only, no
real browser-URL ownership (TanStack Router works in-memory only). A query-param router workaround is
documented (`docs/platform/16`). Opt-in real routing is a standing request (`feature-requests/0008`).

**Build-path gap (audit Gap 9).** There is no *raw* file-level editing — Audos Code's "Full in-platform
editing" (Platform mode) is, on live inspection 2026-07-17, still a scoped natural-language-to-agent
request (select element → describe change → agent executes), not a hand-typable editor; the March-era docs'
"can't edit React code directly, must use Otto" (`docs/platform/02`, `03`, both 2026-03-31) still holds on
that narrow point. The entire 8-endpoint API suite is hand-authored from templates (`docs/platform/09`);
and when Cursor is over its limit, we couldn't reach Audos Code from the external API (needs a signed-in
browser session, `feature-requests/0007`) — leaving an API-driven agent with no build path in that window.

---

### 2.5 Integrations & the undocumented surface (audit Gaps 7, 8)

**I-1 — `proxy()` silent-success (S1).** `platform.integrations.proxy(provider, path, opts)` is a generic
authenticated passthrough to a fixed allowlist (`openai/stripe/twilio/heygen`, read from its own error
message); Audos injects the credential. The body **must** be a `JSON.stringify`'d **string** — passing a
plain object **silently returns the platform's HTML index page with status 200**, indistinguishable from a
real response until you inspect the content (`experiments/0012`). Correct call:
```js
platform.integrations.proxy('openai', '/v1/embeddings', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ model: 'text-embedding-3-small', input: '…' }),
}); // → 200, real 1,536-float embedding, no caller API key
```
*Fix spec:* accept an object body, or return a typed 4xx on a malformed body instead of a 200 + HTML.

**I-2 — false-negative capability probes (S2).** `platform.integrations.isAvailable()` returns `false` for
guessed feature-names (`openai-embeddings`, `vector-search`) *and* for some genuinely-installed
integrations — so it is not a reliable "does X exist" oracle. The embedding path was first reported
**absent** for exactly this reason; it exists (`isAvailable('openai') → true` with the *bare provider
name*, then `proxy()` directly). *Discipline:* probe the generic passthrough with bare provider names, not
`isAvailable()` feature-name guesses (`experiments/0011` → `0012`). *Docs gap:* `platform.integrations`,
`platform.externalFetch`, `platform.getLatestSession`, `session`, `respondRaw`, `db.bulkInsert` are all
live and undocumented; and `eval`/`Function`/`Proxy`/`Reflect`/`WebAssembly` are present despite docs
claiming `eval` is blocked (enumerated, not invoked).

**I-3 — image-upload corruption (S1).** `/api/upload/image` (base64) returns `HTTP 200 {"success":true}`
with a plausible `.png` URL while **corrupting** any non-image payload: a PDF went in md5 `7fc99347…` /
3,031 bytes and came out md5 `ec78847b…` / 3,050 bytes, zero common prefix. A real PNG round-trips
md5-identical; `/api/upload/file` (multipart) round-trips any content type correctly to a ~50 MiB−1 cap
(`bugs/0017`). *Rule:* never send non-images through `/api/upload/image`. *Fix spec:* validate content
type and reject (or correctly store) non-image payloads instead of silently mangling them.

**I-4 — no platform code read (S3).** `search_platform_code` → "no searchable directories";
`read_platform_file` path-denies compiler paths — even for Otto — forcing empirical probe apps for
questions (like the `cdnDependencies` allowlist) that should be a cheap read (`feature-requests/0006`,
`0009`).

---

### 2.6 Scheduler & async (audit Gap 5)

**S-1 — inconsistent firing (S2, Confidence: inconsistent-not-broken).** Two hourly recurring schedules
(`POST /schedules`, valid RRULE, `status: pending`) **never fired** — 2h07m and 20m past `nextRun`,
`runCount: 0`, `lastError: null`, underlying hook confirmed working via direct call (`bugs/0016`,
2026-07-14). A **daily** schedule later **fired within ~9s** of `nextRun` (`experiments/0009`,
2026-07-16). Same mechanism, opposite results; root cause (cadence? since-fixed? per-workspace?) unknown,
unknown to Otto too. There is no dispatcher log, queue position, or disabled flag to diagnose from. *Fix
spec:* observable per-run scheduler state + a root-caused firing guarantee.

**S-2 — no one-time hook scheduling (S3).** `frequency` is required, rejecting one-time payloads; one-time
scheduling exists only on the email endpoint (`bugs/0016`).

**Workaround (proven).** Client-orchestrated sequential hook calls: 5/5 succeeded, ~1–1.8s each, zero
rate-limiting/auth/concurrency friction (`experiments/0010`). Covers anything that can keep a client open;
genuinely unattended fire-and-forget work still has no reliable path, and the ~5–10 min hook ceiling caps
any single synchronous unit of work.

---

### 2.7 UI shell & isolation (audit Gap 4)

**U-1 — the chat-shell escape (S2).** The default signed-in surface is a chat shell; `Desktop.tsx` renders
shell-mode by default and full-screen is a *different state of the same component*. The 4-step escape
(`docs/platform/22`): (1) `config.json` → `desktop.layout.defaultLandingView: "app"` +
`defaultLandingAppId`; (2) route `#app-id` deep links through the full-screen transition, not
`openPanel()`; (3) make the chat affordance a popup, not `returnToAgentView()` (which tears down to the
old three-pane UI); (4) resolve initial state in a **lazy `useState` initializer** reading
`window.location.hash` synchronously — a plain `useState(default)` + `useEffect` correction produces a
visible flash-of-old-shell (FOUC) on first paint, catchable only by a **zero-delay** screenshot.

**U-2 — shared mutable shell (S2).** These fixes live in a **workspace-shared** `Desktop.tsx`, so once
applied they inherit to every app for free (`experiments/0006`, `0017`) — but the same sharing means a
shell change has blast radius across every app, and it **confounded a real experiment** ("One Good Thing"
looked like an upfront-briefing success but was three-quarters inherited — `experiments/0005`). App-level
isolation is genuinely murky; the un-root-caused cross-app "cards due" discrepancy (`BACKLOG.md #9`) is a
thin but live hint. *Fix spec:* a first-class full-screen app mode needing no shell surgery; clearer
per-app isolation.

---

### 2.8 Identity & security (audit Gap 10)

**SEC-1 — open-by-default identity (S1). Re-confirmed live 2026-07-17.** EmailGate default: "anyone who
knows a user's email address can enter it and access that user's data." Reproduced on DoKnow's signed-out
live site — a fabricated, never-used address (`audit-verify-2026-07-17@example.com`) signed in instantly,
no OTP, no verification, landing on a real dashboard. OTP exists but is **off by default and exposed
through no settings UI** — a full settings-surface scan (2026-07-17) found no OTP toggle anywhere; "the
only interface is the API directly," via a browser-console `fetch()` (`docs/platform/17`). *Fix spec:*
secure default (OTP on, or email-only access off) + first-class UI/API config.

**SEC-2 — sessions (S2). No-sign-out re-confirmed live 2026-07-17; session-expiry not re-tested.** There
is **no built-in sign-out primitive** (the avatar menu does nothing; sign-out is done by deleting
`localStorage` keys) and **no OAuth** on the sign-in modal — both re-confirmed live this cycle. Whether
sessions ever expire was **not** re-testable in a quick pass (needs waiting out an unknown TTL), so that
one sub-claim stays sourced to `docs/platform/17` (2026-04-23), not asserted.

---

## 3. Remediation plan — prioritized engineering backlog

Ordered by leverage (capability unlocked per unit of effort), not by severity count.

**P0 — highest leverage, do first**

1. **Fix the DB credential flow (D-5).** It failed three ways this cycle (409 loop, silent 404, "Workspace
   not found"), 0-for-3 — so this is "make the credential step work at all," not just "wire up regenerate."
   *Largest unlock:* it would restore direct transactional SQL, which would ease D-1, D-2, D-3, and D-6 —
   *with the caveat that "direct SQL is full-powered" is an April Throughline observation not re-verified
   this cycle (0-for-3 on even reaching it — see D-5 / §5), so the payoff is a strong inference, not a
   proven one.* The concrete instance of the audit's "open a power-user lane" recommendation, and the
   0-for-3 makes it more urgent, not less.
2. **Finish and re-enable the in-app test runner (J-1 / J-1b).** A trustworthy post-build verification gate
   before "verified," plus surfacing platform-side gate/publish state into the completion context. Framing
   matters here: Audos Code **already ships** this as a browser-test runner, currently *paused* ("disabled
   … while we make its results reliable enough to trust") — so the ask is "get the thing you paused over
   the line," not "build something new." We reproduced the exact self-report/reality mismatch it would catch
   (J-1b), which is the case for prioritizing it.
3. **Authoritative state API (P-1, P-2).** A publish-status endpoint (served bundle hash/timestamp) and a
   file/state API that distinguishes draft from published. Removes the browser-and-eyeball step from every
   change and makes "what exists now" machine-checkable.

**P1 — structural, high value**

4. **Atomic DDL or a `drop_table` cleanup tool (D-1).**
5. **Honor/disclose schema behavior (D-2)** — requested PK types, injected columns, documented.
6. **Install/enable pgvector (+vectorscale) or a managed vector column (D-3)** — the only path to real
   vector search; the brute-force JS scan is a small-scale stopgap, not a capability.
7. **Live job log stream + cancel endpoint (J-3);** persist untruncated transcripts server-side (J-2).
8. **Kill silent-success failure modes (I-1, I-3, P-2):** honest status codes; content-type validation on
   upload; typed 4xx on malformed `proxy()` bodies. "200 OK" must mean success.

**P2 — de-risking & papercuts**

9. **Root-cause the scheduler (S-1)** and expose per-run status; add one-time hook scheduling (S-2).
10. **Dual-React fix (R-1):** pin host React on all importmap entries or fail on duplicate.
11. **Document the undocumented surface (I-2, I-4):** `platform.integrations`/`proxy()`, the live globals,
    the `cdnDependencies` allowlist; give workspace-scoped platform-code read access.
12. **Secure identity defaults + config UI (SEC-1, SEC-2).**
13. **Opt-in modern stack (R-3, R-4):** current React, real dependency resolution, real URL ownership —
    non-default, for teams that want it (`feature-requests/0008`).
14. **First-class full-screen app mode + per-app isolation (U-1, U-2).**

---

## 4. Verification methodology (for reproducers)

The findings above are only as good as how they were checked. An Audos engineer reproducing them should
apply the same disciplines the corpus adopted (often after being burned by not doing so):

- **Verify live, never on self-report.** A job's `Complete`/`verified` is not evidence; reload the actual
  published surface and redo the user action (`docs/platform/19`). This is the standing rule and it roughly
  doubles the cost of each cycle — by design, because the alternative is shipping on a false signal.
- **Check the initial paint, not just settled state.** Screenshot with **zero delay** immediately after
  navigation — a `useEffect`-corrected FOUC is invisible to a screenshot taken a few seconds later
  (`experiments/0016`).
- **Distinguish draft from published.** `GET /files` is published-only; use the Preview panel's Draft/Live
  toggle to see draft work (`bugs/0022`). Empty file API ≠ nothing built.
- **md5 both ends of any upload.** Silent corruption looks like success (`bugs/0017`).
- **Dispatch-then-poll, don't block.** Fire the job in the background, poll a compact status every ~90s up
  to ~20–30 min, log each raw response, keep working meanwhile (`docs/platform/23`).
- **Probe the generic passthrough, not just `isAvailable()`.** Feature-name guesses give false negatives;
  try the bare provider name and call `proxy()` directly (`experiments/0011` → `0012`).
- **Verify the verification.** At least twice in this corpus a correction itself rested on an unverified
  assumption (the file-durability retraction relied on a job's self-report about intermediate state that
  is now unrecoverable — `field-notes/ACTIVITY-LOG.md`, 21:10). Treat a retraction with the same rigor as
  the original claim.

---

## 5. Correction log

- **Vector search (D-3), corrected in the audit and stated precisely here.** The audit's *What's working
  well* section originally read "brute-force vector search is fast enough at real scale" / "semantic search
  is viable." That conflated three separate things and overstated all three: (1) embedding **generation**
  works (real, via `proxy()`); (2) there is **no vector search and no pgvector** — only a brute-force JS
  cosine scan over JSON arrays; (3) the benchmark proved **cheap arithmetic at small row counts**, not
  retrieval quality and not scaling. Corrected in both the `.md` and `.html` audit files, and reflected in
  D-3 above. The user flagged this directly ("there is no vector search"), and they were right.
- **pgvector state.** Also corrected: an earlier audit phrase called pgvector "installed but unenabled." It
  is **available in the catalog but not installed** (`installed_version: null`) — a distinction the corpus
  itself had already corrected once (`CHANGELOG.md`, round-2 entry) and which this annex now matches.
- **Confidence flags carried forward, not smoothed over.** D-6 (rawQuery bypass) is **unconfirmed** — not a
  proven exploit. S-1 (scheduler) is **inconsistent, not broken**. U-2 / `BACKLOG.md #9` (cross-app data
  isolation) is an **un-root-caused hint on thin evidence**. Any of these hardening into a firm claim
  requires a fresh, dedicated reproduction.
- **Audos Code added after a ~25-min live drive (2026-07-17).** A whole newer editing surface — Audos Code
  (Beta 0.3.0, `audos.com/portfolio/code`) — was absent from the original corpus and first draft. Driven
  live, it turned out to be a real improvement in *observability* (token-by-token streaming step-log,
  default-on; live draft preview; an explicit `?env=live` URL that clears up the draft/published ambiguity
  in-UI). It also **runs on a different backend — Anthropic's Claude Code, not Cursor** — so Audos looks
  mid-migration, and several Cursor-specific findings (J-2 truncation, J-5 usage-limit, `cursorModel`
  pinning) are now flagged legacy-path-specific. But the completion-trust defect reproduced there in under
  two minutes on that new backend (the period-to-exclamation-point edit) — which makes Gap 1 *stronger*,
  not weaker. Audos Code was **not** exhaustively evaluated (and failed to load four times across ~40
  minutes over two sessions on 2026-07-17 — young, and a reproducible reliability signal in its own right).
- **"Direct SQL is full-powered" downgraded to 0-for-3 (2026-07-17).** The positive claim was cut from
  "works, only the credential flow is broken" to an honest "designed to work, reportedly worked in April
  (Throughline), but every attempt to *reach* it this cycle failed differently — field-notes 409, DoKnow
  silent 404, Throughline 'Workspace not found' — three workspaces, zero successes." It is not a
  verified-working escape hatch today; it's a documented intention with a 0-for-3 record on the credential
  step. This *raises* the priority of the D-5 fix rather than lowering it.
- **Identity (SEC-1/SEC-2) re-verified live 2026-07-17.** EmailGate open-by-default (fabricated-email
  instant sign-in, no OTP), no OTP settings UI, no built-in sign-out, and no OAuth were all re-confirmed
  live this cycle — staleness flags dropped on those. The single sub-claim still unverified is whether
  sessions ever expire (needs waiting out a TTL); it stays sourced to April.
- **Staleness caveat, narrowed.** After the 2026-07-17 re-verification, the claims still resting on
  March–April docs and **not re-verified this cycle** are: Gap 6 runtime specifics (single-file / React 18
  / hash routing / no-npm — the quick re-check was blocked by Audos Code's repeated load failures) and the
  "sessions never expire" sub-claim. These carry inline "sourced [month], not re-verified" flags; Gap 10's
  other identity claims and Gap 9's file-editing point were re-confirmed live and no longer do.
- **Cross-check complete (2026-07-17).** The load-bearing factual claims here were sent back to the
  originating DOKNOW-DEV session and confirmed: D-3 (vector search) facts correct and the "strength"
  framing rightly dropped; I-1 (embedding generation + JSON-string-body gotcha) accurate and to be kept
  distinct from search; S-1 (scheduler) is "inconsistent, root cause unknown," not broken or fixed; D-6
  (rawQuery) stays explicitly unconfirmed. Two nuances it added are also folded in: the file-durability
  retraction *itself* later rested on an unverified second-job self-report (a deeper Gap-1 example than the
  original bug), and the "shell-escape from birth" verdict came from a second hands-on check after the
  first job's self-reported success was lost to a crash — not from trusting the job.
