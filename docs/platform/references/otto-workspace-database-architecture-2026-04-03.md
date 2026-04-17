# Audos Workspace Database Architecture

Comprehensive guide for off-platform integration with Audos workspace databases.

---

## Table of Contents

1. [Overview & Architecture](#overview--architecture)
2. [Table Lifecycle](#table-lifecycle)
3. [Data Flow & Writers](#data-flow--writers)
4. [Naming & Column Conventions](#naming--column-conventions)
5. [Schema Management](#schema-management)
6. [Hooks & Database Interaction](#hooks--database-interaction)
7. [Outreach Leads Deep Dive](#outreach-leads-deep-dive)
8. [Off-Platform Integration Patterns](#off-platform-integration-patterns)
9. [Complete Table Reference](#complete-table-reference)
10. [Gotchas & Limitations](#gotchas--limitations)

---

## Overview & Architecture

### Database Technology

- **Engine**: PostgreSQL (managed, likely Google Cloud SQL or Supabase)
- **Isolation**: Schema-per-workspace multi-tenancy
- **Your Schema**: `ws_8f1ad824_832f_4af8_b77e_ab931a250625`

### Access Methods

| Method | Access Level | Use Case |
|-------|--------------|----------|
| **Direct SQL** (you have this) | Full read/write | Off-platform integration, reporting |
| **db-api hook** | Full CRUD + DDL | Apps, external services via HTTP |
| **Otto tools** (`db_query`, etc.) | Full CRUD + DDL | Agentic workflows, chat-based ops |
| **Mini-app SDK** (`useWorkspaceDB`) | Session-scoped CRUD | Frontend apps in Audos Space |
| **execute_sql** | SELECT only | Ad-hoc queries via Otto |

### Architecture Diagram

```
┌─────────────────────────────────┐  ┌─────────────────────────────────┐
│  Your Go Daemon              │  │  Audos Platform               │
│  ┌────────────────────────┐  │  │  ┌────────────────────────┐  │
│  │  Your Schema            │  │  │  │  Otto (Agent)          │  │
│  │  (my_app_*)            │  │  │  └───────┬───────────────┘  │
│  └────────────┬───────────┘  |  |          │
│               │               │  │          │ db_* tools    │
│               │               │  │          │              │
└─────────────┼──────────────────┘  └──────────┼───────────────────────┘
               │                            │
               │  Direct SQL                │
               │                            │
               ▼                            ▼
           ┌─────────────────────────────────┐
           │            PostgreSQL                 │
           │  ┌─────────────────────────────┐    │
           │  │ Schema: ws_8f1ad824_...       │    │
           │  │ (your workspace tables)    │    │
           │  └─────────────────────────────┘    │
           └─────────────────────────────────┘
```

---

## Table Lifecycle

### How Tables Are Created

Tables are **NOT** auto-provisioned. They are created explicitly through one of these methods:

| Creation Method | Who/What Uses It | Example |
|----------------|------------------|---------|
| **Otto `db_create_table`** | Agentic workflows via chat | "Create a table for voice profiles" |
| **db-api hook** | External services via HTTP | `POST /hooks/execute/.../db-api` with `action: 'create_table'` |
| **Direct SQL** (you) | Off-platform development | `CREATE TABLE ws_xxx.my_table (...)` |
| **delegate_database_design** | Otto subagent for complex schemas | "Design a database for my podcast studio" |

### What Controls the Schema?

There is **no centralized schema definition**. The platform uses a **discovery-based** approach — it tracks tables via a metadata registry (`__table_registry`) that stores:

- Table name
- Display name
- Description
- Column descriptions
- Creation timestamp

Tables created via direct SQL are **not automatically registered** in this metadata system.

### Table Lifecycle States

```
                    CREATE
                       │
                       ▼
                ┌───────────┐
                │  ACTIVE  │ ───┐ ALTER TABLE
                └─────┬─────┘    │
                      │     < ────┘
                      | TRUNCATE / DROP
                      ▼
                ┌───────────┐
                │  DELETED  │
                └───────────┘
```

---

## Data Flow & Writers

### Who Writes What?

| Table | Primary Writer | Secondary Writers | Notes |
|-------|---------------|-------------------|-------|
| `voice_profiles` | Otto (via app setup) | Mini-apps, db-api | Created during onboarding |
| `speakers` | Otto | Mini-apps, db-api | Parsed from transcripts |
| `reels` | Mini-apps | Otto, db-api | User-created content |
| `reel_captions` | Otto (AI generation) | Mini-apps | Generated content |
| `outreach_leads` | **Platform (Lead Scout)** | Otto (status updates) | Agentic CRM — see deep dive |
| `linked_references` | Otto (web fetch) | db-api | Cached web pages |
| `dashboard_activity` | Otto, Mini-apps, db-api | Any | Activity logging |

### Read-Only Tables?

**None of your workspace tables are read-only.** You have full write access to all of them via direct SQL.

However, be careful with:

- **`outreach_leads`** — The platform's Lead Scout feature may overwrite `status` or `notes` if you're both writing
- **`linked_references`** — Otto's web fetch caching may create duplicates

### Can You Write Directly via SQL?

**Yes, absolutely.** Since you have direct PostgreSQL access, you can:

```sql
-- Insert directly
INSERT INTO ws_8f1ad824_832f_4af8_b77e_ab931a250625.voice_profiles 
  (name, type, description, user_id, org_id)
  VALUES ('New Host', 'host', 'Description here', 'your-user-id', 'your-org-id');

-- Update directly
UPDATE ws_8f1ad824_832f_4af8_b77e_ab931a250625.speakers
  SET notes = 'Updated from Go daemon'
  WHERE id = 1;
```

**Caveat**: Writes via direct SQL bypass any platform validation or side effects (e.g., `updated_at` triggers).

---

## Naming & Column Conventions

### Table Naming

| Pattern | Mandatory? | Usage |
|---------|------------|-------|
| `app_*` prefix | **No** | Your tables don't use this prefix because they were created via Otto |
| `guest_prep_*` prefix | Convention | Groups related tables for GuestPrep app |
| `briefing_*` prefix | Convention | Groups related tables for Briefing app |
| `studio_*` prefix | Convention | Groups related tables for Studio app |
| `snake_case` | **Yes** | All tables use snake_case |

**Recommendation for your Go daemon**: Use a distinct prefix like `go_*` or `ext_*` to clearly identify tables managed by your off-platform system.

### Standard Columns

Every table created via Otto gets these automatically:

| Column | Type | Auto-Added | Purpose |
|--------|------|------------|---------|
| `id` | `serial` (PK) | ✓ Yes | Primary key |
| `created_at` | `timestamp` | ✓ Yes | Row creation time |
| `updated_at` | `timestamp` | ✓ Yes | Last modification time |
| `session_id` | `text` | ✓ Yes | User session scoping |

### `session_id` Explained

The `session_id` column is central to Audos's data scoping model:

- **What it references**: A visitor session in the Audos Space (the frontend app container)
- **Format**: UUID-like string (e.g., `abcd1234-5678-90ab-cdef-1234567890ab`)
- **Scoping**: When mini-apps use `useWorkspaceDB('table')`, they only see rows where `session_id` matches the current user's session
- **Shared data**: Use `useWorkspaceDB('table', { shared: true })` to read all rows regardless of session
- **Otto/API writes**: Data inserted via Otto or the db-api hook has `session_id = NULL`

**Important for your Go daemon**: If you insert data via direct SQL with `session_id = NULL`, mini-apps must use `{ shared: true }` to see it.

### `user_id` / `org_id` Explained

These are **your custom columns**, not platform-managed:

- **Not auto-populated**: The platform doesn't fill these automatically
- **Your IDs**: Use whatever format makes sense for your system
  - Email addresses (`john@merkhetventures.com`)
  - UUIDs (`8586c1a7-4d63-4b5f-9f4a-2a67c18d3b5e`)
  - External system IDs (`user_1234`)

Your current data uses:
- `user_id`: `john@merkhetventures.com` (email)
- `org_id`: `sow-good-to-grow-good` (slug)

### UUID vs Text IDs

| Column | Type | When Used |
|--------|------|-----------|
| `id` | `serial` (integer) | Primary keys — auto-incrementing |
| `related_id` | `uuid` | References to external entities |
| `outreach_batch_id` | `text` (UUID format) | Platform-generated batch IDs |
| `session_id` | `text` (UUID format) | Visitor session tracking |
| `user_id`, `org_id` | `text` | Custom IDs — any format |

### JSONB Columns

Several tables use `json` / `jsonb` columns for flexible data:

| Column | Table | Expected Structure |
|--------|-------|--------------------|
| `metadata` | `dashboard_activity` | Freeform — `{ "test": true, "source": "off-platform" }` |
| `long_form_samples` | `voice_profiles` | Array of samples — `[{ "title": "...", "content": "..." }]` |
| `learned_patterns` | (if exists) | AI training data — structure varies |

**No strict schema enforcement** — these columns accept any valid JSON.

---

## Schema Management

### Can You Add Your Own Tables?

**Yes, absolutely.** You have two options:

#### Option 1: Direct SQL (your Go daemon)

```sql
-- Create your own table
CREATE TABLE ws_8f1ad824_832f_4af8_b77e_ab931a250625.go_sync_state (
  id SERIAL PRIMARY KEY,
  entity_type TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  last_synced_at TIMESTAMP NOT NULL DEFAULT NOW(),
  sync_hash TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

#### Option 2: db-api hook (from your daemon via HTTP)

```bash
curl -X POST https://audos.com/api/hooks/execute/workspace-351699/db-api \
  -H "Content-Type: application/json" \
  -H "x-api-key: your-api-key" \
  -d '{
    "action": "create_table",
    "table": "go_sync_state",
    "columns": [
      { "name": "entity_type", "type": "text", "nullable": false },
      { "name": "entity_id", "type": "text", "nullable": false },
      { "name": "last_synced_at", "type": "timestamp" },
      { "name": "sync_hash", "type": "text" }
    ]
  }'
```

### Can You Alter Existing Tables?

**Yes.** Via direct SQL:

```sql
-- Add a column
ALTER TABLE ws_8f1ad824_832f_4af8_b77e_ab931a250625.voice_profiles
  ADD COLUMN external_id TEXT;

-- Add an index
CREATE INDEX idx_vp_external_id 
  ON ws_8f1ad824_832f_4af8_b77e_ab931a250625.voice_profiles(external_id);
```

### Will Audos Overwrite Your Changes?

**No.** The platform does not:
- Auto-migrate your schema
- Drop or alter existing tables
- Overwrite columns you've added

The schema is **append-only** from the platform's perspective. Otto only creates tables and adds columns when explicitly asked.

### Migration System?

There is **no built-in migration system** in Audos. We previously set up a local development + Atlas migration workflow where:
1. You run local Postgres for development
2. Use Atlas to generate migration files
3. Ask Otto to apply them to the Audos database

Your **Go daemon** could also manage migrations directly via `golang-migrate` or similar.

---

## Hooks & Database Interaction

### How Hooks Access the Database

Server functions (hooks) get a `db` object injected with these methods:

```javascript
// In a hook file
export default async function handler(request, { db, platform, respond }) {
  
  // Query table
  const rows = await db.query('voice_profiles', {
    filters: [{ column: 'type', operator: 'eq', value: 'host' }],
    limit: 10
  });
  
  // Insert rows
  await db.insert('dashboard_activity', [{
    activity_type: 'api_test',
    title: 'Test',
    description: 'From hook',
    metadata: { source: 'hook' }
  }]);
  
  // Update rows
  await db.update('speakers', {
    filters: [{ column: 'id', operator: 'eq', value: 1 }],
    data: { notes: 'Updated from hook' }
  });
  
  // Delete rows
  await db.delete('linked_references', {
    filters: [{ column: 'id', operator: 'eq', value: 99 }]
  });
  
  // List tables
  const tables = await db.listTables();
  
  respond(200, { rows, tables });
}
```

### Does ai-api Hook Touch the Database?

**No.** Your `ai-api` hook is purely for AI text generation — it proxies to OpenAI/Claude and doesn't read/write database tables.

If you wanted to log AI interactions to the database, you could modify the hook to add `db.insert()` calls.

---

## Outreach Leads Deep Dive

### What Is This Table?

`outreach_leads` is part of Audos's **agentic CRM** system — specifically the **Lead Scout** feature.

### How It Works

```
1. You ask Otto: "Find podcast producers for outreach"

2. Otto triggers: `start_lead_scout` tool
   ────────────
   │
   ▼
3. Lead Scout Agent:
   - Searches web/LinkedIn for matching profiles
   - AI scores each lead (relevance_score, ai_reason)
   - Inserts rows into `outreach_leads`
   │
   ▼
4. You review: Leads appear in Outreach window
   ───────────

5. Status transitions:
   new → drafted → contacted → responded → scheduled
                              ↓ not_interested
```

### Who Manages Status Transitions?

- **`new` → `drafted`**: Otto via `draft_outreach_email` tool (generates email draft)
- **`drafted` → `contacted`**: Platform when email is sent
- **`contacted` → `responded`/`scheduled`/`not_interested`**: User manually or via Otto

### Key Columns

| Column | Populated By | Purpose |
|--------|--------------|---------|
| `relevance_score` | Lead Scout AI | 0-100 relevance to your business |
| `ai_reason` | Lead Scout AI  | Why this lead is a good fit |
| `outreach_batch_id` | Lead Scout | Groups leads from same search job |
| `notes` | Otto (draft email) | Draft email content stored here |
| `status` | Platform/Otto | Current state in the funnel |
| `session_id` | NULL (agent-generated) | Not session-scoped |

### Can You Write to This Table?

**Yes.** You can insert your own leads:

```sql
INSERT INTO ws_8f1ad824_832f_4af8_b77e_ab931a250625.outreach_leads
  (name, email, title, company, relevance_score, ai_reason, status)
  VALUES ('Jane Doe', 'jane@example.com', 'Producer', 'Acme Pods', 85, 'Imported from Go daemon', 'new');
```

---

## Off-Platform Integration Patterns

### Recommended Architecture

```
┌────────────────────────────────────────┐
│  Your Go Daemon                           │
│  ┌───────────────────┐  ┌───────────┐ │
│  │  Your Postgres      │  │   Sync    │ │
│  │  (my_app schema)    │  │  Engine  │ │
│  └───────────────────┘  └────┬─────┘ │
└───────────────────────────────┼────────┘
                                       │
                    Direct SQL Read/Write
                                       │
┌───────────────────────────────┼────────┐
│               PostgreSQL                    │
│  ┌───────────────────┐  ┌───────────────────┐  │
│  │  my_app schema    │  │ ws_8f1ad824_... │  │
│  │  (your daemon's)  │  │  (workspace)    │  │
│  └───────────────────┘  └───────────────────┘  │
└─────────────────────────────────────────────┘
```

### Pattern 1: Direct SQL Access (Recommended)

Since you have direct PostgreSQL access, this is the most efficient pattern:

```go
// Go daemon example
package main

import (
    "database/sql"
    _ "github.com/lib/pq"
)

const audosSchema = "ws_8f1ad824_832f_4af8_b77e_ab931a250625"

func main() {
    db, _ := sql.Open("postgres", "postgres://user:pass@host/db?sslmode=require")
    
    // Read from Audos workspace
    rows, _ := db.Query(`SELECT id, name, type FROM ` + audosSchema + `.voice_profiles`)
    
    // Write to Audos workspace
    db.Exec(`INSERT INTO ` + audosSchema + `.dashboard_activity
        (activity_type, title, metadata) VALUES ($1, $2, $3::jsonb)`,
        "go_sync", "Sync completed", `{"source": "go_daemon"}`)
}
```

### Pattern 2: HTTP via db-api Hook

If you can't use direct SQL (e.g., from a serverless function):

```go
// HTTP request to db-api hook
payload := map[string]interface{}{
    "action": "query",
    "table":  "voice_profiles",
    "limit":  10,
}

req, _ := http.NewRequest("POST", 
    "https://audos.com/api/hooks/execute/workspace-351699/db-api", 
    bytes.NewBuffer(jsonPayload))
req.Header.Set("x-api-key", "your-api-key")
req.Header.Set("Content-Type", "application/json")
```

### Pattern 3: Sync Table (Best for Loose Coupling)

Create a dedicated sync table to track what's been synced:

```sql
-- In your schema (or Audos schema)
CREATE TABLE ws_8f1ad824_....go_sync_state (
  id SERIAL PRIMARY KEY,
  entity_type TEXT NOT NULL,        -- 'voice_profile', 'speaker', etc.
  entity_id INTEGER NOT NULL,       -- ID in the source table
  last_synced_at TIMESTAMP NOT NULL,
  sync_hash TEXT,                   -- MD5 of row data to detect changes
  sync_direction TEXT DEFAULT 'both' -- 'to_audos', 'from_audos', 'both'
);
```

### Pattern 4: Webhook Callbacks (If You Need Real-Time)

Create a hook that notifies your daemon when data changes:

```javascript
// notify-daemon hook
export default async function handler(request, { db, respond }) {
  const { eventType, table, rowId } = request.body;
  
  // Call your Go daemon
  await fetch('https://your-daemon.example.com/webhook', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ eventType, table, rowId })
  });
  
  respond(200, { ok: true });
}
```

Then call this hook from other hooks or mini-apps when data changes.

### Polling Strategy

If you prefer polling, use `updated_at` timestamps:

```sql
-- Get all rows modified since last sync
SELECT * FROM ws_8f1ad824_....voice_profiles
  WHERE updated_at > '2026-04-03T14:00:00Z';
```

---

## Complete Table Reference

### Your 17 Workspace Tables

| Table | Rows | Purpose | Key Columns |
|-------|------|---------|-------------|
| `voice_profiles` | 2 | Voice fingerprints for hosts/brands | `name`, `type`, `is_trained`, `long_form_samples` |
| `speakers` | 3 | Speaker registry for transcripts | `name`, `role`, `voice_profile_id` |
| `reels` | 1 | Social media clips | `transcript`, `status`, `scheduled_date` |
| `reel_captions` | 0 | AI-generated captions | `reel_id`, `platform`, `caption` |
| `outreach_leads` | 11 | Lead Scout CRM | `relevance_score`, `ai_reason`, `status` |
| `linked_references` | 2 | Cached web pages | `url`, `content`, `fetched_at` |
| `dashboard_activity` | 3 | Activity log | `activity_type`, `title`, `metadata` |
| `guest_prep_podcast_profiles` | 1 | Podcast identity config | `name`, `tone`, `brand_voice` |
| `guest_prep_research_sessions` | 0 | Guest research data | `guest_name`, `transcript`, `research_package` |
| `guest_prep_ros_versions` | 0 | Run of Show version history | `version`, `content` |
| `briefing_podcast_profiles` | 0 | Briefing app profiles | `name`, `description` |
| `briefing_research_sessions` | 0 | Briefing sessions | `guest_name`, `briefing_data` |
| `voice_refinements` | 0 | Voice model training data | `voice_profile_id`, `feedback` |
| `studio_episodes` | 0 | Episode drops | `title`, `published_at` |
| `studio_time_tracking` | 0 | Automation metrics | `time_saved_minutes` |
| `studio_generated_content` | 0 | Generated platform content | `platform`, `content`, `status` |
| `podcast_setup_profiles` | 0 | Setup wizard data | `name`, `branding` |

---

## Gotchas & Limitations

### 1. `session_id` Scoping

**Problem**: Data inserted with `session_id = NULL` (like your Go daemon would) is invisible to mini-apps using `useWorkspaceDB('table')`.

**Solution**: Mini-apps must use `{ shared: true }`:

```typescript
// Mini-app
const { data } = useWorkspaceDB('voice_profiles', { shared: true });
```

### 2. `updated_at` Not Auto-Updated

**Problem**: Direct SQL UPDATES may not update the `updated_at` column automatically (depends on whether a trigger exists).

**Solution**: Always set it explicitly:

```sql
UPDATE ws_xxx.voice_profiles
  SET name = 'New Name', updated_at = NOW()
  WHERE id = 1;
```

### 3. No Foreign Key Enforcement

**Problem**: Foreign keys are documented but not always enforced at the database level.

**Solution**: Validate in your application code, or add constraints via direct SQL:

```sql
ALTER TABLE ws_xxx.speakers
  ADD CONSTRAINT fk_speakers_voice_profile
  FOREIGN KEY (voice_profile_id) REFERENCES ws_xxx.voice_profiles(id);
```

### 4. No Change Data Capture

**Problem**: There's no built-in CDC or event stream for database changes.

**Solution**: 
- Poll using `updated_at` timestamps
- Build webhook notifications into hooks/apps
- Use PostgreSQL LISTEN/NOTIFY if you have connection access

### 5. Concurrent Writes

**Problem**: If both your Go daemon and Otto/mini-apps write to the same rows, you could get last-write-wins conflicts.

**Solution**:
- Use distinct tables or rows for each system
- Add a `last_modified_by` column to track ownership
- Use optimistic locking with version numbers

---

## Summary: Key Takeaways for Your Go Daemon

1. **Direct SQL is fully supported** — read and write freely
2. **Use a distinct prefix** for your tables (`go_*` or `ext_*`)
3. **Set `session_id = NULL`** for shared data (or a known sentinel value)
4. **Always set `updated_at`** explicitly in UPDATEs
5. **Use sync tables** to track what's been synced
6. **Avoid concurrent writes** to the same rows as Audos

---

*Document generated: 2026-04-03*
*Workspace: Throughline (ws_8f1ad824_832f_4af8_b77e_ab931a250625)*