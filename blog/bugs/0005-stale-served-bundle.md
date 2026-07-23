---
date: 2026-07-12
area: app-build
status: open
filed: no
label: "\"Published: yes\" doesn't guarantee the live bundle reflects the change"
---

`list_apps`/config can report an app as registered and published while the **served/compiled bundle is
stale** — the app missing from the dock, or a route 404ing. A job reporting `published: yes` doesn't
guarantee the live bundle reflects it, especially when a second job runs against the same app scope
concurrently (one can silently hold the other's publish, no queue-position signal).

> Not yet formally filed. See the related feature request for a real publish-status endpoint that would
> make this independently checkable without a browser.
