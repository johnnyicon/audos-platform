---
date: 2026-07-15
product: DoKnow
status: pass
label: What Audos can actually do
---

# What Audos can actually do

Every entry so far has been about the shell — getting Audos to render a real product instead of a
chatbot with apps on the side. That question is answered now, on both ends: build clean from the start,
or migrate an existing app out in place. What's never been tested, until today, is the layer underneath
the shell: can the platform actually run the product DoKnow is supposed to be? Gather content, transcribe
it, embed it, retrieve from that index to build a course. Nobody had checked. So we checked.

The method was the same one this whole build has leaned on by now: don't trust a claim, run it for real.
We wrote five small server-side hooks — genuinely disposable, nothing touching a real app — and executed
each one against the live platform: create a vector column and see what happens, fetch an external URL
and see what comes back, upload a file and find the actual size limit by bracketing up until it breaks,
schedule a task and watch the clock to see if it fires, call the built-in AI generation and read the raw
response. Five direct questions, five direct answers, no reading between the lines.

Three of the five came back solid. File upload works — bracketed the real limit at exactly 50 MiB, one
byte under succeeds, one byte over throws. External HTTP calls from a hook are unrestricted — any host,
no allowlist, so calling a transcription or embedding API from inside the platform is genuinely open.
Built-in AI generation works too, a real endpoint, real response, defaults to a small fast model exactly
where the product plan said to save cost for high-volume work.

The other two are the ones that actually decide whether this pipeline can exist here at all, and both
came back no. Vector storage doesn't exist — not "hard to set up," genuinely absent. The extension is
sitting right there on the server, `pg_available_extensions` shows it installed and ready, but there's no
way to turn it on: table creation is restricted to a fixed list of column types that don't include one
for vectors, and raw SQL is read-only, so there's no door left to walk through. The fallback is storing
an embedding as a JSON array and comparing it against every other row by hand in application code, which
works for a handful of items and stops meaning anything past that. And the scheduler — the one documented
path to running something later or in the background — creates its schedules just fine, accepts a
perfectly valid recurring rule, reports back a pending status. Then nothing happens. We watched one sit
two hours past its fire time with a run count that never left zero, on a hook we'd already confirmed
works when called directly. The schedule exists. It just doesn't do anything.

Put those two together and the ingestion pipeline as designed doesn't fit here. Not because any single
piece is broken — because the two pieces the whole design leans on, retrieval and background processing,
aren't available, and everything else is downstream of them. A course built from real, retrievable
sources instead of the model's own guess needs an index to retrieve from. A pipeline slow enough to need
transcription and embedding needs somewhere to run that isn't a single five-minute request. Neither
exists today.

One more thing came out of this pass that had nothing to do with the five tests: a genuinely dangerous
bug, found by accident while checking upload limits. The image-upload endpoint accepts anything you send
it — a PDF disguised as an image upload returns success and a plausible-looking URL. The bytes it
actually stored don't match what was sent. Not almost — completely different, no shared prefix, silently
wrong. No error, no warning, a URL that looks exactly like success. That's worse than a limit you can
bracket and plan around; it's a failure mode with no signal attached to it at all.

We also, separately, learned something about asking for help here. Escalating a genuinely stuck job led
us first to the wrong tool — a paid human-assistant queue that exists for delegating real-world tasks,
not platform bugs, and its own policy says as much. The right path turned out to be a Priority Support
form sitting in the workspace's own Help panel, and it worked exactly the way the label promised: an
automated pass looked at the report within seconds, decided it needed a person, and said so plainly
rather than leaving us guessing. Worth knowing before you need it, not while you're stuck.

So the honest state of things: the shell is solved, verified twice over now. The data layer underneath
it, for the specific thing this product needs to do, is not there yet — not "difficult," genuinely
missing two load-bearing pieces. The real decision in front of us now isn't a UI question. It's whether
DoKnow's ingestion runs as a small, synchronous, single-hook-sized version of itself on Audos, or lives
as its own service outside the platform entirely, with Audos as the front door and nothing more. That's
a question worth answering deliberately, and now we can answer it with facts instead of a guess.

> **Update, 2026-07-16.** Before writing a formal feature request off the back of this entry, we went
> back to close exactly the two gaps that mattered most. One closed cleanly. The other didn't — we
> initially wrote it up as closed too, and that was wrong; the correction is below. A full experiment
> writeup lives at [`docs/platform/reports/2026-07-16-vector-search-experiment.html`](../docs/platform/reports/2026-07-16-vector-search-experiment.html); the short version follows.

We hadn't actually measured whether "no vector index" was slow — we'd assumed it from the absence of the
tool, not from a stopwatch. So we measured it: 50, 300, and 1,000 fake embeddings, stored as JSON,
compared brute-force in hook JavaScript. Sub-millisecond at every size — 0.02ms average at 50 rows, all
the way to 0.28ms at 1,000. That's a real result: cosine similarity is fixed arithmetic, and it costs the
same whether the numbers mean something or not, so the row-count scaling genuinely holds. But it's a
narrower result than we first gave it credit for — those were fake, 5-number vectors, not real
embeddings, which typically run 384 to 1,536+ numbers. We hadn't tested the size that actually matters,
and we hadn't tested whether Audos can produce a real embedding at all. Both are still open; see the
correction below.

The scheduler finding held up better. We ran the identical test again, changing only the frequency from
hourly to daily, and this time it fired — on time, within nine seconds, hook executed, database written,
run count advanced. Same platform, same account, opposite result from two days earlier. We asked Otto
directly whether the earlier failure was a known issue; it said, plainly, that it has no record either
way. We're not calling the scheduler fixed. We're calling it inconsistent, which is its own kind of
finding — and it doesn't matter as much as it did, because we also tested whether DoKnow's own app could
just do the pacing itself, calling a hook five times in a row instead of trusting a schedule to fire on
its own. Five for five, no rate-limiting, no friction, about a second and a half a call. That works
regardless of what the scheduler does on any given day. This one really is resolved.

The vector-search one isn't. What we actually proved is that comparing arbitrary numeric arrays is fast —
not that semantic search works here. Whether it does depends on two things we haven't tested: real
embedding size, and whether Audos can generate a real embedding in the first place (no built-in
capability for that has ever turned up — `generateText` isn't `generateEmbedding`). Round 3 is a direct
test of both.

> **Update, 2026-07-16 (later the same day).** Round 3 closed both questions — and along the way, corrected
> a mistake we made *inside* round 3 itself. Worth telling honestly, because the mistake is the more
> useful finding of the two.

The dimensionality test was simple and it held up exactly as arithmetic says it should: rerun the same
brute-force comparison with real 1,536-number vectors instead of 5-number placeholders, same row counts.
Scan time barely moved — about 0.003ms per row at 300 and 1,000 rows, a full scan of 1,000 rows in one to
three milliseconds. Three hundred times more math per comparison and it's still nothing; the actual cost
is fetching the rows from the database, not comparing them once they're in memory. That question is
closed.

The embedding-generation test is the one worth slowing down for, because we got it wrong the first time
and the way we got it wrong is a real lesson about this platform. The first pass asked Audos's own runtime
"do you know about anything named `openai-embeddings`, or `vector-search`, or `pinecone`?" — a list of
guessed names — and got `false` back for every single one. We wrote that up as "no embedding capability
exists," and for about an hour, that's what this document said.

It was wrong, and it was wrong for a specific, avoidable reason: we asked the platform about a *feature*
by guessing feature-shaped names, instead of asking it what it actually *has*. What Audos has is a generic
proxy — `platform.integrations.proxy(providerName, path, options)` — that forwards an authenticated
request to a small allowlist of real providers: OpenAI, Stripe, Twilio, Heygen. Nobody built a
`generateEmbedding` method, because nobody needed to. The bare name `openai` is enough. Point that proxy
at `/v1/embeddings` with a model and some text, and it comes back with a real, 1,536-number OpenAI
embedding — and Audos supplied the API key. Not ours. Theirs, already sitting behind the platform,
available to any hook that knows to ask.

That's a better answer than the one we filed an hour earlier, and a worse discovery process than we'd
like: the capability was there the whole time, just not under a name anyone would guess, and not written
down anywhere. The lesson isn't "Audos has more than it documents," though that's also true and we've
said it before. It's narrower and more useful than that: when a capability test comes back negative,
the negative might be about the names you tried, not about the platform. We caught it this time because
round 3 had two separate sub-jobs test overlapping ground and one contradicted the other. That's not a
process we want to depend on going forward — it's worth writing down as its own finding, not just fixing
quietly.

So, plainly: both blockers this entry opened with are closed. A real embedding is one authenticated proxy
call away, no API key of our own required. Comparing embeddings at DoKnow's real scale, at real
embedding size, costs single-digit milliseconds. The ingestion pipeline — upload, transcribe, embed,
retrieve, generate — has no remaining architectural reason to live anywhere but Audos. What's left to
formally raise with Audos isn't a capability gap anymore. It's a documentation gap: `platform.integrations`
and `proxy()` are real, working, and currently discoverable only by reading the runtime's own source
through introspection — which is a strange way to have to find out your platform can already do the thing
you were about to build yourself.
