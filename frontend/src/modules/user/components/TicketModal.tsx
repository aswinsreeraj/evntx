import { ChevronLeft, ChevronRight, ChevronDown, Mail, MapPin, X } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import Modal from "../../../shared/ui/Modal"
import type { EnrichedBooking } from "../userDashboardData"
import type { TicketRecord } from "../userDashboardData"
import { useAuthStore } from "../../auth/store/authStore"

type Props = {
  booking: EnrichedBooking | null
  tickets: TicketRecord[]
  open: boolean
  onClose: () => void
}

function TicketQr({ value }: { value: string }) {
  const cells = useMemo(() => {
    const seed = value.split("").reduce((sum, char) => sum + char.charCodeAt(0), 0)
    return Array.from({ length: 15 * 15 }, (_, index) => ((index * 17 + seed) % 5 === 0 ? 1 : 0))
  }, [value])

  return (
    <div
      className="grid gap-[2px] rounded-lg bg-white p-2"
      style={{ gridTemplateColumns: "repeat(15, minmax(0, 1fr))" }}
    >
      {cells.map((cell, index) => (
        <div
          key={index}
          className={`h-[6px] w-[6px] rounded-[1px] ${cell ? "bg-[#111827]" : "bg-transparent"}`}
        />
      ))}
    </div>
  )
}

export default function TicketModal({ booking, tickets, open, onClose }: Props) {
  const [activeIndex, setActiveIndex] = useState(0)

  useEffect(() => {
    if (open) setActiveIndex(0)
  }, [open, booking?.booking_id])

  const activeTicket = tickets[activeIndex]

  if (!booking) return null

  const shareTitle = encodeURIComponent(`I'm going to ${booking.event_title}`)
  const shareLink = `https://wa.me/?text=${shareTitle}`
  const emailLink = `mailto:?subject=${shareTitle}&body=${shareTitle}`

  return (
    <Modal
      open={open}
      onClose={onClose}
      className="relative w-[min(90vw,520px)] rounded-[24px] bg-white px-6 pb-8 pt-5 shadow-[0_24px_80px_rgba(15,23,42,0.24)]"
    >
      <button
        type="button"
        onClick={onClose}
        className="absolute right-5 top-5 rounded-full p-1 text-[#d9d9d9] transition hover:bg-[#f7f7f7]"
      >
        <X className="h-8 w-8" />
      </button>

      <div className="mx-auto flex max-w-[470px] flex-col items-center gap-6">
        <div className="relative w-full">
          <select
            value={activeTicket?.ticket_id ?? ""}
            onChange={(event) => {
              const nextIndex = tickets.findIndex((ticket) => ticket.ticket_id === event.target.value)
              setActiveIndex(nextIndex >= 0 ? nextIndex : 0)
            }}
            className="w-full appearance-none rounded-xl border border-[#d9e1ea] px-4 py-3 text-lg text-[#2a2f36] outline-none"
          >
            {tickets.map((ticket) => (
              <option key={ticket.ticket_id} value={ticket.ticket_id}>
                {ticket.ticket_type} Ticket #{ticket.ticket_code}
              </option>
            ))}
          </select>
          <ChevronDown className="pointer-events-none absolute right-5 top-1/2 h-6 w-6 -translate-y-1/2 text-[#8d949e]" />
        </div>

        <div className="relative flex w-full items-center justify-center">
          <button
            type="button"
            onClick={() => setActiveIndex((current) => Math.max(0, current - 1))}
            disabled={activeIndex === 0}
            className="absolute -left-16 rounded-full p-2 text-[#e3e3e3] transition hover:bg-[#f7f7f7] disabled:opacity-50"
          >
            <ChevronLeft className="h-12 w-12" />
          </button>

          <div className="relative w-full overflow-hidden rounded-[26px] bg-[#111821] text-white">
            <img
              src={booking.coverImageUrl}
              alt={booking.event_title}
              className="h-[200px] w-full object-cover"
            />

            <div className="space-y-5 px-6 py-6">
              <div>
                <h3 className="text-xl font-semibold leading-tight">{booking.event_title}</h3>
                <div className="mt-2 flex items-center gap-2 text-base text-[#f0f2f5]">
                  <MapPin className="h-5 w-5 text-[#ff4d5f]" />
                  <span>{booking.venue}</span>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-x-8 gap-y-5 text-left">
                <div>
                  <p className="text-xs uppercase tracking-[0.08em] text-[#8d949e]">Date</p>
                  <p className="mt-1 text-lg">{booking.dateBadge}</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-[0.08em] text-[#8d949e]">Time</p>
                  <p className="mt-1 text-lg">{booking.timeLabel}</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-[0.08em] text-[#8d949e]">Ticket Type</p>
                  <p className="mt-1 text-lg font-semibold">{activeTicket?.ticket_type ?? "General"}</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-[0.08em] text-[#8d949e]">Ticket ID</p>
                  <p className="mt-1 text-lg">{activeTicket?.ticket_code ?? booking.booking_id.slice(0, 4)}</p>
                </div>
              </div>
            </div>

            <div className="relative mt-2 border-t border-dotted border-white/90">
              <div className="absolute -left-6 top-1/2 h-12 w-12 -translate-y-1/2 rounded-full bg-white" />
              <div className="absolute -right-6 top-1/2 h-12 w-12 -translate-y-1/2 rounded-full bg-white" />
            </div>

            <div className="flex items-end justify-between px-6 py-6">
              <div>
                <p className="text-lg text-white">{useAuthStore.getState().user?.name || "Attendee"}</p>
                <p className="mt-0.5 text-sm text-[#7b838e]">Booking ID: {booking.booking_id}</p>
              </div>
              <TicketQr value={`${booking.booking_id}-${activeTicket?.ticket_id ?? "ticket"}`} />
            </div>
          </div>

          <button
            type="button"
            onClick={() => setActiveIndex((current) => Math.min(tickets.length - 1, current + 1))}
            disabled={activeIndex === tickets.length - 1}
            className="absolute -right-16 rounded-full p-2 text-[#e3e3e3] transition hover:bg-[#f7f7f7] disabled:opacity-50"
          >
            <ChevronRight className="h-12 w-12" />
          </button>
        </div>

        <button
          type="button"
          onClick={() => window.print()}
          className="w-full rounded-[16px] bg-[#ef3650] px-4 py-3 text-lg font-semibold text-white transition hover:bg-[#d92f47]"
        >
          Download Tickets
        </button>

        <div className="flex items-center gap-3 text-base text-[#111111]">
          <span>Share with Friends</span>
          <a
            href={shareLink}
            target="_blank"
            rel="noreferrer"
            className="flex h-10 w-10 items-center justify-center rounded-full bg-[#4caf50] text-white"
          >
            W
          </a>
          <a
            href={emailLink}
            className="flex h-10 w-10 items-center justify-center rounded-full bg-[#2795f3] text-white"
          >
            <Mail className="h-5 w-5" />
          </a>
        </div>
      </div>
    </Modal>
  )
}
