# `generate_video` — AI video generation (Veo3)

**Status: ✅ verified live, independently.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`). Paid — Google **Veo3** (`veo-3.1-generate-preview`),
async, ~1–2 min generation time.

**Params:** `prompt` (required), `aspectRatio` (16:9/9:16), `saveToWorkspace` (default true).

**Verified 2026-07-23**: asked Otto to generate a real 16:9 test clip ("a simple animated checkmark
appearing over the text DoKnow API test, 2026-07-23"). Returned a real GCS URL
(`storage.googleapis.com/audos-images/generated-videos/models_veo-3.1-generate-preview_operations_*.mp4`).
**Independently downloaded and probed with `ffprobe`** rather than trusting the URL alone: a real,
decodable MP4, 8.00 seconds long, 627,820 bytes — not a placeholder or empty file.

## Cost — could not be cleanly isolated

Asked Otto to check `query_wallet`/wallet expenses before and after generation to measure the real cost.
Wallet balance moved **$2733.50 → $2733.47** (−$0.03) across the whole test session (image + video +
voiceover + all the free lookups). Otto could cleanly itemize the **voiceover** charge alone
(`$0.03 — ElevenLabs TTS: 116 chars` — see `generate-voiceover.md`) as the entire wallet-balance
deduction. The Veo3 video **did not appear as a discrete wallet line item at all** — a separate "AI token
usage (today)" bucket rose from $0.17 to $0.52 (+$0.35), but that bucket also absorbs image generation and
general agent token use from the same session, so **no honest per-video dollar figure could be
attributed**. Per Otto: this would require Audos's own billing internals, which its tools don't expose.
**Practical takeaway: budget for "some non-trivial amount, non-itemized" rather than a per-video unit
price** until Audos exposes clearer billing granularity.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0032-media-generation-live-verification.md`.
