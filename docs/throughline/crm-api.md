# CRM API

> **Endpoint:** `https://audos.com/api/hooks/execute/workspace-351699/crm-api`
>
> **Method:** `POST`
>
> **Content-Type:** `application/json`

Manage contacts, leads, and customer relationships.

---

## Actions

### list

List all contacts with optional filtering.

**Request:**
```json
{
  "action": "list",
  "limit": 50,
  "hasEmail": true
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"list"` |
| `limit` | integer | No | Max contacts to return (default: 50) |
| `hasEmail` | boolean | No | Filter to contacts with email addresses |
| `days` | integer | No | Filter by contacts created in last N days |
| `sourceCategory` | string | No | Filter by source: `all`, `organic`, `meta`, `instagram`, `manual` |
| `tags` | array | No | Filter by tags (contacts must have ALL specified tags) |

**Response:**
```json
{
  "success": true,
  "contacts": [
    {
      "id": "abc123",
      "email": "guest@example.com",
      "name": "Jane Smith",
      "tags": ["guest", "interviewed"],
      "source": "organic",
      "firstSeen": "2026-03-15T10:30:00Z"
    }
  ],
  "count": 1
}
```

---

### create

Create a new contact.

**Request:**
```json
{
  "action": "create",
  "email": "newguest@example.com",
  "name": "John Doe",
  "phone": "+1-555-123-4567",
  "instagram": "@johndoe",
  "source": "manual"
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"create"` |
| `email` | string | Yes | Contact email address |
| `name` | string | No | Full name |
| `phone` | string | No | Phone number |
| `instagram` | string | No | Instagram handle |
| `source` | string | No | Where this contact came from |

**Response:**
```json
{
  "success": true,
  "contact": {
    "id": "def456",
    "email": "newguest@example.com",
    "name": "John Doe"
  },
  "message": "Contact created successfully"
}
```

---

### update

Update an existing contact.

**Request:**
```json
{
  "action": "update",
  "contactId": "abc123",
  "name": "Jane Smith-Johnson",
  "addTags": ["vip", "repeat-guest"],
  "removeTags": ["prospect"]
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"update"` |
| `contactId` | string | Yes | Contact ID to update |
| `name` | string | No | Updated name |
| `email` | string | No | Updated email |
| `phone` | string | No | Updated phone |
| `instagram` | string | No | Updated Instagram handle |
| `notes` | string | No | Notes about the contact |
| `addTags` | array | No | Tags to add |
| `removeTags` | array | No | Tags to remove |

---

### add-tags / remove-tags

Bulk tag operations on multiple contacts.

**Request:**
```json
{
  "action": "add-tags",
  "tags": ["newsletter", "2026-cohort"],
  "filter": {
    "sourceCategory": "organic",
    "startDate": "2026-01-01"
  }
}
```

**Parameters:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | `"add-tags"` or `"remove-tags"` |
| `tags` | array | Yes | Tags to add or remove |
| `filter` | object | Yes | Filter criteria for which contacts to update |

**Filter options:**
- `contactIds` - Specific contact IDs
- `sourceCategory` - Filter by source
- `startDate` / `endDate` - Date range
- `existingTags` - Contacts that have these tags
- `excludeTags` - Exclude contacts with these tags
- `hasEmail` - Only contacts with email

---

## Code Examples

### Python - List Contacts

```python
import requests

response = requests.post(
    "https://audos.com/api/hooks/execute/workspace-351699/crm-api",
    json={
        "action": "list",
        "limit": 100,
        "hasEmail": True
    }
)
contacts = response.json()["contacts"]
for contact in contacts:
    print(f"{contact['name']}: {contact['email']}")
```

### JavaScript - Create a Contact

```javascript
const response = await fetch(
  "https://audos.com/api/hooks/execute/workspace-351699/crm-api",
  {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      action: "create",
      email: "potential-guest@example.com",
      name: "Sarah Connor",
      source: "referral"
    })
  }
);
const data = await response.json();
console.log("Created contact:", data.contact.id);
```

### cURL - Add Tags to Contacts

```bash
curl -X POST "https://audos.com/api/hooks/execute/workspace-351699/crm-api" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "add-tags",
    "tags": ["follow-up-needed"],
    "filter": {
      "sourceCategory": "organic",
      "startDate": "2026-03-01"
    }
  }'
```

---

## Use Cases

1. **Guest management** - Track potential and past podcast guests
2. **Audience segmentation** - Tag contacts by interest or engagement
3. **Outreach tracking** - Mark contacts as "contacted", "responded", "booked"
4. **Newsletter subscribers** - Manage email list from sign-ups
5. **Integration sync** - Push/pull contacts from external CRMs

---

*Part of the [Throughline API](./README.md)*
