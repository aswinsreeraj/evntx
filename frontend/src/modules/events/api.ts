import api from "../../services/axios";
import { organizerApi } from "../organizer/api";
import { adminApi } from "../admin/api";

export const eventsApi = {
  async listEvents(params?: {
    page?: number;
    limit?: number;
    city?: string;
    category?: string;
    search?: string;
    sort?: string;
    start_date?: string;
    end_date?: string;
    min_price?: number;
    max_price?: number;
  }) {
    const response = await api.get("/events", { params });
    return response.data.data;
  },

  async getEvent(eventId: string) {
    const response = await api.get(`/events/${eventId}`);
    return response.data.data;
  },

  async getEventForViewer(eventId: string, isOrganizer: boolean, isAdmin: boolean = false) {
    if (isAdmin) {
      try {
        return await adminApi.getEventBySlug(eventId);
      } catch {
        return this.getEvent(eventId);
      }
    }

    if (!isOrganizer) {
      return this.getEvent(eventId);
    }

    try {
      return await organizerApi.getEventBySlug(eventId);
    } catch {
      return this.getEvent(eventId);
    }
  },

  async reserveTickets(payload: {
    eventId: string;
    tickets: Array<{
      ticket_type_id: string;
      quantity: number;
    }>;
  }) {
    const response = await api.post("/bookings/reserve", {
      event_id: payload.eventId,
      tickets: payload.tickets,
    });
    return response.data.data;
  },
};
