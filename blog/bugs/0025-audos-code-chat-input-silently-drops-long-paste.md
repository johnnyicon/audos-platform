---
date: 2026-07-17
area: app-build
status: open
filed: no
label: "Audos Code's chat input can silently fail to transmit a long, multi-paragraph message"
---

Sent a single message to Audos Code containing six structured content records, formatted as readable
multi-line prose with blank lines between entries (roughly 3,000 characters). The agent's response:
*"I can insert them, but the six row values did not come through in your message."* It had genuinely not
received the row data — it went and checked the workspace and database for a seed source rather than
inventing placeholder content, then asked for the payload again.

> The same content resent as compact single-line JSON objects (same total length, no blank lines,
> no soft-wrapped prose) transmitted correctly on the next attempt and the agent proceeded normally.

Root cause not confirmed — could be a client-side paste/textarea handling issue with long inputs
containing many blank lines, a payload-size or formatting quirk in the chat submission path, or something
else entirely. What's confirmed: a real, reproducible case of a long structured message not reaching the
model, with no error shown to the sender at submit time — the only signal was the agent's own honest
follow-up. Worth noting as a genuine positive alongside the bug: the agent did not fabricate the missing
data to comply with the request, which is exactly the discipline this whole project has been asking
for elsewhere (`feature-requests/0001`). The gap is in the input pipeline, not the model's judgment.
