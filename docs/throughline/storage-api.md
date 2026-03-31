# Storage API

> **Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/storage-api`
>
> **Method:** `POST`
>
> **Content-Type:** `application/json`

Upload, list, and manage media files in the workspace.

---

## Actions

### list

List all media files in the workspace.

**Request:**
```json
{
  "action": "list",
  "category": "all",
  "limit": 20
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"list"` |
| `category` | string | No | Filter: `all`, `attachment`, `generated`, `reference`, `asset` |
| `limit` | integer | No | Max files to return (default: 20) |

**Response:**
```json
{
  "success": true,
  "files": [
    {
      "id": "abc123",
      "filename": "episode-cover.png",
      "url": "https://storage.googleapis.com/...",
      "contentType": "image/png",
      "category": "asset"
    }
  ],
  "count": 1
}
```

---

### upload

Upload a file using base64-encoded content.

**Request:**
```json
{
  "action": "upload",
  "filename": "guest-photo.jpg",
  "contentType": "image/jpeg",
  "base64": "/9j/4AAQSkZJRgABAQAA...",
  "category": "attachment",
  "description": "Photo of guest Jane Smith"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"upload"` |
| `filename` | string | Yes | Filename with extension |
| `contentType` | string | Yes | MIME type (e.g., `image/jpeg`, `application/pdf`) |
| `base64` | string | Yes | Base64-encoded file content |
| `category` | string | No | Category: `attachment`, `generated`, `reference`, `asset` |
| `description` | string | No | Description of the file |

**Response:**
```json
{
  "success": true,
  "url": "https://storage.googleapis.com/audos-images/...",
  "mediaId": "def456",
  "message": "File uploaded successfully"
}
```

---

### upload-from-url

Upload a file by fetching it from a URL.

**Request:**
```json
{
  "action": "upload-from-url",
  "url": "https://example.com/image.jpg",
  "filename": "downloaded-image.jpg",
  "contentType": "image/jpeg",
  "category": "reference"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"upload-from-url"` |
| `url` | string | Yes | URL to fetch the file from |
| `filename` | string | Yes | Filename to save as |
| `contentType` | string | Yes | MIME type |
| `category` | string | No | Category |
| `description` | string | No | Description |

---

## Code Examples

### Python - Upload a File

```python
import requests
import base64

# Read and encode the file
with open("photo.jpg", "rb") as f:
    encoded = base64.b64encode(f.read()).decode()

response = requests.post(
    "https://audos.com/api/hooks/execute/workspace-351699/storage-api",
    json={
        "action": "upload",
        "filename": "photo.jpg",
        "contentType": "image/jpeg",
        "base64": encoded,
        "category": "attachment"
    }
)
print(response.json()["url"])
```

### JavaScript - List Files

```javascript
const response = await fetch(
  "https://audos.com/api/hooks/execute/workspace-351699/storage-api",
  {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      action: "list",
      category: "asset",
      limit: 10
    })
  }
);
const data = await response.json();
console.log(data.files);
```

---

## Use Cases

1. **Upload guest photos** for briefing documents
2. **Store episode artwork** and thumbnails
3. **Manage media assets** for social posts
4. **Archive transcripts and documents**

---

*Part of the [Throughline API](./README.md)*
