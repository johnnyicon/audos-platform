---
name: otto-pilot
description: Create and operate AI-powered startup workspaces via the Audos API — onboarding through ongoing build/chat with Otto. Use when user wants to start a business, build an MVP, validate a startup idea, create a company workspace, launch a product, monitor an in-progress Audos build, or chat with Otto programmatically. Triggers on requests like "I have a business idea", "help me start a company", "create a startup workspace", "I want to build [product]", "check my Audos build status", or "talk to Otto via the API".
---

# Otto-Pilot: Audos Workspace Builder (API v1.2) — augmented

Create startup workspaces with landing pages, brand identity, AI tools, and ad creatives — fully autonomous.

> **Started as Audos's own onboarding skill file, corrected and extended with what six months of real
> building on Audos actually found — and outgrew "onboarding" along the way.** It now covers the whole
> lifecycle: standing up a workspace, monitoring a build, and the ongoing chat/auth patterns for
> operating one after it exists, not just the first ten minutes. See `original.md` (same folder) for
> the unmodified version we started from — this file changes what was wrong, adds what was missing, and
> points out to a deeper reference library for anything beyond the essentials, rather than cramming
> everything into one document. **Last verified:** 2026-07-12 via a live onboarding run (new workspace
> "DoKnow").
>
> This whole folder (`sdk/skills/otto-pilot/`) is meant to be portable — copy it into any project's
> skills directory and it works standalone; the `references/` files it points to travel with it.

## Progressive Disclosure — Read Only What You Need

Everything below this point is what an agent needs on *every* onboarding call. For anything past that,
read the specific file for the situation — don't load all of it upfront.

| Situation | Read this |
|---|---|
| A build job fails, stalls, or you're not sure whether to trust its status | `references/build-execution-and-failures.md` |
| A build/edit job fails with `usage_limit_exceeded`, or you want to understand the execution backend | `references/build-execution-and-failures.md` |
| Auth is behaving unexpectedly, or you need to understand `chatId` | `references/auth-and-chatid.md` |
| Filing a bug report, or handling `aud_live_` tokens safely | `../../OTTO_API_WORKFLOW.md` (wider SDK — needs the full `sdk/` folder, not just this skill) |
| Calling the actual workspace hooks (db/ai/email/web/storage/analytics/crm) once a workspace exists | `../../src/index.ts` (TypeScript) or `../../go/client.go` (Go) — real, working client code |

## Base URL

```
https://audos.com/api/agent/onboard
```

## URL Construction

The API returns URLs using the current deployment domain:

```json
"urls": {
  "landingPage": "https://audos.com/site/184582",
  "workspace": "https://audos.com/space/workspace-184582"
}
```

Use these URLs directly — no domain swapping needed.

## Quick Reference

| Action | Method | Endpoint |
|--------|--------|----------|
| API docs | GET | / |
| Start onboarding | POST | /start |
| Verify OTP | POST | /verify |
| Check build status | GET | /status/:workspaceId |
| Check build status (alt) | POST | /status |
| Chat with Otto | POST | /chat |
| Chat with Otto | POST | /chat/:workspaceId |
| Rebuild (if failed) | POST | /rebuild/:workspaceId |

## Authentication

- **Token format:** `aud_live_xxxx` (48 hex chars after prefix)
- **Auth tokens never expire** — store persistently by email
- **Session tokens** expire in 30 min (only needed during OTP flow)
- **Default to body-auth** (`authToken` in the request body) — this has been reliably confirmed
  working. Bearer-header auth (`Authorization: Bearer <token>`) is documented but was broken for weeks
  before being fixed; verify it live before depending on it rather than assuming the docs are current.
  See `references/auth-and-chatid.md` for the full history.

## Conversation Flow

### Introducing Audos

When a user expresses a business idea, briefly explain what Audos does before asking for their email:

> "I can help you build that with Audos! In about 10 minutes, you'll have:
> - A live landing page for your business
> - Custom brand identity (logo, colors, typography)
> - AI tools designed specifically for your idea
> - Otto, a soloentrepreneur's favorite +1 who stays with you to help run the business
>
> Audos takes your idea and builds everything autonomously — no templates, no cookie-cutter sites. Everything is custom to your business.
>
> To get started, what email should I use for your account?"

### New Users Flow
1. **Collect** user's email + business idea
2. **Start** → `POST /start` (sends 4-digit OTP to email)
3. **Verify** → `POST /verify` with OTP code → returns `authToken` + starts build
4. **Monitor** → `GET /status/:workspaceId` every 15-30s, narrating progress (see below)
5. **Watch for** `landingPageReady: true` (~10 min) — core build done
6. **Introduce Otto** and offer to chat

### Returning Users (have workspace)
1. **Start** → `POST /start` with email
2. **Response includes** `auth_token` + `urls` directly — skip OTP!
3. **Chat** → `POST /chat/:workspaceId` immediately

> **Caveat (verified 2026-07-12):** "skip OTP" only applies when `/start` targets the
> **existing** workspace. Sending `createNew: true` — even with a valid Bearer auth token
> for a returning user — starts a *new* workspace build and **re-triggers a fresh 4-digit
> OTP** (returns a `sessionToken`; you must `POST /verify` before the build begins). Do not
> assume an authenticated returning user can create a second workspace without OTP.

## Polling During Build — UX Guidelines

**Critical:** The build takes ~6–10 minutes. Users MUST see progress updates or they'll think it's stuck.

Poll every 15-20 seconds (not 60s), and send a progress message immediately after each poll — don't wait
until done. Use `parallelBuildStatus[]` for per-step task breakdown (`✅` done, `🔄` in progress, `⏳`
pending), and watch `landingPageReady === true` **and** `status === 'complete'` as the reliable "done"
signal — earlier `*DraftReady` / `coreStepsComplete` flags mean "usable draft," not finished.

## API Reference

### POST /start
Fields: `email` (required), `businessIdea` (required, min 10 chars), `businessName` (optional),
`targetCustomer` (optional), `callbackUrl` (optional, webhook with HMAC signing), `createNew` (optional).

Returns new user: `sessionToken`. Returns returning user: `auth_token`, `workspaceId`, `urls`. Returning
user + `createNew: true`: **`sessionToken`** — a fresh OTP is sent; verify before the new build starts.

### POST /verify
Body: `{ sessionToken, otpCode }`.
Returns: `workspaceId`, `authToken`, `urls`, `buildInfo` (`totalSteps: 11`, `generationMode: "pro"`,
`stepsOverview`), `aboutAudos`, `platformPreview`.
OTP expires **5 minutes** after being sent. Re-issue via `/start` (60s cooldown, 3-per-15-min limit).

### GET /status/:workspaceId (and POST /status with `{authToken}` body)
Header (if using Bearer): `Authorization: Bearer <authToken>` — see the Authentication caveat above.

Key fields: `landingPageReady`, `status` (`running`/`complete`/`failed`), `progress` (0-100),
`estimatedTimeRemaining`, `completedSteps`, `parallelBuildStatus`, `currentStep`/`totalSteps` (11),
`milestones` (keyed by `brandReady`, `landingPageDraftReady`, `landingPageReady`, `workspaceUsable`,
`adsReady`, `launchReady`, `complete`, `failed`, `recoverable`), `whatsHappening`, `runId`, `eventCursor`.
On completion: `nextSteps`, `trialPeriod`, `adCredits`, `platformCapabilities`, `publishingHouse`.

### POST /chat (and /chat/:workspaceId)
Body: `{ authToken, message }` (prefer this over the Bearer-header form).
Returns: `workspaceId`, `chatId`, `response`. See `references/auth-and-chatid.md` before treating
`chatId` as a named-thread handle — it isn't one.

### POST /rebuild/:workspaceId
Header: `Authorization: Bearer <authToken>`. Retry a failed workspace build.

## Build Stages & Execution

The onboarding build is **11 stages** (not 7), ~6–10 min. Stages 1–4 are sequential; 5–10 run in
parallel after branding; 11 finalizes:

1. Identify Customer · 2. Identify Problems · 3. Design AI Tool Suite · 4. Design Brand Identity ·
5. Create Hero Video · 6. Style the Space · 7. Create the Space · 8. Build Landing Page · 9. Set Up Ads ·
10. Prepare Ad Launch · 11. Go Live

> **Execution backend:** landing-page and app/space build/edit jobs are delegated to **Cursor
> Background Agents** (Otto labels the executor "App agent / Cursor"). Build throughput and success
> depend on **Audos's own Cursor account's** usage/billing headroom, not just the Audos subscription.
> When that account is over its hard limit, jobs **instant-fail at the delegation step** with
> `usage_limit_exceeded`. This is account-level and usually transient — retry once headroom is
> restored. Chat (`/chat`) is **not** gated by this and keeps working during a build outage, so use it
> to probe/report status even when builds are blocked. The failure message links to `cursor.com` but the
> account is almost certainly Audos's, not the workspace owner's — confirm via Otto/Audos support rather
> than assuming a self-serve fix. Full detail: `references/build-execution-and-failures.md`.

## Commercials (observed on completion)

The status `nextSteps` payload includes **$50 ad credits**, plus a free trial whose length is stated
**inconsistently in Audos's own copy** — `trialPeriod.duration` says "30 days" while `afterTrial` says
"after 7 days, $20/workspace/mo (capped $100/mo), auto-charged via the Portfolio Wallet." Verify the
real trial window and auto-charge terms before relying on them or repeating them to a user.

## Error Codes

`AUTH_TOKEN_INVALID` (401), `WORKSPACE_NOT_FOUND` (404), `OTP_EXPIRED` (401, code older than 5 min),
`OTP_INVALID` (401), `OTP_MAX_ATTEMPTS` (429), `RATE_LIMITED` (429, OTP send cooldown),
`VALIDATION_ERROR` (400), `CHAT_FAILED` (502), and — surfaced through Otto/task status rather than the
HTTP layer — `usage_limit_exceeded` (see Build Stages & Execution above).

## Notes

- **Never commit/print/log live `aud_live_` tokens** — redact in docs and support tickets. See
  `../../OTTO_API_WORKFLOW.md`'s security boundary.
- Tokens are workspace-scoped. One email can own multiple workspaces, each with its own token.
- Running queued Otto work can execute **all** pending drafts, not just the intended one — instruct
  Otto to run a specific draft **by its task ID**. See `../../OTTO_API_WORKFLOW.md` for the operational
  guardrails around this.
