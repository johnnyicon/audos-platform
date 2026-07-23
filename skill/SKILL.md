---
name: audos-platform
description: "Audos platform development — working with Audos workspaces, server functions (hooks), database tables, or any API at audos.com/api/hooks/execute/. TRIGGER when: user is working in an audos-workspace/ folder, calling Audos APIs, building or modifying server functions, asking about db/platform/fetch runtime, or developing apps on the Audos platform. DO NOT TRIGGER for: general React/TypeScript work unrelated to Audos platform integration."
---

# Audos Platform Skill

You are working inside an **Audos workspace** — a managed no-code/low-code platform that provides a hosted backend (PostgreSQL, React runtime, AI, email, storage, CRM, analytics) accessible via custom HTTP server functions called "hooks".

**You do not have direct access to the platform internals.** All changes go through:
1. **Otto** (the Audos AI assistant) — for apps, tables, landing pages, server functions
2. **HTTP server function APIs** — for reading/writing data from local code

---

## Progressive Disclosure — Read Only What You Need

Do not load all docs upfront. Read the specific file for the task at hand.

| Task | Read This |
|------|-----------|
| Understand the platform architecture | `/Users/kanekoa/Workspace/audos-platform/docs/platform/02-platform-overview.md` |
| Understand workspace folder structure | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/README.md` |
| Understand the development workflow | `/Users/kanekoa/Workspace/audos-platform/docs/platform/04-development-workflow.md` |
| **Deploy a React app via GitHub Dev Mode** (CDN deps, routing, platform detection) | `/Users/kanekoa/Workspace/audos-platform/docs/platform/16-github-dev-mode-app-deployment.md` |
| Create or modify a server function | `/Users/kanekoa/Workspace/audos-platform/docs/platform/07-api-development-guide.md` |
| Use runtime globals (`db`, `platform`, `fetch`, `respond`) | `/Users/kanekoa/Workspace/audos-platform/docs/platform/08-server-function-runtime.md` |
| Copy a server function template | `/Users/kanekoa/Workspace/audos-platform/docs/platform/09-server-function-templates.md` |
| Build a composite/aggregate API | `/Users/kanekoa/Workspace/audos-platform/docs/platform/05-composite-apis.md` |
| Call any Throughline API endpoint | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/throughline-api-reference.md` |
| Work with the database or tables | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/database-api.md` |
| Generate AI content (quick reference) | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/ai-generation-api.md` |
| **AI hook capability matrix** (models, latency, vision/tools/JSON mode) | `/Users/kanekoa/Workspace/audos-platform/docs/AI-HOOK-CAPABILITY-MATRIX.md` |
| Send email | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/email-api.md` |
| Manage contacts / CRM | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/crm-api.md` |
| Read analytics | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/analytics-api.md` |
| Upload or manage files | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/storage-api.md` |
| Schedule tasks or cron jobs | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/scheduler-api.md` |
| Fetch or analyze web pages | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/web-api.md` |
| Check known platform bugs or open issues | `/Users/kanekoa/Workspace/audos-platform/BACKLOG.md` |
| Check recent platform changes | `/Users/kanekoa/Workspace/audos-platform/docs/throughline/changelog/` |
| **Add or edit a Audos SDK blog entry** (findings write-up, cohort-facing) | `/Users/kanekoa/Workspace/audos-platform/blog/HOW-TO-UPDATE.md` |
| **Eliminate the chat shell / get a real full-screen single-app experience** | `/Users/kanekoa/Workspace/audos-platform/docs/platform/22-eliminating-the-chat-shell-playbook.md` |
| **Dispatching any job to Otto/Cursor** (standard operating pattern — read this before writing your first job) | `/Users/kanekoa/Workspace/audos-platform/docs/platform/23-dispatch-then-poll-sop.md` |
| **Deciding where to write up something you just learned** (log vs backlog vs doc vs blog) | `/Users/kanekoa/Workspace/audos-platform/docs/platform/24-where-new-findings-go.md` |
| **A job is stuck/hung, or you need to escalate a real platform bug to a human** | `/Users/kanekoa/Workspace/audos-platform/docs/platform/25-escalation-and-support-paths.md` |
| **Building anything private, invite-only, or with custom auth** (why EmailGate can't just be disabled, and what happened when we tried to replace it) | `/Users/kanekoa/Workspace/audos-platform/docs/platform/26-unified-space-signed-out-view-is-mandatory.md` |
| **Redeploy the SDK blog, or set up Basic Auth on a Cloudflare Worker + Static Assets** (includes a real caching gotcha that silently defeated the auth check on first deploy) | `/Users/kanekoa/Workspace/audos-platform/docs/platform/27-blog-hosting-cloudflare-worker.md`, `blog/cf-worker/redeploy.sh` |
| **See what's changed in the SDK's own guidance recently** (terse delta, not the full story) | `/Users/kanekoa/Workspace/audos-platform/CHANGELOG.md` |
| **Do any of: log an activity, write a blog entry, archive the log, update the changelog, or run a full harvest pass** — these are real Claude Code skills, not just docs; invoke by name or let the trigger phrases fire. **They live in `doknow-kb/.claude/skills/`, not this repo** — the activity log itself is DoKnow-local by design (see `24-where-new-findings-go.md`), so the skills are scoped there too, and operate cross-repo into this one for the blog/changelog/docs steps | `~/doknow-kb/.claude/skills/audos-sdk-{log-activity,write-blog-entry,archive-activity-log,update-changelog,harvest}/SKILL.md` |

---

## Non-Negotiable Rules

### What You Can and Cannot Change
- **Apps, components, landing pages** → You cannot edit these directly. Ask Otto via the Audos platform.
- **Server functions (hooks)** → Created and modified via Otto using `manage_server_functions`. You can write the code locally and ask Otto to deploy it.
- **Database tables** → Created via Otto using `db_create_table`. Queried/written via the `db-api` endpoint.
- **Data files (`audos-workspace/data/`)** → Can be read locally; writes should go through the API.

### Server Function Runtime Limits
The server function runtime is **not Node.js**. These are unavailable:
- `URLSearchParams` — use the `buildQuery()` helper instead
- `response.headers.get()` — read response as text instead
- `require()` / `import` — all code must be inline
- `Buffer`, `process`, `setTimeout`, `setInterval`, `fs`, `path`, `crypto`

### Always Test After Changes
Run the test script after any API or server function change:
```bash
bun run /Users/kanekoa/Workspace/audos-platform/docs/throughline/test-apis.ts
```

---

## Throughline Workspace Constants

| Property | Value |
|----------|-------|
| Workspace number | `351699` |
| Workspace ID | `8f1ad824-832f-4af8-b77e-ab931a250625` |
| Base API URL | `https://audos.com/api/hooks/execute/workspace-351699` |
| Landing page | `https://www.trythroughline.com` |
| App | `https://app.trythroughline.com` |
