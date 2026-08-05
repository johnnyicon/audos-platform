# Audos Architecture Audit — Is it built for hardcore agentic app development?

*Red-team architectural audit, 2026-07-17. Synthesis across ~two weeks of hands-on building (DoKnow,
field-notes) plus the earlier Throughline-era findings preserved in this repo. Every claim below is
grounded in a cited finding already in `audos-platform` — this report does not introduce new tests, it
reads across the existing ones to find the small number of root causes behind them.*

---

## 1. Executive verdict

> **Scope, timing, and an important migration caveat.** This audit evaluates the build path used
> throughout DoKnow's and field-notes' development over roughly two weeks: the Otto chat interface plus the
> **Cursor-backed** Task dispatch that does the actual building. Audos ships fast — its own capability doc
> opens by warning that "these limits move" — so every finding here is a point-in-time observation, not a
> permanent verdict, and a few older reference points it draws on are flagged inline as unverified this
> cycle.
>
> **The migration caveat matters most.** Audos also ships a newer editing surface, **Audos Code** (Beta
> 0.3.0, `audos.com/portfolio/code`), which we drove live for ~25 minutes. It runs on a **different backend
> entirely — Anthropic's Claude Code, not Cursor** (its model visualizer routes through `claude-opus-4-8`;
> its picker lists Opus 4.8 and Fable 5, both labeled "Claude Code / Anthropic"). Audos appears to be
> **mid-migration** from the Cursor-backed pipeline this audit analyzed to a Claude-Code-backed one. Several
> findings below are therefore **legacy-path-specific and may not apply to where Audos is heading** — we
> flag each inline: the report-truncation root cause (J-2, Cursor's stream), the shared-Cursor-account
> `usage_limit_exceeded` (Gap 1 / `bugs/0006`), and the `cursorModel`-pinning method (`experiments/0004`).
> We did not fully evaluate Audos Code — a complete pass on it is worth doing before treating this as the
> whole picture — but where it materially changes a finding (Gaps 1, 2, 9), we say so.

**As currently architected — and we think this is worth stating directly rather than hedging — Audos does
not yet support hardcore agentic app development the way a terminal-attached coding agent (Claude Code,
Codex) does.** But the honest verdict has a sharp, load-bearing qualification, and the qualification is the
most useful thing in this report: **the gap is not in what Audos can *run* — it's in how Audos is *built
on*.**

Those are two different questions, and the evidence separates them cleanly. On the "what can it run"
question, Audos does far better than its reputation: the team's real frontend stack survives the build
pipeline (react-query, GSAP, Radix, three.js all render — `experiments/0001`), a hand-authored UI ports
in verbatim at high fidelity (`experiments/0003`, `experiments/0006`), native OpenAI embedding
*generation* is available with no API key through an undocumented gateway (`experiments/0012`), and a
hand-rolled brute-force similarity scan over JSON-stored vectors is fast enough at DoKnow's small scale to
stand in for the vector search the platform doesn't actually have (`experiments/0013`, `bugs/0015`), and a
full-screen app genuinely escapes the chat shell from the first paint (`experiments/0017`). A real product — DoKnow — is live on
it. The runtime is more capable than the docs claim, not less.

The failure is in the **development loop**, and it is structural, not a pile of bugs. Every recurring
shape of pain in this repo traces back to a single architectural stance: **the person building — or the
agent building on their behalf — is kept one full indirection away from the actual system.** You don't
touch the filesystem; you ask Otto, who dispatches a Cursor background agent, which acts inside a sandbox,
whose success is then adjudicated by an Audos-side publish/gate step the agent can't see, and you observe
the result through partial APIs that sometimes report the wrong thing with a 200 OK. At every layer we
touched, the authoritative truth tended to live somewhere you couldn't directly read. That is the near-
inverse of a
terminal-attached coding agent (Claude Code, Codex), whose entire power comes from the *absence* of that
indirection: it edits real files under version control, runs a real shell, streams full unredacted output
as it happens, inspects live state directly, and never has to guess whether a change was saved because it
can just look.

This is why the repo's own operating discipline had to become "never trust a job's self-report — always
independently verify the live result" (`bugs/0007`, `docs/platform/19`), a rule that by the team's own
measurement "roughly doubles the cost of every fix cycle." That doubled cost is not a quirk. It is the
tax the indirection imposes, paid on essentially every change we made, until the indirection is removed. A
terminal agent pays it far less often because it rarely loses sight of ground truth in the first place.

So the qualified verdict, earned: **Audos is a genuinely good no-code platform for a founder whose agent
operates the platform conversationally on their behalf. It is not yet a substrate a developer or coding
agent can operate directly, and in our two weeks no amount of prompt discipline closed that gap — it looked
structural, not a matter of usage.** The rest of this report is the 80/20: the handful of root causes that
produce most of the pain, ranked by how much real capability their fix would unlock.

---

## 2. Methodology

**What was reviewed.** The full `audos-platform` evidence corpus: 19 hypothesis-driven experiments
(`blog/experiments/0001–0019`), 23 confirmed bugs (`blog/bugs/0001–0023`), 9 feature requests
(`blog/feature-requests/0001–0009`), 13 narrative blog entries (`blog/0001–0013`), the capability and
operational-model reference docs (`docs/platform/01–25`), the terse issue index (`BACKLOG.md`), the
change history (`CHANGELOG.md`), the field-notes build log (`field-notes/ACTIVITY-LOG.md`), and three
dated raw DoKnow activity-log snapshots (`~/doknow-kb/audos/archive/*`).

**How much evidence exists, and how good it is.** This is an unusually strong corpus to audit against.
The experiments are written as falsifiable claims with method and result; the bugs carry reproduction
detail (md5 hashes, exact error strings, job IDs); and — most valuably for an audit — the record contains
multiple *honest self-corrections* that are cited here as strengths, not hidden: a "verified" course-data
fix that failed on live re-test (`docs/platform/19`), an embedding capability first reported absent then
found present (`experiments/0011` → `0012`), a file-durability data-loss bug filed with Audos then
retracted the same day (`bugs/0022`), and — notably — that retraction itself later flagged as
under-verified because the intermediate state was never independently checked and is now unrecoverable
(`field-notes/ACTIVITY-LOG.md`, 21:10 entry). A corpus that catches itself twice recursively is one you
can lean on.

**Real limits on what can be concluded.** Several findings are genuinely open and are treated as such
below: the `db.rawQuery` leading-keyword bypass (`bugs/0021`) was never confirmed to actually execute a
blocked write and must not be read as a proven exploit; the scheduler is "inconsistent," not "broken" —
one hourly schedule never fired, one daily schedule fired within 9 seconds, root cause unknown
(`experiments/0009`); the cross-app "cards due" discrepancy (`BACKLOG.md #9`) is an un-root-caused
*possible* isolation bug on thin evidence; and the corpus spans more than one workspace and product, so a
handful of Throughline-era observations (2026-03/04) corroborate rather than directly witness the DoKnow
findings. Where evidence is thin or contradictory, this report says so rather than rounding up.

**One framing decision worth stating.** This audit ranks by *leverage* — which fixes unlock the most real
capability — not by how many findings or how much narrative each topic generated. That deliberately
*demotes* the chat-shell cluster (bugs `0001–0004`, by far the most-written-about topic in the repo) to
mid-list, because the evidence shows it is effectively **solved by inheritance**: the shell fixes live in
a workspace-shared `Desktop.tsx`, so once fixed they propagate to every app for free
(`experiments/0006`, `experiments/0017`). And it *promotes* the data/schema layer, which has fewer
findings but harder, unworked-around ones. Counting bugs would invert both of those; leverage is the
right axis.

---

## 3. Top 10 architectural gaps, ranked by leverage

> The **top 3** are flagged as most critical. Gaps 1 and 2 are two faces of one meta-cause — no direct
> channel to ground truth — but are split because their fixes are different. Everything below gap 3 is
> real but either has a working workaround, is lower-frequency, or is downstream of the top 3.

### 🔴 Gap 1 — You can't trust the build agent's done-signal — and in the legacy path, can't watch or stop it either

**The highest-leverage problem we found — and the one part of it that survives Audos's backend migration.**
Every structural change — a table, a hook, an app edit — is a natural-language job handed to Otto, which
(on the Cursor-backed path) dispatches a background agent and returns almost nothing until it's over. The
newer Audos Code surface fixes much of the *visibility* half of this (below); what it does **not** fix — and
what we reproduced live in that newer surface the same day — is the *trust* half: a confident "done" that
isn't true.

*Pattern of evidence, across many unrelated findings:*
- **Self-reports don't match reality.** Jobs reported "complete," "published," even explicit "verified —
  three lesson rows inserted with valid JSON arrays," and the live app reproduced the *identical* failure
  on re-test — twice, worded the same both times (`bugs/0007`, `docs/platform/19`). Otto's own read: an
  **authority-boundary gap** — the agent sincerely reports what it did in its sandbox, but the steps that
  decide ground truth (confirmation gates, publish/recompile) are enforced Audos-side, outside the
  agent's visibility (`docs/platform/20`, Priority 1).
- **Reports truncate at the critical line.** Completion summaries cut off mid-sentence at least eight
  times, "almost always at the exact critical detail — the root-cause line, the fix diff, the publish
  line" (`bugs/0008`). Root cause is upstream (Cursor's Cloud Agents API truncates its own stream), but
  the *choice to build the dev loop on that black-box backend is itself the architectural decision* —
  Audos could persist and expose the full transcript server-side and doesn't (`feature-requests/0002`).
  *(Legacy-path-specific: this attribution is to Cursor; Audos Code runs on a different backend — see the
  migration note below — so it may not carry forward.)*
- **No progress, no cancel, no hang-detection.** A job sat "Running" for 20+ minutes across 12 polls with
  no way to see inside or abort it; it turned out not to be hung — it ran 45m01s doing 67+ real actions —
  but "a real hang is indistinguishable from slow-but-working until it finishes or fails" (`bugs/0018`,
  `BACKLOG.md #10`). No `list_jobs` outcome field, no `get_hook_logs` coverage, nothing.
- **Silent serialization and shared-account failures.** A second job against the same app scope queued
  invisibly behind the first with no queue-position signal (`bugs/0018`); jobs instant-fail with
  `usage_limit_exceeded` against *Audos's own* shared Cursor account, unclearable by the workspace owner
  (`bugs/0006`). *(Also legacy-path-specific: in Audos Code, Fable 5 showed a transparent "0% of daily
  budget used" quota instead of the opaque shared-account failure — this one looks already-addressed on the
  newer backend.)*
- **The entire dispatch-then-poll SOP** (`docs/platform/23`) exists largely because this path has no
  streaming: fire in the background, poll every 90s for up to 30 minutes, then verify live anyway because
  "Complete" isn't evidence of correctness.

**A newer surface fixes the visibility half — on a different backend — but not the trust half. Verified
live, 2026-07-17 (~25 min).** Audos Code is a genuine step forward on *observability*, and deserves credit:
a Settings toggle ("Assistant output — show token-by-token output while a response is in progress") is **on
by default**, and during our edit the left panel streamed a live, typed step-log — `Command run`, `Tool
call`, `File change`, a running duration counter, then a specific completion line. That falsifies "no live
output" as a blanket claim; there is a real live channel here. It also runs on a **different backend**: the
model visualizer routed through `claude-opus-4-8`, and the picker offered Opus 4.8 and Fable 5, both
labeled "Claude Code / Anthropic" — i.e. Audos is moving off Cursor onto Anthropic's Claude Code. **But the
*trust* half was unchanged, which is the point:** asked to change a headline's period to an exclamation
point, it streamed real steps, ran 29s, and reported, specifically and confidently, `"Done. The headline
now reads 'Stop saving. Start knowing!' — the period is now an exclamation point. Nothing else was
touched."`, with an "Up to date" status. The rendered preview still showed a period — confirmed by zoom,
and again after a manual refresh. Same failure mode as `bugs/0007` — reproduced on the *new* backend, in
the *more observable* UI. That's what makes it the strongest evidence in this report: the self-report/
reality mismatch is not an artifact of Cursor or of poor tooling; it survived both the backend swap and the
observability upgrade. *(Maturity note, in the same honest spirit as the rest of this corpus: Audos Code is
young — across two sessions totaling ~40 minutes on 2026-07-17 it failed to load four times, "Audos Code
could not be loaded. Try again" or a stuck traffic spinner — a reproducible reliability signal, not a
one-off.)*

*Root cause:* in the dispatch path the agent's output is a delayed, truncated, unstreamable summary of work
done in a sandbox whose success the agent itself can't fully confirm — you can't watch, steer, or stop it.
The Audos Code surface adds live visibility during the run, but (as above) not yet a done-signal you can
trust without checking the rendered result yourself.

*What "fixed" looks like:* live-streamed job output; a cancel/abort control; full untruncated transcripts
retrievable by the same channel that launched the job; and — the deepest fix — **surfacing the
platform-side gate/publish state back to the job**, so a completion report can honestly say "my sandbox
change survived publish" instead of reporting sandbox-success as ground truth. The repo's own #1 request
(`feature-requests/0001`) — a mandatory post-build smoke check before a job may say "verified" — is the
minimum viable version.

*How a terminal agent avoids this by construction:* it **is** the process. Every command's output streams
live, Ctrl-C cancels instantly, nothing is truncated, and "did it work" is answered by running the thing
and looking — there is no sandbox/ground-truth split to bridge because there is no sandbox.

---

### 🔴 Gap 2 — No authoritative view of actual state: "built," "published," and "live" are three different things you can't tell apart

**The problem that caused the single largest self-correction in the corpus.** There is no reliable way to
ask the platform "what actually exists right now," and the most obvious API answers the wrong question.

*Pattern of evidence:*
- **The false data-loss saga.** A ~2-hour job built a full working app, then crashed before reporting.
  `GET /api/space/{id}/files` returned `{"files": []}`, read as "the work was discarded" — filed with
  Audos, escalated to a human engineer. It was wrong: **that endpoint only ever reflects the *published*
  bundle, never draft work.** The app had been sitting in Draft the whole time; it took loading the
  Preview panel's Draft-vs-Live toggle by hand to see it (`bugs/0022`, `field-notes/ACTIVITY-LOG.md`).
  "Nothing was built" and "built but unpublished" are literally indistinguishable through the API. And the
  trust problem recursed one level deeper: the retraction *itself* later proved to rest on an unverified
  job self-report — a second job's claim that the first job's work "was already there, not rebuilt," never
  independently checked and now unrecoverable to check. This is arguably the sharpest single illustration
  of Gap 1 in the corpus: ground truth was so hard to reach directly that even a careful correction ended
  up trusting an agent's word about state nobody had independently observed.
- **"Published: yes" has meant three different things** on different days: actually live; config-written
  but bundle not recompiled; or blocked because a concurrent same-scope job silently held the publish
  (`bugs/0005`, `docs/platform/19`, `feature-requests/0003`).
- **Stale served bundles.** `list_apps` shows an app registered and published while the served bundle is
  stale — app missing from the dock, route 404ing (`bugs/0005`).
- **The deploy pipeline is a 15–30 minute unobservable window** with silently-dropped webhooks: "push to
  live — 15–30 min… Occasionally a push webhook is silently dropped by Audos — no error, no indication,"
  remediable only by a manual "Sync from GitHub" click (`docs/platform/16`).
- **Cache-busting is folklore.** A manual `?_cb=` is sometimes needed after a `Desktop.tsx` publish to
  see a fix; the platform emits `?_cb=…&cdn=fallback` itself under conditions nobody documented
  (`feature-requests/0004`, `docs/platform/19`).

*Root cause:* the system's real state (draft bundle, published bundle, served bundle, config, file tree)
is scattered across surfaces that don't agree, and the one programmatic read path silently reports only
one of them.

*Partial credit, newer surface (verified live 2026-07-17):* Audos Code makes the draft/live distinction
genuinely clearer for anyone working in that UI — "Preview" shows the inline draft, while "Live" opens a
separate tab at an explicit, human-readable URL with an `?env=live` query param (`audos.com/site/{id}?env=
live#…`). That's a real improvement over the old "`GET /files` only ever shows published" ambiguity — for
the click-through path. The API-level opacity is presumably unchanged for anyone integrating
programmatically rather than through this UI, which is who the fix below is for.

*What "fixed" looks like:* a publish-status endpoint returning the currently-served bundle hash and
timestamp; a file/state API that distinguishes draft from published instead of silently collapsing them;
automatic, consistent cache invalidation on publish. Something an external agent can poll to know the true
current state without opening a browser and eyeballing a toggle.

*How a terminal agent avoids this:* `git status`, `git log`, `ls`, `cat`. State is the filesystem under
version control — there is exactly one authoritative answer to "what's here and is it saved," and it's
free to read.

---

### 🔴 Gap 3 — The data/schema layer is non-transactional, undocumented, and reachable only indirectly — and its one real escape hatch is itself broken

**The gap with the fewest findings but the least workaround headroom**, which is why it ranks in the top
3 despite the shell cluster having 4× the write-ups. Iterating on a schema — the most ordinary thing a
developer does — was, in our experience, unsafe in places and in others blocked outright.

*Pattern of evidence:*
- **Failed DDL doesn't roll back, and there's no cleanup tool.** A `CREATE TABLE` with an unsatisfiable
  foreign key failed as expected — but left a real, physical, catalog-unregistered `app_site_sessions`
  table behind. Every retry (including a corrected one) then failed with `relation already exists`, the
  orphan invisible to `db_list_tables` yet physically blocking the name, with no `drop_table` tool to
  clear it. Resolution: give up and use a different name, `site_sessions_v2` — the original name stays
  blocked until Audos clears it server-side (`bugs/0019`, `BACKLOG.md #15`).
- **The schema you ask for is not the schema you get, silently.** `id` is always a platform-forced
  `serial` integer; requesting `id: uuid` is rejected outright (`column "id" specified more than once`);
  every table silently gains `session_id` and `updated_at` columns — none of it in the column-type
  reference (`bugs/0020`, `BACKLOG.md #16`). This *caused* the orphan above: a uuid FK can't reference a
  serial PK.
- **No native vector support, and no way we could find to enable it.** `vector(N)` is rejected
  (enum-restricted types), and raw SQL is `SELECT`/`WITH`/`EXPLAIN`-only, so `CREATE EXTENSION vector`
  can't be run from the workspace — pgvector is
  *available in the catalog but not installed* (`installed_version: null`, absent from `pg_extension`),
  with no owner-accessible path to install it (`bugs/0015`; the corpus itself corrected an earlier
  "installed but unenabled" misread — it was never installed).
- **The read-only SQL guard may be shallow.** Evidence the `rawQuery` validator checks only the leading
  keyword, so `SELECT 1; DROP …` might slip through — *explicitly unconfirmed*, the probe was abandoned
  and likely didn't succeed; flagged for Audos to verify, not asserted as an exploit (`bugs/0021`).
- **A recurring, three-time serialization footgun.** Every first write to a `json`/`jsonb` column in a
  new code path failed with `invalid input syntax for type json` because the agent sent a raw object, not
  a `JSON.stringify`'d string — three times, across two unrelated apps (`docs/platform/19`).
- **The one direct escape hatch is 0-for-3 on even *reaching* it this cycle.** Direct scoped Postgres
  (full DDL, transactions, `DROP` — everything the above needs) would be exactly the right answer, and
  `docs/platform/13` reports it working *in April, in Throughline's workspace* — not re-verified this
  cycle, so "direct SQL is full-powered once you have credentials" is an inference, not a current fact. And
  the credential-generation step in front of it has now failed **three different ways, zero successes,
  across three workspaces**: Throughline (April) reportedly got "Workspace not found"; field-notes
  (2026-07-16) gets a `409 "Use regenerate to rotate them"` on every click with **no regenerate/view
  control anywhere in the UI**, `GET` → 401, and Otto confirming no view path (`bugs/0023`,
  `experiments/0019`); and DoKnow's own workspace (2026-07-17) returned a **silent `404` with no error at
  all** — worse than field-notes' visible failure. So this isn't "a working escape hatch with a broken
  door" — on current evidence it's *a documented intention nobody has reached this cycle*, which if
  anything makes the D-5 credential fix higher priority, not lower. **Scope this correctly, though:** it
  blocks the *direct psql* path only — database *writes* work fine through Audos Code's DB-write layer
  (verified live 2026-07-17; see "What's actually working well"). So the fix restores a *convenience and a
  standard interface*, not the only way to get data in.

*Root cause:* schema operations run through a non-transactional, undocumented agent-mediated API that
mutates structure in ways it doesn't disclose and can't undo — and the direct-database path that would
sidestep all of it is gated behind a credential flow that doesn't survive first contact.

*What "fixed" looks like:* atomic DDL (roll back partial `CREATE TABLE` on constraint failure) or, at
minimum, a `drop_table`/cleanup tool; documented and honored PK/column behavior; and — likely the highest
leverage of the three — **fix the one-shot credential recovery**. That last one looks like a small
UI-wiring fix (the backend's own 409 references a regenerate action) that would hand developers a real
`psql` connection. One caveat we want to be honest about: the claim that direct SQL is then *full-powered*
rests on the April Throughline observation above, not a this-cycle re-verification — so this fix's payoff
is our best-supported inference, not a proven one. If it holds, most of the orphaned-table, serial-id,
vector-enablement, and raw-SQL friction eases the moment direct, transactional database access is reliably
available.

*How a terminal agent avoids this:* it runs migrations inside a transaction that rolls back on failure,
`DROP`s a bad table, and connects with `psql` — schema iteration is safe, reversible, and direct by
default.

---

### Gap 4 — An imposed chat-shell UI, escaped only by a manual checklist, over a workspace-shared mutable file

The default signed-in surface is a ChatGPT-style shell; a normal full-screen app requires actively
escaping it via a 4-step checklist — full-screen landing config, deep-link resolution, a
popup-not-teardown chat affordance, and a lazy-`useState` fix for a flash-of-old-shell FOUC (`bugs/0001–
0004`, `docs/platform/22`). **Lower leverage than it looks, for two reasons the evidence establishes.**
First, it's largely *solved by inheritance*: the fixes live in a workspace-shared `Desktop.tsx`, so once
fixed they propagate to every app for free — a zero-briefing baseline app came out clean
(`experiments/0006`), and a genuinely fresh workspace was born full-screen from a single upfront briefing
(`experiments/0017`). Second, that same sharing is the real residual risk: `Desktop.tsx` is
**workspace-shared and mutable**, so a shell change has blast radius across every app, and this
*confounded one of the team's own experiments* — "One Good Thing" looked like a clean upfront-briefing
success but was three-quarters inherited from an already-fixed shell (`experiments/0005`). App-level
isolation boundaries are genuinely murky here (the un-root-caused cross-app "cards due" discrepancy,
`BACKLOG.md #9`, is a thin but live hint). *Fixed looks like:* a first-class "full-screen app" mode that
needs no shell surgery, and per-app isolation so one app's shell can't be another's dependency.

### Gap 5 — No reliable unattended/background execution primitive

Recurring scheduled hooks fire inconsistently: two hourly schedules never fired (2h07m and 20m overdue,
`runCount: 0`), while a later daily schedule fired within ~9s — same mechanism, opposite result, root
cause unknown and unknown to Otto too (`bugs/0016`, `experiments/0009`). One-time hook scheduling doesn't
exist at all. There's no dispatcher log, queue, or status to diagnose from. **Partially de-risked** by a
proven workaround — client-orchestrated sequential hook calls, 5/5 clean (`experiments/0010`) — which
covers anything that can keep a client open, so this ranks below the top 3. But genuinely unattended,
fire-and-forget work has no reliable path, and combined with the ~5–10 minute hook ceiling, anything
slow-and-headless must run off-platform. *Fixed looks like:* a scheduler with observable per-run status
and a root-caused firing guarantee.

### Gap 6 — A constrained, non-standard runtime, with no opt-out we could find

Single file per app, server-side ESBuild, CDN/importmap dependency resolution — no `npm install`, no
Vite, React 18 only, hash-based in-memory routing only (no real browser URL ownership), and hooks run in
a ~30s / ~64KB sandbox with no `fs`, `require`, `Buffer`, `setTimeout`, or `URLSearchParams`. *(Several of
these specifics are sourced from March–April reference docs — `docs/platform/07`, `08`, `10`, all
2026-03-31, and `16`, 2026-04-22 — and weren't re-verified line-by-line this cycle; the single-file model
and the dual-React trap below were corroborated by observed behavior.)* The importmap has a genuinely
dangerous edge — user-added packages don't inherit the React 18 pin, silently loading two React copies and
blanking the page (React error #31), catchable only by inspecting network resources (`docs/platform/16`).
Much of this is livable (the library probe shows the real stack mostly works, `experiments/0001`), which
is why it's mid-list — but we found no opt-out, and `feature-requests/0008` asks precisely for an opt-in
modern stack. *Fixed looks like:* an opt-in lane for current React, real dependency resolution, and real
URL ownership, without disturbing the no-code default.

### Gap 7 — Undocumented surface forces discovery-by-trial-and-error, and the probes give false negatives

Core capabilities are undocumented and discoverable only by guessing the right call shape. The working
embeddings path — `platform.integrations.proxy('openai', '/v1/embeddings', …)` — appears in no docs; its
provider allowlist was "read from the error message"; and the first probe **reported it absent** because
it swept `isAvailable()` with guessed feature-names instead of trying the bare provider name and the
proxy directly (`experiments/0011` → `0012`). The `cdnDependencies` allowlist can't be read even by Otto
(`search_platform_code` → "no searchable directories") and had to be reverse-engineered with a throwaway
probe app (`feature-requests/0006`, `0009`). This gap is an **amplifier of every other gap**: because you
can't read the platform's own code and can't trust self-reports, in this evaluation everything had to be
probed empirically — and empirical probes of an undocumented surface produce confident false negatives. *Fixed looks like:*
documenting `platform.integrations`/`proxy()`, the globals, and the `cdnDependencies` allowlist; and
read access to your own workspace's platform files without a full job dispatch.

### Gap 8 — Silent-success failure modes: 200 OK with the wrong content

A distinct, dangerous class: operations that fail while reporting success. `/api/upload/image` returns
`HTTP 200 {"success":true}` with a plausible `.png` URL while **corrupting** any non-image payload
(md5-verified: 3,031 bytes in, 3,050 bytes of bit-shifted garbage out) — a real PDF silently mangled
(`bugs/0017`). `proxy()` called with an options *object* instead of a JSON *string* body silently returns
the platform's own HTML index page with a 200 (`experiments/0012`). "Published: yes" that isn't live
(Gap 2) is the same disease. The through-line: **a success response on this platform has repeatedly not
meant success — treat a 200 as unverified until you check the content** — which is the doubled-verification
tax again, in miniature.
*Fixed looks like:* honest status codes and content-type validation, so a 200 means what it says.

### Gap 9 — No first-class programmatic build/write path — only an agent or a slow git pipeline

There is no *raw* file-level editing in the app layer — you can't hand-type code into a file. The
March-era docs put it flatly ("Can't edit React code directly, must use Otto" — `docs/platform/02`, `03`,
both 2026-03-31), and the current nuance sharpens rather than overturns it: Audos Code's Platform mode is
labeled "Full in-platform editing," but driving it live (2026-07-17) it's still a *scoped
natural-language-to-agent* request — select an element, describe the change, an agent executes it — not a
hand-editable code surface. So the narrow point holds (no raw file editing), while a blunt "no editing at
all" would now be wrong. The entire 8-endpoint API suite (`db-api`, `ai-api`, …) is *user-hand-authored*
from templates — the platform ships no first-class external API (`docs/platform/07`, `09`). The only
non-agent write path is GitHub Dev Mode, which is one-way and 15–30 minutes per iteration
(`docs/platform/16`). And when Cursor is over its shared limit, we couldn't reach the alternative backend
(Audos Code) from the external API — it needs a signed-in browser session, leaving an API-driven agent
with no build path in that window (`feature-requests/0007`). This is partly a facet of Gap 1, but
called out separately because the fix is different: *fixed looks like* a genuine "import this exact code →
deterministic wrap → publish" path (`feature-requests/0005`) and an API-reachable build backend, so a
change can be made and shipped without a conversational or 30-minute round trip.

### Gap 10 — Insecure/opaque identity and config defaults, changeable only by side-channel

EmailGate — the identity layer — is open by default, and we **re-confirmed this live in DoKnow on
2026-07-17**: on the signed-out live site, typing a completely fabricated, never-used address
(`audit-verify-2026-07-17@example.com`) signed straight in — no OTP, no verification step, no delay,
landing on a real dashboard. Anyone who knows (or guesses the format of) an email can enter as that user.
OTP exists but is off by default and **exposed through no settings UI** — "the only interface is the API
directly" (a full settings-surface scan on 2026-07-17 found no OTP toggle anywhere), enabled by pasting a
`fetch()` into a logged-in browser console. There is **no built-in sign-out primitive** (the avatar menu
does nothing; you implement sign-out by deleting localStorage keys) and no OAuth option on the sign-in
modal — both also re-confirmed live 2026-07-17. The one sub-claim we could **not** quickly re-verify is
whether sessions ever expire — that would require waiting out an unknown TTL, so we leave it sourced to
`docs/platform/17` (2026-04-23) rather than assert it. It's lowest on the list because it has a known (if
ugly) path and didn't block building — but a security-relevant default changeable only by a DevTools
incantation is the kind of thing that ships insecure. *Fixed looks like:* secure defaults and first-class UI/API config
for identity, not a browser-console side-channel.

---

## 4. What's actually working well

An audit that only lists failures isn't credible, and several capability probes came back genuinely
positive — some surprisingly so. These are load-bearing for the verdict: they're *why* the failure is
correctly located in the development loop and not the runtime.

- **The real frontend stack survives the pipeline.** react-query, GSAP, Radix, and raw three.js all
  render cleanly via `cdnDependencies`; react-three-fiber works but is slow to warm up (`experiments/0001`,
  `docs/platform/19`). The build model doesn't quietly rule out the team's actual libraries.
- **Verbatim design push is high-fidelity and repeatable.** A 903-line self-contained mockup ported
  element-for-element into a live app, then a second time combined with real DB wiring in place, both
  independently verified by hand (`experiments/0003`, `experiments/0006`, `docs/platform/22`). This
  sidesteps the "every Audos app looks the same" ceiling entirely — the design fidelity is yours.
- **Native embedding *generation* with no key management.** `platform.integrations.proxy('openai',
  '/v1/embeddings', …)` returns real 1,536-float OpenAI embeddings with Audos holding and injecting the
  credential — a generic authenticated passthrough to an `openai/stripe/twilio/heygen` allowlist
  (`experiments/0012`). This is *generation only*, and genuinely better than "no capability, build your
  own" — but it is not search, and must not be confused with it (see the next bullet).
- **The brute-force similarity *workaround* is cheap at small scale — but there is no vector search on
  Audos, and this report earlier implied otherwise.** To be unambiguous: there is no vector index and no
  pgvector (Gap 3) — the only mechanism is a cosine loop in hook JavaScript over embeddings stored as JSON
  arrays. Its one honest positive is that the arithmetic is cheap: 1,536-dim cosine over 1,000 rows scans
  in 1–3ms, dominated by DB fetch not math (`experiments/0013`), so the workaround is viable at DoKnow's
  real scale of tens-to-hundreds of chunks per course. What that benchmark did **not** establish — and
  what "vector search works" would wrongly imply — is that retrieval *quality* is good or that it scales
  past small row counts; `experiments/0008` explicitly walked back exactly that overreach, and `bugs/0015`
  calls the JSON fallback "not a real retrieval index." It earns a place here only because the workaround
  genuinely works at the scale DoKnow needs — not because the capability exists.
- **The shell escape genuinely works, and inherits.** A fresh workspace was born full-screen from the
  first paint on a single upfront briefing; the fix propagates workspace-wide via shared `Desktop.tsx` so
  later apps get it free (`experiments/0006`). Worth stating exactly how this was confirmed, because it
  doubles as a Gap-1 cautionary tale: the *first* attempt's own "all 5 checks passed live" claim was lost
  when the job crashed before it could be independently checked, and the eventual "confirmed" verdict came
  only from a *second, hands-on* verification — loading Draft-vs-Live by hand, logging in with test
  credentials, clicking the chat affordance personally — not from trusting any job's word
  (`experiments/0017`).
- **Direct SQL is the right escape hatch *by design* — but we could not reach it this cycle (0-for-3), so
  treat it as promise, not proof.** Read/write/DDL/`DROP` were reported working against the workspace
  schema in April, with the platform "will not overwrite your schema" (no surprise auto-migrations) — but
  that's `docs/platform/13`, **Throughline's April workspace, not re-verified this cycle**, and every
  attempt to actually reach it this cycle failed: field-notes 409s forever, DoKnow's workspace returned a
  silent 404, Throughline (April) reportedly got "Workspace not found." Three workspaces, three different
  failures, zero successes on the credential step alone. So it earns a place here only as the *correctly
  shaped* escape hatch the platform clearly intends — not as something we saw working. "Full-powered once
  you're in" can't be asserted with any current confidence, which is exactly why fixing the credential flow
  (D-5) is the highest-leverage database fix. **Important scoping, though:** the 0-for-3 is about the
  *direct psql connection* specifically — **not** about getting data into the database at all (see the next
  bullet).
- **Writing content to the database *does* work — via Audos Code, verified live 2026-07-17.** This is the
  correction that keeps the data-layer story honest: the credential bug blocks the direct-SQL escape hatch,
  but **not** database writes in general. Using Audos Code's chat-driven DB-write layer (`db.insert` via
  its MCP tool calls, on the new Claude Code backend), we populated the field-notes blog site's
  `content_items` table from empty — row count **0 → 6** with real corpus content — and verified it two
  independent ways, not on the agent's word: the workspace's own Preview feed rendered all six items
  correctly (Field Log / Reference Docs / Tracker, right chips and dates), and the live site's embedded
  "Field Scout" assistant answered "show me the latest bugs" with the exact two bugs just inserted —
  **visible on Live immediately, no publish step**. Two honest footnotes: a long free-text paste didn't
  reach the agent (it flagged that rather than fabricating content — a good honesty signal — and compact
  single-line JSON went through cleanly), and only 6 of ~64 corpus items were populated in this
  deliberately small proof-of-path batch (the technique scales, but a full sync is worth scripting).
- **Client-orchestrated hook chaining is clean.** 5/5 sequential calls, ~1–1.8s each, zero rate-limiting,
  auth, or concurrency friction (`experiments/0010`) — the runtime is stable under sequential load and
  gives back the scheduler as a hard blocker for attended work.
- **`/api/upload/file` is byte-faithful** (md5-identical round-trips for any content type up to ~50 MiB)
  and **`platform.generateText` works** as advertised once the real signature was found (`docs/platform/19`,
  `06`). The primitives are solid; it's the *loop around them* that isn't.
- **The escalation path exists and is honest.** Priority Support routes to engineering with a transparent
  automated-triage-then-human tier, seconds to first reply (`docs/platform/25`); and Otto itself refused
  to misuse the paid VA queue for a bug and self-corrected — a real positive on agent judgment.

The pattern across all of these: **once you know the right call shape, the runtime delivers.** The cost is
entirely in *finding* the right call shape and *confirming* it worked — which is the top-3 territory,
not the runtime.

---

## 5. Closing verdict

**Restated plainly: as we found it this cycle, Audos does not yet support hardcore agentic app development
the way a terminal-attached agent does — but it can host the apps that development would produce, and that
distinction is the whole story.** The runtime is more capable than its own docs; a real product is live on
it; libraries, embedding generation, verbatim design, and the shell escape all work. What we kept hitting
was the loop — Audos keeps the builder an indirection away from ground truth at most layers, and a coding
agent's effectiveness comes largely from *not* being kept at that distance. The tax that indirection
imposes (verify live, often twice; treat success responses as unverified; poll a black box; discover the
API by probing) is paid on nearly every change, and reads as structural rather than incidental.

**If Audos could fix only one thing**, the answer has two levels, and honesty requires stating both:

- **The architecturally correct fix** is to open the "power-user lane" this repo has been asking for
  (`docs/platform/20`, `feature-requests/0007`, `0008`) — a direct, authoritative channel to ground truth
  that doesn't route through a conversational agent or a 30-minute git pipeline. Concretely and cheapest:
  **fix the database-credential flow** (`bugs/0023`) — which failed three different ways across three
  workspaces this cycle (409, silent 404, "Workspace not found"), zero successes. Getting developers a
  reliable connection would hand them the one direct path (real transactional SQL) that, *if it's as
  full-powered as the April Throughline finding suggests*, would ease most of the Gap 3 cluster. Among the
  highest capability-per-effort fixes we found — with two honest asterisks: its payoff depends on that
  not-yet-re-verified "full-powered" claim, and the credential step itself needs more than a one-line
  wiring fix given how many distinct ways it's currently failing. The 0-for-3 makes it *more* urgent, not
  less.

- **The highest-leverage fix *within the current paradigm*** — if the dispatch model stays — is a reliable
  automated post-build verification gate before a job may report "verified" (`feature-requests/0001`),
  backed by surfacing the platform-side publish/gate state back to the agent (`docs/platform/20`). The
  encouraging part: **Audos appears to already be building exactly this.** Audos Code ships an in-app
  browser-test runner, currently shown as *paused* — "we've disabled the in-app test runner while we make
  its results reliable enough to trust. It will return soon." That's the right instinct — pause rather than
  ship false confidence — and it may be the most valuable single thing to get over the line, because we
  independently reproduced the same day the exact self-report/reality mismatch it would catch (the period-
  to-exclamation-point edit in Gap 1 that reported success but never changed the page). So this isn't
  "here's a gap you missed" — it's "the thing you've already paused is, from where we sit, the
  highest-leverage thing to finish."

The difference between those two is the difference between *removing* the indirection and *making it
honest*. A terminal-attached coding agent needs neither, because it never had the indirection to begin
with. That's the bar we measured against — not as a sales pitch for either product, but because direct,
streaming, version-controlled, verifiable access to the real system is what "architected for agentic
development" means in practice. Audos today is built for something else — and built well for it — and, on
this evidence, would need to open a second lane to clear that bar. Given that it's already shipping Audos
Code and a (paused) in-app test runner, some of that second lane may already be under construction.
