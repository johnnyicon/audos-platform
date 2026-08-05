# `search_meta_targeting` — Meta geo-targeting resolution

**Status: ⚠️ verified live, with a real input-format gap.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`). Free, read-only.

**Params:** `query` (required), `countryCode` (default US), `limit`, `types` (country/zip/geo_market/
city/region).

**Verified 2026-07-23** with three input forms against the real tool:

| Input | Result |
|---|---|
| `"Austin, TX"` (comma-separated city+state) | ❌ **No deterministic match.** The tool's own error message suggests a ZIP, a bare city, or a DMA name instead. |
| `"Austin"` (bare city) | ✅ 5 matches — "Austin" exists across TX/MN/IN/AR/CO. Resolves to `type=city, key=2525495` for the Texas one, default 25mi radius — **but you must pick the right one out of 5.** |
| `"78701"` (ZIP) | ✅ 1 unambiguous match — `type=zip, key=US:78701`. |

**Practical rule: never pass `"City, State"` — it silently fails rather than being normalized. Prefer a
ZIP code over a bare city name** when you have one; a bare city name is fine but requires a
disambiguation step that a ZIP doesn't.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0031-ads-and-marketing-live-verification.md`.
