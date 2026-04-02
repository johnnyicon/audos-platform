---
last-updated: 2026-04-02
sources:
  - SDK-09-ARCHITECTURE-ROUND-2.md
  - AUDOS.md (auto-generated, throughline repo root)
  - UX audit (inbox/2026-04-02T1500-ux-audit-report.md)
  - Single-app architecture test (live, 2026-04-02)
---

# Throughline — Architecture Decisions

Confirmed findings from platform investigation. Answers to all architecture questions.

---

## Confirmed: Single App Architecture

**Decision: Consolidate to one entry in config.json.**

Apps are just React components loaded dynamically by Desktop.tsx. No permissions, isolation, or performance difference between 1 app and 5 apps. "Apps" is a scaffolding convention, not a technical requirement.

```json
// config.json — target state
{
  "apps": [
    {
      "id": "throughline",
      "name": "Throughline",
      "icon": "Sparkles",
      "component": "apps/throughline/App.tsx"
    }
  ]
}
```

Internal routing via `useState` (recommended) or `HashRouter`. Do not use BrowserRouter — may conflict with platform URL handling.

---

## Confirmed: audos-workspace/ is the Only Source of Truth

**`src/pages/apps/` should be deleted.** Audos ignores everything outside `audos-workspace/`. The `src/` folder is dead weight.

Correct structure:
```
audos-workspace/
  apps/throughline/App.tsx   ← Single entry point
  components/ui/             ← ShadCN components
  components/                ← Shared app components
  hooks/                     ← Custom hooks
  lib/                       ← Utilities, API wrappers, theme
  data/                      ← JSON data files
  config.json                ← App registry + branding
```

Local dev uses a mock layer in `src/lib/audos-sdk.ts` with env flag switching (`VITE_USE_REMOTE_API`).

---

## Confirmed: No Platform UI Constraints

Apps are pure React components. The platform injects:
- `SpaceRuntimeContext` — session, file ops, event tracking
- `window.__workspaceDb` / `window.useWorkspaceDB` — database SDK
- Global Tailwind CSS (v3)

Minimum requirement: export a default React component. That's it.

ShadCN: copy components to `audos-workspace/components/ui/`. Fully supported.

Radix UI (ShadCN dependency): must be in `cdnDependencies` in config.json, OR use ShadCN components that don't require Radix. To be validated by test.

---

## Confirmed: Compilation Model

- **Bundler**: ESBuild (ES2020, ESM)
- **Entry points**: One TSX file per app — can import from `components/`, `hooks/`, `lib/`
- **CDN deps**: React, ReactDOM, Lucide React pre-configured. Others added via `cdnDependencies` in config.json.
- **Deploy**: Push to GitHub → Audos inbound sync → ESBuild compile → live

No npm packages available inside apps. All dependencies must be in CDN importmap.

---

## Root Cause: Blank White Screens

Studio, Briefing, Signature all show blank white screens because **the app files don't exist**. config.json references:
- `apps/briefing/App.tsx` — folder does not exist in audos-workspace/
- `apps/signature/App.tsx` — folder does not exist
- `apps/studio/App.tsx` — folder does not exist

Only `apps/home/` and `apps/setup/` have real implementations. The `src/pages/apps/Studio.tsx` etc. are in the wrong location and ignored by the platform.

Fix: either create the missing files, or consolidate to single app (recommended).

---

## Hooks Reference

| Hook | Notes |
|------|-------|
| `useSpaceRuntime()` | Session, config, file ops, subscription state |
| `useWorkspaceDB(table, options)` | PostgreSQL via `window.useWorkspaceDB` |
| `useSubscription()` | Subscription status, plan tier, trial state |
| `useBranding` | ❌ Does not exist — use `useSpaceRuntime().config.desktop.branding` |
| `useSession` | ❌ Does not exist — use `useSpaceRuntime().sessionId` |

---

## Confirmed: Single-App Architecture Test Results (2026-04-02)

Live test passed. App running at workspace-351699 as `apps/throughline/App.tsx`.

| Test | Result | Notes |
|------|--------|-------|
| Single app with internal routing | ✅ | `useState<Page>` routing works |
| DB query — guest_prep_podcast_profiles | ✅ | Query ran, 0 rows (empty workspace) |
| DB query — voice_profiles | ✅ | Query ran, 0 rows |
| ShadCN-style components (no Radix) | ✅ | Button, Card, Badge render correctly |
| Radix UI primitives | ⏳ | Not yet tested |
| TanStack Query | ⏳ | Not yet tested |

**Next:** Add `@radix-ui/react-*` and `@tanstack/react-query` to `cdnDependencies` in config.json and test.

---

## Confirmed: New Apps Require Otto Registration

**Adding a new app to config.json via GitHub push does NOT create a dock entry.**

The Audos platform maintains its own app registry (DB records). The `apps` array in config.json maps app IDs to component files — but the platform must have a matching DB record for the app to appear in the dock. That record must be created via Otto.

**Workflow for adding a new app:**
1. Add the entry to `config.json` (so the platform knows the component path)
2. Create the `App.tsx` file in `audos-workspace/apps/<id>/`
3. Push to GitHub (triggers sync + recompile)
4. Ask Otto to register the app record in the platform DB

Nick has partially automated the sync, but step 4 still requires Otto. This may change in future platform updates.

Existing apps (home, briefing, signature, studio, setup) were registered when the workspace was initialized — no manual step required for them.

---

## New Files Referenced (Not Yet in Repo)

SDK-09 referenced these files from the platform runtime — request from Auto if needed:
- `APP_INTEGRATION_MANIFEST.md` — all 25+ platform integrations
- `SPACE_APP_GUIDE.md` — app development guide
- `integrations/workspace-db/docs.md` — WorkspaceDB SDK documentation
