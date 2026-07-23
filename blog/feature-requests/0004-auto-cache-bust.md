---
date: 2026-07-13
priority: 2
status: not filed
label: Automatic cache-busting on every publish
---

The platform already generates `?_cb=<epoch-ms>&cdn=fallback` query params under some conditions — we
saw them appear on their own, then reverse-engineered when a manual `?_cb=` was required after a
`Desktop.tsx` publish. Make that automatic and consistent on every publish, or clearly document when a
manual cache-bust is required, instead of leaving it to be discovered by accident.
