import { Minus, Plus, X } from "lucide-react"
import { useState, useEffect } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { useQueryClient } from "@tanstack/react-query"
import Modal from "../../../shared/ui/Modal"
import { useEvent } from "../hooks"
import { buildDisplayEvent, formatCurrency } from "../eventBookingData"
import { eventsApi } from "../api"
import { useAuthStore } from "../../auth/store/authStore"
import RazorpayButton from "../../payments/components/RazorpayButton"
import { useWallet, usePaymentSettings, walletQueryKey, walletTransactionsQueryKey } from "../../user/hooks"
import { userApi } from "../../user/api"
import { useEngagement } from "../../../shared/hooks/useEngagement";

export default function EventBookingPage() {
  const queryClient = useQueryClient()
  const { eventId } = useParams()
  const navigate = useNavigate()
  const { data, isLoading } = useEvent(eventId!)
  const { user, roles } = useAuthStore()
  const { trackEvent } = useEngagement();

  const displayEvent = buildDisplayEvent(eventId ?? "", data)

  useEffect(() => {
    if (user && roles && roles.some(r => r.toLowerCase() === "organizer" || r.toLowerCase() === "admin")) {
      navigate(`/events/${eventId}`, { replace: true })
    }
  }, [user, roles, navigate, eventId])
  const [quantities, setQuantities] = useState<Record<string, number>>(() =>
    Object.fromEntries(displayEvent.ticketTypes.map((ticket) => [ticket.name, 0])),
  )
  const [checkoutOpen, setCheckoutOpen] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [reservedBookingId, setReservedBookingId] = useState<string | null>(null)
  const { data: wallet } = useWallet()
  const { data: paymentSettings } = usePaymentSettings()
  const [isPayingWithWallet, setIsPayingWithWallet] = useState(false)
  const [paymentSuccess, setPaymentSuccess] = useState(false)
  const [isLatePaymentMessage, setIsLatePaymentMessage] = useState(false)

  const razorpaySetting = paymentSettings?.find((p) => p.provider === "razorpay")
  const isRazorpayEnabled = razorpaySetting ? razorpaySetting.is_enabled : false

  const walletSetting = paymentSettings?.find((p) => p.provider === "wallet")
  const isWalletEnabled = walletSetting ? walletSetting.is_enabled : false

  const ticketRows = displayEvent.ticketTypes.map((ticket) => ({
    ...ticket,
    quantity: quantities[ticket.name] ?? 0,
    amount: ticket.price * (quantities[ticket.name] ?? 0),
  }))

  const totalAmount = ticketRows.reduce((sum, ticket) => sum + ticket.amount, 0)
  const totalTickets = ticketRows.reduce((sum, ticket) => sum + ticket.quantity, 0)
  const platformFee = totalAmount > 0 ? 30 * totalTickets : 0
  const finalAmount = totalAmount + platformFee

  const selectedTickets = ticketRows.filter((ticket) => ticket.quantity > 0)

  const updateQuantity = (ticketName: string, nextQuantity: number, limit?: number) => {
    if (reservedBookingId) return
    const safeQuantity = Math.max(0, Math.min(limit ?? Number.POSITIVE_INFINITY, nextQuantity))
    setQuantities((current) => ({
      ...current,
      [ticketName]: safeQuantity,
    }))
  }

  const handleProceed = () => {
    if (selectedTickets.length === 0) return
    setError(null)
    setCheckoutOpen(true)
    
    if (data?.id) {
      trackEvent('ticket_selected', data.id);
      trackEvent('checkout_started', data.id);
    }
  }

  const handleReservation = async () => {
    if (!eventId || selectedTickets.length === 0) return

    setIsSubmitting(true)
    setError(null)

    try {
      const hasTicketIds = selectedTickets.every((ticket) => ticket.id)
      if (!hasTicketIds) {
        throw new Error("Ticket information is incomplete. Please refresh and try again.")
      }

      const response = await eventsApi.reserveTickets({
        eventId: displayEvent.id,
        tickets: selectedTickets.map((ticket) => ({
          ticket_type_id: ticket.id!,
          quantity: ticket.quantity,
        })),
      })

      setReservedBookingId(response.booking_id)
    } catch (reservationError: any) {
      setError(
        reservationError?.response?.data?.message ??
          reservationError?.message ??
          "We could not reserve tickets right now. Please try again.",
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  const handlePayWithWallet = async () => {
    if (!reservedBookingId) return
    setIsPayingWithWallet(true)
    setError(null)
    try {
      await userApi.payWithWallet(reservedBookingId)
      await queryClient.invalidateQueries({ queryKey: walletQueryKey });
      await queryClient.invalidateQueries({ queryKey: walletTransactionsQueryKey });
      setPaymentSuccess(true);
      setTimeout(() => {
        navigate("/profile/bookings", { replace: true })
      }, 2000)
    } catch (err: any) {
      setError(err?.response?.data?.message || "Wallet payment failed.")
    } finally {
      setIsPayingWithWallet(false)
    }
  }

  return (
    <div className="bg-white">
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-6 px-6 py-8">
        <section className="overflow-hidden rounded-2xl border border-[#e8e8e8] bg-white shadow-sm">
          <div className="grid md:grid-cols-[1.1fr_1fr]">
            <div className="min-h-[180px] bg-[#111827]">
              <img
                src={displayEvent.coverImageUrl}
                alt={displayEvent.title}
                className="h-full w-full object-cover"
              />
            </div>

            <div className="flex flex-col justify-center gap-2 p-6">
              <h1 className="max-w-md text-lg font-semibold tracking-tight text-[#111111]">
                {displayEvent.title}
              </h1>
              <p className="text-sm text-[#8d949e]">{displayEvent.dateLabel}</p>
              <p className="text-sm text-[#8d949e]">{displayEvent.timeLabel}</p>
              <p className="text-sm text-[#8d949e]">{displayEvent.displayLocation}</p>
            </div>
          </div>
        </section>

        <section className="mx-auto flex w-full max-w-3xl flex-col gap-4">
          <div className="text-center">
            <h2 className="text-xl font-semibold tracking-tight text-[#111111]">Ticket Selection</h2>
            {isLoading ? <p className="mt-1 text-xs text-[#8d949e]">Loading event details...</p> : null}
          </div>

          <div className="flex flex-col gap-3">
            {ticketRows.map((ticket) => {
              const isSoldOut = ticket.availableQuantity === 0;
              return (
              <div
                key={ticket.name}
                className={`grid items-center gap-3 rounded-xl border px-4 py-3 text-[#111111] md:grid-cols-[1.4fr_1fr_auto] ${
                  isSoldOut ? "border-gray-200 bg-gray-50 opacity-80" : "border-[#ffc8cf] bg-white"
                }`}
              >
                <div className="text-sm font-medium flex items-center gap-2">
                  {ticket.name}
                  {isSoldOut && (
                     <span className="text-[10px] font-bold uppercase tracking-wider text-white bg-gray-400 px-2 py-0.5 rounded-md">Sold Out</span>
                  )}
                </div>
                <div className={`text-left text-sm font-semibold md:text-center ${isSoldOut ? "text-gray-500" : "text-[#ff445d]"}`}>
                  Price: {formatCurrency(ticket.price)}
                </div>
                <div className="ml-auto flex items-center gap-3 text-sm">
                  {isSoldOut ? (
                    <span className="text-gray-400 font-medium text-sm pr-2">Unavailable</span>
                  ) : (
                    <>
                      <span className="text-[#111111]">Qty</span>
                      <button
                        type="button"
                        aria-label={`Decrease ${ticket.name} tickets`}
                        className="rounded-full p-1 text-[#111111] transition hover:bg-[#f7f7f7]"
                        onClick={() => updateQuantity(ticket.name, ticket.quantity - 1, ticket.availableQuantity)}
                      >
                        <Minus className="h-4 w-4" />
                      </button>
                      <span className="inline-flex min-w-8 items-center justify-center rounded-md bg-[#ffd9df] px-2 py-0.5 text-sm font-medium text-[#111111]">
                        {ticket.quantity}
                      </span>
                      <button
                        type="button"
                        aria-label={`Increase ${ticket.name} tickets`}
                        className="rounded-full p-1 text-[#111111] transition hover:bg-[#f7f7f7]"
                        onClick={() => updateQuantity(ticket.name, ticket.quantity + 1, ticket.availableQuantity)}
                      >
                        <Plus className="h-4 w-4" />
                      </button>
                    </>
                  )}
                </div>
              </div>
            )})}
          </div>

          <div className="flex justify-end">
            <div className="text-right text-base font-semibold text-[#111111]">
              Total Amount: {formatCurrency(totalAmount)}
            </div>
          </div>

          <button
            type="button"
            disabled={selectedTickets.length === 0 || isSubmitting}
            className="rounded-xl bg-[#111827] px-6 py-2.5 text-sm font-medium text-white transition hover:bg-black disabled:cursor-not-allowed disabled:opacity-50"
            onClick={handleProceed}
          >
            {isSubmitting ? "Processing..." : "Reserve Tickets"}
          </button>
        </section>
      </div>

      <Modal
        open={checkoutOpen}
        onClose={() => {
          if (!isSubmitting && !reservedBookingId) {
            setCheckoutOpen(false)
            setError(null)
          }
        }}
        className="relative w-[min(92vw,380px)] rounded-2xl bg-white px-5 pb-5 pt-4 shadow-[0_28px_90px_rgba(15,23,42,0.22)]"
      >
        <button
          type="button"
          aria-label="Close checkout"
          className="absolute right-4 top-4 rounded-full p-1 text-[#111111] transition hover:bg-[#f5f5f5]"
          onClick={() => setCheckoutOpen(false)}
        >
          <X className="h-5 w-5" />
        </button>

        <div className="flex flex-col items-center gap-4">
          <h2 className="text-base font-medium text-[#111111]">Reserve Tickets</h2>

          <div className="w-full rounded-xl border border-[#dfdfdf] bg-white px-4 py-3 shadow-sm">
            <h3 className="text-center text-sm font-medium text-[#111111]">Order Summary</h3>

            <div className="mt-3 flex flex-col items-center gap-3">
              <img
                src={displayEvent.coverImageUrl}
                alt={displayEvent.title}
                className="h-24 w-20 rounded-xl object-cover shadow-sm"
              />
              <div className="text-center text-sm font-medium leading-tight text-[#111111]">
                {displayEvent.title}
              </div>
            </div>

            <div className="mt-3 flex flex-col gap-2 text-sm">
              {selectedTickets.map((ticket) => (
                <div key={ticket.name} className="flex items-center justify-between gap-3 text-[#8d949e]">
                  <span>{ticket.quantity}x {ticket.name}</span>
                  <span>{formatCurrency(ticket.amount)}</span>
                </div>
              ))}
            </div>

            <div className="my-3 border-t border-[#d9d9d9]" />

            <div className="flex flex-col gap-1.5 text-sm text-[#8d949e]">
              <div className="flex items-center justify-between gap-3">
                <span>Order Amount</span>
                <span>{formatCurrency(totalAmount)}</span>
              </div>
              <div className="flex items-center justify-between gap-3 text-sm text-[#5d6573]">
                <span>Platform fee (₹30 per ticket)</span>
                <span>{formatCurrency(platformFee)}</span>
              </div>
            </div>

            <div className="my-3 border-t border-[#d9d9d9]" />

            <div className="flex items-center justify-between gap-3 text-base font-semibold">
              <span className="text-[#111111]">Total Amount</span>
              <span className="text-[#ff445d]">{formatCurrency(finalAmount)}</span>
            </div>
          </div>

          {error ? (
            <div className="w-full rounded-2xl border border-[#ffc8cf] bg-[#fff5f6] px-4 py-3 text-sm text-[#cc334a]">
              {error}
            </div>
          ) : null}

          {reservedBookingId ? (
            <div className="flex w-full flex-col gap-2">
              {isRazorpayEnabled ? (
                <RazorpayButton
                  bookingId={reservedBookingId}
                  eventTitle={displayEvent.title}
                  onSuccess={(isLatePayment) => {
                    queryClient.invalidateQueries({ queryKey: walletQueryKey });
                    queryClient.invalidateQueries({ queryKey: walletTransactionsQueryKey });
                    if (isLatePayment) {
                      setIsLatePaymentMessage(true);
                    } else {
                      setPaymentSuccess(true);
                      setTimeout(() => {
                        navigate("/profile/bookings", { replace: true })
                      }, 2000)
                    }
                  }}
                  onError={(msg) => setError(msg)}
                />
              ) : (
                <button
                  type="button"
                  disabled
                  className="flex w-full items-center justify-center rounded-xl bg-gray-200 px-6 py-2.5 text-[15px] font-semibold text-gray-400 cursor-not-allowed"
                >
                  Pay with Razorpay
                </button>
              )}
              {isWalletEnabled && wallet && wallet.available_balance >= finalAmount && (
                <button
                  type="button"
                  disabled={isPayingWithWallet}
                  onClick={handlePayWithWallet}
                  className="flex w-full items-center justify-center rounded-xl border border-[#111827] bg-white px-6 py-2.5 text-[15px] font-semibold text-[#111827] transition hover:bg-[#f7f7f7] disabled:opacity-50"
                >
                  {isPayingWithWallet ? "Processing Wallet Payment..." : "Pay with Wallet"}
                </button>
              )}
              {(!isWalletEnabled) && wallet && wallet.available_balance >= finalAmount && (
                 <button
                  type="button"
                  disabled
                  className="flex w-full items-center justify-center rounded-xl border border-gray-200 bg-gray-50 px-6 py-2.5 text-[15px] font-semibold text-gray-400 cursor-not-allowed"
                 >
                   Pay with Wallet Disabled
                 </button>
              )}
              {wallet && wallet.available_balance < finalAmount && wallet.available_balance > 0 && (
                <p className="text-center text-[10px] text-gray-400">
                  Insufficient wallet balance ({formatCurrency(wallet.available_balance)})
                </p>
              )}
            </div>
          ) : (
            <button
              type="button"
              className="w-full rounded-xl bg-[#0b101e] px-5 py-2.5 text-sm font-medium text-white transition hover:bg-black disabled:cursor-not-allowed disabled:opacity-70"
              onClick={handleReservation}
              disabled={isSubmitting}
            >
              {isSubmitting ? "Processing..." : "Reserve Booking"}
            </button>
          )}

          {paymentSuccess && (
            <div className="absolute inset-0 flex flex-col items-center justify-center rounded-2xl bg-white/95 text-center backdrop-blur-sm z-50">
              <div className="mb-4 text-[#34c759]">
                 {}
                <svg className="h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <h2 className="text-xl font-bold text-[#111111]">Payment Successful!</h2>
              <p className="mt-2 text-sm text-[#8d949e]">Your tickets have been confirmed.</p>
              <p className="mt-1 text-xs text-[#8d949e]">Redirecting...</p>
            </div>
          )}

          {isLatePaymentMessage && (
            <div className="absolute inset-0 flex flex-col items-center justify-center rounded-2xl bg-[#fff5f6] text-center p-6 z-50">
              <div className="mb-4 text-[#e53e5d]">
                 {}
                <svg className="h-14 w-14" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
              </div>
              <h2 className="text-lg font-bold text-[#111111] mb-2">Payment Received</h2>
              <p className="text-sm text-gray-700 leading-relaxed">
                Your payment was successful, but the booking expiration time had already passed. 
              </p>
              <p className="text-sm text-gray-700 leading-relaxed mt-2 font-medium">
                The amount will be refunded to you within 3-5 working days. Please ensure your payout details are updated in your Profile.
              </p>
              <button 
                 onClick={() => navigate("/user/profile")}
                 className="mt-6 bg-[#0b101e] text-white px-5 py-2 rounded-xl text-sm font-medium hover:bg-black transition-colors"
              >
                Go to Profile
              </button>
            </div>
          )}
        </div>
      </Modal>
    </div>
  )
}
