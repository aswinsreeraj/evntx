import { usePlatformWallet } from "../hooks";
import AdminLayout from "../components/AdminLayout";
import {
  Wallet,
  TrendingUp,
  TrendingDown,
  Clock,
  RefreshCw,
} from "lucide-react";

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

export default function PlatformWalletPage() {
  const { data, isLoading, refetch, isFetching } = usePlatformWallet();

  const wallet = data
    ? {
        available_balance: data.available_balance ?? 0,
        pending_balance: data.pending_balance ?? 0,
        total_credited: data.total_credited ?? 0,
        total_debited: data.total_debited ?? 0,
        updated_at: data.updated_at,
      }
    : null;

  const netBalance = wallet
    ? wallet.total_credited - wallet.total_debited
    : undefined;

  return (
    <AdminLayout title="Platform Wallet">
      <div className="space-y-6 mt-4">

        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-bold text-[#0b101e]">
              Financial Overview
            </h3>
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
            onClick={() => refetch()}
            disabled={isFetching}
            className="flex items-center gap-2 px-4 py-2 rounded-xl border border-gray-200 text-sm font-semibold text-gray-600 hover:bg-gray-50 transition-colors disabled:opacity-50"
          >
            <RefreshCw
              className={`w-4 h-4 ${isFetching ? "animate-spin" : ""}`}
            />
            Refresh
          </button>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {Array.from({ length: 4 }).map((_, i) => (
              <div
                key={i}
                className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6 h-36 animate-pulse"
              >
                <div className="h-4 bg-gray-100 rounded w-1/2 mb-4" />
                <div className="h-8 bg-gray-100 rounded w-3/4" />
              </div>
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            <StatCard
              title="Available Balance"
              value={wallet?.available_balance}
              icon={Wallet}
              color="bg-green-50 text-green-600"
              sub="Ready to settle / payout"
            />
            <StatCard
              title="Pending Balance"
              value={wallet?.pending_balance}
              icon={Clock}
              color="bg-yellow-50 text-yellow-600"
              sub="Awaiting settlement"
            />
            <StatCard
              title="Total Revenue"
              value={wallet?.total_credited}
              icon={TrendingUp}
              color="bg-blue-50 text-blue-600"
              sub="All-time gross payments"
            />
            <StatCard
              title="Total Payouts & Refunds"
              value={wallet?.total_debited}
              icon={TrendingDown}
              color="bg-red-50 text-red-600"
              sub="All-time disbursements"
            />
          </div>
        )}

        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm p-6">
          <h4 className="text-sm font-bold text-gray-700 mb-4">
            Balance Summary
          </h4>

          <div className="space-y-4">
            <div className="flex items-center justify-between py-3 border-b border-gray-50">
              <div className="flex items-center gap-3">
                <span className="w-2 h-2 rounded-full bg-blue-500 inline-block" />
                <span className="text-sm font-medium text-gray-700">
                  Total Revenue Collected
                </span>
              </div>
              <span className="text-sm font-bold text-[#0b101e]">
                {wallet
                  ? new Intl.NumberFormat("en-IN", {
                      style: "currency",
                      currency: "INR",
                      minimumFractionDigits: 2,
                    }).format(wallet.total_credited)
                  : "—"}
              </span>
            </div>

            <div className="flex items-center justify-between py-3 border-b border-gray-50">
              <div className="flex items-center gap-3">
                <span className="w-2 h-2 rounded-full bg-red-500 inline-block" />
                <span className="text-sm font-medium text-gray-700">
                  Total Payouts &amp; Refunds
                </span>
              </div>
              <span className="text-sm font-bold text-[#e53e5d]">
                −{" "}
                {wallet
                  ? new Intl.NumberFormat("en-IN", {
                      style: "currency",
                      currency: "INR",
                      minimumFractionDigits: 2,
                    }).format(wallet.total_debited)
                  : "—"}
              </span>
            </div>

            <div className="flex items-center justify-between py-3">
              <div className="flex items-center gap-3">
                <span className="w-2 h-2 rounded-full bg-green-500 inline-block" />
                <span className="text-sm font-bold text-gray-900">
                  Net Platform Earnings
                </span>
              </div>
              <span
                className={`text-base font-black ${
                  (netBalance ?? 0) >= 0 ? "text-green-600" : "text-[#e53e5d]"
                }`}
              >
                {netBalance !== undefined
                  ? new Intl.NumberFormat("en-IN", {
                      style: "currency",
                      currency: "INR",
                      minimumFractionDigits: 2,
                    }).format(netBalance)
                  : "—"}
              </span>
            </div>
          </div>
        </div>
      </div>
    </AdminLayout>
  );
}
