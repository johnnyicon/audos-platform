# Auth patterns and chatId behavior

Read this when auth is behaving unexpectedly, or before relying on `chatId` as anything more than a
correlation value.

## Auth: default to body-auth

The onboarding API documents `Authorization: Bearer <token>` as a valid auth form for
`GET /status/:workspaceId` and `POST /chat/:workspaceId`. For an extended period it wasn't reliable —
both endpoints returned `401 AUTH_TOKEN_INVALID` consistently, while the documented body-auth form
worked the whole time:

```
POST /status   { "authToken": "..." }
POST /chat     { "authToken": "...", "message": "..." }
```

This was reported and later independently confirmed fixed — but treat that as "confirmed working as of
a specific date," not a permanent guarantee. **Default to body-auth.** It has always worked. Only rely
on the bearer-header routes after verifying them live yourself, the same way this was caught the first
time — don't assume documentation is current just because it's the platform's own.

## chatId: a correlation ID, not a thread you can create or select

`POST /chat` and `POST /chat/:workspaceId` both return a `chatId`. It's tempting to treat this like a
normal chat-thread API — create a named conversation, switch between several. Confirmed empirically, it
doesn't work that way:

- Sending the same `chatId` back on a later call returns the same `chatId` and correctly preserves
  context — real server-side conversation state exists.
- Sending a **made-up, never-seen `chatId`** does not create or select a new thread. Otto just returns
  whatever chat is currently active and keeps its existing context.

There is no documented way to create a named conversation, list existing ones, resume a specific one on
demand, or get a hard rejection when passing an unknown ID. Treat `chatId` purely as a value you echo
back to continue an existing conversation — never as a handle you invent, or rely on to isolate one
conversation from another.
