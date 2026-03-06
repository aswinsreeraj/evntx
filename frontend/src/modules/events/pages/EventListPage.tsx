import { useEvents } from "../hooks"
import EventGrid from "../components/EventGrid"

function EventListPage() {
  const { data, isLoading } = useEvents()

  if (isLoading) return <p>Loading events...</p>

  return (
    <div className="max-w-7xl mx-auto px-6 py-10 grid grid-cols-4 gap-8">

      {/* Sidebar */}
      <div className="col-span-1 bg-white p-6 rounded-xl shadow-sm">
        <h3 className="font-semibold mb-4">Filters</h3>
      </div>

      {/* Event list */}
      <div className="col-span-3">
        <EventGrid events={data.events} />
      </div>
    </div>
  )
}

export default EventListPage