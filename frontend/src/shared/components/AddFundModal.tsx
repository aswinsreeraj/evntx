import { useState, useEffect } from "react";

export interface AddFundFormData {
  amount: number;
}

interface AddFundModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: AddFundFormData) => Promise<void>;
  isRazorpayEnabled?: boolean;
}

export default function AddFundModal({ isOpen, onClose, onSubmit, isRazorpayEnabled = true }: AddFundModalProps) {
  const [amount, setAmount] = useState<string>("");
  const [error, setError] = useState<string>("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (isOpen) {
      setAmount("");
      setError("");
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    const parsedAmount = parseFloat(amount);
    if (!parsedAmount || parsedAmount <= 0) {
      setError("Please enter a valid positive amount.");
      return;
    }

    setIsSubmitting(true);
    try {
      await onSubmit({ amount: parsedAmount });
      onClose();
    } catch (err: any) {
      setError(err?.response?.data?.error?.message || "Failed to initialize payment. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in duration-200">
      <div className="w-full max-w-md rounded-3xl bg-white p-6 md:p-8 shadow-xl">
        <h2 className="text-2xl font-bold text-[#111827]">Add Funds</h2>
        <p className="mt-2 text-sm text-[#6b7280]">
          Deposit money into your wallet securely via Razorpay.
        </p>

        {error && (
          <div className="mt-4 rounded-xl border border-[#ffd7dd] bg-[#fff5f7] p-3 text-sm text-[#e53e5d]">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="mt-6 flex flex-col gap-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-[#111827]">Amount (₹)</label>
            <input
              type="number"
              step="0.01"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              disabled={isSubmitting}
              className="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm outline-none focus:border-[#111827] focus:bg-white transition"
              placeholder="e.g. 1000.00"
            />
          </div>

          <div className="mt-4 flex gap-3">
            <button
              type="button"
              onClick={onClose}
              disabled={isSubmitting}
              className="flex-1 rounded-xl border border-gray-200 px-4 py-2.5 text-sm font-semibold text-[#111827] hover:bg-gray-50 transition"
            >
              Cancel
            </button>
            {isRazorpayEnabled ? (
              <button
                type="submit"
                disabled={isSubmitting}
                className="flex-1 rounded-xl bg-[#111827] px-4 py-2.5 text-sm font-semibold text-white hover:bg-black transition disabled:opacity-70 flex items-center justify-center"
              >
                {isSubmitting ? (
                  <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
                ) : (
                  "Proceed to Pay"
                )}
              </button>
            ) : (
              <button
                type="button"
                disabled
                className="flex-1 rounded-xl bg-gray-200 px-4 py-2.5 text-sm font-semibold text-gray-400 cursor-not-allowed flex items-center justify-center"
              >
                Razorpay Disabled
              </button>
            )}
          </div>
        </form>
      </div>
    </div>
  );
}
