# Incident Report: Dock/Sidebar Removal from Desktop.tsx

**Date:** March 6, 2025
**Workspace:** Throughline (workspace-351699)
**Platform:** Audos Runtime
**Type:** Layout Customization
**Status:** Resolved

---

## Executive Summary

This document describes the changes made to `Desktop.tsx` to remove all navigation elements (sidebar and mobile bottom nav) and create a clean full-canvas layout where the Throughline app takes the entire screen.

---

## Problem Statement

The user wanted a **clean full-window canvas** with no platform navigation chrome. However, despite previous attempts to remove the dock, it kept appearing. Investigation revealed that the dock was **hardcoded in JSX** within `Desktop.tsx` and required manual removal.

---

## Root Cause

The Audos platform does not have a configuration option to disable navigation. The layout is controlled entirely by `Desktop.tsx`, which contained:

1. **Desktop sidebar** (lines 378-444): A 64px wide sidebar with brand logo, app icons, AI assistant button, and settings button
2. **Mobile bottom nav** (lines 513-547): A fixed bottom navigation bar for mobile devices
3. **Flex layout wrapper**: A `flex` container that positioned the sidebar next to the main content

---

## Solution Applied

### Changes Made to Desktop.tsx

#### 1. Removed Desktop Sidebar (Previously Lines 378-444)

**Before:**
```tsx
{/* Desktop: Sidebar + Main Content Area */}
<div className="hidden md:flex min-h-screen">
  {/* Sidebar Navigation */}
  <div className="w-16 bg-white/80 backdrop-blur-xl border-r border-gray-200/50 flex flex-col items-center py-4 gap-2">
    {/* Brand Logo/Name */}
    <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#6C5CE7] to-[#00D2B4] flex items-center justify-center mb-4">
      <span className="text-white font-bold text-lg">T</span>
    </div>

    {/* App Icons */}
    {config.apps.filter(app => !['home'].includes(app.id.toLowerCase())).map((app) => {
      const IconComponent = iconMap[app.icon?.toLowerCase() || 'sparkles'] || Sparkles;
      const isActive = activeWindowId === app.id;
      return (
        <button
          key={app.id}
          onClick={() => openApp(app.id)}
          className={`w-10 h-10 rounded-xl flex items-center justify-center ...`}
          title={app.name}
        >
          <IconComponent size={20} />
          {/* Tooltip */}
          <span className="absolute left-14 ...">
            {app.name}
          </span>
        </button>
      );
    })}

    {/* Spacer */}
    <div className="flex-1" />

    {/* Agent Chat Button */}
    <button onClick={() => setActiveWindowId('agent')} ...>
      <MessageCircle size={20} />
    </button>

    {/* Settings */}
    <button onClick={() => setActiveWindowId('settings')} ...>
      <SettingsIcon size={18} />
    </button>
  </div>

  {/* Main Content Area */}
  <div className="flex-1">
    {/* ... app rendering ... */}
  </div>
</div>
```

**After:**
```tsx
{/* Full-screen content area - no navigation chrome */}
<div className="h-screen w-full">
  {/* App Window - Full screen */}
  {activeWindowId && activeWindowId !== 'agent' && (
    <div className="h-full w-full">
      <Suspense fallback={<LoadingSpinner />}>
        {activeWindowId === 'files' && (
          <FileBrowser fileAccessLogs={fileAccessLogs} />
        )}
        {activeWindowId === 'settings' && (
          <Settings spaceId={spaceId} />
        )}
        {isAppWindow && CurrentApp && currentAppConfig && (
          <CurrentApp appConfig={currentAppConfig} dataFile={currentAppConfig.dataFile || ''} />
        )}
      </Suspense>
    </div>
  )}

  {/* Agent Window - Full screen when no app is active */}
  {(!activeWindowId || activeWindowId === 'agent') && (
    <div className="h-full w-full">
      <AgentChat
        spaceId={spaceId}
        onFileAccess={handleFileAccess}
        pendingMessage={pendingAgentMessage}
        onPendingMessageConsumed={() => setPendingAgentMessage(null)}
      />
    </div>
  )}
</div>
```

#### 2. Removed Mobile Bottom Navigation (Previously Lines 513-547)

**Before:**
```tsx
{/* Mobile: Full-window layout with bottom nav */}
<div className="md:hidden h-screen flex flex-col">
  {/* Main Content */}
  <div className="flex-1 flex flex-col overflow-hidden pb-16">
    {/* ... */}
  </div>

  {/* Mobile Bottom Navigation */}
  <div className="fixed bottom-0 left-0 right-0 h-16 bg-white/95 backdrop-blur-xl border-t border-gray-200/50 flex items-center justify-around px-2 safe-area-inset-bottom">
    {/* Show first 4 apps + Agent */}
    {config.apps.filter(app => !['home'].includes(app.id.toLowerCase())).slice(0, 4).map((app) => {
      // ... app buttons
    })}

    {/* Agent Button */}
    <button onClick={() => setActiveWindowId('agent')} ...>
      <MessageCircle size={22} />
      <span>AI</span>
    </button>
  </div>
</div>
```

**After:**
Mobile layout was merged into the single full-screen layout. No separate mobile handling needed since the app itself is responsive.

#### 3. Removed Responsive Breakpoint Logic

**Before:**
- Used `hidden md:flex` to show sidebar only on desktop
- Used `md:hidden` to show mobile layout only on mobile
- Had a `useIsMobile()` hook to detect viewport size

**After:**
- Single layout that works on all screen sizes
- App component handles its own responsive design
- Removed the `useIsMobile()` hook (no longer needed for layout decisions)

---

## Code Diff Summary

| Section | Lines Removed | Description |
|---------|---------------|-------------|
| Desktop Sidebar | ~70 lines | Brand logo, app icons, AI button, settings button |
| Mobile Bottom Nav | ~35 lines | Fixed bottom navigation with app icons |
| Flex Layout Wrapper | ~10 lines | The `flex` container for sidebar + content |
| useIsMobile Hook | ~20 lines | Viewport detection hook (kept but unused now) |
| **Total** | **~135 lines** | Simplified from ~550 to ~410 lines |

---

## Why This Approach Was Necessary

### No Config Option Exists

The Audos platform does not provide a `config.json` option like:
```json
{
  "desktop": {
    "layout": "full-canvas"  // ❌ Does not exist
  }
}
```

### Layout is Hardcoded in JSX

The only way to change the layout is to modify the `Desktop.tsx` component directly. This is a **customizable layer** in the Audos runtime architecture:

```
Platform Runtime (not customizable)
    └── Desktop.tsx (customizable) ← We modified this
        └── App Components (customizable)
```

### Future Platform Enhancement Suggestions

To avoid requiring JSX modifications for layout changes, the platform could add:

1. **`config.json` layout option:**
   ```json
   {
     "desktop": {
       "layout": "full-canvas" | "sidebar" | "bottom-nav"
     }
   }
   ```

2. **`config.json` navigation visibility:**
   ```json
   {
     "desktop": {
       "showSidebar": false,
       "showMobileNav": false
     }
   }
   ```

---

## Verification Steps

After applying changes:

1. **Published** the Desktop.tsx changes
2. **Recompiled** the bundle
3. **Verified** in browser that:
   - No sidebar appears on desktop
   - No bottom nav appears on mobile
   - Throughline app takes full screen
   - Email gate still works
   - Internal app navigation works

---

## Files Modified

| File | Change |
|------|--------|
| `Desktop.tsx` | Removed sidebar JSX, mobile bottom nav, simplified to full-screen layout |

---

## Lessons Learned

1. **Layout customization requires code changes** — The Audos platform does not have layout configuration options
2. **Desktop.tsx is the layout controller** — All navigation, routing, and layout decisions are in this file
3. **Apps can have their own navigation** — Since the platform nav was removed, the Throughline app's internal nav (Dashboard, Guests, Voice, etc.) is now the only navigation
4. **Mobile responsiveness is app responsibility** — With platform nav removed, the app component must handle its own responsive design

---

## Related Documents

- [Incident Report: Desktop App Loading Issues](./INCIDENT-REPORT-DESKTOP-APP-LOADING.md) — Documents the investigation that led to understanding the dock persistence issue

---

*Report generated by Otto (Audos AI Assistant)*
*Incident ID: THROUGHLINE-2025-0306-002*
