# `generate_image` — AI image generation

**Status: ✅ verified live, independently.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`). Paid — model `gpt-image-1.5` (Gemini fallback per
Otto), per-unit price not exposed to Otto's own tooling.

**Params:** `prompt` (required), `aspectRatio` (1:1/16:9/9:16).

**Verified 2026-07-23**: asked Otto to generate a 16:9 test image with specific text
("DoKnow API test / 2026-07-23" + a checkmark). Returned a real, permanent GCS URL
(`storage.googleapis.com/audos-images/generated-images/image-*.png`). Independently confirmed — not just
"the URL returned a 200" but downloaded and rendered it directly: a real 1536×1024 PNG, correctly
containing the exact requested text and a green checkmark, `last-modified` timestamp matching the moment
of generation. This is the strongest form of verification in this whole pass — the artifact itself was
inspected, not just Otto's claim about it.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0032-media-generation-live-verification.md`.
