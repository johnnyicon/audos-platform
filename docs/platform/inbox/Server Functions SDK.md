# Audos Server Functions SDK & Deployment Guide

**Date:** March 6, 2025
**Platform:** Audos Runtime
**Document Type:** SDK Reference & Workflow Guide

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Server Function Basics](#2-server-function-basics)
3. [Calling External APIs (Your Daemon)](#3-calling-external-apis-your-daemon)
4. [Repo-Based Workflow with Versioning](#4-repo-based-workflow-with-versioning)
5. [Hook Index File Format](#5-hook-index-file-format)
6. [Migration System](#6-migration-system)
7. [Deployment Commands](#7-deployment-commands)
8. [Go Daemon API Contract Template](#8-go-daemon-api-contract-template)
9. [Sample Wrapper Hook Templates](#9-sample-wrapper-hook-templates)

---

## 1. Architecture Overview

### Three-Layer Architecture

```
┌──────────────────────────────────────────────────────────────────────────┐
│                           USER'S BROWSER                                 │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │  FRONTEND LAYER (React Components)                                 │  │
│  │  • UI rendering, user interactions                                 │  │
│  │  • Calls server functions via fetch()                              │  │
│  │  ⚠️  NO secrets, NO sensitive logic                                │  │
│  └────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│                        AUDOS PLATFORM                                    │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │  SERVER FUNCTIONS (Hooks) - Thin Wrappers                          │  │
│  │  • Validate inputs                                                 │  │
│  │  • Call YOUR external daemon                                       │  │
│  │  • Store results in Audos database                                 │  │
│  │  • Return responses to frontend                                    │  │
│  └────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│                     YOUR INFRASTRUCTURE                                  │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │  YOUR DAEMON (Go binary, Python, Node, etc.)                       │  │
│  │  • All complex business logic                                      │  │
│  │  • AI/ML processing                                                │  │
│  │  • External API integrations                                       │  │
│  │  • Hosted on: Railway, Fly.io, AWS, your VPS                       │  │
│  └────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
```

### Why This Architecture?

| Concern | Where It Lives | Why |
|---------|----------------|-----|
| UI/UX | Frontend (React) | Fast, interactive |
| Input validation | Server Functions | Can't bypass server-side checks |
| API keys/secrets | Server Functions | Never exposed to browser |
| Complex business logic | Your Daemon | Full control, any language |
| Data persistence | Audos Database | Managed, backed up |
| Authentication | Audos Email Gate | Built-in, secure |

---

## 2. Server Function Basics

### Naming Conventions

**Yes, you can specify exact names for server functions.** Names should be:

- **Lowercase with hyphens:** `save-voice-profile`, `research-guest`
- **Descriptive:** Name describes the action
- **Versioned (optional):** `save-voice-profile-v2` for breaking changes

### Available Globals in Server Functions

```javascript
// Request data
request.body      // POST body (parsed JSON)
request.query     // URL query parameters
request.method    // HTTP method (GET, POST, etc.)
request.headers   // Request headers

// Database operations
db.query(table, options)    // Read records
db.insert(table, data)      // Insert record
db.update(table, filters, data)  // Update records
db.delete(table, filters)   // Delete records
db.listTables()             // List all tables

// Platform services
platform.generateText({ prompt, maxTokens })  // AI text generation
platform.sendEmail({ to, subject, text, html })  // Send email

// Utilities
fetch(url, options)         // HTTP requests to external APIs
respond(statusCode, body)   // Return response
JSON, Date, Math, console   // Standard JS globals
```

### Basic Server Function Structure

```javascript
// Hook name: my-function-name

// 1. Extract and validate inputs
const { param1, param2 } = request.body;

if (!param1) {
  return respond(400, { error: 'param1 is required' });
}

// 2. Do work (call daemon, query database, etc.)
const result = await someOperation();

// 3. Return response
return respond(200, { success: true, data: result });
```

---

## 3. Calling External APIs (Your Daemon)

### Wrapper Hook Pattern

```javascript
// Hook name: voice-analysis
// Description: Wrapper that calls the Throughline daemon for voice analysis
// Version: 1.0.0

const DAEMON_URL = 'https://api.throughline.example.com';
const API_KEY = 'your-secret-api-key';  // Safe here - server-side only

const { profileId, audioUrl } = request.body;

// Validate inputs
if (!profileId) {
  return respond(400, { error: 'profileId is required' });
}

if (!audioUrl) {
  return respond(400, { error: 'audioUrl is required' });
}

try {
  // Call YOUR daemon
  const response = await fetch(`${DAEMON_URL}/api/v1/analyze-voice`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${API_KEY}`,
      'X-Request-ID': `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
    },
    body: JSON.stringify({
      profileId,
      audioUrl,
      workspaceId: 'workspace-351699'
    })
  });

  if (!response.ok) {
    const errorText = await response.text();
    console.error('Daemon error:', response.status, errorText);
    return respond(response.status, {
      error: 'Voice analysis service error',
      details: errorText
    });
  }

  const analysis = await response.json();

  // Optionally store results in Audos database
  await db.insert('voice_analysis_results', {
    profile_id: profileId,
    audio_url: audioUrl,
    analysis: JSON.stringify(analysis),
    created_at: new Date().toISOString()
  });

  return respond(200, analysis);

} catch (err) {
  console.error('Failed to reach daemon:', err);
  return respond(503, { error: 'Voice analysis service unavailable' });
}
```

---

## 4. Repo-Based Workflow with Versioning

### Recommended Directory Structure

```
throughline/
├── apps/                          # Frontend React apps (syncs to Audos)
│   └── throughline/
│       └── App.tsx
│
├── hooks/                         # Server function source code
│   ├── _index.json                # Hook registry (Otto reads this)
│   ├── _migrations.json           # Migration history
│   │
│   ├── voice-analysis/
│   │   ├── hook.js                # The actual hook code
│   │   ├── README.md              # Documentation
│   │   └── test.json              # Test payloads
│   │
│   ├── guest-research/
│   │   ├── hook.js
│   │   ├── README.md
│   │   └── test.json
│   │
│   └── generate-captions/
│       ├── hook.js
│       ├── README.md
│       └── test.json
│
├── daemon/                        # Your Go/Python backend
│   ├── main.go
│   ├── handlers/
│   ├── Dockerfile
│   └── api-contract.yaml          # OpenAPI spec
│
└── docs/
    └── deployment.md
```

### Hook Source File Format

Each hook lives in its own directory with a `hook.js` file:

```javascript
// hooks/voice-analysis/hook.js

/**
 * @name voice-analysis
 * @version 1.2.0
 * @description Analyzes voice characteristics by calling the Throughline daemon
 * @author Johnny
 * @created 2025-03-06
 * @modified 2025-03-10
 *
 * @param {string} profileId - The voice profile ID
 * @param {string} audioUrl - URL to the audio file
 * @returns {object} Analysis results with tone, pace, keywords
 */

const DAEMON_URL = 'https://api.throughline.example.com';
const API_KEY = 'your-secret-api-key';

const { profileId, audioUrl } = request.body;

if (!profileId || !audioUrl) {
  return respond(400, { error: 'profileId and audioUrl are required' });
}

try {
  const response = await fetch(`${DAEMON_URL}/api/v1/analyze-voice`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${API_KEY}`
    },
    body: JSON.stringify({ profileId, audioUrl })
  });

  if (!response.ok) {
    return respond(response.status, { error: 'Daemon error' });
  }

  const result = await response.json();

  await db.insert('voice_analysis_results', {
    profile_id: profileId,
    analysis: JSON.stringify(result),
    created_at: new Date().toISOString()
  });

  return respond(200, result);
} catch (err) {
  console.error('Daemon call failed:', err);
  return respond(503, { error: 'Service unavailable' });
}
```

---

## 5. Hook Index File Format

The `_index.json` file is the **single source of truth** that Otto reads to understand what hooks exist and their current versions.

### hooks/_index.json

```json
{
  "$schema": "https://throughline.example.com/schemas/hooks-index.json",
  "version": "1.0.0",
  "workspace": "workspace-351699",
  "daemonUrl": "https://api.throughline.example.com",
  "lastDeployed": "2025-03-10T14:30:00Z",

  "hooks": [
    {
      "name": "voice-analysis",
      "version": "1.2.0",
      "file": "voice-analysis/hook.js",
      "description": "Analyzes voice characteristics via daemon",
      "status": "active",
      "deployedVersion": "1.2.0",
      "lastDeployed": "2025-03-10T14:30:00Z",
      "changelog": [
        { "version": "1.2.0", "date": "2025-03-10", "changes": "Added error retry logic" },
        { "version": "1.1.0", "date": "2025-03-08", "changes": "Added database logging" },
        { "version": "1.0.0", "date": "2025-03-06", "changes": "Initial version" }
      ]
    },
    {
      "name": "guest-research",
      "version": "1.0.0",
      "file": "guest-research/hook.js",
      "description": "AI-powered guest research via daemon",
      "status": "active",
      "deployedVersion": "1.0.0",
      "lastDeployed": "2025-03-06T10:00:00Z",
      "changelog": [
        { "version": "1.0.0", "date": "2025-03-06", "changes": "Initial version" }
      ]
    },
    {
      "name": "generate-captions",
      "version": "2.0.0",
      "file": "generate-captions/hook.js",
      "description": "Generates platform-specific social media captions",
      "status": "active",
      "deployedVersion": "1.5.0",
      "needsDeployment": true,
      "lastDeployed": "2025-03-08T16:00:00Z",
      "changelog": [
        { "version": "2.0.0", "date": "2025-03-10", "changes": "BREAKING: New response format" },
        { "version": "1.5.0", "date": "2025-03-08", "changes": "Added TikTok support" },
        { "version": "1.0.0", "date": "2025-03-06", "changes": "Initial version" }
      ]
    },
    {
      "name": "old-function",
      "version": "1.0.0",
      "file": "old-function/hook.js",
      "description": "Deprecated function",
      "status": "deprecated",
      "deprecatedReason": "Replaced by voice-analysis",
      "removeAfter": "2025-04-01"
    }
  ],

  "secrets": {
    "DAEMON_API_KEY": "{{ DAEMON_API_KEY }}",
    "OPENAI_KEY": "{{ OPENAI_KEY }}"
  },

  "notes": "Secrets are placeholders - actual values stored securely in Otto"
}
```

### Key Fields Explained

| Field | Purpose |
|-------|---------|
| `name` | Exact name of the hook (you control this) |
| `version` | Current version in your repo |
| `deployedVersion` | Version currently deployed to Audos |
| `needsDeployment` | True if version > deployedVersion |
| `status` | `active`, `deprecated`, `disabled` |
| `changelog` | Version history for tracking changes |

---

## 6. Migration System

### hooks/_migrations.json

Track all deployments for audit trail:

```json
{
  "migrations": [
    {
      "id": "mig-20250310-001",
      "timestamp": "2025-03-10T14:30:00Z",
      "type": "deploy",
      "hooks": [
        { "name": "voice-analysis", "from": "1.1.0", "to": "1.2.0" },
        { "name": "generate-captions", "from": "1.5.0", "to": "2.0.0" }
      ],
      "deployedBy": "otto",
      "notes": "Routine update with error handling improvements"
    },
    {
      "id": "mig-20250308-001",
      "timestamp": "2025-03-08T16:00:00Z",
      "type": "deploy",
      "hooks": [
        { "name": "generate-captions", "from": "1.0.0", "to": "1.5.0" }
      ],
      "deployedBy": "otto",
      "notes": "Added TikTok caption support"
    },
    {
      "id": "mig-20250306-001",
      "timestamp": "2025-03-06T10:00:00Z",
      "type": "initial",
      "hooks": [
        { "name": "voice-analysis", "from": null, "to": "1.0.0" },
        { "name": "guest-research", "from": null, "to": "1.0.0" },
        { "name": "generate-captions", "from": null, "to": "1.0.0" }
      ],
      "deployedBy": "otto",
      "notes": "Initial deployment of all hooks"
    }
  ]
}
```

---

## 7. Deployment Commands

When you're ready to deploy, just tell Otto what to do. Here are the commands Otto understands:

### Deploy All Pending Updates

```
You: "Read my hooks/_index.json and deploy any hooks that need updating"

Otto: [Reads the index, identifies hooks where version > deployedVersion,
       reads each hook.js file, deploys to platform, updates index]
```

### Deploy Specific Hook

```
You: "Deploy the voice-analysis hook from my repo"

Otto: [Reads hooks/voice-analysis/hook.js, deploys/updates on platform]
```

### Check Deployment Status

```
You: "Compare my hooks/_index.json with what's deployed and tell me what's different"

Otto: [Reads index, queries platform for current hooks, shows diff]
```

### Rollback

```
You: "Rollback voice-analysis to version 1.1.0"

Otto: [Checks if 1.1.0 exists in changelog, asks for confirmation,
       deploys that version, updates index]
```

### List All Deployed Hooks

```
You: "Show me all server functions currently deployed"

Otto: [Lists all hooks with names, versions, status]
```

### Delete Deprecated Hook

```
You: "Remove the old-function hook that's marked as deprecated"

Otto: [Confirms, deletes from platform, updates index]
```

---

## 8. Go Daemon API Contract Template

Define your daemon's API contract so Otto knows what endpoints exist:

### daemon/api-contract.yaml (OpenAPI 3.0)

```yaml
openapi: 3.0.0
info:
  title: Throughline Daemon API
  version: 1.0.0
  description: Backend business logic for Throughline podcast tools

servers:
  - url: https://api.throughline.example.com/api/v1
    description: Production

paths:
  /analyze-voice:
    post:
      summary: Analyze voice characteristics from audio
      operationId: analyzeVoice
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - profileId
                - audioUrl
              properties:
                profileId:
                  type: string
                  description: Voice profile ID
                audioUrl:
                  type: string
                  format: uri
                  description: URL to audio file
                workspaceId:
                  type: string
                  description: Audos workspace ID
      responses:
        '200':
          description: Analysis complete
          content:
            application/json:
              schema:
                type: object
                properties:
                  tone:
                    type: string
                    enum: [formal, casual, enthusiastic, calm, authoritative]
                  pace:
                    type: string
                    enum: [slow, moderate, fast]
                  keywords:
                    type: array
                    items:
                      type: string
                  confidence:
                    type: number
                    minimum: 0
                    maximum: 1
        '400':
          description: Invalid input
        '503':
          description: Service unavailable

  /research-guest:
    post:
      summary: Research a podcast guest using AI
      operationId: researchGuest
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - guestName
              properties:
                guestName:
                  type: string
                context:
                  type: string
                  description: Podcast context/topic
      responses:
        '200':
          description: Research complete
          content:
            application/json:
              schema:
                type: object
                properties:
                  bio:
                    type: string
                  topics:
                    type: array
                    items:
                      type: string
                  recentWork:
                    type: string
                  questions:
                    type: array
                    items:
                      type: string
                  conversationStarters:
                    type: array
                    items:
                      type: string

  /generate-captions:
    post:
      summary: Generate social media captions for content
      operationId: generateCaptions
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - transcript
                - platform
              properties:
                transcript:
                  type: string
                platform:
                  type: string
                  enum: [instagram, linkedin, twitter, tiktok]
                topic:
                  type: string
      responses:
        '200':
          description: Caption generated
          content:
            application/json:
              schema:
                type: object
                properties:
                  caption:
                    type: string
                  hashtags:
                    type: array
                    items:
                      type: string
                  characterCount:
                    type: integer
                  platformMaxLength:
                    type: integer

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer

security:
  - bearerAuth: []
```

---

## 9. Sample Wrapper Hook Templates

### Template 1: Simple Passthrough

```javascript
// hooks/simple-passthrough/hook.js
/**
 * @name simple-passthrough
 * @version 1.0.0
 * @description Simple passthrough to daemon endpoint
 */

const DAEMON_URL = 'https://api.throughline.example.com';
const API_KEY = 'your-api-key';

const { action, ...payload } = request.body;

if (!action) {
  return respond(400, { error: 'action is required' });
}

try {
  const response = await fetch(`${DAEMON_URL}/api/v1/${action}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${API_KEY}`
    },
    body: JSON.stringify(payload)
  });

  const result = await response.json();
  return respond(response.status, result);
} catch (err) {
  return respond(503, { error: 'Service unavailable' });
}
```

### Template 2: With Database Logging

```javascript
// hooks/with-db-logging/hook.js
/**
 * @name with-db-logging
 * @version 1.0.0
 * @description Calls daemon and logs results to database
 */

const DAEMON_URL = 'https://api.throughline.example.com';
const API_KEY = 'your-api-key';

const { endpoint, payload, logTable } = request.body;

if (!endpoint || !payload) {
  return respond(400, { error: 'endpoint and payload are required' });
}

const startTime = Date.now();

try {
  const response = await fetch(`${DAEMON_URL}${endpoint}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${API_KEY}`
    },
    body: JSON.stringify(payload)
  });

  const result = await response.json();
  const duration = Date.now() - startTime;

  // Log to database if table specified
  if (logTable) {
    await db.insert(logTable, {
      endpoint,
      request_payload: JSON.stringify(payload),
      response_payload: JSON.stringify(result),
      status_code: response.status,
      duration_ms: duration,
      created_at: new Date().toISOString()
    });
  }

  return respond(response.status, result);
} catch (err) {
  console.error('Daemon error:', err);
  return respond(503, { error: 'Service unavailable', details: err.message });
}
```

### Template 3: With Retry Logic

```javascript
// hooks/with-retry/hook.js
/**
 * @name with-retry
 * @version 1.0.0
 * @description Calls daemon with automatic retry on failure
 */

const DAEMON_URL = 'https://api.throughline.example.com';
const API_KEY = 'your-api-key';
const MAX_RETRIES = 3;
const RETRY_DELAY_MS = 1000;

const { endpoint, payload } = request.body;

async function callWithRetry(url, options, retries = 0) {
  try {
    const response = await fetch(url, options);
    if (response.ok) {
      return response;
    }
    // Retry on 5xx errors
    if (response.status >= 500 && retries < MAX_RETRIES) {
      console.log(`Retry ${retries + 1}/${MAX_RETRIES} after ${response.status}`);
      await new Promise(resolve => setTimeout(resolve, RETRY_DELAY_MS * (retries + 1)));
      return callWithRetry(url, options, retries + 1);
    }
    return response;
  } catch (err) {
    if (retries < MAX_RETRIES) {
      console.log(`Retry ${retries + 1}/${MAX_RETRIES} after error: ${err.message}`);
      await new Promise(resolve => setTimeout(resolve, RETRY_DELAY_MS * (retries + 1)));
      return callWithRetry(url, options, retries + 1);
    }
    throw err;
  }
}

try {
  const response = await callWithRetry(`${DAEMON_URL}${endpoint}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${API_KEY}`
    },
    body: JSON.stringify(payload)
  });

  const result = await response.json();
  return respond(response.status, result);
} catch (err) {
  return respond(503, { error: 'Service unavailable after retries' });
}
```

---

## Quick Reference

### Deployment Workflow

```
1. Write/update hook code in hooks/{name}/hook.js
2. Update version in hooks/_index.json
3. Commit to GitHub
4. Tell Otto: "Deploy hooks from my repo"
5. Otto reads _index.json, deploys changed hooks
6. Otto updates deployedVersion and _migrations.json
```

### Naming Rules

| Type | Format | Example |
|------|--------|---------|
| Hook name | lowercase-with-hyphens | `voice-analysis` |
| Hook file | `hooks/{name}/hook.js` | `hooks/voice-analysis/hook.js` |
| Version | semver | `1.2.0` |
| Daemon endpoint | /api/v1/{action} | `/api/v1/analyze-voice` |

### Commands Otto Understands

| Command | What It Does |
|---------|--------------|
| "Deploy hooks from my repo" | Syncs all hooks where version > deployedVersion |
| "Deploy {name} hook" | Deploys specific hook |
| "Show hook deployment status" | Compares repo vs deployed |
| "Rollback {name} to {version}" | Reverts to previous version |
| "Delete {name} hook" | Removes hook from platform |
| "List all hooks" | Shows deployed hooks |
| "Test {name} hook with {payload}" | Tests a hook |

---

## Summary

| Question | Answer |
|----------|--------|
| Can I name my own hooks? | **Yes** - you specify the name in `_index.json` |
| Can I version hooks? | **Yes** - use semver in your index file |
| Can hooks call my daemon? | **Yes** - use `fetch()` to call any URL |
| How do I deploy? | Tell Otto to read your `_index.json` and sync |
| Is there migration tracking? | **Yes** - `_migrations.json` tracks all deployments |
| Can I rollback? | **Yes** - Otto can deploy any previous version |

---

*Document generated by Otto (Audos AI Assistant)*
*Last Updated: March 6, 2025*
