import { ChevronLeft, ChevronRight, ChevronDown, Mail, MapPin } from "lucide-react"
import { useEffect, useState } from "react"
import Modal from "../../../shared/ui/Modal"
import { formatDateBadge, formatEventTime, type BookingRecord } from "../userDashboardData"
import type { TicketRecord } from "../userDashboardData"
import { useAuthStore } from "../../auth/store/authStore"

type Props = {
  booking: BookingRecord | null
  tickets: TicketRecord[]
  open: boolean
  onClose: () => void
}

function TicketQr({ value }: { value: string }) {
  return (
    <div className="rounded-lg bg-white p-2">
      <img
        src={`https://api.qrserver.com/v1/create-qr-code/?size=120x120&data=${encodeURIComponent(value)}`}
        alt="Ticket QR Code"
        className="h-[120px] w-[120px]"
      />
    </div>
  )
}



export default function TicketModal({ booking, tickets, open, onClose }: Props) {
  const [activeIndex, setActiveIndex] = useState(0)
  const attendeeName = useAuthStore.getState().user?.name || "Attendee"

  useEffect(() => {
    if (open) setActiveIndex(0)
  }, [open, booking?.booking_id])

  const activeTicket = tickets[activeIndex]

  if (!booking) return null

  const shareTitle = encodeURIComponent(`I'm going to ${booking.event_title}`)
  const shareLink = `https://wa.me/?text=${shareTitle}`
  const emailLink = `mailto:?subject=${shareTitle}&body=${shareTitle}`

  const coverSrc = booking.coverImageUrl?.startsWith("/")
    ? `${import.meta.env.VITE_API_BASE_URL}${booking.coverImageUrl}`
    : booking.coverImageUrl

  const handlePrint = () => {
    const printWindow = document.createElement("iframe")
    printWindow.style.position = "absolute"
    printWindow.style.top = "-1000px"
    printWindow.style.left = "-1000px"
    document.body.appendChild(printWindow)

    const doc = printWindow.contentDocument || printWindow.contentWindow?.document
    if (!doc) return

    const ticketsHtml = tickets
      .map((t) => {
        const ticketCover = booking.coverImageUrl?.startsWith("/")
          ? `${import.meta.env.VITE_API_BASE_URL}${booking.coverImageUrl}`
          : booking.coverImageUrl

        return `
        <div class="ticket-print-page" style="break-after: page; page-break-after: always; width: 100%; min-height: 100vh; display: flex; justify-content: center; align-items: flex-start; padding: 20mm; box-sizing: border-box; background: white;">
          <div style="width: 600px; border-radius: 24px; overflow: hidden; box-shadow: 0 8px 40px rgba(0,0,0,0.18); background: #111821; color: #fff; font-family: system-ui, -apple-system, sans-serif;">
            ${
              ticketCover
                ? `<img src="${ticketCover}" style="width: 100%; height: 200px; object-fit: cover; display: block;" />`
                : `<div style="width: 100%; height: 200px; background: linear-gradient(135deg,#4b5563,#1f2937); display: flex; alignItems: center; justifyContent: center;"><span style="fontSize: 64px; fontWeight: 700; color: rgba(255,255,255,0.2);">${(booking.event_title || "E")[0].toUpperCase()}</span></div>`
            }
            <div style="padding: 24px 32px 20px;">
              <h2 style="margin: 0; fontSize: 22px; fontWeight: 700;">${booking.event_title}</h2>
              <p style="margin: 6px 0 0; fontSize: 14px; color: #a0aec0;">${booking.venue || booking.event_city}</p>
            </div>
            <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px 32px; padding: 0 32px 24px; color: #fff;">
              <div>
                <p style="margin: 0; fontSize: 11px; text-transform: uppercase; letter-spacing: 0.08em; color: #8d949e;">Date</p>
                <p style="margin: 4px 0 0; fontSize: 17px;">${formatDateBadge(booking.event_start_time)}</p>
              </div>
              <div>
                <p style="margin: 0; fontSize: 11px; text-transform: uppercase; letter-spacing: 0.08em; color: #8d949e;">Time</p>
                <p style="margin: 4px 0 0; fontSize: 17px;">${formatEventTime(booking.event_start_time)}</p>
              </div>
              <div>
                <p style="margin: 0; fontSize: 11px; text-transform: uppercase; letter-spacing: 0.08em; color: #8d949e;">Ticket Type</p>
                <p style="margin: 4px 0 0; fontSize: 17px; fontWeight: 600;">${t.ticket_type ?? "General"}</p>
              </div>
              <div>
                <p style="margin: 0; fontSize: 11px; text-transform: uppercase; letter-spacing: 0.08em; color: #8d949e;">Ticket ID</p>
                <p style="margin: 4px 0 0; fontSize: 17px;">${t.ticket_code ?? booking.booking_id.slice(0, 4)}</p>
              </div>
            </div>
            <div style="border-top: 2px dashed rgba(255,255,255,0.2); margin: 0 32px; position: relative;">
              <div style="position: absolute; left: -44px; top: 50%; transform: translateY(-50%); width: 32px; height: 32px; border-radius: 50%; background: #fff;"></div>
              <div style="position: absolute; right: -44px; top: 50%; transform: translateY(-50%); width: 32px; height: 32px; border-radius: 50%; background: #fff;"></div>
            </div>
            <div style="display: flex; align-items: flex-end; justify-content: space-between; padding: 20px 32px 28px;">
              <div>
                <p style="margin: 0; fontSize: 17px; fontWeight: 600; color: #fff;">${attendeeName}</p>
                <p style="margin: 4px 0 0; fontSize: 12px; color: #7b838e;">Booking ID: ${booking.booking_id}</p>
              </div>
              <div style="border-radius: 8px; background: white; padding: 8px;">
                <img src="https://api.qrserver.com/v1/create-qr-code/?size=120x120&data=${encodeURIComponent(`${booking.booking_id}-${t.ticket_id}`)}" style="width: 120px; height: 120px;" />
              </div>
            </div>
          </div>
        </div>`
      })
      .join("")

    doc.open()
    doc.write(`
      <html>
        <head>
          <title>Tickets - ${booking.event_title}</title>
          <style>
            @page { size: A4 portrait; margin: 0; }
            body { margin: 0; padding: 0; }
            * { box-sizing: border-box; }
            .ticket-print-page:last-child { break-after: avoid; page-break-after: avoid; }
          </style>
        </head>
        <body>
          ${ticketsHtml}
          <script>
            window.onload = function() {
              window.print();
              setTimeout(function() { window.frameElement.remove(); }, 100);
            };
          </script>
        </body>
      </html>
    `)
    doc.close()
  }

  return (
    <Modal
      open={open}
      onClose={onClose}
      className="relative w-[min(90vw,520px)] rounded-[24px] bg-white px-6 pb-8 pt-5 shadow-[0_24px_80px_rgba(15,23,42,0.24)]"
    >
      <div className="mx-auto flex max-w-[470px] flex-col items-center gap-6">
        <div className="relative w-full print:hidden">
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
            {coverSrc ? (
              <img
                src={coverSrc}
                alt={booking.event_title}
                className="h-[200px] w-full object-cover"
              />
            ) : (
              <div className="flex h-[200px] w-full items-center justify-center bg-gradient-to-br from-gray-700 to-gray-900">
                <span className="text-4xl font-bold text-white opacity-30">
                  {(booking.event_title || "E")[0].toUpperCase()}
                </span>
              </div>
            )}

            <div className="space-y-5 px-6 py-6">
              <div>
                <h3 className="text-xl font-semibold leading-tight">{booking.event_title}</h3>
                <div className="mt-2 flex items-center gap-2 text-base text-[#f0f2f5]">
                  <MapPin className="h-5 w-5 text-[#ff4d5f]" />
                  <span>{booking.venue || booking.event_city}</span>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-x-8 gap-y-5 text-left">
                <div>
                  <p className="text-xs uppercase tracking-[0.08em] text-[#8d949e]">Date</p>
                  <p className="mt-1 text-lg">{formatDateBadge(booking.event_start_time)}</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-[0.08em] text-[#8d949e]">Time</p>
                  <p className="mt-1 text-lg">{formatEventTime(booking.event_start_time)}</p>
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
                <p className="text-lg text-white">{attendeeName}</p>
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
          onClick={handlePrint}
          className="w-full rounded-[16px] bg-[#ef3650] px-4 py-3 text-lg font-semibold text-white transition hover:bg-[#d92f47] print:hidden"
        >
          Download Tickets
        </button>

        <div className="flex items-center gap-3 text-base text-[#111111] print:hidden">
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
