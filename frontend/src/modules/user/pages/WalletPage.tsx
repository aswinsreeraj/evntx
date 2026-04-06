import { useState } from "react";
import { useWallet, useWalletTransactions, walletQueryKey, walletTransactionsQueryKey } from "../hooks";
import PayoutModal, { type PayoutFormData } from "../../../shared/components/PayoutModal";
import PayoutCredentialsModal from "../../../shared/components/PayoutCredentialsModal";
import AddFundModal, { type AddFundFormData } from "../../../shared/components/AddFundModal";
import { userApi } from "../api";
import { useQueryClient } from "@tanstack/react-query";

function formatCurrency(amount: number) {
  return new Intl.NumberFormat("en-IN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

const formatTransactionDate = (value: string) =>
  new Intl.DateTimeFormat("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value));

const formatReferenceType = (value: string) =>
  value
    .split("_")
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");

function WalletMetric({
  label,
  value,
  helper,
}: {
  label: string;
  value: number;
  helper?: string;
}) {
  return (
    <div className="rounded-3xl border border-gray-100 bg-[#f8fafc] p-6">
      <p className="text-sm font-medium uppercase tracking-[0.18em] text-[#6b7280]">{label}</p>
      <p className="mt-4 text-3xl font-semibold text-[#111827]">₹{formatCurrency(value)}</p>
      {helper ? <p className="mt-2 text-sm text-[#6b7280]">{helper}</p> : null}
    </div>
  );
}

const TransactionRow = ({ transaction }: { transaction: any }) => {
  const [expanded, setExpanded] = useState(false);
  const isCredit = transaction.type === "cr";

  const toggleExpand = () => setExpanded(!expanded);

  const contextHasDetails = !!transaction.context && !!transaction.context.details;

  return (
    <div className="group border-b border-gray-100 bg-white last:border-b-0 hover:bg-gray-50/50 transition">
      <div 
        onClick={contextHasDetails ? toggleExpand : undefined}
        className={`flex flex-col gap-4 px-6 py-5 md:flex-row md:items-center md:justify-between ${contextHasDetails ? "cursor-pointer" : ""}`}
      >
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-3">
            <span className="text-base font-semibold text-[#111827]">
              {formatReferenceType(transaction.reference_type)}
            </span>
            <span
              className={`rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] ${isCredit
                  ? "bg-[#ebfff1] text-[#0f9f4b]"
                  : "bg-[#fff1f3] text-[#e53e5d]"
                }`}
            >
              {isCredit ? "Credit" : "Debit"}
            </span>
            {contextHasDetails && (
               <span className="text-gray-400">
                 {expanded ? (
                   <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" /></svg>
                 ) : (
                   <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" /></svg>
                 )}
               </span>
            )}
          </div>
          <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-[#6b7280]">
            <span>{formatTransactionDate(transaction.created_at)}</span>
            <span className="uppercase tracking-[0.12em]">{transaction.status}</span>
          </div>
        </div>

        <div className="text-left md:text-right">
          <div
            className={`text-xl font-semibold ${isCredit ? "text-[#0f9f4b]" : "text-[#e53e5d]"
              }`}
          >
            {isCredit ? "+" : "-"}₹{formatCurrency(transaction.amount)}
          </div>
          <div className="mt-1 text-xs font-medium uppercase tracking-[0.16em] text-[#94a3b8]">
            {transaction.reference_id}
          </div>
        </div>
      </div>
      
      {expanded && contextHasDetails && (
        <div className="px-6 pb-5 pt-0 text-sm">
           <div className="p-4 bg-gray-50 rounded-2xl border border-gray-100 text-gray-700">
              {["purchase", "user_cancellation", "organizer_cancellation", "earning", "refund", "booking"].includes(transaction.context.type) ? (
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div>
                    <div className="font-semibold text-gray-900 mb-1">Event Details</div>
                    {transaction.context.details?.event?.title || "Unknown Event"}<br/>
                    <span className="text-gray-500">
                      {transaction.context.details?.event?.city ? transaction.context.details.event.city + " • " : ""}
                      {transaction.context.details?.event?.start_time ? new Date(transaction.context.details.event.start_time).toLocaleString("en-IN", {
                         day: "numeric", month: "short", year: "numeric", hour: "numeric", minute: "2-digit"
                      }) : ""}
                    </span>
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1">Ticket Summary</div>
                    {transaction.context.details?.tickets && transaction.context.details.tickets.length > 0 ? (
                       <ul className="list-disc list-inside">
                         {transaction.context.details.tickets.map((t: any, i: number) => (
                           <li key={i}>{t.quantity} × {t.name}</li>
                         ))}
                       </ul>
                    ) : (
                      "No tickets data available"
                    )}
                  </div>
                  <div>
                    <div className="font-semibold text-gray-900 mb-1">Booking Info</div>
                    Status: <span className="uppercase">{transaction.context.details?.status || "Unknown"}</span><br/>
                    Total Amount: ₹{formatCurrency(transaction.context.details?.total_amount || 0)}
                  </div>
                </div>
              ) : transaction.context.type === "payout" ? (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                   <div>
                     <div className="font-semibold text-gray-900 mb-1">Payout Status</div>
                     <span className="uppercase font-medium">{transaction.context.details?.status}</span>
                   </div>
                   <div>
                     <div className="font-semibold text-gray-900 mb-1">Processed At</div>
                     {transaction.context.details?.processed_at ? new Date(transaction.context.details.processed_at).toLocaleString("en-IN", {
                        day: "numeric", month: "short", year: "numeric", hour: "numeric", minute: "2-digit"
                     }) : "Pending"}
                   </div>
                </div>
              ) : (
                <div className="text-gray-500 italic">Contextual details formatted layout not available for type: {transaction.context.type}.</div>
              )}
           </div>
        </div>
      )}
    </div>
  );
};

export default function WalletPage() {
  const [page, setPage] = useState(1);
  const limit = 10;
  const { data: wallet, isLoading, isError, refetch, error } = useWallet();
  const {
    data: transactionsData,
    isLoading: transactionsLoading,
    isError: transactionsError,
    refetch: refetchTransactions,
    error: transactionsErrorDetail,
  } = useWalletTransactions({ page, limit });
  const transactions = transactionsData?.transactions ?? [];
  const pagination = transactionsData?.pagination;
  const totalPages = pagination ? Math.max(1, Math.ceil(pagination.total / pagination.limit)) : 1;
  const queryClient = useQueryClient();
  const [isPayoutModalOpen, setPayoutModalOpen] = useState(false);
  const [isAddFundModalOpen, setAddFundModalOpen] = useState(false);
  const [isPayoutCredentialsModalOpen, setPayoutCredentialsModalOpen] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);

  const handleAddFundSubmit = async (data: AddFundFormData) => {
    setActionError(null);
    try {
      const order = await userApi.createAddFundOrder(data.amount);
      console.log("order", order);
      const options = {
        key: order.razorpay_key,
        amount: order.amount,
        currency: order.currency,
        name: "EvntX",
        description: "Add funds to wallet",
        order_id: order.id,
        handler: async (response: any) => {
          try {
            await userApi.verifyAddFundPayment({
              razorpay_order_id: response.razorpay_order_id,
              razorpay_payment_id: response.razorpay_payment_id,
              razorpay_signature: response.razorpay_signature,
            });
            await queryClient.invalidateQueries({ queryKey: walletQueryKey });
            await queryClient.invalidateQueries({ queryKey: walletTransactionsQueryKey });
          } catch (verifyErr: any) {
            setActionError(verifyErr?.response?.data?.error?.message || "Failed to verify transaction.");
          }
        },
        theme: {
          color: "#111827",
        },
      };

      const razorpay = new (window as any).Razorpay(options);
      razorpay.on("payment.failed", (res: any) => {
        setActionError(res.error.description || "Payment failed.");
      });
      razorpay.open();
    } catch (err: any) {
      console.log(err)
      setActionError(err?.response?.data?.error?.message || "Failed to initiate add fund.");
      throw err;
    }
  };

  const handlePayoutSubmit = async (data: PayoutFormData) => {
    setActionError(null);
    try {
      await userApi.requestPayout(data.amount);
      await queryClient.invalidateQueries({ queryKey: walletQueryKey });
      await queryClient.invalidateQueries({ queryKey: walletTransactionsQueryKey });
    } catch (err: any) {
      setActionError(err?.response?.data?.error?.message || "Failed to submit payout request.");
      throw err;
    }
  };

  const handleCredentialsSubmit = async (data: any) => {
    setActionError(null);
    try {
      await userApi.addPayoutCredentials(data);
    } catch (err: any) {
      setActionError(err?.response?.data?.error?.message || "Failed to save payout credentials.");
      throw err;
    }
  };



  return (
    <div className="min-h-screen bg-gray-50 px-5 py-10">
      <div className="mx-auto max-w-5xl rounded-3xl border border-gray-100 bg-white p-8 shadow-sm md:p-10">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between border-b border-gray-100 pb-6">
          <div className="flex flex-col gap-2">
            <p className="text-sm font-semibold uppercase tracking-[0.24em] text-[#ff445d]">Wallet</p>
            <h1 className="text-3xl font-semibold text-[#111827]">Your wallet overview</h1>
            <p className="text-sm text-[#6b7280]">Track your available and pending balances in one place.</p>
          </div>
          <div className="flex gap-3 mt-4 sm:mt-0">
            <button
              type="button"
              onClick={() => setAddFundModalOpen(true)}
              className="rounded-full border border-gray-200 bg-white px-6 py-2.5 text-sm font-medium text-[#111827] transition hover:bg-gray-50 w-full sm:w-auto"
            >
              Add Fund
            </button>
            <button
              type="button"
              onClick={() => setPayoutCredentialsModalOpen(true)}
              className="rounded-full border border-gray-200 bg-white px-6 py-2.5 text-sm font-medium text-[#111827] transition hover:bg-gray-50 w-full sm:w-auto"
            >
              Payout Settings
            </button>
            <button
              type="button"
              onClick={() => setPayoutModalOpen(true)}
              className="rounded-full bg-[#111827] px-6 py-2.5 text-sm font-medium text-white transition hover:bg-black w-full sm:w-auto"
            >
              Request Payout
            </button>
          </div>
        </div>

        {actionError && (
          <div className="mt-6 rounded-2xl border border-[#ffd7dd] bg-[#fff5f7] px-6 py-4 text-sm font-medium text-[#d22d4c]">
            {actionError}
          </div>
        )}

        {isLoading ? (
          <div className="flex min-h-[360px] items-center justify-center">
            <div className="h-10 w-10 animate-spin rounded-full border-2 border-gray-200 border-t-[#111827]" />
          </div>
        ) : isError ? (
          <div className="mt-8 rounded-3xl border border-[#ffd7dd] bg-[#fff5f7] p-8 text-center">
            <h2 className="text-xl font-semibold text-[#111827]">Unable to load wallet</h2>
            <p className="mt-3 text-sm text-[#6b7280]">
              {error instanceof Error ? error.message : "Please try again in a moment."}
            </p>
            <button
              type="button"
              onClick={() => void refetch()}
              className="mt-6 rounded-full bg-[#111827] px-5 py-2.5 text-sm font-medium text-white transition hover:bg-black"
            >
              Retry
            </button>
          </div>
        ) : wallet ? (
          <>
            <div className="mt-8 grid gap-5 md:grid-cols-2">
              <WalletMetric label="Available Balance" value={wallet.available_balance} />
              <WalletMetric
                label="Pending Balance"
                value={wallet.pending_balance}
                helper="Amount to be settled"
              />
              <WalletMetric label="Total Credited" value={wallet.total_credited} />
              <WalletMetric label="Total Debited" value={wallet.total_debited} />
            </div>

            <div className="mt-10 border-t border-gray-100 pt-8">
              <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
                <div>
                  <h2 className="text-2xl font-semibold text-[#111827]">Transactions</h2>
                  <p className="mt-1 text-sm text-[#6b7280]">
                    View every completed credit and debit recorded in your wallet ledger.
                  </p>
                </div>
                <div className="text-sm text-[#6b7280]">
                  {pagination ? `${pagination.total} transaction${pagination.total === 1 ? "" : "s"}` : "0 transactions"}
                </div>
              </div>

              {transactionsLoading ? (
                <div className="mt-6 flex min-h-[240px] items-center justify-center rounded-3xl border border-gray-100 bg-[#f8fafc]">
                  <div className="h-10 w-10 animate-spin rounded-full border-2 border-gray-200 border-t-[#111827]" />
                </div>
              ) : transactionsError ? (
                <div className="mt-6 rounded-3xl border border-[#ffd7dd] bg-[#fff5f7] p-8 text-center">
                  <h3 className="text-lg font-semibold text-[#111827]">Unable to load transactions</h3>
                  <p className="mt-3 text-sm text-[#6b7280]">
                    {transactionsErrorDetail instanceof Error
                      ? transactionsErrorDetail.message
                      : "Please try again in a moment."}
                  </p>
                  <button
                    type="button"
                    onClick={() => void refetchTransactions()}
                    className="mt-6 rounded-full bg-[#111827] px-5 py-2.5 text-sm font-medium text-white transition hover:bg-black"
                  >
                    Retry
                  </button>
                </div>
              ) : transactions.length === 0 ? (
                <div className="mt-6 rounded-3xl border border-dashed border-gray-200 bg-[#f8fafc] px-6 py-12 text-center">
                  <h3 className="text-lg font-semibold text-[#111827]">No transactions yet</h3>
                  <p className="mt-2 text-sm text-[#6b7280]">
                    Wallet activity will appear here once credits or debits are recorded.
                  </p>
                </div>
              ) : (
                <div className="mt-6 overflow-hidden rounded-3xl border border-gray-100">
                  <div className="divide-y divide-gray-100 bg-white">
                    {transactions.map((transaction) => (
                      <TransactionRow key={transaction.id} transaction={transaction} />
                    ))}
                  </div>

                  <div className="flex flex-col gap-3 border-t border-gray-100 bg-[#fcfcfd] px-6 py-4 sm:flex-row sm:items-center sm:justify-between">
                    <div className="text-sm text-[#6b7280]">
                      Showing {transactions.length} of {pagination?.total ?? 0} transactions
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        type="button"
                        onClick={() => setPage((current) => Math.max(1, current - 1))}
                        disabled={page === 1}
                        className="rounded-full border border-gray-200 px-4 py-2 text-sm font-medium text-[#111827] transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
                      >
                        Prev
                      </button>
                      <span className="text-sm font-medium text-[#111827]">
                        Page {page} of {totalPages}
                      </span>
                      <button
                        type="button"
                        onClick={() => setPage((current) => Math.min(totalPages, current + 1))}
                        disabled={page >= totalPages}
                        className="rounded-full border border-gray-200 px-4 py-2 text-sm font-medium text-[#111827] transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
                      >
                        Next
                      </button>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </>
        ) : null}
      </div>

      <PayoutModal
        isOpen={isPayoutModalOpen}
        onClose={() => setPayoutModalOpen(false)}
        onSubmit={handlePayoutSubmit}
        maxAmount={wallet?.available_balance || 0}
      />

      <PayoutCredentialsModal
        isOpen={isPayoutCredentialsModalOpen}
        onClose={() => setPayoutCredentialsModalOpen(false)}
        onSubmit={handleCredentialsSubmit}
      />

      <AddFundModal
        isOpen={isAddFundModalOpen}
        onClose={() => setAddFundModalOpen(false)}
        onSubmit={handleAddFundSubmit}
      />
    </div>
  );
}
