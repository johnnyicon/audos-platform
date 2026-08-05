# `generate_ad_copy` — ad copy generation

**Status: ✅ verified live.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`). AI-gen (unpriced by the schema, but no wallet
line item was attributed to it separately in the one cost-tracking pass we ran — see
`../media/generate-video.md` for the wallet-delta methodology).

**Params:** `count` (3–10, default 5), `targetAudience`, `analysisContext`.

**Verified 2026-07-23**, count=3, targetAudience="knowledge workers who read a lot but forget what they
learn". Returned three distinct, genuinely on-target variants (headline + primary text each) that
reference the actual product's real pain point — bookmark/read-later backlog, forgotten reading — rather
than generic ad-copy filler. Example: *"Your read-later list has 200 items and you've finished maybe
3."* Ran cold, with **no prior ad history and no prior `delegate_ad_analysis` call**, without erroring —
so a first-ever ad-copy generation doesn't require historical performance data as a precondition.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0031-ads-and-marketing-live-verification.md`.
