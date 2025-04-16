import { CommentCreateRequest, Comment } from '../types/api';
import { api } from './api';

export const commentService = {
  // Get all comments for an event
  getEventComments: async (eventId: number): Promise<Comment[]> => {
    try {
      const response = await api.get(`/events/${eventId}/comments`);
      return response.data;
    } catch (error) {
      console.error('Error fetching comments:', error);
      throw error;
    }
  },

  // Create a new comment
  createComment: async (commentData: CommentCreateRequest): Promise<any> => {
    try {
      const response = await api.post('/comments/create', commentData);
      return response.data;
    } catch (error) {
      console.error('Error creating comment:', error);
      throw error;
    }
  },

  markAsRead: async (commentId: number): Promise<any> => {
    try {
      const response = await api.patch(`/comments/${commentId}/mark-read`);
      return response.data;
    } catch (error) {
      console.error('Error marking comment as read:', error);
      throw error;
    }
  },

  deleteComment: async (commentId: number): Promise<any> => {
    try {
      const response = await api.delete(`/comments/${commentId}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting comment:', error);
      throw error;
    }
  }
};
