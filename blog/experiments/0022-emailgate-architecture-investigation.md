---
date: 2026-07-17
area: identity
status: confirmed
label: "Why did a real visitor see Audos's own sign-in screen instead of ours?"
---

**Hypothesis:** A real visitor to the live field-notes site hit an unexpected email prompt instead of
the custom `site_access` username/password gate we built. Is Audos's own EmailGate sitting in front of
our custom gate, and if so, in what order, with what consequences?

**Method:** Rather than guess, asked Otto to read the actual served source — the real landing page
content and the EmailGate component's mount/step logic — and report the literal code, not an inference.

**Result: confirmed, with a consequence we hadn't anticipated.** The workspace runs on what Otto called
the "unified space model," where the signed-out landing view for *any* workspace is unconditionally the
platform's own `components/EmailGate.tsx` — described in its own step machine as `loading → email →
code → complete`, and nothing else in the workspace mounts until `step === 'complete'`. Our custom
`site_access` gate lives inside the signed-in content, meaning it structurally can only ever be the
*second* gate. Confirmed the email-capture step is not cosmetic: submitting it registers a CRM contact
and fires Meta/Reddit ad "Lead" pixels, because it's Audos's own lead-capture funnel wearing the shape of
a generic auth screen. Otto could not find a config flag or Developer-panel toggle to disable or reorder
it — the only paths are rebuilding `EmailGate.tsx` itself, enabling a partial "continue as guest" mode
that still shows the Audos screen first, or disabling OTP specifically (cosmetic, doesn't remove the
email step).

**Verdict:** confirmed and consequential. A workspace owner who builds their own auth for a private tool
does not actually control the front door by default — Audos's own identity/lead-capture layer does,
silently, with real third-party data-sharing side effects. See `bugs/0024` for the filed finding.
