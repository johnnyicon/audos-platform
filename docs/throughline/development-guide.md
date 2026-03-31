# Throughline Development Guide

**Last Updated:** 2026-03-31 (confirmed via live testing)
**Workspace ID:** `8f1ad824-832f-4af8-b77e-ab931a250625`
**Config ID:** `351699`

---

## Confirmed Development Pipeline

```
GitHub Repo                        Audos Platform
──────────────────────────         ──────────────────────
audos-workspace/               →   / (workspace root)
  apps/home/App.tsx             →   apps/home/App.tsx ✅
  apps/briefing/App.tsx         →   apps/briefing/App.tsx
  components/                   →   components/
  hooks/                        →   hooks/
  lib/                          →   lib/
  Desktop.tsx                   →   Desktop.tsx

src/  (local Vite dev only — NOT synced)
```

**The workflow:**
1. Edit files inside `audos-workspace/` locally
2. Commit and push to `main` on GitHub
3. Platform auto-syncs via webhook — no manual trigger needed

---

## What Syncs vs What Doesn't

| Path | Syncs? | Notes |
|------|--------|-------|
| `audos-workspace/apps/` | ✅ Yes | Mini-app components |
| `audos-workspace/components/` | ✅ Yes | Shared UI components |
| `audos-workspace/hooks/` | ✅ Yes | Custom React hooks |
| `audos-workspace/lib/` | ✅ Yes | Utilities |
| `audos-workspace/tools/` | ✅ Yes | Internal dashboards |
| `audos-workspace/Desktop.tsx` | ✅ Yes | Main space layout |
| `audos-workspace/SpaceRuntimeContext.tsx` | ✅ Yes | Context provider |
| `audos-workspace/config.json` | ✅ Yes | Space configuration |
| `audos-workspace/landing-pages/` | ❌ No | Otto-managed only |
| `src/` | ❌ No | Local Vite dev, not synced |

---

## Code Constraints

### Imports

```tsx
// ✅ Relative imports
import { useSpaceData } from '../../hooks/useSpaceData';
import { colors } from '../lib/colors';

// ✅ Platform-available packages (React, lucide-react, etc.)
import { useState } from 'react';
import { Activity } from 'lucide-react';

// ❌ Node.js / npm packages
import express from 'express';

// ❌ Alias imports
import { Button } from '@/components/ui/button';
```

### Data Persistence

```tsx
// ✅ Platform hook for JSON data files
import { useSpaceData } from '../../hooks/useSpaceData';
const { data, update } = useSpaceData<Item[]>({ dataFile: 'data/items.json', autoFetch: true });

// ✅ WorkspaceDB for database tables
const db = window.useWorkspaceDB();
const results = await db.query('my_table', { filters: [...] });

// ❌ localStorage — breaks mode isolation
localStorage.setItem('items', JSON.stringify(items));
```

### Folder Names

Folder names must be **lowercase**:
- `apps/home/App.tsx` ✅
- `apps/Home/App.tsx` ❌

---

## App Component Pattern

```tsx
interface MyAppProps {
  dataFile: string; // passed from config.json
}

export default function MyApp({ dataFile }: MyAppProps) {
  const { data, update, loading } = useSpaceData<MyData[]>({
    dataFile,
    autoFetch: true
  });

  if (loading) return <div>Loading...</div>;

  return <div className="p-4">{/* UI */}</div>;
}
```

---

## Backend API Calls

All server function endpoints are available from both local dev and from within the platform.

**Base URL:** `https://audos.com/api/hooks/execute/workspace-351699`

See [`throughline-api-reference.md`](./throughline-api-reference.md) for the full endpoint list.

---

## Creating New Apps vs. Editing Existing Apps

This is the most important workflow distinction on the platform.

### Editing existing apps — use GitHub

Any app that already exists in `.published-source/` can be edited freely via GitHub:

1. Edit `audos-workspace/apps/{app-id}/App.tsx` locally
2. Commit and push to `main`
3. Platform auto-syncs — changes are live

### Creating new apps — requires Otto first

New apps must go through a two-step process. GitHub sync alone is **not sufficient** to make a new app appear in the UI because:

- The platform serves apps from `.published-source/`, which is write-protected from GitHub
- `recompile()` does not process new apps from the `apps/` folder
- When GitHub sync is enabled, Otto's code editor is locked ("code editing frozen")

**Correct workflow for new apps:**

```
1. Open Developer tab → disable GitHub sync
2. Tell Otto: "Create a new app called [name] at apps/[id]/App.tsx"
3. Tell Otto: "Add [app] to the dock in Desktop.tsx"
4. Publish changes via Otto
5. Re-enable GitHub sync in Developer tab
6. All future edits to that app can be done via GitHub
```

> **Note:** `Desktop.tsx` has a hardcoded app list. Adding an app to `config.json` alone does not make it appear in the dock — it must also be added to `Desktop.tsx`. Otto handles this when creating the app; after that, `Desktop.tsx` edits can be done via GitHub.

### Summary table

| Task | GitHub | Otto |
|------|--------|------|
| Edit existing app code | ✅ Preferred | ✅ Optional |
| Create new app | ❌ Cannot | ✅ Required |
| Edit `Desktop.tsx` | ✅ Works | ✅ Works |
| Edit `config.json` | ✅ Works | ✅ Works |
| Delete an app | ❌ Cannot | ✅ Required |

---

## What Otto Manages (Do Not Edit Locally)

- **Landing pages** — use `delegate_landing_page_edit` via Otto
- **workspace-branding.json** — use Otto to update
- **Domain config** — platform-managed
- **Database table schemas** — create via Otto using `db_create_table`
- **`.published-source/`** — platform-managed compilation output; do not commit this directory
