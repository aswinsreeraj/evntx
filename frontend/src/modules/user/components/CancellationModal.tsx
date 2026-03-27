import { AlertTriangle, Info, Minus, Plus, X } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import Modal from "../../../shared/ui/Modal"
import type { BookingRecord, TicketRecord } from "../userDashboardData"
import { formatCurrency } from "../userDashboardData"

type CancellationSelection = {
  ticketType: string
  ownedCount: number
  cancelCount: number
  price: number
}

type Props = {
  booking: BookingRecord | null
  tickets: TicketRecord[]
  open: boolean
  onClose: () => void
  onConfirm: (selection: Array<{ ticketType: string; cancelCount: number; refundAmount: number }>) => void
}

export default function CancellationModal({ booking, tickets, open, onClose, onConfirm }: Props) {
  const groupedOptions = useMemo(() => {
    if (!booking) return []

    const grouped = new Map<string, number>()
    tickets.forEach((ticket) => {
      grouped.set(ticket.ticket_type, (grouped.get(ticket.ticket_type) ?? 0) + 1)
    })

    return Array.from(grouped.entries()).map(([ticketType, count]) => {
      // Platform fee is 5%, so base price = total / 1.05
      const baseTotal = booking.total_amount / 1.05
      const price = (baseTotal / tickets.length) || 0

      return {
        ticketType,
        ownedCount: count,
        cancelCount: 0,
        price,
      }
    })
  }, [booking, tickets])

  const [selection, setSelection] = useState<CancellationSelection[]>([])

  useEffect(() => {
    if (open) {
      setSelection(groupedOptions)
    }
  }, [groupedOptions, open])

  const totalRefund = selection.reduce((sum, item) => sum + item.cancelCount * item.price, 0)

  const updateCount = (ticketType: string, nextCount: number) => {
    setSelection((current) =>
      current.map((item) =>
        item.ticketType === ticketType
          ? { ...item, cancelCount: Math.max(0, Math.min(item.ownedCount, nextCount)) }
          : item,
      ),
    )
  }

  return (
    <Modal
      open={open}
      onClose={onClose}
      className="relative w-[min(94vw,480px)] rounded-[20px] bg-white px-6 pb-6 pt-5 shadow-[0_24px_80px_rgba(15,23,42,0.24)]"
    >
      <button
        type="button"
        onClick={onClose}
        className="absolute right-4 top-4 rounded-full p-1 text-[#e6e6e6] transition hover:bg-[#f7f7f7]"
      >
        <X className="h-5 w-5" />
      </button>

      <div className="flex flex-col items-center">
        <div className="mt-2 flex h-16 w-16 items-center justify-center rounded-full">
          <AlertTriangle className="h-14 w-14 fill-[#ff1717] text-[#ff1717]" strokeWidth={1.6} />
        </div>

        <h2 className="mt-2 text-xl font-semibold text-[#111111]">Cancel Tickets</h2>
        <p className="mt-1 text-center text-sm text-[#5b6069]">
          Select the tickets you want to cancel for <span className="font-semibold">{booking?.event_title}</span>
        </p>

        <div className="mt-4 flex w-full flex-col gap-3">
          {selection.map((item) => (
            <div
              key={item.ticketType}
              className="grid items-center gap-3 rounded-xl border border-[#ffc9cf] px-3 py-3 md:grid-cols-[1fr_auto_auto]"
            >
              <div className="text-sm font-medium text-[#111111]">{item.ticketType}</div>

              <div className="flex items-center gap-2 text-sm text-[#111111]">
                <span>Qty</span>
                <button
                  type="button"
                  onClick={() => updateCount(item.ticketType, item.cancelCount - 1)}
                  className="rounded-md bg-[#f7c6cd] px-1.5 py-0.5 text-[#111111] transition hover:bg-[#f2b4be]"
                >
                  <Minus className="h-3.5 w-3.5" />
                </button>
                <span className="min-w-4 text-center">{item.cancelCount}</span>
                <button
                  type="button"
                  onClick={() => updateCount(item.ticketType, item.cancelCount + 1)}
                  className="rounded-md bg-[#f7c6cd] px-1.5 py-0.5 text-[#111111] transition hover:bg-[#f2b4be]"
                >
                  <Plus className="h-3.5 w-3.5" />
                </button>
              </div>

              <div className="justify-self-end text-sm font-semibold text-[#111111]">
                {formatCurrency(item.cancelCount * item.price)}
              </div>
            </div>
          ))}
        </div>

        <div className="mt-5 flex w-full items-center justify-between text-base">
          <span className="text-[#111111]">Estimated Refund</span>
          <span className="font-semibold text-[#ff1717]">{formatCurrency(totalRefund)}</span>
        </div>

        <div className="mt-4 flex w-full items-center gap-2 rounded-xl border border-[#ffb9c1] bg-[#fff3f5] px-3 py-2.5 text-xs text-[#ef3650]">
          <Info className="h-4 w-4 fill-[#ef3650] text-white flex-shrink-0" />
          <span>Refunds are deposited to the wallet.</span>
        </div>

        <div className="mt-5 flex w-full flex-col gap-2">
          <button
            type="button"
            disabled={selection.every(s => s.cancelCount === 0)}
            onClick={() =>
              onConfirm(
                selection
                  .filter((item) => item.cancelCount > 0)
                  .map((item) => ({
                    ticketType: item.ticketType,
                    cancelCount: item.cancelCount,
                    refundAmount: item.cancelCount * item.price,
                  })),
              )
            }
            className="rounded-xl bg-[#ff1717] px-5 py-2.5 text-sm font-semibold text-white transition hover:bg-[#e41212] disabled:cursor-not-allowed disabled:opacity-60"
          >
            Cancel Selected
          </button>

          <button
            type="button"
            onClick={onClose}
            className="rounded-xl border border-[#111111] px-5 py-2.5 text-sm font-semibold text-[#111111] transition hover:bg-[#f7f7f7]"
          >
            Keep Tickets
          </button>
        </div>
      </div>
    </Modal>
  )
}
