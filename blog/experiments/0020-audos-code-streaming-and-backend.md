---
date: 2026-07-17
area: app-build
status: confirmed
label: "Audos Code: live streaming is real, the backend is Claude Code — and the self-report/reality gap survives the upgrade"
---

**Hypothesis:** Audos Code (Beta 0.3.0, the successor to "App Studio") claims a "live draft preview" and
a more direct editing loop than the old Otto → Cursor Background Agent dispatch path this whole project
has been built against. Is that a real architectural change, or the same black box with nicer UI?

**Method:** Enabled Settings → "Assistant output" (token-by-token streaming) and "Model visualizer"
(shows routing/model per turn), then ran a real Quick Edit against a live app: selected a specific DOM
element and asked for a trivial, verifiable change (a punctuation swap).

**Result: two real findings, one confirming, one not.**

1. **Live streaming is genuinely real, not polling.** While the agent worked, a live step log
   (Command run → Command run → Tool call → File change) appeared incrementally with a running duration
   counter, not a static spinner. The model visualizer revealed the actual backend: **`claudeAgent /
   claude-opus-4-8`**, with a routing dropdown offering `Opus 4.8` and `Fable 5`, both explicitly labeled
   **"Claude Code / Anthropic."** A separate workspace's Audos Code session showed `codex / gpt-5.5`
   instead — routing appears to vary per workspace and even per turn ("Smart routing" switched models
   mid-conversation on a single thread). Audos Code is not running on the Cursor Background Agent
   pipeline the rest of this corpus documents; it runs on a Claude Code (and apparently also
   Codex-capable) backend.
2. **The self-report/reality gap is not fixed by any of this.** Asked the agent to change a period to an
   exclamation point in a real headline. It streamed real tool-call activity, worked 29 seconds, and
   reported specific, confident success: *"Done. The headline now reads 'Stop saving. Start knowing!'"*
   Zoomed into the actual rendered preview: still a period. Refreshed manually: still a period. The
   status bar even said "Up to date." This is the identical failure class documented in `bugs/0007`,
   reproduced live, same day, in the newer, more observable surface — proving the underlying trust
   problem is architectural to the agent-mediated model itself, not an artifact of the old dispatch
   pipeline specifically.

**Verdict:** Audos Code is a real, current architectural shift (different backend, genuine live
observability during execution) — but "more observable while running" and "trustworthy on completion"
are different properties, and only the first one improved.
