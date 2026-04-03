# Audos Database SDK

### Comprehensive Database Management Guide

**Document Version:** 1.0  
**Date:** June 2025  
**Workspace:** Throughline  
**Author:** Otto (AI Assistant)

**Purpose:** Complete reference for all database operations available through the Audos platform, including Otto commands, REST API endpoints, app integration, capabilities, limitations, and workarounds.

---

## Table of Contents

1. [Overview](#overview)
2. [Access Methods Summary](#access-methods-summary)
3. [Method 1: Otto (MCP Tools)](#method-1-otto-mcp-tools)
4. [Method 2: REST API via Server Functions](#method-2-rest-api-via-server-functions)
5. [Method 3: App Code (Runtime SDK)](#method-3-app-code-runtime-sdk)
6. [Supported Data Types](#supported-data-types)
7. [Capabilities Matrix](#capabilities-matrix)
8. [Limitations & Workarounds](#limitations--workarounds)
9. [Real-World Examples](#real-world-examples)
10. [Best Practices](#best-practices)

---

## Overview

The Audos platform provides a **managed PostgreSQL database** for each workspace. Unlike traditional database access where you run arbitrary SQL, operations go through structured APIs.

### Key Concepts

- **Workspace Isolation**: Each workspace has its own isolated schema (`ws_{workspace_id}`)
- **Automatic Columns**: Every table gets `id` (serial PK) and `created_at` (timestamp) automatically
- **Session Scoping**: Data can be scoped to user sessions or shared across all sessions

---

## Access Methods Summary

| Method | Schema Management | CRUD Operations | Raw SQL | External Access |
|--------|-------------------|-----------------|---------|-----------------|
| **Otto (MCP TOols)** | ✅✅✅ Full support | ✅✅✅ Full support | ✅ SELECT only | ❌ Otto session only |
| **REST API (via hooks)** | ❌ Not available | ✅ Full support | ✅ SELECT only | ✅ Callable externally |
| **App Code (Runtime)** | ❌ Not available | ✅ Full support | ❌ Not available | ❌ App runtime only |

**Bottom line:**
- **Schema management** (create/alter/drop tables) — **Only through Otto**
- **CRUD operations** (insert/update/delete/query) — **All three methods**
- **External API access** — **REST API via server functions**

---

## Method 1: Otto (MCP Tools)

This is the **most powerful** method — it's the only way to manage schema.

### Schema Management APIs

#### 1. `db_create_table`

Creates a new table with specified columns.

```typescript
db_create_table({
  name: "products",                     // Required: lowercase, underscores
  displayName: "Products",              // Optional: human-friendly name
  description: "Product catalog",       // Optional
  columns: [                            // Required
    {
      name: "name",
      type: "text",
      nullable: false,                  // Default: true
      description: "Product name"
    },
    {
      name: "price",
      type: "decimal",
      nullable: true
    },
    {
      name: "sku",
      type: "text",
      unique: true                      // Default: false
    },
    {
      name: "is_active",
      type: "boolean",
      defaultValue: "true"              // SQL expression as string
    },
    {
      name: "metadata",
      type: "json"
    }
  ],
  foreignKeys: [                        // Optional: ONLY at creation time!
    {
      column: "category_id",
      referencesTable: "categories",
      referencesColumn: "id",           // Default: "id"
      onDelete: "SET NULL"              // CASCADE, SET NULL, RESTRICT, NO ACTION
    }
  ]
})
```

#### 2. `db_alter_table`

Modify an existing table's structure.

```typescript
// Add columns
db_alter_table({
  table: "my_table",
  changes: {
    addColumns: [
      { name: "user_id", type: "text", nullable: true },
      { name: "org_id", type: "text", nullable: true }
    ]
  }
})

// Rename columns
db_alter_table({
  table: "my_table",
  changes: {
    renameColumns: [
      { from: "old_name", to: "new_name" }
    ]
  }
})

// Drop columns (requires confirmation)
db_alter_table({
  table: "my_table",
  changes: {
    dropColumns: ["column_to_remove"]
  }
})

// Add indexes
db_alter_table({
  table: "my_table",
  changes: {
    addIndexes: [
      { columns: ["user_id", "org_id"], unique: false }
    ]
  }
})
```

#### 3. `db_list_tables`

List all tables in the workspace.

```typescript
db_list_tables()
// Returns: [{ name, displayName, rowCount, columnCount }, ...]
```

#### 4. `db_describe_table`

Get full schema details.

```typescript
db_describe_table({ table: "products" })
// Returns: columns, types, constraints, indexes, sample data
```

#### 5. `db_drop_table` ⚠️

Permanently delete a table (two-step confirmation).

```typescript
// Step 1: Request
const result = db_drop_table({ table: "old_table" })
// Returns: { confirmationToken: "abc123" }

// Step 2: Confirm
dbEconfirm_destructive({ token: "abc123" })
```

#### 6. `db_truncate_table` ⚠️

Delete all rows but keep structure (two-step confirmation).

```typescript
db_truncate_table({ table: "log_entries" })
// Then confirm with db_confirm_destructive()
```

---

### Data Operation APIs (Otto)

#### `db_insert`

```typescript
db_insert({
  table: "products",
  rows: [
    { name: "Widget A", price: 29.99, is_active: true },
    { name: "Widget B", price: 39.99, is_active: true }
  ]
})
// Max 100 rows per call
```

#### `db_update`

```typescript
db_update({
  table: "products",
  filters: [
    { column: "id", operator: "eq", value: 42 }
  ],
  data: { price: 34.99, is_active: false }
})
// At least one filter required (safety)
```

#### `db_delete`

```typescript
db_delete({
  table: "products",
  filters: [
    { column: "is_active", operator: "is_false" }
  ]
})
// No filters = delete ALL (requires confirmation)
```

#### `db_query`

```typescript
db_query({
  table: "products",
  columns: ["id", "name", "price"],
  filters: [
    { column: "is_active", operator: "is_true" }
  ],
  orderBy: { column: "price", direction: "desc" },
  limit: 20,
  offset: 0
})

// Aggregate query
db_query({
  table: "orders",
  aggregate: { function: "sum", column: "total" },
  groupBy: ["user_id"]
})
```

#### `workspace_execute_sql`

Run arbitrary **read-only** SQL (SELECT only).

```typescript
workspace_execute_sql({
  query: `
    SELECT p.name, COUNT(o.id) as order_count
    FROM products p
    LEFT JOIN orders o ON o.product_id = p.id
    GROUP BY p.id
    ORDER BY order_count DESC
  `,
  limit: 50
})
```

**Note:** No schema prefix needed — just use table names directly.

---

### Filter Operators

| Operator | Meaning | Example |
|----------|---------|---------|
| `eq` | Equals | `{ column: "id", operator: "eq", value: 42 }` |
| `neq` | Not equals | `{ column: "status", operator: "neq", value: "draft" }` |
| `gt`, `gte` | Greater than (/or equal) | `{ column: "price", operator: "gt", value: 100 }` |
| `lt`, `lte` | Less than (/or equal) | `{ column: "price", operator: "lte", value: 50 }` |
| `like` | Pattern match (case sensitive) | `{ column: "name", operator: "like", value: "%Widget%" }` |
| `ilike` | Pattern match (case insensitive) | `{ column: "name", operator: "ilike", value: "%widget%" }` |
| `in` | In array | `{ column: "status", operator: "in", value: ["active", "pending"] }` |
| `is_null` | Is null | `{ column: "deleted_at", operator: "is_null" }` |
| `not_null` | Is not null | `{ column: "email", operator: "not_null" }` |
| `is_true` | Is true | `{ column: "is_active", operator: "is_true" }` |
| `is_false` | Is false | `{ column: "is_active", operator: "is_false" }` |

---

## Method 2: REST API via Server Functions

This is how you access the database from **external code** (local scripts, other services, etc.).

### The `db-api` Server Function

A server function (hook) called `db-api` already exists in your workspace. It exposes CRUD operations via HTTP.

#### Base URL

```
https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api
```

#### Authentication

No authentication required — the endpoint is public but scoped to your workspace.

---

### Available Actions

#### 1. `list-tables`

```bash
curl -X POST https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api \
  -H "Content-Type: application/json" \
  -d '{"action": "list-tables"}'
```

**Response:**
```json
{
  "success": true,
  "tables": [
    {
      "tableName": "voice_profiles",
      "displayName": "Voice Profiles",
      "columns": [...]
    },
    ...
  ]
}
```

---

#### 2. `describe`

```bash
curl -X POST https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api \
  -H "Content-Type: application/json" \
  -d '{
    "action": "describe",
    "table": "reels"
  }'
```

**Response:**
```json
{
  "success": true,
  "table": "reels",
  "columns": [
    { "name": "id", "type": "serial", "nullable": false, "primaryKey": true },
    { "name": "title", "type": "text", "nullable": false },
    ...
  ]
}
```

---

#### 3. `query`

```bash
curl -X POST https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api \
  -H "Content-Type: application/json" \
  -d '{
    "action": "query",
    "table": "reels",
    "columns": ["id", "title", "status"],
    "filters": [
      { "column": "status", "operator": "eq", "value": "active" }
    ],
    "orderBy": { "column": "created_at", "direction": "desc" },
    "limit": 10
  }'
```

**Response:**
```json
{
  "success": true,
  "data": {
    "rows": [
      { "id": 1, "title": "Why DonorsChoose Exists", "status": "active" },
      ...
    ],
    "rowCount": 5
  }
}
```

---

#### 4. `insert`

```bash
curl -X POST https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api \
  -H "Content-Type: application/json" \
  -d '{
    "action": "insert",
    "table": "reels",
    "rows": [
      { "title": "New Reel", "status": "draft" }
    ]
  }'
```

**Response:**
```json
{
  "success": true,
  "inserted": 1
}
```

---

#### 5. `update`

```bash
curl -X POST https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api \
  -H "Content-Type: application/json" \
  -d '{
    "action": "update",
    "table": "reels",
    "filters": [
      { "column": "id", "operator": "eq", "value": 1 }
    ],
    "data": {
      "status": "published"
    }
  }'
```

**Response:**
```json
{
  "success": true,
  "updated": 1
}
```

---

#### 6. `delete`

```bash
curl -X POST https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api \
  -H "Content-Type: application/json" \
  -d '{
    "action": "delete",
    "table": "reels",
    "filters": [
      { "column": "id", "operator": "eq", "value": 999 }
    ]
  }'
```

**Response:**
```json
{
  "success": true,
  "deleted": 1
}
```

---

### Python SDK Example

```python
import requests

BASE_URL = "https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api"

class ThroughlineDB:
    def __init__(self):
        self.base_url = BASE_URL
    
    def _request(self, payload):
        response = requests.post(self.base_url, json=payload)
        return response.json()
    
    def list_tables(self):
        return self._request({"action": "list-tables"})
    
    def describe(self, table):
        return self._request({"action": "describe", "table": table})
    
    def query(self, table, columns=None, filters=None, order_by=None, limit=50):
        payload = {"action": "query", "table": table, "limit": limit}
        if columns:
            payload["columns"] = columns
        if filters:
            payload["filters"] = filters
        if order_by:
            payload["orderBy"] = order_by
        return self._request(payload)
    
    def insert(self, table, rows):
        return self._request({"action": "insert", "table": table, "rows": rows})
    
    def update(self, table, filters, data):
        return self._request({"action": "update", "table": table, "filters": filters, "data": data})
    
    def delete(self, table, filters):
        return self._request({"action": "delete", "table": table, "filters": filters})


# Usage
db = ThroughlineDB()

# List all tables
tables = db.list_tables()
print(tables)

# Query reels
reels = db.query(
    table="reels",
    columns=["id", "title", "status"],
    filters=[{"column": "status", "operator": "eq", "value": "active"}],
    order_by={"column": "created_at", "direction": "desc"},
    limit=10
)
print(reels)

# Insert a new reel
result = db.insert("reels", [{"title": "My New Reel", "status": "draft"}])
print(result)
```

---

### TypeScript/Node.js SDK Example

```typescript
const BASE_URL = "https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api";

type FilterOperator = "eq" | "neq" | "gt" | "gte" | "lt" | "lte" | "like" | "ilike" | "in" | "is_null" | "not_null" | "is_true" | "is_false";

interface Filter {
  column: string;
  operator: FilterOperator;
  value?: any;
}

interface OrderBy {
  column: string;
  direction: "asc" | "desc";
}

class ThroughlineDB {
  private baseUrl: string;

  constructor() {
    this.baseUrl = BASE_URL;
  }

  private async request<T>(payload: Record<string, any>): Promise<T> {
    const response = await fetch(this.baseUrl, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    return response.json();
  }

  async listTables() {
    return this.request({ action: "list-tables" });
  }

  async describe(table: string) {
    return this.request({ action: "describe", table });
  }

  async query<T = any>(options: {
    table: string;
    columns?: string[];
    filters?: Filter[];
    orderBy?: OrderBy;
    limit?: number;
  }): Promise<{ success: boolean; data: { rows: T[]; rowCount: number } }> {
    return this.request({ action: "query", ...options });
  }

  async insert(table: string, rows: Record<string, any>[]) {
    return this.request({ action: "insert", table, rows });
  }

  async update(table: string, filters: Filter[], data: Record<string, any>) {
    return this.request({ action: "update", table, filters, data });
  }

  async delete(table: string, filters: Filter[]) {
    return this.request({ action: "delete", table, filters });
  }
}

// Usage
const db = new ThroughlineDB();

// Query reels
const reels = await db.query({
  table: "reels",
  columns: ["id", "title", "status"],
  filters: [{ column: "status", operator: "eq", value: "active" }],
  orderBy: { column: "created_at", direction: "desc" },
  limit: 10,
});

console.log(reels.data.rows);
```

---

### What's NOT Available via REST API

The `db-api` hook **does not** support schema management:

| Action | Status |
|--------|--------|
| `create-table` | ❌ Not available |
| `alter-table` | ❌ Not available |
| `drop-table` | ❌ Not available |
| `truncate-table` | ❌ Not available |

Schema management **must** be done through Otto.

### Why?

The underlying `db` object in server functions only exposes:
- `query`
- `insert`
- `update`
- `delete`
- `listTables`
- `rawQuery` (SELECT only — DDL/DML is blocked)

There is no `createTable`, `alterTable`, etc. exposed to server functions.

---

## Method 3: App Code (Runtime SDK)

Workspace apps have an auto-injected SDK for database access.

### React Hook: `useWorkspaceDB`

```typescript
// Read data with reactive updates
const { data, loading, error, refetch } = useWorkspaceDB('reels', {
  shared: true,  // IMPORTANT: read all data, not just current session
  filters: [{ column: 'status', operator: 'eq', value: 'active' }],
  orderBy: { column: 'created_at', direction: 'desc' }
});

if (loading) return <div>Loading...</div>;
if (error) return <div>Error: {error.message}</div>;

return (
  <ul>
    {data?.map(reel => (
      <li key={reel.id}>{reel.title}</li>
    ))}
  </ul>
);
```

#### Hook Options

| Option | Type | Description |
|--------|------|-------------|
| `shared` | `boolean` | If `true`, reads all data; if `false`, only current session's data |
| `filters` | `Filter[]` | Array of filter conditions |
| `orderBy` | `{ column, direction }` | Sort order |
| `limit` | `number` | Max rows to return |

---

### Imperative API: `window.__workspaceDb`

```typescript
// Insert
const newRow = await window.__workspaceDb.from('reels').insert({
  title: 'My New Reel',
  status: 'draft'
});

// Update
await window.__workspaceDb.from('reels').update(
  { id: 42 },  // filter
  { status: 'published' }  // data
);

// Delete
await window.__workspaceDb.from('reels').delete({ id: 42 });

// Query
const rows = await window.__workspaceDb.from('reels').get({
  shared: true,
  filters: [{ column: 'status', operator: 'eq', value: 'active' }]
});
```

#### Important: `shared: true`

Data inserted via Otto or the REST API has `session_id = NULL`. 

By default, apps only see data where `session_id` matches the current user's session.

To read data inserted by Otto or the API, **you must use `shared: true`**.

---

## Supported Data Types

| Type | Description | Example Use Case |
|------|-------------|------------------|
| `text` | Variable-length string | Names, descriptions, IDs |
| `integer` | 32-bit integer | Counts, small numbers |
| `bigint` | 64-bit integer | Large numbers, timestamps as ms |
| `decimal` | Arbitrary precision | Prices, financial data |
| `boolean` | true/false | Flags, status |
| `timestamp` | Date + time | Created at, scheduled at |
| `date` | Date only | Birthdays, release dates |
| `json` | JSON object/array | Metadata, settings |
| `uuid` | Universally unique ID | External references |

---

## Capabilities Matrix

### By Access Method

| Operation | Otto | REST API | App Code |
|-----------|------|----------|----------|
| Create table | ✅ | ❌ | ❌ |
| Alter table | ✅ | ❌ | ❌ |
| Drop table | ✅ | ❌ | ❌ |
| Truncate table | ✅ | ❌ | ❌ |
| List tables | ✅ | ✅ | ❌ |
| Describe table | ✅ | ✅ | ❌ |
| Query | ✅ | ✅ | ✅ |
| Insert | ✅ | ✅ | ✅ |
| Update | ✅ | ✅ | ✅ |
| Delete | ✅ | ✅ | ✅ |
| Raw SELECT | ✅ | ❌ | ❌ |

### By Operation Type

| Operation | Supported | Notes |
|-----------|-----------|-------|
| Create tables | ✅ | Otto only |
| Add columns | ✅ | Otto only |
| Rename columns | ✅ | Otto only |
| Drop columns | ✅ | Otto only, requires confirmation |
| Add indexes | ✅ | Otto only |
| Foreign keys | ✅ | At table creation only |
| **Change column type** | ❌ | Workaround available |
| **Change nullability** | ❌ | Workaround available |
| **Add FK to existing table** | ❌ | Recreate table |
| **Raw DDL** | ❌ | Must use structured APIs |

---

## Limitations & Workarounds

### 1. Changing Column Type

The API does not expose `ALTER COLUMN ... TYPE`. 

**Workaround:**

```
1. Add a new column with the desired type
   db_alter_table({ table: "my_table", changes: { addColumns: [{ name: "new_col", type: "integer" }] }})

2. Migrate data
   "Update all rows in my_table, set new_col = cast old_col as integer"

3. Drop the old column
   db_alter_table({ table: "my_table", changes: { dropColumns: ["old_col"] }})

4. Rename the new column
   db_alter_table({ table: "my_table", changes: { renameColumns: [{ from: "new_col", to: "old_col" }] }})
```

---

### 2. Changing Nullability

Same workaround as above — create a new column with the desired nullability, migrate data, drop old, rename new.

---

### 3. Adding Foreign Keys to Existing Tables

Foreign keys can only be defined at table creation time.

**Workaround:**
1. Create a new table with the FK
2. Migrate all data from the old table
3. Drop the old table
4. (Optional) Rename the new table to the original name

---

### 4. Running Raw DDL/DML

The `db.rawQuery()` method in server functions is **restricted to SELECT only**.

```javascript
// This works:
await db.rawQuery('SELECT * FROM reels');

// This FAILS:
await db.rawQuery('CREATE TABLE foo (...)');
// Error: "Raw queries are restricted to SELECT, WITH (CTE), and EXPLAIN statements only."
```

---

## Real-World Examples

### Example 1: Adding Multi-Tenancy Columns

**Request to Otto:** 
> "Add user_id and org_id columns to these tables: guest_prep_podcast_profiles, voice_profiles, speakers..."

**What Otto Does:**
```typescript
// Runs 11 parallel db_alter_table calls
db_alter_table({ table: "guest_prep_podcast_profiles", changes: { addColumns: [{ name: "user_id", type: "text", nullable: true }, { name: "org_id", type: "text", nullable: true }] }})
// ... repeat for all tables
```

**Result:** All 11 tables updated in seconds. Existing rows have NULL values.

---

### Example 2: Creating a Table with Relationships

**Request to Otto:**
> "Create a reel_captions table that references reels and voice_profiles"

**What Otto Does:**
```typescript
db_create_table({
  name: "reel_captions",
  columns: [
    { name: "reel_id", type: "integer" },
    { name: "platform", type: "text" },
    { name: "voice_profile_id", type: "integer" },
    { name: "caption_text", type: "text" },
    { name: "status", type: "text", defaultValue: "'draft'" }
  ],
  foreignKeys: [
    { column: "reel_id", referencesTable: "reels", onDelete: "CASCADE" },
    { column: "voice_profile_id", referencesTable: "voice_profiles", onDelete: "SET NULL" }
  ]
})
```

---

### Example 3: Complex Query via REST API

```bash
# Get reels with their caption counts
curl -X POST https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api \
  -H "Content-Type: application/json" \
  -d '{
    "action": "query",
    "table": "reels",
    "columns": ["id", "title", "status", "created_at"],
    "filters": [
      { "column": "status", "operator": "in", "value": ["active", "pending"] }
    ],
    "orderBy": { "column": "created_at", "direction": "desc" },
    "limit": 20
  }'
```

---

## Best Practices

### 1. Plan Your Schema Upfront

Since changing column types is cumbersome, try to get your schema right the first time.

### 2. Use Nullable Columns for New Fields

When adding columns to tables with existing data, always make them nullable initially. Backfill later.

### 3. Use `json` for Flexible Data

If you're not sure what structure you'll need, use a `json` column.

### 4. Back Up Before Destructive Operations

Use `db_create_backup` before dropping columns or tables.

### 5. Use `shared: true` in Apps

If you're inserting data via Otto or the REST API, remember that apps need `shared: true` to see it.

---

## Summary

| What You Want | How to Do It |
|---------------|---------------|
| Create/alter/drop tables | Ask Otto |
| CRUD from external code | REST API via `db-api` hook |
| CRUD from apps | `useWorkspaceDB` or window.__workspaceDb` |
| Complex SELECT | Ask Otto to use `workspace_execute_sql` |
| Change column type | Ask Otto (workaround) |

---

## Appendix: API Endpoint Reference

### REST API (`db-api` hook)

**Base URL:**
```
https://audos.ai/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api
```

**Method:** `POST`

**Content-Type:** `application/json`

**Authentication:** None required

**Actions:**
| Action | Required Params | Optional Params |
|--------|-----------------|-----------------|
| `list-tables` | — | — |
| `describe` | `table` | — |
| `query` | `table` | `columns`, `filters`, `orderBy`, `limit` |
| `insert` | `table`, `rows` | — |
| `update` | `table`, `filters`, `data` | — |
| `delete` | `table`, `filters` | — |

---

**SDK Document ID:** SDK-12  
**Related Documents:** SDK-09 (REST APIs), SDK-10 (GitHub Sync), SDK-11 (Bidirectional Sync)