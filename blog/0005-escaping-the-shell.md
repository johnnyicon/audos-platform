---
date: 2026-07-12
product: DoKnow
status: pass
label: Escaping the shell
---

# Escaping the shell: can Audos build a real product dashboard?

Before writing off the chat-first shell as a hard platform ceiling, we ran two clean, disposable
experiments — built, checked, and safe to throw away.

**Experiment 1: the library probe.** We asked whether our team's real stack — TanStack Query, GSAP,
Radix, three.js, react-three-fiber — could even run inside an Audos app, given the build pipeline
compiles a single file per app via server-side ESBuild and resolves dependencies through a CDN, not
`npm install`. We had a throwaway "Lib Probe" app try to import and render each one, live, with a
pass/fail badge and the real error text on failure.

Result: four of five passed cleanly — react-query, GSAP, Radix, and raw three.js all rendered.
React-three-fiber was flaky on a cold start (it hit a watchdog timeout once, then rendered fine on
reload) but does work. So the "no arbitrary npm" constraint is real, but it's narrower than it
sounds — a useful chunk of a modern frontend stack is reachable through the CDN dependency path.

**Experiment 2: escape the shell.** We built a single full-screen "DoKnow Home" app — its own top bar,
its own left-rail navigation, no chat, no dock — and set the workspace config to
`desktop.layout.defaultLandingView: "app"` with a `defaultLandingAppId`. Then we published and loaded
the base workspace URL cold.

> It worked. The base URL now boots straight into a real dashboard: a Coach greeting, one clear
> "up next" lesson card, a shelf of courses with progress bars, an "add content" flow — the actual
> information architecture DoKnow needs, with zero trace of the chat shell. This was the most important
> result of the day: the genesis shell is a *default*, not a wall. You can build past it.

One real caveat we want to flag honestly, because it matters for anyone trying this themselves:
per-app deep links (`#app-id`) still open *inside* the chat shell as a side panel — only the base URL
without a hash honors the full-screen landing app. So the escape works for the primary entry point,
not yet for arbitrary navigation. We've logged that as an open item for Audos (`BACKLOG.md #8`).

Taken together: Audos can run real libraries and it can be pushed into a genuinely single-app product
surface. The shell was never really the ceiling — the next thing we found was.
