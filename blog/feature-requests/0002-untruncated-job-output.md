---
date: 2026-07-13
priority: 1
status: not filed
label: Untruncated job output retrievable via API/chat, not only a UI panel
---

Five separate times, a job's own completion summary got cut off mid-sentence in the exact spot that
mattered — the root-cause line, the fixability verdict, the "what changed" detail — with no API-accessible
way to fetch the rest, only a UI panel a human has to click into. For an API-driven agent, that's a dead
end every time it happens.

The underlying Cursor Cloud Agents API truncates its own stream by design with no documented full-log
endpoint, which may make this a genuine upstream constraint Audos inherited rather than a withheld
feature. That doesn't remove the ask: Audos could still persist and expose the full run transcript
server-side, independent of what Cursor's own stream truncates.
