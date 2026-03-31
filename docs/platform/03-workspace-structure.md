# Audos Workspace Structure

*Last Updated: March 31, 2026*

---

## Folder Overview

```
workspace-{uuid}/
├── apps/                # React mini-applications
├── assets/             # Static files (images, fonts)
├── community/          # Community features config
├── components/         # Shared React components
├── data/               # JSON data files (app state)
├── hooks/              # Server functions (APIs)
├── landing-pages/      # Landing page React code
├── lib/                # Shared utilities, types, helpers
├── tools/              # Internal dashboards/admin tools
├── config.json         # Space configuration
├── workspace-branding.json  # Brand identity
└── integrations/       # Third-party integration docs
```

---

## Detailed Breakdown

### apps/ - Mini-Applications

**Purpose**: Contains React components for each app in your space.

**Structure**:
```
apps/
├── Home.tsx          # Main dashboard
|   Briefing.tsx      # Guest research/briefings
|   Signature.tsx     # Voice profiles
└── Studio.tsx        # Reels/content management
```

**How to edit**: Use `delegate_app_edit` tool via Otto. Do NOT edit directly.

**App capabilities**:
- `useSpaceFiles()` - Persist JSON data
- `useWorkspaceDB()` - Access database tables
- `useSession()` - Get current user session
- Stripe integration for payments
- Fetch to server functions

---

### data/ - App State

**Purpose**: JSON files that store app-specific data.

**Structure**:
```
data/
├── signature.json    # Voice profiles
├── briefing.json     # Research sessions
└── studio.json       # Episodes/reels
```

**How to access**: Use `useSpaceFiles()` hook in apps, or read directly via Otto.

**Note**: For complex data, prefer database tables over JSON files.

---

### hooks/ - Server Functions

**Purpose**: Custom backend logic accessible via HTTP.

**This is where your APIs live!**

**Structure**:
```
hooks/
├── db-api.js         # Database CRUD
├── ai-api.js         # AI generation
├── email-api.js      # Email sending
├── web-api.js        # Web fetching
├── crm-api.js        # Contact management
├── analytics-api.js  # Visitor metrics
├── storage-api.js    # File upload
└── scheduler-api.js  # Scheduled tasks
```

**How to create/edit**: Use `manage_server_functions` tool via Otto.

**How to test**: Use `test_server_function` tool.

**How to call**:
```
POST https://www.audos.com/api/hooks/execute/workspace-351699/{hook-name}
```

---

### components/ - Shared UI

**Purpose**: Reusable React components shared across apps.

**Examples**:
- Navigation bars
- Form elements
- Card layouts
- Modals

**How to edit**: Use `delegate_app_edit` with specific component changes.

---

### landing-pages/ - Websites

**Purpose**: Public-facing landing pages for your business.

**Structure**:
```
landing-pages/
└── LandingPage.tsx  # Main landing page component
```

**How to edit**: Use `delegate_landing_page_edit` tool.

**Features**:
- Email capture forms
- Stripe checkout integration
- Analytics tracking
- Responsive design

---

### lib/ - Utilities

**Purpose**: Shared types, helpers, and utilities.

**Contents**:
- TypeScript type definitions
- Utility functions
- Constants
- API client helpers

---

### tools/ - Internal Dashboards

**Purpose**: Admin panels and internal tools for managing your business.

**How to create**: Use `create_dashboard` tool.

**Examples**:
- Revenue dashboard
- Contact management
- Content calendar
- Analytics views

---

### integrations/ - API Docs

**Purpose**: Documentation for integration endpoints.

**Structure**:
```
integrations/
├── email.md          # Email API docs
├── scheduler.md      # Scheduler API docs
├── stripe.md         # Stripe integration
├── crm.md            # CRM API docs
└── storage.md        # Storage API docs
```

**Useful for**: Understanding how to call internal platform APIs.

---

## Key Configuration Files

### config.json

Defines space structure, apps, and navigation.

```json
{
  "spaceTitle": "Throughline",
  "loadingScreen": {
    "enabled": true,
    "showProgressBar": true
  },
  "apps": [
    {
      "id": "home",
      "name": "Home",
      "icon": "Home",
      "component": "./apps/Home.tsx"
    }
  ]
}
```

### workspace-branding.json

Defines brand identity.

```json
{
  "name": "Throughline",
  "tagline": "The operating system for podcast creators",
  "logoUrl": "https://...",
  "colors": {
    "primary": "#4444FF",
    "secondary": "#101010"
  }
}
```

---

## What Can Be Edited Locally?

| Folder/File | Edit Locally? | How to Change |
|-------------|---------------|---------------|
| apps/ | ❌ No | `delegate_app_edit` |
| components/ | ❌ No | `delegate_app_edit` |
| landing-pages/ | ❌ No | `delegate_landing_page_edit` |
| hooks/ | ⭕ Via Otto | `manage_server_functions` |
| data/ | ⭕ Via API | `useSpaceFiles()` hook |
| tools/ | ⭕ Via Otto | `create_dashboard` |
| config.json | ⭕ Via Otto | Ask Otto to update |
| workspace-branding.json | ⭕ Via Otto | Ask Otto to update |

**Key Insight**: The primary way you extend functionality is through **hooks** (server functions) and **database tables**.

---

## Database Tables (Not in File System)

Your database tables are NOT stored as files — they're in PostgreSQL.

**Throughline tables**:
- `voice_profiles` - Voice signatures
- `speakers` - Podcast speakers
- `studio_episodes` - Episode data
- `studio_generated_content` - AI-generated content
- `reels` - Social media reels
- `captions` - Platform-specific captions
- `guest_prep_research_sessions` - Guest research
- `dashboard_activity` - Activity log
- And more...

**How to manage**:
- Create tables: `db_create_table` tool
- Query data: `db-api` with `query` action
- Insert data: `db-api` with `insert` action
- External access: Via your custom `db-api` server function