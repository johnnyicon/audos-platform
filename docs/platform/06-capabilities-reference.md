# Audos Capabilities Reference

**SDK-level reference documentation for the Audos platform**

This document provides a complete reference of all capabilities available in the Audos platform, including server-side runtime, client-side hooks, database operations, and platform integrations.

> **Verification note (2026-07-14):** this doc was written 2026-03-31 and, until a live capability probe
> on 2026-07-14, had never actually been checked against the running platform for this workspace. Several
> claims below are **confirmed wrong** — see the corrected Storage/Scheduler sections and the new
> "Vector/Embedding Storage" section, all marked with a 2026-07-14 verification callout. Treat any
> unmarked section in this file as still-unverified reference material, not a confirmed finding, per the
> same discipline the rest of this SDK applies to build-job self-reports.

---

## Table of Contents

1. [Server Function Runtime](#server-function-runtime)
2. [Runtime Limitations](#runtime-limitations)
3. [Internal Platform APIs](#internal-platform-apis)
4. [Client-Side Hooks](#client-side-hooks)
5. [Database Capabilities](#database-capabilities)
6. [Platform Integrations](#platform-integrations)
7. [Helper Functions](#helper-functions)

---

## Server Function Runtime

Server functions are JavaScript functions that run server-side when called via API. They execute in a sandboxed Node.js-like environment with specific global objects and capabilities.

### Global Objects

Server functions have access to the following global objects:

- `request` - Incoming HTTP request object
- `respond` - Function to send HTTP responses
- `db` - Database query interface
- `platform` - Platform services (AI, email)
- `fetch` - HTTP client for external requests
- `console` - Logging (console.log, console.error)
- `JSON` - JSON parsing and serialization
- `Date` - Date/time operations
- `Math` - Mathematical operations

### Request Object

The `request` object contains information about the incoming HTTP request.

**Properties:**

```javascript
request.method      // HTTP method: 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'
request.body        // Parsed request body (object or null)
request.query       // Query parameters (object)
request.headers     // HTTP headers (object)
```

**Example:**

```javascript
// Access request data
console.log('Method:', request.method);
console.log('Body:', request.body);
console.log('Query params:', request.query);

// Get specific query parameter
const userId = request.query.userId;

// Get specific header
const authToken = request.headers['authorization'];
```

### Respond Function

The `respond(statusCode, body)` function sends an HTTP response back to the caller.

**Signature:**

```javascript
respond(statusCode: number, body: any): void
```

**Parameters:**

- `statusCode` - HTTP status code (200, 201, 400, 404, 500, etc.)
- `body` - Response body (object, string, number, boolean, or null)

**Examples:**

```javascript
// Success response with data
respond(200, { success: true, data: { id: 123, name: 'John' } });

// Created response
respond(201, { id: 456, message: 'Resource created' });

// Bad request
respond(400, { error: 'Missing required field: email' });

// Not found
respond(404, { error: 'User not found' });

// Server error
respond(500, { error: 'Database query failed' });

// Simple text response
respond(200, 'OK');

// Boolean response
respond(200, true);
```

### Database Interface (`db`)

The `db` object provides methods for database operations. All queries are automatically scoped to your workspace.

#### `db.query(sql)`

Execute a SELECT query and return results.

**Signature:**

```javascript
db.query(sql: string): Promise<Array<object>>
```

**Examples:**

```javascript
// Select all records
const users = await db.query('SELECT * FROM users');

// Select with WHERE clause
const activeUsers = await db.query(
  "SELECT * FROM users WHERE status = 'active'"
);

// Select with ORDER BY and LIMIT
const recentOrders = await db.query(
  'SELECT * FROM orders ORDER BY created_at DESC LIMIT 10'
);

// Join multiple tables
const ordersWithUsers = await db.query(`
  SELECT orders.*, users.name, users.email
  FROM orders
  JOIN users ON orders.user_id = users.id
  WHERE orders.status = 'pending'
`);

// Aggregate functions
const stats = await db.query(`
  SELECT COUNT(*) as total, AVG(amount) as avg_amount
  FROM orders
  WHERE created_at > NOW() - INTERVAL '30 days'
`);
```

#### `db.insert(table, data)`

Insert one or more records into a table.

**Signature:**

```javascript
db.insert(table: string, data: object | Array<object>): Promise<Array<object>>
```

**Examples:**

```javascript
// Insert single record
const newUser = await db.insert('users', {
  name: 'John Doe',
  email: 'john@example.com',
  status: 'active'
});
console.log('Inserted user ID:', newUser[0].id);

// Insert multiple records
const newContacts = await db.insert('contacts', [
  { name: 'Alice', email: 'alice@example.com' },
  { name: 'Bob', email: 'bob@example.com' },
  { name: 'Charlie', email: 'charlie@example.com' }
]);
console.log('Inserted', newContacts.length, 'contacts');

// Insert with JSON data
const newOrder = await db.insert('orders', {
  user_id: 123,
  total: 99.99,
  items: JSON.stringify([
    { id: 1, name: 'Product A', price: 49.99 },
    { id: 2, name: 'Product B', price: 50.00 }
  ])
});
```

#### `db.update(table, data, where)`

Update records that match the WHERE clause.

**Signature:**

```javascript
db.update(table: string, data: object, where: object): Promise<Array<object>>
```

**Examples:**

```javascript
// Update single field
await db.update(
  'users',
  { status: 'inactive' },
  { id: 123 }
);

// Update multiple fields
await db.update(
  'orders',
  { status: 'shipped', shipped_at: new Date().toISOString() },
  { id: 456 }
);

// Update with multiple WHERE conditions
await db.update(
  'subscriptions',
  { status: 'expired' },
  { status: 'active', expires_at: '< NOW()' }
);
```

#### `db.delete(table, where)`

Delete records that match the WHERE clause.

**Signature:**

```javascript
db.delete(table: string, where: object): Promise<number>
```

**Examples:**

```javascript
// Delete by ID
await db.delete('users', { id: 123 });

// Delete with multiple conditions
await db.delete('sessions', {
  status: 'inactive',
  last_active: '< NOW() - INTERVAL \'30 days\''
});

// Delete all records matching condition
await db.delete('logs', { level: 'debug' });
```

#### `db.listTables()`

List all tables in the workspace.

**Signature:**

```javascript
db.listTables(): Promise<Array<string>>
```

**Example:**

```javascript
const tables = await db.listTables();
console.log('Available tables:', tables);
// Returns: ['users', 'orders', 'products', 'sessions']
```

### Platform Services (`platform`)

The `platform` object provides access to platform services like AI generation and email sending.

#### `platform.generateText(prompt)`

Generate text using AI.

> **Verified live 2026-07-14:** the real signature is
> `platform.generateText({ userPrompt, systemPrompt?, model? })`, not a bare string argument. Confirmed
> working — a `"ping"` prompt returned `{"text":"pong","model":"gpt-4o-mini-2024-07-18",
> "usage":{"promptTokens":14,"completionTokens":1,"totalTokens":15}}` in ~1s, both with and without an
> explicit `model`. **Default model is `gpt-4o-mini`** — fine for high-volume, low-judgment work (chunk
> summaries, quiz-question generation) per the tiering strategy in `product-plan.md` §7, not a
> replacement for a stronger model on judgment-heavy generation. There's also a browser-side
> `/proxy/openai/v1/chat/completions` proxy available for other GPT models.

**Signature:**

```javascript
platform.generateText(prompt: string): Promise<string>
```

**Examples:**

```javascript
// Generate product description
const description = await platform.generateText(
  'Write a compelling product description for a wireless mouse with RGB lighting'
);

// Generate email content
const emailBody = await platform.generateText(
  'Write a professional welcome email for a new customer named Sarah'
);

// Generate summarization
const summary = await platform.generateText(
  `Summarize the following text in 2-3 sentences: ${longText}`
);

// Generate structured data
const jsonData = await platform.generateText(
  'Generate a JSON object with 5 sample products (name, price, category)'
);
const products = JSON.parse(jsonData);
```

#### `platform.sendEmail(options)`

Send transactional email.

**Signature:**

```javascript
platform.sendEmail(options: {
  to: string,
  subject: string,
  text: string,
  html?: string
}): Promise<object>
```

**Examples:**

```javascript
// Send plain text email
await platform.sendEmail({
  to: 'user@example.com',
  subject: 'Welcome to our platform',
  text: 'Thank you for signing up! We are excited to have you.'
});

// Send HTML email
await platform.sendEmail({
  to: 'customer@example.com',
  subject: 'Your order has shipped',
  text: 'Your order #12345 has been shipped.',
  html: `
    <h1>Order Shipped</h1>
    <p>Your order <strong>#12345</strong> has been shipped.</p>
    <p>Tracking number: <a href="#">TRACK123</a></p>
  `
});

// Send email with dynamic content
const user = await db.query('SELECT * FROM users WHERE id = 123');
await platform.sendEmail({
  to: user[0].email,
  subject: `Hello ${user[0].name}!`,
  text: `Hi ${user[0].name}, we have an update for you.`,
  html: `<p>Hi <strong>${user[0].name}</strong>, we have an update for you.</p>`
});
```

### HTTP Client (`fetch`)

The `fetch` function makes HTTP requests to external APIs.

**Signature:**

```javascript
fetch(url: string, options?: object): Promise<Response>
```

**Examples:**

```javascript
// GET request
const response = await fetch('https://api.example.com/users');
const data = await response.json();

// POST request with JSON body
const createResponse = await fetch('https://api.example.com/users', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    name: 'John Doe',
    email: 'john@example.com'
  })
});
const result = await createResponse.json();

// PUT request
await fetch('https://api.example.com/users/123', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer token123'
  },
  body: JSON.stringify({ status: 'active' })
});

// DELETE request
await fetch('https://api.example.com/users/123', {
  method: 'DELETE',
  headers: {
    'Authorization': 'Bearer token123'
  }
});

// Handle errors
const apiResponse = await fetch('https://api.example.com/data');
if (!apiResponse.ok) {
  respond(500, { error: 'External API failed' });
  return;
}
const apiData = await apiResponse.json();
```

### Console Logging

Use `console` for debugging and logging.

**Available methods:**

```javascript
console.log('Info message', data);
console.error('Error occurred', error);
console.warn('Warning message');
```

**Example:**

```javascript
console.log('Processing request:', request.method, request.body);

try {
  const result = await db.query('SELECT * FROM users');
  console.log('Query returned', result.length, 'rows');
} catch (error) {
  console.error('Database error:', error);
  respond(500, { error: 'Query failed' });
}
```

---

## Runtime Limitations

Server functions run in a sandboxed environment with specific limitations. Understanding these constraints is critical for writing robust code.

### ❌ NOT Available

#### URLSearchParams

The `URLSearchParams` API is **NOT available**. Use manual query string building instead.

```javascript
// ❌ DOES NOT WORK
const params = new URLSearchParams({ name: 'John', age: 30 });

// ✅ USE THIS INSTEAD
function buildQueryString(params) {
  return Object.entries(params)
    .map(([key, val]) => `${encodeURIComponent(key)}=${encodeURIComponent(val)}`)
    .join('&');
}

const queryString = buildQueryString({ name: 'John', age: 30 });
// Returns: "name=John&age=30"
```

#### Buffer

The `Buffer` class is **NOT available**. Use base64 strings directly.

```javascript
// ❌ DOES NOT WORK
const buffer = Buffer.from('hello', 'utf-8');
const base64 = buffer.toString('base64');

// ✅ USE THIS INSTEAD
// For base64 encoding/decoding, use native methods or work with strings
const base64Encode = (str) => btoa(str);
const base64Decode = (str) => atob(str);
```

#### Module System

**NO** `require()` or `import` statements. Only vanilla JavaScript.

```javascript
// ❌ DOES NOT WORK
const axios = require('axios');
import fetch from 'node-fetch';

// ✅ USE THIS INSTEAD
// Use built-in fetch (available globally)
const response = await fetch('https://api.example.com/data');
```

#### Timers

**NO** `setTimeout` or `setInterval`. Use the scheduler API for delayed/recurring tasks.

```javascript
// ❌ DOES NOT WORK
setTimeout(() => {
  console.log('This will never run');
}, 1000);

// ✅ USE THIS INSTEAD
// Create a scheduled task via the scheduler API
const scheduleResponse = await fetch(
  'https://api.audos.com/api/workspaces/YOUR_WORKSPACE_ID/schedules',
  {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: 'Delayed Task',
      actionType: 'hook',
      hookName: 'my-delayed-function',
      scheduledAt: new Date(Date.now() + 60000).toISOString() // 1 minute from now
    })
  }
);
```

#### Limited Response Object

The `Response` object from `fetch()` has limited methods. **NO** `response.headers.get()`.

```javascript
// ❌ DOES NOT WORK
const response = await fetch('https://api.example.com/data');
const contentType = response.headers.get('content-type');

// ✅ USE THIS INSTEAD
// Access response data directly
const response = await fetch('https://api.example.com/data');
const data = await response.json(); // or response.text()
console.log('Status:', response.status);
console.log('OK:', response.ok);
```

### Execution Constraints

- **Timeout:** ~30 seconds maximum execution time
- **Memory:** ~128MB memory limit
- **No filesystem access:** Cannot read/write files directly
- **No environment variables:** Use database or hardcoded values
- **Synchronous only at top level:** Must use `await` for async operations

**Best Practices:**

```javascript
// ✅ GOOD - Early return on errors
if (!request.body.email) {
  respond(400, { error: 'Email required' });
  return;
}

// ✅ GOOD - Handle timeouts gracefully
try {
  const result = await db.query('SELECT * FROM large_table');
  respond(200, result);
} catch (error) {
  console.error('Query timeout or error:', error);
  respond(500, { error: 'Request timeout' });
}

// ✅ GOOD - Batch operations
const users = await db.insert('users', [
  { name: 'Alice' },
  { name: 'Bob' },
  { name: 'Charlie' }
]); // Single insert vs. 3 separate inserts
```

---

## Internal Platform APIs

Server functions can access internal platform APIs using `fetch()`. These APIs provide access to CRM, analytics, sessions, and more.

### Authentication

All internal API requests must include the workspace ID in the URL path. No additional authentication headers are required when called from server functions.

**URL Pattern:**

```
/api/{service}/{resource}/{workspaceId}
```

### CRM API

Manage contacts and customer data.

#### Get Contacts

```javascript
// GET /api/crm/contacts/{workspaceId}
const response = await fetch(
  'https://api.audos.com/api/crm/contacts/YOUR_WORKSPACE_ID?limit=50'
);
const contacts = await response.json();

// Filter contacts
const response = await fetch(
  'https://api.audos.com/api/crm/contacts/YOUR_WORKSPACE_ID?hasEmail=true&limit=100'
);
```

#### Create Contact

```javascript
// POST /api/crm/contacts/{workspaceId}
const response = await fetch(
  'https://api.audos.com/api/crm/contacts/YOUR_WORKSPACE_ID',
  {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email: 'newuser@example.com',
      name: 'Jane Smith',
      source: 'server-function'
    })
  }
);
const newContact = await response.json();
```

### Analytics API

Retrieve visitor metrics and funnel data.

#### Get Funnel Metrics

```javascript
// GET /api/funnel/metrics/{workspaceId}
const response = await fetch(
  'https://api.audos.com/api/funnel/metrics/YOUR_WORKSPACE_ID?days=30'
);
const metrics = await response.json();

console.log('Visitors:', metrics.visitors);
console.log('Conversions:', metrics.conversions);
console.log('Conversion rate:', metrics.conversionRate);
```

#### Get Events

```javascript
// GET /api/funnel/events/{workspaceId}
const response = await fetch(
  'https://api.audos.com/api/funnel/events/YOUR_WORKSPACE_ID?eventType=email_submit&limit=100'
);
const events = await response.json();
```

### Sessions API

Access visitor session data.

#### Get Sessions

```javascript
// GET /api/funnel/sessions/{workspaceId}
const response = await fetch(
  'https://api.audos.com/api/funnel/sessions/YOUR_WORKSPACE_ID?hasEmail=true&limit=50'
);
const sessions = await response.json();

sessions.forEach(session => {
  console.log('Session ID:', session.id);
  console.log('Email:', session.email);
  console.log('Context:', session.context);
});
```

### Spaces API

Manage spaces and apps.

```javascript
// GET /api/spaces/{workspaceId}
const response = await fetch(
  'https://api.audos.com/api/spaces/YOUR_WORKSPACE_ID'
);
const spaces = await response.json();
```

### App Skills API

Access app backend capabilities (email, scheduler).

#### Send Email

```javascript
// POST /api/app-skills/email/send
const response = await fetch(
  'https://api.audos.com/api/app-skills/email/send',
  {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'x-workspace-id': 'YOUR_WORKSPACE_ID'
    },
    body: JSON.stringify({
      to: 'user@example.com',
      subject: 'Hello from server function',
      text: 'This email was sent from a server function!'
    })
  }
);
```

#### Schedule Email

```javascript
// POST /api/workspaces/{workspaceId}/schedules/email
const response = await fetch(
  'https://api.audos.com/api/workspaces/YOUR_WORKSPACE_ID/schedules/email',
  {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: 'Welcome Email Reminder',
      description: 'Follow-up email 24 hours after signup',
      scheduledAt: new Date(Date.now() + 86400000).toISOString(), // 24 hours
      email: {
        to: 'user@example.com',
        subject: 'Welcome!',
        text: 'Thanks for signing up yesterday!'
      },
      timezone: 'America/New_York'
    })
  }
);
```

#### Create Scheduled Task

```javascript
// POST /api/workspaces/{workspaceId}/schedules
const response = await fetch(
  'https://api.audos.com/api/workspaces/YOUR_WORKSPACE_ID/schedules',
  {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: 'Daily Backup',
      description: 'Run backup every day at 2am',
      frequency: 'daily',
      time: '02:00',
      timezone: 'America/New_York',
      actionType: 'hook',
      hookName: 'backup-data'
    })
  }
);
```

---

## Client-Side Hooks

React apps and mini-apps can use these hooks to interact with the platform.

### useSpaceFiles()

Persist JSON data to space-scoped files.

**Import:**

```javascript
import { useSpaceFiles } from '@/hooks/useSpaceFiles';
```

**Usage:**

```javascript
function MyApp() {
  const { files, saveFile, loadFile, deleteFile } = useSpaceFiles();

  const saveUserData = async () => {
    await saveFile('user-preferences.json', {
      theme: 'dark',
      notifications: true
    });
  };

  const loadUserData = async () => {
    const data = await loadFile('user-preferences.json');
    console.log('User preferences:', data);
  };

  const removeData = async () => {
    await deleteFile('user-preferences.json');
  };

  return (
    <div>
      <button onClick={saveUserData}>Save Preferences</button>
      <button onClick={loadUserData}>Load Preferences</button>
      <button onClick={removeData}>Clear Preferences</button>
    </div>
  );
}
```

### useWorkspaceDB()

Access workspace database tables from React apps.

**Import:**

```javascript
import { useWorkspaceDB } from '@/hooks/useWorkspaceDB';
```

**Usage:**

```javascript
function ProductList() {
  const { query, insert, update, remove } = useWorkspaceDB();
  const [products, setProducts] = useState([]);

  useEffect(() => {
    loadProducts();
  }, []);

  const loadProducts = async () => {
    const result = await query('SELECT * FROM products ORDER BY name');
    setProducts(result);
  };

  const addProduct = async () => {
    await insert('products', {
      name: 'New Product',
      price: 29.99,
      status: 'active'
    });
    await loadProducts();
  };

  const updateProduct = async (id) => {
    await update('products', { price: 39.99 }, { id });
    await loadProducts();
  };

  const deleteProduct = async (id) => {
    await remove('products', { id });
    await loadProducts();
  };

  return (
    <div>
      <button onClick={addProduct}>Add Product</button>
      <ul>
        {products.map(p => (
          <li key={p.id}>
            {p.name} - ${p.price}
            <button onClick={() => updateProduct(p.id)}>Update</button>
            <button onClick={() => deleteProduct(p.id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
}
```

### useSession()

Access the current user session.

**Import:**

```javascript
import { useSession } from '@/hooks/useSession';
```

**Usage:**

```javascript
function UserProfile() {
  const session = useSession();

  if (!session) {
    return <div>Loading...</div>;
  }

  return (
    <div>
      <h1>Welcome!</h1>
      <p>Session ID: {session.id}</p>
      <p>Email: {session.email || 'Not provided'}</p>
      <p>Source: {session.source}</p>
    </div>
  );
}
```

### Stripe Integration Hooks

React apps can integrate Stripe for payments.

**Example:**

```javascript
import { loadStripe } from '@stripe/stripe-js';
import { Elements, CardElement, useStripe, useElements } from '@stripe/react-stripe-js';

const stripePromise = loadStripe('pk_test_YOUR_KEY');

function CheckoutForm() {
  const stripe = useStripe();
  const elements = useElements();

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!stripe || !elements) return;

    const cardElement = elements.getElement(CardElement);
    const { error, paymentMethod } = await stripe.createPaymentMethod({
      type: 'card',
      card: cardElement
    });

    if (error) {
      console.error(error);
    } else {
      console.log('Payment method:', paymentMethod);
      // Process payment via server function
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <CardElement />
      <button type="submit" disabled={!stripe}>Pay</button>
    </form>
  );
}

function App() {
  return (
    <Elements stripe={stripePromise}>
      <CheckoutForm />
    </Elements>
  );
}
```

---

## Database Capabilities

Audos provides PostgreSQL database with workspace isolation.

### Supported Column Types

- `text` - Variable-length text
- `integer` - 32-bit integer
- `bigint` - 64-bit integer
- `decimal` - Arbitrary precision decimal
- `boolean` - True/false
- `timestamp` - Date and time
- `date` - Date only
- `json` - JSON data
- `uuid` - Universally unique identifier

### Auto-Generated Columns

All tables automatically include:

- `id` - Primary key (auto-increment integer)
- `created_at` - Timestamp (defaults to current time)

### Foreign Keys

Tables can reference other workspace tables.

**Example:**

```javascript
// Table: orders
// Columns: id, user_id (foreign key to users.id), total, created_at

// Query with foreign key relationship
const orders = await db.query(`
  SELECT orders.*, users.name as customer_name, users.email
  FROM orders
  JOIN users ON orders.user_id = users.id
  WHERE orders.status = 'pending'
`);
```

### JSON Columns

Store flexible, nested data in JSON columns.

**Example:**

```javascript
// Insert with JSON data
await db.insert('products', {
  name: 'Laptop',
  price: 999.99,
  specifications: JSON.stringify({
    cpu: 'Intel i7',
    ram: '16GB',
    storage: '512GB SSD',
    ports: ['USB-C', 'HDMI', 'Thunderbolt']
  })
});

// Query and parse JSON
const products = await db.query('SELECT * FROM products');
products.forEach(product => {
  const specs = JSON.parse(product.specifications);
  console.log('CPU:', specs.cpu);
  console.log('RAM:', specs.ram);
});

// Query JSON fields (PostgreSQL syntax)
const laptops = await db.query(`
  SELECT * FROM products
  WHERE specifications->>'cpu' LIKE '%i7%'
`);
```

### Vector/Embedding Storage — no native support, but a JSON fallback is proven fast enough (verified live 2026-07-14 and 2026-07-16)

**There is no native vector column type, and no way to enable one — corrected 2026-07-16: pgvector is
*available in the catalog*, not "installed but disabled" as first reported.**

- `workspace_db_create_table` with a `vector(N)` column type is rejected: `{"code":"invalid_enum_value",
  "options":["text","integer","bigint","decimal","boolean","timestamp","date","json","uuid"],
  "message":"Invalid enum value... received 'vector(5)'"}` — `vector` is not in the allowed column-type
  enum, and this enum appears to be the only way to create a table/column.
- A raw query casting to `::vector` fails: `type "vector" does not exist` — the extension isn't enabled
  in this database.
- **Correction (2026-07-16):** a full, unfiltered `SELECT * FROM pg_available_extensions` (70 rows) shows
  `vector` at `default_version: "0.8.1"`, `installed_version: null` — and cross-checking `pg_extension`
  directly shows exactly one row, `plpgsql`. The earlier "installed but not enabled" phrasing was
  imprecise: it is **available but not installed**, full stop, same as every other extension in this
  workspace except `plpgsql`. Also found and missed the first time: **`vectorscale` 0.9.0** (a DiskANN
  ANN-index extension) is in the same available-not-installed state. Other notable available-not-installed
  extensions: `pg_trgm`, `fuzzystrmatch`, `unaccent`, `pg_cron`, `timescaledb`, `postgis`.
- There is no path to enable any of it yourself: raw queries are restricted to `SELECT`/`WITH`/`EXPLAIN`
  only, so `CREATE EXTENSION vector` cannot be run, and DDL escape attempts via `db.rawQuery`/MCP raw
  query are blocked the same way. Enabling it is a platform-side action only Audos can take.

**The workaround — store embeddings as a JSON array, compare brute-force — has one real, load-bearing
result and two open questions, not a closed case.** `db.insert` into a `json` column with an array of
floats works and reads back correctly. A 2026-07-16 benchmark (fake 5-float vectors, real hook execution
against the live platform) found brute-force cosine similarity in hook JavaScript is **sub-millisecond
even at 1,000 rows**: 50 rows → 0.02ms avg, 300 rows → 0.12ms avg, 1,000 rows → 0.28ms avg (single-pass
scan; DB fetch time, not scan time, dominates total request latency). That result is genuinely valid on
its own terms — cosine similarity is fixed arithmetic (multiply, add, divide, sqrt) that costs the same
per number whether the numbers came from a real embedding model or a random-number generator, so
**row-count scaling is proven fast regardless of what the numbers mean.**

What it does *not* prove is that real semantic search works here. Two separate gaps were open as of the
2026-07-16 benchmark; round 3 (same day, later, two sub-jobs) closed both — and the second sub-job
**overturned** the first sub-job's headline finding. Read the correction below before the summary.

**1. Dimensionality — closed, and it doesn't matter.** Round 3 re-ran the identical brute-force cosine
similarity test with genuine 1,536-float vectors (OpenAI embedding size) instead of 5-float placeholders,
same 50/300/1,000 row counts, three runs each for variance:

| rows | dim | avg scan time (cosine loop only) | avg total (incl. DB fetch) |
|---|---|---|---|
| 50 | 1536 | ~0.04–0.14 ms/row (timer-resolution noise at this size) | 342–453 ms total |
| 300 | 1536 | ~0.003 ms/row | 238–248 ms total |
| 1,000 | 1536 | ~0.003 ms/row | 518–541 ms total |

~307x more arithmetic per comparison than the 5-float test barely moved the needle — the scan itself is
1–3ms total even at 1,000 rows. Wall-clock time is now completely dominated by the **database fetch**
(~250–540ms), not the comparison math. The `json` column also returns already-parsed arrays, not strings
needing a second `JSON.parse` pass. **Row-count scaling and dimensionality are both proven fast now, not
estimated.**

**2. Native embedding generation — the first round-3 sub-job's "confirmed absent" finding was a false
negative. Corrected same day.** That sub-job swept `platform.integrations.isAvailable()` with guessed
provider names (`openai-embeddings`, `embeddings`, `vector-search`, etc.) — all `false` — and concluded no
embedding path exists. **It never tried calling `platform.integrations.proxy()` itself, and never tried
the plain provider name.** The second round-3 sub-job did, and found:

- `platform.integrations.isAvailable('openai')` → **`true`** (the correct name is the bare provider name,
  not a feature-specific string — `isAvailable`'s allowlist, read from `proxy`'s own error message, is
  exactly `openai, stripe, twilio, heygen`).
- `platform.integrations.proxy('openai', '/v1/embeddings', { method: 'POST', headers: {...}, body:
  JSON.stringify({ model: 'text-embedding-3-small', input: 'the quick brown fox' }) })` — **succeeded,
  status 200, returned a real 1,536-float OpenAI embedding**, no API key supplied by the calling hook.
  `proxy` is a generic authenticated passthrough to a fixed provider allowlist — Audos holds and injects
  the credential; the hook just names the provider and the path.

**So: Audos does not have a dedicated `platform.generateEmbedding()` method, but it does have a working,
authenticated path to a real embedding provider that doesn't require managing your own API key** — via
the same undocumented `platform.integrations` gateway TEST E's first pass already found but didn't fully
exercise. This is meaningfully better than "no capability, build your own" — it's "the capability exists,
undocumented, and easy to miss if you only test `isAvailable()` with guessed names instead of trying
`proxy()` directly." Call-signature gotchas found empirically: `body` must be a JSON string, not an
object (an options-object 2nd arg silently returns the platform's HTML index page, status 200 — a bad
failure mode, since it looks like success); the response wrapper must be read via `.text()`/`.json()`.

(Side findings from the same probes, filed for awareness: the docs' "Available Globals" list is missing
`platform.integrations`, `platform.externalFetch`, `platform.getLatestSession`, `session`, `respondRaw`,
and `db.bulkInsert` — all live and callable. The "NOT Available (Security)" section claims `eval` is
unavailable; live enumeration shows `eval`, `Function`, `Proxy`, `Reflect`, `WebAssembly`,
`SharedArrayBuffer` all present on `globalThis` — the sandbox is less locked-down than documented, though
`eval` itself was only enumerated, not invoked.)

**Verdict: both gaps closed.** Dimensionality doesn't matter at DoKnow's scale. A real embedding is
obtainable today via `platform.integrations.proxy('openai', '/v1/embeddings', ...)`, no external API key
management required. See `BACKLOG.md #13` (vector storage/search — still open, still worth requesting
native columns for headroom) and `#14` (corrected same day — see that entry for what's now actually being
asked for, since "no embedding capability exists" was wrong).

### Table Creation — `id` is always a platform-generated serial integer, never the uuid you ask for (verified live 2026-07-16)

Building `field-notes`'s own schema (a separate app, separate workspace from DoKnow) surfaced two
undocumented behaviors of `workspace_db_create_table` itself, independent of the vector/embedding
findings above:

- **Requesting `id: uuid` is rejected.** The platform reserves the `id` column and always generates it
  itself as a `serial` (auto-increment integer) — an explicit `id` in the column list returns `column
  "id" specified more than once"`. `uuid` remains a valid type for any *other* column, just not the
  primary key. Design implication: a foreign key meant to reference another table's `id` must be typed
  `integer`, not `uuid`, regardless of what the schema's own logical design calls for.
- **Every table silently gets two extra columns beyond what's requested**: `session_id` (text) and
  `updated_at` (timestamp, default `NOW()`). Neither appears in this doc's column-type reference or in
  the tool's own request/response shape as documented.
- **A failed `CREATE TABLE` call (e.g. an unsatisfiable foreign-key constraint) is not rolled back.** It
  leaves a real, physical, catalog-unregistered table behind that blocks every future attempt to create a
  table with that name (`relation "..." already exists`), with no drop-table tool available to remove it
  — see `BACKLOG.md #15` for the full reproduction and `#16` for the id-type finding.

### Workspace Isolation

All database operations are automatically scoped to your workspace. You cannot access tables from other workspaces.

```javascript
// ✅ This is safe - only queries YOUR workspace
const users = await db.query('SELECT * FROM users');

// ✅ No risk of cross-workspace data leakage
await db.insert('customers', { name: 'John' });
```

---

## Platform Integrations

### Stripe Integration

Process payments and manage subscriptions.

**Server Function Example:**

```javascript
// Create Stripe checkout session
const response = await fetch('https://api.stripe.com/v1/checkout/sessions', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer sk_test_YOUR_SECRET_KEY',
    'Content-Type': 'application/x-www-form-urlencoded'
  },
  body: 'payment_method_types[]=card&line_items[0][price]=price_123&line_items[0][quantity]=1&mode=payment&success_url=https://yoursite.com/success&cancel_url=https://yoursite.com/cancel'
});

const session = await response.json();
respond(200, { checkoutUrl: session.url });
```

### Meta Ads Integration

Manage ad campaigns programmatically (via internal APIs).

```javascript
// Get campaign performance
const response = await fetch(
  'https://api.audos.com/api/ad-campaigns/YOUR_WORKSPACE_ID'
);
const campaigns = await response.json();

campaigns.forEach(campaign => {
  console.log('Campaign:', campaign.name);
  console.log('Impressions:', campaign.impressions);
  console.log('Clicks:', campaign.clicks);
  console.log('CTR:', campaign.ctr);
});
```

### Email Integration

Send transactional emails via platform service or app skills API.

**Via platform.sendEmail:**

```javascript
await platform.sendEmail({
  to: 'customer@example.com',
  subject: 'Order Confirmation',
  text: 'Your order has been received.',
  html: '<h1>Order Confirmed</h1><p>Thank you for your purchase!</p>'
});
```

**Via App Skills API:**

```javascript
const response = await fetch(
  'https://api.audos.com/api/app-skills/email/send',
  {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'x-workspace-id': 'YOUR_WORKSPACE_ID'
    },
    body: JSON.stringify({
      to: 'user@example.com',
      subject: 'Welcome',
      text: 'Thanks for joining!'
    })
  }
);
```

### Storage Integration

Upload and manage files in Google Cloud Storage.

> **Verified live 2026-07-14 — the endpoint above is wrong.** `POST /api/storage/upload` 404s:
> `{"message":"API route not found: POST /api/storage/upload"}`. The two real endpoints are:

**`POST /api/upload/file`** — multipart, any content type. Round-trips byte-identical (md5-verified) for
any file type, including PDFs. **Hard size limit found by bracketing: 52,428,799 bytes (50 MiB − 1)
succeeds; exactly 52,428,800 bytes (50 MiB) fails** with `HTTP 500 {"success":false,"error":"File too
large"}`. This is the endpoint to use for anything that isn't an image — PDFs, audio, video, notes.

**`POST /api/upload/image`** — base64-encoded JSON body. Accepted up to 200 MB raw (266.7 MB as base64
JSON) with no upper limit found — but **it silently corrupts any non-image payload.** A base64-encoded
PDF returned `HTTP 200 {"success":true}` with a plausible `.png` URL, but the stored bytes did not match
what was sent (md5 mismatch, garbled/bit-shifted, zero common prefix with the original). A real
`image/png` payload round-trips correctly. **Only use this endpoint for actual images. For PDFs/audio/
video, use `/api/upload/file` instead, even though this endpoint's name and lack of any error might
suggest it's fine for arbitrary files — it is not.** See `BACKLOG.md #11`.

```javascript
// Correct pattern for a non-image file (e.g. PDF ingestion source):
const form = new FormData();
form.append('file', pdfBlob, 'report.pdf');
const response = await fetch('https://api.audos.com/api/upload/file', {
  method: 'POST',
  body: form
});
const { success, url, bytes } = await response.json();
```

### Scheduler Integration

Create recurring tasks and delayed jobs.

> **Verified live 2026-07-14, retested 2026-07-16 — reliability is genuinely inconsistent, not simply
> broken. Read both dates before relying on this for anything time-critical.**
>
> 1. **One-time hook scheduling doesn't exist** (confirmed both dates). `POST /schedules` with a one-time
>    payload (any of `scheduleType: "one-time"`/`"one_time"`, a bare `scheduledAt`, or a `runAt` field) is
>    rejected: `{"error":"Invalid parameters: [...{\"expected\":\"'hourly' | 'daily' | 'weekly' |
>    'monthly'\",\"received\":\"undefined\",\"path\":[\"frequency\"]}]"}`. `frequency` is **required** —
>    only recurring hourly/daily/weekly/monthly schedules can be created for hooks. One-time scheduling
>    only exists for the separate `/schedules/email` endpoint shown further below.
> 2. **2026-07-14: an hourly schedule never fired.** Created with a valid RRULE, `status: pending` —
>    confirmed still `pending`, `runCount: 0`, `lastError: null` after 2h07m and 20m past two separate
>    `nextRun` times, on two schedules. The target hook worked fine when called manually — the dispatcher
>    simply never invoked it.
> 3. **2026-07-16: a daily schedule fired correctly.** Same mechanism, different frequency — fired within
>    ~9 seconds of `nextRun`, the hook executed and wrote to the database, `runCount` advanced `0→1`,
>    `lastRun`/`nextRun` updated correctly. The complete schedule row was inspected for any hidden
>    status/error field beyond `runCount`/`lastError` that might explain the earlier non-fire — none
>    exists (`id, workspaceId, name, description, actionType, actionPayload, scheduleType, rrule,
>    timezone, status, nextRun, lastRun, runCount, lastError, createdAt, updatedAt` is the complete field
>    list; no dispatcher log, queue position, or disabled flag). Asked Otto directly whether the earlier
>    non-fire is a known platform issue — verbatim answer: *"I do not have any documented or confirmed
>    knowledge that scheduled/recurring hooks failing to fire is a known, tracked platform issue... I am
>    not going to represent this as either 'known and broken' or 'working as intended.'"*
>
> **Practical consequence, revised:** the scheduler is not reliably broken *or* reliably working — treat
> it as genuinely inconsistent until Audos confirms a root cause (possibly cadence-specific — only hourly
> has ever failed, only daily has ever succeeded — but that's one data point each, not a pattern). **This
> stopped being a hard blocker on 2026-07-16**, though: client-orchestrated chaining (the calling app
> making several sequential hook calls itself, pacing the work, rather than trusting the platform's own
> scheduler) was tested as a substitute and found completely clean — 5/5 sequential calls succeeded,
> ~1–1.8s each, zero rate-limiting, zero auth friction, zero concurrency errors. Any job that can tolerate
> the *client* staying alive for its duration has a proven, reliable path that doesn't depend on the
> scheduler at all. Only a true fire-and-forget/unattended background job still needs the scheduler to be
> trustworthy, which it isn't yet, confirmed. See `BACKLOG.md #12`.

**Recurring Task (Daily):**

```javascript
const response = await fetch(
  'https://api.audos.com/api/workspaces/YOUR_WORKSPACE_ID/schedules',
  {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: 'Daily Report',
      frequency: 'daily',
      time: '09:00',
      timezone: 'America/New_York',
      actionType: 'hook',
      hookName: 'generate-report'
    })
  }
);
```

**One-Time Delayed Task:**

```javascript
const futureDate = new Date(Date.now() + 3600000); // 1 hour from now

const response = await fetch(
  'https://api.audos.com/api/workspaces/YOUR_WORKSPACE_ID/schedules/email',
  {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      name: 'Follow-up Email',
      scheduledAt: futureDate.toISOString(),
      email: {
        to: 'user@example.com',
        subject: 'Follow-up',
        text: 'Just checking in!'
      }
    })
  }
);
```

---

## Helper Functions

Copy-paste ready utility functions for common tasks.

### Build Query String

Convert an object to URL query string.

```javascript
/**
 * Build URL query string from object
 * @param {object} params - Key-value pairs
 * @returns {string} Query string (without leading ?)
 */
function buildQueryString(params) {
  return Object.entries(params)
    .filter(([_, val]) => val !== null && val !== undefined)
    .map(([key, val]) => `${encodeURIComponent(key)}=${encodeURIComponent(val)}`)
    .join('&');
}

// Usage
const params = { name: 'John Doe', age: 30, city: 'New York' };
const qs = buildQueryString(params);
console.log(qs); // "name=John%20Doe&age=30&city=New%20York"

const url = `https://api.example.com/search?${qs}`;
```

### Parse Query String

Parse URL query string to object.

```javascript
/**
 * Parse query string to object
 * @param {string} queryString - Query string (with or without leading ?)
 * @returns {object} Parsed parameters
 */
function parseQueryString(queryString) {
  const params = {};
  const cleaned = queryString.startsWith('?') ? queryString.slice(1) : queryString;

  if (!cleaned) return params;

  cleaned.split('&').forEach(pair => {
    const [key, value] = pair.split('=');
    params[decodeURIComponent(key)] = decodeURIComponent(value || '');
  });

  return params;
}

// Usage
const qs = '?name=John%20Doe&age=30&city=New%20York';
const params = parseQueryString(qs);
console.log(params); // { name: 'John Doe', age: '30', city: 'New York' }
```

### Safe JSON Parse

Safely parse JSON with error handling.

```javascript
/**
 * Safely parse JSON string
 * @param {string} jsonString - JSON string to parse
 * @param {any} defaultValue - Value to return on error
 * @returns {any} Parsed object or default value
 */
function safeJsonParse(jsonString, defaultValue = null) {
  try {
    return JSON.parse(jsonString);
  } catch (error) {
    console.error('JSON parse error:', error);
    return defaultValue;
  }
}

// Usage
const data = safeJsonParse(request.body.data, {});
const config = safeJsonParse(user.settings, { theme: 'light' });
```

### Format Date

Format date as ISO string or custom format.

```javascript
/**
 * Format date to ISO string
 * @param {Date|string|number} date - Date to format
 * @returns {string} ISO 8601 formatted string
 */
function formatDate(date) {
  const d = date instanceof Date ? date : new Date(date);
  return d.toISOString();
}

/**
 * Format date for display
 * @param {Date|string|number} date - Date to format
 * @returns {string} Formatted date (YYYY-MM-DD)
 */
function formatDateDisplay(date) {
  const d = date instanceof Date ? date : new Date(date);
  return d.toISOString().split('T')[0];
}

/**
 * Format timestamp for database
 * @param {Date|string|number} date - Date to format
 * @returns {string} PostgreSQL timestamp format
 */
function formatTimestamp(date) {
  const d = date instanceof Date ? date : new Date(date);
  return d.toISOString().replace('T', ' ').replace('Z', '');
}

// Usage
const now = formatDate(new Date()); // "2026-03-31T12:00:00.000Z"
const today = formatDateDisplay(new Date()); // "2026-03-31"
const ts = formatTimestamp(new Date()); // "2026-03-31 12:00:00.000"
```

### Validate Email

Check if a string is a valid email address.

```javascript
/**
 * Validate email address format
 * @param {string} email - Email to validate
 * @returns {boolean} True if valid
 */
function isValidEmail(email) {
  if (!email || typeof email !== 'string') return false;

  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

// Usage
if (!isValidEmail(request.body.email)) {
  respond(400, { error: 'Invalid email address' });
  return;
}
```

### Paginate Array

Paginate an array of results.

```javascript
/**
 * Paginate array
 * @param {Array} array - Array to paginate
 * @param {number} page - Page number (1-based)
 * @param {number} pageSize - Items per page
 * @returns {object} Paginated result
 */
function paginate(array, page = 1, pageSize = 10) {
  const start = (page - 1) * pageSize;
  const end = start + pageSize;

  return {
    data: array.slice(start, end),
    page,
    pageSize,
    total: array.length,
    totalPages: Math.ceil(array.length / pageSize),
    hasNext: end < array.length,
    hasPrev: page > 1
  };
}

// Usage
const users = await db.query('SELECT * FROM users');
const page = parseInt(request.query.page || '1');
const result = paginate(users, page, 20);

respond(200, result);
```

### Retry with Backoff

Retry failed operations with exponential backoff.

```javascript
/**
 * Retry function with exponential backoff
 * @param {Function} fn - Async function to retry
 * @param {number} maxRetries - Maximum retry attempts
 * @param {number} delay - Initial delay in ms
 * @returns {Promise<any>} Result of function
 */
async function retryWithBackoff(fn, maxRetries = 3, delay = 1000) {
  let lastError;

  for (let i = 0; i < maxRetries; i++) {
    try {
      return await fn();
    } catch (error) {
      lastError = error;
      console.error(`Attempt ${i + 1} failed:`, error);

      if (i < maxRetries - 1) {
        const waitTime = delay * Math.pow(2, i);
        console.log(`Retrying in ${waitTime}ms...`);
        await new Promise(resolve => setTimeout(resolve, waitTime));
      }
    }
  }

  throw lastError;
}

// Usage
try {
  const result = await retryWithBackoff(async () => {
    const response = await fetch('https://unreliable-api.com/data');
    if (!response.ok) throw new Error('API request failed');
    return await response.json();
  }, 3, 500);

  respond(200, result);
} catch (error) {
  respond(500, { error: 'Failed after retries' });
}
```

### Batch Operations

Process items in batches to avoid memory issues.

```javascript
/**
 * Process items in batches
 * @param {Array} items - Items to process
 * @param {number} batchSize - Batch size
 * @param {Function} processFn - Async function to process each batch
 * @returns {Promise<Array>} All results
 */
async function batchProcess(items, batchSize, processFn) {
  const results = [];

  for (let i = 0; i < items.length; i += batchSize) {
    const batch = items.slice(i, i + batchSize);
    const batchResults = await processFn(batch);
    results.push(...batchResults);
  }

  return results;
}

// Usage
const emails = ['user1@example.com', 'user2@example.com', /* ... 1000 more */];

await batchProcess(emails, 50, async (batch) => {
  // Insert 50 contacts at a time
  return await db.insert('contacts', batch.map(email => ({ email })));
});
```

---

## Complete Examples

### Example 1: User Registration Endpoint

```javascript
// POST /api/hooks/execute/workspace-{workspaceId}/user-registration
// Handle user registration with validation and email

// Validate request
if (request.method !== 'POST') {
  respond(405, { error: 'Method not allowed' });
  return;
}

const { name, email, password } = request.body;

// Validate inputs
if (!name || !email || !password) {
  respond(400, { error: 'Missing required fields' });
  return;
}

if (!isValidEmail(email)) {
  respond(400, { error: 'Invalid email address' });
  return;
}

// Check if user exists
const existing = await db.query(
  `SELECT id FROM users WHERE email = '${email}'`
);

if (existing.length > 0) {
  respond(409, { error: 'User already exists' });
  return;
}

// Create user
const newUser = await db.insert('users', {
  name,
  email,
  password_hash: password, // In production, hash this!
  status: 'active'
});

// Send welcome email
await platform.sendEmail({
  to: email,
  subject: 'Welcome to our platform!',
  text: `Hi ${name}, welcome aboard!`,
  html: `<h1>Welcome ${name}!</h1><p>Thanks for signing up.</p>`
});

// Track in CRM
await fetch(
  'https://api.audos.com/api/crm/contacts/YOUR_WORKSPACE_ID',
  {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      email,
      name,
      source: 'registration'
    })
  }
);

respond(201, {
  success: true,
  user: {
    id: newUser[0].id,
    name: newUser[0].name,
    email: newUser[0].email
  }
});

function isValidEmail(email) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);
}
```

### Example 2: Daily Report Generator

```javascript
// GET /api/hooks/execute/workspace-{workspaceId}/daily-report
// Generate and email daily analytics report

const today = new Date();
const yesterday = new Date(today);
yesterday.setDate(yesterday.getDate() - 1);

// Get analytics
const metricsResponse = await fetch(
  'https://api.audos.com/api/funnel/metrics/YOUR_WORKSPACE_ID?days=1'
);
const metrics = await metricsResponse.json();

// Get new contacts
const contactsResponse = await fetch(
  'https://api.audos.com/api/crm/contacts/YOUR_WORKSPACE_ID?days=1'
);
const contacts = await contactsResponse.json();

// Get orders from database
const orders = await db.query(`
  SELECT COUNT(*) as count, SUM(total) as revenue
  FROM orders
  WHERE created_at >= '${yesterday.toISOString()}'
`);

// Generate report
const reportHtml = `
  <h1>Daily Report - ${today.toISOString().split('T')[0]}</h1>

  <h2>Visitors</h2>
  <p>Total visitors: ${metrics.visitors}</p>
  <p>New contacts: ${contacts.length}</p>

  <h2>Revenue</h2>
  <p>Orders: ${orders[0].count}</p>
  <p>Revenue: $${orders[0].revenue || 0}</p>

  <h2>Conversion</h2>
  <p>Conversion rate: ${metrics.conversionRate}%</p>
`;

// Send report
await platform.sendEmail({
  to: 'admin@yourcompany.com',
  subject: `Daily Report - ${today.toISOString().split('T')[0]}`,
  text: 'See HTML version',
  html: reportHtml
});

respond(200, { success: true, message: 'Report sent' });
```

### Example 3: Webhook Processor

```javascript
// POST /api/hooks/execute/workspace-{workspaceId}/webhook-processor
// Process incoming webhooks from external services

console.log('Webhook received:', request.method, request.body);

// Verify webhook signature (example for Stripe)
const signature = request.headers['stripe-signature'];
// In production, verify this signature

const event = request.body;

switch (event.type) {
  case 'payment_intent.succeeded':
    const payment = event.data.object;

    // Record payment in database
    await db.insert('payments', {
      stripe_id: payment.id,
      amount: payment.amount / 100,
      currency: payment.currency,
      customer_email: payment.receipt_email,
      status: 'completed'
    });

    // Send confirmation email
    await platform.sendEmail({
      to: payment.receipt_email,
      subject: 'Payment Received',
      text: `Your payment of $${payment.amount / 100} was successful.`
    });

    break;

  case 'customer.subscription.created':
    const subscription = event.data.object;

    // Update user subscription status
    await db.update(
      'users',
      { subscription_status: 'active', subscription_id: subscription.id },
      { email: subscription.customer_email }
    );

    break;

  default:
    console.log('Unhandled event type:', event.type);
}

respond(200, { received: true });
```

---

## Troubleshooting

### Common Errors

**1. "URLSearchParams is not defined"**

Solution: Use the `buildQueryString()` helper function instead.

**2. "Buffer is not defined"**

Solution: Work with base64 strings directly or use alternative encoding.

**3. "Execution timeout"**

Solution: Optimize queries, reduce external API calls, or batch operations.

**4. "Cannot read property 'get' of undefined"**

Solution: Access response data directly with `response.json()` or `response.text()`.

**5. "Database query failed"**

Solution: Check SQL syntax, ensure table exists, verify column names.

### Debugging Tips

```javascript
// Log everything for debugging
console.log('Request:', JSON.stringify(request, null, 2));
console.log('Method:', request.method);
console.log('Body:', request.body);
console.log('Query:', request.query);

// Wrap in try-catch
try {
  const result = await db.query('SELECT * FROM users');
  console.log('Query successful:', result.length, 'rows');
} catch (error) {
  console.error('Database error:', error);
  respond(500, { error: error.message || 'Database query failed' });
  return;
}

// Test with simple responses
respond(200, { debug: true, received: request.body });
```

### Performance Optimization

```javascript
// ✅ GOOD - Single query with JOIN
const data = await db.query(`
  SELECT orders.*, users.name, users.email
  FROM orders
  JOIN users ON orders.user_id = users.id
`);

// ❌ BAD - Multiple queries in loop
const orders = await db.query('SELECT * FROM orders');
for (const order of orders) {
  const user = await db.query(`SELECT * FROM users WHERE id = ${order.user_id}`);
  // This is very slow!
}

// ✅ GOOD - Batch inserts
await db.insert('contacts', [
  { email: 'user1@example.com' },
  { email: 'user2@example.com' },
  { email: 'user3@example.com' }
]);

// ❌ BAD - Individual inserts
await db.insert('contacts', { email: 'user1@example.com' });
await db.insert('contacts', { email: 'user2@example.com' });
await db.insert('contacts', { email: 'user3@example.com' });
```

---

## Additional Resources

- **Server Functions:** Create and test at `/api/hooks/execute/workspace-{workspaceId}/{hookName}`
- **Database:** Use the database tools to create tables and query data
- **Scheduler:** Schedule recurring tasks and delayed jobs
- **Email:** Send transactional emails via `platform.sendEmail()`
- **Analytics:** Access visitor and conversion metrics via internal APIs

---

**Version:** 1.0
**Last Updated:** 2026-03-31
**Platform:** Audos
