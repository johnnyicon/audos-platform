# Scheduler API

> **Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/scheduler-api`
>
> **Method:** `POST`
>
> **Content-Type:** `application/json`

Create and manage scheduled tasks, cron jobs, and timed emails.

---

## Actions

### list

List all scheduled tasks.

**Request:**
```json
{
  "action": "list"
}
```

**Response:**
```json
{
  "success": true,
  "schedules": [
    {
      "id": "abc123",
      "name": "Daily content check",
      "frequency": "daily",
      "time": "09:00",
      "enabled": true
    }
  ]
}
```

---

### create

Create a recurring scheduled task.

**Request:**
```json
{
  "action": "create",
  "name": "Weekly Analytics Report",
  "description": "Generate and email weekly stats",
  "frequency": "weekly",
  "time": "09:00",
  "timezone": "America/New_York",
  "hookName": "generate-weekly-report"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"create"` |
| `name` | string | Yes | Name for the schedule |
| `description` | string | No | Description |
| `frequency` | string | Yes | `hourly`, `daily`, `weekly`, `monthly` |
| `time` | string | No | Time in HH:MM format (for daily/weekly/monthly) |
| `timezone` | string | No | IANA timezone (default: UTC) |
| `hookName` | string | No | Server function to trigger |
| `actionPayload` | object | No | Data to pass to the hook |

---

### create-email

Schedule a one-time email for a future time.

**Request:**
```json
{
  "action": "create-email",
  "name": "Guest reminder",
  "scheduledAt": "2026-04-01T14:00:00Z",
  "timezone": "America/New_York",
  "email": {
    "to": "guest@example.com",
    "subject": "Reminder: Interview Tomorrow",
    "text": "Hi! Just a reminder that we have our interview scheduled for tomorrow."
  }
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"create-email"` |
| `name` | string | Yes | Name for the scheduled email |
| `scheduledAt` | string | Yes | ISO 8601 datetime when to send |
| `timezone` | string | No | IANA timezone |
| `email` | object | Yes | Email details: `{to, subject, text, html}` |

---

### delete

Delete a scheduled task.

**Request:**
```json
{
  "action": "delete",
  "scheduleId": "abc123"
}
```

---

### pause / resume

Pause or resume a scheduled task.

**Request:**
```json
{
  "action": "pause",
  "scheduleId": "abc123"
}
```

```json
{
  "action": "resume",
  "scheduleId": "abc123"
}
```

---

## Code Examples

### Python - Schedule a Daily Task

```python
import requests

response = requests.post(
    "https://audos.com/api/hooks/execute/workspace-351699/scheduler-api",
    json={
        "action": "create",
        "name": "Daily Content Generation",
        "frequency": "daily",
        "time": "08:00",
        "timezone": "America/Los_Angeles",
        "hookName": "generate-daily-content"
    }
)
print(response.json())
```

### JavaScript - Schedule an Email

```javascript
const response = await fetch(
  "https://audos.com/api/hooks/execute/workspace-351699/scheduler-api",
  {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      action: "create-email",
      name: "Episode Launch Notification",
      scheduledAt: "2026-04-15T12:00:00Z",
      email: {
        to: "subscribers@example.com",
        subject: "New Episode Just Dropped!",
        text: "Check out our latest episode..."
      }
    })
  }
);
console.log(await response.json());
```

---

## Use Cases

1. **Daily content generation** - Auto-generate social posts each morning
2. **Weekly reports** - Email analytics summaries
3. **Guest reminders** - Send interview reminders
4. **Episode launch emails** - Schedule announcement emails
5. **Recurring data syncs** - Keep external systems updated

---

*Part of the [Throughline API](./README.md)*
