# Journal: Session Scoping, Auth, and Multi-Tenant Patterns

**Date:** 2026-04-03  
**Session context:** Building Throughline on Audos — post-Briefing app build, pre-production hardening

---

## What We Discovered

### 1. `useWorkspaceDB` session scoping doesn't work for our data

We built the initial app (Setup, Briefing) using direct REST calls to the `db-api` hook. Data inserted this way has `session_id = NULL`. We then learned that `useWorkspaceDB` with `shared: false` (the default) only shows rows where `session_id` matches the current browser session — meaning it shows nothing for REST-inserted data.

**Confirmed pattern going forward:** Always use `shared: true` + manual `org_id` column filter.

→ See **SDK-13-DATABASE-FAQ.md § 1** for full explanation and code examples.

### 2. `window.__spaceContext?.username` is the correct email accessor

We had been using `window.useSubscription?.()?.email` as the fallback for reading the authenticated user's email in app code. This is wrong. The correct platform API is:

```typescript
window.__spaceContext?.username  // returns the EmailGate email
```

This is populated after the EmailGate fires and persists for the session. Sessions are stable per email (same email = same session across devices).

→ See **SDK-13-DATABASE-FAQ.md § 2** for session identity details.

### 3. REST `db-api` is open — API key is possible

The `db-api` endpoint has no authentication. You can add `x-api-key` validation inside the hook's server function code. We haven't done this yet — needs Otto to edit the hook.

→ See **SDK-13-DATABASE-FAQ.md § 3** for the code pattern.

### 4. Column-based org isolation is correct for SaaS

Separate workspaces per org is not the right model for Throughline. `user_id` + `org_id` columns on every table, always filtered at the data layer, is the confirmed approach.

→ See **SDK-13-DATABASE-FAQ.md § 4–5** for recommended code patterns including a `useOrgDB` wrapper hook.

---

## What Changed in Code

- `identity.ts`: Fixed `window.__spaceContext?.username` (was `window.useSubscription?.()?.email`)

---

## What's Still Pending (as of this session)

- [ ] Ask Otto to add API key validation to the `db-api` hook
- [ ] Migrate app reads from direct REST fetch to `useWorkspaceDB` with `shared: true` + `org_id` filter
- [ ] Create `useOrgDB` wrapper hook as recommended in SDK-13
