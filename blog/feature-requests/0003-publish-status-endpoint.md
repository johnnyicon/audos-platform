---
date: 2026-07-13
priority: 2
status: not filed
label: A build/publish status endpoint (served bundle hash/timestamp)
---

"Published: yes" has meant, on different days, (a) actually live, (b) config-written but the bundle not
yet recompiled, and (c) blocked because a second job was running against the same app scope and silently
held the publish. A `GET` endpoint returning the currently-served bundle hash/timestamp — something an
external agent can poll without a browser — would remove the need to manually cache-bust-and-eyeball
every single change.
