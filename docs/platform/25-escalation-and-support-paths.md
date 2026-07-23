# Escalation and support paths — what actually exists, verified live

*Written 2026-07-14, discovered while trying to escalate a genuinely stuck job (`BACKLOG.md #10`). We
didn't know either of these existed until we needed them — worth a dedicated doc since the SDK had no
prior coverage of "what do you do when something is actually broken and you need a human."*

## Two distinct escalation mechanisms — do not confuse them

### `request_human_help` — a paid Human-In-The-Loop task queue, NOT a support channel

What it actually is: a way to dispatch a **real human virtual assistant** to do one-off manual work
(vendor calls, lead research, custom copywriting, that kind of thing) — **not** a bug/support path.

- **Billed out of the workspace wallet**, default $25/hr.
- Lifecycle: task created → **"awaiting owner approval"** → the workspace owner approves it from a
  Human Help dashboard → a VA is dispatched → **wallet charged on completion**.
- Otto cannot fire-and-forget this — it has to quote a scope/budget and get owner confirmation first.
- **Its own usage policy explicitly forbids routing bugs, broken behavior, or platform support through
  it.** A general VA has no access to internal platform systems (can't kill a stuck job, can't inspect
  billing) — so even ignoring the policy, it wouldn't accomplish anything for a platform bug.

We asked Otto to use this for a stuck-job escalation; it correctly refused and self-corrected to the
right tool below, explaining exactly why. Worth noting as a positive finding: Otto caught a wrong
instruction from us rather than complying with a request that violated the tool's own stated policy.

### Priority Support — the actual bug/feature escalation path

This is what actually gets you to a human who can act on a platform problem. Confirmed live, this
workspace has **Priority Support** enabled (skips the normal triage queue, "your bug reports go directly
to our engineering team for same-day review").

**Where it lives:** the **Help** panel in the workspace UI (bottom-right "Help" button on the workspace
dashboard, not the same surface as the Otto chat window). It has two request types:

- **Feature** — a title + free-text description, goes into a general request queue.
- **Bug** — a richer form: one-line summary, optional title, "what did you expect," "what happened
  instead," steps to reproduce, and an optional screenshot attachment.

**Confirmed behavior, live:**
- Submitting posts to a **"Your Requests"** list, each shown with a status chip (`In Progress`
  confirmed observed) and an unread-reply indicator.
- **There's a tiered response**, not straight to a human: our submission got an automated first pass,
  which decided the issue needed more than it could handle and explicitly said so — *"This one needs a
  bit more than the automated fix could handle — we're pulling in an engineer. You'll hear back
  shortly"* — attributed to "Audo Support." So there's some automated triage/fix-attempt layer before
  human engineering gets involved, and it's transparent about handing off when it can't resolve
  something itself.
- Response time from submission to that first (automated) reply was on the order of **seconds**, not
  minutes.

**Important gap found in the same investigation:** Otto, working via the **external chat API** (not the
in-browser session), *also* has a `prepare_feature_request` (type: bug) tool that it described as
rendering "a review-and-submit card in your interface." **This did not work as described** — after Otto
reported preparing and effectively "filing" a bug report through the API-driven session, the workspace's
own Help → Priority Support panel showed **"No requests yet."** The draft never actually reached the
UI-visible queue. We had to submit the same report manually, directly in the browser, for it to actually
land. **Practical rule: don't trust an API-driven Otto session's claim that it "submitted" or "prepared" a
UI-rendered artifact — verify it actually appears where a human would see it, the same discipline this
whole SDK already applies to build-job completion claims.**

## Practical guidance

- If something is broken on the platform (not "I want a VA to do X"), use **Priority Support → Bug**,
  submitted directly in the browser UI — don't rely on an API-driven Otto session to file it for you
  without checking it actually landed.
- `request_human_help` is for outsourcing real-world manual tasks with a budget you've approved, never
  for platform bugs — Otto itself will refuse if asked to misuse it this way, per its own tool policy.
- Priority Support isn't guaranteed for every workspace — we don't know what tier/plan grants it; treat
  its presence here as workspace-specific until confirmed otherwise.
