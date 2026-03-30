export type BookingStatus = "upcoming" | "finished"

export type BookingRecord = {
  booking_id: string
  event_id: string
  event_title: string
  event_city: string
  event_start_time: string
  status: string
  total_amount: number
  ticket_count: number
  created_at: string
  coverImageUrl: string
  venue: string
  tags: string[]
  event_status: string
}

export type TicketRecord = {
  ticket_id: string
  ticket_code: string
  event_id: string
  event_title: string
  ticket_type: string
  status: string
  checked_in_at?: string | null
}

export const formatEventTime = (value: string) =>
  new Date(value).toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  })

export const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    maximumFractionDigits: 0,
  }).format(amount)

export const formatDateBadge = (value: string) =>
  new Date(value).toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  })
