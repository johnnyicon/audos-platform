# `generate_background_music` — AI music generation

**Status: 📄 schema/presets confirmed live, generation itself not yet run.** Otto-chat-triggered tool
(see `../../29-otto-tool-surface-vs-app-callable-hooks.md`). Paid — ElevenLabs Music.

**Params:** `preset` (corporate/tech/uplifting/calm/energetic/playful/cinematic/lofi), `customPrompt`,
`lengthSeconds` (3–600), `instrumental` (default true).

**Verified 2026-07-23** (list only, not generation): `list_music_presets` returned all 8 real presets,
each with a full underlying prompt string specifying BPM and "instrumental only, no vocals." The presets
themselves are confirmed real; an actual `generate_background_music` call (which would cost money) was
not run in this pass — deprioritized in favor of video/voiceover/image, which were the areas of most
interest.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0032-media-generation-live-verification.md`.
