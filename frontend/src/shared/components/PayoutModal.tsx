import { useState, useEffect } from "react";

export interface PayoutFormData {
  amount: number;
}

interface PayoutModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: PayoutFormData) => Promise<void>;
  maxAmount: number;
}

export default function PayoutModal({ isOpen, onClose, onSubmit, maxAmount }: PayoutModalProps) {
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

    if (parsedAmount > maxAmount) {
      setError(`Amount cannot exceed your withdrawable balance (₹${maxAmount}).`);
      return;
    }

    setIsSubmitting(true);
    try {
      await onSubmit({
        amount: parsedAmount,
      });
      onClose();
    } catch (err: any) {
      setError(err?.response?.data?.error?.message || "Failed to submit payout request. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in duration-200">
      <div className="w-full max-w-md rounded-3xl bg-white p-6 md:p-8 shadow-xl">
        <h2 className="text-2xl font-bold text-[#111827]">Request Payout</h2>
        <p className="mt-2 text-sm text-[#6b7280]">
          Transfer funds from your wallet to your bank account. Maximum withdrawable amount: ₹{maxAmount.toFixed(2)}
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
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-1 rounded-xl bg-[#111827] px-4 py-2.5 text-sm font-semibold text-white hover:bg-black transition disabled:opacity-70 flex items-center justify-center"
            >
              {isSubmitting ? (
                <div className="h-4 w-4 animate-spin rounded-full border-2 border-white border-t-transparent" />
              ) : (
                "Submit Request"
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
