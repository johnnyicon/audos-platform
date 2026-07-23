# The unified-space signed-out view is mandatory, not a lead-capture option you can skip

*Written 2026-07-17, after `bugs/0024`, `0027`, `0028`. Corrects the framing in
`17-emailgate-otp-configuration.md`, which was written as if EmailGate is an optional lead-capture
layer you choose to add OTP to. For workspaces on the unified-space model, it is not optional —
this doc is the mechanism `17` assumes but doesn't state.*

## The mechanism

Every workspace on the unified-space model has exactly one signed-out view, and it is not
app-specific — it is the platform's own `components/EmailGate.tsx`. This is not a widget you place
on a landing page; it *is* the landing page. Nothing else in the workspace — your own auth gate,
your app shell, your data — mounts until this component's internal state machine reports
`step === 'complete'`.

Confirmed by reading the actual served source, not inferred from behavior:

```
GateStep = 'loading' | 'email' | 'code' | 'complete'
```

- `loading` — checks for an existing `space_session_<spaceId>` in localStorage.
- `email` — Audos's own email-capture form. Submitting it **registers a CRM contact and fires
  Meta/Reddit "Lead" ad-tracking pixels** — this is the platform's lead-capture funnel, not a
  neutral sign-in screen, regardless of what your app actually is.
- `code` — only if OTP is enabled per `17-emailgate-otp-configuration.md`.
- `complete` — `completeGateEntry()` fires, EmailGate stops rendering, and *only now* does anything
  your app built — including a fully custom `site_access` username/password gate — get a chance to
  mount.

**The practical consequence:** any custom auth you build is structurally a *second* gate. It cannot
be the first thing a visitor sees, and it cannot intercept or replace the CRM/pixel side effects of
the email step, because that step already ran before your code exists.

## What we tried, and what we learned trying it

Since `EmailGate.tsx` literally is the landing page, the only lever available is rewriting that
component directly — replacing its email-capture UI with your own gate logic, calling
`completeGateEntry()` on success. We tried exactly this, twice, and both attempts revealed a second,
separate problem layered on top of the first:

1. **The rewrite doesn't reliably publish.** A job can edit the component, pass its own internal
   syntax checks, and self-report `Complete` with specific, confident detail — while the live site
   keeps serving an old published snapshot with zero trace of the change (`bugs/0027`). The server
   logs showed the workspace being "restored... from GCS," consistent with the edit landing in
   source but never reaching the compiled/published bundle. Treat any EmailGate rewrite's own
   success report as unverified until you've cold-loaded the real URL yourself.

2. **Even the un-rewritten path can silently stop routing anywhere useful.** Independent of the
   rewrite, walking the *original* email-then-password flow cold found that submitting any email
   goes straight through (no OTP) into a generic, unbranded assistant shell — not the app's own
   downstream gate or content at all (`bugs/0028`). The underlying data was intact; only the entry
   path was broken. This means you cannot assume the "old" flow is a safe fallback while a rewrite
   is in progress — verify both, separately, live.

## What we have not confirmed

- No config flag or Developer-panel toggle exists to disable, reorder, or opt out of EmailGate —
  we looked. The only lever is rewriting the component, with the caveats above.
- We have not confirmed a version of this rewrite that both (a) actually publishes and (b) routes a
  successful sign-in to the app's real content, end to end, live. Until one has been cold-verified,
  treat "we rewrote EmailGate" as unproven regardless of what any build job reports.

## Bottom line for anyone building a private, invite-only, or password-gated app on Audos

Access control for the signed-out view is platform-owned by default, not builder-owned, no matter
what auth you build inside the space. If "who can see this before they're signed in" matters for
your app, budget for this being harder than "add a password screen" — and verify every claimed fix
cold, from a real never-visited browser, before trusting it.
