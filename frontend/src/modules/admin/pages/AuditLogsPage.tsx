import { useState } from "react";
import AdminLayout from "../components/AdminLayout";
import { useAuditLogs } from "../hooks";
import { Loader2, Search, Download, X } from "lucide-react";
import type { AuditLog } from "../api";

export default function AuditLogsPage() {
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [search, setSearch] = useState("");
  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null);

  const { data, isLoading } = useAuditLogs(page, limit);

  
  const filteredLogs = data?.logs?.filter(log => 
      log.action.toLowerCase().includes(search.toLowerCase()) || 
      log.admin_name.toLowerCase().includes(search.toLowerCase())
  );

  const getTagStyle = (tag: string) => {
    switch(tag) {
      case "EVENT": return "bg-rose-500 text-white";
      case "USER": return "bg-[#0b132b] text-white";
      case "ORGANIZER": return "bg-indigo-400 text-white";
      case "PAYOUT": return "bg-emerald-500 text-white";
      case "REFUND": return "bg-amber-500 text-white";
      case "SETTINGS": return "bg-sky-500 text-white";
      default: return "bg-gray-500 text-white";
    }
  };

  const handleDownloadCSV = () => {
    if (!data?.logs) return;
    const headers = ["Action", "Admin Name", "Timestamp", "IP Address", "Action Tag", "Details"];
    const rows = data.logs.map(log => [
      `"${log.action}"`,
      `"${log.admin_name}"`,
      `"${new Date(log.timestamp).toLocaleString()}"`,
      `"${log.ip_address}"`,
      `"${log.action_tag}"`,
      `"${log.details.replace(/"/g, '""')}"`
    ]);
    
    let csvContent = "data:text/csv;charset=utf-8," 
        + headers.join(",") + "\n" 
        + rows.map(e => e.join(",")).join("\n");
        
    var encodedUri = encodeURI(csvContent);
    var link = document.createElement("a");
    link.setAttribute("href", encodedUri);
    link.setAttribute("download", "admin_audit_logs.csv");
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  const parseDetails = (detailsStr: string) => {
    try {
      return JSON.parse(detailsStr);
    } catch (e) {
      return ;
    }
  };

  return (
    <AdminLayout title="Track Admin Changes">
      
      
      <div className="flex justify-between items-center mb-6">
        <div className="relative w-full max-w-sm">
           <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
           <input 
             type="text" 
             placeholder="Quick Search" 
             className="w-full pl-9 pr-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
             value={search}
             onChange={(e) => setSearch(e.target.value)}
           />
        </div>
      </div>

      
      <div className="bg-white rounded-[20px] shadow-sm border border-gray-100 overflow-hidden flex flex-col min-h-[500px]">
        {isLoading ? (
          <div className="flex-1 flex justify-center items-center">
            <Loader2 className="w-8 h-8 animate-spin text-gray-300" />
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="w-full text-sm text-center">
                <thead className="bg-[#f8fafc] text-gray-700 text-xs font-bold border-b border-gray-100">
                  <tr>
                    <th className="px-6 py-4 font-bold">Action</th>
                    <th className="px-6 py-4 font-bold">Admin Name</th>
                    <th className="px-6 py-4 font-bold">Timestamp</th>
                    <th className="px-6 py-4 font-bold">IP Address</th>
                    <th className="px-6 py-4 font-bold">Action Tag</th>
                    <th className="px-6 py-4 font-bold">Details</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100">
                  {filteredLogs?.map((log) => (
                    <tr key={log.id} className="hover:bg-gray-50/50 transition">
                      <td className="px-6 py-4 text-gray-800">{log.action}</td>
                      <td className="px-6 py-4 text-gray-700">{log.admin_name}</td>
                      <td className="px-6 py-4 text-gray-600">
                        {new Date(log.timestamp).toLocaleString("en-US", {
                          hour: "numeric", minute: "numeric", hour12: true,
                          month: "2-digit", day: "2-digit", year: "numeric",
                        }).replace(",", "")}
                      </td>
                      <td className="px-6 py-4 text-gray-500 font-mono text-xs">{log.ip_address}</td>
                      <td className="px-6 py-4">
                        <span className={`px-4 py-1.5 rounded-full text-xs font-bold uppercase tracking-wide inline-block min-w-[100px] ${getTagStyle(log.action_tag)}`}>
                          {log.action_tag}
                        </span>
                      </td>
                      <td className="px-6 py-4">
                        <button 
                          onClick={() => setSelectedLog(log)}
                          className="border border-[#cbd5e1] text-[#475569] hover:bg-gray-100 px-4 py-1.5 rounded-xl font-semibold text-xs transition"
                        >
                          View
                        </button>
                      </td>
                    </tr>
                  ))}
                  
                  {filteredLogs?.length === 0 && (
                    <tr>
                      <td colSpan={6} className="py-12 text-gray-400">No logs found</td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>

            
            <div className="mt-auto border-t border-gray-100 px-6 py-4 flex items-center justify-between">
               
               <button 
                 onClick={handleDownloadCSV}
                 className="flex items-center gap-2 border border-black text-black px-4 py-2 rounded-lg font-bold text-sm hover:bg-gray-50 transition"
               >
                 Download as CSV <Download className="w-4 h-4 shadow" strokeWidth={3} />
               </button>

               <div className="flex items-center gap-2">
                 <button 
                   disabled={page === 1}
                   onClick={() => setPage(page - 1)}
                   className="text-gray-500 hover:text-black font-semibold text-sm px-2 disabled:opacity-50"
                 >
                   &lt; Prev
                 </button>
                 <div className="flex items-center gap-1">
                    {[page, page+1, page+2].map(p => {
                       if (p > Math.ceil((data?.pagination?.total || 0)/limit) && p !== 1) return null;
                       return (
                         <button 
                           key={p}
                           onClick={() => setPage(p)}
                           className={`w-8 h-8 rounded-lg text-sm font-semibold flex items-center justify-center ${p === page ? 'bg-gray-200 text-black' : 'text-gray-600 hover:bg-gray-100'}`}
                         >
                           {p}
                         </button>
                       )
                    })}
                 </div>
                 <button 
                   disabled={page >= Math.ceil((data?.pagination?.total || 0) / limit)}
                   onClick={() => setPage(page + 1)}
                   className="text-gray-500 hover:text-black font-semibold text-sm px-2 disabled:opacity-50"
                 >
                   Next &gt;
                 </button>
               </div>

            </div>
          </>
        )}
      </div>

      
      {selectedLog && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40 backdrop-blur-sm animate-in fade-in">
          <div className="bg-white rounded-2xl w-full max-w-md shadow-2xl overflow-hidden">
            <div className="px-6 py-4 border-b border-gray-100 flex justify-between items-center bg-gray-50/50">
              <h3 className="font-bold text-gray-900">Audit Log Details</h3>
              <button 
                onClick={() => setSelectedLog(null)}
                className="text-gray-400 hover:text-gray-900 p-1 rounded-md hover:bg-gray-200 transition"
              >
                <X className="w-5 h-5" />
              </button>
            </div>
            
            <div className="p-6">
               <div className="mb-4">
                  <span className={`px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-wide inline-block mb-3 ${getTagStyle(selectedLog.action_tag)}`}>
                    {selectedLog.action_tag}
                  </span>
                  <p className="font-medium text-gray-900">{selectedLog.action}</p>
               </div>

               <div className="bg-gray-50 rounded-xl p-4 border border-gray-100">
                 <h4 className="text-xs font-bold tracking-widest text-gray-400 uppercase mb-3">Entity Data</h4>
                 {Object.entries(parseDetails(selectedLog.details)).length > 0 ? (
                   <ul className="space-y-3">
                     {Object.entries(parseDetails(selectedLog.details)).map(([k, v]: [string, any]) => (
                       <li key={k} className="flex flex-col">
                         <span className="text-xs text-gray-500 font-medium capitalize">{k.replace(/_/g, ' ')}</span>
                         <span className="text-sm font-semibold text-gray-800 break-all">{v?.toString() || "N/A"}</span>
                       </li>
                     ))}
                   </ul>
                 ) : (
                   <span className="text-sm text-gray-500 italic">No additional details recorded.</span>
                 )}
               </div>

               <div className="mt-6 flex justify-end">
                 <button 
                   onClick={() => setSelectedLog(null)}
                   className="bg-gray-900 text-white font-semibold text-sm px-5 py-2.5 rounded-xl shadow-lg hover:bg-black transition"
                 >
                   Close
                 </button>
               </div>
            </div>
          </div>
        </div>
      )}
    </AdminLayout>
  );
}
