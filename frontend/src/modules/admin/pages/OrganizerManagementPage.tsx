import { useState, useRef, useEffect } from "react";
import { createPortal } from "react-dom";
import { useSearchParams } from "react-router-dom";
import { useApproveOrganizer, useOrganizers, useRejectOrganizer, useToggleUserStatus } from "../hooks";
import AdminLayout from "../components/AdminLayout";
import { ChevronDown, Download, Search, Filter } from "lucide-react";
import { useDebounce } from "../../../shared/hooks/useDebounce";
import { exportToCSV } from "../../../shared/utils/csv";

export default function OrganizerManagementPage() {
  const [searchParams, setSearchParams] = useSearchParams();

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

  const { data } = useOrganizers({
    page,
    limit,
    ...(searchFilter && { search: searchFilter }),
    ...(statusFilter !== "all" && { status: statusFilter })
  });
  const toggleUser = useToggleUserStatus();
  const approveOrganizer = useApproveOrganizer();
  const rejectOrganizer = useRejectOrganizer();

  const [openDropdownId, setOpenDropdownId] = useState<string | null>(null);
  const [dropdownPosition, setDropdownPosition] = useState<{ top: number; left: number } | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

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

  const toggleDropdown = (organizerId: string, target: HTMLButtonElement) => {
    if (openDropdownId === organizerId) {
      setOpenDropdownId(null);
      setDropdownPosition(null);
      return;
    }

    const rect = target.getBoundingClientRect();
    setOpenDropdownId(organizerId);
    setDropdownPosition({
      top: rect.bottom + 8,
      left: rect.right - 160,
    });
  };

  const organizersList = data?.organizers || [];
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
    <AdminLayout title="Organizer List">

      
      <div className="flex flex-col sm:flex-row justify-between items-center gap-4 mb-6 mt-6">
        <div className="relative w-full sm:w-80">
          <input
            type="text"
            placeholder="Search organizers..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2.5 bg-white border border-gray-200 rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all"
          />
          <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        </div>

        <div className="flex items-center gap-3 w-full sm:w-auto">
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
              <option value="active">Active</option>
              <option value="suspended">Suspended</option>
              <option value="pending">Pending Approval</option>
              <option value="approved">Approved</option>
              <option value="rejected">Rejected</option>
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
                <th className="px-6 py-4 text-center">Name</th>
                <th className="px-6 py-4 text-center">Email</th>
                <th className="px-6 py-4 text-center">Total Bookings</th>
                <th className="px-6 py-4 text-center">Total Events</th>
                <th className="px-6 py-4 text-center">Wallet Balance</th>
                <th className="px-6 py-4 text-center">Total Revenue Generated</th>
                <th className="px-6 py-4 text-center">Status</th>
                <th className="px-6 py-4 text-center">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 text-gray-700 font-medium">
              {organizersList.map((org: any) => (
                <tr key={org.id} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-4 text-center">{org.name}</td>
                  <td className="px-6 py-4 text-center">{org.email}</td>
                  <td className="px-6 py-4 text-center">{org.total_bookings ?? 0}</td>
                  <td className="px-6 py-4 text-center">{org.total_events ?? 0}</td>
                  <td className="px-6 py-4 text-center">{org.wallet_balance ?? 0}</td>
                  <td className="px-6 py-4 text-center">{org.total_revenue_generated ?? 0}</td>
                  <td className="px-6 py-4 text-center">
                    {org.approval_status === "pending" ? (
                      <span className="inline-block px-4 py-1.5 rounded-full text-xs text-white bg-amber-500">
                        Pending
                      </span>
                    ) : org.approval_status === "rejected" ? (
                      <span className="inline-block px-4 py-1.5 rounded-full text-xs text-white bg-gray-500">
                        Rejected
                      </span>
                    ) : (
                    <span
                      className={`inline-block px-4 py-1.5 rounded-full text-xs text-white ${
                        org.is_active ? "bg-[#0ec3c5]" : "bg-[#e53e5d]"
                      }`}
                    >
                      {org.is_active ? "Active" : "Suspended"}
                    </span>
                    )}
                  </td>
                  <td className="px-6 py-4 text-center">
                    <button
                      onClick={(e) => toggleDropdown(org.id, e.currentTarget)}
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

        {organizersList.length === 0 && (
          <div className="py-12 text-center text-gray-500">
            No organizers found
          </div>
        )}

        <div className="flex items-center justify-between px-6 py-4 border-t border-gray-100 mb-2">
          <div className="flex items-center gap-2">
            <button className="flex items-center gap-1 text-sm text-gray-500 hover:text-gray-900 transition-colors font-medium">
              <ChevronDown className="w-4 h-4 rotate-90" />
              Prev
            </button>
            <div className="flex gap-1 ml-4">{renderPageNumbers()}</div>
            <button className="flex items-center gap-1 ml-4 text-sm text-gray-500 hover:text-gray-900 transition-colors font-medium">
              Next
              <ChevronDown className="w-4 h-4 -rotate-90" />
            </button>
          </div>
        </div>
      </div>

      <div className="flex justify-end mt-4">
        <button
          onClick={() => exportToCSV(organizersList, "organizers_list", [
            { header: "Name", key: "name" },
            { header: "Email", key: "email" },
            { header: "Total Bookings", key: "total_bookings" },
            { header: "Total Events", key: "total_events" },
            { header: "Wallet Balance", key: "wallet_balance" },
            { header: "Total Revenue Generated", key: "total_revenue_generated" },
            { header: "Status", key: "approval_status" },
          ])}
          className="flex items-center gap-2 border border-gray-900 rounded-xl px-4 py-2.5 text-sm font-bold text-gray-900 hover:bg-gray-50 transition-colors"
        >
          Download as CSV
          <Download className="w-4 h-4" />
        </button>
      </div>

      {openDropdownId && dropdownPosition &&
        createPortal(
          <div
            ref={dropdownRef}
            className="fixed z-[1000] bg-white border border-gray-200 shadow-xl rounded-lg py-1 w-40"
            style={{ top: dropdownPosition.top, left: dropdownPosition.left }}
          >
            {organizersList
              .filter((org: any) => org.id === openDropdownId)
              .map((org: any) => (
                <div key={org.id}>
                  <button
                    className="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm font-medium text-gray-700"
                    onClick={() => {
                      if (org.approval_status === "pending" || org.approval_status === "rejected") {
                        approveOrganizer.mutate(org.id);
                      } else {
                        toggleUser.mutate({
                          userId: org.id,
                          isActive: !org.is_active,
                        });
                      }
                      setOpenDropdownId(null);
                      setDropdownPosition(null);
                    }}
                  >
                    {org.approval_status === "pending" || org.approval_status === "rejected"
                      ? "Approve Organizer"
                      : org.is_active
                        ? "Suspend User"
                        : "Activate User"}
                  </button>
                  {org.approval_status === "pending" && (
                    <button
                      className="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm font-medium text-red-600"
                      onClick={() => {
                        rejectOrganizer.mutate(org.id);
                        setOpenDropdownId(null);
                        setDropdownPosition(null);
                      }}
                    >
                      Reject Request
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
