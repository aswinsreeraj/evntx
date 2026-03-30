import React, { useState } from "react";
import { useCreateRazorpayOrder, useVerifyRazorpayPayment } from "../hooks";
import type { RazorpayResponse } from "../types";
import { useAuthStore } from "../../auth/store/authStore";
import { CreditCard, Loader2 } from "lucide-react";

interface RazorpayButtonProps {
  bookingId: string;
  eventTitle: string;
  onSuccess: () => void;
  onError: (message: string) => void;
  autoOpen?: boolean;
}

const RazorpayButton: React.FC<RazorpayButtonProps> = ({
  bookingId,
  eventTitle,
  onSuccess,
  onError,
  autoOpen = false,
}) => {
  const [isProcessing, setIsProcessing] = useState(false);
  const { user } = useAuthStore();

  const createOrderMutation = useCreateRazorpayOrder();
  const verifyPaymentMutation = useVerifyRazorpayPayment();

  React.useEffect(() => {
    if (autoOpen && !isProcessing && !createOrderMutation.isSuccess && !verifyPaymentMutation.isSuccess) {
      handlePayment();
    }
  }, [autoOpen]);

  const loadScript = (src: string) => {
    return new Promise((resolve) => {
      const script = document.createElement("script");
      script.src = src;
      script.onload = () => resolve(true);
      script.onerror = () => resolve(false);
      document.body.appendChild(script);
    });
  };

  const handlePayment = async () => {
    setIsProcessing(true);

    try {

      const res = await loadScript("https://checkout.razorpay.com/v1/checkout.js");
      if (!res) {
        throw new Error("Razorpay SDK failed to load. Are you online?");
      }

      const order = await createOrderMutation.mutateAsync(bookingId);

      if (order.is_free_booking) {
        onSuccess();
        return;
      }

      const options = {
        key: import.meta.env.VITE_RAZORPAY_KEY_ID,
        amount: order.amount,
        currency: order.currency,
        name: "EVNTX",
        description: `Payment for ${eventTitle}`,
        order_id: order.id,
        handler: async (response: RazorpayResponse) => {
          try {
            setIsProcessing(true);

            await verifyPaymentMutation.mutateAsync({
              razorpay_order_id: response.razorpay_order_id,
              razorpay_payment_id: response.razorpay_payment_id,
              razorpay_signature: response.razorpay_signature,
            });
            onSuccess();
          } catch (err: any) {
            onError(err?.response?.data?.message || "Payment verification failed.");
          } finally {
            setIsProcessing(false);
          }
        },
        prefill: {
          name: user?.name || "",
        },
        theme: {
          color: "#090c44",
        },
        modal: {
          ondismiss: () => {
            setIsProcessing(false);
          },
        },
      };

      const rzp = new (window as any).Razorpay(options);
      rzp.open();
    } catch (err: any) {
      onError(err?.response?.data?.message || err.message || "Something went wrong.");
      setIsProcessing(false);
    }
  };

  return (
    <button
      type="button"
      onClick={handlePayment}
      disabled={isProcessing || createOrderMutation.isPending || verifyPaymentMutation.isPending}
      className="group relative flex w-full items-center justify-center gap-3 overflow-hidden rounded-xl bg-[#090c44] px-6 py-3.5 text-sm font-semibold text-white shadow-lg transition-all hover:bg-[#06082f] hover:shadow-xl active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-70"
    >
      {(isProcessing || createOrderMutation.isPending || verifyPaymentMutation.isPending) ? (
        <>
          <Loader2 className="h-5 w-5 animate-spin" />
          <span>Processing Payment...</span>
        </>
      ) : (
        <>
          <CreditCard className="h-5 w-5 transition-transform group-hover:scale-110" />
          <span>Pay with Razorpay</span>
        </>
      )}

      {}
      <div className="absolute inset-0 translate-x-[-100%] bg-gradient-to-r from-transparent via-white/10 to-transparent transition-transform duration-1000 group-hover:translate-x-[100%]" />
    </button>
  );
};

export default RazorpayButton;
