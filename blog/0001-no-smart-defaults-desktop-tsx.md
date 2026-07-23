---
date: 2026-04-03
product: Throughline
status: fixed
label: No smart defaults
---

# No smart defaults: what Desktop.tsx actually controls

Before DoKnow, there was Throughline — still live today, still running its GTM on Audos. This is the
earliest entry in this log, pulled forward from that build, because the platform lessons from it are
still load-bearing for everything that came after.

First production deploy, and the workspace loaded a blank screen, then the wrong app. The instinct is
to assume the platform is smart about this — that it picks a sensible default app, shows the right
sidebar, routes you somewhere reasonable. It isn't, and it doesn't. Every one of those decisions is
just JSX and config sitting in one file, `Desktop.tsx`, that ships in your own workspace: which app
loads by default (`config.apps.find(app => app.id === 'throughline') || config.apps[0]`), which apps
appear in the sidebar at all (an explicit ID filter), and whether there's a sidebar in the first place.
Earlier attempts to "remove the dock" had been fighting `Desktop.tsx` re-rendering its own sidebar —
there was no separate platform chrome to fight. It was just our own code.

> The second discovery cost more time than it should have: **publishing a change is not the same as it
> being live.** Editing `Desktop.tsx` or `config.json` takes three separate steps — edit, publish, then
> a distinct **recompile** — and skipping the third means the old compiled bundle keeps serving
> regardless of what the config says. Months later, building DoKnow, we'd rediscover almost exactly
> this same gap — a job reporting "published: yes" while the live dock still served stale code. Worth
> knowing it's not a one-off glitch; it's a structural seam in how Audos separates config from compiled
> output.

The third was more mundane but easy to trip on: a placeholder app (`apps/throughline/App.tsx`, with
its own "Dashboard, Guests, Voice, Studio" navigation) had shipped as the default on the platform
before we ever wrote a line of our own. Our real apps — `apps/briefing/`, `apps/setup/` — were correctly
listed in the sidebar the whole time; the *default view* was just still pointing at someone else's
scaffold.

Net takeaway, stated the way we wrote it down at the time: the rendering hierarchy is a compiled,
non-customizable platform shell, then `Desktop.tsx` as our own fully customizable layout layer, then
individual apps underneath that. Everything below the platform shell is ours to shape — but nothing
above it will do that shaping for us.
