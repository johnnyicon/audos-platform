---
date: 2026-03-31
area: web-api
status: open
filed: no
label: "`contentLength` from `fetch` on JS-rendered SPAs returns just the title length"
---

On JS-rendered single-page apps, `fetch`'s `contentLength` returns only the length of the page title, not
real content — because the raw HTML hasn't executed the client-side render yet. Workaround: use
`analyze`/metadata endpoints instead of raw `fetch` for SPA targets.

> Partially resolved by platform changelog `20260331-1119`, but the related `isJsRendered` flag (bug
> 0010) is still broken, so the underlying detection gap isn't fully closed.
