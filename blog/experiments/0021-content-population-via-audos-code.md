---
date: 2026-07-17
area: db-api
status: confirmed
label: "Can Audos Code populate real data when direct Postgres access is completely broken?"
---

**Hypothesis:** `bugs/0023` established direct-Postgres credential generation is 0-for-3 across every
attempt this cycle. Is there a second, working path to write real data into a workspace's database
without that escape hatch — specifically, Audos Code's chat-driven `db.insert`/MCP tool-call layer?

**Method:** Asked Audos Code to insert six real, representative rows (two blog posts, two bugs, two
experiments — actual corpus content, not placeholders) into `content_items`, a table that had sat at
zero rows since it was created. First attempt sent as multi-paragraph prose; the row data never reached
the agent (see `bugs/0025`). Resent as compact JSON; this time it landed. The agent queried the table
back itself before declaring success — not just claiming a number, showing the actual returned rows.

**Result: confirmed pass, independently verified two separate ways, not trusted from the agent's report
alone.** Row count went 0 → 6, matching the exact slugs sent. Verified myself: (1) the workspace's own
Preview panel, Draft, signed in as the test user — all six items rendered correctly, split into the
right sections with the right status chips; (2) the **Live** site's own embedded AI assistant (a
separate feature, unrelated to Audos Code) correctly answered "show me the latest bugs" using the exact
two bugs just inserted — with **no publish action taken**, confirming database writes are visible on
Live immediately, consistent with the existing finding that DB writes are durable the instant they
happen while only the file/app bundle needs an explicit publish.

**Verdict:** the credential-generation bug blocks one specific escape hatch (direct SQL), not database
writes in general. Audos Code's DB-write layer is a second, currently-working path to the same
destination — real, verified, usable today, even while the direct-Postgres route stays broken.
