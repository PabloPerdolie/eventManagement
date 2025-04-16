import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { EventResponse, EventParticipant } from '../../types/api';
import { api } from '../../services/api';

interface EventsState {
  events: EventResponse[];
  currentEvent: EventResponse | null;
  participants: EventParticipant[];
  loading: boolean;
  error: string | null;
  total: number;
}

const initialState: EventsState = {
  events: [],
  currentEvent: null,
  participants: [],
  loading: false,
  error: null,
  total: 0,
};

export const fetchEvents = createAsyncThunk(
  'events/fetchEvents',
  async ({ page = 1, size = 10 }: { page?: number; size?: number }) => {
    const response = await api.get(`/events?page=${page}&size=${size}`);
    return response.data;
  }
);

export const createEvent = createAsyncThunk(
  'events/createEvent',
  async (eventData: {
    title: string;
    description?: string;
    location?: string;
    start_date: string;
    end_date: string;
  }) => {
    const response = await api.post('/events', eventData);
    return response.data;
  }
);

export const fetchEventDetails = createAsyncThunk(
  'events/fetchEventDetails',
  async (eventId: number) => {
    const response = await api.get(`/events/${eventId}`);
    return response.data;
  }
);

export const addParticipant = createAsyncThunk(
  'events/addParticipant',
  async ({ eventId, userId }: { eventId: number; userId: number }) => {
    const response = await api.post(`/events/${eventId}/participants`, {
      user_id: userId,
    });
    return response.data;
  }
);

const eventsSlice = createSlice({
  name: 'events',
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchEvents.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchEvents.fulfilled, (state, action) => {
        state.loading = false;
        state.events = action.payload.events;
        state.total = action.payload.total;
      })
      .addCase(fetchEvents.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch events';
      })
      .addCase(createEvent.fulfilled, (state, action) => {
        state.events.unshift(action.payload);
      })
      .addCase(fetchEventDetails.fulfilled, (state, action) => {
        state.currentEvent = action.payload.eventData;
        state.participants = action.payload.eventParticipants.participants;
      });
  },
});

export default eventsSlice.reducer;