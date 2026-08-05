# Otto's tool surface is a separate thing from the `platform.*` hooks your app can call

*Written 2026-07-23, after asking Otto directly (via the onboarding API's `/chat` endpoint) for the*
*concrete technical surface behind capabilities we'd only seen described secondhand — ad campaigns,*
*analytics, media generation. The answer clarified an architectural split this SDK hadn't stated*
*explicitly before.*

## The split

Audos exposes capability through two genuinely different layers, and it's easy to conflate them because
both get described in the same marketing language ("Audos can run ads," "Audos can generate video"):

1. **`platform.*` hooks** — callable from server-function code *you* write and deploy into a workspace
   (`db`, `platform.generateText`, `platform.sendEmail`, `platform.secretsProxy`, `useWorkspaceDB`,
   `useSession`, `useSpaceFiles`). These are the primitives documented in `docs/platform/06`. They run in
   *your* app's execution context.
2. **Otto-orchestration tools** — an internal MCP-style tool surface Otto itself calls when you ask it to
   do something in chat (`query_analytics`, `generate_image`, `generate_video`, `search_meta_targeting`,
   `get_ad_campaigns`, and dozens more). **As of 2026-07-23, Otto has no verified evidence these are also
   exposed as `platform.*` methods callable from app/hook code.** They are chat-triggered, not
   app-embeddable.

Practical consequence: if you're building a workspace app and want it to, say, generate a product image or
query funnel analytics *from inside the app itself*, don't assume the tool names above are hook calls you
can import. The confirmed app-callable surface is still just what `docs/platform/06` already lists. Ads,
richer analytics, and media generation are things you ask Otto to do on the founder's behalf, not
primitives your own code invokes directly — until Audos documents otherwise.

## `docs/platform/06`'s internal REST endpoints don't work from an external caller

`docs/platform/06-capabilities-reference.md` documents an "Internal Platform APIs" section
(`/api/crm/contacts/{workspaceId}`, `/api/funnel/metrics/{workspaceId}`, `/api/funnel/events/{workspaceId}`,
`/api/funnel/sessions/{workspaceId}`, `/api/spaces/{workspaceId}`) with the claim: *"No additional
authentication headers are required when called from server functions."*

Verified live 2026-07-23: calling these same URLs directly from an external client (plain `curl`, no
session, no server-function context) against the real DoKnow workspace returns a uniform
`401 Unauthorized: You must log in to access this resource.` on every one of the five endpoints tested.

This isn't a contradiction of the doc — the doc's claim is scoped to *"when called from server functions"*,
which implies a trusted execution context (a deployed hook running inside Audos's own infra) rather than
a fully public, workspace-ID-only auth model. But it's worth being explicit about, since it's easy to
misread that line as "these endpoints are unauthenticated" full stop. They are not reachable this way from
outside a real server function. Independent verification pending: whether they work called *from* an
actual deployed server function (as `docs/platform/06` claims) hasn't been separately confirmed by us —
only the external-caller failure mode has been.

## What actually is queryable, and how it was found

Rather than the REST paths in `docs/platform/06`, the verified way to pull structured workspace data as of
this writing is Otto's own `query_data_source` tool, with a `sourceId` enum (confirmed from live tool
schema): `contacts`, `funnel-events`, `ad-campaigns`, `printify-orders`, `stripe-payments`,
`stripe-subscriptions`, `stripe-customers`, `stripe-invoices`, `stripe-coupons`, `analytics-overview`,
`session-recordings`, `community-posts`. See `docs/platform/30-analytics-and-reporting-live-verification.md`
for what was actually pulled through it, including one real discrepancy it surfaced.

## Source

Otto chat, DoKnow workspace (`8a65a4ac-5a22-435f-b55f-c41ea34ca00d`), 2026-07-23, in response to a direct
request for the concrete tool surface behind the capability-matrix additions. Cross-checked the REST-auth
claim independently via direct `curl` against the live endpoints (see above) rather than taking Otto's
"these require login" framing on faith.
