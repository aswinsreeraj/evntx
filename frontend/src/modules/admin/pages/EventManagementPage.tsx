import { useState, useRef, useEffect } from "react";
import { createPortal } from "react-dom";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useEvents, useApproveEvent, useRejectEvent } from "../hooks";
import AdminLayout from "../components/AdminLayout";
import { ChevronDown, Download, Search, Filter, X } from "lucide-react";
import { useDebounce } from "../../../shared/hooks/useDebounce";
import { exportToCSV } from "../../../shared/utils/csv";

function getStatusColor(status: string) {
  switch (status.toLowerCase()) {
    case "approved":
      return "bg-[#0ec3c5] text-white border-transparent";
    case "rejected":
      return "bg-[#e53e5d] text-white border-transparent";
    case "pending":
      return "bg-[#ffb020] text-white border-transparent";
    case "draft":
      return "bg-[#1c2438] text-white border-transparent";
    case "live":
      return "bg-white text-[#e53e5d] border-[#e53e5d] border";
    case "completed":
      return "bg-gray-200 text-gray-700 border-transparent";
    default:
      return "bg-gray-100 text-gray-700 border-transparent";
  }
}

function canApprove(status: string) {
  return status.toLowerCase() === "pending";
}

function canReject(status: string) {
  return status.toLowerCase() === "pending";
}

export default function EventManagementPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();

  const page = parseInt(searchParams.get("page") || "1", 10);
  const statusFilter = searchParams.get("status") || "all";
  const limit = parseInt(searchParams.get("limit") || "10", 10);
  const searchFilter = searchParams.get("search") || "";

  const [searchTerm, setSearchTerm] = useState(searchFilter);
  const debouncedSearchTerm = useDebounce(searchTerm, 300);

  const updateParams = (updates: Record<string, string>) => {
    const newParams = new URLSearchParams(searchParams);
    Object.entries(updates).forEach(([key, value]) => {
      if (!value || (key === "page" && value === "1") || (key === "limit" && value === "10") || (key === "status" && value === "all")) {
        newParams.delete(key);
      } else {
        newParams.set(key, value);
      }
    });
    setSearchParams(newParams);
  };

  useEffect(() => {
    if (debouncedSearchTerm !== searchFilter) {
      updateParams({ search: debouncedSearchTerm, page: "1" });
    }
  }, [debouncedSearchTerm]);

  const setPage = (p: number | ((prev: number) => number)) => {
    const newPage = typeof p === "function" ? p(page) : p;
    updateParams({ page: newPage.toString() });
  };

  const setStatusFilter = (s: string) => updateParams({ status: s, page: "1" });
  const setLimit = (l: number) => updateParams({ limit: l.toString(), page: "1" });

  const { data } = useEvents({
    page,
    limit,
    ...(searchFilter && { search: searchFilter }),
    ...(statusFilter !== "all" && { status: statusFilter })
  });

  const approveEvent = useApproveEvent();
  const rejectEvent = useRejectEvent();

  const [openDropdownId, setOpenDropdownId] = useState<string | null>(null);
  const [dropdownPosition, setDropdownPosition] = useState<{ top: number; left: number } | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const [modalState, setModalState] = useState<{ isOpen: boolean; type: "Approve" | "Reject" | null; eventId: string | null }>({
    isOpen: false,
    type: null,
    eventId: null,
  });
  const [rejectReason, setRejectReason] = useState("");

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setOpenDropdownId(null);
        setDropdownPosition(null);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const toggleDropdown = (eventId: string, target: HTMLButtonElement) => {
    if (openDropdownId === eventId) {
      setOpenDropdownId(null);
      setDropdownPosition(null);
      return;
    }

    const rect = target.getBoundingClientRect();
    setOpenDropdownId(eventId);
    setDropdownPosition({
      top: rect.bottom + 8,
      left: rect.right - 160,
    });
  };

  const handleAction = (eventId: string, type: "Approve" | "Reject") => {
    setModalState({ isOpen: true, type, eventId });
    setOpenDropdownId(null);
    setDropdownPosition(null);
  };

  const confirmAction = () => {
    if (modalState.type === "Approve" && modalState.eventId) {
      approveEvent.mutate(modalState.eventId, {
        onSuccess: () => setModalState({ isOpen: false, type: null, eventId: null })
      });
    } else if (modalState.type === "Reject" && modalState.eventId) {
      rejectEvent.mutate({ eventId: modalState.eventId, reason: rejectReason || "Rejected by admin" }, {
        onSuccess: () => {
          setModalState({ isOpen: false, type: null, eventId: null });
          setRejectReason("");
        }
      });
    }
  };

  const eventsList = data?.events || [];
  const pagination = data?.pagination;
  const totalPages = pagination ? Math.ceil(pagination.total / pagination.limit) : 1;

  const renderPageNumbers = () => {
    const pages = [];
    for (let i = 1; i <= totalPages; i++) {
      pages.push(
        <button
          key={i}
          onClick={() => setPage(i)}
          className={`w-8 h-8 rounded shrink-0 flex items-center justify-center text-sm transition-colors ${
            page === i
              ? "bg-gray-900 text-white"
              : "hover:bg-gray-100 text-gray-600"
          }`}
        >
          {i}
        </button>
      );
    }
    return pages;
  };

  return (
    <AdminLayout title="Manage Events">

      {}
      <div className="flex flex-col sm:flex-row justify-between items-center gap-4 mb-6 mt-6">
        <div className="relative w-full sm:w-80">
          <input
            type="text"
            placeholder="Search events..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          />
          <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        </div>

        <div className="flex items-center gap-3 w-auto">
          <div className="relative w-full sm:w-fit">
            <select
              value={limit}
              onChange={(e) => setLimit(Number(e.target.value))}
              className="w-full pl-4 pr-10 py-2.5 bg-white border border-gray-200 rounded-xl text-sm appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer transition-all"
            >
              <option value={5}>5 per page</option>
              <option value={10}>10 per page</option>
              <option value={20}>20 per page</option>
              <option value={50}>50 per page</option>
            </select>
            <ChevronDown className="absolute right-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
          </div>

          <div className="relative w-full sm:w-42">
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-200 rounded-xl text-sm appearance-none focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer transition-all"
            >
              <option value="all">All Status</option>
              <option value="pending">Pending</option>
              <option value="approved">Approved</option>
              <option value="rejected">Rejected</option>
              <option value="live">Live</option>
              <option value="draft">Draft</option>
              <option value="completed">Completed</option>
            </select>
            <Filter className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
            <ChevronDown className="absolute right-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
          </div>
        </div>
      </div>

      <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-visible">
        <div className="w-full overflow-x-auto overflow-y-visible">
          <table className="w-full text-sm text-left">
            <thead className="bg-[#f8f9fa] text-xs font-bold text-gray-900 border-b border-gray-100">
              <tr>
                <th className="px-6 py-4 text-center">Event Name</th>
                <th className="px-6 py-4 text-center">Organizer</th>
                <th className="px-6 py-4 text-center">Date</th>
                <th className="px-6 py-4 text-center">Tickets Sold</th>
                <th className="px-6 py-4 text-center">Revenue</th>
                <th className="px-6 py-4 text-center">City</th>
                <th className="px-6 py-4 text-center">Status</th>
                <th className="px-6 py-4 text-center">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 text-gray-700 font-medium">
              {eventsList.map((evt: any) => (
                <tr key={evt.id} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-4 text-center">{evt.title}</td>
                  <td className="px-6 py-4 text-center">{evt.organizer_name || "Unknown"}</td>
                  <td className="px-6 py-4 text-center">
                    {evt.start_time
                      ? new Date(evt.start_time).toLocaleDateString("en-GB", { day: "numeric", month: "short", year: "numeric" })
                      : "N/A"}
                  </td>
                  <td className="px-6 py-4 text-center">{evt.tickets_sold ?? 0}</td>
                  <td className="px-6 py-4 text-center">{evt.revenue ?? 0}</td>
                  <td className="px-6 py-4 text-center">{evt.city}</td>
                  <td className="px-6 py-4 text-center">
                    <span
                      className={`inline-block px-4 py-1.5 rounded-full text-xs font-bold ${getStatusColor(evt.status)}`}
                    >
                      {evt.status.charAt(0).toUpperCase() + evt.status.slice(1)}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-center">
                    <button
                      onClick={(e) => toggleDropdown(evt.id, e.currentTarget)}
                      className="inline-flex items-center justify-between w-20 px-3 py-1.5 border border-blue-500 rounded-lg text-[#0b101e] text-xs font-semibold hover:bg-gray-50 transition-colors"
                    >
                      View <ChevronDown className="w-4 h-4 text-blue-500" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {eventsList.length === 0 && (
          <div className="py-12 text-center text-gray-500">
            No events found
          </div>
        )}

        <div className="flex items-center justify-between px-6 py-4 border-t border-gray-100 mb-2">
          <div className="flex items-center gap-2">
            <button
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-2 py-1 text-gray-500 hover:text-gray-900 disabled:opacity-50 transition-colors text-sm font-medium flex items-center"
            >
              <ChevronDown className="w-4 h-4 rotate-90 mr-1" />
              Prev
            </button>
            <div className="flex gap-1 mx-2">{renderPageNumbers()}</div>
            <button
              onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
              disabled={page === totalPages || totalPages === 0}
              className="px-2 py-1 text-gray-500 hover:text-gray-900 disabled:opacity-50 transition-colors text-sm font-medium flex items-center"
            >
              Next
              <ChevronDown className="w-4 h-4 -rotate-90 ml-1" />
            </button>
          </div>
        </div>
      </div>

      <div className="flex justify-end mt-4">
        <button
          onClick={() => exportToCSV(eventsList, "events_list")}
          className="flex items-center gap-2 border border-gray-900 rounded-xl px-4 py-2.5 text-sm font-bold text-gray-900 hover:bg-gray-50 transition-colors"
        >
          Download as CSV
          <Download className="w-4 h-4" />
        </button>
      </div>

      {modalState.isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50">
          <div className="bg-white rounded-2xl shadow-xl w-full max-w-md overflow-hidden animate-in fade-in zoom-in duration-200">
            <div className="flex justify-between items-center p-6 border-b border-gray-100">
              <h2 className="text-xl font-bold">{modalState.type} Event</h2>
              <button
                onClick={() => {
                  setModalState({ isOpen: false, type: null, eventId: null });
                  setRejectReason("");
                }}
                className="text-gray-400 hover:text-gray-600 transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="p-6">
              <p className="text-gray-600 mb-6">
                Are you sure you want to {modalState.type?.toLowerCase()} this event?
                {modalState.type === "Approve" ? " It will be allowed to go live." : " It will be hidden from users."}
              </p>

              {modalState.type === "Reject" && (
                <div className="mb-6">
                  <label className="block text-sm font-semibold text-[#0b101e] mb-2">
                    Reason for Rejection <span className="text-[#e53e5d]">*</span>
                  </label>
                  <textarea
                    value={rejectReason}
                    onChange={(e) => setRejectReason(e.target.value)}
                    placeholder="Provide a reason for the organizer..."
                    className="w-full px-4 py-3 bg-[#f8f9fa] border border-transparent rounded-xl text-sm focus:outline-none focus:bg-white focus:border-blue-500 focus:ring-4 focus:ring-blue-50 transition-all min-h-[100px] resize-none"
                    required
                  />
                </div>
              )}

              <div className="flex gap-4">
                <button
                  onClick={() => {
                    setModalState({ isOpen: false, type: null, eventId: null });
                    setRejectReason("");
                  }}
                  className="flex-1 px-4 py-3 border border-gray-200 text-gray-700 font-bold rounded-xl hover:bg-gray-50 transition-colors"
                  disabled={approveEvent.isPending || rejectEvent.isPending}
                >
                  Cancel
                </button>
                <button
                  onClick={confirmAction}
                  disabled={approveEvent.isPending || rejectEvent.isPending || (modalState.type === "Reject" && !rejectReason.trim())}
                  className={`flex-1 px-4 py-3 text-white font-bold rounded-xl transition-colors ${
                    modalState.type === "Approve"
                      ? "bg-[#0ec3c5] hover:bg-[#0da6a8]"
                      : "bg-[#e53e5d] hover:bg-[#c2344f]"
                  } disabled:opacity-50`}
                >
                  {approveEvent.isPending || rejectEvent.isPending ? "Processing..." : `Confirm ${modalState.type}`}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {openDropdownId && dropdownPosition &&
        createPortal(
          <div
            ref={dropdownRef}
            className="fixed z-[1000] bg-white border border-gray-200 shadow-xl rounded-lg py-1 w-40"
            style={{ top: dropdownPosition.top, left: dropdownPosition.left }}
          >
                  {eventsList
              .filter((evt: any) => evt.id === openDropdownId)
              .map((evt: any) => (
                <div key={evt.id}>
                  <button
                    className="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm font-medium text-gray-700"
                    disabled={!evt.slug && !evt.id}
                    onClick={() => {
                      navigate(`/events/${evt.slug || evt.id}`);
                      setOpenDropdownId(null);
                      setDropdownPosition(null);
                    }}
                  >
                    View Event
                  </button>
                  {canApprove(evt.status) && (
                    <button
                      className="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm font-medium text-gray-700"
                      disabled={!evt.id}
                      onClick={() => handleAction(evt.id, "Approve")}
                    >
                      Approve
                    </button>
                  )}
                  {canReject(evt.status) && (
                    <button
                      className="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm font-medium text-red-600"
                      disabled={!evt.id}
                      onClick={() => handleAction(evt.id, "Reject")}
                    >
                      Reject
                    </button>
                  )}
                </div>
              ))}
          </div>,
          document.body
        )}

    </AdminLayout>
  );
}
