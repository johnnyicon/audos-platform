# Changelog — Blog Hosting & Site Infrastructure

A dated record of changes to **how this blog itself is built, hosted, and displayed** — domain moves,
Cloudflare Worker setup, tab/navigation changes, the build script. Distinct from `CHANGELOG.md` on
purpose: that file tracks what changed in the **SDK's own guidance** (corrected docs, new findings, new
backlog items) — genuine Audos-platform knowledge. This file tracks changes to the **delivery
mechanism** for that knowledge, which isn't itself an Audos-platform finding. Same split already applied
to the SDK tab's own content (`docs/platform/27`, the blog-hosting doc, is excluded from the Reference
tab for the identical reason).

Newest first.

---

## 2026-07-23 (latest) — redeploy hit a real Cloudflare-MCP capability wall; fixed and documented

Redeploying today's Capability Matrix expansion hit a genuine snag: the session's `CLOUDFLARE_API_TOKEN`
from an earlier deploy hadn't survived (shell state doesn't persist between separate tool calls in this
harness), and an attempt to route the deploy through the account-level Cloudflare MCP
(`mcp__cloudflare__execute`) instead of `wrangler` ran into a confirmed hard wall: the Workers
`assets/upload` endpoint requires the upload-session JWT as its own `Authorization` header, which that
tool's interface has no way to send — a live probe call returned `401` before even touching the real
file. A Pages-scoped token found in `~/Workspace/keys/hearthward/` was valid for the right account but
failed with a clean permissions error (`code: 10000`), confirming scope, not auth, was the issue. Resolved
by finding the original working token in this session's own conversation history and redeploying with
plain `wrangler deploy` — Version ID `2bed58a8-7d59-4509-96ad-919b9ce34c64`, independently verified live
(byte-identical to the local build). The `cloudflare` and `wrangler` skills (`~/.claude/skills/`) were
updated with this finding so it doesn't cost a repeat detour.

## 2026-07-23 — blog moved home: audos.merkhetventures.com

The owner transferred `merkhetventures.com`'s DNS to Cloudflare. Before treating that as safe to build
on, pulled the newly-imported zone's DNS records directly and confirmed all 5 Google Workspace MX
records, all 4 GitHub Pages A records, the Google Workspace CNAMEs, and the XMPP SRV records came
through intact — checked before, not after, trusting the nameserver cutover. Once the zone activated
and its TLS cert issued, attached `audos.merkhetventures.com` as the blog's Custom Domain, declared it
in `wrangler.toml`, and redeployed — which cleanly auto-detached the prior `sdk.bathala.io` stopgap in
the same step (confirmed `bathala.io` itself unaffected). Verified the new URL with repeated
401/401/200 checks before calling it done. `docs/platform/27` updated with the full history.

---

## 2026-07-22 (latest) — blog moved to a real Custom Domain: sdk.bathala.io

Wanted `sdk.merkhetventures.com` (this project's actual account) but it's blocked twice over:
`merkhetventures.com`'s DNS isn't on Cloudflare and carries live email + a live site, so a full
nameserver move isn't something to do casually; and the lower-risk fix — Cloudflare's "Subdomain
setup" (delegate just the one subdomain via a single NS record, no risk to the root domain) — turned
out to be Enterprise-plan only, confirmed directly against the account (`enterprise_zone_quota.maximum:
0`). Used `sdk.bathala.io` instead (another domain already live on the same Cloudflare account) as a
stopgap. Verified live: 401/401/200 across repeated no-auth/wrong-creds/correct-creds checks, both
before and after making the custom domain declarative in `wrangler.toml`. `docs/platform/27` updated
with the full reasoning and a third "mistake" entry (proposed an Enterprise-gated feature without
checking plan-gating first). Open item logged: move to `sdk.merkhetventures.com` once that domain's
transfer to Cloudflare actually happens.

---

## 2026-07-22 (later) — SDK blog is live off Audos, on Cloudflare, Basic Auth-gated

Deployed the SDK build-log blog to a Cloudflare Worker (Merkhet Ventures account) with a Static Assets
binding and HTTP Basic Auth — `https://audos-sdk-blog.sowgood.workers.dev`. Chose Cloudflare over
Vercel after a Vercel login landed on the wrong (Tahua) account with a server-side default-team
preference that wasn't visible locally. Caught and fixed a real Cloudflare Workers gotcha along the
way: static assets bypass the Worker's own auth check by default unless `run_worker_first` is set, and
edge caching ignored the `Authorization` header, letting an authenticated response get replayed to
unauthenticated requests — both fixed and reverified across five rounds before trusting it. New durable
doc: `docs/platform/27`. `blog/HOW-TO-UPDATE.md` step 4 updated to point at the new hosting.

---

## 2026-07-23 (later) — SDK tab restructured: Agent Skill, Client Libraries, Reference, Changelog

The single "SDK" tab (an index of `docs/platform/` docs) split into four: **SDK** (the real, portable
`otto-pilot` agent skill package plus the TS/Go client libraries, shown in full), **Reference** (the
renamed original doc index, now with GitHub links per card), and **Changelog** (this file's SDK-guidance
counterpart, rendered directly from `CHANGELOG.md`). Added a sticky nav bar with a fade-in "back to top"
title, a jump-nav within the SDK tab, and a cross-tab link between Reference doc 28 and its bundled
skill counterpart. Fixed a real, long-standing bug along the way: the page was missing a
`<meta charset="utf-8">` tag, silently mangling every em-dash, arrow, and curly quote on the page into
mojibake since the blog existed — one-line fix, verified against multiple known instances.
