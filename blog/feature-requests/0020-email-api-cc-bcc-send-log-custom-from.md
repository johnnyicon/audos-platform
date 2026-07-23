---
date: 2026-04-05
priority: 3
status: "not filed"
label: Audos's email-api has no CC/BCC, no send log, and no custom From address
---

Confirmed against Audos's own SDK documentation: `email-api` supports HTML bodies and `replyTo`, but has
no CC or BCC support, no send log to check what actually went out, and no way to set a custom From
address. For anything shaped like guest outreach — sending on behalf of a named show, cc'ing a producer,
confirming a batch of invites actually sent — these are real, load-bearing gaps rather than nice-to-haves.

Source: `throughline-forge/tmp/handoff-nicholas-audos-sdk.md`.
