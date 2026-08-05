# Media Generation — index

*Written 2026-07-23. Video and voiceover were independently verified by downloading the returned files*
*and probing them with `ffprobe` — not just checking the URL returned 200. One file per capability, in*
*`capabilities/media/` — read only the one you need.*

See `29-otto-tool-surface-vs-app-callable-hooks.md` first: these are Otto-chat-triggered tools, not
`platform.*` hooks your own app code can call. All of these draw on real workspace credits/wallet except
where noted free.

| Capability | Status | File |
|---|---|---|
| `generate_image` | ✅ verified, artifact independently inspected | `capabilities/media/generate-image.md` |
| `generate_video` (Veo3) | ✅ verified, artifact independently inspected | `capabilities/media/generate-video.md` |
| `generate_voiceover` (ElevenLabs TTS) | ⚠️ verified, real bug in `list_voiceover_voices` found | `capabilities/media/generate-voiceover.md` |
| `generate_background_music` | 📄 presets confirmed, generation not run | `capabilities/media/generate-background-music.md` |
| `search_stock_photos` (Unsplash) | ✅ verified, free | `capabilities/media/search-stock-photos.md` |
| `redesign_logo`, carousel tools | 📄 schema only, not run | `capabilities/media/other-media-tools-pending.md` |

**Cost note:** only the voiceover charge was cleanly itemized by Audos's own wallet tooling
(`$0.03 — ElevenLabs TTS: 116 chars`). Image and video generation both drew from a non-itemized "AI token
usage" bucket with no per-generation dollar figure exposed — see `capabilities/media/generate-video.md`
for the full wallet-delta methodology.

Narrative: `../../blog/experiments/0032-media-generation-live-verification.md`.
