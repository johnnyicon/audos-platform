# `query_events` — event breakdown

**Status: ✅ verified live.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`).

**Params:** `aggregation` (none/by_type/by_day/by_space/by_app/summary), `eventType`, `appId`, `spaceId`,
`days`, `startDate`, `endDate`, `limit`. Free.

**Verified 2026-07-23**, `aggregation=by_type, days=30` against the real DoKnow workspace: 187 total
events — `app_opened`: 67 · `space_entered`: 56 · `email_submit`: 6 · `lead`: 6 · `agent_view_opened`: 5 ·
`agent_message`: 2 · `thread_created`: 1. Plausible against known usage (dev/QA passes on this workspace
since 2026-07-12, no real end-user traffic yet) — not independently re-derived from a second source, so
treat the plausibility check as weaker than the `query_analytics` cross-check.

**This is the reliable surface for event data.** Compare `query-data-source.md` — the `funnel-events`
`sourceId` on that more general tool returned a degenerate 1-row/all-blank result for the identical
workspace and time window. Use `query_events`, not `query_data_source(funnel-events)`.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0030-analytics-reporting-live-verification.md`.
