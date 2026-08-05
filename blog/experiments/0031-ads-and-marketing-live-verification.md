---
date: 2026-07-23
area: ads-marketing
status: confirmed
label: Ad-copy, keyword research, and campaign-listing tools all produce real, on-target output — but Meta geo-targeting silently fails on the "City, State" format most people would type first
---

**Hypothesis:** Same as the analytics pass — Otto listed a set of ad/marketing tools it can run, and
rather than trust the description, actually exercise the free/no-launch ones and see what comes back.

**Method:** Asked Otto to run `get_ad_campaigns`, `search_meta_targeting` ("Austin, TX"), `keyword_ideas`
(seeded with two DoKnow-relevant terms), and `generate_ad_copy` (targeted at the actual DoKnow audience) —
deliberately excluding `delegate_ad_generation` → `launch_previewed_campaign`, the one path in this group
that spends real money and reaches a real audience.

**Result:** Three of four came back clean and real: `get_ad_campaigns` correctly reported empty (no
campaigns exist yet — a clean empty result, not an error). `keyword_ideas` returned 10 real Google-Ads-shaped
rows (search volume, trend, competition index, bid range) with plausible, non-templated variation across
rows. `generate_ad_copy` produced three genuinely on-target variants referencing DoKnow's actual pain point
(bookmark backlog, forgotten reading) rather than generic filler — and ran cold, with no prior ad history,
without erroring.

The one real gap: `search_meta_targeting` with `"Austin, TX"` (comma-separated city+state) returned no
match at all. The tool wants a bare city name, a ZIP code, or a DMA name — the most natural first thing to
type ("City, State") silently fails to resolve rather than being normalized or falling back gracefully.
Follow-up confirmed both alternatives work, with a wrinkle: a bare city name ("Austin") returns 5
same-named cities across different states, forcing a disambiguation step, while a ZIP ("78701") resolves
to exactly one unambiguous match — ZIP is the more reliable input of the two.

Source: Otto chat, DoKnow workspace, 2026-07-23. Full detail:
`docs/platform/31-ads-and-marketing-live-verification.md`.
