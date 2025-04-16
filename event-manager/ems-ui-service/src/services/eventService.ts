import { api } from './api';
import { 
  EventResponse, 
  EventsResponse, 
  EventCreateRequest, 
  EventData,
  EventParticipantCreateRequest,
  EventParticipantsResponse
} from '../types/api';

export const eventService = {
  // Get all events with pagination
  getEvents: async (page: number = 1, size: number = 10): Promise<EventsResponse> => {
    try {
      const response = await api.get(`/events?page=${page}&size=${size}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching events:', error);
      throw error;
    }
  },

  // Get a specific event with all details
  getEvent: async (eventId: number): Promise<EventData> => {
    try {
      const response = await api.get(`/events/${eventId}`);
      const data = response.data;
      
      if (!data) {
        throw new Error('Event data is empty');
      }
      
      const transformResponse = (data: any) => ({
        eventData: data.EventData || data,
        eventParticipants: data.EventParticipants || data.participants || { participants: [], total: 0 },
        tasks: data.Tasks || data.tasks || { tasks: [], total: 0 },
        comments: data.Comments ? {
          comments: Array.isArray(data.Comments.Comments) 
            ? data.Comments.Comments.map((comment: any) => ({
                commentId: comment.CommentId,
                content: comment.Content,
                eventId: comment.EventId,
                taskId: comment.TaskId,
                senderId: comment.SenderId,
                createdAt: comment.CreatedAt,
                isRead: comment.IsRead,
                isDeleted: comment.IsDeleted
              }))
            : data.Comments
        } : { comments: [] },
        expenses: data.Expenses || data.expenses || { expenses: [], total: 0 },
        balanceReport: data.BalanceReport || data.balanceReport || { report: [], total: 0 }
      });

      return transformResponse(data);
    } catch (error) {
      console.error('Error fetching event details:', error);
      throw error;
    }
  },

  // Create a new event
  createEvent: async (eventData: EventCreateRequest) => {
    try {
      const response = await api.post('/events', eventData);
      return response.data;
    } catch (error) {
      console.error('Error creating event:', error);
      throw error;
    }
  },

  // Delete an event
  deleteEvent: async (eventId: number) => {
    try {
      const response = await api.delete(`/events/${eventId}`);
      return response.data;
    } catch (error) {
      console.error('Error deleting event:', error);
      throw error;
    }
  },

  // Add a participant to an event
  addParticipant: async (eventId: number, participantData: EventParticipantCreateRequest) => {
    try {
      const response = await api.post(`/events/${eventId}/participants`, participantData);
      return response.data;
    } catch (error) {
      console.error('Error adding participant:', error);
      throw error;
    }
  },

  // Remove a participant from an event
  removeParticipant: async (participantId: number) => {
    try {
      const response = await api.delete(`/events/participants/${participantId}`);
      return response.data;
    } catch (error) {
      console.error('Error removing participant:', error);
      throw error;
    }
  },

  // Get user's events they participate in
  getUserEvents: async (page: number = 1, size: number = 10): Promise<EventParticipantsResponse> => {
    try {
      const response = await api.get(`/participants/user?page=${page}&size=${size}`);
      return response.data;
    } catch (error) {
      console.error('Error fetching user events:', error);
      throw error;
    }
  }
};
