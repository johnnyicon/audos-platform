# Wishlist: a power-user / developer channel for Audos

*Written 2026-07-13, after two days of hands-on building (Throughline, DoKnow). This is a synthesis,
not a complaint list — everything here is grounded in a specific, cited incident from this build.*

## The actual question

Everything Audos gets right — the chat-first Otto interface, the genesis shell, one-sentence "build me
a course on X" generation — is aimed at a **no-code founder** who wants a business without touching
code. That's a real, well-served audience, and nothing here argues for changing it.

But a technical team working the way we work today — writing our own design specs, running verified
experiments, treating job output as a claim to check rather than a fact — keeps hitting the same wall:
**the platform has exactly one interaction mode, tuned for someone who isn't us.** The question isn't
"fix Otto." It's: **should Audos keep the no-code experience exactly as it is, and open a second,
narrower channel — fewer guardrails, more direct access, aimed at teams who already know what they
want built?** We think yes, and this doc is the concrete shape of what that second channel would need.

## Priority 1 — trust the job output (highest leverage of anything here)

This is the one thing that would have saved the most real time today. Every other item below is
smaller than this one.

- **Job completion language should mean what it says.** Twice today (see blog entry 6, the empty
  course mystery) a job reported "complete," "published," and explicitly **"verified"** — and the live
  app contradicted it immediately on re-test. Not a one-off: it happened on two independent fix
  attempts, worded identically both times. **Root cause, per Otto's own analysis (doc 21): this is
  likely an authority-boundary gap, not a lying agent.** The Cursor agent reports sincerely on what it
  did inside its own sandbox; the steps that determine actual ground truth — confirmation gates,
  publish/recompile, what's actually served — are enforced Audos-side, outside the agent's visibility.
  That reframes the ask: don't just ask agents to hedge more; **surface the platform-side gate state
  back to the agent** (or to us) so a job can honestly know whether its own sandbox success survived
  the platform's own approval/publish step, instead of reporting sandbox-success as if it were
  ground-truth.
- **Feature request: an automated post-build smoke check**, even a minimal one — load the affected
  route, check for a thrown error, confirm the target DOM element exists — before a job is allowed to
  say "verified." Would have caught both false "verified" claims immediately, no human re-test needed.
- **Feature request: full, untruncated job output should be retrievable via the same channel that
  launched the job.** Today, five separate times, a job's own summary got cut off mid-sentence in the
  exact spot that mattered (the root-cause line, the fixability verdict, the "what changed" detail),
  and there was no API-accessible way to fetch the rest — only a UI panel a human has to click into.
  For an API-driven agent, that's a dead end every time it happens. *(Update: the underlying Cursor
  Cloud Agents API truncates its own stream by design with no documented full-log endpoint — see
  `docs/platform/21-cursor-backend-research.md` — so this may be a genuine upstream constraint Audos
  inherited, not a withheld feature. Doesn't remove the ask: Audos could still persist and expose the
  full run transcript server-side independent of what Cursor's own stream truncates.)*

## Priority 2 — publish/build reliability

- **Bug, filed:** deep-linking a specific app (`#app-id`) opens inside the chat shell instead of
  respecting `defaultLandingView`/`defaultLandingAppId` — see `BACKLOG.md #8`. Partially fixed
  ourselves in `Desktop.tsx`; confirmed working for one app (`course-builder`), confirmed **not**
  working for another (`doknow-mockup-test`) with no explanation yet found. This inconsistency itself
  is the strongest evidence something is genuinely broken, not just under-documented.
- **Feature request: expose real build/publish status, not just "job complete."** "Published: yes" has
  meant, on different days, (a) actually live, (b) config-written but the bundle not yet recompiled,
  and (c) blocked because a second job was running against the same app scope and silently held the
  publish. A `GET` endpoint that returns the currently-served bundle hash/timestamp — something an
  external agent can poll without a browser — would remove the need to manually cache-bust-and-eyeball
  every single change, which is what we had to do all day.
- **Feature request: auto-bust cache on every publish**, or clearly document when a manual `?_cb=`
  is required. The platform already generates these params under some conditions (we saw
  `?_cb=<epoch-ms>&cdn=fallback` appear on its own) — make that automatic and consistent instead of
  something we reverse-engineered by accident.
- **Feature request: real concurrency control on same-app-scope jobs.** Two jobs targeting the same app
  didn't run in parallel despite both reporting "started" — one silently held the other's publish with
  no queue-position or wait-time signal.

## Priority 3 — the actual developer channel

- **Feature request: a genuine "import this exact code" mode**, not a prompt. Getting our own
  hand-built UI to render faithfully (blog entry 7) required pasting a 59KB file into a chat message
  with the instruction "port this verbatim, do not redesign, restyle, or improve anything" — stated
  twice, because it's an LLM's judgment call whether to comply, not a mechanical guarantee. It worked,
  but it worked because we got lucky with agent compliance, not because there's a real feature for it.
  A literal file-upload → deterministic-wrap → publish path (no LLM interpretation of layout/content,
  only the boilerplate React wrapping that's genuinely mechanical) would make this a guarantee instead
  of a bet.
- **Feature request: document the `cdnDependencies` accepted list.** We had to build a disposable
  throwaway app that tried five libraries live just to find out react-query, GSAP, Radix, and three.js
  work and react-three-fiber is flaky — because neither we nor Otto itself could read this from any
  documented source (`search_platform_code` returns "no searchable directories" even for Otto).
- **Feature request: let API-driven agents use the Audos Code backend**, not just Cursor. Right now,
  when Cursor is over Audos's own shared usage limit, an API-only agent has *no* working build path —
  Audos Code exists specifically to sidestep that, but it requires a signed-in browser session and
  can't be triggered externally. That's the one case today where we were fully blocked with no
  workaround at all except waiting.
- **Feature request: let a workspace opt into React 19** (or whatever current), real npm resolution,
  and — for at least a "single full-screen app" mode — real browser URL ownership instead of hash-only
  in-memory routing. None of these need to be the default; they'd need to be an explicit, documented
  opt-in for workspaces that want it.
- **Feature request: reliable platform code read access.** `search_platform_code` and
  `read_platform_file` were unavailable or path-denied on our own workspace's files (`Desktop.tsx`,
  `config.json`) more often than not, forcing every "just read the code" question into a full Cursor
  job dispatch — slow and expensive for what should be a cheap read.

## What to actually file

**As bugs (something is broken relative to the platform's own documented/intended behavior):**
1. `#app-id` deep-link routing inconsistency — fixed for one app, not another, same code path
   (`BACKLOG.md #8`).
2. "Published: yes" not reflecting the live served bundle when a concurrent job holds the same scope
   (`BACKLOG.md #7`).

**As feature requests (things that don't exist yet, needed for a credible developer channel):**
3. Automated build verification before a job may say "verified."
4. Untruncated job output retrievable via API/chat, not only a UI panel.
5. A build/publish status endpoint (served bundle hash/timestamp).
6. Automatic cache-busting on publish.
7. A genuine, mechanical "import exact code" path.
8. Documented `cdnDependencies` allow-list.
9. Audos Code reachable from the external API, not signed-in-session-only.
10. Opt-in React version / real npm / real URL ownership for workspaces that want it.
11. Reliable, workspace-scoped platform code read access outside a full job dispatch.

None of this requires touching the no-code experience most Audos customers use. It's a parallel, opt-in
lane — which is exactly the ask: keep Otto as it is, and give technical teams a shorter, more honest
path to the same platform.
