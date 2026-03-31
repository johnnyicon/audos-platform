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

---

## Use Cases

1. **Generate social media captions** for podcast episodes
2. **Create show notes** from transcripts
3. **Draft guest outreach emails**
4. **Generate episode titles and descriptions**
5. **Create newsletter content**

---

*Part of the [Throughline API](./README.md)*
