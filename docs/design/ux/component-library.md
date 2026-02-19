# Component Library Inventory

**Version**: v2.0-beta
**Last Updated**: 2025-11-26
**Owner**: Frontend Team
**Status**: Living Document

---

## Introduction

This document inventories all reusable React components in the StreamSpace UI, including Material-UI (MUI) components and custom components. Use this as a reference when building new features to promote consistency and code reuse.

**Conventions**:
-  **Production Ready**: Fully implemented, tested, documented
-  **In Progress**: Implemented but needs refinement
-  **Planned**: Design approved, not yet implemented

---

## Component Categories

### 1. Layout Components
### 2. Display Components (Data)
### 3. Input Components (Forms)
### 4. Feedback Components (Loading, Errors)
### 5. Navigation Components
### 6. Domain-Specific Components

---

## 1. Layout Components

### App Shell

#### **AppLayout** 
- **Location**: `src/layouts/AppLayout.tsx`
- **Purpose**: Main application layout with sidebar and app bar
- **Props**:
  - `children`: React.ReactNode
- **Usage**:
  ```typescript
  <AppLayout>
    <Dashboard />
  </AppLayout>
  ```
- **MUI Components Used**: `Box`, `Drawer`, `AppBar`, `Toolbar`

#### **AdminLayout** 
- **Location**: `src/layouts/AdminLayout.tsx`
- **Purpose**: Layout for admin pages with expanded navigation
- **Props**: Same as AppLayout
- **Differences**: Additional admin nav items, different color scheme

### MUI Layout Components (Used Directly)

- **Box**  - Generic container (replaces `div`)
- **Container**  - Responsive centered container
- **Grid**  - 12-column responsive grid
- **Stack**  - 1-dimensional layout (vertical/horizontal)
- **Paper**  - Card-like container with elevation

---

## 2. Display Components (Data)

### Custom Components

#### **SessionCard** 
- **Location**: `src/components/SessionCard.tsx`
- **Purpose**: Display session information with actions
- **Props**:
  ```typescript
  interface SessionCardProps {
    session: Session;
    onConnect?: (sessionId: string) => void;
    onDelete?: (sessionId: string) => void;
    onHibernate?: (sessionId: string) => void;
  }
  ```
- **Features**:
  - Status badge (running, pending, stopped, failed)
  - Template name and icon
  - Created timestamp (relative format)
  - Action buttons (Connect, Delete, Hibernate)
  - Responsive (card on mobile, row on desktop)
- **MUI Components**: `Card`, `CardContent`, `CardActions`, `Chip`, `Button`
- **Test Coverage**:  85%

#### **TemplateCard** 
- **Location**: `src/components/TemplateCard.tsx` (to be created)
- **Purpose**: Display template in catalog
- **Props**:
  ```typescript
  interface TemplateCardProps {
    template: Template;
    onLaunch: (templateId: string) => void;
  }
  ```
- **Features**:
  - Template name and description
  - Category tags
  - Resource requirements (CPU, memory)
  - Launch button
- **Status**: Needs extraction from inline component

#### **TemplateDetailModal** 
- **Location**: `src/components/TemplateDetailModal.tsx`
- **Purpose**: Show template details in modal
- **Props**: `template: Template`, `open: boolean`, `onClose: () => void`
- **MUI Components**: `Dialog`, `DialogTitle`, `DialogContent`, `DialogActions`

#### **PluginCard** 
- **Location**: `src/components/PluginCard.tsx`
- **Purpose**: Display plugin in catalog
- **Props**:
  ```typescript
  interface PluginCardProps {
    plugin: Plugin;
    onInstall?: (pluginId: string) => void;
  }
  ```
- **Features**:
  - Plugin name, author, version
  - Rating stars
  - Install button
  - Tags/categories
- **Test Coverage**:  78%

#### **PluginCardSkeleton** 
- **Location**: `src/components/PluginCardSkeleton.tsx`
- **Purpose**: Loading placeholder for PluginCard
- **MUI Components**: `Skeleton`, `Card`

#### **PluginDetailModal** 
- **Location**: `src/components/PluginDetailModal.tsx`
- **Purpose**: Plugin details with installation options
- **Props**: `plugin: Plugin`, `open: boolean`, `onClose: () => void`

#### **RepositoryCard** 
- **Location**: `src/components/RepositoryCard.tsx`
- **Purpose**: Display template repository info
- **Props**: `repository: TemplateRepository`

#### **QuotaCard** 
- **Location**: `src/components/QuotaCard.tsx`
- **Purpose**: Display quota usage (sessions, CPU, memory)
- **Props**:
  ```typescript
  interface QuotaCardProps {
    label: string;
    current: number;
    limit: number;
    unit?: string;
  }
  ```
- **Features**:
  - Progress bar (color-coded: green → yellow → red)
  - Percentage display
  - Limit warning at 80%
- **MUI Components**: `Card`, `LinearProgress`, `Typography`

#### **QuotaAlert** 
- **Location**: `src/components/QuotaAlert.tsx`
- **Purpose**: Alert banner when quota exceeded
- **Props**: `quotaType: string`, `current: number`, `limit: number`
- **MUI Components**: `Alert`, `AlertTitle`

#### **RatingStars** 
- **Location**: `src/components/RatingStars.tsx`
- **Purpose**: Display star rating (for plugins)
- **Props**: `rating: number`, `totalRatings?: number`
- **MUI Components**: `Rating` (read-only)

#### **TagChip** 
- **Location**: `src/components/TagChip.tsx`
- **Purpose**: Display tag/category chip
- **Props**: `label: string`, `color?: string`, `onDelete?: () => void`
- **MUI Components**: `Chip`

### MUI Display Components (Used Directly)

- **Typography**  - Text display (h1-h6, body, caption)
- **Chip**  - Compact status/tag display
- **Badge**  - Notification badge
- **Avatar**  - User avatar (future)
- **Divider**  - Section separator
- **List** / **ListItem**  - Vertical lists
- **Table** / **TableRow** / **TableCell**  - Data tables

---

## 3. Input Components (Forms)

### MUI Input Components (Used Directly)

- **TextField**  - Text input
- **Select** / **MenuItem**  - Dropdown selection
- **Checkbox**  - Boolean input
- **Radio** / **RadioGroup**  - Single selection from options
- **Switch**  - Toggle on/off
- **Button**  - Primary action button
  - Variants: `contained`, `outlined`, `text`
  - Colors: `primary`, `secondary`, `error`, `success`
- **IconButton**  - Icon-only button
- **Autocomplete**  - Searchable dropdown

### Form Examples

**Standard Form Pattern**:
```typescript
import { TextField, Button, Box } from '@mui/material';

const CreateSessionForm = () => {
  const [templateId, setTemplateId] = useState('');

  return (
    <Box component="form" onSubmit={handleSubmit}>
      <TextField
        label="Template"
        value={templateId}
        onChange={(e) => setTemplateId(e.target.value)}
        fullWidth
        required
      />
      <Button type="submit" variant="contained" color="primary">
        Create Session
      </Button>
    </Box>
  );
};
```

---

## 4. Feedback Components (Loading, Errors)

### Custom Components

#### **ActivityIndicator** 
- **Location**: `src/components/ActivityIndicator.tsx`
- **Purpose**: Show activity/heartbeat status
- **Props**: `active: boolean`, `label?: string`
- **Features**:
  - Pulsing dot when active
  - Gray when inactive
  - Optional label

#### **NotificationQueue** 
- **Location**: `src/components/NotificationQueue.tsx`
- **Purpose**: Global notification snackbar queue
- **Usage**: Import `useNotificationStore` hook
- **Example**:
  ```typescript
  import { useNotificationStore } from '../store/notificationStore';

  const { addNotification } = useNotificationStore();

  addNotification('Session created successfully', 'success');
  addNotification('Failed to delete session', 'error');
  ```
- **MUI Components**: `Snackbar`, `Alert`

#### **ErrorBoundary** 
- **Location**: `src/components/ErrorBoundary.tsx`
- **Purpose**: Catch React component errors
- **Props**: `children`, `fallback?`
- **Usage**: Wrap app or critical sections
  ```typescript
  <ErrorBoundary fallback={<ErrorFallback />}>
    <App />
  </ErrorBoundary>
  ```

#### **WebSocketErrorBoundary** 
- **Location**: `src/components/WebSocketErrorBoundary.tsx`
- **Purpose**: Handle WebSocket connection errors
- **Features**: Auto-reconnect logic, error display

### MUI Feedback Components (Used Directly)

- **CircularProgress**  - Spinning loader (indeterminate)
- **LinearProgress**  - Progress bar (determinate/indeterminate)
- **Skeleton**  - Loading placeholder (content shimmer)
- **Alert**  - Inline alert (success, info, warning, error)
- **Snackbar**  - Toast notification
- **Dialog**  - Modal dialog
- **Backdrop**  - Overlay background

### Loading Patterns

**Skeleton Loading** (preferred for initial page load):
```typescript
import { Skeleton, Card, CardContent } from '@mui/material';

const SessionCardSkeleton = () => (
  <Card>
    <CardContent>
      <Skeleton variant="text" width="60%" height={30} />
      <Skeleton variant="rectangular" width="100%" height={100} />
    </CardContent>
  </Card>
);
```

**Spinner Loading** (for actions):
```typescript
import { CircularProgress, Button } from '@mui/material';

<Button disabled={loading}>
  {loading ? <CircularProgress size={20} /> : 'Create Session'}
</Button>
```

---

## 5. Navigation Components

### Custom Components

#### **EnhancedWebSocketStatus** 
- **Location**: `src/components/EnhancedWebSocketStatus.tsx`
- **Purpose**: Display WebSocket connection status in app bar
- **Props**: `status: 'connected' | 'disconnected' | 'reconnecting'`
- **Features**:
  - Color-coded indicator (green, red, yellow)
  - Connection latency display
  - Click to reconnect

### MUI Navigation Components (Used Directly)

- **Drawer**  - Sidebar navigation
  - Variants: `permanent`, `persistent`, `temporary`
- **AppBar**  - Top navigation bar
- **Toolbar**  - App bar content container
- **Tabs** / **Tab**  - Tabbed navigation
- **Breadcrumbs**  - Breadcrumb trail
- **Link**  - Navigation link (integrates with React Router)
- **Menu** / **MenuItem**  - Dropdown menu
- **BottomNavigation**  - Mobile bottom nav (future)

---

## 6. Domain-Specific Components

### Session Components

#### **SessionCard** 
(See Display Components above)

#### **SessionViewer** 
- **Location**: `src/pages/SessionViewer.tsx`
- **Purpose**: VNC stream viewer (full page component)
- **Features**:
  - noVNC client integration
  - Fullscreen mode
  - Clipboard sync
  - Keyboard/mouse capture
- **Dependencies**: `@novnc/novnc`

#### **IdleTimer** 
- **Location**: `src/components/IdleTimer.tsx`
- **Purpose**: Track user idle time for session hibernation
- **Props**: `timeout: number`, `onIdle: () => void`
- **Features**: Mouse/keyboard activity detection

### Template Components

#### **TemplateCard** 
(See Display Components above)

#### **TemplateDetailModal** 
(See Display Components above)

### Plugin Components

#### **PluginCard** 
#### **PluginDetailModal** 
#### **PluginCardSkeleton** 
(See Display Components above)

### Admin Components

#### **AgentStatusCard** 
- **Location**: TBD
- **Purpose**: Display agent health in Admin > Agents page
- **Props**: `agent: Agent`
- **Features**:
  - Heartbeat status (online, degraded, offline)
  - Last seen timestamp
  - Session count
  - Region/platform info

#### **AuditLogTable** 
- **Location**: TBD
- **Purpose**: Display audit logs in Admin > Audit page
- **Props**: `logs: AuditLog[]`
- **Features**:
  - Searchable, filterable, sortable
  - Pagination
  - Export to CSV

---

## WebSocket Providers

### **EnterpriseWebSocketProvider** 
- **Location**: `src/components/EnterpriseWebSocketProvider.tsx`
- **Purpose**: Global WebSocket connection manager
- **Features**:
  - Auto-reconnect with exponential backoff
  - Connection state management
  - Real-time session/metric updates
  - Org-scoped subscriptions
- **Usage**: Wrap app at root level
  ```typescript
  <EnterpriseWebSocketProvider wsUrl="wss://api/ws/ui">
    <App />
  </EnterpriseWebSocketProvider>
  ```

---

## Theming

### MUI Theme Configuration

**Location**: `src/theme.ts`

**Color Palette**:
```typescript
const theme = createTheme({
  palette: {
    mode: 'dark', // or 'light'
    primary: {
      main: '#1976d2', // Blue
    },
    secondary: {
      main: '#dc004e', // Pink
    },
    success: {
      main: '#4caf50', // Green
    },
    error: {
      main: '#f44336', // Red
    },
    warning: {
      main: '#ff9800', // Orange
    },
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
  },
});
```

**Theme Provider** :
```typescript
import { ThemeProvider, createTheme } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';

<ThemeProvider theme={theme}>
  <CssBaseline /> {/* Normalize CSS */}
  <App />
</ThemeProvider>
```

### Dark Mode Toggle

**Implementation**:
```typescript
import { useThemeMode } from './App'; // Context hook

const ThemeToggle = () => {
  const { mode, toggleTheme } = useThemeMode();

  return (
    <IconButton onClick={toggleTheme}>
      {mode === 'dark' ? <LightModeIcon /> : <DarkModeIcon />}
    </IconButton>
  );
};
```

---

## Icon Library

### MUI Icons

**Import**:
```typescript
import {
  DashboardIcon,
  ComputerIcon,
  SettingsIcon,
  PersonIcon,
  LogoutIcon,
  // ... 2000+ icons
} from '@mui/icons-material';
```

**Commonly Used Icons**:
- `DashboardIcon` - Dashboard page
- `ComputerIcon` - Sessions
- `ViewListIcon` - Templates
- `ExtensionIcon` - Plugins
- `SettingsIcon` - Settings
- `PersonIcon` - User profile
- `AdminPanelSettingsIcon` - Admin area
- `KeyIcon` - API keys
- `MonitorHeartIcon` - Monitoring
- `HistoryIcon` - Audit logs

---

## Component Usage Guidelines

### When to Create a New Component

**Create a new component when**:
- Used in 2+ places (DRY principle)
- Complex logic that can be isolated
- Testable unit (props in, UI out)
- Part of design system (consistent styling)

**Don't create a component when**:
- Used only once (inline is fine)
- Trivial (e.g., `<Box>` wrapper)
- Premature abstraction

### Component File Structure

```
src/components/
├── SessionCard.tsx        # Component implementation
├── SessionCard.test.tsx   # Unit tests
└── index.ts               # Barrel export (optional)
```

**Barrel Export** (`index.ts`):
```typescript
export { default as SessionCard } from './SessionCard';
export { default as TemplateCard } from './TemplateCard';
// ... allows: import { SessionCard, TemplateCard } from '@/components';
```

### Component Documentation

**JSDoc Comments**:
```typescript
/**
 * Displays session information with action buttons.
 *
 * @param session - Session object with id, status, template
 * @param onConnect - Callback when Connect button clicked
 * @param onDelete - Callback when Delete button clicked
 *
 * @example
 * <SessionCard
 *   session={mySession}
 *   onConnect={(id) => console.log('Connect', id)}
 *   onDelete={(id) => console.log('Delete', id)}
 * />
 */
export const SessionCard: React.FC<SessionCardProps> = ({ ... }) => { ... };
```

---

## Testing

### Component Testing (React Testing Library)

**Pattern**:
```typescript
import { render, screen, fireEvent } from '@testing-library/react';
import SessionCard from './SessionCard';

describe('SessionCard', () => {
  const mockSession = { id: 'sess-123', status: 'running', ... };

  it('renders session information', () => {
    render(<SessionCard session={mockSession} />);
    expect(screen.getByText('sess-123')).toBeInTheDocument();
  });

  it('calls onConnect when button clicked', () => {
    const handleConnect = jest.fn();
    render(<SessionCard session={mockSession} onConnect={handleConnect} />);
    fireEvent.click(screen.getByRole('button', { name: /connect/i }));
    expect(handleConnect).toHaveBeenCalledWith('sess-123');
  });
});
```

---

## MUI Component Reference

**Official Docs**: https://mui.com/material-ui/

**Most Used Components** (by frequency in codebase):
1. **Box** - ~500 usages (generic container)
2. **Typography** - ~300 usages (text)
3. **Button** - ~200 usages (actions)
4. **Card** / **CardContent** - ~150 usages (content containers)
5. **Grid** - ~100 usages (layout)
6. **TextField** - ~80 usages (forms)
7. **Dialog** - ~50 usages (modals)
8. **Chip** - ~40 usages (status badges)
9. **CircularProgress** - ~30 usages (loading)
10. **Alert** - ~20 usages (notifications)

---

## Future Component Additions (v2.1+)

### Planned Components

1. **UserAvatarMenu** 
   - User avatar with dropdown menu
   - Profile, settings, logout
   - Location: App bar (top right)

2. **SessionMetricsChart** 
   - Real-time CPU/memory chart for session
   - Uses Chart.js or Recharts
   - Location: Session viewer sidebar

3. **TemplateImportWizard** 
   - Multi-step wizard for importing templates
   - Validation, preview, confirmation steps
   - Location: Admin > Templates

4. **AccessibilityPanel** 
   - Accessibility settings panel
   - Font size, contrast, keyboard shortcuts
   - Location: User settings

5. **MultiSelectTable** 
   - Table with checkbox selection and bulk actions
   - For user management, session management
   - Reusable across admin pages

---

## References

- **Material-UI Docs**: https://mui.com/material-ui/
- **React Component Patterns**: https://react.dev/learn/thinking-in-react
- **Accessibility**: https://www.w3.org/WAI/ARIA/apg/

---

**Version History**:
- **v1.0** (2025-11-26): Initial component inventory for v2.0-beta
- **Next Review**: v2.1 release (Q1 2026)
