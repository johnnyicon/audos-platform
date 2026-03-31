# Email API

> **Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/email-api`
>
> **Method:** `POST`
>
> **Content-Type:** `application/json`

Send transactional emails from your workspace.

---

## Actions

### send

Send an email to a recipient.

**Request:**
```json
{
  "action": "send",
  "to": "guest@example.com",
  "subject": "Your Episode Briefing is Ready",
  "text": "Hi! Your briefing document for the upcoming episode is ready to review.",
  "html": "<h1>Your Briefing is Ready</h1><p>Click here to review...</p>"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | Always `"send"` |
| `to` | string | Yes | Recipient email address |
| `subject` | string | Yes | Email subject line |
| `text` | string | Yes | Plain text email body |
| `html` | string | No | HTML email body (optional) |
| `replyTo` | string | No | Reply-to email address |

**Response:**
```json
{
  "success": true,
  "message": "Email sent successfully",
  "to": "guest@example.com",
  "subject": "Your Episode Briefing is Ready"
}
```

---

## Code Examples

### Python

```python
import requests

response = requests.post(
    "https://audos.com/api/hooks/execute/workspace-351699/email-api",
    json={
        "action": "send",
        "to": "guest@example.com",
        "subject": "Interview Confirmation",
        "text": "Hi! Just confirming our interview for next Tuesday at 2pm EST.",
        "replyTo": "john@throughline.com"
    }
)
print(response.json())
```

### JavaScript

```javascript
const response = await fetch(
  "https://audos.com/api/hooks/execute/workspace-351699/email-api",
  {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      action: "send",
      to: "guest@example.com",
      subject: "Your Episode is Live!",
      text: "Great news! Your episode just went live.",
      html: "<h1>Your Episode is Live!</h1><p>Listen now...</p>"
    })
  }
);
console.log(await response.json());
```

### cURL

```bash
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/email-api" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "send",
    "to": "guest@example.com",
    "subject": "Thank you for being on the show!",
    "text": "It was great having you on So Good to Grow Good!"
  }'
```

---

## Use Cases

1. **Guest notifications** - Briefing ready, episode live, etc.
2. **Follow-up emails** - Post-interview thank yous
3. **Content alerts** - New episode announcements
4. **Internal notifications** - Task completions, reminders

---

*Part of the [Throughline API](./README.md)*
