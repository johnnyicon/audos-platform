# Otto API Workflow

Last verified: 2026-07-12 (build-execution backend + task-draft behavior; see new sections
"Build & Edit Execution Runs on Cursor" and "Task/Draft Execution Gotcha"). Prior: 2026-07-08.

This note documents the current external-agent workflow for the Audos/Otto
onboarding API. It covers conversation continuity and how to file precise Help
bugs when the API, DNS, or platform behavior does not match the documented
contract.

## Current Throughline Workspace

- Workspace UUID: `8f1ad824-832f-4af8-b77e-ab931a250625`
- Workspace slug: `workspace-351699`
- Landing page: `https://www.trythroughline.com`
- Audos workspace: `https://audos.com/workspace/8f1ad824-832f-4af8-b77e-ab931a250625`
- Current Railway app: `https://throughline-web-production.up.railway.app/`
- Desired app subdomain: `https://app.trythroughline.com/`

Do not print, commit, paste into chat, or log live `aud_live_` tokens. Treat
them as secrets even when the API returns them through `/start`.

## Throughline Domain Split Clarification

Last checked through the Otto API on 2026-07-10. The API call returned
`chatId=ae7c2e48-f751-4e43-9803-3c69cd046e67`; no auth token was printed or
stored.

Desired architecture:

- keep `trythroughline.com` and `www.trythroughline.com` hosted/managed by
  Audos for marketing, sales, lead capture, and nurture;
- point only `app.trythroughline.com` to the externally hosted Railway product
  app;
- do not move the whole domain away from Audos unless the split is not
  supportable.

Otto's read-only answer:

- The current DNS zone is managed inside Audos's DNSimple account and fronted by
  Audos Cloudflare-for-SaaS custom hostnames.
- `@`, `www`, and `app` all currently point at `fallback.audohost.com`.
- `app.trythroughline.com` is still claimed by Audos through its CNAME,
  `_cf-custom-hostname.app`, and ACME validation records.
- At the DNS-record level, splitting the architecture is technically coherent:
  apex and `www` can remain on Audos while only `app` points to Railway.
- The unresolved platform question is automation/reconciliation: Audos must
  confirm or configure that its custom-domain automation will not recreate or
  reclaim `app -> fallback.audohost.com` after the manual cutover.

Required app-subdomain cutover if Audos supports the split:

```text
1. Keep apex/@ and www records unchanged on Audos.
2. Keep marketing/sales site hosted by Audos.
3. In Railway, keep app.trythroughline.com attached to the production
   throughline-web service.
4. In Audos/DNSimple, repoint only:
   CNAME app -> <Railway-provided target>
5. Add Railway's verification TXT if Railway requires it.
6. Remove stale Audos-only records for the app hostname:
   _cf-custom-hostname.app
   _acme-challenge.app
7. Confirm Audos domain automation will not re-add app -> fallback.audohost.com.
```

Expected impact:

- Repointing `app.trythroughline.com` should not affect
  `trythroughline.com` or `www.trythroughline.com`; they are separate DNS
  records and separate SSL/custom-hostname claims.
- Because Audos controls the DNS zone, all app-subdomain changes still need to
  go through Audos unless the domain/nameservers are transferred.

Alternatives if Audos cannot guarantee the split:

- keep the product app on the Railway-provided URL;
- use a new subdomain Audos has never onboarded, such as `go.` or `product.`;
- transfer DNS authority to a domain account controlled by Throughline, then
  point apex/www back to Audos and app to Railway manually.

## API Conversation Continuity

The production-working API path is body-auth:

```text
POST /api/agent/onboard/start
POST /api/agent/onboard/status
POST /api/agent/onboard/chat
```

Use `/start` with `workspaceId` and `createNew=false` to target an existing
workspace. Read the returned token in memory, then use that token in body-auth
`/status` and `/chat` calls.

The `/chat` response includes:

```json
{
  "workspaceId": "...",
  "chatId": "...",
  "response": "..."
}
```

Observed behavior on 2026-07-08:

- Two consecutive body-auth `/chat` calls for the Throughline workspace returned
  the same `chatId`.
- Otto remembered a nonce from the first API call in the second API call.
- The public skill docs do not clearly define a `chatId`, `threadId`, or
  `conversationId` request contract.
- Otto was asked directly whether a request-side continuation field exists, but
  could not verify it from an authoritative route schema/OpenAPI contract.
- A follow-up probe sent the current `chatId` back in the body. The call
  returned the same `chatId` and preserved context.
- The same probe sent a random UUID as `chatId`. The call ignored that random
  value, returned the active `chatId`, and still remembered prior context.
- A follow-up no-id/blank/null creation probe was inconclusive: `/start`
  returned a workspace-scoped token, but `/chat` timed out before returning a
  result. Do not treat `chatId: ""` as a confirmed way to create a fresh topic
  thread.
- A later smoke test on the same date confirmed `/start` still returned a
  Throughline-scoped token, but both body-auth `POST /chat` and bearer
  `POST /chat/:workspaceId` timed out after 90 seconds on a minimal "OK only"
  prompt. A 2026-07-09 retest after the Help tickets were marked complete
  passed: `/start` returned a real `aud_live_` token in memory, body-auth
  `/status` and `/chat` returned HTTP 200, and bearer-header `/status/:workspaceId`
  and `/chat/:workspaceId` returned HTTP 200.

Practical rule: treat `chatId` as a returned continuation/correlation ID, not
as a reliable create/select-thread API. To keep iterating in the same
conversation, reuse the same workspace auth token, same workspace, and pass the
returned `chatId` back when available. To start fresh, omit `chatId` unless and
until Audos documents a stronger create-thread route. Do not invent random IDs
or assume `chatId: ""` creates a new topic thread. Store `chatId` only in
operational notes, Gomanan tasks, or support tickets as a non-secret correlation
value.

If a future Audos doc adds an explicit `chatId` or `threadId` request field,
update the SDK before relying on it.

## Build & Edit Execution Runs on Cursor

Verified 2026-07-12 during a live onboarding run (new workspace `DoKnow`).

When Otto runs a landing-page edit or an app/space build, the actual code work is
**delegated to Cursor Background Agents** — Otto surfaces the executor as
"App agent / Cursor". Practical consequences:

- **Throughput and success depend on the Audos-side Cursor account's usage/billing
  headroom, not just the Audos plan.** When that account is over its hard limit, jobs
  **instant-fail at the delegation step** with:

  ```text
  usage_limit_exceeded — "You need to increase your team hard limit.
  Background Agent requires at least $2 remaining until your hard limit.
  Manage it at https://www.cursor.com/dashboard?tab=settings"
  ```

- The failure is **account-level and typically transient.** In the same session, jobs
  that instant-failed with `usage_limit_exceeded` succeeded on retry a short time later
  (the next job got *past* delegation and ran) — i.e. the limit was restored upstream.
- **`/chat` is not gated by this.** Otto keeps answering while builds are blocked, so use
  chat to probe task status and read the exact error even during a build outage.
- **Ownership caveat:** the error links to `cursor.com`, but the account is almost
  certainly **Audos's shared build infrastructure**, not the workspace owner's — so the
  owner usually **cannot** clear it self-serve. Confirm via Otto/Audos support before
  telling a user to "go raise the Cursor limit."
- **When to file a Help bug:** if `usage_limit_exceeded` persists across retries over time,
  it's an Audos platform-capacity issue — file a Help bug (template below) with the exact
  error, task IDs, and timestamps; do not tell the user to fix a Cursor account they don't own.

**Alternate backend — Audos Code (verified 2026-07-12):** builds can also run on **Audos Code**
instead of Cursor, which sidesteps the Cursor usage limit. **However, Audos Code requires a
signed-in workspace session's user email for attribution and CANNOT be launched from the external
onboarding/chat API** — both the draft-run and direct-handoff routes fail with
"no user email for attribution — retry from a signed-in workspace session," and there is no request
field to inject the email. **Consequence for external agents:** when Cursor is over its limit, an
API-only agent has **no working build path** — it can *create* the drafts, but a human must run them
from inside the signed-in workspace (the in-browser Otto/Tasks panel). Otto can create the Audos Code
drafts and report their IDs (`wjd_…`) for the human to run.

## Task/Draft Execution Gotcha

Verified 2026-07-12. Otto stages work as **drafts/tasks** in a Tasks panel, then runs them.

- A generic "run" / `runAll` can execute **every pending draft in the panel**, not just the
  one you asked for. In one run this launched three jobs — the intended edit plus a stale
  generic landing-copy draft and an unrelated app-build draft — which would have collided on
  the same landing page (last writer wins).
- **Always instruct Otto to run a specific draft by its task ID**, e.g. "run only job #80224,
  do not run #80186 or #80185."
- **Drafts are consumed when they run** — after a run (success or fail) they leave the panel,
  so a failed batch clears the conflicting drafts as a side effect. Re-create only the draft
  you want before re-running.
- Task IDs (e.g. `#80224`) and `chatId` are non-secret correlation values; record them in
  Gomanan/operational notes for traceability.

## DoKnow Workspace (onboarding trial, 2026-07-12)

Second workspace created under `john@merkhetventures.com` while exercising the onboarding API.
Recorded for platform-behavior traceability (token stored out-of-band, **not** in this repo).

- Workspace UUID: `8a65a4ac-5a22-435f-b55f-c41ea34ca00d`
- Slug: `workspace-156396`
- Dashboard: `https://audos.com/workspace/156396`
- Landing page: `https://audos.com/site/156396`
- App space: `https://audos.com/space/workspace-156396`
- Build: 11 stages, `generationMode: "pro"`, completed ~6 min. Initial app scaffold: one app
  ("Coach Queue"); landing hero shipped with generic template copy (DoKnow-specific edit run
  separately). Do not print/commit its `aud_live_` token — treat as a secret.

## Recommended Agent Pattern

1. Call `/start` with the owner email, `workspaceId`, `businessIdea`, and
   `createNew=false`.
2. Extract the returned token in memory only.
3. Call `/status` with `{ "authToken": token }` to verify the workspace, or
   call `GET /status/:workspaceId` with `Authorization: Bearer <token>`.
4. Call `/chat` with `{ "authToken": token, "message": "..." }`, or call
   `POST /chat/:workspaceId` with `Authorization: Bearer <token>`.
5. For follow-ups on the same topic, pass the returned `chatId` back in the
   next `/chat` request.
6. Record the returned `chatId`, task IDs, and status in Gomanan.
7. For every mutating request, include explicit guardrails:
   - do not publish unless asked;
   - do not change DNS unless asked;
   - do not create a new workspace;
   - **run only the specific draft by its task ID — never a blanket "run all"** (avoids
     executing stale/conflicting drafts; see "Task/Draft Execution Gotcha");
   - return task IDs and draft/publish status;
   - state whether a task is new, resumed, or superseded;
   - if a build/edit fails with `usage_limit_exceeded`, report it as an Audos-side Cursor
     capacity limit (not the user's account) and retry rather than escalating to the user.

## Potential Feature Request: Explicit Topic Threads

Current state: the API returns `chatId`, but the public contract does not clearly
support creating, naming, listing, selecting, or isolating topic-specific chats.
This matters when agents need separate long-running workstreams such as
branding, DNS, SDK, product strategy, and support bugs.

Proposed API capability:

```text
POST /api/agent/onboard/chats
Body: { authToken, name, purpose, metadata? }
Returns: { workspaceId, chatId, name, purpose, createdAt }

GET /api/agent/onboard/chats
Body or header auth: { authToken }
Returns: [{ chatId, name, purpose, updatedAt, latestTaskIds? }]

POST /api/agent/onboard/chat
Body: { authToken, chatId, message }
Returns: { workspaceId, chatId, response }
```

Acceptance criteria:

- External agents can create a named chat such as `Throughline Branding System`.
- External agents can list existing chats for the workspace without using the
  browser.
- Passing `chatId` continues that exact chat.
- Passing an unknown `chatId` returns a clear 404/validation error instead of
  silently falling back to another chat.
- Separate chats do not share conversation history unless Otto explicitly
  retrieves workspace-level facts.
- Responses identify which workspace-level tools/history were read.
- The browser UI and API expose the same `chatId` for the same chat thread.
- No bearer tokens or private workspace data are exposed in support transcripts.

## Help Bug Report Template

Use Help for platform/API/DNS behavior that Audos must fix. File bugs with
specific verification criteria so the support AI and QA team can tell whether
the issue is actually resolved.

```text
Title:
<short action + failing surface>

Workspace:
- UUID: <workspace uuid>
- Slug: <workspace slug>
- Owner email: <email>
- Relevant URL: <workspace/site/app/support url>

Intent:
I am trying to <goal>, so that <user/business outcome>.

Expected behavior:
<specific documented or required behavior>

Actual behavior:
<specific failure, status code, wrong response, missing record, stuck task, etc.>

Reproduction steps:
1. <step>
2. <step>
3. <step>

Evidence:
- Endpoint(s): <method + path only, no secrets>
- Request shape: <redacted JSON shape>
- Response/status: <redacted status/body excerpt>
- Task ID / chatId / support thread ID: <ids if available>
- Timestamp/time zone: <when tested>

Impact and priority:
<what is blocked, who is affected, and why this should be priority>

Security boundary:
Do not send bearer tokens through chat or support transcripts. If a token is
needed, provide it only through an approved secure path.

Success criteria:
- <criterion 1, externally verifiable>
- <criterion 2, externally verifiable>
- <criterion 3, externally verifiable>

Requested outcome:
Please either <fix option A> or <document/confirm option B>. If this is expected
behavior, please update the public docs and confirm the supported integration
path.
```

## Example: Bearer Header Auth Bug

```text
Title:
Otto external API bearer-header auth fails for existing Throughline workspace

Workspace:
- UUID: 8f1ad824-832f-4af8-b77e-ab931a250625
- Slug: workspace-351699
- Owner email: john@merkhetventures.com

Intent:
I am trying to let an external coding agent use the documented Otto API for the
existing Throughline workspace without using the browser UI.

Expected behavior:
The documented bearer-header endpoints should authenticate with the scoped
workspace token:
- GET /api/agent/onboard/status/:workspaceId
- POST /api/agent/onboard/chat/:workspaceId

Actual behavior:
The body-auth endpoints work, but the bearer-header endpoints return 401.

Reproduction steps:
1. POST /api/agent/onboard/start with the owner email, existing workspace UUID,
   and createNew=false.
2. Confirm the response targets workspace-351699 and returns a token.
3. POST /api/agent/onboard/status with { authToken } succeeds.
4. POST /api/agent/onboard/chat with { authToken, message } succeeds.
5. GET /api/agent/onboard/status/:workspaceId with Authorization: Bearer ...
   returns 401.
6. POST /api/agent/onboard/chat/:workspaceId with Authorization: Bearer ...
   returns 401.

Security boundary:
The token is intentionally omitted from this report. Please provide a secure
channel if you need a live token for debugging.

Success criteria:
- Bearer-header status returns HTTP 200 for the existing Throughline workspace.
- Bearer-header chat returns HTTP 200 and an Otto response for that workspace.
- The public API docs clearly identify whether bearer-header auth or body-auth
  is the supported production path.
- No token is printed into chat, support logs, or public docs.
```

## Example: DNS Cutover Bug

```text
Title:
Apply DNS cutover for app.trythroughline.com to Railway

Workspace:
- UUID: 8f1ad824-832f-4af8-b77e-ab931a250625
- Slug: workspace-351699
- Domain: app.trythroughline.com

Intent:
I am moving the production app subdomain from Audos-hosted legacy app behavior
to the Railway-hosted Throughline app.

Expected behavior:
The Audos-managed DNS zone should point app.trythroughline.com at Railway and
release the old Audos custom-hostname claim.

Actual behavior:
Public DNS still points app.trythroughline.com at fallback.audohost.com, and
Railway still reports the custom domain as unverified.

Required DNS changes:
- CNAME app -> <Railway CNAME target>
- TXT _railway-verify.app -> <Railway verification value>
- Remove old _cf-custom-hostname.app TXT record after Railway value is present.
- Remove old _acme-challenge.app TXT record if it is only for the old Audos
  custom-hostname claim.
- Preserve apex, www, MX, SPF/DKIM/DMARC, and unrelated records.

Success criteria:
- Public DNS for app.trythroughline.com no longer resolves through
  fallback.audohost.com.
- Railway shows app.trythroughline.com verified/active on the production
  throughline-web service.
- https://app.trythroughline.com serves the same app as the Railway production
  frontend URL.
- No unrelated DNS records are changed.
```
