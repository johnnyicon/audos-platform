# `keyword_ideas` — real Google Ads keyword research

**Status: ✅ verified live.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`). Free (draws on Google Ads data, no ad spend).

**Params:** `seed` — one of `{type: keyword, keywords: [...]}`, `{type: url, url}`, `{type: site, site}` —
plus optional `pageSize` (1–50). A companion tool, `keyword_historical_metrics`, takes 1–50 keywords
directly and is rate-limited to ~1 QPS batch (schema-confirmed, not independently run).

**Verified 2026-07-23**, seed=`{type: keyword, keywords: ["personal knowledge management",
"spaced repetition learning app"]}`, pageSize=10. Returned 10 real rows with internally-consistent,
non-templated Google Ads metrics — e.g. "spaced repetition app": 1,000 avg/mo, flat trend, LOW competition
(index 8), $0.83–$3.44 top-of-page bid range. Volume, trend %, and competition index vary independently
across rows in a way that matches real ad-platform data rather than a fixed pattern (higher volume doesn't
mechanically imply a higher or lower bid, for instance).

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0031-ads-and-marketing-live-verification.md`.
