import { ChevronLeft, ChevronRight, MapPin } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import UserDashboardShell from "../components/UserDashboardShell"
import { userApi } from "../api"
import { enrichBookings, getMonthLabel, type BookingRecord } from "../userDashboardData"

const monthOptions = Array.from({ length: 12 }, (_, month) => month)

export default function CalendarPage() {
  const [bookings, setBookings] = useState<BookingRecord[]>([])
  const [loading, setLoading] = useState(true)
  const now = new Date()
  const [activeMonth, setActiveMonth] = useState(now.getMonth())
  const [activeYear, setActiveYear] = useState(now.getFullYear())
  const [selectedDay, setSelectedDay] = useState(now.getDate())

  useEffect(() => {
    const loadBookings = async () => {
      try {
        const data = await userApi.getMyBookings()
        setBookings(data)
      } catch {
        setBookings([])
      } finally {
        setLoading(false)
      }
    }

    void loadBookings()
  }, [])

  const enrichedBookings = useMemo(() => enrichBookings(bookings), [bookings])

  const monthBookings = enrichedBookings.filter((booking) => {
    const date = new Date(booking.event_start_time)
    return date.getMonth() === activeMonth && date.getFullYear() === activeYear
  })

  const daysInMonth = new Date(activeYear, activeMonth + 1, 0).getDate()
  const selectedDateBookings = monthBookings.filter(
    (booking) => new Date(booking.event_start_time).getDate() === selectedDay,
  )

  const firstDayOfMonth = new Date(activeYear, activeMonth, 1).getDay()
  const emptyDaysOffset = firstDayOfMonth === 0 ? 6 : firstDayOfMonth - 1

  const eventDays = Array.from(new Set(monthBookings.map((booking) => new Date(booking.event_start_time).getDate())))

  const weekdays = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]

  return (
    <UserDashboardShell>
      <div className="pb-10">
        <h1 className="mb-5 text-[2rem] tracking-[0.02em] text-[#111111]">Your Schedule</h1>

        <div className="grid gap-5 xl:grid-cols-[1fr_260px]">
          <section className="rounded-[20px] border border-[#ececec] bg-white p-5 shadow-[0_12px_28px_rgba(15,23,42,0.06)]">
            <div className="mb-5 flex items-center justify-between">
              <button
                type="button"
                onClick={() => {
                  if (activeMonth === 0) {
                    setActiveMonth(11)
                    setActiveYear((y) => y - 1)
                  } else {
                    setActiveMonth((m) => m - 1)
                  }
                }}
                className="rounded-full p-2 text-[#ff5561] transition hover:bg-[#fff4f5]"
              >
                <ChevronLeft className="h-6 w-6" />
              </button>

              <div className="flex gap-4">
                <select
                  value={activeMonth}
                  onChange={(event) => setActiveMonth(Number(event.target.value))}
                  className="rounded-2xl border border-[#e7e7e7] px-5 py-2.5 text-base shadow-sm outline-none"
                >
                  {monthOptions.map((month) => (
                    <option key={month} value={month}>
                      {getMonthLabel(month)}
                    </option>
                  ))}
                </select>

                <select
                  value={activeYear}
                  onChange={(event) => setActiveYear(Number(event.target.value))}
                  className="rounded-2xl border border-[#e7e7e7] px-5 py-2.5 text-base shadow-sm outline-none"
                >
                  {[2025, 2026, 2027].map((year) => (
                    <option key={year} value={year}>
                      {year}
                    </option>
                  ))}
                </select>
              </div>

              <button
                type="button"
                onClick={() => {
                  if (activeMonth === 11) {
                    setActiveMonth(0)
                    setActiveYear((y) => y + 1)
                  } else {
                    setActiveMonth((m) => m + 1)
                  }
                }}
                className="rounded-full p-2 text-[#ff5561] transition hover:bg-[#fff4f5]"
              >
                <ChevronRight className="h-6 w-6" />
              </button>
            </div>

            <div className="grid grid-cols-7 gap-y-3 px-4 text-center text-base text-[#111111]">
              {weekdays.map((day) => (
                <div key={day} className="mb-2 text-sm font-semibold text-gray-500">{day}</div>
              ))}
              {Array.from({ length: emptyDaysOffset }, (_, i) => (
                <div key={`empty-${i}`} />
              ))}
              {Array.from({ length: daysInMonth }, (_, index) => index + 1).map((day) => {
                const hasEvent = eventDays.includes(day)
                const isSelected = selectedDay === day
                const isToday = day === now.getDate() && activeMonth === now.getMonth() && activeYear === now.getFullYear()

                let colorClass = ""
                if (hasEvent) {
                  colorClass = "bg-[#ef3650] text-white"
                }

                return (
                  <button
                    key={day}
                    type="button"
                    onClick={() => setSelectedDay(day)}
                    className={`mx-auto flex h-10 w-10 items-center justify-center rounded-2xl transition ${
                      isSelected
                        ? "bg-[#111827] text-white"
                        : isToday && !hasEvent
                          ? "ring-2 ring-[#ef3650] text-[#ef3650] font-semibold"
                          : colorClass || "text-[#111111] hover:bg-gray-100"
                    }`}
                  >
                    {day}
                  </button>
                )
              })}
            </div>
          </section>

          <aside>
            <h2 className="mb-4 text-lg font-medium text-[#202020]">Your events for this month</h2>
            <div className="flex flex-col gap-4">
              {monthBookings.map((booking, index) => (
                <button
                  key={booking.booking_id}
                  type="button"
                  onClick={() => setSelectedDay(new Date(booking.event_start_time).getDate())}
                  className="flex items-center gap-3 rounded-[20px] border border-[#efefef] bg-white px-4 py-3 text-left shadow-[0_10px_24px_rgba(15,23,42,0.05)]"
                >
                  <div className={`rounded-2xl px-3 py-2 text-center text-white ${index % 3 === 0 ? "bg-[#c9c9c9] text-[#111111]" : index % 3 === 1 ? "bg-[#ef3650]" : "bg-[#111827]"}`}>
                    <div className="text-xs">{getMonthLabel(new Date(booking.event_start_time).getMonth())}</div>
                    <div className="text-xl">{new Date(booking.event_start_time).getDate()}</div>
                  </div>
                  <div className="min-w-0">
                    <div className="truncate text-sm font-medium text-[#202020]">{booking.event_title}</div>
                    <div className="mt-1 text-xs text-[#73809b]">
                      {booking.timeLabel} • {booking.event_city}
                    </div>
                  </div>
                </button>
              ))}
            </div>
          </aside>
        </div>

        <section className="mt-8">
          <h2 className="mb-4 text-[1.7rem] tracking-[0.02em] text-[#111111]">
            Events for {selectedDay} {new Date(activeYear, activeMonth, 1).toLocaleDateString("en-GB", { month: "long", year: "numeric" })}
          </h2>

          <div className="flex flex-col gap-4">
            {loading ? (
              <div className="flex h-32 items-center justify-center rounded-[24px] border border-[#ececec] bg-white shadow-sm">
                <div className="h-8 w-8 animate-spin rounded-full border-2 border-gray-200 border-t-[#111827]" />
              </div>
            ) : null}

            {!loading && selectedDateBookings.map((booking) => (
              <div
                key={booking.booking_id}
                className="grid overflow-hidden rounded-[20px] border border-[#efefef] bg-white shadow-[0_10px_24px_rgba(15,23,42,0.05)] md:grid-cols-[220px_1fr]"
              >
                <div className="relative min-h-[136px]">
                  <img src={booking.coverImageUrl} alt={booking.event_title} className="h-full w-full object-cover" />
                  <div className="absolute left-3 top-3 rounded-xl bg-[#a77da5]/80 px-3 py-1.5 text-sm text-white backdrop-blur-sm">
                    {booking.dateBadge}
                  </div>
                </div>
                <div className="p-5">
                  <h3 className="text-[1.5rem] font-semibold text-[#111827]">{booking.event_title}</h3>
                  <p className="mt-1 text-base text-[#111827]">{booking.timeLabel}</p>
                  <div className="mt-2 flex items-center gap-2 text-sm text-[#5d6573]">
                    <MapPin className="h-4 w-4 text-[#ff445d]" />
                    <span>{booking.venue}</span>
                  </div>
                  <div className="mt-3 flex flex-wrap gap-2">
                    {booking.tags.map((tag) => (
                      <span key={tag} className="rounded-full bg-[#efefef] px-3 py-1 text-xs text-[#5b6069]">
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            ))}

            {!loading && selectedDateBookings.length === 0 ? (
              <div className="rounded-[20px] border border-dashed border-[#e2e2e2] bg-white px-6 py-10 text-center text-base text-[#7b8392]">
                No booked events for this day.
              </div>
            ) : null}
          </div>
        </section>
      </div>
    </UserDashboardShell>
  )
}
