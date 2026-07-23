# Updating the Audos SDK blog

The blog is generated from Markdown. Nobody hand-edits HTML. This works the same for any agent
(Claude, Codex, or a human) — it's a plain Python 3 stdlib script, no dependencies to install.

## Add a new entry

1. Write `blog/NNNN-your-slug.md`, one number higher than the last entry, with this exact frontmatter
   shape at the top:

   ```markdown
   ---
   date: 2026-07-20
   product: DoKnow
   status: pass
   label: Short scorecard label
   ---

   # The full post title

   Body paragraphs, plain Markdown. **Bold**, *italics*, and `code` all work.

   > A paragraph starting with "> " renders as a highlighted callout — use it for the one or two
   > moments in the post that matter most (a key finding, a "we verified this ourselves" moment).
   ```

   - `status` drives the chip color: `pass` / `fixed` / `shipped` → green, `open` → rust. Add a new
     status word in `CHIP_CLASS` in `scripts/build_blog.py` if you need a third color.
   - `label` is the short scorecard/index text (a few words), not the full title.
   - Entries render in filename order — that's the only thing that controls sequence.

2. Add one line to `blog/INDEX.md`'s entry list (for humans browsing the repo; not used by the build).

3. Run the build:

   ```bash
   python3 scripts/build_blog.py
   ```

   This writes `blog/dist/index.html` — a single self-contained file, font embedded, no external
   requests. Open it directly in a browser to check it before publishing.

4. Publish it wherever the blog is currently hosted. If hosting changes later, this step is the only
   thing that changes; steps 1–3 stay the same regardless of where it's published.

   **As of 2026-07-23, live at `https://audos.merkhetventures.com`** (Cloudflare Worker, Merkhet
   Ventures account, Custom Domain on the `merkhetventures.com` zone), gated by HTTP Basic Auth
   (username `team`, password `audos` — change via `wrangler secret put BASIC_AUTH_USER` /
   `BASIC_AUTH_PASS` in `blog/cf-worker/`). The old `audos-sdk-blog.sowgood.workers.dev` URL still
   resolves (same Worker) but isn't advertised; a brief earlier stopgap on `sdk.bathala.io` has been
   retired. To republish after a build: `blog/cf-worker/redeploy.sh` — rebuilds and redeploys in one
   step. See `docs/platform/27-blog-hosting-cloudflare-worker.md` for the full setup, including the
   caching gotcha that made the first deploy silently leak content unauthenticated, and the
   `bathala.io` → `merkhetventures.com` migration history.

   Superseded: a Claude Artifact was used earlier in the project; no longer the canonical copy.

## Add a new bug

Write `blog/bugs/NNNN-your-slug.md` (one number higher than the last entry), with this frontmatter:

```markdown
---
date: 2026-07-20
area: db-api
status: open
filed: no
label: Short bug title
---

Body paragraphs, same as a blog post. Use a `>` callout for the verbatim technical detail (exact error
strings, byte counts, whatever proves the finding rather than just asserting it).
```

- `status`: `open` or `fixed` (only tracks whether *we* consider it resolved — see `filed`/`audos_status`
  below for whether Audos has actually done anything about it).
- `filed`: `yes` or `no` — have we actually submitted this to Audos (Priority Support, a filed ticket,
  etc.), not just written it down internally.
- `filed_ref` (optional): a path or reference to where it was filed, if `filed: yes`.
- `audos_status` (optional, only meaningful when `filed: yes`): tracks what Audos has actually *done*
  with a filed bug, independent of our own `status`. One of:
  - `pending` (default if `filed: yes` and this is omitted) — submitted, no response yet.
  - `acknowledged` — Audos has confirmed receipt / is looking at it.
  - `fixed` — Audos resolved it on their end.
  - `wontfix` — Audos looked at it and isn't planning to change it.

  This is deliberately separate from `status`: a bug can be `status: open` (still broken for us) and
  `audos_status: acknowledged` (they know, haven't shipped a fix yet) at the same time — don't conflate
  "we've closed the loop on our workaround" with "Audos has actually fixed the platform."

## Add a new feature request

Write `blog/feature-requests/NNNN-your-slug.md` (one number higher than the last entry), with this
frontmatter:

```markdown
---
date: 2026-07-20
priority: 2
status: not filed
label: Short feature title
---

Body paragraphs, same as a blog post. Explain what's missing and why it matters — distinct from a bug,
nothing here is "broken," it's unbuilt.
```

- `priority`: `1` (highest) through `3` — how much this would actually change our day-to-day building.
- `status`: `not filed` or `filed` — whether this has actually been raised with Audos, not just written
  down internally.

## Add a new experiment

Write `blog/experiments/NNNN-your-slug.md` (one number higher than the last entry), with this
frontmatter:

```markdown
---
date: 2026-07-20
area: db-api
status: confirmed
label: Short experiment title
---

Body paragraphs, same as a blog post. Structure it as **Hypothesis:**, **Method:**, **Result:** —
bolded lead-ins within paragraphs, same pattern used throughout the bug write-ups. State what we
expected before running it, what we actually did, and what came back — including the honest version if
the first pass got it wrong and needed a re-test.
```

An experiment is a deliberate, hypothesis-driven test of a specific platform capability — distinct from
a bug (something broken, unplanned, discovered) and a feature request (something missing, never built).
If the entry doesn't have a falsifiable "we expected X, we got Y" shape, it's probably a bug or a blog
post instead.

- `status`: one of
  - `confirmed` — hypothesis tested and the result held up under independent verification.
  - `corrected` — the first pass got it wrong (a bad test, an unverified job self-report, a wrong
    assumption) and a later, more rigorous re-test corrected it. Keep the original wrong conclusion
    visible in the body, don't just delete it — the correction is part of the finding.
  - `inconclusive` — run, but the result doesn't clearly answer the hypothesis either way (e.g. the
    evidence needed to confirm it was lost before it could be checked).
  - `open` — designed but not yet run, or blocked partway through.

## Editing an existing entry

Just edit the `.md` file and re-run `python3 scripts/build_blog.py`. Frontmatter fields (`status`,
`label`, `date`, `product`) all update the page automatically.

## Why Markdown → script → static HTML, not a CMS

Every alternative we considered (Notion, GitBook, ClickUp) meant a second copy of the content that
needed manual re-syncing on every post. This way the Markdown in `blog/` is the only copy that ever
exists — the generated page is a pure function of it. Any agent that can write a Markdown file and run
`python3` can update the blog correctly, without knowing anything about the page's HTML/CSS.
