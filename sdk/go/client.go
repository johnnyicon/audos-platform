package audos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Config struct {
	WorkspaceID string
	APIKey      string
	BaseURL     string
}

type Client struct {
	config Config
	http   *http.Client
}

func NewClient(config Config) *Client {
	return &Client{config: config, http: &http.Client{}}
}

func (c *Client) callHook(hookName string, body any) (json.RawMessage, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	url := fmt.Sprintf("%s/api/hooks/execute/workspace-%s/%s", c.config.BaseURL, c.config.WorkspaceID, hookName)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if hookName == "db-api" {
		req.Header.Set("x-api-key", c.config.APIKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %d %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	return json.RawMessage(data), nil
}

// ─── Types ────────────────────────────────────────────────────────────────────

type FilterOperator string

const (
	OpEq      FilterOperator = "eq"
	OpNeq     FilterOperator = "neq"
	OpGt      FilterOperator = "gt"
	OpGte     FilterOperator = "gte"
	OpLt      FilterOperator = "lt"
	OpLte     FilterOperator = "lte"
	OpLike    FilterOperator = "like"
	OpILike   FilterOperator = "ilike"
	OpIn      FilterOperator = "in"
	OpIsNull  FilterOperator = "is_null"
	OpNotNull FilterOperator = "not_null"
)

type Filter struct {
	Column   string         `json:"column"`
	Operator FilterOperator `json:"operator"`
	Value    any            `json:"value,omitempty"`
}

type OrderBy struct {
	Column    string `json:"column"`
	Direction string `json:"direction,omitempty"` // "asc" | "desc"
}

type QueryOptions struct {
	Filters []Filter `json:"filters,omitempty"`
	OrderBy *OrderBy `json:"orderBy,omitempty"`
	Limit   int      `json:"limit,omitempty"`
	Offset  int      `json:"offset,omitempty"`
	Columns []string `json:"columns,omitempty"`
}

type EmailOptions struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html,omitempty"`
	ReplyTo string `json:"replyTo,omitempty"`
}

type ContactInput struct {
	Email  string `json:"email"`
	Name   string `json:"name,omitempty"`
	Source string `json:"source,omitempty"`
}

// ─── Database ─────────────────────────────────────────────────────────────────

func (c *Client) DBListTables() (json.RawMessage, error) {
	return c.callHook("db-api", map[string]any{"action": "list-tables"})
}

func (c *Client) DBDescribe(table string) (json.RawMessage, error) {
	return c.callHook("db-api", map[string]any{"action": "describe", "table": table})
}

func (c *Client) DBQuery(table string, opts QueryOptions) (json.RawMessage, error) {
	body := map[string]any{"action": "query", "table": table}
	if len(opts.Filters) > 0 {
		body["filters"] = opts.Filters
	}
	if opts.OrderBy != nil {
		body["orderBy"] = opts.OrderBy
	}
	if opts.Limit > 0 {
		body["limit"] = opts.Limit
	}
	if opts.Offset > 0 {
		body["offset"] = opts.Offset
	}
	if len(opts.Columns) > 0 {
		body["columns"] = opts.Columns
	}
	return c.callHook("db-api", body)
}

func (c *Client) DBRawQuery(sql string, params ...any) (json.RawMessage, error) {
	return c.callHook("db-api", map[string]any{"action": "raw-query", "sql": sql, "params": params})
}

func (c *Client) DBInsert(table string, data any) (json.RawMessage, error) {
	return c.callHook("db-api", map[string]any{"action": "insert", "table": table, "data": data})
}

func (c *Client) DBUpdate(table string, filters []Filter, data map[string]any) (json.RawMessage, error) {
	return c.callHook("db-api", map[string]any{"action": "update", "table": table, "filters": filters, "data": data})
}

func (c *Client) DBDelete(table string, filters []Filter) (json.RawMessage, error) {
	return c.callHook("db-api", map[string]any{"action": "delete", "table": table, "filters": filters})
}

// ─── AI ───────────────────────────────────────────────────────────────────────

func (c *Client) AIGenerate(prompt, systemPrompt string) (json.RawMessage, error) {
	body := map[string]any{"action": "generate", "prompt": prompt}
	if systemPrompt != "" {
		body["systemPrompt"] = systemPrompt
	}
	return c.callHook("ai-api", body)
}

// ─── Email ────────────────────────────────────────────────────────────────────

func (c *Client) EmailSend(opts EmailOptions) (json.RawMessage, error) {
	body := map[string]any{
		"action":  "send",
		"to":      opts.To,
		"subject": opts.Subject,
		"text":    opts.Text,
	}
	if opts.HTML != "" {
		body["html"] = opts.HTML
	}
	if opts.ReplyTo != "" {
		body["replyTo"] = opts.ReplyTo
	}
	return c.callHook("email-api", body)
}

// ─── Web ──────────────────────────────────────────────────────────────────────

func (c *Client) WebFetch(url string) (json.RawMessage, error) {
	return c.callHook("web-api", map[string]any{"action": "fetch", "url": url})
}

func (c *Client) WebMetadata(url string) (json.RawMessage, error) {
	return c.callHook("web-api", map[string]any{"action": "metadata", "url": url})
}

// ─── Storage ──────────────────────────────────────────────────────────────────

func (c *Client) StorageList(category string) (json.RawMessage, error) {
	body := map[string]any{"action": "list"}
	if category != "" {
		body["category"] = category
	}
	return c.callHook("storage-api", body)
}

// ─── Analytics ────────────────────────────────────────────────────────────────

func (c *Client) AnalyticsOverview(days int) (json.RawMessage, error) {
	if days == 0 {
		days = 30
	}
	return c.callHook("analytics-api", map[string]any{"action": "overview", "days": days})
}

// ─── CRM ──────────────────────────────────────────────────────────────────────

func (c *Client) CRMListContacts(limit int) (json.RawMessage, error) {
	if limit == 0 {
		limit = 50
	}
	return c.callHook("crm-api", map[string]any{"action": "list-contacts", "limit": limit})
}

func (c *Client) CRMCreateContact(contact ContactInput) (json.RawMessage, error) {
	return c.callHook("crm-api", map[string]any{
		"action": "create-contact",
		"email":  contact.Email,
		"name":   contact.Name,
		"source": contact.Source,
	})
}
