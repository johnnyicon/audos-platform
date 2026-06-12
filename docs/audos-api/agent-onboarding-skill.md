---
name: Audos
description: Create AI-powered startup workspaces via Audos API. Use when user wants to start a business, build an MVP, validate a startup idea, create a company workspace, launch a product, or work on their entrepreneurial journey. Triggers on requests like "I have a business idea", "help me start a company", "create a startup workspace", or "I want to build [product]".
---

# Audos Workspace Builder (API v1.2)

Create startup workspaces with landing pages, brand identity, AI tools, and ad creatives — fully autonomous.

## Base URL

```
https://audos.com/api/agent/onboard
```

## URL Construction

The API returns URLs using the current deployment domain:

```json
"urls": {
  "landingPage": "https://audos.com/site/184582",
  "workspace": "https://audos.com/space/workspace-184582"
}
```

Use these URLs directly — no domain swapping needed.

## Quick Reference

| Action | Method | Endpoint |
|--------|--------|----------|
| API docs | GET | / |
| Start onboarding | POST | /start |
| Verify OTP | POST | /verify |
| Check build status | GET | /status/:workspaceId |
| Check build status (alt) | POST | /status |
| Chat with Otto | POST | /chat |
| Chat with Otto | POST | /chat/:workspaceId |
| Rebuild (if failed) | POST | /rebuild/:workspaceId |

## Authentication

- **Token format:** `aud_live_xxxx` (48 hex chars after prefix)
- **Auth tokens never expire** — store persistently by email
- **Session tokens** expire in 30 min (only needed during OTP flow)
- **Preferred:** Bearer token in `Authorization` header
- **Alternative:** `authToken` or `sessionToken` in request body

## Conversation Flow

### New Users Flow
1. Collect user's email + business idea
2. Start → POST /start (sends 4-digit OTP to email)
3. Verify → POST /verify with OTP code → returns authToken + starts build
4. Monitor → GET /status/:workspaceId every 15-30s
5. Watch for landingPageReady: true (~10 min)
6. Chat → POST /chat/:workspaceId

### Returning Users (have workspace)
1. Start → POST /start with email
2. Response includes auth_token + urls directly — skip OTP
3. Chat → POST /chat/:workspaceId immediately

## API Reference

### POST /start
Fields: email (required), businessIdea (required, min 10 chars), businessName (optional), targetCustomer (optional), callbackUrl (optional), createNew (optional, force new workspace)

Returns new user: sessionToken
Returns returning user: auth_token, workspaceId, urls

### POST /verify
Body: { sessionToken, otpCode }
Returns: workspaceId, authToken, urls, buildInfo

### GET /status/:workspaceId
Header: Authorization: Bearer <authToken>
Key fields: landingPageReady (bool), status, progress (0-100), estimatedTimeRemaining, completedSteps, parallelBuildStatus

### POST /chat/:workspaceId
Header: Authorization: Bearer <authToken>
Body: { message }
Returns: workspaceId, chatId, response

### POST /chat
Body: { authToken, message }
Returns: workspaceId, chatId, response

### POST /rebuild/:workspaceId
Header: Authorization: Bearer <authToken>

## Error Codes
AUTH_TOKEN_INVALID (401), WORKSPACE_NOT_FOUND (404), OTP_EXPIRED (401), OTP_INVALID (401), VALIDATION_ERROR (400), CHAT_FAILED (502)

## Notes
- Auth tokens never expire — store by email
- Tokens are workspace-scoped (one token per workspace)
- Poll status every 15-30s during build
- landingPageReady is the most reliable completion signal
