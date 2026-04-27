import { useState } from "react";
import { usePlatformWallet, usePlatformTransactions } from "../hooks";
import AdminLayout from "../components/AdminLayout";
import {
  Wallet,
  TrendingUp,
  TrendingDown,
  RefreshCw,
  ArrowRightLeft,
  ChevronDown,
  ChevronUp
} from "lucide-react";
import type { PlatformWalletTransaction } from "../api";

function StatCard({
  title,
  value,
  icon: Icon,
  color,
  sub,
}: {
  title: string;
  value: number | undefined;
  icon: React.ElementType;
  color: string;
  sub?: string;
}) {
  const formatted =
    value !== undefined
      ? new Intl.NumberFormat("en-IN", {
          style: "currency",
          currency: "INR",
          minimumFractionDigits: 2,
        }).format(value)
      : "—";

  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <span className="text-sm font-semibold text-gray-500">{title}</span>
        <span className={`p-2 rounded-xl ${color}`}>
          <Icon className="w-5 h-5" />
        </span>
      </div>
      <div>
        <p className="text-3xl font-black text-[#0b101e] tracking-tight">
          {formatted}
        </p>
        {sub && <p className="text-xs text-gray-400 mt-1">{sub}</p>}
      </div>
    </div>
  );
}

const formatTransactionDate = (value: string) =>
  new Intl.DateTimeFormat("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(new Date(value));

const mapReferenceType = (type: string) => {
  switch (type) {
    case "payment":
      return "Booking Fee";
    case "earning":
      return "Platform Fee";
    case "fund_addition":
      return "Fund Addition";
    case "refund":
      return "Refund Issued";
    case "payout":
      return "Payout Disbursed";
    default:
      return type
        .split("_")
        .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
        .join(" ");
  }
};

const TransactionRow = ({ transaction }: { transaction: PlatformWalletTransaction }) => {
  const [expanded, setExpanded] = useState(false);
  const isCredit = transaction.type === "cr";

  return (
    <div className="group border-b border-gray-100 bg-white last:border-b-0 hover:bg-gray-50/50 transition">
      <div 
        onClick={() => setExpanded(!expanded)}
        className="flex flex-col gap-4 px-6 py-5 md:flex-row md:items-center md:justify-between cursor-pointer"
      >
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-3">
            <span className="text-base font-semibold text-[#111827]">
              {mapReferenceType(transaction.reference_type)}
            </span>
            <span
              className={`rounded-full px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] ${
                isCredit ? "bg-[#ebfff1] text-[#0f9f4b]" : "bg-[#fff1f3] text-[#e53e5d]"
              }`}
            >
              {isCredit ? "Credit" : "Debit"}
            </span>
            <span className="text-gray-400">
              {expanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
            </span>
          </div>
          <div className="mt-2 flex flex-wrap items-center gap-x-4 gap-y-2 text-sm text-[#6b7280]">
            <span>{formatTransactionDate(transaction.created_at)}</span>
          </div>
        </div>

        <div className="text-left md:text-right">
          <div className={`text-xl font-semibold ${isCredit ? "text-[#0f9f4b]" : "text-[#e53e5d]"}`}>
            {isCredit ? "+" : "−"}
            {new Intl.NumberFormat("en-IN", {
              style: "currency",
              currency: "INR",
              minimumFractionDigits: 2,
            }).format(transaction.amount)}
          </div>
        </div>
      </div>
      
      {expanded && (
        <div className="px-6 pb-5 pt-0 text-sm">
           <div className="p-4 bg-gray-50 rounded-2xl border border-gray-100 text-gray-700">
              <div className="flex flex-col gap-2">
                <div>
                  <span className="font-semibold text-gray-900 mr-2">Reference ID:</span>
                  <span className="font-mono text-xs">{transaction.reference_id}</span>
                </div>
                <div>
                  <span className="font-semibold text-gray-900 mr-2">Transaction ID:</span>
                  <span className="font-mono text-xs">{transaction.id}</span>
                </div>
              </div>
           </div>
        </div>
      )}
    </div>
  );
};

export default function PlatformWalletPage() {
  const [page, setPage] = useState(1);
  const limit = 10;

  const { data: wallet, isLoading: walletLoading, refetch: refetchWallet, isFetching: isWalletFetching } = usePlatformWallet();
  const { data: txnsData, isLoading: txnsLoading, refetch: refetchTxns, isFetching: isTxnsFetching } = usePlatformTransactions(page, limit);

  const transactions = txnsData?.transactions ?? [];
  const pagination = txnsData?.pagination;
  const totalPages = pagination ? Math.max(1, Math.ceil(pagination.total / pagination.limit)) : 1;

  const isFetching = isWalletFetching || isTxnsFetching;

  const handleRefetch = () => {
    refetchWallet();
    refetchTxns();
  };

  const netBalance = wallet
    ? wallet.total_revenue + wallet.total_fees - wallet.total_refunds
    : undefined;

  return (
    <AdminLayout title="Platform Wallet">
      <div className="space-y-8 mt-4">
        {/* Financial Overview Section */}
        <div>
          <div className="flex items-center justify-between mb-6">
            <div>
              <h3 className="text-lg font-bold text-[#0b101e]">Financial Overview</h3>
              {wallet?.updated_at && (
                <p className="text-xs text-gray-400 mt-1">
                  Last updated:{" "}
                  {new Date(wallet.updated_at).toLocaleString("en-IN", {
                    day: "numeric",
                    month: "short",
                    year: "numeric",
                    hour: "2-digit",
                    minute: "2-digit",
                  })}
                </p>
              )}
            </div>
            <button
              onClick={handleRefetch}
              disabled={isFetching}
              className="flex items-center gap-2 px-4 py-2 rounded-xl border border-gray-200 text-sm font-semibold text-gray-600 hover:bg-gray-50 transition-colors disabled:opacity-50"
            >
              <RefreshCw className={`w-4 h-4 ${isFetching ? "animate-spin" : ""}`} />
              Refresh
            </button>
          </div>

          {walletLoading ? (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {Array.from({ length: 4 }).map((_, i) => (
                <div key={i} className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 h-36 animate-pulse">
                  <div className="h-4 bg-gray-100 rounded w-1/2 mb-4" />
                  <div className="h-8 bg-gray-100 rounded w-3/4" />
                </div>
              ))}
            </div>
          ) : (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              <StatCard
                title="Platform Fees Collected"
                value={wallet?.available_balance}
                icon={Wallet}
                color="bg-green-50 text-green-600"
                sub="Money earned as platform fees"
              />
              <StatCard
                title="Total Revenue"
                value={wallet?.total_revenue}
                icon={TrendingUp}
                color="bg-blue-50 text-blue-600"
                sub="Gross booking payments"
              />
              <StatCard
                title="Total Payouts"
                value={wallet?.total_payouts}
                icon={ArrowRightLeft}
                color="bg-orange-50 text-orange-600"
                sub="Disbursed to organizers/users"
              />
              <StatCard
                title="Total Refunds"
                value={wallet?.total_refunds}
                icon={TrendingDown}
                color="bg-red-50 text-red-600"
                sub="Returned to users"
              />
            </div>
          )}
        </div>

        {/* Balance Summary Section */}
        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
          <h4 className="text-sm font-bold text-gray-700 mb-4">Balance Summary</h4>

          <div className="space-y-4">
            <div className="flex items-center justify-between py-3 border-b border-gray-50">
              <div className="flex items-center gap-3">
                <span className="w-2 h-2 rounded-full bg-blue-500 inline-block" />
                <span className="text-sm font-medium text-gray-700">Total Revenue Collected</span>
              </div>
              <span className="text-sm font-bold text-[#0b101e]">
                {wallet ? new Intl.NumberFormat("en-IN", { style: "currency", currency: "INR" }).format(wallet.total_revenue) : "—"}
              </span>
            </div>

            <div className="flex items-center justify-between py-3 border-b border-gray-50">
              <div className="flex items-center gap-3">
                <span className="w-2 h-2 rounded-full bg-green-500 inline-block" />
                <span className="text-sm font-medium text-gray-700">Platform Fees Earned</span>
              </div>
              <span className="text-sm font-bold text-green-600">
                {wallet ? new Intl.NumberFormat("en-IN", { style: "currency", currency: "INR" }).format(wallet.total_fees) : "—"}
              </span>
            </div>



            <div className="flex items-center justify-between py-3 border-b border-gray-50">
              <div className="flex items-center gap-3">
                <span className="w-2 h-2 rounded-full bg-red-500 inline-block" />
                <span className="text-sm font-medium text-gray-700">Total Refunds</span>
              </div>
              <span className="text-sm font-bold text-[#e53e5d]">
                − {wallet ? new Intl.NumberFormat("en-IN", { style: "currency", currency: "INR" }).format(wallet.total_refunds) : "—"}
              </span>
            </div>

            <div className="flex items-center justify-between py-3">
              <div className="flex items-center gap-3">
                <span className="w-2 h-2 rounded-full bg-gray-900 inline-block" />
                <span className="text-sm font-bold text-gray-900">Net Platform Earnings</span>
              </div>
              <span className={`text-base font-black ${(netBalance ?? 0) >= 0 ? "text-green-600" : "text-[#e53e5d]"}`}>
                {netBalance !== undefined ? new Intl.NumberFormat("en-IN", { style: "currency", currency: "INR" }).format(netBalance) : "—"}
              </span>
            </div>
          </div>
        </div>

        {/* Transaction History Section */}
        <div>
          <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between mb-4">
            <div>
              <h3 className="text-lg font-bold text-[#0b101e]">Platform Ledger</h3>
              <p className="mt-1 text-sm text-[#6b7280]">
                All credits and debits to the platform wallet.
              </p>
            </div>
            <div className="text-sm text-[#6b7280]">
              {pagination ? `${pagination.total} transaction${pagination.total === 1 ? "" : "s"}` : "0 transactions"}
            </div>
          </div>

          {txnsLoading ? (
            <div className="mt-6 flex min-h-[240px] items-center justify-center rounded-3xl border border-gray-100 bg-[#f8fafc]">
              <div className="h-10 w-10 animate-spin rounded-full border-2 border-gray-200 border-t-[#111827]" />
            </div>
          ) : transactions.length === 0 ? (
            <div className="mt-6 rounded-3xl border border-dashed border-gray-200 bg-[#f8fafc] px-6 py-12 text-center">
              <h3 className="text-lg font-semibold text-[#111827]">No transactions yet</h3>
              <p className="mt-2 text-sm text-[#6b7280]">
                Activity will appear here once credits or debits are recorded.
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
                    disabled={page === 1 || txnsLoading}
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
                    disabled={page >= totalPages || txnsLoading}
                    className="rounded-full border border-gray-200 px-4 py-2 text-sm font-medium text-[#111827] transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
                  >
                    Next
                  </button>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </AdminLayout>
  );
}
