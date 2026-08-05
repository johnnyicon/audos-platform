---
date: 2026-08-05
area: identity
status: open
filed: no
label: The Otto reachable over the external API has no manage_server_functions tool, contradicting the documented creation path
---

`docs/platform/07-api-development-guide.md` documents `manage_server_functions` as the way server
functions get created — both "ask the AI assistant" (which uses the tool) and "use the
`manage_server_functions` MCP tool directly" if you're automating.

Asked the Otto reachable through the external onboarding API (`POST /api/agent/onboard/chat/:workspaceId`)
to create a small read-only diagnostic hook. It searched its own toolset twice and reported plainly that
**`manage_server_functions` does not exist in its available tools.** What it does have that's hook-adjacent
is `get_hook_logs` (read-only) and `create_dashboard` (spawns a subagent that builds a dashboard UI, not a
one-shot query). It noted from its own build log that prior probe tasks which created and called hooks were
done through **Cursor delegation jobs**, not by Otto directly.

This is the same tool-surface split recorded in `docs/platform/29` — the tools Otto has in a signed-in
workspace session are not the same set as the ones reachable from outside. `manage_server_functions` looks
to be another instance: probably present for the in-workspace Otto, absent for the external one.

**Why it matters:** doc 07 reads as though there's a programmatic creation path for anyone automating
against the platform. From outside there isn't one. The only route we found to create a server function
externally is to stage and dispatch a Cursor delegation job — heavyweight for what should be a single API
call, and dependent on Audos's own Cursor account having headroom (`docs/platform/21`).

**Not independently verified**, and worth stating as such: this is Otto reporting on its own toolset. It
could be wrong about what it has. But it searched twice, was specific about what it *does* have, and
declined to attempt a workaround rather than fabricate — the same care it showed when refusing to invent
SQL output in the same session.

**Ask:** either expose hook management as a documented external API (the backing REST endpoints appear to
exist — Otto cited build logs showing the Cursor agent listing, creating and deleting hooks over REST with
`DELETE` returning 204), or correct doc 07 to state that creation is a signed-in-session-only path.

Source: Otto chat, DoKnow workspace, 2026-08-05, during the pgvector re-verification
(`blog/experiments/0033`).
