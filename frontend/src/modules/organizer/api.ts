import api from "../../services/axios";

export interface TicketInput {
  id?: string;
  name: string;
  price: number;
  total_quantity: number;
}

export interface PersonnelInput {
  id?: string;
  name: string;
  role: string;
  image?: string;
  profile_link?: string;
}

export interface EventDetailsInput {
  description: string;
  venue_address: string;
  map_url?: string;
  total_capacity: number;
  terms_and_conditions?: string;
}

export interface CreateEventPayload {
  title: string;
  city: string;
  venue_name: string;
  category?: string;
  start_time: string;
  end_time: string;
  tags?: string[];
  cover_image_url?: string;
  details: EventDetailsInput;
  ticket_types: TicketInput[];
  key_personnel?: PersonnelInput[];
}

export interface UpdateEventPayload {
  title?: string;
  city?: string;
  venue_name?: string;
  category?: string;
  start_time?: string;
  end_time?: string;
  tags?: string[];
  cover_image_url?: string;
  details?: Partial<EventDetailsInput>;
  ticket_types?: TicketInput[];
  key_personnel?: PersonnelInput[];
}

export const organizerApi = {
  async createEvent(payload: CreateEventPayload) {
    const res = await api.post("/organizer/events", payload);
    return res.data;
  },

  async updateEvent(eventId: string, payload: UpdateEventPayload) {
    const res = await api.put(`/organizer/events/${eventId}`, payload);
    return res.data;
  },

  async getOrganizerEvents(status?: string) {
    const res = await api.get(`/organizer/events`, { params: { status } });
    return res.data;
  },

  async deleteEvent(eventId: string) {
    const res = await api.delete(`/organizer/events/${eventId}`);
    return res.data;
  },

  async submitEvent(eventId: string) {
    const res = await api.post(`/organizer/events/${eventId}/submit`);
    return res.data;
  },

  async uploadImage(file: File) {
    const formData = new FormData();
    formData.append("image", file);
    const res = await api.post("/organizer/upload", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return res.data.data;
  },
};
