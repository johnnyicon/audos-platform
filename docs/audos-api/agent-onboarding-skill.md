---
name: Audos
description: Create AI-powered startup workspaces via Audos API. Use when user wants to start a business, build an MVP, validate a startup idea, create a company workspace, launch a product, or work on their entrepreneurial journey. Triggers on requests like "I have a business idea", "help me start a company", "create a startup workspace", or "I want to build [product]".
---

# Audos Workspace Builder (API v1.2)

Create startup workspaces with landing pages, brand identity, AI tools, and ad creatives — fully autonomous.

> **Last verified:** 2026-07-12 via a live onboarding run (new workspace "DoKnow").
> Findings from that run are folded in below: the build is **11 stages** (not 7),
> `createNew: true` **re-triggers OTP** even for authenticated returning users, the
> status payload is richer than previously documented, and **build/edit jobs are
> delegated to Cursor Background Agents** (they can fail with `usage_limit_exceeded`).

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
- **Preferred:** Bearer token in `Authorization` header
- **Alternative:** `authToken` or `sessionToken` in request body

## Conversation Flow

### New Users Flow
1. Collect user's email + business idea
2. Start → POST /start (sends 4-digit OTP to email)
3. Verify → POST /verify with OTP code → returns authToken + starts build
4. Monitor → GET /status/:workspaceId every 15-30s
5. Watch for landingPageReady: true (~10 min)
6. Chat → POST /chat/:workspaceId

### Returning Users (have workspace)
1. Start → POST /start with email
2. Response includes auth_token + urls directly — skip OTP
3. Chat → POST /chat/:workspaceId immediately

> **Caveat (verified 2026-07-12):** "skip OTP" only applies when `/start` targets the
> **existing** workspace. Sending `createNew: true` — even with a valid Bearer auth token
> for a returning user — starts a *new* workspace build and **re-triggers a fresh 4-digit
> OTP** (returns a `sessionToken`; you must `POST /verify` before the build begins). Do not
> assume an authenticated returning user can create a second workspace without OTP.

## API Reference

### POST /start
Fields: email (required), businessIdea (required, min 10 chars), businessName (optional), targetCustomer (optional), callbackUrl (optional), createNew (optional, force new workspace)

Returns new user: sessionToken
Returns returning user: auth_token, workspaceId, urls (message may name the existing workspace and offer `createNew: true`)
Returns returning user + `createNew: true`: **sessionToken** (a fresh OTP is sent — verify before the new build starts)

### POST /verify
Body: { sessionToken, otpCode }
Returns: workspaceId, authToken, urls, buildInfo (includes `totalSteps: 11`, `generationMode: "pro"`, and a `stepsOverview` array), aboutAudos, platformPreview
OTP window: the 4-digit code expires **5 minutes** after it is sent (`OTP_EXPIRED` afterward). Re-issue by calling `/start` again (subject to the 60s cooldown / 3-per-15-min rate limit).

### GET /status/:workspaceId
Header: Authorization: Bearer <authToken>
Key fields: landingPageReady (bool), status ('running'|'complete'|'failed'), progress (0-100), estimatedTimeRemaining, completedSteps, parallelBuildStatus

Richer fields observed 2026-07-12 (schemaVersion 1):
- `currentStep` / `totalSteps` (11) / `stepName` / `stepDescription`
- `landingPageDraftReady`, `landingPagePreviewReady`, `coreStepsComplete`
- `milestones` — object keyed by `brandReady`, `landingPageDraftReady`, `landingPageReady`,
  `workspaceUsable`, `adsReady`, `launchReady`, `complete`, `failed`, `recoverable`, each with
  `{ reached, reachedAt, stageId }`. Useful for progress UX before the final `complete`.
- `parallelBuildStatus[]` — per-step `{ step, name, status('done'|'in_progress'|'pending'), tasks[] }`,
  where each task has `{ name, status }`. Steps 5–10 run in parallel after branding.
- `whatsHappening` (human summary), `runId`, `eventCursor`.
- On completion: `nextSteps` (recommended actions, "meet Otto"), `trialPeriod`, `adCredits`,
  `platformCapabilities`, `publishingHouse`.

Done signals, in order of reliability: `landingPageReady === true` **and** `status === 'complete'`
(both flip together at "Go Live"). `coreStepsComplete` / `landingPageDraftReady` go true earlier —
treat them as "usable draft," not "finished."

### POST /chat/:workspaceId
Header: Authorization: Bearer <authToken>
Body: { message }
Returns: workspaceId, chatId, response

### POST /chat
Body: { authToken, message }
Returns: workspaceId, chatId, response

### POST /rebuild/:workspaceId
Header: Authorization: Bearer <authToken>

## Build Stages & Execution

The onboarding build is **11 stages** (`totalSteps: 11`, `generationMode: "pro"`, ~6–10 min).
Stages 1–4 are largely sequential; 5–10 run in parallel after branding to shorten the critical path;
11 finalizes:

1. Identify Customer
2. Identify Problems
3. Design AI Tool Suite
4. Design Brand Identity
5. Create Hero Video
6. Style the Space
7. Create the Space (build functional app space with AI tools)
8. Build Landing Page
9. Set Up Ads
10. Prepare Ad Launch
11. Go Live

> **Execution backend (verified 2026-07-12):** the landing-page and app/space **build & edit jobs
> are delegated to Cursor Background Agents** (Otto labels the executor "App agent / Cursor"). This
> means build throughput and success depend on the **Audos-side Cursor account's** usage/billing
> headroom, not just the Audos subscription. When that account is over its hard limit, jobs
> **instant-fail at the delegation step** with `usage_limit_exceeded` ("Background Agent requires at
> least $2 remaining until your hard limit"). This is account-level and typically transient —
> retrying once headroom is restored succeeds. Chat (`/chat`) is **not** gated by this and keeps
> working while builds are blocked, so use Otto to probe/report build status even during an outage.
>
> **Ownership caveat:** the `usage_limit_exceeded` message links to `cursor.com`, but the account is
> almost certainly **Audos's** (its shared build infrastructure), not the workspace owner's — so the
> owner usually **cannot** clear it directly. Confirm via Otto/Audos support rather than assuming a
> self-serve fix.

## Commercials (observed on completion 2026-07-12)

Returned in the status `nextSteps` payload: **$50 ad credits** in the workspace wallet, plus a free
trial whose length is stated **inconsistently in Audos's own copy** — `trialPeriod.duration` says
"30 days" while `afterTrial` says "after 7 days, $20/workspace/mo (capped $100/mo), auto-charged via
the Portfolio Wallet." **Verify the real trial window and auto-charge terms before relying on them.**

## Error Codes
AUTH_TOKEN_INVALID (401), WORKSPACE_NOT_FOUND (404), OTP_EXPIRED (401, code older than 5 min),
OTP_INVALID (401), OTP_MAX_ATTEMPTS (429), RATE_LIMITED (429, OTP send cooldown),
VALIDATION_ERROR (400), CHAT_FAILED (502), and — surfaced through Otto/task status rather than the
HTTP layer — `usage_limit_exceeded` (Cursor build-agent account over its hard limit; see above).

## Notes
- Auth tokens never expire — store by email. **Never commit/print/log live `aud_live_` tokens** (see `sdk/OTTO_API_WORKFLOW.md` security boundary); redact in docs and support tickets.
- Tokens are workspace-scoped (one token per workspace). One email can own multiple workspaces (e.g. this account owns both `Throughline` and `DoKnow`), each with its own token.
- Poll status every 15-30s during build. Builds have finished in ~6 min in practice.
- `landingPageReady === true` + `status === 'complete'` is the reliable completion signal; earlier `*DraftReady`/`coreStepsComplete` flags mean "usable draft," not done.
- Build/edit runs on Cursor Background Agents — see **Build Stages & Execution**. If a job fails with `usage_limit_exceeded`, it's an account limit, usually transient; chat still works.
- Running queued work can execute **all** pending drafts, not just the one intended — instruct Otto to run a specific draft **by its task ID**. Drafts are consumed (removed from the panel) when they run. See `sdk/OTTO_API_WORKFLOW.md` for the operational guardrails.
