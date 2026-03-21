import { MapPin, Star } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { Link } from "react-router-dom"
import CancellationModal from "../components/CancellationModal"
import UserDashboardShell from "../components/UserDashboardShell"
import TicketModal from "../components/TicketModal"
import { userApi } from "../api"
import { getFeedbackMap, saveFeedbackMap, type FeedbackRecord } from "../feedbackStorage"
import {
  enrichBookings,
  fallbackBookings,
  fallbackTicketsByEvent,
  type BookingRecord,
  type EnrichedBooking,
  type TicketRecord,
} from "../userDashboardData"

type TicketModalState = {
  booking: EnrichedBooking
  tickets: TicketRecord[]
}

type CancellationModalState = {
  booking: EnrichedBooking
  tickets: TicketRecord[]
}

const buildFallbackTicket = (booking: EnrichedBooking): TicketRecord => ({
  ticket_id: `${booking.booking_id}-ticket`,
  ticket_code: booking.booking_id.slice(0, 4).toUpperCase(),
  event_id: booking.event_id,
  event_title: booking.event_title,
  ticket_type: "Standard",
  status: "active",
})

export default function MyBookingsPage() {
  const [activeTab, setActiveTab] = useState<"upcoming" | "finished">("upcoming")
  const [bookings, setBookings] = useState<BookingRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [feedbackMap, setFeedbackMap] = useState<Record<string, FeedbackRecord>>({})
  const [draftFeedback, setDraftFeedback] = useState<Record<string, FeedbackRecord>>({})
  const [ticketModal, setTicketModal] = useState<TicketModalState | null>(null)
  const [cancellationModal, setCancellationModal] = useState<CancellationModalState | null>(null)
  const [loadingTickets, setLoadingTickets] = useState<string | null>(null)
  const [loadingCancellation, setLoadingCancellation] = useState<string | null>(null)
  const [eventTickets, setEventTickets] = useState<Record<string, TicketRecord[]>>({})

  useEffect(() => {
    setFeedbackMap(getFeedbackMap())
  }, [])

  useEffect(() => {
    const loadBookings = async () => {
      try {
        const data = await userApi.getMyBookings()
        setBookings(data.length > 0 ? data : fallbackBookings)
      } catch {
        setBookings(fallbackBookings)
      } finally {
        setLoading(false)
      }
    }

    void loadBookings()
  }, [])

  const enrichedBookings = useMemo(() => enrichBookings(bookings), [bookings])
  const filteredBookings = enrichedBookings.filter((booking) => booking.bookingStatus === activeTab)

  const ensureDraft = (bookingId: string) => {
    const current = draftFeedback[bookingId] ?? feedbackMap[bookingId] ?? { rating: 0, comment: "" }
    setDraftFeedback((state) => ({ ...state, [bookingId]: current }))
  }

  const updateDraft = (bookingId: string, nextDraft: FeedbackRecord) => {
    setDraftFeedback((state) => ({ ...state, [bookingId]: nextDraft }))
  }

  const submitFeedback = (bookingId: string) => {
    const current = draftFeedback[bookingId]
    if (!current || current.rating === 0 || !current.comment.trim()) return

    const next = { ...feedbackMap, [bookingId]: current }
    setFeedbackMap(next)
    saveFeedbackMap(next)
  }

  const openTickets = async (booking: EnrichedBooking) => {
    setLoadingTickets(booking.booking_id)

    try {
      const apiTickets = await userApi.getMyTickets(booking.event_id)
      const tickets =
        apiTickets.length > 0
          ? apiTickets
          : fallbackTicketsByEvent[booking.event_id] ?? [buildFallbackTicket(booking)]
      setEventTickets((current) => ({ ...current, [booking.booking_id]: tickets }))
      setTicketModal({ booking, tickets })
    } catch {
      const tickets = fallbackTicketsByEvent[booking.event_id] ?? [buildFallbackTicket(booking)]
      setEventTickets((current) => ({ ...current, [booking.booking_id]: tickets }))
      setTicketModal({
        booking,
        tickets,
      })
    } finally {
      setLoadingTickets(null)
    }
  }

  const openCancellation = async (booking: EnrichedBooking) => {
    setLoadingCancellation(booking.booking_id)

    try {
      const apiTickets = await userApi.getMyTickets(booking.event_id)
      const tickets =
        apiTickets.length > 0
          ? apiTickets
          : eventTickets[booking.booking_id] ?? fallbackTicketsByEvent[booking.event_id] ?? [buildFallbackTicket(booking)]
      setEventTickets((current) => ({ ...current, [booking.booking_id]: tickets }))
      setCancellationModal({ booking, tickets })
    } catch {
      const tickets = eventTickets[booking.booking_id] ?? fallbackTicketsByEvent[booking.event_id] ?? [buildFallbackTicket(booking)]
      setEventTickets((current) => ({ ...current, [booking.booking_id]: tickets }))
      setCancellationModal({ booking, tickets })
    } finally {
      setLoadingCancellation(null)
    }
  }

  const handleConfirmCancellation = (
    selection: Array<{ ticketType: string; cancelCount: number; refundAmount: number }>,
  ) => {
    if (!cancellationModal) return

    const { booking, tickets } = cancellationModal
    const remainingTickets = [...tickets]

    selection.forEach((item) => {
      let remainingToRemove = item.cancelCount

      for (let index = remainingTickets.length - 1; index >= 0 && remainingToRemove > 0; index -= 1) {
        if (remainingTickets[index].ticket_type === item.ticketType) {
          remainingTickets.splice(index, 1)
          remainingToRemove -= 1
        }
      }
    })

    const totalRefund = selection.reduce((sum, item) => sum + item.refundAmount, 0)

    setEventTickets((current) => ({
      ...current,
      [booking.booking_id]: remainingTickets,
    }))

    setBookings((current) =>
      current
        .map((item) =>
          item.booking_id === booking.booking_id
            ? {
                ...item,
                ticket_count: Math.max(0, item.ticket_count - selection.reduce((sum, entry) => sum + entry.cancelCount, 0)),
                total_amount: Math.max(0, item.total_amount - totalRefund),
              }
            : item,
        )
        .filter((item) => item.ticket_count > 0),
    )

    setCancellationModal(null)
  }

  return (
    <UserDashboardShell>
      <div className="pb-10">
        <div className="mb-8 flex items-center gap-8 text-[2rem]">
          {(["upcoming", "finished"] as const).map((tab) => (
            <button
              key={tab}
              type="button"
              onClick={() => setActiveTab(tab)}
              className={`border-b-2 px-2 pb-3 capitalize transition ${
                activeTab === tab ? "border-[#ff445d] text-[#111111]" : "border-transparent text-[#6d6d6d]"
              }`}
            >
              {tab}
            </button>
          ))}
        </div>

        <div className="flex flex-col gap-6">
          {loading ? (
            <div className="flex h-40 items-center justify-center rounded-[28px] border border-gray-100 bg-white shadow-sm">
              <div className="h-8 w-8 animate-spin rounded-full border-2 border-gray-200 border-t-[#111827]" />
            </div>
          ) : null}

          {!loading && filteredBookings.map((booking) => {
            const savedFeedback = feedbackMap[booking.booking_id]
            const currentDraft = draftFeedback[booking.booking_id] ?? savedFeedback ?? { rating: 0, comment: "" }
            const showFeedbackForm = booking.bookingStatus === "finished" && (!savedFeedback || draftFeedback[booking.booking_id])

            return (
              <div
                key={booking.booking_id}
                className="overflow-hidden rounded-[28px] border border-[#ececec] bg-white shadow-[0_12px_36px_rgba(15,23,42,0.08)]"
              >
                <div className="grid md:grid-cols-[370px_1fr]">
                  <div className="relative min-h-[190px]">
                    <img src={booking.coverImageUrl} alt={booking.event_title} className="h-full w-full object-cover" />
                    <div className="absolute left-5 top-5 rounded-2xl bg-[#9d7a7a]/80 px-4 py-2 text-[1.55rem] text-white backdrop-blur-sm">
                      {booking.dateBadge}
                    </div>
                  </div>

                  <div className="flex flex-col justify-between gap-6 p-8">
                    <div>
                      <h2 className="text-[2rem] font-semibold text-[#111827]">{booking.event_title}</h2>
                      <p className="mt-1 text-[1.45rem] text-[#111827]">{booking.timeLabel}</p>
                      <div className="mt-2 flex items-center gap-2 text-[1.35rem] text-[#5d6573]">
                        <MapPin className="h-5 w-5 text-[#ff445d]" />
                        <span>{booking.venue}</span>
                      </div>
                      <div className="mt-4 flex flex-wrap gap-3">
                        {booking.tags.map((tag) => (
                          <span key={tag} className="rounded-full bg-[#efefef] px-4 py-1.5 text-[1.15rem] text-[#5b6069]">
                            {tag}
                          </span>
                        ))}
                      </div>
                    </div>

                    <div className="flex items-center justify-end gap-4">
                      {booking.bookingStatus === "upcoming" ? (
                        <button
                          type="button"
                          onClick={() => void openCancellation(booking)}
                          className="text-[1.45rem] text-[#ff2020] transition hover:opacity-80"
                        >
                          {loadingCancellation === booking.booking_id ? "Loading..." : "Cancel Ticket"}
                        </button>
                      ) : savedFeedback ? (
                        <button
                          type="button"
                          className="rounded-[18px] border border-[#111827] px-6 py-3 text-[1.4rem] text-[#111827] transition hover:bg-[#f7f7f7]"
                          onClick={() => ensureDraft(booking.booking_id)}
                        >
                          Edit Review
                        </button>
                      ) : (
                        <button
                          type="button"
                          className="rounded-[18px] border border-[#111827] px-6 py-3 text-[1.4rem] text-[#111827] transition hover:bg-[#f7f7f7]"
                          onClick={() => ensureDraft(booking.booking_id)}
                        >
                          Write Review
                        </button>
                      )}

                      <button
                        type="button"
                        onClick={() => void openTickets(booking)}
                        className="rounded-[18px] bg-[#111827] px-6 py-3 text-[1.4rem] text-white transition hover:bg-black"
                      >
                        {loadingTickets === booking.booking_id ? "Loading..." : booking.ticket_count > 1 ? "View Tickets" : "View Ticket"}
                      </button>

                      <Link
                        to={`/events/${booking.event_id}`}
                        className="rounded-[18px] bg-[#111827] px-6 py-3 text-[1.4rem] text-white transition hover:bg-black"
                      >
                        View Event
                      </Link>
                    </div>
                  </div>
                </div>

                {showFeedbackForm ? (
                  <div className="border-t border-[#ececec] px-8 py-6">
                    <div className="mb-5 flex items-center gap-2">
                      {Array.from({ length: 5 }, (_, index) => index + 1).map((value) => (
                        <button
                          key={value}
                          type="button"
                          onClick={() => updateDraft(booking.booking_id, { ...currentDraft, rating: value })}
                          className="text-[#111111]"
                        >
                          <Star
                            className={`h-8 w-8 ${value <= currentDraft.rating ? "fill-[#111111] text-[#111111]" : "text-[#111111]"}`}
                          />
                        </button>
                      ))}
                    </div>

                    <textarea
                      value={currentDraft.comment}
                      onChange={(event) =>
                        updateDraft(booking.booking_id, { ...currentDraft, comment: event.target.value })
                      }
                      placeholder="Write your feedback here..."
                      className="min-h-[120px] w-full resize-none rounded-2xl border border-transparent px-2 py-2 text-[1.5rem] text-[#6e7480] outline-none"
                    />

                    <div className="mt-4 flex justify-end">
                      <button
                        type="button"
                        onClick={() => submitFeedback(booking.booking_id)}
                        className="rounded-[18px] bg-[#111827] px-8 py-3 text-[1.4rem] text-white transition hover:bg-black"
                      >
                        Submit Feedback
                      </button>
                    </div>
                  </div>
                ) : null}
              </div>
            )
          })}
        </div>
      </div>

      <TicketModal
        open={ticketModal !== null}
        booking={ticketModal?.booking ?? null}
        tickets={ticketModal?.tickets ?? []}
        onClose={() => setTicketModal(null)}
      />
      <CancellationModal
        open={cancellationModal !== null}
        booking={cancellationModal?.booking ?? null}
        tickets={cancellationModal?.tickets ?? []}
        onClose={() => setCancellationModal(null)}
        onConfirm={handleConfirmCancellation}
      />
    </UserDashboardShell>
  )
}
