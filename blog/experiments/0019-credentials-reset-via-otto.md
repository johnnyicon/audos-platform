---
date: 2026-07-16
area: db-api
status: confirmed
label: "Can Otto reset or view stuck direct-Postgres credentials from chat?"
---

**Hypothesis:** The Developer panel's "Generate Credentials" button was stuck in a permanent 409
("Credentials already exist. Use regenerate to rotate them.") with no UI control to act on it. Does Otto
have a chat-accessible tool to view or rotate these credentials, as a way around the missing UI control?

**Method:** Asked Otto directly, in a fresh chat inside the affected workspace, whether it had any way to
view the existing connection string or trigger a rotation. Otto searched its own API docs before
answering rather than guessing.

**Result: confirmed no, on both counts — and this independently corroborated the UI-side finding.** Otto
has no tool to manage this specific scoped Postgres role; its own database tools (`db_query`,
`db_list_tables`) run through a separate, short-lived workspace token unrelated to the Developer-panel
credential. More importantly, Otto confirmed there is genuinely no "view" path for anyone — the
connection string is only ever shown once at generation time and isn't stored in plaintext anywhere
retrievable afterward, by design. That second, independent confirmation (from the platform's own side,
not just our own UI reproduction) is what justified actually filing this with Priority Support — see
`blog/bugs/0023-db-credentials-generate-once-no-recovery.md`.

See `field-notes/ACTIVITY-LOG.md`, rows on the Otto credentials-reset chat and the subsequent filing.
