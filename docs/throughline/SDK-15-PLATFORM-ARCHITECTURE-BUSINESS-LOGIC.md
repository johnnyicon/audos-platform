# Audos Platform Architecture: Where Business Logic Lives

**Date:** March 6, 2025
**Platform:** Audos Runtime
**Document Type:** Architecture Reference

---

## Executive Summary

The Audos platform uses a **three-layer architecture** for separating concerns:

1. **Frontend (React Components)** — UI rendering and user interactions
2. **Server Functions (Hooks)** — Custom backend logic you write
3. **Platform APIs** — Built-in endpoints for CRM, analytics, payments, etc.

This document explains where different types of business logic should live, with examples specific to the Throughline podcast creator app.

---

## Architecture Overview

```
┌──────────────────────────────────────────────────────────────────────────┐
│                           USER'S BROWSER                                 │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │                    FRONTEND LAYER                                  │  │
│  │                    (React Components)                              │  │
│  │                                                                    │  │
│  │  • UI rendering (JSX/TSX)                                         │  │
│  │  • User interactions (onClick, onChange, etc.)                    │  │
│  │  • Local state management (useState, useReducer)                  │  │
│  │  • HTTP calls to backend APIs                                     │  │
│  │                                                                    │  │
│  │  ⚠️  NO secrets, NO sensitive business logic                      │  │
│  │  ⚠️  Code is visible in browser DevTools                          │  │
│  └────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ HTTP Requests (fetch/axios)
                                    ▼
┌──────────────────────────────────────────────────────────────────────────┐
│                           AUDOS SERVER                                   │
│                                                                          │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │                 SERVER FUNCTIONS LAYER                             │  │
│  │                 (Hooks - Custom Backend Code)                      │  │
│  │                                                                    │  │
│  │  Endpoint: /api/hooks/execute/workspace-{id}/{hookName}           │  │
│  │                                                                    │  │
│  │  Available in your code:                                          │  │
│  │  • request.body     — POST body data                              │  │
│  │  • request.query    — URL query parameters                        │  │
│  │  • request.method   — HTTP method (GET, POST, etc.)               │  │
│  │  • db.query()       — Read from database                          │  │
│  │  • db.insert()      — Insert into database                        │  │
│  │  • db.update()      — Update database records                     │  │
│  │  • db.delete()      — Delete from database                        │  │
│  │  • platform.generateText() — AI text generation                   │  │
│  │  • platform.sendEmail()    — Send emails                          │  │
│  │  • fetch()          — Call external APIs                          │  │
│  │  • respond(status, body) — Return response                        │  │
│  │                                                                    │  │
│  │  ✅ Secrets stay server-side                                      │  │
│  │  ✅ Business logic is protected                                   │  │
│  │  ✅ Validation cannot be bypassed                                 │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                                                                          │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │                 PLATFORM APIs LAYER                                │  │
│  │                 (Built-in Endpoints)                               │  │
│  │                                                                    │  │
│  │  CRM & Contacts:                                                  │  │
│  │  • /api/crm/contacts/{workspaceId}                                │  │
│  │  • /api/crm/contacts/{workspaceId}/{contactId}                    │  │
│  │                                                                    │  │
│  │  Workspace Database:                                              │  │
│  │  • /api/workspace-db/{workspaceId}/tables                         │  │
│  │  • /api/workspace-db/{workspaceId}/query                          │  │
│  │  • /api/workspace-db/{workspaceId}/insert                         │  │
│  │                                                                    │  │
│  │  Analytics:                                                       │  │
│  │  • /api/funnel/events/{workspaceId}                               │  │
│  │  • /api/analytics/{workspaceId}                                   │  │
│  │                                                                    │  │
│  │  Payments (Stripe):                                               │  │
│  │  • /api/stripe/checkout                                           │  │
│  │  • /api/stripe/subscriptions                                      │  │
│  │                                                                    │  │
│  │  ✅ Already secured by platform                                   │  │
│  │  ✅ Workspace-scoped automatically                                │  │
│  └────────────────────────────────────────────────────────────────────┘  │
│                                                                          │
│  ┌────────────────────────────────────────────────────────────────────┐  │
│  │                 DATABASE LAYER                                     │  │
│  │                 (PostgreSQL)                                       │  │
│  │                                                                    │  │
│  │  System Tables (managed by platform):                             │  │
│  │  • funnel_contacts    — CRM contacts                              │  │
│  │  • funnel_events      — Analytics events                          │  │
│  │  • ad_campaigns       — Ad campaign data                          │  │
│  │                                                                    │  │
│  │  Workspace Tables (created by you):                               │  │
│  │  • voice_profiles     — Voice training data                       │  │
│  │  • podcast_profiles   — Podcast configurations                    │  │
│  │  • speakers           — Guest/host profiles                       │  │
│  │  • reels              — Content clips                             │  │
│  │  • captions           — Generated captions                        │  │
│  └────────────────────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
```

---

## Where Different Logic Types Should Live

| Logic Type | Layer | Why |
|------------|-------|-----|
| UI rendering (buttons, forms, layouts) | Frontend | It's presentation |
| Form validation (required fields, formats) | Frontend | Immediate UX feedback |
| Form validation (security, authorization) | Server Functions | Can't trust client-side validation |
| Database reads (simple queries) | Platform APIs | Built-in, secure, workspace-scoped |
| Database writes (with business rules) | Server Functions | Enforce rules server-side |
| API key usage (OpenAI, external services) | Server Functions | Never expose keys in frontend |
| Complex business logic | Server Functions | Keep algorithms protected |
| AI text generation | Server Functions | Uses platform.generateText() |
| Email sending | Server Functions or Platform APIs | Server-side only |
| Payment processing | Platform APIs (Stripe) | Already PCI compliant |
| File uploads | Platform APIs | Handled by platform |
| User authentication | Platform (Email Gate) | Managed by platform |

---

## Example: Throughline App Architecture

### Use Case 1: Saving a Voice Profile

**Frontend (apps/throughline/App.tsx):**
```tsx
// User fills out voice profile form
const [voiceData, setVoiceData] = useState({
  name: '',
  description: '',
  traits: [],
  sampleText: ''
});

const handleSaveVoiceProfile = async () => {
  // Basic validation (UX only - server will re-validate)
  if (!voiceData.name) {
    setError('Name is required');
    return;
  }

  try {
    // Call YOUR server function (hook)
    const response = await fetch(
      '/api/hooks/execute/workspace-8f1ad824-832f-4af8-b77e-ab931a250625/save-voice-profile',
      {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(voiceData)
      }
    );

    const result = await response.json();

    if (result.success) {
      // Update UI
      setProfiles([...profiles, result.profile]);
      setVoiceData({ name: '', description: '', traits: [], sampleText: '' });
    } else {
      setError(result.error);
    }
  } catch (err) {
    setError('Failed to save profile');
  }
};
```

**Server Function (Hook: save-voice-profile):**
```javascript
// This runs on the server - code is NOT visible to users

// Validate the request
const { name, description, traits, sampleText } = request.body;

if (!name || name.length < 2) {
  return respond(400, { success: false, error: 'Name must be at least 2 characters' });
}

if (!description || description.length < 10) {
  return respond(400, { success: false, error: 'Description must be at least 10 characters' });
}

// Check for duplicate names (business rule)
const existing = await db.query('voice_profiles', {
  filters: [{ column: 'name', operator: 'eq', value: name }]
});

if (existing.rows.length > 0) {
  return respond(400, { success: false, error: 'A profile with this name already exists' });
}

// Insert into database
const result = await db.insert('voice_profiles', {
  name,
  description,
  traits: JSON.stringify(traits),
  sample_text: sampleText,
  created_at: new Date().toISOString()
});

return respond(200, {
  success: true,
  profile: {
    id: result.id,
    name,
    description,
    traits,
    sampleText
  }
});
```

---

### Use Case 2: AI-Powered Guest Research

**Frontend:**
```tsx
const handleResearchGuest = async (guestName: string, context: string) => {
  setLoading(true);

  const response = await fetch(
    '/api/hooks/execute/workspace-8f1ad824.../research-guest',
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ guestName, context })
    }
  );

  const result = await response.json();
  setGuestResearch(result.research);
  setLoading(false);
};
```

**Server Function (Hook: research-guest):**
```javascript
const { guestName, context } = request.body;

if (!guestName) {
  return respond(400, { error: 'Guest name is required' });
}

// Use AI to generate research (this uses YOUR platform's AI credits)
const research = await platform.generateText({
  prompt: `Research the following podcast guest and provide:
1. Brief bio (2-3 sentences)
2. Key topics they're known for
3. Recent work or projects
4. 5 potential interview questions
5. Conversation starters based on their interests

Guest: ${guestName}
Podcast context: ${context || 'General interest podcast'}

Format the response as JSON with keys: bio, topics, recentWork, questions, conversationStarters`,
  maxTokens: 1000
});

// Parse the AI response
let parsed;
try {
  parsed = JSON.parse(research);
} catch {
  parsed = { raw: research };
}

// Save to database for future reference
await db.insert('guest_research', {
  guest_name: guestName,
  context,
  research: JSON.stringify(parsed),
  created_at: new Date().toISOString()
});

return respond(200, { research: parsed });
```

---

### Use Case 3: Generating Captions with External API

**Frontend:**
```tsx
const handleGenerateCaptions = async (reelId: string, platform: string) => {
  const response = await fetch(
    '/api/hooks/execute/workspace-8f1ad824.../generate-captions',
    {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ reelId, platform })
    }
  );

  const result = await response.json();
  setCaptions(result.captions);
};
```

**Server Function (Hook: generate-captions):**
```javascript
const { reelId, platform } = request.body;

// Fetch the reel from database
const reelResult = await db.query('reels', {
  filters: [{ column: 'id', operator: 'eq', value: reelId }]
});

if (reelResult.rows.length === 0) {
  return respond(404, { error: 'Reel not found' });
}

const reel = reelResult.rows[0];

// Platform-specific formatting rules (business logic!)
const platformRules = {
  instagram: {
    maxLength: 2200,
    hashtagStyle: 'inline',
    emojiDensity: 'high'
  },
  linkedin: {
    maxLength: 3000,
    hashtagStyle: 'bottom',
    emojiDensity: 'low'
  },
  twitter: {
    maxLength: 280,
    hashtagStyle: 'inline',
    emojiDensity: 'medium'
  }
};

const rules = platformRules[platform] || platformRules.instagram;

// Generate caption using AI
const caption = await platform.generateText({
  prompt: `Generate a ${platform} caption for this podcast clip:

Transcript: ${reel.transcript}
Topic: ${reel.topic}

Rules:
- Maximum ${rules.maxLength} characters
- Hashtag style: ${rules.hashtagStyle}
- Emoji usage: ${rules.emojiDensity}
- Include a call-to-action
- Match the tone of the transcript`,
  maxTokens: 500
});

// Save caption to database
await db.insert('captions', {
  reel_id: reelId,
  platform,
  content: caption,
  created_at: new Date().toISOString()
});

return respond(200, {
  captions: {
    platform,
    content: caption,
    characterCount: caption.length,
    maxLength: rules.maxLength
  }
});
```

---

## Security Comparison

### ❌ WRONG: Business Logic in Frontend

```tsx
// apps/throughline/App.tsx - INSECURE!

const OPENAI_API_KEY = 'sk-...'; // ❌ Exposed in browser!

const generateCaption = async (transcript: string) => {
  // ❌ Anyone can see this API call in DevTools
  const response = await fetch('https://api.openai.com/v1/completions', {
    headers: {
      'Authorization': `Bearer ${OPENAI_API_KEY}` // ❌ Key visible!
    },
    body: JSON.stringify({
      model: 'gpt-4',
      prompt: transcript
    })
  });
  // ...
};

// ❌ Business rules can be bypassed
const canDeleteProfile = (profile) => {
  return profile.createdBy === currentUser.id; // User can modify this in console!
};
```

### ✅ CORRECT: Business Logic in Server Functions

```tsx
// apps/throughline/App.tsx - SECURE

const generateCaption = async (transcript: string) => {
  // ✅ Calls YOUR server function - no keys exposed
  const response = await fetch(
    '/api/hooks/execute/workspace-8f1ad824.../generate-caption',
    {
      method: 'POST',
      body: JSON.stringify({ transcript })
    }
  );
  return response.json();
};

const deleteProfile = async (profileId: string) => {
  // ✅ Server validates ownership - can't be bypassed
  const response = await fetch(
    '/api/hooks/execute/workspace-8f1ad824.../delete-profile',
    {
      method: 'POST',
      body: JSON.stringify({ profileId })
    }
  );
  return response.json();
};
```

```javascript
// Server Function: delete-profile - SECURE

const { profileId } = request.body;

// Get profile from database
const profile = await db.query('voice_profiles', {
  filters: [{ column: 'id', operator: 'eq', value: profileId }]
});

if (profile.rows.length === 0) {
  return respond(404, { error: 'Profile not found' });
}

// ✅ Server-side authorization check - CANNOT be bypassed
if (profile.rows[0].user_id !== request.session.userId) {
  return respond(403, { error: 'Not authorized to delete this profile' });
}

// Delete the profile
await db.delete('voice_profiles', {
  filters: [{ column: 'id', operator: 'eq', value: profileId }]
});

return respond(200, { success: true });
```

---

## How to Create Server Functions

### Using Otto (This AI Assistant)

```
You: "Create a server function called 'save-voice-profile' that validates
     and saves voice profile data to the database"

Otto: [Uses manage_server_functions tool to create the hook]
```

### Server Function Template

```javascript
// Hook Name: my-function-name
// Available globals: request, db, platform, fetch, respond, JSON, Date, Math, console

// Get request data
const { param1, param2 } = request.body;  // POST body
const { queryParam } = request.query;      // URL query params
const method = request.method;             // GET, POST, etc.

// Validate
if (!param1) {
  return respond(400, { error: 'param1 is required' });
}

// Database operations
const results = await db.query('table_name', {
  filters: [{ column: 'field', operator: 'eq', value: 'something' }],
  limit: 10
});

await db.insert('table_name', { field1: 'value1', field2: 'value2' });

await db.update('table_name',
  { filters: [{ column: 'id', operator: 'eq', value: 123 }] },
  { field1: 'new_value' }
);

await db.delete('table_name', {
  filters: [{ column: 'id', operator: 'eq', value: 123 }]
});

// AI generation
const text = await platform.generateText({
  prompt: 'Your prompt here',
  maxTokens: 500
});

// Email sending
await platform.sendEmail({
  to: 'user@example.com',
  subject: 'Hello',
  text: 'Plain text body',
  html: '<p>HTML body</p>'
});

// External API calls
const externalData = await fetch('https://api.example.com/data', {
  method: 'GET',
  headers: { 'Authorization': 'Bearer token' }
}).then(r => r.json());

// Return response
return respond(200, { success: true, data: results });
```

---

## Testing Server Functions

After creating a server function, test it using the `test_server_function` tool:

```
You: "Test the save-voice-profile hook with sample data"

Otto: [Uses test_server_function tool to execute and show results]
```

This returns:
- Status code
- Response body
- Console logs
- Errors (if any)
- Execution time

---

## Summary: The Decision Tree

```
Is this logic about how things LOOK?
  └── YES → Frontend (React component)
  └── NO ↓

Does it need API keys or secrets?
  └── YES → Server Function (Hook)
  └── NO ↓

Is it a security-sensitive operation (delete, update, payment)?
  └── YES → Server Function (Hook)
  └── NO ↓

Is it a simple CRUD operation?
  └── YES → Platform API (workspace-db, crm, etc.)
  └── NO ↓

Is it complex business logic or AI processing?
  └── YES → Server Function (Hook)
  └── NO → Frontend (but consider Server Function anyway)
```

---

## Related Documents

- [Incident Report: Desktop App Loading](./INCIDENT-REPORT-DESKTOP-APP-LOADING.md)
- [Incident Report: Dock Removal](./INCIDENT-REPORT-DOCK-REMOVAL.md)
- [Audos Database SDK](./AUDOS-DATABASE-SDK.md)

---

*Document generated by Otto (Audos AI Assistant)*
*Last Updated: March 6, 2025*
