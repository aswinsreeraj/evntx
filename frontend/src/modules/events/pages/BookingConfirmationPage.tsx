import { Check } from "lucide-react"
import { Link, Navigate, useParams } from "react-router-dom"
import { getBookingConfirmation } from "../bookingStorage"
import { formatCurrency } from "../eventBookingData"

export default function BookingConfirmationPage() {
  const { eventId } = useParams()
  const booking = getBookingConfirmation()

  if (!booking || booking.eventId !== eventId) {
    return <Navigate to={eventId ? `/events/${eventId}/book` : "/events"} replace />
  }

  const bookingDateTime = new Date(booking.createdAt).toLocaleString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  })

  return (
    <div className="bg-white">
      <div className="mx-auto flex w-full max-w-4xl flex-col items-center gap-6 px-6 py-8">
        <div className="flex h-14 w-14 items-center justify-center rounded-full bg-[#19bcc2]">
          <Check className="h-7 w-7 text-white" strokeWidth={4} />
        </div>

        <div className="text-center">
          <h1 className="text-xl font-semibold tracking-tight text-[#111111]">Payment Successful</h1>
        </div>

        <section className="grid w-full gap-5 rounded-2xl bg-[#f7f7f7] px-6 py-6 text-[#111111] shadow-sm md:grid-cols-2">
          <div className="space-y-4">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.08em] text-[#8d949e]">Booking Reference</p>
              <p className="mt-1 text-base">{booking.bookingId}</p>
            </div>
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.08em] text-[#8d949e]">Date & Time</p>
              <p className="mt-1 text-base">{bookingDateTime}</p>
            </div>
          </div>

          <div className="space-y-4">
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.08em] text-[#8d949e]">Event Name</p>
              <p className="mt-1 text-base leading-tight">{booking.eventTitle}</p>
            </div>
            <div>
              <p className="text-xs font-semibold uppercase tracking-[0.08em] text-[#8d949e]">Total Paid</p>
              <p className="mt-1 text-base font-semibold text-[#ff445d]">{formatCurrency(booking.finalAmount)}</p>
            </div>
          </div>
        </section>

        <div className="w-full rounded-xl border border-[#ffc8cf] bg-[#fff7f8] px-4 py-3 text-sm text-[#333333]">
          Tickets and receipt have been sent to <span className="font-semibold">{booking.email}</span>
        </div>

        <div className="grid w-full gap-3 md:grid-cols-2">
          <Link
            to="/profile"
            className="rounded-xl bg-[#111827] px-5 py-2.5 text-center text-sm font-medium text-white transition hover:bg-black"
          >
            View Dashboard
          </Link>
          <button
            type="button"
            className="rounded-xl border border-[#111827] bg-white px-5 py-2.5 text-sm font-medium text-[#111827] transition hover:bg-[#f7f7f7]"
            onClick={() => window.print()}
          >
            Download Ticket
          </button>
        </div>
      </div>
    </div>
  )
}
