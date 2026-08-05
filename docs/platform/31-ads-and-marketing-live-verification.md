# Ads & Marketing — index

*Written 2026-07-23. Every capability below was independently exercised — deliberately excluding*
*`delegate_ad_generation` → `launch_previewed_campaign`, the one path in this tool group that spends real*
*money and reaches a real audience. One file per tool, in `capabilities/ads/` — read only the one you need.*

See `29-otto-tool-surface-vs-app-callable-hooks.md` first: these are Otto-chat-triggered tools, not
`platform.*` hooks your own app code can call.

| Tool | Status | File |
|---|---|---|
| `get_ad_campaigns` (list) | ✅ verified | `capabilities/ads/get-ad-campaigns.md` |
| `search_meta_targeting` (Meta geo-targeting) | ⚠️ verified, real format gap found | `capabilities/ads/search-meta-targeting.md` |
| `keyword_ideas` (Google Ads keyword research) | ✅ verified | `capabilities/ads/keyword-ideas.md` |
| `generate_ad_copy` | ✅ verified | `capabilities/ads/generate-ad-copy.md` |
| `get_campaign_insights`, `get_dm_campaign_status`, `get_dm_conversations`, `delegate_ad_generation`→`launch_previewed_campaign` | 📄 schema only / deliberately unrun | `capabilities/ads/other-ads-tools-pending.md` |

Narrative: `../../blog/experiments/0031-ads-and-marketing-live-verification.md`.
