# Audos Capabilities & Limitations — Lessons Learned

> **⚠️ Audos is actively developed. These limits move.** Every row below is a point-in-time
> observation, not a permanent truth. Audos ships fast and the team is expanding the platform, so
> **re-verify before relying on any limitation** — it may have been removed, raised, or changed.
> Each entry carries a *verified* date. **Last full verification: 2026-07-12.**
>
> This doc is the running knowledge base of *how to use Audos and where its ceilings are*. When you
> re-test and a limit has changed, update the row + the date and note what changed.

## How to re-verify (don't trust this doc blind)

Two cheap probes resolve almost everything here empirically:

1. **The lib test.** Have Otto build a throwaway app that `import`s each candidate library via the
   `cdnDependencies` path and renders a minimal proof-of-life, showing PASS/FAIL + error text per lib.
   Audos uses single-file + server ESBuild + importmap/CDN (no `npm install`), so support is per-lib
   and can only be confirmed by trying. Delete the probe after.
2. **The build probe.** Ask Otto to attempt the specific thing you're unsure about (a full-screen app,
   a shell edit, a new dependency) and report exactly what compiled/rendered vs what the platform
   blocked, with error text.

Otto can inspect the served workspace state but **cannot read the platform/compiler source**
(`search_platform_code` → "no searchable directories"; `read_platform_file` denies compiler paths),
so the exact `cdnDependencies` accept-list is not enumerable from chat — only empirically.

## The app build pipeline (verified 2026-07-12)

- **Single file per app** (`apps/<Name>/App.tsx`), compiled by **server-side ESBuild**, dependencies
  resolved via **importmap / CDN** — **not** `npm install`. No local `tailwind.config`/PostCSS, no CLI
  codegen step.
- **React 18** (the shared host instance). *(House template elsewhere uses React 19 — mismatch.)*
- **Default dependency allow-list:** `react`, `react-dom`, `lucide-react`, **Tailwind CSS v3.4.1**
  (loaded via CDN — this is the default app styling; `DesktopThemeTokens` theme the *shell* only).
- **Pre-registered extras:** **Rive** is the paved animation path.
- **`cdnDependencies` path:** additional ESM / React-18-compatible libs *may* be addable; the accepted
  set is not documented and must be probed (see the lib test).
- **No arbitrary npm packages**, no Vite, no pnpm, no bundler config, no Vitest/Playwright test stack.

## Capability matrix (verified 2026-07-12 — RE-CHECK)

| Concern | Status | Notes | Verified |
|---|---|---|---|
| Tailwind CSS | ✅ Usable (default) | v3.4.1 via CDN; default app styling | 2026-07-12 |
| shadcn/ui (CLI + `components.json`) | ❌ Not usable | No `npm install`/CLI codegen. Hand-authored shadcn-style TSX+Tailwind works instead | 2026-07-12 |
| Radix primitives | ✅ **PASS** (verified) | `@radix-ui/react-dialog@1.1.4` opened + portal-mounted; via `cdnDependencies` (esm.sh, deps pinned react@18.3.1) | 2026-07-12 |
| TanStack Router | ⚠️ In-memory only | Space runtime owns real URL/history/mount; app can't own browser URLs. (Not a render probe — a routing boundary) | 2026-07-12 |
| TanStack Query | ✅ **PASS** (verified) | `useQuery` resolved + returned its value; `react-query@5.62.7` via `cdnDependencies` | 2026-07-12 |
| GSAP | ✅ **PASS** (verified) | tween ticking + box animating; runtime v3.12.5. (Paved path is still Rive) | 2026-07-12 |
| three.js (raw) | ✅ **PASS** (verified) | WebGLRenderer created, rotating cube rendered; r170 | 2026-07-12 |
| @react-three/fiber | ⚠️ **Flaky PASS** (verified) | Renders (spinning box, 10 `useFrame` ticks) on a warm load, but a cold first run exceeded the probe's 15s watchdog and read as FAIL. It works — just slow/inconsistent to init. Prefer **raw three.js** where reliability matters | 2026-07-12 |
| Arbitrary npm / Vite / pnpm / React 19 / Vitest | ❌ Not usable | Pipeline is single-file + ESBuild + CDN, React 18.3.1 | 2026-07-12 |

*(Rows above marked "verified" were empirically confirmed 2026-07-12 by building + publishing a throwaway **Lib Probe** app and reading its live PASS/FAIL badges. `cdnDependencies` resolves ESM libs via esm.sh with `deps`/`external` pinned to the host react@18.3.1. Compile-success alone is NOT render-success. Note: r3f DOES render (verified on reload), but its
cold-start init is flaky and tripped a 15s watchdog on the first run — treat it as works-but-flaky, not
a hard fail. So all five libraries actually render; r3f is just slow/inconsistent to initialize.)*

## Structural limitations (verified 2026-07-12)

- **Chat-first shell / genesis template.** Every workspace clones the same genesis space; the default
  signed-in surface is a **ChatGPT-style chat** (conversation list + "Ask me anything" composer) with
  features opening as **right-docked side-panel apps** the chatbot points you to. There is no default
  home/shelf/dashboard. See `18-genesis-space-and-ui-ceiling.md` for the clone model and the
  fixed-vs-flexible breakdown. **A single full-screen, non-chat app is CONFIRMED possible (verified
  in-browser 2026-07-12).** Setting `desktop.layout.defaultLandingView: "app"` +
  `defaultLandingAppId: "<app>"` makes the **base space URL boot directly into one full-screen app** —
  no chat, no conversation list, no dock (verified: a "DoKnow Home" app rendered as a full dashboard
  with its own top bar + left rail). **Caveat:** deep-linking a specific app hash (`#app-id`) still
  opens it *inside* the chat shell as a right-docked side panel — only the **base URL** honors the
  default-landing-app. So the escape works for the primary entry point, not for per-app deep links.
- **Builds complete as drafts; you must publish to go live.** A "complete" build registers the app in
  `config.json` but does **not** recompile the live space — the app is absent from the dock and its
  deep link 404s until you **publish/recompile** (a recompile, not a rebuild). Don't mistake a
  completed-but-unpublished build for a broken one.
- **Build backends & who can launch them.** Build/edit jobs run on **Cursor Background Agents**
  (default) — gated by *Audos's* Cursor account usage limit (`usage_limit_exceeded`, transient) — or
  on **Audos Code**, which **requires a signed-in workspace session** and so **cannot be launched from
  the external onboarding/chat API**. Net: an API-only agent is limited to Cursor; Audos Code needs a
  human session (or a browser-driven session). See `docs/audos-api/agent-onboarding-skill.md` and
  `sdk/OTTO_API_WORKFLOW.md`.

## "Verify live" has a blind spot too: check the initial paint, not just the settled state (2026-07-13)

A correction to our own method, not just a platform finding. After "fixing" and independently verifying
(by clicking) that a deep-linked app rendered full-screen with no old shell, the user caught something
our checks had missed: **on a fresh load, the old shell (dock/sidebar) renders visibly for a brief
moment before disappearing** — a flash-of-old-shell, not something present in the settled state. Every
one of our "verify live" checks that day waited a few seconds before screenshotting, which is exactly
long enough to miss this. Reproduced directly: navigating and screenshotting with near-zero delay shows
the old dock/sidebar rendered alongside the full-screen content; waiting the usual few seconds shows a
clean page with no trace of it.

**Rule, refined: "verify live" must include checking the initial paint, not only the settled state.**
Screenshot immediately after navigation (no artificial wait) at least once per verification, in addition
to the settled-state check — a fix that only corrects the settled DOM after a flash is not the same as
a fix that renders correctly from the start, and a normal "wait then check" pass will not catch the
difference.

**Root cause, confirmed by Otto directly (2026-07-13):** the chat shell isn't a separate page — it's
`Desktop.tsx`'s **default state**; full-screen is just a different state of the same component. The
deep-link logic that switches to full-screen runs inside a `useEffect` on mount (~lines 556–607), and
`useEffect` fires **after** React's first render/paint. So the actual sequence is: (1) component mounts
with default state → shell renders → React paints it (the flash); (2) the effect runs → flips state to
full-screen → re-renders → shell disappears. Classic effect-driven-correction FOUC, not a platform
mystery.

**The fix, per Otto:** compute the initial state **synchronously, before first render** — read
`window.location.hash` inside a lazy `useState` initializer (`useState(() => computeFromHash())`)
instead of defaulting to shell-mode and correcting via an effect. No effect-driven correction, no flash.
The one caveat (config/apps availability at first render) turned out not to apply — config was already
passed synchronously as a prop, so no gating was needed.

**Fixed and independently verified 2026-07-13 (job #81755), using the corrected zero-delay method:**
navigated and screenshotted with **no wait at all** — the exact method that caught the bug — across two
different apps (`#model-probe`, `#course-builder`), each checked twice. Clean full-screen from the very
first frame both times, no sidebar, no flash. This is the one item on this page confirmed by the
strictest version of our own verification standard, not just a settled-state check.

## Build-agent self-reports are unreliable — verify live, every time (verified 2026-07-12)

**The single most important operating rule from this session.** Twice, a Cursor build job reported
success in terms that sounded conclusive — "complete," "published: yes," even an explicit
"**verified** — three lesson rows inserted with valid JSON arrays" — and in both cases, re-testing the
*actual live app in the browser immediately after* reproduced the **identical failure**
(`Insert failed: invalid input syntax for type json`) that the fix was supposed to resolve. This was
not a regression (no prior working feature broke) — it was a **verification gap**: the agent's
self-check did not exercise the real user-facing path (or checked something narrower than the actual
bug), so "verified" in the job report did not mean "verified in the live app."

**Rule: never trust a build/fix job's self-reported success or "verified" language. Always
independently re-test the actual live, published surface (reload the page, redo the user action) before
believing a fix landed.** This roughly doubles the cost of every fix cycle (build, then independently
re-verify) but is not optional — the alternative is shipping on a false signal.

## You can push your OWN hand-authored UI code verbatim — high fidelity confirmed (verified 2026-07-12)

**This may be the single highest-leverage finding of the session.** Rather than describing a UI in a
text brief and hoping the build agent's interpretation matches, you can hand Audos an **already-built,
self-contained HTML/CSS/JS mockup** (ours was a 903-line, ~55KB single file with inline `<style>` +
markup + vanilla JS) and instruct it to port the file **verbatim** into an app, with explicit
instructions not to redesign/restyle/rewrite anything. Result, independently verified in-browser (not
from a self-report): the live app matched our mockup with **high fidelity** — exact wordmark, exact
streak pill, exact 4-icon left rail, exact greeting copy, exact hero card (title, format chips, source
clip citation), exact course-shelf cards and progress values, exact lesson body copy — and in-app
navigation between views worked with no page reload (true SPA behavior), all rendered from OUR content,
not the platform's own generated copy.

**Required adaptation** (minimal, mechanical, not a redesign): wrap the markup in a default-exported
`function App() {...}`, convert the vanilla-JS behavior to work inside React (e.g. `useEffect`), and
inject the original CSS verbatim via `<style dangerouslySetInnerHTML>`. This is a boilerplate shim, not
a reinterpretation.

**Implication for how to work with Audos:** prefer **"design it exactly as wanted first (mockup/HTML,
by us or a coding subagent), then port verbatim"** over **"describe the UI in a prompt and hope the
build agent's design judgment matches."** This sidesteps the "every Audos app looks the same" problem
(see the genesis-space doc) entirely for apps built this way — the design fidelity is ours, not the
platform's default theme/agent judgment.

**Caveat:** this specific test ported *static, non-interactive-with-real-data* markup (sample/fixture
content, not live DB-backed data). Wiring a verbatim-ported UI to real `useWorkspaceDB` data (as the
"DoKnow Home" full-screen app required) is a separate, already-confirmed-possible step (see the
escape-the-shell section above) — combining both (verbatim design + live data) was not tested in one
app in this session, but nothing observed suggests it wouldn't compose.

## "Published" can mean config-written, not bundle-compiled (verified 2026-07-12)

An app can be **registered in the space config** (visible via `list_apps`) while the **served/compiled
bundle the dock actually renders is stale** — the app doesn't appear in the dock and its route 404s,
even though a job reported `published: yes`. Cause observed: **a second job running against the same
app scope holds/rebuilds the space**, so an earlier job's publish doesn't propagate until the later job
finishes. Practical rule: **never run two jobs against the same app concurrently** (reinforces "run by
specific task ID, no run-all" from the API doc), and treat "published: yes" as unconfirmed until you
see the app live in the dock yourself.

## Recurring pattern: json/jsonb columns get written unserialized (three occurrences, three different code paths)

Not a one-off bug — a pattern worth watching for in **any new write path** on this platform. Every time
a build agent has written to a `json`/`jsonb` column for the first time in a given code path, it has
initially sent a raw JS object/array instead of a `JSON.stringify`'d string, failing with
`Insert failed: invalid input syntax for type json`:

1. Course Builder's original lesson insert (`concepts`, `quiz` fields) — see blog entry 6 / `BACKLOG.md #5`.
2. The same bug's first follow-up attempt, same app, same root cause not fully caught the first time.
3. `doknow-app`'s `learner_progress` write (`completed_lessons` field) — a completely different app,
   built independently, same exact failure mode.

**Takeaway:** when briefing any job that writes to the database for the first time in a new code path,
explicitly instruct it to check every `json`/`jsonb` column and confirm `JSON.stringify` is applied
before the write — don't wait to discover it via a live `Insert failed` error. This is now a standing
item in any new-feature brief, not something to rediscover per app.

## Known quality issues observed (re-test — may already be fixed)

- **Course generation produced 0-lesson skeletons** (2026-07-12): generated courses had a title,
  description, one tier and one module but **no lessons**; the Lesson Player dead-ended at "This course
  has no lessons yet." Whether this is a prompt/build-agent limitation or a transient failure needs
  re-testing.

## Re-verification checklist (run periodically)

- [ ] Re-run the **lib test** and update the ⚠️ rows (react-query, GSAP, Radix, three, r3f).
- [ ] Re-probe whether a **full-screen non-chat app** / single-app mode is supported or easier now.
- [ ] Re-check the **Cursor vs Audos Code** launch constraints (is Audos Code API-launchable yet?).
- [ ] Re-test **course-generation depth** (does it produce real lessons now?).
- [ ] Confirm React version / default dependency list / whether `cdnDependencies` is documented.
- [ ] Update each row's *verified* date and note what changed.
