import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { adminApi, type AdminPayoutDetail } from "../api";
import AdminLayout from "../components/AdminLayout";
import { CheckCircle, XCircle, RefreshCw, ChevronDown, ChevronUp } from "lucide-react";

const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    minimumFractionDigits: 2,
  }).format(amount);

const formatDate = (value?: string) =>
  value
    ? new Intl.DateTimeFormat("en-IN", {
        day: "2-digit",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      }).format(new Date(value))
    : "—";

const STATUS_STYLES: Record<string, string> = {
  pending: "bg-yellow-50 text-yellow-700 border-yellow-200",
  approved: "bg-green-50 text-green-700 border-green-200",
  rejected: "bg-red-50 text-red-700 border-red-200",
  processing: "bg-blue-50 text-blue-700 border-blue-200",
  completed: "bg-emerald-50 text-emerald-700 border-emerald-200",
  failed: "bg-gray-50 text-gray-700 border-gray-200",
};

function PayoutRow({
  payout,
  selected,
  onToggleSelect,
  onApprove,
  onReject,
  isApproving,
  isRejecting,
}: {
  payout: AdminPayoutDetail;
  selected: boolean;
  onToggleSelect: () => void;
  onApprove: () => void;
  onReject: (reason: string) => void;
  isApproving: boolean;
  isRejecting: boolean;
}) {
  const [expanded, setExpanded] = useState(false);
  const [rejectReason, setRejectReason] = useState("");
  const [showRejectForm, setShowRejectForm] = useState(false);

  const statusStyle = STATUS_STYLES[payout.status] || "bg-gray-50 text-gray-700 border-gray-200";
  const isPending = payout.status === "pending";

  return (
    <div className="border-b border-gray-100 last:border-b-0">
      <div className="flex items-center gap-4 px-6 py-4 hover:bg-gray-50/50 transition">
        {isPending && (
          <input
            type="checkbox"
            checked={selected}
            onChange={onToggleSelect}
            className="h-4 w-4 rounded border-gray-300 accent-[#111827] cursor-pointer"
          />
        )}
        {!isPending && <div className="w-4" />}

        <div className="flex-1 min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <span className="font-semibold text-[#111827] text-sm">{payout.user_name || "Unknown"}</span>
            <span className="text-xs text-gray-400">{payout.user_email}</span>
          </div>
          <div className="mt-1 flex flex-wrap items-center gap-3 text-xs text-[#6b7280]">
            <span>Requested: {formatDate(payout.requested_at)}</span>
            {payout.reviewed_at && <span>Reviewed: {formatDate(payout.reviewed_at)}</span>}
          </div>
        </div>

        <div className="text-right shrink-0">
          <div className="text-lg font-bold text-[#111827]">{formatCurrency(payout.amount)}</div>
          <span className={`mt-1 inline-block rounded-full border px-2.5 py-0.5 text-xs font-semibold uppercase tracking-wide ${statusStyle}`}>
            {payout.status}
          </span>
        </div>

        {isPending && (
          <div className="flex gap-2 shrink-0">
            <button
              type="button"
              onClick={onApprove}
              disabled={isApproving || isRejecting}
              title="Approve"
              className="rounded-full bg-green-600 p-2 text-white hover:bg-green-700 transition disabled:opacity-50"
            >
              <CheckCircle className="h-4 w-4" />
            </button>
            <button
              type="button"
              onClick={() => setShowRejectForm((v) => !v)}
              disabled={isApproving || isRejecting}
              title="Reject"
              className="rounded-full bg-red-600 p-2 text-white hover:bg-red-700 transition disabled:opacity-50"
            >
              <XCircle className="h-4 w-4" />
            </button>
          </div>
        )}

        <button
          type="button"
          onClick={() => setExpanded((v) => !v)}
          className="ml-2 text-gray-400 hover:text-gray-600 transition shrink-0"
        >
          {expanded ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
        </button>
      </div>

      {showRejectForm && isPending && (
        <div className="px-6 py-4 bg-red-50 border-t border-red-100">
          <p className="text-sm font-medium text-red-700 mb-2">Rejection Reason</p>
          <div className="flex gap-2">
            <input
              type="text"
              value={rejectReason}
              onChange={(e) => setRejectReason(e.target.value)}
              placeholder="Enter reason for rejection..."
              className="flex-1 rounded-xl border border-red-200 bg-white px-4 py-2 text-sm outline-none focus:border-red-400 transition"
            />
            <button
              type="button"
              onClick={() => {
                if (rejectReason.trim()) {
                  onReject(rejectReason.trim());
                  setShowRejectForm(false);
                }
              }}
              disabled={!rejectReason.trim() || isRejecting}
              className="rounded-xl bg-red-600 px-4 py-2 text-sm font-semibold text-white hover:bg-red-700 transition disabled:opacity-50"
            >
              Confirm
            </button>
            <button
              type="button"
              onClick={() => setShowRejectForm(false)}
              className="rounded-xl border border-gray-200 px-4 py-2 text-sm font-semibold text-gray-600 hover:bg-gray-50 transition"
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {expanded && (
        <div className="px-6 py-4 bg-gray-50 border-t border-gray-100 text-sm text-[#6b7280] grid grid-cols-2 gap-x-8 gap-y-2">
          <div><span className="font-semibold text-[#111827]">Payout ID:</span> {payout.id}</div>
          <div><span className="font-semibold text-[#111827]">User ID:</span> {payout.user_id}</div>
          {payout.admin_id && <div><span className="font-semibold text-[#111827]">Admin ID:</span> {payout.admin_id}</div>}
          {payout.failure_reason && (
            <div className="col-span-2">
              <span className="font-semibold text-red-600">Failure Reason:</span> {payout.failure_reason}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

const payoutsQueryKey = (status: string) => ["admin-payouts", status];

export default function PayoutsPage() {
  const queryClient = useQueryClient();
  const [statusFilter, setStatusFilter] = useState("pending");
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [approveTarget, setApproveTarget] = useState<AdminPayoutDetail | null>(null);

  const { data, isLoading, refetch, isFetching } = useQuery({
    queryKey: payoutsQueryKey(statusFilter),
    queryFn: () => adminApi.getPayouts({ status: statusFilter || undefined }),
  });

  const payouts = data?.payouts ?? [];
  const total = data?.total ?? 0;

  const approveMutation = useMutation({
    mutationFn: (id: string) => adminApi.approvePayout(id),
    onSuccess: () => {
      setApproveTarget(null);
      queryClient.invalidateQueries({ queryKey: payoutsQueryKey(statusFilter) });
    },
  });

  const rejectMutation = useMutation({
    mutationFn: ({ id, reason }: { id: string; reason: string }) => adminApi.rejectPayout(id, reason),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: payoutsQueryKey(statusFilter) }),
  });

  const bulkApproveMutation = useMutation({
    mutationFn: (ids: string[]) => adminApi.bulkApprovePayouts(ids),
    onSuccess: () => {
      setSelectedIds(new Set());
      queryClient.invalidateQueries({ queryKey: payoutsQueryKey(statusFilter) });
    },
  });

  const handleToggleSelect = (id: string) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleSelectAll = () => {
    const pendingIds = payouts.filter((p) => p.status === "pending").map((p) => p.id);
    if (selectedIds.size === pendingIds.length) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(pendingIds));
    }
  };

  const pendingPayouts = payouts.filter((p) => p.status === "pending");

  return (
    <AdminLayout title="Payout Requests">
      <div className="space-y-6 mt-4">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h3 className="text-lg font-bold text-[#0b101e]">Payout Management</h3>
            <p className="text-sm text-gray-500 mt-1">
              Review and approve or reject user and organizer payout requests.
            </p>
          </div>
          <div className="flex gap-2">
            {selectedIds.size > 0 && (
              <button
                type="button"
                onClick={() => bulkApproveMutation.mutate([...selectedIds])}
                disabled={bulkApproveMutation.isPending}
                className="flex items-center gap-2 rounded-xl bg-green-600 px-4 py-2 text-sm font-semibold text-white hover:bg-green-700 transition disabled:opacity-50"
              >
                <CheckCircle className="h-4 w-4" />
                Approve Selected ({selectedIds.size})
              </button>
            )}
            <button
              type="button"
              onClick={() => refetch()}
              disabled={isFetching}
              className="flex items-center gap-2 rounded-xl border border-gray-200 px-4 py-2 text-sm font-semibold text-gray-600 hover:bg-gray-50 transition disabled:opacity-50"
            >
              <RefreshCw className={`h-4 w-4 ${isFetching ? "animate-spin" : ""}`} />
              Refresh
            </button>
          </div>
        </div>

        <div className="flex gap-2 flex-wrap">
          {["pending", "approved", "rejected", "processing", "completed", "failed", ""].map((s) => (
            <button
              key={s || "all"}
              type="button"
              onClick={() => { setStatusFilter(s); setSelectedIds(new Set()); }}
              className={`rounded-full px-4 py-1.5 text-sm font-semibold capitalize transition border ${
                statusFilter === s
                  ? "bg-[#111827] text-white border-[#111827]"
                  : "bg-white text-gray-600 border-gray-200 hover:bg-gray-50"
              }`}
            >
              {s || "All"}
            </button>
          ))}
        </div>

        <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
          <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100">
            <div className="flex items-center gap-3">
              {statusFilter === "pending" && pendingPayouts.length > 0 && (
                <input
                  type="checkbox"
                  checked={selectedIds.size === pendingPayouts.length && pendingPayouts.length > 0}
                  onChange={handleSelectAll}
                  className="h-4 w-4 rounded border-gray-300 accent-[#111827] cursor-pointer"
                />
              )}
              <span className="text-sm font-semibold text-[#0b101e]">
                {total} payout{total !== 1 ? "s" : ""}
              </span>
            </div>
          </div>

          {isLoading ? (
            <div className="flex items-center justify-center py-20">
              <div className="h-8 w-8 animate-spin rounded-full border-2 border-gray-200 border-t-[#111827]" />
            </div>
          ) : payouts.length === 0 ? (
            <div className="py-20 text-center">
              <p className="text-lg font-semibold text-[#111827]">No payouts found</p>
              <p className="mt-2 text-sm text-gray-500">
                {statusFilter ? `No ${statusFilter} payout requests at the moment.` : "No payout requests found."}
              </p>
            </div>
          ) : (
            <div className="divide-y divide-gray-100">
              {payouts.map((payout) => (
                <PayoutRow
                  key={payout.id}
                  payout={payout}
                  selected={selectedIds.has(payout.id)}
                  onToggleSelect={() => handleToggleSelect(payout.id)}
                  onApprove={() => setApproveTarget(payout)}
                  onReject={(reason) => rejectMutation.mutate({ id: payout.id, reason })}
                  isApproving={approveMutation.isPending && approveMutation.variables === payout.id}
                  isRejecting={rejectMutation.isPending && rejectMutation.variables?.id === payout.id}
                />
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Approve Confirmation Modal */}
      {approveTarget && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div
            className="absolute inset-0 bg-black/40 backdrop-blur-sm"
            onClick={() => !approveMutation.isPending && setApproveTarget(null)}
          />
          <div className="relative z-10 w-full max-w-sm rounded-2xl bg-white shadow-2xl border border-gray-100 p-6">
            <div className="flex items-center gap-3 mb-4">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-green-100">
                <CheckCircle className="h-5 w-5 text-green-600" />
              </div>
              <h2 className="text-base font-bold text-[#111827]">Confirm Approval</h2>
            </div>

            <p className="text-sm text-gray-600 mb-1">
              You are about to approve the payout request from:
            </p>
            <p className="text-sm font-semibold text-[#111827] mb-1">
              {approveTarget.user_name || "Unknown"}{" "}
              <span className="font-normal text-gray-400">({approveTarget.user_email})</span>
            </p>
            <p className="text-2xl font-bold text-green-600 mt-3 mb-5">
              {new Intl.NumberFormat("en-IN", { style: "currency", currency: "INR", minimumFractionDigits: 2 }).format(approveTarget.amount)}
            </p>

            <p className="text-xs text-gray-400 mb-5">
              This action will mark the payout as approved and process the transfer. This cannot be undone.
            </p>

            <div className="flex gap-3">
              <button
                type="button"
                onClick={() => setApproveTarget(null)}
                disabled={approveMutation.isPending}
                className="flex-1 rounded-xl border border-gray-200 px-4 py-2.5 text-sm font-semibold text-gray-600 hover:bg-gray-50 transition disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                type="button"
                onClick={() => approveMutation.mutate(approveTarget.id)}
                disabled={approveMutation.isPending}
                className="flex-1 rounded-xl bg-green-600 px-4 py-2.5 text-sm font-semibold text-white hover:bg-green-700 transition disabled:opacity-50 flex items-center justify-center gap-2"
              >
                {approveMutation.isPending ? (
                  <>
                    <span className="h-4 w-4 animate-spin rounded-full border-2 border-white/30 border-t-white" />
                    Approving…
                  </>
                ) : (
                  <>
                    <CheckCircle className="h-4 w-4" />
                    Confirm Approve
                  </>
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </AdminLayout>
  );
}
