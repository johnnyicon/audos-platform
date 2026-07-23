---
date: 2026-04-09
priority: 1
status: "not filed"
label: OTP is the only Audos-native security lever, and it's only reachable via an undocumented raw API call — no OAuth/SSO at all
---

EmailGate ships with OTP (one-time-passcode) verification, but there's no UI to enable it anywhere in
the workspace settings panel — the only way is a raw, undocumented API call:

```
PUT /api/auth/otp/space/config/{workspaceId}
{ "enabled": true, "trigger": "always" }
```

A developer has to discover this from source code, make the call from a browser console using session
cookies, and hope the request/response shape stays stable across platform updates. OTP is, today, the
*only* Audos-native security enhancement available — there is no email allowlist, no invite-only mode,
and no OAuth/SSO support of any kind. Confirmed working when called directly (sessions correctly gained
a `verified: true` flag) — so the underlying capability is real, just badly surfaced. A single toggle in
the workspace settings UI ("Require email verification") would close the gap entirely.

Distinct from this project's own earlier finding that every unified-space workspace's signed-out view is
unconditionally the platform's own EmailGate (`docs/platform/26`) — that's about EmailGate being
mandatory and unavoidable; this is about the one real security lever *inside* EmailGate being invisible
unless you already know the API exists.

Source: `throughline-forge/tmp/audos-feature-request-off-platform-development.md`.
