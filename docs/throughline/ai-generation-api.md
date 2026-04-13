# AI Generation API

> **Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/ai-api`
>
> **Method:** `POST`
>
> **Content-Type:** `application/json`

Generate text content using GPT-4o-mini.

---

## Actions

### generate

Generate text from a prompt.

**Request:**
```json
{
  "action": "generate",
  "prompt": "Write a LinkedIn post about my latest podcast episode on sustainability",
  "systemPrompt": "You are a social media expert for podcasters. Write in a conversational, engaging tone."
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | Always `"generate"` |
| `prompt` | string | Yes | The user prompt / what you want generated |
| `systemPrompt` | string | No | System instructions to guide the AI |
| `maxTokens` | integer | No | Maximum response length |
| `temperature` | number | No | Creativity (0-1, default ~0.7) |

**Response:**
```json
{
  "success": true,
  "text": "Just dropped a new episode about sustainability...",
  "model": "gpt-4o-mini-2024-07-18",
  "usage": {
    "promptTokens": 45,
    "completionTokens": 150,
    "totalTokens": 195
  }
}
```

---

### chat

Multi-turn conversation. Use this when you need to pass feedback history — each prior response and the user's reaction to it — as real conversation turns rather than simulated text. The model sees full turn context and weights the most recent user message highest, making iterative refinement significantly more reliable than single-turn with packed history.

**Request:**
```json
{
  "action": "chat",
  "messages": [
    { "role": "user", "content": "Write a caption for this clip about slow pandemics." },
    { "role": "assistant", "content": "Climate change is a slow pandemic — creeping in silently but powerfully." },
    { "role": "user", "content": "Good direction. Make it two sentences only." }
  ],
  "systemPrompt": "You are a social media copywriter for the podcast Throughline.",
  "maxTokens": 1200
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | Always `"chat"` |
| `messages` | array | Yes | Alternating `user`/`assistant` turns. Last message must be `role: "user"`. |
| `messages[].role` | string | Yes | `"user"` or `"assistant"` |
| `messages[].content` | string | Yes | The turn content |
| `systemPrompt` | string | No | System instructions prepended to the conversation |
| `maxTokens` | integer | No | Maximum response length (default 1200) |
| `temperature` | number | No | Creativity (0–1, default 0.7) |

**Response:** Same shape as `generate`.

```json
{
  "success": true,
  "text": "Climate change is a slow pandemic, creeping in silently. Abhinav puts it plainly: nature speaks — and we need to listen.",
  "model": "gpt-4o-mini-2024-07-18",
  "usage": {
    "promptTokens": 120,
    "completionTokens": 32,
    "totalTokens": 152
  }
}
```

**When to use `chat` vs `generate`:**

- Use `generate` for one-shot generation with no prior context.
- Use `chat` any time the user has already seen and reacted to a prior output. Packing conversation history into a single `prompt` string degrades instruction-following — the model treats simulated history as background text rather than conversation turns, so late-arriving instructions like "make it shorter" lose priority against earlier system prompt rules.

---

## Comparison

| | `generate` | `chat` |
|---|---|---|
| `action` | `"generate"` | `"chat"` |
| `prompt` | Required | Not used |
| `messages` | Not used | Required |
| `systemPrompt` | Optional | Optional |
| `maxTokens` | Optional (default 1200) | Optional (default 1200) |
| `temperature` | Optional (default 0.7) | Optional (default 0.7) |

---

## Code Examples

### Python

```python
import requests

response = requests.post(
    "https://audos.com/api/hooks/execute/workspace-351699/ai-api",
    json={
        "action": "generate",
        "prompt": "Write a tweet about my podcast",
        "systemPrompt": "Keep it under 280 characters, engaging, with relevant hashtags"
    }
)
print(response.json()["text"])
```

### JavaScript

```javascript
const response = await fetch(
  "https://audos.com/api/hooks/execute/workspace-351699/ai-api",
  {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      action: "generate",
      prompt: "Write a tweet about my podcast",
      systemPrompt: "Keep it under 280 characters"
    })
  }
);
const data = await response.json();
console.log(data.text);
```

### cURL

```bash
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/ai-api" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "generate",
    "prompt": "Write a LinkedIn post about my podcast",
    "systemPrompt": "Professional tone, include a call to action"
  }'
```

### Multi-turn (JavaScript)

```javascript
// Round 1 — initial generation
const round1 = await fetch(AI_API_URL, {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    action: "chat",
    messages: [
      { role: "user", content: "Write a caption for this clip about slow pandemics." }
    ],
    systemPrompt: "You are a social media copywriter for Throughline."
  })
});
const r1 = await round1.json();

// Round 2 — user gives feedback, new turn appended
const round2 = await fetch(AI_API_URL, {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    action: "chat",
    messages: [
      { role: "user",      content: "Write a caption for this clip about slow pandemics." },
      { role: "assistant", content: r1.text },
      { role: "user",      content: "Good direction. Make it two sentences only." }
    ],
    systemPrompt: "You are a social media copywriter for Throughline."
  })
});
const r2 = await round2.json();
console.log(r2.text); // Two-sentence caption
```

---

## Use Cases

1. **Generate social media captions** for podcast episodes
2. **Create show notes** from transcripts
3. **Draft guest outreach emails**
4. **Generate episode titles and descriptions**
5. **Create newsletter content**
6. **Iterative refinement loops** — use `chat` when building training flows where users give feedback across multiple rounds. Each round's output becomes an `assistant` turn; the user's feedback becomes the next `user` turn. The model weights the final user message highest, so instructions like "make it two sentences" reliably override earlier defaults.

---

*Part of the [Throughline API](./README.md)*
