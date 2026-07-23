---
date: 2026-07-10
area: app-build
status: open
filed: no
label: Audos support chat flattens structured technical replies (line breaks, code blocks) into a single unreadable block
---

Posting a detailed support reply — browser evidence, DNS checks, HTTP header comparisons, exact
requested CNAME/TXT changes, formatted with line breaks and code-like blocks for readability — got
flattened by the Audos support chat UI into a single hard-to-read block, with all the structure lost.
Technical evidence that depends on being read as distinct, separated items (a list of DNS records, a
before/after header diff) becomes much harder for a human reviewer to parse once collapsed.

Worked around by asking, in a short follow-up message, for an email address or alternate channel to
send the same evidence properly formatted. The underlying request — file attachments and a structured
ticket API for support tickets, so technical reports don't have to survive being pasted into a plain
chat box at all — was filed as its own feature request; see the matching entry.

Source: `throughline/docs/working/audos-otto-browser-approach.md`.
