import EventCard from "./EventCard"

function EventGrid({ events }: any) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {events.map((event: any) => (
        <EventCard key={event.id} event={event} />
      ))}
    </div>
  )
}

export default EventGrid