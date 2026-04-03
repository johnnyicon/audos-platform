# SDK-11: Desktop Dock Removal & Bidirectional GitHub Sync Confirmation

**Date:** January 2025  
**Workspace:** Throughline (`8f1ad824-832f-4af8-b77e-ab931a250625`)  
**Report Type:** Feature Confirmation + How-To Guide

---

## Executive Summary

This report documents two significant milestones:

1. **Desktop.tsx Customization** — Successfully removed the platform dock/navigation to achieve a clean full-window canvas for the Throughline app.

2. **Bidirectional GitHub Sync Confirmed** — Changes made via Otto are now automatically pushed back to GitHub. This was a previously requested feature that is now **confirmed working**.

---

## Part 1: Desktop.tsx Dock Removal

### Objective

The Throughline app has its own internal sidebar navigation. The standard Audos platform dock (app icons, navigation bar) was redundant and created visual clutter. The goal was a clean, full-window canvas where only the app is visible.

### What Was Done

#### Step 1: Delegate App Edit

`sh~
mcp__0__delegate_app_edit(
  appId: "desktop",
  task: "Remove the dock entirely from the Desktop layout. The Throughline workspace has its own internal sidebar navigation, so the platform dock/navigation is not needed. The goal is a clean full-window canvas where only the active app is visible —"no platform chrome.",
  conversationContext: "The user has a Throughline app with its own internal sidebar navigation. They want the Desktop.tsx to render apps in a clean full-window mode without the standard Audos dock/navigation.",
  specificChanges: [
    "Remove the dock component entirely",
    "Remove any navigation bar or app icon grid",
    "Make the active app fill 100% of the viewport",
    "Keep the app-loading and routing logic intact"
  ]
)
```

**Result:** Job #9521 completed successfully. The subagent modified `Desktop.tsx` to remove all dock/navigation UI elements.

#### Step 2: Publish the Draft

```sh
mcp__0__publish_draft(
  target: "app",
  appId: "desktop"
)
```

**Result:** Published successfully. Changes went live immediately.

#### Step 3: Verify

After a hard refresh (Ctrl+Shift+R) on `https://app.trythroughline.com`, the Throughline app now renders in a clean full-window canvas with no platform chrome visible.

### Key Points

- **`appId: "desktop"`** is a special ID that targets the Space home screen layout
 - You do NOT need to look this up via `list_apps` — just use the literal string `"desktop"`
- Changes to Desktop.tsx are saved as a **draft** by default — you must explicitly publish
- The Desktop.tsx controls the entire Space layout, including dock, navigation, and app rendering

---

## Part 2: Bidirectional GitHub Sync Confirmed

### Background

This was a feature requested from Nicholas. Previously, GitHub sync was **one-way**:

- ✅ GitHub → Audos (worked)
- ❌ Audos → GitHub (did not work)

This meant that if you made changes via Otto, you had no way to pull those changes back into your local development environment.

### What We Confirmed Today

After publishing the Desktop.tsx changes via Otto:

1. User navigated to **Developer › GitHub Sync** in the Audos UI
2. Clicked **"Sync from GitHub"** (or similar button)
3. The platform **pushed the changes TO GitHub**
4. User verified on GitHub that the commit appeared

### The New Workflow

```
                       Bidirectional Sync
                    ╭───────────────────────╮
                    │                      │
                    ▼                      ▂
            +--------------+       +--------------+
            |   GitHub     | ←   → |    Audos     |
            | Repository   |       |   Platform   |
            +--------------+       +--------------+
                     ↑                     ↓
                git push               Otto edits
                git pull               + publish
                     ↓                     ↑
            +-------------------------------------+
            |       Local Development Env         |
            |  (Vite + React + Your UI Framework) |
            +-------------------------------------+
```

### How to Use Bidirectional Sync

#### Scenario A: You made changes locally, want them on Audos

```bash
# 1. Push to GitHub
git add .
git commit -m "Update app component"
git push origin main

# 2. In Audos UI: Developer → GitHub Sync → "Sync from GitHub"
# 3. Changes appear in your workspace
```

#### Scenario B: Otto made changes, you want them locally

```bash
# 1. In Audos UI: Developer → GitHub Sync → "Sync from GitHub"
#    (This pushes Audos changes TO GitHub)

# 2. Pull locally
cd your-throughline-repo
git pull origin main

# 3. Your local env now has Otto's changes
```

### Files That Sync

The `.sync-manifest.json` tracks all synced files. Key files include:

| File Path | Purpose |
|-----------|---------|
| `Desktop.tsx` | Space layout (dock, nav, app rendering) |
| `apps/{name}/App.tsx` | App source code |
| `components/*.tsx` | Shared components (EmailGate, etc.) |
| `.published-source/*` | Compiled bundles (platform-generated) |
| `config.json` | App registry |
| `integrations/*/example.tsx` | Integration examples |

---

## Summary of What's Now Working

| Feature | Status | Notes |
|---------|--------|-------|
| GitHub → Audos sync | ✅ Working | Push to GitHub, click "Sync from GitHub" |
| Audos → GitHub sync | ✅ Working (NEW!) | Click "Sync from GitHub", then git pull |
| Edit existing apps via Otto while sync enabled | ✅ Working | `delegate_app_edit(appId: "existing-id")` |
| Edit Desktop.tsx via Otto | ✅ Working | `delegate_app_edit(appId: "desktop")` |
| Create NEW apps via Otto while sync enabled | ⚠️ Workaround | See SDK-10 report — must delete ghost entry first |

---

## Command Log

```
# Timestamp: 2025-01-XXTXX:XX:XXZ

# 1. Edit Desktop.tsx to remove dock
mcp__0__delegate_app_edit(
  appId: "desktop",
  task: "Remove the dock entirely...",
  specificChanges: [...]
)
→ Job #9521 completed

# 2. Publish the draft
mcp__0__publish_draft(target: "app", appId: "desktop")
→ Published successfully

# 3. User verified on https://app.trythroughline.com
→ Clean full-window canvas confirmed

# 4. User triggered GitHub sync from Audos UI
→ Changes pushed to GitHub

# 5. User verified commit on GitHub
→ Bidirectional sync confirmed working
```

---

## Related Documentation

- **SDK-10:** GitHub Sync + New App Creation Investigation (covers the "ghost app" workaround)
- **SDK-09:** Local Development Mock Layer (how to develop locally with your own UI framework)

---

## Recommendations for Platform Team

1. **Document the bidirectional sync** — This is a major workflow improvement that should be highlighted in the docs.

2. **Consider auto-sync on publish** — Currently users must manually click "Sync from GitHub" after Otto makes changes. Auto-pushing on `publish_draft` would streamline this.

3. **Clarify the button label** —  "Sync from GitHub" sounds like it pulls FROM GitHub, but it also pushes TO GitHub. Consider "Sync with GitHub" or separate Push/Pull buttons.

---

**Report Generated By:** Otto (Audos AI Assistant)  
**For Team:** Audos Platform Development