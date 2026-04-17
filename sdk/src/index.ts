export interface AudosClientConfig {
  workspaceId: string;
  apiKey: string;
  baseUrl: string;
}

export type FilterOperator =
  | 'eq' | 'neq' | 'gt' | 'gte' | 'lt' | 'lte'
  | 'like' | 'ilike' | 'in' | 'is_null' | 'not_null';

export interface Filter {
  column: string;
  operator: FilterOperator;
  value?: unknown;
}

export interface QueryOptions {
  filters?: Filter[];
  orderBy?: { column: string; direction?: 'asc' | 'desc' };
  limit?: number;
  offset?: number;
  columns?: string[];
}

export function createClient(config: AudosClientConfig) {
  async function callHook(hookName: string, body: unknown): Promise<any> {
    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (hookName === 'db-api') headers['x-api-key'] = config.apiKey;
    const response = await fetch(
      `${config.baseUrl}/api/hooks/execute/workspace-${config.workspaceId}/${hookName}`,
      { method: 'POST', headers, body: JSON.stringify(body) }
    );
    if (!response.ok) throw new Error(`API error: ${response.status} ${response.statusText}`);
    return response.json();
  }

  const db = {
    async listTables() {
      return callHook('db-api', { action: 'list-tables' });
    },

    async describe(table: string) {
      return callHook('db-api', { action: 'describe', table });
    },

    async query<T = any>(table: string, options: QueryOptions = {}): Promise<{ success: boolean; data: { rows: T[]; rowCount: number } }> {
      return callHook('db-api', { action: 'query', table, ...options });
    },

    async rawQuery<T = any>(sql: string, params?: unknown[]): Promise<{ success: boolean; data: { rows: T[]; rowCount: number } }> {
      return callHook('db-api', { action: 'raw-query', sql, params });
    },

    async insert(table: string, data: Record<string, unknown> | Record<string, unknown>[]) {
      return callHook('db-api', { action: 'insert', table, data });
    },

    async update(table: string, filters: Filter[], data: Record<string, unknown>) {
      return callHook('db-api', { action: 'update', table, filters, data });
    },

    async delete(table: string, filters: Filter[]) {
      return callHook('db-api', { action: 'delete', table, filters });
    },
  };

  const ai = {
    async generate(prompt: string, systemPrompt?: string): Promise<{ success: boolean; text: string }> {
      return callHook('ai-api', { action: 'generate', prompt, systemPrompt });
    },
  };

  const email = {
    async send(options: { to: string; subject: string; text: string; html?: string; replyTo?: string }) {
      return callHook('email-api', { action: 'send', ...options });
    },
  };

  const web = {
    async fetch(url: string): Promise<{ success: boolean; content: string; title: string }> {
      return callHook('web-api', { action: 'fetch', url });
    },
    async metadata(url: string) {
      return callHook('web-api', { action: 'metadata', url });
    },
  };

  const storage = {
    async list(category?: string) {
      return callHook('storage-api', { action: 'list', category });
    },
  };

  const analytics = {
    async overview(days: number = 30) {
      return callHook('analytics-api', { action: 'overview', days });
    },
  };

  const crm = {
    async listContacts(limit: number = 50) {
      return callHook('crm-api', { action: 'list-contacts', limit });
    },
    async createContact(contact: { email: string; name?: string; source?: string }) {
      return callHook('crm-api', { action: 'create-contact', ...contact });
    },
  };

  return { db, ai, email, web, storage, analytics, crm, callHook };
}
