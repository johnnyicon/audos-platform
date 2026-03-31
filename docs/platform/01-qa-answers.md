# Audos Platform Q&A - Your Questions Answered

*Last Updated: March 31, 2026*

**This document answers specific questions raised during development.**

---

## Q1: Where Does Content Get Stored?

There are **two separate storage systems** in Audos:

### Database Tables (Structured Data)
- **What**: Rows and columns in PostgreSQL
- **Access**: `db-api` with `query`, `insert`, `update`, `delete`
- **Examples**: Transcripts, voice profiles, episode metadata, guest research notes

### File Storage (Binary Files)
- **What**: Google Cloud Storage (GCS) bucket
- **Access**: `storage-api` with `upload`, `list`, `get`
- **Examples**: Audio files, images, PDFs, exported documents

### In Your Throughline Workspace

| Data Type | Where It Lives |
|-----------|----------------|
| Episode transcripts | `studio_episodes.transcript` (text column) |
| Voice profiles | `voice_profiles.description` (JSON column) |
| Guest research | `guest_prep_research_sessions.research_notes` |
| Uploaded audio | GCS via `storage-api` → returns URL |

---

## Q2: What Are the Audos Primitives (Foundational Components)?

### Server-Side Primitives (Available in Server Functions)

| Primitive | Description | How to Use |
|-----------|-------------|------------|
| **`db`** | Database operations | `db.query()`, `db.insert()`, `db.update()`, `db.delete()` |
| **`platform.generateText()`** | AI text generation (GPT-4o-mini) | Pass a prompt, get text back |
| **`platform.sendEmail()`** | Transactional email | to, subject, body (html or text) |
| **`fetch()`** | HTTP requests | Standard fetch API (with limitations) |
| **`respond()`** | Return HTTP response | `respond(200, { data })` |
| **`console`** | Logging | `console.log()` - visible in hook logs |

### Internal Platform APIs (Accessed via fetch)

| API Route | Purpose |
|-----------|---------|
| `/api/crm/contacts/{workspaceId}` | Contact management (CRM) |
| `/api/funnel/metrics/{workspaceId}` | Visitor analytics |
| `/api/funnel/sessions/{workspaceId}` | Session data |
| `/api/app-skills/...` | App backend capabilities |
| `/api/spaces/...` | Space management |

### Client-Side Primitives (Available in Apps)

| Primitive | Description |
|-----------|-------------|
| **`useSpaceFiles()`** | Persist JSON data to workspace storage |
| **`useWorkspaceDB()`** | Access workspace database tables |
| **`useSession()`** | Get current user session data |
| **Stripe integration** | Payments, subscriptions, checkout |

---

## Q3: What Are Composite/Aggregate APIs?

Yes! The term you're looking for is **composite APIs** (or "aggregate APIs", "orchestration APIs").

### The Pattern

```
Composite API = Primitive 1 + Primitive 2 + ... + Business Logic
```

### Examples

| Composite API | Primitives Used | What It Does |
|---------------|-----------------|---------------|
| `guest-research-api` | `fetch()` + `platform.generateText()` + `db.insert()` | Fetch guest info, AI summarize, save to DB |
| `transcript-api` | `fetch()` → external service + `db.insert()` | Send audio to Whisper/Deepgram, save result |
| `social-clips-api` | `db.query()` + `platform.generateText()` | Load transcript, AI identify clip moments |
| `email-digest-api` | `db.query()` + `platform.generateText()` + `platform.sendEmail()` | Query data, AI format, send email |

---

## Q4: Can I Build Composite Logic Locally instead of Server Functions?

**Yes and No.** It depends on what you're trying to do.

### YES - You Can Do Locally:

✅ Orchestrate multiple API calls from your local code

```python
# Local Python script that orchestrates multiple APIs
async def research_guest(guest_name, url):
    # Step 1: Fetch web content
    web_response = await api.web_fetch(url)
    
    # Step 2: Generate AI summary
    summary = await api.ai_generate(f"Summarize this: {web_response.content}")
    
    # Step 3: Save to database
    await api.db_insert("guest_prep_research_sessions", {
        "guest_name": guest_name,
        "research_notes": summary
    })
```

⌓ This works because you're calling the existing primitive APIs (`web-api`, `ai-api`, `db-api`).

### NO - You Need a Server Function When:

1. **Apps need to call it** - React apps in Audos can only call server functions, not your local machine
2. **Webhooks need to trigger it** - External services need a public URL to call
3. **Scheduled tasks** - Cron jobs run on the platform, not your local machine
4. **Single atomic operation** - Reduce latency by doing multiple steps in one call

### Decision Flowchart

```
Do I need this functionality from...

                 +───────────+
                 | Local dev? |
                 +─────┬─────+
                       |
        Yes ┌─────┴────┐ No
            |               |
    ✅ Use existing     Skip server function
    primitive APIs     (no HTTP wrapper needed)
            |
            v
    Need it from Audos apps?
            |
       Yes ┌┴┐ No
            |     |
      Create   Orchestrate
      Server   locally
      Function
```

---

## Q5: When Do I Need to Go Back to Audos/Otto to Develop Capabilities?

### Develop Locally When:
- Building scripts that USE the APIs
- Testing data flows
- Prototyping new features
- Building CLI tools
- Analytics and reporting scripts

### Go Back to Audos/Otto When:
- Need a new **HTTP endpoint** (server function)
- Need a new **database table**
- Need to **update the React apps**
- Need to **change the landing page**
- Need to access **internal platform APIs** that aren't wrapped yet
- Setting up **scheduled tasks**
- Configuring **webhooks**

### Workflow Example

```
1. You want to build a guest research feature

2. Locally, you prototype it:
   - Use web-api to fetch guest websites
   - Use ai-api to generate summary
   - Use db-api to save results
   
3. It works! But now you want it in the app...

4. Go back to Otto and ask:
   "Create a guest-research-api server function that:
    - Takes a guest name and URL
    - Fetches the URL content
    - Generates an AI summary
    - Saves to guest_prep_research_sessions"
    
5. Now the React app can call it directly
```

---

## Q6: What Are the Workspace Folders and What Are They For?

Based on the actual Throughline workspace structure:

| Folder | Purpose | Editable Locally? |
|--------|---------|-------------------|
| **apps/** | React components for each mini-app | Via delegate_app_edit only |
| **assets/** | Static files (images, fonts) | Upload via storage-api |
| **community/** | Community features config | Via platform UI |
| **components/** | Shared React components | Via delegate_app_edit only |
| **data/** | JSON data files (app state) | Via useSpaceFiles() hook |
| **hooks/** | Server functions (APIs) | **Yes** - create via Otto |
| **landing-pages/** | Landing page React code | Via delegate_landing_page_edit |
| **lib/** | Shared utilities, types, helpers | Via Otto or direct edit |
| **tools/** | Internal dashboards & admin tools | Via create_dashboard |

### Key Insight

The **hooks/** folder (server functions) is the primary way you extend the platform's capabilities. This is where all 8 of your APIs live.

---

## Q7: How Do I Develop With This Framework Locally?

### What You Can Do Locally

1. **Call all HTTP APIs** - db, ai, email, web, crm, analytics, storage, scheduler
2. **Build scripts** that orchestrate multiple API calls
3. **Test data flows** before implementing in apps
4. **Build CLI tools** for your podcast workflow

### What Requires Otto/Platform

1. **New server functions** → `manage_server_functions`
2. **New database tables** → `db_create_table`
3. **App changes** → `delegate_app_edit`
4. **Landing page changes** → `delegate_landing_page_edit`
5. **Branding changes** → workspace-branding.json

### The Hybrid Workflow

```
┌─────────────────────┐     ┌─────────────────────┐
│  Local Development     │     │   Audos Platform      │
│                      │     │                      │
│  Python/scripts that   |     |  Apps (React)         │
│  call HTTP APIs       │     |  Server Functions      │
│                      │     │  Database Tables      │
│  (Prototyping, CLI)   │     |  (Production)          │
└──────────┬──────────┘     └──────────┬──────────┘
           │                          │
           └───────── HTTP ────────┘
                      APIs
```

---

## Summary

| Question | Short Answer |
|----------|--------------|
| Where is content stored? | Database tables (structured) + GCS (files) |
| What are primitives? | db, platform.generateText, fetch, email + internal APIs |
| What are composite APIs? | Server functions that combine multiple primitives |
| Can I build locally? | Yes for scripts, no for app/API changes |
| When to go back to Otto? | New endpoints, tables, apps, landing pages |