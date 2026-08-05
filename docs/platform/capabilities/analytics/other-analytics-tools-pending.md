# Other analytics tools — schema confirmed, not yet run

**Status: 📄 documented, unverified output.** Otto-chat-triggered tools (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`).

These exist with real parameter shapes, pulled from Otto's live tool definitions — but the *outputs* have
not been independently exercised, unlike `query-analytics.md`, `query-events.md`, and
`query-data-source.md`. Treat as "known to exist," not "known to work."

| Tool | Required | Optional |
|---|---|---|
| `get_funnel_metrics` | — | `appId`, `spaceId`, `days`, `startDate`, `endDate` |
| `query_sessions` | — | `hasEmail` (bool), `limit` (default 50) |
| `delegate_analytics_insight` | `question` | `timeframe {days, startDate, endDate, comparison}` — spawns an async subagent |

Source: Otto chat, 2026-07-23.
