import { MapPin, Star } from "lucide-react"
import { useEffect, useState } from "react"
import { useQueryClient } from "@tanstack/react-query"
import { Link } from "react-router-dom"
import CancellationModal from "../components/CancellationModal"
import UserDashboardShell from "../components/UserDashboardShell"
import TicketModal from "../components/TicketModal"
import { userApi } from "../api"
import { getFeedbackMap, saveFeedbackMap, type FeedbackRecord } from "../feedbackStorage"
import { walletQueryKey, walletTransactionsQueryKey } from "../hooks"
import {
  formatDateBadge,
  formatEventTime,
  type BookingRecord,
  type TicketRecord,
} from "../userDashboardData"

type TicketModalState = {
  booking: BookingRecord
  tickets: TicketRecord[]
}

type CancellationModalState = {
  booking: BookingRecord
  tickets: TicketRecord[]
}

export default function MyBookingsPage() {
  const queryClient = useQueryClient()
  const [activeTab, setActiveTab] = useState<"upcoming" | "finished">("upcoming")
  const [bookings, setBookings] = useState<BookingRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [refundNotice, setRefundNotice] = useState("")
  const [bookingActionError, setBookingActionError] = useState("")
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
        const validBookings = (data || []).filter(b => b.ticket_count > 0)
        setBookings(validBookings)
      } catch {
        setBookings([])
      } finally {
        setLoading(false)
      }
    }

    void loadBookings()
  }, [])

  const filteredBookings = bookings.filter((booking) => {
    const bookingDate = new Date(booking.event_start_time)
    const bookingStatus = bookingDate.getTime() >= Date.now() ? "upcoming" : "finished"
    return bookingStatus === activeTab
  })

  const ensureDraft = (bookingId: string) => {
    setDraftFeedback((state) => {
      if (state[bookingId]) {
        const nextState = { ...state }
        delete nextState[bookingId]
        return nextState
      }

      return {
        ...state,
        [bookingId]: feedbackMap[bookingId] ?? { rating: 0, comment: "" },
      }
    })
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
    setDraftFeedback((state) => {
      const nextState = { ...state }
      delete nextState[bookingId]
      return nextState
    })
  }

  const openTickets = async (booking: BookingRecord) => {
    setLoadingTickets(booking.booking_id)

    try {
      const tickets = await userApi.getMyTickets(undefined, booking.booking_id)
      setEventTickets((current) => ({ ...current, [booking.booking_id]: tickets || [] }))
      setTicketModal({ booking, tickets: tickets || [] })
    } catch {
      setTicketModal({ booking, tickets: [] })
    } finally {
      setLoadingTickets(null)
    }
  }

  const openCancellation = async (booking: BookingRecord) => {
    setLoadingCancellation(booking.booking_id)

    try {
      const apiTickets = await userApi.getMyTickets(undefined, booking.booking_id)
      const tickets = apiTickets.length > 0 ? apiTickets : eventTickets[booking.booking_id] ?? []
      setEventTickets((current) => ({ ...current, [booking.booking_id]: tickets }))
      setCancellationModal({ booking, tickets })
    } catch {
      const tickets = eventTickets[booking.booking_id] ?? []
      setCancellationModal({ booking, tickets })
    } finally {
      setLoadingCancellation(null)
    }
  }

  const handleConfirmCancellation = async (
    selection: Array<{ ticketType: string; cancelCount: number; refundAmount: number }>,
  ) => {
    if (!cancellationModal) return

    const { booking } = cancellationModal
    const totalCancelled = selection.reduce((sum, item) => sum + item.cancelCount, 0)
    const shouldRefundToWallet = totalCancelled > 0 && totalCancelled === cancellationModal.tickets.length

    const items = selection.map(s => ({
      ticket_type: s.ticketType,
      quantity: s.cancelCount
    }))

    setRefundNotice("")
    setBookingActionError("")

    try {
      await userApi.cancelBooking(booking.booking_id, { items })

      if (shouldRefundToWallet) {
        await userApi.refundBooking(booking.booking_id)
        setRefundNotice("Refunded to wallet")
        await Promise.all([
          queryClient.invalidateQueries({ queryKey: walletQueryKey }),
          queryClient.invalidateQueries({ queryKey: walletTransactionsQueryKey }),
        ])
      }

      const data = await userApi.getMyBookings()
      const validBookings = (data || []).filter(b => b.ticket_count > 0)
      setBookings(validBookings)
    } catch (error: any) {
      setBookingActionError(
        error?.response?.data?.error?.message || "Failed to process cancellation.",
      )
    }

    setCancellationModal(null)
  }

  return (
    <UserDashboardShell>
      <div className="pb-10">
        <div className="mb-5 flex items-center gap-5 text-xl">
          {(["upcoming", "finished"] as const).map((tab) => (
            <button
              key={tab}
              type="button"
              onClick={() => setActiveTab(tab)}
              className={`border-b-2 px-2 pb-2 capitalize transition ${activeTab === tab ? "border-[#ff445d] text-[#111111]" : "border-transparent text-[#6d6d6d]"
                }`}
            >
              {tab}
            </button>
          ))}
        </div>

        {refundNotice ? (
          <div className="mb-5 rounded-2xl border border-[#d9f3e3] bg-[#f1fbf5] px-4 py-3 text-sm font-medium text-[#118a43]">
            {refundNotice}
          </div>
        ) : null}

        {bookingActionError ? (
          <div className="mb-5 rounded-2xl border border-[#ffd7dd] bg-[#fff5f7] px-4 py-3 text-sm font-medium text-[#d22d4c]">
            {bookingActionError}
          </div>
        ) : null}

        <div className="flex flex-col gap-6">
          {loading ? (
            <div className="flex h-40 items-center justify-center rounded-[28px] border border-gray-100 bg-white shadow-sm">
              <div className="h-8 w-8 animate-spin rounded-full border-2 border-gray-200 border-t-[#111827]" />
            </div>
          ) : null}

          {!loading && filteredBookings.length === 0 && (
            <div className="flex flex-col items-center justify-center py-20 bg-white rounded-3xl border border-gray-100 shadow-sm text-center">
              <h3 className="text-xl font-medium text-gray-900 mb-2">
                {activeTab === "upcoming" ? "No upcoming events" : "No completed events"}
              </h3>
              <p className="text-gray-500">
                {activeTab === "upcoming"
                  ? "You don't have any upcoming bookings yet."
                  : "You haven't attended any events yet."}
              </p>
            </div>
          )}

          {!loading && filteredBookings.map((booking) => {
            const savedFeedback = feedbackMap[booking.booking_id]
            const currentDraft = draftFeedback[booking.booking_id] ?? savedFeedback ?? { rating: 0, comment: "" }
            const isFinished = new Date(booking.event_start_time).getTime() < Date.now()
            const showFeedbackForm = isFinished && !!draftFeedback[booking.booking_id]

            return (
              <div
                key={booking.booking_id}
                className="overflow-hidden rounded-[20px] border border-[#ececec] bg-white shadow-[0_10px_24px_rgba(15,23,42,0.06)]"
              >
                <div className="grid md:grid-cols-[200px_1fr]">
                  <div className="relative min-h-[150px]">
                    {booking.coverImageUrl ? (
                      <img
                        src={booking.coverImageUrl?.startsWith("/") ? `${import.meta.env.VITE_API_BASE_URL}${booking.coverImageUrl}` : booking.coverImageUrl}
                        alt={booking.event_title}
                        className="h-full w-full object-cover"
                      />
                    ) : (
                      <div className="flex h-full w-full items-center justify-center bg-gradient-to-br from-gray-700 to-gray-900">
                        <span className="text-4xl font-bold text-white opacity-30">
                          {(booking.event_title || "E")[0].toUpperCase()}
                        </span>
                      </div>
                    )}
                    <div className="absolute left-3 top-3 rounded-xl bg-[#9d7a7a]/80 px-3 py-1.5 text-sm text-white backdrop-blur-sm">
                      {formatDateBadge(booking.event_start_time)}
                    </div>
                  </div>

                  <div className="flex flex-col justify-between gap-4 p-5">
                    <div>
                      <h2 className="text-xl font-semibold leading-tight text-[#111827]">{booking.event_title}</h2>
                      <p className="mt-1 text-base text-[#111827]">{formatEventTime(booking.event_start_time)}</p>
                      <div className="mt-2 flex items-center gap-2 text-sm text-[#5d6573]">
                        <MapPin className="h-4 w-4 text-[#ff445d]" />
                        <span>{booking.venue || booking.event_city}</span>
                      </div>
                      <div className="mt-3 flex flex-wrap gap-2">
                        {Array.isArray(booking.tags) && booking.tags.map((tag: string) => (
                          <span key={tag} className="rounded-full bg-[#efefef] px-3 py-1 text-xs text-[#5b6069]">
                            {tag}
                          </span>
                        ))}
                        <span className="rounded-full bg-[#ff445d]/10 px-3 py-1 text-xs font-medium text-[#ff445d]">
                          {booking.ticket_count} Ticket{booking.ticket_count !== 1 ? "s" : ""}
                        </span>
                      </div>
                    </div>

                    <div className="flex flex-wrap items-center justify-end gap-3">
                      {!isFinished ? (
                        <button
                          type="button"
                          onClick={() => void openCancellation(booking)}
                          className="text-sm font-medium text-[#ff2020] transition hover:opacity-80"
                        >
                          {loadingCancellation === booking.booking_id ? "Loading..." : "Cancel Ticket"}
                        </button>
                      ) : savedFeedback ? (
                        <button
                          type="button"
                          className="rounded-2xl border border-[#111827] px-5 py-2.5 text-sm font-medium text-[#111827] transition hover:bg-[#f7f7f7]"
                          onClick={() => ensureDraft(booking.booking_id)}
                        >
                          Edit Review
                        </button>
                      ) : (
                        <button
                          type="button"
                          className="rounded-2xl border border-[#111827] px-5 py-2.5 text-sm font-medium text-[#111827] transition hover:bg-[#f7f7f7]"
                          onClick={() => ensureDraft(booking.booking_id)}
                        >
                          Write Review
                        </button>
                      )}

                      <button
                        type="button"
                        onClick={() => void openTickets(booking)}
                        className="rounded-xl bg-[#111827] px-4 py-2 text-sm font-medium text-white transition hover:bg-black"
                      >
                        {loadingTickets === booking.booking_id ? "Loading..." : booking.ticket_count > 1 ? "View Tickets" : "View Ticket"}
                      </button>

                      <Link
                        to={`/events/${booking.event_id}`}
                        className="rounded-2xl bg-[#111827] px-5 py-2.5 text-sm font-medium text-white transition hover:bg-black"
                      >
                        View Event
                      </Link>
                    </div>
                  </div>
                </div>

                {showFeedbackForm ? (
                  <div className="border-t border-[#ececec] px-6 py-5">
                    <div className="mb-4 flex items-center gap-2">
                      {Array.from({ length: 5 }, (_, index) => index + 1).map((value) => (
                        <button
                          key={value}
                          type="button"
                          onClick={() => updateDraft(booking.booking_id, { ...currentDraft, rating: value })}
                          className="text-[#111111]"
                        >
                          <Star
                            className={`h-6 w-6 ${value <= currentDraft.rating ? "fill-[#111111] text-[#111111]" : "text-[#111111]"}`}
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
                      className="min-h-[96px] w-full resize-none rounded-2xl border border-transparent px-2 py-2 text-sm text-[#6e7480] outline-none"
                    />

                    <div className="mt-4 flex justify-end">
                      <button
                        type="button"
                        onClick={() => submitFeedback(booking.booking_id)}
                        className="rounded-xl bg-[#111827] px-5 py-2 text-sm font-medium text-white transition hover:bg-black"
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
