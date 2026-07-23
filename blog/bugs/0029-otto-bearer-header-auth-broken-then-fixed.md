---
date: 2026-07-08
area: identity
status: fixed
filed: yes
filed_ref: "Priority Support, Throughline workspace, 2026-07-08, via Otto onboarding API"
audos_status: fixed
label: Otto's onboarding API rejected Bearer-header auth for weeks while body-auth worked — independently reverified after Audos marked it fixed
---

Otto's external onboarding API (`GET /api/agent/onboard/status/:workspaceId`,
`POST /api/agent/onboard/chat/:workspaceId`) documents `Authorization: Bearer <token>` as a valid auth
form. It wasn't — both endpoints returned `401 AUTH_TOKEN_INVALID` consistently, while the same
operations worked fine through the documented body-auth form (`POST /status` / `POST /chat` with
`{"authToken": "..."}` in the JSON body). Not a one-off — reproduced across multiple sessions over
several days.

Filed as a Priority Bug, specifically calling out the split: body-auth confirmed working, bearer-header
confirmed broken, asking Audos to either fix the header route or update the docs to say body-auth is the
supported form for external agents.

> **Didn't trust the "Completed" status either.** The Help panel marked both the original ticket and its
> follow-up `Completed` on 2026-07-09 — support's own note in-thread was explicit that "the help panel
> status alone did not prove the bearer-header route was fixed." Direct re-verification followed: on
> 2026-07-09, a fresh bearer-header check against the live workspace passed (`GET /status/:workspaceId`
> → 200, `POST /chat/:workspaceId` → 200 with a real `chatId` and Otto response). Retested again on
> 2026-07-10 after reopening the same ticket thread, same clean result both times. Only then treated as
> genuinely fixed.

Same discipline this whole SDK log keeps independently arriving at on the DoKnow/field-notes side: a
platform status label saying "Completed" is not evidence on its own — it's a claim, and the only thing
that closes the loop is testing the actual behavior yourself, after the fact, not before.

Source: `throughline/docs/working/audos-otto-browser-approach.md`.
