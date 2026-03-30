import { useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import OrganizerLayout from "../components/OrganizerLayout";
import { organizerApi, organizerWalletSummaryQueryKey } from "../api";
import { userApi } from "../../user/api";
import { walletTransactionsQueryKey } from "../../user/hooks";
import PayoutModal, { type PayoutFormData } from "../../../shared/components/PayoutModal";
import AddFundModal, { type AddFundFormData } from "../../../shared/components/AddFundModal";

function formatCurrency(amount: number) {
  return new Intl.NumberFormat("en-IN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

function formatTransactionDate(value: string) {
  return new Intl.DateTimeFormat("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value));
}

function formatReferenceType(value: string) {
  return value
    .split("_")
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function SummaryCard({
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

export default function OrganizerWalletPage() {
  const queryClient = useQueryClient();
  const [page, setPage] = useState(1);
  const [isPayoutModalOpen, setPayoutModalOpen] = useState(false);
  const [isAddFundModalOpen, setAddFundModalOpen] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);
  const limit = 10;

  const {
    data: wallet,
    isLoading,
    isError,
    error,
    refetch,
  } = useQuery({
    queryKey: organizerWalletSummaryQueryKey,
    queryFn: () => organizerApi.getWalletSummary(),
  });

  const {
    data: transactionsData,
    isLoading: transactionsLoading,
    isError: transactionsError,
    error: transactionsErrorDetail,
    refetch: refetchTransactions,
  } = useQuery({
    queryKey: [...walletTransactionsQueryKey, "organizer", page, limit],
    queryFn: () => userApi.getWalletTransactions({ page, limit }),
  });

  const handleAddFundSubmit = async (data: AddFundFormData) => {
    setActionError(null);
    try {
      const order = await userApi.createAddFundOrder(data.amount);

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
            await queryClient.invalidateQueries({ queryKey: organizerWalletSummaryQueryKey });
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
      setActionError(err?.response?.data?.error?.message || "Failed to initiate add fund.");
      throw err;
    }
  };

  const handlePayoutSubmit = async (data: PayoutFormData) => {
    setActionError(null);
    await organizerApi.requestPayout({
      amount: data.amount,
      account_name: data.account,
      account_number: data.accountNumber,
      ifsc_code: data.ifsc,
    });
    await Promise.all([
      queryClient.invalidateQueries({ queryKey: organizerWalletSummaryQueryKey }),
      queryClient.invalidateQueries({ queryKey: walletTransactionsQueryKey }),
    ]);
  };

  const transactions = transactionsData?.transactions ?? [];
  const pagination = transactionsData?.pagination;
  const totalPages = pagination ? Math.max(1, Math.ceil(pagination.total / pagination.limit)) : 1;

  return (
    <OrganizerLayout activeTab="Wallet">
      <div className="px-8 py-10">
        <div className="rounded-3xl border border-gray-100 bg-white p-8 shadow-sm md:p-10">
          <div className="border-b border-gray-100 pb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.24em] text-[#ff445d]">Wallet</p>
              <h1 className="mt-2 text-3xl font-semibold text-[#111827]">Organizer wallet</h1>
              <p className="mt-2 text-sm text-[#6b7280]">Pending balance is the amount to be settled.</p>
            </div>
            <div className="mt-4 sm:mt-0 flex gap-3">
              <button
                type="button"
                onClick={() => setAddFundModalOpen(true)}
                className="rounded-full border border-gray-200 bg-white px-6 py-2.5 text-sm font-medium text-[#111827] transition hover:bg-gray-50 focus:outline-none focus:ring-4 focus:ring-gray-100 whitespace-nowrap"
              >
                Add Fund
              </button>
            </div>
          </div>

          {isLoading ? (
            <div className="flex min-h-[280px] items-center justify-center">
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
              <div className="mt-8 grid gap-5 lg:grid-cols-4">
                <SummaryCard label="Available Balance" value={wallet.available_balance} />
                <SummaryCard
                  label="Pending Balance"
                  value={wallet.pending_balance}
                  helper="Amount to be settled"
                />
                <SummaryCard
                  label="Reserve Balance"
                  value={wallet.reserve_balance}
                  helper="Refund reserves"
                />
                <SummaryCard label="Total Earnings" value={wallet.total_credited - wallet.total_debited} />
              </div>

              <div className="mt-8 rounded-3xl border border-gray-100 bg-[#f8fafc] p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div className="flex flex-col gap-2">
                  <h3 className="text-lg font-semibold text-[#111827]">Wallet Payout</h3>
                  <p className="text-sm text-[#6b7280]">
                    Transfer your available balance to your registered bank account.
                  </p>
                </div>

                <div className="flex gap-3">
                  <button
                    type="button"
                    onClick={() => setPayoutModalOpen(true)}
                    className="rounded-full bg-[#111827] px-6 py-2.5 text-sm font-medium text-white transition hover:bg-[#1f2937] focus:outline-none focus:ring-4 focus:ring-gray-100 whitespace-nowrap"
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

              <div className="mt-10 border-t border-gray-100 pt-8">
                <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
                  <div>
                    <h2 className="text-2xl font-semibold text-[#111827]">Transactions</h2>
                    <p className="mt-1 text-sm text-[#6b7280]">
                      Credits include earnings, and debits include payout requests.
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
                      Earnings and payouts will appear here once they are recorded.
                    </p>
                  </div>
                ) : (
                  <div className="mt-6 overflow-hidden rounded-3xl border border-gray-100">
                    <div className="divide-y divide-gray-100 bg-white">
                      {transactions.map((transaction) => {
                        const isCredit = transaction.type === "cr";

                        return (
                          <div
                            key={transaction.id}
                            className="flex flex-col gap-4 px-6 py-5 md:flex-row md:items-center md:justify-between"
                          >
                            <div className="min-w-0">
                              <div className="flex flex-wrap items-center gap-3">
                                <span className="text-base font-semibold text-[#111827]">
                                  {formatReferenceType(transaction.reference_type)}
                                </span>
                                <span
                                  className={`rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] ${
                                    isCredit
                                      ? "bg-[#ebfff1] text-[#0f9f4b]"
                                      : "bg-[#fff1f3] text-[#e53e5d]"
                                  }`}
                                >
                                  {isCredit ? "Credit" : "Debit"}
                                </span>
                              </div>
                              <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-[#6b7280]">
                                <span>{formatTransactionDate(transaction.created_at)}</span>
                                <span className="uppercase tracking-[0.12em]">{transaction.status}</span>
                              </div>
                            </div>

                            <div className="text-left md:text-right">
                              <div
                                className={`text-xl font-semibold ${
                                  isCredit ? "text-[#0f9f4b]" : "text-[#e53e5d]"
                                }`}
                              >
                                {isCredit ? "+" : "-"}₹{formatCurrency(transaction.amount)}
                              </div>
                              <div className="mt-1 text-xs uppercase tracking-[0.16em] text-[#94a3b8]">
                                {transaction.reference_id}
                              </div>
                            </div>
                          </div>
                        );
                      })}
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
      </div>

      <PayoutModal
        isOpen={isPayoutModalOpen}
        onClose={() => setPayoutModalOpen(false)}
        onSubmit={handlePayoutSubmit}
        maxAmount={Math.max(0, (wallet?.available_balance || 0) - Math.abs(Math.min(0, wallet?.reserve_balance || 0)))}
      />

      <AddFundModal
        isOpen={isAddFundModalOpen}
        onClose={() => setAddFundModalOpen(false)}
        onSubmit={handleAddFundSubmit}
      />
    </OrganizerLayout>
  );
}
