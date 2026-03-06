import api from "../../services/axios";

export const eventsApi = {
  async listEvents(params?: {
    page?: number;
    limit?: number;
    city?: string;
    category?: string;
  }) {
    const response = await api.get("/events", { params });
    return response.data.data;
  },

  async getEvent(eventId: string) {
    const response = await api.get(`/events/${eventId}`);
    return response.data.data;
  },
};