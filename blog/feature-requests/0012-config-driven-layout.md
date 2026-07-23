---
date: 2026-04-09
priority: 1
status: "not filed"
label: No config-driven way to declare a space's shape (single-app full-canvas vs. multi-app shell)
---

`config.json` has no field to declare the shape of a space — whether it should render as a full-canvas
single-app experience or the default multi-app shell with sidebar/dock. The only way to get full-canvas
today is imperative JSX surgery directly on `Desktop.tsx`, which means an architectural decision ends up
living inside a file that's also the single most common target for silent platform overwrites (see
`bugs/0031`, `feature-requests/0011`).

A declarative flag would remove the need for that surgery entirely — something like
`{"desktop": {"layout": "full-canvas"}}` or `{"spaceType": "single-app"}` in `config.json`, read by the
shell at render time instead of requiring a hand-edited component. This project's own earlier playbook
for eliminating the chat shell (`docs/platform/22`) hit the identical underlying gap independently —
this isn't a one-off Throughline complaint, it's a recurring shape of the same missing platform primitive.

Source: `throughline-forge/tmp/audos-feature-request-off-platform-development.md`.
