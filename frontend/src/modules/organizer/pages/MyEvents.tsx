import { useState, useEffect, useCallback } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import OrganizerLayout from "../components/OrganizerLayout";
import { organizerApi } from "../api";
import { X, ChevronDown, MapPin, Loader2, Ticket, ChevronUp } from "lucide-react";

interface TicketType {
  id: string;
  name: string;
  price: number;
  total_quantity: number;
  available_quantity: number;
}

function EventTicketDetails({ slug }: { slug: string }) {
  const [loading, setLoading] = useState(true);
  const [tickets, setTickets] = useState<TicketType[]>([]);
  const [totalCapacity, setTotalCapacity] = useState(0);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const data = await organizerApi.getEventBySlug(slug);
        if (!cancelled) {
          const ticketTypes: TicketType[] = data?.ticket_types ?? [];
          setTickets(ticketTypes);
          setTotalCapacity(
            data?.details?.total_capacity ?? ticketTypes.reduce((s: number, t: TicketType) => s + t.total_quantity, 0)
          );
        }
      } catch {
        if (!cancelled) setError("Failed to load ticket details.");
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();
    return () => { cancelled = true; };
  }, [slug]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-6">
        <Loader2 className="w-5 h-5 animate-spin text-gray-400" />
      </div>
    );
  }

  if (error) {
    return <p className="text-sm text-red-500 py-4 px-1">{error}</p>;
  }

  if (tickets.length === 0) {
    return <p className="text-sm text-gray-400 py-4 px-1">No ticket types defined.</p>;
  }

  const totalAvailable = tickets.reduce((s, t) => s + t.available_quantity, 0);
  const totalSold = totalCapacity - totalAvailable;

  return (
    <div className="mt-1 space-y-3">
      <div className="flex flex-wrap gap-4 text-xs font-semibold text-gray-500 uppercase tracking-wide px-1 pb-1 border-b border-gray-100">
        <span>Total Capacity: <span className="text-gray-900">{totalCapacity}</span></span>
        <span>Sold: <span className="text-[#e53e5d]">{totalSold}</span></span>
        <span>Available: <span className="text-emerald-600">{totalAvailable}</span></span>
      </div>
      <div className="overflow-hidden rounded-xl border border-gray-100">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-xs text-gray-500 uppercase tracking-wide">
            <tr>
              <th className="text-left px-4 py-2.5 font-semibold">Ticket Type</th>
              <th className="text-right px-4 py-2.5 font-semibold">Price</th>
              <th className="text-right px-4 py-2.5 font-semibold">Total</th>
              <th className="text-right px-4 py-2.5 font-semibold">Available</th>
              <th className="text-right px-4 py-2.5 font-semibold">Sold</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50 bg-white">
            {tickets.map((t) => {
              const sold = t.total_quantity - t.available_quantity;
              const pct = t.total_quantity > 0 ? Math.round((sold / t.total_quantity) * 100) : 0;
              const isSoldOut = t.available_quantity === 0;
              const isLow = !isSoldOut && t.available_quantity <= Math.ceil(t.total_quantity * 0.2);
              return (
                <tr key={t.id} className="hover:bg-gray-50/50 transition">
                  <td className="px-4 py-3 font-medium text-gray-800">
                    <div className="flex items-center gap-2">
                      {t.name}
                      {isSoldOut && (
                        <span className="text-[10px] font-bold bg-red-100 text-red-600 px-1.5 py-0.5 rounded-full uppercase tracking-wide">Sold Out</span>
                      )}
                      {isLow && !isSoldOut && (
                        <span className="text-[10px] font-bold bg-amber-100 text-amber-700 px-1.5 py-0.5 rounded-full uppercase tracking-wide">Low</span>
                      )}
                    </div>
                    <div className="mt-1.5 h-1 w-full max-w-[120px] rounded-full bg-gray-100 overflow-hidden">
                      <div
                        className={`h-1 rounded-full transition-all ${pct >= 90 ? "bg-red-500" : pct >= 60 ? "bg-amber-400" : "bg-emerald-500"}`}
                        style={{ width: `${pct}%` }}
                      />
                    </div>
                  </td>
                  <td className="px-4 py-3 text-right text-gray-700">
                    {t.price === 0 ? <span className="text-emerald-600 font-semibold">Free</span> : `\u20b9${t.price.toLocaleString("en-IN")}`}
                  </td>
                  <td className="px-4 py-3 text-right text-gray-600">{t.total_quantity}</td>
                  <td className={`px-4 py-3 text-right font-semibold ${isSoldOut ? "text-red-500" : isLow ? "text-amber-600" : "text-emerald-600"}`}>
                    {t.available_quantity}
                  </td>
                  <td className="px-4 py-3 text-right text-gray-500">{sold}</td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default function MyEvents() {
   const navigate = useNavigate();
   const location = useLocation();
   const [events, setEvents] = useState<any[]>([]);
   const [loading, setLoading] = useState(true);
   const [statusFilter, setStatusFilter] = useState("All");
   const [toast, setToast] = useState<string | null>(null);
   const [rejectionModalOpen, setRejectionModalOpen] = useState(false);
   const [selectedReason, setSelectedReason] = useState("");
   const [expandedTickets, setExpandedTickets] = useState<Set<string>>(new Set());
   const [cancellationModalOpen, setCancellationModalOpen] = useState(false);
   const [selectedEventId, setSelectedEventId] = useState<string | null>(null);
   const [cancellationReason, setCancellationReason] = useState("");
   const [allowEventCancellation, setAllowEventCancellation] = useState<boolean>(true);
   const statuses = ["All", "Draft", "Pending", "Approved", "Rejected", "Live", "Cancellation Pending", "Completed"];

   useEffect(() => {
      organizerApi.getPlatformSettings()
         .then((s) => setAllowEventCancellation(s.allow_event_cancellation))
         .catch(() => setAllowEventCancellation(true));
   }, []);

   useEffect(() => {
      if (location.state?.toastMessage) {
         setToast(location.state.toastMessage);
         window.history.replaceState({}, document.title);
         const timer = setTimeout(() => setToast(null), 4000);
         return () => clearTimeout(timer);
      }
   }, [location]);

   const loadEvents = useCallback(async () => {
      setLoading(true);
      try {
         const normalizedStatus = statusFilter === "Cancellation Pending" ? "cancellation_pending" : statusFilter;
         const data = await organizerApi.getOrganizerEvents(normalizedStatus === "All" ? "" : normalizedStatus);
         setEvents(data.events || []);
      } catch (err) {
         console.error(err);
      } finally {
         setLoading(false);
      }
   }, [statusFilter]);

   useEffect(() => {
      loadEvents();
   }, [loadEvents]);

   const handleDelete = async (id: string) => {
      if (!confirm("Are you sure you want to delete this event? This action cannot be undone.")) return;
      try {
         await organizerApi.deleteEvent(id);
         loadEvents();
      } catch (err) {
         console.error(err);
         alert("Failed to delete event.");
      }
   };

   const handleSubmit = async (id: string) => {
      if (!confirm("Submit this event for admin approval?")) return;
      try {
         await organizerApi.submitEvent(id);
         loadEvents();
      } catch (err) {
         console.error(err);
         alert("Failed to submit event for approval.");
      }
   };

   const toggleTickets = (eventId: string) => {
      setExpandedTickets((prev) => {
         const next = new Set(prev);
         if (next.has(eventId)) next.delete(eventId);
         else next.add(eventId);
         return next;
      });
   };

   const getStatusBadge = (status: string) => {
      const s = (status || "").toLowerCase();
      switch (s) {
         case "approved":
            return <div className="bg-[#1ccbce] text-white px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase flex items-center gap-2"><span className="w-1.5 h-1.5 rounded-full bg-white" /> APPROVED</div>;
         case "rejected":
            return <div className="bg-[#ff0000] text-white px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase flex items-center gap-1.5"><X className="w-3 h-3" /> REJECTED</div>;
         case "pending":
            return <div className="bg-[#facc15] text-[#854d0e] px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase flex items-center gap-2"><span className="w-1.5 h-1.5 rounded-full bg-[#854d0e]" /> PENDING</div>;
         case "draft":
            return <div className="bg-[#0b101e] text-white px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase flex items-center gap-2"> DRAFT</div>;
         case "live":
            return <div className="bg-white border text-red-500 border-red-500 px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase flex items-center gap-2"><span className="w-1.5 h-1.5 rounded-full md:bg-red-500 bg-red-500" /> LIVE</div>;
         case "completed":
            return <div className="bg-gray-300 text-gray-700 px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase flex items-center gap-2"><span className="w-1.5 h-1.5 rounded-full bg-gray-500" /> COMPLETED</div>;
         case "cancellation_pending":
            return <div className="bg-amber-100 text-amber-700 px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase flex items-center gap-2"><span className="w-1.5 h-1.5 rounded-full bg-amber-600" /> CANCELLATION PENDING</div>;
         default:
            return <div className="bg-gray-100 text-gray-700 px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase">{status}</div>;
      }
   };

   return (
      <OrganizerLayout activeTab="My Events">
         {toast && (
            <div className="fixed bottom-6 right-6 bg-gray-900 text-white px-6 py-3 rounded-xl shadow-lg flex items-center gap-3 z-50">
               <span className="w-2 h-2 rounded-full bg-green-400"></span>
               <span className="text-sm font-medium">{toast}</span>
               <button onClick={() => setToast(null)} className="ml-4 text-gray-400 hover:text-white">
                  <X className="w-4 h-4" />
               </button>
            </div>
         )}
         <div className="py-10 px-4 lg:px-10 max-w-6xl mx-auto">
            <div className="flex items-center justify-between mb-8">
               <h1 className="text-2xl font-bold text-gray-900">Browse through your events</h1>
               <button
                  onClick={() => navigate("/organizer/events/create")}
                  className="text-[#e53e5d] font-semibold text-sm hover:underline"
               >
                  + Create Event
               </button>
            </div>

            <div className="flex justify-end mb-6">
               <div className="relative inline-block w-40">
                  <select
                     value={statusFilter}
                     onChange={e => setStatusFilter(e.target.value)}
                     className="w-full appearance-none border border-gray-200 rounded-xl px-4 py-2 text-sm text-gray-900 bg-white outline-none focus:border-gray-400 font-medium cursor-pointer"
                  >
                     {statuses.map(s => <option key={s} value={s}>Status: {s}</option>)}
                  </select>
                  <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500 pointer-events-none" />
               </div>
            </div>

            <div className="flex flex-col gap-6">
               {loading ? (
                  <div className="flex justify-center py-20"><Loader2 className="w-8 h-8 animate-spin text-gray-400" /></div>
               ) : events.length === 0 ? (
                  <div className="text-center py-20 text-gray-500 bg-white rounded-3xl border border-gray-100">No events found.</div>
               ) : events.map(event => {
                  const isExpanded = expandedTickets.has(event.id);
                  return (
                     <div key={event.id} className="bg-white border border-gray-100 rounded-[24px] overflow-hidden hover:shadow-sm transition-shadow">
                        <div className="p-4 flex gap-6">
                           <div className="w-[300px] shrink-0 aspect-[16/9] rounded-xl overflow-hidden bg-gray-100 relative">
                              {event.cover_image_url ? (
                                 <img src={`${import.meta.env.VITE_API_BASE_URL}${event.cover_image_url}`} alt={event.title} className="w-full h-full object-cover" />
                              ) : (
                                 <div className="w-full h-full flex items-center justify-center text-gray-400 text-sm">No Poster Available</div>
                              )}
                           </div>
                           <div className="flex-1 flex flex-col justify-center py-2">
                              <div>
                                 <h2 className="text-[20px] font-bold tracking-tight text-gray-900 mb-2 leading-tight">{event.title}</h2>
                                 <p className="text-[13px] font-medium text-gray-600 mb-1">
                                    {new Date(event.start_time).toLocaleString("en-US", {
                                       day: '2-digit', month: 'short', year: 'numeric',
                                       hour: '2-digit', minute: '2-digit'
                                    })}
                                 </p>
                                 <p className="text-[13px] text-gray-500 flex items-center gap-1.5 mt-2"><MapPin className="w-3.5 h-3.5" /> {event.venue_name}, {event.city}</p>
                                 {event.available_capacity !== undefined && (
                                    <p className="text-[12px] text-[#e53e5d] flex items-center gap-1.5 mt-2 font-semibold bg-red-50 w-fit px-2.5 py-1 rounded-lg border border-red-100">
                                       Tickets left: {event.available_capacity}
                                    </p>
                                 )}
                              </div>
                              <div className="mt-4 flex items-center justify-between">
                                 <div className="flex items-center gap-4">
                                    {getStatusBadge(event.status)}
                                 </div>
                                 <button
                                    onClick={() => toggleTickets(event.id)}
                                    className="flex items-center gap-1.5 text-xs font-semibold text-gray-500 hover:text-gray-800 border border-gray-200 rounded-lg px-3 py-1.5 hover:bg-gray-50 transition"
                                 >
                                    <Ticket className="w-3.5 h-3.5" />
                                    Ticket Details
                                    {isExpanded ? <ChevronUp className="w-3 h-3" /> : <ChevronDown className="w-3 h-3" />}
                                 </button>
                              </div>
                           </div>
                           <div className="w-[180px] shrink-0 flex flex-col justify-center gap-2 border-l border-gray-100 pl-6 py-2">
                              <button
                                 onClick={() => navigate(`/events/${event.slug}`)}
                                 className="w-full bg-[#0b101e] hover:bg-black text-white px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors shadow-sm"
                              >
                                 View Event
                              </button>
                              {(event.status || "").toLowerCase() === "completed" ? (
                                 <button className="w-full border border-gray-300 text-gray-700 hover:bg-gray-50 px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors mt-1">
                                    User Feedbacks
                                 </button>
                              ) : (
                                 <button
                                    onClick={() => navigate(`/organizer/events/${event.id}/edit`)}
                                    disabled={event.status?.toLowerCase() === 'live'}
                                    className="w-full border border-gray-800 text-gray-900 hover:bg-gray-50 px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors mt-1 disabled:opacity-50 disabled:cursor-not-allowed"
                                 >
                                    Edit Event
                                 </button>
                              )}
                              {["draft", "rejected"].includes((event.status || "").toLowerCase()) && (
                                 <button
                                    onClick={() => handleSubmit(event.id)}
                                    className="w-full bg-[#e53e5d] hover:bg-[#d03550] text-white px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors mt-1"
                                 >
                                    Submit for Approval
                                 </button>
                              )}
                              {(event.status || "").toLowerCase() === "live" && (
                                 <button
                                    onClick={() => navigate(`/organizer/events/${event.id}/check-in`)}
                                    className="w-full border border-emerald-500 text-emerald-700 hover:bg-emerald-50 px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors mt-1"
                                 >
                                    Check In Tickets
                                 </button>
                              )}
                              {(event.status || "").toLowerCase() === "live" && allowEventCancellation && (
                                 <button
                                    onClick={() => {
                                       setSelectedEventId(event.id);
                                       setCancellationReason("");
                                       setCancellationModalOpen(true);
                                    }}
                                    className="w-full border border-amber-500 text-amber-700 hover:bg-amber-50 px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors mt-1"
                                 >
                                    Request Cancellation
                                 </button>
                              )}
                              {(event.status || "").toLowerCase() === "cancellation_pending" && event.cancellation_request_reason && (
                                 <button
                                    onClick={() => {
                                       setSelectedReason(event.cancellation_request_reason);
                                       setRejectionModalOpen(true);
                                    }}
                                    className="w-full border border-amber-500 text-amber-700 hover:bg-amber-50 px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors mt-1"
                                 >
                                    View Cancellation Request
                                 </button>
                              )}
                              {(event.status || "").toLowerCase() === "rejected" && event.rejection_reason && (
                                 <button
                                    onClick={() => {
                                       setSelectedReason(event.rejection_reason);
                                       setRejectionModalOpen(true);
                                    }}
                                    className="w-full border border-red-500 text-red-500 hover:bg-red-50 px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors mt-1"
                                 >
                                    Show Reason
                                 </button>
                              )}
                               {!["live", "cancellation_pending"].includes((event.status || "").toLowerCase()) && (
                                 <button
                                    onClick={() => handleDelete(event.id)}
                                    className="w-full text-[#e53e5d] hover:text-[#d03550] text-sm font-semibold mt-3 p-1"
                                 >
                                    Delete Event
                                 </button>
                              )}
                           </div>
                        </div>

                        {isExpanded && (
                           <div className="px-6 pb-5 border-t border-gray-100 pt-4 bg-gray-50/50">
                              <h3 className="text-xs font-bold text-gray-500 uppercase tracking-widest mb-3 flex items-center gap-1.5">
                                 <Ticket className="w-3.5 h-3.5" /> Ticket Types &amp; Capacity
                              </h3>
                              <EventTicketDetails slug={event.slug} />
                           </div>
                        )}
                     </div>
                  );
               })}
            </div>
         </div>

         {rejectionModalOpen && (
            <div className="fixed inset-0 bg-black/50 z-[100] flex items-center justify-center p-4">
               <div className="bg-white rounded-2xl w-full max-w-md shadow-xl overflow-hidden animate-in fade-in zoom-in duration-200">
                   <div className="flex items-center justify-between p-6 border-b border-gray-100">
                       <h3 className="text-lg font-bold text-gray-900">Rejection Reason</h3>
                       <button onClick={() => setRejectionModalOpen(false)} className="text-gray-400 hover:text-gray-600 transition-colors">
                           <X className="w-5 h-5" />
                       </button>
                   </div>
                   <div className="p-6">
                       <div className="bg-red-50 text-red-700 p-4 rounded-xl text-sm leading-relaxed border border-red-100">
                           {selectedReason}
                       </div>
                   </div>
                   <div className="p-6 bg-gray-50 border-t border-gray-100 flex justify-end">
                       <button onClick={() => setRejectionModalOpen(false)} className="bg-gray-900 text-white px-6 py-2 rounded-xl text-sm font-medium hover:bg-black transition-colors">
                           Close
                       </button>
                   </div>
               </div>
            </div>
         )}

         {cancellationModalOpen && (
            <div className="fixed inset-0 bg-black/50 z-[100] flex items-center justify-center p-4">
               <div className="bg-white rounded-2xl w-full max-w-md shadow-xl overflow-hidden">
                  <div className="flex items-center justify-between p-6 border-b border-gray-100">
                     <h3 className="text-lg font-bold text-gray-900">Request Event Cancellation</h3>
                     <button onClick={() => setCancellationModalOpen(false)} className="text-gray-400 hover:text-gray-600">
                        <X className="w-5 h-5" />
                     </button>
                  </div>
                  <div className="p-6">
                     <label className="block text-sm font-semibold text-gray-900 mb-2">Reason <span className="text-red-500">*</span></label>
                     <textarea
                        value={cancellationReason}
                        onChange={(e) => setCancellationReason(e.target.value)}
                        className="w-full min-h-[120px] border border-gray-200 rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-amber-300"
                        placeholder="Explain why this live event should be cancelled..."
                     />
                  </div>
                  <div className="p-6 border-t border-gray-100 flex justify-end gap-2">
                     <button onClick={() => setCancellationModalOpen(false)} className="px-4 py-2 rounded-lg border border-gray-200 text-sm">Close</button>
                     <button
                        onClick={async () => {
                           if (!selectedEventId || !cancellationReason.trim()) return;
                           try {
                              await organizerApi.requestEventCancellation(selectedEventId, cancellationReason.trim());
                              setCancellationModalOpen(false);
                              setToast("Cancellation request sent to admin.");
                              loadEvents();
                           } catch (err) {
                               const msg = (err as any)?.response?.data?.error?.message
                                  || (err as any)?.response?.data?.message
                                  || "Failed to submit cancellation request.";
                               alert(msg);
                           }
                        }}
                        className="px-4 py-2 rounded-lg bg-amber-600 text-white text-sm font-semibold hover:bg-amber-700"
                     >
                        Submit Request
                     </button>
                  </div>
               </div>
            </div>
         )}
      </OrganizerLayout>
   );
}
