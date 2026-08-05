---
date: 2026-07-23
area: media-generation
status: confirmed
label: Real video (Veo3), voiceover (ElevenLabs), and image generation all verified by downloading and probing the actual files — plus a real bug where 82 of 85 listed voices have unusable IDs
---

**Hypothesis:** Otto listed video/voiceover/music/image generation as capabilities the SDK's matrix hadn't
captured. A URL in a chat response isn't proof of anything by itself — the standing rule on this project
is to independently verify a claimed artifact, not just trust that the tool "said" it succeeded.

**Method:** Asked Otto to run `generate_image`, `generate_video` (Veo3), and `generate_voiceover`
(ElevenLabs) for real, plus list the free lookup tools (`list_voiceover_voices`, `list_music_presets`,
`search_stock_photos`). For every generated artifact, downloaded the file directly and inspected it —
`file`/dimensions for the image, `ffprobe` for video/audio duration and validity — rather than treating a
200 response on the URL as sufficient. Also asked Otto to check wallet balance before and after, to see
whether the real cost could be attributed per-generation.

**Result:** All three generations were real. The image was a correctly-rendered 1536×1024 PNG containing
the exact requested text. The video was a genuine 8-second, 627KB MP4 (`ffprobe`-verified) from Veo3. The
voiceover was a genuine 8.6-second, 139KB MP3 (`ffprobe`-verified) from ElevenLabs. Cost was only cleanly
attributable for the voiceover — a discrete `$0.03 — ElevenLabs TTS: 116 chars` wallet line item. Image and
video both drew from a shared, non-itemized "AI token usage" bucket with no way to isolate a per-generation
dollar figure.

One real bug surfaced along the way: `list_voiceover_voices` returns ~85 voices with real names/accents,
but every one of their `ID` fields comes back as the literal string `undefined` — only the 3 hardcoded
fallback voices (Sarah, Rachel, Josh) are actually selectable via the documented `voiceId` parameter. The
other ~82 are real but functionally inert.

Source: Otto chat, DoKnow workspace (`8a65a4ac-5a22-435f-b55f-c41ea34ca00d`), 2026-07-23. Full detail:
`docs/platform/32-media-generation-live-verification.md`.
