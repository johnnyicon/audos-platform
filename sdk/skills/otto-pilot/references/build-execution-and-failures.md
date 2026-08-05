# Build execution and failure modes

Read this when a build job fails, stalls, or you're deciding whether to trust its reported status.

## The execution backend

Landing-page and app/space build/edit jobs triggered by the onboarding API are delegated to **Cursor
Background Agents** — Otto's own status labels the executor "App agent / Cursor." This matters because
build throughput and success depend on **Audos's own Cursor account's** usage/billing headroom, not
just your Audos subscription.

When that shared account is over its hard limit, jobs **instant-fail at the delegation step** with
`usage_limit_exceeded` ("Background Agent requires at least $2 remaining until your hard limit"). This
is account-level and typically transient — retrying once headroom is restored usually succeeds.

**Chat is not gated by this.** `/chat` and `/chat/:workspaceId` keep working even while builds are
blocked — use Otto to probe or report build status during an outage rather than assuming everything is
down.

The failure message links to `cursor.com`, which reads as if it's the workspace owner's own account
issue. It almost certainly isn't — it's Audos's shared build infrastructure. Confirm via Otto or Audos
support rather than assuming there's a self-serve fix on the owner's side.

## Don't trust a job's own status report — verify independently

A job marking itself "Complete" (or a status endpoint reporting success) is a claim, not proof. The
platform has repeatedly self-reported success on builds/edits that, checked directly (loading the real
URL cold, reading the actually-served file, querying the actual database), turned out not to have
landed. Two concrete patterns worth watching for:

- **A job reports success from its own internal testing, then errors on something unrelated before
  filing a completion report** — whatever it tested live during the run may not have actually been
  committed. If a job's last visible state is an error, don't assume its earlier "it works" claims
  survived.
- **Two jobs dispatched at overlapping times against the same file/target** can silently clobber each
  other, or the live site can keep serving a stale published snapshot underneath both.

The reliable pattern: **dispatch, then poll, then verify independently** — check the actual live
behavior yourself (a fresh page load, a direct query, an unauthenticated request) rather than trusting
the job's or the platform's own report of what happened.
