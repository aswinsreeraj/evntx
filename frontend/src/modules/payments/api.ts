import api from "../../services/axios";
import type { RazorpayOrder, RazorpayVerificationPayload } from "./types";

export const paymentsApi = {
  async createRazorpayOrder(bookingId: string): Promise<RazorpayOrder> {
    const response = await api.post("/payments/razorpay/order", {
      booking_id: bookingId,
    });
    return response.data.data;
  },

  async verifyRazorpayPayment(payload: RazorpayVerificationPayload) {
    const response = await api.post("/payments/razorpay/verify", payload);
    return response.data;
  },
};
