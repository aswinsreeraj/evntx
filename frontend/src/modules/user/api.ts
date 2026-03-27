import api from "../../services/axios";

export type UpdateProfilePayload = {
  name: string;
  mobile: string;
  dob: string;
  gender: string;
  locations: string[];
  organization_name?: string;
  address?: string;
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
  coverImageUrl: string;
  venue: string;
  tags: string[];
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

export type WalletData = {
  available_balance: number;
  pending_balance: number;
  total_credited: number;
  total_debited: number;
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

  async getWallet(): Promise<WalletData> {
    const res = await api.get("/users/me/wallet");
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

  async getMyTickets(eventId?: string, bookingId?: string, status?: string): Promise<UserTicket[]> {
    const res = await api.get("/users/me/tickets", {
      params: {
        event_id: eventId,
        booking_id: bookingId,
        status,
      },
    });
    return res.data.data.tickets;
  },

  async cancelBooking(bookingId: string, payload: { items: { ticket_type: string; quantity: number }[] }): Promise<void> {
    await api.post(`/bookings/${bookingId}/cancel`, payload);
  },
};
