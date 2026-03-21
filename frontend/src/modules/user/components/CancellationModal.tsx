import { AlertTriangle, Info, Minus, Plus, X } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import Modal from "../../../shared/ui/Modal"
import type { EnrichedBooking, TicketRecord } from "../userDashboardData"
import { formatCurrency, getTicketPricingForEvent } from "../userDashboardData"

type CancellationSelection = {
  ticketType: string
  ownedCount: number
  cancelCount: number
  price: number
}

type Props = {
  booking: EnrichedBooking | null
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

    const pricing = getTicketPricingForEvent(booking.event_title)
    const baseTypes = Array.from(new Set([...pricing.map((item) => item.ticketType), ...grouped.keys()]))

    return baseTypes.map((ticketType) => ({
      ticketType,
      ownedCount: grouped.get(ticketType) ?? 0,
      cancelCount: 0,
      price: pricing.find((item) => item.ticketType === ticketType)?.price ?? 0,
    }))
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
      className="relative w-[min(94vw,640px)] rounded-[28px] bg-white px-9 pb-10 pt-6 shadow-[0_24px_80px_rgba(15,23,42,0.24)]"
    >
      <button
        type="button"
        onClick={onClose}
        className="absolute right-6 top-6 rounded-full p-1 text-[#e6e6e6] transition hover:bg-[#f7f7f7]"
      >
        <X className="h-8 w-8" />
      </button>

      <div className="flex flex-col items-center">
        <div className="mt-5 flex h-28 w-28 items-center justify-center rounded-full">
          <AlertTriangle className="h-24 w-24 fill-[#ff1717] text-[#ff1717]" strokeWidth={1.6} />
        </div>

        <h2 className="mt-3 text-[2.5rem] font-semibold text-[#111111]">Cancel Tickets</h2>
        <p className="mt-2 text-center text-[1.35rem] text-[#5b6069]">
          Select the tickets you want to cancel for <span className="font-semibold">{booking?.event_title}</span>
        </p>

        <div className="mt-7 flex w-full flex-col gap-4">
          {selection.map((item) => (
            <div
              key={item.ticketType}
              className="grid items-center gap-4 rounded-[18px] border border-[#ffc9cf] px-4 py-5 md:grid-cols-[1fr_auto_auto]"
            >
              <div className="text-[1.7rem] font-medium text-[#111111]">{item.ticketType}</div>

              <div className="flex items-center gap-3 text-[1.6rem] text-[#111111]">
                <span>Qty</span>
                <button
                  type="button"
                  onClick={() => updateCount(item.ticketType, item.cancelCount - 1)}
                  className="rounded-md bg-[#f7c6cd] px-2 py-1 text-[#111111] transition hover:bg-[#f2b4be]"
                >
                  <Minus className="h-4 w-4" />
                </button>
                <span className="min-w-5 text-center">{item.cancelCount}</span>
                <button
                  type="button"
                  onClick={() => updateCount(item.ticketType, item.cancelCount + 1)}
                  className="rounded-md bg-[#f7c6cd] px-2 py-1 text-[#111111] transition hover:bg-[#f2b4be]"
                >
                  <Plus className="h-4 w-4" />
                </button>
              </div>

              <div className="justify-self-end text-[1.65rem] font-semibold text-[#111111]">
                {formatCurrency(item.cancelCount * item.price)}
              </div>
            </div>
          ))}
        </div>

        <div className="mt-8 flex w-full items-center justify-between text-[1.95rem]">
          <span className="text-[#111111]">Total Refund</span>
          <span className="font-semibold text-[#ff1717]">{formatCurrency(totalRefund)}</span>
        </div>

        <div className="mt-6 flex w-full items-center gap-3 rounded-[18px] border border-[#ffb9c1] bg-[#fff3f5] px-4 py-4 text-[1.35rem] text-[#ef3650]">
          <Info className="h-5 w-5 fill-[#ef3650] text-white" />
          <span>Refunds are deposited to the wallet.</span>
        </div>

        <div className="mt-8 flex w-full flex-col gap-3">
          <button
            type="button"
            disabled={totalRefund === 0}
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
            className="rounded-[18px] bg-[#ff1717] px-6 py-4 text-[1.6rem] font-semibold text-white transition hover:bg-[#e41212] disabled:cursor-not-allowed disabled:opacity-60"
          >
            Proceed to Cancel
          </button>

          <button
            type="button"
            onClick={onClose}
            className="rounded-[18px] border border-[#111111] px-6 py-4 text-[1.6rem] font-semibold text-[#111111] transition hover:bg-[#f7f7f7]"
          >
            Keep Tickets
          </button>
        </div>
      </div>
    </Modal>
  )
}
