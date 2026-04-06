import { useState } from "react";
import AdminLayout from "../components/AdminLayout";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { adminApi, type AdminRefundDetail } from "../api";
import { CheckCircle, Clock } from "lucide-react";

export default function RefundsPage() {
  const [filterStatus, setFilterStatus] = useState<string>("pending");
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ["admin_refunds", filterStatus],
    queryFn: () => adminApi.getRefunds(filterStatus !== "all" ? { status: filterStatus } : undefined),
  });

  const processMutation = useMutation({
    mutationFn: (id: string) => adminApi.processRefund(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin_refunds"] });
      queryClient.invalidateQueries({ queryKey: ["admin_platformWallet"] });
    },
  });

  const handleProcess = (id: string) => {
    if (confirm("Are you sure you want to mark this refund as processed? Ensure the amount is transferred.")) {
      processMutation.mutate(id);
    }
  };

  return (
    <AdminLayout title="Refund Management">
      <div className="flex flex-col gap-6">
        <div className="flex items-center justify-between">
          <div className="flex gap-2">
            {["pending", "processed", "all"].map((status) => (
              <button
                key={status}
                onClick={() => setFilterStatus(status)}
                className={`px-4 py-2 rounded-xl text-sm font-medium transition-colors capitalize ${
                  filterStatus === status
                    ? "bg-[#0b101e] text-white"
                    : "bg-white text-gray-600 border border-gray-200 hover:bg-gray-50"
                }`}
              >
                {status}
              </button>
            ))}
          </div>
        </div>

        <div className="bg-white rounded-3xl border border-gray-200 shadow-sm overflow-hidden">
          {isLoading ? (
            <div className="p-12 flex justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
            </div>
          ) : !data || data.refunds.length === 0 ? (
            <div className="p-12 text-center text-gray-500">
              No {filterStatus !== "all" ? filterStatus : ""} refund requests found.
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-left text-sm text-gray-600">
                <thead className="bg-gray-50 text-gray-900 font-semibold border-b border-gray-200 text-xs uppercase tracking-wider">
                  <tr>
                    <th className="px-6 py-4">Request Details</th>
                    <th className="px-6 py-4">User</th>
                    <th className="px-6 py-4">Amount</th>
                    <th className="px-6 py-4">Status</th>
                    <th className="px-6 py-4 text-right">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {data.refunds.map((refund: AdminRefundDetail) => (
                    <tr key={refund.id} className="hover:bg-gray-50 transition-colors">
                      <td className="px-6 py-4">
                        <p className="font-medium text-gray-900">
                          {new Date(refund.requested_at).toLocaleDateString()}
                        </p>
                        <p className="text-xs text-gray-500 mt-0.5">ID: {refund.id.slice(0, 8)}</p>
                        <p className="text-xs text-gray-400">Booking: {refund.booking_id.slice(0, 8)}</p>
                      </td>
                      <td className="px-6 py-4">
                        <p className="font-medium text-gray-900">{refund.user_name}</p>
                        <p className="text-xs text-gray-500 mt-0.5">{refund.user_email}</p>
                      </td>
                      <td className="px-6 py-4">
                        <span className="font-semibold text-gray-900">
                          ₹{refund.amount.toFixed(2)}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <span
                          className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold ${
                            refund.status === "processed"
                              ? "bg-green-100 text-green-700"
                              : "bg-orange-100 text-orange-700"
                          }`}
                        >
                          {refund.status === "processed" ? (
                            <CheckCircle className="w-3.5 h-3.5" />
                          ) : (
                            <Clock className="w-3.5 h-3.5" />
                          )}
                          <span className="capitalize">{refund.status}</span>
                        </span>
                      </td>
                      <td className="px-6 py-4 text-right">
                        {refund.status === "pending" && (
                          <div className="flex items-center justify-end gap-2">
                            <button
                              onClick={() => handleProcess(refund.id)}
                              disabled={processMutation.isPending}
                              className="px-4 py-2 bg-green-500 text-white font-medium text-xs rounded-lg hover:bg-green-600 transition-colors disabled:opacity-50"
                            >
                              Mark Processed
                            </button>
                          </div>
                        )}
                        {refund.status === "processed" && refund.processed_at && (
                          <p className="text-xs text-gray-400">
                            Processed: {new Date(refund.processed_at).toLocaleDateString()}
                          </p>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </AdminLayout>
  );
}
