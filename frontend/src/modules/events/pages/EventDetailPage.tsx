import { useNavigate, useParams } from "react-router-dom";
import { useEvent } from "../hooks";
import { CalendarDays, MapPin, Clock, Hourglass, X } from "lucide-react";
import { useState, useEffect } from "react";
import { buildDisplayEvent, formatCurrency } from "../eventBookingData";
import { useAuthStore } from "../../auth/store/authStore";
import Modal from "../../../shared/ui/Modal";
import { useEngagement } from "../../../shared/hooks/useEngagement";

export default function EventDetailPage() {
  const { eventId } = useParams();
  const { roles } = useAuthStore();
  const isOrganizer = roles.includes("organizer");
  const isAdmin = roles.includes("admin");
  const { data, isLoading, isError } = useEvent(eventId!, isOrganizer, isAdmin);
  const navigate = useNavigate();
  const { trackEvent } = useEngagement();

  useEffect(() => {
    if (data?.event?.id && !isLoading && !isError) {
      trackEvent('event_view', data.event.id);
    }
  }, [data?.event?.id, isLoading, isError, trackEvent]);

  const [activeTab, setActiveTab] = useState("About");
  const [isHostModalOpen, setIsHostModalOpen] = useState(false);

  if (isLoading) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center bg-gray-50">
        <div className="h-8 w-8 animate-spin rounded-full border-b-2 border-gray-900"></div>
      </div>
    );
  }

  if (isError || !data) {
    return (
      <div className="min-h-[60vh] flex items-center justify-center bg-gray-50 px-6">
        <div className="w-full max-w-xl rounded-[24px] border border-[#ececec] bg-white p-8 text-center shadow-[0_12px_32px_rgba(15,23,42,0.08)]">
          <div className="text-sm font-semibold uppercase tracking-[0.24em] text-[#ff445d]">404</div>
          <h1 className="mt-3 text-2xl font-semibold text-[#111827]">Event Unavailable</h1>
          <p className="mt-3 text-sm leading-6 text-[#6b7280]">
            The event you're looking for is either pending approval, has been removed, or just doesn't exist.
          </p>
        </div>
      </div>
    );
  }

  const displayEvent = buildDisplayEvent(eventId ?? "", data);

  return (
    <div className="min-h-screen bg-gray-50 pb-20">
      <div className="max-w-7xl mx-auto px-6 py-10 grid grid-cols-1 lg:grid-cols-3 gap-8">
        <div className="lg:col-span-2 flex flex-col gap-6">
          {isAdmin && (
            <button onClick={() => navigate('/admin/events')} className="flex items-center gap-2 text-sm font-medium text-gray-500 hover:text-gray-900 transition w-fit mb-2">
              <span aria-hidden="true">&larr;</span> Back to Admin Events
            </button>
          )}
          <div className="w-full h-[400px] rounded-2xl overflow-hidden shadow-sm">
            <img
              src={displayEvent.coverImageUrl}
              alt={displayEvent.title}
              className="w-full h-full object-cover"
            />
          </div>

          <h1 className="text-2xl md:text-3xl font-bold text-gray-900 mt-2">
            {displayEvent.title}
          </h1>

          <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mt-2">
            <div className="flex gap-8 border-b border-gray-100 mb-6">
              {["About", "Venue", "Terms & Conditions"].map((tab) => (
                <button
                  key={tab}
                  className={`pb-3 text-sm font-medium transition-colors relative ${
                    activeTab === tab ? "text-[#e53e5d]" : "text-gray-600 hover:text-gray-900"
                  }`}
                  onClick={() => setActiveTab(tab)}
                >
                  {tab}
                  {activeTab === tab && (
                    <div className="absolute bottom-0 left-0 w-full h-0.5 bg-[#e53e5d] rounded-t-full"></div>
                  )}
                </button>
              ))}
            </div>

            {activeTab === "About" && (
              <div className="animate-in fade-in duration-300">
                <h3 className="font-bold text-gray-900 text-sm mb-4">{displayEvent.about?.subtitle || displayEvent.title}</h3>
                <div className="flex flex-col gap-4 text-sm text-gray-700 leading-relaxed">
                  {displayEvent.about?.content ? displayEvent.about.content.map((paragraph: string, index: number) => (
                    <p key={index}>{paragraph}</p>
                  )) : (
                    <p>{displayEvent.description}</p>
                  )}
                </div>
              </div>
            )}
            {activeTab === "Venue" && (
              <div className="animate-in fade-in duration-300 text-sm text-gray-700">
                <div className="flex flex-col gap-4">
                  <div className="font-bold text-base text-gray-900">{displayEvent.venueName}</div>
                  {displayEvent.venueAddress && <p>{displayEvent.venueAddress}</p>}
                  {displayEvent.city && <p>{displayEvent.city}</p>}
                  {displayEvent.mapUrl && (
                    <div className="mt-2 flex flex-col gap-3">
                      <div className="w-full overflow-hidden rounded-xl bg-gray-100 border border-gray-100 shadow-sm">
                        <iframe
                          title="Event Venue Map"
                          src={
                            displayEvent.mapUrl.includes("output=embed") || displayEvent.mapUrl.includes("/embed")
                              ? displayEvent.mapUrl
                              : `https://maps.google.com/maps?q=${encodeURIComponent(
                                  displayEvent.venueAddress ? `${displayEvent.venueName}, ${displayEvent.venueAddress}, ${displayEvent.city || ''}` : `${displayEvent.venueName}, ${displayEvent.city || ''}`
                                )}&output=embed`
                          }
                          width="100%"
                          height="280"
                          style={{ border: 0 }}
                          allowFullScreen
                          loading="lazy"
                          referrerPolicy="no-referrer-when-downgrade"
                        />
                      </div>
                      <a href={displayEvent.mapUrl} target="_blank" rel="noopener noreferrer" className="text-[#e53e5d] font-medium hover:underline inline-block">
                        View Larger Map
                      </a>
                    </div>
                  )}
                </div>
              </div>
            )}
            {activeTab === "Terms & Conditions" && (
              <div className="animate-in fade-in duration-300 text-sm text-gray-700">
                <div className="prose prose-sm max-w-none">
                  {displayEvent.termsAndConditions ? (
                    <div className="whitespace-pre-line">{displayEvent.termsAndConditions}</div>
                  ) : (
                    <p>No terms and conditions specified for this event.</p>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>

        <div className="lg:col-span-1 flex flex-col gap-6">

          {!(useAuthStore.getState().roles || []).some(r => r.toLowerCase() === "organizer" || r.toLowerCase() === "admin") ? (
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
              <h3 className="text-lg font-bold text-gray-900 mb-6">Book Tickets</h3>

              <div className="flex flex-col gap-4 mb-8">
                <div className="flex items-start gap-4 text-sm text-gray-700">
                  <CalendarDays className="w-5 h-5 text-[#e53e5d] shrink-0 mt-0.5" />
                  <span>{displayEvent.dateLabel}</span>
                </div>
                <div className="flex items-start gap-4 text-sm text-gray-700">
                  <Clock className="w-5 h-5 text-[#e53e5d] shrink-0 mt-0.5" />
                  <span>{displayEvent.timeLabel}</span>
                </div>
                <div className="flex items-start gap-4 text-sm text-gray-700">
                  <MapPin className="w-5 h-5 text-[#e53e5d] shrink-0 mt-0.5" />
                  <span>{displayEvent.displayLocation}</span>
                </div>
                <div className="flex items-start gap-4 text-sm text-gray-700">
                  <Hourglass className="w-5 h-5 text-[#e53e5d] shrink-0 mt-0.5" />
                  <span>{displayEvent.durationLabel || "N/A"}</span>
                </div>
              </div>

              <div className="bg-[#fcf3f4] rounded-xl p-4 mb-4 flex justify-between items-center">
                <span className="text-sm font-medium text-gray-700">Price From</span>
                <span className="text-sm font-bold text-[#e53e5d]">
                  {displayEvent.ticketTypes && displayEvent.ticketTypes.length > 0
                    ? formatCurrency(Math.min(...displayEvent.ticketTypes.map((t: any) => t.price)))
                    : `₹ ${displayEvent.priceLabel || "N/A"}`}
                </span>
              </div>

              {(() => {
                const isSoldOut = displayEvent.ticketTypes?.length > 0 &&
                                 displayEvent.ticketTypes.every((t: any) => (t.availableQuantity ?? 0) <= 0);

                const isSellingFast = displayEvent.availableCapacity !== undefined &&
                                      displayEvent.totalCapacity !== undefined &&
                                      displayEvent.availableCapacity > 0 &&
                                      displayEvent.availableCapacity <= displayEvent.totalCapacity * 0.2;

                return (
                  <div className="flex flex-col gap-3">
                    {isSellingFast && !isSoldOut && displayEvent.status !== "completed" && (
                      <div className="bg-orange-50 border border-orange-200 text-orange-700 px-4 py-2.5 rounded-xl text-sm font-semibold flex items-center justify-between">
                         <span className="flex items-center gap-2">
                            <span className="w-2 h-2 rounded-full bg-orange-500 animate-pulse"></span>
                            Selling Fast!
                         </span>
                         <span>Only {displayEvent.availableCapacity} left</span>
                      </div>
                    )}
                    {displayEvent.status === "completed" ? (
                      <button
                        disabled
                        className="w-full py-3.5 rounded-xl text-sm font-medium bg-gray-100 text-gray-500 cursor-not-allowed"
                      >
                        Event Completed
                      </button>
                    ) : (
                      <button
                        disabled={isSoldOut}
                        className={`w-full py-3.5 rounded-xl text-sm font-medium transition-colors ${
                          isSoldOut
                            ? "bg-gray-200 text-gray-500 cursor-not-allowed"
                            : "bg-[#0b101e] hover:bg-black text-white"
                        }`}
                        onClick={() => {
                          if (useAuthStore.getState().isAuthenticated) {
                            navigate(`/events/${eventId}/book`)
                          } else {
                            useAuthStore.getState().openAuthModal("goer")
                          }
                        }}
                      >
                        {isSoldOut ? "Sold Out" : "Continue to Booking"}
                      </button>
                    )}
                  </div>
                );
              })()}
            </div>
          ) : (
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
              <h3 className="text-lg font-bold text-[#e53e5d] mb-6">Organizer Ticket Details</h3>
              <div className="flex flex-col gap-4">
                {displayEvent.ticketTypes && displayEvent.ticketTypes.length > 0 ? (
                  displayEvent.ticketTypes.map((t: any) => (
                    <div key={t.id} className="flex justify-between items-center bg-gray-50 p-3 rounded-lg border border-gray-100">
                      <div>
                        <p className="font-semibold text-gray-900 text-sm">{t.name}</p>
                        <p className="text-xs text-gray-500">Price: {formatCurrency(t.price)}</p>
                      </div>
                      <div className="text-right">
                        <p className="text-sm font-bold text-gray-900">{t.availableQuantity} / {t.totalQuantity}</p>
                        <p className="text-[10px] uppercase tracking-wide text-gray-500">Available</p>
                      </div>
                    </div>
                  ))
                ) : (
                  <p className="text-sm text-gray-500">No ticket types defined.</p>
                )}
                {displayEvent.totalCapacity !== undefined && (
                  <div className="mt-4 pt-4 border-t border-gray-100 flex justify-between font-bold text-gray-900">
                    <span>Total Capacity</span>
                    <span>{displayEvent.availableCapacity} / {displayEvent.totalCapacity} Available</span>
                  </div>
                )}
              </div>
            </div>
          )}

          {displayEvent.host && (
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
              <h3 className="text-lg font-bold text-gray-900 mb-6">Host</h3>
              <div className="flex items-center gap-4 mb-6">
                <img src={displayEvent.host.avatar} alt={displayEvent.host.name} className="w-12 h-12 rounded-full object-cover" />
                <div>
                  <h4 className="font-bold text-gray-900 text-sm">{displayEvent.host.name}</h4>
                  <p className="text-xs text-gray-500">{displayEvent.host.role}</p>
                </div>
              </div>
              <button onClick={() => setIsHostModalOpen(true)} className="w-full border border-gray-300 hover:bg-gray-50 text-gray-700 py-2.5 rounded-xl text-sm font-medium transition-colors">
                View Profile
              </button>
            </div>
          )}

          {displayEvent.personnel && (
            <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6">
              <h3 className="text-lg font-bold text-gray-900 mb-6">Key Personnel</h3>
              <div className="flex flex-col gap-6">
                {displayEvent.personnel.map((person: any, index: number) => (
                  <div key={index} className="flex items-center gap-4">
                    <img src={person.avatar} alt={person.name} className="w-12 h-12 rounded-full object-cover" />
                    <div>
                      <h4 className="font-bold text-gray-900 text-sm">{person.name}</h4>
                      <p className="text-xs text-gray-500">{person.role}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}

        </div>
      </div>

      <Modal open={isHostModalOpen} onClose={() => setIsHostModalOpen(false)} className="relative w-[min(92vw,400px)] rounded-2xl bg-white p-6 shadow-[0_28px_90px_rgba(15,23,42,0.22)]">
        <button
          type="button"
          aria-label="Close modal"
          className="absolute right-4 top-4 rounded-full p-1 text-[#111111] transition hover:bg-[#f5f5f5]"
          onClick={() => setIsHostModalOpen(false)}
        >
          <X className="h-5 w-5" />
        </button>
        {displayEvent.host && (
          <div className="flex flex-col items-center text-center">
             <img src={displayEvent.host.avatar} alt={displayEvent.host.name} className="w-24 h-24 rounded-full object-cover shadow-sm mb-4" />
             <h2 className="text-xl font-bold text-gray-900">{displayEvent.host.name}</h2>
             {displayEvent.host.organization && (
               <p className="text-[13px] font-semibold text-gray-700 bg-gray-100 px-3 py-1 rounded-md mt-2 mb-1">{displayEvent.host.organization}</p>
             )}
             <p className={`text-sm font-medium text-[#e53e5d] uppercase tracking-wider mb-4 ${!displayEvent.host.organization && 'mt-2'}`}>{displayEvent.host.role}</p>
             <p className="text-sm text-gray-600 mb-6 leading-relaxed">
               This organizer manages events on EVNTX. Stay tuned for their upcoming events and more information about their organization profile.
             </p>
             <button
                onClick={() => setIsHostModalOpen(false)}
                className="bg-[#0b101e] text-white py-3 px-6 rounded-xl hover:bg-black transition-colors font-medium text-sm w-full"
             >
                Close
             </button>
          </div>
        )}
      </Modal>
    </div>
  );
}
