---
date: 2026-07-13
priority: 3
status: not filed
label: Reliable, workspace-scoped platform code read access
---

`search_platform_code` and `read_platform_file` were unavailable or path-denied on our own workspace's own
files (`Desktop.tsx`, `config.json`) more often than not, forcing every "just read the code" question into
a full Cursor job dispatch — slow and expensive for what should be a cheap read.
