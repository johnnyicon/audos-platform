# Audos Development Workflow Guide

> **SDK-Level Documentation** for developers working with the Audos platform

This guide explains how to develop effectively in the Audos ecosystem by understanding when to work locally versus when to use the platform, and how to make both environments work together seamlessly.

---

## Table of Contents

1. [Overview: Two Development Environments](#overview-two-development-environments)
2. [Local vs Platform Development](#local-vs-platform-development)
3. [The Hybrid Workflow](#the-hybrid-workflow)
4. [Setting Up Local Development](#setting-up-local-development)
5. [Working with Platform APIs](#working-with-platform-apis)
6. [When to Go Back to Otto](#when-to-go-back-to-otto)
7. [Testing Strategies](#testing-strategies)
8. [Best Practices](#best-practices)
9. [Common Patterns](#common-patterns)

---

## Overview: Two Development Environments

Audos development happens in **two complementary environments**:

1. **Local Development** - Your machine, your IDE, your debugging tools
2. **Platform Development** - The Audos platform, managed by Otto (the AI agent)

These environments are **not competing alternatives** - they work together. The platform hosts your production infrastructure (database, apps, server functions), while local development gives you the speed and control to build, test, and iterate on business logic.

### The Mental Model

Think of the Audos platform as your **backend infrastructure** and local development as your **business logic layer**:

```
┌─────────────────────────────────────────────────────────┐
│  Local Development (Your Machine)                       │
│  • Scripts & business logic                             │
│  • Data analysis & processing                           │
│  • Testing & debugging                                  │
│  • CLI tools & automation                               │
└──────────────────┬──────────────────────────────────────┘
                   │ HTTP APIs
                   │
┌──────────────────▼──────────────────────────────────────┐
│  Audos Platform (Cloud)                                 │
│  • Database tables (PostgreSQL)                         │
│  • Server functions (webhooks, endpoints)               │
│  • React apps (hosted, compiled)                        │
│  • Landing pages                                        │
│  • CRM, analytics, payments                             │
└─────────────────────────────────────────────────────────┘
```

---

## Local vs Platform Development

### Use Local Development For:

**Scripts that call APIs**
- Data pipelines that read from your database and write results
- Automation scripts that sync with external services
- CLI tools for managing workspace data
- Testing and prototyping new features

**Rapid iteration**
- Hot reload during development
- Full IDE debugging (breakpoints, inspectors)
- Git version control on your terms
- No compilation wait times

**Data processing and analysis**
- Python scripts for analytics
- Data transformations and exports
- Scheduled jobs (that call platform APIs)

**Integration testing**
- Test API behavior without affecting production
- Mock external services
- Validate data flows end-to-end

### Use Platform Development (via Otto) For:

**New server functions (hooks)**
- Creating webhook receivers
- Building API endpoints for your apps
- Scheduled automations (via Task Scheduler + hooks)
- Phone agent tools (mid-call data lookup)

**Database changes**
- Creating new tables (`db_create_table`)
- Modifying table schemas (`db_alter_table`)
- Setting up foreign keys and indexes

**React app changes**
- Building or editing Space apps (mini-apps)
- Modifying the landing page
- Updating the EmailGate (entry screen)
- Changing the Desktop layout

**Infrastructure configuration**
- Setting up custom domains
- Configuring email sending
- Managing team members
- Connecting integrations (Stripe, Meta Ads, etc.)

---

## The Hybrid Workflow

The most powerful Audos development pattern combines both environments.

### How It Works

1. **Platform hosts the infrastructure** - Otto creates database tables, server functions, and React apps
2. **Local scripts call HTTP APIs** - Your Python/TypeScript code reads and writes data via REST APIs
3. **GitHub sync is ONE-WAY** - Changes pushed to GitHub flow into Audos, but Audos doesn't push back
4. **Server functions bridge the gap** - Otto creates endpoints that your local scripts can call

### Example Workflow: Building a Podcast Guest Management System

Let's walk through a real example:

#### Step 1: Design the data model (Platform)

Talk to Otto:
```
"I need a database table for podcast guests with name, email, status, episode_date, and notes."
```

Otto creates the table using `db_create_table`:
```typescript
// Otto runs this behind the scenes
await db_create_table({
  name: "podcast_guests",
  displayName: "Podcast Guests",
  columns: [
    { name: "name", type: "text", nullable: false },
    { name: "email", type: "text", nullable: false },
    { name: "status", type: "text", nullable: false },
    { name: "episode_date", type: "date", nullable: true },
    { name: "notes", type: "text", nullable: true }
  ]
});
```

#### Step 2: Create a server function for guest lookup (Platform)

Talk to Otto:
```
"Create a server function called 'lookup-guest' that takes an email and returns the guest's record."
```

Otto creates the hook:
```javascript
// Server function: lookup-guest
const { email } = request.body;

const result = await db.query('podcast_guests', {
  where: { email },
  limit: 1
});

if (result.rowCount > 0) {
  respond(200, { guest: result.rows[0] });
} else {
  respond(404, { error: 'Guest not found' });
}
```

#### Step 3: Write local scripts to use these APIs (Local)

Now on your machine:

```python
# local_scripts/guest_importer.py
import requests
import csv

BASE_URL = "https://audos.com/api/hooks/execute/workspace-351699"
DB_API = f"{BASE_URL}/db-api"
LOOKUP_API = f"{BASE_URL}/lookup-guest"

def import_guests_from_csv(filepath):
    """Import guests from a CSV file into the platform."""
    with open(filepath, 'r') as f:
        reader = csv.DictReader(f)

        for row in reader:
            # Check if guest already exists
            response = requests.post(LOOKUP_API, json={
                "email": row['email']
            })

            if response.status_code == 404:
                # Guest doesn't exist, insert them
                requests.post(DB_API, json={
                    "action": "insert",
                    "table": "podcast_guests",
                    "data": {
                        "name": row['name'],
                        "email": row['email'],
                        "status": "pending",
                        "notes": row.get('notes', '')
                    }
                })
                print(f"Added {row['name']}")
            else:
                print(f"Skipped {row['name']} (already exists)")

if __name__ == "__main__":
    import_guests_from_csv("guests.csv")
```

```typescript
// local_scripts/guest_report.ts
interface Guest {
  id: number;
  name: string;
  email: string;
  status: string;
  episode_date: string | null;
}

const BASE_URL = "https://audos.com/api/hooks/execute/workspace-351699";

async function generateGuestReport() {
  const response = await fetch(`${BASE_URL}/db-api`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      action: "query",
      table: "podcast_guests",
      filters: [
        { column: "status", operator: "eq", value: "confirmed" }
      ],
      orderBy: { column: "episode_date", direction: "asc" }
    })
  });

  const { rows } = await response.json();
  const guests = rows as Guest[];

  console.log("\n=== Confirmed Guests ===\n");
  guests.forEach(guest => {
    console.log(`${guest.name} - ${guest.email}`);
    console.log(`  Episode Date: ${guest.episode_date || 'TBD'}\n`);
  });
}

generateGuestReport();
```

#### Step 4: Run local scripts whenever you need them

```bash
# Import a batch of guests
python local_scripts/guest_importer.py

# Generate a report
npx tsx local_scripts/guest_report.ts

# Set up a cron job for daily reports
0 9 * * * cd /path/to/project && npx tsx local_scripts/guest_report.ts
```

### Key Benefits of This Workflow

1. **Speed** - Write and test business logic locally with instant feedback
2. **Flexibility** - Use any language, any tools, any libraries
3. **Reliability** - Platform manages database, hosting, and uptime
4. **Collaboration** - Team members can run the same scripts against the same data
5. **Version control** - Git tracks your local scripts, platform auto-syncs via GitHub

---

## Setting Up Local Development

### Prerequisites

- Node.js 18+ (for TypeScript) or Python 3.8+ (for Python scripts)
- Git (for version control)
- Your favorite IDE (VS Code, PyCharm, etc.)

### Step 1: Get Your Workspace Credentials

Every workspace has a unique **workspace number** that identifies it. Find yours:

1. Look at your workspace URL: `https://audos.com/workspace/{number}`
2. Or ask Otto: "What's my workspace number?"

Example workspace info:
```
Workspace ID: 8f1ad824-832f-4af8-b77e-ab931a250625
Workspace Number: 351699
Base API URL: https://audos.com/api/hooks/execute/workspace-351699
```

### Step 2: Set Up Environment Variables

Create a `.env` file in your project:

```bash
# .env
WORKSPACE_NUMBER=351699
BASE_URL=https://audos.com/api/hooks/execute/workspace-351699
```

Add `.env` to your `.gitignore`:
```bash
echo ".env" >> .gitignore
```

### Step 3: Create an API Client Wrapper

#### Python Version

```python
# audos_client.py
import os
import requests
from typing import Optional, Dict, Any, List
from dotenv import load_dotenv

load_dotenv()

class AudosClient:
    """Client for interacting with Audos platform APIs."""

    def __init__(self):
        self.workspace_number = os.getenv("WORKSPACE_NUMBER")
        self.base_url = os.getenv("BASE_URL")

        if not self.workspace_number or not self.base_url:
            raise ValueError("Missing WORKSPACE_NUMBER or BASE_URL in environment")

    def db_query(
        self,
        table: str,
        columns: Optional[List[str]] = None,
        filters: Optional[List[Dict]] = None,
        order_by: Optional[Dict] = None,
        limit: int = 50,
        offset: int = 0
    ) -> List[Dict[str, Any]]:
        """Query data from a database table."""
        response = requests.post(
            f"{self.base_url}/db-api",
            json={
                "action": "query",
                "table": table,
                "columns": columns,
                "filters": filters,
                "orderBy": order_by,
                "limit": limit,
                "offset": offset
            }
        )
        response.raise_for_status()
        return response.json()["rows"]

    def db_insert(self, table: str, data: Dict[str, Any]) -> Dict[str, Any]:
        """Insert a row into a database table."""
        response = requests.post(
            f"{self.base_url}/db-api",
            json={
                "action": "insert",
                "table": table,
                "data": data
            }
        )
        response.raise_for_status()
        result = response.json()
        return result["rows"][0]

    def db_update(
        self,
        table: str,
        filters: List[Dict],
        data: Dict[str, Any]
    ) -> int:
        """Update rows in a database table."""
        response = requests.post(
            f"{self.base_url}/db-api",
            json={
                "action": "update",
                "table": table,
                "filters": filters,
                "data": data
            }
        )
        response.raise_for_status()
        return response.json()["updated"]

    def db_delete(self, table: str, filters: List[Dict]) -> int:
        """Delete rows from a database table."""
        response = requests.post(
            f"{self.base_url}/db-api",
            json={
                "action": "delete",
                "table": table,
                "filters": filters
            }
        )
        response.raise_for_status()
        return response.json()["deleted"]

    def generate_text(
        self,
        prompt: str,
        system_prompt: Optional[str] = None
    ) -> str:
        """Generate AI text content."""
        response = requests.post(
            f"{self.base_url}/ai-api",
            json={
                "action": "generate",
                "prompt": prompt,
                "systemPrompt": system_prompt
            }
        )
        response.raise_for_status()
        return response.json()["text"]

    def send_email(
        self,
        to: str,
        subject: str,
        text: str,
        html: Optional[str] = None
    ) -> Dict[str, Any]:
        """Send an email."""
        response = requests.post(
            f"{self.base_url}/email-api",
            json={
                "action": "send",
                "to": to,
                "subject": subject,
                "text": text,
                "html": html
            }
        )
        response.raise_for_status()
        return response.json()

    def list_contacts(
        self,
        limit: int = 50,
        has_email: bool = True
    ) -> List[Dict[str, Any]]:
        """List CRM contacts."""
        response = requests.post(
            f"{self.base_url}/crm-api",
            json={
                "action": "list",
                "limit": limit,
                "hasEmail": has_email
            }
        )
        response.raise_for_status()
        return response.json()["contacts"]

    def create_contact(
        self,
        email: str,
        name: Optional[str] = None,
        phone: Optional[str] = None
    ) -> Dict[str, Any]:
        """Create a new CRM contact."""
        response = requests.post(
            f"{self.base_url}/crm-api",
            json={
                "action": "create",
                "email": email,
                "name": name,
                "phone": phone
            }
        )
        response.raise_for_status()
        return response.json()["contact"]

# Usage example
if __name__ == "__main__":
    client = AudosClient()

    # Query database
    guests = client.db_query("podcast_guests", limit=10)
    print(f"Found {len(guests)} guests")

    # Generate AI content
    caption = client.generate_text(
        prompt="Write an Instagram caption for our latest episode",
        system_prompt="Keep it under 150 characters, engaging tone"
    )
    print(f"Generated caption: {caption}")
```

#### TypeScript Version

```typescript
// audos-client.ts
import 'dotenv/config';

interface QueryOptions {
  columns?: string[];
  filters?: Array<{
    column: string;
    operator: string;
    value: any;
  }>;
  orderBy?: {
    column: string;
    direction: 'asc' | 'desc';
  };
  limit?: number;
  offset?: number;
}

export class AudosClient {
  private workspaceNumber: string;
  private baseUrl: string;

  constructor() {
    this.workspaceNumber = process.env.WORKSPACE_NUMBER!;
    this.baseUrl = process.env.BASE_URL!;

    if (!this.workspaceNumber || !this.baseUrl) {
      throw new Error('Missing WORKSPACE_NUMBER or BASE_URL in environment');
    }
  }

  async dbQuery<T = any>(
    table: string,
    options: QueryOptions = {}
  ): Promise<T[]> {
    const response = await fetch(`${this.baseUrl}/db-api`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: 'query',
        table,
        ...options
      })
    });

    if (!response.ok) {
      throw new Error(`DB query failed: ${response.statusText}`);
    }

    const data = await response.json();
    return data.rows;
  }

  async dbInsert<T = any>(
    table: string,
    data: Record<string, any>
  ): Promise<T> {
    const response = await fetch(`${this.baseUrl}/db-api`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: 'insert',
        table,
        data
      })
    });

    if (!response.ok) {
      throw new Error(`DB insert failed: ${response.statusText}`);
    }

    const result = await response.json();
    return result.rows[0];
  }

  async dbUpdate(
    table: string,
    filters: Array<{ column: string; operator: string; value: any }>,
    data: Record<string, any>
  ): Promise<number> {
    const response = await fetch(`${this.baseUrl}/db-api`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: 'update',
        table,
        filters,
        data
      })
    });

    if (!response.ok) {
      throw new Error(`DB update failed: ${response.statusText}`);
    }

    const result = await response.json();
    return result.updated;
  }

  async dbDelete(
    table: string,
    filters: Array<{ column: string; operator: string; value: any }>
  ): Promise<number> {
    const response = await fetch(`${this.baseUrl}/db-api`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: 'delete',
        table,
        filters
      })
    });

    if (!response.ok) {
      throw new Error(`DB delete failed: ${response.statusText}`);
    }

    const result = await response.json();
    return result.deleted;
  }

  async generateText(
    prompt: string,
    systemPrompt?: string
  ): Promise<string> {
    const response = await fetch(`${this.baseUrl}/ai-api`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: 'generate',
        prompt,
        systemPrompt
      })
    });

    if (!response.ok) {
      throw new Error(`AI generation failed: ${response.statusText}`);
    }

    const result = await response.json();
    return result.text;
  }

  async sendEmail(options: {
    to: string;
    subject: string;
    text: string;
    html?: string;
  }): Promise<void> {
    const response = await fetch(`${this.baseUrl}/email-api`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: 'send',
        ...options
      })
    });

    if (!response.ok) {
      throw new Error(`Email send failed: ${response.statusText}`);
    }
  }

  async listContacts(options?: {
    limit?: number;
    hasEmail?: boolean;
  }): Promise<any[]> {
    const response = await fetch(`${this.baseUrl}/crm-api`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: 'list',
        limit: options?.limit ?? 50,
        hasEmail: options?.hasEmail ?? true
      })
    });

    if (!response.ok) {
      throw new Error(`CRM list failed: ${response.statusText}`);
    }

    const result = await response.json();
    return result.contacts;
  }

  async createContact(options: {
    email: string;
    name?: string;
    phone?: string;
  }): Promise<any> {
    const response = await fetch(`${this.baseUrl}/crm-api`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        action: 'create',
        ...options
      })
    });

    if (!response.ok) {
      throw new Error(`CRM create failed: ${response.statusText}`);
    }

    const result = await response.json();
    return result.contact;
  }
}

// Usage example
async function main() {
  const client = new AudosClient();

  // Query database
  const guests = await client.dbQuery('podcast_guests', { limit: 10 });
  console.log(`Found ${guests.length} guests`);

  // Generate AI content
  const caption = await client.generateText(
    'Write an Instagram caption for our latest episode',
    'Keep it under 150 characters, engaging tone'
  );
  console.log(`Generated caption: ${caption}`);
}

if (require.main === module) {
  main();
}
```

### Step 4: Install Dependencies

**For Python:**
```bash
pip install requests python-dotenv
```

**For TypeScript:**
```bash
npm install dotenv
npm install -D @types/node tsx
```

### Step 5: Test Your Setup

```python
# test_connection.py
from audos_client import AudosClient

client = AudosClient()

# List all database tables
response = requests.post(
    f"{client.base_url}/db-api",
    json={"action": "list-tables"}
)
tables = response.json()["tables"]
print(f"Found {len(tables)} tables:")
for table in tables:
    print(f"  - {table['name']} ({table['rowCount']} rows)")
```

```typescript
// test-connection.ts
import { AudosClient } from './audos-client';

async function test() {
  const client = new AudosClient();

  const response = await fetch(`${client['baseUrl']}/db-api`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ action: 'list-tables' })
  });

  const { tables } = await response.json();
  console.log(`Found ${tables.length} tables:`);
  tables.forEach((table: any) => {
    console.log(`  - ${table.name} (${table.rowCount} rows)`);
  });
}

test();
```

---

## Working with Platform APIs

### Available APIs

The Audos platform exposes several REST APIs for local development:

| API | Endpoint | Purpose |
|-----|----------|---------|
| **Database API** | `/db-api` | Full CRUD access to workspace tables |
| **AI API** | `/ai-api` | Text generation (GPT-4o-mini) |
| **Email API** | `/email-api` | Send transactional emails |
| **CRM API** | `/crm-api` | Manage contacts and leads |
| **Analytics API** | `/analytics-api` | Query visitor metrics |
| **Storage API** | `/storage-api` | Upload and retrieve files |
| **Web API** | `/web-api` | Web scraping and search |
| **Scheduler API** | `/scheduler-api` | Create cron jobs |

### Database API Patterns

#### Querying Data

```python
# Simple query
guests = client.db_query("podcast_guests")

# With filters
active_guests = client.db_query(
    "podcast_guests",
    filters=[
        {"column": "status", "operator": "eq", "value": "active"}
    ]
)

# With sorting and pagination
recent_guests = client.db_query(
    "podcast_guests",
    order_by={"column": "created_at", "direction": "desc"},
    limit=20,
    offset=0
)

# Complex filters
confirmed_this_month = client.db_query(
    "podcast_guests",
    filters=[
        {"column": "status", "operator": "eq", "value": "confirmed"},
        {"column": "episode_date", "operator": "gte", "value": "2026-03-01"}
    ]
)
```

#### Inserting Data

```python
# Insert a single row
new_guest = client.db_insert("podcast_guests", {
    "name": "Jane Doe",
    "email": "jane@example.com",
    "status": "pending"
})
print(f"Created guest with ID: {new_guest['id']}")

# Insert multiple rows (batch)
guests_data = [
    {"name": "Guest 1", "email": "guest1@example.com", "status": "pending"},
    {"name": "Guest 2", "email": "guest2@example.com", "status": "pending"}
]
for data in guests_data:
    client.db_insert("podcast_guests", data)
```

#### Updating Data

```python
# Update by ID
client.db_update(
    "podcast_guests",
    filters=[{"column": "id", "operator": "eq", "value": 123}],
    data={"status": "confirmed", "episode_date": "2026-04-15"}
)

# Bulk update
client.db_update(
    "podcast_guests",
    filters=[{"column": "status", "operator": "eq", "value": "pending"}],
    data={"status": "contacted"}
)
```

#### Deleting Data

```python
# Delete by ID
client.db_delete(
    "podcast_guests",
    filters=[{"column": "id", "operator": "eq", "value": 999}]
)

# Bulk delete
client.db_delete(
    "podcast_guests",
    filters=[{"column": "status", "operator": "eq", "value": "cancelled"}]
)
```

### AI API Patterns

```python
# Generate social media content
caption = client.generate_text(
    prompt="Write an Instagram caption for our episode with Dr. Jane Smith about climate change",
    system_prompt="Keep it engaging, under 150 characters, with 3 hashtags"
)

# Generate email drafts
email_body = client.generate_text(
    prompt=f"Draft a thank you email for {guest_name} who appeared on our podcast",
    system_prompt="Professional but warm tone, include a call to share the episode"
)

# Summarize transcripts
summary = client.generate_text(
    prompt=f"Summarize this podcast transcript in 3 key takeaways:\n\n{transcript}",
    system_prompt="Concise, bullet points, focus on actionable insights"
)
```

### Email API Patterns

```python
# Send guest confirmation
client.send_email(
    to=guest_email,
    subject="Your Podcast Episode is Scheduled!",
    text=f"Hi {guest_name},\n\nGreat news! Your episode is scheduled for {episode_date}.",
    html=f"<h1>Hi {guest_name}!</h1><p>Great news! Your episode is scheduled for {episode_date}.</p>"
)

# Send batch notifications
for guest in confirmed_guests:
    client.send_email(
        to=guest["email"],
        subject="Your Episode Goes Live Tomorrow",
        text=f"Hi {guest['name']},\n\nYour episode drops tomorrow!"
    )
```

### CRM API Patterns

```python
# List and filter contacts
newsletter_subscribers = requests.post(
    f"{client.base_url}/crm-api",
    json={
        "action": "list",
        "tags": ["newsletter"],
        "limit": 1000
    }
).json()["contacts"]

# Create new contact
new_contact = client.create_contact(
    email="prospect@example.com",
    name="Sarah Johnson",
    phone="+1-555-987-6543"
)

# Add tags to contacts
requests.post(
    f"{client.base_url}/crm-api",
    json={
        "action": "add-tags",
        "tags": ["vip", "repeat-guest"],
        "filter": {
            "contactIds": [contact_id]
        }
    }
)
```

---

## When to Go Back to Otto

While local development is powerful, certain operations **must be done via the platform**.

### You Need Otto When:

#### 1. Creating New Server Functions

**Local development can't do:** Create new HTTP endpoints

**What to do:** Ask Otto to create a server function (hook)

```
User: "Create a server function called 'guest-stats' that returns stats
about podcast guests - total count, confirmed count, and next 5 scheduled episodes."

Otto: [Creates the hook with db.query calls and respond()]
```

Then in your local code:
```python
response = requests.post(f"{base_url}/guest-stats")
stats = response.json()
print(f"Total guests: {stats['total']}")
```

#### 2. Database Schema Changes

**Local development can't do:** Create or modify database tables

**What to do:** Ask Otto to create or alter tables

```
User: "Add a 'social_media_handles' column to podcast_guests table, type JSON."

Otto: [Uses db_alter_table to add the column]
```

Then in your local code:
```python
# Now you can use the new column
client.db_update(
    "podcast_guests",
    filters=[{"column": "id", "operator": "eq", "value": 123}],
    data={"social_media_handles": {"twitter": "@janedoe", "linkedin": "janedoe"}}
)
```

#### 3. React App Changes

**Local development can't do:** Edit Space apps, landing pages, or the EmailGate

**What to do:** Ask Otto to make the changes

```
User: "Update the Guest Prep app to show episode dates in a calendar view instead of a list."

Otto: [Uses delegate_app_edit to modify the React component]
```

#### 4. Infrastructure Configuration

**Local development can't do:** Set up domains, email sending, integrations

**What to do:** Ask Otto

```
User: "Connect a custom domain: podcast.example.com"
User: "Set up email sending from no-reply@podcast.example.com"
User: "Connect my Stripe account for accepting payments"
```

### The Pattern: Otto for Structure, Local for Logic

Think of Otto as your **infrastructure engineer** and your local environment as your **application layer**.

```
Otto creates:                     You build locally:
├── Database tables               ├── Data import scripts
├── Server functions              ├── Analytics reports
├── React apps                    ├── Automation workflows
├── API endpoints                 ├── CLI tools
└── Infrastructure config         └── Integration tests
```

---

## Testing Strategies

### Testing Locally with Mock Responses

For fast iteration without hitting live APIs:

```python
# mock_audos_client.py
class MockAudosClient:
    """Mock client for testing without hitting live APIs."""

    def __init__(self):
        self.data = {
            "podcast_guests": [
                {"id": 1, "name": "Jane Doe", "email": "jane@example.com", "status": "confirmed"},
                {"id": 2, "name": "John Smith", "email": "john@example.com", "status": "pending"}
            ]
        }

    def db_query(self, table, **kwargs):
        return self.data.get(table, [])

    def db_insert(self, table, data):
        new_id = len(self.data.get(table, [])) + 1
        row = {"id": new_id, **data}
        if table not in self.data:
            self.data[table] = []
        self.data[table].append(row)
        return row

    def generate_text(self, prompt, system_prompt=None):
        return f"[Mock AI response for: {prompt[:50]}...]"

# test_my_script.py
from mock_audos_client import MockAudosClient

def test_guest_import():
    client = MockAudosClient()

    # Test your logic without hitting live API
    result = client.db_insert("podcast_guests", {
        "name": "Test Guest",
        "email": "test@example.com",
        "status": "pending"
    })

    assert result["id"] == 3
    assert result["name"] == "Test Guest"
```

### Testing Against Live APIs

Use environment flags to switch between mock and live:

```python
# main.py
import os
from audos_client import AudosClient
from mock_audos_client import MockAudosClient

def get_client():
    if os.getenv("USE_MOCK") == "true":
        return MockAudosClient()
    return AudosClient()

client = get_client()

# Your code works with either client
guests = client.db_query("podcast_guests")
```

Run with mock:
```bash
USE_MOCK=true python main.py
```

Run against live:
```bash
python main.py
```

### Testing Server Functions

**Option 1: Via Otto's test_server_function tool**

Ask Otto:
```
User: "Test the 'guest-stats' server function"

Otto: [Uses test_server_function tool and shows you the response]
```

**Option 2: Direct HTTP call from local**

```python
response = requests.post(
    f"{base_url}/guest-stats",
    json={"test": True}
)
print(response.json())
```

### Integration Testing Strategy

1. **Unit test business logic with mocks** (fast, no network)
2. **Integration test with live APIs** (slower, validates end-to-end)
3. **Use separate test data** (create a `test_*` table or use status flags)

```python
# integration_test.py
def test_full_workflow():
    client = AudosClient()

    # Create a test guest
    guest = client.db_insert("podcast_guests", {
        "name": "Test Guest",
        "email": "test@example.com",
        "status": "test"
    })

    # Generate AI content for them
    caption = client.generate_text(
        f"Write a caption for our episode with {guest['name']}"
    )
    assert len(caption) > 0

    # Clean up test data
    client.db_delete("podcast_guests", [
        {"column": "status", "operator": "eq", "value": "test"}
    ])
```

---

## Best Practices

### 1. Keep Business Logic in Server Functions (Not Just Local)

**Bad:**
```python
# All logic in local script - not reusable
def calculate_guest_score(guest):
    # Complex scoring logic here...
    return score

for guest in guests:
    score = calculate_guest_score(guest)
    print(f"{guest['name']}: {score}")
```

**Good:**
```javascript
// Server function: calculate-guest-score
const { guestId } = request.body;

const guest = await db.query('podcast_guests', {
  where: { id: guestId },
  limit: 1
});

if (guest.rowCount === 0) {
  respond(404, { error: 'Guest not found' });
  return;
}

// Complex scoring logic here - now reusable from anywhere
const score = calculateScore(guest.rows[0]);
respond(200, { score });
```

```python
# Local script - simple and focused
for guest in guests:
    response = requests.post(f"{base_url}/calculate-guest-score",
                           json={"guestId": guest["id"]})
    score = response.json()["score"]
    print(f"{guest['name']}: {score}")
```

**Why?** Server functions are accessible from:
- Local scripts
- React apps
- Other server functions
- External webhooks
- Phone agents

### 2. Use Database Tables for Persistent Data (Not JSON Files)

**Bad:**
```python
# Storing data in local JSON files
with open("guests.json", "r") as f:
    guests = json.load(f)

guests.append(new_guest)

with open("guests.json", "w") as f:
    json.dump(guests, f)
```

**Good:**
```python
# Using platform database
client.db_insert("podcast_guests", new_guest)
```

**Why?**
- Database is accessible from apps, server functions, and analytics
- Automatic backups and recovery
- Query performance and indexing
- Multi-user safe (no file locking issues)

### 3. Design APIs to be Reusable

**Bad:**
```javascript
// Server function: specific-report-for-john
const guests = await db.query('podcast_guests', {
  where: { status: 'confirmed' },
  orderBy: { column: 'episode_date', direction: 'asc' }
});
respond(200, { guests });
```

**Good:**
```javascript
// Server function: guest-report
const { status, orderBy, limit } = request.body;

const guests = await db.query('podcast_guests', {
  where: status ? { status } : undefined,
  orderBy: orderBy || { column: 'created_at', direction: 'desc' },
  limit: limit || 50
});

respond(200, { guests, count: guests.length });
```

**Why?** Flexible APIs can serve multiple use cases:
- Different users with different needs
- Different React apps
- Different automation scripts

### 4. Document Your Server Functions

Add descriptions when Otto creates them:

```
User: "Create a server function called 'guest-report' that returns podcast
guests filtered by status. Document what parameters it accepts."

Otto: [Creates the hook with this description]
Description: Returns podcast guests filtered by status. Accepts:
  - status (string): Filter by guest status (confirmed, pending, etc.)
  - orderBy (object): Sort order {column, direction}
  - limit (number): Max guests to return
```

### 5. Version Control Your Local Scripts

```bash
# Recommended project structure
my-podcast-automation/
├── .env                    # Environment config (gitignored)
├── .gitignore
├── audos_client.py         # API client wrapper
├── scripts/
│   ├── guest_importer.py   # Import guests from CSV
│   ├── generate_reports.py # Generate analytics reports
│   └── send_reminders.py   # Send episode reminders
├── tests/
│   ├── test_client.py      # Unit tests
│   └── integration_test.py # Integration tests
└── README.md               # Project documentation
```

Commit your local scripts to Git, but **not your .env file**.

### 6. Use Consistent Naming Conventions

**Tables:** `snake_case` (matches SQL conventions)
```
podcast_guests
episode_transcripts
voice_profiles
```

**Server Functions:** `kebab-case` (matches URL conventions)
```
guest-report
send-reminder
calculate-score
```

**Local Scripts:** `snake_case.py` or `kebab-case.ts`
```
guest_importer.py
generate-reports.ts
```

---

## Common Patterns

### Pattern 1: Scheduled Report Generation

**Scenario:** Generate a weekly guest report and email it

**Implementation:**

1. Create a server function (via Otto):
```javascript
// Server function: weekly-guest-report
const { includeAll } = request.body;

// Query guests
const confirmedGuests = await db.query('podcast_guests', {
  where: { status: 'confirmed' },
  orderBy: { column: 'episode_date', direction: 'asc' }
});

const pendingGuests = await db.query('podcast_guests', {
  where: { status: 'pending' }
});

// Generate report with AI
const reportText = await platform.generateText({
  userPrompt: `Generate a weekly podcast guest report. Confirmed: ${confirmedGuests.rowCount}, Pending: ${pendingGuests.rowCount}`,
  model: 'gpt-4o-mini'
});

// Email it
await platform.sendEmail({
  to: 'host@podcast.com',
  subject: 'Weekly Guest Report',
  text: reportText
});

respond(200, { sent: true, reportText });
```

2. Schedule it (via Otto):
```
User: "Schedule the weekly-guest-report to run every Monday at 9am"

Otto: [Creates a scheduled task that calls the hook every Monday]
```

3. Or call it manually from local:
```python
# scripts/run_report.py
client = AudosClient()
response = requests.post(f"{client.base_url}/weekly-guest-report")
print(response.json()["reportText"])
```

### Pattern 2: CSV Import Pipeline

**Scenario:** Import guest data from CSV, validate, dedupe, and insert

**Implementation:**

```python
# scripts/import_guests.py
import csv
from audos_client import AudosClient

def import_guests(filepath):
    client = AudosClient()

    # Read CSV
    with open(filepath, 'r') as f:
        reader = csv.DictReader(f)
        rows = list(reader)

    imported = 0
    skipped = 0
    errors = []

    for row in rows:
        # Validate
        if not row.get('email'):
            errors.append(f"Row missing email: {row}")
            continue

        # Check for duplicates
        existing = client.db_query("podcast_guests", filters=[
            {"column": "email", "operator": "eq", "value": row['email']}
        ])

        if existing:
            skipped += 1
            continue

        # Insert
        try:
            client.db_insert("podcast_guests", {
                "name": row.get('name', ''),
                "email": row['email'],
                "status": row.get('status', 'pending'),
                "notes": row.get('notes', '')
            })
            imported += 1
        except Exception as e:
            errors.append(f"Failed to insert {row['email']}: {str(e)}")

    print(f"\nImport Complete:")
    print(f"  Imported: {imported}")
    print(f"  Skipped (duplicates): {skipped}")
    print(f"  Errors: {len(errors)}")

    if errors:
        print("\nErrors:")
        for error in errors:
            print(f"  - {error}")

if __name__ == "__main__":
    import sys
    if len(sys.argv) < 2:
        print("Usage: python import_guests.py <csv_file>")
        sys.exit(1)

    import_guests(sys.argv[1])
```

Run it:
```bash
python scripts/import_guests.py guests.csv
```

### Pattern 3: AI Content Generation Pipeline

**Scenario:** Generate social media captions for all recent episodes

**Implementation:**

```typescript
// scripts/generate-captions.ts
import { AudosClient } from '../audos-client';

interface Episode {
  id: number;
  title: string;
  description: string;
  transcript: string;
}

async function generateCaptions() {
  const client = new AudosClient();

  // Get recent episodes without captions
  const episodes = await client.dbQuery<Episode>('studio_episodes', {
    filters: [
      { column: 'status', operator: 'eq', value: 'ready' }
    ],
    limit: 10
  });

  console.log(`Generating captions for ${episodes.length} episodes...\n`);

  for (const episode of episodes) {
    // Generate Instagram caption
    const instagramCaption = await client.generateText(
      `Write an Instagram caption for this podcast episode:\nTitle: ${episode.title}\nDescription: ${episode.description}`,
      'Engaging tone, under 150 characters, include 3 relevant hashtags'
    );

    // Generate LinkedIn post
    const linkedinPost = await client.generateText(
      `Write a LinkedIn post for this podcast episode:\nTitle: ${episode.title}\nDescription: ${episode.description}`,
      'Professional tone, 2-3 paragraphs, include a call to action'
    );

    // Save generated content to database
    await client.dbInsert('studio_generated_content', {
      episode_id: episode.id,
      platform: 'instagram',
      content_type: 'caption',
      content: instagramCaption,
      status: 'draft'
    });

    await client.dbInsert('studio_generated_content', {
      episode_id: episode.id,
      platform: 'linkedin',
      content_type: 'post',
      content: linkedinPost,
      status: 'draft'
    });

    console.log(`✓ Generated captions for: ${episode.title}`);
  }

  console.log('\nDone! Captions saved to studio_generated_content table.');
}

generateCaptions();
```

Run it:
```bash
npx tsx scripts/generate-captions.ts
```

### Pattern 4: Webhook Handler + Local Processing

**Scenario:** Receive Stripe webhooks, validate, and update database

**Implementation:**

1. Create webhook receiver (via Otto):
```javascript
// Server function: stripe-webhook
const event = request.body;
const signature = request.headers['stripe-signature'];

// Validate webhook (Stripe signature validation)
// ... validation logic ...

if (event.type === 'checkout.session.completed') {
  const session = event.data.object;

  // Update guest record with payment info
  await db.update('podcast_guests',
    [{ column: 'email', operator: 'eq', value: session.customer_email }],
    {
      status: 'paid',
      payment_status: 'completed',
      stripe_session_id: session.id
    }
  );

  // Post to agent chat
  await platform.postAgentMessage({
    message: `New payment received from ${session.customer_email} - $${(session.amount_total / 100).toFixed(2)}`
  });
}

respond(200, { received: true });
```

2. Optional: Local monitoring script:
```python
# scripts/monitor_payments.py
import time
from audos_client import AudosClient

def monitor_payments():
    client = AudosClient()

    while True:
        # Check for recent payments
        recent = client.db_query("podcast_guests", filters=[
            {"column": "payment_status", "operator": "eq", "value": "completed"},
            {"column": "updated_at", "operator": "gte", "value": "last_hour"}
        ])

        if recent:
            print(f"[{time.strftime('%H:%M:%S')}] {len(recent)} new payments")
            for guest in recent:
                print(f"  - {guest['name']} ({guest['email']})")

        time.sleep(60)  # Check every minute

if __name__ == "__main__":
    monitor_payments()
```

### Pattern 5: Multi-Step Data Pipeline

**Scenario:** Scrape podcast directories, analyze, score, and import leads

**Implementation:**

```python
# scripts/lead_pipeline.py
import requests
from audos_client import AudosClient

def scrape_podcast_directory():
    """Scrape external podcast directory."""
    # ... web scraping logic ...
    return [
        {"name": "Podcast 1", "email": "host1@example.com", "category": "tech"},
        {"name": "Podcast 2", "email": "host2@example.com", "category": "business"}
    ]

def score_lead(lead, client):
    """Use AI to score lead relevance."""
    prompt = f"""
    Score this podcast lead for potential collaboration (0-100):
    Name: {lead['name']}
    Category: {lead['category']}

    Our podcast focuses on entrepreneurship and technology.
    """

    score_text = client.generate_text(prompt, "Return just a number between 0-100")
    return int(score_text.strip())

def run_pipeline():
    client = AudosClient()

    # Step 1: Scrape
    print("Step 1: Scraping podcast directories...")
    raw_leads = scrape_podcast_directory()
    print(f"  Found {len(raw_leads)} leads")

    # Step 2: Score
    print("\nStep 2: Scoring leads with AI...")
    scored_leads = []
    for lead in raw_leads:
        score = score_lead(lead, client)
        scored_leads.append({**lead, "score": score})
        print(f"  {lead['name']}: {score}")

    # Step 3: Filter and import
    print("\nStep 3: Importing high-quality leads...")
    imported = 0
    for lead in scored_leads:
        if lead['score'] >= 70:
            client.db_insert("outreach_leads", {
                "name": lead['name'],
                "email": lead['email'],
                "podcast_name": lead['name'],
                "source": "directory_scrape",
                "status": "new",
                "notes": f"AI Score: {lead['score']}"
            })
            imported += 1
            print(f"  ✓ Imported {lead['name']}")

    print(f"\nPipeline complete! Imported {imported}/{len(raw_leads)} leads.")

if __name__ == "__main__":
    run_pipeline()
```

---

## Conclusion

The Audos development workflow gives you the best of both worlds:

- **Platform (Otto)** manages your infrastructure - database, hosting, server functions, React apps
- **Local development** gives you speed, flexibility, and full control over business logic
- **HTTP APIs** bridge the gap - your local code calls the platform, platform code can be triggered by schedules or webhooks

### Quick Decision Tree

**Need to...**
- Create/modify a database table? → Ask Otto
- Create/modify a server function? → Ask Otto
- Edit a React app or landing page? → Ask Otto
- Set up infrastructure (domain, email, etc.)? → Ask Otto

**Need to...**
- Query data and generate reports? → Local script
- Import/export data from CSV/JSON? → Local script
- Test API behavior? → Local script
- Build CLI tools? → Local script
- Schedule jobs that call APIs? → Local script + cron

### Next Steps

1. **Set up your local environment** - Follow the setup guide above
2. **Create your API client wrapper** - Use Python or TypeScript examples
3. **Build your first script** - Query your database and print results
4. **Ask Otto to create a server function** - Make a custom API endpoint
5. **Call it from your local script** - See the hybrid workflow in action

### Additional Resources

- [Database API Docs](/tmp/workspace-workspace-351699/docs/database-api.md)
- [Server Functions Docs](/tmp/workspace-workspace-351699/integrations/server-functions/docs.md)
- [CRM API Docs](/tmp/workspace-workspace-351699/docs/crm-api.md)
- [AI API Docs](/tmp/workspace-workspace-351699/docs/ai-api.md)
- [Space App Development Guide](/tmp/workspace-workspace-351699/SPACE_APP_GUIDE.md)

---

**Questions?** Ask Otto: "How do I [specific task]?" - Otto will guide you to the right approach (local vs platform).
