# Playbook: getting a real single-app full-screen experience on Audos

*Written 2026-07-13, distilled from a full day fixing this on the DoKnow workspace. This is the
step-by-step recipe — not a debugging narrative (see `21-cursor-backend-research.md` and `BACKLOG.md #8`
for that). Follow this in order if you're standing up a new workspace and don't want the default
chat-first genesis shell.*

**Scope note, important:** this playbook fixes **routing/shell mechanics** — whether the page is
full-screen or chat-wrapped, and whether it flashes on load. It does **not** fix **visual branding** —
whether the page looks like your actual product design vs. generic Audos-default styling. Those are two
separate problems with two separate fixes. See the last section.

## Step 1 — make the base URL boot full-screen

In `config.json`, set:
```json
"desktop": { "layout": { "defaultLandingView": "app", "defaultLandingAppId": "<your-app-id>" } }
```
This makes the bare workspace URL (no hash) load your chosen app full-screen, no chat, no dock. On its
own this only covers the base URL — the next three steps close every other route.

## Step 2 — make `#app-id` deep links full-screen too, not just the base URL

By default, `Desktop.tsx`'s deep-link effect (a `useEffect` on mount, reads `window.location.hash`)
routes any `#app-id` hash into `openPanel(matchingApp.id)` — the old chat-wrapped side-panel view,
regardless of your `defaultLandingView` setting. Fix: change that branch to run the **same full-screen
state transition** the default-landing path uses, gated on `defaultLandingView === 'app'`, instead of
`openPanel(...)`.

**Known gap, not yet resolved:** this fix did not apply consistently to every app in our workspace —
one app never picked it up despite an identical code path and a fresh cache-bust. Root cause not
established from our side; filed as a bug (`docs/platform/bug-reports/2026-07-13-deep-link-fullscreen-inconsistency.md`).
If you hit this, don't assume you did something wrong — verify per-app, and file your own report if it
recurs.

## Step 3 — fix the "return to chat" affordance so it doesn't undo everything

You'll want *some* way back to Otto/chat — don't remove it entirely. But the default button for this
(an "Assistant" pill) calls `returnToAgentView()`, which tears down the full-screen state and returns
the user to the **complete old three-pane chat interface** — chat thread, dock, and your app demoted to
a side panel. That's not a small affordance, it's a full regression. Fix: change the handler so it opens
a small, self-contained popup/drawer instead of calling the state-teardown function — the full-screen
app stays mounted and visible underneath.

## Step 4 — kill the flash-of-old-shell on initial load

Even after steps 1–3, expect a **visible flash of the old dock/sidebar on every fresh page load**,
before the corrected view takes over. This is invisible to a normal check (screenshot a few seconds
after navigating shows a clean page) — you will only catch it by screenshotting **immediately after
navigation, with no delay**.

**Root cause:** `Desktop.tsx`'s relevant `useState` calls default to shell-mode
(`isSidebarOpen = useState(true)`, `activePanelId = useState(null)`, `mobileView = useState('chat')`),
and the deep-link `useEffect` only corrects this **after** React's first paint. Every load renders shell
→ paints it → effect fires → corrects to full-screen. Classic effect-driven-correction flash-of-wrong-
content, not a platform mystery.

**Fix:** replace the plain `useState(defaultValue)` calls with **lazy initializers that compute from
`window.location.hash` (and the landing config) synchronously**, e.g.
`useState(() => resolveFromHash())`, so the correct state is already set on the very first render — no
effect, no correction pass, no flash. Confirm the apps/config list is available synchronously at that
point (as a prop, not an async fetch) — if it's async, you'll need to gate rendering rather than assume
this fix applies as-is.

## Verification checklist (don't skip any of these — each one caught something the others missed)

1. **Base URL, settled state.** Load the bare workspace URL, wait, screenshot. Confirms step 1.
2. **`#app-id`, settled state, cache-busted.** Load a deep link with a fresh `?_cb=<value>`, wait,
   screenshot. A stale cached bundle can mask an otherwise-working fix — always cache-bust after any
   `Desktop.tsx` publish, don't just reload.
3. **`#app-id`, settled state, plain reload (no cache-bust).** Confirms the fix survives normal use,
   not just a forced-fresh load.
4. **Click the "return to chat" affordance.** Confirms step 3 didn't regress — should open a small
   popup, not the full old UI.
5. **`#app-id`, zero-delay screenshot immediately after navigation, no wait at all.** The only check
   that catches step 4's flash. Do this **in addition to**, never instead of, steps 1–3 — a page that
   passes the settled-state check can still fail this one.
6. **Repeat 1–5 for more than one app.** We found real per-app inconsistency (step 2's known gap) —
   one clean pass on one app does not confirm the fix generally.

Skipping step 5 specifically is how a real, user-visible bug survived multiple rounds of "verified"
fixes on this workspace. Don't skip it.

## The separate problem this playbook does NOT solve: visual branding

Getting a route to render full-screen with no chat shell does **not** make it look like your product.
Every app fixed by this playbook (Course Builder, Lesson Player, etc. in our workspace) still renders
with **generic, unbranded Audos-default styling** — a bare white page, a minimal breadcrumb, no teal,
no left nav, no product identity — because they were built via *describe-then-generate*, which never
carried our design.

The one technique that reliably produces your actual design is the **verbatim design push** from
`docs/platform/20-power-user-wishlist.md` / blog entry 7: hand a build agent your own complete,
already-built HTML/CSS/JS and instruct it to port that file exactly, not redesign from a text brief.
That's a separate exercise per app — this playbook gets the *mechanics* right; the verbatim push gets
the *look* right. Do both if you want an app that's both full-screen and actually yours.

## Worked example: consolidating into one fully-owned app (2026-07-13)

After fixing the mechanics above app-by-app, we went further: rather than keep multiple
describe-then-generate apps (each individually full-screen but each still generically styled), we
consolidated into **one single app** (`doknow-app`) built entirely from the verbatim-push technique,
then wired to real data (not the static sample content a plain verbatim push gives you):

- Ported our own complete HTML/CSS/JS mockup verbatim, as in blog entry 7.
- Added real `useWorkspaceDB`/DB-API reads for the home shelf (courses/modules/lessons/
  `learner_progress`) so it reflects actual state instead of fixture data.
- Ported the AI course-generation flow (topic → leveled tiers/modules/lessons, editable preview,
  approve-to-save) into this same app instead of leaving it in a separate generated app.
- Ported the lesson-playing flow (Read/Quiz/Apply tabs, real scoring, `learner_progress` writes on
  mark-complete) the same way.
- Set this one app as `defaultLandingAppId`, then deleted the now-superseded apps.

**Result, independently verified (not just reported):** a single app, full-screen from the first paint,
carrying our actual design, with a real generate → learn → track loop confirmed end-to-end by hand —
generated a real course, approved it with no error, took the quiz, marked a lesson complete, and watched
the streak and progress actually update live. This is the fullest realization of "escape the shell" we
reached: not just fixed mechanics on generated apps, but one app that behaves like a normal,
fully-owned React app sitting inside a single Audos app slot.

**One recurring gotcha hit again during this build**, worth restating because it's now happened three
times across two unrelated apps: **any new database write path needs every `json`/`jsonb` column
explicitly `JSON.stringify`'d** — see the "recurring pattern" section in
`19-capabilities-and-limitations.md`. Brief this explicitly every time; don't wait to rediscover it via
a live `Insert failed` error.

## "Avoidable from day one" — original claim and a correction (2026-07-14)

Every fix above was a retrofit — found because something broke on an app that was already built the
default way. The open question was whether a **brand-new** app, briefed with everything above stated
up front, could just... not have any of these problems. Tested it directly (the "One Good Thing"
experiment): built a small new app from scratch with steps 1–4 stated explicitly in the initial brief,
then verified by hand — zero-delay screenshot (no flash), a real database write (no JSON error), and
clicking the chat-access affordance (opened a small popup, not the old shell). At the time this was
reported as **all four issues avoided from the start, purely because they were briefed upfront.**

**That conclusion is only one-quarter correct, and the mistake is worth naming plainly.** See the
terminology note in `18-genesis-space-and-ui-ceiling.md` (workspace vs. app) — it explains why. Steps 1,
2, and 4 (full-screen landing, deep-link resolution, chat-popup-not-teardown) all live in `Desktop.tsx`,
which is **shared across the whole workspace**, not owned per app. By the time "One Good Thing" was
built, this workspace's `Desktop.tsx` had *already* been fixed, days earlier, while retrofitting other
apps. So that app inherited clean shell behavior automatically — the brief calling out steps 1, 2, and 4
had **no effect either way**, because there was nothing left in the shared shell for it to get wrong.
This wasn't isolated as a real test until the Red Pill migration (see the section below) built *another*
new app with **no** special briefing at all and watched it come out just as clean, from inheritance
alone.

**Only step 3 — `JSON.stringify` every write — is genuinely per-app code**, freshly generated by
whatever agent builds that specific app's own file, not shared plumbing. That one really was a clean,
unconfounded test: briefed upfront, and the bug that had appeared three separate times in three other
apps' code simply didn't happen here. That result stands.

**What's still actually untested:** whether steps 1, 2, and 4 can be avoided *from a genuinely fresh
workspace* — one where `Desktop.tsx` has never been touched — purely by briefing them upfront, with no
prior fix to inherit. That's a different, harder experiment (a new workspace, not a new app in this one)
and hasn't been run. Worth running later if the question matters again.

**Attempted, 2026-07-16, inconclusive — not a pass or a fail.** A genuinely fresh workspace was created
for a separate app (`field-notes`, workspace `1d30572d-2ced-4dd1-872f-3e67a74891dd`, `Desktop.tsx` never
touched by any prior fix) and a build job was briefed with all four checklist items explicitly, from the
first message, before any file existed. Per the job's own transcript, the result looked like a clean
pass: the bare workspace URL loaded the app full-screen with no chat or dock, a full password-gate auth
flow worked end-to-end, and a returning-user reload was reported "flash-free." **But the job errored out
during a later cosmetic fix without ever posting a completion report, and a direct check of the raw file
tree afterward (`GET /api/space/{id}/files`) came back completely empty — no app file, `config.json`
content empty.** None of the described work actually exists in the workspace now. See `BACKLOG.md #18`:
file-level changes made during a job appear to not be durably committed until some final step, and an
error before that step discards everything, even work already live-tested within the same run.

**Resolved, 2026-07-16 (later the same day) — confirmed YES, independently, not on a job's word alone.**
The "inconclusive" verdict above turned out to be based on a mistake in our own verification, not a real
loss. What actually happened: `GET /api/space/{id}/files` only ever reflects the *published* bundle — it
has no visibility into draft/unpublished work at all. Checking it and finding it empty proved nothing had
been *published*; it said nothing about whether the app had actually been built. A follow-up job,
scoped down to skip the (already-resolved) database work and told to stop immediately once verification
passed rather than continuing into further polish, reported that the original job's app file, `config.json`
edits, and server-function hooks were **already present** when it started — untouched, not rebuilt. It
made one small fix (a session-expiry mismatch) and re-ran the same 5-check verification.

We did not accept that report on its own. We checked directly: the workspace's own Preview panel, "Live"
vs "Draft" toggle, side by side. "Live" shows the same old, unchanged marketing page the file-tree API had
implied. **"Draft" shows the real, working app** — a genuine password-gate screen, styled correctly, no
chat shell, no dock. We logged in ourselves with the test credentials (not the job's screenshot — our own
click) and reached the actual content feed, three sections rendering correctly with proper empty states.
We clicked the "Ask Field Scout" affordance ourselves and watched it open a small, self-contained popup
over the still-visible feed, not a teardown to the old three-pane chat view.

**So: yes.** Briefed upfront, in a genuinely fresh workspace whose `Desktop.tsx` had never been touched by
any prior fix, the app was born full-screen from the first paint — no chat-shell flash, no retrofit
needed, exactly as the checklist promised. The `id`/`generateEmbedding`-style false-negative pattern this
SDK keeps re-learning applies here too: our first "it failed" conclusion was wrong not because the
platform failed, but because we checked the wrong signal. See `BACKLOG.md #18` for the retraction — the
"file writes aren't durable" bug was withdrawn the same day it was filed, once we found our own mistake,
and we followed up directly with Audos so their engineer wouldn't spend time chasing it.

**One real, smaller finding survives the retraction:** there is no obvious way to distinguish "nothing was
built" from "something was built but exists only as an unpublished draft" without specifically knowing to
check the Preview panel's Draft toggle — the file-tree API returns empty in both cases. That's a genuine
discoverability gap worth Audos documenting, even though it isn't data loss.

**The checklist, if you're starting a new app rather than fixing an old one — still the right thing to
brief, even though 1/2/4 may already be moot in a workspace whose shell was fixed previously:**
1. Set `desktop.layout.defaultLandingView`/`defaultLandingAppId` in `config.json` for full-screen intent
   from the start (note: only one app can be the *default*; other apps still need step 2 for their own
   deep links to be full-screen).
2. If the app has its own deep-link/full-screen logic, use a lazy `useState` initializer that resolves
   from `window.location.hash` synchronously — never a plain default + `useEffect` correction.
3. `JSON.stringify` every `json`/`jsonb` column on every write, from the first line of database code —
   **this is the one item confirmed to need re-briefing on every new app, unconfounded by shell
   inheritance.**
4. If you include a "back to chat" affordance, make it open a small popup — never a handler that tears
   down the full-screen state to return to the old three-pane view.

## Migrating an EXISTING app in place, not replacing it (2026-07-14)

Everything above answers "how do I avoid the matrix building something new." A separate, harder question:
if you already have a generated app you don't want to throw away and rebuild from scratch — can you pull
*that specific app* out of the matrix, in place, keeping its id, its file, its data — rather than building
a parallel app and cutting over? Tested directly on `doknow-home`.

**First finding: the mechanics layer may already be fixed for you, permanently, workspace-wide.**
`Desktop.tsx` — where the shell/deep-link/flash fixes from steps 1–4 above live — is **shared across every
app in a workspace, not owned per-app**. Once you've fixed it once (this workspace fixed it days earlier,
retrofitting Course Builder and friends), **every subsequently created app in that workspace inherits the
fix automatically**, even a plain describe-then-generate app built with zero special instructions. We
confirmed this by rebuilding a fresh baseline app from scratch with deliberately no fixes stated in the
brief, expecting it to reproduce the old shell-flash bug — it didn't. It rendered full-screen, no flash,
from the very first frame, entirely by inheritance. **Practical implication: check whether your workspace's
`Desktop.tsx` has already been fixed before assuming a "new" app needs the steps-1–4 treatment at all** —
it may not.

**Second finding: visual design is NOT workspace-shared — it's still generic per app, and still needs the
verbatim push, done in place.** The same freshly-built baseline app, despite inheriting clean mechanics,
rendered with entirely generic, self-invented styling — different copy, different layout, different card
treatment from our actual product design (it did pick up the workspace's teal brand token automatically,
which made it look closer than it was — don't let a matching accent color fool you into thinking the design
matches).

**The in-place migration recipe:**
1. Identify the app's real file (e.g. `apps/<AppName>/App.tsx`) rather than treating it as a black box.
2. Brief a build agent to port your verbatim mockup source (full CSS **and** full HTML body **and** full
   `<script>` — pasted inline, not referenced/described) into that **same file**, replacing its existing
   UI, explicitly instructed not to redesign.
3. **Explicitly instruct the agent to preserve whatever real-data wiring the existing file already has**
   (DB reads/writes) — a verbatim UI port can accidentally revert to static sample content if not told to
   graft onto the existing data logic rather than overwrite it wholesale.
4. Do not touch `config.json`'s `defaultLandingAppId` or any other app's files — this is a single-app,
   in-place operation, not a cutover.
5. Verify live with the full checklist above, **plus one more check specific to migrations**: load the
   app you did NOT touch (the previous default) side by side and confirm its data still matches — this
   catches any accidental cross-app data-scoping change the migration might have introduced.

**Gotcha hit during this exact migration:** a build agent's first pass claimed a verbatim port but had
only embedded the CSS literally — the HTML body and script were "referenced from the founder's brief"
rather than pasted in, meaning the agent was reconstructing markup from a description, not porting it
byte-for-byte. This defeats the entire point of the verbatim technique and would not have been caught
without inspecting what was actually embedded in the brief, not just trusting the "ported verbatim" claim
in the completion report. Always paste the complete literal source into the brief yourself — don't
describe it or point to "the file discussed earlier," even in the same chat thread.

**Result, independently verified live:** `doknow-home` now renders our actual design (left icon rail, teal
brand mark, streak pill, real greeting copy, matching card layout) — not a parallel app, the same app id
and file that started out generic — while still reading the same real, shared data as `doknow-app`
(matching streak, matching course progress, matching up-next lesson). The app that was NOT touched
(`doknow-app`, still the default) was confirmed unaffected. One small residual discrepancy: a "cards due
for review" count differed between the two apps post-migration (likely a spaced-repetition query difference
introduced somewhere in one app's history, not a shell or design issue) — worth tracking down later, not
blocking the conclusion that the migration itself worked.
