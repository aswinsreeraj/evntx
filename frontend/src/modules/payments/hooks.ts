import { useMutation } from "@tanstack/react-query";
import { paymentsApi } from "./api";
import type { RazorpayVerificationPayload } from "./types";

export const useCreateRazorpayOrder = () => {
  return useMutation({
    mutationFn: (bookingId: string) => paymentsApi.createRazorpayOrder(bookingId),
  });
};

export const useVerifyRazorpayPayment = () => {
  return useMutation({
    mutationFn: (payload: RazorpayVerificationPayload) =>
      paymentsApi.verifyRazorpayPayment(payload),
  });
};
