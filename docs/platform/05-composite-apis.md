# Audos Composite/Aggregate APIs Guide

## What Are Composite APIs?

### Definition

**Composite APIs** (also called Aggregate APIs) are server functions that combine multiple Audos primitives into a single, cohesive operation. They orchestrate several lower-level operations to accomplish a higher-level business goal.

**Pattern:**
```
Primitive 1 + Primitive 2 + Business Logic = Composite API
```

For example:
```
fetch() + platform.generateText() + db.insert() = Guest Research API
```

### Why Use Composite APIs?

1. **Single Atomic Operation**: Multiple steps execute together on the server, reducing round-trips
2. **Reduced Latency**: No network overhead between steps—everything runs server-side
3. **Reusable Across Apps**: Multiple mini-apps can call the same composite API
4. **Simplified Client Code**: Apps call one endpoint instead of orchestrating multiple calls
5. **Consistent Business Logic**: Logic lives in one place, not duplicated across apps
6. **Better Error Handling**: Server can handle failures and retry logic internally

### When NOT to Use Composite APIs

- **One-off scripts**: If you're running a script once locally, orchestrate there
- **Rapid prototyping**: Build locally first, extract to composite API when pattern solidifies
- **UI-driven workflows**: When the user needs to see/approve intermediate steps

---

## Available Primitives to Combine

Composite APIs orchestrate these building blocks:

### Database Operations
```javascript
// Query data
const rows = await db.query('users', {
  filters: [{ column: 'status', operator: 'eq', value: 'active' }],
  limit: 50
});

// Insert data
await db.insert('episodes', [
  { title: 'Episode 1', guest_name: 'Alice', recorded_at: new Date().toISOString() }
]);

// Update data
await db.update('episodes',
  [{ column: 'id', operator: 'eq', value: episodeId }],
  { published: true, published_at: new Date().toISOString() }
);

// Delete data
await db.delete('episodes',
  [{ column: 'id', operator: 'eq', value: episodeId }]
);
```

### AI Text Generation
```javascript
// Generate text with GPT-4o-mini
const result = await platform.generateText({
  prompt: 'Summarize this transcript in 3 bullet points:\n\n' + transcript,
  maxTokens: 500
});

console.log(result.text); // Generated content
```

### Transactional Email
```javascript
// Send email
await platform.sendEmail({
  to: 'user@example.com',
  subject: 'Your Daily Digest',
  text: 'Plain text version',
  html: '<h1>HTML version</h1>'
});
```

### HTTP Requests
```javascript
// External API
const response = await fetch('https://api.example.com/data');
const data = await response.json();

// Internal Audos API
const contactsResponse = await fetch(
  `https://audos.com/api/crm/contacts/workspace-${request.workspaceId}`,
  {
    headers: {
      'x-workspace-id': request.workspaceId,
      'x-api-key': request.apiKey
    }
  }
);
const contacts = await contactsResponse.json();
```

### Internal APIs Available
- `/api/crm/contacts/{workspaceId}` - CRM contacts
- `/api/funnel/metrics/{workspaceId}` - Funnel analytics
- `/api/funnel/events/{workspaceId}` - Event stream
- `/api/spaces/{spaceId}/sessions` - Session data
- `/api/workspaces/{workspaceId}/analytics` - Workspace analytics

---

## Common Patterns

### 1. Fetch-Transform-Store
Fetch external content → AI processes it → Store in database

**Use Cases:**
- Research automation
- Content ingestion
- Data enrichment

**Flow:**
```
fetch(url) → platform.generateText() → db.insert()
```

---

### 2. Query-Enrich-Respond
Query database → AI enhances results → Return enriched response

**Use Cases:**
- Smart search
- Recommendation engines
- Personalized responses

**Flow:**
```
db.query() → platform.generateText() → respond()
```

---

### 3. Multi-Source Aggregation
Fetch from multiple sources → Combine data → Return unified response

**Use Cases:**
- Dashboard APIs
- Report generation
- Cross-system queries

**Flow:**
```
fetch(source1) + fetch(source2) + db.query() → combine → respond()
```

---

### 4. Scheduled Processing
Cron trigger → Query/Fetch → Process → Store/Email

**Use Cases:**
- Daily digests
- Automated reporting
- Background sync jobs

**Flow:**
```
cron trigger → db.query() → platform.generateText() → platform.sendEmail()
```

---

## Complete Examples

### Example 1: Guest Research API

**Purpose**: Automatically research podcast guests by fetching content from provided URLs and generating AI summaries.

**Pattern**: Fetch-Transform-Store

**Code**:
```javascript
// Hook name: guest-research-api
// Endpoint: POST /api/hooks/execute/workspace-{configId}/guest-research-api

const { guestName, urls, episodeId } = request.body;

if (!guestName || !urls || !Array.isArray(urls)) {
  return respond(400, { error: 'Missing required fields: guestName, urls (array)' });
}

console.log(`[Guest Research] Starting research for ${guestName}`);
console.log(`[Guest Research] URLs to fetch: ${urls.length}`);

// Step 1: Fetch content from all URLs
const fetchedContent = [];
for (const url of urls) {
  try {
    console.log(`[Guest Research] Fetching ${url}`);
    const response = await fetch(url);
    const html = await response.text();

    // Extract text content (basic HTML stripping)
    const text = html
      .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
      .replace(/<style\b[^<]*(?:(?!<\/style>)<[^<]*)*<\/style>/gi, '')
      .replace(/<[^>]+>/g, ' ')
      .replace(/\s+/g, ' ')
      .trim()
      .slice(0, 5000); // Limit to 5000 chars per URL

    fetchedContent.push({
      url,
      content: text,
      success: true
    });
  } catch (error) {
    console.error(`[Guest Research] Failed to fetch ${url}:`, error.message);
    fetchedContent.push({
      url,
      error: error.message,
      success: false
    });
  }
}

// Step 2: Combine all successful content
const combinedContent = fetchedContent
  .filter(item => item.success)
  .map(item => `Source: ${item.url}\n\n${item.content}`)
  .join('\n\n---\n\n');

if (!combinedContent) {
  return respond(500, { error: 'Failed to fetch any content from provided URLs' });
}

// Step 3: Generate AI summary
console.log('[Guest Research] Generating AI summary');
const summaryPrompt = `You are researching a podcast guest named "${guestName}".

Based on the following web content, create a comprehensive research summary:

${combinedContent}

Generate a structured summary with:
1. **Background** (2-3 sentences about who they are)
2. **Notable Work** (bullet points of key achievements/projects)
3. **Expertise** (areas they're known for)
4. **Interesting Facts** (unique angles for podcast conversation)
5. **Suggested Questions** (3-5 questions to ask)

Keep the tone professional but conversational.`;

const summaryResult = await platform.generateText({
  prompt: summaryPrompt,
  maxTokens: 1000
});

// Step 4: Store research in database
console.log('[Guest Research] Saving to database');

// Ensure guest_research table exists (you'd create this with db_create_table)
// Columns: guest_name (text), summary (text), sources (json), episode_id (integer), created_at (timestamp)

await db.insert('guest_research', [{
  guest_name: guestName,
  summary: summaryResult.text,
  sources: JSON.stringify(fetchedContent),
  episode_id: episodeId || null,
  created_at: new Date().toISOString()
}]);

console.log('[Guest Research] Research complete');

return respond(200, {
  success: true,
  guestName,
  summary: summaryResult.text,
  sourcesProcessed: fetchedContent.filter(f => f.success).length,
  sourcesFailed: fetchedContent.filter(f => !f.success).length
});
```

**Client Usage**:
```typescript
// From a mini-app
const result = await fetch('/api/hooks/execute/workspace-{configId}/guest-research-api', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'x-workspace-id': workspaceId,
    'x-api-key': apiKey
  },
  body: JSON.stringify({
    guestName: 'Jane Doe',
    urls: [
      'https://janedoe.com/about',
      'https://twitter.com/janedoe',
      'https://medium.com/@janedoe/latest-article'
    ],
    episodeId: 42
  })
});

const data = await result.json();
console.log(data.summary); // AI-generated research
```

---

### Example 2: Episode Processor API

**Purpose**: Process podcast transcripts to generate show notes, timestamps, and suggested clips.

**Pattern**: Query-Enrich-Respond

**Code**:
```javascript
// Hook name: episode-processor-api
// Endpoint: POST /api/hooks/execute/workspace-{configId}/episode-processor-api

const { episodeId } = request.body;

if (!episodeId) {
  return respond(400, { error: 'Missing required field: episodeId' });
}

console.log(`[Episode Processor] Processing episode ${episodeId}`);

// Step 1: Query episode and transcript
const episodes = await db.query('episodes', {
  filters: [{ column: 'id', operator: 'eq', value: episodeId }],
  limit: 1
});

if (episodes.length === 0) {
  return respond(404, { error: 'Episode not found' });
}

const episode = episodes[0];
const transcript = episode.transcript;

if (!transcript) {
  return respond(400, { error: 'Episode has no transcript' });
}

console.log(`[Episode Processor] Found transcript (${transcript.length} chars)`);

// Step 2: Generate show notes
console.log('[Episode Processor] Generating show notes');
const showNotesPrompt = `You are a podcast editor. Generate engaging show notes for this episode.

Episode Title: ${episode.title}
Guest: ${episode.guest_name || 'N/A'}
Transcript:
${transcript}

Create show notes with:
1. **Episode Summary** (2-3 sentences)
2. **Key Topics** (bullet points)
3. **Timestamps** (major topic shifts, format as MM:SS - Topic)
4. **Quotes** (2-3 memorable quotes)
5. **Resources Mentioned** (links, books, tools mentioned)

Format in Markdown.`;

const showNotesResult = await platform.generateText({
  prompt: showNotesPrompt,
  maxTokens: 1500
});

// Step 3: Generate suggested clips
console.log('[Episode Processor] Identifying clip-worthy moments');
const clipsPrompt = `You are identifying viral-worthy clips from this podcast transcript.

Transcript:
${transcript.slice(0, 10000)} // Limit for token efficiency

Identify 3-5 short segments (30-90 seconds each) that would make great social media clips.

For each clip, provide:
- **Start Time** (estimate in MM:SS format)
- **Duration** (in seconds)
- **Hook** (why it's compelling)
- **Suggested Title** (catchy, under 60 chars)

Format as JSON array:
[
  {
    "startTime": "12:34",
    "duration": 45,
    "hook": "...",
    "title": "..."
  }
]`;

const clipsResult = await platform.generateText({
  prompt: clipsPrompt,
  maxTokens: 800
});

// Parse clips (handle potential JSON formatting issues)
let suggestedClips = [];
try {
  const jsonMatch = clipsResult.text.match(/\[[\s\S]*\]/);
  if (jsonMatch) {
    suggestedClips = JSON.parse(jsonMatch[0]);
  }
} catch (parseError) {
  console.error('[Episode Processor] Failed to parse clips JSON:', parseError.message);
  // Continue without clips
}

// Step 4: Update episode with generated content
console.log('[Episode Processor] Updating episode');
await db.update('episodes',
  [{ column: 'id', operator: 'eq', value: episodeId }],
  {
    show_notes: showNotesResult.text,
    suggested_clips: JSON.stringify(suggestedClips),
    processed_at: new Date().toISOString()
  }
);

console.log('[Episode Processor] Processing complete');

return respond(200, {
  success: true,
  episodeId,
  showNotes: showNotesResult.text,
  suggestedClips,
  clipCount: suggestedClips.length
});
```

**Client Usage**:
```typescript
// Trigger from mini-app after uploading transcript
const result = await fetch('/api/hooks/execute/workspace-{configId}/episode-processor-api', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'x-workspace-id': workspaceId,
    'x-api-key': apiKey
  },
  body: JSON.stringify({
    episodeId: 42
  })
});

const data = await result.json();
console.log(data.showNotes); // Generated show notes
console.log(data.suggestedClips); // Clip recommendations
```

---

### Example 3: Daily Digest API

**Purpose**: Scheduled function that queries recent activity, formats with AI, and emails a daily digest.

**Pattern**: Scheduled Processing

**Code**:
```javascript
// Hook name: daily-digest-api
// Endpoint: POST /api/hooks/execute/workspace-{configId}/daily-digest-api
// Scheduled: Daily at 9 AM (configure with booster or external cron)

console.log('[Daily Digest] Starting daily digest generation');

// Step 1: Query recent episodes (last 7 days)
const sevenDaysAgo = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString();

const recentEpisodes = await db.query('episodes', {
  filters: [
    { column: 'created_at', operator: 'gte', value: sevenDaysAgo }
  ],
  orderBy: { column: 'created_at', direction: 'desc' },
  limit: 10
});

console.log(`[Daily Digest] Found ${recentEpisodes.length} recent episodes`);

// Step 2: Query recent guests
const recentGuests = await db.query('guest_research', {
  filters: [
    { column: 'created_at', operator: 'gte', value: sevenDaysAgo }
  ],
  orderBy: { column: 'created_at', direction: 'desc' },
  limit: 5
});

console.log(`[Daily Digest] Found ${recentGuests.length} recent guest research entries`);

// Step 3: Fetch analytics from internal API (if available)
let analytics = null;
try {
  const analyticsResponse = await fetch(
    `https://audos.com/api/funnel/metrics/workspace-${request.workspaceId}`,
    {
      headers: {
        'x-workspace-id': request.workspaceId,
        'x-api-key': request.apiKey
      }
    }
  );
  analytics = await analyticsResponse.json();
} catch (error) {
  console.error('[Daily Digest] Failed to fetch analytics:', error.message);
}

// Step 4: Generate digest with AI
console.log('[Daily Digest] Generating AI digest');

const digestPrompt = `You are creating a daily digest email for a podcast production workspace.

**Recent Episodes (last 7 days):**
${recentEpisodes.map(ep => `- ${ep.title} (Guest: ${ep.guest_name || 'TBD'})`).join('\n')}

**Recent Guest Research:**
${recentGuests.map(g => `- ${g.guest_name}`).join('\n')}

**Analytics:**
${analytics ? `- Visitors: ${analytics.visitors}\n- Conversions: ${analytics.conversions}` : 'Not available'}

Create an engaging daily digest email with:
1. **Greeting** (friendly, energetic)
2. **This Week's Progress** (highlight episodes and guests)
3. **Key Metrics** (if available)
4. **Today's Focus** (1-2 suggested priorities)
5. **Motivational Close**

Format in HTML with good styling. Keep it under 500 words.`;

const digestResult = await platform.generateText({
  prompt: digestPrompt,
  maxTokens: 1000
});

// Step 5: Send email to workspace owner
// Note: You'd need to query workspace owner email or configure recipient
const recipientEmail = request.body.recipientEmail || 'owner@example.com';

console.log(`[Daily Digest] Sending email to ${recipientEmail}`);

await platform.sendEmail({
  to: recipientEmail,
  subject: `🎙️ Your Daily Podcast Digest - ${new Date().toLocaleDateString()}`,
  html: digestResult.text,
  text: digestResult.text.replace(/<[^>]+>/g, '') // Strip HTML for plain text
});

console.log('[Daily Digest] Digest sent successfully');

return respond(200, {
  success: true,
  episodesIncluded: recentEpisodes.length,
  guestsIncluded: recentGuests.length,
  sentTo: recipientEmail,
  timestamp: new Date().toISOString()
});
```

**Scheduling Setup**:
```javascript
// Create a booster to trigger this daily
// From Otto or via API:
await createBooster({
  name: 'Daily Digest Sender',
  target: { type: 'all_sessions' }, // Trigger context
  when: {
    type: 'recurring',
    interval: { every: 1, unit: 'days' },
    at: '09:00',
    timezone: 'America/New_York'
  },
  delivery: { channel: 'current' },
  content: {
    type: 'text',
    value: 'Triggering daily digest...'
  },
  // Or call via external cron service hitting the hook endpoint
});
```

---

### Example 4: Content Calendar API

**Purpose**: Query episodes and generate an AI-optimized posting schedule for social media.

**Pattern**: Query-Enrich-Respond

**Code**:
```javascript
// Hook name: content-calendar-api
// Endpoint: POST /api/hooks/execute/workspace-{configId}/content-calendar-api

const { startDate, endDate, platformType } = request.body;

if (!startDate || !endDate) {
  return respond(400, { error: 'Missing required fields: startDate, endDate (YYYY-MM-DD)' });
}

console.log(`[Content Calendar] Generating calendar from ${startDate} to ${endDate}`);

// Step 1: Query all episodes with clips
const episodes = await db.query('episodes', {
  filters: [
    { column: 'suggested_clips', operator: 'not_null' }
  ],
  orderBy: { column: 'created_at', direction: 'desc' },
  limit: 50
});

console.log(`[Content Calendar] Found ${episodes.length} episodes with clips`);

// Step 2: Extract all clips
const allClips = [];
for (const episode of episodes) {
  try {
    const clips = JSON.parse(episode.suggested_clips || '[]');
    clips.forEach(clip => {
      allClips.push({
        episodeId: episode.id,
        episodeTitle: episode.title,
        guestName: episode.guest_name,
        ...clip
      });
    });
  } catch (error) {
    console.error(`[Content Calendar] Failed to parse clips for episode ${episode.id}`);
  }
}

console.log(`[Content Calendar] Total clips available: ${allClips.length}`);

// Step 3: Generate optimal posting schedule with AI
const calendarPrompt = `You are a social media strategist creating a content calendar.

Platform: ${platformType || 'Instagram/TikTok'}
Date Range: ${startDate} to ${endDate}
Available Clips: ${allClips.length}

Clips:
${allClips.slice(0, 20).map((clip, i) =>
  `${i+1}. "${clip.title}" from episode "${clip.episodeTitle}" (Guest: ${clip.guestName})`
).join('\n')}

Create an optimal posting schedule with:
- Post frequency: 3-5 posts per week
- Best times: Weekdays 10 AM, 2 PM, 6 PM; Weekends 11 AM, 4 PM
- Variety: Mix different guests and topics
- Engagement optimization: Spread viral-worthy clips evenly

Return as JSON array:
[
  {
    "date": "YYYY-MM-DD",
    "time": "HH:MM",
    "clipIndex": 0,
    "clipTitle": "...",
    "episodeTitle": "...",
    "rationale": "Why this time/clip"
  }
]`;

const calendarResult = await platform.generateText({
  prompt: calendarPrompt,
  maxTokens: 2000
});

// Parse schedule
let schedule = [];
try {
  const jsonMatch = calendarResult.text.match(/\[[\s\S]*\]/);
  if (jsonMatch) {
    schedule = JSON.parse(jsonMatch[0]);
  }
} catch (parseError) {
  console.error('[Content Calendar] Failed to parse schedule JSON:', parseError.message);
  return respond(500, { error: 'Failed to generate valid schedule' });
}

// Step 4: Enrich schedule with full clip data
const enrichedSchedule = schedule.map(item => ({
  ...item,
  clip: allClips[item.clipIndex] || null,
  episodeId: allClips[item.clipIndex]?.episodeId
}));

console.log(`[Content Calendar] Generated ${enrichedSchedule.length} scheduled posts`);

return respond(200, {
  success: true,
  schedule: enrichedSchedule,
  totalPosts: enrichedSchedule.length,
  dateRange: { startDate, endDate },
  platformType: platformType || 'default'
});
```

**Client Usage**:
```typescript
// From content calendar mini-app
const result = await fetch('/api/hooks/execute/workspace-{configId}/content-calendar-api', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'x-workspace-id': workspaceId,
    'x-api-key': apiKey
  },
  body: JSON.stringify({
    startDate: '2026-04-01',
    endDate: '2026-04-30',
    platformType: 'Instagram'
  })
});

const data = await result.json();

// Render calendar in UI
data.schedule.forEach(post => {
  console.log(`${post.date} at ${post.time}: "${post.clipTitle}"`);
  console.log(`  Rationale: ${post.rationale}`);
});
```

---

## Decision Framework

### When to Build a Composite API

✅ **Build a composite API when:**

- Multiple mini-apps need the same multi-step operation
- External webhooks trigger the workflow
- The operation should be atomic (all-or-nothing)
- You want to reduce client-side complexity
- The workflow is well-defined and stable
- You need consistent error handling across calls

**Example**: Guest research that multiple apps (episode planner, guest manager) need to call.

### When to Orchestrate Locally

✅ **Orchestrate locally when:**

- You're prototyping and the pattern isn't stable yet
- Only one script or app needs it (not reusable)
- The user needs to approve intermediate steps
- You're running a one-off migration or data fix
- Real-time feedback is important

**Example**: A one-time script to backfill all episode transcripts with AI summaries.

### Decision Tree

```
Does more than one client need to call this?
├─ Yes → Build composite API
└─ No
   └─ Is it triggered by webhook/schedule?
      ├─ Yes → Build composite API
      └─ No
         └─ Is the pattern stable and reusable?
            ├─ Yes → Consider building composite API (future-proofing)
            └─ No → Orchestrate locally, extract later if needed
```

---

## Best Practices

### 1. Keep Composite APIs Focused

**DO**: Single responsibility
```javascript
// GOOD: Does one thing well
// Hook: generate-episode-summary
const summary = await platform.generateText({
  prompt: `Summarize: ${transcript}`
});
return respond(200, { summary: summary.text });
```

**DON'T**: Kitchen sink
```javascript
// BAD: Does too many unrelated things
// Hook: do-everything
const summary = await platform.generateText(...);
await db.insert('summaries', ...);
await platform.sendEmail(...);
await fetch('https://external-api.com/notify');
const analytics = await fetch('/api/analytics');
// ... 100 more lines
```

### 2. Return Structured Responses

**DO**: Consistent, typed responses
```javascript
return respond(200, {
  success: true,
  data: {
    episodeId: 42,
    summary: "...",
    clipsGenerated: 5
  },
  meta: {
    processingTime: 1.2,
    tokensUsed: 450
  }
});
```

**DON'T**: Unstructured text blobs
```javascript
return respond(200, "Done! Generated summary and 5 clips.");
```

### 3. Log Important Steps

**DO**: Log for debugging
```javascript
console.log('[Guest Research] Starting research for', guestName);
console.log('[Guest Research] Fetching', urls.length, 'URLs');
console.log('[Guest Research] Generated', summary.length, 'char summary');
console.log('[Guest Research] Saved to database');
```

**DON'T**: Silent execution
```javascript
// No logs = impossible to debug
await fetch(url);
const summary = await platform.generateText(...);
await db.insert(...);
```

### 4. Handle Errors Gracefully

**DO**: Try/catch with fallbacks
```javascript
try {
  const response = await fetch(url);
  const data = await response.json();
} catch (error) {
  console.error(`[API] Failed to fetch ${url}:`, error.message);
  // Continue with default or skip this source
  return respond(500, {
    error: 'Failed to fetch external data',
    details: error.message
  });
}
```

**DON'T**: Let errors crash
```javascript
// BAD: Uncaught errors kill the function
const response = await fetch(url); // May throw
const data = await response.json(); // May throw
```

### 5. Validate Inputs

**DO**: Check required fields early
```javascript
const { episodeId, transcript } = request.body;

if (!episodeId || !transcript) {
  return respond(400, {
    error: 'Missing required fields',
    required: ['episodeId', 'transcript']
  });
}

if (transcript.length < 100) {
  return respond(400, { error: 'Transcript too short (min 100 chars)' });
}
```

### 6. Use JSON for Complex Data

**DO**: JSON for arrays/objects
```javascript
await db.insert('episodes', [{
  title: 'Episode 1',
  metadata: JSON.stringify({
    duration: 3600,
    fileSize: 52428800,
    format: 'mp3'
  })
}]);
```

**DON'T**: String concatenation
```javascript
// BAD: Hard to parse later
await db.insert('episodes', [{
  metadata: `duration:3600,fileSize:52428800,format:mp3`
}]);
```

### 7. Document Your API

Include a header comment in each hook:

```javascript
/**
 * Guest Research API
 *
 * Fetches content from URLs, generates AI summary, stores in database.
 *
 * POST /api/hooks/execute/workspace-{configId}/guest-research-api
 *
 * Body:
 *   - guestName (string, required): Name of the guest
 *   - urls (array, required): URLs to research
 *   - episodeId (integer, optional): Link to episode
 *
 * Returns:
 *   - success (boolean)
 *   - guestName (string)
 *   - summary (string): AI-generated research summary
 *   - sourcesProcessed (integer)
 *   - sourcesFailed (integer)
 *
 * Example:
 *   POST { "guestName": "Jane Doe", "urls": ["https://janedoe.com"] }
 */

const { guestName, urls, episodeId } = request.body;
// ... implementation
```

---

## Debugging Composite APIs

### View Logs
Use `get_hook_logs` to see execution history:

```bash
# From Otto
get_hook_logs --hookName guest-research-api --limit 5
```

### Test in Isolation
Use `test_server_function` to run with sample data:

```javascript
test_server_function({
  hookName: 'guest-research-api',
  body: {
    guestName: 'Test Guest',
    urls: ['https://example.com']
  }
});
```

### Add Debug Logging
Temporary verbose logging:

```javascript
console.log('[DEBUG] Request:', JSON.stringify(request.body, null, 2));
console.log('[DEBUG] Query result:', rows);
console.log('[DEBUG] AI response:', result.text.slice(0, 200));
```

---

## Performance Tips

### 1. Limit Data Fetched
```javascript
// GOOD: Limit queries
const episodes = await db.query('episodes', { limit: 50 });

// BAD: Fetch everything
const episodes = await db.query('episodes'); // Could be 10,000+ rows
```

### 2. Truncate Text for AI
```javascript
// GOOD: Limit tokens
const summary = await platform.generateText({
  prompt: `Summarize: ${transcript.slice(0, 8000)}`, // ~2k tokens
  maxTokens: 500
});

// BAD: Pass entire 50k word transcript
const summary = await platform.generateText({
  prompt: `Summarize: ${transcript}` // May exceed context window
});
```

### 3. Parallelize Independent Operations
```javascript
// GOOD: Parallel fetches
const [episodes, guests, analytics] = await Promise.all([
  db.query('episodes', { limit: 10 }),
  db.query('guests', { limit: 10 }),
  fetch('/api/analytics')
]);

// BAD: Sequential (slower)
const episodes = await db.query('episodes', { limit: 10 });
const guests = await db.query('guests', { limit: 10 });
const analytics = await fetch('/api/analytics');
```

### 4. Cache Expensive Operations
Store results to avoid recomputation:

```javascript
// Check if already processed
const existing = await db.query('processed_episodes', {
  filters: [{ column: 'episode_id', operator: 'eq', value: episodeId }],
  limit: 1
});

if (existing.length > 0) {
  console.log('[Cache] Using cached result');
  return respond(200, { cached: true, data: existing[0] });
}

// ... expensive processing
```

---

## Conclusion

Composite APIs are powerful tools for building sophisticated workflows on Audos. They:

- ✅ Combine primitives into reusable operations
- ✅ Reduce latency by keeping logic server-side
- ✅ Simplify client code
- ✅ Enable atomic transactions
- ✅ Centralize business logic

**Start simple**: Build one composite API that combines 2-3 primitives. Test it. Then build more complex aggregations as patterns emerge.

**Iterate**: Begin by orchestrating locally in a script. Once the pattern stabilizes, extract it into a composite API that multiple clients can call.

**Document**: Always include clear comments on what your API does, what it expects, and what it returns.

Happy building! 🚀
