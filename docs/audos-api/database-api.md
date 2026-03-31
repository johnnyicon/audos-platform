# Database API

> **Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/db-api`
>
> **Method:** `POST`
>
> **Content-Type:** `application/json`

Full CRUD access to all workspace database tables. PostgreSQL under the hood with workspace-isolated schemas.

---

## Quick Start

```bash
# List all tables
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/db-api" \
  -H "Content-Type: application/json" \
  -d '{"action": "list-tables"}'

# Query a table
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/db-api" \
  -H "Content-Type: application/json" \
  -d '{"action": "query", "table": "voice_profiles"}'
```

---

## Actions

### list-tables

List all available tables in the workspace.

**Request:**
```json
{
  "action": "list-tables"
}
```

**Response:**
```json
{
  "success": true,
  "tables": [
    {
      "name": "voice_profiles",
      "displayName": "Voice Profiles",
      "description": "Voice fingerprints for hosts and brand",
      "rowCount": 3
    }
  ],
  "count": 15
}
```

---

### describe

Get the schema of a specific table.

**Request:**
```json
{
  "action": "describe",
  "table": "voice_profiles"
}
```

**Response:**
```json
{
  "success": true,
  "table": "voice_profiles",
  "columns": [
    { "name": "id", "type": "integer", "nullable": false, "description": "Auto-generated ID" },
    { "name": "name", "type": "text", "nullable": true, "description": "Profile name" },
    { "name": "type", "type": "text", "nullable": true, "description": "host, guest, or brand" }
  ]
}
```

---

### query

Query data from a table with filtering, sorting, and pagination.

**Request:**
```json
{
  "action": "query",
  "table": "voice_profiles",
  "columns": ["id", "name", "type"],
  "filters": [
    { "column": "type", "operator": "eq", "value": "host" }
  ],
  "orderBy": { "column": "created_at", "direction": "desc" },
  "limit": 10,
  "offset": 0
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"query"` |
| `table` | string | Yes | Table name |
| `columns` | array | No | Columns to return (default: all) |
| `filters` | array | No | Filter conditions |
| `orderBy` | object | No | `{ column, direction: "asc" \| "desc" }` |
| `limit` | integer | No | Max rows to return (default: 50) |
| `offset` | integer | No | Rows to skip for pagination |

**Filter Operators:**

| Operator | Description | Example Value |
|----------|-------------|---------------|
| `eq` | Equals | `"published"` |
| `neq` | Not equals | `"draft"` |
| `gt` | Greater than | `100` |
| `gte` | Greater than or equal | `100` |
| `lt` | Less than | `50` |
| `lte` | Less than or equal | `50` |
| `like` | Pattern match (case-sensitive) | `"%podcast%"` |
| `ilike` | Pattern match (case-insensitive) | `"%Podcast%"` |
| `in` | In list | `["draft", "published"]` |
| `is_null` | Is null | (no value needed) |
| `not_null` | Is not null | (no value needed) |

**Response:**
```json
{
  "success": true,
  "table": "voice_profiles",
  "rows": [
    { "id": 1, "name": "John's Voice", "type": "host" }
  ],
  "count": 1
}
```

---

### insert

Insert one or more rows into a table.

**Request:**
```json
{
  "action": "insert",
  "table": "voice_profiles",
  "data": {
    "name": "Guest Voice Profile",
    "type": "guest",
    "description": "Voice characteristics for guest speakers"
  }
}
```

Or insert multiple rows:
```json
{
  "action": "insert",
  "table": "voice_profiles",
  "data": [
    { "name": "Profile 1", "type": "host" },
    { "name": "Profile 2", "type": "guest" }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "table": "voice_profiles",
  "inserted": 1,
  "rows": [
    { "id": 4, "name": "Guest Voice Profile", "type": "guest" }
  ]
}
```

---

### update

Update existing rows in a table.

**Request:**
```json
{
  "action": "update",
  "table": "voice_profiles",
  "filters": [
    { "column": "id", "operator": "eq", "value": 1 }
  ],
  "data": {
    "description": "Updated description"
  }
}
```

**Response:**
```json
{
  "success": true,
  "table": "voice_profiles",
  "updated": 1
}
```

---

### delete

Delete rows from a table.

**Request:**
```json
{
  "action": "delete",
  "table": "voice_profiles",
  "filters": [
    { "column": "id", "operator": "eq", "value": 99 }
  ]
}
```

**Response:**
```json
{
  "success": true,
  "table": "voice_profiles",
  "deleted": 1
}
```

---

## Table Schemas

### voice_profiles

Voice fingerprints for hosts and brand identity.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `name` | text | Profile name |
| `type` | text | `host`, `guest`, or `brand` |
| `description` | text | Voice characteristics description |
| `tone_keywords` | json | Array of tone descriptors |
| `vocabulary_preferences` | json | Preferred words/phrases |
| `formatting_rules` | json | Output formatting preferences |
| `example_outputs` | json | Sample content in this voice |

---

### speakers

Speaker registry for transcript parsing.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `name` | text | Speaker name |
| `role` | text | `host` or `guest` |
| `voice_profile_id` | integer | FK to voice_profiles |
| `aliases` | json | Alternative names/labels |

---

### voice_refinements

Training data for voice model refinement.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `voice_profile_id` | integer | FK to voice_profiles |
| `original_text` | text | AI-generated text |
| `refined_text` | text | User-edited version |
| `feedback_type` | text | Type of refinement |
| `context` | text | Where this was used |

---

### studio_episodes

Episode drops for content generation.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `title` | text | Episode title |
| `description` | text | Episode description |
| `transcript` | text | Full transcript |
| `audio_url` | text | Link to audio file |
| `publish_date` | date | Publication date |
| `status` | text | `draft`, `processing`, `ready` |
| `metadata` | json | Additional episode data |

---

### studio_generated_content

Generated content for each platform.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `episode_id` | integer | FK to studio_episodes |
| `platform` | text | `linkedin`, `twitter`, `instagram`, etc. |
| `content_type` | text | `post`, `thread`, `caption`, `story` |
| `content` | text | The generated content |
| `status` | text | `draft`, `approved`, `posted` |
| `voice_profile_id` | integer | FK to voice_profiles |

---

### studio_time_tracking

Time saved through automation.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `task_type` | text | Type of task automated |
| `estimated_manual_minutes` | integer | Time if done manually |
| `actual_minutes` | integer | Time with automation |
| `episode_id` | integer | FK to studio_episodes |

---

### reels

Content pieces for social posting.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `title` | text | Reel title |
| `transcript_segment` | text | Source transcript portion |
| `start_time` | text | Timestamp in episode |
| `end_time` | text | End timestamp |
| `episode_id` | integer | FK to studio_episodes |
| `status` | text | `draft`, `ready`, `posted` |

---

### reel_captions

Generated captions per platform.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `reel_id` | integer | FK to reels |
| `platform` | text | Target platform |
| `caption` | text | Platform-specific caption |
| `hashtags` | json | Array of hashtags |
| `status` | text | `draft`, `approved`, `posted` |

---

### guest_prep_podcast_profiles

Podcast identity configuration.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `podcast_name` | text | Name of the podcast |
| `host_name` | text | Host name(s) |
| `description` | text | Podcast description |
| `target_audience` | text | Who listens |
| `typical_episode_length` | integer | Minutes |
| `interview_style` | text | Conversational, structured, etc. |

---

### guest_prep_research_sessions

Guest research data.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `guest_name` | text | Guest's name |
| `guest_bio` | text | Biography |
| `guest_links` | json | Social/website links |
| `research_notes` | text | Research findings |
| `suggested_questions` | json | AI-suggested questions |
| `episode_date` | date | Scheduled recording date |
| `status` | text | `researching`, `ready`, `completed` |

---

### guest_prep_ros_versions

Run of show version history.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `research_session_id` | integer | FK to research_sessions |
| `version_number` | integer | Version number |
| `content` | json | Full run of show content |
| `changes_summary` | text | What changed |

---

### briefing_podcast_profiles

Briefing app podcast profiles.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `podcast_name` | text | Podcast name |
| `description` | text | Description |
| `host_info` | json | Host details |
| `branding` | json | Brand guidelines |

---

### briefing_research_sessions

Briefing research sessions.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `profile_id` | integer | FK to podcast_profiles |
| `guest_info` | json | Guest details |
| `research_data` | json | Compiled research |
| `briefing_doc` | text | Generated briefing |
| `status` | text | Session status |

---

### dashboard_activity

Activity feed for the dashboard.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `activity_type` | text | Type of activity |
| `description` | text | Activity description |
| `metadata` | json | Additional data |
| `user_id` | text | Who performed it |

---

### outreach_leads

Discovered podcast creator leads.

| Column | Type | Description |
|--------|------|-------------|
| `id` | integer | Auto-generated ID |
| `created_at` | timestamp | Creation timestamp |
| `name` | text | Lead name |
| `email` | text | Email address |
| `podcast_name` | text | Their podcast |
| `source` | text | Where found |
| `status` | text | `new`, `contacted`, `responded` |
| `notes` | text | Notes about the lead |

---

## Code Examples

### Python

```python
import requests

BASE_URL = "https://audos.com/api/hooks/execute/workspace-351699/db-api"

# List all tables
response = requests.post(BASE_URL, json={"action": "list-tables"})
tables = response.json()["tables"]

# Query with filters
response = requests.post(BASE_URL, json={
    "action": "query",
    "table": "studio_episodes",
    "filters": [
        {"column": "status", "operator": "eq", "value": "ready"}
    ],
    "orderBy": {"column": "publish_date", "direction": "desc"},
    "limit": 5
})
episodes = response.json()["rows"]

# Insert a new record
response = requests.post(BASE_URL, json={
    "action": "insert",
    "table": "studio_episodes",
    "data": {
        "title": "New Episode",
        "description": "About something cool",
        "status": "draft"
    }
})
new_id = response.json()["rows"][0]["id"]

# Update a record
response = requests.post(BASE_URL, json={
    "action": "update",
    "table": "studio_episodes",
    "filters": [{"column": "id", "operator": "eq", "value": new_id}],
    "data": {"status": "ready"}
})
```

### JavaScript

```javascript
const BASE_URL = "https://audos.com/api/hooks/execute/workspace-351699/db-api";

async function dbRequest(body) {
  const response = await fetch(BASE_URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });
  return response.json();
}

// List tables
const { tables } = await dbRequest({ action: "list-tables" });

// Query with filters
const { rows: episodes } = await dbRequest({
  action: "query",
  table: "studio_episodes",
  filters: [{ column: "status", operator: "eq", value: "ready" }],
  orderBy: { column: "publish_date", direction: "desc" },
  limit: 5
});

// Insert
const { rows: [newEpisode] } = await dbRequest({
  action: "insert",
  table: "studio_episodes",
  data: { title: "New Episode", status: "draft" }
});

// Update
await dbRequest({
  action: "update",
  table: "studio_episodes",
  filters: [{ column: "id", operator: "eq", value: newEpisode.id }],
  data: { status: "ready" }
});
```

---

## Creating New Tables

New tables you create (via the platform or by asking the AI assistant) are **immediately available** through this same API. No new endpoints needed.

```bash
# After creating a new table called "my_custom_table"
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/db-api" \
  -H "Content-Type: application/json" \
  -d '{"action": "query", "table": "my_custom_table"}'
```

Use `list-tables` to see all available tables at any time.

---

*Part of the [Throughline API](./README.md)*
