export type BookingSelection = {
  ticketName: string
  quantity: number
  unitPrice: number
}

export type BookingConfirmation = {
  bookingId: string
  eventId: string
  eventTitle: string
  eventImage: string
  eventDate: string
  eventTime: string
  venue: string
  totalAmount: number
  platformFee: number
  finalAmount: number
  email: string
  tickets: BookingSelection[]
  createdAt: string
}

const STORAGE_KEY = "event-booking-confirmation"

export const saveBookingConfirmation = (booking: BookingConfirmation) => {
  sessionStorage.setItem(STORAGE_KEY, JSON.stringify(booking))
}

export const getBookingConfirmation = (): BookingConfirmation | null => {
  const raw = sessionStorage.getItem(STORAGE_KEY)
  if (!raw) return null

  try {
    return JSON.parse(raw) as BookingConfirmation
  } catch {
    sessionStorage.removeItem(STORAGE_KEY)
    return null
  }
}
