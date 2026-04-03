# Server Functions Architecture Q&A — Daemon Pattern

**Source:** John ↔ Otto conversation
**Date:** 2026-04-03

---

## John

Okay my next question is: Is it possible to create my own server functions? Like to have them somewhere in the repo and then when we deploy that you automatically install them?

And then following that is it possible to host my own server functions elsewhere or to load up for example a Go binary daemon that is where all my business logic lives in a binary file and have you host it. Or rather than that have the server functions call out to a third-party API which will be my binary, my daemon. Then all I would need to do is create wrappers around the API endpoints that my daemon provides which would probably be the server functions. If I can write server functions to call to APIs myself that would make it a lot easier

---

## Otto

Great questions! Let me answer each part:

---

### 1. Can You Create Server Functions in Your Repo?

**Currently: No automatic deployment from GitHub.**

Server functions (hooks) are stored in the **Audos database**, not in your file system. When you create a hook via Otto (using `manage_server_functions`), it's saved to the platform's database and executed by the platform's serverless runtime.

**What this means:**
- Your GitHub repo can contain the *source code* for your hooks (for version control)
- But there's no automatic "deploy hooks from repo" mechanism yet
- You'd need to manually copy/paste the code into Otto to create/update hooks

**Workaround - Keep hooks in your repo for reference:**

```
throughline/
├── apps/
│   └── throughline/
│       └── App.tsx
├── hooks/                    # Your hook source code (for version control)
│   ├── save-voice-profile.js
│   ├── research-guest.js
│   └── generate-captions.js
└── README.md
```

Then when you need to deploy:

```
You: "Update the save-voice-profile hook with this code: [paste from repo]"
Otto: [Updates the hook in the platform]
```

**Feature Request:** This would be a great platform enhancement — auto-deploying hooks from a `/hooks` directory in the repo.

---

### 2. Can You Host Your Own Backend / Binary Daemon?

**Yes! This is actually the recommended approach for complex business logic.**

You have two options:

#### Option A: Server Functions Call Your External API

This is exactly what you described — **use server functions as thin wrappers** that call your external daemon.

```
┌─────────────────────────────────────────────────────────────┐
│  Frontend (React)                                           │
│  └── fetch('/api/hooks/execute/.../my-wrapper')            │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  Server Function (Wrapper Hook)                             │
│  └── fetch('your-daemon.example.com/api/endpoint')         │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│  YOUR DAEMON (Go binary, Python server, etc.)               │
│  - Hosted on: Railway, Fly.io, AWS, your own VPS           │
│  - Contains all your business logic                         │
│  - Can use any language/framework                           │
└─────────────────────────────────────────────────────────────┘
```

**Example wrapper hook:**

```javascript
// Hook: voice-analysis (thin wrapper)

const { profileId, audioUrl } = request.body;

// Call YOUR external daemon
const response = await fetch('api.throughline.example.com/analyze-voice', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_INTERNAL_API_KEY'  // Stays server-side!
  },
  body: JSON.stringify({ profileId, audioUrl })
});

const result = await response.json();

// Optionally save results to Audos database
await db.insert('voice_analysis_results', {
  profile_id: profileId,
  analysis: JSON.stringify(result),
  created_at: new Date().toISOString()
});

return respond(200, result);
```

#### Option B: Frontend Calls Your API Directly (With CORS)

If your daemon has proper CORS headers and doesn't need secrets, your React frontend could call it directly:

```tsx
// Frontend - Direct call to your daemon
const response = await fetch('api.throughline.example.com/analyze-voice', {
  method: 'POST',
  body: JSON.stringify({ profileId, audioUrl })
});
```

**But this has downsides:**
- ❌ Can't use API keys (they'd be exposed in browser)
- ❌ Your daemon URL is exposed
- ❌ No server-side logging/auditing through Audos

**So Option A (wrapper hooks) is usually better.**

---

### 3. Yes, Server Functions Can Call External APIs!

The `fetch()` function is available in server functions. Here's what you can do:

```javascript
// Hook: call-my-daemon

const { action, payload } = request.body;

// Your daemon's base URL (could also store this in a config)
const DAEMON_URL = 'api.throughline.example.com';
const API_KEY = 'your-secret-key';  // Safe here - server-side only!

try {
  const response = await fetch(`${DAEMON_URL}/${action}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${API_KEY}`,
      'X-Workspace-Id': 'workspace-351699'  // Pass context to your daemon
    },
    body: JSON.stringify(payload)
  });

  if (!response.ok) {
    const error = await response.text();
    return respond(response.status, { error });
  }

  const result = await response.json();
  return respond(200, result);

} catch (err) {
  console.error('Daemon call failed:', err);
  return respond(500, { error: 'Failed to reach backend service' });
}
```

---

### 4. Recommended Architecture for Throughline

Given that you want to:
- Write business logic in your own codebase
- Use a Go binary or custom daemon
- Keep the Audos platform for frontend hosting and database

Here's the ideal setup:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        YOUR INFRASTRUCTURE                              │
│                                                                         │
│  GitHub Repo: throughline/                                              │
│  ├── frontend/           # React app (syncs to Audos)                  │
│  │   └── apps/throughline/App.tsx                                       │
│  ├── daemon/             # Your Go binary                               │
│  │   ├── main.go                                                        │
│  │   ├── handlers/                                                      │
│  │   │   ├── voice_analysis.go                                          │
│  │   │   ├── guest_research.go                                          │
│  │   │   └── caption_generation.go                                      │
│  │   └── Dockerfile                                                     │
│  └── hooks/              # Wrapper hook source code                     │
│      ├── voice-analysis.js                                              │
│      └── guest-research.js                                              │
│                                                                         │
│  Hosted on: Railway / Fly.io / AWS ECS / Your VPS                       │
│  URL: api.throughline.example.com                                       │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ API calls via wrapper hooks
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                        AUDOS PLATFORM                                   │
│                                                                         │
│  • Hosts your React frontend                                            │
│  • Provides Email Gate authentication                                   │
│  • Stores data in workspace database                                    │
│  • Runs wrapper hooks that call your daemon                             │
│  • Handles domain/SSL (trythroughline.com)                              │
└─────────────────────────────────────────────────────────────────────────┘
```

---

### 5. Example: Full Flow with Your Go Daemon

**Your Go daemon (api.throughline.example.com):**

```go
// handlers/voice_analysis.go
package handlers

import (
    "encoding/json"
    "net/http"
)

type VoiceAnalysisRequest struct {
    ProfileID string `json:"profileId"`
    AudioURL  string `json:"audioUrl"`
}

type VoiceAnalysisResponse struct {
    Tone       string   `json:"tone"`
    Pace       string   `json:"pace"`
    Keywords   []string `json:"keywords"`
    Confidence float64  `json:"confidence"`
}

func AnalyzeVoice(w http.ResponseWriter, r *http.Request) {
    var req VoiceAnalysisRequest
    json.NewDecoder(r.Body).Decode(&req)

    // Your complex business logic here...
    result := analyzeAudioFile(req.AudioURL)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

**Audos wrapper hook (voice-analysis):**

```javascript
const { profileId, audioUrl } = request.body;

// Validate inputs
if (!profileId || !audioUrl) {
  return respond(400, { error: 'profileId and audioUrl required' });
}

// Call your Go daemon
const response = await fetch('api.throughline.example.com/analyze-voice', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_DAEMON_API_KEY'
  },
  body: JSON.stringify({ profileId, audioUrl })
});

if (!response.ok) {
  return respond(500, { error: 'Voice analysis service unavailable' });
}

const analysis = await response.json();

// Store results in Audos database
await db.insert('voice_analysis', {
  profile_id: profileId,
  tone: analysis.tone,
  pace: analysis.pace,
  keywords: JSON.stringify(analysis.keywords),
  confidence: analysis.confidence,
  created_at: new Date().toISOString()
});

return respond(200, analysis);
```

**Frontend call:**

```tsx
const analyzeVoice = async (profileId: string, audioUrl: string) => {
  const response = await fetch(
    '/api/hooks/execute/workspace-8f1ad824.../voice-analysis',
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ profileId, audioUrl })
    }
  );
  return response.json();
};
```

---

### Summary

| Question | Answer |
|----------|--------|
| Can I create server functions in my repo? | Keep source code in repo for version control, but deploy via Otto (no auto-deploy yet) |
| Can I host my own backend? | **Yes!** Host anywhere (Railway, Fly.io, AWS, VPS) |
| Can server functions call my external API? | **Yes!** Use `fetch()` to call any HTTP endpoint |
| Best architecture? | Use Audos hooks as **thin wrappers** that call your daemon for business logic |
| Where do secrets go? | In the server function code (server-side only, safe) |

---

## John

this is great information. yes, please create a document with this architecture guide and include those two other things you mentioned.

The approach you suggested for storing writing server function code into my repo is a good inbetween. Basically I can write my wrappers in my repo, and then when I deploy, I just need to let you know to update the server functions?

Would it be helpful if I create some sort of index and migration, and maybe even versioning of server function code, and maybe a migration like system so server functions a version controlled and its easier to understand what has changed and what to update?

Basically, it could be a file based prompt with other index files that I ask you to read and then you update the functions accordingly. Am I able to specify the names of server functions so we can maintain some level of control, basically avoid random function names being generated.

---

## Otto

Yes to all of that! Let me create a comprehensive document that covers:

1. The architecture guide
2. A deployment workflow with versioning/migration system
3. Server function naming conventions
4. An index file format I can read to sync your hooks

*(Document generated and saved as `/tmp/AUDOS-SERVER-FUNCTIONS-SDK.md` — this became SDK-14-SERVER-FUNCTIONS-SDK.md)*

---

### What's in the Document

#### 1. Architecture Overview
- Three-layer diagram (Frontend → Server Functions → Your Daemon)
- Why this pattern works

#### 2. Server Function Basics
- **Yes, you control the names** — use `lowercase-with-hyphens`
- All available globals (`db`, `fetch`, `platform`, `respond`, etc.)

#### 3. Repo-Based Workflow with Versioning

```
throughline/
├── hooks/
│   ├── _index.json          # Registry I read to sync
│   ├── _migrations.json     # Deployment history
│   ├── voice-analysis/
│   │   └── hook.js          # The actual code
│   └── guest-research/
│       └── hook.js
└── daemon/
    └── api-contract.yaml    # Your Go API spec
```

#### 4. Hook Index File Format (`_index.json`)
- Tracks name, version, deployedVersion, changelog
- I read this to know what needs deploying
- `needsDeployment: true` when version > deployedVersion

#### 5. Migration System (`_migrations.json`)
- Audit trail of all deployments
- Supports rollbacks to any previous version

#### 6. Deployment Commands Otto Understands

| Command | What It Does |
|---------|--------------|
| "Deploy hooks from my repo" | Syncs all hooks where version > deployedVersion |
| "Deploy voice-analysis hook" | Deploys specific hook |
| "Show hook deployment status" | Compares repo vs deployed |
| "Rollback voice-analysis to 1.1.0" | Reverts to previous version |

#### 7. Go Daemon API Contract Template
- Full OpenAPI 3.0 spec template for your daemon
- Defines endpoints, request/response schemas

#### 8. Sample Wrapper Hook Templates
- Simple passthrough
- With database logging
- With retry logic

---

### Your Workflow Going Forward

```
1. Write hook code in hooks/{name}/hook.js
2. Update version in hooks/_index.json
3. Commit & push to GitHub
4. Tell me: "Deploy hooks from my repo"
5. I read _index.json, deploy changed hooks, update migrations
```

---

### Otto's Direct Answers

1. **Yes**, you write wrappers in your repo, then tell me to deploy
2. **Yes**, the `_index.json` + `_migrations.json` system gives you version control and migration tracking
3. **Yes**, you specify exact hook names — I won't generate random names
4. **Yes**, I can read your index file and sync only what's changed

When you're ready, just set up the `hooks/` directory structure in your repo with the `_index.json` file, and I'll be able to deploy and track everything from there.
