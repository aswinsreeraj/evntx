import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import OrganizerLayout from "../components/OrganizerLayout";
import { organizerApi } from "../api";
import { X, ChevronDown, MapPin, Loader2 } from "lucide-react";

export default function MyEvents() {
   const navigate = useNavigate();
   const [events, setEvents] = useState<any[]>([]);
   const [loading, setLoading] = useState(true);
   const [statusFilter, setStatusFilter] = useState("All");
   const statuses = ["All", "Draft", "Pending", "Approved", "Rejected", "Live", "Completed"];

   const loadEvents = async () => {
      setLoading(true);
      try {
         const data = await organizerApi.getOrganizerEvents(statusFilter === "All" ? "" : statusFilter);
         setEvents(data.data?.events || []);
      } catch (err) {
         console.error(err);
      } finally {
         setLoading(false);
      }
   };

   useEffect(() => {
      loadEvents();
   }, [statusFilter]);

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
         default:
            return <div className="bg-gray-100 text-gray-700 px-4 py-1.5 rounded-full text-[11px] font-bold tracking-wide uppercase">{status}</div>;
      }
   };

   return (
      <OrganizerLayout activeTab="My Events">
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
               ) : events.map(event => (
                  <div key={event.id} className="bg-white border border-gray-100 rounded-[24px] p-4 flex gap-6 hover:shadow-sm transition-shadow">
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
                        </div>
                        <div className="mt-4 flex items-center justify-between">
                           <div className="flex items-center gap-4">
                              {getStatusBadge(event.status)}
                           </div>
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
                              className="w-full border border-gray-800 text-gray-900 hover:bg-gray-50 px-4 py-2.5 rounded-xl text-sm font-semibold transition-colors mt-1"
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
                        <button
                           onClick={() => handleDelete(event.id)}
                           className="w-full text-[#e53e5d] hover:text-[#d03550] text-sm font-semibold mt-3 p-1"
                        >
                           Delete Event
                        </button>
                     </div>
                  </div>
               ))}
            </div>
         </div>
      </OrganizerLayout>
   );
}
