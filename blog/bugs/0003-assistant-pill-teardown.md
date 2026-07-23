---
date: 2026-07-13
area: app-build
status: fixed
filed: no
label: "\"Assistant\" pill tore down the full-screen app back to the old chat UI"
---

The small "Assistant" pill meant to keep chat reachable from a full-screen app called
`returnToAgentView()` in `Desktop.tsx` — which tore down the full-screen shell entirely and returned the
user to the **complete old three-pane chat interface**, not a small popup. Confirmed by direct
observation (a first look at a cropped screenshot missed this; a fuller screenshot from the user caught
it) after initially misreading it as harmless.

> Fixed ourselves: changed the handler to open a small bottom-right popup instead of calling the
> state-teardown function. Verified live by clicking it post-publish — the full-screen app stays mounted
> and visible underneath.
