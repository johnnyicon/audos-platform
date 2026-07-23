---
date: 2026-03-31
area: web-api
status: open
filed: no
label: "`isJsRendered` field missing from `fetch` response"
---

The `isJsRendered` boolean field is undefined on the `fetch` action's response, even though the
associated warning text fires correctly when a page is JS-rendered — the underlying detection works, the
flag surfacing it to the caller doesn't.
