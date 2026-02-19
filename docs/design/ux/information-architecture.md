# Information Architecture

**Version**: v2.0-beta
**Last Updated**: 2025-11-26
**Owner**: UX/Frontend Team
**Status**: Living Document

---

## Introduction

This document defines the information architecture (IA) for the StreamSpace Web UI, including site structure, navigation hierarchy, URL routing, and page organization.

**Goals**:
- Clear, intuitive navigation for all user roles
- Scalable structure for future features
- Consistent URL patterns
- Accessibility and discoverability

---

## User Roles

### 1. End User
- Access and manage personal sessions
- Browse template catalog
- View usage metrics

### 2. Organization Admin
- Manage org users and groups
- Configure templates and policies
- View org-wide metrics

### 3. Platform Admin
- System configuration
- Agent management
- Platform monitoring
- Compliance and audit

---

## Site Map

```
StreamSpace
│
├── Public (Unauthenticated)
│   ├── /login                    # Login page
│   └── /setup                    # Setup wizard (first-time deployment)
│
├── User Area (Authenticated)
│   ├── /                         # Dashboard (default landing)
│   ├── /sessions                 # Session list
│   ├── /sessions/:id             # Session viewer (VNC)
│   ├── /templates                # Template catalog
│   ├── /plugins                  # Plugin catalog
│   └── /plugins/installed        # Installed plugins
│
└── Admin Area (Admin Role)
    ├── /admin                    # Admin dashboard
    ├── /admin/users              # User management
    ├── /admin/groups             # Group management
    ├── /admin/groups/create      # Create group
    ├── /admin/groups/:id         # Group detail
    ├── /admin/templates          # Template management
    ├── /admin/agents             # Agent status & config
    ├── /admin/api-keys           # API key management
    ├── /admin/settings           # System settings
    ├── /admin/monitoring         # System monitoring
    ├── /admin/audit              # Audit logs
    ├── /admin/recordings         # Session recordings
    ├── /admin/compliance         # Compliance reports
    └── /admin/plugins            # Plugin management
```

---

## Navigation Structure

### Primary Navigation (Authenticated Users)

Located in left sidebar (Material-UI Drawer):

```
┌─────────────────────────┐
│ StreamSpace Logo        │
├─────────────────────────┤
│ � Dashboard            │
│ � Sessions             │
│  Templates            │
│ � Plugins              │
├─────────────────────────┤
│  Settings             │ (User settings)
│ � Profile              │
│ � Logout               │
└─────────────────────────┘
```

### Admin Navigation (Admin Users Only)

Additional section in sidebar:

```
┌─────────────────────────┐
│  Admin                │ (Expandable section)
│   ├─ Dashboard          │
│   ├─ Users              │
│   ├─ Groups             │
│   ├─ Templates          │
│   ├─ Agents             │
│   ├─ API Keys           │
│   ├─ Settings           │
│   ├─ Monitoring         │
│   ├─ Audit Logs         │
│   ├─ Recordings         │
│   ├─ Compliance         │
│   └─ Plugins            │
└─────────────────────────┘
```

---

## Page Hierarchy

### 1. Public Pages

#### `/login` - Login Page
- **Purpose**: User authentication
- **Components**: LoginForm, SSOButtons, MFAInput
- **Layout**: Centered, no sidebar
- **Routes**:
  - Success → `/` (Dashboard)
  - First-time setup → `/setup`

#### `/setup` - Setup Wizard
- **Purpose**: First-time platform configuration
- **Steps**:
  1. Welcome
  2. Admin account creation
  3. Database configuration
  4. SSO configuration (optional)
  5. Agent deployment instructions
- **Layout**: Wizard stepper, no sidebar
- **Routes**: Complete → `/login`

---

### 2. User Pages

#### `/` - Dashboard
- **Purpose**: Overview of user's sessions and activity
- **Components**:
  - ActiveSessionsCard (quick access to running sessions)
  - RecentActivityTimeline
  - QuickActionsPanel (Create Session button)
  - UsageMetricsChart (if enabled)
- **Permissions**: All authenticated users

#### `/sessions` - Session List
- **Purpose**: View and manage personal sessions
- **Components**:
  - SessionFilter (status, template, date)
  - SessionList (table or grid)
  - SessionCard (with Connect/Delete actions)
  - CreateSessionButton
- **Permissions**: All authenticated users
- **URL Params**: `?status=running&template=ubuntu`

#### `/sessions/:id` - Session Viewer
- **Purpose**: VNC stream viewer for active session
- **Components**:
  - VNCViewer (noVNC client)
  - SessionToolbar (fullscreen, keyboard, clipboard)
  - SessionInfo (sidebar with metadata)
- **Permissions**: Session owner only (org-scoped)
- **URL Example**: `/sessions/sess-abc-123`

#### `/templates` - Template Catalog
- **Purpose**: Browse and filter available templates
- **Components**:
  - TemplateGrid
  - TemplateCard (with Launch button)
  - TemplateFilter (category, tags, search)
- **Permissions**: All authenticated users (org-scoped)

#### `/plugins` - Plugin Catalog
- **Purpose**: Browse available plugins
- **Components**:
  - PluginGrid
  - PluginCard (with Install button)
  - PluginFilter
- **Permissions**: All authenticated users

#### `/plugins/installed` - Installed Plugins
- **Purpose**: Manage installed plugins
- **Components**:
  - InstalledPluginList
  - PluginSettings
  - UninstallButton
- **Permissions**: All authenticated users

---

### 3. Admin Pages

#### `/admin` - Admin Dashboard
- **Purpose**: Platform overview for admins
- **Components**:
  - PlatformMetrics (total sessions, users, orgs)
  - AgentHealthStatus
  - RecentAuditEvents
  - SystemAlertsPanel
- **Permissions**: Admin role only

#### `/admin/users` - User Management
- **Purpose**: Manage platform users
- **Components**:
  - UserTable (searchable, filterable)
  - CreateUserButton
  - BulkActionsMenu (enable/disable, delete)
- **Permissions**: Org Admin or Platform Admin

#### `/admin/groups` - Group Management
- **Purpose**: Manage user groups for RBAC
- **Components**:
  - GroupList
  - CreateGroupButton
- **Permissions**: Org Admin or Platform Admin
- **Routes**:
  - Create → `/admin/groups/create`
  - View/Edit → `/admin/groups/:id`

#### `/admin/templates` - Template Management
- **Purpose**: Create and configure session templates
- **Components**:
  - TemplateTable
  - CreateTemplateButton
  - TemplateEditor (YAML/JSON)
- **Permissions**: Org Admin or Platform Admin

#### `/admin/agents` - Agent Management
- **Purpose**: Monitor and configure execution agents
- **Components**:
  - AgentList (status, heartbeat, region)
  - AgentDetailCard
  - AgentConfigEditor
- **Permissions**: Platform Admin only

#### `/admin/api-keys` - API Key Management
- **Purpose**: Generate and revoke API keys
- **Components**:
  - APIKeyTable
  - CreateAPIKeyButton
  - RevokeAPIKeyButton
- **Permissions**: Org Admin or Platform Admin

#### `/admin/settings` - System Settings
- **Purpose**: Platform configuration
- **Sections**:
  - General (platform name, URL)
  - Authentication (SSO, MFA)
  - Quotas (session limits, resource limits)
  - Security (IP allow/deny lists)
  - Storage (home directory backend)
- **Permissions**: Platform Admin only

#### `/admin/monitoring` - System Monitoring
- **Purpose**: Real-time platform health
- **Components**:
  - MetricsDashboard (CPU, memory, sessions/sec)
  - AlertsPanel
  - LogViewer
- **Permissions**: Platform Admin only

#### `/admin/audit` - Audit Logs
- **Purpose**: Security and compliance audit trail
- **Components**:
  - AuditLogTable (searchable by user, action, date)
  - AuditLogFilter
  - ExportButton (CSV, JSON)
- **Permissions**: Org Admin or Platform Admin

#### `/admin/recordings` - Session Recordings
- **Purpose**: View and manage session recordings
- **Components**:
  - RecordingTable
  - RecordingPlayer
- **Permissions**: Org Admin or Platform Admin

#### `/admin/compliance` - Compliance Reports
- **Purpose**: SOC2, HIPAA, PCI compliance reports
- **Components**:
  - ComplianceChecklist
  - ComplianceReport (PDF export)
- **Permissions**: Platform Admin only

#### `/admin/plugins` - Plugin Management (Admin)
- **Purpose**: Configure plugin policies, approve plugins
- **Components**:
  - PluginPolicyEditor
  - PluginApprovalQueue
- **Permissions**: Platform Admin only

---

## URL Routing

### Route Patterns

**RESTful Conventions**:
- List: `/resources` (GET)
- Detail: `/resources/:id` (GET)
- Create: `/resources/create` (GET form)
- Edit: `/resources/:id/edit` (GET form)
- Actions: `/resources/:id/:action` (POST)

**Examples**:
```
GET  /sessions              # List sessions
GET  /sessions/sess-123     # View session
POST /sessions              # Create session (API)
GET  /sessions/create       # Create session form (UI)
GET  /sessions/sess-123/edit# Edit session (future)
POST /sessions/sess-123/hibernate # Hibernate action
```

### Route Guards

**Authentication**:
- Public routes: `/login`, `/setup`
- Protected routes: All others (redirect to `/login` if unauthenticated)

**Authorization**:
- User routes: All authenticated users
- Admin routes: `role = "admin"` or `role = "org_admin"`
- Org scoping: Filter resources by `org_id` from JWT

**Implementation** (React Router):
```typescript
<Route
  path="/admin/*"
  element={
    <RequireAuth requireRole="admin">
      <AdminLayout />
    </RequireAuth>
  }
/>
```

---

## Breadcrumbs

**Pattern**: Home > Section > Page > Detail

**Examples**:
```
Home > Sessions
Home > Sessions > sess-123
Home > Templates
Home > Admin > Users
Home > Admin > Groups > Create Group
Home > Admin > Groups > Engineering Team
```

**Implementation**:
- Auto-generated from route hierarchy
- Clickable links for navigation
- Located below app bar (top of content area)

---

## Search and Navigation

### Global Search (Future v2.1)

Location: App bar (top right)

**Searchable Entities**:
- Sessions (by ID, template, status)
- Templates (by name, description, tags)
- Users (by name, email) - Admin only
- API Keys (by name) - Admin only

**Search Results**:
- Grouped by entity type
- Top 5 results per type
- "View all" link to dedicated search page

---

## Mobile Responsiveness

### Breakpoints (Material-UI)

- **xs** (0-600px): Phone portrait
- **sm** (600-960px): Phone landscape, small tablet
- **md** (960-1280px): Tablet landscape
- **lg** (1280-1920px): Desktop
- **xl** (1920px+): Large desktop

### Mobile Adaptations

**Sidebar Navigation**:
- xs/sm: Drawer hidden by default, hamburger menu
- md+: Permanent drawer (always visible)

**Session List**:
- xs/sm: Card layout (stacked)
- md+: Table layout (grid)

**Admin Pages**:
- xs/sm: Simplified layout, hide less critical info
- md+: Full dashboard with all widgets

---

## Accessibility

### Navigation

- **Keyboard Navigation**: All interactive elements accessible via keyboard (Tab, Enter, Escape)
- **ARIA Labels**: Descriptive labels for screen readers
- **Focus Indicators**: Clear visual focus states
- **Skip Links**: "Skip to main content" for screen readers

### URL Structure

- **Meaningful URLs**: `/sessions` not `/s`, `/admin/users` not `/a/u`
- **Persistent URLs**: Session URLs remain valid (bookmarkable)
- **No State in URLs**: Avoid encoding complex state in query params

---

## URL Examples

### User Flows

**Create Session**:
1. User clicks "Create Session" on Dashboard
2. Navigate to `/templates` (or inline modal)
3. User selects template "Ubuntu Desktop"
4. POST to `/api/v1/sessions` (API call)
5. Navigate to `/sessions/sess-abc-123` (new session viewer)

**Browse Templates**:
1. User clicks "Templates" in sidebar
2. Navigate to `/templates`
3. User filters by category "Development"
4. URL updates to `/templates?category=development`

**Admin Manage Users**:
1. Admin clicks "Admin > Users"
2. Navigate to `/admin/users`
3. Admin searches for "alice"
4. URL updates to `/admin/users?search=alice`
5. Admin clicks user row
6. Navigate to `/admin/users/user-456` (user detail)

---

## Page Layout Components

### Standard Layout

All authenticated pages use consistent layout:

```
┌────────────────────────────────────────────┐
│ App Bar (Logo, Breadcrumbs, User Menu)    │
├──────────┬─────────────────────────────────┤
│          │                                 │
│ Sidebar  │ Content Area                    │
│ Nav      │ (Page-specific components)      │
│          │                                 │
│          │                                 │
│          │                                 │
└──────────┴─────────────────────────────────┘
```

### Exception Layouts

- **Login Page**: Centered, no sidebar
- **Setup Wizard**: Wizard stepper, no sidebar
- **Session Viewer**: Fullscreen VNC, minimal chrome (optional hide controls)

---

## Future Enhancements (v2.1+)

### Planned IA Changes

1. **User Profile Page** (`/profile`)
   - Edit user settings, avatar, preferences
   - MFA configuration

2. **Session History** (`/sessions/history`)
   - Archive of stopped/deleted sessions
   - Usage reports

3. **Favorites/Starred Templates** (`/templates/favorites`)
   - Quick access to frequently used templates

4. **Notifications Center** (`/notifications`)
   - Session events, quota alerts, system messages

5. **Multi-Org Switcher** (if user belongs to multiple orgs)
   - Org switcher in app bar
   - URL structure: `/org/:org_id/sessions`

---

## References

- **Material-UI Navigation**: [MUI Drawer](https://mui.com/material-ui/react-drawer/)
- **React Router**: [React Router v6](https://reactrouter.com/)
- **URL Design**: [RESTful URL Best Practices](https://restfulapi.net/resource-naming/)
- **IA Best Practices**: [IA Institute](https://www.iainstitute.org/)

---

**Version History**:
- **v1.0** (2025-11-26): Initial IA for v2.0-beta
- **Next Review**: v2.1 release (Q1 2026)
