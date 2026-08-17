# Audos Capability Matrix

A living, at-a-glance answer to "can Audos actually do X" — synthesized from every bug, experiment, and
finding in this SDK, so nobody has to read 76 individual entries to piece it together. **Update this
file whenever a finding changes a row's status** — that's the whole point of keeping it separate from
the narrative log. See `docs/platform/24-where-new-findings-go.md` for the general rule; this file is
the "durable synthesis" version of that rule, one row per capability instead of one file per finding.

## Status key

| Symbol | Meaning |
|---|---|
| ✅ | Verified working, by us, live |
| ⚠️ | Works, but with a real caveat — read the note before relying on it |
| ❌ | Confirmed not possible, or broken |
| 📄 | Documented by Audos; **not yet independently verified by us** — treat as a claim, not a fact |
| ❓ | Open question — tested but inconclusive, or asked and never answered |

*Last synthesized: 2026-08-05. Every row links to its source finding — if a row looks wrong, check the
source's own date before trusting this table over it; this file can drift, the source findings are the
ground truth.*

---

## Database

| Capability | Status | Note | Source |
|---|---|---|---|
| Generic REST CRUD (`db-api`: query/insert/update/delete) | ✅ | Works against any workspace table, no per-table setup. Update/delete require ≥1 filter. | `docs/platform/25`, `sdk/src/index.ts` |
| Direct Postgres access (dev role) | ⚠️ | DML only — `CREATE`/`ALTER`/`DROP TABLE` all rejected with permission errors. Schema changes still require Otto. | `bugs/0034` |
| DB Access credential rotation | ❌ | One-shot "Generate Credentials," no regenerate/view control. Audos claimed fixed 2026-07-17; independently re-tested 2026-07-22 — **still broken**, claim is false. | `bugs/0023` |
| `id` as UUID primary key | ❌ | Always platform-generated `serial` integer; explicit `uuid` id rejected. `session_id`/`updated_at` auto-injected on every table. | `bugs/0020` |
| Failed `CREATE TABLE` rollback | ❌ | A failed call (e.g. bad FK constraint) leaves an orphaned, undroppable table under that name. No cleanup tool exists. | `bugs/0019` |
| Native vector/embedding column (pgvector) | ❌ | Not installed. **Re-verified live 2026-08-05:** `pg_extension` returns exactly one row, `plpgsql`. `vectorscale` 0.9.0 sits in `pg_available_extensions` with `installed_version: null`. No enable path. JSON-array + brute-force fallback is fine at current scale — a ceiling on headroom, not a blocker today. | `bugs/0015`, `blog/experiments/0033` |
| Vector similarity search (JSON-array fallback) | ✅ | Brute-force cosine similarity in a hook, benchmarked at real embedding size (1,536 floats): sub-3ms at 1,000 rows. Wall-clock dominated by DB fetch, not math. | `experiments/0013` |
| Embedding generation | ✅ | `platform.integrations.proxy('openai', '/v1/embeddings', ...)` — real embeddings, no API key needed from the caller. Undocumented; no dedicated convenience method. | `docs/platform/06` (Vector/Embedding Storage section, corrected) |
| Raw SQL restriction (`SELECT`/`WITH`/`EXPLAIN` only) | ❓ | Validator may only check the leading keyword, not the full statement — a possible bypass. Never confirmed to actually execute a blocked write; flagged for Audos to verify directly. | `blog/experiments/0015` |
| Workspace schema conventions | ✅ | `ws_<workspace-uuid>` schema, `app_` table prefix, platform infra lives in a separate DB entirely (not workspace-touchable). | `blog/experiments/0025` |

## AI / Generation

| Capability | Status | Note | Source |
|---|---|---|---|
| Text generation (`ai-api` / `Chat()`) | ✅ | Works, `systemPrompt` supported, uncapped response length. | `sdk/skills/otto-pilot/`, `docs/platform/06` |
| Text generation (`Generate()` call type) | ⚠️ | Same model, but a hard ~1,000-token response cap. Use `Chat()` instead for anything longer. | `bugs/0035` |
| Model selection (in-app AI hook) | ✅ | **Corrected 2026-08-10, then solved 2026-08-14.** Whole OpenAI GPT-4/3.5 family works. **GPT-5 (all sizes), o1/o3/o3-mini/o4-mini also work** via a corrected hook that sends `max_completion_tokens` — verified live returning real output on Audos's own OpenAI key. The stock `ai-api` hook still blocks them (`bugs/0040`); the fix is a new hook alongside it (`sdk/hooks/ai-api.reference.js`, live as `ai2-api` in `workspace-156396`). Claude/Fable/Gemini genuinely 404 (OpenAI proxy only; they run on the Cursor/Otto build surface, not this API). | live audit 2026-08-10, job #108698 verified 2026-08-14, `bugs/0040` |
| Reasoning-model token budget | ⚠️ | Reasoning models spend the budget on **internal reasoning** before any visible text. `gpt-5` at `maxTokens: 2000` returned **empty text** (all 2000 consumed reasoning); at 4000 it produced real output. Use ≥4000. The SDK client defaults reasoning models to 4000 and errors clearly instead of returning "". | job #108698, measured 2026-08-14 |
| Vision / image understanding (in-app hook) | ✅ | **Corrected 2026-08-10** — works on `gpt-4o`/`gpt-4.1`. Pass `messages[].content` as an OpenAI-style array with an `image_url` part (URL or `data:` URI). Verified: red PNG → "Red". | live audit 2026-08-10 |
| Reasoning / effort control (in-app hook) | ⚠️ | The reasoning models that support effort (o3, gpt-5) are reachable-but-hook-blocked (above); once a hook sends `max_completion_tokens` they should be usable. Effort params on the current hook are ignored. | `bugs/0040` |
| Streaming responses | ❌ | Not supported anywhere in the AI hook. Use loading states. | `experiments/0028` |
| PDF text extraction via AI hook | ✅ | Routes to OpenAI's Chat Completions endpoint, extracts PDF text natively. | `bugs/0035` |
| DOCX (or other non-PDF) extraction via the same path | ❌ | Rejected — OpenAI's own API only accepts PDF on that content-block type. | `bugs/0035` |
| Audos Code (`portfolio/code`) — Claude Code backend | ✅ | A second, distinct build surface from the Cursor pipeline — genuine token-by-token streaming during execution, model-routing panel visible. | `blog/experiments/0020` |

## Auth / Identity

| Capability | Status | Note | Source |
|---|---|---|---|
| Custom app-level auth as the *first* thing a visitor sees | ❌ | Every unified-space workspace's signed-out view is unconditionally Audos's own `EmailGate.tsx`. Custom auth can only ever be a second, downstream gate. No config flag exists to disable/reorder it. | `docs/platform/26` |
| EmailGate email-capture side effects | ⚠️ | Submitting the email step registers a CRM contact and fires Meta/Reddit ad "Lead" pixels — not a neutral auth screen. | `bugs/0024` |
| OTP verification | ⚠️ | Real, works once enabled — but only reachable via an undocumented raw `PUT` API call, no settings UI. | `feature-requests/0015` |
| OAuth / SSO | ❌ | Not supported at all. OTP is the only native security lever. | `feature-requests/0015` |
| Platform-detection global (`window.__WORKSPACE_ID__`) | ✅ | Always injected — the reliable "is this running on Audos" signal. `window.__spaceContext` (bare, for this purpose) is never actually set. | `blog/experiments/0024` |
| User-identity accessor (`window.__spaceContext?.username`) | ✅ | Populated once EmailGate completes — a *different* property from the platform-detection check above, don't conflate the two. | `blog/0002` |
| `useWorkspaceDB`/`useSession` React hooks | ⚠️ | Work, but default `shared: false` silently scopes queries to the current browser session — a common footgun. Audos-hosted-app only; don't work once code is ported off-platform. | `blog/0002`, `experiments/0026` |
| Multi-tenant isolation pattern | ✅ | Column-based (`org_id`/`user_id` on every table, filtered at the data layer) is the workable pattern — not workspace-per-customer, which the platform nudges you toward by default. | `blog/0002` |

## Email

| Capability | Status | Note | Source |
|---|---|---|---|
| Send email (HTML + plain text, `replyTo`) | ✅ | Works as documented. | `docs/platform/06` |
| CC / BCC | ❌ | Not supported. | `feature-requests/0020` |
| Send log / delivery confirmation | ❌ | Not supported. | `feature-requests/0020` |
| Custom From address | ❌ | Not supported. | `feature-requests/0020` |
| Scheduled email (`/schedules/email`) | 📄 | Documented by Audos; not independently tested by us. | `docs/platform/06` |

## Storage / Files

| Capability | Status | Note | Source |
|---|---|---|---|
| Multipart file upload (`/api/upload/file`) | ✅ | Round-trips correctly for any content type. | `bugs/0017` |
| Base64 image upload (`/api/upload/image`) | ❌ | **Silently corrupts non-image payloads** — returns `200 success` with a plausible URL, but the stored bytes don't match what was sent. Only use this endpoint for real images. | `bugs/0017` |
| `useSpaceFiles()` client hook | 📄 | Documented (persist JSON to workspace storage); not independently verified by us. | `docs/platform/06` |

## Scheduling / Background Jobs

| Capability | Status | Note | Source |
|---|---|---|---|
| Recurring scheduled hooks | ⚠️ | Genuinely inconsistent — a daily schedule fired correctly, an hourly one never fired at all in the same test window. Root cause unknown. | `docs/platform/06` (Scheduler Integration, corrected) |
| Client-orchestrated sequential hook chaining | ✅ | A reliable substitute regardless of scheduler reliability — 5/5 succeeded, zero friction. | `blog/experiments/0010` |
| One-time (non-recurring) scheduled hook | ❌ | Doesn't exist as a primitive. | `docs/platform/06` (corrected 2026-07-15) |

## App Build & Deploy

| Capability | Status | Note | Source |
|---|---|---|---|
| Full-screen, single-app landing (no chat shell) | ✅ | `desktop.layout.defaultLandingView: "app"` + `defaultLandingAppId` makes the **base space URL** boot straight into one app. | `docs/platform/19`, `docs/platform/22` |
| Deep-linking a specific app via hash (`#app-id`) | ❌ | Still opens inside the chat shell as a side panel — the full-screen escape only works for the base URL, not per-app deep links. Separately, `Desktop.tsx`'s hash router matches the *entire* hash, so `#app/subpath` doesn't even match `#app`. | `docs/platform/19`, `bugs/0032` |
| Config-driven layout (declare shape in `config.json`) | ❌ | No such field exists — achieving full-canvas requires imperative `Desktop.tsx` surgery, which is then vulnerable to the regression risk below. | `feature-requests/0012` |
| File preservation (mark a file "don't regenerate") | ❌ | No `.audosignore` equivalent, no manifest flag. A hand-customized `Desktop.tsx` got silently reverted by a platform-initiated commit **at least three times**. | `bugs/0031`, `feature-requests/0011` |
| GitHub Dev Mode sync → live | ⚠️ | Source lands and Sync Activity says "done" well before the actual compile finishes — an undocumented ~15–30 min latency window. Misdiagnosing this as "broken" cost a full day once. | `feature-requests/0016` |
| New app registration via GitHub push → auto-compile | ❓ | One test showed a newly-registered app never compiled without a manual trigger — but the original report explicitly flagged this as possibly just "didn't wait long enough," not confirmed as a separate code path. | `experiments/0029` |
| Adding a manually-chosen CDN package with a React peer-dependency | ❌ | Audos only auto-pins `?deps=react@...` for its own pre-configured packages. Anything added by hand silently loads a second React copy → blank screen, error #31. | `bugs/0030` |
| ESBuild compile scope | ⚠️ | Compiles **every** `.tsx` file in the app directory, not just the entry point's import tree — an unused sibling import still loads at runtime. | `experiments/0023` |
| `import.meta.env` (Vite convention) | ❌ | Not exposed by Audos's ESBuild pipeline — Vite-only dev fixtures don't port as-is. | `experiments/0027` |
| Frontend libraries: Tailwind, Radix primitives, TanStack Query, GSAP, raw three.js | ✅ | All confirmed rendering live via `cdnDependencies` (esm.sh, React 18.3.1 pinned). | `docs/platform/19` |
| `@react-three/fiber` | ⚠️ | Renders, but cold-start init is slow/inconsistent — flaky, not a hard fail. Prefer raw three.js where reliability matters. | `docs/platform/19` |
| shadcn/ui (CLI codegen) | ❌ | No `npm install`/CLI available — hand-authored shadcn-style TSX+Tailwind works as a substitute. | `docs/platform/19` |
| TanStack Router (browser URL ownership) | ❌ | The space runtime owns real URL/history/mount — an app can use TanStack Router in-memory, but can't own the actual browser URL bar. | `docs/platform/19` |
| Arbitrary npm / Vite / pnpm / React 19 / Vitest | ❌ | Pipeline is single-file + ESBuild + CDN only, React pinned to 18.3.1. | `docs/platform/19` |
| Build execution backend | ⚠️ | Build/edit jobs delegated to **Cursor Background Agents** — gated by *Audos's own* Cursor account usage limit (`usage_limit_exceeded`, transient, not the workspace owner's problem to fix). Chat stays up during a build outage. | `docs/platform/21`, `sdk/skills/otto-pilot/` |
| Launching a build from the external onboarding/chat API | ⚠️ | Limited to the Cursor path — Audos Code (the other build surface) requires a signed-in human/browser session and can't be triggered from the API alone. | `docs/platform/19` |
| DB_API_KEY exposure in client-side apps | ❌ | Found hardcoded in compiled client bundle for any CDN-deployed app using direct DB-hook access — fully readable in devtools. Pre-existing platform pattern, not something porting introduced. | `bugs/0037` |
| Hardcoded `aud_live_` client secrets | ❌ | Same shape of exposure — never rely on client-side secrecy for these tokens. | `sdk/skills/otto-pilot/` |
| A job's own "Complete"/success self-report | ❌ | **Cannot be trusted without independent verification.** Repeatedly reproduced: builds that report success while nothing shipped, platform support tickets marked "Completed" with false claims. The one standing rule everything else in this matrix was checked against. | `docs/platform/23`, `bugs/0027`, `0028`, `0029` |

## Analytics & Reporting

*Otto-chat-triggered tools, not `platform.*` hooks callable from app code — see `docs/platform/29`.*

| Capability | Status | Note | Source |
|---|---|---|---|
| `query_analytics` (overview: sessions/contacts/conversion) | ✅ | Verified live against DoKnow: 5 sessions, 4 contacts, 80% conversion. Contact count cross-checked against `query_data_source(contacts)` independently — matched. | `docs/platform/30`, `blog/experiments/0030` |
| `query_events` (by type/day/space/app) | ✅ | Verified live: 187-event breakdown by type, plausible against known dev/QA usage. | `docs/platform/30`, `blog/experiments/0030` |
| `query_data_source(sourceId=contacts)` | ✅ | Returned 4 real records (UUIDs, emails), matched the `query_analytics` count exactly. | `docs/platform/30` |
| `query_data_source(sourceId=funnel-events)` | ⚠️ | Returned 1 record, every field blank — contradicts `query_events`'s 187-row result for the same workspace/window. Use `query_events`, not this, for event data. Root cause unconfirmed. | `docs/platform/30` |
| `query_data_source` (other sourceIds: `ad-campaigns`, `printify-orders`, `stripe-*`, `analytics-overview`, `session-recordings`, `community-posts`) | 📄 | Enum confirmed from live tool schema; not independently exercised yet. | `docs/platform/29` |
| `get_funnel_metrics`, `query_sessions`, `delegate_analytics_insight` | 📄 | Real params confirmed from live schema; outputs not yet run. | `docs/platform/30` |

## Ads & Marketing

*Otto-chat-triggered tools, not `platform.*` hooks callable from app code — see `docs/platform/29`.*

| Capability | Status | Note | Source |
|---|---|---|---|
| `get_ad_campaigns` (list) | ✅ | Correctly returned empty for a workspace with no campaigns — clean empty result, not an error. | `docs/platform/31`, `blog/experiments/0031` |
| `keyword_ideas` (Google Ads data) | ✅ | Returned 10 real, internally-consistent keyword rows (volume/trend/competition/bid) for DoKnow-relevant seed terms. | `docs/platform/31` |
| `generate_ad_copy` | ✅ | Produced genuinely on-target, non-generic ad copy referencing the real product's pain point; works cold with no prior ad history. | `docs/platform/31` |
| `search_meta_targeting` | ⚠️ | `"City, State"` format (e.g. "Austin, TX") silently fails to resolve. Bare city name works but can return 5+ same-named cities across states (disambiguation required); ZIP code resolves to exactly one match. **Prefer ZIP.** | `docs/platform/31`, `blog/experiments/0031` |
| `get_campaign_insights`, `get_dm_campaign_status`, `get_dm_conversations` | 📄 | Real param shapes confirmed from schema; nothing to query yet (no campaigns exist on this workspace). | `docs/platform/31` |
| `delegate_ad_generation` → `launch_previewed_campaign` (actually launch a paid campaign) | ❓ | Deliberately not exercised — real money, real audience. Preview-then-confirm gate exists per Otto's description; not independently verified. | `docs/platform/31` |

## Media Generation

*Otto-chat-triggered tools, not `platform.*` hooks callable from app code — see `docs/platform/29`.*

| Capability | Status | Note | Source |
|---|---|---|---|
| `generate_image` | ✅ | Real GCS URL, independently downloaded and rendered — a correct 1536×1024 PNG matching the exact prompt. Paid (`gpt-image-1.5`). | `docs/platform/32`, `capabilities/media/generate-image.md` |
| `generate_video` (Veo3) | ✅ | Real GCS URL, independently downloaded and `ffprobe`-verified — genuine 8s MP4. Paid, but cost couldn't be isolated (see note below). | `docs/platform/32`, `capabilities/media/generate-video.md` |
| `generate_voiceover` (ElevenLabs) | ✅ | Real GCS URL, independently downloaded and `ffprobe`-verified — genuine 8.6s MP3. Only media-gen cost cleanly itemized: `$0.03 — ElevenLabs TTS: 116 chars`. | `docs/platform/32`, `capabilities/media/generate-voiceover.md` |
| `list_voiceover_voices` | ❌ | Returns ~85 real voices but **every `ID` field is literally `undefined`** — only 3 hardcoded fallback voices (Sarah/Rachel/Josh) are actually selectable via `voiceId`. | `capabilities/media/generate-voiceover.md` |
| `generate_background_music` | 📄 | 8 real presets confirmed (`list_music_presets`); generation itself not run. | `capabilities/media/generate-background-music.md` |
| `search_stock_photos` (Unsplash) | ✅ | Real photo IDs + 4 size variants each, free. | `capabilities/media/search-stock-photos.md` |
| `redesign_logo`, Instagram carousel tools | 📄 | Schema confirmed; not run. | `capabilities/media/other-media-tools-pending.md` |
| Media-gen cost itemization | ⚠️ | Only voiceover (ElevenLabs TTS) gets a discrete wallet line item. Image and video draw from a shared, non-itemized "AI token usage" bucket — no per-generation dollar figure exposed by Audos's own tooling. | `capabilities/media/generate-video.md` |

## Otto / Onboarding API

| Capability | Status | Note | Source |
|---|---|---|---|
| Onboarding flow (`/start` → `/verify` → build) | ✅ | Works; 11 build stages in practice (Audos's own docs say 7 — undercounted). | `sdk/skills/otto-pilot/` |
| Auth via request body (`{authToken: ...}`) | ✅ | The reliable form — always worked. | `bugs/0029` |
| Auth via `Authorization: Bearer` header | ⚠️ | Was broken for weeks (401 on working endpoints), independently confirmed fixed 2026-07-10. Verify live before depending on it rather than trusting docs. | `bugs/0029` |
| `createNew: true` for a returning, authenticated user | ⚠️ | Re-triggers a fresh OTP even with a valid auth token — "skip OTP for returning users" only applies when *not* forcing a new workspace. | `sdk/skills/otto-pilot/` |
| `chatId` as a named conversation thread | ❌ | It's a correlation ID only. A made-up `chatId` doesn't create/select a thread — Otto just returns whatever's currently active. No documented create/list/resume API. | `feature-requests/0018` |
| Creating a server function from outside a chat session | ❌ | The externally-reachable Otto has no `manage_server_functions` tool, contradicting `docs/platform/07` which documents it as the creation path. Only route found is staging a Cursor delegation job. Otto's own self-report, not independently verified. | `bugs/0038` |
| Retrieving a build job's full output | ❌ | The jobs list truncates output **mid-JSON**, not at a record boundary and with no marker. Deleting a temp hook destroys the `get_hook_logs` fallback. Same shape as `BACKLOG.md #4`. | `bugs/0039` |
| Support ticket structured API / attachments | ❌ | Doesn't exist — tickets are plain chat text, structured evidence gets flattened into an unreadable block. | `bugs/0036`, `feature-requests/0017` |
| Self-serve DNS changes for custom domains | ❌ | Requires a support ticket every time, even for one record. Otto's own DNS tools are read-only. | `feature-requests/0019` |

---

## How to keep this current

Whenever a new bug/experiment/feature-request resolves or changes a capability question, update the
matching row here in the same pass — same discipline as promoting a finding into `docs/platform/`. If a
capability doesn't have a row yet, add one. If Audos ships a fix that changes a row's status, don't just
trust their word for it — verify live the way every ⚠️/❌ row above already was, then update the status
and the source link together.

**For a whole new area of tools** (like Analytics, Ads, or Media Generation): write a short index doc in
`docs/platform/` (one row per tool, linking onward) plus one focused file per real, distinct capability in
`docs/platform/capabilities/<area>/` — not a giant single doc, and not one file per raw helper/list tool.
This is the same progressive-disclosure pattern `sdk/skills/otto-pilot/` uses for its own `references/`
folder: cheap to skim the index, load only the one capability file you actually need. See
`docs/platform/29-otto-tool-surface-vs-app-callable-hooks.md` for the architecture note this pattern grew
out of, and `docs/platform/30`–`32` for worked examples.
