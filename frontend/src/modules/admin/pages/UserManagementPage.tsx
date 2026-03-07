import { useState, useRef, useEffect } from "react";
import { useUsers, useToggleUserStatus } from "../hooks";
import AdminLayout from "../components/AdminLayout";
import { ChevronDown, Download } from "lucide-react";

export default function UserManagementPage() {
  const [page, setPage] = useState(1);
  const { data } = useUsers({ page, limit: 10 });
  const toggleUser = useToggleUserStatus();

  // State for the "dropdown" action menu
  const [openDropdownId, setOpenDropdownId] = useState<string | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setOpenDropdownId(null);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  const usersList = data?.users || [];

  return (
    <AdminLayout title="User List">
      
      <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">

        <div className="w-full overflow-x-auto">
          <table className="w-full text-sm text-left">
            <thead className="bg-[#f8f9fa] text-xs font-bold text-gray-900 border-b border-gray-100">
              <tr>
                <th className="px-6 py-4 text-center">Name</th>
                <th className="px-6 py-4 text-center">Email</th>
                <th className="px-6 py-4 text-center">Total Bookings</th>
                <th className="px-6 py-4 text-center">Wallet Balance</th>
                <th className="px-6 py-4 text-center">Status</th>
                <th className="px-6 py-4 text-center">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 text-gray-700 font-medium">
              {usersList.map((user: any) => (
                <tr key={user.id} className="hover:bg-gray-50 transition-colors">
                  <td className="px-6 py-4 text-center">{user.name}</td>
                  <td className="px-6 py-4 text-center">{user.email}</td>
                  {/* Since these properties might not exist in the real API yet, we fallback to 0 */}
                  <td className="px-6 py-4 text-center">{user.total_bookings ?? 0}</td>
                  <td className="px-6 py-4 text-center">{user.wallet_balance ?? 0}</td>
                  <td className="px-6 py-4 text-center">
                    <span 
                      className={`inline-block px-4 py-1.5 rounded-full text-xs text-white ${
                        user.is_active ? "bg-[#0ec3c5]" : "bg-[#e53e5d]"
                      }`}
                    >
                      {user.is_active ? "Active" : "Suspended"}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-center relative">
                    <button
                      onClick={() => setOpenDropdownId(openDropdownId === user.id ? null : user.id)}
                      className="inline-flex items-center justify-between w-20 px-3 py-1.5 border border-blue-500 rounded-lg text-[#0b101e] text-xs font-semibold hover:bg-gray-50 transition-colors"
                    >
                      View <ChevronDown className="w-4 h-4 text-blue-500" />
                    </button>
                    
                    {/* Minimal Context Menu for toggling status */}
                    {openDropdownId === user.id && (
                      <div ref={dropdownRef} className="absolute z-10 right-10 top-12 bg-white border border-gray-200 shadow-xl rounded-lg py-1 w-32">
                        <button
                          className="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm font-medium text-gray-700"
                          onClick={() => {
                            toggleUser.mutate({
                              userId: user.id,
                              isActive: !user.is_active,
                            });
                            setOpenDropdownId(null);
                          }}
                        >
                          {user.is_active ? "Suspend User" : "Activate User"}
                        </button>
                      </div>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* Footer actions of table */}
        <div className="flex items-center justify-end px-6 py-4 border-t border-gray-100 mb-2">
          <div className="flex items-center gap-1 text-sm font-medium">
            <button 
              onClick={() => setPage((p) => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-2 py-1 text-gray-500 hover:text-gray-900 disabled:opacity-50"
            >
              &lt; Prev
            </button>
            <button className="w-8 h-8 rounded shrink-0 bg-gray-300 text-gray-800 flex items-center justify-center">1</button>
            <button className="w-8 h-8 rounded shrink-0 hover:bg-gray-100 text-gray-600 flex items-center justify-center">2</button>
            <button className="w-8 h-8 rounded shrink-0 hover:bg-gray-100 text-gray-600 flex items-center justify-center">3</button>
            <button 
              onClick={() => setPage((p) => p + 1)}
              className="px-2 py-1 text-gray-500 hover:text-gray-900"
            >
              Next &gt;
            </button>
          </div>
        </div>
      </div>

      <div className="flex justify-end mt-4">
        <button className="flex items-center gap-2 border border-gray-900 rounded-xl px-4 py-2.5 text-sm font-bold text-gray-900 hover:bg-gray-50 transition-colors">
          Download as CSV
          <Download className="w-4 h-4" />
        </button>
      </div>

    </AdminLayout>
  );
}