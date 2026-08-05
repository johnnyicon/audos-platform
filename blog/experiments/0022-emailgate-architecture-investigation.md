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

**Re-verified live, 2026-07-23 — still true, and slightly worse than first found.** Opened the real
field-notes site cold (no session) and went through the flow again rather than trusting the original
finding to still hold six days later. Same result: the custom `site_access` username/password gate never
appears at all. This time, typing in a fresh, never-used, made-up email address didn't even trigger an
OTP prompt — it walked straight into the full app with zero verification of any kind. Whatever gate
exists today is Audos's own EmailGate, still unconditional, still the only thing a visitor sees, six days
after the original finding and with no change in behavior.
