# Web API

> **Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/web-api`
>
> **Method:** `POST`
>
> **Content-Type:** `application/json`

Fetch web pages and search the web for content.

---

## Actions

### fetch

Fetch and extract content from a URL.

**Request:**
```json
{
  "action": "fetch",
  "url": "https://example.com/article"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"fetch"` |
| `url` | string | Yes | URL to fetch content from |

**Response:**
```json
{
  "success": true,
  "url": "https://example.com/article",
  "title": "Article Title",
  "content": "The extracted text content of the page...",
  "contentLength": 5432
}
```

---

### search

Search the web using Google.

**Request:**
```json
{
  "action": "search",
  "query": "podcast marketing strategies 2026",
  "num": 5
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"search"` |
| `query` | string | Yes | Search query |
| `num` | integer | No | Number of results (1-10, default: 5) |

**Response:**
```json
{
  "success": true,
  "query": "podcast marketing strategies 2026",
  "results": [
    {
      "title": "Top Podcast Marketing Strategies",
      "url": "https://example.com/marketing",
      "snippet": "Learn the best strategies for growing your podcast..."
    }
  ],
  "count": 5
}
```

---

## Code Examples

### Python - Fetch a Page

```python
import requests

response = requests.post(
    "https://audos.com/api/hooks/execute/workspace-351699/web-api",
    json={
        "action": "fetch",
        "url": "https://example.com/guest-bio"
    }
)
data = response.json()
print(data["content"])
```

### JavaScript - Search the Web

```javascript
const response = await fetch(
  "https://audos.com/api/hooks/execute/workspace-351699/web-api",
  {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      action: "search",
      query: "podcast guest research tips",
      num: 10
    })
  }
);
const data = await response.json();
console.log(data.results);
```

### cURL - Fetch Content

```bash
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/web-api" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "fetch",
    "url": "https://example.com/article"
  }'
```

---

## Use Cases

1. **Guest research** - Fetch LinkedIn profiles, company pages, past interviews
2. **Content inspiration** - Search for trending topics in your niche
3. **Competitive analysis** - Fetch competitor podcast descriptions
4. **Show notes research** - Pull relevant articles for episode topics
5. **Link verification** - Check that links in show notes are valid

---

*Part of the [Throughline API](./README.md)*
