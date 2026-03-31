# Analytics API

> **Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/analytics-api`
>
> **Method:** `POST`
>
> **Content-Type:** `application/json`

Access visitor metrics, funnel data, and engagement analytics.

---

## Actions

### overview

Get a high-level analytics overview.

**Request:**
```json
{
  "action": "overview",
  "days": 30
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"overview"` |
| `days` | integer | No | Number of days to analyze (default: 30) |
| `startDate` | string | No | Start date (ISO format, overrides `days`) |
| `endDate` | string | No | End date (ISO format) |

**Response:**
```json
{
  "success": true,
  "period": {
    "start": "2026-03-01",
    "end": "2026-03-31",
    "days": 30
  },
  "metrics": {
    "visitors": 1250,
    "uniqueVisitors": 980,
    "pageViews": 3400,
    "emailCaptures": 45,
    "conversionRate": 4.6
  }
}
```

---

### funnel

Get conversion funnel metrics.

**Request:**
```json
{
  "action": "funnel",
  "days": 30
}
```

**Response:**
```json
{
  "success": true,
  "funnel": {
    "landingPageViews": 1500,
    "emailsCaptured": 120,
    "spaceEntered": 85,
    "appOpened": 60,
    "purchases": 5
  },
  "conversionRates": {
    "viewToEmail": 8.0,
    "emailToSpace": 70.8,
    "spaceToPurchase": 5.9
  }
}
```

---

### events

Query specific event types.

**Request:**
```json
{
  "action": "events",
  "eventType": "email_submit",
  "days": 7,
  "aggregation": "by_day"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"events"` |
| `eventType` | string | No | Filter: `page_view`, `email_submit`, `cta_click`, `purchase`, etc. |
| `days` | integer | No | Days to query |
| `aggregation` | string | No | `none`, `by_type`, `by_day`, `summary` |
| `limit` | integer | No | Max events to return (for `aggregation: none`) |

**Response (aggregation: by_day):**
```json
{
  "success": true,
  "eventType": "email_submit",
  "data": [
    { "date": "2026-03-25", "count": 8 },
    { "date": "2026-03-26", "count": 12 },
    { "date": "2026-03-27", "count": 5 }
  ]
}
```

---

### sessions

Get visitor session data.

**Request:**
```json
{
  "action": "sessions",
  "limit": 20,
  "hasEmail": true
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"sessions"` |
| `limit` | integer | No | Max sessions to return (default: 20) |
| `hasEmail` | boolean | No | Filter to sessions with email captured |

**Response:**
```json
{
  "success": true,
  "sessions": [
    {
      "id": "sess_abc123",
      "email": "visitor@example.com",
      "firstSeen": "2026-03-27T14:30:00Z",
      "lastActivity": "2026-03-27T14:45:00Z",
      "pageViews": 5,
      "source": "organic"
    }
  ],
  "count": 20
}
```

---

### visitors

Get unique visitor counts by time period.

**Request:**
```json
{
  "action": "visitors",
  "days": 30,
  "groupBy": "day"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"visitors"` |
| `days` | integer | No | Days to analyze |
| `groupBy` | string | No | `day`, `week`, `month` |

---

## Code Examples

### Python - Get Analytics Overview

```python
import requests

response = requests.post(
    "https://audos.com/api/hooks/execute/workspace-351699/analytics-api",
    json={
        "action": "overview",
        "days": 30
    }
)
data = response.json()
metrics = data["metrics"]
print(f"Visitors: {metrics['visitors']}")
print(f"Conversion Rate: {metrics['conversionRate']}%")
```

### JavaScript - Get Funnel Data

```javascript
const response = await fetch(
  "https://audos.com/api/hooks/execute/workspace-351699/analytics-api",
  {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      action: "funnel",
      days: 7
    })
  }
);
const data = await response.json();
console.log("Email capture rate:", data.conversionRates.viewToEmail + "%");
```

### cURL - Query Events

```bash
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/analytics-api" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "events",
    "eventType": "page_view",
    "days": 7,
    "aggregation": "by_day"
  }'
```

---

## Event Types

| Event Type | Description |
|------------|-------------|
| `page_view` | Landing page or space page viewed |
| `landing_page_view` | Specifically landing page views |
| `email_submit` | Email captured via form |
| `email_captured` | Same as email_submit |
| `cta_click` | Call-to-action button clicked |
| `space_entered` | User entered the app/space |
| `app_opened` | Specific app opened within space |
| `app_action` | Action taken within an app |
| `agent_message` | Message sent to AI agent |
| `purchase` | Purchase completed |
| `custom` | Custom tracked events |

---

## Use Cases

1. **Track landing page performance** - Monitor views and email captures
2. **Measure content engagement** - See which apps get the most usage
3. **Analyze conversion funnel** - Identify drop-off points
4. **Monitor growth trends** - Track visitors over time
5. **Attribution analysis** - See where traffic is coming from

---

*Part of the [Throughline API](./README.md)*
