export type EventTicketOption = {
  id?: string
  name: string
  price: number
  availableQuantity?: number
  totalQuantity?: number
}

export type DisplayEvent = {
  id: string
  title: string
  status?: string
  coverImageUrl: string
  city: string
  venue: string
  startTime: string
  endTime?: string
  dateLabel: string
  timeLabel: string
  displayLocation: string
  durationLabel: string
  priceLabel: string
  description?: string
  venueName: string
  venueAddress: string
  mapUrl: string
  termsAndConditions: string
  about?: {
    subtitle: string
    content: string[]
  }
  host?: {
    name: string
    organization?: string
    role: string
    avatar: string
  }
  personnel?: Array<{
    name: string
    role: string
    avatar: string
  }>
  ticketTypes: EventTicketOption[]
  availableCapacity?: number
  totalCapacity?: number
}

const fallbackTicketTypes: EventTicketOption[] = [
  { name: "Standard", price: 5000 },
  { name: "Premium", price: 7500 },
  { name: "VIP Access", price: 10000 },
]

const fallbackEvent = {
  title: "Friday Night at Vapour Ladies Night",
  cover_image_url: "/assets/images/badass-bollywood.png",
  date: "Saturday, 21 February 2026",
  time: "07:00 PM",
  venue: "JLN Stadium, Kochi",
  city: "Kochi",
  duration: "5 hours 30 minutes",
  price: "5,000",
  description:
    "Duis placerat nisl at nisi luctus in rhoncus felis condimentum. Vivamus in augue et sem porttitor scelerisque at ac ex.",
  about: {
    subtitle: "Saturday Bollywood Dhamaka",
    content: [
      "Duis placerat nisl at nisi luctus in rhoncus felis condimentum. Vivamus in augue et sem porttitor scelerisque at ac ex. Nam vel gravida lorem.",
      "Aliquam ultrices pretium odio nec hendrerit. Curabitur quis massa interdum, condimentum purus eu, bibendum felis. Proin libero ex, maximus et quam ut, volutpat condimentum tellus. Aliquam erat volutpat.",
      "Ut ipsum eros venenatis eu velit vitae landit bibendum massa.",
      "Cras id urna a quam viverra egestas sit amet et ante. In hac habitasse platea dictumst. Cras nec blandit nisi. Sed ac massa arcu.",
    ],
  },
  host: {
    name: "Jane Doe",
    role: "Event Organizer",
    avatar: "/assets/images/host.jpg",
  },
  personnel: [
    {
      name: "Joe Smith",
      role: "Lead Performer",
      avatar: "/assets/images/perfomer.jpg",
    },
    {
      name: "DJ Jazee",
      role: "Professional DJ",
      avatar: "/assets/images/dj.jpg",
    },
  ],
  ticket_types: fallbackTicketTypes,
}

const resolveAssetUrl = (value?: string) => {
  if (!value) return value
  if (value.startsWith("http://") || value.startsWith("https://") || value.startsWith("data:")) {
    return value
  }

  const baseUrl = import.meta.env.VITE_API_BASE_URL?.replace(/\/$/, "") ?? ""
  if (!baseUrl) return value

  if (value.startsWith("/")) {
    return `${baseUrl}${value}`
  }

  return `${baseUrl}/${value}`
}

export const formatCurrency = (amount: number) =>
  new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    maximumFractionDigits: 0,
  }).format(amount)

const formatDateLabel = (date: string) =>
  new Date(date).toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "long",
    year: "numeric",
  })

const formatTimeLabel = (date: string) =>
  new Date(date).toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
    hour12: true,
  })

const formatDurationLabel = (startTime?: string, endTime?: string) => {
  if (!startTime || !endTime) return fallbackEvent.duration

  const diff = new Date(endTime).getTime() - new Date(startTime).getTime()
  if (Number.isNaN(diff) || diff <= 0) return fallbackEvent.duration

  const totalMinutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(totalMinutes / 60)
  const minutes = totalMinutes % 60

  if (!hours) return `${minutes} minutes`
  if (!minutes) return `${hours} hours`
  return `${hours} hours ${minutes} minutes`
}

export const buildDisplayEvent = (eventId: string, payload: any): DisplayEvent => {
  const rawEvent = payload?.event ?? payload ?? {}
  const rawDetails = payload?.details ?? {}
  const rawPersonnels = payload?.personnels ?? payload?.personnel ?? rawEvent.personnel ?? rawEvent.Personnels ?? []
  const rawTicketTypes = payload?.ticket_types ?? rawEvent.ticket_types ?? rawEvent.TicketTypes ?? []

  const startTime = rawEvent.start_time ?? rawEvent.StartTime ?? rawEvent.startTime ?? "2026-02-21T19:00:00Z"
  const endTime = rawEvent.end_time ?? rawEvent.EndTime ?? rawEvent.endTime
  const venue =
    rawEvent.venue_name ??
    rawEvent.VenueName ??
    rawDetails.venue_name ??
    rawDetails.venue_address ??
    rawDetails.VenueAddress ??
    rawEvent.venue ??
    fallbackEvent.venue
  const city = rawEvent.city ?? rawEvent.City ?? rawDetails.city ?? fallbackEvent.city
  const hostSource = payload?.host ?? rawEvent.host ?? null
  const personnel =
    Array.isArray(rawPersonnels) && rawPersonnels.length > 0
      ? rawPersonnels.map((person: any) => ({
          name: person.name ?? person.Name,
          role: person.role ?? person.Role,
          avatar:
            resolveAssetUrl(person.avatar ?? person.image ?? person.Image ?? fallbackEvent.personnel[0]?.avatar) ??
            fallbackEvent.personnel[0]?.avatar ??
            fallbackEvent.host.avatar,
        }))
      : []
  const ticketTypes =
    Array.isArray(rawTicketTypes) && rawTicketTypes.length > 0
      ? rawTicketTypes.map((ticket: any) => ({
          id: ticket.id ?? ticket.ID,
          name: ticket.name ?? ticket.Name,
          price: Number(ticket.price ?? ticket.Price ?? 0),
          availableQuantity: ticket.available_quantity ?? ticket.AvailableQuantity,
          totalQuantity: ticket.total_quantity ?? ticket.TotalQuantity,
        }))
      : fallbackTicketTypes

  return {
    id: rawEvent.id ?? rawEvent.ID ?? rawEvent.slug ?? rawEvent.Slug ?? eventId,
    title: rawEvent.title ?? rawEvent.Title ?? fallbackEvent.title,
    status: rawEvent.status ?? rawEvent.Status,
    coverImageUrl: resolveAssetUrl(rawEvent.cover_image_url ?? rawEvent.CoverImageURL ?? rawEvent.coverImageUrl ?? fallbackEvent.cover_image_url) ?? fallbackEvent.cover_image_url,
    city,
    venue,
    startTime,
    endTime,
    dateLabel: rawEvent.date ?? formatDateLabel(startTime),
    timeLabel: rawEvent.time ?? formatTimeLabel(startTime),
    displayLocation: venue || city,
    durationLabel: rawEvent.duration ?? formatDurationLabel(startTime, endTime),
    priceLabel: rawEvent.price ?? fallbackEvent.price,
    description: rawDetails.description ?? rawDetails.Description ?? rawEvent.description ?? rawEvent.Description ?? fallbackEvent.description,
    venueName: rawEvent.venue_name ?? rawEvent.VenueName ?? venue,
    venueAddress: rawDetails.venue_address ?? rawDetails.VenueAddress ?? "",
    mapUrl: rawDetails.map_url ?? rawDetails.MapURL ?? "",
    termsAndConditions: rawDetails.terms_and_conditions ?? rawDetails.TermsAndConditions ?? "",
    about: {
      subtitle: rawEvent.title ?? rawEvent.Title ?? fallbackEvent.about.subtitle,
      content: [
        rawDetails.description ?? rawDetails.Description ?? rawEvent.description ?? rawEvent.Description ?? fallbackEvent.description,
      ],
    },
    host: hostSource
      ? {
          name: hostSource.name ?? hostSource.Name ?? fallbackEvent.host.name,
          organization: hostSource.organization ?? hostSource.Organization,
          role: hostSource.role ?? hostSource.Role ?? fallbackEvent.host.role,
          avatar:
            resolveAssetUrl(hostSource.avatar ?? hostSource.image ?? hostSource.Image ?? fallbackEvent.host.avatar) ??
            fallbackEvent.host.avatar,
        }
      : undefined,
    personnel: personnel.length > 0 ? personnel : undefined,
    ticketTypes,
    availableCapacity: rawDetails.available_capacity ?? rawDetails.AvailableCapacity,
    totalCapacity: rawDetails.total_capacity ?? rawDetails.TotalCapacity,
  }
}
