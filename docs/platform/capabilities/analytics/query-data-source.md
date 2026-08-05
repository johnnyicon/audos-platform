# `query_data_source` — generalized structured query, with a real internal inconsistency

**Status: ⚠️ verified live, with a bug.** Otto-chat-triggered tool (see
`../../29-otto-tool-surface-vs-app-callable-hooks.md`).

**Params:** `sourceId` (required), `params` (optional object of string k/v). `sourceId` enum, confirmed
from the live tool schema: `contacts`, `funnel-events`, `ad-campaigns`, `printify-orders`,
`stripe-payments`, `stripe-subscriptions`, `stripe-customers`, `stripe-invoices`, `stripe-coupons`,
`analytics-overview`, `session-recordings`, `community-posts`. This is the closest thing to a general
"raw query" surface Otto exposes. Free.

## `sourceId=contacts` — ✅ verified

Returned 4 real records for the DoKnow workspace: `contact_dfc44bd0-…`, `contact_f4eb21ed-…`,
`contact_8b5023ff-…` (all QA/test addresses from prior verification passes), and one real contact
(`john@merkhetventures.com`). Shape: `id | email | (name) | source | (—)`. Count matches `query_analytics`
exactly (see `query-analytics.md`) — a genuine independent cross-check.

## `sourceId=funnel-events` — ❌ broken / degenerate, as tested

Returned exactly **1 record with every field blank** (id/type/space/visitor/data/date all empty) — for
the same workspace where `query_events` (see `query-events.md`) reported **187 events** over the identical
30-day window in the same conversation. Not root-caused: could be a different backing view/adapter for
this specific `sourceId`, a pagination default of 1, or something broken specifically in this
`query_data_source` path. **Don't use this for event data — use `query_events` instead.**

## Other `sourceId` values — 📄 not yet run

`ad-campaigns`, `printify-orders`, `stripe-payments`, `stripe-subscriptions`, `stripe-customers`,
`stripe-invoices`, `stripe-coupons`, `analytics-overview`, `session-recordings`, `community-posts` — enum
confirmed real from the live schema, but none independently exercised yet. Given the `funnel-events`
discrepancy found above, don't assume any of these are reliable until actually run.

Source: Otto chat, 2026-07-23. Narrative: `../../../../blog/experiments/0030-analytics-reporting-live-verification.md`.
