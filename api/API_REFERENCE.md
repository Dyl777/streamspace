# StreamSpace API Reference

**Version**: v1
**Base URL**: `/api/v1`
**Authentication**: JWT Bearer Token (except auth endpoints)

All requests include a `X-Request-ID` header for distributed tracing.

---

## Table of Contents

- [Authentication](#authentication)
- [Sessions](#sessions)
- [Templates](#templates)
- [Users](#users)
- [Groups](#groups)
- [Plugins](#plugins)
- [Catalog](#catalog)
- [Activity Tracking](#activity-tracking)
- [Sharing](#sharing)
- [Kubernetes Resources](#kubernetes-resources)
- [System](#system)

---

## Authentication

### POST /api/v1/auth/login

Local username/password authentication.

**Request Body**:
```json
{
  "username": "string",
  "password": "string"
}
```

**Response** (200 OK):
```json
{
  "token": "jwt-token-string",
  "expiresAt": "2025-01-15T12:00:00Z",
  "user": {
    "id": "user-id",
    "username": "user1",
    "email": "user@example.com",
    "fullName": "User Name",
    "role": "user",
    "active": true
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Invalid credentials
- `403 Forbidden`: Account disabled

---

### POST /api/v1/auth/refresh

Refresh an expiring JWT token.

**Request Body**:
```json
{
  "token": "current-jwt-token"
}
```

**Response** (200 OK):
```json
{
  "token": "new-jwt-token",
  "expiresAt": "2025-01-16T12:00:00Z",
  "user": { ... }
}
```

**Errors**:
- `401 Unauthorized`: Token invalid or not eligible for refresh

---

### GET /api/v1/auth/saml/login

Initiate SAML SSO authentication flow.

**Query Parameters**:
- `return_url` (optional): URL to redirect after authentication (default: `/`)

**Response**: Redirects to SAML Identity Provider

**Errors**:
- `503 Service Unavailable`: SAML not configured

---

### POST /api/v1/auth/saml/acs

SAML Assertion Consumer Service (callback endpoint).

**Response** (200 OK):
```json
{
  "token": "jwt-token",
  "expiresAt": "2025-01-15T12:00:00Z",
  "user": { ... },
  "returnUrl": "/"
}
```

**Errors**:
- `400 Bad Request`: Missing required SAML attributes
- `401 Unauthorized`: No SAML assertion
- `403 Forbidden`: Account disabled

---

### GET /api/v1/auth/saml/metadata

Returns SAML Service Provider metadata XML for IdP configuration.

**Response** (200 OK):
```xml
<?xml version="1.0"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" ...>
  ...
</EntityDescriptor>
```

**Headers**:
- `Content-Type: application/samlmetadata+xml`

---

## Sessions

### GET /api/v1/sessions

List all sessions (admin/operator) or user's own sessions.

**Query Parameters**:
- `user` (optional): Filter by username
- `template` (optional): Filter by template name
- `state` (optional): Filter by state (running, hibernated, terminated)

**Response** (200 OK):
```json
[
  {
    "id": "session-id",
    "name": "user1-firefox",
    "user": "user1",
    "template": "firefox-browser",
    "state": "running",
    "url": "https://user1-firefox.streamspace.local",
    "createdAt": "2025-01-15T10:00:00Z",
    "lastActivity": "2025-01-15T11:30:00Z"
  }
]
```

---

### POST /api/v1/sessions

Create a new session from a template.

**Request Body**:
```json
{
  "template": "firefox-browser",
  "resources": {
    "memory": "2Gi",
    "cpu": "1000m"
  },
  "persistentHome": true,
  "idleTimeout": "30m"
}
```

**Response** (201 Created):
```json
{
  "id": "session-id",
  "name": "user1-firefox-abc123",
  "state": "pending",
  ...
}
```

---

### GET /api/v1/sessions/:id

Get session details.

**Response** (200 OK):
```json
{
  "id": "session-id",
  "name": "user1-firefox",
  "user": "user1",
  "template": "firefox-browser",
  "state": "running",
  "url": "https://user1-firefox.streamspace.local",
  "podName": "ss-user1-firefox-abc123",
  "resourceUsage": {
    "memory": "1.2Gi",
    "cpu": "450m"
  },
  "createdAt": "2025-01-15T10:00:00Z",
  "lastActivity": "2025-01-15T11:30:00Z"
}
```

---

### PATCH /api/v1/sessions/:id

Update session state (hibernate, wake, terminate).

**Request Body**:
```json
{
  "state": "hibernated"
}
```

**Response** (200 OK):
```json
{
  "id": "session-id",
  "state": "hibernated",
  ...
}
```

---

### DELETE /api/v1/sessions/:id

Terminate and delete a session.

**Response** (204 No Content)

---

## Templates

### GET /api/v1/templates

List all available templates.

**Query Parameters**:
- `category` (optional): Filter by category

**Response** (200 OK):
```json
[
  {
    "name": "firefox-browser",
    "displayName": "Firefox Web Browser",
    "description": "Modern, privacy-focused web browser",
    "category": "Web Browsers",
    "icon": "https://...",
    "defaultResources": {
      "memory": "2Gi",
      "cpu": "1000m"
    },
    "capabilities": ["Network", "Audio", "Clipboard"],
    "tags": ["browser", "web", "privacy"]
  }
]
```

---

### GET /api/v1/templates/:name

Get template details.

**Response** (200 OK):
```json
{
  "name": "firefox-browser",
  "displayName": "Firefox Web Browser",
  "spec": {
    "baseImage": "lscr.io/linuxserver/firefox:latest",
    "ports": [...],
    "env": [...],
    ...
  }
}
```

---

### PUT /api/v1/templates/:name

Update template configuration (admin only).

**Request Body**:
```json
{
  "displayName": "Updated Name",
  "description": "Updated description",
  "defaultResources": {
    "memory": "4Gi",
    "cpu": "2000m"
  }
}
```

**Response** (200 OK)

---

## Kubernetes Resources

### POST /api/v1/resources

Create a generic Kubernetes resource (admin only).

**Request Body**:
```json
{
  "apiVersion": "v1",
  "kind": "ConfigMap",
  "metadata": {
    "name": "my-config",
    "namespace": "streamspace"
  },
  "data": {
    "key": "value"
  }
}
```

**Response** (201 Created):
```json
{
  "apiVersion": "v1",
  "kind": "ConfigMap",
  "metadata": {...},
  "data": {...}
}
```

---

### PUT /api/v1/resources/:type/:name

Update a Kubernetes resource (admin only).

**Query Parameters**:
- `namespace` (optional): Target namespace

**Request Body**: Full resource definition

**Response** (200 OK)

---

### DELETE /api/v1/resources/:type/:name

Delete a Kubernetes resource (admin only).

**Query Parameters**:
- `apiVersion` (required): e.g., "apps/v1"
- `kind` (required): e.g., "Deployment"
- `namespace` (optional): Target namespace

**Response** (200 OK):
```json
{
  "message": "Resource deleted successfully",
  "name": "resource-name",
  "type": "deployment"
}
```

---

## Catalog

### GET /api/v1/catalog

List available catalog templates.

**Response** (200 OK):
```json
[
  {
    "id": "catalog-1",
    "name": "Firefox ESR",
    "description": "Extended Support Release",
    "category": "Browsers",
    "manifest": "...",
    "repository": {
      "id": "repo-1",
      "name": "LinuxServer.io",
      "url": "https://..."
    }
  }
]
```

---

### POST /api/v1/catalog/:id/install

Install a template from the catalog (admin only).

**Response** (201 Created):
```json
{
  "template": "firefox-esr",
  "status": "installed"
}
```

---

## Plugins

### GET /api/v1/plugins

List installed plugins.

**Response** (200 OK):
```json
[
  {
    "id": "plugin-1",
    "name": "backup-plugin",
    "version": "1.0.0",
    "enabled": true,
    "description": "Automated backup plugin"
  }
]
```

---

### POST /api/v1/plugins

Install a new plugin (admin only).

**Request Body**:
```json
{
  "name": "backup-plugin",
  "version": "1.0.0",
  "config": {...}
}
```

---

## System

### GET /api/v1/health

Health check endpoint.

**Response** (200 OK):
```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T12:00:00Z",
  "version": "v0.1.0",
  "checks": {
    "database": "ok",
    "kubernetes": "ok",
    "redis": "ok"
  }
}
```

---

### GET /api/v1/metrics

Prometheus metrics endpoint.

**Response** (200 OK): Prometheus text format

---

## Error Responses

All errors follow this format:

```json
{
  "error": "Short error message",
  "message": "Detailed explanation",
  "requestId": "uuid-request-id"
}
```

**Common Status Codes**:
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Authentication required or failed
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `408 Request Timeout`: Request took too long
- `409 Conflict`: Resource conflict (e.g., duplicate name)
- `422 Unprocessable Entity`: Validation failed
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error
- `503 Service Unavailable`: Service temporarily unavailable

---

## Rate Limiting

**IP-based**: 100 requests/second per IP (burst: 200)
**User-based**: 1000 requests/hour per authenticated user (burst: 50)
**Auth endpoints**: 5 requests/second (burst: 10)

**Headers**:
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Requests remaining
- `X-RateLimit-Reset`: Unix timestamp when limit resets

---

## Request Tracing

All requests include a `X-Request-ID` header for distributed tracing:

```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

Use this ID when reporting issues or searching logs.

---

## Security

- All API endpoints require HTTPS in production
- JWT tokens expire after 24 hours
- Refresh tokens are valid for 7 days before expiry
- CSRF protection enabled for state-changing operations
- Rate limiting enforced per IP and per user
- Request timeouts prevent slow loris attacks
- Input validation and sanitization on all endpoints

---

## Examples

### Create a Session with cURL

```bash
curl -X POST https://streamspace.example.com/api/v1/sessions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "template": "firefox-browser",
    "resources": {
      "memory": "2Gi",
      "cpu": "1000m"
    }
  }'
```

### Hibernate a Session

```bash
curl -X PATCH https://streamspace.example.com/api/v1/sessions/session-id \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"state": "hibernated"}'
```

### Get Kubernetes Resource

```bash
curl -X GET "https://streamspace.example.com/api/v1/resources/deployment/my-app?apiVersion=apps/v1&kind=Deployment" \
  -H "Authorization: Bearer $TOKEN"
```

---

**For more information**, see the [StreamSpace Documentation](https://github.com/yourusername/streamspace/docs).
