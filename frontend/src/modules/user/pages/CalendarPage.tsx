import { ChevronLeft, ChevronRight, MapPin } from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import UserDashboardShell from "../components/UserDashboardShell"
import { userApi } from "../api"
import { enrichBookings, fallbackBookings, getMonthLabel, type BookingRecord } from "../userDashboardData"

const monthOptions = Array.from({ length: 12 }, (_, month) => month)

export default function CalendarPage() {
  const [bookings, setBookings] = useState<BookingRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [activeMonth, setActiveMonth] = useState(1)
  const [activeYear, setActiveYear] = useState(2026)
  const [selectedDay, setSelectedDay] = useState(17)

  useEffect(() => {
    const loadBookings = async () => {
      try {
        const data = await userApi.getMyBookings()
        setBookings(data.length > 0 ? data : fallbackBookings)
      } catch {
        setBookings(fallbackBookings)
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

  const eventDays = Array.from(new Set(monthBookings.map((booking) => new Date(booking.event_start_time).getDate())))

  return (
    <UserDashboardShell>
      <div className="pb-10">
        <h1 className="mb-8 text-[2.2rem] tracking-[0.02em] text-[#111111]">Your Schedule</h1>

        <div className="grid gap-6 xl:grid-cols-[1fr_300px]">
          <section className="rounded-[28px] border border-[#ececec] bg-white p-8 shadow-[0_12px_36px_rgba(15,23,42,0.08)]">
            <div className="mb-8 flex items-center justify-between">
              <button
                type="button"
                onClick={() => setActiveMonth((current) => (current === 0 ? 11 : current - 1))}
                className="rounded-full p-2 text-[#ff5561] transition hover:bg-[#fff4f5]"
              >
                <ChevronLeft className="h-8 w-8" />
              </button>

              <div className="flex gap-8">
                <select
                  value={activeMonth}
                  onChange={(event) => setActiveMonth(Number(event.target.value))}
                  className="rounded-2xl border border-[#e7e7e7] px-8 py-3 text-[1.6rem] shadow-sm outline-none"
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
                  className="rounded-2xl border border-[#e7e7e7] px-8 py-3 text-[1.6rem] shadow-sm outline-none"
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
                onClick={() => setActiveMonth((current) => (current === 11 ? 0 : current + 1))}
                className="rounded-full p-2 text-[#ff5561] transition hover:bg-[#fff4f5]"
              >
                <ChevronRight className="h-8 w-8" />
              </button>
            </div>

            <div className="grid grid-cols-7 gap-y-8 px-10 text-center text-[2rem] text-[#111111]">
              {Array.from({ length: daysInMonth }, (_, index) => index + 1).map((day, index) => {
                const hasEvent = eventDays.includes(day)
                const isSelected = selectedDay === day
                const colorClass = hasEvent
                  ? index % 3 === 0
                    ? "bg-[#c9c9c9] text-[#111111]"
                    : index % 3 === 1
                      ? "bg-[#ef3650] text-white"
                      : "bg-[#111827] text-white"
                  : ""

                return (
                  <button
                    key={day}
                    type="button"
                    onClick={() => setSelectedDay(day)}
                    className={`mx-auto flex h-14 w-14 items-center justify-center rounded-2xl transition ${
                      isSelected ? colorClass || "bg-[#111827] text-white" : colorClass || "text-[#111111]"
                    }`}
                  >
                    {day}
                  </button>
                )
              })}
            </div>
          </section>

          <aside>
            <h2 className="mb-5 text-[1.7rem] text-[#202020]">Your events for this month</h2>
            <div className="flex flex-col gap-4">
              {monthBookings.map((booking, index) => (
                <button
                  key={booking.booking_id}
                  type="button"
                  onClick={() => setSelectedDay(new Date(booking.event_start_time).getDate())}
                  className="flex items-center gap-4 rounded-[24px] border border-[#efefef] bg-white px-5 py-4 text-left shadow-[0_12px_30px_rgba(15,23,42,0.06)]"
                >
                  <div className={`rounded-2xl px-4 py-3 text-center text-white ${index % 3 === 0 ? "bg-[#c9c9c9] text-[#111111]" : index % 3 === 1 ? "bg-[#ef3650]" : "bg-[#111827]"}`}>
                    <div className="text-[1.2rem]">{getMonthLabel(new Date(booking.event_start_time).getMonth())}</div>
                    <div className="text-[2rem]">{new Date(booking.event_start_time).getDate()}</div>
                  </div>
                  <div className="min-w-0">
                    <div className="truncate text-[1.45rem] font-medium text-[#202020]">{booking.event_title}</div>
                    <div className="mt-1 text-[1.2rem] text-[#73809b]">
                      {booking.timeLabel} • {booking.event_city}
                    </div>
                  </div>
                </button>
              ))}
            </div>
          </aside>
        </div>

        <section className="mt-8">
          <h2 className="mb-5 text-[2rem] tracking-[0.02em] text-[#111111]">
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
                className="grid overflow-hidden rounded-[28px] border border-[#efefef] bg-white shadow-[0_12px_36px_rgba(15,23,42,0.06)] md:grid-cols-[280px_1fr]"
              >
                <div className="relative min-h-[160px]">
                  <img src={booking.coverImageUrl} alt={booking.event_title} className="h-full w-full object-cover" />
                  <div className="absolute left-4 top-4 rounded-2xl bg-[#a77da5]/80 px-4 py-2 text-[1.4rem] text-white backdrop-blur-sm">
                    {booking.dateBadge}
                  </div>
                </div>
                <div className="p-6">
                  <h3 className="text-[1.9rem] font-semibold text-[#111827]">{booking.event_title}</h3>
                  <p className="mt-1 text-[1.4rem] text-[#111827]">{booking.timeLabel}</p>
                  <div className="mt-2 flex items-center gap-2 text-[1.3rem] text-[#5d6573]">
                    <MapPin className="h-5 w-5 text-[#ff445d]" />
                    <span>{booking.venue}</span>
                  </div>
                  <div className="mt-4 flex flex-wrap gap-3">
                    {booking.tags.map((tag) => (
                      <span key={tag} className="rounded-full bg-[#efefef] px-4 py-1.5 text-[1.1rem] text-[#5b6069]">
                        {tag}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            ))}

            {!loading && selectedDateBookings.length === 0 ? (
              <div className="rounded-[24px] border border-dashed border-[#e2e2e2] bg-white px-6 py-10 text-center text-[1.3rem] text-[#7b8392]">
                No booked events for this day.
              </div>
            ) : null}
          </div>
        </section>
      </div>
    </UserDashboardShell>
  )
}
