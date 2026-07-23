---
date: 2026-07-13
area: app-build
status: open
filed: yes
filed_ref: docs/platform/bug-reports/2026-07-13-deep-link-fullscreen-inconsistency.md
label: Deep-link full-screen fix didn't apply consistently across apps
---

The same `Desktop.tsx` fix that made `#course-builder` deep-link full-screen (see bug 0001) **did not**
apply to `#doknow-mockup-test`, despite an identical code path and a fresh cache-bust. Three separate
diagnostic attempts were made; each was cut off by the platform's own report-truncation before reaching
the differing config field. Root cause never established from our side.

> The affected app has since been deleted (its purpose — proving the verbatim design push — was already
> fulfilled), so this couldn't be root-caused further. Filed formally with Audos, since per-app
> inconsistency on an identical code path is the strongest evidence something is genuinely broken rather
> than under-documented.
