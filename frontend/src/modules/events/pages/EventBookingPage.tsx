import { Minus, Plus, X } from "lucide-react"
import { useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import Modal from "../../../shared/ui/Modal"
import { useEvent } from "../hooks"
import { buildDisplayEvent, formatCurrency } from "../eventBookingData"
import { eventsApi } from "../api"
import { saveBookingConfirmation } from "../bookingStorage"
import { useAuthStore } from "../../auth/store/authStore"

const PLATFORM_FEE_RATE = 0.05

export default function EventBookingPage() {
  const { eventId } = useParams()
  const navigate = useNavigate()
  const { data, isLoading } = useEvent(eventId!)
  const { user } = useAuthStore()

  const displayEvent = buildDisplayEvent(eventId ?? "", data)
  const [quantities, setQuantities] = useState<Record<string, number>>(() =>
    Object.fromEntries(displayEvent.ticketTypes.map((ticket) => [ticket.name, 0])),
  )
  const [checkoutOpen, setCheckoutOpen] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const ticketRows = displayEvent.ticketTypes.map((ticket) => ({
    ...ticket,
    quantity: quantities[ticket.name] ?? 0,
    amount: ticket.price * (quantities[ticket.name] ?? 0),
  }))

  const totalAmount = ticketRows.reduce((sum, ticket) => sum + ticket.amount, 0)
  const platformFee = totalAmount > 0 ? Math.round(totalAmount * PLATFORM_FEE_RATE) : 0
  const finalAmount = totalAmount + platformFee

  const selectedTickets = ticketRows.filter((ticket) => ticket.quantity > 0)

  const updateQuantity = (ticketName: string, nextQuantity: number, limit?: number) => {
    const safeQuantity = Math.max(0, Math.min(limit ?? Number.POSITIVE_INFINITY, nextQuantity))
    setQuantities((current) => ({
      ...current,
      [ticketName]: safeQuantity,
    }))
  }

  const handleProceed = () => {
    if (totalAmount === 0) return
    setError(null)
    setCheckoutOpen(true)
  }

  const handlePayment = async () => {
    if (!eventId || selectedTickets.length === 0) return

    setIsSubmitting(true)
    setError(null)

    try {
      const hasTicketIds = selectedTickets.every((ticket) => ticket.id)
      let bookingId = `BK-${Date.now()}`

      if (hasTicketIds) {
        const response = await eventsApi.reserveTickets({
          eventId,
          tickets: selectedTickets.map((ticket) => ({
            ticket_type_id: ticket.id!,
            quantity: ticket.quantity,
          })),
        })

        bookingId = response.booking_id
      }

      saveBookingConfirmation({
        bookingId,
        eventId,
        eventTitle: displayEvent.title,
        eventImage: displayEvent.coverImageUrl,
        eventDate: displayEvent.dateLabel,
        eventTime: displayEvent.timeLabel,
        venue: displayEvent.displayLocation,
        totalAmount,
        platformFee,
        finalAmount,
        email: user?.name ? `${user.name.toLowerCase().replace(/\s+/g, "")}@example.com` : "johnsmith@example.com",
        tickets: selectedTickets.map((ticket) => ({
          ticketName: ticket.name,
          quantity: ticket.quantity,
          unitPrice: ticket.price,
        })),
        createdAt: new Date().toISOString(),
      })

      navigate(`/events/${eventId}/confirmation`)
    } catch (paymentError: any) {
      setError(
        paymentError?.response?.data?.message ??
          "We could not complete this booking right now. Please try again.",
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="bg-white">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-10 px-6 py-10 md:py-12">
        <section className="overflow-hidden rounded-[28px] border border-[#e8e8e8] bg-white shadow-[0_18px_48px_rgba(15,23,42,0.12)]">
          <div className="grid md:grid-cols-[1.1fr_1fr]">
            <div className="min-h-[240px] bg-[#111827]">
              <img
                src={displayEvent.coverImageUrl}
                alt={displayEvent.title}
                className="h-full w-full object-cover"
              />
            </div>

            <div className="flex flex-col justify-center gap-3 p-8">
              <h1 className="max-w-md text-3xl font-semibold tracking-tight text-[#111111]">
                {displayEvent.title}
              </h1>
              <p className="text-xl text-[#8d949e]">{displayEvent.dateLabel}</p>
              <p className="text-xl text-[#8d949e]">{displayEvent.timeLabel}</p>
              <p className="text-xl text-[#8d949e]">{displayEvent.displayLocation}</p>
            </div>
          </div>
        </section>

        <section className="mx-auto flex w-full max-w-4xl flex-col gap-6">
          <div className="text-center">
            <h2 className="text-4xl font-semibold tracking-tight text-[#111111]">Ticket Selection</h2>
            {isLoading ? <p className="mt-2 text-sm text-[#8d949e]">Loading event details...</p> : null}
          </div>

          <div className="flex flex-col gap-4">
            {ticketRows.map((ticket) => (
              <div
                key={ticket.name}
                className="grid items-center gap-4 rounded-[18px] border border-[#ffc8cf] bg-white px-6 py-5 text-[#111111] md:grid-cols-[1.4fr_1fr_auto]"
              >
                <div className="text-[1.75rem] font-medium">{ticket.name}</div>
                <div className="text-left text-[1.75rem] font-semibold text-[#ff445d] md:text-center">
                  Price: {formatCurrency(ticket.price)}
                </div>
                <div className="ml-auto flex items-center gap-4 text-2xl">
                  <span className="text-[#111111]">Qty</span>
                  <button
                    type="button"
                    aria-label={`Decrease ${ticket.name} tickets`}
                    className="rounded-full p-1 text-[#111111] transition hover:bg-[#f7f7f7]"
                    onClick={() => updateQuantity(ticket.name, ticket.quantity - 1, ticket.availableQuantity)}
                  >
                    <Minus className="h-5 w-5" />
                  </button>
                  <span className="inline-flex min-w-10 items-center justify-center rounded-md bg-[#ffd9df] px-3 py-1 text-[1.75rem] font-medium text-[#111111]">
                    {ticket.quantity}
                  </span>
                  <button
                    type="button"
                    aria-label={`Increase ${ticket.name} tickets`}
                    className="rounded-full p-1 text-[#111111] transition hover:bg-[#f7f7f7]"
                    onClick={() => updateQuantity(ticket.name, ticket.quantity + 1, ticket.availableQuantity)}
                  >
                    <Plus className="h-5 w-5" />
                  </button>
                </div>
              </div>
            ))}
          </div>

          <div className="flex justify-end">
            <div className="text-right text-2xl font-semibold text-[#111111]">
              Total Amount: {formatCurrency(totalAmount)}
            </div>
          </div>

          <button
            type="button"
            disabled={totalAmount === 0}
            className="rounded-[18px] bg-[#111827] px-8 py-5 text-2xl font-medium text-white transition hover:bg-black disabled:cursor-not-allowed disabled:opacity-50"
            onClick={handleProceed}
          >
            Proceed to Payment
          </button>
        </section>
      </div>

      <Modal
        open={checkoutOpen}
        onClose={() => {
          if (!isSubmitting) {
            setCheckoutOpen(false)
            setError(null)
          }
        }}
        className="relative w-[min(92vw,400px)] rounded-[28px] bg-white px-8 pb-8 pt-6 shadow-[0_28px_90px_rgba(15,23,42,0.22)]"
      >
        <button
          type="button"
          aria-label="Close checkout"
          className="absolute right-5 top-5 rounded-full p-1 text-[#111111] transition hover:bg-[#f5f5f5]"
          onClick={() => setCheckoutOpen(false)}
        >
          <X className="h-7 w-7" />
        </button>

        <div className="flex flex-col items-center gap-6">
          <h2 className="text-[2rem] font-medium text-[#111111]">Checkout</h2>

          <div className="w-full rounded-[22px] border border-[#dfdfdf] bg-white px-4 py-5 shadow-[0_10px_30px_rgba(17,24,39,0.08)]">
            <h3 className="text-center text-[1.8rem] font-medium text-[#111111]">Order Summary</h3>

            <div className="mt-4 flex flex-col items-center gap-4">
              <img
                src={displayEvent.coverImageUrl}
                alt={displayEvent.title}
                className="h-32 w-28 rounded-[18px] object-cover shadow-sm"
              />
              <div className="text-center text-[1.95rem] font-medium leading-tight text-[#111111]">
                {displayEvent.title}
              </div>
            </div>

            <div className="mt-5 flex flex-col gap-3 text-[1.35rem]">
              {selectedTickets.map((ticket) => (
                <div key={ticket.name} className="flex items-center justify-between gap-4 text-[#8d949e]">
                  <span>{ticket.quantity}x {ticket.name}</span>
                  <span>{formatCurrency(ticket.amount)}</span>
                </div>
              ))}
            </div>

            <div className="my-4 border-t border-[#d9d9d9]" />

            <div className="flex flex-col gap-2 text-[1.35rem] text-[#8d949e]">
              <div className="flex items-center justify-between gap-4">
                <span>Order Amount</span>
                <span>{formatCurrency(totalAmount)}</span>
              </div>
              <div className="flex items-center justify-between gap-4">
                <span>Platform fee (5%)</span>
                <span>{formatCurrency(platformFee)}</span>
              </div>
            </div>

            <div className="my-4 border-t border-[#d9d9d9]" />

            <div className="flex items-center justify-between gap-4 text-[1.7rem] font-semibold">
              <span className="text-[#111111]">Total Amount</span>
              <span className="text-[#ff445d]">{formatCurrency(finalAmount)}</span>
            </div>
          </div>

          {error ? (
            <div className="w-full rounded-2xl border border-[#ffc8cf] bg-[#fff5f6] px-4 py-3 text-sm text-[#cc334a]">
              {error}
            </div>
          ) : null}

          <button
            type="button"
            className="w-full rounded-[18px] bg-[#090c44] px-6 py-4 text-[1.7rem] font-medium text-white transition hover:bg-[#06082f] disabled:cursor-not-allowed disabled:opacity-70"
            onClick={handlePayment}
            disabled={isSubmitting}
          >
            {isSubmitting ? "Processing..." : "Pay using Razorpay"}
          </button>
        </div>
      </Modal>
    </div>
  )
}
