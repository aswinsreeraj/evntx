import api from "../../services/axios";

export const organizerWalletSummaryQueryKey = ["organizer-wallet-summary"] as const;

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
  status?: string;
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
  status?: string;
  details?: Partial<EventDetailsInput>;
  ticket_types?: TicketInput[];
  key_personnel?: PersonnelInput[];
}

export interface OrganizerWalletSummary {
  available_balance: number;
  pending_balance: number;
  total_credited: number;
  total_debited: number;
}

export interface CheckInResponse {
  ticket_id: string;
  ticket_code: string;
  status: string;
  checked_in_at: string;
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

  async getWalletSummary(): Promise<OrganizerWalletSummary> {
    const res = await api.get("/organizer/wallet");
    return res.data.data;
  },

  async requestPayout(amount: number) {
    const res = await api.post("/organizer/wallet/payout", { amount });
    return res.data.data;
  },

  async getEventBySlug(slug: string) {
    const res = await api.get(`/organizer/events/slug/${slug}`);
    return res.data.data;
  },

  async deleteEvent(eventId: string) {
    const res = await api.delete(`/organizer/events/${eventId}`);
    return res.data;
  },

  async submitEvent(eventId: string) {
    const res = await api.post(`/organizer/events/${eventId}/submit`);
    return res.data;
  },

  async checkInTicket(eventId: string, ticketCode: string): Promise<CheckInResponse> {
    const res = await api.post(`/events/${eventId}/check-in`, {
      ticket_code: ticketCode,
    });
    return res.data.data;
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
