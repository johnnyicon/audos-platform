# `generate_voiceover` — AI text-to-speech, and a real voice-selection bug

**Status: ⚠️ verified live, with a real bug in the companion voice-listing tool.** Otto-chat-triggered
tool (see `../../29-otto-tool-surface-vs-app-callable-hooks.md`). Paid — ElevenLabs TTS, itemized and
wallet-deducted (the one cleanly-attributable cost in this whole pass).

**Params:** `script` (required), `voiceId` (default Sarah, `EXAVITQu4vr4xnSDxMaL`). A segmented variant,
`generate_voiceover_segments`, takes `segments: [{text, sceneLabel?}]` + `voiceId` (schema-confirmed, not
independently run).

**Verified 2026-07-23**: generated a real test line ("This is a real test of the Audos voiceover tool...").
Returned a real, per-workspace-scoped GCS URL
(`.../workspaces/8a65a4ac-…/audio/voiceover-full-*.mp3`). **Independently downloaded and probed with
`ffprobe`**: a real, decodable MP3, 8.6 seconds, 139,224 bytes, 128kbps/44.1kHz mono — not a placeholder.

**Cost — cleanly itemized, unlike video/image:** wallet showed a discrete line item,
`$0.03 — ElevenLabs TTS: 116 chars`, i.e. roughly **$0.00026/char**. This was the entire wallet-balance
deduction for the whole test session (image + video + voiceover combined) — meaning voiceover is the one
media-gen cost Audos's own tooling itemizes cleanly; video and image are not (see `generate-video.md`).

## Bug: `list_voiceover_voices` returns unusable `undefined` IDs for ~96% of voices

Asked Otto to list all available voices to see the real selection surface. It returned **~85 voices**
with real names/accents/demographics — but **every single `ID` field came back literally the string
`undefined`**. The only voice IDs that actually work are the 3 hardcoded as defaults/fallbacks in the
tool's own footer: **Sarah** (`EXAVITQu4vr4xnSDxMaL`), **Rachel** (`21m00Tcm4TlvDq8ikWAM`), **Josh**
(`TxGEqnHWrfWFTfGW9XjX`). In practice, **only these 3 voices are selectable via the documented `voiceId`
param** — the other ~82 listed voices are real but inert, since there's no working ID to pass for them.
Not yet filed as a formal bug report with Audos.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0032-media-generation-live-verification.md`.
