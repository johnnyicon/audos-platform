# Audos SDK — Build Log

The dedicated knowledge base for our real-world work on the Audos platform, across both products
we've built on it: **Throughline** (still live, still running its GTM on Audos) and **DoKnow** (our
current flagship end-to-end build, and our formal agreement with Audos to surface what we find). This
is one continuous, chronological record — not two separate logs — because the platform lessons compound
across both.

Source of truth for the underlying reference material: `docs/platform/` (capability findings) and
`BACKLOG.md` (issues raised with Audos). This blog is the narrative layer on top of that data.

## Entries

1. [No smart defaults: what Desktop.tsx actually controls](0001-no-smart-defaults-desktop-tsx.md) — 2026-04-03 · Throughline
2. [Multi-tenant on Audos: session scoping, auth, and the org_id pattern](0002-multi-tenant-session-scoping-auth.md) — 2026-04-03 · Throughline
3. [Kickoff: standing up DoKnow on Audos](0003-kickoff-standing-up-doknow-on-audos.md) — 2026-07-12 · DoKnow
4. [It felt like ChatGPT with apps on the side](0004-it-felt-like-chatgpt-with-apps-on-the-side.md) — 2026-07-12 · DoKnow
5. [Escaping the shell: can Audos build a real product dashboard?](0005-escaping-the-shell.md) — 2026-07-12 · DoKnow
6. [The empty course mystery: three rounds to fix lesson persistence](0006-the-empty-course-mystery.md) — 2026-07-12 · DoKnow
7. [Design fidelity: can we push our own UI and have it stick?](0007-design-fidelity-can-we-push-our-own-ui.md) — 2026-07-12 · DoKnow
8. [Three shell bugs, and a verification method we didn't know we needed](0008-three-shell-bugs-and-a-verification-blind-spot.md) — 2026-07-13 · DoKnow
9. [Escaping the matrix: one app, no remnants, real data](0009-escaping-the-matrix.md) — 2026-07-13 · DoKnow
10. [Building it right the first time](0010-building-it-right-the-first-time.md) — 2026-07-13 · DoKnow
11. [Taking the red pill](0011-taking-the-red-pill.md) — 2026-07-14 · DoKnow
12. [What Audos can actually do](0012-what-audos-can-actually-do.md) — 2026-07-15 · DoKnow
13. [Building field-notes in the open](0013-building-field-notes-in-the-open.md) — 2026-07-16 · field-notes
14. [The audit, the upgrade, and the open door](0014-the-audit-the-upgrade-and-the-open-door.md) — 2026-07-17 · field-notes
15. [No live path to our own content](0015-no-live-path-to-our-own-content.md) — 2026-07-17 · field-notes
16. [Audos says fixed — one confirmed false, two in flight](0016-audos-says-fixed-three-real-checks-in-flight.md) — 2026-07-22 · field-notes
17. [We ported Throughline onto Audos, then mostly left anyway](0017-we-ported-throughline-onto-audos-then-mostly-left.md) — 2026-04-20 · Throughline
18. [Delegating a real website edit to Otto over the API, not the chat UI](0018-delegating-a-real-website-edit-to-otto-over-the-api.md) — 2026-07-08 · Throughline

More entries land as either build continues. The rendered page is generated from these Markdown files
by `scripts/build_blog.py` — see `HOW-TO-UPDATE.md` for the procedure. This index is a convenience for
browsing the repo; only the `.md` files and the script drive the actual published page.

The published page has five tabs, not just this narrative one: **Blog** (this list), **Experiments**
(`experiments/` — hypothesis-driven capability tests, distinct from bugs), **Bugs** (`bugs/`),
**Feature Requests** (`feature-requests/`), and **SDK** (an index of `docs/platform/`'s durable reference
docs, distinct from all of the above — see `scripts/build_blog.py`'s `PLATFORM_DOC_BLURBS`).
See `HOW-TO-UPDATE.md` for how to add to each.
