# Throughline API Documentation

This folder contains documentation for all external APIs available for the Throughline workspace.

## Available APIs

| API | Endpoint | Description |
|-----|----------|-------------|
| [Database API](./database-api.md) | `/db-api` | Full CRUD access to all workspace tables |
| [AI API](./ai-api.md) | `/ai-api` | AI text generation (GPT-4o-mini) |
| [Email API](./email-api.md) | `/email-api` | Send transactional emails |
| [Storage API](./storage-api.md) | `/storage-api` | File upload and management |
| [Scheduler API](./scheduler-api.md) | `/scheduler-api` | Cron jobs and scheduled tasks |
| [Web API](./web-api.md) | `/web-api` | Web scraping and search |
| [CRM API](./crm-api.md) | `/crm-api` | Contact and lead management |
| [Analytics API](./analytics-api.md) | `/analytics-api` | Visitor metrics and funnel data |

## Base URL

All APIs use the same base URL pattern:

```
https://audos.com/api/hooks/execute/workspace-351699/{endpoint}
```

## Authentication

Currently, these endpoints are open (no API key required). They are scoped to the Throughline workspace.

## Quick Start

```bash
# List all database tables
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/db-api" \
  -H "Content-Type: application/json" \
  -d '{"action": "list-tables"}'

# Generate AI content
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/ai-api" \
  -H "Content-Type: application/json" \
  -d '{"action": "generate", "prompt": "Write a LinkedIn post about podcasting"}'

# List contacts
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/crm-api" \
  -H "Content-Type: application/json" \
  -d '{"action": "list", "limit": 10}'
```

## Workspace Info

- **Workspace ID:** `8f1ad824-832f-4af8-b77e-ab931a250625`
- **Workspace Number:** `351699`
- **Live URLs:**
  - Landing Page: https://www.trythroughline.com
  - App/Space: https://app.trythroughline.com

---

*Last updated: 2026-03-31*
