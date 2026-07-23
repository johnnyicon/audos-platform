---
date: 2026-07-16
area: db-api
status: open
filed: yes
filed_ref: "Priority Support, field-notes workspace, 2026-07-16, via Otto chat"
audos_status: pending
label: "Database Access \"Generate Credentials\" is one-shot with no view/rotate path in the UI"
---

The workspace Developer panel's Database Access section offers direct scoped PostgreSQL access (schema
owner role, full DDL inside the workspace's own schema) via a single "Generate Credentials" button. The
UI text implies this is a normal, repeatable action ("Generate a scoped PostgreSQL role...").

In practice: credentials had already been generated once earlier in the same workspace session (the
generation click happened during an earlier detour into this panel, and the returned connection string
was never captured). Every subsequent click of the same button returns `409 Conflict` with the message
`"Credentials already exist. Use regenerate to rotate them."` — but **there is no regenerate/rotate/view
button anywhere in the panel.** A full interactive-element scan of the Developer tab found exactly one
control for this feature: `Generate Credentials`. No secondary action appears after the 409, before it,
or on hover/focus.

A `GET` to the same endpoint (`/api/workspaces/{id}/db-credentials`) returns `401`, so it's not a
"fetch what you already generated" read path either — whatever auth that endpoint expects, the normal
session cookie doesn't satisfy it.

Net effect: once credentials are generated once and the connection string isn't saved at that moment,
the workspace's direct-DB-access feature becomes permanently unusable through the UI — no way to see,
regenerate, or reset it. The API error message itself (`"Use regenerate to rotate them"`) references a
capability the frontend doesn't expose, suggesting the backend supports rotation but the UI was never
wired up to call it.

**Filed with Audos 2026-07-16, via Otto chat inside the field-notes workspace (Priority Support).** Held
back initially per our own standing rule to confirm findings with full independent certainty before
posting, after an earlier finding needed a same-day correction. Before filing, asked Otto directly whether
it had any tool to view or rotate these credentials from chat — it confirmed independently, from the
platform's own side, that it does not (its DB tools use a separate short-lived token, unrelated to the
Developer-panel role) and that there is genuinely no "view" path for anyone, since the connection string
is only ever shown once at generation and never stored in plaintext. That second, independent confirmation
— on top of our own UI reproduction — is what cleared the bar to file. Otto prepared and submitted a
Priority bug report; response confirmed "engineering team notified," with a $27 bug-bounty offer attached
if it turns out to be a real fix. A near-duplicate-match check against an unrelated existing report
(serial-id-not-uuid, `#0020`) flagged as a possible dupe but was a false positive on unrelated shared text
— submitted as a new report instead.

> **Update, 2026-07-22.** Audos's own Priority Support panel now marks this ticket **Completed**,
> with a named fix ("Database Access credentials 409 UI recovery") and a message timestamped
> 2026-07-17, 6:15 PM: "We've verified the fix and it has now been published to production. Please
> try it again." We tried it again, live, ourselves — not on their word. `POST
> /api/workspaces/{id}/db-credentials` still returns `409`. `GET` on the same endpoint still
> returns `401`. The Developer panel still shows exactly one control, `Generate Credentials`, with
> no regenerate/rotate/view option anywhere. **The claimed fix is not live.** This is the same
> failure this whole log keeps finding — a completion report that doesn't match the served
> behavior — except this time the false "Complete" came from Audos's own support process, not a
> build job. See `blog/0016` for the full account, including the two other bugs claimed fixed
> alongside this one.
