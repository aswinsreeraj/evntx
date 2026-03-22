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

export type EnrichedBooking = BookingRecord & {
  coverImageUrl: string
  venue: string
  timeLabel: string
  dateBadge: string
  tags: string[]
  bookingStatus: BookingStatus
}

type EventMetadata = {
  title: string
  coverImageUrl: string
  venue: string
  tags: string[]
  ticketPricing: Array<{
    ticketType: string
    price: number
  }>
}

const eventMetadata: EventMetadata[] = [
  {
    title: "Sand Castle Workshop",
    coverImageUrl: "/assets/images/sand-castle.png",
    venue: "XYZ Conference hall, Pune",
    tags: ["Art", "Workshop"],
    ticketPricing: [
      { ticketType: "Standard", price: 5000 },
      { ticketType: "Premium", price: 7500 },
      { ticketType: "VIP Access", price: 10000 },
    ],
  },
  {
    title: "Premium Roy by Shreya",
    coverImageUrl: "/assets/images/premium-roy.png",
    venue: "ABC Cafe, Chennai",
    tags: ["Comedy", "Live Show"],
    ticketPricing: [
      { ticketType: "Standard", price: 5000 },
      { ticketType: "Premium", price: 7500 },
      { ticketType: "VIP Access", price: 10000 },
    ],
  },
  {
    title: "Scorpions Coming Home Live 2026",
    coverImageUrl: "/assets/images/scorpions.png",
    venue: "MNO Stadium, Pune",
    tags: ["Music", "Live Show"],
    ticketPricing: [
      { ticketType: "Standard", price: 5000 },
      { ticketType: "Premium", price: 7500 },
      { ticketType: "VIP Access", price: 10000 },
    ],
  },
  {
    title: "Advancing Passive Fire Protection",
    coverImageUrl: "/assets/images/fire-protect.png",
    venue: "XYZ Conference hall, Pune",
    tags: ["Art", "Workshop"],
    ticketPricing: [
      { ticketType: "Standard", price: 2200 },
    ],
  },
  {
    title: "If I'm Not Wrong By Tarang Hardikar",
    coverImageUrl: "/assets/images/if-im-not-wrong.png",
    venue: "ABC Cafe, Chennai",
    tags: ["Comedy", "Live Show"],
    ticketPricing: [
      { ticketType: "Standard", price: 1500 },
    ],
  },
]

export const fallbackBookings: BookingRecord[] = [
  {
    booking_id: "booking-sand-1",
    event_id: "sand-castle-workshop",
    event_title: "Sand Castle Workshop",
    event_city: "Pune",
    event_start_time: "2026-02-17T12:00:00Z",
    status: "confirmed",
    total_amount: 5000,
    ticket_count: 1,
    created_at: "2026-02-01T10:00:00Z",
  },
  {
    booking_id: "booking-premium-1",
    event_id: "premium-roy-by-shreya",
    event_title: "Premium Roy by Shreya",
    event_city: "Chennai",
    event_start_time: "2026-02-17T18:00:00Z",
    status: "confirmed",
    total_amount: 7500,
    ticket_count: 1,
    created_at: "2026-02-03T10:00:00Z",
  },
  {
    booking_id: "booking-scorpions-1",
    event_id: "scorpions-coming-home-live-2026",
    event_title: "Scorpions Coming Home Live 2026",
    event_city: "Mumbai",
    event_start_time: "2026-02-28T18:00:00Z",
    status: "confirmed",
    total_amount: 10000,
    ticket_count: 1,
    created_at: "2026-02-04T10:00:00Z",
  },
  {
    booking_id: "booking-fire-1",
    event_id: "advancing-passive-fire-protection",
    event_title: "Advancing Passive Fire Protection",
    event_city: "Pune",
    event_start_time: "2026-02-21T12:00:00Z",
    status: "completed",
    total_amount: 2200,
    ticket_count: 1,
    created_at: "2026-02-02T10:00:00Z",
  },
  {
    booking_id: "booking-tarang-1",
    event_id: "if-im-not-wrong-by-tarang-hardikar",
    event_title: "If I'm Not Wrong By Tarang Hardikar",
    event_city: "Chennai",
    event_start_time: "2026-03-03T18:00:00Z",
    status: "completed",
    total_amount: 1500,
    ticket_count: 1,
    created_at: "2026-02-06T10:00:00Z",
  },
]

export const fallbackTicketsByEvent: Record<string, TicketRecord[]> = {
  "sand-castle-workshop": [
    {
      ticket_id: "ticket-sand-1",
      ticket_code: "1101",
      event_id: "sand-castle-workshop",
      event_title: "Sand Castle Workshop",
      ticket_type: "Standard",
      status: "active",
    },
  ],
  "premium-roy-by-shreya": [
    {
      ticket_id: "ticket-premium-1",
      ticket_code: "2201",
      event_id: "premium-roy-by-shreya",
      event_title: "Premium Roy by Shreya",
      ticket_type: "Premium",
      status: "active",
    },
  ],
  "scorpions-coming-home-live-2026": [
    {
      ticket_id: "ticket-scorpions-1",
      ticket_code: "6301",
      event_id: "scorpions-coming-home-live-2026",
      event_title: "Scorpions Coming Home Live 2026",
      ticket_type: "Premium",
      status: "active",
    },
    {
      ticket_id: "ticket-scorpions-2",
      ticket_code: "6302",
      event_id: "scorpions-coming-home-live-2026",
      event_title: "Scorpions Coming Home Live 2026",
      ticket_type: "Premium",
      status: "active",
    },
  ],
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

const getEventMetadata = (title: string, city: string) => {
  const exactMatch = eventMetadata.find((event) => event.title.toLowerCase() === title.toLowerCase())
  if (exactMatch) return exactMatch

  return {
    title,
    coverImageUrl: "/assets/images/badass-bollywood.png",
    venue: city ? `${city} City Arena` : "Main Event Venue",
    tags: ["Live Event"],
    ticketPricing: [
      { ticketType: "Standard", price: 5000 },
    ],
  }
}

export const getTicketPricingForEvent = (title: string) =>
  getEventMetadata(title, "").ticketPricing

export const enrichBookings = (bookings: BookingRecord[]): EnrichedBooking[] =>
  bookings.map((booking) => {
    const metadata = getEventMetadata(booking.event_title, booking.event_city)
    const bookingDate = new Date(booking.event_start_time)
    const bookingStatus: BookingStatus = bookingDate.getTime() >= Date.now() ? "upcoming" : "finished"

    return {
      ...booking,
      coverImageUrl: metadata.coverImageUrl,
      venue: metadata.venue,
      tags: metadata.tags,
      timeLabel: formatEventTime(booking.event_start_time),
      dateBadge: formatDateBadge(booking.event_start_time),
      bookingStatus,
    }
  })

export const getMonthLabel = (monthIndex: number) =>
  new Date(2026, monthIndex, 1).toLocaleString("en-US", { month: "short" })
