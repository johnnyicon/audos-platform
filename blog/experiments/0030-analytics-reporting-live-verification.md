---
date: 2026-07-23
area: analytics-api
status: confirmed
label: Otto's analytics tools return real, cross-verifiable data — but one data-source path (funnel-events) contradicts another (query_events) for the same workspace
---

**Hypothesis:** Otto surfaced a list of analytics/reporting tools when asked what Audos capabilities we
hadn't captured yet in the SDK's capability matrix. Rather than record Otto's description as fact — this
project's standing rule since the "never trust a job's self-report" finding — actually run each tool
against the real DoKnow workspace and check the outputs against each other.

**Method:** Asked Otto (via the onboarding API's `/chat` endpoint) to run `query_analytics` (overview,
30d), `query_events` (by_type, 30d), and `query_data_source` against two `sourceId` values (`contacts`,
`funnel-events`) — all read-only, no side effects — and return the raw output, not a paraphrase.

**Result:** `query_analytics` returned 5 sessions / 4 contacts / 80% conversion. `query_events` returned a
187-event breakdown by type. `query_data_source(contacts)` independently returned 4 contact records with
real UUIDs and emails — matching the `query_analytics` contact count exactly, a genuine cross-check across
two separate tool calls, not the same number repeated. But `query_data_source(funnel-events)` — which
should be a structured view of the same event data `query_events` just reported 187 rows for — returned
exactly 1 record with every field blank. `query_events` is the reliable surface for event data;
`query_data_source(funnel-events)` is not, at least not as exercised here. Root cause not yet identified.

Source: Otto chat, DoKnow workspace (`8a65a4ac-5a22-435f-b55f-c41ea34ca00d`), 2026-07-23. Full detail:
`docs/platform/30-analytics-and-reporting-live-verification.md`.
