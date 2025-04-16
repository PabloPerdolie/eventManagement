import { api } from './api';
import { 
  TaskResponse, 
  TasksResponse, 
  TaskCreateRequest, 
  TaskUpdateRequest,
  SuccessResponse
} from '../types/api';

export const taskService = {
  // Get all tasks with pagination and filtering
  getTasks: async (page: number = 1, size: number = 10, eventId?: number): Promise<TasksResponse> => {
    try {
      let url = `/tasks?page=${page}&size=${size}`;
      if (eventId) {
        url += `&event_id=${eventId}`;
      }
      const response = await api.get(url);
      return response.data;
    } catch (error) {
      console.error('Error fetching tasks:', error);
      throw error;
    }
  },

  // Create a new task
  createTask: async (taskData: TaskCreateRequest): Promise<TaskResponse> => {
    try {
      const response = await api.post('/tasks', taskData);
      return response.data;
    } catch (error) {
      console.error('Error creating task:', error);
      throw error;
    }
  },

  // Update an existing task
  updateTask: async (taskId: number, taskData: TaskUpdateRequest): Promise<TaskResponse> => {
    try {
      const response = await api.put(`/tasks/${taskId}`, taskData);
      return response.data;
    } catch (error) {
      console.error('Error updating task:', error);
      throw error;
    }
  },

  // Delete a task
  deleteTask: async (taskId: number): Promise<SuccessResponse> => {
    try {
      const response = await api.delete(`/tasks/${taskId}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting task:', error);
      throw error;
    }
  },
  
  // Get all tasks assigned to the current user across all events
  getUserTasks: async (page: number = 1, size: number = 10): Promise<TasksResponse> => {
    try {
      const response = await api.get(`/tasks/assigned?page=${page}&size=${size}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching user tasks:', error);
      throw error;
    }
  },

  // Get task details by ID
  getTaskById: async (taskId: number): Promise<TaskResponse> => {
    try {
      const response = await api.get(`/tasks/${taskId}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching task details:', error);
      throw error;
    }
  },

  // Check if task can be updated based on parent task status
  canUpdateTask: async (taskId: number): Promise<boolean> => {
    try {
      const task = await taskService.getTaskById(taskId);
      
      // If task has no parent, it can always be updated
      if (!task.parent_id) return true;
      
      // Get parent task details
      const parentTask = await taskService.getTaskById(task.parent_id);
      
      // Task can only be updated if parent task is completed
      return parentTask.status === 'completed';
    } catch (error) {
      console.error('Error checking if task can be updated:', error);
      // In case of error, default to false to prevent unauthorized updates
      return false;
    }
  }
};
