import api from "../../services/axios";

export type UpdateProfilePayload = {
  name: string;
  mobile: string;
  dob: string;
  gender: string;
  locations: string[];
  organization_name?: string;
};

export type UserBooking = {
  booking_id: string;
  event_id: string;
  event_title: string;
  event_city: string;
  event_start_time: string;
  status: string;
  total_amount: number;
  ticket_count: number;
  created_at: string;
};

export type UserTicket = {
  ticket_id: string;
  ticket_code: string;
  event_id: string;
  event_title: string;
  ticket_type: string;
  status: string;
  checked_in_at?: string | null;
};

export const userApi = {
  async getProfile() {
    const res = await api.get("/users/me");
    return res.data.data;
  },

  async updateProfile(payload: UpdateProfilePayload) {
    const res = await api.put("/users/me", payload);
    return res.data.data;
  },

  async uploadProfileImage(file: File) {
    const formData = new FormData();
    formData.append("profile_image", file);
    const res = await api.post("/users/me/image", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    });
    return res.data.data;
  },

  async getMyBookings(status?: string): Promise<UserBooking[]> {
    const res = await api.get("/users/me/bookings", {
      params: { status },
    });
    return res.data.data.bookings;
  },

  async getMyTickets(eventId?: string, status?: string): Promise<UserTicket[]> {
    const res = await api.get("/users/me/tickets", {
      params: {
        event_id: eventId,
        status,
      },
    });
    return res.data.data.tickets;
  },
};
