import React, { useEffect, useState } from 'react';
import { list, create, update, remove } from '../services/taskService';
import {
  Box, Typography, Button, List, ListItem, ListItemText,
  IconButton, TextField, Dialog, DialogTitle, DialogContent,
  DialogActions, Checkbox
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';

function Todo() {
  const [tasks, setTasks] = useState([]);
  const [newTask, setNewTask] = useState({ title: '', description: '' });
  const [editTask, setEditTask] = useState(null);
  const [open, setOpen] = useState(false);

  const fetchTasks = async () => {
    try {
      const response = await list();
      const taskArray = Array.isArray(response.data) ? response.data : response.data.tasks || response.data.data || [];
      setTasks(taskArray);
    } catch (error) {
      console.error("Tasklarni olishda xatolik:", error);
    }
  };

  const handleCreate = async () => {
    if (!newTask.title) return;
    try {
      await create({ ...newTask, done: false });
      setNewTask({ title: '', description: '' });
      fetchTasks();
    } catch (error) {
      console.error("Yaratishda xatolik:", error);
    }
  };

  const handleDelete = async (id) => {
    try {
      await remove(id);
      setTasks(tasks.filter(t => t.id !== id));
    } catch (error) {
      console.error("O‘chirishda xatolik:", error);
    }
  };

  const handleUpdate = async () => {
    try {
      await update(editTask.id, editTask);
      setEditTask(null);
      setOpen(false);
      fetchTasks();
    } catch (error) {
      console.error("Tahrirlashda xatolik:", error);
    }
  };

  const handleToggleDone = async (task) => {
    try {
      await update(task.id, { ...task, done: !task.done });
      fetchTasks();
    } catch (error) {
      console.error("Done statusni o‘zgartirishda xatolik:", error);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, []);

  return (
    <Box sx={{ backgroundColor: '#1e1e1e', p: 3, borderRadius: 3 }}>
      <Typography variant="h5" color="white" textAlign="center" mb={2}>
        Yangi vazifa qo‘shish
      </Typography>

      <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
        <TextField
          label="Sarlavha"
          variant="outlined"
          size="small"
          value={newTask.title}
          onChange={(e) => setNewTask({ ...newTask, title: e.target.value })}
          fullWidth
        />
        <TextField
          label="Izoh"
          variant="outlined"
          size="small"
          value={newTask.description}
          onChange={(e) => setNewTask({ ...newTask, description: e.target.value })}
          fullWidth
        />
        <Button
          variant="contained"
          onClick={handleCreate}
          sx={{ backgroundColor: '#00ff2a', '&:hover': { backgroundColor: '#00cc22' } }}
        >
          Qo‘shish
        </Button>
      </Box>

      <List sx={{ mt: 4 }}>
        {tasks.length === 0 && (
          <Typography color="white" textAlign="center" fontStyle="italic">
            Hozircha hech qanday vazifa yo‘q.
          </Typography>
        )}
        {tasks.map(task => (
          <ListItem
            key={task.id}
            sx={{
              backgroundColor: '#2e2e2e',
              borderRadius: 2,
              mb: 1,
              color: 'white'
            }}
            secondaryAction={
              <Box sx={{ display: 'flex', gap: 1 }}>
                <IconButton onClick={() => { setEditTask(task); setOpen(true); }}>
                  <EditIcon sx={{ color: '#00e676' }} />
                </IconButton>
                <IconButton onClick={() => handleDelete(task.id)}>
                  <DeleteIcon sx={{ color: '#ff5252' }} />
                </IconButton>
              </Box>
            }
          >
            <Checkbox
              checked={task.done}
              onChange={() => handleToggleDone(task)}
              sx={{ color: '#00e676' }}
            />
            <ListItemText
              primary={task.title}
              secondary={task.description}
              sx={{
                textDecoration: task.done ? 'line-through' : 'none',
                color: task.done ? '#888' : '#fff'
              }}
            />
          </ListItem>
        ))}
      </List>

      <Dialog open={open} onClose={() => setOpen(false)}>
        <DialogTitle>Taskni tahrirlash</DialogTitle>
        <DialogContent>
          <TextField
            margin="dense"
            label="Sarlavha"
            fullWidth
            value={editTask?.title || ''}
            onChange={(e) => setEditTask({ ...editTask, title: e.target.value })}
          />
          <TextField
            margin="dense"
            label="Izoh"
            fullWidth
            multiline
            rows={3}
            value={editTask?.description || ''}
            onChange={(e) => setEditTask({ ...editTask, description: e.target.value })}
            sx={{ mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Bekor qilish</Button>
          <Button onClick={handleUpdate} variant="contained">Saqlash</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}

export default Todo;
