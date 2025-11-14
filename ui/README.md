# StreamSpace Web UI

React TypeScript frontend for StreamSpace - Stream any containerized application to your browser.

## Features

### ✅ Implemented

- **User Authentication** (Demo mode - username only)
- **Dashboard**
  - Session statistics
  - Active connections count
  - Recent sessions overview
- **Session Management**
  - View all user sessions
  - Start/Stop (Running ↔ Hibernated)
  - Connect to running sessions
  - Delete sessions
  - Real-time status updates (polling)
- **Template Catalog**
  - Browse installed templates
  - Browse marketplace templates
  - Create sessions from templates
  - Filter by category and tags
- **Repository Management**
  - View all template repositories
  - Add new repositories
  - Sync repositories (individual or all)
  - Delete repositories
  - Repository status tracking

## Tech Stack

- **Framework**: React 18
- **Language**: TypeScript
- **Build Tool**: Vite
- **UI Library**: Material-UI (MUI) v5
- **State Management**: Zustand
- **Data Fetching**: TanStack Query (React Query)
- **Routing**: React Router v6
- **HTTP Client**: Axios

## Project Structure

```
ui/
├── public/                     # Static assets
├── src/
│   ├── components/
│   │   └── Layout.tsx         # Main layout with sidebar and app bar
│   ├── hooks/
│   │   └── useApi.ts          # React Query hooks for API
│   ├── lib/
│   │   └── api.ts             # API client (Axios)
│   ├── pages/
│   │   ├── Login.tsx          # Login page (demo mode)
│   │   ├── Dashboard.tsx      # Overview dashboard
│   │   ├── Sessions.tsx       # Session management
│   │   ├── Catalog.tsx        # Template catalog browser
│   │   └── Repositories.tsx   # Repository management
│   ├── store/
│   │   └── userStore.ts       # Zustand user state
│   ├── App.tsx                # Main app component with routing
│   ├── main.tsx               # Entry point
│   └── index.css              # Global styles
├── index.html                 # HTML template
├── package.json               # Dependencies
├── tsconfig.json              # TypeScript configuration
├── vite.config.ts             # Vite configuration
└── README.md                  # This file
```

## Development

### Prerequisites

- Node.js 18+
- npm or yarn
- StreamSpace API backend running on `http://localhost:8000`

### Installation

```bash
cd ui
npm install
```

### Running Locally

```bash
npm run dev
```

The UI will start on `http://localhost:3000` with proxy to API backend.

### Building for Production

```bash
npm run build
```

Build output will be in `dist/` directory.

### Preview Production Build

```bash
npm run preview
```

## Configuration

### Environment Variables

Create `.env.local` for environment-specific configuration:

```bash
# API Base URL (optional, uses proxy in development)
VITE_API_URL=http://localhost:8000/api/v1
```

### Vite Proxy

Development proxy is configured in `vite.config.ts`:

```typescript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8000',
      changeOrigin: true,
    }
  }
}
```

## Features Overview

### Login (Demo Mode)

- Simple username entry
- No password required (for prototype)
- Users can enter any username
- Username "admin" gets admin role
- TODO: Full OIDC authentication in Phase 2.3

### Dashboard

- Session count by state (running, hibernated)
- Active connections count
- Template and repository counts
- Recent sessions list
- Real-time metrics

### Sessions Page

- View all user sessions
- Session cards with:
  - Template name
  - State and phase status
  - Resource allocation
  - Active connections
  - Access URL
- Actions:
  - **Connect**: Open session in new tab
  - **Play/Pause**: Toggle running ↔ hibernated
  - **Delete**: Remove session
- Auto-refresh every 5 seconds

### Catalog Page

Two tabs:
1. **Installed Templates**: Templates ready to use
2. **Marketplace**: Templates from repositories

Features:
- Browse templates by category
- View template details (description, tags, app type)
- Create session from template (one-click)
- Filter by category and tags

### Repositories Page

- Table view of all repositories
- Repository details:
  - Name, URL, branch
  - Sync status (pending, syncing, synced, failed)
  - Template count
  - Last sync timestamp
- Actions:
  - **Add Repository**: Add new Git repository
  - **Sync**: Trigger sync for repository
  - **Sync All**: Sync all repositories
  - **Delete**: Remove repository

## API Integration

All API calls go through `src/lib/api.ts` which provides:

### Session Management
- `listSessions(user?)` - List sessions
- `getSession(id)` - Get session details
- `createSession(data)` - Create new session
- `updateSessionState(id, state)` - Update session state
- `deleteSession(id)` - Delete session
- `connectSession(id, user)` - Connect to session
- `sendHeartbeat(id, connectionId)` - Send connection heartbeat

### Template Management
- `listTemplates(category?)` - List templates
- `getTemplate(id)` - Get template details
- `createTemplate(data)` - Create template
- `deleteTemplate(id)` - Delete template

### Catalog & Repositories
- `listCatalogTemplates(category?, tag?)` - Browse catalog
- `installCatalogTemplate(id)` - Install from catalog
- `listRepositories()` - List repositories
- `addRepository(data)` - Add repository
- `syncRepository(id)` - Sync repository
- `deleteRepository(id)` - Delete repository

### React Query Hooks

All hooks auto-refresh and provide loading/error states:

```typescript
// Example usage
const { data: sessions, isLoading, error } = useSessions(username);
const createSession = useCreateSession();

// Mutations automatically invalidate related queries
createSession.mutate(sessionData);
```

## Roadmap

### Phase 2 UI (Current - Complete)
- ✅ Login page (demo mode)
- ✅ Dashboard with stats
- ✅ Session management
- ✅ Template catalog browser
- ✅ Repository management
- ✅ Responsive layout
- ✅ Dark theme

### Phase 3 (Future)
- [ ] Full OIDC authentication (Authentik/Keycloak)
- [ ] User profile and settings
- [ ] Admin panel (all users, all sessions)
- [ ] Advanced filtering and search
- [ ] Session resource customization UI

### Phase 4 (Future)
- [ ] WebSocket real-time updates
- [ ] Pod logs viewer
- [ ] Terminal (exec into containers)
- [ ] Session sharing
- [ ] Metrics and analytics charts

### Phase 5 (Future)
- [ ] Cluster management UI
- [ ] YAML editor for resources
- [ ] Node management
- [ ] Deployment management
- [ ] Service management

## Contributing

See the main [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

## License

MIT License - See [LICENSE](../LICENSE)
