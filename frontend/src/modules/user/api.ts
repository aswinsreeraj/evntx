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
  event_status: string;
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

export type WalletTransaction = {
  id: string;
  wallet_id: string;
  type: "cr" | "dr";
  amount: number;
  reference_type: string;
  reference_id: string;
  status: "pending" | "completed" | "failed";
  created_at: string;
};

export type WalletTransactionsResponse = {
  transactions: WalletTransaction[];
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
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

  async getWalletTransactions(params?: {
    page?: number;
    limit?: number;
    type?: "cr" | "dr";
    status?: "pending" | "completed" | "failed";
  }): Promise<WalletTransactionsResponse> {
    const res = await api.get("/users/me/wallet/transactions", { params });
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

  async refundBooking(bookingId: string): Promise<void> {
    await api.post(`/bookings/${bookingId}/refund`);
  },

  async requestPayout(payload: { amount: number; account_name: string; account_number: string; ifsc_code: string }): Promise<void> {
    await api.post("/users/me/wallet/payout", payload);
  },

  async createAddFundOrder(amount: number): Promise<{ id: string; amount: number; currency: string; razorpay_key: string }> {
    const res = await api.post("/users/me/wallet/add-fund", { amount });
    return res.data.data;
  },

  async verifyAddFundPayment(payload: { razorpay_order_id: string; razorpay_payment_id: string; razorpay_signature: string }): Promise<void> {
    await api.post("/users/me/wallet/add-fund/verify", payload);
  },

  async payWithWallet(bookingId: string): Promise<void> {
    await api.post(`/bookings/${bookingId}/pay-with-wallet`);
  },
};
