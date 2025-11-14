import { Grid, Paper, Typography, Box, Card, CardContent, Chip } from '@mui/material';
import {
  Computer as ComputerIcon,
  Apps as AppsIcon,
  Folder as FolderIcon,
  Timeline as TimelineIcon,
} from '@mui/icons-material';
import Layout from '../components/Layout';
import { useSessions, useTemplates, useRepositories, useMetrics } from '../hooks/useApi';
import { useUserStore } from '../store/userStore';

export default function Dashboard() {
  const username = useUserStore((state) => state.username);
  const { data: sessions = [], isLoading: sessionsLoading } = useSessions(username || undefined);
  const { data: templates = [], isLoading: templatesLoading } = useTemplates();
  const { data: repositories = [], isLoading: reposLoading } = useRepositories();
  const { data: metrics } = useMetrics();

  const stats = [
    {
      title: 'My Sessions',
      value: sessions.length,
      icon: <ComputerIcon sx={{ fontSize: 40 }} />,
      color: '#3f51b5',
      loading: sessionsLoading,
    },
    {
      title: 'Available Templates',
      value: templates.length,
      icon: <AppsIcon sx={{ fontSize: 40 }} />,
      color: '#f50057',
      loading: templatesLoading,
    },
    {
      title: 'Repositories',
      value: repositories.length,
      icon: <FolderIcon sx={{ fontSize: 40 }} />,
      color: '#4caf50',
      loading: reposLoading,
    },
    {
      title: 'Active Connections',
      value: metrics?.activeConnections || 0,
      icon: <TimelineIcon sx={{ fontSize: 40 }} />,
      color: '#ff9800',
      loading: false,
    },
  ];

  const runningSessions = sessions.filter((s) => s.state === 'running');
  const hibernatedSessions = sessions.filter((s) => s.state === 'hibernated');

  return (
    <Layout>
      <Box>
        <Typography variant="h4" sx={{ mb: 3, fontWeight: 700 }}>
          Welcome back, {username}!
        </Typography>

        <Grid container spacing={3} sx={{ mb: 4 }}>
          {stats.map((stat) => (
            <Grid item xs={12} sm={6} md={3} key={stat.title}>
              <Card>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Box>
                      <Typography color="text.secondary" variant="body2" sx={{ mb: 1 }}>
                        {stat.title}
                      </Typography>
                      <Typography variant="h4" sx={{ fontWeight: 700 }}>
                        {stat.loading ? '...' : stat.value}
                      </Typography>
                    </Box>
                    <Box sx={{ color: stat.color }}>{stat.icon}</Box>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>

        <Grid container spacing={3}>
          <Grid item xs={12} md={6}>
            <Paper sx={{ p: 3 }}>
              <Typography variant="h6" sx={{ mb: 2, fontWeight: 600 }}>
                Session Overview
              </Typography>
              <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Typography variant="body2" color="text.secondary">
                    Running
                  </Typography>
                  <Chip label={runningSessions.length} color="success" size="small" />
                </Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Typography variant="body2" color="text.secondary">
                    Hibernated
                  </Typography>
                  <Chip label={hibernatedSessions.length} color="warning" size="small" />
                </Box>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                  <Typography variant="body2" color="text.secondary">
                    Total
                  </Typography>
                  <Chip label={sessions.length} color="primary" size="small" />
                </Box>
              </Box>
            </Paper>
          </Grid>

          <Grid item xs={12} md={6}>
            <Paper sx={{ p: 3 }}>
              <Typography variant="h6" sx={{ mb: 2, fontWeight: 600 }}>
                Recent Sessions
              </Typography>
              {sessions.length === 0 ? (
                <Typography variant="body2" color="text.secondary">
                  No sessions yet. Create one from the Template Catalog!
                </Typography>
              ) : (
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                  {sessions.slice(0, 5).map((session) => (
                    <Box
                      key={session.name}
                      sx={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        p: 1,
                        borderRadius: 1,
                        '&:hover': { bgcolor: 'action.hover' },
                      }}
                    >
                      <Box>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {session.template}
                        </Typography>
                        <Typography variant="caption" color="text.secondary">
                          {session.name}
                        </Typography>
                      </Box>
                      <Chip
                        label={session.state}
                        size="small"
                        color={session.state === 'running' ? 'success' : 'default'}
                      />
                    </Box>
                  ))}
                </Box>
              )}
            </Paper>
          </Grid>
        </Grid>
      </Box>
    </Layout>
  );
}
