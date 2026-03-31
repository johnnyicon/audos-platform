# Audos REST API Reference

All APIs are accessible via server functions (hooks). These work from anywhere - local development, external apps, or mobile apps.

## Base URL Pattern

```
POST https://audos.app/api/hooks/execute/workspace-{WORKSPACE_ID}/{HOOK_NAME}
```

**Workspace ID:** `8f1ad824-832f-4af8-b77e-ab931a250625`

**Or use your custom domain:**
```
POST https://trythroughline.com/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/{HOOK_NAME}
```

---

## 1. Database API (`db-api`)

Full CRUD operations on workspace tables.

### List Tables

```json
POST /api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/db-api

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
      "tableName": "voice_profiles",
      "displayName": "Voice Profiles",
      "description": "Stores voice fingerprints...",
      "columns": [...]
    }
  ]
}
```

### Query Data

```json
{
  "action": "query",
  "table": "voice_profiles",
  "filters": [
    { "column": "type", "operator": "eq", "value": "host" }
  ],
  "orderBy": { "column": "created_at", "direction": "desc" },
  "limit": 10
}
```

**Filter Operators:** `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `like`, `ilike`, `in`, `is_null`, `not_null`

### Insert Data

```json
{
  "action": "insert",
  "table": "voice_profiles",
  "rows": [
    {
      "name": "John Gonzales",
      "type": "host",
      "description": "Personal voice for John's content"
    }
  ]
}
```

### Update Data

```json
{
  "action": "update",
  "table": "voice_profiles",
  "filters": [
    { "column": "id", "operator": "eq", "value": 1 }
  ],
  "data": {
    "is_trained": true,
    "last_trained_at": "2025-06-10T12:00:00Z"
  }
}
```

### Delete Data

```json
{
  "action": "delete",
  "table": "voice_profiles",
  "filters": [
    { "column": "id", "operator": "eq", "value": 1 }
  ]
}
```

### Describe Table

```json
{
  "action": "describe",
  "table": "voice_profiles"
}
```

---

## 2. AI API (`ai-api`)

Generate text using AI.

```json
POST /api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/ai-api

{
  "action": "generate",
  "prompt": "Write a short podcast intro for a show about AI technology",
  "system": "You are a podcast script writer. Keep it conversational."
}
```

**Response:**
```json
{
  "success": true,
  "text": "Hey everyone, welcome back to...",
  "model": "gpt-4o-mini-2024-07-18",
  "usage": {
    "promptTokens": 25,
    "completionTokens": 50,
    "totalTokens": 75
  }
}
```

---

## 3. Email API (`email-api`)

Send transactional emails.

```json
POST /api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/email-api

{
  "action": "send",
  "to": "user@example.com",
  "subject": "Your Guest Research is Ready",
  "text": "Hi there, your guest research packet has been generated...",
  "html": "<h1>Your Guest Research is Ready</h1><p>...</p>"
}
```

---

## 4. Web API (`web-api`)

Fetch and parse web pages.

```json
POST /api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/web-api

{
  "action": "fetch",
  "url": "https://example.com/guest-bio"
}
```

**Response:**
```json
{
  "success": true,
  "content": "...parsed markdown content...",
  "title": "Page Title",
  "meta": { ... }
}
```

---

## 5. Storage API (`storage-api`)

Upload and manage files.

```json
POST /api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/storage-api

{
  "action": "list",
  "category": "attachment"
}
```

---

## 6. Scheduler API (`scheduler-api`)

Create cron jobs and scheduled tasks.

```json
POST /api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/scheduler-api

{
  "action": "create",
  "name": "weekly-digest",
  "frequency": "weekly",
  "time": "09:00",
  "timezone": "America/Los_Angeles",
  "hookName": "send-digest"
}
```

---

## 7. Analytics API (`analytics-api`)

Get visitor metrics and funnel data.

```json
POST /api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/analytics-api

{
  "action": "overview",
  "days": 30
}
```

---

## 8. CRM API (`crm-api`)

Manage contacts and leads.

```json
POST /api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/crm-api

{
  "action": "list-contacts",
  "limit": 50
}
```

---

## TypeScript Helper Class

```typescript
// lib/audos-api.ts

const WORKSPACE_ID = '8f1ad824-832f-4af8-b77e-ab931a250625';
const API_BASE = process.env.NERX_PUBLIC_API_BASE || 'https://audos.app';

export class AudosAPI {
  private async callHook(hookName: string, body: any) {
    const response = await fetch(
      `${API_BASE}/api/hooks/execute/workspace-${WORKSPACE_ID}/${hookName}`,
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      }
    );
    return response.json();
  }

  // Database
  async listTables() {
    return this.callHook('db-api', { action: 'list-tables' });
  }

  async query(table: string, options: { filters?: any[]; limit?: number; orderBy?: any } = {}) {
    return this.callHook('db-api', { action: 'query', table, ...options });
  }

  async insert(table: string, rows: any[]) {
    return this.callHook('db-api', { action: 'insert', table, rows });
  }

  async update(table: string, filters: any[], data: any) {
    return this.callHook('db-api', { action: 'update', table, filters, data });
  }

  async delete(table: string, filters: any[]) {
    return this.callHook('db-api', { action: 'delete', table, filters });
  }

  // AI
  async generateText(prompt: string, system?: string) {
    return this.callHook('ai-api', { action: 'generate', prompt, system });
  }

  // Email
  async sendEmail(to: string, subject: string, text: string, html?: string) {
    return this.callHook('email-api', { action: 'send', to, subject, text, html });
  }

  // Web
  async fetchUrl(url: string) {
    return this.callHook('web-api', { action: 'fetch', url });
  }
}

export const audos = new AudosAPI();
```

## Usage Example

```typescript
import { audos } from './lib/audos-api';

// Query voice profiles
const { data } = await audos.query('voice_profiles', {
  filters: [{ column: 'type', operator: 'eq', value: 'host' }],
  limit: 10
});

// Generate content
const { text } = await audos.generateText(
  'Write a podcast intro for a guest named Jane Doe',
  'You are a podcast script writer.'
);

// Send email
const result = await audos.sendEmail(
  'guest@example.com',
  'Your Guest Research is Ready',
  'Hi, your research packet is ready!'
);
```