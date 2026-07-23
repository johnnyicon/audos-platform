---
date: 2026-07-17
area: identity
status: open
filed: no
label: "In the \"unified space\" model, Audos's own EmailGate unconditionally wraps any custom app-level auth — with undocumented CRM/ad-pixel side effects"
---

We built a workspace with our own password gate (a `site_access` table, username/password, custom
sign-in screen) specifically so *we* controlled who could see the content — not Audos. A real visitor
hitting the live URL cold didn't see our gate at all. They saw Audos's own email-capture prompt first.

Asked Otto to read the actual served code rather than guess, and got a precise answer: on this platform,
the "landing page" a signed-out visitor sees for any workspace **is not a page you own** — it's the
platform's own `components/EmailGate.tsx`, described as "the signed-out view of the unified space." It is
unconditionally the entry point. Nothing else in the workspace — including a fully custom, independently
built sign-in screen — mounts until this component reports its internal `step === 'complete'`. A custom
app-level gate, however it's built, can only ever be the *second* gate, reachable after Audos's own first.

> Confirmed flow for a cold, never-visited third party, straight from the served source: `loading` →
> `email` (Audos's own email-capture form — submitting it **registers a CRM contact, fires Meta/Reddit ad
> "Lead" pixels, and stores marketing attribution**, because this step is Audos's own lead-capture
> funnel, repurposed as the front door for every workspace) → `code` (only if OTP is enabled for the
> workspace) → `complete`, at which point the signed-in space content finally mounts and our custom gate
> can render.

Two real problems, not one: (1) a workspace owner cannot make their own custom auth the actual front
door without doing a full rebuild of `EmailGate.tsx` itself — there's no config flag, no "use my own
gate" toggle, nothing in the Developer panel; Otto's own words: "you can't just delete it." (2) The
default behavior of that mandatory front door has real, undisclosed side effects — a private internal
tool we built specifically to avoid depending on any platform's own sharing/identity model was, without
our knowledge, registering visitor emails as CRM contacts and firing third-party ad-tracking pixels on
every sign-in attempt, simply because EmailGate is the platform's own lead-capture mechanism wearing the
costume of a generic auth gate.

Fix in progress on our side (rebuild `EmailGate.tsx` to be our own invite-link gate instead), but the
underlying platform behavior — no way to opt a workspace out of EmailGate's CRM/pixel side effects
short of overwriting the component entirely — is the actual bug. See `experiments/0020` for how this was
discovered and confirmed.
