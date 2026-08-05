# Analytics & Reporting — index

*Written 2026-07-23. Otto surfaced these tools when asked "what can Audos do that we haven't captured?" —*
*each was independently exercised against the real DoKnow workspace rather than taken on Otto's word.*
*One file per tool, in `capabilities/analytics/` — read only the one you need.*

See `29-otto-tool-surface-vs-app-callable-hooks.md` first: these are Otto-chat-triggered tools, not
`platform.*` hooks your own app code can call.

| Tool | Status | File |
|---|---|---|
| `query_analytics` (sessions/contacts/conversion overview) | ✅ verified | `capabilities/analytics/query-analytics.md` |
| `query_events` (event breakdown) | ✅ verified | `capabilities/analytics/query-events.md` |
| `query_data_source` (contacts, funnel-events, + 10 more sourceIds) | ⚠️ verified, real bug found | `capabilities/analytics/query-data-source.md` |
| `get_funnel_metrics`, `query_sessions`, `delegate_analytics_insight` | 📄 schema only | `capabilities/analytics/other-analytics-tools-pending.md` |

Narrative: `../../blog/experiments/0030-analytics-reporting-live-verification.md`.
