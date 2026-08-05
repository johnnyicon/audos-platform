# Otto's onboarding API: auth patterns and chatId behavior

*Written 2026-07-23, consolidating `bugs/0029` and `feature-requests/0018` into one durable reference.
Both are about the same external surface — `https://audos.com/api/agent/onboard` — and belong together
rather than scattered across separate bug/feature-request write-ups. A condensed, agent-facing copy of
this same content is bundled with the portable skill at
`sdk/skills/audos-onboarding/references/auth-and-chatid.md` — this doc is the fuller version, for a
human reading the SDK reference library.*

## Auth: use body-auth, not Bearer-header, until you've verified otherwise

The onboarding API documents `Authorization: Bearer <token>` as a valid auth form for
`GET /status/:workspaceId` and `POST /chat/:workspaceId`. For an extended period it wasn't reliable —
both endpoints returned `401 AUTH_TOKEN_INVALID` consistently, while the same operations worked through
the documented body-auth form:

```
POST /status   { "authToken": "..." }
POST /chat     { "authToken": "...", "message": "..." }
```

This was filed with Audos and independently re-verified fixed on 2026-07-10 — but treat that as "was
confirmed working as of that date," not a permanent guarantee. **Practical rule: default to body-auth.**
It has always worked. If you specifically need the bearer-header routes (e.g. for a client library that
expects standard header auth), verify them live before depending on them, the same way this was caught
the first time.

## chatId: a correlation ID, not a thread you can create or select

`POST /chat` and `POST /chat/:workspaceId` both return a `chatId` in the response. It's tempting to treat
this like a normal chat-thread API — create a named conversation, switch between several. Confirmed
empirically, it doesn't work that way:

- Sending the same `chatId` back on a later call returns the same `chatId` and correctly preserves
  context (real server-side conversation state exists).
- Sending a **made-up, never-seen `chatId`** does not create or select a new thread — Otto just returns
  whatever chat is currently active and keeps its existing context.

There is no documented way to create a named conversation, list existing ones, resume a specific one on
demand, or get a hard rejection when passing an unknown ID. Treat `chatId` purely as a correlation value
you echo back to continue an existing conversation — never as a handle you invent or rely on to isolate
one conversation from another.

## Why this is one doc, not two bug write-ups

Both findings are about the same onboarding API surface, discovered around the same
work (`throughline/docs/decision-journals/2026-07-08-auto-api-website-execution-workstream.md`), and an
agent calling this API needs both pieces of context at the same time — trusting the wrong auth form and
misusing `chatId` are the two most likely ways to get confusing, hard-to-debug behavior from this
specific endpoint set. Kept together here as the single reference to point to from the onboarding skill
file, rather than making an agent piece it together from two separately-filed items.
