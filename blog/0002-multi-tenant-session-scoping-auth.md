---
date: 2026-04-03
product: Throughline
status: fixed
label: Multi-tenant & auth
---

# Multi-tenant on Audos: session scoping, auth, and the org_id pattern

Second entry from the original Throughline build, same week as the first. Where that one was about the
platform shell, this one is about data — specifically, the two ways we watched Throughline's own
records disappear on us before we understood why.

We'd built the first apps (Setup, Briefing) against `db-api` with plain REST calls. Rows went in fine.
Then we switched some reads to the platform's `useWorkspaceDB` hook and those same rows vanished from
the UI. The cause: `useWorkspaceDB` defaults to `shared: false`, which silently scopes every query to
rows whose `session_id` matches the *current browser session*. Data written by a direct REST call has
`session_id = NULL` — it was never invisible, it just never matched. The fix going forward was a fixed
rule, not a one-off patch: always pass `shared: true` and filter by an explicit `org_id` column
ourselves.

Second: reading *who's* logged in had been going through the wrong door. We'd been using
`window.useSubscription?.()?.email` as a fallback — plausible-looking, undocumented, and wrong. The
actual accessor is `window.__spaceContext?.username`, populated once the EmailGate fires and stable for
the session (same email, same session, across devices). Small API, easy to get backwards, and nothing
in the platform surfaces the correct one over the incorrect one — you just have to know.

Third, worth flagging for anyone building a server function: the `db-api` endpoint itself has no
authentication by default. Anyone with the URL can call it. Adding an `x-api-key` check inside the hook
is possible but has to be done by hand — it isn't a toggle.

> The bigger decision this entry settles: **column-based isolation, not workspace-per-org.** Audos makes
> it easy to reach for a separate workspace per customer, since that's the unit the platform itself is
> built around. For a real multi-tenant SaaS, that's the wrong shape — `user_id` and `org_id` columns on
> every table, filtered at the data layer on every read, is the pattern we committed to, wrapped later in
> a `useOrgDB` helper so the filter can't be forgotten call-site by call-site.

None of this is exotic. It's the ordinary shape of multi-tenant auth bugs anywhere. What's specific to
Audos is that the platform's own convenience hook (`useWorkspaceDB`) has a default that looks like it's
helping and is actually just quietly filtering your data out from under you.
