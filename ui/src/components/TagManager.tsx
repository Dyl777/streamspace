import { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Box,
  Typography,
  IconButton,
} from '@mui/material';
import { Add, Close } from '@mui/icons-material';
import TagChip from './TagChip';
import type { Session } from '../lib/api';

interface TagManagerProps {
  open: boolean;
  session: Session;
  onClose: () => void;
  onSave: (tags: string[]) => Promise<void>;
}

export default function TagManager({ open, session, onClose, onSave }: TagManagerProps) {
  const [tags, setTags] = useState<string[]>(session.tags || []);
  const [newTag, setNewTag] = useState('');
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');

  const handleAddTag = () => {
    const trimmedTag = newTag.trim().toLowerCase();

    if (!trimmedTag) {
      setError('Tag cannot be empty');
      return;
    }

    if (tags.includes(trimmedTag)) {
      setError('Tag already exists');
      return;
    }

    // Validate tag format (alphanumeric and hyphens only)
    if (!/^[a-z0-9-]+$/.test(trimmedTag)) {
      setError('Tag can only contain lowercase letters, numbers, and hyphens');
      return;
    }

    setTags([...tags, trimmedTag]);
    setNewTag('');
    setError('');
  };

  const handleRemoveTag = (tagToRemove: string) => {
    setTags(tags.filter(tag => tag !== tagToRemove));
  };

  const handleSave = async () => {
    setSaving(true);
    setError('');

    try {
      await onSave(tags);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update tags');
    } finally {
      setSaving(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddTag();
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h6">Manage Tags</Typography>
          <IconButton onClick={onClose} size="small">
            <Close />
          </IconButton>
        </Box>
      </DialogTitle>

      <DialogContent>
        <Typography variant="body2" color="text.secondary" gutterBottom>
          Session: <strong>{session.name}</strong>
        </Typography>

        <Box mt={2}>
          <Typography variant="subtitle2" gutterBottom>
            Current Tags
          </Typography>

          <Box display="flex" flexWrap="wrap" gap={0.5} mb={2}>
            {tags.length === 0 ? (
              <Typography variant="body2" color="text.secondary" fontStyle="italic">
                No tags yet
              </Typography>
            ) : (
              tags.map(tag => (
                <TagChip
                  key={tag}
                  tag={tag}
                  onDelete={() => handleRemoveTag(tag)}
                />
              ))
            )}
          </Box>
        </Box>

        <Box mt={2}>
          <Typography variant="subtitle2" gutterBottom>
            Add New Tag
          </Typography>

          <Box display="flex" gap={1} alignItems="flex-start">
            <TextField
              value={newTag}
              onChange={(e) => setNewTag(e.target.value)}
              onKeyPress={handleKeyPress}
              placeholder="e.g., development, testing, production"
              size="small"
              fullWidth
              error={!!error}
              helperText={error || 'Use lowercase letters, numbers, and hyphens'}
            />

            <Button
              variant="contained"
              onClick={handleAddTag}
              disabled={!newTag.trim()}
              startIcon={<Add />}
            >
              Add
            </Button>
          </Box>
        </Box>
      </DialogContent>

      <DialogActions>
        <Button onClick={onClose} disabled={saving}>
          Cancel
        </Button>
        <Button
          onClick={handleSave}
          variant="contained"
          disabled={saving}
        >
          {saving ? 'Saving...' : 'Save Changes'}
        </Button>
      </DialogActions>
    </Dialog>
  );
}
