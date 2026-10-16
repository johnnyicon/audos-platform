Title:
The `ai-api` hook refuses 10–30% of small chat requests at random with an instant HTTP 500 "AI generation failed", and refuses everything for ~1 minute after a burst — with no 429, no Retry-After and nothing in the body saying it was rate limiting

Workspace:
- Slug: `workspace-351699`
- Owner email: `john@merkhetventures.com`
- Hook: `POST https://audos.com/api/hooks/execute/workspace-351699/ai-api`, action `chat`, model `gpt-5.5`, `temperature: 1`, `maxTokens: 1500`

Intent:
We call the AI hook from our own server (Go daemon) for an agent-style script checker: many small chat requests (1–5 s of generation each, one action per turn), about 55 requests per check, four conversations in parallel. Total volume on 15 Oct 2026: about 1,170 requests over the afternoon.

Expected behavior:
Either the requests succeed, or the hook tells us it is rate limiting: HTTP 429 with `Retry-After` (or the limit in the body), so a client can pace itself. A refusal should be distinguishable from a model failure.

Actual behavior:
1. After a burst (24 requests in 11 s) the hook refused every request for about a minute: HTTP 500, body `{"error":"AI generation failed"}`, returned instantly (not a timeout).
2. At a steady 20 requests a minute it still refused 10% of requests at random (24 of 241 in one run), and 30% in another run at the same pace within the hour (87 of 288). At 30 a minute, about one refusal a minute. The refusals are spread evenly, not tied to request size or content: the same request succeeds on retry 3–12 s later.
3. The response carries no `Retry-After`, no rate-limit headers, no 429, and the body text is the same one a genuine upstream model error would produce, so the client cannot tell the two apart.
4. Long requests are unaffected: our one-shot design makes 11–29 long requests (10–40 s each) per check and saw none of this on the same day. The behaviour is a function of request count, not load.

Reproduction steps:
1. POST to the hook 25 small chat requests within ~10 s (`{"action":"chat","model":"gpt-5.5","maxTokens":1500,"temperature":1,"systemPrompt":"Reply with the JSON {\"ok\":true} only.","messages":[{"role":"user","content":"ping"}]}`).
2. Keep sending one every 2–3 s: for about a minute every response is HTTP 500 `{"error":"AI generation failed"}` in under a second.
3. Then send one request every 3 s for ten minutes (20/min): between 10% and 30% of them return the same instant 500; immediate retries of the same body succeed.

Evidence:
- Run logs with timestamps, request ids and the error bodies, 15 Oct 2026 12:06–12:50 local: `run-1-rate-limited.log` (burst), `run-3-paced-20-per-min/run.log` (10% refused), `replay.stdout.log` (30% refused), in the Throughline repo under `docs/plans/script-loop/evidence/replay/2026-10-15-agentic/`. Example lines: `12:06:30 HTTP 500 kind=facts job=replay-ep3-facts#take/2 error="AI generation failed"`, `12:06:31 … #take/2-retry error="AI generation failed"` (retry 1 s later also refused, inside the burst lock-out).
- Same hook, same hour, long requests: 0 refusals in 40 requests.

Impact and priority:
Medium–high for anyone building an agent loop on the hook. We had to add client-side pacing (20 requests a minute) and up to three retries with growing pauses (3, 6, 12 s) before a single run completed with every check intact; four of five live runs were lost to this, none to our own code. Because the refusal looks like a model failure, the first day was spent debugging the model.

Security boundary:
None. Our own workspace, our own content.

Success criteria:
- The hook's rate limit (requests per minute, burst size, per workspace or per hook) is documented in the AI hook reference.
- A refusal for rate is reported as one: HTTP 429 with `Retry-After`, or at least `{"error":"rate_limited","retry_after_ms":…}` in the body, never the generic "AI generation failed".
- Ideally the random 10–30% refusals at 20/min (well below the burst limit) are fixed, since they look like a leaky bucket mis-sized for a single workspace.

Related: ticket 351699-84HKUSAA (the ~40 s execution cutoff on the same hook), filed earlier the same day.
