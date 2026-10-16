Title:
The `ai-api` hook cuts off model calls at ~40 s with "No response from AI model" — too short for reasoning models, and the body doesn't say it was a timeout

Workspace:
- Slug: `workspace-351699`
- Owner email: `john@merkhetventures.com`
- Hook: `POST https://audos.com/api/hooks/execute/workspace-351699/ai-api`, action `chat`

Intent:
We call the AI hook from our own server (Go daemon) to check podcast scripts with `gpt-5.5`: one request per script section, asking for a JSON list of suggested edits. The answers are long, and `gpt-5.5` reasons before it writes.

Expected behavior:
A chat request that takes 60–90 s to generate completes, or the hook tells us the limit it enforces so we can design around it. A timeout is reported as a timeout.

Actual behavior:
Requests that run past about 40 s return HTTP 500 with `{"error":"No response from AI model","_meta":{"success":true,"durationMs":41481,"logs":[]}}` (also seen at 44551 ms and ~55 s). `_meta.success` is `true`, `logs` is empty, and nothing in the body says the request was cut off for time. The upstream call may complete after the hook stops waiting, so we are billed with no result. Measured on 15 Oct 2026: with one script section per request, 2 of 9 "house rules" requests were cut off (retries succeeded); a whole-script request was cut off every time.

Reproduction steps:
1. POST to the hook with `{"action":"chat","model":"gpt-5.5","maxTokens":4000,"temperature":1,"systemPrompt":"<an instruction asking for up to 25 JSON suggestions>","messages":[{"role":"user","content":"<about 600 words of text plus 26 rules>"}]}`.
2. Repeat 5–10 times; the longer generations (over ~40 s) return the 500 above.
3. Compare: the same request asking for "at most 5 suggestions" returns in 10–20 s.

Evidence:
- Hook responses with `durationMs` 41481 and 44551 on the error, 15 Oct 2026 (our run log).
- A direct ping (`"Reply with the JSON {"ok":true} only."`) returns in ~2 s with `reasoning_tokens: 9`.
- Three concurrent requests each complete in 2 s, so concurrency is not the limit; per-request wall time is.

Impact and priority:
Medium. We work around it by splitting work into per-section requests, running them in parallel and retrying a cut-off once, which roughly halves wall time but keeps a floor of ~40 s per request. Any Audos app using a reasoning model for a long structured answer will hit this, and the error text sends people looking at the model rather than the hook.

Security boundary:
None. Our own workspace, our own content.

Success criteria:
- The AI hook's execution limit is documented (number of seconds) in the hook's reference.
- A timeout is reported as one (`"error":"timeout"`, with the limit and `durationMs`), not as "No response from AI model".
- Ideally the AI hook's limit is longer than the general server-function limit (90–120 s), since model calls are the one case where waiting is the whole job.

Also noticed, same day: the hook's default `temperature` (0.7) is forwarded to `gpt-5.5`, which refuses anything but 1 (`400 unsupported_value`); the caller has to send `temperature: 1` explicitly. Worth a model-aware default.
