---
date: 2026-03-31
area: web-api
status: fixed
filed: no
label: "`response.headers.get is not a function`"
---

Calling `.get` on a `fetch` response's `headers` threw — the returned object wasn't a real `Headers`
instance in this execution environment.

> Fixed by reading the full response as text instead of relying on the `Headers` API. Resolved 2026-03-31,
> same day it was raised.
