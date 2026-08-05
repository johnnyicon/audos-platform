Title:
`#app-id` deep-link full-screen routing works for one app, not another, on the same fix

Workspace:
- UUID: `8a65a4ac-5a22-435f-b55f-c41ea34ca00d`
- Slug: `workspace-156396`
- Owner email: `john@merkhetventures.com`
- Relevant URL: `https://audos.com/space/workspace-156396`

Intent:
I am trying to make every app in this workspace load full-screen (no chat shell) whether reached via
the base URL or a direct `#app-id` link, so that the workspace behaves as one consistent product
surface, not a chat assistant with some routes escaping it and others not.

Expected behavior:
With `desktop.layout.defaultLandingView: "app"` set, and after editing `Desktop.tsx`'s deep-link effect
(~line 572) so the `matchingApp` branch renders the same full-screen state path as the default-landing
route (instead of `openPanel(...)`, which docks the app inside the chat shell), **every** `#app-id`
hash route should render full-screen.

Actual behavior:
The fix works for `#course-builder` — verified full-screen, no chat shell, on both a cache-busted load
and a plain reload. It does **not** work for `#doknow-mockup-test` — that route still opens inside the
chat shell as a right-docked panel, on the identical Desktop.tsx code path, even after a cache-busted
reload (ruling out stale-bundle caching as the cause). We could not obtain the `doknow-mockup-test`
entry from `config.json` to compare against `course-builder`'s entry — every diagnostic job asking for
that diff returned a response truncated exactly at the point the differing field would appear.

Reproduction steps:
1. Open `https://audos.com/space/workspace-156396#course-builder` → renders full-screen, confirmed.
2. Open `https://audos.com/space/workspace-156396#doknow-mockup-test` → renders inside the chat shell
   as a docked panel, not full-screen.
3. Both hashes route through the same `Desktop.tsx` deep-link effect and the same
   `defaultLandingView: "app"` config.

Evidence:
- File/line: `Desktop.tsx`, initial-mount deep-link effect, ~line 572 — reads
  `window.location.hash.slice(1)` into `deepLinkId`, looks up `matchingApp` in `config.apps`.
- `course-builder` app entry (partial, confirmed): `{ "id": "course-builder", "name": "Course Builder",
  "icon": "GraduationCap", "component": "apps/CourseBuilder/App.tsx" }`.
- `doknow-mockup-test` app entry: could not be retrieved — every diagnostic response truncated at this
  exact point across three separate attempts.
- Timestamp: 2026-07-13, ~20:00–20:31 UTC.

Impact and priority:
Low urgency for us specifically (we're removing the `doknow-mockup-test` app regardless), but the
underlying inconsistency is a real product concern: the same fix, on the same code path, produces
different results per app with no visible reason. Anyone trying to build a genuinely single-shell
product on Audos will hit this the moment they have more than one app.

Security boundary:
No bearer token included in this report.

Success criteria:
- Confirm whether `course-builder` and `doknow-mockup-test`'s `config.json` entries differ in any field
  relevant to full-screen eligibility (e.g. an app-type flag, an allowlist tied to `defaultLandingAppId`
  specifically rather than any deep-linked app).
- If they differ: document what field controls this so we (and future builders) can set it correctly.
- If they don't differ: this is a `Desktop.tsx`/runtime bug independent of config, worth Audos's own
  engineering attention.

Requested outcome:
Please either explain the actual differentiator (so we can replicate it for future apps), or confirm
this is a genuine inconsistency in how the deep-link effect resolves full-screen eligibility across
apps in a workspace.
