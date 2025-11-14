import { useEffect, useState } from 'react';
import { Typography, Box, LinearProgress, Tooltip } from '@mui/material';
import { AccessTime as TimeIcon } from '@mui/icons-material';

interface IdleTimerProps {
  lastActivity?: string;
  idleDuration?: number; // seconds
  idleThreshold?: number; // seconds
  showProgress?: boolean;
  compact?: boolean;
}

export default function IdleTimer({
  lastActivity,
  idleDuration = 0,
  idleThreshold = 0,
  showProgress = false,
  compact = false,
}: IdleTimerProps) {
  const [, setTick] = useState(0);

  // Update every second for live countdown
  useEffect(() => {
    const interval = setInterval(() => {
      setTick((prev) => prev + 1);
    }, 1000);

    return () => clearInterval(interval);
  }, []);

  if (!lastActivity && idleDuration === 0) {
    return (
      <Box display="flex" alignItems="center" gap={0.5}>
        <TimeIcon sx={{ fontSize: 16, color: 'text.secondary' }} />
        <Typography variant="caption" color="text.secondary">
          No activity data
        </Typography>
      </Box>
    );
  }

  // Calculate time since last activity
  let secondsSinceActivity = idleDuration;
  if (lastActivity && !idleDuration) {
    const lastActivityDate = new Date(lastActivity);
    const now = new Date();
    secondsSinceActivity = Math.floor((now.getTime() - lastActivityDate.getTime()) / 1000);
  }

  const formattedTime = formatDuration(secondsSinceActivity);

  // Calculate progress percentage (if threshold provided)
  const progressPercentage = idleThreshold > 0
    ? Math.min((secondsSinceActivity / idleThreshold) * 100, 100)
    : 0;

  const isNearingThreshold = progressPercentage >= 75;
  const exceededThreshold = progressPercentage >= 100;

  const getColor = () => {
    if (exceededThreshold) return 'error';
    if (isNearingThreshold) return 'warning';
    return 'primary';
  };

  const tooltip = idleThreshold > 0
    ? `Idle for ${formattedTime} of ${formatDuration(idleThreshold)} threshold`
    : `Last activity: ${formattedTime} ago`;

  if (compact) {
    return (
      <Tooltip title={tooltip}>
        <Box display="flex" alignItems="center" gap={0.5}>
          <TimeIcon
            sx={{
              fontSize: 16,
              color: exceededThreshold ? 'error.main' : isNearingThreshold ? 'warning.main' : 'text.secondary'
            }}
          />
          <Typography
            variant="caption"
            color={exceededThreshold ? 'error' : isNearingThreshold ? 'warning.main' : 'text.secondary'}
          >
            {formattedTime}
          </Typography>
        </Box>
      </Tooltip>
    );
  }

  return (
    <Box>
      <Box display="flex" alignItems="center" gap={1} mb={0.5}>
        <TimeIcon sx={{ fontSize: 18, color: 'text.secondary' }} />
        <Typography variant="body2" color="text.secondary">
          Idle: <strong style={{ color: exceededThreshold ? 'error' : isNearingThreshold ? 'warning' : 'inherit' }}>
            {formattedTime}
          </strong>
        </Typography>
      </Box>

      {showProgress && idleThreshold > 0 && (
        <Box mt={1}>
          <LinearProgress
            variant="determinate"
            value={progressPercentage}
            color={getColor()}
            sx={{ height: 6, borderRadius: 1 }}
          />
          <Typography variant="caption" color="text.secondary" display="block" mt={0.5}>
            {formatDuration(Math.max(0, idleThreshold - secondsSinceActivity))} until auto-hibernation
          </Typography>
        </Box>
      )}
    </Box>
  );
}

// Format seconds into human-readable duration
function formatDuration(seconds: number): string {
  if (seconds < 60) {
    return `${seconds}s`;
  }

  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (days > 0) {
    const remainingHours = hours % 24;
    return remainingHours > 0 ? `${days}d ${remainingHours}h` : `${days}d`;
  }

  if (hours > 0) {
    const remainingMinutes = minutes % 60;
    return remainingMinutes > 0 ? `${hours}h ${remainingMinutes}m` : `${hours}h`;
  }

  const remainingSeconds = seconds % 60;
  return remainingSeconds > 0 ? `${minutes}m ${remainingSeconds}s` : `${minutes}m`;
}
