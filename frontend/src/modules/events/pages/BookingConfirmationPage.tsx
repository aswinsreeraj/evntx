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
      <div className="mx-auto flex w-full max-w-5xl flex-col items-center gap-10 px-6 py-12">
        <div className="flex h-28 w-28 items-center justify-center rounded-full bg-[#19bcc2]">
          <Check className="h-16 w-16 text-white" strokeWidth={4} />
        </div>

        <div className="text-center">
          <h1 className="text-4xl font-semibold tracking-tight text-[#111111]">Payment Successful</h1>
        </div>

        <section className="grid w-full gap-8 rounded-[24px] bg-[#f7f7f7] px-8 py-8 text-[#111111] shadow-[0_10px_30px_rgba(15,23,42,0.08)] md:grid-cols-2">
          <div className="space-y-5">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.08em] text-[#8d949e]">Booking Reference</p>
              <p className="mt-1 text-[2rem]">{booking.bookingId}</p>
            </div>
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.08em] text-[#8d949e]">Date & Time</p>
              <p className="mt-1 text-[2rem]">{bookingDateTime}</p>
            </div>
          </div>

          <div className="space-y-5">
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.08em] text-[#8d949e]">Event Name</p>
              <p className="mt-1 text-[2rem] leading-tight">{booking.eventTitle}</p>
            </div>
            <div>
              <p className="text-sm font-semibold uppercase tracking-[0.08em] text-[#8d949e]">Total Paid</p>
              <p className="mt-1 text-[2rem] font-semibold text-[#ff445d]">{formatCurrency(booking.finalAmount)}</p>
            </div>
          </div>
        </section>

        <div className="w-full rounded-[18px] border border-[#ffc8cf] bg-[#fff7f8] px-5 py-4 text-lg text-[#333333]">
          Tickets and receipt have been sent to <span className="font-semibold">{booking.email}</span>
        </div>

        <div className="grid w-full gap-4 md:grid-cols-2">
          <Link
            to="/profile"
            className="rounded-[18px] bg-[#111827] px-6 py-4 text-center text-[1.6rem] font-medium text-white transition hover:bg-black"
          >
            View Dashboard
          </Link>
          <button
            type="button"
            className="rounded-[18px] border border-[#111827] bg-white px-6 py-4 text-[1.6rem] font-medium text-[#111827] transition hover:bg-[#f7f7f7]"
            onClick={() => window.print()}
          >
            Download Ticket
          </button>
        </div>
      </div>
    </div>
  )
}
