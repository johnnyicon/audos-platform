---
date: 2026-07-08
priority: 2
status: filed
label: Otto's chatId is a correlation ID, not a real named-thread conversation API
---

Confirmed empirically, not just assumed: two consecutive body-auth `/chat` calls returned the same
`chatId` and Otto correctly remembered a nonce passed in the first call — real server-side conversation
state exists. But sending back a made-up, never-seen `chatId` doesn't create or select a new thread; Otto
just returns whatever the currently-active chat is and keeps its existing context. There's no documented
way to create a named conversation, list existing ones, resume a specific one by ID, or get a hard
rejection when passing an unknown ID — `chatId` behaves as a correlation value, not an addressable
thread handle.

Filed as a feature request: create/list named chats, continue a specific one by `chatId`, reject unknown
IDs instead of silently falling back to the active chat, and expose the same thread IDs consistently
across the browser UI and the API — so a session started in one surface can actually be resumed from the
other.

Source: `throughline/docs/decision-journals/2026-07-08-auto-api-website-execution-workstream.md`.
