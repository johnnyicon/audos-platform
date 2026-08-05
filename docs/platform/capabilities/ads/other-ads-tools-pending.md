# Other ads tools — schema confirmed, mostly not yet run

**Status: 📄 documented, unverified** (except the launch path, which is deliberately unrun for a different
reason — see below). Otto-chat-triggered tools (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`).

| Tool | Required | Optional | Why not run |
|---|---|---|---|
| `get_campaign_insights` | `campaignId` (number) | — | No campaigns exist on this workspace to query |
| `get_dm_campaign_status` | — | `campaignId` | No DM campaigns exist |
| `get_dm_conversations` | — | `campaignId`, `limit` (default 20) | No DM conversations exist |
| `delegate_ad_generation` → `launch_previewed_campaign` | preview then confirm | budget, duration, geoTargeting, creatives | **Deliberately not run** — this is the one real-money, real-audience path in the ads tool group. Meta bills whatever budget you approve. Not exercised as part of a capability-matrix probe; would need a separate, explicit go-ahead tied to an actual campaign someone wants to run. |

Source: Otto chat, 2026-07-23.
