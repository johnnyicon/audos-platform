---
date: 2026-04-16
area: app-build
status: open
filed: no
label: GitHub App connection broke entirely (no OAuth redirect), and Otto's "New Chat" doesn't start a clean session
---

Two separate, unrelated breaks hit in the same session. The Developer panel showed "No GitHub
repository connected" despite a previously-working GitHub Dev Mode connection; both "Install GitHub
App" and "Relink existing installation" triggered an "Authentication required" toast with no OAuth
redirect and no other visible action — the connection flow simply didn't run. This blocked verifying an
unrelated fix (a marker commit pushed to confirm compile behavior) mid-test, leaving that test
inconclusive rather than confirmed either way.

Separately: clicking "New Chat" in Otto doesn't actually start a clean session — it loads old
conversation history instead of a blank one, making it unreliable as a way to get a fresh, uncontaminated
context for a new question.

Neither issue was independently re-verified as fixed as of this write-up — flagged as a real, current
gap at time of writing rather than a resolved one.

Source: `throughline-forge/tmp/nick-preservation-test-report.md`.
