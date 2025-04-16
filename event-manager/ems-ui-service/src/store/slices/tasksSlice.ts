import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { TaskResponse, TaskStatus } from '../../types/api';
import { api } from '../../services/api';

interface TasksState {
  tasks: TaskResponse[];
  loading: boolean;
  error: string | null;
  total: number;
}

const initialState: TasksState = {
  tasks: [],
  loading: false,
  error: null,
  total: 0,
};

// Define async thunks
export const fetchTasks = createAsyncThunk(
  'tasks/fetchTasks',
  async ({ eventId, page = 1, size = 10 }: { eventId?: number; page?: number; size?: number }) => {
    const response = await api.get(`/tasks?event_id=${eventId}&page=${page}&size=${size}`);
    return response.data;
  }
);

export const fetchTasksByUserId = createAsyncThunk(
  'tasks/fetchTasksByUserId',
  async ({ page = 1, size = 10 }: { eventId?: number; page?: number; size?: number }) => {
    const response = await api.get(`/tasks?page=${page}&size=${size}`);
    return response.data;
  }
);

export const createTask = createAsyncThunk(
  'tasks/createTask',
  async (taskData: {
    title: string;
    description?: string;
    event_id: number;
    assigned_to?: number;
    priority?: string;
    story_points?: number;
  }) => {
    const response = await api.post('/tasks', taskData);
    return response.data;
  }
);

export const updateTaskStatus = createAsyncThunk(
  'tasks/updateTaskStatus',
  async ({ taskId, status }: { taskId: number; status: TaskStatus }) => {
    const response = await api.put(`/tasks/${taskId}`, { status });
    return response.data;
  }
);

// Create the slice with reducers and extra reducers
const tasksSlice = createSlice({
  name: 'tasks',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    // Handle fetchTasks actions
    builder
      .addCase(fetchTasks.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTasks.fulfilled, (state, action) => {
        state.loading = false;
        state.tasks = action.payload.tasks;
        state.total = action.payload.total;
      })
      .addCase(fetchTasks.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch tasks';
      })
      // Handle other thunk actions
      .addCase(createTask.fulfilled, (state, action) => {
        state.tasks.unshift(action.payload);
      })
      .addCase(updateTaskStatus.fulfilled, (state, action) => {
        const index = state.tasks.findIndex(task => task.id === action.payload.id);
        if (index !== -1) {
          state.tasks[index] = action.payload;
        }
      });
  },
});

export default tasksSlice.reducer;