import { useWallet } from "../hooks";

function formatCurrency(amount: number) {
  return new Intl.NumberFormat("en-IN", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);
}

function WalletMetric({
  label,
  value,
}: {
  label: string;
  value: number;
}) {
  return (
    <div className="rounded-3xl border border-gray-100 bg-[#f8fafc] p-6">
      <p className="text-sm font-medium uppercase tracking-[0.18em] text-[#6b7280]">{label}</p>
      <p className="mt-4 text-3xl font-semibold text-[#111827]">₹{formatCurrency(value)}</p>
    </div>
  );
}

export default function WalletPage() {
  const { data: wallet, isLoading, isError, refetch, error } = useWallet();

  return (
    <div className="min-h-screen bg-gray-50 px-5 py-10">
      <div className="mx-auto max-w-5xl rounded-3xl border border-gray-100 bg-white p-8 shadow-sm md:p-10">
        <div className="flex flex-col gap-2 border-b border-gray-100 pb-6">
          <p className="text-sm font-semibold uppercase tracking-[0.24em] text-[#ff445d]">Wallet</p>
          <h1 className="text-3xl font-semibold text-[#111827]">Your wallet overview</h1>
          <p className="text-sm text-[#6b7280]">Track your available and pending balances in one place.</p>
        </div>

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
          <div className="mt-8 grid gap-5 md:grid-cols-2">
            <WalletMetric label="Available Balance" value={wallet.available_balance} />
            <WalletMetric label="Pending Balance" value={wallet.pending_balance} />
            <WalletMetric label="Total Credited" value={wallet.total_credited} />
            <WalletMetric label="Total Debited" value={wallet.total_debited} />
          </div>
        ) : null}
      </div>
    </div>
  );
}
