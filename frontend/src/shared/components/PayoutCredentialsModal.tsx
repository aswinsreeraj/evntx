import { useState, useEffect } from "react";
import type { PayoutCredentialPayload } from "../../modules/user/api";

interface PayoutCredentialsModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: PayoutCredentialPayload) => Promise<void>;
}

export default function PayoutCredentialsModal({ isOpen, onClose, onSubmit }: PayoutCredentialsModalProps) {
  const [accountName, setAccountName] = useState("");
  const [accountNumber, setAccountNumber] = useState("");
  const [ifsc, setIfsc] = useState("");
  const [upiId, setUpiId] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (isOpen) {
      setAccountName("");
      setAccountNumber("");
      setIfsc("");
      setUpiId("");
      setError("");
    }
  }, [isOpen]);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!accountName.trim() || !accountNumber.trim() || !ifsc.trim()) {
      setError("Please fill in all required bank details.");
      return;
    }

    setIsSubmitting(true);
    try {
      await onSubmit({
        account_holder_name: accountName.trim(),
        account_number: accountNumber.trim(),
        ifsc_code: ifsc.trim(),
        upi_id: upiId.trim() || undefined,
      });
      onClose();
    } catch (err: any) {
      setError(err?.response?.data?.error?.message || "Failed to save payout credentials. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4 animate-in fade-in duration-200">
      <div className="w-full max-w-md rounded-3xl bg-white p-6 md:p-8 shadow-xl">
        <h2 className="text-2xl font-bold text-[#111827]">Payout Details</h2>
        <p className="mt-2 text-sm text-[#6b7280]">
          Enter your bank account details securely. These will be used for your automated payouts.
        </p>

        {error && (
          <div className="mt-4 rounded-xl border border-[#ffd7dd] bg-[#fff5f7] p-3 text-sm text-[#e53e5d]">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="mt-6 flex flex-col gap-4">
          <div>
            <label className="mb-1 block text-sm font-medium text-[#111827]">Account Holder Name *</label>
            <input
              type="text"
              value={accountName}
              onChange={(e) => setAccountName(e.target.value)}
              disabled={isSubmitting}
              className="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm outline-none focus:border-[#111827] focus:bg-white transition"
              placeholder="e.g. John Doe"
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-[#111827]">Account Number *</label>
            <input
              type="text"
              value={accountNumber}
              onChange={(e) => setAccountNumber(e.target.value)}
              disabled={isSubmitting}
              className="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm outline-none focus:border-[#111827] focus:bg-white transition flex-1"
              placeholder="Bank Account Number"
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-[#111827]">IFSC Code *</label>
            <input
              type="text"
              value={ifsc}
              onChange={(e) => setIfsc(e.target.value)}
              disabled={isSubmitting}
              className="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm outline-none focus:border-[#111827] focus:bg-white transition uppercase font-mono"
              placeholder="e.g. SBIN0001234"
            />
          </div>

          <div>
            <label className="mb-1 block text-sm font-medium text-[#111827]">UPI ID (Optional)</label>
            <input
              type="text"
              value={upiId}
              onChange={(e) => setUpiId(e.target.value)}
              disabled={isSubmitting}
              className="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm outline-none focus:border-[#111827] focus:bg-white transition"
              placeholder="e.g. user@okhdfc"
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
                "Save Credentials"
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
